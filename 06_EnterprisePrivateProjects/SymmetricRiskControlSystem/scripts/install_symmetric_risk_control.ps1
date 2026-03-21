# 对称风险管控系统安装脚本
# 自动装配到ELR容器中

# 版本信息
$VERSION = "1.0.0"
$PROJECT_NAME = "Symmetric Risk Control System"
$CONTAINER_NAME = "symmetric-risk-control"
$CONTAINER_IMAGE = "python:3.9"

# 安装目录
$INSTALL_DIR = "$env:ProgramFiles\EnlightenmentLighthouse\SymmetricRiskControlSystem"
$ELR_DIR = "$env:ProgramFiles\EnlightenmentLighthouse\Runtime"

# 安装状态
$SUCCESS = $false

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
    Write-Host "安装失败，请检查上述错误信息。" -ForegroundColor Red
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

# 检查ELR容器是否存在
function Check-ELR {
    Write-Info "检查ELR容器是否存在..."
    
    # 检查ELR目录
    if (-not (Test-Path $ELR_DIR)) {
        Handle-Error "ELR容器未安装，请先安装Enlightenment Lighthouse Runtime。"
    }
    
    # 检查ELR可执行文件
    $elrExe = "$ELR_DIR\elr.ps1"
    if (-not (Test-Path $elrExe)) {
        $elrExe = "$ELR_DIR\elr.bat"
        if (-not (Test-Path $elrExe)) {
            Handle-Error "找不到ELR可执行文件，请确认ELR安装正确。"
        }
    }
    
    return $elrExe
}

# 检查系统需求
function Check-System-Requirements {
    Write-Info "检查系统需求..."
    
    # 检查PowerShell版本
    $psVersion = $PSVersionTable.PSVersion
    if ($psVersion.Major -lt 5) {
        Handle-Error "PowerShell版本过低，需要PowerShell 5.0或更高版本。"
    }
    
    # 检查.NET Framework版本
    try {
        $netVersion = (Get-ItemProperty "HKLM:\Software\Microsoft\NET Framework Setup\NDP\v4\Full" -ErrorAction Stop).Version
        if ([version]$netVersion -lt [version]"4.5") {
            Handle-Error ".NET Framework版本过低，需要.NET Framework 4.5或更高版本。"
        }
    } catch {
        Handle-Error "无法检测.NET Framework版本，请确保已安装.NET Framework 4.5或更高版本。"
    }
    
    # 检查系统架构
    if ([System.Environment]::Is64BitOperatingSystem) {
        Write-Info "检测到64位操作系统，符合要求。"
    } else {
        Handle-Error "本系统仅支持64位操作系统。"
    }
    
    Handle-Success "系统需求检查通过。"
}

