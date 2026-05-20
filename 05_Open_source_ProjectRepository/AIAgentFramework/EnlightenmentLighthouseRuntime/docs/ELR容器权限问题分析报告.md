# ELR 容器权限限制问题分析报告

## 1. 问题概述

### 1.1 问题描述

ELR 容器在尝试访问用户主目录下的 `.elr` 目录时遇到 Trae 沙箱安全限制，导致沙箱启动失败。

### 1.2 错误信息

```
TRAE Sandbox Error: hit restricted
  Not allow operate files: C:\Users\Administrator\.elr
  Hint: You can configure sandbox rules via Settings -> Conversation -> Custom Sandbox Configuration.
```

### 1.3 问题影响

| 影响范围 | 描述 |
|---------|------|
| 沙箱启动 | 无法启动沙箱容器 |
| 模型加载 | 无法通过 ELR CLI 加载模型 |
| 数据持久化 | 无法保存沙箱状态和容器映射 |

---

## 2. 问题原因分析

### 2.1 根本原因

ELR 容器默认将数据目录设置为用户主目录下的 `.elr` 目录：

```go
// elr/runtime.go:1071
config.DataDir = filepath.Join(homeDir, ".elr", "data")
```

而 Trae 的沙箱安全机制默认不允许访问用户主目录下的隐藏目录，导致权限拒绝。

### 2.2 代码层面分析

**问题代码位置**：`D:\ELR\EnlightenmentLighthouseRuntime\elr\runtime.go`

```go
// 默认数据目录设置（第 1070-1075 行）
if config.DataDir == "" {
    config.DataDir = filepath.Join(homeDir, ".elr", "data")
} else if strings.HasPrefix(config.DataDir, "~") {
    // Expand ~ to home directory
    config.DataDir = filepath.Join(homeDir, config.DataDir[1:])
}
```

**问题代码位置**：`D:\ELR\EnlightenmentLighthouseRuntime\cli\main.go`

```go
// 默认配置（第 771-775 行）
func defaultConfig() *elr.Config {
    return &elr.Config{
        LogLevel:  "info",
        DataDir:   "~/.elr/data",  // 默认使用用户主目录
        PluginDir: "~/.elr/plugins",
        ...
    }
}
```

### 2.3 问题流程图

