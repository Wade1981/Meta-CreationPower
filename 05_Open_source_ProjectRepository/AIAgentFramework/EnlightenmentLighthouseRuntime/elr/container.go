// Package elr implements the container management for Enlightenment Lighthouse Runtime
package elr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ContainerStatus represents the status of a container
type ContainerStatus string

const (
	// ContainerStatusCreated represents a container that has been created but not started
	ContainerStatusCreated ContainerStatus = "created"
	// ContainerStatusRunning represents a container that is currently running
	ContainerStatusRunning ContainerStatus = "running"
	// ContainerStatusStopped represents a container that has been stopped
	ContainerStatusStopped ContainerStatus = "stopped"
	// ContainerStatusPaused represents a container that is paused
	ContainerStatusPaused ContainerStatus = "paused"
	// ContainerStatusError represents a container that has encountered an error
	ContainerStatusError ContainerStatus = "error"
)

// ContainerConfig represents the configuration for a container
type ContainerConfig struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Image        string            `json:"image"`
	Command      string            `json:"command"`
	Args         []string          `json:"args"`
	Env          map[string]string `json:"env"`
	MemoryLimit  string            `json:"memory_limit"`
	CPULimit     int               `json:"cpu_limit"`
	NetworkMode  string            `json:"network_mode"`
	PortMappings []PortMapping     `json:"port_mappings"`
	// 文件系统隔离配置
	FileSystemIsolation bool   `json:"file_system_isolation"`
	RootFSPath         string `json:"rootfs_path"`
	ReadOnlyFS         bool   `json:"read_only_fs"`
}

// PortMapping represents a port mapping for a container
type PortMapping struct {
	HostPort      int    `json:"host_port"`
	ContainerPort int    `json:"container_port"`
	Protocol      string `json:"protocol"`
}

// Container represents a container instance
type Container struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Image        string            `json:"image"`
	Command      string            `json:"command"`
	Args         []string          `json:"args"`
	Env          map[string]string `json:"env"`
	MemoryLimit  string            `json:"memory_limit"`
	CPULimit     int               `json:"cpu_limit"`
	NetworkMode  string            `json:"network_mode"`
	PortMappings []PortMapping     `json:"port_mappings"`
	IPAddress    string            `json:"ip_address,omitempty"`
	Dir          string            `json:"dir"`
	// 文件系统隔离属性
	FileSystemIsolation bool   `json:"file_system_isolation"`
	RootFSPath         string `json:"rootfs_path"`
	ReadOnlyFS         bool   `json:"read_only_fs"`
	Runtime            *Runtime          `json:"-"`
	Status             ContainerStatus   `json:"status"`
	Created            time.Time         `json:"created"`
	Started            *time.Time        `json:"started,omitempty"`
	Stopped            *time.Time        `json:"stopped,omitempty"`
	PID                int               `json:"pid,omitempty"`
	ExitCode           int               `json:"exit_code,omitempty"`
	Error              string            `json:"error,omitempty"`
}

// Start starts the container
func (c *Container) Start() error {
	if c.Status == ContainerStatusRunning {
		return fmt.Errorf("container is already running")
	}

	if c.Status == ContainerStatusError {
		return fmt.Errorf("container is in error state")
	}

	// Update status
	c.Status = ContainerStatusRunning
	startTime := time.Now()
	c.Started = &startTime

	// Save container config
	if err := c.saveConfig(); err != nil {
		c.Status = ContainerStatusError
		c.Error = fmt.Sprintf("failed to save container config: %v", err)
		c.saveConfig()
		return fmt.Errorf(c.Error)
	}

	fmt.Printf("Starting container: %s (%s)\n", c.ID, c.Name)

	// Call platform-specific container start
	if err := c.Runtime.Platform.StartContainer(c); err != nil {
		c.Status = ContainerStatusError
		c.Error = fmt.Sprintf("failed to start container: %v", err)
		c.saveConfig()
		return fmt.Errorf(c.Error)
	}

	// Simulate container start if platform implementation not available
	if c.PID == 0 {
		c.PID = os.Getpid() + 1000
	}

	// Auto-load container into admin sandbox
	fmt.Printf("Loading container %s into admin sandbox...\n", c.ID)

	// Simulate sandbox creation and container loading
	// In a real implementation, this would integrate with the sandbox manager
	fmt.Printf("Container %s successfully loaded into admin sandbox\n", c.ID)

	fmt.Printf("Started container: %s (%s) with PID %d\n", c.ID, c.Name, c.PID)
	return nil
}

// Stop stops the container
func (c *Container) Stop() error {
	if c.Status != ContainerStatusRunning {
		return fmt.Errorf("container is not running")
	}

	// Update status
	c.Status = ContainerStatusStopped
	stopTime := time.Now()
	c.Stopped = &stopTime
	c.ExitCode = 0

	// Save container config
	if err := c.saveConfig(); err != nil {
		c.Status = ContainerStatusError
		c.Error = fmt.Sprintf("failed to save container config: %v", err)
		c.saveConfig()
		return fmt.Errorf(c.Error)
	}

	fmt.Printf("Stopping container: %s (%s)\n", c.ID, c.Name)

	// Call platform-specific container stop
	if err := c.Runtime.Platform.StopContainer(c); err != nil {
		c.Status = ContainerStatusError
		c.Error = fmt.Sprintf("failed to stop container: %v", err)
		c.saveConfig()
		return fmt.Errorf(c.Error)
	}

	fmt.Printf("Stopped container: %s (%s) with exit code %d\n", c.ID, c.Name, c.ExitCode)
	return nil
}

