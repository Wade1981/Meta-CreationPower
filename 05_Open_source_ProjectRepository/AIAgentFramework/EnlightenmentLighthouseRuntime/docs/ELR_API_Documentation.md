# ELR (Enlightenment Lighthouse Runtime) API文档

## 1. 概述

ELR提供了一套完整的RESTful API，用于管理容器、模型和网络服务。API设计遵循RESTful原则，使用JSON格式进行数据交换。

### 1.1 基础URL

所有API请求的基础URL格式为：

```
http://localhost:{port}/api
```

其中 `{port}` 是配置的API端口，默认值为：
- Public API: 8080
- Desktop API: 8081
- Model API: 8082

### 1.2 响应格式

所有API响应都使用JSON格式，包含以下字段：

- `success`：布尔值，表示请求是否成功
- `data`：响应数据（成功时）
- `error`：错误信息（失败时）
- `message`：提示信息
- `timestamp`：时间戳

## 2. 容器管理API

### 2.1 列出容器

**请求**：
- 方法：GET
- 路径：`/api/container/list`

**响应**：
```json
[
  {
    "id": "container-1",
    "name": "Test Container",
    "image": "ubuntu:latest",
    "status": "running",
    "created": "2026-04-11T10:00:00Z",
    "started": "2026-04-11T10:01:00Z",
    "ip_address": "172.16.0.2"
  }
]
```

### 2.2 获取容器状态

**请求**：
- 方法：GET
- 路径：`/api/container/status`
- 查询参数：`id`（容器ID）

**响应**：
```json
{
  "id": "container-1",
  "name": "Test Container",
  "image": "ubuntu:latest",
  "status": "running",
  "created": "2026-04-11T10:00:00Z",
  "started": "2026-04-11T10:01:00Z",
  "ip_address": "172.16.0.2",
  "pid": 12345
}
```

## 3. 模型管理API

### 3.1 运行模型

**请求**：
- 方法：POST
- 路径：`/api/model/run`
- 请求体：
```json
{
  "container_id": "container-1",
  "model_id": "elr-chat",
  "input": "Hello, world!"
}
```

**响应**：
```json
{
  "output": "Hello! How can I help you today?",
  "timestamp": 1712832000
}
```

### 3.2 列出模型

**请求**：
- 方法：GET
- 路径：`/api/model/list`

**响应**：
```json
[
  {
    "id": "elr-chat",
    "name": "ELR Chat Model",
    "version": "1.0",
    "type": "text"
  },
  {
    "id": "fish-speech",
    "name": "Fish Speech Model",
    "version": "1.0",
    "type": "speech"
  }
]
```

## 4. 网络管理API

### 4.1 获取网络状态

**请求**：
- 方法：GET
- 路径：`/api/network/status`

**响应**：
```json
{
  "runtime_network_enabled": true,
  "api_ports": {
    "desktop_api": 8081,
    "public_api": 8080,
    "model_api": 8082
  },
  "desktop_api": {
    "address": "http://localhost:8081",
    "port": 8081,
    "status": "running"
  },
  "public_api": {
    "address": "http://localhost:8080",
    "port": 8080,
    "status": "running"
  },
  "model_api": {
    "address": "http://localhost:8082/api/model",
    "port": 8082,
    "status": "running"
  },
  "containers": [
    {
      "id": "container-1",
      "name": "Test Container",
      "network_enabled": true,
      "network_mode": "nat",
      "ip_address": "172.16.0.2",
      "port_mappings": [],
      "status": "running"
    }
  ],
  "container_count": 1,
  "timestamp": 1712832000
}
```

### 4.2 隔离容器网络

**请求**：
- 方法：POST
- 路径：`/api/network/isolate`
- 请求体：
```json
{
  "container_id": "container-1"
}
```

**响应**：
```json
{
  "message": "Network isolated successfully",
  "container_id": "container-1",
  "ip_address": "172.18.0.2",
  "timestamp": 1712832000
}
```

### 4.3 取消网络隔离

**请求**：
- 方法：POST
- 路径：`/api/network/unisolate`
- 请求体：
```json
{
  "container_id": "container-1"
}
```

**响应**：
```json
{
  "message": "Network isolation removed successfully",
  "container_id": "container-1",
  "timestamp": 1712832000
}
```

### 4.4 获取网络配置

**请求**：
- 方法：GET
- 路径：`/api/network/config`
- 查询参数：`container_id`（容器ID）

**响应**：
```json
{
  "network_id": "net-container-1",
  "container_id": "container-1",
  "ip_address": "172.18.0.2",
  "subnet": "172.18.0.0/16",
  "allowed_ports": [80, 443, 8080],
  "blocked_ports": [22, 3389],
  "enabled": true
}
```

## 5. 令牌管理API

### 5.1 创建令牌

**请求**：
- 方法：POST
- 路径：`/api/token/create`
- 请求体：
```json
{
  "description": "Test Token"
}
```

**响应**：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "Token generated successfully",
  "timestamp": 1712832000
}
```

### 5.2 验证令牌

**请求**：
- 方法：POST
- 路径：`/api/token/validate`
- 请求体：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应**：
```json
{
  "valid": true,
  "message": "Token is valid",
  "timestamp": 1712832000
}
```

### 5.3 刷新令牌

**请求**：
- 方法：POST
- 路径：`/api/token/refresh`
- 请求体：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "description": "Refreshed Token"
}
```

