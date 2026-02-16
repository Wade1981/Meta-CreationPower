# Enlightenment Lighthouse Runtime (ELR) 使用指南

## 1. 项目介绍

Enlightenment Lighthouse Runtime (ELR) 是启蒙灯塔起源团队开发的轻量级、跨平台容器运行环境，专为碳硅协同场景设计。它不依赖Docker，而是使用系统原生的隔离机制，提供了一个统一的、可扩展的平台，支持主流编程语言。

### 核心特性

- **轻量级设计**：核心运行时小于10MB，启动时间毫秒级
- **跨平台支持**：支持Windows、Linux和macOS
- **多语言支持**：内置支持C语言和Python，可扩展支持其他语言
- **零依赖部署**：在Windows系统上实现了完全独立的运行环境
- **完整容器管理**：支持容器的创建、启动、停止、删除和检查
- **状态持久化**：实现了跨命令的状态管理，确保运行时一致性

## 2. 安装指南

### 2.1 PowerShell实现（零依赖版本）

ELR提供了一个纯PowerShell实现的版本，无需任何外部依赖即可在Windows系统上运行。

**获取方式**：
1. 克隆GitHub仓库：
   ```bash
   git clone https://github.com/Wade1981/Meta-CreationPower.git
   ```
2. 进入ELR目录：
   ```bash
   cd Meta-CreationPower/05_Open_source_ProjectRepository/AIAgentFramework/EnlightenmentLighthouseRuntime
   ```

### 2.2 二进制版本（推荐用于生产环境）

**获取方式**：
1. 从GitHub Releases页面下载对应平台的二进制文件：
   - Windows: `elr-windows-amd64.exe`
   - Linux: `elr-linux-amd64`
   - macOS: `elr-darwin-amd64`

2. 安装：
   - **Windows**：将二进制文件复制到系统路径（如 `C:\Windows\System32`）
   - **Linux/macOS**：将二进制文件复制到系统路径（如 `/usr/local/bin`）并添加执行权限
     ```bash
     chmod +x elr-linux-amd64
     sudo mv elr-linux-amd64 /usr/local/bin/elr
     ```

## 3. 基本使用

### 3.1 启动ELR运行时

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 start

# 二进制版本
elr start
```

### 3.2 查看容器列表

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 list

# 二进制版本
elr list
```

### 3.3 创建容器

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 create --name my-container --image ubuntu:latest

# 二进制版本
elr create --name my-container --image ubuntu:latest
```

### 3.4 运行容器

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 run --name my-container --image ubuntu:latest

# 二进制版本
elr run --name my-container --image ubuntu:latest
```

### 3.5 启动容器

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 start-container --id <容器ID>

# 二进制版本
elr start-container --id <容器ID>
```

### 3.6 停止容器

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 stop-container --id <容器ID>

# 二进制版本
elr stop-container --id <容器ID>
```

### 3.7 删除容器

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 delete --id <容器ID>

# 二进制版本
elr delete --id <容器ID>
```

### 3.8 检查容器

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 inspect --id <容器ID>

# 二进制版本
elr inspect --id <容器ID>
```

### 3.9 执行命令

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 exec --id <容器ID> --command <命令>

# 二进制版本
elr exec --id <容器ID> --command <命令>
```

### 3.10 停止ELR运行时

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 stop

# 二进制版本
elr stop
```

## 4. C语言程序支持

ELR提供了专门的命令来编译和运行C语言程序。

### 4.1 基本用法

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 run-c --source hello.c

# 二进制版本
elr run-c --source hello.c
```

### 4.2 指定输出文件名

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 run-c --source hello.c --output hello.exe

# 二进制版本
elr run-c --source hello.c --output hello.exe
```

### 4.3 添加编译参数

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 run-c --source hello.c --args '-Wall -O2'

# 二进制版本
elr run-c --source hello.c --args '-Wall -O2'
```

### 4.4 前提条件

要使用run-c命令，需要安装gcc编译器：

**Windows系统**：
- MinGW-w64: https://www.mingw-w64.org/
- MSYS2: https://www.msys2.org/
- Cygwin: https://www.cygwin.com/

**Linux/macOS系统**：
```bash
# Ubuntu/Debian
sudo apt update && sudo apt install gcc

# CentOS/RHEL
sudo yum install gcc

# macOS (使用Homebrew)
brew install gcc
```

## 5. Python支持

ELR提供了专门的命令来运行Python脚本或直接执行Python代码。

### 5.1 运行Python脚本

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 run-python --source script.py

# 二进制版本
elr run-python --source script.py
```

### 5.2 直接执行Python代码

```powershell
# PowerShell实现版本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 run-python --code 'print("Hello from Python!")'

# 二进制版本
elr run-python --code 'print("Hello from Python!")'
```

### 5.3 前提条件

要使用run-python命令，需要安装Python 3.8或更高版本：

**安装方法**：
1. 从官方网站下载并安装：https://www.python.org/downloads/
2. 或使用Python便携版：https://www.python.org/downloads/windows/（选择Windows embeddable package）

### 5.4 Windows Store Python占位符检测

ELR会自动检测并拒绝Windows Store的Python占位符，确保使用的是实际的Python解释器。如果检测到占位符，会提供详细的安装建议。

## 6. 与Meta-CreationPower的结合

### 6.1 运行项目代码

```powershell
# 运行主程序
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 run-python --source src/main.py

# 运行测试脚本
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 run-python --source tests/test_main.py
```

