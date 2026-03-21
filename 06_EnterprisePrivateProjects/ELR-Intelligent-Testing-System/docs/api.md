# ELR 智能测试系统 - API 文档

## 1. 概述

本文档描述了 ELR 智能测试系统的 RESTful API 接口，用于系统的远程调用和集成。API 基于 FastAPI 框架实现，支持自动 API 文档生成，以及请求验证和错误处理。

## 2. 基础信息

- **API 基础 URL**：`http://localhost:8000/api`
- **认证方式**：JWT 令牌认证
- **请求格式**：JSON
- **响应格式**：JSON
- **错误处理**：统一的错误响应格式

## 3. 认证 API

### 3.1 登录

- **路径**：`/auth/login`
- **方法**：`POST`
- **请求体**：
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **响应**：
  ```json
  {
    "access_token": "string",
    "token_type": "bearer",
    "expires_in": 3600,
    "user": {
      "id": "string",
      "username": "string",
      "role": "string"
    }
  }
  ```

### 3.2 刷新令牌

- **路径**：`/auth/refresh`
- **方法**：`POST`
- **请求头**：`Authorization: Bearer <refresh_token>`
- **响应**：
  ```json
  {
    "access_token": "string",
    "token_type": "bearer",
    "expires_in": 3600
  }
  ```

## 4. 测试用例 API

### 4.1 获取测试用例列表

- **路径**：`/testcases`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer <access_token>`
- **查询参数**：
  - `page`: 页码，默认 1
  - `page_size`: 每页大小，默认 10
  - `category`: 分类筛选
  - `status`: 状态筛选
- **响应**：
  ```json
  {
    "total": 100,
    "page": 1,
    "page_size": 10,
    "items": [
      {
        "id": "string",
        "name": "string",
        "description": "string",
        "category": "string",
        "status": "string",
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
      }
    ]
  }
  ```

### 4.2 获取测试用例详情

- **路径**：`/testcases/{id}`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer <access_token>`
- **响应**：
  ```json
  {
    "id": "string",
    "name": "string",
    "description": "string",
    "category": "string",
    "status": "string",
    "steps": ["string"],
    "expected_result": "string",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
  ```

### 4.3 创建测试用例

- **路径**：`/testcases`
- **方法**：`POST`
- **请求头**：`Authorization: Bearer <access_token>`
- **请求体**：
  ```json
  {
    "name": "string",
    "description": "string",
    "category": "string",
    "steps": ["string"],
    "expected_result": "string"
  }
  ```
- **响应**：
  ```json
  {
    "id": "string",
    "name": "string",
    "description": "string",
    "category": "string",
    "status": "created",
    "steps": ["string"],
    "expected_result": "string",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
  ```

### 4.4 更新测试用例

- **路径**：`/testcases/{id}`
- **方法**：`PUT`
- **请求头**：`Authorization: Bearer <access_token>`
- **请求体**：
  ```json
  {
    "name": "string",
    "description": "string",
    "category": "string",
    "steps": ["string"],
    "expected_result": "string",
    "status": "string"
  }
  ```
- **响应**：
  ```json
  {
    "id": "string",
    "name": "string",
    "description": "string",
    "category": "string",
    "status": "string",
    "steps": ["string"],
    "expected_result": "string",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
  ```

### 4.5 删除测试用例

- **路径**：`/testcases/{id}`
- **方法**：`DELETE`
- **请求头**：`Authorization: Bearer <access_token>`
- **响应**：
  ```json
  {
    "message": "Test case deleted successfully"
  }
  ```

## 5. 测试计划 API

### 5.1 获取测试计划列表

- **路径**：`/testplans`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer <access_token>`
- **查询参数**：
  - `page`: 页码，默认 1
  - `page_size`: 每页大小，默认 10
  - `status`: 状态筛选
- **响应**：
  ```json
  {
    "total": 50,
    "page": 1,
    "page_size": 10,
    "items": [
      {
        "id": "string",
        "name": "string",
        "description": "string",
        "status": "string",
        "created_at": "2023-01-01T00:00:00Z",
        "updated_at": "2023-01-01T00:00:00Z"
      }
    ]
  }
  ```

### 5.2 获取测试计划详情

- **路径**：`/testplans/{id}`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer <access_token>`
- **响应**：
  ```json
  {
    "id": "string",
    "name": "string",
    "description": "string",
    "status": "string",
    "testcases": [
      {
        "id": "string",
        "name": "string",
        "status": "string"
      }
    ],
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
  ```

