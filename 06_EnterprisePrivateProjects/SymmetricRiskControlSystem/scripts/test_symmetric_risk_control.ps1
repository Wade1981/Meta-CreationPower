# 对称风险管控系统测试脚本
# 用于验证系统的ELR集成和功能

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

# 测试结果
$testResults = @()
$totalTests = 0
$passedTests = 0
$failedTests = 0

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

# 测试结果记录函数
function Record-TestResult {
    param(
        [string]$TestName,
        [bool]$Passed,
        [string]$Message = "",
        [string]$ErrorInfo = $null
    )
    
    $global:totalTests++
    if ($Passed) {
        $global:passedTests++
        Write-Host "[PASS] $TestName" -ForegroundColor Green
    } else {
        $global:failedTests++
        Write-Host "[FAIL] $TestName" -ForegroundColor Red
        if ($Message) {
            Write-Host "       $Message" -ForegroundColor Yellow
        }
        if ($ErrorInfo) {
            Write-Host "       Error: $ErrorInfo" -ForegroundColor DarkYellow
        }
    }
    
    $result = @{
        TestName = $TestName
        Passed = $Passed
        Message = $Message
        ErrorInfo = $ErrorInfo
        Timestamp = Get-Date
    }
    
    $global:testResults += $result
}

# 检查ELR可执行文件
function Get-ELR-Exe {
    $elrExe = "$ELR_DIR\elr.ps1"
    if (-not (Test-Path $elrExe)) {
        $elrExe = "$ELR_DIR\elr.bat"
        if (-not (Test-Path $elrExe)) {
            return $null
        }
    }
    return $elrExe
}

# 测试安装状态
function Test-Installation {
    Write-Title "测试 1: 安装状态检查"
    
    # 检查安装目录
    if (Test-Path $INSTALL_DIR) {
        Record-TestResult "安装目录存在" $true "安装目录: $INSTALL_DIR"
    } else {
        Record-TestResult "安装目录存在" $false "安装目录不存在: $INSTALL_DIR"
        return $false
    }
    
    # 检查配置文件
    $configFile = "$INSTALL_DIR\config\system_config.json"
    if (Test-Path $configFile) {
        Record-TestResult "配置文件存在" $true "配置文件: $configFile"
    } else {
        Record-TestResult "配置文件存在" $false "配置文件不存在: $configFile"
    }
    
    # 检查启动脚本
    $startScript = "$INSTALL_DIR\start_symmetric_risk_control.ps1"
    if (Test-Path $startScript) {
        Record-TestResult "启动脚本存在" $true "启动脚本: $startScript"
    } else {
        Record-TestResult "启动脚本存在" $false "启动脚本不存在: $startScript"
    }
    
    # 检查停止脚本
    $stopScript = "$INSTALL_DIR\stop_symmetric_risk_control.ps1"
    if (Test-Path $stopScript) {
        Record-TestResult "停止脚本存在" $true "停止脚本: $stopScript"
    } else {
        Record-TestResult "停止脚本存在" $false "停止脚本不存在: $stopScript"
    }
    
    # 检查状态脚本
    $statusScript = "$INSTALL_DIR\status_symmetric_risk_control.ps1"
    if (Test-Path $statusScript) {
        Record-TestResult "状态脚本存在" $true "状态脚本: $statusScript"
    } else {
        Record-TestResult "状态脚本存在" $false "状态脚本不存在: $statusScript"
    }
    
    # 检查对话式交互模块
    $dialogueModule = "$INSTALL_DIR\bin\dialogue\__init__.py"
    if (Test-Path $dialogueModule) {
        Record-TestResult "对话式交互模块存在" $true "对话模块: $dialogueModule"
    } else {
        Record-TestResult "对话式交互模块存在" $false "对话模块不存在: $dialogueModule"
    }
    
    return $true
}