### 6.2 容器化部署

```powershell
# 创建并运行容器
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 run --name meta-creationpower --image python:3.9

# 在容器中执行项目
powershell -ExecutionPolicy RemoteSigned -File elr.ps1 exec --id <容器ID> --command 'python src/main.py'
```

## 7. 高级功能

### 7.1 运行时状态管理

ELR使用JSON文件保存运行时和容器状态，确保跨命令会话的一致性。

**状态文件**：`elr-state.json`

### 7.2 错误处理

ELR提供了完善的错误处理机制，包括：
- 容器不存在时的错误提示
- 容器未运行时的错误提示
- 命令执行失败时的错误提示
- Python解释器未找到时的安装建议
- C编译器未找到时的安装建议

### 7.3 扩展功能

ELR采用插件架构，支持通过插件扩展功能：
- **语言插件**：添加对新编程语言的支持
- **服务插件**：添加新的服务和功能
- **平台插件**：支持新的操作系统

## 8. 最佳实践

### 8.1 容器管理最佳实践

- **使用有意义的容器名称**：便于识别和管理
- **定期清理未使用的容器**：释放系统资源
- **使用合适的镜像**：根据项目需求选择合适的基础镜像
- **监控容器状态**：定期检查容器运行状态

### 8.2 Python开发最佳实践

- **使用虚拟环境**：隔离项目依赖
- **遵循PEP 8编码规范**：提高代码可读性
- **添加类型提示**：提高代码可维护性
- **编写单元测试**：确保代码质量
- **使用版本控制**：跟踪代码变更

### 8.3 C语言开发最佳实践

- **使用头文件保护**：防止头文件重复包含
- **遵循C语言编码规范**：提高代码可读性
- **使用适当的编译器警告**：`-Wall -Wextra -Werror`
- **进行内存管理**：避免内存泄漏
- **编写单元测试**：确保代码质量

## 9. 常见问题

### 9.1 Python相关问题

**问题**：Error: Python interpreter not found
**解决方案**：安装Python 3.8或更高版本，并确保添加到系统路径

**问题**：Error: Found Windows Store Python placeholder
**解决方案**：从官方网站安装Python，而不是使用Windows Store版本

### 9.2 C语言相关问题

**问题**：Error: gcc compiler not found
**解决方案**：安装gcc编译器，如MinGW-w64、MSYS2或Cygwin

**问题**：Error: Compilation failed
**解决方案**：检查C代码语法，确保没有错误

### 9.3 容器相关问题

**问题**：Error: Container not found
**解决方案**：使用正确的容器ID，可通过`elr list`查看所有容器

**问题**：Error: Container is not running
**解决方案**：先启动容器，再执行命令

### 9.4 运行时问题

**问题**：Error: ELR runtime is not running
**解决方案**：先启动ELR运行时，再执行其他命令

**问题**：Error: ELR runtime is already running
**解决方案**：ELR运行时已经在运行，无需重复启动

## 10. 命令参考

### 10.1 基本命令

| 命令 | 描述 | 示例 |
|------|------|------|
| `version` | 显示版本信息 | `elr version` |
| `help` | 显示帮助信息 | `elr help` |
| `start` | 启动ELR运行时 | `elr start` |
| `stop` | 停止ELR运行时 | `elr stop` |
| `status` | 检查运行时状态 | `elr status` |

### 10.2 容器管理命令

| 命令 | 描述 | 示例 |
|------|------|------|
| `create` | 创建新容器 | `elr create --name my-container --image ubuntu:latest` |
| `run` | 创建并启动新容器 | `elr run --name my-container --image ubuntu:latest` |
| `start-container` | 启动容器 | `elr start-container --id elr-1234567890` |
| `stop-container` | 停止容器 | `elr stop-container --id elr-1234567890` |
| `list` | 列出所有容器 | `elr list` |
| `delete` | 删除容器 | `elr delete --id elr-1234567890` |
| `inspect` | 检查容器 | `elr inspect --id elr-1234567890` |
| `exec` | 在容器中执行命令 | `elr exec --id elr-1234567890 --command 'ls -la'` |

### 10.3 语言支持命令

| 命令 | 描述 | 示例 |
|------|------|------|
| `run-c` | 编译并运行C程序 | `elr run-c --source hello.c` |
| `run-python` | 运行Python脚本或代码 | `elr run-python --source script.py` |

## 11. 总结

Enlightenment Lighthouse Runtime (ELR) 是一个轻量级、跨平台的容器运行环境，为碳硅协同场景提供了强大的支持。它不仅支持基本的容器管理，还内置了对C语言和Python的支持，可通过插件扩展支持其他语言。

通过本指南，您应该已经了解了ELR的基本使用方法和高级功能。ELR的设计理念是提供一个简单、高效、可扩展的平台，为碳硅协同创新提供基础设施。

### 未来展望

- **更多语言支持**：计划添加对JavaScript、Java、Go等语言的支持
- **网络功能**：添加容器网络支持，实现容器间通信
- **存储功能**：添加持久化存储支持
- **分布式能力**：支持多节点部署和集群管理
- **图形界面**：提供可视化的管理界面

ELR将继续发展，成为碳硅协同创新的重要基础设施，推动人文价值与科技理性的平衡共生。

---

**版本**：1.0.0
**更新时间**：2026年2月12日
**项目**：Meta-CreationPower
**作者**：启蒙灯塔起源团队
