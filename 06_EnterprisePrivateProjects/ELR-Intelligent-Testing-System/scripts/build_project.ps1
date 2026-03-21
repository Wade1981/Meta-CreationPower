#!/usr/bin/env pwsh

# 构建项目脚本

Write-Host "开始构建 ELR 智能测试系统..."

# 创建构建目录
if (-not (Test-Path "build")) {
    New-Item -ItemType Directory -Path "build" | Out-Null
}

# 进入构建目录
Set-Location "build"

# 运行 CMake 配置
Write-Host "配置 CMake..."
cmake .. -G "Visual Studio 16 2019" -A x64

if ($LASTEXITCODE -ne 0) {
    Write-Host "CMake 配置失败!" -ForegroundColor Red
    Set-Location ..
    exit 1
}

# 构建项目
Write-Host "构建项目..."
cmake --build . --config Release

if ($LASTEXITCODE -ne 0) {
    Write-Host "构建失败!" -ForegroundColor Red
    Set-Location ..
    exit 1
}

Write-Host "构建成功!" -ForegroundColor Green

# 复制配置文件
Write-Host "复制配置文件..."
if (-not (Test-Path "configs")) {
    New-Item -ItemType Directory -Path "configs" | Out-Null
}

Copy-Item "..\configs\*" "configs\" -Recurse

# 回到项目根目录
Set-Location ..

Write-Host "构建完成!" -ForegroundColor Green