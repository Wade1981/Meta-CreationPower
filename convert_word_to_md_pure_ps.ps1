#!/usr/bin/env powershell
<#
.SYNOPSIS
    纯PowerShell实现的Word转Markdown转换器

.DESCRIPTION
    此脚本完全使用PowerShell内置功能，无需Python或Microsoft Word，
    直接解析.docx文件的XML结构并转换为Markdown格式。

.PARAMETER InputFile
    输入Word文档路径

.PARAMETER OutputFile
    输出Markdown文件路径

.EXAMPLE
    .\convert_word_to_md_pure_ps.ps1 -InputFile "docs\程序员称谓衍进数智工程师.docx"

.EXAMPLE
    .\convert_word_to_md_pure_ps.ps1 -InputFile "docs\程序员称谓衍进数智工程师.docx" -OutputFile "docs\程序员称谓衍进数智工程师.md"
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
Write-Host "纯PowerShell Word转Markdown转换器"
Write-Host "========================================"
Write-Host "输入文件: $InputFile"
Write-Host "输出文件: $OutputFile"
Write-Host "========================================"

# 检查文件扩展名
$extension = [System.IO.Path]::GetExtension($InputFile).ToLower()
if ($extension -ne ".docx") {
    Write-Host "错误: 仅支持.docx格式的文件"
    return 1
}

# 主转换函数
function Convert-DocxToMarkdown {
    param(
        [string]$DocxPath,
        [string]$OutputPath
    )
    
    try {
        Write-Host "正在处理Word文档..."
        
        # 创建临时目录
        $tempDir = Join-Path $env:TEMP "docx2md_$(Get-Random)"
        New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
        
        # 复制.docx文件并重命名为.zip
        $zipFile = Join-Path $tempDir "temp.zip"
        Copy-Item $DocxPath $zipFile -Force
        
        # 解压缩ZIP文件
        $extractDir = Join-Path $tempDir "extract"
        New-Item -ItemType Directory -Path $extractDir -Force | Out-Null
        
        try {
            Expand-Archive -Path $zipFile -DestinationPath $extractDir -Force
        } catch {
            Write-Host "错误: 无法解压缩Word文档 - $($_.Exception.Message)"
            Remove-Item $tempDir -Recurse -Force
            return $false
        }
        
        # 检查document.xml文件是否存在
        $documentXmlPath = Join-Path $extractDir "word" "document.xml"
        if (-not (Test-Path $documentXmlPath)) {
            Write-Host "错误: 无法找到文档内容"
            Remove-Item $tempDir -Recurse -Force
            return $false
        }
        
        Write-Host "成功打开文档，正在提取内容..."
        
        # 读取XML内容
        try {
            [xml]$xmlContent = Get-Content $documentXmlPath -Encoding UTF8 -Raw
        } catch {
            Write-Host "错误: 无法读取文档内容 - $($_.Exception.Message)"
            Remove-Item $tempDir -Recurse -Force
            return $false
        }
        
        # 提取文本内容
        $markdownContent = Extract-TextFromXml -XmlContent $xmlContent
        
        # 清理临时目录
        Remove-Item $tempDir -Recurse -Force
        
        # 写入Markdown文件
        Write-Host "正在写入Markdown文件..."
        try {
            $markdownContent | Set-Content -Path $OutputPath -Encoding UTF8
        } catch {
            Write-Host "错误: 无法写入输出文件 - $($_.Exception.Message)"
            return $false
        }
        
        # 检查输出文件是否存在
        if (Test-Path $OutputPath) {
            $fileSize = (Get-Item $OutputPath).Length
            Write-Host "========================================"
            Write-Host "转换成功!"
            Write-Host "输入文件: $DocxPath"
            Write-Host "输出文件: $OutputPath"
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
        Write-Host "错误: 转换文档时出错 - $($_.Exception.Message)"
        Write-Host "========================================"
        return $false
    }
}

# 从XML中提取文本
function Extract-TextFromXml {
    param(
        [xml]$XmlContent
    )
    
    $lines = @()
    
    # 注册命名空间
    $ns = @{
        w = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"
    }
    
    # 查找所有段落
    $paragraphs = $XmlContent.document.body.SelectNodes('//w:p', $ns)
    
    foreach ($para in $paragraphs) {
        # 提取段落文本
        $paraText = Extract-TextFromParagraph -Paragraph $para -Namespace $ns
        
        if ($paraText -and $paraText.Trim()) {
            # 检查是否为标题
            $isHeading = $false
            $headingLevel = 0
            
            # 检查段落样式
            $pPr = $para.SelectSingleNode('.//w:pPr', $ns)
            if ($pPr) {
                $pStyle = $pPr.SelectSingleNode('.//w:pStyle', $ns)
                if ($pStyle -and $pStyle.Val) {
                    $styleName = $pStyle.Val
                    if ($styleName -like 'Heading*') {
                        $isHeading = $true
                        $levelStr = $styleName -replace 'Heading', ''
                        if ([int]::TryParse($levelStr, [ref]$headingLevel) -and $headingLevel -ge 1 -and $headingLevel -le 6) {
                            # 添加标题格式
                            $lines += "$("#" * $headingLevel) $($paraText.Trim())"
                        } else {
                            $lines += $paraText.Trim()
                        }
                    } else {
                        $lines += $paraText.Trim()
                    }
                } else {
                    $lines += $paraText.Trim()
                }
            } else {
                $lines += $paraText.Trim()
            }
        }
    }
    
    # 合并行，添加适当的空行
    $result = @()
    foreach ($line in $lines) {
        if ($line -match '^#+ ') {
            # 标题前添加空行
            if ($result.Count -gt 0 -and $result[-1] -ne '') {
                $result += ''
            }
            $result += $line
            # 标题后添加空行
            $result += ''
        } else {
            $result += $line
        }
    }
    
    return $result -join "`n"
}

# 从段落中提取文本
function Extract-TextFromParagraph {
    param(
        [System.Xml.XmlNode]$Paragraph,
        [hashtable]$Namespace
    )
    
    $textParts = @()
    
    # 查找所有文本运行
    $runs = $Paragraph.SelectNodes('.//w:r', $Namespace)
    
    foreach ($run in $runs) {
        # 查找文本节点
        $textNodes = $run.SelectNodes('.//w:t', $Namespace)
        foreach ($textNode in $textNodes) {
            if ($textNode.InnerText) {
                $textParts += $textNode.InnerText
            }
        }
        
        # 处理制表符
        $tabNodes = $run.SelectNodes('.//w:tab', $Namespace)
        if ($tabNodes.Count -gt 0) {
            $textParts += "`t" * $tabNodes.Count
        }
        
        # 处理换行
        $brNodes = $run.SelectNodes('.//w:br', $Namespace)
        if ($brNodes.Count -gt 0) {
            $textParts += "`n" * $brNodes.Count
        }
    }
    
    return $textParts -join ''
}

# 执行转换
$success = Convert-DocxToMarkdown -DocxPath $InputFile -OutputPath $OutputFile

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