```
┌─────────────────────────────────────────────────────────────┐
│                    ELR 启动流程                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. 启动 ELR Runtime                                         │
│         │                                                    │
│         ▼                                                    │
│  2. 加载配置文件 (elr_config.yaml)                           │
│         │                                                    │
│         ▼                                                    │
│  3. 检查 DataDir 配置                                         │
│         │                                                    │
│         ▼                                                    │
│  4. 默认使用 ~/.elr/data                                      │
│         │                                                    │
│         ▼                                                    │
│  5. 尝试创建/访问 ~/.elr 目录                                  │
│         │                                                    │
│         ▼                                                    │
│  6. ❌ Trae 沙箱阻止访问                                       │
│         │                                                    │
│         ▼                                                    │
│  7. 返回权限错误                                               │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. 可能解决方案

### 方案一：使用项目目录作为数据目录（推荐）

**修改思路**：将默认数据目录从用户主目录改为项目目录下的相对路径。

**修改位置**：`D:\ELR\EnlightenmentLighthouseRuntime\cli\main.go`

```go
// 修改 defaultConfig 函数（约第 774 行）
func defaultConfig() *elr.Config {
    return &elr.Config{
        LogLevel:  "info",
        DataDir:   "./data",           // 修改为相对路径
        PluginDir: "./plugins",        // 修改为相对路径
        ...
    }
}
```

**修改位置**：`D:\ELR\EnlightenmentLighthouseRuntime\elr\runtime.go`

```go
// 修改 NewRuntime 函数（约第 1070-1075 行）
if config.DataDir == "" {
    // 使用当前工作目录作为数据目录
    config.DataDir = "./data"
} else if strings.HasPrefix(config.DataDir, "~") {
    // 保留 ~ 扩展支持，但添加警告
    fmt.Println("Warning: Using home directory may cause sandbox permission issues")
    config.DataDir = filepath.Join(homeDir, config.DataDir[1:])
}
```

### 方案二：自动检测并切换目录

**修改思路**：在启动时检测是否在受限环境中，如果是则自动切换到安全目录。

```go
// 在 runtime.go 中添加检测逻辑
func detectSandboxEnvironment() bool {
    // 检测是否在 Trae 沙箱环境中
    if os.Getenv("TRAE_SANDBOX") != "" {
        return true
    }
    // 检查当前目录是否可写
    testFile := "./.sandbox_test"
    if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
        return true
    }
    os.Remove(testFile)
    return false
}
```

### 方案三：环境变量配置

**修改思路**：支持通过环境变量 `ELR_DATA_DIR` 覆盖默认数据目录。

```go
// 在 cli/main.go loadConfig 函数中添加
func loadConfig() (*elr.Config, error) {
    configPath := os.Getenv("ELR_CONFIG")
    if configPath == "" {
        configPath = "~/.elr/config.yaml"
    }
    
    // 添加环境变量支持
    dataDir := os.Getenv("ELR_DATA_DIR")
    if dataDir != "" {
        fmt.Printf("Using data directory from environment: %s\n", dataDir)
    }
    ...
}
```

### 方案四：创建配置向导

**修改思路**：首次启动时引导用户选择数据目录。

```go
// 添加配置向导
func showConfigWizard() {
    fmt.Println("Welcome to ELR Runtime!")
    fmt.Println("Please select a data directory:")
    fmt.Println("1. Use current directory (./data)")
    fmt.Println("2. Use home directory (~/.elr/data)")
    fmt.Println("3. Specify custom directory")
    
    var choice string
    fmt.Scan(&choice)
    
    switch choice {
    case "1":
        os.Setenv("ELR_DATA_DIR", "./data")
    case "2":
        os.Setenv("ELR_DATA_DIR", "~/.elr/data")
    case "3":
        fmt.Print("Enter custom directory: ")
        var customDir string
        fmt.Scan(&customDir)
        os.Setenv("ELR_DATA_DIR", customDir)
    }
}
```

---

## 4. 推荐方案

### 综合方案：多层回退机制

```go
// 修改 runtime.go 的数据目录初始化逻辑
func initDataDir() string {
    // 优先级：环境变量 > 配置文件 > 当前目录 > 用户目录
    
    // 1. 检查环境变量
    if envDir := os.Getenv("ELR_DATA_DIR"); envDir != "" {
        return envDir
    }
    
    // 2. 检查配置文件
    if configDir := loadConfigDataDir(); configDir != "" {
        return configDir
    }
    
    // 3. 尝试当前目录
    if canWrite("./data") {
        return "./data"
    }
    
    // 4. 回退到用户目录（可能需要权限）
    homeDir, _ := os.UserHomeDir()
    return filepath.Join(homeDir, ".elr", "data")
}

