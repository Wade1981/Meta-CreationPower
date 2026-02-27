# PowerShell 5 语法指南：替代 && 操作符

## 问题分析

您遇到的错误：`标记"&&"不是此版本中的有效语句分隔符` 是因为您正在使用 PowerShell 5.1，而 `&&` 操作符是在 PowerShell 7+ 中才引入的。

## 解决方案

### 1. 使用分号 (`;`) 替代 &&

在 PowerShell 5 中，最简单的替代方法是使用分号来分隔多个命令：

**错误示例**：
```powershell
mkdir build && cd build && cmake ..
```

**正确示例**：
```powershell
mkdir build; cd build; cmake ..
```

### 2. 保持错误停止行为

如果您需要像 `&&` 一样在命令失败时停止执行后续命令，可以使用 `$ErrorActionPreference` 或 `Try-Catch` 块：

**方法 1：使用 `$ErrorActionPreference`**
```powershell
$ErrorActionPreference = "Stop"
mkdir build; cd build; cmake ..
```

**方法 2：使用 `Try-Catch`**
```powershell
try {
    mkdir build
    cd build
    cmake ..
} catch {
    Write-Host "Error: $($_.Exception.Message)"
    return
}
```

### 3. 常用命令转换示例

#### 示例 1：创建目录并进入
```powershell
# 错误
mkdir myproject && cd myproject

# 正确
mkdir myproject; cd myproject
```

#### 示例 2：编译并运行
```powershell
# 错误
gcc hello.c -o hello && ./hello

# 正确
gcc hello.c -o hello; ./hello
```

#### 示例 3：运行多个 PowerShell 脚本
```powershell
# 错误
.cript1.ps1 && .cript2.ps1

# 正确
.cript1.ps1; .cript2.ps1
```

## 临时解决方案：使用 cmd.exe

如果您需要使用 `&&` 操作符，可以临时切换到 cmd.exe：

```powershell
cmd.exe /c "mkdir build && cd build && cmake .."
```

## 永久解决方案：升级到 PowerShell 7+

如果您经常需要使用现代 PowerShell 功能，建议升级到 PowerShell 7+：

1. 从 [PowerShell 官方网站](https://github.com/PowerShell/PowerShell/releases) 下载安装包
2. 运行安装程序
3. 安装完成后，使用 `pwsh` 命令启动 PowerShell 7+

## 如何检查 PowerShell 版本

```powershell
$PSVersionTable.PSVersion
```

## 注意事项

- PowerShell 5 不支持以下现代操作符：
  - `&&`（逻辑与）
  - `||`（逻辑或）
  - `??`（空值合并）
  - `??=`（空值合并赋值）

- 所有这些操作符都需要使用 PowerShell 7+ 或使用替代语法。

## 推荐做法

对于跨版本兼容性，建议：
1. 使用分号 (`;`) 分隔命令
2. 添加适当的错误处理
3. 避免使用 PowerShell 7+ 特有的语法
4. 在脚本开头添加版本检查和提示