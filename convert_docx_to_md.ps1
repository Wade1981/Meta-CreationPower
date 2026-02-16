#!/usr/bin/env powershell
<#
.SYNOPSIS
    将Word文档(.docx)转换为Markdown格式

.DESCRIPTION
    此脚本使用PowerShell的COM对象来读取Word文档并转换为Markdown格式，
    无需安装Python，可直接在Windows系统上运行。

.PARAMETER InputFile
    输入Word文档路径

.PARAMETER OutputFile
    输出Markdown文件路径

.EXAMPLE
    .\convert_docx_to_md.ps1 -InputFile "docs\程序员称谓衍进数智工程师.docx"

.EXAMPLE
    .\convert_docx_to_md.ps1 -InputFile "docs\程序员称谓衍进数智工程师.docx" -OutputFile "docs\程序员称谓衍进数智工程师.md"
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

# 尝试使用COM对象打开Word文档
try {
    Write-Host "正在打开Word文档..."
    $word = New-Object -ComObject Word.Application
    $word.Visible = $false
    
    $doc = $word.Documents.Open($InputFile)
    Write-Host "成功打开文档: $InputFile"
    
    # 创建Markdown内容
    $mdContent = @()
    
    # 遍历文档中的所有段落
    Write-Host "正在处理文档内容..."
    foreach ($para in $doc.Paragraphs) {
        if ($para.Range.Text -ne "\r\n") {
            # 获取段落样式
            $styleName = $para.Style.NameLocal
            
            # 根据样式转换为Markdown
            switch ($styleName) {
                "标题 1" {
                    $mdContent += "# $($para.Range.Text.Trim())\n"
                }
                "标题 2" {
                    $mdContent += "## $($para.Range.Text.Trim())\n"
                }
                "标题 3" {
                    $mdContent += "### $($para.Range.Text.Trim())\n"
                }
                "标题 4" {
                    $mdContent += "#### $($para.Range.Text.Trim())\n"
                }
                "标题 5" {
                    $mdContent += "##### $($para.Range.Text.Trim())\n"
                }
                "标题 6" {
                    $mdContent += "###### $($para.Range.Text.Trim())\n"
                }
                default {
                    $mdContent += "$($para.Range.Text.Trim())\n"
                }
            }
        }
    }
    
    # 遍历文档中的所有表格
    foreach ($table in $doc.Tables) {
        # 处理表头
        $headerRow = $table.Rows[1]
        $headerCells = @()
        
        for ($i = 1; $i -le $headerRow.Cells.Count; $i++) {
            $headerCells += $headerRow.Cells[$i].Range.Text.Trim()
        }
        
        # 添加表头
        $mdContent += "| " + ($headerCells -join " | ") + " |\n"
        
        # 添加分隔线
        $mdContent += "| " + ($headerCells | ForEach-Object { "---" }) -join " | " + " |\n"
        
        # 处理表格内容
        for ($rowIndex = 2; $rowIndex -le $table.Rows.Count; $rowIndex++) {
            $row = $table.Rows[$rowIndex]
            $rowCells = @()
            
            for ($i = 1; $i -le $row.Cells.Count; $i++) {
                $rowCells += $row.Cells[$i].Range.Text.Trim()
            }
            
            $mdContent += "| " + ($rowCells -join " | ") + " |\n"
        }
        
        $mdContent += "\n"
    }
    
    # 关闭文档和Word应用
    $doc.Close()
    $word.Quit()
    
    # 释放COM对象
    [System.Runtime.Interopservices.Marshal]::ReleaseComObject($doc) | Out-Null
    [System.Runtime.Interopservices.Marshal]::ReleaseComObject($word) | Out-Null
    [System.GC]::Collect()
    [System.GC]::WaitForPendingFinalizers()
    
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
    
    # 尝试清理COM对象
    try {
        if ($doc) {
            $doc.Close()
            [System.Runtime.Interopservices.Marshal]::ReleaseComObject($doc) | Out-Null
        }
        if ($word) {
            $word.Quit()
            [System.Runtime.Interopservices.Marshal]::ReleaseComObject($word) | Out-Null
        }
        [System.GC]::Collect()
        [System.GC]::WaitForPendingFinalizers()
    } catch {}
    
    return 1
}
