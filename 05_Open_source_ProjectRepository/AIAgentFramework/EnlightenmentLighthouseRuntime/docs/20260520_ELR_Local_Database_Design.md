# ELR 本地化数据库设计方案

**日期：** 2026年5月20日
**作者：** 代码织梦者 & X54先生
**项目：** Enlightenment Lighthouse Runtime (ELR)

---

## 一、需求分析

### 1.1 背景
当前ELR使用JSON文件管理容器、沙箱、模型等数据，存在以下问题：
- 数据一致性难以保证
- 查询效率低，不支持复杂查询
- 缺乏事务支持
- 难以实现数据关联和约束
- 不便于数据备份和迁移

### 1.2 核心需求

| 需求类别 | 需求描述 | 优先级 |
|---------|---------|--------|
| 容器管理 | 存储容器的基本信息、状态、资源配置 | 高 |
| 沙箱管理 | 存储沙箱信息、关联容器、运行状态 | 高 |
| 模型管理 | 存储模型元数据、版本、路径、依赖 | 高 |
| 项目管理 | 存储项目信息、部署状态、运行记录 | 高 |
| 资源监控 | 记录容器/沙箱的资源使用情况（CPU、内存、磁盘、网络） | 高 |
| 监控告警 | 基于资源阈值的告警机制 | 中 |
| 访问日志 | 记录容器被外界访问的情况（请求次数、来源、时间） | 中 |
| 运行时数据 | 记录运行过程中产生的有价值数据 | 中 |
| 令牌管理 | 令牌生成、验证、刷新、撤销、权限分级 | 高 |
| 管理员系统 | 管理员创建、权限管理、角色控制、分组管理 | 高 |
| 安全审计 | 记录操作日志、权限变更、异常事件 | 高 |

### 1.3 设计原则

| 原则 | 说明 |
|-----|------|
| **安全性** | 数据加密存储、访问控制、审计日志、令牌认证 |
| **低依赖性** | 使用SQLite，无需额外服务，开箱即用 |
| **低性能适配** | 优化查询语句、合理索引设计、数据分片 |
| **可扩展性** | 模块化设计、预留扩展字段、支持插件化 |
| **未来智能性** | 预留AI分析字段、支持数据挖掘、可扩展为时序数据库 |

---

## 二、数据库架构设计

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                     ELR Database Layer                          │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐    │
│  │  Container  │  │   Sandbox   │  │      Model          │    │
│  │   Manager   │  │   Manager   │  │      Manager        │    │
│  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘    │
│         │                │                     │               │
│  ┌──────┴──────┐  ┌──────┴──────┐  ┌──────────┴──────────┐    │
│  │   Token     │  │  Admin      │  │    Monitor          │    │
│  │   Manager   │  │   Manager   │  │    Manager          │    │
│  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘    │
│         │                │                     │               │
├─────────┼────────────────┼─────────────────────┼───────────────┤
│         ▼                ▼                     ▼               │
│  ┌─────────────────────────────────────────────────────┐       │
│  │                    SQLite Database                  │       │
│  │  (elr_data.db - 核心数据)                           │       │
│  │  (elr_metrics.db - 监控数据)                        │       │
│  │  (elr_logs.db - 日志数据)                           │       │
│  └─────────────────────────────────────────────────────┘       │
├─────────────────────────────────────────────────────────────────┤
│                     数据访问层 (DAO)                           │
│  - CRUD操作封装                                               │
│  - 事务管理                                                   │
│  - 连接池管理                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 数据库文件规划

| 数据库文件 | 用途 | 存储策略 |
|-----------|------|---------|
| `elr_data.db` | 核心业务数据（容器、沙箱、模型、项目、用户、令牌、管理员） | 持久化，定期备份 |
| `elr_metrics.db` | 监控指标数据（资源使用、性能指标） | 按时间分片，定期清理 |
| `elr_logs.db` | 操作日志、审计日志、访问日志 | 按时间分片，定期清理 |

---

## 三、核心数据表设计

