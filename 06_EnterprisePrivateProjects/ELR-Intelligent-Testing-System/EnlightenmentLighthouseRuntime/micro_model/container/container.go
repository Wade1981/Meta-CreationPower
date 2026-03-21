package container

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"micro-model/config"
)

// ContainerManager 容器管理器
type ContainerManager struct {
	config *config.ContainerConfig
	client *client.Client
}

// Container 容器信息
type Container struct {
	Name      string    `json:"name"`
	ID        string    `json:"id"`
	Image     string    `json:"image"`
	Status    string    `json:"status"`
	ModelID   string    `json:"model_id"`
	CreatedAt time.Time `json:"created_at"`
	StartedAt time.Time `json:"started_at"`
}

// NewContainerManager 创建容器管理器
func NewContainerManager(config *config.ContainerConfig) (*ContainerManager, error) {
	// 初始化Docker客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}

	// 测试Docker连接
	_, err = cli.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Docker: %v", err)
	}

	return &ContainerManager{
		config: config,
		client: cli,
	}, nil
}

// CreateContainer 创建容器
func (c *ContainerManager) CreateContainer(name string, modelID string, resources map[string]interface{}) error {
	ctx := context.Background()

	// 准备容器配置
	containerConfig := &container.Config{
		Image: c.config.BaseImage,
		Cmd:   []string{"python", "-c", "while True: import time; time.sleep(1)"},
		Env:   []string{fmt.Sprintf("MODEL_ID=%s", modelID)},
	}

	// 准备主机配置
	hostConfig := &container.HostConfig{
		NetworkMode: container.NetworkMode(c.config.NetworkMode),
		Resources: container.Resources{
			Memory:     1024 * 1024 * 1024, // 1GB
			MemorySwap: -1,
			CPUShares:  1024,
			CPUQuota:   int64(c.config.CPULimit * 100000),
		},
	}

	// 覆盖资源限制（如果提供）
	if memory, ok := resources["memory"].(string); ok {
		// 这里可以解析内存字符串，如 "2G"，并设置相应的内存限制
	}

	if cpu, ok := resources["cpu"].(int); ok {
		hostConfig.Resources.CPUQuota = int64(cpu * 100000)
	}

	// 创建容器
	resp, err := c.client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, name)
	if err != nil {
		return fmt.Errorf("failed to create container: %v", err)
	}

	fmt.Printf("Container %s created successfully with ID: %s\n", name, resp.ID)
	return nil
}

// StartContainer 启动容器
func (c *ContainerManager) StartContainer(name string) error {
	ctx := context.Background()

	// 启动容器
	if err := c.client.ContainerStart(ctx, name, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	fmt.Printf("Container %s started successfully\n", name)
	return nil
}

// StopContainer 停止容器
func (c *ContainerManager) StopContainer(name string) error {
	ctx := context.Background()

	// 停止容器
	if err := c.client.ContainerStop(ctx, name, container.StopOptions{}); err != nil {
		return fmt.Errorf("failed to stop container: %v", err)
	}

	fmt.Printf("Container %s stopped successfully\n", name)
	return nil
}

// RemoveContainer 删除容器
func (c *ContainerManager) RemoveContainer(name string) error {
	ctx := context.Background()

	// 删除容器
	if err := c.client.ContainerRemove(ctx, name, types.ContainerRemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("failed to remove container: %v", err)
	}

	fmt.Printf("Container %s removed successfully\n", name)
	return nil
}

// GetContainer 获取容器信息
func (c *ContainerManager) GetContainer(name string) (*Container, error) {
	ctx := context.Background()

	// 获取容器信息
	containerJSON, err := c.client.ContainerInspect(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %v", err)
	}

	// 提取模型ID（从环境变量中）
	var modelID string
	for _, env := range containerJSON.Config.Env {
		if len(env) > 8 && env[:8] == "MODEL_ID=" {
			modelID = env[8:]
			break
		}
	}

	return &Container{
		Name:      containerJSON.Name[1:], // 移除开头的斜杠
		ID:        containerJSON.ID,
		Image:     containerJSON.Config.Image,
		Status:    containerJSON.State.Status,
		ModelID:   modelID,
		CreatedAt: containerJSON.Created,
		StartedAt: containerJSON.State.StartedAt,
	}, nil
}

// ListContainers 列出所有容器
func (c *ContainerManager) ListContainers() ([]*Container, error) {
	ctx := context.Background()

	// 列出所有容器
	containers, err := c.client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %v", err)
	}

	var result []*Container
	for _, container := range containers {
		// 获取完整的容器信息
		containerInfo, err := c.GetContainer(container.ID)
		if err == nil {
			result = append(result, containerInfo)
		}
	}

	return result, nil
}

// Cleanup 清理资源
func (c *ContainerManager) Cleanup() error {
	// 关闭Docker客户端连接
	if err := c.client.Close(); err != nil {
		return fmt.Errorf("failed to close Docker client: %v", err)
	}

	fmt.Println("Container manager cleanup completed")
	return nil
}
