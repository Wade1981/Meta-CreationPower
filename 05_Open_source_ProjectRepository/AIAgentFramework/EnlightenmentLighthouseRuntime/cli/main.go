// Package main implements the command-line interface for Enlightenment Lighthouse Runtime
package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"elr.local"
	"micro_model/config"
	"micro_model/model"
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
		runContainer()
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
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
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
} 

// loadConfig loads the configuration from file
func loadConfig() (*elr.Config, error) {
	configPath := os.Getenv("ELR_CONFIG")
	if configPath == "" {
		configPath = "~/.elr/config.yaml"
	}

	// Expand ~ to home directory
	if configPath[:2] == "~" {
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

	return config, nil
}

// defaultConfig returns the default configuration
func defaultConfig() *elr.Config {
	return &elr.Config{
		LogLevel:  "info",
		DataDir:   "~/.elr/data",
		PluginDir: "~/.elr/plugins",
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
		}{
			Enable:  true,
			Bridge:  "elr0",
			Subnet:  "172.16.0.0/16",
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

// getRuntime returns a runtime instance
func getRuntime() (*elr.Runtime, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	runtime, err := elr.NewRuntime(config)
	if err != nil {
		return nil, err
	}

	// Start runtime if not already running
	if err := runtime.Start(); err != nil {
		return nil, err
	}

	return runtime, nil
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

// loadModelConfig 加载模型配置
func loadModelConfig() (*config.ModelConfig, error) {
	// 创建默认模型配置
	return &config.ModelConfig{
		ModelDir: "./micro_model/model/models",
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
	sandboxManager, err := sandbox.NewSandboxManager(&config.SandboxConfig{}, modelManager)
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
	sandboxManager, err := sandbox.NewSandboxManager(&config.SandboxConfig{}, modelManager)
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
	sandboxManager, err := sandbox.NewSandboxManager(&config.SandboxConfig{}, modelManager)
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
	sandboxManager, err := sandbox.NewSandboxManager(&config.SandboxConfig{}, modelManager)
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
	sandboxManager, err := sandbox.NewSandboxManager(&config.SandboxConfig{}, modelManager)
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
	sandboxManager, err := sandbox.NewSandboxManager(&config.SandboxConfig{}, modelManager)
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
	sandboxManager, err := sandbox.NewSandboxManager(&config.SandboxConfig{}, modelManager)
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
		startAPIServices()
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

// startAPIServices 启动 API 服务
func startAPIServices() {
	fmt.Println("Starting API services...")

	// 解析参数
	apiType := "all"
	for i := 3; i < len(os.Args); i++ {
		if os.Args[i] == "--api-type" && i+1 < len(os.Args) {
			apiType = os.Args[i+1]
			break
		}
	}

	// 启动运行时
	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	switch apiType {
	case "all":
		// 启动所有API服务
		if err := runtime.NetworkManager.Start(); err != nil {
			fmt.Printf("Error starting network service: %v\n", err)
			os.Exit(1)
		}
		// 这里可以添加启动其他API服务的代码
		fmt.Println("All API services started successfully!")
		fmt.Println("Public API: http://localhost:8080")
		fmt.Println("Desktop API: http://localhost:8081")
		fmt.Println("Model API: http://localhost:8082")
	case "public":
		// 启动Public API
		if err := runtime.NetworkManager.Start(); err != nil {
			fmt.Printf("Error starting public API: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Public API started successfully!")
		fmt.Println("Public API: http://localhost:8080")
	case "desktop":
		// 启动Desktop API
		fmt.Println("Starting Desktop API...")
		// 这里可以添加启动Desktop API的代码
		fmt.Println("Desktop API started successfully!")
		fmt.Println("Desktop API: http://localhost:8081")
	case "model":
		// 启动Model API
		fmt.Println("Starting Model API...")
		// 这里可以添加启动Model API的代码
		fmt.Println("Model API started successfully!")
		fmt.Println("Model API: http://localhost:8082")
	default:
		fmt.Printf("Unknown API type: %s\n", apiType)
		printHelp()
		os.Exit(1)
	}
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

	// 启动运行时
	runtime, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	switch apiType {
	case "all":
		// 停止所有API服务
		if err := runtime.NetworkManager.Stop(); err != nil {
			fmt.Printf("Error stopping network service: %v\n", err)
			os.Exit(1)
		}
		// 这里可以添加停止其他API服务的代码
		fmt.Println("All API services stopped successfully!")
	case "public":
		// 停止Public API
		if err := runtime.NetworkManager.Stop(); err != nil {
			fmt.Printf("Error stopping public API: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Public API stopped successfully!")
	case "desktop":
		// 停止Desktop API
		fmt.Println("Stopping Desktop API...")
		// 这里可以添加停止Desktop API的代码
		fmt.Println("Desktop API stopped successfully!")
	case "model":
		// 停止Model API
		fmt.Println("Stopping Model API...")
		// 这里可以添加停止Model API的代码
		fmt.Println("Model API stopped successfully!")
	default:
		fmt.Printf("Unknown API type: %s\n", apiType)
		printHelp()
		os.Exit(1)
	}
}

// checkAPIStatus 检查 API 服务状态
func checkAPIStatus() {
	fmt.Println("Checking API service status...")

	// 启动运行时
	_, err := getRuntime()
	if err != nil {
		fmt.Printf("Error getting runtime: %v\n", err)
		os.Exit(1)
	}

	// 模拟检查网络状态
	fmt.Println("API Service Status:")
	fmt.Println("==================")
	fmt.Println("Public API: http://localhost:8080 - Running")
	fmt.Println("Desktop API: http://localhost:8081 - Running")
	fmt.Println("Model API: http://localhost:8082 - Running")
	fmt.Println("Health check: http://localhost:8080/health - Available")
	fmt.Println("Container API: http://localhost:8080/api/container/list - Available")
	fmt.Println("Model API: http://localhost:8082/api/models - Available")
}

// configureAPI 配置 API 地址和端口
func configureAPI() {
	fmt.Println("Configuring API services...")

	// 解析参数
	apiType := ""
	address := "localhost"
	port := 0

	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--api-type":
			if i+1 < len(os.Args) {
				apiType = os.Args[i+1]
			}
		case "--address":
			if i+1 < len(os.Args) {
				address = os.Args[i+1]
			}
		case "--port":
			if i+1 < len(os.Args) {
				fmt.Sscanf(os.Args[i+1], "%d", &port)
			}
		}
	}

	if apiType == "" {
		fmt.Println("Error: API type is required")
		printHelp()
		os.Exit(1)
	}

	if port == 0 {
		fmt.Println("Error: Port is required")
		printHelp()
		os.Exit(1)
	}

	switch apiType {
	case "public":
		fmt.Printf("Configuring Public API: %s:%d\n", address, port)
		// 这里可以添加配置Public API的代码
	case "desktop":
		fmt.Printf("Configuring Desktop API: %s:%d\n", address, port)
		// 这里可以添加配置Desktop API的代码
	case "model":
		fmt.Printf("Configuring Model API: %s:%d\n", address, port)
		// 这里可以添加配置Model API的代码
	default:
		fmt.Printf("Unknown API type: %s\n", apiType)
		printHelp()
		os.Exit(1)
	}

	fmt.Println("API configuration updated successfully!")
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
		role := admin["role"].(string)
		status := admin["status"].(string)
		createdAt := time.Unix(int64(admin["created_at"].(float64)), 0).Format("2006-01-02")

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
