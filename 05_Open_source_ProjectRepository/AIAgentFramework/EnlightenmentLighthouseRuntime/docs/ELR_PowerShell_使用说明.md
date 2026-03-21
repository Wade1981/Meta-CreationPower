# ELR PowerShell 命令使用说明

## 概述

Enlightenment Lighthouse Runtime (ELR) 是启蒙灯塔起源团队开发的运行时环境，提供容器管理、令牌管理和网络服务管理等功能。本文档介绍 ELR PowerShell 命令的使用方法。

## 版本信息

- **当前版本**: v1.6.0
- **平台**: Windows
- **实现方式**: PowerShell
- **依赖**: 无外部依赖

## 安装与配置

### 系统要求

- Windows 10 或更高版本
- PowerShell 5.1 或更高版本

### 安装步骤

1. 将 `elr.ps1` 文件放置到目标目录
2. 确保 `elr\token_manager.ps1` 文件存在于同一目录下的 `elr` 文件夹中
3. 以管理员权限运行 PowerShell

### 执行策略设置

首次使用前，需要设置 PowerShell 执行策略：

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## 命令列表

### 基础命令

| 命令 | 描述 |
|------|------|
| `version` | 打印版本信息 |
| `help` | 显示帮助信息 |
| `start` | 启动 ELR 运行时 |
| `stop` | 停止 ELR 运行时 |
| `status` | 检查运行时状态 |
| `list` | 列出所有容器 |
| `stats` | 显示容器资源使用统计 |
| `tray` | 启动 ELR 托盘应用 |

### 网络服务管理命令

| 命令 | 描述 |
|------|------|
| `start-all` | 启动所有网络服务 |
| `stop-all` | 停止所有网络服务 |
| `start-desktop` | 启动 Desktop API (端口 8081) |
| `stop-desktop` | 停止 Desktop API |
| `start-public` | 启动 Public API (端口 8080) |
| `stop-public` | 停止 Public API |
| `start-model` | 启动 Model Service (端口 8082) |
| `stop-model` | 停止 Model Service |
| `start-micro` | 启动 Micro Model Server (端口 8083) |
| `stop-micro` | 停止 Micro Model Server |
| `network-list` | 查看可用的 IP 地址和端口 |

### 令牌管理命令

| 命令 | 描述 |
|------|------|
| `token` | 管理 ELR 容器令牌 |
| `network-status` | 检查网络状态 |

### 容器管理命令

| 命令 | 描述 |
|------|------|
| `exec` | 在容器中执行命令 |
| `upload` | 上传文件到容器 |

## 详细命令说明

### 1. 版本信息

显示 ELR 版本和平台信息。

```powershell
.\elr.ps1 version
```

**输出示例：**

```
Enlightenment Lighthouse Runtime v1.4.0
Platform: Windows
PowerShell Implementation
No external dependencies required
```

### 2. 帮助信息

显示所有可用命令和用法。

```powershell
.\elr.ps1 help
```

**输出示例：**

```
Enlightenment Lighthouse Runtime (ELR)
Usage: elr [command] [options]

Basic Commands:
  version           Print version information
  help              Print this help message
  start             Start the ELR runtime
  stop              Stop the ELR runtime
  status            Check the runtime status
  list              List all containers

Network Service Commands:
  start-all         Start all network services
  stop-all          Stop all network services
  start-desktop     Start Desktop API (port 8081)
  stop-desktop      Stop Desktop API
  start-public      Start Public API (port 8080)
  stop-public       Stop Public API
  start-model       Start Model Service (port 8082)
  stop-model        Stop Model Service
  start-micro       Start Micro Model Server (port 8083)
  stop-micro        Stop Micro Model Server

Token Commands:
  token             Manage ELR container tokens
  network-status    Check network status
```

### 3. 启动运行时

启动 ELR 运行时环境和容器。

```powershell
.\elr.ps1 start
```

**输出示例：**

