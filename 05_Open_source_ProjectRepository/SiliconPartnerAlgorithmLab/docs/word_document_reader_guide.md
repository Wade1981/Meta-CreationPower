# Word 文档读取指南

本文档详细介绍如何使用启蒙灯塔运行时（ELR）读取和处理 Word 文档（.docx 文件）。

## 1. 概述

ELR 提供了强大的 Word 文档读取能力，支持以下功能：

- 提取文档文本内容
- 分析文档结构
- 提取特定类型的内容（如标题）
- 将 Word 文档转换为 Markdown 格式

## 2. 系统要求

- Windows 操作系统
- PowerShell 5.0 或更高版本
- 互联网连接（首次运行时需要下载便携式 Python）

## 3. 使用方法

### 3.1 启动 ELR 运行时

首先，确保 ELR 运行时处于运行状态：

```powershell
.lr-word-reader.ps1 start
```

### 3.2 读取 Word 文档

#### 3.2.1 基本文本提取

提取文档的纯文本内容：

```powershell
.lr-word-reader.ps1 read-word --input docs\程序员称谓衍进数智工程师.docx
```

#### 3.2.2 文档结构分析

分析文档的结构，包括段落数、标题数、单词数等：

```powershell
.lr-word-reader.ps1 read-word --input docs\程序员称谓衍进数智工程师.docx --analyze
```

输出示例：

```
Document Analysis Results:
=============================
Total Paragraphs: 25
Total Headings: 5
Total Words: 1200

Headings:
Level 1:
  - 程序员称谓的历史演变
Level 2:
  - 传统程序员阶段
  - 软件工程师阶段
  - 全栈工程师阶段
  - 数智工程师阶段
```

#### 3.2.3 提取特定内容

提取文档中的标题：

```powershell
.lr-word-reader.ps1 read-word --input docs\程序员称谓衍进数智工程师.docx --extract headings
```

提取文档中的普通文本：

```powershell
.lr-word-reader.ps1 read-word --input docs\程序员称谓衍进数智工程师.docx --extract normal
```

#### 3.2.4 不同格式输出

以 JSON 格式输出文档分析结果：

```powershell
.lr-word-reader.ps1 read-word --input docs\程序员称谓衍进数智工程师.docx --analyze --format json
```

## 4. 转换 Word 文档为 Markdown

将 Word 文档转换为 Markdown 格式：

```powershell
.lr-word-reader.ps1 convert-docx --input docs\程序员称谓衍进数智工程师.docx --output docs\程序员称谓衍进数智工程师.md
```

如果不指定输出文件路径，默认会在同一目录下创建同名的 .md 文件：

```powershell
.lr-word-reader.ps1 convert-docx --input docs\程序员称谓衍进数智工程师.docx
```

## 5. 使用专门的 Word 文档处理容器

创建并运行一个专门用于处理 Word 文档的容器：

```powershell
.lr-word-reader.ps1 run-word-container --name doc-analyzer
```

在容器中执行文档处理命令：

```powershell
.lr-word-reader.ps1 exec --id <container-id> --command 'python -c "import zipfile; print(\"Word document processing ready\")"'
```

## 6. 技术原理

### 6.1 Word 文档结构

Word 文档（.docx）本质上是一个 ZIP 文件，包含了以下主要部分：

- `word/document.xml`：包含文档的主要内容和结构
- `word/styles.xml`：包含文档的样式定义
- `word/numbering.xml`：包含文档的编号定义

### 6.2 文档解析过程

ELR 使用以下步骤解析 Word 文档：

1. 将 .docx 文件作为 ZIP 文件打开
2. 读取并解析 `word/document.xml` 文件
3. 遍历 XML 结构，提取段落、文本和样式信息
4. 根据需要转换为不同的输出格式

### 6.3 便携式 Python

为了确保在没有安装 Python 的系统上也能运行，ELR 实现了便携式 Python 支持：

1. 检测系统是否已安装 Python
2. 如果没有，自动下载便携式 Python 3.9
3. 提取并配置便携式 Python 环境
4. 使用便携式 Python 执行文档解析脚本

## 7. 示例代码

### 7.1 基本文档读取

```python
#!/usr/bin/env python3
import zipfile
from xml.etree import ElementTree as ET

# 注册命名空间
nsmap = {'w': 'http://schemas.openxmlformats.org/wordprocessingml/2006/main'}

def read_word_document(docx_path):
    """读取 Word 文档并返回文本内容"""
    text = []
    
    with zipfile.ZipFile(docx_path, 'r') as zf:
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
                
                # 合并段落文本
                para_text_str = ''.join(para_text)
                if para_text_str.strip():
                    text.append(para_text_str.strip())
    
    return '\n'.join(text)

# 使用示例
if __name__ == "__main__":
    docx_path = "docs/程序员称谓衍进数智工程师.docx"
    content = read_word_document(docx_path)
    print(content)
```

