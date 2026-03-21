# 对称风险管控系统启动脚本

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

# 启动系统
function Start-System {
    Write-Title "启动 $PROJECT_NAME v$VERSION"
    
    # 检查安装目录
    if (-not (Test-Path $INSTALL_DIR)) {
        Handle-Error "系统未安装，请先运行安装脚本。"
    }
    
    # 检查配置文件
    $configFile = "$INSTALL_DIR\config\system_config.json"
    if (-not (Test-Path $configFile)) {
        Handle-Error "配置文件不存在，请重新安装系统。"
    }
    
    # 获取ELR可执行文件
    $elrExe = Get-ELR-Exe
    
    # 切换到ELR目录
    Push-Location $ELR_DIR
    
    try {
        # 检查ELR运行状态
        Write-Info "检查ELR运行状态..."
        $status = & $elrExe status 2>&1
        
        # 如果ELR未运行，启动它
        if ($status -like "*Error: ELR runtime is not running*" -or $status -like "*未运行*") {
            Write-Info "启动ELR运行时..."
            & $elrExe start
            Start-Sleep -Seconds 2
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
            Handle-Error "容器 $CONTAINER_NAME 不存在，请重新安装系统。"
        }
        
        # 如果容器未运行，启动它
        if ($containerStatus -ne "running" -and $containerStatus -ne "RUNNING") {
            Write-Info "启动容器 $CONTAINER_NAME..."
            $startResult = & $elrExe start-container --id $containerId 2>&1
            
            if ($startResult -like "*Error:*") {
                Handle-Error "启动容器失败" $startResult
            }
        } else {
            Write-Info "容器 $CONTAINER_NAME 已在运行中。"
        }
        
        # 启动对称风险管控系统服务
        Write-Info "启动对称风险管控系统服务..."
        
        # 在容器中启动系统服务
        $startServiceResult = & $elrExe exec --id $containerId --command "python -m symmetric_risk_control start" 2>&1
        
        if ($startServiceResult -like "*Error:*" -and $startServiceResult -notlike "*服务已在运行*" -and $startServiceResult -notlike "*already running*") {
            Write-Host "[WARNING] 启动服务时出现警告，但系统将继续: $startServiceResult" -ForegroundColor Yellow
        }
        
        # 等待服务启动
        Start-Sleep -Seconds 3
        
        # 验证服务状态
        Write-Info "验证服务状态..."
        $statusResult = & $elrExe exec --id $containerId --command "python -m symmetric_risk_control status" 2>&1
        
        if ($statusResult -like "*Error:*" -and $statusResult -notlike "*服务未运行*" -and $statusResult -notlike "*not running*") {
            Write-Host "[WARNING] 验证服务状态时出现警告: $statusResult" -ForegroundColor Yellow
        }
        
        Handle-Success "$PROJECT_NAME 启动成功！"
        
        # 显示服务状态
        Write-Host ""
        Write-Host "服务状态:"
        Write-Host "- 容器名称: $CONTAINER_NAME"
        Write-Host "- 容器ID: $containerId"
        Write-Host "- 安装目录: $INSTALL_DIR"
        Write-Host ""
        Write-Host "使用方法:"
        Write-Host "  - 停止服务: 运行 stop_symmetric_risk_control.ps1"
        Write-Host "  - 查看状态: 运行 status_symmetric_risk_control.ps1"
        Write-Host "  - 对话式交互: elr exec --id $containerId --command \"python -m symmetric_risk_control <command>\""
        Write-Host ""
        
    } catch {
        Handle-Error "启动系统时出现错误" $_.Exception.Message
    } finally {
        Pop-Location
    }
}

# 主函数
function Main {
    try {
        Start-System
    } catch {
        Handle-Error "启动过程中出现未预期的错误" $_.Exception.Message
    }
}

# 开始启动
Main
