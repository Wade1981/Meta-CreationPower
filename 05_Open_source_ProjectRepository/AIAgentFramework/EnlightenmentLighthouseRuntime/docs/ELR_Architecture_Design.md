# ELR (Enlightenment Lighthouse Runtime) 架构设计文档

## 1. 架构概述

ELR是一个轻量级的容器运行时，基于"和清寂静"核心原则，提供安全、高效的容器管理和模型运行环境。ELR采用模块化设计，分为核心运行时、容器管理、网络服务和模型管理等核心组件。

### 1.1 设计原则

- **模块化设计**：各个组件独立封装，便于维护和扩展
- **跨平台兼容**：支持Windows、Linux和macOS平台
- **安全优先**：内置网络隔离和安全策略
- **性能优化**：并行处理和资源管理
- **易于使用**：提供简洁的命令行界面和API

## 2. 核心组件

### 2.1 运行时管理 (Runtime)

运行时是ELR的核心组件，负责初始化和管理整个系统。

**主要功能**：
- 初始化平台环境
- 管理容器生命周期
- 加载和管理插件
- 启动网络服务
- 资源管理和监控

**核心文件**：
- `elr/runtime.go`：运行时核心实现
- `elr/platforms/`：平台特定实现

### 2.2 容器管理 (Container)

容器管理组件负责容器的创建、启动、停止、删除等操作。

**主要功能**：
- 容器生命周期管理
- 资源限制和监控
- 文件系统隔离
- 网络隔离
- 容器状态管理

**核心文件**：
- `elr/container.go`：容器核心实现
- `elr/platforms/windows/windows.go`：Windows平台容器实现

### 2.3 网络服务 (Network)

网络服务组件提供API服务和网络隔离功能。

**主要功能**：
- RESTful API服务
- 网络隔离和安全策略
- CORS和速率限制
- 容器网络管理

**核心文件**：
- `elr/network.go`：网络服务核心实现
- `api/`：API实现

### 2.4 模型管理 (Model)

模型管理组件负责模型的加载、切换和管理。

**主要功能**：
- 模型动态加载和卸载
- 模型切换和状态管理
- 资源使用监控
- 模型依赖管理

**核心文件**：
- `micro_model/model/model.go`：模型管理核心实现
- `micro_model/model/model_adapter.go`：模型适配器

### 2.5 安全管理 (Security)

安全管理组件提供安全策略和访问控制。

**主要功能**：
- CORS策略
- 速率限制
- 令牌管理
- 权限控制

**核心文件**：
- `elr/network.go`：安全管理器实现
- `elr/token_manager.go`：令牌管理

## 3. 架构流程图

```
┌─────────────────────────────────────────────────────────────┐
│                     ELR Runtime                           │
├─────────────┬─────────────┬─────────────┬──────────────┐
│             │             │             │              │
▼             ▼             ▼             ▼              ▼
┌─────────┐ ┌─────────┐ ┌─────────┐ ┌──────────┐ ┌────────┐
│Container│ │ Network │ │  Model  │ │ Platform │ │ Security│
│ 管理    │ │ 服务    │ │ 管理    │ │ 抽象层  │ │ 管理    │
└─────────┘ └─────────┘ └─────────┘ └──────────┘ └────────┘
```

## 4. 核心API

### 4.1 容器管理API

| API方法 | 功能 | 参数 | 返回值 |
|--------|------|------|--------|
| `CreateContainer` | 创建容器 | ContainerConfig | Container, error |
| `GetContainer` | 获取容器 | containerID string | Container, error |
| `ListContainers` | 列出容器 | 无 | []*Container |
| `DeleteContainer` | 删除容器 | containerID string | error |
| `Container.Start` | 启动容器 | 无 | error |
| `Container.Stop` | 停止容器 | 无 | error |
| `Container.Pause` | 暂停容器 | 无 | error |
| `Container.Unpause` | 恢复容器 | 无 | error |
| `Container.UploadFile` | 上传文件 | localPath, containerPath, token string | error |
| `Container.DownloadFile` | 下载文件 | containerPath, localPath, token string | error |

### 4.2 网络服务API

| API路径 | 方法 | 功能 | 权限 |
|--------|------|------|------|
| `/health` | GET | 健康检查 | 公开 |
| `/api/container/list` | GET | 列出容器 | 公开 |
| `/api/container/status` | GET | 获取容器状态 | 公开 |
| `/api/model/run` | POST | 运行模型 | 公开 |
| `/api/model/list` | GET | 列出模型 | 公开 |
| `/api/network/status` | GET | 获取网络状态 | 公开 |
| `/api/network/isolate` | POST | 隔离容器网络 | 公开 |
| `/api/network/unisolate` | POST | 取消网络隔离 | 公开 |
| `/api/network/config` | GET | 获取网络配置 | 公开 |
| `/api/token/create` | POST | 创建令牌 | 公开 |
| `/api/token/validate` | POST | 验证令牌 | 公开 |
| `/api/token/refresh` | POST | 刷新令牌 | 公开 |
| `/api/token/list` | GET | 列出令牌 | 公开 |
| `/api/token/revoke` | POST | 撤销令牌 | 公开 |
| `/api/desktop/health` | GET | 桌面API健康检查 | 公开 |
| `/api/desktop/status` | GET | 获取ELR状态 | 公开 |
| `/api/desktop/containers` | GET | 获取容器列表 | 公开 |
| `/api/desktop/resources` | GET | 获取系统资源 | 公开 |
| `/api/desktop/files` | GET | 列出上传文件 | 公开 |
| `/api/desktop/upload` | POST | 上传文件 | 公开 |