**响应**：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "Token refreshed successfully",
  "timestamp": 1712832000
}
```

### 5.4 列出令牌

**请求**：
- 方法：GET
- 路径：`/api/token/list`

**响应**：
```json
{
  "tokens": [
    {
      "id": "token-1",
      "description": "Test Token",
      "created_at": "2026-04-11T10:00:00Z",
      "expires_at": "2026-05-11T10:00:00Z"
    }
  ],
  "count": 1,
  "timestamp": 1712832000
}
```

### 5.5 撤销令牌

**请求**：
- 方法：POST
- 路径：`/api/token/revoke`
- 请求体：
```json
{
  "token_id": "token-1"
}
```

**响应**：
```json
{
  "message": "Token revoked successfully",
  "timestamp": 1712832000
}
```

## 6. 桌面API

### 6.1 健康检查

**请求**：
- 方法：GET
- 路径：`/api/desktop/health`

**响应**：
```json
{
  "status": "ok",
  "timestamp": 1712832000,
  "service": "elr-desktop-api",
  "version": "1.0.0"
}
```

### 6.2 获取ELR状态

**请求**：
- 方法：GET
- 路径：`/api/desktop/status`

**响应**：
```json
{
  "status": "running",
  "message": "ELR Desktop API服务运行正常",
  "timestamp": 1712832000,
  "containers": 1,
  "api_version": "1.0.0"
}
```

### 6.3 获取容器列表

**请求**：
- 方法：GET
- 路径：`/api/desktop/containers`

**响应**：
```json
[
  {
    "id": "container-1",
    "name": "Test Container",
    "image": "ubuntu:latest",
    "status": "running",
    "created": "2026-04-11 10:00:00",
    "started": "2026-04-11 10:01:00",
    "ip_address": "172.16.0.2"
  }
]
```

### 6.4 获取系统资源

**请求**：
- 方法：GET
- 路径：`/api/desktop/resources`

**响应**：
```json
{
  "success": true,
  "resources": {
    "memory": {
      "total": 17179869184,
      "used": 4294967296,
      "free": 12884901888,
      "usage_percent": 25.0
    },
    "cpu": {
      "usage_percent": 15.5,
      "cores": 8
    },
    "disk": {
      "total": 536870912000,
      "used": 107374182400,
      "free": 429496729600,
      "usage_percent": 20.0
    },
    "system": {
      "platform": "windows",
      "version": "10.0.19045",
      "timestamp": 1712832000
    }
  },
  "timestamp": 1712832000
}
```

### 6.5 列出上传文件

**请求**：
- 方法：GET
- 路径：`/api/desktop/files`

**响应**：
```json
{
  "success": true,
  "files": [
    {
      "name": "test.py",
      "type": "python",
      "size": 1024,
      "path": "C:\\Users\\User\\.elr\\data\\uploads\\test.py",
      "created": 1712832000
    }
  ],
  "timestamp": 1712832000
}
```

### 6.6 上传文件

**请求**：
- 方法：POST
- 路径：`/api/desktop/upload`
- 内容类型：`multipart/form-data`
- 表单字段：`file`（文件）

**响应**：
```json
{
  "success": true,
  "message": "File uploaded successfully: test.py",
  "file_type": "python",
  "filepath": "C:\\Users\\User\\.elr\\data\\uploads\\test.py",
  "file_size": 1024,
  "timestamp": 1712832000
}
```

## 7. 健康检查API

### 7.1 系统健康检查

**请求**：
- 方法：GET
- 路径：`/health`

**响应**：
```json
{
  "status": "ok",
  "timestamp": 1712832000,
  "service": "elr-network",
  "version": "1.0.0"
}
```

## 8. 错误处理

### 8.1 常见错误代码

| 状态码 | 描述 | 示例 |
|--------|------|------|
| 400 | 错误请求 | 缺少必要参数 |
| 404 | 资源不存在 | 容器未找到 |
| 405 | 方法不允许 | 使用GET方法访问POST端点 |
| 429 | 请求过多 | 速率限制 exceeded |
| 500 | 服务器错误 | 内部服务器错误 |

### 8.2 错误响应格式

```json
{
  "error": "Container not found",
  "timestamp": 1712832000
}
```

## 9. 使用示例

### 9.1 使用cURL调用API

```bash
# 列出容器
curl http://localhost:8080/api/container/list

# 运行模型
curl -X POST http://localhost:8080/api/model/run \
  -H "Content-Type: application/json" \
  -d '{"container_id": "container-1", "model_id": "elr-chat", "input": "Hello"}'

# 隔离容器网络
curl -X POST http://localhost:8080/api/network/isolate \
  -H "Content-Type: application/json" \
  -d '{"container_id": "container-1"}'
```

### 9.2 使用Python调用API

```python
import requests
import json

# 列出容器
response = requests.get('http://localhost:8080/api/container/list')
print(response.json())

# 运行模型
payload = {
    'container_id': 'container-1',
    'model_id': 'elr-chat',
    'input': 'Hello'
}
response = requests.post('http://localhost:8080/api/model/run', json=payload)
print(response.json())

# 隔离容器网络
payload = {'container_id': 'container-1'}
response = requests.post('http://localhost:8080/api/network/isolate', json=payload)
print(response.json())
```

## 10. 安全注意事项

- **API访问控制**：生产环境应限制API访问
- **令牌管理**：定期轮换令牌
- **网络隔离**：为容器启用网络隔离
- **输入验证**：所有用户输入都应验证
- **HTTPS**：生产环境应使用HTTPS

## 11. 版本控制

API版本通过URL路径进行控制，当前版本为v1。未来版本将通过路径如`/api/v2/`进行区分。

## 12. 结论

ELR API提供了一套完整的接口，用于管理容器、模型和网络服务。API设计简洁明了，易于使用，支持多种客户端语言。通过这些API，用户可以：

- 管理容器生命周期
- 运行和管理AI模型
- 配置网络隔离和安全策略
- 监控系统资源和状态
- 上传和管理文件

ELR API为构建自动化工具和集成第三方系统提供了强大的基础。
