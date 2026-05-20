package model

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
	"micro_model/config"
	"gopkg.in/yaml.v3"
)

// ELRConfig ELR 配置结构体
type ELRConfig struct {
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

// ModelManager 模型管理器
type ModelManager struct {
	config       *config.Config
	loadedModels map[string]*LoadedModel
	modelMutex   sync.RWMutex
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

// LoadedModel 已加载的模型
type LoadedModel struct {
	Model        *Model
	Adapter      *ModelAdapter
	LoadedAt     time.Time
	LastUsedAt   time.Time
	ResourceUsage map[string]interface{}
	IsActive     bool
}

// loadConfig 加载 ELR 配置
func loadConfig() (*ELRConfig, error) {
	configPath := os.Getenv("ELR_CONFIG")
	if configPath == "" {
		configPath = "~/.elr/config.yaml"
	}

	// 扩展 ~ 为 home 目录
	if len(configPath) > 0 && configPath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %v", err)
		}
		configPath = homeDir + configPath[1:]
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 返回默认配置
		return defaultConfig(), nil
	}

	// 读取配置文件
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// 解析配置
	config := &ELRConfig{}
	if err := yaml.Unmarshal(configBytes, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return config, nil
}

// defaultConfig 返回默认配置
func defaultConfig() *ELRConfig {
	return &ELRConfig{
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

// NewModelManager 创建模型管理器
func NewModelManager(config *config.Config) (*ModelManager, error) {
	// 确保模型目录存在
	if err := os.MkdirAll(config.Model.ModelDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create model directory: %v", err)
	}

	return &ModelManager{
		config:       config,
		loadedModels: make(map[string]*LoadedModel),
	}, nil
}

// GetModel 获取模型信息
func (m *ModelManager) GetModel(modelID string) (*Model, error) {
	// 首先尝试在当前模型目录中查找
	modelPath := filepath.Join(m.config.Model.ModelDir, modelID)
	if _, err := os.Stat(modelPath); err == nil {
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

	// 如果在当前模型目录中没有找到，尝试从 ELR Settings 中获取所有模型类型目录
	elrConfig, err := loadConfig()
	if err == nil {
		// 尝试在所有模型类型目录中查找
		for modelType, modelConfig := range elrConfig.Resources.ModelTypes {
			if modelConfig.Enable && modelConfig.Dir != "" {
				modelPath := filepath.Join(modelConfig.Dir, modelID)
				if _, err := os.Stat(modelPath); err == nil {
					// 加载模型属性
					properties, err := LoadModelProperties(modelPath)
					if err != nil {
						fmt.Printf("Warning: Failed to load model properties: %v\n", err)
						// 继续执行，使用默认值
					}

					// 从模型属性中获取信息，或使用默认值
					modelName := modelID
					modelVersion := "1.0.0"

					if properties != nil {
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
			}
		}

		// 尝试在所有资源类型目录中查找
		for _, resourceConfig := range elrConfig.Resources.Types {
			if resourceConfig.Enable && resourceConfig.Dir != "" {
				modelPath := filepath.Join(resourceConfig.Dir, modelID)
				if _, err := os.Stat(modelPath); err == nil {
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
			}
		}
	}

	// 如果都没有找到，返回错误
	return nil, fmt.Errorf("model %s not found", modelID)
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
	// 检查模型是否正在运行
	m.modelMutex.Lock()
	if _, exists := m.loadedModels[modelID]; exists {
		// 模型正在运行，先卸载
		delete(m.loadedModels, modelID)
		fmt.Printf("Stopped running model: %s\n", modelID)
	}
	m.modelMutex.Unlock()

	// 检查模型目录是否存在
	modelPath := filepath.Join(m.config.Model.ModelDir, modelID)
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
	// 首先尝试在当前模型目录中查找
	modelPath := filepath.Join(m.config.Model.ModelDir, modelID)
	if _, err := os.Stat(modelPath); err == nil {
		return true
	}

	// 如果在当前模型目录中没有找到，尝试从 ELR Settings 中获取所有模型类型目录
	elrConfig, err := loadConfig()
	if err == nil {
		// 尝试在所有模型类型目录中查找
		for _, modelConfig := range elrConfig.Resources.ModelTypes {
			if modelConfig.Enable && modelConfig.Dir != "" {
				modelPath := filepath.Join(modelConfig.Dir, modelID)
				if _, err := os.Stat(modelPath); err == nil {
					return true
				}
			}
		}

		// 尝试在所有资源类型目录中查找
		for _, resourceConfig := range elrConfig.Resources.Types {
			if resourceConfig.Enable && resourceConfig.Dir != "" {
				modelPath := filepath.Join(resourceConfig.Dir, modelID)
				if _, err := os.Stat(modelPath); err == nil {
					return true
				}
			}
		}
	}

	// 如果都没有找到，返回 false
	return false
}

// LoadModel 加载模型
func (m *ModelManager) LoadModel(modelID string) error {
	m.modelMutex.Lock()
	defer m.modelMutex.Unlock()

	// 检查模型是否已加载
	if _, exists := m.loadedModels[modelID]; exists {
		fmt.Printf("Model %s is already loaded\n", modelID)
		return nil
	}

	// 获取模型信息
	model, err := m.GetModel(modelID)
	if err != nil {
		return err
	}

	// 获取模型适配器
	adapter, err := m.GetModelAdapter(modelID)
	if err != nil {
		return err
	}

	// 创建已加载模型
	loadedModel := &LoadedModel{
		Model:        model,
		Adapter:      adapter,
		LoadedAt:     time.Now(),
		LastUsedAt:   time.Now(),
		ResourceUsage: make(map[string]interface{}),
		IsActive:     false,
	}

	// 模拟模型加载过程
	fmt.Printf("Loading model %s...\n", modelID)
	time.Sleep(1 * time.Second)

	// 记录资源使用情况
	loadedModel.ResourceUsage["memory"] = 256 // 假设使用256MB内存
	loadedModel.ResourceUsage["cpu"] = 10     // 假设使用10% CPU

	// 添加到已加载模型列表
	m.loadedModels[modelID] = loadedModel

	fmt.Printf("Model %s loaded successfully\n", modelID)
	return nil
}

// UnloadModel 卸载模型
func (m *ModelManager) UnloadModel(modelID string) error {
	m.modelMutex.Lock()
	defer m.modelMutex.Unlock()

	// 检查模型是否已加载
	loadedModel, exists := m.loadedModels[modelID]
	if !exists {
		return fmt.Errorf("model %s is not loaded", modelID)
	}

	// 检查模型是否为活动状态
	if loadedModel.IsActive {
		return fmt.Errorf("cannot unload active model %s", modelID)
	}

	// 模拟模型卸载过程
	fmt.Printf("Unloading model %s...\n", modelID)
	time.Sleep(500 * time.Millisecond)

	// 从已加载模型列表中移除
	delete(m.loadedModels, modelID)

	fmt.Printf("Model %s unloaded successfully\n", modelID)
	return nil
}

// SwitchModel 切换模型
func (m *ModelManager) SwitchModel(modelID string) error {
	m.modelMutex.Lock()
	defer m.modelMutex.Unlock()

	// 检查模型是否存在
	if !m.Exists(modelID) {
		return fmt.Errorf("model %s not found", modelID)
	}

	// 检查模型是否已加载，如果没有则加载
	if _, exists := m.loadedModels[modelID]; !exists {
		if err := m.LoadModel(modelID); err != nil {
			return err
		}
	}

	// 先将所有模型设置为非活动状态
	for id, model := range m.loadedModels {
		model.IsActive = (id == modelID)
	}

	// 更新活动模型的最后使用时间
	if loadedModel, exists := m.loadedModels[modelID]; exists {
		loadedModel.LastUsedAt = time.Now()
		loadedModel.IsActive = true
	}

	fmt.Printf("Switched to model %s\n", modelID)
	return nil
}

// GetLoadedModels 获取已加载的模型
func (m *ModelManager) GetLoadedModels() map[string]*LoadedModel {
	m.modelMutex.RLock()
	defer m.modelMutex.RUnlock()

	// 创建副本以避免并发修改
	loadedModels := make(map[string]*LoadedModel)
	for id, model := range m.loadedModels {
		loadedModels[id] = model
	}

	return loadedModels
}

// IsModelLoaded 检查模型是否已加载
func (m *ModelManager) IsModelLoaded(modelID string) bool {
	m.modelMutex.RLock()
	defer m.modelMutex.RUnlock()

	_, exists := m.loadedModels[modelID]
	return exists
}

// GetActiveModel 获取当前活动的模型
func (m *ModelManager) GetActiveModel() *LoadedModel {
	m.modelMutex.RLock()
	defer m.modelMutex.RUnlock()

	for _, model := range m.loadedModels {
		if model.IsActive {
			return model
		}
	}

	return nil
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
	adapter := NewModelAdapter(model.Properties, model)

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
