# 对称风险管控系统停止脚本

# 版本信息
$VERSION = "1.0.0"
$PROJECT_NAME = "Symmetric Risk Control System"

# 安装目录
$INSTALL_DIR = $env:SYMMETRIC_RISK_CONTROL_HOME
if (-not $INSTALL_DIR) {
    $INSTALL_DIR = "$env:ProgramFiles\EnlightenmentLighthouse\SymmetricRiskControlSystem"
}

# ELR目录
$ELR_DIR = "$env:ProgramFiles\EnlightenmentLighthouse\Runtime"

# 容器名称
$CONTAINER_NAME = "symmetric-risk-control"

# 错误处理函数
function Handle-Error {
    param(
        [string]$Message,
        [string]$ErrorInfo = $null
    )
    Write-Host "[ERROR] $Message" -ForegroundColor Red
    if ($ErrorInfo) {
        Write-Host "[DETAIL] $ErrorInfo" -ForegroundColor Yellow
    }
    exit 1
}

# 成功处理函数
function Handle-Success {
    param(
        [string]$Message
    )
    Write-Host "[SUCCESS] $Message" -ForegroundColor Green
}

# 信息输出函数
function Write-Info {
    param(
        [string]$Message
    )
    Write-Host "[INFO] $Message" -ForegroundColor Cyan
}

# 标题输出
function Write-Title {
    param(
        [string]$Title
    )
    Write-Host ""
    Write-Host "========================================"
    Write-Host "$Title"
    Write-Host "========================================"
    Write-Host ""
}

# 检查ELR可执行文件
function Get-ELR-Exe {
    $elrExe = "$ELR_DIR\elr.ps1"
    if (-not (Test-Path $elrExe)) {
        $elrExe = "$ELR_DIR\elr.bat"
        if (-not (Test-Path $elrExe)) {
            Handle-Error "找不到ELR可执行文件，请确认ELR安装正确。"
        }
    }
    return $elrExe
}

# 停止系统
function Stop-System {
    Write-Title "停止 $PROJECT_NAME v$VERSION"
    
    # 检查安装目录
    if (-not (Test-Path $INSTALL_DIR)) {
        Write-Info "系统未安装，跳过停止步骤。"
        return
    }
    
    # 获取ELR可执行文件
    $elrExe = Get-ELR-Exe
    
    # 切换到ELR目录
    Push-Location $ELR_DIR
    
    try {
        # 检查ELR运行状态
        Write-Info "检查ELR运行状态..."
        $status = & $elrExe status 2>&1
        
        # 如果ELR未运行，直接返回
        if ($status -like "*Error: ELR runtime is not running*" -or $status -like "*未运行*") {
            Write-Info "ELR运行时未运行，跳过停止步骤。"
            Pop-Location
            return
        }
        
        # 检查容器状态
        Write-Info "检查容器状态..."
        $containers = & $elrExe list 2>&1
        
        $containerId = $null
        $containerStatus = $null
        
        foreach ($line in $containers) {
            if ($line -like "*$CONTAINER_NAME*" -and $line -notlike "*NAME*" -and $line -notlike "*--*" -and $line -notlike "*====*" -and $line -notlike "*Containers:*") {
                $parts = $line.Trim().Split()
                $containerId = $parts[0]
                $containerStatus = $parts[3]
                break
            }
        }
        
        if (-not $containerId) {
            Write-Info "容器 $CONTAINER_NAME 不存在，跳过停止步骤。"
            Pop-Location
            return
        }
        
        # 停止对称风险管控系统服务
        Write-Info "停止对称风险管控系统服务..."
        
        # 在容器中停止系统服务
        $stopServiceResult = & $elrExe exec --id $containerId --command "python -m symmetric_risk_control stop" 2>&1
        
        if ($stopServiceResult -like "*Error:*" -and $stopServiceResult -notlike "*服务未运行*" -and $stopServiceResult -notlike "*not running*") {
            Write-Host "[WARNING] 停止服务时出现警告，但系统将继续: $stopServiceResult" -ForegroundColor Yellow
        }
        
        # 等待服务停止
        Start-Sleep -Seconds 2
        
        # 停止容器
        if ($containerStatus -eq "running" -or $containerStatus -eq "RUNNING") {
            Write-Info "停止容器 $CONTAINER_NAME..."
            $stopResult = & $elrExe stop-container --id $containerId 2>&1
            
            if ($stopResult -like "*Error:*" -and $stopResult -notlike "*容器未运行*" -and $stopResult -notlike "*not running*") {
                Write-Host "[WARNING] 停止容器时出现警告，但系统将继续: $stopResult" -ForegroundColor Yellow
            }
        } else {
            Write-Info "容器 $CONTAINER_NAME 未在运行中，跳过停止步骤。"
        }
        
        Handle-Success "$PROJECT_NAME 停止成功！"
        
        # 显示停止状态
        Write-Host ""
        Write-Host "停止状态:"
        Write-Host "- 容器名称: $CONTAINER_NAME"
        Write-Host "- 容器ID: $containerId"
        Write-Host "- 服务状态: 已停止"
        Write-Host ""
        Write-Host "使用方法:"
        Write-Host "  - 启动服务: 运行 start_symmetric_risk_control.ps1"
        Write-Host "  - 查看状态: 运行 status_symmetric_risk_control.ps1"
        Write-Host ""
        
    } catch {
        Handle-Error "停止系统时出现错误" $_.Exception.Message
    } finally {
        Pop-Location
    }
}

# 主函数
function Main {
    try {
        Stop-System
    } catch {
        Handle-Error "停止过程中出现未预期的错误" $_.Exception.Message
    }
}

# 开始停止
Main
