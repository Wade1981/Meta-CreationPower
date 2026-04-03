package model

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ModelAdapter 模型适配器
type ModelAdapter struct {
	Properties *ModelProperties
}

// NewModelAdapter 创建模型适配器
func NewModelAdapter(properties *ModelProperties) *ModelAdapter {
	return &ModelAdapter{
		Properties: properties,
	}
}

// SetupEnvironment 设置环境
func (a *ModelAdapter) SetupEnvironment() error {
	// 设置环境变量
	for key, value := range a.Properties.EnvironmentVariables {
		os.Setenv(key, value)
	}

	return nil
}

// InstallDependencies 安装依赖
func (a *ModelAdapter) InstallDependencies() error {
	// 安装系统依赖
	if len(a.Properties.Dependencies.System) > 0 {
		fmt.Println("Installing system dependencies...")
		// 这里需要根据不同的操作系统执行不同的命令
		// 暂时只支持Linux
		if isLinux() {
			cmd := exec.Command("apt", append([]string{"update"}, a.Properties.Dependencies.System...)...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install system dependencies: %v", err)
			}
		} else {
			fmt.Println("System dependencies installation not supported on this platform")
		}
	}

	// 安装Python依赖
	if len(a.Properties.Dependencies.Pip) > 0 {
		fmt.Println("Installing Python dependencies...")
		pipCmd := "pip"
		if isWindows() {
			pipCmd = "pip.exe"
		}
		cmd := exec.Command(pipCmd, append([]string{"install"}, a.Properties.Dependencies.Pip...)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install Python dependencies: %v", err)
		}
	}

	return nil
}

// GetContainerConfig 获取容器配置
func (a *ModelAdapter) GetContainerConfig() map[string]interface{} {
	config := map[string]interface{}{
		"base_image":   fmt.Sprintf("python:%s-slim", a.Properties.Dependencies.Python),
		"memory_limit": a.Properties.Resources.Memory.MinRAM,
		"cpu_limit":    a.Properties.Resources.CPU.MinCores,
	}

	return config
}

// RunEntryPoint 运行入口点
func (a *ModelAdapter) RunEntryPoint(entryPoint string) error {
	if cmd, ok := a.Properties.EntryPoints[entryPoint]; ok {
		fmt.Printf("Running entry point: %s\n", cmd)
		
		// 解析命令
		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			return fmt.Errorf("invalid entry point command")
		}

		// 执行命令
		command := parts[0]
		args := parts[1:]
		execCmd := exec.Command(command, args...)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		
		if err := execCmd.Run(); err != nil {
			return fmt.Errorf("failed to run entry point: %v", err)
		}

		return nil
	}

	return fmt.Errorf("entry point %s not found", entryPoint)
}



// Run 运行模型
func (a *ModelAdapter) Run(input string) (string, error) {
	// 检查推理入口点是否存在
	inferenceCmd, ok := a.Properties.EntryPoints["inference"]
	if !ok {
		// 如果没有推理入口点，使用默认实现
		return fmt.Sprintf("Processed input: %s", input), nil
	}

	// 解析命令
	parts := strings.Fields(inferenceCmd)
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid inference command")
	}

	// 执行命令
	command := parts[0]
	args := parts[1:]
	
	// 添加输入作为参数
	args = append(args, input)
	
	execCmd := exec.Command(command, args...)
	
	// 捕获输出
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run model: %v, output: %s", err, string(output))
	}

	return string(output), nil
}