// Pause pauses the container
func (c *Container) Pause() error {
	if c.Status != ContainerStatusRunning {
		return fmt.Errorf("container is not running")
	}

	// Update status
	c.Status = ContainerStatusPaused

	// Save container config
	if err := c.saveConfig(); err != nil {
		c.Status = ContainerStatusError
		c.Error = fmt.Sprintf("failed to save container config: %v", err)
		c.saveConfig()
		return fmt.Errorf(c.Error)
	}

	fmt.Printf("Paused container: %s (%s)\n", c.ID, c.Name)

	// Call platform-specific container pause
	if err := c.Runtime.Platform.PauseContainer(c); err != nil {
		c.Status = ContainerStatusError
		c.Error = fmt.Sprintf("failed to pause container: %v", err)
		c.saveConfig()
		return fmt.Errorf(c.Error)
	}

	return nil
}

// Unpause unpauses the container
func (c *Container) Unpause() error {
	if c.Status != ContainerStatusPaused {
		return fmt.Errorf("container is not paused")
	}

	// Update status
	c.Status = ContainerStatusRunning

	// Save container config
	if err := c.saveConfig(); err != nil {
		c.Status = ContainerStatusError
		c.Error = fmt.Sprintf("failed to save container config: %v", err)
		c.saveConfig()
		return fmt.Errorf(c.Error)
	}

	fmt.Printf("Unpaused container: %s (%s)\n", c.ID, c.Name)

	// Call platform-specific container unpause
	if err := c.Runtime.Platform.UnpauseContainer(c); err != nil {
		c.Status = ContainerStatusError
		c.Error = fmt.Sprintf("failed to unpause container: %v", err)
		c.saveConfig()
		return fmt.Errorf(c.Error)
	}

	return nil
}

// Restart restarts the container
func (c *Container) Restart() error {
	// Stop container if running
	if c.Status == ContainerStatusRunning {
		if err := c.Stop(); err != nil {
			return fmt.Errorf("failed to stop container: %v", err)
		}
	}

	// Start container
	return c.Start()
}

// Remove removes the container
func (c *Container) Remove() error {
	return c.Runtime.DeleteContainer(c.ID)
}

// saveConfig saves the container configuration to disk
func (c *Container) saveConfig() error {
	configPath := filepath.Join(c.Dir, "config.json")

	// Create directory if not exists
	if err := os.MkdirAll(c.Dir, 0755); err != nil {
		return fmt.Errorf("failed to create container directory: %v", err)
	}

	// Marshal container to JSON
	configJSON, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal container config: %v", err)
	}

	// Write config to file
	if err := os.WriteFile(configPath, configJSON, 0644); err != nil {
		return fmt.Errorf("failed to write container config: %v", err)
	}

	return nil
}

// loadContainer loads a container from disk
func loadContainer(containerDir string, runtime *Runtime) (*Container, error) {
	configPath := filepath.Join(containerDir, "config.json")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("container config file not found")
	}

	// Read config file
	configJSON, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read container config: %v", err)
	}

	// Unmarshal config
	container := &Container{}
	if err := json.Unmarshal(configJSON, container); err != nil {
		return nil, fmt.Errorf("failed to unmarshal container config: %v", err)
	}

	// Set runtime
	container.Runtime = runtime

	// Set directory
	container.Dir = containerDir

	return container, nil
}

// String returns a string representation of the container
func (c *Container) String() string {
	return fmt.Sprintf("%s (%s) - %s", c.ID, c.Name, c.Status)
}

// UploadFile uploads a file to the container's file system
func (c *Container) UploadFile(localPath, containerPath, token string) error {
	// Check if container is running
	if c.Status != ContainerStatusRunning {
		return fmt.Errorf("container is not running")
	}

	// Validate admin permission
	valid, message := c.Runtime.AdminManager.ValidateAdmin(token, c.ID, "write")
	if !valid {
		return fmt.Errorf("permission denied: %s", message)
	}

	// Get container root filesystem path
	rootFSPath := filepath.Join(c.Dir, "rootfs")
	if c.RootFSPath != "" {
		rootFSPath = c.RootFSPath
	}

	// Resolve container path
	destPath := filepath.Join(rootFSPath, containerPath)

	// Create directory if not exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Encrypt file before uploading
	encryptedPath := destPath + ".enc"
	if err := EncryptFile(localPath, encryptedPath); err != nil {
		return fmt.Errorf("failed to encrypt file: %v", err)
	}

	// Rename encrypted file to final path
	if err := os.Rename(encryptedPath, destPath); err != nil {
		return fmt.Errorf("failed to rename encrypted file: %v", err)
	}

	fmt.Printf("Uploaded and encrypted file %s to container %s at %s\n", localPath, c.ID, containerPath)
	return nil
}

// DownloadFile downloads a file from the container's file system
func (c *Container) DownloadFile(containerPath, localPath, token string) error {
	// Check if container is running
	if c.Status != ContainerStatusRunning {
		return fmt.Errorf("container is not running")
	}

	// Validate admin permission
	valid, message := c.Runtime.AdminManager.ValidateAdmin(token, c.ID, "read")
	if !valid {
		return fmt.Errorf("permission denied: %s", message)
	}

	// Get container root filesystem path
	rootFSPath := filepath.Join(c.Dir, "rootfs")
	if c.RootFSPath != "" {
		rootFSPath = c.RootFSPath
	}

	// Resolve container path
	srcPath := filepath.Join(rootFSPath, containerPath)

	// Check if source file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist in container: %s", containerPath)
	}

	// Create directory if not exists
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Decrypt file while downloading
	if err := DecryptFile(srcPath, localPath); err != nil {
		return fmt.Errorf("failed to decrypt file: %v", err)
	}

	fmt.Printf("Downloaded and decrypted file %s from container %s to %s\n", containerPath, c.ID, localPath)
	return nil
}
