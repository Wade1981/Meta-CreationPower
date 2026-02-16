#!/usr/bin/env powershell
<#
.SYNOPSIS
    自动Word转Markdown转换器

.DESCRIPTION
    此脚本会自动检查Python环境，如果没有则下载便携式Python，
    然后运行Word转Markdown转换器，支持在干净的机器上运行。

.PARAMETER InputFile
    输入Word文档路径

.PARAMETER OutputFile
    输出Markdown文件路径

.EXAMPLE
    .\auto_convert_word_to_md.ps1 -InputFile "docs\程序员称谓衍进数智工程师.docx"

.EXAMPLE
    .\auto_convert_word_to_md.ps1 -InputFile "docs\程序员称谓衍进数智工程师.docx" -OutputFile "docs\程序员称谓衍进数智工程师.md"
#>

param(
    [Parameter(Mandatory=$true, HelpMessage="输入Word文档路径")]
    [string]$InputFile,
    
    [Parameter(Mandatory=$false, HelpMessage="输出Markdown文件路径")]
    [string]$OutputFile
)

# 检查输入文件是否存在
if (-not (Test-Path $InputFile)) {
    Write-Host "错误: 输入文件不存在 - $InputFile"
    return 1
}

# 如果未指定输出路径，使用默认路径
if (-not $OutputFile) {
    $baseName = [System.IO.Path]::GetFileNameWithoutExtension($InputFile)
    $OutputFile = "$baseName.md"
}

Write-Host "========================================"
Write-Host "自动Word转Markdown转换器"
Write-Host "========================================"
Write-Host "输入文件: $InputFile"
Write-Host "输出文件: $OutputFile"
Write-Host "========================================"

# 检查Python环境
function Test-Python {
    try {
        $pythonVersion = python --version 2>&1
        Write-Host "发现系统Python: $pythonVersion"
        return $true
    } catch {
        Write-Host "系统中未找到Python环境"
        return $false
    }
}

# 下载便携式Python
function Download-PortablePython {
    $pythonUrl = "https://www.python.org/ftp/python/3.9.13/python-3.9.13-embed-amd64.zip"
    $pythonZip = "python-portable.zip"
    $pythonDir = "python-portable"
    $pythonExe = "$pythonDir\python.exe"
    
    if (Test-Path $pythonExe) {
        Write-Host "便携式Python已存在，跳过下载"
        return $pythonExe
    }
    
    Write-Host "正在下载便携式Python..."
    Write-Host "URL: $pythonUrl"
    
    try {
        # 创建Python目录
        New-Item -ItemType Directory -Path $pythonDir -Force | Out-Null
        
        # 下载Python
        Invoke-WebRequest -Uri $pythonUrl -OutFile $pythonZip -ErrorAction Stop
        Write-Host "下载完成: $pythonZip"
        
        # 解压缩
        Write-Host "正在解压缩便携式Python..."
        Expand-Archive -Path $pythonZip -DestinationPath $pythonDir -Force -ErrorAction Stop
        Write-Host "解压缩完成"
        
        # 清理
        Remove-Item $pythonZip -Force
        
        # 检查Python是否可用
        if (Test-Path $pythonExe) {
            Write-Host "便携式Python准备就绪: $pythonExe"
            return $pythonExe
        } else {
            Write-Host "错误: 无法找到便携式Python可执行文件"
            return $null
        }
        
    } catch {
        Write-Host "错误下载便携式Python: $($_.Exception.Message)"
        # 清理
        if (Test-Path $pythonDir) {
            Remove-Item $pythonDir -Recurse -Force
        }
        if (Test-Path $pythonZip) {
            Remove-Item $pythonZip -Force
        }
        return $null
    }
}

# 检查Word转Markdown转换器脚本是否存在
function Test-ConverterScript {
    $scriptPath = "word_to_md.py"
    if (Test-Path $scriptPath) {
        Write-Host "发现Word转Markdown转换器脚本"
        return $scriptPath
    } else {
        Write-Host "错误: 未找到Word转Markdown转换器脚本"
        return $null
    }
}

# 执行转换
function Convert-WordToMarkdown {
    param(
        [string]$PythonPath,
        [string]$ConverterScript,
        [string]$InputFile,
        [string]$OutputFile
    )
    
    Write-Host "========================================"
    Write-Host "开始执行Word转Markdown转换"
    Write-Host "使用Python: $PythonPath"
    Write-Host "使用转换器: $ConverterScript"
    Write-Host "========================================"
    
    try {
        # 构建命令
        $command = "$PythonPath $ConverterScript '$InputFile' -o '$OutputFile'"
        Write-Host "执行命令: $command"
        
        # 执行转换
        $output = & $PythonPath $ConverterScript $InputFile -o $OutputFile 2>&1
        Write-Host $output
        
        # 检查转换结果
        if (Test-Path $OutputFile) {
            $fileSize = (Get-Item $OutputFile).Length
            Write-Host "========================================"
            Write-Host "转换成功!"
            Write-Host "输入文件: $InputFile"
            Write-Host "输出文件: $OutputFile"
            Write-Host "输出文件大小: $fileSize 字节"
            Write-Host "========================================"
            return $true
        } else {
            Write-Host "========================================"
            Write-Host "错误: 转换失败，输出文件不存在"
            Write-Host "========================================"
            return $false
        }
        
    } catch {
        Write-Host "========================================"
        Write-Host "错误执行转换: $($_.Exception.Message)"
        Write-Host "========================================"
        return $false
    }
}

# 主执行流程
Write-Host "检查Python环境..."
$pythonPath = $null

if (Test-Python) {
    $pythonPath = "python"
} else {
    Write-Host "系统Python不可用，尝试使用便携式Python..."
    $pythonPath = Download-PortablePython
    if (-not $pythonPath) {
        Write-Host "错误: 无法获取Python环境"
        return 1
    }
}

# 检查转换器脚本
$converterScript = Test-ConverterScript
if (-not $converterScript) {
    Write-Host "错误: 缺少Word转Markdown转换器脚本"
    return 1
}

# 执行转换
$success = Convert-WordToMarkdown -PythonPath $pythonPath -ConverterScript $converterScript -InputFile $InputFile -OutputFile $OutputFile

if ($success) {
    Write-Host "========================================"
    Write-Host "Word转Markdown转换完成!"
    Write-Host "========================================"
    return 0
} else {
    Write-Host "========================================"
    Write-Host "Word转Markdown转换失败!"
    Write-Host "========================================"
    return 1
}
