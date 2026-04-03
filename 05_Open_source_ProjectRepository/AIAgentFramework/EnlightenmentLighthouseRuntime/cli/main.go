// Package main implements the command-line interface for Enlightenment Lighthouse Runtime
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"elr.local"
	"api"
	"micro_model/config"
	"micro_model/container"
	"micro_model/model"
	"micro_model/monitor"
	"micro_model/sandbox"
	"gopkg.in/yaml.v2"
)

func main() {
	// Parse command-line arguments
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "version":
		printVersion()
	case "help":
		printHelp()
	case "start":
		startRuntime()
	case "stop":
		stopRuntime()
	case "create":
		createContainer()
	case "run":
		if len(os.Args) >= 3 && os.Args[2] == "python" {
			runPython()
		} else {
			runContainer()
		}
	case "install":
		if len(os.Args) >= 3 && os.Args[2] == "python" {
			installPython()
		} else {
			fmt.Println("Error: Unknown install command")
			fmt.Println("Usage: elr install python [version] [path]")
			os.Exit(1)
		}
	case "start-container":
		startContainer()
	case "stop-container":
		stopContainer()
	case "list":
		listContainers()
	case "delete":
		deleteContainer()
	case "inspect":
		inspectContainer()
	// 系统设置命令
	case "setup":
		setupCommand()
	// 资源配置命令
	case "Settings":
		settingsCommand()
	// 模型管理命令
	case "model":
		modelCommand()
	// 沙箱管理命令
	case "sandbox":
		sandboxCommand()
	// API 服务命令
	case "api":
		apiCommand()
	// 文件系统管理命令
	case "fs":
		fsCommand()
	// 管理员管理命令
	case "admin":
		adminCommand()
	// 上传命令
	case "Upload":
		uploadCommand()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

// runPython runs a Python script
func runPython() {
	if len(os.Args) < 4 {
		fmt.Println("Error: Python script path is required")
		fmt.Println("Usage: elr run python <script.py>")
		os.Exit(1)
	}

	scriptPath := os.Args[3]
	
	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		fmt.Printf("Error: Script file not found: %s\n", scriptPath)
		os.Exit(1)
	}

	// Load config to get Python path
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Get Python runtime path
	pythonRuntime := config.Languages["python"].Runtime
	
	// If Python path is not set in config, use default
	if pythonRuntime == "" {
		pythonRuntime = "python"
	}

	// Check if Python is available
	_, err = exec.LookPath(pythonRuntime)
	if err != nil {
		// Try to use Python from D:\Python
		pythonRuntime = "D:\\Python\\python.exe"
		if _, err := os.Stat(pythonRuntime); os.IsNotExist(err) {
			fmt.Println("Error: Python interpreter not found")
			fmt.Println("Please install Python or set Python path in config")
			os.Exit(1)
		}
	}

	fmt.Printf("Running Python script: %s\n", scriptPath)
	fmt.Printf("Using Python: %s\n", pythonRuntime)

	// Prepare arguments
	args := []string{scriptPath}
	if len(os.Args) > 4 {
		args = append(args, os.Args[4:]...)
	}

	// Run Python script
	cmd := exec.Command(pythonRuntime, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running Python script: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Python script execution completed successfully!")
}

// installPython installs Python to the specified path
func installPython() {
	// Parse arguments
	version := "3.13.12"
	installPath := "D:\\Python"

	if len(os.Args) >= 4 {
		version = os.Args[3]
	}

	if len(os.Args) >= 5 {
		installPath = os.Args[4]
	}

	fmt.Printf("Installing Python %s to %s\n", version, installPath)

	// Check if Python version is already recorded in config
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Check if the requested Python version is already installed
	if existingPath, exists := config.PythonVersions[version]; exists {
		// Check if the existing path still exists
		if _, err := os.Stat(existingPath); err == nil {
			fmt.Printf("Python %s is already installed at: %s\n", version, existingPath)
			// Test the existing installation
			testCmd := exec.Command(existingPath, "--version")
			var testOutput strings.Builder
			testCmd.Stdout = &testOutput
			testCmd.Stderr = &testOutput
			if err := testCmd.Run(); err != nil {
				fmt.Printf("Error testing existing Python installation: %v\n", err)
				os.Exit(1)
			}
			existingVersion := strings.TrimSpace(testOutput.String())
			fmt.Printf("Existing Python version: %s\n", existingVersion)
			
			// Update config with Python path
			pythonLang := config.Languages["python"]
			pythonLang.Runtime = existingPath
			config.Languages["python"] = pythonLang
			if err := saveConfig(config); err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("Python path updated in config: %s\n", existingPath)
			os.Exit(0)
		} else {
			fmt.Printf("Python %s was recorded in config but the path no longer exists, proceeding with installation...\n", version)
			// Remove the invalid entry from config
			delete(config.PythonVersions, version)
			if err := saveConfig(config); err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				os.Exit(1)
			}
		}
	}

	// Check if Python is already installed at the specified path
	pythonExe := filepath.Join(installPath, "python.exe")
	if _, err := os.Stat(pythonExe); err == nil {
		fmt.Printf("Python is already installed at: %s\n", pythonExe)
		// Test the existing installation
		testCmd := exec.Command(pythonExe, "--version")
		var testOutput strings.Builder
		testCmd.Stdout = &testOutput
		testCmd.Stderr = &testOutput
		if err := testCmd.Run(); err != nil {
			fmt.Printf("Error testing existing Python installation: %v\n", err)
			os.Exit(1)
		}
		existingVersion := strings.TrimSpace(testOutput.String())
		fmt.Printf("Existing Python version: %s\n", existingVersion)
		
		// Check if the installed version matches the requested version
		if strings.Contains(existingVersion, version) {
			fmt.Println("Skipping installation as the requested Python version is already installed.")
			
			// Update config with Python path and version
			pythonLang := config.Languages["python"]
			pythonLang.Runtime = pythonExe
			config.Languages["python"] = pythonLang
			config.PythonVersions[version] = pythonExe
			if err := saveConfig(config); err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				os.Exit(1)
			}
			
			fmt.Printf("Python path updated in config: %s\n", pythonExe)
			fmt.Printf("Python version %s recorded in config\n", version)
			os.Exit(0)
		} else {
			fmt.Printf("Installed version (%s) does not match requested version (%s), proceeding with installation...\n", existingVersion, version)
		}
	}

	// Create installation directory if it doesn't exist
	fmt.Printf("Creating installation directory: %s\n", installPath)
	if err := os.MkdirAll(installPath, 0755); err != nil {
		fmt.Printf("Error creating installation directory: %v\n", err)
		os.Exit(1)
	}

	// Check if directory was created
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		fmt.Printf("Error: Installation directory was not created: %s\n", installPath)
		os.Exit(1)
	}
	fmt.Printf("Installation directory created successfully\n")

	// Download Python installer
	installerURL := fmt.Sprintf("https://www.python.org/ftp/python/%s/python-%s-amd64.exe", version, version)
	installerPath := filepath.Join(installPath, "python-installer.exe")

	fmt.Printf("Downloading Python installer from: %s\n", installerURL)
	fmt.Printf("Saving to: %s\n", installerPath)

	// Download the installer with progress
	data, err := downloadFileFromURL(installerURL)
	if err != nil {
		fmt.Printf("Error downloading Python installer: %v\n", err)
		os.Exit(1)
	}

	// Save the installer to disk
	fmt.Printf("Saving installer to disk...\n")
	if err := os.WriteFile(installerPath, data, 0644); err != nil {
		fmt.Printf("Error saving Python installer: %v\n", err)
		os.Exit(1)
	}

	// Check if installer was saved
	if _, err := os.Stat(installerPath); os.IsNotExist(err) {
		fmt.Printf("Error: Installer was not saved: %s\n", installerPath)
		os.Exit(1)
	}
	fmt.Printf("Installer saved successfully (%.2f MB)\n", float64(len(data))/1024/1024)

	fmt.Println("Download completed, starting installation...")
	fmt.Println("This may take a few minutes, please wait...")

	// Run the installer with progress
	fmt.Println("Running Python installer with progress...")
	fmt.Println("You may see a Python installation window appear...")

	// Run the installer with PowerShell to get admin privileges
	fmt.Println("Running installer with administrative privileges...")
	
	// Create a PowerShell command to run the installer
	psCommand := fmt.Sprintf(`
		Write-Host "Running Python installer as administrator..."
		$installerPath = "%s"
		$installPath = "%s"
		
		Write-Host "Installer path: $installerPath"
		Write-Host "Target directory: $installPath"
		
		# Create target directory if it doesn't exist
		if (-not (Test-Path $installPath)) {
			Write-Host "Creating target directory: $installPath"
			New-Item -ItemType Directory -Path $installPath -Force | Out-Null
		}
		
		# Verify installer exists
		if (-not (Test-Path $installerPath)) {
			Write-Host "Error: Installer not found at $installerPath"
			exit 1
		}
		
		# Run the installer with verbose output
		Write-Host "Starting Python installer..."
		Write-Host "You may need to click 'Yes' on the UAC prompt..."
		Write-Host "Installer arguments: /passive InstallAllUsers=1 DefaultAllUsersTargetDir=$installPath PrependPath=1 Include_test=0"
		
		# Use different installation parameters
		try {
			$process = Start-Process -FilePath $installerPath -ArgumentList "/passive", "InstallAllUsers=1", "DefaultAllUsersTargetDir=$installPath", "PrependPath=1", "Include_test=0" -Verb RunAs -PassThru
			
			# Wait for installation to complete
			Write-Host "Waiting for installation to complete..."
			$process.WaitForExit()
			
			Write-Host "Installation process exited with code: $($process.ExitCode)"
		} catch {
			Write-Host "Error running installer: $($_.Exception.Message)"
			exit 1
		}
		
		# Check if Python was installed
		$pythonExe = Join-Path $installPath "python.exe"
		if (Test-Path $pythonExe) {
			Write-Host "Python installed successfully!"
			Write-Host "Python executable: $pythonExe"
			Write-Host "Python version: $(& $pythonExe --version)"
		} else {
			Write-Host "Error: Python executable not found at $pythonExe"
			Write-Host "Directory contents:"
			Get-ChildItem $installPath -Force | ForEach-Object { Write-Host "  $($_.Name)" }
			Write-Host "Installer log:"
			Get-ChildItem $env:TEMP -Filter "python*.log" | Sort-Object LastWriteTime -Descending | Select-Object -First 1 | ForEach-Object { Get-Content $_.FullName | Select-Object -Last 50 }
			exit 1
		}
	`, installerPath, installPath)

	// Write the PowerShell script to a temporary file
	tempScriptPath := filepath.Join(os.TempDir(), "install-python-admin.ps1")
	if err := os.WriteFile(tempScriptPath, []byte(psCommand), 0644); err != nil {
		fmt.Printf("Error creating temporary script: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(tempScriptPath)

	// Run the PowerShell script
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", tempScriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running installation script: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Installation completed successfully!")

	// List directory contents before cleanup
	fmt.Println("Checking installation directory contents...")
	if files, err := os.ReadDir(installPath); err == nil {
		fmt.Println("Directory contents:")
		for _, file := range files {
			fmt.Printf("  %s\n", file.Name())
		}
	} else {
		fmt.Printf("Error reading directory: %v\n", err)
	}

	// Clean up installer
	fmt.Println("Cleaning up installer...")
	if err := os.Remove(installerPath); err != nil {
		fmt.Printf("Warning: Failed to clean up installer: %v\n", err)
	}

	// Verify installation
	pythonExe = filepath.Join(installPath, "python.exe")
	fmt.Printf("Verifying installation at: %s\n", pythonExe)
	if _, err := os.Stat(pythonExe); os.IsNotExist(err) {
		fmt.Println("Error: Python installation failed, executable not found")
		// List directory contents again
		fmt.Println("Final directory contents:")
		if files, err := os.ReadDir(installPath); err == nil {
			for _, file := range files {
				fmt.Printf("  %s\n", file.Name())
			}
		} else {
			fmt.Printf("Error reading directory: %v\n", err)
		}
		os.Exit(1)
	}

	// Test Python executable
	fmt.Println("Testing Python executable...")
	testCmd := exec.Command(pythonExe, "--version")
	var testOutput strings.Builder
	testCmd.Stdout = &testOutput
	testCmd.Stderr = &testOutput
	if err := testCmd.Run(); err != nil {
		fmt.Printf("Error testing Python executable: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Python version: %s\n", strings.TrimSpace(testOutput.String()))

	// Update config with Python path
	config, err = loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Fix: Get the struct, modify it, then put it back
	pythonLang := config.Languages["python"]
	pythonLang.Runtime = pythonExe
	config.Languages["python"] = pythonLang

	// Record Python version and installation path
	config.PythonVersions[version] = pythonExe
	if err := saveConfig(config); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Python %s installed successfully to %s\n", version, installPath)
	fmt.Printf("Python executable: %s\n", pythonExe)
	fmt.Printf("Python version %s recorded in config\n", version)
}

// downloadFileFromURL downloads a file from the specified URL with progress
func downloadFileFromURL(url string) ([]byte, error) {
	// Create channel to receive download result
	resultCh := make(chan struct {
		data []byte
		err  error
	})

	// Start download in a goroutine
	go func() {
		fmt.Printf("Starting download from: %s\n", url)
		
		resp, err := http.Get(url)
		if err != nil {
			resultCh <- struct {
				data []byte
				err  error
			}{nil, fmt.Errorf("failed to start download: %v", err)}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			resultCh <- struct {
				data []byte
				err  error
			}{nil, fmt.Errorf("HTTP error: %s", resp.Status)}
			return
		}

		contentLength := resp.ContentLength
		fmt.Printf("File size: %.2f MB\n", float64(contentLength)/1024/1024)
		
		// Check if content length is reasonable for Python installer
		if contentLength < 20*1024*1024 { // Less than 20MB
			resultCh <- struct {
				data []byte
				err  error
			}{nil, fmt.Errorf("downloaded file size seems too small: %.2f MB", float64(contentLength)/1024/1024)}
			return
		}
		
		// Warn if content length seems smaller than expected
		if contentLength < 50*1024*1024 { // Less than 50MB
			fmt.Printf("Warning: Downloaded file size (%.2f MB) seems smaller than expected for Python installer\n", float64(contentLength)/1024/1024)
			fmt.Println("Proceeding with installation...")
		}

		var totalDownloaded int64 = 0
		var data []byte
		buffer := make([]byte, 1024*1024) // 1MB buffer

		startTime := time.Now()

		for {
			n, err := resp.Body.Read(buffer)
			if err != nil && err != io.EOF {
				resultCh <- struct {
					data []byte
					err  error
				}{nil, fmt.Errorf("error during download: %v", err)}
				return
			}
			if n == 0 {
				break
			}

			data = append(data, buffer[:n]...)
			totalDownloaded += int64(n)

			// Show download progress
			if contentLength > 0 {
				percentage := float64(totalDownloaded) / float64(contentLength) * 100
				timeElapsed := time.Since(startTime).Seconds()
				if timeElapsed > 0 {
					speed := float64(totalDownloaded) / timeElapsed / 1024 / 1024 // MB/s
					fmt.Printf("Downloading: %.2f%% (%.2f MB / %.2f MB) | Speed: %.2f MB/s\r", percentage, float64(totalDownloaded)/1024/1024, float64(contentLength)/1024/1024, speed)
				} else {
					fmt.Printf("Downloading: %.2f%% (%.2f MB / %.2f MB)\r", percentage, float64(totalDownloaded)/1024/1024, float64(contentLength)/1024/1024)
				}
			} else {
				timeElapsed := time.Since(startTime).Seconds()
				if timeElapsed > 0 {
					speed := float64(totalDownloaded) / timeElapsed / 1024 / 1024 // MB/s
					fmt.Printf("Downloaded: %.2f MB | Speed: %.2f MB/s\r", float64(totalDownloaded)/1024/1024, speed)
				} else {
					fmt.Printf("Downloaded: %.2f MB\r", float64(totalDownloaded)/1024/1024)
				}
			}
		}

		fmt.Println() // New line after progress
		
		// Verify download completion
		if contentLength > 0 && totalDownloaded != contentLength {
			resultCh <- struct {
				data []byte
				err  error
			}{nil, fmt.Errorf("download incomplete: expected %.2f MB, got %.2f MB", float64(contentLength)/1024/1024, float64(totalDownloaded)/1024/1024)}
			return
		}
		
		fmt.Printf("Download completed successfully! Total size: %.2f MB\n", float64(totalDownloaded)/1024/1024)
		resultCh <- struct {
			data []byte
			err  error
		}{data, nil}
	}()

	// Wait for download to complete
	result := <-resultCh
	return result.data, result.err
}

// printVersion prints the version information
func printVersion() {
	fmt.Printf("Enlightenment Lighthouse Runtime v%s\n", elr.Version)
	fmt.Printf("Platform: %s\n", runtime.GOOS)
}

// printHelp prints the help information
func printHelp() {
	fmt.Println("Enlightenment Lighthouse Runtime (ELR)")
	fmt.Println("Usage: elr [command] [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  version           Print version information")
	fmt.Println("  help              Print this help message")
	fmt.Println("  start             Start the ELR runtime")
	fmt.Println("  stop              Stop the ELR runtime")
	fmt.Println("  create            Create a new container")
	fmt.Println("  run               Create and start a new container")
	fmt.Println("  start-container   Start a container")
	fmt.Println("  stop-container    Stop a container")
	fmt.Println("  list              List all containers")
	fmt.Println("  delete            Delete a container")
	fmt.Println("  inspect           Inspect a container")
	// 系统设置命令
	fmt.Println("  setup             Setup ELR system (e.g., isolation)")
	// 模型管理命令
	fmt.Println("  model list        List all models")
	fmt.Println("  model get         Get model information")
	fmt.Println("  model download    Download a model")
	fmt.Println("  model delete      Delete a model")
	fmt.Println("  model install-deps Install model dependencies")
	// 沙箱管理命令
	fmt.Println("  sandbox list      List all sandboxes")
	fmt.Println("  sandbox create    Create a new sandbox")
	fmt.Println("  sandbox start     Start a sandbox")
	fmt.Println("  sandbox stop      Stop a sandbox")
	fmt.Println("  sandbox delete    Delete a sandbox")
	fmt.Println("  sandbox load-model Load model into sandbox")
	fmt.Println("  sandbox unload-model Unload model from sandbox")
	// API 服务命令
	fmt.Println("  api start         Start API services (all or specific)")
	fmt.Println("  api stop          Stop API services (all or specific)")
	fmt.Println("  api status        Check API service status")
	fmt.Println("  api config        Configure API addresses and ports")
	// 文件系统管理命令
	fmt.Println("  fs upload         Upload file to container")
	fmt.Println("  fs download       Download file from container")
	fmt.Println("  fs set-dir        Set directory for file type")
	fmt.Println("  fs get-dir        Get directory for file type")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --name            Container name")
	fmt.Println("  --image           Container image")
	fmt.Println("  --command         Command to run")
	fmt.Println("  --arg             Command argument")
	fmt.Println("  --env             Environment variable")
	fmt.Println("  --id              Container ID")
	fmt.Println("  --model-id        Model ID")
	fmt.Println("  --sandbox-id      Sandbox ID")
	fmt.Println("  --container       Container name")
	fmt.Println("  --type            Model type")
	fmt.Println("  --url             Download URL")
	// 系统设置选项
	fmt.Println("  --isolation       Isolation type (windows-container, wsl, basic)")
	// API 服务选项
	fmt.Println("  --api-type        API type (desktop, public, model)")
	fmt.Println("  --address         API address")
	fmt.Println("  --port            API port")
	// 文件系统选项
	fmt.Println("  --file-type       File type (e.g., model, data, config)")
	fmt.Println("  --directory       Directory path")
	fmt.Println("  --local-path      Local file path")
	fmt.Println("  --container-path  Container file path")
	fmt.Println("  --token           Authentication token")
} 

// loadConfig loads the configuration from file
func loadConfig() (*elr.Config, error) {
	configPath := os.Getenv("ELR_CONFIG")
	if configPath == "" {
		configPath = "~/.elr/config.yaml"
	}

	// Expand ~ to home directory
	if len(configPath) > 0 && configPath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %v", err)
		}
		configPath = homeDir + configPath[1:]
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config
		return defaultConfig(), nil
	}

	// Read config file
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse config
	config := &elr.Config{}
	if err := yaml.Unmarshal(configBytes, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Initialize FileDirectories if it's nil
	if config.FileDirectories == nil {
		config.FileDirectories = make(map[string]string)
	}

	// Initialize PythonVersions if it's nil
	if config.PythonVersions == nil {
		config.PythonVersions = make(map[string]string)
	}

	return config, nil
}

