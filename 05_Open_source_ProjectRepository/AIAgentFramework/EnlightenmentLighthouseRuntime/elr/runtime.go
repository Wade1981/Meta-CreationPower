// Package elr implements the core runtime for Enlightenment Lighthouse Runtime
package elr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Runtime represents the core runtime of ELR
type Runtime struct {
	Config         *Config
	Platform       Platform
	Containers     map[string]*Container
	ContainerMutex sync.RWMutex
	Plugins        map[string]Plugin
	PluginMutex    sync.RWMutex
	Stopped        bool
	StopCh         chan struct{}
}

// Config represents the runtime configuration
type Config struct {
	LogLevel  string `yaml:"log_level"`
	DataDir   string `yaml:"data_dir"`
	PluginDir string `yaml:"plugin_dir"`
	Platform  struct {
		Linux struct {
			UseNamespaces bool `yaml:"use_namespaces"`
			UseCgroups    bool `yaml:"use_cgroups"`
		} `yaml:"linux"`
		Windows struct {
			UseJobObjects bool `yaml:"use_job_objects"`
			UseWSL        bool `yaml:"use_wsl"`
		} `yaml:"windows"`
		Darwin struct {
			UseSandbox bool `yaml:"use_sandbox"`
			UseSpctl   bool `yaml:"use_spctl"`
		} `yaml:"darwin"`
	} `yaml:"platform"`
	Network struct {
		Enable  bool   `yaml:"enable"`
		Bridge  string `yaml:"bridge"`
		Subnet  string `yaml:"subnet"`
	} `yaml:"network"`
	Storage struct {
		Enable  bool   `yaml:"enable"`
		Driver  string `yaml:"driver"`
		BaseDir string `yaml:"base_dir"`
	} `yaml:"storage"`
	Languages map[string]struct {
		Enable  bool   `yaml:"enable"`
		Runtime string `yaml:"runtime"`
	} `yaml:"languages"`
}

// Platform represents the platform abstraction layer
type Platform interface {
	Name() string
	Version() string
	Init() error
	Cleanup() error
	CreateContainer(id string, config ContainerConfig) (Container, error)
	DestroyContainer(container Container) error
}

// Plugin represents a runtime plugin
type Plugin interface {
	Name() string
	Version() string
	Init(runtime *Runtime) error
	Cleanup() error
}

// NewRuntime creates a new runtime instance
func NewRuntime(config *Config) (*Runtime, error) {
	// Initialize data directory
	if config.DataDir == "" {
		config.DataDir = filepath.Join(os.Getenv("HOME"), ".elr", "data")
	}
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	// Initialize plugin directory
	if config.PluginDir == "" {
		config.PluginDir = filepath.Join(os.Getenv("HOME"), ".elr", "plugins")
	}
	if err := os.MkdirAll(config.PluginDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %v", err)
	}

	// Initialize platform
	var platform Platform
	var err error

	switch runtime.GOOS {
	case "linux":
		platform, err = NewLinuxPlatform(config)
	case "windows":
		platform, err = NewWindowsPlatform(config)
	case "darwin":
		platform, err = NewDarwinPlatform(config)
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize platform: %v", err)
	}

	if err := platform.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize platform: %v", err)
	}

	// Create runtime
	runtime := &Runtime{
		Config:     config,
		Platform:   platform,
		Containers: make(map[string]*Container),
		Plugins:    make(map[string]Plugin),
		StopCh:     make(chan struct{}),
	}

	// Load plugins
	if err := runtime.loadPlugins(); err != nil {
		return nil, fmt.Errorf("failed to load plugins: %v", err)
	}

	return runtime, nil
}

// Start starts the runtime
func (r *Runtime) Start() error {
	fmt.Printf("Starting Enlightenment Lighthouse Runtime v%s\n", Version)
	fmt.Printf("Platform: %s %s\n", r.Platform.Name(), r.Platform.Version())
	fmt.Printf("Data directory: %s\n", r.Config.DataDir)
	fmt.Printf("Plugin directory: %s\n", r.Config.PluginDir)

	// Start plugins
	for name, plugin := range r.Plugins {
		if err := plugin.Init(r); err != nil {
			fmt.Printf("Warning: failed to initialize plugin %s: %v\n", name, err)
		} else {
			fmt.Printf("Initialized plugin: %s v%s\n", name, plugin.Version())
		}
	}

	// Load existing containers
	if err := r.loadContainers(); err != nil {
		fmt.Printf("Warning: failed to load existing containers: %v\n", err)
	}

	fmt.Println("Enlightenment Lighthouse Runtime started successfully!")
	return nil
}

// Stop stops the runtime
func (r *Runtime) Stop() error {
	if r.Stopped {
		return nil
	}

	r.Stopped = true
	close(r.StopCh)

	fmt.Println("Stopping Enlightenment Lighthouse Runtime...")

	// Stop all containers
	r.ContainerMutex.Lock()
	containers := make([]*Container, 0, len(r.Containers))
	for _, container := range r.Containers {
		containers = append(containers, container)
	}
	r.ContainerMutex.Unlock()

	for _, container := range containers {
		if err := container.Stop(); err != nil {
			fmt.Printf("Warning: failed to stop container %s: %v\n", container.ID, err)
		}
	}

	// Cleanup plugins
	for name, plugin := range r.Plugins {
		if err := plugin.Cleanup(); err != nil {
			fmt.Printf("Warning: failed to cleanup plugin %s: %v\n", name, err)
		} else {
			fmt.Printf("Cleaned up plugin: %s\n", name)
		}
	}

	// Cleanup platform
	if err := r.Platform.Cleanup(); err != nil {
		fmt.Printf("Warning: failed to cleanup platform: %v\n", err)
	}

	fmt.Println("Enlightenment Lighthouse Runtime stopped successfully!")
	return nil
}

