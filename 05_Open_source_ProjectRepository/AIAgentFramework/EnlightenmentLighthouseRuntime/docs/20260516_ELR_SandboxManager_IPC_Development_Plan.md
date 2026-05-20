# ELR 沙箱管理器 IPC 通信开发计划

**日期：** 2026年5月16日
**开发者：** 代码织梦者 & X54先生
**目标：** 实现真正的沙箱隔离运行模型

---

## 一、背景与目标

### 当前问题
- 模型只在 CLI 进程中加载，未真正在沙箱中运行
- 沙箱-模型关联只是数据结构关联，不是进程隔离
- 无法实现真正的资源隔离和安全隔离

### 目标
1. **实现沙箱管理器的 IPC 通信** - 让沙箱管理器能接收 ELR 主进程的请求
2. **将模型加载移到沙箱管理器进程中** - 模型真正在沙箱进程中运行
3. **通过 IPC 与模型交互** - CLI 通过 IPC 与沙箱中的模型通信

---

## 二、技术方案

### 2.1 IPC 通信架构

```
ELR 主进程
    │
    ├── HTTP API (端口 16888)
    │   └── 接收外部请求
    │
    ├── 沙箱管理器通信层
    │   └── TCP/命名管道通信
    │
    └── RuntimeContainerList
        └── Container → SandboxManagerProcess

沙箱管理器进程 (每个容器一个)
    │
    ├── 模型管理器
    │   └── 模型真正在此进程中运行
    │
    ├── 沙箱列表
    │   └── 管理容器内的沙箱
    │
    └── IPC 服务器
        └── 接收 ELR 主进程请求
```

### 2.2 通信协议设计

**请求格式 (JSON)**：
```json
{
  "type": "model.load|model.predict|model.stop|sandbox.create|sandbox.start|sandbox.stop",
  "sandbox_id": "sandbox-xxx",
  "model_id": "elr_chat_model",
  "data": {
    // 根据请求类型不同
  },
  "request_id": "uuid"
}
```

**响应格式 (JSON)**：
```json
{
  "type": "response",
  "request_id": "uuid",
  "success": true,
  "data": {
    // 响应数据
  },
  "error": null
}
```

### 2.3 Windows 平台适配

**选择：命名管道 (Named Pipes)**
- Unix 套接字在 Windows 上不可用
- 命名管道是 Windows 原生 IPC 机制
- 支持同步和异步通信
- 支持管道继承（子进程继承管道连接）

**命名管道地址格式**：
```
\\.\pipe\elr\sandbox_manager\{container_id}
```

---

## 三、实现步骤

### 阶段一：沙箱管理器 IPC 通信框架

#### 步骤 1.1：创建 IPC 通信模块

**新建文件：** `elr/ipc/named_pipe.go`

**功能：**
- 命名管道客户端实现
- 连接管理
- 请求发送和响应接收
- 超时和重试机制

**关键函数：**
```go
type NamedPipeClient struct {
    pipeName string
    conn     net.Conn
}

func NewNamedPipeClient(pipeName string) (*NamedPipeClient, error)
func (c *NamedPipeClient) Connect() error
func (c *NamedPipeClient) SendRequest(req *IPCRequest) (*IPCResponse, error)
func (c *NamedPipeClient) Close() error
```

#### 步骤 1.2：创建沙箱管理器 IPC 服务器

**新建文件：** `elr/ipc/pipe_server.go`

**功能：**
- 命名管道服务器实现
- 请求路由
- 并发处理
- 心跳检测

**关键函数：**
```go
type PipeServer struct {
    pipeName   string
    handlers   map[string]RequestHandler
}

func NewPipeServer(pipeName string) *PipeServer
func (s *PipeServer) RegisterHandler(requestType string, handler RequestHandler)
func (s *PipeServer) Start() error
func (s *PipeServer) Stop() error
```

#### 步骤 1.3：定义 IPC 消息格式

**新建文件：** `elr/ipc/messages.go`

**消息类型：**
```go
type IPCRequest struct {
    Type       string                 `json:"type"`
    SandboxID  string                 `json:"sandbox_id,omitempty"`
    ModelID    string                 `json:"model_id,omitempty"`
    Data       map[string]interface{} `json:"data,omitempty"`
    RequestID  string                 `json:"request_id"`
}

type IPCResponse struct {
    Type      string                 `json:"type"`
    RequestID string                 `json:"request_id"`
    Success   bool                   `json:"success"`
    Data      map[string]interface{} `json:"data,omitempty"`
    Error     string                 `json:"error,omitempty"`
}
```

