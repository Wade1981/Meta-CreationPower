// Package elr implements the core runtime for Enlightenment Lighthouse Runtime
package elr

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"gopkg.in/yaml.v2"
)

// Runtime represents the core runtime of ELR
type Runtime struct {
	Config         *Config
	Platform       Platform
	Containers     map[string]*Container
	ContainerMutex sync.RWMutex
	Plugins        map[string]Plugin
	PluginMutex    sync.RWMutex
	NetworkManager *NetworkManager
	TokenManager   *TokenManager
	AdminManager   *AdminManager
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
			UseJobObjects bool   `yaml:"use_job_objects"`
			UseWSL        bool   `yaml:"use_wsl"`
			UseContainers bool   `yaml:"use_containers"`
			IsolationType string `yaml:"isolation_type"` // Options: "windows-container", "wsl", "basic"
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
		// API服务端口配置
		APIPorts struct {
			DesktopAPI int `yaml:"desktop_api"`
			PublicAPI  int `yaml:"public_api"`
			ModelAPI   int `yaml:"model_api"`
		} `yaml:"api_ports"`
	} `yaml:"network"`
	Storage struct {
		Enable  bool   `yaml:"enable"`
		Driver  string `yaml:"driver"`
		BaseDir string `yaml:"base_dir"`
	} `yaml:"storage"`
	FileDirectories map[string]string `yaml:"file_directories"` // File type to directory mapping
	Languages map[string]struct {
		Enable  bool   `yaml:"enable"`
		Runtime string `yaml:"runtime"`
	} `yaml:"languages"`
	PythonVersions map[string]string `yaml:"python_versions"` // Python version to installation path mapping
	Resources struct {
		Types map[string]struct {
			Enable bool   `yaml:"enable"`
			Dir    string `yaml:"dir"`
		} `yaml:"types"`
		ModelTypes map[string]struct {
			Enable bool   `yaml:"enable"`
			Dir    string `yaml:"dir"`
		} `yaml:"model_types"`
	} `yaml:"resources"`
}

// Platform represents the platform abstraction layer
type Platform interface {
	Name() string
	Version() string
	Init() error
	Cleanup() error
	// Container isolation methods
	CreateContainer(c *Container) error
	StartContainer(c *Container) error
	StopContainer(c *Container) error
	PauseContainer(c *Container) error
	UnpauseContainer(c *Container) error
	DeleteContainer(c *Container) error
	// 文件系统隔离方法
	SetupFileSystemIsolation(c *Container) error
	MountFileSystem(c *Container) error
	UnmountFileSystem(c *Container) error
}

// Plugin represents a runtime plugin
type Plugin interface {
	Name() string
	Version() string
	Init(runtime *Runtime) error
	Cleanup() error
}

// WindowsConfig represents the Windows platform configuration
type WindowsConfig struct {
	UseJobObjects bool   `yaml:"use_job_objects"`
	UseWSL        bool   `yaml:"use_wsl"`
	UseContainers bool   `yaml:"use_containers"`
	IsolationType string `yaml:"isolation_type"` // Options: "windows-container", "wsl", "basic"
}

// WindowsPlatform implements the Platform interface for Windows
type WindowsPlatform struct {
	Config *WindowsConfig
}

// Name returns the platform name
func (p *WindowsPlatform) Name() string {
	return "windows"
}

// Version returns the platform version
func (p *WindowsPlatform) Version() string {
	var versionInfo struct {
		OSVersionInfoSize uint32
		MajorVersion      uint32
		MinorVersion      uint32
		BuildNumber       uint32
		PlatformID        uint32
		CSDVersion        [128]uint16
	}

	versionInfo.OSVersionInfoSize = uint32(unsafe.Sizeof(versionInfo))
	syscall.Syscall(syscall.NewLazyDLL("kernel32.dll").NewProc("GetVersionExW").Addr(), 1, uintptr(unsafe.Pointer(&versionInfo)), 0, 0)

	return fmt.Sprintf("%d.%d.%d", versionInfo.MajorVersion, versionInfo.MinorVersion, versionInfo.BuildNumber)
}

// Init initializes the platform
func (p *WindowsPlatform) Init() error {
	fmt.Println("Initializing Windows platform...")

	// Check if we're running as administrator
	if !p.isAdmin() {
		fmt.Println("Warning: Not running as administrator, some features may be limited")
	}

	// Check if WSL is available (if enabled)
	if p.Config.UseWSL {
		if err := p.checkWSL(); err != nil {
			fmt.Printf("Warning: WSL not available: %v\n", err)
		}
	}

	fmt.Println("Windows platform initialized successfully!")
	return nil
}

// Cleanup cleans up the platform
func (p *WindowsPlatform) Cleanup() error {
	fmt.Println("Cleaning up Windows platform...")
	// No cleanup needed for Windows platform
	return nil
}