### 7.2 文档结构分析

```python
#!/usr/bin/env python3
import zipfile
from xml.etree import ElementTree as ET
import json

# 注册命名空间
nsmap = {'w': 'http://schemas.openxmlformats.org/wordprocessingml/2006/main'}

def analyze_document_structure(docx_path):
    """分析文档结构并返回统计信息"""
    structure = {
        'paragraphs': [],
        'headings': {},
        'statistics': {
            'total_paragraphs': 0,
            'total_headings': 0,
            'total_words': 0
        }
    }
    
    with zipfile.ZipFile(docx_path, 'r') as zf:
        with zf.open('word/document.xml') as f:
            tree = ET.parse(f)
            root = tree.getroot()
            
            # 遍历所有段落
            for para in root.findall('.//w:p', namespaces=nsmap):
                para_info = {
                    'text': '',
                    'style': '',
                    'word_count': 0
                }
                
                # 遍历段落中的所有文本运行
                for run in para.findall('.//w:r', namespaces=nsmap):
                    for text_elem in run.findall('.//w:t', namespaces=nsmap):
                        if text_elem.text:
                            para_info['text'] += text_elem.text
                
                # 检查段落样式
                pPr = para.find('.//w:pPr', namespaces=nsmap)
                if pPr:
                    pStyle = pPr.find('.//w:pStyle', namespaces=nsmap)
                    if pStyle is not None and 'w:val' in pStyle.attrib:
                        para_info['style'] = pStyle.attrib['w:val']
                
                # 计算单词数
                para_info['word_count'] = len(para_info['text'].split())
                
                # 添加到结构中
                structure['paragraphs'].append(para_info)
                structure['statistics']['total_paragraphs'] += 1
                structure['statistics']['total_words'] += para_info['word_count']
                
                # 处理标题
                if para_info['style'].startswith('Heading'):
                    structure['statistics']['total_headings'] += 1
                    heading_level = para_info['style'][-1]
                    if heading_level not in structure['headings']:
                        structure['headings'][heading_level] = []
                    structure['headings'][heading_level].append(para_info['text'])
    
    return structure

# 使用示例
if __name__ == "__main__":
    docx_path = "docs/程序员称谓衍进数智工程师.docx"
    structure = analyze_document_structure(docx_path)
    print(json.dumps(structure, ensure_ascii=False, indent=2))
```

## 8. 常见问题

### 8.1 运行时错误

#### 8.1.1 Python 下载失败

**问题**：首次运行时，便携式 Python 下载失败。

**解决方案**：
- 检查网络连接
- 手动下载 Python 3.9 便携式版本并解压到 `python-portable` 目录
- 确保下载的 Python 版本与系统架构匹配（32位或64位）

#### 8.1.2 文档解析错误

**问题**：解析某些 Word 文档时出现错误。

**解决方案**：
- 确保文档格式正确，没有损坏
- 尝试将文档另存为新的 .docx 文件
- 检查文档是否包含特殊格式或元素

### 8.2 性能问题

#### 8.2.1 大型文档处理缓慢

**问题**：处理大型 Word 文档时速度缓慢。

**解决方案**：
- 对于非常大的文档，可以考虑拆分为多个较小的文档
- 增加系统内存以提高处理速度
- 避免同时处理多个大型文档

## 9. 高级功能

### 9.1 自定义输出格式

ELR 支持自定义输出格式，可以根据需要修改文档解析脚本：

1. 复制 `word_reader.py` 脚本
2. 修改输出格式相关代码
3. 使用 `run-python` 命令执行自定义脚本

### 9.2 批量处理

可以使用 PowerShell 脚本批量处理多个 Word 文档：

```powershell
# 批量转换文档
Get-ChildItem -Path "docs" -Filter "*.docx" | ForEach-Object {
    .\elr-word-reader.ps1 convert-docx --input $_.FullName
}
```

### 9.3 集成到其他系统

ELR 的 Word 文档读取功能可以集成到其他系统中：

1. 作为命令行工具调用
2. 通过 PowerShell 脚本集成
3. 作为服务提供给其他应用程序

## 10. 总结

ELR 提供了强大而灵活的 Word 文档读取能力，通过以下特点为用户提供便捷的文档处理体验：

- **无依赖执行**：内置便携式 Python 支持，无需预先安装 Python
- **容器化隔离**：通过 ELR 容器机制确保处理过程的稳定性
- **多格式支持**：支持文本、JSON 和 Markdown 输出格式
- **扩展性强**：模块化设计，易于添加新功能
- **跨平台潜力**：基于标准库实现，可扩展到其他平台

---

**作者**：启蒙灯塔起源团队
**版本**：1.0.0
**更新日期**：2026-02-16
