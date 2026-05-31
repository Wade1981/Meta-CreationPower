# ELR 容器级 API 访问设计方案

**日期：** 2026年5月20日
**作者：** 代码织梦者 & X54先生
**项目：** Enlightenment Lighthouse Runtime (ELR)

---

## 一、问题分析

### 1.1 当前架构问题

| 问题点 | 描述 | 影响 |
|-------|------|------|
| **模型与容器解耦** | Model API 是全局的，模型未绑定到容器/沙箱 | 模型在 ELR 外运行，无法利用容器隔离特性 |
| **资源浪费** | 模型独立运行，未利用容器的资源配置 | CPU、内存限制无法生效 |
| **安全性缺失** | 模型 API 缺乏容器级访问控制 | 无法实现细粒度权限管理 |
| **监控盲区** | 模型运行状态无法关联到容器 | 资源监控、告警无法精准定位 |

### 1.2 用户需求

1. **模型必须装载到容器沙箱**：模型配置/上传后，必须指定运行的目标容器
2. **容器级 API 配置**：每个容器可以独立配置 API 访问方式
3. **隔离运行**：模型在容器沙箱内运行，享受隔离保护
4. **统一管理**：API 访问与容器生命周期绑定

---

## 二、设计原则

| 原则 | 说明 |
|-----|------|
| **容器作为唯一入口** | 模型和项目的所有访问必须先经过容器，容器是资源访问的唯一边界 |
| **容器绑定** | 模型必须绑定到特定容器才能运行，未绑定的模型无法被访问 |
| **API 隔离** | 每个容器的 API 独立配置、独立端口，互不干扰 |
| **资源受控** | 模型运行受容器资源限制（CPU、内存、磁盘） |
| **安全边界** | 容器级访问控制，支持令牌-端点类型匹配 |
| **可观测性** | 模型运行状态与容器监控统一 |

---

## 三、架构设计

