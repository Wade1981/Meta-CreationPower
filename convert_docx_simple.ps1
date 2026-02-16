#!/usr/bin/env powershell
<#
.SYNOPSIS
    将Word文档(.docx)转换为Markdown格式

.DESCRIPTION
    此脚本使用PowerShell的内置功能来解压缩Word文档并解析XML内容，
    无需安装Python或Microsoft Word，适合在干净的机器上运行。

.PARAMETER InputFile
    输入Word文档路径

.PARAMETER OutputFile
    输出Markdown文件路径

.EXAMPLE
    .\convert_docx_simple.ps1 -InputFile "docs\程序员称谓衍进数智工程师.docx"

.EXAMPLE
    .\convert_docx_simple.ps1 -InputFile "docs\程序员称谓衍进数智工程师.docx" -OutputFile "docs\程序员称谓衍进数智工程师.md"
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

Write-Host "开始转换Word文档到Markdown"
Write-Host "输入文件: $InputFile"
Write-Host "输出文件: $OutputFile"

# 尝试读取和解析Word文档
try {
    # 检查文件扩展名
    $extension = [System.IO.Path]::GetExtension($InputFile).ToLower()
    if ($extension -ne ".docx") {
        Write-Host "错误: 仅支持.docx格式的文件"
        return 1
    }
    
    # 创建临时目录
    $tempDir = Join-Path $env:TEMP "docx2md_$(Get-Random)"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    
    # 复制.docx文件并重命名为.zip
    $zipFile = Join-Path $tempDir "temp.zip"
    Copy-Item $InputFile $zipFile -Force
    
    # 解压缩ZIP文件
    $extractDir = Join-Path $tempDir "extract"
    New-Item -ItemType Directory -Path $extractDir -Force | Out-Null
    
    try {
        Expand-Archive -Path $zipFile -DestinationPath $extractDir -Force
    } catch {
        Write-Host "错误: 无法解压缩Word文档 - $($_.Exception.Message)"
        Remove-Item $tempDir -Recurse -Force
        return 1
    }
    
    # 检查document.xml文件是否存在
    $documentXmlPath = Join-Path $extractDir "word" "document.xml"
    if (-not (Test-Path $documentXmlPath)) {
        Write-Host "错误: 无法找到文档内容"
        Remove-Item $tempDir -Recurse -Force
        return 1
    }
    
    Write-Host "成功打开文档，正在提取内容..."
    
    # 读取XML内容
    try {
        $xmlContent = Get-Content $documentXmlPath -Encoding UTF8 -Raw
    } catch {
        Write-Host "错误: 无法读取文档内容 - $($_.Exception.Message)"
        Remove-Item $tempDir -Recurse -Force
        return 1
    }
    
    # 提取文本内容
    # 移除XML标签
    $plainText = $xmlContent -replace '<[^>]+>', ''
    # 移除多余的空白字符
    $plainText = $plainText -replace '\s+', ' ' -replace '^\s+|\s+$', ''
    # 分割成段落
    $paragraphs = $plainText -split '\s{2,}' | Where-Object { $_.Trim() -ne '' }
    
    # 清理临时目录
    Remove-Item $tempDir -Recurse -Force
    
    # 将内容写入Markdown文件
    Write-Host "正在写入Markdown文件..."
    $paragraphs | Set-Content -Path $OutputFile -Encoding UTF8
    
    # 检查输出文件是否存在
    if (Test-Path $OutputFile) {
        $fileSize = (Get-Item $OutputFile).Length
        Write-Host "成功: 文档转换完成"
        Write-Host "输入文件: $InputFile"
        Write-Host "输出文件: $OutputFile"
        Write-Host "输出文件大小: $fileSize 字节"
        return 0
    } else {
        Write-Host "错误: 转换后文件不存在 - $OutputFile"
        return 1
    }
    
} catch {
    Write-Host "错误: 转换文档时出错 - $($_.Exception.Message)"
    return 1
}