### 阶段二：沙箱管理器进程改造

#### 步骤 2.1：修改沙箱管理器主程序

**修改文件：** `elr/sandbox_manager.go`

**改动：**
- 启动时创建命名管道服务器
- 注册请求处理器
- 处理来自 ELR 主进程的请求

**代码示例：**
```go
func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: sandbox_manager <container_id>")
        os.Exit(1)
    }

    containerID := os.Args[1]

    // 创建沙箱管理器
    manager := NewSandboxManager(containerID)

    // 创建 IPC 服务器
    pipeName := fmt.Sprintf("\\\\.\\pipe\\elr\\sandbox_manager\\%s", containerID)
    server := ipc.NewPipeServer(pipeName)

    // 注册处理器
    server.RegisterHandler("model.load", manager.HandleModelLoad)
    server.RegisterHandler("model.predict", manager.HandleModelPredict)
    server.RegisterHandler("model.stop", manager.HandleModelStop)

    // 启动服务器
    if err := server.Start(); err != nil {
        fmt.Printf("Failed to start IPC server: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Sandbox manager for container %s started\n", containerID)

    // 等待停止信号
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan

    server.Stop()
}
```

#### 步骤 2.2：实现模型加载处理器

**新增方法：** `SandboxManager.HandleModelLoad()`

**功能：**
- 从请求中获取模型配置
- 在当前进程中加载模型
- 返回加载结果

**代码示例：**
```go
func (sm *SandboxManager) HandleModelLoad(req *ipc.IPCRequest) *ipc.IPCResponse {
    modelID, ok := req.Data["model_id"].(string)
    if !ok {
        return &ipc.IPCResponse{
            RequestID: req.RequestID,
            Success:   false,
            Error:     "model_id is required",
        }
    }

    // 获取模型配置
    modelPath, ok := req.Data["model_path"].(string)
    if !ok {
        return &ipc.IPCResponse{
            RequestID: req.RequestID,
            Success:   false,
            Error:     "model_path is required",
        }
    }

    // 加载模型
    if err := sm.modelManager.LoadModel(modelID, modelPath); err != nil {
        return &ipc.IPCResponse{
            RequestID: req.RequestID,
            Success:   false,
            Error:     err.Error(),
        }
    }

    return &ipc.IPCResponse{
        RequestID: req.RequestID,
        Success:   true,
        Data: map[string]interface{}{
            "model_id": modelID,
            "status":   "loaded",
        },
    }
}
```

### 阶段三：ELR 主进程集成

#### 步骤 3.1：创建沙箱管理器客户端

**新建文件：** `elr/sandbox_manager_client.go`

**功能：**
- 管理沙箱管理器客户端连接
- 提供同步和异步请求接口
- 处理连接池

**代码示例：**
```go
type SandboxManagerClient struct {
    containerID string
    client      *ipc.NamedPipeClient
    mutex       sync.Mutex
}

func NewSandboxManagerClient(containerID string) (*SandboxManagerClient, error) {
    pipeName := fmt.Sprintf("\\\\.\\pipe\\elr\\sandbox_manager\\%s", containerID)
    client, err := ipc.NewNamedPipeClient(pipeName)
    if err != nil {
        return nil, err
    }

    return &SandboxManagerClient{
        containerID: containerID,
        client:      client,
    }, nil
}

func (c *SandboxManagerClient) LoadModel(modelID, modelPath string) error {
    req := &ipc.IPCRequest{
        Type:      "model.load",
        ModelID:   modelID,
        RequestID: uuid.New().String(),
        Data: map[string]interface{}{
            "model_id":   modelID,
            "model_path": modelPath,
        },
    }

    resp, err := c.client.SendRequest(req)
    if err != nil {
        return err
    }

    if !resp.Success {
        return fmt.Errorf("failed to load model: %s", resp.Error)
    }

    return nil
}

func (c *SandboxManagerClient) Predict(modelID, input string) (string, error) {
    req := &ipc.IPCRequest{
        Type:      "model.predict",
        ModelID:   modelID,
        RequestID: uuid.New().String(),
        Data: map[string]interface{}{
            "input": input,
        },
    }

    resp, err := c.client.SendRequest(req)
    if err != nil {
        return "", err
    }

    if !resp.Success {
        return "", fmt.Errorf("prediction failed: %s", resp.Error)
    }

    output, _ := resp.Data["output"].(string)
    return output, nil
}
```