```
====================================
Starting Enlightenment Lighthouse Runtime v1.4.0
Platform: Windows
====================================
Initializing platform...
Loading plugins...
Loading containers...
====================================
Starting ELR container...
ELR container started successfully!

Enlightenment Lighthouse Runtime started successfully!
====================================
```

**说明：**

- 如果 `elr\elr-container.exe` 存在，将启动实际的容器进程
- 如果容器可执行文件不存在，将使用 PowerShell 模拟模式
- 启动后状态将保存到 `elr-state.json` 文件

### 4. 停止运行时

停止 ELR 运行时环境和容器。

```powershell
.\elr.ps1 stop
```

**输出示例：**

```
====================================
Stopping Enlightenment Lighthouse Runtime...
Stopping containers...
Cleaning up plugins...
Cleaning up platform...
====================================
Stopping ELR container...
Stopped ELR container process: 12345
Enlightenment Lighthouse Runtime stopped successfully!
====================================
```

**说明：**

- 将自动查找并停止名为 `elr-container` 的进程
- 状态将更新到 `elr-state.json` 文件

### 5. 检查状态

检查 ELR 运行时的当前状态。

```powershell
.\elr.ps1 status
```

**输出示例（运行中）：**

```
Enlightenment Lighthouse Runtime is RUNNING
Started: 2026-03-19 14:30:00
Containers: 2
Running containers: 1
```

**输出示例（未运行）：**

```
Error: ELR runtime is not running
```

### 6. 列出容器

列出所有容器及其状态。

```powershell
.lr.ps1 list
```

**输出示例：**

```
====================================
Containers:
====================================
ID                 NAME            IMAGE           STATUS    CREATED
--                 ----            -----           ------    -------
elr-1234567890     test-container  ubuntu:latest   created   2026-03-19 14:30:00
elr-0987654321     python-app      python:3.9      running   2026-03-19 14:30:00
====================================
```

### 7. 容器资源使用统计

显示容器的资源使用情况，包括内存、CPU和GPU使用。

```powershell
.lr.ps1 stats
```

**输出示例：**

```
====================================
Container Stats:
====================================
ID                 NAME            MEMORY    CPU     GPU        
--                 ----            ------    ---     ---
elr-1234567890 test-container   0MB       0%     0%
elr-0987654321 python-app       8MB       0%     0%
elr-1234567900 micro-model      28MB      0%     0%
====================================
```

**说明：**

- 显示硬编码容器的真实资源消耗
- 显示实际运行进程的资源消耗
- 对于多个进程对应同一容器的情况，会汇总资源消耗

### 8. 启动托盘应用

启动 ELR 托盘应用，提供图形化界面管理 ELR 容器和服务。

```powershell
.\elr.ps1 tray
```

**输出示例：**

```
====================================
Starting ELR Tray Application...
====================================
Starting ELR Tray Application in background...
Path: E:\X54\github\Meta-CreationPower\ELR-Tray-App.ps1
ELR Tray Application started successfully!
You can find the ELR icon in the system tray.
====================================
```

**说明：**

- 托盘应用会在系统托盘中显示 ELR 图标
- 右键点击图标可访问各种功能
- 托盘应用提供图形化界面，包括容器管理、服务管理、设置配置等功能
- 支持多语言界面，自动根据系统语言切换中英文

---

## 网络服务管理

ELR 提供四个网络服务，可以通过 PowerShell 命令进行启动和停止。

### 服务概览

| 服务名称 | 端口 | 描述 |
|----------|------|------|
| Desktop API | 8081 | 桌面应用程序接口，为 ELR-Desktop-Assistant 提供专属通信通道 |
| Public API | 8080 | 公共访问接口，提供通用 API 服务 |
| Model Service | 8082 | 模型服务，提供模型推理功能 |
| Micro Model Server | 8083 | 微模型服务器，提供轻量级模型服务 |

### 1. 启动所有网络服务

一键启动所有网络服务。

```powershell
.\elr.ps1 start-all
```

**输出示例：**