# 测试ELR集成
function Test-ELR-Integration {
    Write-Title "测试 2: ELR集成检查"
    
    # 检查ELR目录
    if (-not (Test-Path $ELR_DIR)) {
        Record-TestResult "ELR目录存在" $false "ELR目录不存在: $ELR_DIR"
        return $false
    }
    
    # 检查ELR可执行文件
    $elrExe = Get-ELR-Exe
    if ($elrExe) {
        Record-TestResult "ELR可执行文件存在" $true "ELR可执行文件: $elrExe"
    } else {
        Record-TestResult "ELR可执行文件存在" $false "ELR可执行文件不存在"
        return $false
    }
    
    # 切换到ELR目录
    Push-Location $ELR_DIR
    
    try {
        # 检查ELR运行状态
        Write-Host "  检查ELR运行状态..."
        $status = & $elrExe status 2>&1
        
        if ($status -like "*Error: ELR runtime is not running*" -or $status -like "*未运行*") {
            # 启动ELR运行时
            Write-Host "  启动ELR运行时..."
            & $elrExe start
            Start-Sleep -Seconds 2
            Record-TestResult "ELR运行时启动" $true "ELR运行时已启动"
        } else {
            Record-TestResult "ELR运行时状态" $true "ELR运行时已在运行"
        }
        
        # 检查容器状态
        Write-Host "  检查容器状态..."
        $containers = & $elrExe list 2>&1
        
        $containerExists = $false
        foreach ($line in $containers) {
            if ($line -like "*$CONTAINER_NAME*" -and $line -notlike "*NAME*" -and $line -notlike "*--*" -and $line -notlike "*====*" -and $line -notlike "*Containers:*") {
                $containerExists = $true
                break
            }
        }
        
        if ($containerExists) {
            Record-TestResult "ELR容器存在" $true "容器 $CONTAINER_NAME 已存在"
        } else {
            Record-TestResult "ELR容器存在" $false "容器 $CONTAINER_NAME 不存在"
        }
        
    } catch {
        Record-TestResult "ELR集成测试" $false "测试过程中出现错误" $_.Exception.Message
    } finally {
        Pop-Location
    }
    
    return $true
}

# 测试对话式交互
function Test-Dialogue-Interface {
    Write-Title "测试 3: 对话式交互检查"
    
    # 检查对话式交互模块
    $dialogueModule = "$INSTALL_DIR\bin\dialogue\__init__.py"
    if (-not (Test-Path $dialogueModule)) {
        Record-TestResult "对话式交互模块存在" $false "对话模块不存在: $dialogueModule"
        return $false
    }
    
    # 测试对话式交互命令
    try {
        # 测试status命令
        Write-Host "  测试status命令..."
        $statusResult = python -m dialogue status 2>&1
        if ($statusResult -like "*对称风险管控系统 - 状态信息*" -and $statusResult -notlike "*Error:*") {
            Record-TestResult "status命令测试" $true "status命令执行成功"
        } else {
            Record-TestResult "status命令测试" $false "status命令执行失败" $statusResult
        }
        
        # 测试mechanism命令
        Write-Host "  测试mechanism命令..."
        $mechanismResult = python -m dialogue mechanism 2>&1
        if ($mechanismResult -like "*对称风险管控系统 - 运行机制*" -and $mechanismResult -notlike "*Error:*") {
            Record-TestResult "mechanism命令测试" $true "mechanism命令执行成功"
        } else {
            Record-TestResult "mechanism命令测试" $false "mechanism命令执行失败" $mechanismResult
        }
        
        # 测试report命令
        Write-Host "  测试report命令..."
        $reportResult = python -m dialogue report 2>&1
        if ($reportResult -like "*风险分析报告*" -and $reportResult -notlike "*Error:*") {
            Record-TestResult "report命令测试" $true "report命令执行成功"
        } else {
            Record-TestResult "report命令测试" $false "report命令执行失败" $reportResult
        }
        
        # 测试config命令
        Write-Host "  测试config命令..."
        $configResult = python -m dialogue config 2>&1
        if ($configResult -like "*配置信息*" -and $configResult -notlike "*Error:*") {
            Record-TestResult "config命令测试" $true "config命令执行成功"
        } else {
            Record-TestResult "config命令测试" $false "config命令执行失败" $configResult
        }
        
    } catch {
        Record-TestResult "对话式交互测试" $false "测试过程中出现错误" $_.Exception.Message
    }
    
    return $true
}

