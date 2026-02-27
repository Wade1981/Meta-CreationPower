# 启动后端服务器脚本

# 设置工作目录
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location -Path $scriptDir

# 检查Python是否安装
if (-not (Get-Command python -ErrorAction SilentlyContinue)) {
    Write-Host "错误: 未找到Python。请确保Python已安装并添加到系统路径中。" -ForegroundColor Red
    pause
    exit 1
}

# 检查Flask是否安装
try {
    python -c "import flask"
} catch {
    Write-Host "安装Flask..." -ForegroundColor Yellow
    python -m pip install flask
}

# 启动后端服务器
Write-Host "启动后端服务器..." -ForegroundColor Green
Write-Host "服务器将运行在 http://localhost:5000" -ForegroundColor Cyan
Write-Host "按 Ctrl+C 停止服务器" -ForegroundColor Yellow

python backend.py