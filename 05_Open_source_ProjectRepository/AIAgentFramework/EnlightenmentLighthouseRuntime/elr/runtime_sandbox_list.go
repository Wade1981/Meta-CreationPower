package elr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// RunningSandbox 正在运行的沙箱信息
type RunningSandbox struct {
	SandboxID   string    `json:"sandbox_id"`
	ContainerID string    `json:"container_id"`
	StartedAt   time.Time `json:"started_at"`
}

// RuntimeSandboxList 运行时沙箱列表管理器
type RuntimeSandboxList struct {
	sandboxes map[string]*RunningSandbox // sandboxID -> sandbox
	mutex     sync.RWMutex
}

// NewRuntimeSandboxList 创建新的运行时沙箱列表管理器
func NewRuntimeSandboxList() *RuntimeSandboxList {
	return &RuntimeSandboxList{
		sandboxes: make(map[string]*RunningSandbox),
	}
}

// AddSandbox 添加沙箱到运行时列表
func (rsl *RuntimeSandboxList) AddSandbox(sandboxID, containerID string) error {
	rsl.mutex.Lock()
	defer rsl.mutex.Unlock()

	if _, exists := rsl.sandboxes[sandboxID]; exists {
		return fmt.Errorf("sandbox %s is already running", sandboxID)
	}

	rs := &RunningSandbox{
		SandboxID:   sandboxID,
		ContainerID: containerID,
		StartedAt:   time.Now(),
	}

	rsl.sandboxes[sandboxID] = rs
	fmt.Printf("Sandbox %s added to runtime list (container: %s)\n", sandboxID, containerID)
	return nil
}

// RemoveSandbox 从运行时列表移除沙箱
func (rsl *RuntimeSandboxList) RemoveSandbox(sandboxID string) error {
	rsl.mutex.Lock()
	defer rsl.mutex.Unlock()

	if _, exists := rsl.sandboxes[sandboxID]; !exists {
		return fmt.Errorf("sandbox %s is not running", sandboxID)
	}

	delete(rsl.sandboxes, sandboxID)
	fmt.Printf("Sandbox %s removed from runtime list\n", sandboxID)
	return nil
}

// GetSandbox 获取沙箱信息
func (rsl *RuntimeSandboxList) GetSandbox(sandboxID string) (*RunningSandbox, bool) {
	rsl.mutex.RLock()
	defer rsl.mutex.RUnlock()

	rs, exists := rsl.sandboxes[sandboxID]
	return rs, exists
}

// IsSandboxRunning 检查沙箱是否在运行
func (rsl *RuntimeSandboxList) IsSandboxRunning(sandboxID string) bool {
	rsl.mutex.RLock()
	defer rsl.mutex.RUnlock()

	_, exists := rsl.sandboxes[sandboxID]
	return exists
}

// ListSandboxes 列出所有运行的沙箱
func (rsl *RuntimeSandboxList) ListSandboxes() []*RunningSandbox {
	rsl.mutex.RLock()
	defer rsl.mutex.RUnlock()

	var sandboxes []*RunningSandbox
	for _, rs := range rsl.sandboxes {
		sandboxes = append(sandboxes, rs)
	}
	return sandboxes
}

// GetSandboxesByContainer 获取指定容器的所有沙箱
func (rsl *RuntimeSandboxList) GetSandboxesByContainer(containerID string) []*RunningSandbox {
	rsl.mutex.RLock()
	defer rsl.mutex.RUnlock()

	var sandboxes []*RunningSandbox
	for _, rs := range rsl.sandboxes {
		if rs.ContainerID == containerID {
			sandboxes = append(sandboxes, rs)
		}
	}
	return sandboxes
}

// Clear 清除所有沙箱
func (rsl *RuntimeSandboxList) Clear() {
	rsl.mutex.Lock()
	defer rsl.mutex.Unlock()

	rsl.sandboxes = make(map[string]*RunningSandbox)
	fmt.Println("Runtime sandbox list cleared")
}

// 全局运行时沙箱列表实例
var globalRuntimeSandboxList *RuntimeSandboxList

// InitRuntimeSandboxList 初始化全局运行时沙箱列表
func InitRuntimeSandboxList() {
	if globalRuntimeSandboxList == nil {
		globalRuntimeSandboxList = NewRuntimeSandboxList()
	}
}

// GetRuntimeSandboxList 获取全局运行时沙箱列表
func GetRuntimeSandboxList() *RuntimeSandboxList {
	if globalRuntimeSandboxList == nil {
		InitRuntimeSandboxList()
	}
	return globalRuntimeSandboxList
}

// GetSandboxStatus 获取沙箱的运行状态
// 返回格式：has run|running 或 has run|Norunning 或 stopped
func GetSandboxStatus(sandboxID string) string {
	rsl := GetRuntimeSandboxList()
	if rsl.IsSandboxRunning(sandboxID) {
		rs, _ := rsl.GetSandbox(sandboxID)
		rcl := GetRuntimeContainerList()

		// 检查容器是否在运行时列表中
		// rs.ContainerID 可能是容器名称或ID
		if _, exists := rcl.GetContainer(rs.ContainerID); exists {
			return "has run|running"
		}

		// 如果 ContainerID 是容器名称，尝试查找对应的容器ID
		containerID := findContainerIDByName(rs.ContainerID)
		if containerID != "" {
			if _, exists := rcl.GetContainer(containerID); exists {
				return "has run|running"
			}
		}

		return "has run|Norunning"
	}

	return "stopped"
}

// FindContainerIDByName 通过容器名称查找容器ID
func FindContainerIDByName(containerName string) string {
	return findContainerIDByName(containerName)
}

// findContainerIDByName 通过容器名称查找容器ID（内部函数）
func findContainerIDByName(containerName string) string {
	// 获取用户主目录
	homeDir, _ := os.UserHomeDir()
	if homeDir == "" {
		return ""
	}

	// 构建容器目录
	containersDir := filepath.Join(homeDir, ".elr", "data", "containers")

	// 遍历容器目录
	entries, err := os.ReadDir(containersDir)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if entry.IsDir() {
			configPath := filepath.Join(containersDir, entry.Name(), "config.json")
			data, err := os.ReadFile(configPath)
			if err != nil {
				continue
			}

			var config struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
			if err := json.Unmarshal(data, &config); err != nil {
				continue
			}

			// 查找名称匹配的容器
			if config.Name == containerName {
				return config.ID
			}
		}
	}

	return ""
}

// findContainerFromSandboxState 从 sandbox-state.json 读取容器信息
func findContainerFromSandboxState(sandboxID, dataDir string) (string, error) {
	sandboxDir := filepath.Join(dataDir, "sandboxes", sandboxID)
	stateFile := filepath.Join(sandboxDir, "sandbox-state.json")
	
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return "", fmt.Errorf("failed to read sandbox state: %w", err)
	}
	
	var state struct {
		ContainerID string `json:"container_id"`
		ContainerName string `json:"container_name"`
	}
	
	if err := json.Unmarshal(data, &state); err != nil {
		return "", fmt.Errorf("failed to parse sandbox state: %w", err)
	}
	
	// 优先返回 container_id，如果为空则返回 container_name
	if state.ContainerID != "" {
		return state.ContainerID, nil
	}
	if state.ContainerName != "" {
		return state.ContainerName, nil
	}
	
	return "", fmt.Errorf("no container information found in sandbox state")
}
