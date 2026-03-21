#!/usr/bin/env pwsh

# 启动测试脚本

Write-Host "开始启动 ELR 智能测试系统..."

# 检查构建目录是否存在
if (-not (Test-Path "build")) {
    Write-Host "构建目录不存在，请先运行 build_project.ps1 构建项目!" -ForegroundColor Red
    exit 1
}

# 进入构建目录
Set-Location "build"

# 检查可执行文件是否存在
if (-not (Test-Path "Release\elr_test_engine.exe")) {
    Write-Host "可执行文件不存在，请先运行 build_project.ps1 构建项目!" -ForegroundColor Red
    Set-Location ..
    exit 1
}

# 启动 ELR 容器（如果需要）
Write-Host "启动 ELR 容器..."
docker-compose -f "..\ELR-Containers\docker-compose.yml" up -d

if ($LASTEXITCODE -ne 0) {
    Write-Host "启动 ELR 容器失败!" -ForegroundColor Red
    Set-Location ..
    exit 1
}

# 等待 ELR 服务启动
Write-Host "等待 ELR 服务启动..."
Start-Sleep -Seconds 10

# 运行测试
Write-Host "运行测试..."
./Release/elr_test_engine.exe

if ($LASTEXITCODE -ne 0) {
    Write-Host "测试执行失败!" -ForegroundColor Red
    # 停止 ELR 容器
    docker-compose -f "..\ELR-Containers\docker-compose.yml" down
    Set-Location ..
    exit 1
}

# 停止 ELR 容器
Write-Host "停止 ELR 容器..."
docker-compose -f "..\ELR-Containers\docker-compose.yml" down

if ($LASTEXITCODE -ne 0) {
    Write-Host "停止 ELR 容器失败!" -ForegroundColor Red
    Set-Location ..
    exit 1
}

# 回到项目根目录
Set-Location ..

Write-Host "测试完成!" -ForegroundColor Green