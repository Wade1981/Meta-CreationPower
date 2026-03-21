# 对称数智财务系统启动脚本

# 设置工作目录
$workingDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$rootDir = Join-Path $workingDir ".."
Set-Location $rootDir

# 创建必要的目录
New-Item -Path "data\input" -ItemType Directory -Force
New-Item -Path "data\output" -ItemType Directory -Force
New-Item -Path "logs" -ItemType Directory -Force

# 检查Python是否安装
if (-not (Get-Command python -ErrorAction SilentlyContinue)) {
    Write-Host "Python未安装，请先安装Python 3.8+" -ForegroundColor Red
    exit 1
}

# 检查pip是否安装
if (-not (Get-Command pip -ErrorAction SilentlyContinue)) {
    Write-Host "pip未安装，请先安装pip" -ForegroundColor Red
    exit 1
}

# 安装依赖
Write-Host "正在安装依赖..." -ForegroundColor Green
pip install -r requirements.txt 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "依赖安装失败，请检查网络连接" -ForegroundColor Red
    exit 1
}

# 运行健康检查
Write-Host "正在运行健康检查..." -ForegroundColor Green
python src/main.py

# 提示用户如何使用系统
Write-Host "\n对称数智财务系统启动成功！" -ForegroundColor Green
Write-Host "使用方法: python src/main.py <输入数据文件路径>" -ForegroundColor Yellow
Write-Host "示例输入数据文件格式:"
Write-Host "{
  \"internal_entropy\": 0.1,
  \"external_entropy\": 0.1,
  \"financial_indicators\": {
    \"cash_flow\": 100000,
    \"revenue\": 500000,
    \"expenses\": 400000,
    \"assets\": 1000000,
    \"liabilities\": 500000
  },
  \"dimension_data\": {
    \"宏观调控\": [0.1, 0.2, 0.3],
    \"市场动态\": [0.2, 0.3, 0.4],
    \"决策模型\": [0.3, 0.4, 0.5],
    \"风控模型\": [0.4, 0.5, 0.6],
    \"执行管理\": [0.5, 0.6, 0.7]
  }
}" -ForegroundColor Cyan
