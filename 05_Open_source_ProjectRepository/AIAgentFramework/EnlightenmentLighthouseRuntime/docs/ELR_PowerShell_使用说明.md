# ELR PowerShell 使用说明

## 概述

ELR (Enlightenment Lighthouse Runtime) 是一个轻量级的运行时环境，提供了容器管理、模型管理、沙箱管理、API 服务等功能。本说明文档详细介绍了如何使用 PowerShell 脚本 `elr.ps1` 来管理 ELR 环境。

## 基本使用

### 查看帮助信息

```powershell
.\elr.ps1 help
```

### 查看版本信息

```powershell
.\elr.ps1 version
```

## 容器管理

### 创建容器

```powershell
# 基本创建
.\elr.ps1 create <容器名称>

# 完整参数创建
.\elr.ps1 create --name <容器名称> --image <镜像名称> [--fs-isolation] [--rootfs <根文件系统路径>] [--read-only]
```

### 运行容器

```powershell
# 创建并启动容器
.\elr.ps1 run --name <容器名称> --image <镜像名称>

# 运行 Python 脚本
.\elr.ps1 run python <脚本路径> [参数...]
```

### 启动/停止容器

```powershell
# 启动容器
.\elr.ps1 start-container --id <容器ID>

# 停止容器
.\elr.ps1 stop-container --id <容器ID>
```

### 列出容器

```powershell
.\elr.ps1 list
```

### 删除容器

```powershell
.\elr.ps1 delete --id <容器ID>
```

### 查看容器详情

```powershell
.\elr.ps1 inspect --id <容器ID>
```

## 模型管理

### 列出所有模型

```powershell
.\elr.ps1 model list
```

### 获取模型信息

```powershell
.\elr.ps1 model get --model-id <模型ID>
```

### 下载模型

```powershell
.\elr.ps1 model download --model-id <模型ID> --type <模型类型> --url <下载地址>
```

### 删除模型

```powershell
.\elr.ps1 model delete --model-id <模型ID>
```

### 安装模型依赖

```powershell
.\elr.ps1 model install-deps --model-id <模型ID> --type <依赖类型>
```

## 沙箱管理

### 列出所有沙箱

```powershell
.\elr.ps1 sandbox list
```

### 创建沙箱

```powershell
.\elr.ps1 sandbox create --container <容器名称>
```

### 启动/停止沙箱

```powershell
# 启动沙箱
.\elr.ps1 sandbox start --sandbox-id <沙箱ID>

# 停止沙箱
.\elr.ps1 sandbox stop --sandbox-id <沙箱ID>
```

### 删除沙箱

```powershell
.\elr.ps1 sandbox delete --sandbox-id <沙箱ID>
```

### 加载模型到沙箱

```powershell
.\elr.ps1 sandbox load-model --sandbox-id <沙箱ID> --model-id <模型ID>
```

### 从沙箱卸载模型

```powershell
.\elr.ps1 sandbox unload-model --sandbox-id <沙箱ID> --model-id <模型ID>
```

### 在沙箱中运行模型

```powershell
.\elr.ps1 sandbox run-model --sandbox-id <沙箱ID> --model-id <模型ID> --input <输入内容>
```

## API 管理

### 启动 API 服务

```powershell
# 启动所有 API 服务
.\elr.ps1 api start

# 启动特定 API 服务
.\elr.ps1 api start desktop
.\elr.ps1 api start public
.\elr.ps1 api start model
```

### 停止 API 服务

```powershell
# 停止所有 API 服务
.\elr.ps1 api stop

# 停止特定 API 服务
.\elr.ps1 api stop desktop
.\elr.ps1 api stop public
.\elr.ps1 api stop model
```

### 检查 API 服务状态

```powershell
.\elr.ps1 api status
```

### 配置 API 设置

