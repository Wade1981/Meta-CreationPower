#!/usr/bin/env powershell

# 启动所有ELR服务，包括Golang运行时和Python API服务器

Write-Host "Starting ELR services..." -ForegroundColor Green

# 切换到ELR根目录
cd "$(Split-Path -Parent $MyInvocation.MyCommand.Path)"

# 启动Python API服务器
Write-Host "Starting Python API server..." -ForegroundColor Cyan
Start-Process python -ArgumentList "elr_api_server.py" -WindowStyle Hidden -WorkingDirectory "."

# 等待API服务器启动
Start-Sleep -Seconds 2

# 切换到micro_model目录
cd micro_model

# 构建并启动Golang服务
Write-Host "Building and starting Golang server..." -ForegroundColor Cyan
go build -o micro_model_server.exe .

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful. Starting server..." -ForegroundColor Green
    .\micro_model_server.exe
} else {
    Write-Host "Build failed." -ForegroundColor Red
    exit 1
}
