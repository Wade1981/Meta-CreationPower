package model

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ModelAdapter 模型适配器
type ModelAdapter struct {
	Properties *ModelProperties
	Model      *Model
}

// NewModelAdapter 创建模型适配器
func NewModelAdapter(properties *ModelProperties, model *Model) *ModelAdapter {
	return &ModelAdapter{
		Properties: properties,
		Model:      model,
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
		// 尝试查找模型目录中的Python脚本
		var modelScript string
		var modelDir string
		
		if a.Model != nil && a.Model.ID != "" {
			// 使用模型ID作为脚本名称
			modelScript = a.Model.ID + ".py"
			modelDir = a.Model.ID
		} else {
			// 使用模型名称作为脚本名称
			modelName := a.Properties.ModelName
			modelScript = modelName + ".py"
			modelDir = modelName
		}
		
		// 检查模型路径是否存在
		if a.Model != nil && a.Model.Path != "" {
			// 检查模型路径中是否存在模型脚本
			modelScriptPath := filepath.Join(a.Model.Path, modelScript)
			if _, err := os.Stat(modelScriptPath); err == nil {
				// 执行Python脚本的predict方法，而不是直接执行脚本
				pythonCode := fmt.Sprintf(`
import sys
import os

# 设置编码为UTF-8
sys.stdout.reconfigure(encoding='utf-8')
sys.stderr.reconfigure(encoding='utf-8')

# 隐藏模型初始化的输出
sys.stdout = open(os.devnull, 'w', encoding='utf-8')

# 添加模型目录到Python路径
sys.path.insert(0, r"%s")

# 导入模型模块
from %s import ELRChatModel

# 初始化模型
model = ELRChatModel()

# 恢复标准输出
sys.stdout = sys.__stdout__

# 执行预测
input_text = sys.argv[1]
response = model.predict(input_text)
print(response)
`, a.Model.Path, strings.TrimSuffix(modelScript, ".py"))
				cmd := exec.Command("python", "-c", pythonCode, input)
				// 设置环境变量，确保Python使用UTF-8编码
				cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
				output, err := cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("failed to run Python model: %v, output: %s", err, string(output))
				}
				return string(output), nil
			}
		}
		
		// 检查当前目录是否存在模型脚本
		if _, err := os.Stat(modelScript); err == nil {
			// 执行Python脚本的predict方法
			pythonCode := fmt.Sprintf(`
import sys
import os

# 设置编码为UTF-8
sys.stdout.reconfigure(encoding='utf-8')
sys.stderr.reconfigure(encoding='utf-8')

# 隐藏模型初始化的输出
sys.stdout = open(os.devnull, 'w', encoding='utf-8')

from %s import ELRChatModel

model = ELRChatModel()

# 恢复标准输出
sys.stdout = sys.__stdout__

input_text = sys.argv[1]
response = model.predict(input_text)
print(response)
`, strings.TrimSuffix(modelScript, ".py"))
			cmd := exec.Command("python", "-c", pythonCode, input)
			// 设置环境变量，确保Python使用UTF-8编码
			cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
			output, err := cmd.CombinedOutput()
			if err != nil {
				return "", fmt.Errorf("failed to run Python model: %v, output: %s", err, string(output))
			}
			return string(output), nil
		}
		
		// 检查模型目录是否存在模型脚本
		if _, err := os.Stat(modelDir); err == nil {
			modelScriptPath := filepath.Join(modelDir, modelScript)
			if _, err := os.Stat(modelScriptPath); err == nil {
				// 执行Python脚本的predict方法
				pythonCode := fmt.Sprintf(`
import sys
import os

# 设置编码为UTF-8
sys.stdout.reconfigure(encoding='utf-8')
sys.stderr.reconfigure(encoding='utf-8')

# 隐藏模型初始化的输出
sys.stdout = open(os.devnull, 'w', encoding='utf-8')

sys.path.insert(0, r"%s")
from %s import ELRChatModel

model = ELRChatModel()

# 恢复标准输出
sys.stdout = sys.__stdout__

input_text = sys.argv[1]
response = model.predict(input_text)
print(response)
`, modelDir, strings.TrimSuffix(modelScript, ".py"))
				cmd := exec.Command("python", "-c", pythonCode, input)
				// 设置环境变量，确保Python使用UTF-8编码
				cmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
				output, err := cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("failed to run Python model: %v, output: %s", err, string(output))
				}
				return string(output), nil
			}
		}
		
		// 如果没有推理入口点且没有找到模型脚本，使用默认实现
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
	// 设置环境变量，确保命令使用UTF-8编码
	execCmd.Env = append(os.Environ(), "PYTHONIOENCODING=utf-8")
	
	// 捕获输出
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run model: %v, output: %s", err, string(output))
	}

	return string(output), nil
}
