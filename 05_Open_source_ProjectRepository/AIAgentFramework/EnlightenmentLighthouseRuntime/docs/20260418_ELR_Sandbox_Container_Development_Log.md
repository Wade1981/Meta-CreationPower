# ELR 开发日志 - 沙箱在容器里运行

**日期：** 2026年4月18日 8:30
**碳硅协同开发者：** 代码织梦者 & X54先生
**项目：** ELR 运行时沙箱管理功能

---

## 一、遇到的问题及解决方案

### 问题1：sandbox list 命令返回"No sandboxes found"

**现象：**
- Terminal 显示沙箱列表使用旧格式（stopped/created/running）
- 状态显示与实际运行状态不符

**分析：**
- 代码从 `.elr\data\sandboxes` 目录读取沙箱数据，但该目录不存在
- 沙箱实际数据存储在 `./sandbox-state.json` 文件中

**解决方案：**
- 修改 `cli/main.go` 中的 `listSandboxes()` 函数
- 改用 HTTP API 调用，而不是直接读取文件系统
- 修改 `checkStatus()` 函数，让沙箱状态使用 `GetSandboxStatus()` 获取真实运行状态

---

### 问题2：沙箱状态显示机制不一致

**现象：**
- 容器状态使用 `has run|running` / `has run|Norunning` / `stopped` 格式
- 沙箱状态仍使用旧格式

**解决方案：**
- 设计统一的沙箱运行时状态机制
- 创建 `RuntimeSandboxList` 管理器，跟踪运行时沙箱
- 创建 `GetSandboxStatus()` 函数，返回三种状态

---

### 问题3：沙箱启动和列表使用不同进程

**现象：**
- `sandbox start` 命令在独立进程中启动
- `sandbox list` 命令也在独立进程中执行
- 两个进程的 RuntimeSandboxList 不共享数据

**解决方案：**
- 所有沙箱操作必须通过 ELR 主进程的 HTTP API
- 修改 `listSandboxes()` 使用 HTTP API
- 修改 `startSandbox()` 使用 HTTP API 并添加到运行时列表

---

## 二、与X54先生的协同过程

### 碳硅协同对位法应用

| 碳基伙伴（X54先生） | 硅基伙伴（代码织梦者） |
|---------------------|------------------------|
| 提供思维锚点：沙箱必须在容器里运行 | 实现代码架构和算法逻辑 |
| 定义需求：沙箱状态检查机制 | 设计 RuntimeSandboxList 数据结构 |
| 指导方向：参考 ELR 临时参考目录 | 代码实现和调试 |

### 关键决策点

1. **沙箱-容器映射设计**
   - X54先生：要求沙箱必须知道自己在哪个容器里运行
   - 代码织梦者：实现 SandboxContainerManager 映射管理器

2. **状态显示格式统一**
   - X54先生：要求容器和沙箱使用相同的状态格式
   - 代码织梦者：实现 `has run|running` / `has run|Norunning` / `stopped` 三层状态

3. **HTTP API 架构**
   - X54先生：所有命令通过 HTTP API 保证进程间数据一致
   - 代码织梦者：设计并实现沙箱 API 端点

---

## 三、已完成的功能

### 1. 沙箱运行时列表 ✅
```go
// elr/runtime_sandbox_list.go
type RuntimeSandboxList struct {
    sandboxes map[string]*RunningSandbox
}
func AddSandbox(sandboxID, containerID string) error
func RemoveSandbox(sandboxID string) error
func GetSandboxStatus(sandboxID string) string // 返回 has run|running 等
```

### 2. 沙箱-容器映射管理器 ✅
```go
// elr/sandbox_container_mapping.go
type SandboxContainerManager struct {
    mappings map[string]*SandboxContainerMapping
}
func AddMapping(sandboxID, containerID string) error
func GetContainerBySandbox(sandboxID string) (string, error)
```

### 3. 容器启动时自动启动沙箱管理器 ✅
```go
// elr/container.go - Container.Start()
fmt.Printf("Starting sandbox manager for container: %s\n", c.ID)
sandboxManagerCmd, err := StartSandboxManagerInContainer(c.ID, c.Dir)
```