```powershell
# 配置 API 地址和端口
.\elr.ps1 api config set --api-type <api类型> --address <地址> --port <端口>

# 查看 API 配置
.\elr.ps1 api config list
```

## 文件系统管理

### 上传文件到容器

```powershell
.\elr.ps1 fs upload --local-path <本地文件路径> --container-path <容器文件路径>
```

### 从容器下载文件

```powershell
.\elr.ps1 fs download --container-path <容器文件路径> --local-path <本地文件路径>
```

### 设置文件类型目录

```powershell
.\elr.ps1 fs set-dir --file-type <文件类型> --directory <目录路径>
```

### 获取文件类型目录

```powershell
.\elr.ps1 fs get-dir --file-type <文件类型>
```

## 系统管理

### 启动/停止 ELR 运行时

```powershell
# 启动 ELR 运行时
.\elr.ps1 start

# 停止 ELR 运行时
.\elr.ps1 stop
```

### 系统设置

```powershell
.\elr.ps1 setup --isolation <隔离类型>
```

### 资源配置

```powershell
.\elr.ps1 Settings
```

## 依赖管理

### 安装 Python

```powershell
# 安装默认版本的 Python
.\elr.ps1 install python

# 安装指定版本的 Python
.\elr.ps1 install python <版本> <安装路径>
```

## GUI 管理

### 启动 ELR 托盘应用

```powershell
# 启动 ELR 托盘应用
.\elr.ps1 gui
# 或
.\elr.ps1 tray
```

## 示例

### 示例 1: 创建并启动容器

```powershell
# 创建容器
.\elr.ps1 create my-container

# 启动容器
.\elr.ps1 start-container --id <容器ID>
```

### 示例 2: 管理模型

```powershell
# 下载模型
.\elr.ps1 model download --model-id gpt2 --type text --url https://example.com/gpt2

# 加载模型到沙箱
.\elr.ps1 sandbox create --container my-container
.\elr.ps1 sandbox load-model --sandbox-id <沙箱ID> --model-id gpt2

# 在沙箱中运行模型
.\elr.ps1 sandbox run-model --sandbox-id <沙箱ID> --model-id gpt2 --input "Hello, world!"
```

### 示例 3: 管理 API 服务

```powershell
# 配置 API 地址和端口
.\elr.ps1 api config set --api-type desktop --address localhost --port 8081

# 启动 API 服务
.\elr.ps1 api start desktop

# 检查 API 服务状态
.\elr.ps1 api status

# 停止 API 服务
.\elr.ps1 api stop desktop
```

## 注意事项

1. 所有命令都需要在 ELR 目录下执行
2. 部分命令需要管理员权限
3. 容器和沙箱操作需要 ELR 运行时处于启动状态
4. API 服务启动后会在后台运行
5. 更多命令详情请使用 `.\elr.ps1 help` 查看

## 故障排除

### 常见错误及解决方案

1. **Error: elr.exe not found**
   - 解决方案: 编译 Go 代码生成 elr.exe
   ```powershell
   go build -o elr.exe cli/main.go
   ```

2. **Error: Container name is required**
   - 解决方案: 使用 `--name` 参数指定容器名称
   ```powershell
   .\elr.ps1 create --name <容器名称> --image <镜像名称>
   ```

3. **Error: Python interpreter not found**
   - 解决方案: 安装 Python
   ```powershell
   .\elr.ps1 install python
   ```

4. **API 服务无法启动**
   - 解决方案: 检查端口是否被占用，或使用不同的端口
   ```powershell
   .\elr.ps1 api config set --api-type <api类型> --address <地址> --port <不同端口>
   ```

## 总结

ELR PowerShell 脚本提供了丰富的命令来管理 ELR 环境，包括容器管理、模型管理、沙箱管理、API 服务管理等。通过本说明文档，您应该能够掌握基本的 ELR 管理操作，为您的开发和部署工作提供便利。

如有任何问题，请参考 ELR 官方文档或联系技术支持。