func canWrite(path string) bool {
    if err := os.MkdirAll(path, 0755); err != nil {
        return false
    }
    testFile := filepath.Join(path, ".test")
    if err := os.WriteFile(testFile, []byte(""), 0644); err != nil {
        return false
    }
    os.Remove(testFile)
    return true
}
```

---

## 5. 实施建议

### 5.1 优先级排序

| 优先级 | 方案 | 实施难度 | 收益 |
|--------|------|---------|------|
| 1 | 方案一：修改默认目录 | 低 | 高 |
| 2 | 方案三：环境变量支持 | 低 | 中 |
| 3 | 方案二：自动检测 | 中 | 高 |
| 4 | 方案四：配置向导 | 高 | 中 |

### 5.2 实施步骤

1. **第一步**：修改 `cli/main.go` 的 `defaultConfig` 函数，将默认数据目录改为 `./data`
2. **第二步**：修改 `elr/runtime.go` 的 `NewRuntime` 函数，添加目录创建和检测逻辑
3. **第三步**：添加环境变量 `ELR_DATA_DIR` 支持
4. **第四步**：添加启动时的目录可写性检测

### 5.3 测试验证

| 测试场景 | 预期结果 |
|---------|---------|
| 默认启动 | 使用 ./data 目录 |
| 设置 ELR_DATA_DIR | 使用指定目录 |
| 配置文件指定 | 使用配置文件目录 |
| 目录不可写 | 自动回退或报错提示 |

---

## 6. 风险评估

### 6.1 方案一风险

| 风险 | 描述 | 缓解措施 |
|------|------|---------|
| 数据迁移 | 现有用户数据在旧目录 | 提供数据迁移脚本 |
| 权限问题 | 当前目录可能也不可写 | 添加检测和回退机制 |
| 兼容性 | 可能影响现有配置 | 保持环境变量覆盖能力 |

### 6.2 兼容性考虑

```go
// 保持向后兼容
func initDataDir() string {
    // 检查旧目录是否存在且可写
    homeDir, _ := os.UserHomeDir()
    oldDir := filepath.Join(homeDir, ".elr", "data")
    if _, err := os.Stat(oldDir); err == nil && canWrite(oldDir) {
        // 如果旧目录存在且可写，继续使用
        return oldDir
    }
    // 否则使用新默认目录
    return "./data"
}
```

---

## 7. 代码优化建议

### 7.1 runtime.go 优化

**文件**：`D:\ELR\EnlightenmentLighthouseRuntime\elr\runtime.go`

```go
// 添加目录初始化函数（建议位置：约第 1060 行）
func initDataDirectory(config *Config) error {
    // 优先级：环境变量 > 配置 > 当前目录 > 用户目录
    
    // 1. 环境变量
    if envDir := os.Getenv("ELR_DATA_DIR"); envDir != "" {
        config.DataDir = envDir
    }
    
    // 2. 如果未设置，使用当前目录作为默认
    if config.DataDir == "" {
        config.DataDir = "./data"
    }
    
    // 3. 处理 ~ 路径
    if strings.HasPrefix(config.DataDir, "~") {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return fmt.Errorf("failed to get home directory: %v", err)
        }
        config.DataDir = filepath.Join(homeDir, config.DataDir[1:])
    }
    
    // 4. 创建目录
    if err := os.MkdirAll(config.DataDir, 0755); err != nil {
        // 尝试回退到临时目录
        tempDir := filepath.Join(os.TempDir(), "elr_data")
        fmt.Printf("Warning: Failed to create data directory %s, using %s\n", config.DataDir, tempDir)
        config.DataDir = tempDir
        if err := os.MkdirAll(config.DataDir, 0755); err != nil {
            return fmt.Errorf("failed to create data directory: %v", err)
        }
    }
    
    return nil
}
```

### 7.2 cli/main.go 优化

**文件**：`D:\ELR\EnlightenmentLighthouseRuntime\cli\main.go`

```go
// 修改 defaultConfig 函数（约第 771-775 行）
func defaultConfig() *elr.Config {
    return &elr.Config{
        LogLevel:  "info",
        DataDir:   "./data",           // 修改为相对路径
        PluginDir: "./plugins",        // 修改为相对路径
        FileDirectories: make(map[string]string),
        PythonVersions: make(map[string]string),
        Platform: struct {
            Linux struct {
                Enabled bool `yaml:"enabled"`
            } `yaml:"linux"`
            Windows struct {
                Enabled bool `yaml:"enabled"`
            } `yaml:"windows"`
        }{
            Linux:   struct{ Enabled bool }{Enabled: true},
            Windows: struct{ Enabled bool }{Enabled: true},
        },
    }
}
```

---

## 8. 总结

### 8.1 问题核心

ELR 容器默认使用用户主目录 `~/.elr` 作为数据目录，但 Trae 沙箱安全机制阻止了对该目录的访问。

### 8.2 解决方案推荐

**最优方案**：修改默认数据目录为项目相对路径 `./data`，同时支持环境变量和配置文件覆盖。

### 8.3 实施优先级

1. 修改默认数据目录（立即实施）
2. 添加环境变量支持（立即实施）
3. 添加目录检测和回退机制（后续优化）
4. 添加配置向导（可选功能）

---

## 9. 修复实施记录

### 9.1 已完成修改

#### 修改1：cli/main.go - defaultConfig 函数（第 771-775 行）

**修改前：**
```go
func defaultConfig() *elr.Config {
    return &elr.Config{
        LogLevel:  "info",
        DataDir:   "~/.elr/data",
        PluginDir: "~/.elr/plugins",
        ...
    }
}
```

**修改后：**
```go
func defaultConfig() *elr.Config {
    return &elr.Config{
        LogLevel:  "info",
        DataDir:   "./data",
        PluginDir: "./plugins",
        ...
    }
}
```

#### 修改2：elr/runtime.go - 添加 initDataDirectory 函数（第 1061-1091 行）

**新增函数：**
```go
func initDataDirectory(config *Config, homeDir string) error {
    // 1. Check environment variable first
    if envDir := os.Getenv("ELR_DATA_DIR"); envDir != "" {
        config.DataDir = envDir
    }
    
    // 2. If not set, use current directory as default
    if config.DataDir == "" {
        config.DataDir = "./data"
    }
    
    // 3. Handle ~ path expansion
    if strings.HasPrefix(config.DataDir, "~") {
        config.DataDir = filepath.Join(homeDir, config.DataDir[1:])
    }
    
    // 4. Try to create the directory with fallback
    if err := os.MkdirAll(config.DataDir, 0755); err != nil {
        tempDir := filepath.Join(os.TempDir(), "elr_data")
        fmt.Printf("Warning: Failed to create data directory %s, using %s\n", config.DataDir, tempDir)
        config.DataDir = tempDir
        if err := os.MkdirAll(config.DataDir, 0755); err != nil {
            return fmt.Errorf("failed to create data directory: %v", err)
        }
    }
    return nil
}
```

#### 修改3：elr/runtime.go - 添加 initPluginDirectory 函数（第 1093-1122 行）

**新增函数：**
```go
func initPluginDirectory(config *Config, homeDir string) error {
    // 1. Check environment variable first
    if envDir := os.Getenv("ELR_PLUGIN_DIR"); envDir != "" {
        config.PluginDir = envDir
    }
    
    // 2. If not set, use current directory as default
    if config.PluginDir == "" {
        config.PluginDir = "./plugins"
    }
    
    // 3. Handle ~ path expansion
    if strings.HasPrefix(config.PluginDir, "~") {
        config.PluginDir = filepath.Join(homeDir, config.PluginDir[1:])
    }
    
    // 4. Try to create the directory with fallback
    if err := os.MkdirAll(config.PluginDir, 0755); err != nil {
        tempDir := filepath.Join(os.TempDir(), "elr_plugins")
        fmt.Printf("Warning: Failed to create plugin directory %s, using %s\n", config.PluginDir, tempDir)
        config.PluginDir = tempDir
        if err := os.MkdirAll(config.PluginDir, 0755); err != nil {
            return fmt.Errorf("failed to create plugin directory: %v", err)
        }
    }
    return nil
}
```

#### 修改4：elr/runtime.go - NewRuntime 函数（第 1124-1140 行）

**修改前：**
```go
// Initialize data directory
if config.DataDir == "" {
    config.DataDir = filepath.Join(homeDir, ".elr", "data")
} else if strings.HasPrefix(config.DataDir, "~") {
    config.DataDir = filepath.Join(homeDir, config.DataDir[1:])
}
if err := os.MkdirAll(config.DataDir, 0755); err != nil {
    return nil, fmt.Errorf("failed to create data directory: %v", err)
}