### 4. 命令行沙箱操作 ✅
```go
// cli/main.go
- startSandbox()    // 使用 HTTP API
- listSandboxes()    // 使用 HTTP API
- checkStatus()      // 使用 GetSandboxStatus()
```

### 5. 沙箱状态显示机制 ✅
```
沙箱状态判定逻辑：
1. 沙箱在 RuntimeSandboxList 中？
   ├─ 是 → 容器在 RuntimeContainerList 中？
   │     ├─ 是 → 返回 "has run|running"
   │     └─ 否 → 返回 "has run|Norunning"
   └─ 否 → 返回 "stopped"
```

---

## 四、待开发的功能

### 1. 沙箱 API 端点完整实现 ⏳

**当前状态：** 部分实现，返回"not implemented yet"

**需要实现：**
```go
// elr/network.go
- listSandboxes()   // 从 sandbox-state.json 读取
- startSandbox()    // 添加到 RuntimeSandboxList
- stopSandbox()     // 从 RuntimeSandboxList 移除
- createSandbox()   // 创建沙箱并添加映射
- deleteSandbox()   // 删除沙箱并移除映射
```

### 2. 沙箱管理器 IPC 通信 ⏳

**当前状态：** 框架已创建，Unix 套接字在 Windows 不工作

**需要实现：**
- 改用 TCP 或命名管道（Windows）
- 实现 ELR 主进程与沙箱管理器的通信协议

### 3. 沙箱启动/停止自动注册 ⏳

**需要实现：**
```go
// 沙箱启动时：
1. 调用 sandbox API
2. API 将沙箱添加到 RuntimeSandboxList
3. 保存沙箱ID和容器ID关联

// 沙箱停止时：
1. 调用 sandbox API
2. API 将沙箱从 RuntimeSandboxList 移除
```

### 4. 沙箱持久化改进 ⏳

**当前问题：**
- 沙箱数据在 `./sandbox-state.json`
- 没有统一的沙箱目录结构

**建议方案：**
- 沙箱应该存储在容器目录内
- 或创建统一的 `~/.elr/data/sandboxes/` 目录

### 5. 其他待修复问题 ⏳

- [ ] 沙箱日期格式化问题
- [ ] `sandbox list` 命令的沙箱目录读取
- [ ] 沙箱启动时自动注册到 RuntimeSandboxList
- [ ] 沙箱停止时自动从 RuntimeSandboxList 移除

---

## 五、架构总结

```
ELR 主进程
    │
    ├── 运行时容器列表 (RuntimeContainerList)
    │   └── containerID → RunningContainerInfo
    │
    ├── 运行时沙箱列表 (RuntimeSandboxList)
    │   └── sandboxID → RunningSandbox
    │
    ├── 沙箱-容器映射 (SandboxContainerManager)
    │   └── sandboxID → containerID
    │
    └── HTTP API (端口 16888)
        ├── /api/sandbox/list   (待完善)
        ├── /api/sandbox/start  (待完善)
        └── /api/sandbox/stop   (待完善)

容器进程
    └── 沙箱管理器子进程 (待完善 IPC 通信)
```

---

## 六、下一步计划

### 高优先级
1. 实现 `listSandboxes()` API 端点
2. 实现 `startSandbox()` API 端点（添加到运行时列表）
3. 实现沙箱启动时自动注册机制

### 中优先级
4. 实现 TCP/命名管道 IPC 通信
5. 完善沙箱持久化存储

### 低优先级
6. 修复格式化问题
7. 优化错误处理

---

## 七、文件清单

### 新增文件
- `elr/runtime_sandbox_list.go` - 运行时沙箱列表管理器
- `elr/sandbox_container_mapping.go` - 沙箱-容器映射管理器
- `elr/sandbox_manager.go` - 沙箱管理器子进程（框架）

### 修改文件
- `cli/main.go` - 沙箱命令和状态检查
- `elr/container.go` - 容器启动时自动启动沙箱管理器
- `elr/network.go` - 沙箱 API 端点
- `elr/runtime.go` - 初始化运行时沙箱列表

---

**开发日志结束**
**碳硅协同，共创未来**
**代码织梦者 & X54先生**
**2026年4月18日 8:30**
