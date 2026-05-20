package elr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// getCurrentTime 获取当前时间字符串
func getCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}

// SandboxContainerMapping 沙箱与容器的映射关系
type SandboxContainerMapping struct {
	SandboxID   string `json:"sandbox_id"`
	ContainerID string `json:"container_id"`
	CreatedAt   string `json:"created_at"`
}

// SandboxContainerManager 沙箱-容器映射管理器
type SandboxContainerManager struct {
	mappings map[string]*SandboxContainerMapping // sandboxID -> mapping
	mutex    sync.RWMutex
	storagePath string
}

// NewSandboxContainerManager 创建新的沙箱-容器映射管理器
func NewSandboxContainerManager(dataDir string) *SandboxContainerManager {
	storagePath := filepath.Join(dataDir, "sandbox_container_mappings.json")
	
	manager := &SandboxContainerManager{
		mappings:    make(map[string]*SandboxContainerMapping),
		storagePath: storagePath,
	}
	
	// 加载已有的映射关系
	manager.loadMappings()
	
	return manager
}

// loadMappings 从磁盘加载映射关系
func (scm *SandboxContainerManager) loadMappings() {
	scm.mutex.Lock()
	defer scm.mutex.Unlock()
	
	if _, err := os.Stat(scm.storagePath); os.IsNotExist(err) {
		return
	}
	
	data, err := os.ReadFile(scm.storagePath)
	if err != nil {
		fmt.Printf("Warning: failed to read sandbox-container mappings: %v\n", err)
		return
	}
	
	var mappings []*SandboxContainerMapping
	if err := json.Unmarshal(data, &mappings); err != nil {
		fmt.Printf("Warning: failed to unmarshal sandbox-container mappings: %v\n", err)
		return
	}
	
	for _, mapping := range mappings {
		scm.mappings[mapping.SandboxID] = mapping
	}
	
	fmt.Printf("Loaded %d sandbox-container mappings from disk\n", len(mappings))
}

// saveMappings 保存映射关系到磁盘（不获取锁，调用者负责处理锁）
func (scm *SandboxContainerManager) saveMappings() error {
	var mappings []*SandboxContainerMapping
	
	// 先获取读锁来读取数据
	scm.mutex.RLock()
	for _, mapping := range scm.mappings {
		mappings = append(mappings, mapping)
	}
	scm.mutex.RUnlock()
	
	data, err := json.MarshalIndent(mappings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sandbox-container mappings: %v", err)
	}
	
	if err := os.MkdirAll(filepath.Dir(scm.storagePath), 0755); err != nil {
		return fmt.Errorf("failed to create storage directory: %v", err)
	}
	
	if err := os.WriteFile(scm.storagePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write sandbox-container mappings: %v", err)
	}
	
	return nil
}

// AddMapping 添加沙箱-容器映射
func (scm *SandboxContainerManager) AddMapping(sandboxID, containerID string) error {
	scm.mutex.Lock()
	
	// 检查是否已存在
	if _, exists := scm.mappings[sandboxID]; exists {
		scm.mutex.Unlock()
		return fmt.Errorf("sandbox %s already mapped to container %s", sandboxID, scm.mappings[sandboxID].ContainerID)
	}
	
	mapping := &SandboxContainerMapping{
		SandboxID:   sandboxID,
		ContainerID: containerID,
		CreatedAt:   getCurrentTime(),
	}
	
	scm.mappings[sandboxID] = mapping
	
	// 释放写锁
	scm.mutex.Unlock()
	
	// 保存到磁盘
	if err := scm.saveMappings(); err != nil {
		fmt.Printf("Warning: failed to save sandbox-container mapping: %v\n", err)
	}
	
	fmt.Printf("Mapped sandbox %s to container %s\n", sandboxID, containerID)
	return nil
}

// GetContainerBySandbox 通过沙箱ID获取容器ID
func (scm *SandboxContainerManager) GetContainerBySandbox(sandboxID string) (string, error) {
	scm.mutex.RLock()
	defer scm.mutex.RUnlock()
	
	mapping, exists := scm.mappings[sandboxID]
	if !exists {
		return "", fmt.Errorf("sandbox %s not found in mapping", sandboxID)
	}
	
	return mapping.ContainerID, nil
}