### 5.3 创建测试计划

- **路径**：`/testplans`
- **方法**：`POST`
- **请求头**：`Authorization: Bearer <access_token>`
- **请求体**：
  ```json
  {
    "name": "string",
    "description": "string",
    "testcase_ids": ["string"]
  }
  ```
- **响应**：
  ```json
  {
    "id": "string",
    "name": "string",
    "description": "string",
    "status": "created",
    "testcases": [
      {
        "id": "string",
        "name": "string",
        "status": "pending"
      }
    ],
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
  ```

### 5.4 执行测试计划

- **路径**：`/testplans/{id}/execute`
- **方法**：`POST`
- **请求头**：`Authorization: Bearer <access_token>`
- **响应**：
  ```json
  {
    "testplan_id": "string",
    "execution_id": "string",
    "status": "running",
    "started_at": "2023-01-01T00:00:00Z"
  }
  ```

## 6. 测试执行 API

### 6.1 获取执行列表

- **路径**：`/executions`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer <access_token>`
- **查询参数**：
  - `page`: 页码，默认 1
  - `page_size`: 每页大小，默认 10
  - `status`: 状态筛选
  - `testplan_id`: 测试计划 ID 筛选
- **响应**：
  ```json
  {
    "total": 200,
    "page": 1,
    "page_size": 10,
    "items": [
      {
        "id": "string",
        "testplan_id": "string",
        "testplan_name": "string",
        "status": "string",
        "started_at": "2023-01-01T00:00:00Z",
        "completed_at": "2023-01-01T00:00:00Z",
        "duration": 60
      }
    ]
  }
  ```

### 6.2 获取执行详情

- **路径**：`/executions/{id}`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer <access_token>`
- **响应**：
  ```json
  {
    "id": "string",
    "testplan_id": "string",
    "testplan_name": "string",
    "status": "string",
    "started_at": "2023-01-01T00:00:00Z",
    "completed_at": "2023-01-01T00:00:00Z",
    "duration": 60,
    "results": [
      {
        "testcase_id": "string",
        "testcase_name": "string",
        "status": "string",
        "message": "string",
        "duration": 10
      }
    ]
  }
  ```

### 6.3 取消执行

- **路径**：`/executions/{id}/cancel`
- **方法**：`POST`
- **请求头**：`Authorization: Bearer <access_token>`
- **响应**：
  ```json
  {
    "execution_id": "string",
    "status": "cancelled"
  }
  ```

## 7. ELR 容器 API

### 7.1 获取容器列表

- **路径**：`/containers`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer <access_token>`
- **查询参数**：
  - `page`: 页码，默认 1
  - `page_size`: 每页大小，默认 10
  - `status`: 状态筛选
- **响应**：
  ```json
  {
    "total": 50,
    "page": 1,
    "page_size": 10,
    "items": [
      {
        "id": "string",
        "name": "string",
        "status": "string",
        "image": "string",
        "created_at": "2023-01-01T00:00:00Z",
        "started_at": "2023-01-01T00:00:00Z"
      }
    ]
  }
  ```

### 7.2 获取容器详情

- **路径**：`/containers/{id}`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer <access_token>`
- **响应**：
  ```json
  {
    "id": "string",
    "name": "string",
    "status": "string",
    "image": "string",
    "ports": [
      {
        "container_port": 8080,
        "host_port": 8080,
        "protocol": "tcp"
      }
    ],
    "environment": {
      "key": "value"
    },
    "created_at": "2023-01-01T00:00:00Z",
    "started_at": "2023-01-01T00:00:00Z"
  }
  ```

### 7.3 创建容器

- **路径**：`/containers`
- **方法**：`POST`
- **请求头**：`Authorization: Bearer <access_token>`
- **请求体**：
  ```json
  {
    "name": "string",
    "image": "string",
    "ports": [
      {
        "container_port": 8080,
        "host_port": 8080,
        "protocol": "tcp"
      }
    ],
    "environment": {
      "key": "value"
    },
    "resources": {
      "cpu": "1",
      "memory": "1g"
    }
  }
  ```
- **响应**：
  ```json
  {
    "id": "string",
    "name": "string",
    "status": "created",
    "image": "string",
    "created_at": "2023-01-01T00:00:00Z"
  }
  ```