```
====================================
Starting All Network Services...
====================================
====================================
Starting Desktop API (port 8081)...
====================================
Starting Desktop API in a new window...
Desktop API started!
Address: http://localhost:8081
====================================
====================================
Starting Public API (port 8080)...
====================================
Starting Public API in a new window...
Public API started!
Address: http://localhost:8080
====================================
====================================
Starting Model Service (port 8082)...
====================================
Starting Model Service in a new window...
Model Service started!
Address: http://localhost:8082
====================================
====================================
Starting Micro Model Server (port 8083)...
====================================
Starting Micro Model Server in a new window...
Micro Model Server started!
Address: http://localhost:8083
====================================
====================================
All services started!
====================================
```

### 2. 停止所有网络服务

一键停止所有网络服务。

```powershell
.\elr.ps1 stop-all
```

**输出示例：**

```
====================================
Stopping All Network Services...
====================================
====================================
Stopping Desktop API...
====================================
Stopped process on port 8081: 12345
====================================
====================================
Stopping Public API...
====================================
Stopped process on port 8080: 12346
====================================
====================================
Stopping Model Service...
====================================
Stopped process on port 8082: 12347
====================================
====================================
Stopping Micro Model Server...
====================================
Stopped process on port 8083: 12348
====================================
====================================
All services stopped!
====================================
```

### 3. Desktop API 管理

**启动 Desktop API：**

```powershell
.\elr.ps1 start-desktop
```

**停止 Desktop API：**

```powershell
.\elr.ps1 stop-desktop
```

**可用端点：**

| 方法 | 端点 | 描述 |
|------|------|------|
| GET | `/api/desktop/status` | 获取 ELR 状态 |
| GET | `/api/desktop/containers` | 获取容器列表 |
| POST | `/api/desktop/upload` | 上传文件 |
| GET | `/api/desktop/files` | 列出上传的文件 |
| DELETE | `/api/desktop/files/{name}` | 删除文件 |
| GET | `/api/desktop/resources` | 获取系统资源使用情况 |
| GET | `/api/desktop/health` | 健康检查 |

### 4. Public API 管理

**启动 Public API：**

```powershell
.\elr.ps1 start-public
```

**停止 Public API：**

```powershell
.\elr.ps1 stop-public
```

### 5. Model Service 管理

**启动 Model Service：**

```powershell
.\elr.ps1 start-model
```

**停止 Model Service：**

```powershell
.\elr.ps1 stop-model
```

### 6. Micro Model Server 管理

**启动 Micro Model Server：**

```powershell
.\elr.ps1 start-micro
```

**停止 Micro Model Server：**

```powershell
.\elr.ps1 stop-micro
```

**说明：**

- 如果 `micro_model_server.exe` 不存在，将自动尝试编译 Go 代码
- 需要安装 Go 语言环境才能编译

### 7. 查看网络状态

检查所有网络服务的运行状态。

```powershell
.\elr.ps1 network-status
```

**输出示例：**

```
====================================
ELR Container Network Status
====================================
Desktop API: Running
  Address: http://localhost:8081
Public API: Running
  Address: http://localhost:8080
Model Service: Running
  Address: http://localhost:8082
Micro Model Server: Running
  Address: http://localhost:8083
====================================
```

### 8. 查看网络列表

查看可用的 IP 地址和端口配置。

```powershell
.\elr.ps1 network-list
```

**输出示例：**

```
====================================
ELR Network - Available IPs and Ports
====================================

Available IP Addresses:
----------------------
  localhost
  127.0.0.1
  192.168.1.100

Default Port Configuration:
--------------------------
  Desktop API:      localhost:8081
  Public API:       localhost:8080
  Model Service:    localhost:8082
  Micro Model:      localhost:8083

Currently Listening Ports:
-------------------------
  127.0.0.1:8081 - python
  127.0.0.1:8080 - elr-container
  127.0.0.1:8082 - python
  127.0.0.1:8083 - micro_model_server

Usage Examples:
---------------
  .\elr.ps1 start-desktop                    # Use default localhost:8081
  .\elr.ps1 start-desktop 192.168.1.100:9081 # Custom IP and Port
  .\elr.ps1 start-desktop 0.0.0.0:8081       # Listen on all interfaces
  .\elr.ps1 start-desktop :9081              # Custom Port only
====================================
```

