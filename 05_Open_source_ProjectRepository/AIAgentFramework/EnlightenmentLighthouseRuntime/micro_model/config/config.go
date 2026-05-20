package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Config 全局配置结构
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Container  ContainerConfig  `mapstructure:"container"`
	Model      ModelConfig      `mapstructure:"model"`
	Sandbox    SandboxConfig    `mapstructure:"sandbox"`
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
	Python     PythonConfig     `mapstructure:"python"`
}

// PythonConfig Python配置
type PythonConfig struct {
	Path        string `mapstructure:"path"`
	Version     string `mapstructure:"version"`
	AutoInstall bool   `mapstructure:"auto_install"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// ContainerConfig 容器配置
type ContainerConfig struct {
	BaseImage   string `mapstructure:"base_image"`
	NetworkMode string `mapstructure:"network_mode"`
	MemoryLimit string `mapstructure:"memory_limit"`
	CPULimit    int    `mapstructure:"cpu_limit"`
}

// ModelConfig 模型配置
type ModelConfig struct {
	ModelDir string `mapstructure:"model_dir"`
	DefaultModel string `mapstructure:"default_model"`
	DownloadTimeout int `mapstructure:"download_timeout"`
}

// SandboxConfig 沙箱配置
type SandboxConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	IsolationLevel string `mapstructure:"isolation_level"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Interval int `mapstructure:"interval"`
	PrometheusPort int `mapstructure:"prometheus_port"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 设置默认配置
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("container.base_image", "python:3.8-slim")
	viper.SetDefault("container.network_mode", "bridge")
	viper.SetDefault("container.memory_limit", "1G")
	viper.SetDefault("container.cpu_limit", 1)
	viper.SetDefault("model.model_dir", "./model/models")
	viper.SetDefault("model.default_model", "gpt2")
	viper.SetDefault("model.download_timeout", 300)
	viper.SetDefault("sandbox.enabled", true)
	viper.SetDefault("sandbox.isolation_level", "container")
	viper.SetDefault("monitoring.enabled", true)
	viper.SetDefault("monitoring.interval", 5)
	viper.SetDefault("monitoring.prometheus_port", 9090)
	viper.SetDefault("python.path", "")
	viper.SetDefault("python.version", "3.8+")
	viper.SetDefault("python.auto_install", false)

	// 尝试加载配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在，使用默认配置
			fmt.Println("Config file not found, using default config")
		} else {
			// 配置文件存在但有错误
			return nil, fmt.Errorf("failed to read config file: %v", err)
		}
	}

	// 尝试读取 ELR Settings 配置文件
	elrModelDir := getELRModelDir()
	if elrModelDir != "" {
		fmt.Printf("Using model directory from ELR Settings: %s\n", elrModelDir)
		viper.Set("model.model_dir", elrModelDir)
	}

	// 确保模型目录存在
	modelDir := viper.GetString("model.model_dir")
	absModelDir, err := filepath.Abs(modelDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for model directory: %v", err)
	}
	viper.Set("model.model_dir", absModelDir)

	// 解析配置
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}

// getELRModelDir 从 ELR Settings 配置文件中获取模型目录
func getELRModelDir() string {
	// ELR 配置文件路径
	configPath := os.Getenv("ELR_CONFIG")
	if configPath == "" {
		// 默认路径
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configPath = filepath.Join(homeDir, ".elr", "config.yaml")
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return ""
	}

	// 读取配置文件
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}

	// 解析配置
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

	var elrConfig ELRConfig
	if err := yaml.Unmarshal(configBytes, &elrConfig); err != nil {
		return ""
	}

	// 查找模型类型目录
	// 优先使用任何启用的模型类型目录
	for modelType, modelConfig := range elrConfig.Resources.ModelTypes {
		if modelConfig.Enable && modelConfig.Dir != "" {
			fmt.Printf("Using model directory from ELR Settings for model type '%s': %s\n", modelType, modelConfig.Dir)
			return modelConfig.Dir
		}
	}

	// 其次使用任何启用的资源类型目录
	for resourceType, resourceConfig := range elrConfig.Resources.Types {
		if resourceConfig.Enable && resourceConfig.Dir != "" {
			fmt.Printf("Using resource directory from ELR Settings for resource type '%s': %s\n", resourceType, resourceConfig.Dir)
			return resourceConfig.Dir
		}
	}

	return ""
}

// GetConfig 获取配置实例
func GetConfig() (*Config, error) {
	return LoadConfig()
}
