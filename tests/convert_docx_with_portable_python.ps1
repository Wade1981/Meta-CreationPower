#!/usr/bin/env powershell
<#
.SYNOPSIS
    使用便携式Python将Word文档(.docx)转换为Markdown格式

.DESCRIPTION
    此脚本会自动下载便携式Python，然后使用它来运行文档转换工具，
    无需安装Python或Microsoft Word，可在任何Windows系统上运行。

.PARAMETER InputFile
    输入Word文档路径

.PARAMETER OutputFile
    输出Markdown文件路径

.EXAMPLE
    .\convert_docx_with_portable_python.ps1 -InputFile "docs\程序员称谓衍进数智工程师.docx"

.EXAMPLE
    .\convert_docx_with_portable_python.ps1 -InputFile "docs\程序员称谓衍进数智工程师.docx" -OutputFile "docs\程序员称谓衍进数智工程师.md"
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
    $OutputFile = [System.IO.Path]::ChangeExtension($InputFile, ".md")
}

# Python portable version information
$PYTHON_PORTABLE_URL = "https://www.python.org/ftp/python/3.9.13/python-3.9.13-embed-amd64.zip"
$PYTHON_PORTABLE_ZIP = "python-portable.zip"
$PYTHON_DIR = "python-portable"
$PYTHON_EXE = "$PYTHON_DIR\python.exe"

# Function: Download portable Python
function Download-PortablePython {
    if (Test-Path $PYTHON_EXE) {
        Write-Host "便携式Python已存在: $PYTHON_EXE"
        return $true
    }

    Write-Host "正在下载便携式Python..."
    Write-Host "下载地址: $PYTHON_PORTABLE_URL"
    
    try {
        # 创建临时目录
        New-Item -ItemType Directory -Path $PYTHON_DIR -Force | Out-Null
        
        # 下载Python便携式版本
        Invoke-WebRequest -Uri $PYTHON_PORTABLE_URL -OutFile $PYTHON_PORTABLE_ZIP -ErrorAction Stop
        Write-Host "下载完成: $PYTHON_PORTABLE_ZIP"
        
        # 解压缩
        Write-Host "正在解压缩便携式Python..."
        Expand-Archive -Path $PYTHON_PORTABLE_ZIP -DestinationPath $PYTHON_DIR -Force -ErrorAction Stop
        Write-Host "解压缩完成"
        
        # 清理
        Remove-Item $PYTHON_PORTABLE_ZIP -Force
        
        # 检查Python是否可用
        if (Test-Path $PYTHON_EXE) {
            Write-Host "便携式Python准备就绪: $PYTHON_EXE"
            return $true
        } else {
            Write-Host "错误: 未找到便携式Python可执行文件"
            return $false
        }
        
    } catch {
        Write-Host "错误: 下载便携式Python时出错 - $($_.Exception.Message)"
        # 清理
        if (Test-Path $PYTHON_DIR) {
            Remove-Item $PYTHON_DIR -Recurse -Force
        }
        if (Test-Path $PYTHON_PORTABLE_ZIP) {
            Remove-Item $PYTHON_PORTABLE_ZIP -Force
        }
        return $false
    }
}