// CreateContainer creates a new container for Windows
func (p *WindowsPlatform) CreateContainer(c *Container) error {
	fmt.Printf("Creating container %s on Windows...\n", c.ID)
	
	// Create container directory structure
	if err := p.createContainerDirs(c); err != nil {
		return fmt.Errorf("failed to create container directories: %v", err)
	}
	
	// Create file system isolation based on configuration
	isolationType := p.Config.IsolationType
	switch isolationType {
	case "windows-container":
		// Use Windows Container isolation
		if err := p.setupWindowsContainerIsolation(c); err != nil {
			fmt.Printf("Warning: Failed to setup Windows Container isolation: %v\n", err)
			// Fallback to WSL isolation
			if err := p.setupWSLIsolation(c); err != nil {
				fmt.Printf("Warning: Failed to setup WSL isolation: %v\n", err)
				// Fallback to AppContainers isolation
				if err := p.setupAppContainerIsolation(c); err != nil {
					fmt.Printf("Warning: Failed to setup AppContainers isolation: %v\n", err)
					// Fallback to basic isolation
					if err := p.setupFileSystemIsolation(c); err != nil {
						return fmt.Errorf("failed to setup file system isolation: %v", err)
					}
				}
			}
		}
	case "wsl":
		// Use WSL isolation
		if err := p.setupWSLIsolation(c); err != nil {
			fmt.Printf("Warning: Failed to setup WSL isolation: %v\n", err)
			// Fallback to AppContainers isolation
			if err := p.setupAppContainerIsolation(c); err != nil {
				fmt.Printf("Warning: Failed to setup AppContainers isolation: %v\n", err)
				// Fallback to basic isolation
				if err := p.setupFileSystemIsolation(c); err != nil {
					return fmt.Errorf("failed to setup file system isolation: %v", err)
				}
			}
		}
	case "appcontainer":
		// Use AppContainers isolation
		if err := p.setupAppContainerIsolation(c); err != nil {
			fmt.Printf("Warning: Failed to setup AppContainers isolation: %v\n", err)
			// Fallback to basic isolation
			if err := p.setupFileSystemIsolation(c); err != nil {
				return fmt.Errorf("failed to setup file system isolation: %v", err)
			}
		}
	default:
		// Use basic file system isolation
		if err := p.setupFileSystemIsolation(c); err != nil {
			return fmt.Errorf("failed to setup file system isolation: %v", err)
		}
	}
	
	// Set up network isolation
	if err := p.setupNetworkIsolation(c); err != nil {
		return fmt.Errorf("failed to setup network isolation: %v", err)
	}
	
	return nil
}

// createContainerDirs creates the directory structure for a container
func (p *WindowsPlatform) createContainerDirs(c *Container) error {
	// Create root directory
	rootDir := filepath.Join(c.Dir, "rootfs")
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return err
	}
	
	// Create standard directories
	dirs := []string{
		"bin", "etc", "home", "lib", "lib64", "proc", "sys", "tmp", "usr", "var",
	}
	
	for _, dir := range dirs {
		dirPath := filepath.Join(rootDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}
	
	return nil
}

// setupFileSystemIsolation sets up file system isolation for a container
func (p *WindowsPlatform) setupFileSystemIsolation(c *Container) error {
	// For Windows, we'll create a basic file system structure
	// In a real implementation, this would use Windows-specific isolation mechanisms
	rootDir := filepath.Join(c.Dir, "rootfs")
	
	// Create a simple hosts file
	hostsFile := filepath.Join(rootDir, "etc", "hosts")
	hostsContent := "127.0.0.1 localhost\n::1 localhost\n"
	if err := os.WriteFile(hostsFile, []byte(hostsContent), 0644); err != nil {
		return err
	}
	
	// Create a simple passwd file
	passwdFile := filepath.Join(rootDir, "etc", "passwd")
	passwdContent := "root:x:0:0:root:/root:/bin/sh\n"
	if err := os.WriteFile(passwdFile, []byte(passwdContent), 0644); err != nil {
		return err
	}
	
	return nil
}