### 3.1 容器表 (containers)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | TEXT | PRIMARY KEY | 容器唯一标识 |
| `name` | TEXT | NOT NULL | 容器名称 |
| `image` | TEXT | NOT NULL | 镜像名称 |
| `status` | TEXT | NOT NULL | 状态：created/running/stopped/deleted |
| `isolation_type` | TEXT | | 隔离类型：windows-container/wsl/basic |
| `cpu_limit` | INTEGER | | CPU核心数限制 |
| `memory_limit` | INTEGER | | 内存限制（MB） |
| `disk_limit` | INTEGER | | 磁盘限制（MB） |
| `network_mode` | TEXT | | 网络模式：bridge/host/nat |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `started_at` | TIMESTAMP | | 启动时间 |
| `stopped_at` | TIMESTAMP | | 停止时间 |
| `config` | TEXT | | JSON格式配置 |
| `metadata` | TEXT | | JSON格式元数据 |

### 3.2 沙箱表 (sandboxes)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | TEXT | PRIMARY KEY | 沙箱唯一标识 |
| `name` | TEXT | NOT NULL | 沙箱名称 |
| `container_id` | TEXT | FOREIGN KEY | 关联容器ID |
| `status` | TEXT | NOT NULL | 状态：created/running/stopped |
| `process_id` | INTEGER | | 沙箱进程ID |
| `pipe_name` | TEXT | | IPC命名管道名称 |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `started_at` | TIMESTAMP | | 启动时间 |
| `uptime` | INTEGER | | 运行时长（秒） |
| `config` | TEXT | | JSON格式配置 |

### 3.3 模型表 (models)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | TEXT | PRIMARY KEY | 模型唯一标识 |
| `name` | TEXT | NOT NULL | 模型名称 |
| `type` | TEXT | NOT NULL | 模型类型：llm/embedding/image/audio |
| `version` | TEXT | NOT NULL | 版本号 |
| `path` | TEXT | NOT NULL | 本地路径 |
| `url` | TEXT | | 下载地址 |
| `size` | INTEGER | | 文件大小（字节） |
| `status` | TEXT | NOT NULL | 状态：downloading/installed/ready/error |
| `dependencies` | TEXT | | JSON格式依赖列表 |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `updated_at` | TIMESTAMP | | 更新时间 |
| `metadata` | TEXT | | JSON格式元数据（作者、描述、许可证等） |

### 3.4 项目表 (projects)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | TEXT | PRIMARY KEY | 项目唯一标识 |
| `name` | TEXT | NOT NULL | 项目名称 |
| `type` | TEXT | NOT NULL | 项目类型：nodejs/php/java/python |
| `version` | TEXT | | 版本号 |
| `description` | TEXT | | 项目描述 |
| `path` | TEXT | NOT NULL | 本地路径 |
| `status` | TEXT | NOT NULL | 状态：created/deployed/running/stopped |
| `sandbox_id` | TEXT | FOREIGN KEY | 关联沙箱ID |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `deployed_at` | TIMESTAMP | | 部署时间 |

### 3.5 资源监控表 (resource_metrics)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | INTEGER | PRIMARY KEY AUTOINCREMENT | 记录ID |
| `target_id` | TEXT | NOT NULL | 目标ID（容器/沙箱） |
| `target_type` | TEXT | NOT NULL | 目标类型：container/sandbox |
| `timestamp` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 采集时间 |
| `cpu_usage` | REAL | | CPU使用率（0-100） |
| `memory_usage` | INTEGER | | 内存使用量（MB） |
| `memory_percent` | REAL | | 内存使用率（0-100） |
| `disk_usage` | INTEGER | | 磁盘使用量（MB） |
| `disk_percent` | REAL | | 磁盘使用率（0-100） |
| `network_in` | INTEGER | | 网络入站流量（字节） |
| `network_out` | INTEGER | | 网络出站流量（字节） |
| `gpu_usage` | REAL | | GPU使用率（0-100） |
| `gpu_memory` | INTEGER | | GPU内存使用量（MB） |
| `process_count` | INTEGER | | 进程数量 |
| `thread_count` | INTEGER | | 线程数量 |

