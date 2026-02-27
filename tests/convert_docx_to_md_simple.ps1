#!/usr/bin/env powershell
<#
.SYNOPSIS
    将Word文档(.docx)转换为Markdown格式

.DESCRIPTION
    此脚本使用PowerShell的内置功能直接读取和解析Word文档的XML结构，
    无需安装Python或Microsoft Word，可直接在Windows系统上运行。

.PARAMETER InputFile
    输入Word文档路径

.PARAMETER OutputFile
    输出Markdown文件路径

.EXAMPLE
    .\convert_docx_to_md_simple.ps1 -InputFile "docs\程序员称谓衍进数智工程师.docx"

.EXAMPLE
    .\convert_docx_to_md_simple.ps1 -InputFile "docs\程序员称谓衍进数智工程师.docx" -OutputFile "docs\程序员称谓衍进数智工程师.md"
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

# 尝试读取和解析Word文档
try {
    Write-Host "正在打开Word文档..."
    
    # 检查文件扩展名
    $extension = [System.IO.Path]::GetExtension($InputFile).ToLower()
    if ($extension -ne ".docx") {
        Write-Host "错误: 仅支持.docx格式的文件"
        return 1
    }
    
    # 创建临时目录
    $tempDir = Join-Path $env:TEMP "docx2md_$(Get-Random)"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    
    # 尝试将.docx文件作为ZIP文件解压缩
    try {
        Add-Type -AssemblyName System.IO.Compression.FileSystem
        [System.IO.Compression.ZipFile]::ExtractToDirectory($InputFile, $tempDir)
    } catch {
        Write-Host "错误: 无法解压缩Word文档 - $($_.Exception.Message)"
        Remove-Item $tempDir -Recurse -Force
        return 1
    }
    
    # 检查document.xml文件是否存在
    $documentXmlPath = Join-Path $tempDir "word" "document.xml"
    if (-not (Test-Path $documentXmlPath)) {
        Write-Host "错误: 无法找到文档内容"
        Remove-Item $tempDir -Recurse -Force
        return 1
    }
    
    Write-Host "成功打开文档: $InputFile"
    Write-Host "正在处理文档内容..."
    
    # 读取并解析XML内容
    [xml]$xmlContent = Get-Content $documentXmlPath -Encoding UTF8
    
    # 创建Markdown内容
    $mdContent = @()
    
    # 遍历所有段落
    foreach ($p in $xmlContent.document.body.GetElementsByTagName("w:p")) {
        $paraText = ""
        
        # 遍历段落中的所有文本运行
        foreach ($r in $p.GetElementsByTagName("w:r")) {
            foreach ($t in $r.GetElementsByTagName("w:t")) {
                if ($t.InnerText) {
                    $paraText += $t.InnerText
                }
            }
        }
        
        # 检查段落样式
        $styleName = ""
        $pPr = $p.GetElementsByTagName("w:pPr") | Select-Object -First 1
        if ($pPr) {
            $pStyle = $pPr.GetElementsByTagName("w:pStyle") | Select-Object -First 1
            if ($pStyle -and $pStyle."w:val") {
                $styleName = $pStyle."w:val"
            }
        }
        
        # 根据样式转换为Markdown
        if ($paraText.Trim()) {
            switch ($styleName) {
                "Heading1" {
                    $mdContent += "# $($paraText.Trim())\n"
                }
                "Heading2" {
                    $mdContent += "## $($paraText.Trim())\n"
                }
                "Heading3" {
                    $mdContent += "### $($paraText.Trim())\n"
                }
                "Heading4" {
                    $mdContent += "#### $($paraText.Trim())\n"
                }
                "Heading5" {
                    $mdContent += "##### $($paraText.Trim())\n"
                }
                "Heading6" {
                    $mdContent += "###### $($paraText.Trim())\n"
                }
                default {
                    $mdContent += "$($paraText.Trim())\n"
                }
            }
        }
    }
    
    # 清理临时目录
    Remove-Item $tempDir -Recurse -Force
    
    # 将内容写入Markdown文件
    Write-Host "正在写入Markdown文件..."
    $mdContent | Set-Content -Path $OutputFile -Encoding UTF8
    
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