// CreateContainer creates a new container
func (r *Runtime) CreateContainer(config ContainerConfig) (*Container, error) {
	if r.Stopped {
		return nil, fmt.Errorf("runtime is stopped")
	}

	// Generate container ID if not provided
	if config.ID == "" {
		config.ID = fmt.Sprintf("elr-%d", time.Now().UnixNano())
	}

	// Check if container already exists
	r.ContainerMutex.RLock()
	if _, exists := r.Containers[config.ID]; exists {
		r.ContainerMutex.RUnlock()
		return nil, fmt.Errorf("container with ID %s already exists", config.ID)
	}
	r.ContainerMutex.RUnlock()

	// Create container directory
	containerDir := filepath.Join(r.Config.DataDir, "containers", config.ID)
	if err := os.MkdirAll(containerDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create container directory: %v", err)
	}

	// Create container
	container := &Container{
		ID:      config.ID,
		Name:    config.Name,
		Image:   config.Image,
		Command: config.Command,
		Args:    config.Args,
		Env:     config.Env,
		Dir:     containerDir,
		Runtime: r,
		Status:  ContainerStatusCreated,
		Created: time.Now(),
	}

	// Save container config
	if err := container.saveConfig(); err != nil {
		return nil, fmt.Errorf("failed to save container config: %v", err)
	}

	// Add container to runtime
	r.ContainerMutex.Lock()
	r.Containers[config.ID] = container
	r.ContainerMutex.Unlock()

	fmt.Printf("Created container: %s (%s)\n", container.ID, container.Name)
	return container, nil
}

// GetContainer gets a container by ID
func (r *Runtime) GetContainer(id string) (*Container, error) {
	r.ContainerMutex.RLock()
	container, exists := r.Containers[id]
	r.ContainerMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("container with ID %s not found", id)
	}

	return container, nil
}

// ListContainers lists all containers
func (r *Runtime) ListContainers() []*Container {
	r.ContainerMutex.RLock()
	containers := make([]*Container, 0, len(r.Containers))
	for _, container := range r.Containers {
		containers = append(containers, container)
	}
	r.ContainerMutex.RUnlock()

	return containers
}

// DeleteContainer deletes a container
func (r *Runtime) DeleteContainer(id string) error {
	// Get container
	container, err := r.GetContainer(id)
	if err != nil {
		return err
	}

	// Stop container if running
	if container.Status == ContainerStatusRunning {
		if err := container.Stop(); err != nil {
			return fmt.Errorf("failed to stop container: %v", err)
		}
	}

	// Remove container from runtime
	r.ContainerMutex.Lock()
	delete(r.Containers, id)
	r.ContainerMutex.Unlock()

	// Delete container directory
	if err := os.RemoveAll(container.Dir); err != nil {
		return fmt.Errorf("failed to delete container directory: %v", err)
	}

	fmt.Printf("Deleted container: %s (%s)\n", container.ID, container.Name)
	return nil
}

// loadPlugins loads plugins from the plugin directory
func (r *Runtime) loadPlugins() error {
	pluginDir := r.Config.PluginDir

	// Check if plugin directory exists
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		return nil
	}

	// Read plugin directory
	files, err := os.ReadDir(pluginDir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %v", err)
	}

	// Load plugins
	for _, file := range files {
		if file.IsDir() {
			pluginPath := filepath.Join(pluginDir, file.Name())
			if err := r.loadPlugin(pluginPath); err != nil {
				fmt.Printf("Warning: failed to load plugin %s: %v\n", file.Name(), err)
			}
		}
	}

	return nil
}

// loadPlugin loads a single plugin
func (r *Runtime) loadPlugin(pluginPath string) error {
	// TODO: Implement plugin loading
	// This is a placeholder for now
	return nil
}

// loadContainers loads existing containers from the data directory
func (r *Runtime) loadContainers() error {
	containersDir := filepath.Join(r.Config.DataDir, "containers")

	// Check if containers directory exists
	if _, err := os.Stat(containersDir); os.IsNotExist(err) {
		return nil
	}

	// Read containers directory
	files, err := os.ReadDir(containersDir)
	if err != nil {
		return fmt.Errorf("failed to read containers directory: %v", err)
	}

	// Load containers
	for _, file := range files {
		if file.IsDir() {
			containerID := file.Name()
			containerDir := filepath.Join(containersDir, containerID)
			
			// Load container config
			container, err := loadContainer(containerDir, r)
			if err != nil {
				fmt.Printf("Warning: failed to load container %s: %v\n", containerID, err)
				continue
			}

			// Add container to runtime
			r.ContainerMutex.Lock()
			r.Containers[containerID] = container
			r.ContainerMutex.Unlock()

			fmt.Printf("Loaded container: %s (%s)\n", container.ID, container.Name)
		}
	}

	return nil
}

// Version represents the runtime version
const Version = "1.0.0"