// StartContainer starts a container for Windows
func (p *WindowsPlatform) StartContainer(c *Container) error {
	fmt.Printf("Starting container %s on Windows...\n", c.ID)
	
	// Start container based on configuration
	isolationType := p.Config.IsolationType
	switch isolationType {
	case "windows-container":
		// Use Windows Container isolation
		if err := p.startWindowsContainer(c); err != nil {
			fmt.Printf("Warning: Failed to start Windows Container: %v\n", err)
			// Fallback to WSL container
			if err := p.startWSLContainer(c); err != nil {
				fmt.Printf("Warning: Failed to start WSL container: %v\n", err)
				// Fallback to AppContainer
				if err := p.startAppContainer(c); err != nil {
					fmt.Printf("Warning: Failed to start AppContainer: %v\n", err)
					// Fallback to basic process creation
					// Create job object for process isolation
					jobHandle, err := p.createJobObject()
					if err != nil {
						return fmt.Errorf("failed to create job object: %v", err)
					}
					
					// Set job object limits
					if err := p.setJobObjectLimits(jobHandle, 100, 512); err != nil { // 100% CPU, 512MB memory
						return fmt.Errorf("failed to set job object limits: %v", err)
					}
					
					// Create process for container
					if err := p.createContainerProcess(c, jobHandle); err != nil {
						return fmt.Errorf("failed to create container process: %v", err)
					}
				}
			}
		}
	case "wsl":
		// Use WSL isolation
		if err := p.startWSLContainer(c); err != nil {
			fmt.Printf("Warning: Failed to start WSL container: %v\n", err)
			// Fallback to AppContainer
			if err := p.startAppContainer(c); err != nil {
				fmt.Printf("Warning: Failed to start AppContainer: %v\n", err)
				// Fallback to basic process creation
				// Create job object for process isolation
				jobHandle, err := p.createJobObject()
				if err != nil {
					return fmt.Errorf("failed to create job object: %v", err)
				}
				
				// Set job object limits
				if err := p.setJobObjectLimits(jobHandle, 100, 512); err != nil { // 100% CPU, 512MB memory
					return fmt.Errorf("failed to set job object limits: %v", err)
				}
				
				// Create process for container
				if err := p.createContainerProcess(c, jobHandle); err != nil {
					return fmt.Errorf("failed to create container process: %v", err)
				}
			}
		}
	case "appcontainer":
		// Use AppContainers isolation
		if err := p.startAppContainer(c); err != nil {
			fmt.Printf("Warning: Failed to start AppContainer: %v\n", err)
			// Fallback to basic process creation
			// Create job object for process isolation
			jobHandle, err := p.createJobObject()
			if err != nil {
				return fmt.Errorf("failed to create job object: %v", err)
			}
			
			// Set job object limits
			if err := p.setJobObjectLimits(jobHandle, 100, 512); err != nil { // 100% CPU, 512MB memory
				return fmt.Errorf("failed to set job object limits: %v", err)
			}
			
			// Create process for container
			if err := p.createContainerProcess(c, jobHandle); err != nil {
				return fmt.Errorf("failed to create container process: %v", err)
			}
		}
	default:
		// Use basic process creation
		// Create job object for process isolation
		jobHandle, err := p.createJobObject()
		if err != nil {
			return fmt.Errorf("failed to create job object: %v", err)
		}
		
		// Set job object limits
		if err := p.setJobObjectLimits(jobHandle, 100, 512); err != nil { // 100% CPU, 512MB memory
			return fmt.Errorf("failed to set job object limits: %v", err)
		}
		
		// Create process for container
		if err := p.createContainerProcess(c, jobHandle); err != nil {
			return fmt.Errorf("failed to create container process: %v", err)
		}
	}
	
	return nil
}

// StopContainer stops a container for Windows
func (p *WindowsPlatform) StopContainer(c *Container) error {
	fmt.Printf("Stopping container %s on Windows...\n", c.ID)
	
	// Stop container based on configuration
	isolationType := p.Config.IsolationType
	switch isolationType {
	case "windows-container":
		// Use Windows Container isolation
		return p.stopWindowsContainer(c)
	case "wsl":
		// Use WSL isolation
		// Stop WSL container using wsl.exe
		cmd := exec.Command("wsl", "--terminate", c.ID)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Warning: Failed to stop WSL container: %v\n", err)
			fmt.Printf("Output: %s\n", string(output))
		}
	case "appcontainer":
		// Use AppContainers isolation
		// Stop AppContainer process
		if c.PID > 0 {
			process, err := os.FindProcess(c.PID)
			if err == nil {
				process.Kill()
			}
		}
	default:
		// For basic isolation, we'll just simulate stopping the process
		// In a real implementation, we would terminate the job object
	}
	
	return nil
}

// PauseContainer pauses a container for Windows
func (p *WindowsPlatform) PauseContainer(c *Container) error {
	fmt.Printf("Pausing container %s on Windows...\n", c.ID)
	// TODO: Implement Windows-specific container pause
	return nil
}

// UnpauseContainer unpauses a container for Windows
func (p *WindowsPlatform) UnpauseContainer(c *Container) error {
	fmt.Printf("Unpausing container %s on Windows...\n", c.ID)
	// TODO: Implement Windows-specific container unpause
	return nil
}

// DeleteContainer deletes a container for Windows
func (p *WindowsPlatform) DeleteContainer(c *Container) error {
	fmt.Printf("Deleting container %s on Windows...\n", c.ID)
	
	// Delete container based on configuration
	isolationType := p.Config.IsolationType
	switch isolationType {
	case "windows-container":
		// Use Windows Container isolation
		return p.deleteWindowsContainer(c)
	case "wsl":
		// Use WSL isolation
		// Unregister WSL distribution using wsl.exe
		cmd := exec.Command("wsl", "--unregister", c.ID)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Warning: Failed to unregister WSL distribution: %v\n", err)
			fmt.Printf("Output: %s\n", string(output))
		}
	case "appcontainer":
		// Use AppContainers isolation
		// Delete AppContainer using PowerShell
		psScript := fmt.Sprintf(`
		# Delete AppContainer
		$appContainerName = "%s"
		Remove-AppContainer -Name $appContainerName -ErrorAction SilentlyContinue
		`, c.ID)
		
		// Save PowerShell script to file
		psScriptPath := filepath.Join(c.Dir, "delete-appcontainer.ps1")
		if err := os.WriteFile(psScriptPath, []byte(psScript), 0644); err != nil {
			fmt.Printf("Warning: Failed to write PowerShell script: %v\n", err)
		}
		
		// Execute PowerShell script
		cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", psScriptPath)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Warning: Failed to delete AppContainer: %v\n", err)
			fmt.Printf("Output: %s\n", string(output))
		}
	default:
		// For basic isolation, we'll just simulate deleting the container
		// In a real implementation, we would delete the container directory
	}
	
	// Delete container directory
	if err := os.RemoveAll(c.Dir); err != nil {
		fmt.Printf("Warning: Failed to delete container directory: %v\n", err)
	}
	
	return nil
}

