package model

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
	"micro_model/config"
)

// ModelManager 模型管理器
type ModelManager struct {
	config *config.Config
}

// Model 模型信息
type Model struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Name        string          `json:"name"`
	Version     string          `json:"version"`
	Path        string          `json:"path"`
	DownloadURL string          `json:"download_url"`
	Properties  *ModelProperties `json:"properties,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// NewModelManager 创建模型管理器
func NewModelManager(config *config.Config) (*ModelManager, error) {
	// 确保模型目录存在
	if err := os.MkdirAll(config.Model.ModelDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create model directory: %v", err)
	}

	return &ModelManager{
		config: config,
	}, nil
}

// GetModel 获取模型信息
func (m *ModelManager) GetModel(modelID string) (*Model, error) {
	modelPath := filepath.Join(m.config.Model.ModelDir, modelID)

	// 检查模型目录是否存在
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("model %s not found", modelID)
	}

	// 加载模型属性
	properties, err := LoadModelProperties(modelPath)
	if err != nil {
		fmt.Printf("Warning: Failed to load model properties: %v\n", err)
		// 继续执行，使用默认值
	}

	// 从模型属性中获取信息，或使用默认值
	modelType := "unknown"
	modelName := modelID
	modelVersion := "1.0.0"

	if properties != nil {
		if properties.Type != "" {
			modelType = properties.Type
		}
		if properties.ModelName != "" {
			modelName = properties.ModelName
		}
		if properties.Version != "" {
			modelVersion = properties.Version
		}
	}

	// 返回模型信息
	return &Model{
		ID:          modelID,
		Type:        modelType,
		Name:        modelName,
		Version:     modelVersion,
		Path:        modelPath,
		Properties:  properties,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// ListModels 列出所有模型
func (m *ModelManager) ListModels() ([]*Model, error) {
	var models []*Model

	// 读取模型目录
	entries, err := os.ReadDir(m.config.Model.ModelDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read model directory: %v", err)
	}

	// 遍历所有模型目录
	for _, entry := range entries {
		if entry.IsDir() {
			model, err := m.GetModel(entry.Name())
			if err == nil {
				models = append(models, model)
			}
		}
	}

	return models, nil
}

// DownloadModel 下载模型
func (m *ModelManager) DownloadModel(modelID string, modelType string, downloadURL string) error {
	// Create channel to receive download result
	errCh := make(chan error)

	// Start download in a goroutine
	go func() {
		modelPath := filepath.Join(m.config.Model.ModelDir, modelID)

		// 创建模型目录
		if err := os.MkdirAll(modelPath, 0755); err != nil {
			errCh <- fmt.Errorf("failed to create model directory: %v", err)
			return
		}

		// 下载模型文件
		// 这里只是一个简单的示例，实际下载逻辑会更复杂
		// 例如，对于大型模型，可能需要分块下载、校验等
		fmt.Printf("Downloading model %s from %s...\n", modelID, downloadURL)

		// 模拟下载过程
		time.Sleep(2 * time.Second)

		// 创建一个简单的模型文件作为示例
		modelFile := filepath.Join(modelPath, "model.txt")
		if err := os.WriteFile(modelFile, []byte(fmt.Sprintf("Model: %s\nType: %s\nDownloaded from: %s\n", modelID, modelType, downloadURL)), 0644); err != nil {
			errCh <- fmt.Errorf("failed to create model file: %v", err)
			return
		}

		fmt.Printf("Model %s downloaded successfully\n", modelID)
		errCh <- nil
	}()

	// Wait for download to complete
	return <-errCh
}

// DeleteModel 删除模型
func (m *ModelManager) DeleteModel(modelID string) error {
	modelPath := filepath.Join(m.config.Model.ModelDir, modelID)

	// 检查模型目录是否存在
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return fmt.Errorf("model %s not found", modelID)
	}

	// 删除模型目录
	if err := os.RemoveAll(modelPath); err != nil {
		return fmt.Errorf("failed to delete model: %v", err)
	}

	fmt.Printf("Model %s deleted successfully\n", modelID)
	return nil
}

// UpdateModel 更新模型
func (m *ModelManager) UpdateModel(modelID string, downloadURL string) error {
	// 先删除旧模型
	if err := m.DeleteModel(modelID); err != nil {
		return err
	}

	// 重新下载模型
	return m.DownloadModel(modelID, "unknown", downloadURL)
}

// Exists 检查模型是否存在
func (m *ModelManager) Exists(modelID string) bool {
	modelPath := filepath.Join(m.config.Model.ModelDir, modelID)
	_, err := os.Stat(modelPath)
	return !os.IsNotExist(err)
}

// GetModelAdapter 获取模型适配器
func (m *ModelManager) GetModelAdapter(modelID string) (*ModelAdapter, error) {
	// 获取模型信息
	model, err := m.GetModel(modelID)
	if err != nil {
		return nil, err
	}

	// 检查模型属性是否存在
	if model.Properties == nil {
		return nil, fmt.Errorf("model properties not found for model %s", modelID)
	}

	// 创建模型适配器
	adapter := NewModelAdapter(model.Properties)

	return adapter, nil
}

// InstallModelDependencies 安装模型依赖
func (m *ModelManager) InstallModelDependencies(modelID string, depType string) error {
	// 获取模型适配器
	adapter, err := m.GetModelAdapter(modelID)
	if err != nil {
		return err
	}

	// 安装依赖
	switch depType {
	case "pip":
		// 安装Python依赖
		if len(adapter.Properties.Dependencies.Pip) > 0 {
			fmt.Println("Installing Python dependencies...")
			pipCmd := "pip"
			if isWindows() {
				pipCmd = "pip.exe"
			}
			cmd := exec.Command(pipCmd, append([]string{"install"}, adapter.Properties.Dependencies.Pip...)...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install Python dependencies: %v", err)
			}
		}
	case "system":
		// 安装系统依赖
		if len(adapter.Properties.Dependencies.System) > 0 {
			fmt.Println("Installing system dependencies...")
			// 这里需要根据不同的操作系统执行不同的命令
			// 暂时只支持Linux
			if isLinux() {
				cmd := exec.Command("apt", append([]string{"update"}, adapter.Properties.Dependencies.System...)...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("failed to install system dependencies: %v", err)
				}
			} else {
				fmt.Println("System dependencies installation not supported on this platform")
			}
		}
	default:
		return fmt.Errorf("unknown dependency type: %s", depType)
	}

	return nil
}

// isLinux 检查是否为Linux系统
func isLinux() bool {
	return runtime.GOOS == "linux"
}

// isWindows 检查是否为Windows系统
func isWindows() bool {
	return runtime.GOOS == "windows"
}