// defaultConfig returns the default configuration
func defaultConfig() *elr.Config {
	return &elr.Config{
		LogLevel:  "info",
		DataDir:   "~/.elr/data",
		PluginDir: "~/.elr/plugins",
		FileDirectories: make(map[string]string),
		PythonVersions: make(map[string]string),
		Platform: struct {
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
		}{
			Linux: struct {
				UseNamespaces bool `yaml:"use_namespaces"`
				UseCgroups    bool `yaml:"use_cgroups"`
			}{
				UseNamespaces: true,
				UseCgroups:    true,
			},
			Windows: struct {
				UseJobObjects bool   `yaml:"use_job_objects"`
				UseWSL        bool   `yaml:"use_wsl"`
				UseContainers bool   `yaml:"use_containers"`
				IsolationType string `yaml:"isolation_type"` // Options: "windows-container", "wsl", "basic"
			}{
				UseJobObjects: true,
				UseWSL:        false,
				UseContainers: false,
				IsolationType: "basic",
			},
			Darwin: struct {
				UseSandbox bool `yaml:"use_sandbox"`
				UseSpctl   bool `yaml:"use_spctl"`
			}{
				UseSandbox: true,
				UseSpctl:   true,
			},
		},
		Network: struct {
			Enable  bool   `yaml:"enable"`
			Bridge  string `yaml:"bridge"`
			Subnet  string `yaml:"subnet"`
			APIPorts struct {
				DesktopAPI int `yaml:"desktop_api"`
				PublicAPI  int `yaml:"public_api"`
				ModelAPI   int `yaml:"model_api"`
			} `yaml:"api_ports"`
		}{
			Enable:  false, // 默认禁用网络
			Bridge:  "elr0",
			Subnet:  "172.16.0.0/16",
			APIPorts: struct {
				DesktopAPI int `yaml:"desktop_api"`
				PublicAPI  int `yaml:"public_api"`
				ModelAPI   int `yaml:"model_api"`
			}{
				DesktopAPI: 8081,
				PublicAPI:  8080,
				ModelAPI:   8082,
			},
		},
		Storage: struct {
			Enable  bool   `yaml:"enable"`
			Driver  string `yaml:"driver"`
			BaseDir string `yaml:"base_dir"`
		}{
			Enable:  true,
			Driver:  "overlay",
			BaseDir: "~/.elr/storage",
		},
		Languages: map[string]struct {
			Enable  bool   `yaml:"enable"`
			Runtime string `yaml:"runtime"`
		}{
			"cpp": {
				Enable:  true,
				Runtime: "/usr/bin/gcc",
			},
			"python": {
				Enable:  true,
				Runtime: "/usr/bin/python3",
			},
			"nodejs": {
				Enable:  true,
				Runtime: "/usr/bin/node",
			},
			"java": {
				Enable:  true,
				Runtime: "/usr/bin/java",
			},
			"go": {
				Enable:  true,
				Runtime: "/usr/bin/go",
			},
		},
		Resources: struct {
			Types map[string]struct {
				Enable bool   `yaml:"enable"`
				Dir    string `yaml:"dir"`
			} `yaml:"types"`
			ModelTypes map[string]struct {
				Enable bool   `yaml:"enable"`
				Dir    string `yaml:"dir"`
			} `yaml:"model_types"`
		}{
			Types: map[string]struct {
				Enable bool   `yaml:"enable"`
				Dir    string `yaml:"dir"`
			}{
				"component": {
					Enable: true,
					Dir:    "~/.elr/resources/components",
				},
				"model": {
					Enable: true,
					Dir:    "~/.elr/resources/models",
				},
				"project": {
					Enable: true,
					Dir:    "~/.elr/resources/projects",
				},
			},
			ModelTypes: map[string]struct {
				Enable bool   `yaml:"enable"`
				Dir    string `yaml:"dir"`
			}{
				"text": {
					Enable: true,
					Dir:    "~/.elr/resources/models/text",
				},
				"image": {
					Enable: true,
					Dir:    "~/.elr/resources/models/image",
				},
				"audio": {
					Enable: true,
					Dir:    "~/.elr/resources/models/audio",
				},
				"video": {
					Enable: true,
					Dir:    "~/.elr/resources/models/video",
				},
			},
		},
	}
}

// startRuntime starts the ELR runtime
func startRuntime() {
	// Load config
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Create runtime
	runtime, err := elr.NewRuntime(config)
	if err != nil {
		fmt.Printf("Error creating runtime: %v\n", err)
		os.Exit(1)
	}

	// Start runtime
	if err := runtime.Start(); err != nil {
		fmt.Printf("Error starting runtime: %v\n", err)
		os.Exit(1)
	}

	// Wait for signal to stop
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh

	// Stop runtime
	if err := runtime.Stop(); err != nil {
		fmt.Printf("Error stopping runtime: %v\n", err)
		os.Exit(1)
	}
}

// Global runtime instance
var globalRuntime *elr.Runtime
var runtimeOnce sync.Once

// getRuntime returns a runtime instance
func getRuntime() (*elr.Runtime, error) {
	var err error
	runtimeOnce.Do(func() {
		config, loadErr := loadConfig()
		if loadErr != nil {
			err = loadErr
			return
		}

		globalRuntime, loadErr = elr.NewRuntime(config)
		if loadErr != nil {
			err = loadErr
			return
		}

		// Start runtime if not already running
		if loadErr := globalRuntime.Start(); loadErr != nil {
			err = loadErr
			return
		}
	})

	if err != nil {
		return nil, err
	}

	return globalRuntime, nil
}

// stopRuntime stops the ELR runtime
func stopRuntime() {
	fmt.Println("Stopping ELR runtime...")
	
	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	if err := runtime.Stop(); err != nil {
		fmt.Printf("Error stopping runtime: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("ELR runtime stopped successfully!")
}

// createContainer creates a new container
func createContainer() {
	fmt.Println("Creating container...")
	
	// Parse arguments
	name := ""
	image := ""
	fileSystemIsolation := false
	rootFSPath := ""
	readOnlyFS := false
	
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--name" && i+1 < len(os.Args) {
			name = os.Args[i+1]
		} else if os.Args[i] == "--image" && i+1 < len(os.Args) {
			image = os.Args[i+1]
		} else if os.Args[i] == "--fs-isolation" {
			fileSystemIsolation = true
		} else if os.Args[i] == "--rootfs" && i+1 < len(os.Args) {
			rootFSPath = os.Args[i+1]
		} else if os.Args[i] == "--read-only" {
			readOnlyFS = true
		}
	}

	if name == "" {
		fmt.Println("Error: Container name is required")
		os.Exit(1)
	}

	if image == "" {
		fmt.Println("Error: Container image is required")
		os.Exit(1)
	}

	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	config := elr.ContainerConfig{
		Name:               name,
		Image:              image,
		FileSystemIsolation: fileSystemIsolation,
		RootFSPath:         rootFSPath,
		ReadOnlyFS:         readOnlyFS,
	}

	container, err := runtime.CreateContainer(config)
	if err != nil {
		fmt.Printf("Error creating container: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Container created successfully! ID: %s, Name: %s, Image: %s\n", container.ID, container.Name, container.Image)
	if fileSystemIsolation {
		fmt.Printf("File system isolation: enabled\n")
		if rootFSPath != "" {
			fmt.Printf("Root FS path: %s\n", rootFSPath)
		}
		if readOnlyFS {
			fmt.Printf("Read-only filesystem: enabled\n")
		}
	}
}

// runContainer creates and starts a new container
func runContainer() {
	fmt.Println("Running container...")
	
	// Parse arguments
	name := ""
	image := ""
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--name" && i+1 < len(os.Args) {
			name = os.Args[i+1]
		} else if os.Args[i] == "--image" && i+1 < len(os.Args) {
			image = os.Args[i+1]
		}
	}

	if name == "" {
		fmt.Println("Error: Container name is required")
		os.Exit(1)
	}

	if image == "" {
		fmt.Println("Error: Container image is required")
		os.Exit(1)
	}

	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	config := elr.ContainerConfig{
		Name:  name,
		Image: image,
	}

	container, err := runtime.CreateContainer(config)
	if err != nil {
		fmt.Printf("Error creating container: %v\n", err)
		os.Exit(1)
	}

	if err := container.Start(); err != nil {
		fmt.Printf("Error starting container: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Container started successfully! ID: %s, Name: %s, Image: %s\n", container.ID, container.Name, container.Image)
}

// startContainer starts a container
func startContainer() {
	fmt.Println("Starting container...")
	
	// Parse arguments
	id := ""
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--id" && i+1 < len(os.Args) {
			id = os.Args[i+1]
			break
		}
	}

	if id == "" {
		fmt.Println("Error: Container ID is required")
		os.Exit(1)
	}

	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	container, err := runtime.GetContainer(id)
	if err != nil {
		fmt.Printf("Error getting container: %v\n", err)
		os.Exit(1)
	}

	if err := container.Start(); err != nil {
		fmt.Printf("Error starting container: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Container started successfully! ID: %s, Name: %s\n", container.ID, container.Name)
}

// stopContainer stops a container
func stopContainer() {
	fmt.Println("Stopping container...")
	
	// Parse arguments
	id := ""
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--id" && i+1 < len(os.Args) {
			id = os.Args[i+1]
			break
		}
	}

	if id == "" {
		fmt.Println("Error: Container ID is required")
		os.Exit(1)
	}

	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	container, err := runtime.GetContainer(id)
	if err != nil {
		fmt.Printf("Error getting container: %v\n", err)
		os.Exit(1)
	}

	if err := container.Stop(); err != nil {
		fmt.Printf("Error stopping container: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Container stopped successfully! ID: %s, Name: %s\n", container.ID, container.Name)
}

// listContainers lists all containers
func listContainers() {
	fmt.Println("Listing containers...")
	
	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	containers := runtime.ListContainers()

	if len(containers) == 0 {
		fmt.Println("No containers found")
		return
	}

	fmt.Printf("%-20s %-15s %-15s %-10s %-20s\n", "ID", "Name", "Image", "Status", "Created")
	fmt.Printf("%-20s %-15s %-15s %-10s %-20s\n", "--", "----", "-----", "------", "-------")

	for _, container := range containers {
		created := container.Created.Format("2006-01-02 15:04:05")
		fmt.Printf("%-20s %-15s %-15s %-10s %-20s\n", container.ID, container.Name, container.Image, container.Status, created)
	}

	fmt.Printf("\nTotal containers: %d\n", len(containers))
}

// deleteContainer deletes a container
func deleteContainer() {
	fmt.Println("Deleting container...")
	
	// Parse arguments
	id := ""
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--id" && i+1 < len(os.Args) {
			id = os.Args[i+1]
			break
		}
	}

	if id == "" {
		fmt.Println("Error: Container ID is required")
		os.Exit(1)
	}

	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	if err := runtime.DeleteContainer(id); err != nil {
		fmt.Printf("Error deleting container: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Container deleted successfully! ID: %s\n", id)
}

// inspectContainer inspects a container
func inspectContainer() {
	fmt.Println("Inspecting container...")
	
	// Parse arguments
	id := ""
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--id" && i+1 < len(os.Args) {
			id = os.Args[i+1]
			break
		}
	}

	if id == "" {
		fmt.Println("Error: Container ID is required")
		os.Exit(1)
	}

	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	container, err := runtime.GetContainer(id)
	if err != nil {
		fmt.Printf("Error getting container: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Container Details:\n")
	fmt.Printf("==================\n")
	fmt.Printf("ID: %s\n", container.ID)
	fmt.Printf("Name: %s\n", container.Name)
	fmt.Printf("Image: %s\n", container.Image)
	fmt.Printf("Status: %s\n", container.Status)
	fmt.Printf("Created: %s\n", container.Created.Format("2006-01-02 15:04:05"))
	if container.Started != nil {
		fmt.Printf("Started: %s\n", container.Started.Format("2006-01-02 15:04:05"))
	}
	if container.Stopped != nil {
		fmt.Printf("Stopped: %s\n", container.Stopped.Format("2006-01-02 15:04:05"))
	}
	if container.PID > 0 {
		fmt.Printf("PID: %d\n", container.PID)
	}
	if container.ExitCode != 0 {
		fmt.Printf("Exit Code: %d\n", container.ExitCode)
	}
	if container.Error != "" {
		fmt.Printf("Error: %s\n", container.Error)
	}
	fmt.Printf("Directory: %s\n", container.Dir)
}

// 模型管理命令处理函数
func modelCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Model subcommand is required")
		printHelp()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
		case "list":
			listModels()
		case "get":
			getModel()
		case "download":
			downloadModel()
		case "delete":
			deleteModel()
		case "install-deps":
			installModelDependencies()
		default:
			fmt.Printf("Unknown model subcommand: %s\n", subcommand)
			printHelp()
			os.Exit(1)
		}
}

// listModels 列出所有模型
func listModels() {
	fmt.Println("Listing models...")

	// 加载模型配置
	modelConfig, err := loadModelConfig()
	if err != nil {
		fmt.Printf("Error loading model config: %v\n", err)
		os.Exit(1)
	}

	// 创建模型管理器
	modelManager, err := model.NewModelManager(modelConfig)
	if err != nil {
		fmt.Printf("Error creating model manager: %v\n", err)
		os.Exit(1)
	}

	// 列出模型
	models, err := modelManager.ListModels()
	if err != nil {
		fmt.Printf("Error listing models: %v\n", err)
		os.Exit(1)
	}

	if len(models) == 0 {
		fmt.Println("No models found")
		return
	}

	fmt.Printf("%-20s %-15s %-15s %-10s %-20s\n", "ID", "Name", "Type", "Version", "Path")
	fmt.Printf("%-20s %-15s %-15s %-10s %-20s\n", "--", "----", "----", "-------", "----")

	for _, m := range models {
		fmt.Printf("%-20s %-15s %-15s %-10s %-20s\n", m.ID, m.Name, m.Type, m.Version, m.Path)
	}

	fmt.Printf("\nTotal models: %d\n", len(models))
}

// getModel 获取模型信息
func getModel() {
	fmt.Println("Getting model information...")

	// 解析参数
	modelID := ""
	for i := 3; i < len(os.Args); i++ {
		if os.Args[i] == "--model-id" && i+1 < len(os.Args) {
			modelID = os.Args[i+1]
			break
		}
	}

	if modelID == "" {
		fmt.Println("Error: Model ID is required")
		os.Exit(1)
	}

	// 加载模型配置
	modelConfig, err := loadModelConfig()
	if err != nil {
		fmt.Printf("Error loading model config: %v\n", err)
		os.Exit(1)
	}

	// 创建模型管理器
	modelManager, err := model.NewModelManager(modelConfig)
	if err != nil {
		fmt.Printf("Error creating model manager: %v\n", err)
		os.Exit(1)
	}

	// 获取模型信息
	m, err := modelManager.GetModel(modelID)
	if err != nil {
		fmt.Printf("Error getting model: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Model ID: %s\n", m.ID)
	fmt.Printf("Name: %s\n", m.Name)
	fmt.Printf("Type: %s\n", m.Type)
	fmt.Printf("Version: %s\n", m.Version)
	fmt.Printf("Path: %s\n", m.Path)
	if m.Properties != nil {
		fmt.Printf("Description: %s\n", m.Properties.Description)
	}
}

// downloadModel 下载模型
func downloadModel() {
	fmt.Println("Downloading model...")

	// 解析参数
	modelID := ""
	modelType := ""
	url := ""
	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--model-id":
			if i+1 < len(os.Args) {
				modelID = os.Args[i+1]
			}
		case "--type":
			if i+1 < len(os.Args) {
				modelType = os.Args[i+1]
			}
		case "--url":
			if i+1 < len(os.Args) {
				url = os.Args[i+1]
			}
		}
	}

	if modelID == "" || url == "" {
		fmt.Println("Error: Model ID and URL are required")
		os.Exit(1)
	}

	// 加载模型配置
	modelConfig, err := loadModelConfig()
	if err != nil {
		fmt.Printf("Error loading model config: %v\n", err)
		os.Exit(1)
	}

	// 创建模型管理器
	modelManager, err := model.NewModelManager(modelConfig)
	if err != nil {
		fmt.Printf("Error creating model manager: %v\n", err)
		os.Exit(1)
	}

	// 下载模型
	if err := modelManager.DownloadModel(modelID, modelType, url); err != nil {
		fmt.Printf("Error downloading model: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Model downloaded successfully! ID: %s\n", modelID)
}

// deleteModel 删除模型
func deleteModel() {
	fmt.Println("Deleting model...")

	// 解析参数
	modelID := ""
	for i := 3; i < len(os.Args); i++ {
		if os.Args[i] == "--model-id" && i+1 < len(os.Args) {
			modelID = os.Args[i+1]
			break
		}
	}

	if modelID == "" {
		fmt.Println("Error: Model ID is required")
		os.Exit(1)
	}

	// 加载模型配置
	modelConfig, err := loadModelConfig()
	if err != nil {
		fmt.Printf("Error loading model config: %v\n", err)
		os.Exit(1)
	}

	// 创建模型管理器
	modelManager, err := model.NewModelManager(modelConfig)
	if err != nil {
		fmt.Printf("Error creating model manager: %v\n", err)
		os.Exit(1)
	}

	// 删除模型
	if err := modelManager.DeleteModel(modelID); err != nil {
		fmt.Printf("Error deleting model: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Model deleted successfully! ID: %s\n", modelID)
}

// installModelDependencies 安装模型依赖
func installModelDependencies() {
	fmt.Println("Installing model dependencies...")

	// 解析参数
	modelID := ""
	depType := ""
	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--model-id":
			if i+1 < len(os.Args) {
				modelID = os.Args[i+1]
			}
		case "--type":
			if i+1 < len(os.Args) {
				depType = os.Args[i+1]
			}
		}
	}

	if modelID == "" {
		fmt.Println("Error: Model ID is required")
		os.Exit(1)
	}

	if depType == "" {
		fmt.Println("Error: Dependency type is required")
		os.Exit(1)
	}

	// 加载模型配置
	modelConfig, err := loadModelConfig()
	if err != nil {
		fmt.Printf("Error loading model config: %v\n", err)
		os.Exit(1)
	}

	// 创建模型管理器
	modelManager, err := model.NewModelManager(modelConfig)
	if err != nil {
		fmt.Printf("Error creating model manager: %v\n", err)
		os.Exit(1)
	}

	// 安装模型依赖
	if err := modelManager.InstallModelDependencies(modelID, depType); err != nil {
		fmt.Printf("Error installing model dependencies: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Model dependencies installed successfully! Model ID: %s, Type: %s\n", modelID, depType)
}