// SetupFileSystemIsolation sets up file system isolation for a container
func (p *WindowsPlatform) SetupFileSystemIsolation(c *Container) error {
	// For Windows, we'll create a basic file system structure
	// In a real implementation, this would use Windows-specific isolation mechanisms
	rootDir := filepath.Join(c.Dir, "rootfs")
	
	// Create a simple hosts file
	hostsFile := filepath.Join(rootDir, "etc", "hosts")
	hostsContent := "127.0.0.1 localhost\n::1 localhost\n"
	if err := os.WriteFile(hostsFile, []byte(hostsContent), 0644); err != nil {
		return err
	}
	
	// Create a simple passwd file
	passwdFile := filepath.Join(rootDir, "etc", "passwd")
	passwdContent := "root:x:0:0:root:/root:/bin/sh\n"
	if err := os.WriteFile(passwdFile, []byte(passwdContent), 0644); err != nil {
		return err
	}
	
	return nil
}

// MountFileSystem mounts the file system for a container
func (p *WindowsPlatform) MountFileSystem(c *Container) error {
	// For Windows, we'll just simulate mounting the file system
	// In a real implementation, this would use Windows-specific mounting mechanisms
	return nil
}

// UnmountFileSystem unmounts the file system for a container
func (p *WindowsPlatform) UnmountFileSystem(c *Container) error {
	// For Windows, we'll just simulate unmounting the file system
	// In a real implementation, this would use Windows-specific unmounting mechanisms
	return nil
}

// createJobObject creates a new job object for process isolation
func (p *WindowsPlatform) createJobObject() (syscall.Handle, error) {
	// Simplified implementation to avoid undefined constants
	// In a real implementation, we would use the appropriate Windows API calls
	// For now, we'll just return a dummy handle
	return 0, nil
}

// setJobObjectLimits sets limits for a job object
func (p *WindowsPlatform) setJobObjectLimits(jobHandle syscall.Handle, cpuLimit, memLimit int) error {
	// Simplified implementation to avoid undefined constants
	// In a real implementation, we would use the appropriate Windows API calls
	// For now, we'll just log the limits
	fmt.Printf("Setting job object limits: CPU=%d%%, Memory=%dMB\n", cpuLimit, memLimit)
	return nil
}

// createContainerProcess creates a new process for the container
func (p *WindowsPlatform) createContainerProcess(c *Container, jobHandle syscall.Handle) error {
	// Get command to run
	command := c.Command
	if command == "" {
		command = "cmd.exe"
	}
	
	// Create process using os/exec instead of syscall to avoid undefined constants
	cmd := exec.Command(command, c.Args...)
	
	// Start the process
	if err := cmd.Start(); err != nil {
		return err
	}
	
	// Store the process ID in the container
	c.PID = cmd.Process.Pid
	
	return nil
}

// setupNetworkIsolation sets up network isolation for a container
func (p *WindowsPlatform) setupNetworkIsolation(c *Container) error {
	// Only set up network isolation if container network is enabled
	if !c.NetworkEnabled {
		fmt.Printf("Container %s network is disabled, skipping network isolation setup\n", c.ID)
		return nil
	}

	fmt.Printf("Setting up network isolation for container %s...\n", c.ID)
	
	// For Windows, we'll create a basic network configuration
	// In a real implementation, this would use Windows-specific network isolation mechanisms
	
	// Create network configuration directory
	networkDir := filepath.Join(c.Dir, "network")
	if err := os.MkdirAll(networkDir, 0755); err != nil {
		return err
	}
	
	// Create network configuration file
	networkConfig := filepath.Join(networkDir, "config.json")
	configContent := map[string]interface{}{
		"container_id": c.ID,
		"network_mode": "nat",
		"ip_address": fmt.Sprintf("172.16.0.%d", (c.PID % 254) + 1),
		"port_mappings": []map[string]interface{}{},
	}
	
	configJSON, err := json.MarshalIndent(configContent, "", "  ")
	if err != nil {
		return err
	}
	
	if err := os.WriteFile(networkConfig, configJSON, 0644); err != nil {
		return err
	}
	
	return nil
}

// isAdmin checks if the current process is running as administrator
func (p *WindowsPlatform) isAdmin() bool {
	// Simplified implementation to avoid undefined constants
	// In a real implementation, we would use the appropriate Windows API calls
	// For now, we'll just return false as a placeholder
	return false
}

// checkWSL checks if WSL is available
func (p *WindowsPlatform) checkWSL() error {
	// Run wsl --list to check if WSL is available
	cmd := exec.Command("wsl", "--list")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("WSL not available: %v", err)
	}

	return nil
}

// checkWindowsContainer checks if Windows Container feature is available
func (p *WindowsPlatform) checkWindowsContainer() error {
	// Run powershell command to check if Containers feature is enabled
	cmd := exec.Command("powershell", "Get-WindowsOptionalFeature", "-FeatureName", "Containers", "-Online")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Windows Container feature not available: %v", err)
	}

	return nil
}

