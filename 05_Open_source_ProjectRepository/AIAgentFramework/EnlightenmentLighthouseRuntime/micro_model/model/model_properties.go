package model

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ModelProperties 模型属性结构
type ModelProperties struct {
	ModelName           string                 `json:"model_name"`
	Version             string                 `json:"version"`
	Description         string                 `json:"description"`
	Type                string                 `json:"type"`
	Dependencies        Dependencies           `json:"dependencies"`
	Resources           Resources              `json:"resources"`
	EnvironmentVariables map[string]string     `json:"environment_variables"`
	EntryPoints         map[string]string     `json:"entry_points"`
}

// Dependencies 依赖结构
type Dependencies struct {
	Python string   `json:"python"`
	Pip    []string `json:"pip"`
	System []string `json:"system"`
}

// Resources 资源需求结构
type Resources struct {
	CPU    CPUResources    `json:"cpu"`
	Memory MemoryResources `json:"memory"`
	GPU    GPUResources    `json:"gpu"`
}

// CPUResources CPU资源需求
type CPUResources struct {
	MinCores        int `json:"min_cores"`
	RecommendedCores int `json:"recommended_cores"`
}

// MemoryResources 内存资源需求
type MemoryResources struct {
	MinRAM        string `json:"min_ram"`
	RecommendedRAM string `json:"recommended_ram"`
}

// GPUResources GPU资源需求
type GPUResources struct {
	Required           bool   `json:"required"`
	MinMemory          string `json:"min_memory"`
	RecommendedMemory  string `json:"recommended_memory"`
	CUDAVersion        string `json:"cuda_version"`
}

// LoadModelProperties 加载模型属性
func LoadModelProperties(modelDir string) (*ModelProperties, error) {
	// 尝试加载模型属性文件
	propertiesPath := filepath.Join(modelDir, "model_properties.json")
	if _, err := os.Stat(propertiesPath); os.IsNotExist(err) {
		// 尝试加载examples目录中的示例属性文件
		propertiesPath = filepath.Join("examples", "model_properties.json")
		if _, err := os.Stat(propertiesPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("model properties file not found")
		}
	}

	// 读取文件内容
	content, err := os.ReadFile(propertiesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read model properties file: %v", err)
	}

	// 解析JSON
	var properties ModelProperties
	if err := json.Unmarshal(content, &properties); err != nil {
		return nil, fmt.Errorf("failed to unmarshal model properties: %v", err)
	}

	return &properties, nil
}

// Validate 验证模型属性
func (p *ModelProperties) Validate() error {
	if p.ModelName == "" {
		return fmt.Errorf("model name is required")
	}

	if p.Dependencies.Python == "" {
		return fmt.Errorf("python version is required")
	}

	return nil
}

// GetResourceRequirements 获取资源需求
func (p *ModelProperties) GetResourceRequirements() string {
	return fmt.Sprintf("CPU: %d-%d cores, Memory: %s-%s, GPU: %v (min %s)",
		p.Resources.CPU.MinCores,
		p.Resources.CPU.RecommendedCores,
		p.Resources.Memory.MinRAM,
		p.Resources.Memory.RecommendedRAM,
		p.Resources.GPU.Required,
		p.Resources.GPU.MinMemory)
}