// loadModelConfig 加载模型配置
func loadModelConfig() (*config.Config, error) {
	// 加载 ELR 配置
	elrConfig, err := loadConfig()
	if err != nil {
		// 如果加载失败，使用默认配置
		return &config.Config{
			Model: config.ModelConfig{
				ModelDir: "../micro_model/model/models",
			},
			Sandbox: config.SandboxConfig{},
		}, nil
	}
	
	// 从 ELR 配置中获取模型目录
	modelDir := "../micro_model/model/models"
	if chatModelConfig, exists := elrConfig.Resources.ModelTypes["chat"]; exists {
		modelDir = chatModelConfig.Dir
	}
	
	// 创建配置
	return &config.Config{
		Model: config.ModelConfig{
			ModelDir: modelDir,
		},
		Sandbox: config.SandboxConfig{},
	}, nil
}

// 沙箱管理命令处理函数
func sandboxCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Sandbox subcommand is required")
		printHelp()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "list":
		listSandboxes()
	case "create":
		createSandbox()
	case "start":
		startSandbox()
	case "stop":
		stopSandbox()
	case "delete":
		deleteSandbox()
	case "load-model":
		loadModelToSandbox()
	case "unload-model":
		unloadModelFromSandbox()
	case "run-model":
		runModelInSandbox()
	default:
		fmt.Printf("Unknown sandbox subcommand: %s\n", subcommand)
		printHelp()
		os.Exit(1)
	}
}

// listSandboxes 列出所有沙箱
func listSandboxes() {
	fmt.Println("Listing sandboxes...")

	// 加载配置
	modelConfig, err := loadModelConfig()
	if err != nil {
		fmt.Printf("Error loading model config: %v\n", err)
		os.Exit(1)
	}

	// 创建模型管理器
	modelManager, err := model.NewModelManager(modelConfig)
	if err != nil {
		fmt.Printf("Error creating model manager: %v\n", err)
		os.Exit(1)
	}

	// 创建沙箱管理器
	sandboxManager, err := sandbox.NewSandboxManager(modelConfig, modelManager)
	if err != nil {
		fmt.Printf("Error creating sandbox manager: %v\n", err)
		os.Exit(1)
	}

	// 列出沙箱
	sandboxes := sandboxManager.ListSandboxes()

	if len(sandboxes) == 0 {
		fmt.Println("No sandboxes found")
		return
	}

	fmt.Printf("%-20s %-15s %-15s %-10s %-20s\n", "ID", "Name", "Container", "Status", "Created")
	fmt.Printf("%-20s %-15s %-15s %-10s %-20s\n", "--", "----", "---------", "------", "-------")

	for _, s := range sandboxes {
		created := s.CreatedAt.Format("2006-01-02 15:04:05")
		fmt.Printf("%-20s %-15s %-15s %-10s %-20s\n", s.ID, s.ID, s.Container, s.Status, created)
	}

	fmt.Printf("\nTotal sandboxes: %d\n", len(sandboxes))
}

// createSandbox 创建新沙箱
func createSandbox() {
	fmt.Println("Creating sandbox...")

	// 解析参数
	container := ""
	for i := 3; i < len(os.Args); i++ {
		if os.Args[i] == "--container" && i+1 < len(os.Args) {
			container = os.Args[i+1]
			break
		}
	}

	if container == "" {
		fmt.Println("Error: Container name is required")
		os.Exit(1)
	}

	// 加载配置
	modelConfig, err := loadModelConfig()
	if err != nil {
		fmt.Printf("Error loading model config: %v\n", err)
		os.Exit(1)
	}

	// 创建模型管理器
	modelManager, err := model.NewModelManager(modelConfig)
	if err != nil {
		fmt.Printf("Error creating model manager: %v\n", err)
		os.Exit(1)
	}

	// 创建沙箱管理器
	sandboxManager, err := sandbox.NewSandboxManager(modelConfig, modelManager)
	if err != nil {
		fmt.Printf("Error creating sandbox manager: %v\n", err)
		os.Exit(1)
	}

	// 创建沙箱
	s, err := sandboxManager.CreateSandbox(container)
	if err != nil {
		fmt.Printf("Error creating sandbox: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sandbox created successfully! ID: %s, Container: %s\n", s.ID, s.Container)
}

// startSandbox 启动沙箱
func startSandbox() {
	fmt.Println("Starting sandbox...")

	// 解析参数
	sandboxID := ""
	for i := 3; i < len(os.Args); i++ {
		if os.Args[i] == "--sandbox-id" && i+1 < len(os.Args) {
			sandboxID = os.Args[i+1]
			break
		}
	}

	if sandboxID == "" {
		fmt.Println("Error: Sandbox ID is required")
		os.Exit(1)
	}

	// 加载配置
	modelConfig, err := loadModelConfig()
	if err != nil {
		fmt.Printf("Error loading model config: %v\n", err)
		os.Exit(1)
	}

	// 创建模型管理器
	modelManager, err := model.NewModelManager(modelConfig)
	if err != nil {
		fmt.Printf("Error creating model manager: %v\n", err)
		os.Exit(1)
	}

	// 创建沙箱管理器
	sandboxManager, err := sandbox.NewSandboxManager(modelConfig, modelManager)
	if err != nil {
		fmt.Printf("Error creating sandbox manager: %v\n", err)
		os.Exit(1)
	}

	// 启动沙箱
	if err := sandboxManager.StartSandbox(sandboxID); err != nil {
		fmt.Printf("Error starting sandbox: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sandbox started successfully! ID: %s\n", sandboxID)
}

// stopSandbox 停止沙箱
func stopSandbox() {
	fmt.Println("Stopping sandbox...")

	// 解析参数
	sandboxID := ""
	for i := 3; i < len(os.Args); i++ {
		if os.Args[i] == "--sandbox-id" && i+1 < len(os.Args) {
			sandboxID = os.Args[i+1]
			break
		}
	}

	if sandboxID == "" {
		fmt.Println("Error: Sandbox ID is required")
		os.Exit(1)
	}

	// 加载配置
	modelConfig, err := loadModelConfig()
	if err != nil {
		fmt.Printf("Error loading model config: %v\n", err)
		os.Exit(1)
	}

	// 创建模型管理器
	modelManager, err := model.NewModelManager(modelConfig)
	if err != nil {
		fmt.Printf("Error creating model manager: %v\n", err)
		os.Exit(1)
	}

	// 创建沙箱管理器
	sandboxManager, err := sandbox.NewSandboxManager(modelConfig, modelManager)
	if err != nil {
		fmt.Printf("Error creating sandbox manager: %v\n", err)
		os.Exit(1)
	}

	// 停止沙箱
	if err := sandboxManager.StopSandbox(sandboxID); err != nil {
		fmt.Printf("Error stopping sandbox: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sandbox stopped successfully! ID: %s\n", sandboxID)
}

// deleteSandbox 删除沙箱
func deleteSandbox() {
	fmt.Println("Deleting sandbox...")

	// 解析参数
	sandboxID := ""
	for i := 3; i < len(os.Args); i++ {
		if os.Args[i] == "--sandbox-id" && i+1 < len(os.Args) {
			sandboxID = os.Args[i+1]
			break
		}
	}

	if sandboxID == "" {
		fmt.Println("Error: Sandbox ID is required")
		os.Exit(1)
	}

	// 加载配置
	modelConfig, err := loadModelConfig()
	if err != nil {
		fmt.Printf("Error loading model config: %v\n", err)
		os.Exit(1)
	}

	// 创建模型管理器
	modelManager, err := model.NewModelManager(modelConfig)
	if err != nil {
		fmt.Printf("Error creating model manager: %v\n", err)
		os.Exit(1)
	}

	// 创建沙箱管理器
	sandboxManager, err := sandbox.NewSandboxManager(modelConfig, modelManager)
	if err != nil {
		fmt.Printf("Error creating sandbox manager: %v\n", err)
		os.Exit(1)
	}

	// 删除沙箱
	if err := sandboxManager.DeleteSandbox(sandboxID); err != nil {
		fmt.Printf("Error deleting sandbox: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sandbox deleted successfully! ID: %s\n", sandboxID)
}

// loadModelToSandbox 加载模型到沙箱
func loadModelToSandbox() {
	fmt.Println("Loading model into sandbox...")

	// 解析参数
	sandboxID := ""
	modelID := ""
	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--sandbox-id":
			if i+1 < len(os.Args) {
				sandboxID = os.Args[i+1]
			}
		case "--model-id":
			if i+1 < len(os.Args) {
				modelID = os.Args[i+1]
			}
		}
	}

	if sandboxID == "" || modelID == "" {
		fmt.Println("Error: Sandbox ID and Model ID are required")
		os.Exit(1)
	}

	// 加载配置
	modelConfig, err := loadModelConfig()
	if err != nil {
		fmt.Printf("Error loading model config: %v\n", err)
		os.Exit(1)
	}

	// 创建模型管理器
	modelManager, err := model.NewModelManager(modelConfig)
	if err != nil {
		fmt.Printf("Error creating model manager: %v\n", err)
		os.Exit(1)
	}

	// 创建沙箱管理器
	sandboxManager, err := sandbox.NewSandboxManager(modelConfig, modelManager)
	if err != nil {
		fmt.Printf("Error creating sandbox manager: %v\n", err)
		os.Exit(1)
	}

	// 加载模型到沙箱
	if err := sandboxManager.LoadModel(sandboxID, modelID); err != nil {
		fmt.Printf("Error loading model into sandbox: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Model loaded successfully! Model ID: %s, Sandbox ID: %s\n", modelID, sandboxID)
}

// unloadModelFromSandbox 从沙箱卸载模型
func unloadModelFromSandbox() {
	fmt.Println("Unloading model from sandbox...")

	// 解析参数
	sandboxID := ""
	modelID := ""
	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--sandbox-id":
			if i+1 < len(os.Args) {
				sandboxID = os.Args[i+1]
			}
		case "--model-id":
			if i+1 < len(os.Args) {
				modelID = os.Args[i+1]
			}
		}
	}

	if sandboxID == "" || modelID == "" {
		fmt.Println("Error: Sandbox ID and Model ID are required")
		os.Exit(1)
	}

	// 加载配置
	modelConfig, err := loadModelConfig()
	if err != nil {
		fmt.Printf("Error loading model config: %v\n", err)
		os.Exit(1)
	}

	// 创建模型管理器
	modelManager, err := model.NewModelManager(modelConfig)
	if err != nil {
		fmt.Printf("Error creating model manager: %v\n", err)
		os.Exit(1)
	}

	// 创建沙箱管理器
	sandboxManager, err := sandbox.NewSandboxManager(modelConfig, modelManager)
	if err != nil {
		fmt.Printf("Error creating sandbox manager: %v\n", err)
		os.Exit(1)
	}

	// 从沙箱卸载模型
	if err := sandboxManager.UnloadModel(sandboxID, modelID); err != nil {
		fmt.Printf("Error unloading model from sandbox: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Model unloaded successfully! Model ID: %s, Sandbox ID: %s\n", modelID, sandboxID)
}

// runModelInSandbox 在沙箱中运行模型
func runModelInSandbox() {
	fmt.Println("Running model in sandbox...")

	// 解析参数
	sandboxID := ""
	modelID := ""
	input := ""
	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--sandbox-id":
			if i+1 < len(os.Args) {
				sandboxID = os.Args[i+1]
			}
		case "--model-id":
			if i+1 < len(os.Args) {
				modelID = os.Args[i+1]
			}
		case "--input":
			if i+1 < len(os.Args) {
				input = os.Args[i+1]
			}
		}
	}

	if sandboxID == "" || modelID == "" || input == "" {
		fmt.Println("Error: Sandbox ID, Model ID, and Input are required")
		os.Exit(1)
	}

	// 加载配置
	modelConfig, err := loadModelConfig()
	if err != nil {
		fmt.Printf("Error loading model config: %v\n", err)
		os.Exit(1)
	}

	// 创建模型管理器
	modelManager, err := model.NewModelManager(modelConfig)
	if err != nil {
		fmt.Printf("Error creating model manager: %v\n", err)
		os.Exit(1)
	}

	// 创建沙箱管理器
	sandboxManager, err := sandbox.NewSandboxManager(modelConfig, modelManager)
	if err != nil {
		fmt.Printf("Error creating sandbox manager: %v\n", err)
		os.Exit(1)
	}

	// 运行模型
	output, err := sandboxManager.RunModel(sandboxID, modelID, input)
	if err != nil {
		fmt.Printf("Error running model: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Model run successfully!\n")
	fmt.Printf("Input: %s\n", input)
	fmt.Printf("Output: %s\n", output)
}

// API 服务命令处理函数
func apiCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Error: API subcommand is required")
		printHelp()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "start":
		// 检查是否是后台运行模式
		isBackground := false
		for i := 3; i < len(os.Args); i++ {
			if os.Args[i] == "--background" {
				isBackground = true
				break
			}
		}

		if isBackground {
			// 后台模式：运行API服务
			startAPIServicesBackground()
		} else {
			// 前台模式：启动后台进程
			fmt.Println("Testing API startup and access...")

			// 获取当前可执行文件路径
			execPath, err := os.Executable()
			if err != nil {
				fmt.Printf("Error getting executable path: %v\n", err)
				return
			}

			// 构建后台命令参数
			args := []string{"api", "start", "--background"}
			// 传递其他参数
			for i := 3; i < len(os.Args); i++ {
				if os.Args[i] != "--background" {
					args = append(args, os.Args[i])
				}
			}

			// 启动后台进程
			cmd := exec.Command(execPath, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err = cmd.Start()
			if err != nil {
				fmt.Printf("Error starting background process: %v\n", err)
				return
			}

			// 等待后台进程启动
			time.Sleep(3 * time.Second)

			// 测试API访问
			testAPIAccess()

			// 显示API服务状态
			fmt.Println("\nAPI Service Status:")
			fmt.Println("==================")

			// 加载API服务状态
			status, err := loadAPIServiceStatus()
			if err != nil {
				fmt.Printf("Warning: Failed to load API service status: %v\n", err)
			} else {
				// 显示API服务状态
				// 检查Public API是否在运行
				if isPortListening(status.Public.Current.Port) {
					fmt.Printf("Public API: http://%s:%d - Running\n", status.Public.Current.Address, status.Public.Current.Port)
				} else {
					fmt.Printf("Public API: http://%s:%d - Stopped\n", status.Public.Current.Address, status.Public.Current.Port)
				}
				// 检查Desktop API是否在运行
				if isPortListening(status.Desktop.Current.Port) {
					fmt.Printf("Desktop API: http://%s:%d - Running\n", status.Desktop.Current.Address, status.Desktop.Current.Port)
				} else {
					fmt.Printf("Desktop API: http://%s:%d - Stopped\n", status.Desktop.Current.Address, status.Desktop.Current.Port)
				}
				// 检查Model API是否在运行
				if isPortListening(status.Model.Current.Port) {
					fmt.Printf("Model API: http://%s:%d - Running\n", status.Model.Current.Address, status.Model.Current.Port)
				} else {
					fmt.Printf("Model API: http://%s:%d - Stopped\n", status.Model.Current.Address, status.Model.Current.Port)
				}

				// 检查API服务端点
				if isPortListening(status.Public.Current.Port) {
					// 检查Public API健康端点
					resp, err := http.Get(fmt.Sprintf("http://%s:%d/health", status.Public.Current.Address, status.Public.Current.Port))
					if err == nil && resp.StatusCode == http.StatusOK {
						fmt.Printf("Health check: http://%s:%d/health - Available\n", status.Public.Current.Address, status.Public.Current.Port)
					} else {
						fmt.Printf("Health check: http://%s:%d/health - Unavailable\n", status.Public.Current.Address, status.Public.Current.Port)
					}
					if resp != nil {
						resp.Body.Close()
					}

					// 检查Container API端点
					resp, err = http.Get(fmt.Sprintf("http://%s:%d/api/container/list", status.Public.Current.Address, status.Public.Current.Port))
					if err == nil && resp.StatusCode == http.StatusOK {
						fmt.Printf("Container API: http://%s:%d/api/container/list - Available\n", status.Public.Current.Address, status.Public.Current.Port)
					} else {
						fmt.Printf("Container API: http://%s:%d/api/container/list - Unavailable\n", status.Public.Current.Address, status.Public.Current.Port)
					}
					if resp != nil {
						resp.Body.Close()
					}
				}

				// 检查Desktop API端点
				if isPortListening(status.Desktop.Current.Port) {
					resp, err := http.Get(fmt.Sprintf("http://%s:%d/health", status.Desktop.Current.Address, status.Desktop.Current.Port))
					if err == nil && resp.StatusCode == http.StatusOK {
						fmt.Printf("Desktop API: http://%s:%d/health - Available\n", status.Desktop.Current.Address, status.Desktop.Current.Port)
					} else {
						fmt.Printf("Desktop API: http://%s:%d/health - Unavailable\n", status.Desktop.Current.Address, status.Desktop.Current.Port)
					}
					if resp != nil {
						resp.Body.Close()
					}
				}

				// 检查Model API端点
				if isPortListening(status.Model.Current.Port) {
					// 检查Model API健康端点
					resp, err := http.Get(fmt.Sprintf("http://%s:%d/health", status.Model.Current.Address, status.Model.Current.Port))
					if err == nil && resp.StatusCode == http.StatusOK {
						fmt.Printf("Model API: http://%s:%d/health - Available\n", status.Model.Current.Address, status.Model.Current.Port)
					} else {
						fmt.Printf("Model API: http://%s:%d/health - Unavailable\n", status.Model.Current.Address, status.Model.Current.Port)
					}
					if resp != nil {
						resp.Body.Close()
					}

					// 检查Model API models端点
					resp, err = http.Get(fmt.Sprintf("http://%s:%d/api/models", status.Model.Current.Address, status.Model.Current.Port))
					if err == nil && resp.StatusCode == http.StatusOK {
						fmt.Printf("Model API: http://%s:%d/api/models - Available\n", status.Model.Current.Address, status.Model.Current.Port)
					} else {
						fmt.Printf("Model API: http://%s:%d/api/models - Unavailable\n", status.Model.Current.Address, status.Model.Current.Port)
					}
					if resp != nil {
						resp.Body.Close()
					}
				}
			}

			// 显示提示信息
			fmt.Println("API services are running in the background.")
			fmt.Println("Use 'elr api stop' command to stop services...")
			fmt.Println("You can now continue using this terminal for other operations.")
		}
	case "stop":
		stopAPIServices()
	case "status":
		checkAPIStatus()
	case "config":
		configureAPI()
	default:
		fmt.Printf("Unknown API subcommand: %s\n", subcommand)
		printHelp()
		os.Exit(1)
	}
}

// testAPIAccess 测试API访问
func testAPIAccess() {
	fmt.Println("\nTesting API access...")

	// 加载API服务状态，获取当前配置
	status, err := loadAPIServiceStatus()
	if err != nil {
		fmt.Printf("Warning: Failed to load API service status: %v\n", err)
		// 使用默认状态
		status = &APIServiceStatus{
			Public: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8080,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8080,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
			Desktop: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8081,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8081,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
			Model: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8082,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8082,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
		}
	}

	// 测试Public API
	testAPICall("Public API", fmt.Sprintf("http://%s:%d/health", status.Public.Current.Address, status.Public.Current.Port))

	// 测试Desktop API
	testAPICall("Desktop API", fmt.Sprintf("http://%s:%d/health", status.Desktop.Current.Address, status.Desktop.Current.Port))

	// 测试Model API
	testAPICall("Model API", fmt.Sprintf("http://%s:%d/health", status.Model.Current.Address, status.Model.Current.Port))

	fmt.Println("API access test completed!")
}

