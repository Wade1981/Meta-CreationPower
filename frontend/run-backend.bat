@echo off

REM 设置工作目录
cd /d %~dp0

REM 检查Python是否安装
python --version >nul 2>&1
if %errorlevel% neq 0 (
    echo 错误: 未找到Python。请确保Python已安装并添加到系统路径中。
    pause
    exit /b 1
)

REM 检查Flask是否安装
python -c "import flask" >nul 2>&1
if %errorlevel% neq 0 (
    echo 安装Flask...
    python -m pip install flask
)

REM 启动后端服务器
echo 启动后端服务器...
echo 服务器将运行在 http://localhost:5000
echo 按 Ctrl+C 停止服务器

python backend.py