### 7.4 启动容器

- **路径**：`/containers/{id}/start`
- **方法**：`POST`
- **请求头**：`Authorization: Bearer <access_token>`
- **响应**：
  ```json
  {
    "container_id": "string",
    "status": "running",
    "started_at": "2023-01-01T00:00:00Z"
  }
  ```

### 7.5 停止容器

- **路径**：`/containers/{id}/stop`
- **方法**：`POST`
- **请求头**：`Authorization: Bearer <access_token>`
- **响应**：
  ```json
  {
    "container_id": "string",
    "status": "stopped",
    "stopped_at": "2023-01-01T00:00:00Z"
  }
  ```

### 7.6 删除容器

- **路径**：`/containers/{id}`
- **方法**：`DELETE`
- **请求头**：`Authorization: Bearer <access_token>`
- **响应**：
  ```json
  {
    "message": "Container deleted successfully"
  }
  ```

## 8. 监控 API

### 8.1 获取监控指标

- **路径**：`/monitoring/metrics`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer <access_token>`
- **查询参数**：
  - `metric`: 指标名称
  - `start_time`: 开始时间
  - `end_time`: 结束时间
  - `interval`: 时间间隔
- **响应**：
  ```json
  {
    "metric": "string",
    "data": [
      {
        "timestamp": "2023-01-01T00:00:00Z",
        "value": 100
      }
    ]
  }
  ```

### 8.2 获取系统状态

- **路径**：`/monitoring/status`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer <access_token>`
- **响应**：
  ```json
  {
    "system": {
      "status": "healthy",
      "uptime": 3600,
      "cpu_usage": 20,
      "memory_usage": 50
    },
    "containers": [
      {
        "id": "string",
        "name": "string",
        "status": "running",
        "cpu_usage": 10,
        "memory_usage": 30
      }
    ]
  }
  ```

## 9. 报告 API

### 9.1 生成测试报告

- **路径**：`/reports/generate`
- **方法**：`POST`
- **请求头**：`Authorization: Bearer <access_token>`
- **请求体**：
  ```json
  {
    "execution_id": "string",
    "format": "pdf",
    "include_screenshots": true
  }
  ```
- **响应**：
  ```json
  {
    "report_id": "string",
    "execution_id": "string",
    "status": "generated",
    "url": "string"
  }
  ```

### 9.2 获取报告列表

- **路径**：`/reports`
- **方法**：`GET`
- **请求头**：`Authorization: Bearer <access_token>`
- **查询参数**：
  - `page`: 页码，默认 1
  - `page_size`: 每页大小，默认 10
- **响应**：
  ```json
  {
    "total": 50,
    "page": 1,
    "page_size": 10,
    "items": [
      {
        "id": "string",
        "execution_id": "string",
        "format": "pdf",
        "status": "generated",
        "created_at": "2023-01-01T00:00:00Z",
        "url": "string"
      }
    ]
  }
  ```

## 10. 错误响应

所有 API 错误都返回统一的错误响应格式：

```json
{
  "error": {
    "code": "string",
    "message": "string",
    "details": "string"
  }
}
```

常见错误代码：

- `400`: 请求参数错误
- `401`: 未授权
- `403`: 禁止访问
- `404`: 资源不存在
- `500`: 内部服务器错误

## 11. 示例请求

### 11.1 登录示例

```bash
curl -X POST "http://localhost:8000/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'
```

### 11.2 创建测试用例示例

```bash
curl -X POST "http://localhost:8000/api/testcases" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "name": "测试 ELR 启动",
    "description": "测试 ELR 服务是否能正常启动",
    "category": "功能测试",
    "steps": [
      "启动 ELR 容器",
      "等待 5 秒",
      "检查 ELR 服务状态"
    ],
    "expected_result": "ELR 服务状态为运行中"
  }'
```

### 11.3 执行测试计划示例

```bash
curl -X POST "http://localhost:8000/api/testplans/{id}/execute" \
  -H "Authorization: Bearer <access_token>"
```

## 12. 总结

ELR 智能测试系统的 API 接口提供了全面的测试管理、执行和监控功能，支持系统的远程调用和集成。通过这些 API，可以实现测试用例管理、测试计划执行、ELR 容器管理、监控和报告生成等功能，为 ELR 项目的开发和维护提供有力支持。