# 安装对称风险管控系统
function Install-System {
    param(
        [string]$ElrExe
    )
    
    Write-Title "开始安装 $PROJECT_NAME v$VERSION"
    
    # 创建安装目录
    Write-Info "创建安装目录..."
    try {
        New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
        New-Item -ItemType Directory -Path "$INSTALL_DIR\bin" -Force | Out-Null
        New-Item -ItemType Directory -Path "$INSTALL_DIR\config" -Force | Out-Null
        New-Item -ItemType Directory -Path "$INSTALL_DIR\logs" -Force | Out-Null
        New-Item -ItemType Directory -Path "$INSTALL_DIR\data" -Force | Out-Null
    } catch {
        Handle-Error "无法创建安装目录" $_.Exception.Message
    }
    
    # 复制系统文件
    Write-Info "复制系统文件..."
    try {
        # 复制核心模块
        Copy-Item -Path "$PSScriptRoot\..\src\core" -Destination "$INSTALL_DIR\bin" -Recurse -Force -ErrorAction Stop
        
        # 复制对话式交互模块
        Copy-Item -Path "$PSScriptRoot\..\src\dialogue" -Destination "$INSTALL_DIR\bin" -Recurse -Force -ErrorAction Stop
        
        # 复制ELR集成模块
        Copy-Item -Path "$PSScriptRoot\..\src\elr_integration" -Destination "$INSTALL_DIR\bin" -Recurse -Force -ErrorAction Stop
        
        # 复制配置文件
        Copy-Item -Path "$PSScriptRoot\..\config" -Destination "$INSTALL_DIR" -Recurse -Force -ErrorAction Stop
        
        # 复制启动脚本
        Copy-Item -Path "$PSScriptRoot\start_symmetric_risk_control.ps1" -Destination "$INSTALL_DIR" -Force -ErrorAction Stop
        Copy-Item -Path "$PSScriptRoot\stop_symmetric_risk_control.ps1" -Destination "$INSTALL_DIR" -Force -ErrorAction Stop
        Copy-Item -Path "$PSScriptRoot\status_symmetric_risk_control.ps1" -Destination "$INSTALL_DIR" -Force -ErrorAction Stop
    } catch {
        Handle-Error "无法复制系统文件" $_.Exception.Message
    }
    
    # 编译Python代码为.pyc文件（保护源代码）
    Write-Info "编译Python代码为.pyc文件..."
    try {
        $pythonFiles = Get-ChildItem -Path "$INSTALL_DIR\bin" -Recurse -Filter "*.py"
        foreach ($file in $pythonFiles) {
            $pycFile = $file.FullName.Replace(".py", ".pyc")
            Write-Host "  编译: $($file.Name)"
            # 使用Python编译
            python -m py_compile $file.FullName
            if ($LASTEXITCODE -eq 0) {
                # 删除.py文件，只保留.pyc文件
                Remove-Item $file.FullName -Force
            } else {
                Write-Host "  编译失败: $($file.Name)" -ForegroundColor Yellow
            }
        }
    } catch {
        Write-Host "[WARNING] 编译Python代码时出现错误，但安装将继续: $($_.Exception.Message)" -ForegroundColor Yellow
    }
    
    Handle-Success "系统文件安装完成。"
}

# 创建ELR容器
function Create-ELR-Container {
    param(
        [string]$ElrExe
    )
    
    Write-Info "创建ELR容器..."
    
    # 切换到ELR目录
    Push-Location $ELR_DIR
    
    try {
        # 检查ELR运行状态
        Write-Host "  检查ELR运行状态..."
        $status = & $ElrExe status 2>&1
        
        # 如果ELR未运行，启动它
        if ($status -like "*Error: ELR runtime is not running*" -or $status -like "*未运行*") {
            Write-Host "  启动ELR运行时..."
            & $ElrExe start
            Start-Sleep -Seconds 2
        }
        
        # 检查是否已存在同名容器
        Write-Host "  检查容器是否已存在..."
        $containers = & $ElrExe list 2>&1
        
        if ($containers -like "*$CONTAINER_NAME*" -and $containers -notlike "*Error:*") {
            Write-Info "容器 $CONTAINER_NAME 已存在，跳过创建步骤。"
        } else {
            # 创建新容器
            Write-Host "  创建容器 $CONTAINER_NAME..."
            $createResult = & $ElrExe create --name $CONTAINER_NAME --image $CONTAINER_IMAGE 2>&1
            
            if ($createResult -like "*Error:*") {
                Handle-Error "创建ELR容器失败" $createResult
            }
            
            # 启动容器
            Write-Host "  启动容器 $CONTAINER_NAME..."
            # 获取容器ID
            $containers = & $ElrExe list 2>&1
            $containerId = $null
            foreach ($line in $containers) {
                if ($line -like "*$CONTAINER_NAME*" -and $line -notlike "*NAME*" -and $line -notlike "*--*" -and $line -notlike "*====*" -and $line -notlike "*Containers:*") {
                    $containerId = $line.Trim().Split()[0]
                    break
                }
            }
            
            if ($containerId) {
                $startResult = & $ElrExe start-container --id $containerId 2>&1
                if ($startResult -like "*Error:*") {
                    Handle-Error "启动ELR容器失败" $startResult
                }
            } else {
                Handle-Error "无法获取容器ID"
            }
        }
        
    } catch {
        Handle-Error "创建ELR容器时出现错误" $_.Exception.Message
    } finally {
        Pop-Location
    }
    
    Handle-Success "ELR容器创建完成。"
}

