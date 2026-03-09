package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"micro-model/model"
)

func main() {
	// 测试模型属性加载
	fmt.Println("Testing model properties loading...")
	
	// 创建临时模型目录
	modelDir := "test_model"
	if err := os.MkdirAll(modelDir, 0755); err != nil {
		log.Fatalf("Failed to create test model directory: %v", err)
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
		},
		"resources": {
			"cpu": {
				"min_cores": 2,
				"recommended_cores": 4
			},
			"memory": {
				"min_ram": "8G",
				"recommended_ram": "16G"
			},
			"gpu": {
				"required": false,
				"min_memory": "4G",
				"recommended_memory": "8G",
				"cuda_version": "12.0"
			}
		},
		"environment_variables": {
			"TEST_VAR": "test_value"
		},
		"entry_points": {
			"test": "echo test"
		}
	}`
	
	if err := os.WriteFile(propertiesPath, []byte(propertiesContent), 0644); err != nil {
		log.Fatalf("Failed to create test model properties file: %v", err)
	}
	
	// 测试加载模型属性
	properties, err := model.LoadModelProperties(modelDir)
	if err != nil {
		log.Fatalf("Failed to load model properties: %v", err)
	}
	
	fmt.Println("Model properties loaded successfully:")
	fmt.Printf("Model name: %s\n", properties.ModelName)
	fmt.Printf("Version: %s\n", properties.Version)
	fmt.Printf("Description: %s\n", properties.Description)
	fmt.Printf("Type: %s\n", properties.Type)
	fmt.Printf("Python version: %s\n", properties.Dependencies.Python)
	fmt.Printf("Pip dependencies: %v\n", properties.Dependencies.Pip)
	fmt.Printf("System dependencies: %v\n", properties.Dependencies.System)
	fmt.Printf("Resource requirements: %s\n", properties.GetResourceRequirements())
	fmt.Printf("Environment variables: %v\n", properties.EnvironmentVariables)
	fmt.Printf("Entry points: %v\n", properties.EntryPoints)
	
	// 测试模型适配器
	fmt.Println("\nTesting model adapter...")
	adapter := model.NewModelAdapter(properties)
	
	// 测试环境设置
	if err := adapter.SetupEnvironment(); err != nil {
		log.Fatalf("Failed to setup environment: %v", err)
	}
	fmt.Println("Environment setup successful")
	
	// 测试容器配置
	containerConfig := adapter.GetContainerConfig()
	fmt.Println("Container config:")
	for key, value := range containerConfig {
		fmt.Printf("%s: %v\n", key, value)
	}
	
	fmt.Println("\nAll tests passed! The adaptation changes are working correctly.")
}
