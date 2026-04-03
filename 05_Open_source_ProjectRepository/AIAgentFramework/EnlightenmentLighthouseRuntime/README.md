# Enlightenment Lighthouse Runtime (ELR)

[![GitHub](https://img.shields.io/github/stars/Wade1981/Meta-CreationPower?style=social)](https://github.com/Wade1981/Meta-CreationPower)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Enlightenment Lighthouse Runtime (ELR) 是启蒙灯塔起源团队开发的轻量级、跨平台容器运行环境，专为碳硅协同场景设计。它不依赖Docker，而是使用系统原生的隔离机制，提供了一个统一的、可扩展的平台，支持主流编程语言。

## 项目愿景

构建一个轻量、高效、安全的容器运行环境，成为碳硅协同创新的基础设施，推动人文价值与科技理性的平衡共生。

## 核心特性

### 轻量级设计
- **无外部依赖**：核心使用Go语言开发，编译为静态二进制文件
- **小巧高效**：核心运行时小于10MB，启动时间毫秒级
- **低资源占用**：内存占用低，CPU使用率高

### 跨平台支持
- **Windows**：支持Windows 10/11
- **Linux**：支持主流Linux发行版
- **macOS**：支持macOS 10.15+
- **架构支持**：x86-64, ARM64

### 多语言支持
- **C/C++**：高性能计算和系统级编程
- **Python**：AI模型训练和数据科学
- **JavaScript/Node.js**：Web服务和前端开发
- **Java**：企业级应用和大型系统
- **Go**：高性能后端服务和微服务

### 安全隔离
- **系统级隔离**：使用操作系统原生隔离机制
  - Linux：namespace和cgroup
  - Windows：Job Objects和WSL
  - macOS：sandbox和spctl
- **权限管理**：细粒度的权限控制
- **网络隔离**：容器间网络隔离

### 可扩展性
- **插件架构**：基于插件的模块化设计
- **语言插件**：易于添加新语言支持
- **服务插件**：易于添加新服务和功能
- **平台插件**：易于支持新操作系统

### 分布式能力
- **轻量级服务发现**：基于mDNS的服务发现
- **简单负载均衡**：内置基本负载均衡功能
- **状态管理**：使用分布式配置存储

### 碳硅协同
- **优化的人机交互**：为碳基和硅基智能提供友好接口
- **智能调度**：基于任务类型的资源分配
- **元协议集成**：内置启蒙灯塔元协议支持

### 网络状态查询
- **专属通信地址查询**：查询与ELR-Desktop-Assistant的专属通信地址和端口
- **公众API查询**：查询与公众接触应用的通信API地址和端口
- **模型API查询**：查询模型API的地址、端口和网络状态
- **实时状态**：提供网络服务的实时运行状态

### 令牌管理
- **令牌生成**：生成ELR容器访问令牌
- **令牌验证**：验证令牌的有效性和过期状态
- **令牌刷新**：更新过期的令牌
- **令牌撤销**：撤销不再需要的令牌
- **令牌列表**：查看所有令牌的状态和详情

### 管理员系统
- **管理员创建**：创建具有不同角色的管理员账户
- **权限管理**：为管理员分配容器管理权限
- **角色控制**：支持超级管理员和普通管理员角色
- **令牌管理**：为管理员自动生成访问令牌

### 文件系统管理
- **文件上传**：上传文件到容器文件系统
- **文件下载**：从容器文件系统下载文件
- **权限控制**：基于管理员权限的文件操作控制
- **文件加密**：所有存储的文件自动加密

### 沙箱管理
- **沙箱创建**：为容器创建独立的沙箱环境
- **模型加载**：在沙箱中加载和运行AI模型
- **资源管理**：监控和管理沙箱资源使用
- **自动装载**：容器启动时自动装载到管理员沙箱

### API 服务
- **Desktop API**：桌面应用专用API
- **Public API**：公众应用专用API
- **Model API**：模型服务专用API
- **RESTful接口**：标准化的API接口

## 目录结构

```
EnlightenmentLighthouseRuntime/
├── README.md          # 项目说明文档
├── api/               # API服务实现
│   ├── api_manager.go # API管理器
│   ├── desktop_api.go # 桌面API
│   ├── model_api.go   # 模型API
│   ├── public_api.go  # 公众API
│   └── go.mod         # Go模块依赖
├── cli/               # 命令行工具
│   ├── main.go        # 主入口
│   ├── go.mod         # Go模块依赖
│   └── go.sum         # Go模块依赖校验
├── commands/          # PowerShell命令实现
│   ├── chat.ps1       # 聊天命令
│   ├── container.ps1  # 容器命令
│   ├── help.ps1       # 帮助命令
│   └── ...            # 其他命令
├── components/        # 组件目录
├── docs/              # 文档目录
│   ├── ELR_PowerShell_使用说明.md # PowerShell使用说明
│   └── ...            # 其他文档
├── elr/               # 核心运行时
│   ├── container.go   # 容器管理
│   ├── network.go     # 网络管理
│   ├── runtime.go     # 运行时核心
│   ├── token_manager.go # 令牌管理
│   ├── admin_manager.go # 管理员管理
│   └── go.mod         # Go模块依赖
├── examples/          # 示例
│   ├── cpp/           # C++示例
│   └── python/        # Python示例
├── icons/             # 图标目录
├── micro_model/       # 微型模型系统
│   ├── api/           # API接口
│   ├── config/        # 配置管理
│   ├── container/     # 容器管理
│   ├── model/         # 模型核心
│   ├── monitor/       # 监控服务
│   └── sandbox/       # 沙箱运行时
├── model/             # 模型目录
├── models/            # 模型文件
├── platforms/         # 平台特定实现
│   ├── linux/         # Linux实现
│   ├── windows/       # Windows实现
│   └── darwin/        # macOS实现
├── resource/          # 资源目录
│   └── language.json  # 语言配置
├── scripts/           # 辅助脚本
├── uploads/           # 上传文件目录
├── ELR-Tray-App.ps1   # 托盘应用
├── elr.ps1            # PowerShell主脚本
├── elr.exe            # Go编译的可执行文件
├── go.mod             # Go模块依赖
└── go.sum             # Go模块依赖校验
```

## 技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| 开发语言 | Go | 1.17+ |
| 脚本语言 | PowerShell | 5.1+ |
| 系统隔离 |
| Linux | namespace + cgroup | 内核 4.0+ |
| Windows | Job Objects + WSL | Windows 10+ |
| macOS | sandbox + spctl | macOS 10.15+ |
| 网络 | TCP/IP + HTTP | - |
| 配置 | JSON + YAML | - |
| API | RESTful | - |
| 模型管理 | Python | 3.8+ |

## 快速开始

### 前提条件
- **Go语言**：版本 1.17+（仅用于开发）
- **操作系统**：Windows 10+, Linux, macOS 10.15+
- **系统权限**：需要管理员/root权限（仅用于安装和配置）

### 安装

#### 二进制安装

1. **下载二进制文件**

从GitHub Releases页面下载对应平台的二进制文件：
- Windows: `elr-windows-amd64.exe`
- Linux: `elr-linux-amd64`
- macOS: `elr-darwin-amd64`

2. **安装**

- **Windows**：将二进制文件复制到系统路径（如 `C:\Windows\System32`）
- **Linux/macOS**：将二进制文件复制到系统路径（如 `/usr/local/bin`）并添加执行权限

```bash
chmod +x elr-linux-amd64
sudo mv elr-linux-amd64 /usr/local/bin/elr
```

#### 源码安装

1. **克隆仓库**

```bash
git clone https://github.com/Wade1981/Meta-CreationPower.git
cd Meta-CreationPower/05_Open_source_ProjectRepository/AIAgentFramework/EnlightenmentLighthouseRuntime
```

2. **构建**

```bash
go build -o elr ./cli
```

3. **安装**

```bash
# Linux/macOS
sudo mv elr /usr/local/bin/

# Windows
# 将elr.exe复制到系统路径
```

### 使用

#### 基本命令

```bash
# 查看版本
elr version

# 查看帮助
elr help

# 创建容器
elr create --name my-container --image ubuntu:latest

# 启动容器
elr start-container --id container-id

# 停止容器
elr stop-container --id container-id

# 列出容器
elr list

# 删除容器
elr delete --id container-id
```

#### API 管理

```bash
# 启动 API 服务
elr api start

# 停止 API 服务
elr api stop

# 检查 API 服务状态
elr api status

# 配置 API 设置
elr api config set --api-type desktop --address localhost --port 8081
```

#### 模型管理

```bash
# 列出所有模型
elr model list

# 获取模型信息
elr model get --model-id model-id

# 下载模型
elr model download --model-id model-id --type model-type --url download-url

# 删除模型
elr model delete --model-id model-id

# 安装模型依赖
elr model install-deps --model-id model-id --type dep-type
```

#### 沙箱管理

```bash
# 列出沙箱
elr sandbox list

# 创建沙箱
elr sandbox create --container container-name

# 启动沙箱
elr sandbox start --sandbox-id sandbox-id

# 停止沙箱
elr sandbox stop --sandbox-id sandbox-id

# 删除沙箱
elr sandbox delete --sandbox-id sandbox-id

# 加载模型到沙箱
elr sandbox load-model --sandbox-id sandbox-id --model-id model-id

# 从沙箱卸载模型
elr sandbox unload-model --sandbox-id sandbox-id --model-id model-id

# 在沙箱中运行模型
elr sandbox run-model --sandbox-id sandbox-id --model-id model-id --input input-text
```

#### 文件系统管理

```bash
# 上传文件到容器
elr fs upload --local-path local-file --container-path container-path

# 从容器下载文件
elr fs download --container-path container-path --local-path local-file

# 设置文件类型目录
elr fs set-dir --file-type file-type --directory directory-path

# 获取文件类型目录
elr fs get-dir --file-type file-type
```

#### 运行应用

```bash
# 运行C++应用
elr run --name cpp-app --language cpp --command "./app"

# 运行Python应用
elr run --name python-app --language python --command "python app.py"

# 运行Node.js应用
elr run --name nodejs-app --language nodejs --command "node app.js"

# 运行Java应用
elr run --name java-app --language java --command "java -jar app.jar"

# 运行Go应用
elr run --name go-app --language go --command "./app"
```

## PowerShell 实现（零依赖版本）

ELR提供了一个纯PowerShell实现的版本，无需任何外部依赖即可在Windows系统上运行。这个版本是X54先生与代码织梦者通过对话式协作开发的，特别适合：

- 快速原型开发和测试
- 无法安装Go语言环境的场景
- 需要立即使用ELR功能的用户
- 学习和理解ELR架构的开发者

### 特点

- **零依赖**：无需安装Go、Docker或其他任何外部工具
- **即开即用**：下载后直接运行，无需编译
- **完整功能**：支持所有核心容器管理功能
- **状态持久化**：自动保存运行时状态，支持多命令会话

### 使用方法

1. **下载脚本**

```bash
git clone https://github.com/Wade1981/Meta-CreationPower.git
cd Meta-CreationPower/05_Open_source_ProjectRepository/AIAgentFramework/EnlightenmentLighthouseRuntime
```

2. **运行命令**

```bash
# 查看版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 version

# 启动运行时
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 start

# 列出容器
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 list

# 检查状态
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 status

# 停止运行时
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 stop
```

3. **创建和管理容器**

```bash
# 创建容器
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 create my-container

# 启动容器
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 start-container --id container-id

# 停止容器
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 stop-container --id container-id

# 检查容器
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 inspect --id container-id

# 删除容器
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 delete --id container-id
```

## 托盘应用

ELR提供了一个系统托盘应用，方便用户管理ELR服务和容器。

### 启动托盘应用

```bash
# 使用PowerShell脚本启动
powershell -ExecutionPolicy RemoteSigned -File ELR-Tray-App.ps1

# 或使用VBS脚本后台启动
wscript ELR-Tray-App.vbs
```

### 托盘应用功能

- **打开 ELR 容器设置**：配置API地址和端口，启动/停止API服务
- **容器管理**：管理ELR容器（开发中）
- **域名管理**：管理ELR域名（开发中）
- **打开终端**：打开ELR命令行终端
- **退出**：退出托盘应用

## 配置

ELR使用YAML格式的配置文件，默认路径为 `~/.elr/config.yaml`。

```yaml
# 基本配置
runtime:
  log_level: info
  data_dir: ~/.elr/data
  plugin_dir: ~/.elr/plugins

# 平台配置
platform:
  linux:
    use_namespaces: true
    use_cgroups: true
  windows:
    use_job_objects: true
    use_wsl: false
  darwin:
    use_sandbox: true
    use_spctl: true

# 网络配置
network:
  enable: true
  bridge: elr0
  subnet: 172.16.0.0/16

# 存储配置
storage:
  enable: true
  driver: overlay
  base_dir: ~/.elr/storage

# 语言配置
languages:
  cpp:
    enable: true
    runtime: /usr/bin/gcc
  python:
    enable: true
    runtime: /usr/bin/python3
  nodejs:
    enable: true
    runtime: /usr/bin/node
  java:
    enable: true
    runtime: /usr/bin/java
  go:
    enable: true
    runtime: /usr/bin/go
```

## 多语言支持

### 语言插件

ELR使用插件架构支持多种编程语言。每种语言都有一个对应的插件，负责语言运行时的管理和隔离。

#### 内置语言插件

- **cpp**：C/C++语言支持
- **python**：Python语言支持
- **nodejs**：JavaScript/Node.js语言支持
- **java**：Java语言支持
- **go**：Go语言支持

#### 自定义语言插件

您可以通过创建自定义语言插件来支持其他编程语言。语言插件需要实现以下接口：

```go
// LanguagePlugin 语言插件接口
type LanguagePlugin interface {
    // Name 返回语言名称
    Name() string
    
    // Version 返回语言版本
    Version() string
    
    // Validate 验证语言环境
    Validate() error
    
    // CreateEnvironment 创建语言环境
    CreateEnvironment(config map[string]interface{}) (Environment, error)
    
    // DestroyEnvironment 销毁语言环境
    DestroyEnvironment(env Environment) error
}

// Environment 语言环境接口
type Environment interface {
    // Run 运行命令
    Run(cmd string, args []string, env map[string]string) error
    
    // Exec 执行命令并返回输出
    Exec(cmd string, args []string, env map[string]string) (string, error)
    
    // Path 返回环境路径
    Path() string
    
    // Close 关闭环境
    Close() error
}
```

## 开发指南

### 项目结构

```
EnlightenmentLighthouseRuntime/
├── elr/               # 核心运行时
│   ├── runtime.go     # 运行时核心
│   ├── container.go   # 容器管理
│   ├── network.go     # 网络管理
│   └── storage.go     # 存储管理
├── cli/               # 命令行工具
│   ├── main.go        # 主入口
│   ├── command.go     # 命令处理
│   └── subcommand/    # 子命令
├── api/               # API服务
│   ├── api_manager.go # API管理器
│   ├── desktop_api.go # 桌面API
│   ├── model_api.go   # 模型API
│   └── public_api.go  # 公众API
├── plugins/           # 插件目录
│   ├── languages/     # 语言插件
│   └── services/      # 服务插件
├── platforms/         # 平台特定实现
│   ├── linux/         # Linux实现
│   ├── windows/       # Windows实现
│   └── darwin/        # macOS实现
├── examples/          # 示例
└── scripts/           # 辅助脚本
```

### 核心概念

#### 容器（Container）

容器是ELR的基本运行单位，包含一个或多个进程，运行在隔离的环境中。

#### 镜像（Image）

镜像是容器的模板，包含了运行容器所需的文件系统和配置。

#### 环境（Environment）

环境是语言运行时的隔离实例，用于运行特定语言的应用。

#### 插件（Plugin）

插件是ELR的扩展机制，用于添加新功能和支持新语言。

### 开发流程

1. **克隆仓库**

```bash
git clone https://github.com/Wade1981/Meta-CreationPower.git
cd Meta-CreationPower/05_Open_source_ProjectRepository/AIAgentFramework/EnlightenmentLighthouseRuntime
```

2. **安装依赖**

```bash
go mod tidy
```

3. **构建**

```bash
go build -o elr ./cli
```

4. **测试**

```bash
go test ./...
```

5. **运行**

```bash
./elr help
```

### 贡献指南

我们欢迎社区贡献，包括但不限于：

- **代码贡献**：修复bug、添加新功能
- **文档改进**：完善文档和示例
- **测试覆盖**：添加测试用例
- **问题反馈**：报告bug和提出建议

### 贡献流程

1. **Fork 仓库**
2. **创建分支**：`git checkout -b feature/your-feature`
3. **提交更改**：