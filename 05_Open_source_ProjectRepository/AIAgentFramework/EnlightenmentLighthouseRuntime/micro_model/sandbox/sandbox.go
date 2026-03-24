package sandbox

import (
	"fmt"
	"sync"
	"time"

	"micro_model/config"
	"micro_model/model"
)

// Sandbox 沙箱实例
type Sandbox struct {
	ID           string            `json:"id"`
	Container    string            `json:"container"`
	Status       string            `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	StartedAt    *time.Time        `json:"started_at,omitempty"`
	StoppedAt    *time.Time        `json:"stopped_at,omitempty"`
	Models       map[string]*Model `json:"models"` // 模型ID -> 模型信息
	Resources    Resources         `json:"resources"`
}

// Model 沙箱中的模型信息
type Model struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // running, stopped, error
	Resources   Resources `json:"resources"`
	LoadedAt    time.Time `json:"loaded_at"`
}

// RuntimeStatus 运行时状态
type RuntimeStatus struct {
	Status    string    `json:"status"`
	Container string    `json:"container"`
	ModelID   string    `json:"model_id"`
	StartedAt time.Time `json:"started_at"`
	Uptime    string    `json:"uptime"`
	Resources Resources  `json:"resources"`
}

// Resources 资源使用情况
type Resources struct {
	CPU    float64 `json:"cpu"`
	Memory int64   `json:"memory"`
	Disk   int64   `json:"disk"`
}

// SandboxManager 沙箱管理器
type SandboxManager struct {
	config   *config.SandboxConfig
	sandboxes map[string]*Sandbox // 沙箱ID -> 沙箱实例
	modelManager *model.ModelManager
	mutex    sync.RWMutex
}

// NewSandboxManager 创建沙箱管理器
func NewSandboxManager(config *config.SandboxConfig, modelManager *model.ModelManager) (*SandboxManager, error) {
	return &SandboxManager{
		config:   config,
		sandboxes: make(map[string]*Sandbox),
		modelManager: modelManager,
	}, nil
}

// CreateSandbox 创建新沙箱
func (sm *SandboxManager) CreateSandbox(containerName string) (*Sandbox, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// 生成沙箱ID
	sandboxID := fmt.Sprintf("sandbox-%d", time.Now().UnixNano()/1000000)

	// 创建沙箱实例
	sandbox := &Sandbox{
		ID:        sandboxID,
		Container: containerName,
		Status:    "created",
		CreatedAt: time.Now(),
		Models:    make(map[string]*Model),
		Resources: Resources{
			CPU:    0,
			Memory: 0,
			Disk:   0,
		},
	}

	// 保存沙箱实例
	sm.sandboxes[sandboxID] = sandbox

	fmt.Printf("Created sandbox: %s in container: %s\n", sandboxID, containerName)

	return sandbox, nil
}

// GetSandbox 获取沙箱信息
func (sm *SandboxManager) GetSandbox(sandboxID string) (*Sandbox, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	sandbox, exists := sm.sandboxes[sandboxID]
	if !exists {
		return nil, fmt.Errorf("sandbox %s not found", sandboxID)
	}

	return sandbox, nil
}

// ListSandboxes 列出所有沙箱
func (sm *SandboxManager) ListSandboxes() []*Sandbox {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	sandboxes := make([]*Sandbox, 0, len(sm.sandboxes))
	for _, sandbox := range sm.sandboxes {
		sandboxes = append(sandboxes, sandbox)
	}

	return sandboxes
}

// StartSandbox 启动沙箱
func (sm *SandboxManager) StartSandbox(sandboxID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sandbox, exists := sm.sandboxes[sandboxID]
	if !exists {
		return fmt.Errorf("sandbox %s not found", sandboxID)
	}

	if sandbox.Status == "running" {
		return fmt.Errorf("sandbox %s is already running", sandboxID)
	}

	// 启动沙箱
	startTime := time.Now()
	sandbox.StartedAt = &startTime
	sandbox.Status = "running"

	fmt.Printf("Started sandbox: %s\n", sandboxID)

	return nil
}

// StopSandbox 停止沙箱
func (sm *SandboxManager) StopSandbox(sandboxID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sandbox, exists := sm.sandboxes[sandboxID]
	if !exists {
		return fmt.Errorf("sandbox %s not found", sandboxID)
	}

	if sandbox.Status != "running" {
		return fmt.Errorf("sandbox %s is not running", sandboxID)
	}

	// 停止沙箱中的所有模型
	for modelID := range sandbox.Models {
		if err := sm.UnloadModel(sandboxID, modelID); err != nil {
			fmt.Printf("Error unloading model %s: %v\n", modelID, err)
		}
	}

	// 停止沙箱
	stopTime := time.Now()
	sandbox.StoppedAt = &stopTime
	sandbox.Status = "stopped"

	fmt.Printf("Stopped sandbox: %s\n", sandboxID)

	return nil
}

// DeleteSandbox 删除沙箱
func (sm *SandboxManager) DeleteSandbox(sandboxID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// 检查沙箱是否存在
	_, exists := sm.sandboxes[sandboxID]
	if !exists {
		return fmt.Errorf("sandbox %s not found", sandboxID)
	}

	// 如果沙箱正在运行，先停止
	if err := sm.StopSandbox(sandboxID); err != nil && err.Error() != fmt.Sprintf("sandbox %s is not running", sandboxID) {
		return err
	}

	// 删除沙箱
	delete(sm.sandboxes, sandboxID)

	fmt.Printf("Deleted sandbox: %s\n", sandboxID)

	return nil
}

// LoadModel 加载模型到沙箱
func (sm *SandboxManager) LoadModel(sandboxID string, modelID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// 检查沙箱是否存在
	sandbox, exists := sm.sandboxes[sandboxID]
	if !exists {
		return fmt.Errorf("sandbox %s not found", sandboxID)
	}

	// 检查沙箱是否运行中
	if sandbox.Status != "running" {
		return fmt.Errorf("sandbox %s is not running", sandboxID)
	}

	// 检查模型是否已加载
	if _, exists := sandbox.Models[modelID]; exists {
		return fmt.Errorf("model %s is already loaded in sandbox %s", modelID, sandboxID)
	}

	// 检查模型是否存在
	if !sm.modelManager.Exists(modelID) {
		return fmt.Errorf("model %s not found", modelID)
	}

	// 获取模型信息
	modelInfo, err := sm.modelManager.GetModel(modelID)
	if err != nil {
		return fmt.Errorf("failed to get model info: %v", err)
	}

	// 创建模型实例
	model := &Model{
		ID:          modelID,
		Name:        modelInfo.Name,
		Description: modelInfo.Properties.Description,
		Status:      "running",
		Resources: Resources{
			CPU:    10.0,  // 模拟值
			Memory: 256 * 1024 * 1024, // 256MB
			Disk:   512 * 1024 * 1024, // 512MB
		},
		LoadedAt: time.Now(),
	}

	// 加载模型到沙箱
	sandbox.Models[modelID] = model

	// 更新沙箱资源使用情况
	sandbox.Resources.CPU += model.Resources.CPU
	sandbox.Resources.Memory += model.Resources.Memory
	sandbox.Resources.Disk += model.Resources.Disk

	fmt.Printf("Loaded model %s into sandbox %s\n", modelID, sandboxID)

	return nil
}

// UnloadModel 从沙箱卸载模型
func (sm *SandboxManager) UnloadModel(sandboxID string, modelID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// 检查沙箱是否存在
	sandbox, exists := sm.sandboxes[sandboxID]
	if !exists {
		return fmt.Errorf("sandbox %s not found", sandboxID)
	}

	// 检查模型是否已加载
	model, exists := sandbox.Models[modelID]
	if !exists {
		return fmt.Errorf("model %s is not loaded in sandbox %s", modelID, sandboxID)
	}

	// 从沙箱卸载模型
	delete(sandbox.Models, modelID)

	// 更新沙箱资源使用情况
	sandbox.Resources.CPU -= model.Resources.CPU
	sandbox.Resources.Memory -= model.Resources.Memory
	sandbox.Resources.Disk -= model.Resources.Disk

	fmt.Printf("Unloaded model %s from sandbox %s\n", modelID, sandboxID)

	return nil
}

// GetModelCount 获取沙箱中模型数量
func (sm *SandboxManager) GetModelCount(sandboxID string) (int, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// 检查沙箱是否存在
	sandbox, exists := sm.sandboxes[sandboxID]
	if !exists {
		return 0, fmt.Errorf("sandbox %s not found", sandboxID)
	}

	return len(sandbox.Models), nil
}

// GetModels 获取沙箱中的模型列表
func (sm *SandboxManager) GetModels(sandboxID string) ([]*Model, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	// 检查沙箱是否存在
	sandbox, exists := sm.sandboxes[sandboxID]
	if !exists {
		return nil, fmt.Errorf("sandbox %s not found", sandboxID)
	}

	models := make([]*Model, 0, len(sandbox.Models))
	for _, model := range sandbox.Models {
		models = append(models, model)
	}

	return models, nil
}

// RunModel 运行模型
func (sm *SandboxManager) RunModel(sandboxID string, modelID string, input string) (string, error) {
	sm.mutex.RLock()
	// 检查沙箱是否存在
	sandbox, exists := sm.sandboxes[sandboxID]
	if !exists {
		sm.mutex.RUnlock()
		return "", fmt.Errorf("sandbox %s not found", sandboxID)
	}

	// 检查模型是否已加载
	model, exists := sandbox.Models[modelID]
	if !exists {
		sm.mutex.RUnlock()
		return "", fmt.Errorf("model %s is not loaded in sandbox %s", modelID, sandboxID)
	}
	sm.mutex.RUnlock()

	// 检查模型状态
	if model.Status != "running" {
		return "", fmt.Errorf("model %s is not running", modelID)
	}

	// 获取模型适配器
	adapter, err := sm.modelManager.GetModelAdapter(modelID)
	if err != nil {
		return "", fmt.Errorf("failed to get model adapter: %v", err)
	}

	// 运行模型
	output, err := adapter.Run(input)
	if err != nil {
		return "", fmt.Errorf("failed to run model: %v", err)
	}

	fmt.Printf("Running model %s in sandbox %s with input: %s\n", modelID, sandboxID, input)
	fmt.Printf("Model output: %s\n", output)

	return output, nil
}

// GetRuntimeStatus 获取运行时状态
func (sm *SandboxManager) GetRuntimeStatus(sandboxID string) (*RuntimeStatus, error) {
	sm.mutex.RLock()
	// 检查沙箱是否存在
	sandbox, exists := sm.sandboxes[sandboxID]
	if !exists {
		sm.mutex.RUnlock()
		return nil, fmt.Errorf("sandbox %s not found", sandboxID)
	}
	sm.mutex.RUnlock()

	// 构建运行时状态
	var uptime string
	if sandbox.StartedAt != nil {
		duration := time.Since(*sandbox.StartedAt)
		uptime = fmt.Sprintf("%vm%vs", int(duration.Minutes()), int(duration.Seconds())%60)
	} else {
		uptime = "0s"
	}

	// 获取第一个模型ID作为示例
	var modelID string
	sm.mutex.RLock()
	for id := range sandbox.Models {
		modelID = id
		break
	}
	sm.mutex.RUnlock()

	if modelID == "" {
		modelID = "none"
	}

	return &RuntimeStatus{
		Status:    sandbox.Status,
		Container: sandbox.Container,
		ModelID:   modelID,
		StartedAt: sandbox.CreatedAt,
		Uptime:    uptime,
		Resources: sandbox.Resources,
	}, nil
}

// Cleanup 清理沙箱环境
func (sm *SandboxManager) Cleanup() error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// 停止并删除所有沙箱
	for sandboxID := range sm.sandboxes {
		if err := sm.StopSandbox(sandboxID); err != nil && err.Error() != fmt.Sprintf("sandbox %s is not running", sandboxID) {
			fmt.Printf("Error stopping sandbox %s: %v\n", sandboxID, err)
		}
		delete(sm.sandboxes, sandboxID)
	}

	fmt.Println("Cleaned up sandbox environment")

	return nil
}

// ExecuteCommand 在沙箱中执行命令
func (sm *SandboxManager) ExecuteCommand(sandboxID string, command []string) (string, error) {
	sm.mutex.RLock()
	// 检查沙箱是否存在
	_, exists := sm.sandboxes[sandboxID]
	if !exists {
		sm.mutex.RUnlock()
		return "", fmt.Errorf("sandbox %s not found", sandboxID)
	}
	sm.mutex.RUnlock()

	// 执行命令
	output := fmt.Sprintf("Executed command %v in sandbox %s", command, sandboxID)

	fmt.Printf("Executing command %v in sandbox %s\n", command, sandboxID)
	fmt.Printf("Command output: %s\n", output)

	return output, nil
}

// SandboxRuntime 沙箱运行时
type SandboxRuntime struct {
	manager *SandboxManager
}

// NewSandboxRuntime 创建沙箱运行时
func NewSandboxRuntime(config *config.Config) (*SandboxRuntime, error) {
	// 创建模型管理器
	modelManager, err := model.NewModelManager(&config.Model)
	if err != nil {
		return nil, fmt.Errorf("failed to create model manager: %v", err)
	}

	// 创建沙箱管理器
	manager, err := NewSandboxManager(&config.Sandbox, modelManager)
	if err != nil {
		return nil, err
	}

	return &SandboxRuntime{
		manager: manager,
	}, nil
}

// RunModel 运行模型
func (sr *SandboxRuntime) RunModel(containerName string, modelID string, input string) (string, error) {
	// 这里可以根据容器名称找到对应的沙箱ID
	// 简化实现，直接使用容器名称作为沙箱ID
	sandboxID := containerName

	// 确保沙箱存在
	_, err := sr.manager.GetSandbox(sandboxID)
	if err != nil {
		// 沙箱不存在，创建一个
		_, err = sr.manager.CreateSandbox(containerName)
		if err != nil {
			return "", fmt.Errorf("failed to create sandbox: %v", err)
		}

		// 启动沙箱
		if err := sr.manager.StartSandbox(sandboxID); err != nil {
			return "", fmt.Errorf("failed to start sandbox: %v", err)
		}

		// 加载模型
		if err := sr.manager.LoadModel(sandboxID, modelID); err != nil {
			return "", fmt.Errorf("failed to load model: %v", err)
		}
	}

	// 运行模型
	return sr.manager.RunModel(sandboxID, modelID, input)
}

// GetRuntimeStatus 获取运行时状态
func (sr *SandboxRuntime) GetRuntimeStatus(containerName string) (*RuntimeStatus, error) {
	// 简化实现，直接使用容器名称作为沙箱ID
	sandboxID := containerName
	return sr.manager.GetRuntimeStatus(sandboxID)
}

// ExecuteCommand 执行命令
func (sr *SandboxRuntime) ExecuteCommand(containerName string, command []string) (string, error) {
	// 简化实现，直接使用容器名称作为沙箱ID
	sandboxID := containerName
	return sr.manager.ExecuteCommand(sandboxID, command)
}

// StopRuntime 停止运行时
func (sr *SandboxRuntime) StopRuntime(containerName string) error {
	// 简化实现，直接使用容器名称作为沙箱ID
	sandboxID := containerName
	return sr.manager.StopSandbox(sandboxID)
}