---

## 容器管理命令

### 1. 在容器中执行命令

在 ELR 容器中执行指定命令。

```powershell
.\elr.ps1 exec --command "命令内容"
```

**示例：**

```powershell
.\elr.ps1 exec --command "echo Hello, ELR!"
```

**输出示例：**

```
====================================
Executing command in ELR container:
echo Hello, ELR!
====================================
Hello, ELR!
====================================
Command executed successfully!
====================================
```

**说明：**

- 命令会在本地执行，模拟容器执行环境
- 支持任何有效的 PowerShell 命令

### 2. 上传文件到容器

上传文件到 ELR 容器中。

```powershell
.\elr.ps1 upload --file "文件路径"
```

**示例：**

```powershell
.\elr.ps1 upload --file "C:\path\to\file.txt"
```

**输出示例：**

```
====================================
Uploading file to ELR container:
C:\path\to\file.txt
====================================
File uploaded successfully!
Destination: E:\X54\github\Meta-CreationPower\05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\elr\uploads\file.txt
====================================
```

**说明：**

- 文件会被上传到 `elr\uploads` 目录
- 如果上传目录不存在，会自动创建
- 支持上传任何类型的文件

---

## 令牌管理

令牌管理功能用于创建、验证、刷新和撤销 ELR 容器的访问令牌。

### 令牌管理命令

| 动作 | 描述 |
|------|------|
| `create` | 创建新令牌 |
| `validate` | 验证令牌 |
| `refresh` | 刷新令牌 |
| `list` | 列出所有令牌 |
| `revoke` | 撤销令牌 |

### 1. 创建令牌

创建一个新的 ELR 容器访问令牌。

```powershell
.\elr.ps1 token --action create --description "令牌描述"
```

**示例：**

```powershell
.\elr.ps1 token --action create --description "Admin Token"
```

**输出示例：**

```
====================================
New Token Generated:
48f2706b617b43edb9a2394c710b115615cecc66
====================================
Token ID: ee9ba256
Valid for: 7 days
Please save this token securely
```

**说明：**

- 令牌由两个 GUID 组合生成，共 40 个字符
- 令牌默认有效期为 7 天
- 令牌数据保存在 `elr\elr_token.json` 文件中
- 请妥善保存生成的令牌，丢失后无法恢复

### 2. 验证令牌

验证令牌是否有效。

```powershell
.\elr.ps1 token --action validate --token "令牌值"
```

**示例：**

```powershell
.\elr.ps1 token --action validate --token "48f2706b617b43edb9a2394c710b115615cecc66"
```

**输出示例（有效）：**

```
====================================
Token Validation Result:
Status: Valid
Message: Token is valid
====================================
```

**输出示例（已过期）：**

```
====================================
Token Validation Result:
Status: Invalid
Message: Token has expired
====================================
```

**输出示例（已撤销）：**

```
====================================
Token Validation Result:
Status: Invalid
Message: Token has been revoked
====================================
```

### 3. 列出所有令牌

列出所有已创建的令牌及其状态。

```powershell
.\elr.ps1 token --action list
```

**输出示例：**

```
====================================
Token List:
====================================
ID       | Description         | Status  | Created
-------- | ------------------- | ------- | --------
ee9ba256 | Admin Token         | Valid   | 2026-03-19
a1b2c3d4 | User Token          | Expired | 2026-03-12
f5e6d7c8 | Test Token          | Revoked | 2026-03-10
====================================
```

**状态说明：**

| 状态 | 描述 |
|------|------|
| Valid | 令牌有效，可以正常使用 |
| Expired | 令牌已过期，需要刷新或重新创建 |
| Revoked | 令牌已被撤销，无法使用 |

### 4. 刷新令牌

刷新现有令牌，生成新的令牌值并延长有效期。

```powershell
.\elr.ps1 token --action refresh --token "旧令牌" --description "新描述"
```

**示例：**