# 测试系统功能
function Test-System-Functionality {
    Write-Title "测试 4: 系统功能检查"
    
    # 测试启动脚本
    try {
        Write-Host "  测试启动脚本..."
        $startScript = "$INSTALL_DIR\start_symmetric_risk_control.ps1"
        if (Test-Path $startScript) {
            # 运行启动脚本
            $startResult = & $startScript 2>&1
            if ($startResult -like "*启动成功*" -or $startResult -like "*已在运行中*" -or $startResult -notlike "*Error:*") {
                Record-TestResult "启动脚本测试" $true "启动脚本执行成功"
            } else {
                Record-TestResult "启动脚本测试" $false "启动脚本执行失败" $startResult
            }
        } else {
            Record-TestResult "启动脚本测试" $false "启动脚本不存在"
        }
        
    } catch {
        Record-TestResult "启动脚本测试" $false "测试过程中出现错误" $_.Exception.Message
    }
    
    # 测试状态脚本
    try {
        Write-Host "  测试状态脚本..."
        $statusScript = "$INSTALL_DIR\status_symmetric_risk_control.ps1"
        if (Test-Path $statusScript) {
            # 运行状态脚本
            $statusResult = & $statusScript 2>&1
            if ($statusResult -like "*状态检查*" -and $statusResult -notlike "*Error:*") {
                Record-TestResult "状态脚本测试" $true "状态脚本执行成功"
            } else {
                Record-TestResult "状态脚本测试" $false "状态脚本执行失败" $statusResult
            }
        } else {
            Record-TestResult "状态脚本测试" $false "状态脚本不存在"
        }
        
    } catch {
        Record-TestResult "状态脚本测试" $false "测试过程中出现错误" $_.Exception.Message
    }
    
    # 测试停止脚本
    try {
        Write-Host "  测试停止脚本..."
        $stopScript = "$INSTALL_DIR\stop_symmetric_risk_control.ps1"
        if (Test-Path $stopScript) {
            # 运行停止脚本
            $stopResult = & $stopScript 2>&1
            if ($stopResult -like "*停止成功*" -or $stopResult -like "*未运行*" -or $stopResult -notlike "*Error:*") {
                Record-TestResult "停止脚本测试" $true "停止脚本执行成功"
            } else {
                Record-TestResult "停止脚本测试" $false "停止脚本执行失败" $stopResult
            }
        } else {
            Record-TestResult "停止脚本测试" $false "停止脚本不存在"
        }
        
    } catch {
        Record-TestResult "停止脚本测试" $false "测试过程中出现错误" $_.Exception.Message
    }
    
    return $true
}

# 生成测试报告
function Generate-Test-Report {
    Write-Title "测试报告: $PROJECT_NAME v$VERSION"
    
    Write-Host "测试结果汇总:"
    Write-Host "========================================"
    Write-Host "总测试数: $totalTests"
    Write-Host "通过测试: $passedTests" -ForegroundColor Green
    Write-Host "失败测试: $failedTests" -ForegroundColor Red
    
    $passRate = if ($totalTests -gt 0) { [math]::Round(($passedTests / $totalTests) * 100, 2) } else { 0 }
    Write-Host "通过率: $passRate%"
    
    if ($failedTests -eq 0) {
        Write-Host ""
        Write-Host "🎉 所有测试通过！系统运行正常。" -ForegroundColor Green
    } else {
        Write-Host ""
        Write-Host "⚠️  存在失败的测试，请检查上述错误信息。" -ForegroundColor Yellow
    }
    
    Write-Host ""
    Write-Host "详细测试结果:"
    Write-Host "========================================"
    
    foreach ($result in $testResults) {
        $status = if ($result.Passed) { "PASS" } else { "FAIL" }
        $color = if ($result.Passed) { "Green" } else { "Red" }
        
        Write-Host "[$status] $($result.TestName)" -ForegroundColor $color
        if ($result.Message) {
            Write-Host "       $($result.Message)" -ForegroundColor Gray
        }
        if ($result.ErrorInfo) {
            Write-Host "       Error: $($result.ErrorInfo)" -ForegroundColor DarkYellow
        }
    }
    
    Write-Host ""
    Write-Host "测试完成时间: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
    Write-Host "========================================"
}

# 主测试函数
function Run-Tests {
    Write-Title "开始测试 $PROJECT_NAME v$VERSION"
    
    Write-Host "测试环境信息:"
    Write-Host "========================================"
    Write-Host "操作系统: $([System.Environment]::OSVersion.VersionString)"
    Write-Host "PowerShell: $($PSVersionTable.PSVersion.ToString())"
    Write-Host ".NET版本: $([System.Environment]::Version.ToString())"
    Write-Host "安装目录: $INSTALL_DIR"
    Write-Host "ELR目录: $ELR_DIR"
    Write-Host "容器名称: $CONTAINER_NAME"
    Write-Host "========================================"
    Write-Host ""
    
    # 测试安装状态
    $installTest = Test-Installation
    
    # 测试ELR集成
    if ($installTest) {
        Test-ELR-Integration
    }
    
    # 测试对话式交互
    if ($installTest) {
        Test-Dialogue-Interface
    }
    
    # 测试系统功能
    if ($installTest) {
        Test-System-Functionality
    }
    
    # 生成测试报告
    Generate-Test-Report
}

# 开始测试
Run-Tests
