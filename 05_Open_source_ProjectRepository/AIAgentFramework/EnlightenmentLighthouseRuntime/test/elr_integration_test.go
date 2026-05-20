package test

import (
	"fmt"
	"testing"
	"time"

	"elr"
)

// TestContainerPerformance 测试容器性能优化
func TestContainerPerformance(t *testing.T) {
	// 创建运行时
	runtime, err := elr.NewRuntime(&elr.Config{
		DataDir: "./test_data",
		Platform: struct {
			Linux struct {
				UseNamespaces bool `yaml:"use_namespaces"`
				UseCgroups    bool `yaml:"use_cgroups"`
			} `yaml:"linux"`
			Windows struct {
				UseJobObjects bool   `yaml:"use_job_objects"`
				UseWSL        bool   `yaml:"use_wsl"`
				UseContainers bool   `yaml:"use_containers"`
				IsolationType string `yaml:"isolation_type"`
			} `yaml:"windows"`
			Darwin struct {
				UseSandbox bool `yaml:"use_sandbox"`
				UseSpctl   bool `yaml:"use_spctl"`
			} `yaml:"darwin"`
		}{},
		Network: struct {
			Enable  bool   `yaml:"enable"`
			Bridge  string `yaml:"bridge"`
			Subnet  string `yaml:"subnet"`
			APIPorts struct {
				DesktopAPI int `yaml:"desktop_api"`
				PublicAPI  int `yaml:"public_api"`
				ModelAPI   int `yaml:"model_api"`
			} `yaml:"api_ports"`
		}{},
	})

	if err != nil {
		t.Fatalf("Failed to create runtime: %v", err)
	}

	// 启动运行时
	if err := runtime.Start(); err != nil {
		t.Fatalf("Failed to start runtime: %v", err)
	}
	defer runtime.Stop()

	// 创建容器配置
	containerConfig := elr.ContainerConfig{
		ID:        "test-container-1",
		Name:      "Test Container 1",
		Image:     "ubuntu:latest",
		Command:   "cmd.exe",
		Args:      []string{"/c", "echo Hello World"},
		MemoryLimit: "512MB",
		CPULimit:    50,
	}

	// 测试容器创建和启动性能
	startTime := time.Now()
	container, err := runtime.CreateContainer(containerConfig)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}

	if err := container.Start(); err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	elapsed := time.Since(startTime)

	fmt.Printf("Container creation and startup time: %v\n", elapsed)
	if elapsed > 5*time.Second {
		t.Errorf("Container startup took too long: %v", elapsed)
	}

	// 测试容器资源监控
	if container.ResourceMonitor == nil {
		t.Error("ResourceMonitor should not be nil")
	}

	// 测试容器停止
	if err := container.Stop(); err != nil {
		t.Fatalf("Failed to stop container: %v", err)
	}

	// 测试容器删除
	if err := runtime.DeleteContainer(container.ID); err != nil {
		t.Fatalf("Failed to delete container: %v", err)
	}
}

// TestNetworkService 测试网络服务增强功能
func TestNetworkService(t *testing.T) {
	// 创建运行时
	runtime, err := elr.NewRuntime(&elr.Config{
		DataDir: "./test_data",
		Platform: struct {
			Linux struct {
				UseNamespaces bool `yaml:"use_namespaces"`
				UseCgroups    bool `yaml:"use_cgroups"`
			} `yaml:"linux"`
			Windows struct {
				UseJobObjects bool   `yaml:"use_job_objects"`
				UseWSL        bool   `yaml:"use_wsl"`
				UseContainers bool   `yaml:"use_containers"`
				IsolationType string `yaml:"isolation_type"`
			} `yaml:"windows"`
			Darwin struct {
				UseSandbox bool `yaml:"use_sandbox"`
				UseSpctl   bool `yaml:"use_spctl"`
			} `yaml:"darwin"`
		}{},
		Network: struct {
			Enable  bool   `yaml:"enable"`
			Bridge  string `yaml:"bridge"`
			Subnet  string `yaml:"subnet"`
			APIPorts struct {
				DesktopAPI int `yaml:"desktop_api"`
				PublicAPI  int `yaml:"public_api"`
				ModelAPI   int `yaml:"model_api"`
			} `yaml:"api_ports"`
		}{
			Enable: true,
			APIPorts: struct {
				DesktopAPI int `yaml:"desktop_api"`
				PublicAPI  int `yaml:"public_api"`
				ModelAPI   int `yaml:"model_api"`
			}{
				DesktopAPI: 8081,
				PublicAPI:  8080,
				ModelAPI:   8082,
			},
		},
	})

	if err != nil {
		t.Fatalf("Failed to create runtime: %v", err)
	}

	// 启动运行时
	if err := runtime.Start(); err != nil {
		t.Fatalf("Failed to start runtime: %v", err)
	}
	defer runtime.Stop()

	// 测试网络管理器初始化
	if runtime.NetworkManager == nil {
		t.Error("NetworkManager should not be nil")
	}

	// 测试安全管理器和网络隔离器
	if runtime.NetworkManager.(*elr.NetworkManager).SecurityManager == nil {
		t.Error("SecurityManager should not be nil")
	}

	if runtime.NetworkManager.(*elr.NetworkManager).NetworkIsolator == nil {
		t.Error("NetworkIsolator should not be nil")
	}

	// 创建测试容器
	containerConfig := elr.ContainerConfig{
		ID:        "test-container-2",
		Name:      "Test Container 2",
		Image:     "ubuntu:latest",
		Command:   "cmd.exe",
		Args:      []string{"/c", "echo Hello World"},
	}

	container, err := runtime.CreateContainer(containerConfig)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}
	defer runtime.DeleteContainer(container.ID)

	// 测试网络隔离
	networkIsolator := runtime.NetworkManager.(*elr.NetworkManager).NetworkIsolator
	if err := networkIsolator.ApplyNetworkIsolation(container); err != nil {
		t.Fatalf("Failed to apply network isolation: %v", err)
	}

	// 测试网络配置获取
	config, exists := networkIsolator.GetNetworkConfig(container.ID)
	if !exists {
		t.Error("Network config should exist")
	}

	if config.ContainerID != container.ID {
		t.Errorf("Container ID mismatch: expected %s, got %s", container.ID, config.ContainerID)
	}

	// 测试网络隔离移除
	if err := networkIsolator.RemoveNetworkIsolation(container.ID); err != nil {
		t.Fatalf("Failed to remove network isolation: %v", err)
	}
}

