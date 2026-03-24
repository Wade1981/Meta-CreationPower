// Package windows implements the Windows platform for Enlightenment Lighthouse Runtime
package windows

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

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

// Container represents a container
type Container struct {
	ID                 string
	Name               string
	Image              string
	Command            string
	Args               []string
	Env                []string
	MemoryLimit        int
	CPULimit           int
	NetworkMode        string
	PortMappings       []string
	IPAddress          string
	Dir                string
	FileSystemIsolation bool
	RootFSPath         string
	ReadOnlyFS         bool
	PID                int
	Status             string
	Created            string
}

// NewWindowsPlatform creates a new Windows platform instance
func NewWindowsPlatform(config map[string]interface{}) (Platform, error) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf("not running on Windows")
	}

	// Extract Windows configuration from the map
	windowsConfig := &WindowsConfig{}
	if platformConfig, ok := config["platform"].(map[string]interface{}); ok {
		if windowsPlatformConfig, ok := platformConfig["windows"].(map[string]interface{}); ok {
			if useJobObjects, ok := windowsPlatformConfig["use_job_objects"].(bool); ok {
				windowsConfig.UseJobObjects = useJobObjects
			}
			if useWSL, ok := windowsPlatformConfig["use_wsl"].(bool); ok {
				windowsConfig.UseWSL = useWSL
			}
			if useContainers, ok := windowsPlatformConfig["use_containers"].(bool); ok {
				windowsConfig.UseContainers = useContainers
			}
			if isolationType, ok := windowsPlatformConfig["isolation_type"].(string); ok {
				windowsConfig.IsolationType = isolationType
			}
		}
	}

	return &WindowsPlatform{
		Config: windowsConfig,
	}, nil
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
	if p.Config.Platform.Windows.UseWSL {
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

// createJobObject creates a new job object for process isolation
func (p *WindowsPlatform) createJobObject() (syscall.Handle, error) {
	jobHandle, err := syscall.CreateJobObject(nil, nil)
	if err != nil {
		return 0, err
	}
	
	// Set basic job object information
	info := syscall.JOBOBJECT_BASIC_LIMIT_INFORMATION{
		LimitFlags: syscall.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
	}
	
	exinfo := syscall.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: info,
	}
	
	size := uint32(unsafe.Sizeof(exinfo))
	err = syscall.SetInformationJobObject(jobHandle, syscall.JobObjectExtendedLimitInformation, (*byte)(unsafe.Pointer(&exinfo)), size)
	if err != nil {
		return 0, err
	}
	
	return jobHandle, nil
}

// setJobObjectLimits sets limits for a job object
func (p *WindowsPlatform) setJobObjectLimits(jobHandle syscall.Handle, cpuLimit, memLimit int) error {
	// Set CPU limit
	cpuLimitInfo := syscall.JOBOBJECT_CPU_RATE_CONTROL_INFORMATION{
		ControlFlags: syscall.JOB_OBJECT_CPU_RATE_CONTROL_ENABLE | syscall.JOB_OBJECT_CPU_RATE_CONTROL_HARD_CAP,
		CpuRate:      uint32(cpuLimit * 100), // Convert percentage to 1/100th of a percent
	}
	
	size := uint32(unsafe.Sizeof(cpuLimitInfo))
	err := syscall.SetInformationJobObject(jobHandle, syscall.JobObjectCpuRateControlInformation, (*byte)(unsafe.Pointer(&cpuLimitInfo)), size)
	if err != nil {
		return err
	}
	
	// Set memory limit
	memLimitInfo := syscall.JOBOBJECT_BASIC_LIMIT_INFORMATION{
		LimitFlags: syscall.JOB_OBJECT_LIMIT_PROCESS_MEMORY | syscall.JOB_OBJECT_LIMIT_WORKINGSET,
		ProcessMemoryLimit: uintptr(memLimit * 1024 * 1024), // Convert MB to bytes
		MinimumWorkingSetSize: 0,
		MaximumWorkingSetSize: uintptr(memLimit * 1024 * 1024),
	}
	
	exinfo := syscall.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: memLimitInfo,
	}
	
	size = uint32(unsafe.Sizeof(exinfo))
	err = syscall.SetInformationJobObject(jobHandle, syscall.JobObjectExtendedLimitInformation, (*byte)(unsafe.Pointer(&exinfo)), size)
	if err != nil {
		return err
	}
	
	// Set process limit
	processLimitInfo := syscall.JOBOBJECT_BASIC_LIMIT_INFORMATION{
		LimitFlags: syscall.JOB_OBJECT_LIMIT_ACTIVE_PROCESS,
		ActiveProcessLimit: 10, // Limit to 10 processes per container
	}
	
	exinfo2 := syscall.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: processLimitInfo,
	}
	
	size = uint32(unsafe.Sizeof(exinfo2))
	err = syscall.SetInformationJobObject(jobHandle, syscall.JobObjectExtendedLimitInformation, (*byte)(unsafe.Pointer(&exinfo2)), size)
	if err != nil {
		return err
	}
	
	return nil
}

