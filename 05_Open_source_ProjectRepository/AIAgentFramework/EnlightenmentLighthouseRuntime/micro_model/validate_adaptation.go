package main

import (
	"fmt"
	"os"
	"path/filepath"

	"micro-model/model"
)

func main() {
	// 测试模型属性结构
	fmt.Println("Testing model properties structure...")
	
	// 创建测试模型属性
	properties := &model.ModelProperties{
		ModelName: "test-model",
		Version:   "1.0",
		Description: "Test model",
		Type:      "test",
		Dependencies: model.Dependencies{
			Python: "3.12",
			Pip:    []string{"numpy", "pandas"},
			System: []string{"curl"},
		},
		Resources: model.Resources{
			CPU: model.CPUResources{
				MinCores:        2,
				RecommendedCores: 4,
			},
			Memory: model.MemoryResources{
				MinRAM:        "8G",
				RecommendedRAM: "16G",
			},
			GPU: model.GPUResources{
				Required:           false,
				MinMemory:          "4G",
				RecommendedMemory:  "8G",
				CUDAVersion:        "12.0",
			},
		},
		EnvironmentVariables: map[string]string{
			"TEST_VAR": "test_value",
		},
		EntryPoints: map[string]string{
			"test": "echo test",
		},
	}
	
	// 测试模型属性验证
	if err := properties.Validate(); err != nil {
		fmt.Printf("Model properties validation failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Model properties validation passed")
	
	// 测试资源需求获取
	resourceReq := properties.GetResourceRequirements()
	fmt.Printf("Resource requirements: %s\n", resourceReq)
	
	// 测试模型适配器
	fmt.Println("\nTesting model adapter...")
	adapter := model.NewModelAdapter(properties)
	
	// 测试环境设置
	if err := adapter.SetupEnvironment(); err != nil {
		fmt.Printf("Environment setup failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Environment setup successful")
	
	// 测试容器配置
	containerConfig := adapter.GetContainerConfig()
	fmt.Println("Container config:")
	for key, value := range containerConfig {
		fmt.Printf("%s: %v\n", key, value)
	}
	
	// 测试模型属性文件加载
	fmt.Println("\nTesting model properties file loading...")
	
	// 创建临时模型目录
	modelDir := "test_model"
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		fmt.Printf("Failed to create test model directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(modelDir)
	
	// 创建测试模型属性文件
	propertiesPath := filepath.Join(modelDir, "model_properties.json")
	propertiesContent := `{
		"model_name": "test-model",
		"version": "1.0",
		"description": "Test model",
		"type": "test",
		"dependencies": {
			"python": "3.12",
			"pip": ["numpy", "pandas"],
			"system": ["curl"]
		}
	}`
	
	if err := os.WriteFile(propertiesPath, []byte(propertiesContent), 0644); err != nil {
		fmt.Printf("Failed to create test model properties file: %v\n", err)
		os.Exit(1)
	}
	
	// 测试加载模型属性
	loadedProperties, err := model.LoadModelProperties(modelDir)
	if err != nil {
		fmt.Printf("Failed to load model properties: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Loaded model name: %s\n", loadedProperties.ModelName)
	fmt.Printf("Loaded Python version: %s\n", loadedProperties.Dependencies.Python)
	
	fmt.Println("\nAll tests passed! The adaptation changes are working correctly.")
}