```powershell
.\elr.ps1 token --action refresh --token "48f2706b617b43edb9a2394c710b115615cecc66" --description "Refreshed Admin Token"
```

**输出示例：**

```
====================================
Token refreshed successfully:
a1b2c3d4e5f6789012345678901234567890abcd
====================================
Valid for: 7 days
Please save this token securely
```

**说明：**

- 刷新后会生成新的令牌值
- 原令牌将失效
- 有效期重置为 7 天
- 已撤销的令牌无法刷新

### 5. 撤销令牌

撤销指定令牌，使其无法再使用。

```powershell
.\elr.ps1 token --action revoke --token "令牌ID"
```

**示例：**

```powershell
.\elr.ps1 token --action revoke --token "ee9ba256"
```

**输出示例：**

```
====================================
Token revocation result:
Status: Success
Message: Token has been revoked
====================================
```

**说明：**

- 撤销操作使用令牌 ID（8 位字符），而非完整令牌值
- 撤销后令牌无法恢复
- 如需重新使用，需要创建新令牌

---

## 文件结构

```
EnlightenmentLighthouseRuntime/
├── elr.ps1                    # ELR 主脚本
├── elr-state.json             # 运行时状态文件
├── elr-new.ps1                # ELR 新版本脚本（备用）
├── elr_api_server.py          # Public API 服务脚本
└── elr/
    ├── elr-container.exe      # ELR 容器可执行文件
    ├── token_manager.ps1      # 令牌管理脚本
    ├── elr_token.json         # 令牌数据文件
    ├── desktop_api.py         # Desktop API 服务脚本
    └── start_desktop_api.ps1  # Desktop API 启动脚本
└── micro_model/
    ├── main.go                # Micro Model Server 主程序
    ├── micro_model_server.exe # Micro Model Server 可执行文件
    └── python_server.py       # Model Service Python 脚本
```

---

## 状态持久化

ELR 使用 JSON 文件保存运行时状态，确保在不同 PowerShell 会话之间保持状态一致性。

### 状态文件

| 文件 | 描述 |
|------|------|
| `elr-state.json` | 保存运行时启动状态和时间 |
| `elr\elr_token.json` | 保存所有令牌数据 |

### 状态文件格式

**elr-state.json 示例：**

```json
{
    "RUNTIME_STARTED": true,
    "RUNTIME_START_TIME": "2026-03-19T14:30:00.0000000+08:00"
}
```

**elr_token.json 示例：**

```json
{
    "tokens": [
        {
            "id": "ee9ba256",
            "token": "48f2706b617b43edb9a2394c710b115615cecc66",
            "description": "Admin Token",
            "created": 1710835800,
            "expires": 1711439400,
            "revoked": false
        }
    ],
    "last_updated": 1710835800
}
```

---

## 常见问题

### 1. 执行策略错误

**问题：** 运行脚本时提示"无法加载文件，因为在此系统上禁止运行脚本"

**解决方案：**

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### 2. 令牌管理脚本未找到

**问题：** 提示"Token manager script not found"

**解决方案：**

确保 `elr\token_manager.ps1` 文件存在于正确位置。

### 3. 容器无法启动

**问题：** 启动时提示"ELR container executable not found"

**解决方案：**

- 这是正常提示，表示使用 PowerShell 模拟模式
- 如需实际容器功能，请确保 `elr\elr-container.exe` 存在

### 4. 令牌验证失败

**问题：** 验证令牌时提示"Token not found"

**解决方案：**

- 检查令牌值是否正确复制
- 确认令牌未被撤销
- 检查 `elr\elr_token.json` 文件是否存在

### 5. 网络服务启动失败

**问题：** 启动网络服务时提示脚本未找到

**解决方案：**

- Desktop API: 确保 `elr\desktop_api.py` 存在
- Public API: 确保 `elr_api_server.py` 存在
- Model Service: 确保 `micro_model\python_server.py` 存在
- Micro Model Server: 确保 `micro_model\main.go` 存在或已编译

### 6. Python 未找到

**问题：** 启动服务时提示"Python not found"

**解决方案：**