// createContainerProcess creates a new process for the container
func (p *WindowsPlatform) createContainerProcess(c *Container, jobHandle syscall.Handle) error {
	// Get command to run
	command := c.Command
	if command == "" {
		command = "cmd.exe"
	}
	
	// Set up process creation flags
	flags := uint32(syscall.CREATE_NEW_CONSOLE | syscall.CREATE_SUSPENDED)
	
	// Create process
	var si syscall.StartupInfo
	var pi syscall.ProcessInformation
	
	err := syscall.CreateProcess(nil, &command[0], nil, nil, false, flags, nil, nil, &si, &pi)
	if err != nil {
		return err
	}
	defer syscall.CloseHandle(pi.Thread)
	
	// Assign process to job object
	err = syscall.AssignProcessToJobObject(jobHandle, pi.Process)
	if err != nil {
		// If assignment fails, terminate the process
		syscall.TerminateProcess(pi.Process, 1)
		syscall.CloseHandle(pi.Process)
		return err
	}
	
	// Resume the process
	syscall.ResumeThread(pi.Thread)
	
	// Store the process ID in the container
	c.PID = int(pi.ProcessId)
	
	// Close the process handle
	syscall.CloseHandle(pi.Process)
	
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

// setupNetworkIsolation sets up network isolation for a container
func (p *WindowsPlatform) setupNetworkIsolation(c *Container) error {
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
	var tokenHandle syscall.Token
	if err := syscall.OpenProcessToken(syscall.GetCurrentProcess(), syscall.TOKEN_QUERY, &tokenHandle); err != nil {
		return false
	}
	defer syscall.CloseHandle(tokenHandle)

	var elevation uint32
	if err := syscall.GetTokenInformation(tokenHandle, syscall.TokenElevation, (*byte)(unsafe.Pointer(&elevation)), uint32(unsafe.Sizeof(elevation)), nil); err != nil {
		return false
	}

	return elevation != 0
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
func (p *WindowsPlatform) setupWindowsContainerIsolation(c *elr.Container) error {
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
func (p *WindowsPlatform) startWindowsContainer(c *elr.Container) error {
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
func (p *WindowsPlatform) stopWindowsContainer(c *elr.Container) error {
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
func (p *WindowsPlatform) deleteWindowsContainer(c *elr.Container) error {
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
func (p *WindowsPlatform) setupWSLIsolation(c *elr.Container) error {
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
func (p *WindowsPlatform) startWSLContainer(c *elr.Container) error {
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
func (p *WindowsPlatform) setupAppContainerIsolation(c *elr.Container) error {
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
func (p *WindowsPlatform) startAppContainer(c *elr.Container) error {
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

// SetupFileSystemIsolation sets up file system isolation for a container
func (p *WindowsPlatform) SetupFileSystemIsolation(c *elr.Container) error {
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
func (p *WindowsPlatform) MountFileSystem(c *elr.Container) error {
	// For Windows, we'll just simulate mounting the file system
	// In a real implementation, this would use Windows-specific mounting mechanisms
	return nil
}

// UnmountFileSystem unmounts the file system for a container
func (p *WindowsPlatform) UnmountFileSystem(c *elr.Container) error {
	// For Windows, we'll just simulate unmounting the file system
	// In a real implementation, this would use Windows-specific unmounting mechanisms
	return nil
}