// TestModelManagement 测试模型管理系统
func TestModelManagement(t *testing.T) {
	// 测试模型管理器初始化
	config := &elr.Config{
		Model: struct {
			ModelDir string `yaml:"model_dir"`
		}{
			ModelDir: "./test_models",
		},
	}

	// 这里需要导入模型包并创建模型管理器
	// 由于模型管理器在micro_model包中，这里简化处理
	// 实际测试时需要导入相应的包

	// 测试模型加载、切换和卸载功能
	// 由于模型管理器在micro_model包中，这里简化处理
	// 实际测试时需要实现完整的测试逻辑

	fmt.Println("Model management test placeholder")
}

// TestIntegration 测试完整集成
func TestIntegration(t *testing.T) {
	// 创建运行时
	runtime, err := elr.NewRuntime(&elr.Config{
		DataDir: "./test_data",
		Platform: struct {
			Linux struct {
				UseNamespaces bool `yaml:"use_namespaces"`
				UseCgroups    bool `yaml:"use_cgroups"`
			} `yaml:"linux"`
			Windows struct {
				UseJobObjects bool   `yaml:"use_job_objects"`
				UseWSL        bool   `yaml:"use_wsl"`
				UseContainers bool   `yaml:"use_containers"`
				IsolationType string `yaml:"isolation_type"`
			} `yaml:"windows"`
			Darwin struct {
				UseSandbox bool `yaml:"use_sandbox"`
				UseSpctl   bool `yaml:"use_spctl"`
			} `yaml:"darwin"`
		}{},
		Network: struct {
			Enable  bool   `yaml:"enable"`
			Bridge  string `yaml:"bridge"`
			Subnet  string `yaml:"subnet"`
			APIPorts struct {
				DesktopAPI int `yaml:"desktop_api"`
				PublicAPI  int `yaml:"public_api"`
				ModelAPI   int `yaml:"model_api"`
			} `yaml:"api_ports"`
		}{
			Enable: true,
			APIPorts: struct {
				DesktopAPI int `yaml:"desktop_api"`
				PublicAPI  int `yaml:"public_api"`
				ModelAPI   int `yaml:"model_api"`
			}{
				DesktopAPI: 8081,
				PublicAPI:  8080,
				ModelAPI:   8082,
			},
		},
	})

	if err != nil {
		t.Fatalf("Failed to create runtime: %v", err)
	}

	// 启动运行时
	if err := runtime.Start(); err != nil {
		t.Fatalf("Failed to start runtime: %v", err)
	}
	defer runtime.Stop()

	// 测试完整的容器生命周期
	containerConfig := elr.ContainerConfig{
		ID:        "test-container-3",
		Name:      "Test Container 3",
		Image:     "ubuntu:latest",
		Command:   "cmd.exe",
		Args:      []string{"/c", "echo Hello World"},
		MemoryLimit: "512MB",
		CPULimit:    50,
	}

	// 创建容器
	container, err := runtime.CreateContainer(containerConfig)
	if err != nil {
		t.Fatalf("Failed to create container: %v", err)
	}
	defer runtime.DeleteContainer(container.ID)

	// 启动容器
	if err := container.Start(); err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	// 测试网络隔离
	networkIsolator := runtime.NetworkManager.(*elr.NetworkManager).NetworkIsolator
	if err := networkIsolator.ApplyNetworkIsolation(container); err != nil {
		t.Fatalf("Failed to apply network isolation: %v", err)
	}

	// 测试容器状态
	if container.Status != elr.ContainerStatusRunning {
		t.Errorf("Expected container status to be running, got %s", container.Status)
	}

	// 停止容器
	if err := container.Stop(); err != nil {
		t.Fatalf("Failed to stop container: %v", err)
	}

	// 测试容器状态
	if container.Status != elr.ContainerStatusStopped {
		t.Errorf("Expected container status to be stopped, got %s", container.Status)
	}

	fmt.Println("Integration test completed successfully")
}
