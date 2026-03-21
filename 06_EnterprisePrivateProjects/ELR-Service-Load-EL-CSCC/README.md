# ELR 服务负载 - EL-CSCC 档案库

本项目实现了一个 ELR（启蒙灯塔运行时）服务负载，用于 EL-CSCC（启蒙灯塔・碳硅协同共识晶体）档案库数字档案管理系统。

## 项目概述

ELR 服务负载 - EL-CSCC 档案库提供了一个 RESTful API 用于管理数字档案，包括：
- 档案的创建、检索、更新和删除
- 档案列表和搜索
- 服务监控的健康检查端点
- 使用 JSON 格式永久存储档案数据

## 架构

项目采用模块化架构：

```
ELR-Service-Load-EL-CSCC/
├── main.py              # 主入口点
├── src/
│   ├── config/          # 配置管理
│   ├── service/         # 档案服务实现
│   ├── api/             # HTTP 服务器和 API 端点
│   └── utils/           # 工具函数
└── README.md            # 本文档
```

## API 端点

| 方法 | 端点 | 描述 | 认证要求 |
|------|------|------|----------|
| GET    | /health  | 健康检查端点 | 否 |
| POST   | /login | 用户登录获取令牌 | 否 |
| GET    | /archives | 列出所有档案 | 否 |
| GET    | /archives/{id} | 获取指定档案 | 否 |
| POST   | /archives | 创建新档案 | 是 |
| PUT    | /archives/{id} | 更新现有档案 | 是 |
| DELETE | /archives/{id} | 删除指定档案 | 是 |
| GET    | /search?q={query} | 按查询词搜索档案 | 否 |
| WS     | /ws | WebSocket 连接 | 否 |

## 环境变量

服务可以通过以下环境变量进行配置：

| 变量 | 默认值 | 描述 |
|------|--------|------|
| ELR_SERVICE_HOST | 0.0.0.0 | HTTP 服务器绑定的主机 |
| ELR_SERVICE_PORT | 8000 | HTTP 服务器绑定的端口 |
| ELR_ARCHIVE_FILE | el_cscc_archive.json | 档案 JSON 文件路径 |
| ELR_LOG_LEVEL | INFO | 日志级别（DEBUG, INFO, WARNING, ERROR） |
| ELR_WEBSOCKET_HOST | 与 ELR_SERVICE_HOST 相同 | WebSocket 服务器绑定的主机 |
| ELR_WEBSOCKET_PORT | 8001 | WebSocket 服务器绑定的端口 |
| ELR_SECRET_KEY | default_secret_key | 用于生成和验证令牌的密钥 |
| ELR_CERT_FILE | 无 | SSL 证书文件路径 |
| ELR_KEY_FILE | 无 | SSL 私钥文件路径 |

## 使用方法

### 运行服务

1. 导航到项目目录：
   ```
   cd E:\X54\github\Meta-CreationPower\06_EnterprisePrivateProjects\ELR-Service-Load-EL-CSCC
   ```

2. 运行服务：
   ```
   python main.py
   ```

3. 服务将默认在 `http://0.0.0.0:8000` 上启动。

### 使用 API

#### 健康检查
```bash
curl http://localhost:8000/health
```

#### 用户登录
```bash
curl -X POST http://localhost:8000/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'
```

#### 列出档案
```bash
curl http://localhost:8000/archives
```

#### 创建档案（需要认证）
```bash
# 首先获取令牌
TOKEN=$(curl -s -X POST http://localhost:8000/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}' \
  | jq -r '.token')

# 使用令牌创建档案
curl -X POST http://localhost:8000/archives \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "测试档案",
    "description": "一个测试档案",
    "carbon_partner": "X54SIR",
    "silicon_partner": "CODE-WEAVER",
    "archive_type": "test"
  }'
```

#### 获取档案
```bash
curl http://localhost:8000/archives/{archive_id}
```

#### 更新档案（需要认证）
```bash
curl -X PUT http://localhost:8000/archives/{archive_id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"description": "更新的测试档案"}'
```

#### 删除档案（需要认证）
```bash
curl -X DELETE http://localhost:8000/archives/{archive_id} \
  -H "Authorization: Bearer $TOKEN"
```

#### 搜索档案
```bash
curl http://localhost:8000/search?q=测试
```

#### 使用 WebSocket

```javascript
// 连接 WebSocket 服务器
const socket = new WebSocket('ws://localhost:8001');

// 连接建立
 socket.onopen = function(event) {
    console.log('WebSocket 连接已建立');
    
    // 发送 ping 消息
    socket.send(JSON.stringify({ type: 'ping' }));
    
    // 获取档案列表
    socket.send(JSON.stringify({ type: 'list_archives' }));
};

// 接收消息
 socket.onmessage = function(event) {
    const message = JSON.parse(event.data);
    console.log('收到消息:', message);
    
    // 处理不同类型的消息
    if (message.type === 'pong') {
        console.log('收到 pong 响应');
    } else if (message.type === 'archives') {
        console.log('档案列表:', message.data);
    } else if (message.type === 'error') {
        console.error('错误:', message.message);
    }
};

// 连接关闭
 socket.onclose = function(event) {
    console.log('WebSocket 连接已关闭');
};

// 发生错误
 socket.onerror = function(error) {
    console.error('WebSocket 错误:', error);
};
```

## 依赖

项目使用以下 Python 库：
- **标准库**：
  - `http.server` 用于 HTTP 服务器
  - `json` 用于 JSON 解析和序列化
  - `urllib.parse` 用于 URL 解析
  - `logging` 用于日志记录
  - `os` 用于环境变量和文件操作
  - `datetime` 用于时间戳生成
  - `asyncio` 用于异步 I/O 操作
  - `ssl` 用于 SSL/TLS 加密（可选）

- **外部依赖**：
  - `websockets` 用于 WebSocket 服务器实现

### 安装依赖

```bash
pip install websockets
```

## ELR 兼容性

此服务负载设计为与 ELR 容器兼容，遵循以下原则：
- 轻量级且最小依赖
- 带有健康检查端点的长期运行服务
- 通过环境变量可配置
- 可在不同环境中移植
- 支持容器化和编排

## 数据存储

档案数据存储在 JSON 文件中（默认名为 `el_cscc_archive.json`），在服务重启之间提供永久存储。如果文件不存在，会自动创建。

## 日志记录

服务会将启动、关闭和关键操作的信息记录到控制台。日志级别可以通过 `ELR_LOG_LEVEL` 环境变量配置。

## 许可证

本项目是启蒙灯塔起源团队 Meta-CreationPower 框架的一部分。