### 4.3 模型管理API

| API方法 | 功能 | 参数 | 返回值 |
|--------|------|------|--------|
| `GetModel` | 获取模型信息 | modelID string | *Model, error |
| `ListModels` | 列出所有模型 | 无 | []*Model, error |
| `DownloadModel` | 下载模型 | modelID, modelType, downloadURL string | error |
| `DeleteModel` | 删除模型 | modelID string | error |
| `UpdateModel` | 更新模型 | modelID, downloadURL string | error |
| `LoadModel` | 加载模型 | modelID string | error |
| `UnloadModel` | 卸载模型 | modelID string | error |
| `SwitchModel` | 切换模型 | modelID string | error |
| `GetLoadedModels` | 获取已加载模型 | 无 | map[string]*LoadedModel |
| `GetActiveModel` | 获取活动模型 | 无 | *LoadedModel |

## 5. 配置系统

### 5.1 配置结构

ELR使用YAML配置文件，主要配置项包括：

- **数据目录**：存储容器和模型数据
- **插件目录**：存储插件
- **平台配置**：平台特定设置
- **网络配置**：网络服务和API端口
- **存储配置**：存储驱动和目录
- **语言配置**：支持的编程语言
- **资源配置**：资源类型和目录

### 5.2 默认配置

```yaml
log_level: info
data_dir: ~/.elr/data
plugin_dir: ~/.elr/plugins
platform:
  windows:
    use_job_objects: true
    use_wsl: false
    use_containers: false
    isolation_type: basic
network:
  enable: false
  bridge: elr0
  subnet: 172.16.0.0/16
  api_ports:
    desktop_api: 8081
    public_api: 8080
    model_api: 8082
storage:
  enable: true
  driver: local
  base_dir: ~/.elr/storage
languages:
  python:
    enable: true
    runtime: python3
resources:
  types:
    models:
      enable: true
      dir: ~/.elr/models
```

## 6. 平台支持

### 6.1 Windows平台

- **隔离方式**：Job Objects, WSL, Windows Containers, AppContainers
- **文件系统**：NTFS, ReFS
- **网络**：Windows Firewall, Hyper-V Network

### 6.2 Linux平台

- **隔离方式**：Namespaces, Cgroups
- **文件系统**：OverlayFS, Btrfs
- **网络**：Docker Network, CNI

### 6.3 macOS平台

- **隔离方式**：Sandbox, spctl
- **文件系统**：APFS
- **网络**：pf, ipfw

## 7. 安全策略

### 7.1 网络安全

- **CORS策略**：允许跨域请求
- **速率限制**：防止API滥用
- **网络隔离**：容器网络隔离
- **防火墙规则**：限制容器网络访问

### 7.2 容器安全

- **文件系统隔离**：只读文件系统选项
- **资源限制**：CPU和内存限制
- **权限控制**：基于令牌的访问控制
- **加密存储**：文件加密

### 7.3 模型安全

- **模型隔离**：模型运行在隔离环境
- **依赖管理**：安全的依赖安装
- **访问控制**：模型访问权限

## 8. 性能优化

### 8.1 容器启动优化

- **并行初始化**：使用goroutines并行初始化容器组件
- **资源预分配**：提前分配资源
- **缓存机制**：缓存容器状态

### 8.2 网络服务优化

- **并发处理**：支持并发API请求
- **连接池**：复用HTTP连接
- **负载均衡**：分布式服务支持

### 8.3 模型管理优化

- **模型缓存**：缓存常用模型
- **资源监控**：实时监控模型资源使用
- **动态加载**：按需加载模型

## 9. 插件系统

ELR支持插件扩展，插件可以：

- 扩展容器功能
- 添加新的隔离方式
- 集成外部服务
- 提供自定义API

## 10. 部署和运维

### 10.1 部署方式

- **本地部署**：直接运行elr.exe
- **服务部署**：作为系统服务运行
- **容器部署**：在容器中运行ELR

### 10.2 监控和日志

- **系统监控**：资源使用监控
- **容器监控**：容器状态监控
- **API监控**：API请求监控
- **日志管理**：结构化日志

### 10.3 故障处理

- **自动恢复**：容器故障自动恢复
- **错误处理**：优雅的错误处理
- **备份恢复**：配置和数据备份

## 11. 未来规划

### 11.1 功能增强

- **Kubernetes集成**：与K8s集成
- **GPU支持**：GPU加速模型运行
- **多集群管理**：管理多个ELR实例
- **AI模型市场**：模型分享和交易

### 11.2 性能优化

- **容器镜像管理**：高效镜像管理
- **网络优化**：网络性能优化
- **存储优化**：存储性能优化

### 11.3 安全增强

- **零信任架构**：零信任安全模型
- **漏洞扫描**：容器和模型漏洞扫描
- **合规性**：符合行业合规标准

## 12. 结论

ELR是一个设计先进、功能丰富的容器运行时，具有以下优势：

- **模块化设计**：易于扩展和维护
- **跨平台兼容**：支持多种操作系统
- **安全可靠**：内置安全策略
- **性能优化**：高效的资源管理
- **易于使用**：简洁的命令行界面

ELR为AI模型运行和容器管理提供了一个安全、高效的环境，适合从开发到生产的各种场景。
