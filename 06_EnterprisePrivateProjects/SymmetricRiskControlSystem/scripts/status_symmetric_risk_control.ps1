# 对称风险管控系统状态检查脚本

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

# 检查系统状态
function Check-System-Status {
    Write-Title "$PROJECT_NAME v$VERSION 状态检查"
    
    # 检查安装状态
    $installStatus = "未安装"
    $installDetails = ""
    
    if (Test-Path $INSTALL_DIR) {
        $installStatus = "已安装"
        $installDetails = "安装目录: $INSTALL_DIR"
        
        # 检查配置文件
        $configFile = "$INSTALL_DIR\config\system_config.json"
        if (Test-Path $configFile) {
            $installDetails += "\n配置文件: 存在"
        } else {
            $installDetails += "\n配置文件: 缺失"
        }
        
        # 检查启动脚本
        $startScript = "$INSTALL_DIR\start_symmetric_risk_control.ps1"
        if (Test-Path $startScript) {
            $installDetails += "\n启动脚本: 存在"
        } else {
            $installDetails += "\n启动脚本: 缺失"
        }
    }
    
    Write-Host "[INSTALL] 状态: $installStatus" -ForegroundColor Magenta
    if ($installDetails) {
        Write-Host "          $installDetails" -ForegroundColor Gray
    }
    
    # 检查ELR状态
    $elrStatus = "未安装"
    $elrDetails = ""
    
    if (Test-Path $ELR_DIR) {
        $elrExe = Get-ELR-Exe
        $elrStatus = "已安装"
        $elrDetails = "ELR目录: $ELR_DIR"
        
        # 切换到ELR目录
        Push-Location $ELR_DIR
        
        try {
            # 检查ELR运行状态
            $status = & $elrExe status 2>&1
            
            if ($status -like "*Error: ELR runtime is not running*" -or $status -like "*未运行*") {
                $elrStatus = "已安装 (未运行)"
                $elrDetails += "\n运行状态: 未运行"
            } else {
                $elrStatus = "已安装 (运行中)"
                $elrDetails += "\n运行状态: 运行中"
                
                # 检查容器状态
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
                
                if ($containerId) {
                    $elrDetails += "\n容器名称: $CONTAINER_NAME"
                    $elrDetails += "\n容器ID: $containerId"
                    $elrDetails += "\n容器状态: $containerStatus"
                    
                    # 检查服务状态
                    if ($containerStatus -eq "running" -or $containerStatus -eq "RUNNING") {
                        $serviceStatus = & $elrExe exec --id $containerId --command "python -m symmetric_risk_control status" 2>&1
                        
                        if ($serviceStatus -like "*Error:*" -and $serviceStatus -like "*服务未运行*" -or $serviceStatus -like "*not running*") {
                            $elrDetails += "\n服务状态: 未运行"
                        } elseif ($serviceStatus -like "*Error:*") {
                            $elrDetails += "\n服务状态: 未知 (错误: $serviceStatus)"
                        } else {
                            $elrDetails += "\n服务状态: 运行中"
                            # 显示服务状态详情
                            Write-Host ""
                            Write-Host "[SERVICE] 服务状态详情:" -ForegroundColor Blue
                            foreach ($line in $serviceStatus) {
                                if ($line -notlike "*=======================================*" -and $line -notlike "*Running Python...*" -and $line -notlike "*Executing:*") {
                                    Write-Host "          $line" -ForegroundColor Gray
                                }
                            }
                        }
                    } else {
                        $elrDetails += "\n服务状态: 容器未运行"
                    }
                } else {
                    $elrDetails += "\n容器状态: 容器不存在"
                    $elrDetails += "\n服务状态: 未安装"
                }
            }
        } catch {
            $elrDetails += "\n运行状态: 未知 (错误: $($_.Exception.Message))"
        } finally {
            Pop-Location
        }
    }
    
    Write-Host "[ELR] 状态: $elrStatus" -ForegroundColor Blue
    if ($elrDetails) {
        Write-Host "          $elrDetails" -ForegroundColor Gray
    }
    
    # 检查系统环境
    Write-Host ""
    Write-Host "[SYSTEM] 环境信息:" -ForegroundColor Green
    Write-Host "          操作系统: $([System.Environment]::OSVersion.VersionString)"
    Write-Host "          PowerShell: $($PSVersionTable.PSVersion.ToString())"
    Write-Host "          .NET版本: $([System.Environment]::Version.ToString())"
    
    # 检查Python
    try {
        $pythonVersion = python --version 2>&1
        Write-Host "          Python: $pythonVersion"
    } catch {
        Write-Host "          Python: 未安装"
    }
    
    # 显示使用方法
    Write-Host ""
    Write-Host "[USAGE] 使用方法:" -ForegroundColor Cyan
    Write-Host "          1. 启动服务: 运行 start_symmetric_risk_control.ps1"
    Write-Host "          2. 停止服务: 运行 stop_symmetric_risk_control.ps1"
    Write-Host "          3. 对话式交互: elr exec --id <container-id> --command \"python -m symmetric_risk_control <command>\""
    Write-Host ""
    Write-Host "          可用的对话式命令:"
    Write-Host "          - status: 查看系统状态"
    Write-Host "          - mechanism: 查看运行机制"
    Write-Host "          - report: 获取风险分析报告"
    Write-Host "          - config: 查看/调整配置"
    Write-Host ""
    
    # 显示健康状态
    Write-Host ""
    Write-Host "[HEALTH] 健康状态:" -ForegroundColor Yellow
    
    if ($installStatus -eq "未安装") {
        Write-Host "          状态: 需要安装"
        Write-Host "          建议: 运行 install_symmetric_risk_control.ps1 进行安装"
    } elseif ($elrStatus -like "*未安装*") {
        Write-Host "          状态: ELR未安装"
        Write-Host "          建议: 先安装Enlightenment Lighthouse Runtime"
    } elseif ($elrStatus -like "*未运行*") {
        Write-Host "          状态: ELR未运行"
        Write-Host "          建议: 启动ELR运行时后再启动服务"
    } elseif ($elrDetails -like "*容器不存在*") {
        Write-Host "          状态: 容器不存在"
        Write-Host "          建议: 重新运行安装脚本创建容器"
    } elseif ($elrDetails -like "*服务状态: 未运行*") {
        Write-Host "          状态: 服务未运行"
        Write-Host "          建议: 运行 start_symmetric_risk_control.ps1 启动服务"
    } else {
        Write-Host "          状态: 正常运行"
        Write-Host "          建议: 系统运行正常，无需操作"
    }
    
    Write-Host ""
}

# 主函数
function Main {
    try {
        Check-System-Status
    } catch {
        Handle-Error "检查状态时出现未预期的错误" $_.Exception.Message
    }
}

# 开始检查
Main