# 配置系统
function Configure-System {
    param(
        [string]$ElrExe
    )
    
    Write-Info "配置系统..."
    
    try {
        # 生成配置文件
        $configContent = @"
{
    "version": "$VERSION",
    "container_name": "$CONTAINER_NAME",
    "elr_path": "$ELR_DIR",
    "install_dir": "$INSTALL_DIR",
    "log_level": "info",
    "risk_assessment": {
        "enabled": true,
        "threshold": 0.7,
        "update_interval": 60
    },
    "risk_control": {
        "enabled": true,
        "max_exposure": 0.5,
        "hedging_ratio": 0.8
    },
    "data_sources": {
        "enabled": [],
        "configs": {}
    },
    "dialogue": {
        "enabled": true,
        "language": "zh-CN"
    }
}
"@
        
        $configContent | Set-Content -Path "$INSTALL_DIR\config\system_config.json" -Encoding UTF8 -Force
        
        # 创建环境变量
        [Environment]::SetEnvironmentVariable("SYMMETRIC_RISK_CONTROL_HOME", $INSTALL_DIR, "Machine")
        [Environment]::SetEnvironmentVariable("PATH", "$env:PATH;$INSTALL_DIR", "Machine")
        
    } catch {
        Handle-Error "配置系统时出现错误" $_.Exception.Message
    }
    
    Handle-Success "系统配置完成。"
}

# 验证安装
function Verify-Installation {
    param(
        [string]$ElrExe
    )
    
    Write-Info "验证安装..."
    
    try {
        # 检查安装目录
        if (-not (Test-Path $INSTALL_DIR)) {
            Handle-Error "安装目录不存在"
        }
        
        # 检查核心文件
        $requiredFiles = @(
            "$INSTALL_DIR\config\system_config.json",
            "$INSTALL_DIR\start_symmetric_risk_control.ps1",
            "$INSTALL_DIR\stop_symmetric_risk_control.ps1",
            "$INSTALL_DIR\status_symmetric_risk_control.ps1"
        )
        
        foreach ($file in $requiredFiles) {
            if (-not (Test-Path $file)) {
                Handle-Error "缺少必要文件: $file"
            }
        }
        
        # 检查ELR容器状态
        Push-Location $ELR_DIR
        try {
            $containers = & $ElrExe list 2>&1
            if ($containers -notlike "*$CONTAINER_NAME*" -or $containers -like "*Error:*") {
                Handle-Error "ELR容器创建失败"
            }
        } finally {
            Pop-Location
        }
        
    } catch {
        Handle-Error "验证安装时出现错误" $_.Exception.Message
    }
    
    Handle-Success "安装验证完成。"
}

# 主安装函数
function Main-Install {
    try {
        # 显示安装信息
        Write-Title "$PROJECT_NAME 安装向导"
        Write-Host "版本: $VERSION"
        Write-Host "目标: 自动装配到ELR容器中"
        Write-Host ""
        
        # 检查系统需求
        Check-System-Requirements
        
        # 检查ELR
        $elrExe = Check-ELR
        
        # 安装系统
        Install-System
        
        # 创建ELR容器
        Create-ELR-Container -ElrExe $elrExe
        
        # 配置系统
        Configure-System -ElrExe $elrExe
        
        # 验证安装
        Verify-Installation -ElrExe $elrExe
        
        # 安装成功
        Write-Title "安装完成"
        Write-Host "$PROJECT_NAME v$VERSION 已成功安装！"
        Write-Host ""
        Write-Host "安装位置: $INSTALL_DIR"
        Write-Host "ELR容器: $CONTAINER_NAME"
        Write-Host ""
        Write-Host "使用方法:"
        Write-Host "  1. 启动系统: 运行 start_symmetric_risk_control.ps1"
        Write-Host "  2. 停止系统: 运行 stop_symmetric_risk_control.ps1"
        Write-Host "  3. 查看状态: 运行 status_symmetric_risk_control.ps1"
        Write-Host "  4. 对话式交互: elr exec --id <container-id> --command "python -m symmetric_risk_control <command>""
        Write-Host ""
        Write-Host "如有问题，请联系技术支持。"
        Write-Host ""
        
        $SUCCESS = $true
        
    } catch {
        Handle-Error "安装过程中出现未预期的错误" $_.Exception.Message
    } finally {
        if (-not $SUCCESS) {
            Write-Host "安装失败，请检查上述错误信息。" -ForegroundColor Red
            exit 1
        }
    }
}

# 开始安装
Main-Install