### 3.6 监控告警规则表 (alert_rules)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | TEXT | PRIMARY KEY | 规则唯一标识 |
| `name` | TEXT | NOT NULL | 规则名称 |
| `target_type` | TEXT | NOT NULL | 目标类型：container/sandbox/system |
| `metric_type` | TEXT | NOT NULL | 指标类型：cpu/memory/disk/network/gpu |
| `operator` | TEXT | NOT NULL | 比较运算符：> / < / >= / <= / == |
| `threshold` | REAL | NOT NULL | 阈值 |
| `duration` | INTEGER | DEFAULT 60 | 持续时间（秒） |
| `severity` | TEXT | NOT NULL | 严重级别：info/warning/critical |
| `enabled` | INTEGER | DEFAULT 1 | 是否启用（0/1） |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `updated_at` | TIMESTAMP | | 更新时间 |

### 3.7 监控告警记录表 (alerts)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | INTEGER | PRIMARY KEY AUTOINCREMENT | 记录ID |
| `rule_id` | TEXT | FOREIGN KEY | 关联规则ID |
| `target_id` | TEXT | NOT NULL | 目标ID |
| `target_type` | TEXT | NOT NULL | 目标类型 |
| `metric_type` | TEXT | NOT NULL | 指标类型 |
| `current_value` | REAL | NOT NULL | 当前值 |
| `threshold` | REAL | NOT NULL | 阈值 |
| `severity` | TEXT | NOT NULL | 严重级别 |
| `status` | TEXT | NOT NULL | 状态：firing/resolved |
| `message` | TEXT | | 告警消息 |
| `fired_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 触发时间 |
| `resolved_at` | TIMESTAMP | | 解决时间 |

### 3.8 访问日志表 (access_logs)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | INTEGER | PRIMARY KEY AUTOINCREMENT | 记录ID |
| `container_id` | TEXT | FOREIGN KEY | 关联容器ID |
| `client_ip` | TEXT | NOT NULL | 客户端IP |
| `request_method` | TEXT | | 请求方法：GET/POST/PUT/DELETE |
| `request_path` | TEXT | | 请求路径 |
| `status_code` | INTEGER | | HTTP状态码 |
| `response_time` | REAL | | 响应时间（毫秒） |
| `user_agent` | TEXT | | 用户代理 |
| `token_id` | TEXT | FOREIGN KEY | 关联令牌ID |
| `endpoint_type` | TEXT | | 访问端类型：cli/api/desktop/web |
| `timestamp` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 请求时间 |

### 3.9 操作审计表 (audit_logs)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | INTEGER | PRIMARY KEY AUTOINCREMENT | 记录ID |
| `user_id` | TEXT | FOREIGN KEY | 操作用户ID |
| `admin_id` | TEXT | FOREIGN KEY | 操作管理员ID |
| `operation` | TEXT | NOT NULL | 操作类型：create/update/delete/start/stop |
| `target_type` | TEXT | NOT NULL | 目标类型：container/sandbox/model/project/token/admin |
| `target_id` | TEXT | | 目标ID |
| `target_name` | TEXT | | 目标名称 |
| `details` | TEXT | | JSON格式操作详情 |
| `result` | TEXT | NOT NULL | 操作结果：success/failure |
| `error_message` | TEXT | | 错误信息 |
| `ip_address` | TEXT | | 操作IP地址 |
| `timestamp` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 操作时间 |

### 3.10 运行时数据表 (runtime_data)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | INTEGER | PRIMARY KEY AUTOINCREMENT | 记录ID |
| `sandbox_id` | TEXT | FOREIGN KEY | 关联沙箱ID |
| `model_id` | TEXT | FOREIGN KEY | 关联模型ID |
| `data_type` | TEXT | NOT NULL | 数据类型：prompt/completion/embedding/log |
| `content` | TEXT | | 数据内容 |
| `metadata` | TEXT | | JSON格式元数据 |
| `timestamp` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 产生时间 |

### 3.11 用户表 (users)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | TEXT | PRIMARY KEY | 用户唯一标识 |
| `username` | TEXT | NOT NULL UNIQUE | 用户名 |
| `password_hash` | TEXT | NOT NULL | 密码哈希 |
| `role` | TEXT | NOT NULL | 角色：admin/user/guest |
| `email` | TEXT | | 邮箱 |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `last_login_at` | TIMESTAMP | | 最后登录时间 |
| `is_active` | INTEGER | DEFAULT 1 | 是否激活（0/1） |

### 3.12 管理员表 (admins)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | TEXT | PRIMARY KEY | 管理员唯一标识 |
| `username` | TEXT | NOT NULL UNIQUE | 管理员用户名 |
| `password_hash` | TEXT | NOT NULL | 密码哈希 |
| `role` | TEXT | NOT NULL | 角色：super_admin/admin |
| `email` | TEXT | | 邮箱 |
| `group_id` | TEXT | FOREIGN KEY | 关联管理员组ID |
| `token_id` | TEXT | FOREIGN KEY | 关联令牌ID |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `last_login_at` | TIMESTAMP | | 最后登录时间 |
| `is_active` | INTEGER | DEFAULT 1 | 是否激活（0/1） |
| `metadata` | TEXT | | JSON格式元数据 |

### 3.13 管理员组表 (admin_groups)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | TEXT | PRIMARY KEY | 组唯一标识 |
| `name` | TEXT | NOT NULL UNIQUE | 组名称 |
| `description` | TEXT | | 组描述 |
| `permissions` | TEXT | NOT NULL | JSON格式权限列表 |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `updated_at` | TIMESTAMP | | 更新时间 |

### 3.14 令牌表 (tokens)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | TEXT | PRIMARY KEY | 令牌唯一标识 |
| `user_id` | TEXT | FOREIGN KEY | 关联用户ID |
| `admin_id` | TEXT | FOREIGN KEY | 关联管理员ID |
| `token_hash` | TEXT | NOT NULL UNIQUE | 令牌哈希值 |
| `token_type` | TEXT | NOT NULL | 令牌类型：cli/api/desktop/web/system |
| `endpoint_type` | TEXT | NOT NULL | 访问端类型：cli/api/desktop/web |
| `permissions` | TEXT | NOT NULL | JSON格式权限列表 |
| `expires_at` | TIMESTAMP | | 过期时间 |
| `issued_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 颁发时间 |
| `last_used_at` | TIMESTAMP | | 最后使用时间 |
| `status` | TEXT | NOT NULL | 状态：active/revoked/expired |
| `refresh_token_hash` | TEXT | | 刷新令牌哈希 |
| `usage_count` | INTEGER | DEFAULT 0 | 使用次数 |
| `max_usage` | INTEGER | | 最大使用次数（0表示无限制） |
| `allowed_ips` | TEXT | | JSON格式允许的IP列表 |
| `metadata` | TEXT | | JSON格式元数据 |