// setupWindowsContainerIsolation sets up Windows Container isolation
func (p *WindowsPlatform) setupWindowsContainerIsolation(c *Container) error {
	fmt.Printf("Setting up Windows Container isolation for container %s...\n", c.ID)
	
	// Check if Windows Container feature is available
	if err := p.checkWindowsContainer(); err != nil {
		return fmt.Errorf("Windows Container feature not available: %v", err)
	}
	
	// Create Windows Container configuration
	containerConfig := map[string]interface{}{
		"id":      c.ID,
		"name":    c.Name,
		"image":   c.Image,
		"command": c.Command,
		"args":    c.Args,
		"env":     c.Env,
	}
	
	// Save container configuration
	configPath := filepath.Join(c.Dir, "container.json")
	configJSON, err := json.MarshalIndent(containerConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal container config: %v", err)
	}
	
	if err := os.WriteFile(configPath, configJSON, 0644); err != nil {
		return fmt.Errorf("failed to write container config: %v", err)
	}
	
	// Create Windows Container using PowerShell and Windows Container API
	// Note: This uses the Windows Container API directly without Docker dependency
	psScript := fmt.Sprintf(`
	# Import the container module
	Import-Module ContainerImage -ErrorAction SilentlyContinue
	
	# Create a new container
	$container = New-Container -Name "%s" -ContainerImageName "%s" -ErrorAction Stop
	
	# Start the container
	Start-Container -Container $container -ErrorAction Stop
	
	# Get container ID
	$containerId = $container.Id
	
	# Return the container ID
	Write-Output $containerId
	`, c.ID, c.Image)
	
	// Save PowerShell script to file
	psScriptPath := filepath.Join(c.Dir, "create-container.ps1")
	if err := os.WriteFile(psScriptPath, []byte(psScript), 0644); err != nil {
		return fmt.Errorf("failed to write PowerShell script: %v", err)
	}
	
	// Execute PowerShell script
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", psScriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Warning: Failed to create Windows Container: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		// Fallback to WSL isolation if Windows Container is not available
		return p.setupWSLIsolation(c)
	}
	
	return nil
}

// startWindowsContainer starts a Windows Container
func (p *WindowsPlatform) startWindowsContainer(c *Container) error {
	fmt.Printf("Starting Windows Container %s...\n", c.ID)
	
	// Start Windows Container using PowerShell and Windows Container API
	psScript := fmt.Sprintf(`
	# Import the container module
	Import-Module ContainerImage -ErrorAction SilentlyContinue
	
	# Get the container
	$container = Get-Container -Name "%s" -ErrorAction Stop
	
	# Start the container
	Start-Container -Container $container -ErrorAction Stop
	
	# Get container process ID
	$containerProcess = Get-ContainerProcess -Container $container -ErrorAction Stop
	$pid = $containerProcess.ProcessId
	
	# Return the PID
	Write-Output $pid
	`, c.ID)
	
	// Save PowerShell script to file
	psScriptPath := filepath.Join(c.Dir, "start-container.ps1")
	if err := os.WriteFile(psScriptPath, []byte(psScript), 0644); err != nil {
		return fmt.Errorf("failed to write PowerShell script: %v", err)
	}
	
	// Execute PowerShell script
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", psScriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Warning: Failed to start Windows Container: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		// Fallback to WSL isolation if Windows Container is not available
		return p.startWSLContainer(c)
	}
	
	// Parse PID
	var pid int
	if _, err := fmt.Sscanf(string(output), "%d", &pid); err != nil {
		fmt.Printf("Warning: Failed to parse container PID: %v\n", err)
		// Fallback to WSL isolation if Windows Container is not available
		return p.startWSLContainer(c)
	}
	
	c.PID = pid
	return nil
}

// stopWindowsContainer stops a Windows Container
func (p *WindowsPlatform) stopWindowsContainer(c *Container) error {
	fmt.Printf("Stopping Windows Container %s...\n", c.ID)
	
	// Stop Windows Container using PowerShell and Windows Container API
	psScript := fmt.Sprintf(`
	# Import the container module
	Import-Module ContainerImage -ErrorAction SilentlyContinue
	
	# Get the container
	$container = Get-Container -Name "%s" -ErrorAction Stop
	
	# Stop the container
	Stop-Container -Container $container -ErrorAction Stop
	`, c.ID)
	
	// Save PowerShell script to file
	psScriptPath := filepath.Join(c.Dir, "stop-container.ps1")
	if err := os.WriteFile(psScriptPath, []byte(psScript), 0644); err != nil {
		return fmt.Errorf("failed to write PowerShell script: %v", err)
	}
	
	// Execute PowerShell script
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", psScriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Warning: Failed to stop Windows Container: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
	}
	
	return nil
}