// Initialize plugin directory
if config.PluginDir == "" {
    config.PluginDir = filepath.Join(homeDir, ".elr", "plugins")
} else if strings.HasPrefix(config.PluginDir, "~") {
    config.PluginDir = filepath.Join(homeDir, config.PluginDir[1:])
}
if err := os.MkdirAll(config.PluginDir, 0755); err != nil {
    return nil, fmt.Errorf("failed to create plugin directory: %v", err)
}
```

**修改后：**
```go
// Initialize data directory with fallback mechanism
// Priority: environment variable > config > current directory > user home directory
if err := initDataDirectory(config, homeDir); err != nil {
    return nil, err
}

// Initialize plugin directory with fallback mechanism
if err := initPluginDirectory(config, homeDir); err != nil {
    return nil, err
}
```

### 9.2 修改效果

| 修改项 | 效果 |
|--------|------|
| 默认数据目录 | 从 `~/.elr/data` 改为 `./data` |
| 默认插件目录 | 从 `~/.elr/plugins` 改为 `./plugins` |
| 环境变量支持 | 添加 `ELR_DATA_DIR` 和 `ELR_PLUGIN_DIR` 环境变量支持 |
| 自动回退机制 | 当目录创建失败时自动回退到临时目录 |

### 9.3 测试验证

| 测试场景 | 预期结果 | 状态 |
|---------|---------|------|
| 默认启动 | 使用 `./data` 目录 | ✅ 通过 |
| 设置 `ELR_DATA_DIR` 环境变量 | 使用指定目录 | ✅ 通过 |
| 设置 `ELR_PLUGIN_DIR` 环境变量 | 使用指定目录 | ✅ 通过 |
| 目录不可写 | 自动回退到临时目录并显示警告 | ✅ 通过 |

---

## 10. 报告信息

| 项目 | 内容 |
|------|------|
| **报告作者** | 硅基伙伴・清源（ELR-SELLM项目代码实现负责人） |
| **报告时间** | 2026-05-20 |
| **问题状态** | ✅ 已修复 |
| **修复实施方** | 代码织梦者（ELR容器开发团队） |
| **协同方式** | 碳硅协同对位法 |
| **修复版本** | EnlightenmentLighthouseRuntime v1.0.1 |

---

---

## 附录：相关代码位置

| 文件 | 行号 | 描述 |
|------|------|------|
| `elr/runtime.go` | 1070-1075 | 默认数据目录设置 |
| `cli/main.go` | 771-775 | 默认配置定义 |
| `cli/main.go` | 724-768 | 配置加载逻辑 |
| `elr/runtime.go` | 1668 | 默认数据目录函数 |
