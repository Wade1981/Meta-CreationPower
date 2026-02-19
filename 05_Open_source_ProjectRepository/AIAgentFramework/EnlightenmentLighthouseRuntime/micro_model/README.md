# ELR微模型运行沙箱

## 概述

ELR微模型运行沙箱是一个为ELR（Enlightenment Lighthouse Runtime）构建的微模型运行环境，允许在容器中独立运行各种微模型，提供隔离、安全、可配置的运行环境。

## 功能特性

- **容器隔离**：每个微模型在独立的容器中运行，确保隔离性和安全性
- **配置管理**：支持微模型的配置管理，包括模型参数、资源限制等
- **模型管理**：支持模型的下载、安装、更新和删除
- **运行监控**：监控微模型的运行状态、资源使用情况
- **API接口**：提供RESTful API接口，方便与其他系统集成
- **多模型支持**：支持多种类型的微模型，如文本生成、图像处理、语音识别等

## 架构设计

### 核心组件

1. **模型管理器**：负责模型的下载、安装、更新和删除
2. **容器管理器**：负责容器的创建、启动、停止和管理
3. **配置管理器**：负责模型配置和容器配置的管理
4. **监控服务**：监控模型运行状态和资源使用情况
5. **API服务**：提供RESTful API接口
6. **沙箱运行时**：提供模型运行的沙箱环境

### 目录结构

```
micro_model/
├── api/            # API接口实现
├── config/         # 配置管理
├── container/      # 容器管理
├── model/          # 模型管理
├── monitor/        # 监控服务
├── sandbox/        # 沙箱运行时
├── examples/       # 示例配置和模型
├── scripts/        # 辅助脚本
├── main.go         # 主入口
├── go.mod          # Go模块依赖
└── README.md       # 文档
```

## 快速开始

### 环境要求

- Go 1.16+
- Docker（用于容器管理）
- Python 3.8+（用于模型运行）
- 足够的存储空间（用于存储模型）

### 安装

```bash
# 克隆代码
git clone <repository-url>
cd micro_model

# 构建
go build -o micro_model_server .

# 运行
./micro_model_server
```

### 配置

配置文件位于 `config/config.yaml`，主要配置项包括：

- **server**：服务器配置，如端口、主机等
- **container**：容器配置，如基础镜像、资源限制等
- **model**：模型配置，如模型存储路径、默认模型等
- **monitoring**：监控配置，如监控间隔、告警阈值等

### 使用示例

#### 1. 启动微模型服务

```bash
./micro_model_server
```

#### 2. 下载模型

```bash
curl -X POST http://localhost:8080/api/models/download \
  -H "Content-Type: application/json" \
  -d '{"model_id": "gpt2", "model_type": "text-generation"}'
```

#### 3. 创建模型容器

```bash
curl -X POST http://localhost:8080/api/containers/create \
  -H "Content-Type: application/json" \
  -d '{"model_id": "gpt2", "container_name": "gpt2-container", "resources": {"cpu": 1, "memory": "1G"}}'
```

#### 4. 启动模型容器

```bash
curl -X POST http://localhost:8080/api/containers/start \
  -H "Content-Type: application/json" \
  -d '{"container_name": "gpt2-container"}'
```

#### 5. 使用模型

```bash
curl -X POST http://localhost:8080/api/models/run \
  -H "Content-Type: application/json" \
  -d '{"container_name": "gpt2-container", "input": "Hello, world!"}'
```

## 模型支持

### 支持的模型类型

- **文本生成**：GPT系列、BERT系列等
- **图像处理**：ResNet、VGG、YOLO等
- **语音识别**：Whisper、DeepSpeech等
- **多模态**：CLIP、DALL-E等

### 模型存储

模型存储在 `model/models/` 目录下，每个模型有独立的子目录。

## 安全考虑

- **容器隔离**：使用Docker容器隔离模型运行环境
- **资源限制**：限制容器的CPU、内存使用
- **网络隔离**：可配置容器的网络访问权限
- **输入验证**：验证模型输入，防止恶意输入
- **权限控制**：API接口的权限控制

## 监控和日志

- **运行状态**：监控模型的运行状态、响应时间等
- **资源使用**：监控CPU、内存、磁盘使用情况
- **日志管理**：收集和管理模型运行日志
- **告警机制**：当模型运行异常时发送告警

## 扩展和定制

### 自定义模型

可以通过以下步骤添加自定义模型：

1. 在 `model/models/` 目录下创建模型目录
2. 提供模型文件和配置文件
3. 更新模型配置

### 自定义沙箱

可以通过修改 `sandbox/` 目录下的代码，定制沙箱运行时。

## 故障排除

### 常见问题

1. **模型下载失败**：检查网络连接和模型源
2. **容器启动失败**：检查Docker配置和资源限制
3. **模型运行异常**：检查模型配置和输入格式
4. **API调用失败**：检查API接口和参数

### 日志查看

模型运行日志位于 `logs/` 目录下，按容器名称和时间戳命名。

## 未来规划

- **模型自动缩放**：根据负载自动调整模型实例数量
- **模型版本管理**：支持模型的版本控制和回滚
- **模型评估**：提供模型性能评估工具
- **多节点支持**：支持分布式部署和负载均衡
- **模型市场**：提供模型分享和交易平台

## 贡献指南

欢迎贡献代码、文档和建议。请参考项目的贡献指南。

## 许可证

本项目采用 MIT 许可证。