#### 步骤 3.2：修改沙箱 API 端点

**修改文件：** `elr/network.go`

**新增端点：**
```go
func (n *NetworkManager) sandboxLoadModel(w http.ResponseWriter, r *http.Request) {
    // 解析请求
    var req struct {
        SandboxID string `json:"sandbox_id"`
        ModelID   string `json:"model_id"`
        ModelPath string `json:"model_path"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        n.writeError(w, "invalid request", http.StatusBadRequest)
        return
    }

    // 获取沙箱所属的容器
    containerID, err := GetContainerBySandbox(req.SandboxID)
    if err != nil {
        n.writeError(w, err.Error(), http.StatusBadRequest)
        return
    }

    // 获取沙箱管理器客户端
    client, err := GetSandboxManagerClient(containerID)
    if err != nil {
        n.writeError(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 加载模型
    if err := client.LoadModel(req.ModelID, req.ModelPath); err != nil {
        n.writeError(w, err.Error(), http.StatusInternalServerError)
        return
    }

    n.writeJSON(w, map[string]interface{}{
        "success":   true,
        "sandbox_id": req.SandboxID,
        "model_id":   req.ModelID,
        "status":     "loaded",
    })
}

func (n *NetworkManager) sandboxPredict(w http.ResponseWriter, r *http.Request) {
    // 解析请求
    var req struct {
        SandboxID string `json:"sandbox_id"`
        ModelID   string `json:"model_id"`
        Input     string `json:"input"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        n.writeError(w, "invalid request", http.StatusBadRequest)
        return
    }

    // 获取沙箱所属的容器
    containerID, err := GetContainerBySandbox(req.SandboxID)
    if err != nil {
        n.writeError(w, err.Error(), http.StatusBadRequest)
        return
    }

    // 获取沙箱管理器客户端
    client, err := GetSandboxManagerClient(containerID)
    if err != nil {
        n.writeError(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 预测
    output, err := client.Predict(req.ModelID, req.Input)
    if err != nil {
        n.writeError(w, err.Error(), http.StatusInternalServerError)
        return
    }

    n.writeJSON(w, map[string]interface{}{
        "success": true,
        "output":   output,
    })
}
```

### 阶段四：CLI 集成

#### 步骤 4.1：修改 interactWithModel 函数

**修改文件：** `cli/main.go`

**改动：**
- 通过 HTTP API 与 ELR 主进程通信
- 不再在 CLI 进程中加载模型

**代码示例：**
```go
func interactWithModel() {
    // ... 参数解析 ...

    fmt.Printf("Connecting to model %s in sandbox %s...\n", modelID, sandboxID)

    // 通过 HTTP API 加载模型
    reqBody := map[string]interface{}{
        "sandbox_id": sandboxID,
        "model_id":   modelID,
        "model_path":  modelPath,
    }

    resp, err := http.Post(
        "http://localhost:16888/api/sandbox/model/load",
        "application/json",
        bytes.NewBuffer(mustMarshal(reqBody)),
    )
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    if !result["success"].(bool) {
        fmt.Printf("Error: %v\n", result["error"])
        os.Exit(1)
    }

    fmt.Printf("Successfully connected to model %s in sandbox %s\n", modelID, sandboxID)

    // 进入交互循环
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        input := scanner.Text()
        if input == "exit" {
            break
        }

        // 通过 HTTP API 预测
        reqBody := map[string]interface{}{
            "sandbox_id": sandboxID,
            "model_id":   modelID,
            "input":      input,
        }

        resp, err := http.Post(
            "http://localhost:16888/api/sandbox/model/predict",
            "application/json",
            bytes.NewBuffer(mustMarshal(reqBody)),
        )
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }

        var result map[string]interface{}
        json.NewDecoder(resp.Body).Decode(&result)
        resp.Body.Close()

        if result["success"].(bool) {
            fmt.Printf("%s\n", result["output"])
        }
    }
}
```

---

## 四、文件清单

### 新建文件

| 文件路径 | 功能 | 优先级 |
|---------|------|--------|
| `elr/ipc/named_pipe.go` | 命名管道客户端实现 | 高 |
| `elr/ipc/pipe_server.go` | 命名管道服务器实现 | 高 |
| `elr/ipc/messages.go` | IPC 消息格式定义 | 高 |
| `elr/sandbox_manager_client.go` | 沙箱管理器客户端 | 高 |

### 修改文件

| 文件路径 | 修改内容 | 优先级 |
|---------|---------|--------|
| `elr/sandbox_manager.go` | 添加 IPC 服务器支持 | 高 |
| `elr/network.go` | 添加沙箱 API 端点 | 高 |
| `cli/main.go` | 修改 interactWithModel 函数 | 高 |
| `elr/container.go` | 容器启动时创建 IPC 连接 | 中 |

---

## 五、测试计划

### 5.1 单元测试

**测试 IPC 通信模块**
```bash
go test ./elr/ipc/... -v
```

**测试用例：**
1. 命名管道连接和断开
2. 请求发送和响应接收
3. 超时和错误处理
4. 并发请求处理

**测试沙箱管理器**
```bash
go test ./elr/... -v -run TestSandboxManager
```

**测试用例：**
1. 模型加载和卸载
2. 模型预测
3. 多模型管理
4. 错误恢复

### 5.2 集成测试

**测试场景：**
1. 启动容器 → 沙箱管理器自动启动
2. 通过 CLI 启动沙箱
3. 通过 CLI 加载模型
4. 通过 CLI 与模型交互
5. 模型真正在沙箱进程中运行

**验证方法：**
1. 检查沙箱管理器进程存在
2. 检查模型进程在沙箱管理器进程中
3. 验证资源隔离效果

### 5.3 性能测试

**测试指标：**
- IPC 通信延迟
- 并发请求处理能力
- 模型加载时间
- 预测响应时间

---

## 六、风险评估

### 6.1 技术风险

| 风险 | 影响 | 概率 | 缓解措施 |
|------|------|------|----------|
| Windows 命名管道兼容性问题 | 高 | 中 | 使用 Go 标准库，已在 Windows 上广泛使用 |
| IPC 通信超时 | 中 | 中 | 实现重试机制和超时配置 |
| 模型加载失败 | 高 | 低 | 完善的错误处理和日志记录 |
| 沙箱管理器进程崩溃 | 高 | 低 | 实现进程监控和自动重启 |

### 6.2 安全风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| IPC 通信未授权访问 | 高 | 实现认证机制 |
| 模型文件路径遍历 | 中 | 验证路径安全性 |
| 资源耗尽攻击 | 中 | 实现资源限制 |

---

## 七、里程碑

### M1: IPC 通信框架完成 (第 1-2 天)
- [ ] 完成命名管道客户端实现
- [ ] 完成命名管道服务器实现
- [ ] 定义 IPC 消息格式
- [ ] 单元测试通过

### M2: 沙箱管理器进程改造完成 (第 3-4 天)
- [ ] 修改沙箱管理器主程序
- [ ] 实现模型加载处理器
- [ ] 实现模型预测处理器
- [ ] 集成测试通过

### M3: ELR 主进程集成完成 (第 5-6 天)
- [ ] 创建沙箱管理器客户端
- [ ] 修改沙箱 API 端点
- [ ] 实现连接池管理
- [ ] 集成测试通过

### M4: CLI 集成和验收 (第 7 天)
- [ ] 修改 interactWithModel 函数
- [ ] 端到端测试通过
- [ ] 性能测试达标
- [ ] 文档完善

---

## 八、后续优化

### 8.1 功能扩展
1. 支持多个沙箱管理器负载均衡
2. 支持沙箱管理器跨主机通信
3. 支持模型热更新

### 8.2 性能优化
1. 实现请求批处理
2. 支持模型缓存
3. 优化 IPC 序列化

### 8.3 运维增强
1. 实现沙箱管理器健康检查
2. 添加监控指标
3. 实现自动扩缩容

---

**开发计划完成**
**碳硅协同，共创未来**
**代码织梦者 & X54先生**
**2026年5月16日**