// deleteWindowsContainer deletes a Windows Container
func (p *WindowsPlatform) deleteWindowsContainer(c *Container) error {
	fmt.Printf("Deleting Windows Container %s...\n", c.ID)
	
	// Delete Windows Container using PowerShell and Windows Container API
	psScript := fmt.Sprintf(`
	# Import the container module
	Import-Module ContainerImage -ErrorAction SilentlyContinue
	
	# Get the container
	$container = Get-Container -Name "%s" -ErrorAction Stop
	
	# Stop the container if it's running
	if ($container.State -eq "Running") {
		Stop-Container -Container $container -ErrorAction Stop
	}
	
	# Delete the container
	Remove-Container -Container $container -ErrorAction Stop
	`, c.ID)
	
	// Save PowerShell script to file
	psScriptPath := filepath.Join(c.Dir, "delete-container.ps1")
	if err := os.WriteFile(psScriptPath, []byte(psScript), 0644); err != nil {
		return fmt.Errorf("failed to write PowerShell script: %v", err)
	}
	
	// Execute PowerShell script
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", psScriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Warning: Failed to delete Windows Container: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
	}
	
	return nil
}

// setupWSLIsolation sets up WSL isolation
func (p *WindowsPlatform) setupWSLIsolation(c *Container) error {
	fmt.Printf("Setting up WSL isolation for container %s...\n", c.ID)
	
	// Check if WSL is available
	if err := p.checkWSL(); err != nil {
		return fmt.Errorf("WSL not available: %v", err)
	}
	
	// Create WSL configuration
	wslConfig := map[string]interface{}{
		"id":      c.ID,
		"name":    c.Name,
		"image":   c.Image,
		"command": c.Command,
		"args":    c.Args,
		"env":     c.Env,
	}
	
	// Save WSL configuration
	configPath := filepath.Join(c.Dir, "wsl-config.json")
	configJSON, err := json.MarshalIndent(wslConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal WSL config: %v", err)
	}
	
	if err := os.WriteFile(configPath, configJSON, 0644); err != nil {
		return fmt.Errorf("failed to write WSL config: %v", err)
	}
	
	// Create WSL distribution using wsl.exe
	// Note: This creates a new WSL distribution for the container
	cmd := exec.Command("wsl", "--import", c.ID, c.Dir, c.Image)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Warning: Failed to create WSL distribution: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		// Fallback to AppContainers isolation if WSL is not available
		return p.setupAppContainerIsolation(c)
	}
	
	return nil
}

// startWSLContainer starts a WSL container
func (p *WindowsPlatform) startWSLContainer(c *Container) error {
	fmt.Printf("Starting WSL container %s...\n", c.ID)
	
	// Start WSL container using wsl.exe
	cmd := exec.Command("wsl", "-d", c.ID, c.Command)
	
	// Start the process
	if err := cmd.Start(); err != nil {
		fmt.Printf("Warning: Failed to start WSL container: %v\n", err)
		// Fallback to AppContainers isolation if WSL is not available
		return p.startAppContainer(c)
	}
	
	// Store the process ID in the container
	c.PID = cmd.Process.Pid
	
	return nil
}

// setupAppContainerIsolation sets up AppContainers isolation
func (p *WindowsPlatform) setupAppContainerIsolation(c *Container) error {
	fmt.Printf("Setting up AppContainers isolation for container %s...\n", c.ID)
	
	// Create AppContainer configuration
	appContainerConfig := map[string]interface{}{
		"id":      c.ID,
		"name":    c.Name,
		"command": c.Command,
		"args":    c.Args,
		"env":     c.Env,
	}
	
	// Save AppContainer configuration
	configPath := filepath.Join(c.Dir, "appcontainer-config.json")
	configJSON, err := json.MarshalIndent(appContainerConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal AppContainer config: %v", err)
	}
	
	if err := os.WriteFile(configPath, configJSON, 0644); err != nil {
		return fmt.Errorf("failed to write AppContainer config: %v", err)
	}
	
	// Create AppContainer using PowerShell
	psScript := fmt.Sprintf(`
	# Create AppContainer
	$appContainerName = "%s"
	$appContainerSid = New-AppContainer -Name $appContainerName -ErrorAction Stop
	
	# Return the AppContainer SID
	Write-Output $appContainerSid
	`, c.ID)
	
	// Save PowerShell script to file
	psScriptPath := filepath.Join(c.Dir, "create-appcontainer.ps1")
	if err := os.WriteFile(psScriptPath, []byte(psScript), 0644); err != nil {
		return fmt.Errorf("failed to write PowerShell script: %v", err)
	}
	
	// Execute PowerShell script
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", psScriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Warning: Failed to create AppContainer: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		// Fallback to basic isolation if AppContainers is not available
		return p.setupFileSystemIsolation(c)
	}
	
	return nil
}

// startAppContainer starts an AppContainer
func (p *WindowsPlatform) startAppContainer(c *Container) error {
	fmt.Printf("Starting AppContainer %s...\n", c.ID)
	
	// Start AppContainer using PowerShell
	psScript := fmt.Sprintf(`
	# Start process in AppContainer
	$appContainerName = "%s"
	$command = "%s"
	$args = "%s"
	
	# Start process in AppContainer
	$process = Start-Process -FilePath $command -ArgumentList $args -AppContainer $appContainerName -PassThru -ErrorAction Stop
	
	# Return the process ID
	Write-Output $process.Id
	`, c.ID, c.Command, strings.Join(c.Args, " "))
	
	// Save PowerShell script to file
	psScriptPath := filepath.Join(c.Dir, "start-appcontainer.ps1")
	if err := os.WriteFile(psScriptPath, []byte(psScript), 0644); err != nil {
		return fmt.Errorf("failed to write PowerShell script: %v", err)
	}
	
	// Execute PowerShell script
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", psScriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Warning: Failed to start AppContainer: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		// Fallback to basic process creation if AppContainers is not available
		return p.createContainerProcess(c, syscall.Handle(0))
	}
	
	// Parse PID
	var pid int
	if _, err := fmt.Sscanf(string(output), "%d", &pid); err != nil {
		fmt.Printf("Warning: Failed to parse AppContainer PID: %v\n", err)
		// Fallback to basic process creation if AppContainers is not available
		return p.createContainerProcess(c, syscall.Handle(0))
	}

	c.PID = pid
	return nil
}