### 3.15 令牌使用统计表 (token_stats)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | INTEGER | PRIMARY KEY AUTOINCREMENT | 记录ID |
| `token_id` | TEXT | FOREIGN KEY | 关联令牌ID |
| `date` | DATE | NOT NULL | 统计日期 |
| `request_count` | INTEGER | DEFAULT 0 | 请求次数 |
| `success_count` | INTEGER | DEFAULT 0 | 成功次数 |
| `failure_count` | INTEGER | DEFAULT 0 | 失败次数 |
| `total_response_time` | REAL | DEFAULT 0 | 总响应时间（毫秒） |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

### 3.16 沙箱-容器映射表 (sandbox_container_mapping)

| 字段名 | 类型 | 约束 | 说明 |
|-------|------|------|------|
| `id` | INTEGER | PRIMARY KEY AUTOINCREMENT | 记录ID |
| `sandbox_id` | TEXT | NOT NULL UNIQUE | 沙箱ID |
| `container_id` | TEXT | NOT NULL | 容器ID |
| `created_at` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

---

## 四、索引设计

### 4.1 核心表索引

```sql
-- containers 表索引
CREATE INDEX idx_containers_status ON containers(status);
CREATE INDEX idx_containers_created_at ON containers(created_at);

-- sandboxes 表索引
CREATE INDEX idx_sandboxes_container_id ON sandboxes(container_id);
CREATE INDEX idx_sandboxes_status ON sandboxes(status);
CREATE INDEX idx_sandboxes_process_id ON sandboxes(process_id);

-- models 表索引
CREATE INDEX idx_models_type ON models(type);
CREATE INDEX idx_models_status ON models(status);

-- projects 表索引
CREATE INDEX idx_projects_sandbox_id ON projects(sandbox_id);
CREATE INDEX idx_projects_status ON projects(status);

-- admins 表索引
CREATE INDEX idx_admins_group_id ON admins(group_id);
CREATE INDEX idx_admins_role ON admins(role);

-- tokens 表索引
CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_admin_id ON tokens(admin_id);
CREATE INDEX idx_tokens_status ON tokens(status);
CREATE INDEX idx_tokens_token_type ON tokens(token_type);
CREATE INDEX idx_tokens_endpoint_type ON tokens(endpoint_type);
CREATE INDEX idx_tokens_expires_at ON tokens(expires_at);
```

