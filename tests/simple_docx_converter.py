#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
简单的Word文档转Markdown工具

使用Python内置库，无需安装任何依赖
"""

import os
import sys
import zipfile
import xml.etree.ElementTree as ET


def convert_docx_to_md(docx_path, md_path):
    """
    将Word文档(.docx)转换为Markdown格式
    
    Args:
        docx_path (str): Word文档路径
        md_path (str): 输出Markdown文件路径
    """
    print(f"转换文档: {docx_path} -> {md_path}")
    
    try:
        # 注册命名空间
        ET.register_namespace('', 'http://schemas.openxmlformats.org/wordprocessingml/2006/main')
        nsmap = {'w': 'http://schemas.openxmlformats.org/wordprocessingml/2006/main'}
        
        # 打开docx文件（本质是zip文件）
        with zipfile.ZipFile(docx_path, 'r') as zf:
            # 检查document.xml是否存在
            if 'word/document.xml' not in zf.namelist():
                print("错误: 无法找到文档内容")
                return False
            
            # 读取document.xml文件
            with zf.open('word/document.xml') as f:
                tree = ET.parse(f)
                root = tree.getroot()
                
                # 提取文本
                text = []
                
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
        
        # 写入Markdown文件
        with open(md_path, 'w', encoding='utf-8') as f:
            f.write('\n'.join(text))
        
        print(f"成功: 转换完成，输出文件 - {md_path}")
        return True
        
    except Exception as e:
        print(f"错误: 转换文档时出错 - {str(e)}")
        return False


def main():
    """
    主函数
    """
    if len(sys.argv) < 2:
        print("用法: python simple_docx_converter.py <input.docx> [output.md]")
        return 1
    
    input_file = sys.argv[1]
    output_file = sys.argv[2] if len(sys.argv) > 2 else os.path.splitext(input_file)[0] + '.md'
    
    success = convert_docx_to_md(input_file, output_file)
    return 0 if success else 1


if __name__ == "__main__":
    sys.exit(main())
