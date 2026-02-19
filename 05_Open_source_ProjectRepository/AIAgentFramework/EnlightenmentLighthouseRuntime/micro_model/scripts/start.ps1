#!/usr/bin/env powershell

# 启动微模型服务器

Write-Host "Starting micro model server..."

# 检查是否在正确的目录
if (-not (Test-Path "main.go")) {
    Write-Host "Error: main.go not found. Please run this script in the micro_model directory." -ForegroundColor Red
    exit 1
}

# 构建并运行服务器
Write-Host "Building server..."
go build -o micro_model_server.exe .

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful. Starting server..."
    .\micro_model_server.exe
} else {
    Write-Host "Build failed." -ForegroundColor Red
    exit 1
}