- 安装 Python 3.8 或更高版本
- 确保 Python 已添加到系统 PATH
- 避免使用 Windows Store 版本的 Python

---

## 安全建议

1. **令牌保护**
   - 不要在公共场合分享令牌
   - 定期刷新重要令牌
   - 不再使用的令牌应及时撤销

2. **文件权限**
   - 确保 `elr_token.json` 文件权限设置正确
   - 避免将令牌文件提交到版本控制系统

3. **运行时安全**
   - 仅在可信环境中启动 ELR 运行时
   - 定期检查运行时状态

---

## 更新日志

### v1.6.0 (2026-03-21)

- 新增 `stats` 命令，显示容器资源使用统计
- 新增 `network-list` 命令，查看可用的 IP 地址和端口
- 新增 `exec` 命令，在容器中执行命令
- 新增 `upload` 命令，上传文件到容器
- 优化容器资源使用统计功能，支持真实进程资源监控
- 支持多个进程对应同一容器的资源汇总

### v1.4.0 (2026-03-19)

- 添加网络服务管理功能
- 新增 `start-all` 和 `stop-all` 命令
- 新增单独服务启动/停止命令
- 新增 Model Service 和 Micro Model Server 状态检查
- 优化网络状态查询功能

### v1.3.0 (2026-03-19)

- 添加状态持久化功能
- 添加令牌管理功能
- 添加网络状态查询功能
- 修复语法兼容性问题
- 支持 PowerShell 5.1

---

## 联系与支持

如有问题或建议，请联系启蒙灯塔起源团队。

---

## ELR 托盘应用

ELR 托盘应用（ELR-Tray-App.ps1）是一个系统托盘工具，提供图形化界面管理 ELR 容器和服务。

### 启动方法

1. **直接运行**
   
   在 PowerShell 中执行：
   
   ```powershell
   # 在项目根目录执行
   .\ELR-Tray-App.ps1
   
   # 或使用完整路径
   & "E:\X54\github\Meta-CreationPower\ELR-Tray-App.ps1"
   ```

2. **设置执行策略**
   
   如果遇到执行策略限制，请先设置执行策略：
   
   ```powershell
   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
   ```

### 功能特性

- **系统托盘集成**：在系统托盘中显示 ELR 图标，提供快速访问菜单
- **多语言支持**：自动根据系统语言切换中英文界面
- **容器管理**：查看容器状态和管理容器
- **服务管理**：启动/停止所有网络服务
- **设置界面**：配置 API 地址和端口，查看管理员令牌信息
- **容器对话**：与 ELR 容器进行交互对话

### 使用说明

1. **系统托盘菜单**
   
   右键点击系统托盘中的 ELR 图标，可访问以下功能：
   
   - **打开 ELR 助手**：显示主窗口，查看运行状态和容器列表
   - **ELR 容器对话**：打开对话窗口，与容器进行交互
   - **ELR 容器设置**：打开设置窗口，配置 API 地址和查看令牌信息
   - **网络状态**：检查网络服务运行状态
   - **启动所有服务**：一键启动所有网络服务
   - **停止所有服务**：一键停止所有网络服务
   - **退出**：退出托盘应用

2. **设置窗口**
   
   在设置窗口中，您可以：
   
   - 查看当前管理员令牌信息
   - 配置 Desktop API、Public API、Model Service 和 Micro Model 的地址和端口
   - 一键启动/停止各个网络服务
   - 查看管理员数量和管理列表

3. **对话窗口**
   
   在对话窗口中，您可以：
   
   - 与 ELR 容器进行文本对话
   - 上传文件到容器
   - 查看对话历史

### 配置文件

- **语言配置**：`language.json` 文件存储界面文本的中英文翻译
- **设置持久化**：API 配置会保存到配置文件中

### 系统要求

- Windows 10 或更高版本
- PowerShell 5.1 或更高版本
- .NET Framework 4.7.2 或更高版本（用于 Windows Forms）

---

**文档版本**: v1.3  
**最后更新**: 2026-03-21  
**作者**: 代码织梦者 (Code Weaver)