# 创建文档转换脚本
function Create-DocumentConverter {
    $converterScript = @'
import os
import sys
from zipfile import ZipFile
from xml.etree import ElementTree as ET

# 注册命名空间
ET.register_namespace('', 'http://schemas.openxmlformats.org/wordprocessingml/2006/main')
nsmap = {'w': 'http://schemas.openxmlformats.org/wordprocessingml/2006/main'}

def extract_text_from_docx(docx_path):
    """
    从docx文件中提取文本内容
    """
    text = []
    
    try:
        # 打开docx文件（本质是zip文件）
        with ZipFile(docx_path, 'r') as zf:
            # 读取document.xml文件
            with zf.open('word/document.xml') as f:
                tree = ET.parse(f)
                root = tree.getroot()
                
                # 遍历所有段落
                for para in root.findall('.//w:p', namespaces=nsmap):
                    para_text = []
                    
                    # 遍历段落中的所有文本运行
                    for run in para.findall('.//w:r', namespaces=nsmap):
                        for text_elem in run.findall('.//w:t', namespaces=nsmap):
                            if text_elem.text:
                                para_text.append(text_elem.text)
                    
                    # 检查段落样式
                    style_name = ""
                    pPr = para.find('.//w:pPr', namespaces=nsmap)
                    if pPr:
                        pStyle = pPr.find('.//w:pStyle', namespaces=nsmap)
                        if pStyle is not None and 'w:val' in pStyle.attrib:
                            style_name = pStyle.attrib['w:val']
                    
                    # 合并段落文本
                    para_text_str = ''.join(para_text)
                    if para_text_str.strip():
                        # 根据样式添加Markdown格式
                        if style_name == 'Heading1':
                            text.append(f"# {para_text_str.strip()}")
                        elif style_name == 'Heading2':
                            text.append(f"## {para_text_str.strip()}")
                        elif style_name == 'Heading3':
                            text.append(f"### {para_text_str.strip()}")
                        elif style_name == 'Heading4':
                            text.append(f"#### {para_text_str.strip()}")
                        elif style_name == 'Heading5':
                            text.append(f"##### {para_text_str.strip()}")
                        elif style_name == 'Heading6':
                            text.append(f"###### {para_text_str.strip()}")
                        else:
                            text.append(para_text_str.strip())
        
        return '\n'.join(text)
        
    except Exception as e:
        print(f"Error extracting text: {e}")
        return ""

def convert_docx_to_md(docx_path, md_path):
    """
    将docx文件转换为Markdown格式
    """
    print(f"Converting: {docx_path} -> {md_path}")
    
    # 提取文本
    text = extract_text_from_docx(docx_path)
    
    if not text:
        print("Error: Failed to extract text from document")
        return False
    
    # 写入Markdown文件
    try:
        with open(md_path, 'w', encoding='utf-8') as f:
            f.write(text)
        print(f"Successfully converted to: {md_path}")
        return True
    except Exception as e:
        print(f"Error writing Markdown file: {e}")
        return False

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python docx_to_md.py <input.docx> [output.md]")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2] if len(sys.argv) > 2 else os.path.splitext(input_file)[0] + '.md'
    
    success = convert_docx_to_md(input_file, output_file)
    sys.exit(0 if success else 1)
'@
    
    $converterScriptPath = "docx_to_md.py"
    $converterScript | Set-Content -Path $converterScriptPath -Encoding UTF8
    
    return $converterScriptPath
}

# 主函数
function Main {
    Write-Host "使用便携式Python将Word文档转换为Markdown格式"
    Write-Host "=============================================="
    Write-Host "输入文件: $InputFile"
    Write-Host "输出文件: $OutputFile"
    Write-Host "=============================================="
    
    # 下载便携式Python
    if (-not (Download-PortablePython)) {
        Write-Host "错误: 无法下载便携式Python"
        return 1
    }
    
    # 创建文档转换脚本
    $converterScript = Create-DocumentConverter
    
    # 运行文档转换工具
    if ($converterScript) {
        Write-Host "正在运行文档转换工具..."
        
        try {
            # 构建命令
            $command = "& '$PYTHON_EXE' '$converterScript' '$InputFile' '$OutputFile'"
            Write-Host "执行命令: $command"
            
            # 执行命令
            Invoke-Expression $command
            $exitCode = $LASTEXITCODE
            
            Write-Host "执行完成，退出代码: $exitCode"
            
            # 检查转换是否成功
            if (Test-Path $OutputFile) {
                $fileSize = (Get-Item $OutputFile).Length
                Write-Host "=============================================="
                Write-Host "转换成功!"
                Write-Host "输出文件: $OutputFile"
                Write-Host "文件大小: $fileSize 字节"
                Write-Host "=============================================="
                
                # 清理
                if (Test-Path $converterScript) {
                    Remove-Item $converterScript -Force
                }
                
                return 0
            } else {
                Write-Host "=============================================="
                Write-Host "错误: 转换失败，输出文件未找到"
                Write-Host "=============================================="
                
                # 清理
                if (Test-Path $converterScript) {
                    Remove-Item $converterScript -Force
                }
                
                return 1
            }
            
        } catch {
            Write-Host "错误: $($_.Exception.Message)"
            
            # 清理
            if (Test-Path $converterScript) {
                Remove-Item $converterScript -Force
            }
            
            return 1
        }
    } else {
        Write-Host "错误: 无法创建文档转换脚本"
        return 1
    }
}

# 执行主函数
exit (Main)