// SimplePlatform is a simple platform implementation for testing
type SimplePlatform struct {
	os string
}

func (p *SimplePlatform) Name() string {
	return p.os
}

func (p *SimplePlatform) Version() string {
	return "1.0.0"
}

func (p *SimplePlatform) Init() error {
	return nil
}

func (p *SimplePlatform) Cleanup() error {
	return nil
}

// Container isolation methods for SimplePlatform
func (p *SimplePlatform) CreateContainer(c *Container) error {
	return nil
}

func (p *SimplePlatform) StartContainer(c *Container) error {
	return nil
}

func (p *SimplePlatform) StopContainer(c *Container) error {
	return nil
}

func (p *SimplePlatform) PauseContainer(c *Container) error {
	return nil
}

func (p *SimplePlatform) UnpauseContainer(c *Container) error {
	return nil
}

func (p *SimplePlatform) DeleteContainer(c *Container) error {
	return nil
}

// 文件系统隔离方法 for SimplePlatform
func (p *SimplePlatform) SetupFileSystemIsolation(c *Container) error {
	return nil
}

func (p *SimplePlatform) MountFileSystem(c *Container) error {
	return nil
}

func (p *SimplePlatform) UnmountFileSystem(c *Container) error {
	return nil
}



// NewRuntime creates a new runtime instance
func NewRuntime(config *Config) (*Runtime, error) {
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	// Initialize data directory
	if config.DataDir == "" {
		config.DataDir = filepath.Join(homeDir, ".elr", "data")
	} else if strings.HasPrefix(config.DataDir, "~") {
		// Expand ~ to home directory
		config.DataDir = filepath.Join(homeDir, config.DataDir[1:])
	}
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	// Initialize plugin directory
	if config.PluginDir == "" {
		config.PluginDir = filepath.Join(homeDir, ".elr", "plugins")
	} else if strings.HasPrefix(config.PluginDir, "~") {
		// Expand ~ to home directory
		config.PluginDir = filepath.Join(homeDir, config.PluginDir[1:])
	}
	if err := os.MkdirAll(config.PluginDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %v", err)
	}

	// Set default network configuration
	if config.Network.Enable == false {
		// 默认禁用网络
		config.Network.Enable = false
	}

	// Set default API ports to avoid conflicts
	if config.Network.APIPorts.DesktopAPI == 0 {
		config.Network.APIPorts.DesktopAPI = 8081
	}
	if config.Network.APIPorts.PublicAPI == 0 {
		config.Network.APIPorts.PublicAPI = 8080
	}
	if config.Network.APIPorts.ModelAPI == 0 {
		config.Network.APIPorts.ModelAPI = 8082
	}

	// Create platform based on OS
	var platform Platform
	switch runtime.GOOS {
	case "windows":
		// Use Windows-specific platform
		platform = &WindowsPlatform{
			Config: &WindowsConfig{
				UseJobObjects: config.Platform.Windows.UseJobObjects,
				UseWSL:        config.Platform.Windows.UseWSL,
				UseContainers: config.Platform.Windows.UseContainers,
				IsolationType: config.Platform.Windows.IsolationType,
			},
		}
	default:
		// Use simple platform for other OS
		platform = &SimplePlatform{
			os: runtime.GOOS,
		}
	}

	// Create runtime
	runtime := &Runtime{
		Config:     config,
		Platform:   platform,
		Containers: make(map[string]*Container),
		Plugins:    make(map[string]Plugin),
		StopCh:     make(chan struct{}),
	}

	// Create network manager with public API port
	runtime.NetworkManager = NewNetworkManager(runtime, config.Network.APIPorts.PublicAPI)

	// Create token manager
	tokenFile := filepath.Join(config.DataDir, "tokens.json")
	runtime.TokenManager = NewTokenManager(tokenFile)
	if err := runtime.TokenManager.LoadTokens(); err != nil {
		fmt.Printf("Warning: failed to load tokens: %v\n", err)
	}

	// Create admin manager
	adminFile := filepath.Join(config.DataDir, "admins.json")
	runtime.AdminManager = NewAdminManager(adminFile)
	if err := runtime.AdminManager.LoadAdmins(); err != nil {
		fmt.Printf("Warning: failed to load admins: %v\n", err)
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
	fmt.Printf("Network enabled: %v\n", r.Config.Network.Enable)

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

	// Start network service only if network is enabled
	if r.Config.Network.Enable {
		if err := r.NetworkManager.Start(); err != nil {
			fmt.Printf("Warning: failed to start network service: %v\n", err)
		} else {
			fmt.Println("Network service started successfully!")
		}
	} else {
		fmt.Println("Network service is disabled. Use 'elr network enable' to enable it.")
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

	// Stop network service
	if err := r.NetworkManager.Stop(); err != nil {
		fmt.Printf("Warning: failed to stop network service: %v\n", err)
	}

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

	// Generate IP address for container
	ipAddress := fmt.Sprintf("172.16.0.%d", (os.Getpid()%254)+1)

	// Create container with network disabled by default
	container := &Container{
		ID:                 config.ID,
		Name:               config.Name,
		Image:              config.Image,
		Command:            config.Command,
		Args:               config.Args,
		Env:                config.Env,
		MemoryLimit:        config.MemoryLimit,
		CPULimit:           config.CPULimit,
		NetworkMode:        config.NetworkMode,
		PortMappings:       config.PortMappings,
		IPAddress:          ipAddress,
		Dir:                containerDir,
		FileSystemIsolation: config.FileSystemIsolation,
		RootFSPath:         config.RootFSPath,
		ReadOnlyFS:         config.ReadOnlyFS,
		NetworkEnabled:     false, // 默认禁用网络
		Runtime:            r,
		Status:             ContainerStatusCreated,
		Created:            time.Now(),
	}

	// Create channel to receive creation result
	resultCh := make(chan struct {
		container *Container
		err       error
	})

	// Start container creation in a goroutine
	go func() {
		// Call platform-specific container creation
		if err := r.Platform.CreateContainer(container); err != nil {
			resultCh <- struct {
				container *Container
				err       error
			}{nil, fmt.Errorf("failed to create container on platform: %v", err)}
			return
		}

		// Save container config
		if err := container.saveConfig(); err != nil {
			resultCh <- struct {
				container *Container
				err       error
			}{nil, fmt.Errorf("failed to save container config: %v", err)}
			return
		}

		// Add container to runtime
		r.ContainerMutex.Lock()
		r.Containers[config.ID] = container
		r.ContainerMutex.Unlock()

		fmt.Printf("Created container: %s (%s)\n", container.ID, container.Name)
		fmt.Printf("Container network is disabled by default. Use 'elr container network enable %s' to enable it.\n", container.ID)
		resultCh <- struct {
			container *Container
			err       error
		}{container, nil}
	}()

	// Wait for creation to complete
	result := <-resultCh
	return result.container, result.err
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

	// Call platform-specific container deletion
	if err := r.Platform.DeleteContainer(container); err != nil {
		return fmt.Errorf("failed to delete container on platform: %v", err)
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
const Version = "1.1"

// GetVersion returns the runtime version
func GetVersion() string {
	return Version
}

// GetConfig returns the runtime configuration
func (r *Runtime) GetConfig() *Config {
	return r.Config
}

// GetFileDirectory returns the directory for a specific file type
func (r *Runtime) GetFileDirectory(fileType string) (string, error) {
	// Initialize FileDirectories if it's nil
	if r.Config.FileDirectories == nil {
		r.Config.FileDirectories = make(map[string]string)
	}
	
	// Check if directory is already set
	if dir, exists := r.Config.FileDirectories[fileType]; exists {
		// Ensure the directory exists
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %v", err)
		}
		return dir, nil
	}
	
	// Create default directory if not set
	defaultDir := filepath.Join(r.Config.DataDir, "files", fileType)
	if err := os.MkdirAll(defaultDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create default directory: %v", err)
	}
	
	// Set the default directory
	r.Config.FileDirectories[fileType] = defaultDir
	return defaultDir, nil
}

// SetFileDirectory sets the directory for a specific file type
func (r *Runtime) SetFileDirectory(fileType string, directory string) error {
	// Initialize FileDirectories if it's nil
	if r.Config.FileDirectories == nil {
		r.Config.FileDirectories = make(map[string]string)
	}
	
	// Ensure the directory exists
	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	
	// Set the directory
	r.Config.FileDirectories[fileType] = directory
	
	// Save config to file
	if err := r.SaveConfig(); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}
	
	return nil
}

// SaveConfig saves the configuration to file
func (r *Runtime) SaveConfig() error {
	configPath := os.Getenv("ELR_CONFIG")
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configPath = filepath.Join(homeDir, ".elr", "config.yaml")
	}

	// Create config directory
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Serialize config
	configBytes, err := yaml.Marshal(r.Config)
	if err != nil {
		return err
	}

	// Write config file
	return os.WriteFile(configPath, configBytes, 0644)
}

// SaveFile saves a file to the specified file type directory
func (r *Runtime) SaveFile(fileType string, fileName string, content []byte) error {
	// Get the directory for this file type
	dir, err := r.GetFileDirectory(fileType)
	if err != nil {
		return err
	}
	
	// Create the full file path
	filePath := filepath.Join(dir, fileName)
	
	// Write the file
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	
	return nil
}

// DefaultDataDir returns the default data directory
func DefaultDataDir(homeDir string) string {
	return filepath.Join(homeDir, ".elr", "data")
}

// DefaultPluginDir returns the default plugin directory
func DefaultPluginDir(homeDir string) string {
	return filepath.Join(homeDir, ".elr", "plugins")
}