### 4.2 监控数据表索引

```sql
-- resource_metrics 表索引
CREATE INDEX idx_metrics_target ON resource_metrics(target_id, target_type);
CREATE INDEX idx_metrics_timestamp ON resource_metrics(timestamp);

-- alert_rules 表索引
CREATE INDEX idx_alert_rules_target ON alert_rules(target_type, metric_type);
CREATE INDEX idx_alert_rules_enabled ON alert_rules(enabled);

-- alerts 表索引
CREATE INDEX idx_alerts_rule_id ON alerts(rule_id);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_severity ON alerts(severity);
CREATE INDEX idx_alerts_timestamp ON alerts(fired_at);

-- access_logs 表索引
CREATE INDEX idx_access_container ON access_logs(container_id);
CREATE INDEX idx_access_timestamp ON access_logs(timestamp);
CREATE INDEX idx_access_client_ip ON access_logs(client_ip);
CREATE INDEX idx_access_token_id ON access_logs(token_id);
CREATE INDEX idx_access_endpoint_type ON access_logs(endpoint_type);

-- audit_logs 表索引
CREATE INDEX idx_audit_user ON audit_logs(user_id);
CREATE INDEX idx_audit_admin ON audit_logs(admin_id);
CREATE INDEX idx_audit_target ON audit_logs(target_type, target_id);
CREATE INDEX idx_audit_timestamp ON audit_logs(timestamp);

-- token_stats 表索引
CREATE INDEX idx_token_stats_token_id ON token_stats(token_id);
CREATE INDEX idx_token_stats_date ON token_stats(date);
```

---

## 五、安全性设计

### 5.1 数据加密

| 加密层级 | 说明 | 实现方式 |
|---------|------|---------|
| **传输加密** | API通信加密 | HTTPS/TLS |
| **存储加密** | 敏感字段加密（密码、令牌） | AES-256 |
| **数据库加密** | 整个数据库文件加密 | SQLCipher |

### 5.2 访问控制

| 角色 | 权限 |
|-----|------|
| **super_admin** | 全部权限，包括管理员管理 |
| **admin** | 容器、沙箱、模型、项目管理 |
| **user** | 查看、创建、修改自己的资源 |
| **guest** | 仅查看权限 |

### 5.3 令牌-端点类型匹配

| 令牌类型 | 允许的访问端 | 说明 |
|---------|------------|------|
| `cli` | CLI | 命令行工具 |
| `api` | API | RESTful API |
| `desktop` | Desktop | 桌面客户端 |
| `web` | Web | Web浏览器 |
| `system` | 所有 | 系统内部使用 |

### 5.4 审计机制

- 所有敏感操作记录到 `audit_logs` 表
- 记录操作时间、操作者、操作对象、操作结果
- 保留90天审计日志
- 管理员操作单独记录关联管理员ID

---

## 六、性能优化策略

### 6.1 数据分片

| 数据表 | 分片策略 | 保留周期 |
|-------|---------|---------|
| `resource_metrics` | 按天分片 | 30天 |
| `access_logs` | 按天分片 | 30天 |
| `audit_logs` | 按月分片 | 90天 |
| `runtime_data` | 按周分片 | 60天 |
| `alerts` | 按月分片 | 60天 |
| `token_stats` | 按月分片 | 90天 |

### 6.2 查询优化

