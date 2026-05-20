package container

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"micro_model/config"
)

// LightweightContainerManager 轻量级容器管理器
type LightweightContainerManager struct {
	config     *config.ContainerConfig
	containers map[string]*Container // 容器ID -> 容器信息
}

// NewLightweightContainerManager 创建轻量级容器管理器
func NewLightweightContainerManager(config *config.ContainerConfig) (*LightweightContainerManager, error) {
	return &LightweightContainerManager{
		config:     config,
		containers: make(map[string]*Container),
	}, nil
}

// CreateContainer 创建容器
func (c *LightweightContainerManager) CreateContainer(name string, modelID string, resources map[string]interface{}) error {
	// 生成容器ID
	containerID := fmt.Sprintf("elr-%d", time.Now().UnixNano()/1000000)
	
	// 创建容器信息
	container := &Container{
		Name:      name,
		ID:        containerID,
		Image:     c.config.BaseImage,
		Status:    "created",
		ModelID:   modelID,
		CreatedAt: time.Now().Format(time.RFC3339),
		StartedAt: "",
	}
	
	// 保存容器信息
	c.containers[containerID] = container
	
	// 创建容器目录
	containerDir := filepath.Join(".", "containers", containerID)
	if err := os.MkdirAll(containerDir, 0755); err != nil {
		return fmt.Errorf("failed to create container directory: %v", err)
	}
	
	// 创建容器配置文件
	configPath := filepath.Join(containerDir, "config.json")
	configData := fmt.Sprintf(`{
	"id": "%s",
	"name": "%s",
	"image": "%s",
	"status": "%s",
	"model_id": "%s",
	"created_at": "%s"
}`, containerID, name, c.config.BaseImage, "created", modelID, container.CreatedAt)
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		return fmt.Errorf("failed to write container config: %v", err)
	}
	
	fmt.Printf("Lightweight container %s created successfully with ID: %s\n", name, containerID)
	return nil
}

// StartContainer 启动容器
func (c *LightweightContainerManager) StartContainer(name string) error {
	// 查找容器
	var container *Container
	for _, c := range c.containers {
		if c.Name == name || c.ID == name {
			container = c
			break
		}
	}
	
	if container == nil {
		return fmt.Errorf("container %s not found", name)
	}
	
	// 更新容器状态
	container.Status = "running"
	container.StartedAt = time.Now().Format(time.RFC3339)
	
	// 更新容器配置文件
	containerDir := filepath.Join(".", "containers", container.ID)
	configPath := filepath.Join(containerDir, "config.json")
	configData := fmt.Sprintf(`{
	"id": "%s",
	"name": "%s",
	"image": "%s",
	"status": "%s",
	"model_id": "%s",
	"created_at": "%s",
	"started_at": "%s"
}`, container.ID, container.Name, container.Image, "running", container.ModelID, container.CreatedAt, container.StartedAt)
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		return fmt.Errorf("failed to update container config: %v", err)
	}
	
	fmt.Printf("Lightweight container %s started successfully\n", name)
	return nil
}

// StopContainer 停止容器
func (c *LightweightContainerManager) StopContainer(name string) error {
	// 查找容器
	var container *Container
	for _, c := range c.containers {
		if c.Name == name || c.ID == name {
			container = c
			break
		}
	}
	
	if container == nil {
		return fmt.Errorf("container %s not found", name)
	}
	
	// 更新容器状态
	container.Status = "stopped"
	
	// 更新容器配置文件
	containerDir := filepath.Join(".", "containers", container.ID)
	configPath := filepath.Join(containerDir, "config.json")
	configData := fmt.Sprintf(`{
	"id": "%s",
	"name": "%s",
	"image": "%s",
	"status": "%s",
	"model_id": "%s",
	"created_at": "%s",
	"started_at": "%s"
}`, container.ID, container.Name, container.Image, "stopped", container.ModelID, container.CreatedAt, container.StartedAt)
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		return fmt.Errorf("failed to update container config: %v", err)
	}
	
	fmt.Printf("Lightweight container %s stopped successfully\n", name)
	return nil
}

// RemoveContainer 删除容器
func (c *LightweightContainerManager) RemoveContainer(name string) error {
	// 查找容器
	var container *Container
	var containerID string
	for id, c := range c.containers {
		if c.Name == name || c.ID == name {
			container = c
			containerID = id
			break
		}
	}
	
	if containerID == "" {
		return fmt.Errorf("container %s not found", name)
	}
	
	// 如果容器正在运行，先停止
	if container.Status == "running" {
		if err := c.StopContainer(name); err != nil {
			fmt.Printf("Warning: failed to stop container %s: %v\n", name, err)
		}
	}
	
	// 删除容器目录
	containerDir := filepath.Join(".", "containers", containerID)
	if err := os.RemoveAll(containerDir); err != nil {
		return fmt.Errorf("failed to remove container directory: %v", err)
	}
	
	// 从容器列表中删除
	delete(c.containers, containerID)
	
	fmt.Printf("Lightweight container %s removed successfully\n", name)
	return nil
}

// GetContainer 获取容器信息
func (c *LightweightContainerManager) GetContainer(name string) (*Container, error) {
	// 查找容器
	for _, container := range c.containers {
		if container.Name == name || container.ID == name {
			return container, nil
		}
	}
	
	return nil, fmt.Errorf("container %s not found", name)
}

// ListContainers 列出所有容器
func (c *LightweightContainerManager) ListContainers() ([]*Container, error) {
	containers := make([]*Container, 0, len(c.containers))
	for _, container := range c.containers {
		containers = append(containers, container)
	}
	return containers, nil
}

// Cleanup 清理资源
func (c *LightweightContainerManager) Cleanup() error {
	// 清理容器目录
	containersDir := filepath.Join(".", "containers")
	if err := os.RemoveAll(containersDir); err != nil {
		return fmt.Errorf("failed to cleanup containers directory: %v", err)
	}
	
	// 清空容器列表
	c.containers = make(map[string]*Container)
	
	fmt.Println("Lightweight container manager cleanup completed")
	return nil
}