### 3.1 整体架构

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          ELR Runtime                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    │
│  │   Container A   │    │   Container B   │    │   Container C   │    │
│  │  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │    │
│  │  │  Sandbox  │  │    │  │  Sandbox  │  │    │  │  Sandbox  │  │    │
│  │  │  ┌─────┐  │  │    │  │  ┌─────┐  │  │    │  │  ┌─────┐  │  │    │
│  │  │  │Model │  │  │    │  │  │Model │  │  │    │  │  │Model │  │  │    │
│  │  │  │ API  │  │  │    │  │  │ API  │  │  │    │  │  │ API  │  │  │    │
│  │  │  └─────┘  │  │    │  │  └─────┘  │  │    │  │  └─────┘  │  │    │
│  │  │  Port: 8001│  │    │  │  Port: 8002│  │    │  │  Port: 8003│  │    │
│  │  └───────────┘  │    │  └───────────┘  │    │  └───────────┘  │    │
│  │  API Config:    │    │  API Config:    │    │  API Config:    │    │
│  │  - Auth: Token  │    │  - Auth: None   │    │  - Auth: Token  │    │
│  │  - CORS: *      │    │  - CORS: Local  │    │  - CORS: Domain│    │
│  └─────────────────┘    └─────────────────┘    └─────────────────┘    │
│         │                       │                       │              │
│         └───────────────────────┼───────────────────────┘              │
│                                 ▼                                      │
│                   ┌───────────────────────┐                            │
│                   │     API Gateway       │                            │
│                   │  (请求路由、负载均衡)  │                            │
│                   └───────────────────────┘                            │
│                                 │                                      │
│                                 ▼                                      │
│                         ┌───────────┐                                  │
│                         │  Clients  │                                  │
│                         └───────────┘                                  │
└─────────────────────────────────────────────────────────────────────────┘
```

### 3.2 核心组件

| 组件 | 职责 |
|-----|------|
| **Container API Manager** | 管理每个容器的 API 配置和生命周期 |
| **API Gateway** | 统一入口，路由请求到对应容器 |
| **Sandbox Runtime** | 沙箱内模型运行时环境 |
| **Token Validator** | 容器级令牌验证 |

---

## 四、数据模型设计

### 4.1 容器 API 配置表 (container_api_config)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | TEXT | PRIMARY KEY | 配置唯一标识 |
| `container_id` | TEXT | FOREIGN KEY, UNIQUE | 关联容器ID |
| `enabled` | INTEGER | DEFAULT 1 | 是否启用 API |
| `port` | INTEGER | UNIQUE | 绑定端口（如 8001, 8002...） |
| `auth_type` | TEXT | NOT NULL | 认证类型：none/token/password |
| `allowed_tokens` | TEXT | | JSON格式允许的令牌ID列表 |
| `cors_enabled` | INTEGER | DEFAULT 0 | 是否启用 CORS |
| `cors_origin` | TEXT | | 允许的来源域名（逗号分隔或 *） |
| `rate_limit` | INTEGER | DEFAULT 100 | 每分钟请求限制 |
| `request_timeout` | INTEGER | DEFAULT 30 | 请求超时时间（秒） |
| `max_request_size` | INTEGER | DEFAULT 10485760 | 最大请求大小（字节） |
| `api_prefix` | TEXT | DEFAULT "/api" | API 路径前缀 |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `updated_at` | TIMESTAMP | | 更新时间 |

### 4.2 模型-沙箱绑定表 (model_sandbox_binding)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | INTEGER | PRIMARY KEY AUTOINCREMENT | 记录ID |
| `model_id` | TEXT | FOREIGN KEY | 关联模型ID |
| `sandbox_id` | TEXT | FOREIGN KEY | 关联沙箱ID |
| `container_id` | TEXT | FOREIGN KEY | 关联容器ID |
| `status` | TEXT | NOT NULL | 状态：loaded/unloaded/error |
| `load_time` | TIMESTAMP | | 装载时间 |
| `unload_time` | TIMESTAMP | | 卸载时间 |
| `metadata` | TEXT | | JSON格式元数据 |

### 4.3 容器 API 端点表 (container_api_endpoints)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | TEXT | PRIMARY KEY | 端点唯一标识 |
| `config_id` | TEXT | FOREIGN KEY | 关联配置ID |
| `endpoint` | TEXT | NOT NULL | 端点路径（如 /v1/completions） |
| `method` | TEXT | NOT NULL | HTTP方法：GET/POST/PUT/DELETE |
| `enabled` | INTEGER | DEFAULT 1 | 是否启用 |
| `permissions` | TEXT | | JSON格式权限列表 |
| `handler` | TEXT | NOT NULL | 处理程序名称 |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

### 4.4 API 请求日志表 (api_request_logs)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | INTEGER | PRIMARY KEY AUTOINCREMENT | 记录ID |
| `config_id` | TEXT | FOREIGN KEY | 关联配置ID |
| `container_id` | TEXT | FOREIGN KEY | 关联容器ID |
| `sandbox_id` | TEXT | FOREIGN KEY | 关联沙箱ID |
| `model_id` | TEXT | FOREIGN KEY | 关联模型ID |
| `endpoint` | TEXT | NOT NULL | 请求端点 |
| `method` | TEXT | NOT NULL | HTTP方法 |
| `client_ip` | TEXT | NOT NULL | 客户端IP |
| `token_id` | TEXT | FOREIGN KEY | 关联令牌ID |
| `status_code` | INTEGER | | HTTP状态码 |
| `response_time` | REAL | | 响应时间（毫秒） |
| `request_size` | INTEGER | | 请求大小（字节） |
| `response_size` | INTEGER | | 响应大小（字节） |
| `error_message` | TEXT | | 错误信息 |
| `timestamp` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 请求时间 |

---

## 五、API 架构设计

### 5.1 容器级 API 端点

| 端点 | 方法 | 功能 |
|-----|------|------|
| `/{container-id}/api/v1/completions` | POST | 文本补全 |
| `/{container-id}/api/v1/chat/completions` | POST | 聊天补全 |
| `/{container-id}/api/v1/embeddings` | POST | 向量嵌入 |
| `/{container-id}/api/v1/models` | GET | 列出已装载模型 |
| `/{container-id}/api/v1/models/{model-id}` | GET | 获取模型详情 |
| `/{container-id}/api/v1/models/{model-id}/load` | POST | 装载模型 |
| `/{container-id}/api/v1/models/{model-id}/unload` | POST | 卸载模型 |
| `/{container-id}/api/v1/status` | GET | 获取容器 API 状态 |
| `/{container-id}/api/v1/metrics` | GET | 获取性能指标 |

### 5.2 管理 API 端点

| 端点 | 方法 | 功能 |
|-----|------|------|
| `/api/container/{id}/api/config` | GET | 获取容器 API 配置 |
| `/api/container/{id}/api/config` | PUT | 更新容器 API 配置 |
| `/api/container/{id}/api/config` | DELETE | 删除容器 API 配置 |
| `/api/container/{id}/api/config` | POST | 创建容器 API 配置 |
| `/api/container/{id}/api/start` | POST | 启动容器 API |
| `/api/container/{id}/api/stop` | POST | 停止容器 API |
| `/api/container/{id}/api/logs` | GET | 获取 API 日志 |

---

## 六、认证与授权设计

### 6.1 认证方式

| 认证类型 | 适用场景 | 说明 |
|---------|---------|------|
| **None** | 开发环境、内网访问 | 无需认证 |
| **Token** | 生产环境、API 访问 | 令牌验证，支持端点类型匹配 |
| **Password** | 管理界面、管理员访问 | 用户名+密码 |

### 6.2 令牌-端点类型匹配

| 令牌类型 | 允许访问的容器 API |
|---------|-------------------|
| `cli` | 仅限 CLI 命令调用 |
| `api` | 仅限 API 调用 |
| `desktop` | 仅限桌面客户端 |
| `web` | 仅限 Web 浏览器 |
| `system` | 所有端点 |

### 6.3 权限分级

| 权限级别 | 说明 |
|---------|------|
| `read` | 仅查看（获取状态、列出模型） |
| `write` | 读写（装载/卸载模型） |
| `execute` | 执行（调用模型 API） |
| `admin` | 管理（配置 API、查看日志） |

---

## 七、流程设计

### 7.1 模型装载流程（必须经过容器）

```
1. 用户上传模型 → 模型存储到本地（仅元数据，不可访问）
2. 用户指定目标容器 → 验证容器存在且运行中
3. 创建 model_sandbox_binding 记录 → 建立模型与容器的绑定关系
4. 在容器沙箱内启动模型进程 → 模型进入容器隔离环境
5. 更新绑定状态为 loaded → 模型变为可访问状态
6. 容器 API 自动注册模型端点 → 仅通过容器 API 可访问
```

**关键点**：未绑定到容器的模型无法被访问，容器是模型访问的唯一入口。

### 7.2 API 请求流程（必须经过容器）

```
1. Client → API Gateway → 解析容器 ID（必须提供）
2. 验证容器状态 → 容器必须处于 running 状态
3. 查找容器 API 配置 → 验证端口、认证方式
4. 验证令牌（如启用）→ 检查令牌类型匹配、权限验证
5. 路由请求到容器沙箱 → 沙箱内模型处理（享受容器资源限制）
6. 返回响应 → 记录请求日志（关联容器ID）
```

**安全检查点**：
- 无容器 ID → 拒绝请求
- 容器未运行 → 拒绝请求  
- 容器 API 未启用 → 拒绝请求
- 令牌不匹配 → 拒绝请求

### 7.3 项目部署流程（必须经过容器）

```
1. 用户创建项目 → 项目存储到本地（仅元数据）
2. 用户指定目标容器 → 验证容器存在且运行中
3. 创建项目部署记录 → 关联项目与容器
4. 在容器沙箱内部署项目 → 项目进入容器隔离环境
5. 更新部署状态为 deployed → 项目变为可访问状态
6. 容器 API 注册项目端点 → 仅通过容器 API 可访问
```

### 7.4 API 配置流程

```
1. 创建容器时 → 自动创建默认 API 配置（禁用状态）
2. 用户启用 API → 设置端口、认证方式
3. 启动容器时 → 启动容器内 API 服务
4. 停止容器时 → 停止容器内 API 服务
5. 删除容器时 → 删除 API 配置和相关数据
```

---

## 八、资源管理

### 8.1 端口分配策略

| 端口范围 | 用途 |
|---------|------|
| 8001-8100 | 容器 API 端口（自动分配） |
| 8101-8200 | 预留扩展 |

### 8.2 资源限制

| 资源类型 | 限制方式 |
|---------|---------|
| **CPU** | 继承容器 CPU 限制 |
| **内存** | 继承容器内存限制 |
| **并发连接** | 可配置（默认 100） |
| **请求速率** | 可配置（默认 100/分钟） |

---

## 九、监控与日志

### 9.1 API 指标

| 指标 | 说明 |
|-----|------|
| `requests_total` | 总请求数 |
| `requests_success` | 成功请求数 |
| `requests_failed` | 失败请求数 |
| `response_time_avg` | 平均响应时间 |
| `response_time_p95` | 95% 响应时间 |
| `active_connections` | 活跃连接数 |
| `models_loaded` | 已装载模型数 |

### 9.2 日志记录

| 日志类型 | 存储位置 | 保留周期 |
|---------|---------|---------|
| 请求日志 | `elr_logs.db` | 30天 |
| 错误日志 | `elr_logs.db` | 90天 |
| 访问日志 | `access_logs` | 30天 |

---

## 十、安全设计

### 10.1 安全边界

| 安全层 | 措施 |
|-------|------|
| **网络隔离** | 容器独立网络命名空间 |
| **进程隔离** | 沙箱进程隔离 |
| **文件系统隔离** | 容器独立文件系统 |
| **API 认证** | 令牌验证、端点类型匹配 |
| **访问控制** | 容器级权限配置 |

### 10.2 数据保护

- API 请求体加密传输（HTTPS）
- 敏感日志脱敏处理
- 令牌使用完毕后自动清理
- 请求日志定期清理

---

## 十一、总结

### 11.1 设计亮点

1. **容器作为唯一入口**：模型和项目的所有访问必须先经过容器，未绑定容器的模型/项目无法被访问
2. **容器绑定**：模型必须装载到容器沙箱才能运行，容器是资源访问的唯一边界
3. **API 隔离**：每个容器独立配置、独立端口，互不干扰
4. **资源受控**：模型运行受容器资源限制（CPU、内存、磁盘）
5. **安全边界**：容器级访问控制，支持令牌-端点类型匹配，多重安全检查点
6. **可观测性**：统一的监控指标和日志记录，所有请求关联容器ID

### 11.2 与现有设计的集成

| 现有组件 | 集成方式 |
|---------|---------|
| **容器管理** | 容器创建时自动创建 API 配置 |
| **沙箱管理** | 沙箱启动时启动 API 服务 |
| **模型管理** | 模型装载到沙箱后注册 API 端点 |
| **令牌管理** | 令牌类型匹配验证 |
| **监控系统** | 收集容器 API 指标 |

### 11.3 后续工作

1. 实现容器 API 配置管理
2. 实现模型-沙箱绑定机制
3. 实现 API Gateway 路由
4. 实现容器级令牌验证
5. 实现监控指标收集
6. 编写单元测试和集成测试