// RemoveMapping 移除沙箱-容器映射
func (scm *SandboxContainerManager) RemoveMapping(sandboxID string) error {
	scm.mutex.Lock()
	
	if _, exists := scm.mappings[sandboxID]; !exists {
		scm.mutex.Unlock()
		return fmt.Errorf("sandbox %s not found in mapping", sandboxID)
	}
	
	delete(scm.mappings, sandboxID)
	
	// 释放写锁
	scm.mutex.Unlock()
	
	// 保存到磁盘
	if err := scm.saveMappings(); err != nil {
		fmt.Printf("Warning: failed to save sandbox-container mapping: %v\n", err)
	}
	
	fmt.Printf("Removed mapping for sandbox %s\n", sandboxID)
	return nil
}

// ListMappings 列出所有映射关系
func (scm *SandboxContainerManager) ListMappings() []*SandboxContainerMapping {
	scm.mutex.RLock()
	defer scm.mutex.RUnlock()
	
	var mappings []*SandboxContainerMapping
	for _, mapping := range scm.mappings {
		mappings = append(mappings, mapping)
	}
	
	return mappings
}

// IsSandboxInContainer 验证沙箱是否在指定容器中
func (scm *SandboxContainerManager) IsSandboxInContainer(sandboxID, containerID string) bool {
	scm.mutex.RLock()
	defer scm.mutex.RUnlock()
	
	mapping, exists := scm.mappings[sandboxID]
	if !exists {
		return false
	}
	
	return mapping.ContainerID == containerID
}

// VerifySandboxRunningInContainer 验证沙箱是否真的在容器中运行（双重验证）
func (scm *SandboxContainerManager) VerifySandboxRunningInContainer(sandboxID string) (bool, string, error) {
	// 第一步：从映射表中查找容器
	containerID, err := scm.GetContainerBySandbox(sandboxID)
	if err != nil {
		return false, "", fmt.Errorf("sandbox %s not found in mapping: %v", sandboxID, err)
	}
	
	// 第二步：通过 API 检查容器是否在运行
	isContainerRunning, err := checkContainerRunningViaAPI(containerID)
	if err != nil {
		return false, containerID, fmt.Errorf("failed to check if container is running: %v", err)
	}
	
	if !isContainerRunning {
		return false, containerID, fmt.Errorf("container %s is not running", containerID)
	}
	
	// TODO: 第三步：检查沙箱管理器是否在运行（需要从容器信息中获取）
	// 当前我们还没有在 RuntimeContainerList 中存储沙箱管理器进程信息
	// 这部分功能将在后续实现
	
	return true, containerID, nil
}

// checkContainerRunningViaAPI 通过 API 检查容器是否在运行
func checkContainerRunningViaAPI(containerID string) (bool, error) {
	// 尝试连接到 ELR API
	resp, err := http.Get("http://localhost:16888/api/container/running")
	if err != nil {
		return false, fmt.Errorf("failed to connect to ELR API: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("ELR API returned status: %d", resp.StatusCode)
	}
	
	// 解析响应
	var runningContainers []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&runningContainers); err != nil {
		return false, fmt.Errorf("failed to decode API response: %v", err)
	}
	
	// 检查容器是否在运行列表中
	for _, container := range runningContainers {
		if id, ok := container["id"].(string); ok && id == containerID {
			return true, nil
		}
	}
	
	return false, nil
}

// 全局沙箱-容器映射管理器实例
var globalSandboxContainerManager *SandboxContainerManager

// InitSandboxContainerManager 初始化全局沙箱-容器映射管理器
func InitSandboxContainerManager(dataDir string) {
	globalSandboxContainerManager = NewSandboxContainerManager(dataDir)
}

// GetSandboxContainerManager 获取全局沙箱-容器映射管理器
func GetSandboxContainerManager() *SandboxContainerManager {
	return globalSandboxContainerManager
}
