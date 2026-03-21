# 在ELR容器中运行RootPulseOS的脚本

Write-Host "=====================================" -ForegroundColor Green
Write-Host "RootPulseOS ELR Container Runner" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green

# 检查ELR是否存在
$elrPath = "E:\X54\github\Meta-CreationPower\05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\elr.ps1"
if (-not (Test-Path $elrPath)) {
    Write-Host "ELR脚本不存在，请检查路径是否正确。" -ForegroundColor Red
    Write-Host "预期路径: $elrPath" -ForegroundColor Yellow
    exit 1
}

Write-Host "ELR脚本路径: $elrPath" -ForegroundColor Cyan

# 导航到项目根目录
Set-Location "$(Split-Path -Parent $MyInvocation.MyCommand.Path)\.."
$projectRoot = Get-Location
Write-Host "项目根目录: $projectRoot" -ForegroundColor Cyan

# 启动ELR运行时
Write-Host "启动ELR运行时..." -ForegroundColor Yellow
try {
    & powershell -ExecutionPolicy RemoteSigned -File $elrPath start
    Write-Host "ELR运行时启动成功!" -ForegroundColor Green
} catch {
    Write-Host "ELR运行时启动失败: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 等待ELR运行时初始化
Write-Host "等待ELR运行时初始化..." -ForegroundColor Yellow
Start-Sleep -Seconds 3

# 检查ELR运行时状态
Write-Host "检查ELR运行时状态..." -ForegroundColor Yellow
try {
    $status = & powershell -ExecutionPolicy RemoteSigned -File $elrPath status
    Write-Host "ELR运行时状态: $status" -ForegroundColor Cyan
} catch {
    Write-Host "检查ELR运行时状态失败: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 检查容器是否已存在
Write-Host "检查RootPulseOS容器是否已存在..." -ForegroundColor Yellow
try {
    $containers = & powershell -ExecutionPolicy RemoteSigned -File $elrPath list
    if ($containers -like "*rootpulseos-container*") {
        Write-Host "RootPulseOS容器已存在，正在停止并删除..." -ForegroundColor Yellow
        & powershell -ExecutionPolicy RemoteSigned -File $elrPath stop-container --id rootpulseos-container
        & powershell -ExecutionPolicy RemoteSigned -File $elrPath delete --id rootpulseos-container
        Start-Sleep -Seconds 2
    }
} catch {
    Write-Host "检查容器状态失败: $($_.Exception.Message)" -ForegroundColor Yellow
    # 继续执行，可能容器不存在
}

# 创建并运行RootPulseOS容器
Write-Host "创建并运行RootPulseOS容器..." -ForegroundColor Yellow
try {
    # 使用ELR的run命令运行Python应用
    & powershell -ExecutionPolicy RemoteSigned -File $elrPath run `
        --name rootpulseos-container `
        --language python `
        --command "python main.py" `
        --port 8000:8000 `
        --env ROOTPULSEOS_ENV=production `
        --env ROOTPULSEOS_LOG_LEVEL=INFO `
        --env ELR_TEST_API_URL=http://lumina-runtime:8080
    
    Write-Host "RootPulseOS容器创建并运行成功!" -ForegroundColor Green
} catch {
    Write-Host "创建容器失败: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 等待容器启动
Write-Host "等待RootPulseOS容器启动..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# 检查容器状态
Write-Host "检查RootPulseOS容器状态..." -ForegroundColor Yellow
try {
    $containerStatus = & powershell -ExecutionPolicy RemoteSigned -File $elrPath inspect --id rootpulseos-container
    Write-Host "容器状态: $containerStatus" -ForegroundColor Cyan
} catch {
    Write-Host "检查容器状态失败: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 显示访问信息
Write-Host "=====================================" -ForegroundColor Green
Write-Host "RootPulseOS容器已成功运行!" -ForegroundColor Green
Write-Host "访问地址: http://localhost:8000" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Green

# 提示如何停止容器
Write-Host "要停止RootPulseOS容器，请运行:" -ForegroundColor Yellow
Write-Host "powershell -ExecutionPolicy RemoteSigned -File $elrPath stop-container --id rootpulseos-container" -ForegroundColor Cyan

Write-Host "" -ForegroundColor White
Write-Host "按任意键退出..." -ForegroundColor White
$null = $Host.UI.RawUI.ReadKey('NoEcho,IncludeKeyDown')