// testAPICall 测试单个API调用
func testAPICall(apiName, url string) {
	fmt.Printf("Testing %s at %s... ", apiName, url)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("SUCCESS")
	} else {
		fmt.Printf("FAILED with status code: %d\n", resp.StatusCode)
	}
}

// findAvailablePort 查找可用端口
func findAvailablePort(startPort int) int {
	for port := startPort; port < startPort+100; port++ {
		if !isPortListening(port) {
			return port
		}
	}
	return -1
}

// max 返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// API服务管理器，用于管理所有API服务的状态
var apiServiceManager *api.APIServiceManager

// 初始化API服务管理器
func init() {
	apiServiceManager = api.NewAPIServiceManager()
}

// startAPIServices 启动 API 服务
func startAPIServices() {
	fmt.Println("Starting API services...")

	// 解析参数
	apiType := "all"
	address := ""
	publicPort := 0
	desktopPort := 0
	modelPort := 0
	
	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--api-type":
			if i+1 < len(os.Args) {
				apiType = os.Args[i+1]
				i++
			}
		case "--address":
			if i+1 < len(os.Args) {
				address = os.Args[i+1]
				i++
			}
		case "--public-port":
			if i+1 < len(os.Args) {
				fmt.Sscanf(os.Args[i+1], "%d", &publicPort)
				i++
			}
		case "--desktop-port":
			if i+1 < len(os.Args) {
				fmt.Sscanf(os.Args[i+1], "%d", &desktopPort)
				i++
			}
		case "--model-port":
			if i+1 < len(os.Args) {
				fmt.Sscanf(os.Args[i+1], "%d", &modelPort)
				i++
			}
		}
	}
	
	// 加载API服务状态，获取当前配置
	status, err := loadAPIServiceStatus()
	if err != nil {
		fmt.Printf("Error: Failed to load API service status: %v\n", err)
		os.Exit(1)
	}
	
	// 使用配置中的地址和端口
	switch apiType {
	case "all":
		if publicPort == 0 {
			publicPort = status.Public.Current.Port
			if address == "" {
				address = status.Public.Current.Address
			}
		}
		if desktopPort == 0 {
			desktopPort = status.Desktop.Current.Port
			if address == "" {
				address = status.Desktop.Current.Address
			}
		}
		if modelPort == 0 {
			modelPort = status.Model.Current.Port
			if address == "" {
				address = status.Model.Current.Address
			}
		}
	case "public":
		if publicPort == 0 {
			publicPort = status.Public.Current.Port
			if address == "" {
				address = status.Public.Current.Address
			}
		}
	case "desktop":
		if desktopPort == 0 {
			desktopPort = status.Desktop.Current.Port
			if address == "" {
				address = status.Desktop.Current.Address
			}
		}
	case "model":
		if modelPort == 0 {
			modelPort = status.Model.Current.Port
			if address == "" {
				address = status.Model.Current.Address
			}
		}
	}
	
	// 检查地址和端口是否有效
	if address == "" {
		fmt.Println("Error: No address specified and no enabled API address found")
		os.Exit(1)
	}
	
	switch apiType {
	case "all":
		if publicPort == 0 || desktopPort == 0 || modelPort == 0 {
			fmt.Println("Error: No port specified and no enabled API port found")
			os.Exit(1)
		}
	case "public":
		if publicPort == 0 {
			fmt.Println("Error: No port specified and no enabled Public API port found")
			os.Exit(1)
		}
	case "desktop":
		if desktopPort == 0 {
			fmt.Println("Error: No port specified and no enabled Desktop API port found")
			os.Exit(1)
		}
	case "model":
		if modelPort == 0 {
			fmt.Println("Error: No port specified and no enabled Model API port found")
			os.Exit(1)
		}
	}



	// 检查端口是否被占用
	switch apiType {
	case "all":
		if isPortListening(publicPort) {
			fmt.Printf("Error: Port %d is already in use for Public API\n", publicPort)
			os.Exit(1)
		}
		if isPortListening(desktopPort) {
			fmt.Printf("Error: Port %d is already in use for Desktop API\n", desktopPort)
			os.Exit(1)
		}
		if isPortListening(modelPort) {
			fmt.Printf("Error: Port %d is already in use for Model API\n", modelPort)
			os.Exit(1)
		}
	case "public":
		if isPortListening(publicPort) {
			fmt.Printf("Error: Port %d is already in use for Public API\n", publicPort)
			os.Exit(1)
		}
	case "desktop":
		if isPortListening(desktopPort) {
			fmt.Printf("Error: Port %d is already in use for Desktop API\n", desktopPort)
			os.Exit(1)
		}
	case "model":
		if isPortListening(modelPort) {
			fmt.Printf("Error: Port %d is already in use for Model API\n", modelPort)
			os.Exit(1)
		}
	}

	// 检查地址可访问性
	switch apiType {
	case "all":
		// Check if Public API address is accessible
		if !isAddressAccessible(address, publicPort) {
			fmt.Printf("Error: Public API address %s:%d is not accessible\n", address, publicPort)
			os.Exit(1)
		}
		// Check if Desktop API address is accessible
		if !isAddressAccessible(address, desktopPort) {
			fmt.Printf("Error: Desktop API address %s:%d is not accessible\n", address, desktopPort)
			os.Exit(1)
		}
		// Check if Model API address is accessible
		if !isAddressAccessible(address, modelPort) {
			fmt.Printf("Error: Model API address %s:%d is not accessible\n", address, modelPort)
			os.Exit(1)
		}
	case "public":
		// Check if Public API address is accessible
		if !isAddressAccessible(address, publicPort) {
			fmt.Printf("Error: Public API address %s:%d is not accessible\n", address, publicPort)
			os.Exit(1)
		}
	case "desktop":
		// Check if Desktop API address is accessible
		if !isAddressAccessible(address, desktopPort) {
			fmt.Printf("Error: Desktop API address %s:%d is not accessible\n", address, desktopPort)
			os.Exit(1)
		}
	case "model":
		// Check if Model API address is accessible
		if !isAddressAccessible(address, modelPort) {
			fmt.Printf("Error: Model API address %s:%d is not accessible\n", address, modelPort)
			os.Exit(1)
		}
	}

	// 保存当前进程PID
	currentPID := os.Getpid()
	switch apiType {
	case "all":
		saveAPIServicePID("public", currentPID)
		saveAPIServicePID("desktop", currentPID)
		saveAPIServicePID("model", currentPID)
	case "public":
		saveAPIServicePID("public", currentPID)
	case "desktop":
		saveAPIServicePID("desktop", currentPID)
	case "model":
		saveAPIServicePID("model", currentPID)
	}

	// 启动API服务
	switch apiType {
	case "all":
		// 启动Public API
		fmt.Println("Starting Public API...")
		publicServer := api.NewPublicAPIServer(publicPort)
		apiServiceManager.RegisterService("public", publicServer)
		if err := publicServer.Start(); err != nil {
			fmt.Printf("Error starting Public API: %v\n", err)
			os.Exit(1)
		}
		
		// 启动Desktop API
		fmt.Println("Starting Desktop API...")
		desktopServer := api.NewDesktopAPIServer(desktopPort)
		apiServiceManager.RegisterService("desktop", desktopServer)
		if err := desktopServer.Start(); err != nil {
			fmt.Printf("Error starting Desktop API: %v\n", err)
			os.Exit(1)
		}
		
		// 启动Model API
		fmt.Println("Starting Model API...")
		// 加载模型配置
		modelConfig, err := loadModelConfig()
		if err != nil {
			fmt.Printf("Error loading model config: %v\n", err)
			os.Exit(1)
		}
		
		// 创建模型管理器
		modelManager, err := model.NewModelManager(modelConfig)
		if err != nil {
			fmt.Printf("Error creating model manager: %v\n", err)
			os.Exit(1)
		}
		
		// 创建容器管理器
		containerManager, err := container.NewContainerManager(&modelConfig.Container)
		if err != nil {
			fmt.Printf("Error creating container manager: %v\n", err)
			os.Exit(1)
		}
		
		// 创建沙箱运行时
		sandboxRuntime, err := sandbox.NewSandboxRuntime(modelConfig)
		if err != nil {
			fmt.Printf("Error creating sandbox runtime: %v\n", err)
			os.Exit(1)
		}
		
		// 创建监控服务
		monitorService, err := monitor.NewMonitorService(&modelConfig.Monitoring)
		if err != nil {
			fmt.Printf("Error creating monitor service: %v\n", err)
			os.Exit(1)
		}
		
		// 创建模型API服务器配置
		serverConfig := &config.ServerConfig{
			Host: address,
			Port: modelPort,
		}
		
		// 创建并启动Model API服务器
		modelServer := api.NewModelAPIServer(serverConfig, modelManager, containerManager, sandboxRuntime, monitorService)
		apiServiceManager.RegisterService("model", modelServer)
		if err := modelServer.Start(); err != nil {
			fmt.Printf("Error starting Model API: %v\n", err)
			os.Exit(1)
		}
	case "public":
		// 启动Public API
		fmt.Println("Starting Public API...")
		publicServer := api.NewPublicAPIServer(publicPort)
		apiServiceManager.RegisterService("public", publicServer)
		if err := publicServer.Start(); err != nil {
			fmt.Printf("Error starting Public API: %v\n", err)
			os.Exit(1)
		}
	case "desktop":
		// 启动Desktop API
		fmt.Println("Starting Desktop API...")
		desktopServer := api.NewDesktopAPIServer(desktopPort)
		apiServiceManager.RegisterService("desktop", desktopServer)
		if err := desktopServer.Start(); err != nil {
			fmt.Printf("Error starting Desktop API: %v\n", err)
			os.Exit(1)
		}
	case "model":
		// 启动Model API
		fmt.Println("Starting Model API...")
		// 加载模型配置
		modelConfig, err := loadModelConfig()
		if err != nil {
			fmt.Printf("Error loading model config: %v\n", err)
			os.Exit(1)
		}
		
		// 创建模型管理器
		modelManager, err := model.NewModelManager(modelConfig)
		if err != nil {
			fmt.Printf("Error creating model manager: %v\n", err)
			os.Exit(1)
		}
		
		// 创建容器管理器
		containerManager, err := container.NewContainerManager(&modelConfig.Container)
		if err != nil {
			fmt.Printf("Error creating container manager: %v\n", err)
			os.Exit(1)
		}
		
		// 创建沙箱运行时
		sandboxRuntime, err := sandbox.NewSandboxRuntime(modelConfig)
		if err != nil {
			fmt.Printf("Error creating sandbox runtime: %v\n", err)
			os.Exit(1)
		}
		
		// 创建监控服务
		monitorService, err := monitor.NewMonitorService(&modelConfig.Monitoring)
		if err != nil {
			fmt.Printf("Error creating monitor service: %v\n", err)
			os.Exit(1)
		}
		
		// 创建模型API服务器配置
		serverConfig := &config.ServerConfig{
			Host: address,
			Port: modelPort,
		}
		
		// 创建并启动Model API服务器
		modelServer := api.NewModelAPIServer(serverConfig, modelManager, containerManager, sandboxRuntime, monitorService)
		apiServiceManager.RegisterService("model", modelServer)
		if err := modelServer.Start(); err != nil {
			fmt.Printf("Error starting Model API: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown API type: %s\n", apiType)
		printHelp()
		os.Exit(1)
	}

	// 保存API服务状态
	switch apiType {
	case "all":
		savedStatus, err := loadAPIServiceStatus()
		if err != nil {
			fmt.Printf("Warning: Failed to load API service status: %v\n", err)
			savedStatus = &APIServiceStatus{
				Public: APIServiceConfig{
					Current: APIServiceInfo{
						Address: address,
						Port:    publicPort,
						Running: true,
					},
					Available: APIServiceInfo{
						Address: address,
						Port:    publicPort,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Desktop: APIServiceConfig{
					Current: APIServiceInfo{
						Address: address,
						Port:    desktopPort,
						Running: true,
					},
					Available: APIServiceInfo{
						Address: address,
						Port:    desktopPort,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Model: APIServiceConfig{
					Current: APIServiceInfo{
						Address: address,
						Port:    modelPort,
						Running: true,
					},
					Available: APIServiceInfo{
						Address: address,
						Port:    modelPort,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
			}
		} else {
			savedStatus.Public.Current.Address = address
			savedStatus.Public.Current.Port = publicPort
			savedStatus.Public.Current.Running = true
			savedStatus.Desktop.Current.Address = address
			savedStatus.Desktop.Current.Port = desktopPort
			savedStatus.Desktop.Current.Running = true
			savedStatus.Model.Current.Address = address
			savedStatus.Model.Current.Port = modelPort
			savedStatus.Model.Current.Running = true
		}
		if err := saveAPIServiceStatus(savedStatus); err != nil {
			fmt.Printf("Warning: Failed to save API service status: %v\n", err)
		}
		
		fmt.Println("All API services started successfully!")
		fmt.Printf("Public API: http://%s:%d\n", address, publicPort)
		fmt.Printf("Desktop API: http://%s:%d\n", address, desktopPort)
		fmt.Printf("Model API: http://%s:%d\n", address, modelPort)
	case "public":
		status, err := loadAPIServiceStatus()
		if err != nil {
			fmt.Printf("Warning: Failed to load API service status: %v\n", err)
			status = &APIServiceStatus{
				Public: APIServiceConfig{
					Current: APIServiceInfo{
						Address: address,
						Port:    publicPort,
						Running: true,
					},
					Available: APIServiceInfo{
						Address: address,
						Port:    publicPort,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Desktop: APIServiceConfig{
					Current: APIServiceInfo{
						Address: "localhost",
						Port:    8081,
						Running: false,
					},
					Available: APIServiceInfo{
						Address: "localhost",
						Port:    8081,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Model: APIServiceConfig{
					Current: APIServiceInfo{
						Address: "localhost",
						Port:    8083,
						Running: false,
					},
					Available: APIServiceInfo{
						Address: "localhost",
						Port:    8083,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
			}
		} else {
			status.Public.Current.Address = address
			status.Public.Current.Port = publicPort
			status.Public.Current.Running = true
		}
		if err := saveAPIServiceStatus(status); err != nil {
			fmt.Printf("Warning: Failed to save API service status: %v\n", err)
		}
		
		fmt.Println("Public API started successfully!")
		fmt.Printf("Public API: http://%s:%d\n", address, publicPort)
	case "desktop":
		status, err := loadAPIServiceStatus()
		if err != nil {
			fmt.Printf("Warning: Failed to load API service status: %v\n", err)
			status = &APIServiceStatus{
				Public: APIServiceConfig{
					Current: APIServiceInfo{
						Address: "localhost",
						Port:    8080,
						Running: false,
					},
					Available: APIServiceInfo{
						Address: "localhost",
						Port:    8080,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Desktop: APIServiceConfig{
					Current: APIServiceInfo{
						Address: address,
						Port:    desktopPort,
						Running: true,
					},
					Available: APIServiceInfo{
						Address: address,
						Port:    desktopPort,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Model: APIServiceConfig{
					Current: APIServiceInfo{
						Address: "localhost",
						Port:    8083,
						Running: false,
					},
					Available: APIServiceInfo{
						Address: "localhost",
						Port:    8083,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
			}
		} else {
			status.Desktop.Current.Address = address
			status.Desktop.Current.Port = desktopPort
			status.Desktop.Current.Running = true
		}
		if err := saveAPIServiceStatus(status); err != nil {
			fmt.Printf("Warning: Failed to save API service status: %v\n", err)
		}
		
		fmt.Println("Desktop API started successfully!")
		fmt.Printf("Desktop API: http://%s:%d\n", address, desktopPort)
	case "model":
		status, err := loadAPIServiceStatus()
		if err != nil {
			fmt.Printf("Warning: Failed to load API service status: %v\n", err)
			status = &APIServiceStatus{
				Public: APIServiceConfig{
					Current: APIServiceInfo{
						Address: "localhost",
						Port:    8080,
						Running: false,
					},
					Available: APIServiceInfo{
						Address: "localhost",
						Port:    8080,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Desktop: APIServiceConfig{
					Current: APIServiceInfo{
						Address: "localhost",
						Port:    8081,
						Running: false,
					},
					Available: APIServiceInfo{
						Address: "localhost",
						Port:    8081,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Model: APIServiceConfig{
					Current: APIServiceInfo{
						Address: address,
						Port:    modelPort,
						Running: true,
					},
					Available: APIServiceInfo{
						Address: address,
						Port:    modelPort,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
			}
		} else {
			status.Model.Current.Address = address
			status.Model.Current.Port = modelPort
			status.Model.Current.Running = true
		}
		if err := saveAPIServiceStatus(status); err != nil {
			fmt.Printf("Warning: Failed to save API service status: %v\n", err)
		}
		
		fmt.Println("Model API started successfully!")
		fmt.Printf("Model API: http://%s:%d\n", address, modelPort)
	default:
		fmt.Printf("Unknown API type: %s\n", apiType)
		printHelp()
		os.Exit(1)
	}

	// 保持主进程运行
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// 收到信号后停止所有服务
	fmt.Println("Stopping API services...")

	// 使用API服务管理器停止所有服务
	if err := apiServiceManager.StopAllServices(); err != nil {
		fmt.Printf("Error stopping API services: %v\n", err)
	}

	// 清理PID文件
	switch apiType {
	case "all":
		deleteAPIServicePID("public")
		deleteAPIServicePID("desktop")
		deleteAPIServicePID("model")
	case "public":
		deleteAPIServicePID("public")
	case "desktop":
		deleteAPIServicePID("desktop")
	case "model":
		deleteAPIServicePID("model")
	}

	// 更新API服务状态
	status, err = loadAPIServiceStatus()
	if err == nil {
		status.Public.Current.Running = false
		status.Desktop.Current.Running = false
		status.Model.Current.Running = false
		if err := saveAPIServiceStatus(status); err != nil {
			fmt.Printf("Warning: Failed to save API service status: %v\n", err)
		}
	}

	fmt.Println("All API services stopped successfully!")
	// 退出进程
	os.Exit(0)
}

// startAPIServicesBackground 后台启动API服务
func startAPIServicesBackground() {
	// 解析参数
	apiType := "all"
	address := ""
	publicPort := 0
	desktopPort := 0
	modelPort := 0
	
	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--api-type":
			if i+1 < len(os.Args) {
				apiType = os.Args[i+1]
				i++
			}
		case "--address":
			if i+1 < len(os.Args) {
				address = os.Args[i+1]
				i++
			}
		case "--public-port":
			if i+1 < len(os.Args) {
				fmt.Sscanf(os.Args[i+1], "%d", &publicPort)
				i++
			}
		case "--desktop-port":
			if i+1 < len(os.Args) {
				fmt.Sscanf(os.Args[i+1], "%d", &desktopPort)
				i++
			}
		case "--model-port":
			if i+1 < len(os.Args) {
				fmt.Sscanf(os.Args[i+1], "%d", &modelPort)
				i++
			}
		case "public", "desktop", "model":
			// 支持直接指定API类型，如 "elr api start desktop"
			apiType = os.Args[i]
		}
	}
	
	// 加载配置文件
	elrConfig, err := loadConfig()
	if err != nil {
		fmt.Printf("Warning: Failed to load config: %v\n", err)
	}
	
	// 加载API服务状态，获取当前配置
	status, err := loadAPIServiceStatus()
	if err != nil {
		fmt.Printf("Warning: Failed to load API service status: %v\n", err)
	} else {
		// 使用配置中的地址和端口
		switch apiType {
		case "all":
			if publicPort == 0 {
				publicPort = status.Public.Current.Port
				if address == "" {
					address = status.Public.Current.Address
				}
			}
			if desktopPort == 0 {
				desktopPort = status.Desktop.Current.Port
				if address == "" {
					address = status.Desktop.Current.Address
				}
			}
			if modelPort == 0 {
				modelPort = status.Model.Current.Port
				if address == "" {
					address = status.Model.Current.Address
				}
			}
		case "public":
			if publicPort == 0 {
				publicPort = status.Public.Current.Port
				if address == "" {
					address = status.Public.Current.Address
				}
			}
		case "desktop":
			if desktopPort == 0 {
				desktopPort = status.Desktop.Current.Port
				if address == "" {
					address = status.Desktop.Current.Address
				}
			}
		case "model":
			if modelPort == 0 {
				modelPort = status.Model.Current.Port
				if address == "" {
					address = status.Model.Current.Address
				}
			}
		}
	}
	
	// 如果地址为空，使用配置文件中的默认地址
	if address == "" {
		address = "localhost"
	}
	
	// 如果端口为0，使用配置文件中的默认端口
	if publicPort == 0 && elrConfig != nil {
		publicPort = elrConfig.Network.APIPorts.PublicAPI
	}
	if desktopPort == 0 && elrConfig != nil {
		desktopPort = elrConfig.Network.APIPorts.DesktopAPI
	}
	if modelPort == 0 && elrConfig != nil {
		modelPort = elrConfig.Network.APIPorts.ModelAPI
	}
	
	// 如果端口仍然为0，使用默认值
	if publicPort == 0 {
		publicPort = 8080
	}
	if desktopPort == 0 {
		desktopPort = 8081
	}
	if modelPort == 0 {
		modelPort = 8082
	}

	// 根据API类型检查并查找可用端口
	switch apiType {
	case "all":
		// 检查并查找可用端口，确保三个端口都不同
		// 先检查Public API端口
		if isPortListening(publicPort) {
			fmt.Printf("Port %d is already in use, finding alternative port...\n", publicPort)
			publicPort = findAvailablePort(publicPort)
			if publicPort == -1 {
				fmt.Println("Error: No available port found for Public API")
				os.Exit(1)
			}
			fmt.Printf("Using alternative port %d for Public API\n", publicPort)
		}

		// 检查Desktop API端口，确保不与Public API端口冲突
		if isPortListening(desktopPort) || desktopPort == publicPort {
			fmt.Printf("Port %d is already in use or conflicts with Public API, finding alternative port...\n", desktopPort)
			// 从desktopPort+1开始查找，避免与publicPort冲突
			startPort := desktopPort + 1
			if startPort <= publicPort {
				startPort = publicPort + 1
			}
			desktopPort = findAvailablePort(startPort)
			if desktopPort == -1 {
				fmt.Println("Error: No available port found for Desktop API")
				os.Exit(1)
			}
			fmt.Printf("Using alternative port %d for Desktop API\n", desktopPort)
		}

		// 检查Model API端口，确保不与其他两个端口冲突
		if isPortListening(modelPort) || modelPort == publicPort || modelPort == desktopPort {
			fmt.Printf("Port %d is already in use or conflicts with other APIs, finding alternative port...\n", modelPort)
			// 从modelPort+1开始查找，避免与其他端口冲突
			startPort := modelPort + 1
			if startPort <= publicPort || startPort <= desktopPort {
				startPort = max(publicPort, desktopPort) + 1
			}
			modelPort = findAvailablePort(startPort)
			if modelPort == -1 {
				fmt.Println("Error: No available port found for Model API")
				os.Exit(1)
			}
			fmt.Printf("Using alternative port %d for Model API\n", modelPort)
		}
	case "public":
		// 只检查Public API端口
		if isPortListening(publicPort) {
			fmt.Printf("Port %d is already in use, finding alternative port...\n", publicPort)
			publicPort = findAvailablePort(publicPort)
			if publicPort == -1 {
				fmt.Println("Error: No available port found for Public API")
				os.Exit(1)
			}
			fmt.Printf("Using alternative port %d for Public API\n", publicPort)
		}
	case "desktop":
		// 只检查Desktop API端口
		if isPortListening(desktopPort) {
			fmt.Printf("Port %d is already in use, finding alternative port...\n", desktopPort)
			desktopPort = findAvailablePort(desktopPort)
			if desktopPort == -1 {
				fmt.Println("Error: No available port found for Desktop API")
				os.Exit(1)
			}
			fmt.Printf("Using alternative port %d for Desktop API\n", desktopPort)
		}
	case "model":
		// 只检查Model API端口
		if isPortListening(modelPort) {
			fmt.Printf("Port %d is already in use, finding alternative port...\n", modelPort)
			modelPort = findAvailablePort(modelPort)
			if modelPort == -1 {
				fmt.Println("Error: No available port found for Model API")
				os.Exit(1)
			}
			fmt.Printf("Using alternative port %d for Model API\n", modelPort)
		}
	}

	// 检查地址可访问性
	switch apiType {
	case "all":
		// Check if Public API address is accessible
		if !isAddressAccessible(address, publicPort) {
			fmt.Printf("Error: Public API address %s:%d is not accessible\n", address, publicPort)
			os.Exit(1)
		}
		// Check if Desktop API address is accessible
		if !isAddressAccessible(address, desktopPort) {
			fmt.Printf("Error: Desktop API address %s:%d is not accessible\n", address, desktopPort)
			os.Exit(1)
		}
		// Check if Model API address is accessible
		if !isAddressAccessible(address, modelPort) {
			fmt.Printf("Error: Model API address %s:%d is not accessible\n", address, modelPort)
			os.Exit(1)
		}
	case "public":
		// Check if Public API address is accessible
		if !isAddressAccessible(address, publicPort) {
			fmt.Printf("Error: Public API address %s:%d is not accessible\n", address, publicPort)
			os.Exit(1)
		}
	case "desktop":
		// Check if Desktop API address is accessible
		if !isAddressAccessible(address, desktopPort) {
			fmt.Printf("Error: Desktop API address %s:%d is not accessible\n", address, desktopPort)
			os.Exit(1)
		}
	case "model":
		// Check if Model API address is accessible
		if !isAddressAccessible(address, modelPort) {
			fmt.Printf("Error: Model API address %s:%d is not accessible\n", address, modelPort)
			os.Exit(1)
		}
	}

	// 保存当前进程PID
	currentPID := os.Getpid()
	switch apiType {
	case "all":
		saveAPIServicePID("public", currentPID)
		saveAPIServicePID("desktop", currentPID)
		saveAPIServicePID("model", currentPID)
	case "public":
		saveAPIServicePID("public", currentPID)
	case "desktop":
		saveAPIServicePID("desktop", currentPID)
	case "model":
		saveAPIServicePID("model", currentPID)
	}

	switch apiType {
	case "all":
		// 启动Public API
		fmt.Println("Starting Public API...")
		publicServer := api.NewPublicAPIServer(publicPort)
		apiServiceManager.RegisterService("public", publicServer)
		if err := publicServer.Start(); err != nil {
			fmt.Printf("Error starting Public API: %v\n", err)
			os.Exit(1)
		}
		
		// 启动Desktop API
		fmt.Println("Starting Desktop API...")
		desktopServer := api.NewDesktopAPIServer(desktopPort)
		apiServiceManager.RegisterService("desktop", desktopServer)
		if err := desktopServer.Start(); err != nil {
			fmt.Printf("Error starting Desktop API: %v\n", err)
			os.Exit(1)
		}
		
		// 启动Model API
		fmt.Println("Starting Model API...")
		// 加载模型配置
		modelConfig, err := loadModelConfig()
		if err != nil {
			fmt.Printf("Error loading model config: %v\n", err)
			os.Exit(1)
		}
		
		// 创建模型管理器
		modelManager, err := model.NewModelManager(modelConfig)
		if err != nil {
			fmt.Printf("Error creating model manager: %v\n", err)
			os.Exit(1)
		}
		
		// 创建容器管理器
		containerManager, err := container.NewContainerManager(&modelConfig.Container)
		if err != nil {
			fmt.Printf("Error creating container manager: %v\n", err)
			os.Exit(1)
		}
		
		// 创建沙箱运行时
		sandboxRuntime, err := sandbox.NewSandboxRuntime(modelConfig)
		if err != nil {
			fmt.Printf("Error creating sandbox runtime: %v\n", err)
			os.Exit(1)
		}
		
		// 创建监控服务
		monitorService, err := monitor.NewMonitorService(&modelConfig.Monitoring)
		if err != nil {
			fmt.Printf("Error creating monitor service: %v\n", err)
			os.Exit(1)
		}
		
		// 创建模型API服务器配置
		serverConfig := &config.ServerConfig{
			Host: address,
			Port: modelPort,
		}
		
		// 创建并启动Model API服务器
		modelServer := api.NewModelAPIServer(serverConfig, modelManager, containerManager, sandboxRuntime, monitorService)
		apiServiceManager.RegisterService("model", modelServer)
		if err := modelServer.Start(); err != nil {
			fmt.Printf("Error starting Model API: %v\n", err)
			os.Exit(1)
		}
	case "public":
		// 启动Public API
		fmt.Println("Starting Public API...")
		publicServer := api.NewPublicAPIServer(publicPort)
		apiServiceManager.RegisterService("public", publicServer)
		if err := publicServer.Start(); err != nil {
			fmt.Printf("Error starting Public API: %v\n", err)
			os.Exit(1)
		}
	case "desktop":
		// 启动Desktop API
		fmt.Println("Starting Desktop API...")
		desktopServer := api.NewDesktopAPIServer(desktopPort)
		apiServiceManager.RegisterService("desktop", desktopServer)
		if err := desktopServer.Start(); err != nil {
			fmt.Printf("Error starting Desktop API: %v\n", err)
			os.Exit(1)
		}
	case "model":
		// 启动Model API
		fmt.Println("Starting Model API...")
		// 加载模型配置
		modelConfig, err := loadModelConfig()
		if err != nil {
			fmt.Printf("Error loading model config: %v\n", err)
			os.Exit(1)
		}
		
		// 创建模型管理器
		modelManager, err := model.NewModelManager(modelConfig)
		if err != nil {
			fmt.Printf("Error creating model manager: %v\n", err)
			os.Exit(1)
		}
		
		// 创建容器管理器
		containerManager, err := container.NewContainerManager(&modelConfig.Container)
		if err != nil {
			fmt.Printf("Error creating container manager: %v\n", err)
			os.Exit(1)
		}
		
		// 创建沙箱运行时
		sandboxRuntime, err := sandbox.NewSandboxRuntime(modelConfig)
		if err != nil {
			fmt.Printf("Error creating sandbox runtime: %v\n", err)
			os.Exit(1)
		}
		
		// 创建监控服务
		monitorService, err := monitor.NewMonitorService(&modelConfig.Monitoring)
		if err != nil {
			fmt.Printf("Error creating monitor service: %v\n", err)
			os.Exit(1)
		}
		
		// 创建模型API服务器配置
		serverConfig := &config.ServerConfig{
			Host: address,
			Port: modelPort,
		}
		
		// 创建并启动Model API服务器
		modelServer := api.NewModelAPIServer(serverConfig, modelManager, containerManager, sandboxRuntime, monitorService)
		apiServiceManager.RegisterService("model", modelServer)
		if err := modelServer.Start(); err != nil {
			fmt.Printf("Error starting Model API: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown API type: %s\n", apiType)
		printHelp()
		os.Exit(1)
	}

	// 保存API服务状态
	switch apiType {
	case "all":
		savedStatus, err := loadAPIServiceStatus()
		if err != nil {
			fmt.Printf("Warning: Failed to load API service status: %v\n", err)
			savedStatus = &APIServiceStatus{
				Public: APIServiceConfig{
					Current: APIServiceInfo{
						Address: address,
						Port:    publicPort,
						Running: true,
					},
					Available: APIServiceInfo{
						Address: address,
						Port:    publicPort,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Desktop: APIServiceConfig{
					Current: APIServiceInfo{
						Address: address,
						Port:    desktopPort,
						Running: true,
					},
					Available: APIServiceInfo{
						Address: address,
						Port:    desktopPort,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Model: APIServiceConfig{
					Current: APIServiceInfo{
						Address: address,
						Port:    modelPort,
						Running: true,
					},
					Available: APIServiceInfo{
						Address: address,
						Port:    modelPort,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
			}
		} else {
			savedStatus.Public.Current.Address = address
			savedStatus.Public.Current.Port = publicPort
			savedStatus.Public.Current.Running = true
			savedStatus.Desktop.Current.Address = address
			savedStatus.Desktop.Current.Port = desktopPort
			savedStatus.Desktop.Current.Running = true
			savedStatus.Model.Current.Address = address
			savedStatus.Model.Current.Port = modelPort
			savedStatus.Model.Current.Running = true
		}
		if err := saveAPIServiceStatus(savedStatus); err != nil {
			fmt.Printf("Warning: Failed to save API service status: %v\n", err)
		}
		
		fmt.Println("All API services started successfully!")
		fmt.Printf("Public API: http://%s:%d\n", address, publicPort)
		fmt.Printf("Desktop API: http://%s:%d\n", address, desktopPort)
		fmt.Printf("Model API: http://%s:%d\n", address, modelPort)
	case "public":
		status, err := loadAPIServiceStatus()
		if err != nil {
			fmt.Printf("Warning: Failed to load API service status: %v\n", err)
			status = &APIServiceStatus{
				Public: APIServiceConfig{
					Current: APIServiceInfo{
						Address: address,
						Port:    publicPort,
						Running: true,
					},
					Available: APIServiceInfo{
						Address: address,
						Port:    publicPort,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Desktop: APIServiceConfig{
					Current: APIServiceInfo{
						Address: "localhost",
						Port:    8081,
						Running: false,
					},
					Available: APIServiceInfo{
						Address: "localhost",
						Port:    8081,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Model: APIServiceConfig{
					Current: APIServiceInfo{
						Address: "localhost",
						Port:    8083,
						Running: false,
					},
					Available: APIServiceInfo{
						Address: "localhost",
						Port:    8083,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
			}
		} else {
			status.Public.Current.Address = address
			status.Public.Current.Port = publicPort
			status.Public.Current.Running = true
		}
		if err := saveAPIServiceStatus(status); err != nil {
			fmt.Printf("Warning: Failed to save API service status: %v\n", err)
		}
		
		fmt.Println("Public API started successfully!")
		fmt.Printf("Public API: http://%s:%d\n", address, publicPort)
	case "desktop":
		status, err := loadAPIServiceStatus()
		if err != nil {
			fmt.Printf("Warning: Failed to load API service status: %v\n", err)
			status = &APIServiceStatus{
				Public: APIServiceConfig{
					Current: APIServiceInfo{
						Address: "localhost",
						Port:    8080,
						Running: false,
					},
					Available: APIServiceInfo{
						Address: "localhost",
						Port:    8080,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Desktop: APIServiceConfig{
					Current: APIServiceInfo{
						Address: address,
						Port:    desktopPort,
						Running: true,
					},
					Available: APIServiceInfo{
						Address: address,
						Port:    desktopPort,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Model: APIServiceConfig{
					Current: APIServiceInfo{
						Address: "localhost",
						Port:    8083,
						Running: false,
					},
					Available: APIServiceInfo{
						Address: "localhost",
						Port:    8083,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
			}
		} else {
			status.Desktop.Current.Address = address
			status.Desktop.Current.Port = desktopPort
			status.Desktop.Current.Running = true
		}
		if err := saveAPIServiceStatus(status); err != nil {
			fmt.Printf("Warning: Failed to save API service status: %v\n", err)
		}
		
		fmt.Println("Desktop API started successfully!")
		fmt.Printf("Desktop API: http://%s:%d\n", address, desktopPort)
	case "model":
		status, err := loadAPIServiceStatus()
		if err != nil {
			fmt.Printf("Warning: Failed to load API service status: %v\n", err)
			status = &APIServiceStatus{
				Public: APIServiceConfig{
					Current: APIServiceInfo{
						Address: "localhost",
						Port:    8080,
						Running: false,
					},
					Available: APIServiceInfo{
						Address: "localhost",
						Port:    8080,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Desktop: APIServiceConfig{
					Current: APIServiceInfo{
						Address: "localhost",
						Port:    8081,
						Running: false,
					},
					Available: APIServiceInfo{
						Address: "localhost",
						Port:    8081,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
				Model: APIServiceConfig{
					Current: APIServiceInfo{
						Address: address,
						Port:    modelPort,
						Running: true,
					},
					Available: APIServiceInfo{
						Address: address,
						Port:    modelPort,
						Running: false,
					},
					Alternates: []APIServiceInfo{},
				},
			}
		} else {
			status.Model.Current.Address = address
			status.Model.Current.Port = modelPort
			status.Model.Current.Running = true
		}
		if err := saveAPIServiceStatus(status); err != nil {
			fmt.Printf("Warning: Failed to save API service status: %v\n", err)
		}
		
		fmt.Println("Model API started successfully!")
		fmt.Printf("Model API: http://%s:%d\n", address, modelPort)
	default:
		fmt.Printf("Unknown API type: %s\n", apiType)
		printHelp()
		os.Exit(1)
	}

	// 保持主进程运行
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// 收到信号后停止所有服务
	fmt.Println("Stopping API services...")

	// 使用API服务管理器停止所有服务
	if err := apiServiceManager.StopAllServices(); err != nil {
		fmt.Printf("Error stopping API services: %v\n", err)
	}

	// 清理PID文件
	switch apiType {
	case "all":
		deleteAPIServicePID("public")
		deleteAPIServicePID("desktop")
		deleteAPIServicePID("model")
	case "public":
		deleteAPIServicePID("public")
	case "desktop":
		deleteAPIServicePID("desktop")
	case "model":
		deleteAPIServicePID("model")
	}

	// 更新API服务状态
	status, err = loadAPIServiceStatus()
	if err == nil {
		status.Public.Current.Running = false
		status.Desktop.Current.Running = false
		status.Model.Current.Running = false
		if err := saveAPIServiceStatus(status); err != nil {
			fmt.Printf("Warning: Failed to save API service status: %v\n", err)
		}
	}

	fmt.Println("All API services stopped successfully!")
	// 退出进程
	os.Exit(0)
}

// stopAPIServices 停止 API 服务
func stopAPIServices() {
	fmt.Println("Stopping API services...")

	// 解析参数
	apiType := "all"
	for i := 3; i < len(os.Args); i++ {
		if os.Args[i] == "--api-type" && i+1 < len(os.Args) {
			apiType = os.Args[i+1]
			break
		}
	}

	// 加载API服务状态
	status, err := loadAPIServiceStatus()
	if err != nil {
		fmt.Printf("Warning: Failed to load API service status: %v\n", err)
		// 使用默认状态
		status = &APIServiceStatus{
			Public: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8080,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8080,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
			Desktop: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8081,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8081,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
			Model: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8082,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8082,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
		}
	}

	// 获取当前进程ID，避免停止自己
	currentPID := os.Getpid()

	// 定义停止API服务的函数
	stopAPIService := func(serviceType string, port int) bool {
		fmt.Printf("Stopping %s API on port %d...\n", serviceType, port)
		if isPortListening(port) {
			// 尝试通过PID文件停止服务
			pid, err := loadAPIServicePID(serviceType)
			if err == nil && pid > 0 && pid != currentPID {
				process, err := os.FindProcess(pid)
				if err == nil {
					// 发送终止信号
					if err := process.Signal(os.Interrupt); err != nil {
						// 如果发送中断信号失败，尝试强制终止
						if err := process.Kill(); err != nil {
							fmt.Printf("Error killing %s API process: %v\n", serviceType, err)
							// 尝试使用taskkill命令
							cmd := exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", pid))
							output, err := cmd.CombinedOutput()
							if err != nil {
								fmt.Printf("Error using taskkill: %v\n", err)
								fmt.Printf("Command output: %s\n", string(output))
								return false
							} else {
								fmt.Printf("%s API process killed successfully using taskkill\n", serviceType)
								return true
							}
						} else {
							fmt.Printf("%s API process killed successfully\n", serviceType)
							return true
						}
					} else {
						fmt.Printf("%s API process interrupted successfully\n", serviceType)
						return true
					}
				} else {
					fmt.Printf("Error finding %s API process: %v\n", serviceType, err)
					// 尝试使用PowerShell命令查找占用指定端口的进程并停止
					cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`netstat -ano | findstr :%d | Select-String '\d+$' | ForEach-Object { taskkill /F /PID $_.Matches.Value }`, port))
					output, err := cmd.CombinedOutput()
					if err != nil {
						fmt.Printf("Error stopping %s API: %v\n", serviceType, err)
						fmt.Printf("Command output: %s\n", string(output))
						return false
					} else {
						fmt.Printf("%s API stopped successfully using taskkill\n", serviceType)
						return true
					}
				}
			} else {
				fmt.Printf("No PID file found for %s API, trying alternative methods...\n", serviceType)
				// 尝试使用PowerShell命令查找占用指定端口的进程并停止
				cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`netstat -ano | findstr :%d | Select-String '\d+$' | ForEach-Object { taskkill /F /PID $_.Matches.Value }`, port))
				output, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Printf("Error stopping %s API: %v\n", serviceType, err)
					fmt.Printf("Command output: %s\n", string(output))
					return false
				} else {
					fmt.Printf("%s API stopped successfully using taskkill\n", serviceType)
					return true
				}
			}
		}
		return true
	}

	switch apiType {
	case "all":
		// 停止所有API服务
		fmt.Println("Stopping all API services...")
		// 直接使用taskkill命令停止所有elr.exe进程（除了自己）
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`
			Get-Process elr | Where-Object {$_.Id -ne %d} | ForEach-Object { Stop-Process -Id $_.Id -Force }
		`, currentPID))
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error stopping all API services: %v\n", err)
			fmt.Printf("Command output: %s\n", string(output))
		} else {
			fmt.Println("All API services stopped successfully using taskkill")
		}
		
		// 清理PID文件
		deleteAPIServicePID("public")
		deleteAPIServicePID("desktop")
		deleteAPIServicePID("model")
		
		// 更新API服务状态
		status.Public.Current.Running = false
		status.Desktop.Current.Running = false
		status.Model.Current.Running = false
	case "public":
		// 停止Public API
		stopAPIService("public", status.Public.Current.Port)
		
		// 清理PID文件
		deleteAPIServicePID("public")
		
		// 更新API服务状态
		status.Public.Current.Running = false
	case "desktop":
		// 停止Desktop API
		stopAPIService("desktop", status.Desktop.Current.Port)
		
		// 清理PID文件
		deleteAPIServicePID("desktop")
		
		// 更新API服务状态
		status.Desktop.Current.Running = false
	case "model":
		// 停止Model API
		stopAPIService("model", status.Model.Current.Port)
		
		// 清理PID文件
		deleteAPIServicePID("model")
		
		// 更新API服务状态
		status.Model.Current.Running = false
	default:
		fmt.Printf("Unknown API type: %s\n", apiType)
		printHelp()
		os.Exit(1)
	}
	
	// 保存API服务状态
	if err := saveAPIServiceStatus(status); err != nil {
		fmt.Printf("Warning: Failed to save API service status: %v\n", err)
	}
	
	// 等待一点时间让服务完全停止
	time.Sleep(1 * time.Second)
	
	// 检查服务是否真正停止
	fmt.Println("Verifying API services are stopped...")
	if !isPortListening(status.Public.Current.Port) {
		fmt.Printf("Public API on port %d is stopped\n", status.Public.Current.Port)
	} else {
		fmt.Printf("Warning: Public API on port %d may still be running\n", status.Public.Current.Port)
	}
	
	if !isPortListening(status.Desktop.Current.Port) {
		fmt.Printf("Desktop API on port %d is stopped\n", status.Desktop.Current.Port)
	} else {
		fmt.Printf("Warning: Desktop API on port %d may still be running\n", status.Desktop.Current.Port)
	}
	
	if !isPortListening(status.Model.Current.Port) {
		fmt.Printf("Model API on port %d is stopped\n", status.Model.Current.Port)
	} else {
		fmt.Printf("Warning: Model API on port %d may still be running\n", status.Model.Current.Port)
	}
	
	fmt.Println("API services stop command completed!")
}

// stopServiceByPort 通过端口查找并停止进程
func stopServiceByPort(port int, serviceType string) bool {
	// 构建命令，使用PowerShell查找并停止占用指定端口的进程
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf(`
		netstat -ano | findstr :%d | ForEach-Object {
			$parts = $_.Split(' ')
			$pid = $parts[-1]
			taskkill /F /PID $pid
		}
	`, port))
	
	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error stopping %s API: %v\n", serviceType, err)
		fmt.Printf("Command output: %s\n", string(output))
		return false
	}
	
	fmt.Printf("%s API stopped successfully using taskkill\n", serviceType)
	return true
}

// isPortListening checks if a port is listening
func isPortListening(port int) bool {
	// Try to bind to 0.0.0.0:%d
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		// Port is already in use (listening)
		return true
	}
	listener.Close()

	// Try to bind to 127.0.0.1:%d
	addr = fmt.Sprintf("127.0.0.1:%d", port)
	listener, err = net.Listen("tcp", addr)
	if err != nil {
		// Port is already in use (listening)
		return true
	}
	listener.Close()

	// Port is available (not listening)
	return false
}

// isAddressAccessible checks if an address is accessible
func isAddressAccessible(address string, port int) bool {
	// Check if address is localhost or loopback
	if address == "localhost" || address == "127.0.0.1" || address == "::1" {
		// Local addresses are always accessible
		return true
	}
	
	// Try to resolve the address
	ips, err := net.LookupIP(address)
	if err != nil {
		// Address cannot be resolved
		return false
	}
	
	// Try to connect to the first resolved IP on the specified port
	for _, ip := range ips {
		// For IPv4 addresses
		if ip.To4() != nil {
			addr := fmt.Sprintf("%s:%d", ip.String(), port)
			conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
			if err != nil {
				// Cannot connect to the address:port
				return false
			}
			conn.Close()
			// Connection successful, address is accessible
			return true
		} else {
			// For IPv6 addresses
			addr := fmt.Sprintf("[%s]:%d", ip.String(), port)
			conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
			if err != nil {
				// Cannot connect to the address:port
				return false
			}
			conn.Close()
			// Connection successful, address is accessible
			return true
		}
	}
	
	// Address is not accessible
	return false
}

// APIServiceStatus stores the status of API services
type APIServiceStatus struct {
	Public  APIServiceConfig `json:"public"`
	Desktop APIServiceConfig `json:"desktop"`
	Model   APIServiceConfig `json:"model"`
}

// APIServiceConfig stores configuration information for an API service
type APIServiceConfig struct {
	Current    APIServiceInfo   `json:"current"`    // Currently enabled configuration
	Available  APIServiceInfo   `json:"available"`  // Configured but not started
	Alternates []APIServiceInfo `json:"alternates"` // Alternative configurations
}

// APIServiceInfo stores information about an API service
type APIServiceInfo struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Running bool   `json:"running"`
}





// getStatusText 获取状态文本
func getStatusText(running bool) string {
	if running {
		return "Running"
	}
	return "Stopped"
}

// getAPIServiceStatusPath returns the path to the API service status file
func getAPIServiceStatusPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".api_status.json"
	}
	return filepath.Join(homeDir, ".elr", "api_status.json")
}

// getAPIServicePIDPath returns the path to the API service PID file
func getAPIServicePIDPath(apiType string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Sprintf(".%s.pid", apiType)
	}
	return filepath.Join(homeDir, ".elr", fmt.Sprintf("%s.pid", apiType))
}

// saveAPIServicePID saves the API service PID to file
func saveAPIServicePID(apiType string, pid int) error {
	pidPath := getAPIServicePIDPath(apiType)
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(pidPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// Write PID to file
	return os.WriteFile(pidPath, []byte(fmt.Sprintf("%d", pid)), 0644)
}

// loadAPIServicePID loads the API service PID from file
func loadAPIServicePID(apiType string) (int, error) {
	pidPath := getAPIServicePIDPath(apiType)
	
	// Check if file exists
	if _, err := os.Stat(pidPath); os.IsNotExist(err) {
		return 0, err
	}
	
	// Read file
	data, err := os.ReadFile(pidPath)
	if err != nil {
		return 0, err
	}
	
	// Parse PID
	var pid int
	if _, err := fmt.Sscanf(string(data), "%d", &pid); err != nil {
		return 0, err
	}
	
	return pid, nil
}

// deleteAPIServicePID deletes the API service PID file
func deleteAPIServicePID(apiType string) error {
	pidPath := getAPIServicePIDPath(apiType)
	return os.Remove(pidPath)
}

// saveAPIServiceStatus saves the API service status to file
func saveAPIServiceStatus(status *APIServiceStatus) error {
	statusPath := getAPIServiceStatusPath()
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(statusPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	// Marshal to JSON
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return err
	}
	
	// Write to file
	return os.WriteFile(statusPath, data, 0644)
}

// loadAPIServiceStatus loads the API service status from file
func loadAPIServiceStatus() (*APIServiceStatus, error) {
	statusPath := getAPIServiceStatusPath()
	
	// Check if file exists
	if _, err := os.Stat(statusPath); os.IsNotExist(err) {
		// Return default status
		return &APIServiceStatus{
			Public: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8080,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8080,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
			Desktop: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8081,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8081,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
			Model: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8083,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8083,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
		}, nil
	}
	
	// Read file
	data, err := os.ReadFile(statusPath)
	if err != nil {
		return nil, err
	}
	
	// Try to unmarshal into new format first
	status := &APIServiceStatus{}
	if err := json.Unmarshal(data, status); err == nil {
		// Successfully loaded new format
		return status, nil
	}
	
	// If new format fails, try old format
	type OldAPIServiceStatus struct {
		Public  APIServiceInfo `json:"public"`
		Desktop APIServiceInfo `json:"desktop"`
		Model   APIServiceInfo `json:"model"`
	}
	
	oldStatus := &OldAPIServiceStatus{}
	if err := json.Unmarshal(data, oldStatus); err != nil {
		return nil, err
	}
	
	// Convert old format to new format
	newStatus := &APIServiceStatus{
		Public: APIServiceConfig{
			Current:    oldStatus.Public,
			Available:  oldStatus.Public,
			Alternates: []APIServiceInfo{},
		},
		Desktop: APIServiceConfig{
			Current:    oldStatus.Desktop,
			Available:  oldStatus.Desktop,
			Alternates: []APIServiceInfo{},
		},
		Model: APIServiceConfig{
			Current:    oldStatus.Model,
			Available:  oldStatus.Model,
			Alternates: []APIServiceInfo{},
		},
	}
	
	return newStatus, nil
}

// checkAPIStatus 检查 API 服务状态
func checkAPIStatus() {
	fmt.Println("Checking API service status...")

	// 加载API服务状态
	status, err := loadAPIServiceStatus()
	if err != nil {
		fmt.Printf("Warning: Failed to load API service status: %v\n", err)
		// 使用默认状态
		status = &APIServiceStatus{
			Public: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8080,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8080,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
			Desktop: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8081,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8081,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
			Model: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8083,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8083,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
		}
	}

	// 实际检查网络状态，不启动运行时
	fmt.Println("API Service Status:")
	fmt.Println("==================")
	
	// Check public API
	publicStatus := "Stopped"
	if isPortListening(status.Public.Current.Port) {
		publicStatus = "Running"
	}
	fmt.Printf("Public API: http://%s:%d - %s\n", status.Public.Current.Address, status.Public.Current.Port, publicStatus)
	
	// Check desktop API
	desktopStatus := "Stopped"
	if isPortListening(status.Desktop.Current.Port) {
		desktopStatus = "Running"
	}
	fmt.Printf("Desktop API: http://%s:%d - %s\n", status.Desktop.Current.Address, status.Desktop.Current.Port, desktopStatus)
	
	// Check model API
	modelStatus := "Stopped"
	if isPortListening(status.Model.Current.Port) {
		modelStatus = "Running"
	}
	fmt.Printf("Model API: http://%s:%d - %s\n", status.Model.Current.Address, status.Model.Current.Port, modelStatus)
	
	// Check specific endpoints if services are running
	if isPortListening(status.Public.Current.Port) {
		// Check health endpoint
		resp, err := http.Get(fmt.Sprintf("http://%s:%d/health", status.Public.Current.Address, status.Public.Current.Port))
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Printf("Health check: http://%s:%d/health - Available\n", status.Public.Current.Address, status.Public.Current.Port)
		} else {
			fmt.Printf("Health check: http://%s:%d/health - Unavailable\n", status.Public.Current.Address, status.Public.Current.Port)
		}
		
		// Check container API endpoint
		resp, err = http.Get(fmt.Sprintf("http://%s:%d/api/container/list", status.Public.Current.Address, status.Public.Current.Port))
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Printf("Container API: http://%s:%d/api/container/list - Available\n", status.Public.Current.Address, status.Public.Current.Port)
		} else {
			fmt.Printf("Container API: http://%s:%d/api/container/list - Unavailable\n", status.Public.Current.Address, status.Public.Current.Port)
		}
	}
	
	// Check Desktop API endpoints if running
	if isPortListening(status.Desktop.Current.Port) {
		// Check Desktop API health endpoint
		resp, err := http.Get(fmt.Sprintf("http://%s:%d/api/desktop/health", status.Desktop.Current.Address, status.Desktop.Current.Port))
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Printf("Desktop API: http://%s:%d/api/desktop/health - Available\n", status.Desktop.Current.Address, status.Desktop.Current.Port)
		} else {
			fmt.Printf("Desktop API: http://%s:%d/api/desktop/health - Unavailable\n", status.Desktop.Current.Address, status.Desktop.Current.Port)
		}
	}
	
	// Check Model API endpoints if running
	if isPortListening(status.Model.Current.Port) {
		// Check Model API health endpoint
		resp, err := http.Get(fmt.Sprintf("http://%s:%d/health", status.Model.Current.Address, status.Model.Current.Port))
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Printf("Model API: http://%s:%d/health - Available\n", status.Model.Current.Address, status.Model.Current.Port)
		} else {
			fmt.Printf("Model API: http://%s:%d/health - Unavailable\n", status.Model.Current.Address, status.Model.Current.Port)
		}
		
		// Check Model API models endpoint
		resp, err = http.Get(fmt.Sprintf("http://%s:%d/api/models", status.Model.Current.Address, status.Model.Current.Port))
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Printf("Model API: http://%s:%d/api/models - Available\n", status.Model.Current.Address, status.Model.Current.Port)
		} else {
			fmt.Printf("Model API: http://%s:%d/api/models - Unavailable\n", status.Model.Current.Address, status.Model.Current.Port)
		}
	}
}

// configureAPI 配置 API 地址和端口
func configureAPI() {
	fmt.Println("Configuring API services...")

	// 解析参数
	if len(os.Args) < 4 {
		fmt.Println("Error: API config subcommand is required")
		fmt.Println("Usage: elr api config <subcommand> [options]")
		fmt.Println("Subcommands:")
		fmt.Println("  set      Set alternative API address and port")
		fmt.Println("  enable   Enable alternative API address and port")
		fmt.Println("  disable  Disable and remove API address and port from alternatives")
		fmt.Println("  clear    Clear alternative configurations, only keep enabled ones")
		fmt.Println("  list     List API configurations")
		os.Exit(1)
	}

	subcommand := os.Args[3]
	apiType := ""
	address := "localhost"
	port := 0
	index := -1

	for i := 4; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--api-type":
			if i+1 < len(os.Args) {
				apiType = os.Args[i+1]
				i++
			}
		case "--address":
			if i+1 < len(os.Args) {
				address = os.Args[i+1]
				i++
			}
		case "--port":
			if i+1 < len(os.Args) {
				fmt.Sscanf(os.Args[i+1], "%d", &port)
				i++
			}
		case "--index":
			if i+1 < len(os.Args) {
				fmt.Sscanf(os.Args[i+1], "%d", &index)
				i++
			}
		}
	}

	// 加载API服务状态
	status, err := loadAPIServiceStatus()
	if err != nil {
		fmt.Printf("Warning: Failed to load API service status: %v\n", err)
		status = &APIServiceStatus{
			Public: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8080,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8080,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
			Desktop: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8081,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8081,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
			Model: APIServiceConfig{
				Current: APIServiceInfo{
					Address: "localhost",
					Port:    8083,
					Running: false,
				},
				Available: APIServiceInfo{
					Address: "localhost",
					Port:    8083,
					Running: false,
				},
				Alternates: []APIServiceInfo{},
			},
		}
	}

	switch subcommand {
	case "set":
		if apiType == "" {
			fmt.Println("Error: API type is required")
			os.Exit(1)
		}
		if port == 0 {
			fmt.Println("Error: Port is required")
			os.Exit(1)
		}

		// 检查地址和端口可用性
		if isPortListening(port) {
			fmt.Printf("Warning: Port %d is already in use, but will be added to alternative configurations\n", port)
		}

		// 检查备选列表中是否已存在相同的地址和端口
		config := APIServiceInfo{
			Address: address,
			Port:    port,
			Running: false,
		}

		duplicateFound := false
		switch apiType {
		case "public":
			for _, alt := range status.Public.Alternates {
				if alt.Address == address && alt.Port == port {
					fmt.Printf("Error: Address and port %s:%d already exists in Public API alternatives\n", address, port)
					duplicateFound = true
					break
				}
			}
			if !duplicateFound {
				status.Public.Alternates = append(status.Public.Alternates, config)
				fmt.Printf("Added alternative Public API configuration: %s:%d\n", address, port)
			}
		case "desktop":
			for _, alt := range status.Desktop.Alternates {
				if alt.Address == address && alt.Port == port {
					fmt.Printf("Error: Address and port %s:%d already exists in Desktop API alternatives\n", address, port)
					duplicateFound = true
					break
				}
			}
			if !duplicateFound {
				status.Desktop.Alternates = append(status.Desktop.Alternates, config)
				fmt.Printf("Added alternative Desktop API configuration: %s:%d\n", address, port)
			}
		case "model":
			for _, alt := range status.Model.Alternates {
				if alt.Address == address && alt.Port == port {
					fmt.Printf("Error: Address and port %s:%d already exists in Model API alternatives\n", address, port)
					duplicateFound = true
					break
				}
			}
			if !duplicateFound {
				status.Model.Alternates = append(status.Model.Alternates, config)
				fmt.Printf("Added alternative Model API configuration: %s:%d\n", address, port)
			}
		default:
			fmt.Printf("Unknown API type: %s\n", apiType)
			os.Exit(1)
		}

		if duplicateFound {
			os.Exit(1)
		}

		// 保存配置
		if err := saveAPIServiceStatus(status); err != nil {
			fmt.Printf("Error saving API service status: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("API configuration updated successfully!")

	case "enable":
		if apiType == "" {
			fmt.Println("Error: API type is required")
			os.Exit(1)
		}
		if index < 0 {
			fmt.Println("Error: Alternative index is required")
			os.Exit(1)
		}

		switch apiType {
		case "public":
			if index >= len(status.Public.Alternates) {
				fmt.Printf("Error: Invalid alternative index. Public API has %d alternatives\n", len(status.Public.Alternates))
				os.Exit(1)
			}
			// 检查备选配置的可用性
			altConfig := status.Public.Alternates[index]
			if isPortListening(altConfig.Port) {
				fmt.Printf("Error: Port %d is already in use\n", altConfig.Port)
				os.Exit(1)
			}
			// 将当前配置移至备选列表
			status.Public.Alternates = append(status.Public.Alternates, status.Public.Current)
			// 启用备选配置
			status.Public.Current = altConfig
			// 从备选列表中移除
			status.Public.Alternates = append(status.Public.Alternates[:index], status.Public.Alternates[index+1:]...)
			fmt.Printf("Enabled Public API configuration: %s:%d\n", altConfig.Address, altConfig.Port)
		case "desktop":
			if index >= len(status.Desktop.Alternates) {
				fmt.Printf("Error: Invalid alternative index. Desktop API has %d alternatives\n", len(status.Desktop.Alternates))
				os.Exit(1)
			}
			// 检查备选配置的可用性
			altConfig := status.Desktop.Alternates[index]
			if isPortListening(altConfig.Port) {
				fmt.Printf("Error: Port %d is already in use\n", altConfig.Port)
				os.Exit(1)
			}
			// 将当前配置移至备选列表
			status.Desktop.Alternates = append(status.Desktop.Alternates, status.Desktop.Current)
			// 启用备选配置
			status.Desktop.Current = altConfig
			// 从备选列表中移除
			status.Desktop.Alternates = append(status.Desktop.Alternates[:index], status.Desktop.Alternates[index+1:]...)
			fmt.Printf("Enabled Desktop API configuration: %s:%d\n", altConfig.Address, altConfig.Port)
		case "model":
			if index >= len(status.Model.Alternates) {
				fmt.Printf("Error: Invalid alternative index. Model API has %d alternatives\n", len(status.Model.Alternates))
				os.Exit(1)
			}
			// 检查备选配置的可用性
			altConfig := status.Model.Alternates[index]
			if isPortListening(altConfig.Port) {
				fmt.Printf("Error: Port %d is already in use\n", altConfig.Port)
				os.Exit(1)
			}
			// 将当前配置移至备选列表
			status.Model.Alternates = append(status.Model.Alternates, status.Model.Current)
			// 启用备选配置
			status.Model.Current = altConfig
			// 从备选列表中移除
			status.Model.Alternates = append(status.Model.Alternates[:index], status.Model.Alternates[index+1:]...)
			fmt.Printf("Enabled Model API configuration: %s:%d\n", altConfig.Address, altConfig.Port)
		default:
			fmt.Printf("Unknown API type: %s\n", apiType)
			os.Exit(1)
		}

		// 保存配置
		if err := saveAPIServiceStatus(status); err != nil {
			fmt.Printf("Error saving API service status: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("API configuration enabled successfully!")

	case "disable":
		if apiType == "" {
			fmt.Println("Error: API type is required")
			os.Exit(1)
		}
		if index < 0 {
			fmt.Println("Error: Alternative index is required")
			os.Exit(1)
		}

		switch apiType {
		case "public":
			if index >= len(status.Public.Alternates) {
				fmt.Printf("Error: Invalid alternative index. Public API has %d alternatives\n", len(status.Public.Alternates))
				os.Exit(1)
			}
			// 获取要禁用的配置
			altConfig := status.Public.Alternates[index]
			// 检查API是否在运行，如果是则停止它
			if isPortListening(altConfig.Port) {
				fmt.Printf("Stopping Public API on port %d...\n", altConfig.Port)
				// 停止API服务
				stopAPIServices()
			}
			// 检查当前配置是否是要禁用的配置
			if status.Public.Current.Address == altConfig.Address && status.Public.Current.Port == altConfig.Port {
				// 如果有其他备选配置，则启用第一个
				if len(status.Public.Alternates) > 1 {
					// 移除要禁用的配置
					status.Public.Alternates = append(status.Public.Alternates[:index], status.Public.Alternates[index+1:]...)
					// 启用第一个备选配置
					newConfig := status.Public.Alternates[0]
					// 将当前配置移至备选列表
					status.Public.Alternates = append(status.Public.Alternates, status.Public.Current)
					// 启用备选配置
					status.Public.Current = newConfig
					// 从备选列表中移除
					status.Public.Alternates = status.Public.Alternates[1:]
					fmt.Printf("Enabled alternative Public API configuration: %s:%d\n", newConfig.Address, newConfig.Port)
				} else {
					// 如果没有其他备选配置，则使用默认配置
					fmt.Println("No alternative configurations available, using default")
					status.Public.Current = APIServiceInfo{
						Address: "localhost",
						Port:    8080,
						Running: false,
					}
					// 清空备选列表
					status.Public.Alternates = []APIServiceInfo{}
				}
			} else {
				// 从备选列表中删除指定的配置
				status.Public.Alternates = append(status.Public.Alternates[:index], status.Public.Alternates[index+1:]...)
			}
			fmt.Printf("Disabled and removed Public API configuration: %s:%d\n", altConfig.Address, altConfig.Port)
		case "desktop":
			if index >= len(status.Desktop.Alternates) {
				fmt.Printf("Error: Invalid alternative index. Desktop API has %d alternatives\n", len(status.Desktop.Alternates))
				os.Exit(1)
			}
			// 获取要禁用的配置
			altConfig := status.Desktop.Alternates[index]
			// 检查API是否在运行，如果是则停止它
			if isPortListening(altConfig.Port) {
				fmt.Printf("Stopping Desktop API on port %d...\n", altConfig.Port)
				// 停止API服务
				stopAPIServices()
			}
			// 检查当前配置是否是要禁用的配置
			if status.Desktop.Current.Address == altConfig.Address && status.Desktop.Current.Port == altConfig.Port {
				// 如果有其他备选配置，则启用第一个
				if len(status.Desktop.Alternates) > 1 {
					// 移除要禁用的配置
					status.Desktop.Alternates = append(status.Desktop.Alternates[:index], status.Desktop.Alternates[index+1:]...)
					// 启用第一个备选配置
					newConfig := status.Desktop.Alternates[0]
					// 将当前配置移至备选列表
					status.Desktop.Alternates = append(status.Desktop.Alternates, status.Desktop.Current)
					// 启用备选配置
					status.Desktop.Current = newConfig
					// 从备选列表中移除
					status.Desktop.Alternates = status.Desktop.Alternates[1:]
					fmt.Printf("Enabled alternative Desktop API configuration: %s:%d\n", newConfig.Address, newConfig.Port)
				} else {
					// 如果没有其他备选配置，则使用默认配置
					fmt.Println("No alternative configurations available, using default")
					status.Desktop.Current = APIServiceInfo{
						Address: "localhost",
						Port:    8081,
						Running: false,
					}
					// 清空备选列表
					status.Desktop.Alternates = []APIServiceInfo{}
				}
			} else {
				// 从备选列表中删除指定的配置
				status.Desktop.Alternates = append(status.Desktop.Alternates[:index], status.Desktop.Alternates[index+1:]...)
			}
			fmt.Printf("Disabled and removed Desktop API configuration: %s:%d\n", altConfig.Address, altConfig.Port)
		case "model":
			if index >= len(status.Model.Alternates) {
				fmt.Printf("Error: Invalid alternative index. Model API has %d alternatives\n", len(status.Model.Alternates))
				os.Exit(1)
			}
			// 获取要禁用的配置
			altConfig := status.Model.Alternates[index]
			// 检查API是否在运行，如果是则停止它
			if isPortListening(altConfig.Port) {
				fmt.Printf("Stopping Model API on port %d...\n", altConfig.Port)
				// 停止API服务
				stopAPIServices()
			}
			// 检查当前配置是否是要禁用的配置
			if status.Model.Current.Address == altConfig.Address && status.Model.Current.Port == altConfig.Port {
				// 如果有其他备选配置，则启用第一个
				if len(status.Model.Alternates) > 1 {
					// 移除要禁用的配置
					status.Model.Alternates = append(status.Model.Alternates[:index], status.Model.Alternates[index+1:]...)
					// 启用第一个备选配置
					newConfig := status.Model.Alternates[0]
					// 将当前配置移至备选列表
					status.Model.Alternates = append(status.Model.Alternates, status.Model.Current)
					// 启用备选配置
					status.Model.Current = newConfig
					// 从备选列表中移除
					status.Model.Alternates = status.Model.Alternates[1:]
					fmt.Printf("Enabled alternative Model API configuration: %s:%d\n", newConfig.Address, newConfig.Port)
				} else {
					// 如果没有其他备选配置，则使用默认配置
					fmt.Println("No alternative configurations available, using default")
					status.Model.Current = APIServiceInfo{
						Address: "localhost",
						Port:    8083,
						Running: false,
					}
					// 清空备选列表
					status.Model.Alternates = []APIServiceInfo{}
				}
			} else {
				// 从备选列表中删除指定的配置
				status.Model.Alternates = append(status.Model.Alternates[:index], status.Model.Alternates[index+1:]...)
			}
			fmt.Printf("Disabled and removed Model API configuration: %s:%d\n", altConfig.Address, altConfig.Port)
		default:
			fmt.Printf("Unknown API type: %s\n", apiType)
			os.Exit(1)
		}

		// 保存配置
		if err := saveAPIServiceStatus(status); err != nil {
			fmt.Printf("Error saving API service status: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("API configuration disabled successfully!")

	case "clear":
		if apiType == "" {
			// 清除所有API的备选配置
			status.Public.Alternates = []APIServiceInfo{}
			status.Desktop.Alternates = []APIServiceInfo{}
			status.Model.Alternates = []APIServiceInfo{}
			fmt.Println("Cleared all alternative configurations")
		} else {
			// 清除指定API的备选配置
			switch apiType {
			case "public":
				status.Public.Alternates = []APIServiceInfo{}
				fmt.Println("Cleared Public API alternative configurations")
			case "desktop":
				status.Desktop.Alternates = []APIServiceInfo{}
				fmt.Println("Cleared Desktop API alternative configurations")
			case "model":
				status.Model.Alternates = []APIServiceInfo{}
				fmt.Println("Cleared Model API alternative configurations")
			default:
				fmt.Printf("Unknown API type: %s\n", apiType)
				os.Exit(1)
			}
		}

		// 保存配置
		if err := saveAPIServiceStatus(status); err != nil {
			fmt.Printf("Error saving API service status: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("API configuration cleared successfully!")

	case "list":
		if apiType == "" {
			// 列出所有API配置
			fmt.Println("API Configurations:")
			fmt.Println("==================")
			
			// Public API
			fmt.Println("Public API:")
			fmt.Printf("  Current: %s:%d (Running: %t)\n", status.Public.Current.Address, status.Public.Current.Port, status.Public.Current.Running)
			fmt.Printf("  Available: %s:%d\n", status.Public.Available.Address, status.Public.Available.Port)
			if len(status.Public.Alternates) > 0 {
				fmt.Println("  Alternatives:")
				for i, alt := range status.Public.Alternates {
					fmt.Printf("    %d: %s:%d\n", i, alt.Address, alt.Port)
				}
			} else {
				fmt.Println("  Alternatives: None")
			}
			
			// Desktop API
			fmt.Println("\nDesktop API:")
			fmt.Printf("  Current: %s:%d (Running: %t)\n", status.Desktop.Current.Address, status.Desktop.Current.Port, status.Desktop.Current.Running)
			fmt.Printf("  Available: %s:%d\n", status.Desktop.Available.Address, status.Desktop.Available.Port)
			if len(status.Desktop.Alternates) > 0 {
				fmt.Println("  Alternatives:")
				for i, alt := range status.Desktop.Alternates {
					fmt.Printf("    %d: %s:%d\n", i, alt.Address, alt.Port)
				}
			} else {
				fmt.Println("  Alternatives: None")
			}
			
			// Model API
			fmt.Println("\nModel API:")
			fmt.Printf("  Current: %s:%d (Running: %t)\n", status.Model.Current.Address, status.Model.Current.Port, status.Model.Current.Running)
			fmt.Printf("  Available: %s:%d\n", status.Model.Available.Address, status.Model.Available.Port)
			if len(status.Model.Alternates) > 0 {
				fmt.Println("  Alternatives:")
				for i, alt := range status.Model.Alternates {
					fmt.Printf("    %d: %s:%d\n", i, alt.Address, alt.Port)
				}
			} else {
				fmt.Println("  Alternatives: None")
			}
		} else {
			// 列出指定API配置
			switch apiType {
			case "public":
				fmt.Println("Public API Configurations:")
				fmt.Println("========================")
				fmt.Printf("Current: %s:%d (Running: %t)\n", status.Public.Current.Address, status.Public.Current.Port, status.Public.Current.Running)
				fmt.Printf("Available: %s:%d\n", status.Public.Available.Address, status.Public.Available.Port)
				if len(status.Public.Alternates) > 0 {
					fmt.Println("Alternatives:")
					for i, alt := range status.Public.Alternates {
						fmt.Printf("  %d: %s:%d\n", i, alt.Address, alt.Port)
					}
				} else {
					fmt.Println("Alternatives: None")
				}
			case "desktop":
				fmt.Println("Desktop API Configurations:")
				fmt.Println("==========================")
				fmt.Printf("Current: %s:%d (Running: %t)\n", status.Desktop.Current.Address, status.Desktop.Current.Port, status.Desktop.Current.Running)
				fmt.Printf("Available: %s:%d\n", status.Desktop.Available.Address, status.Desktop.Available.Port)
				if len(status.Desktop.Alternates) > 0 {
					fmt.Println("Alternatives:")
					for i, alt := range status.Desktop.Alternates {
						fmt.Printf("  %d: %s:%d\n", i, alt.Address, alt.Port)
					}
				} else {
					fmt.Println("Alternatives: None")
				}
			case "model":
				fmt.Println("Model API Configurations:")
				fmt.Println("=======================")
				fmt.Printf("Current: %s:%d (Running: %t)\n", status.Model.Current.Address, status.Model.Current.Port, status.Model.Current.Running)
				fmt.Printf("Available: %s:%d\n", status.Model.Available.Address, status.Model.Available.Port)
				if len(status.Model.Alternates) > 0 {
					fmt.Println("Alternatives:")
					for i, alt := range status.Model.Alternates {
						fmt.Printf("  %d: %s:%d\n", i, alt.Address, alt.Port)
					}
				} else {
					fmt.Println("Alternatives: None")
				}
			default:
				fmt.Printf("Unknown API type: %s\n", apiType)
				os.Exit(1)
			}
		}

	default:
		fmt.Printf("Unknown API config subcommand: %s\n", subcommand)
		fmt.Println("Usage: elr api config <subcommand> [options]")
		fmt.Println("Subcommands:")
		fmt.Println("  set      Set alternative API address and port")
		fmt.Println("  enable   Enable alternative API address and port")
		fmt.Println("  disable  Disable and remove API address and port from alternatives")
		fmt.Println("  clear    Clear alternative configurations, only keep enabled ones")
		fmt.Println("  list     List API configurations")
		os.Exit(1)
	}
}

// setupCommand 处理系统设置命令
func setupCommand() {
	fmt.Println("Setting up ELR system...")

	// 解析参数
	isolationType := ""

	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--isolation" && i+1 < len(os.Args) {
			isolationType = os.Args[i+1]
			break
		}
	}

	if isolationType == "" {
		// 没有指定隔离类型，显示可用选项
		fmt.Println("Available isolation options:")
		fmt.Println("1. windows-container: Windows Containers (requires Windows container feature)")
		fmt.Println("2. wsl: Windows Subsystem for Linux (requires WSL feature)")
		fmt.Println("3. basic: Basic file system isolation (no additional requirements)")
		fmt.Println()
		fmt.Println("Usage: elr setup --isolation <isolation-type>")
		os.Exit(1)
	}

	// 检查隔离类型是否有效
	validIsolationTypes := map[string]bool{
		"windows-container": true,
		"wsl":              true,
		"basic":            true,
	}

	if !validIsolationTypes[isolationType] {
		fmt.Printf("Invalid isolation type: %s\n", isolationType)
		fmt.Println("Valid options: windows-container, wsl, basic")
		os.Exit(1)
	}

	// 检查系统环境
	switch isolationType {
	case "windows-container":
		fmt.Println("Checking Windows Container feature...")
		// 检查 Windows Container 功能是否启用
		if !isWindowsContainerAvailable() {
			fmt.Println("Windows Container feature is not available.")
			fmt.Println("To enable Windows Container feature:")
			fmt.Println("1. Open 'Turn Windows features on or off'")
			fmt.Println("2. Enable 'Containers' feature")
			fmt.Println("3. Restart your computer")
			os.Exit(1)
		}
	case "wsl":
		fmt.Println("Checking WSL feature...")
		// 检查 WSL 功能是否启用
		if !isWSLAvailable() {
			fmt.Println("WSL feature is not available.")
			fmt.Println("To enable WSL feature:")
			fmt.Println("1. Open PowerShell as administrator")
			fmt.Println("2. Run: wsl --install")
			fmt.Println("3. Restart your computer")
			os.Exit(1)
		}
	case "basic":
		fmt.Println("Basic isolation requires no additional features.")
	}

	// 更新配置文件
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// 更新隔离类型
	config.Platform.Windows.IsolationType = isolationType
	if isolationType == "windows-container" {
		config.Platform.Windows.UseContainers = true
	} else if isolationType == "wsl" {
		config.Platform.Windows.UseWSL = true
	}

	// 保存配置
	if err := saveConfig(config); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Isolation type set to: %s\n", isolationType)
	fmt.Println("Setup completed successfully!")
}

// isWindowsContainerAvailable 检查 Windows Container 功能是否可用
func isWindowsContainerAvailable() bool {
	// 尝试执行 powershell 命令检查容器功能
	cmd := exec.Command("powershell", "Get-WindowsOptionalFeature", "-FeatureName", "Containers", "-Online")
	err := cmd.Run()
	return err == nil
}

// isWSLAvailable 检查 WSL 功能是否可用
func isWSLAvailable() bool {
	// 尝试执行 wsl 命令
	cmd := exec.Command("wsl", "--version")
	err := cmd.Run()
	return err == nil
}

// saveConfig 保存配置到文件
func saveConfig(config *elr.Config) error {
	configPath := os.Getenv("ELR_CONFIG")
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configPath = filepath.Join(homeDir, ".elr", "config.yaml")
	}

	// 创建配置目录
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// 序列化配置
	configBytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// 写入配置文件
	return os.WriteFile(configPath, configBytes, 0644)
}

// 文件系统管理命令处理函数
func fsCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Error: FS subcommand is required")
		printHelp()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "upload":
		uploadFile()
	case "download":
		downloadFile()
	case "set-dir":
		setFileDirectory()
	case "get-dir":
		getFileDirectory()
	default:
		fmt.Printf("Unknown fs subcommand: %s\n", subcommand)
		printHelp()
		os.Exit(1)
	}
}

// uploadFile 上传文件到容器
func uploadFile() {
	fmt.Println("Uploading file to container...")

	// 解析参数
	containerID := ""
	localPath := ""
	containerPath := ""
	token := ""

	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--id":
			if i+1 < len(os.Args) {
				containerID = os.Args[i+1]
			}
		case "--local-path":
			if i+1 < len(os.Args) {
				localPath = os.Args[i+1]
			}
		case "--container-path":
			if i+1 < len(os.Args) {
				containerPath = os.Args[i+1]
			}
		case "--token":
			if i+1 < len(os.Args) {
				token = os.Args[i+1]
			}
		}
	}

	if containerID == "" || localPath == "" || containerPath == "" {
		fmt.Println("Error: Container ID, local path, and container path are required")
		os.Exit(1)
	}

	if token == "" {
		fmt.Println("Error: Token is required for authentication")
		os.Exit(1)
	}

	// 获取运行时
	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	// 验证管理员权限
	valid, message := runtime.AdminManager.ValidateAdmin(token, containerID, "write")
	if !valid {
		fmt.Printf("Error: %s\n", message)
		os.Exit(1)
	}

	// 获取容器
	container, err := runtime.GetContainer(containerID)
	if err != nil {
		fmt.Printf("Error getting container: %v\n", err)
		os.Exit(1)
	}

	// 上传文件
	if err := container.UploadFile(localPath, containerPath, token); err != nil {
		fmt.Printf("Error uploading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File uploaded successfully!")
}

// downloadFile 从容器下载文件
func downloadFile() {
	fmt.Println("Downloading file from container...")

	// 解析参数
	containerID := ""
	containerPath := ""
	localPath := ""
	token := ""

	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--id":
			if i+1 < len(os.Args) {
				containerID = os.Args[i+1]
			}
		case "--container-path":
			if i+1 < len(os.Args) {
				containerPath = os.Args[i+1]
			}
		case "--local-path":
			if i+1 < len(os.Args) {
				localPath = os.Args[i+1]
			}
		case "--token":
			if i+1 < len(os.Args) {
				token = os.Args[i+1]
			}
		}
	}

	if containerID == "" || containerPath == "" || localPath == "" {
		fmt.Println("Error: Container ID, container path, and local path are required")
		os.Exit(1)
	}

	if token == "" {
		fmt.Println("Error: Token is required for authentication")
		os.Exit(1)
	}

	// 获取运行时
	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	// 验证管理员权限
	valid, message := runtime.AdminManager.ValidateAdmin(token, containerID, "read")
	if !valid {
		fmt.Printf("Error: %s\n", message)
		os.Exit(1)
	}

	// 获取容器
	container, err := runtime.GetContainer(containerID)
	if err != nil {
		fmt.Printf("Error getting container: %v\n", err)
		os.Exit(1)
	}

	// 下载文件
	if err := container.DownloadFile(containerPath, localPath, token); err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File downloaded successfully!")
}

// setFileDirectory 设置文件类型的目录
func setFileDirectory() {
	fmt.Println("Setting directory for file type...")

	// 解析参数
	fileType := ""
	directory := ""
	token := ""

	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--file-type":
			if i+1 < len(os.Args) {
				fileType = os.Args[i+1]
			}
		case "--directory":
			if i+1 < len(os.Args) {
				directory = os.Args[i+1]
			}
		case "--token":
			if i+1 < len(os.Args) {
				token = os.Args[i+1]
			}
		}
	}

	if fileType == "" || directory == "" {
		fmt.Println("Error: File type and directory are required")
		os.Exit(1)
	}

	if token == "" {
		fmt.Println("Error: Token is required for authentication")
		os.Exit(1)
	}

	// 获取运行时
	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	// 验证管理员权限
	valid, message := runtime.AdminManager.ValidateAdmin(token, "", "manage")
	if !valid {
		fmt.Printf("Error: %s\n", message)
		os.Exit(1)
	}

	// 设置文件类型目录
	if err := runtime.SetFileDirectory(fileType, directory); err != nil {
		fmt.Printf("Error setting file directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Directory set successfully! File type: %s, Directory: %s\n", fileType, directory)
}

// getFileDirectory 获取文件类型的目录
func getFileDirectory() {
	fmt.Println("Getting directory for file type...")

	// 解析参数
	fileType := ""
	token := ""

	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--file-type":
			if i+1 < len(os.Args) {
				fileType = os.Args[i+1]
			}
		case "--token":
			if i+1 < len(os.Args) {
				token = os.Args[i+1]
			}
		}
	}

	if fileType == "" {
		fmt.Println("Error: File type is required")
		os.Exit(1)
	}

	if token == "" {
		fmt.Println("Error: Token is required for authentication")
		os.Exit(1)
	}

	// 获取运行时
	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	// 验证管理员权限
	valid, message := runtime.AdminManager.ValidateAdmin(token, "", "read")

	if !valid {
		fmt.Printf("Error: %s\n", message)
		os.Exit(1)
	}

	// 获取文件类型目录
	dir, err := runtime.GetFileDirectory(fileType)
	if err != nil {
		fmt.Printf("Error getting file directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Directory for file type %s: %s\n", fileType, dir)
}

// 管理员管理命令处理函数
func adminCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Admin subcommand is required")
		printHelp()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "create":
		createAdmin()
	case "list":
		listAdmins()
	case "add-permission":
		addAdminPermission()
	case "remove-permission":
		removeAdminPermission()
	default:
		fmt.Printf("Unknown admin subcommand: %s\n", subcommand)
		printHelp()
		os.Exit(1)
	}
}

// createAdmin 创建新管理员
func createAdmin() {
	fmt.Println("Creating admin...")

	// 解析参数
	username := ""
	role := "regular"

	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--username":
			if i+1 < len(os.Args) {
				username = os.Args[i+1]
			}
		case "--role":
			if i+1 < len(os.Args) {
				role = os.Args[i+1]
			}
		}
	}

	if username == "" {
		fmt.Println("Error: Username is required")
		os.Exit(1)
	}

	// 获取运行时
	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	// 创建管理员
	token, err := runtime.AdminManager.CreateAdmin(username, elr.AdminRole(role))
	if err != nil {
		fmt.Printf("Error creating admin: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Admin created successfully! Username: %s, Role: %s, Token: %s\n", username, role, token)
}

// listAdmins 列出所有管理员
func listAdmins() {
	fmt.Println("Listing admins...")

	// 获取运行时
	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	// 列出管理员
	admins := runtime.AdminManager.ListAdmins()

	if len(admins) == 0 {
		fmt.Println("No admins found")
		return
	}

	fmt.Printf("%-20s %-15s %-15s %-10s\n", "Username", "Role", "Status", "Created At")
	fmt.Printf("%-20s %-15s %-15s %-10s\n", "--------", "----", "------", "----------")

	for _, admin := range admins {
		username := admin["username"].(string)
		role := admin["role"].(elr.AdminRole)
		status := admin["status"].(elr.AdminStatus)
		createdAt := time.Unix(admin["created_at"].(int64), 0).Format("2006-01-02")

		fmt.Printf("%-20s %-15s %-15s %-10s\n", username, role, status, createdAt)
	}

	fmt.Printf("\nTotal admins: %d\n", len(admins))
}

// addAdminPermission 为管理员添加容器权限
func addAdminPermission() {
	fmt.Println("Adding container permission to admin...")

	// 解析参数
	username := ""
	containerID := ""
	canManage := false
	canRead := false
	canWrite := false

	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--username":
			if i+1 < len(os.Args) {
				username = os.Args[i+1]
			}
		case "--container":
			if i+1 < len(os.Args) {
				containerID = os.Args[i+1]
			}
		case "--manage":
			canManage = true
		case "--read":
			canRead = true
		case "--write":
			canWrite = true
		}
	}

	if username == "" || containerID == "" {
		fmt.Println("Error: Username and container ID are required")
		os.Exit(1)
	}

	// 获取运行时
	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	// 添加权限
	if err := runtime.AdminManager.AddContainerPermission(username, containerID, canManage, canRead, canWrite); err != nil {
		fmt.Printf("Error adding permission: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Permission added successfully! Admin: %s, Container: %s\n", username, containerID)
}

// removeAdminPermission 从管理员移除容器权限
func removeAdminPermission() {
	fmt.Println("Removing container permission from admin...")

	// 解析参数
	username := ""
	containerID := ""

	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--username":
			if i+1 < len(os.Args) {
				username = os.Args[i+1]
			}
		case "--container":
			if i+1 < len(os.Args) {
				containerID = os.Args[i+1]
			}
		}
	}

	if username == "" || containerID == "" {
		fmt.Println("Error: Username and container ID are required")
		os.Exit(1)
	}

	// 获取运行时
	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	// 移除权限
	if err := runtime.AdminManager.RemoveContainerPermission(username, containerID); err != nil {
		fmt.Printf("Error removing permission: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Permission removed successfully! Admin: %s, Container: %s\n", username, containerID)
}

// settingsCommand 处理资源配置命令
func settingsCommand() {
	fmt.Println("Configuring ELR resources...")

	// 解析参数
	resourceType := ""
	modelType := ""
	directory := ""
	action := "list"

	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--resource-type":
			if i+1 < len(os.Args) {
				resourceType = os.Args[i+1]
				action = "set"
			}
		case "--model-type":
			if i+1 < len(os.Args) {
				modelType = os.Args[i+1]
				action = "set"
			}
		case "--directory":
			if i+1 < len(os.Args) {
				directory = os.Args[i+1]
			}
		case "list":
			action = "list"
		}
	}

	// 加载配置
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// 初始化资源配置
	if config.Resources.Types == nil {
		config.Resources.Types = make(map[string]struct {
			Enable bool   `yaml:"enable"`
			Dir    string `yaml:"dir"`
		})
	}

	if config.Resources.ModelTypes == nil {
		config.Resources.ModelTypes = make(map[string]struct {
			Enable bool   `yaml:"enable"`
			Dir    string `yaml:"dir"`
		})
	}

	switch action {
	case "set":
		if resourceType != "" && directory != "" {
			// 设置资源类型目录
			config.Resources.Types[resourceType] = struct {
				Enable bool   `yaml:"enable"`
				Dir    string `yaml:"dir"`
			}{
				Enable: true,
				Dir:    directory,
			}
			fmt.Printf("Resource type '%s' directory set to: %s\n", resourceType, directory)
		} else if modelType != "" && directory != "" {
			// 设置模型类型目录
			config.Resources.ModelTypes[modelType] = struct {
				Enable bool   `yaml:"enable"`
				Dir    string `yaml:"dir"`
			}{
				Enable: true,
				Dir:    directory,
			}
			fmt.Printf("Model type '%s' directory set to: %s\n", modelType, directory)
		} else {
			fmt.Println("Error: Resource type/model type and directory are required")
			fmt.Println("Usage:")
			fmt.Println("  elr Settings --resource-type <type> --directory <path>")
			fmt.Println("  elr Settings --model-type <type> --directory <path>")
			os.Exit(1)
		}

		// 保存配置
		if err := saveConfig(config); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Configuration saved successfully!")

	case "list":
		// 列出所有资源配置
		fmt.Println("Resource Configuration:")
		fmt.Println("=====================")

		// 列出资源类型
		fmt.Println("Resource Types:")
		fmt.Println("---------------")
		for rType, rConfig := range config.Resources.Types {
			fmt.Printf("%s: %s (Enabled: %v)\n", rType, rConfig.Dir, rConfig.Enable)
		}

		// 列出模型类型
		fmt.Println("\nModel Types:")
		fmt.Println("------------")
		for mType, mConfig := range config.Resources.ModelTypes {
			fmt.Printf("%s: %s (Enabled: %v)\n", mType, mConfig.Dir, mConfig.Enable)
		}

	default:
		fmt.Println("Error: Unknown action")
		fmt.Println("Usage:")
		fmt.Println("  elr Settings list - List all resource configurations")
		fmt.Println("  elr Settings --resource-type <type> --directory <path> - Set resource type directory")
		fmt.Println("  elr Settings --model-type <type> --directory <path> - Set model type directory")
		os.Exit(1)
	}
}

// uploadCommand 处理上传命令
func uploadCommand() {
	fmt.Println("Uploading resource...")

	// 解析参数
	resourceType := ""
	filePath := ""

	for i := 2; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "Settings":
			// 跳过 Settings 关键字
		case "type":
			if i+1 < len(os.Args) {
				resourceType = os.Args[i+1]
			}
		case "path:":
			if i+1 < len(os.Args) {
				filePath = os.Args[i+1]
			}
		default:
			// 处理没有关键字的参数
			if resourceType == "" {
				resourceType = os.Args[i]
			} else if filePath == "" {
				filePath = os.Args[i]
			}
		}
	}

	if resourceType == "" || filePath == "" {
		fmt.Println("Error: Resource type and file path are required")
		fmt.Println("Usage:")
		fmt.Println("  elr Upload Settings type <type> path: <path>")
		os.Exit(1)
	}

	// 加载配置
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// 初始化资源配置
	if config.Resources.Types == nil {
		config.Resources.Types = make(map[string]struct {
			Enable bool   `yaml:"enable"`
			Dir    string `yaml:"dir"`
		})
	}

	if config.Resources.ModelTypes == nil {
		config.Resources.ModelTypes = make(map[string]struct {
			Enable bool   `yaml:"enable"`
			Dir    string `yaml:"dir"`
		})
	}

	// 确定目标目录
	targetDir := ""
	if resourceConfig, exists := config.Resources.Types[resourceType]; exists {
		targetDir = resourceConfig.Dir
	} else if modelConfig, exists := config.Resources.ModelTypes[resourceType]; exists {
		targetDir = modelConfig.Dir
	} else {
		// 如果资源类型不存在，使用默认目录
		targetDir = filepath.Join(config.DataDir, "resources", resourceType)
		// 添加到配置中
		config.Resources.Types[resourceType] = struct {
			Enable bool   `yaml:"enable"`
			Dir    string `yaml:"dir"`
		}{
			Enable: true,
			Dir:    targetDir,
		}
		// 保存配置
		if err := saveConfig(config); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}
	}

	// 确保目标目录存在
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Printf("Error creating target directory: %v\n", err)
		os.Exit(1)
	}

	// 检查文件路径是本地路径还是远程URL
	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		// 远程URL，下载文件
		fmt.Printf("Downloading from URL: %s\n", filePath)
		fmt.Printf("Saving to: %s\n", targetDir)
		
		// 提取文件名
		fileName := filepath.Base(filePath)
		targetPath := filepath.Join(targetDir, fileName)
		
		// 下载文件
		data, err := downloadFileFromURL(filePath)
		if err != nil {
			fmt.Printf("Error downloading file: %v\n", err)
			os.Exit(1)
		}
		
		// 保存文件
		if err := os.WriteFile(targetPath, data, 0644); err != nil {
			fmt.Printf("Error saving file: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("File downloaded successfully: %s\n", targetPath)
	} else {
		// 本地路径，复制文件或文件夹
		fmt.Printf("Uploading from local path: %s\n", filePath)
		fmt.Printf("Saving to: %s\n", targetDir)
		
		// 检查本地路径是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Error: Local path does not exist: %s\n", filePath)
			os.Exit(1)
		}
		
		// 检查是文件还是文件夹
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Printf("Error getting file info: %v\n", err)
			os.Exit(1)
		}
		
		if fileInfo.IsDir() {
			// 复制文件夹
			destDir := filepath.Join(targetDir, fileInfo.Name())
			if err := copyDirectory(filePath, destDir); err != nil {
				fmt.Printf("Error copying directory: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Directory uploaded successfully: %s\n", destDir)
		} else {
			// 复制文件
			destFile := filepath.Join(targetDir, fileInfo.Name())
			fileSize := fileInfo.Size()
			fmt.Printf("Copying file (%d bytes)...\n", fileSize)
			
			// 打开源文件
			srcFile, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("Error opening file: %v\n", err)
				os.Exit(1)
			}
			defer srcFile.Close()
			
			// 创建目标文件
			dstFile, err := os.Create(destFile)
			if err != nil {
				fmt.Printf("Error creating file: %v\n", err)
				os.Exit(1)
			}
			defer dstFile.Close()
			
			// 复制文件并显示进度
			buffer := make([]byte, 1024*1024) // 1MB buffer
			totalCopied := int64(0)
			startTime := time.Now()
			
			for {
				n, err := srcFile.Read(buffer)
				if err != nil && err != io.EOF {
					fmt.Printf("Error reading file: %v\n", err)
					os.Exit(1)
				}
				if n == 0 {
					break
				}
				
				if _, err := dstFile.Write(buffer[:n]); err != nil {
					fmt.Printf("Error writing file: %v\n", err)
					os.Exit(1)
				}
				
				totalCopied += int64(n)
				
				// 显示进度
				if fileSize > 0 {
					percentage := float64(totalCopied) / float64(fileSize) * 100
					timeElapsed := time.Since(startTime).Seconds()
					if timeElapsed > 0 {
						speed := float64(totalCopied) / timeElapsed / 1024 / 1024 // MB/s
						fmt.Printf("Copying: %.2f%% (%.2f MB / %.2f MB) | Speed: %.2f MB/s\r", percentage, float64(totalCopied)/1024/1024, float64(fileSize)/1024/1024, speed)
					} else {
						fmt.Printf("Copying: %.2f%% (%.2f MB / %.2f MB)\r", percentage, float64(totalCopied)/1024/1024, float64(fileSize)/1024/1024)
					}
				}
			}
			
			fmt.Println() // 换行
			fmt.Printf("File uploaded successfully: %s\n", destFile)
		}
	}
}

// getDirectorySize 计算目录大小
func getDirectorySize(path string) (int64, error) {
	var size int64
	
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	
	return size, err
}

// copyDirectory 复制目录
func copyDirectory(src, dst string) error {
	// 计算目录总大小
	totalSize, err := getDirectorySize(src)
	if err != nil {
		return err
	}
	
	// 创建目标目录
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	
	// 复制文件并显示进度
	var totalCopied int64
	startTime := time.Now()
	
	err = filepath.Walk(src, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			// 创建目标子目录
			dstPath := filepath.Join(dst, strings.TrimPrefix(srcPath, src))
			return os.MkdirAll(dstPath, 0755)
		} else {
			// 复制文件
			dstPath := filepath.Join(dst, strings.TrimPrefix(srcPath, src))
			
			// 打开源文件
			srcFile, err := os.Open(srcPath)
			if err != nil {
				return err
			}
			defer srcFile.Close()
			
			// 创建目标文件
			dstFile, err := os.Create(dstPath)
			if err != nil {
				return err
			}
			defer dstFile.Close()
			
			// 复制文件内容
			buffer := make([]byte, 1024*1024) // 1MB buffer
			
			for {
				n, err := srcFile.Read(buffer)
				if err != nil && err != io.EOF {
					return err
				}
				if n == 0 {
					break
				}
				
				if _, err := dstFile.Write(buffer[:n]); err != nil {
					return err
				}
				
				totalCopied += int64(n)
				
				// 显示进度
				if totalSize > 0 {
					percentage := float64(totalCopied) / float64(totalSize) * 100
					timeElapsed := time.Since(startTime).Seconds()
					if timeElapsed > 0 {
						speed := float64(totalCopied) / timeElapsed / 1024 / 1024 // MB/s
						fmt.Printf("Copying: %.2f%% (%.2f MB / %.2f MB) | Speed: %.2f MB/s\r", percentage, float64(totalCopied)/1024/1024, float64(totalSize)/1024/1024, speed)
					} else {
						fmt.Printf("Copying: %.2f%% (%.2f MB / %.2f MB)\r", percentage, float64(totalCopied)/1024/1024, float64(totalSize)/1024/1024)
					}
				}
			}
			
			return nil
		}
	})
	
	if err == nil {
		fmt.Println() // 换行
		fmt.Printf("Directory copied successfully! Total size: %.2f MB\n", float64(totalCopied)/1024/1024)
	}
	
	return err
}
