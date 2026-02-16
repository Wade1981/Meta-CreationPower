#!/usr/bin/env powershell
<#
.SYNOPSIS
    在ELR中运行文档转换工具

.DESCRIPTION
    此脚本用于在Enlightenment Lighthouse Runtime (ELR)中运行文档转换工具，
    将Word文档(.doc/.docx)和PowerPoint文档(.ppt)转换为Markdown格式。

.PARAMETER InputFile
    输入文件路径

.PARAMETER OutputFile
    输出Markdown文件路径

.EXAMPLE
    .\run_doc_converter.ps1 -InputFile "docs\程序员称谓衍进数智工程师.docx"

.EXAMPLE
    .\run_doc_converter.ps1 -InputFile "docs\程序员称谓衍进数智工程师.docx" -OutputFile "docs\程序员称谓衍进数智工程师.md"
#>

param(
    [Parameter(Mandatory=$true, HelpMessage="输入文件路径")]
    [string]$InputFile,
    
    [Parameter(Mandatory=$false, HelpMessage="输出Markdown文件路径")]
    [string]$OutputFile
)

# 检查Python是否安装
function Check-Python {
    try {
        $pythonVersion = python --version 2>&1
        Write-Host "Python 已安装: $pythonVersion"
        return $true
    } catch {
        Write-Host "错误: Python 未安装或不在PATH中"
        return $false
    }
}

# 安装必要的依赖
function Install-Dependencies {
    Write-Host "正在检查和安装必要的依赖..."
    
    try {
        # 检查python-docx
        $docxInstalled = python -c "import docx; print('Installed')" 2>$null
        if (-not $docxInstalled) {
            Write-Host "安装 python-docx..."
            python -m pip install python-docx
        } else {
            Write-Host "python-docx 已安装"
        }
        
        # 检查python-pptx
        $pptxInstalled = python -c "import pptx; print('Installed')" 2>$null
        if (-not $pptxInstalled) {
            Write-Host "安装 python-pptx..."
            python -m pip install python-pptx
        } else {
            Write-Host "python-pptx 已安装"
        }
        
        return $true
    } catch {
        Write-Host "错误: 安装依赖时出错 - $($_.Exception.Message)"
        return $false
    }
}

# 运行文档转换工具
function Run-Converter {
    param(
        [string]$InputPath,
        [string]$OutputPath
    )
    
    Write-Host "正在运行文档转换工具..."
    
    try {
        if (-not $OutputPath) {
            # 如果未指定输出路径，使用默认路径
            $baseName = [System.IO.Path]::GetFileNameWithoutExtension($InputPath)
            $OutputPath = "$baseName.md"
        }
        
        # 构建命令
        $command = "python src\utils\doc_to_md_converter.py '$InputPath' -o '$OutputPath'"
        
        # 运行命令
        Write-Host "执行命令: $command"
        Invoke-Expression $command
        
        # 检查转换是否成功
        if (Test-Path $OutputPath) {
            Write-Host "成功: 文档转换完成，输出文件 - $OutputPath"
            return $true
        } else {
            Write-Host "错误: 转换后文件不存在 - $OutputPath"
            return $false
        }
    } catch {
        Write-Host "错误: 运行转换器时出错 - $($_.Exception.Message)"
        return $false
    }
}

# 主函数
function Main {
    Write-Host "在ELR中运行文档转换工具"
    Write-Host "=" * 70
    
    # 检查输入文件是否存在
    if (-not (Test-Path $InputFile)) {
        Write-Host "错误: 输入文件不存在 - $InputFile"
        return 1
    }
    
    # 检查Python
    if (-not (Check-Python)) {
        return 1
    }
    
    # 安装依赖
    if (-not (Install-Dependencies)) {
        return 1
    }
    
    # 运行转换器
    if (-not (Run-Converter -InputPath $InputFile -OutputPath $OutputFile)) {
        return 1
    }
    
    Write-Host "=" * 70
    Write-Host "成功: 文档转换工具运行完成"
    return 0
}

# 执行主函数
exit (Main)
