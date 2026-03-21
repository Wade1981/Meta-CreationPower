# RootPulseOS启动脚本

Write-Host "Starting RootPulseOS..." -ForegroundColor Green

# 导航到项目根目录
Set-Location "$(Split-Path -Parent $MyInvocation.MyCommand.Path)\.."

# 检查Python是否安装
if (-not (Get-Command python -ErrorAction SilentlyContinue)) {
    Write-Host "Python is not installed. Please install Python 3.8+ first." -ForegroundColor Red
    exit 1
}

# 检查Python版本
$pythonVersion = python --version
Write-Host "Python version: $pythonVersion" -ForegroundColor Cyan

# 检查依赖
Write-Host "Checking dependencies..." -ForegroundColor Yellow
try {
    python -c "import requests"
    Write-Host "requests library is installed" -ForegroundColor Green
} catch {
    Write-Host "Installing requests library..." -ForegroundColor Yellow
    pip install requests
}

# 启动RootPulseOS
Write-Host "Starting RootPulseOS system..." -ForegroundColor Green
try {
    python main.py
} catch {
    Write-Host "Failed to start RootPulseOS: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