- 使用参数化查询避免SQL注入
- 合理设计索引覆盖常用查询
- 使用连接池管理数据库连接
- 批量操作使用事务减少IO
- 监控数据使用时间范围查询优化

---

## 七、可扩展性设计

### 7.1 预留扩展字段

- 所有核心表预留 `metadata` 字段存储JSON格式的扩展数据
- 支持动态添加自定义属性
- 支持插件化的数据模型扩展

### 7.2 模块化架构

```
┌────────────────────────────────────────────────┐
│              Database Module                   │
├────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐            │
│  │  Container  │  │    Base     │            │
│  │   DAO       │  │   DAO       │            │
│  └──────┬──────┘  └──────┬──────┘            │
│         │                │                    │
│  ┌──────┴──────┐  ┌──────┴──────┐            │
│  │   Token     │  │   Admin     │            │
│  │   DAO       │  │   DAO       │            │
│  └──────┬──────┘  └──────┬──────┘            │
│         │                │                    │
│         └────────┬───────┘                    │
│                  ▼                            │
│         ┌─────────────┐                       │
│         │  Connection │                       │
│         │    Pool     │                       │
│         └─────────────┘                       │
└────────────────────────────────────────────────┘
```

---

## 八、未来智能性支持

### 8.1 AI分析字段预留

| 字段类型 | 用途 |
|---------|------|
| `embedding_vector` | 存储内容的向量表示，支持相似度搜索 |
| `analysis_score` | AI分析评分 |
| `prediction_result` | AI预测结果 |
| `anomaly_flag` | 异常检测标记 |

### 8.2 时序数据分析

- `resource_metrics` 表支持时序查询
- 预留趋势分析字段
- 支持数据聚合和统计

---

## 九、数据库操作接口设计

### 9.1 DAO层接口

| 接口 | 功能 |
|-----|------|
| `ContainerDAO` | 容器数据访问 |
| `SandboxDAO` | 沙箱数据访问 |
| `ModelDAO` | 模型数据访问 |
| `ProjectDAO` | 项目数据访问 |
| `MetricsDAO` | 监控数据访问 |
| `AlertDAO` | 告警数据访问 |
| `AuditDAO` | 审计日志访问 |
| `TokenDAO` | 令牌数据访问 |
| `AdminDAO` | 管理员数据访问 |

### 9.2 事务管理

- 支持声明式事务
- 自动回滚机制
- 事务嵌套支持

---

## 十、部署与迁移方案

### 10.1 数据库初始化

```sql
-- 初始化脚本示例
CREATE TABLE IF NOT EXISTS containers (...);
CREATE TABLE IF NOT EXISTS sandboxes (...);
CREATE TABLE IF NOT EXISTS tokens (...);
CREATE TABLE IF NOT EXISTS admins (...);
-- ... 其他表

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_containers_status ON containers(status);
-- ... 其他索引
```

### 10.2 数据迁移

| 步骤 | 操作 |
|-----|------|
| 1 | 备份现有JSON文件 |
| 2 | 创建数据库表结构 |
| 3 | 导入JSON数据到数据库 |
| 4 | 验证数据一致性 |
| 5 | 切换应用使用数据库 |

### 10.3 备份策略

| 备份类型 | 频率 | 保留周期 |
|---------|------|---------|
| 每日增量备份 | 每天凌晨2点 | 7天 |
| 每周全量备份 | 每周日凌晨2点 | 30天 |
| 每月归档备份 | 每月1号凌晨2点 | 1年 |

---

## 十一、总结

### 11.1 设计亮点

1. **安全性**：多层加密、细粒度权限控制、完整审计日志、令牌-端点类型匹配
2. **低依赖性**：SQLite嵌入式数据库，无需额外服务
3. **高性能**：合理索引设计、数据分片、连接池管理
4. **可扩展**：模块化设计、预留扩展字段、插件化支持
5. **智能性**：预留AI分析字段、时序数据支持、向量存储

### 11.2 后续工作

1. 实现DAO层代码
2. 编写数据迁移脚本
3. 实现监控数据采集服务
4. 实现监控告警服务
5. 实现令牌管理服务
6. 实现管理员系统服务
7. 编写单元测试和集成测试