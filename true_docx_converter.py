#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
真正的Word文档转Markdown转换器

此脚本能够真正读取Word文档的原文内容，保持格式和结构，
并准确转换为Markdown格式。支持.docx文件格式。

使用方法：
    python true_docx_converter.py input.docx output.md
"""

import os
import sys
import zipfile
import xml.etree.ElementTree as ET
from typing import List, Dict, Optional

class TrueDocxConverter:
    """
    真正的Word文档转Markdown转换器
    能够准确读取和转换Word文档的内容和格式
    """
    
    def __init__(self):
        """初始化转换器"""
        # 注册命名空间
        self.namespaces = {
            'w': 'http://schemas.openxmlformats.org/wordprocessingml/2006/main',
            'wp': 'http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing',
            'a': 'http://schemas.openxmlformats.org/drawingml/2006/main',
            'pic': 'http://schemas.openxmlformats.org/drawingml/2006/picture'
        }
    
    def convert(self, input_path: str, output_path: str) -> bool:
        """
        转换Word文档为Markdown
        
        Args:
            input_path: 输入Word文档路径
            output_path: 输出Markdown文件路径
            
        Returns:
            bool: 转换是否成功
        """
        try:
            print(f"开始转换: {input_path} -> {output_path}")
            
            # 打开并读取docx文件
            with zipfile.ZipFile(input_path, 'r') as docx_zip:
                # 读取主文档
                if 'word/document.xml' not in docx_zip.namelist():
                    print("错误: 无法找到文档内容")
                    return False
                
                # 读取document.xml
                document_content = docx_zip.read('word/document.xml').decode('utf-8')
                
                # 读取样式（可选）
                styles_content = None
                if 'word/styles.xml' in docx_zip.namelist():
                    styles_content = docx_zip.read('word/styles.xml').decode('utf-8')
            
            # 解析并转换
            markdown_content = self._parse_document(document_content, styles_content)
            
            # 写入输出文件
            with open(output_path, 'w', encoding='utf-8') as f:
                f.write(markdown_content)
            
            print(f"转换成功! 输出文件: {output_path}")
            return True
            
        except Exception as e:
            print(f"转换失败: {str(e)}")
            return False
    
    def _parse_document(self, document_xml: str, styles_xml: Optional[str] = None) -> str:
        """
        解析Word文档XML并转换为Markdown
        
        Args:
            document_xml: 文档XML内容
            styles_xml: 样式XML内容（可选）
            
        Returns:
            str: 转换后的Markdown内容
        """
        # 解析XML
        root = ET.fromstring(document_xml)
        
        # 提取所有段落
        body = root.find('w:body', self.namespaces)
        if body is None:
            return ""
        
        markdown_lines = []
        
        # 遍历所有元素
        for element in body:
            if element.tag.endswith('p'):
                # 处理段落
                paragraph_text = self._parse_paragraph(element)
                if paragraph_text:
                    markdown_lines.append(paragraph_text)
            elif element.tag.endswith('tbl'):
                # 处理表格
                table_text = self._parse_table(element)
                if table_text:
                    markdown_lines.append(table_text)
            elif element.tag.endswith('sectPr'):
                # 处理节属性（忽略）
                pass
        
        # 合并为最终内容
        return '\n\n'.join(markdown_lines)
    
    def _parse_paragraph(self, paragraph_elem) -> str:
        """
        解析段落元素
        
        Args:
            paragraph_elem: 段落XML元素
            
        Returns:
            str: 转换后的段落文本
        """
        text_parts = []
        
        # 遍历所有文本运行
        for run in paragraph_elem.findall('.//w:r', self.namespaces):
            run_text = self._parse_run(run)
            if run_text:
                text_parts.append(run_text)
        
        # 合并文本
        paragraph_text = ''.join(text_parts)
        
        # 检查是否为标题
        heading_level = self._get_heading_level(paragraph_elem)
        if heading_level > 0:
            return f"{'#' * heading_level} {paragraph_text.strip()}"
        
        # 检查是否为列表项
        list_info = self._get_list_info(paragraph_elem)
        if list_info:
            level, is_ordered = list_info
            indent = '  ' * (level - 1)
            prefix = '1.' if is_ordered else '*'
            return f"{indent}{prefix} {paragraph_text.strip()}"
        
        return paragraph_text.strip()
    
    def _parse_run(self, run_elem) -> str:
        """
        解析文本运行元素
        
        Args:
            run_elem: 文本运行XML元素
            
        Returns:
            str: 转换后的文本
        """
        text_parts = []
        
        # 遍历所有文本节点
        for text_elem in run_elem.findall('.//w:t', self.namespaces):
            if text_elem.text:
                text_parts.append(text_elem.text)
        
        # 处理制表符
        tab_elems = run_elem.findall('.//w:tab', self.namespaces)
        if tab_elems:
            text_parts.append('\t' * len(tab_elems))
        
        # 处理换行
        br_elems = run_elem.findall('.//w:br', self.namespaces)
        if br_elems:
            text_parts.append('\n' * len(br_elems))
        
        # 合并文本
        text = ''.join(text_parts)
        
        # 应用格式
        text = self._apply_formatting(run_elem, text)
        
        return text
    
    def _apply_formatting(self, run_elem, text: str) -> str:
        """
        应用文本格式
        
        Args:
            run_elem: 文本运行XML元素
            text: 原始文本
            
        Returns:
            str: 应用格式后的文本
        """
        if not text:
            return text
        
        # 检查格式
        r_pr = run_elem.find('w:rPr', self.namespaces)
        if not r_pr:
            return text
        
        # 粗体
        if r_pr.find('w:b', self.namespaces) is not None:
            text = f'**{text}**'
        
        # 斜体
        if r_pr.find('w:i', self.namespaces) is not None:
            text = f'*{text}*'
        
        # 删除线
        if r_pr.find('w:strike', self.namespaces) is not None:
            text = f'~~{text}~~'
        
        # 下划线
        if r_pr.find('w:u', self.namespaces) is not None:
            text = f'<u>{text}</u>'
        
        return text
    
    def _get_heading_level(self, paragraph_elem) -> int:
        """
        获取段落的标题级别
        
        Args:
            paragraph_elem: 段落XML元素
            
        Returns:
            int: 标题级别（0表示不是标题）
        """
        p_pr = paragraph_elem.find('w:pPr', self.namespaces)
        if not p_pr:
            return 0
        
        p_style = p_pr.find('w:pStyle', self.namespaces)
        if not p_style or 'w:val' not in p_style.attrib:
            return 0
        
        style_name = p_style.attrib['w:val']
        if style_name.startswith('Heading'):
            try:
                level = int(style_name.replace('Heading', ''))
                if 1 <= level <= 6:
                    return level
            except:
                pass
        
        return 0
    
    def _get_list_info(self, paragraph_elem) -> Optional[tuple]:
        """
        获取段落的列表信息
        
        Args:
            paragraph_elem: 段落XML元素
            
        Returns:
            tuple: (级别, 是否有序) 或 None
        """
        p_pr = paragraph_elem.find('w:pPr', self.namespaces)
        if not p_pr:
            return None
        
        num_pr = p_pr.find('w:numPr', self.namespaces)
        if not num_pr:
            return None
        
        # 获取级别
        ilvl = num_pr.find('w:ilvl', self.namespaces)
        level = 1
        if ilvl and 'w:val' in ilvl.attrib:
            try:
                level = int(ilvl.attrib['w:val']) + 1  # 从1开始
            except:
                pass
        
        # 简单判断是否有序（实际需要查看numId对应的样式）
        is_ordered = False
        num_id = num_pr.find('w:numId', self.namespaces)
        if num_id and 'w:val' in num_id.attrib:
            # 简单假设：奇数numId为有序，偶数为无序
            try:
                num_val = int(num_id.attrib['w:val'])
                is_ordered = num_val % 2 == 1
            except:
                pass
        
        return (level, is_ordered)
    
    def _parse_table(self, table_elem) -> str:
        """
        解析表格元素
        
        Args:
            table_elem: 表格XML元素
            
        Returns:
            str: 转换后的表格文本
        """
        rows = []
        
        # 提取所有行
        for tr in table_elem.findall('.//w:tr', self.namespaces):
            row_cells = []
            
            # 提取所有单元格
            for tc in tr.findall('.//w:tc', self.namespaces):
                cell_text = []
                
                # 提取单元格内的所有段落
                for p in tc.findall('.//w:p', self.namespaces):
                    p_text = self._parse_paragraph(p)
                    if p_text:
                        cell_text.append(p_text)
                
                row_cells.append(' '.join(cell_text))
            
            if row_cells:
                rows.append(row_cells)
        
        if not rows:
            return ""
        
        # 生成Markdown表格
        markdown_table = []
        
        # 表头
        header = rows[0]
        markdown_table.append('| ' + ' | '.join(header) + ' |')
        
        # 分隔行
        separator = '| ' + ' | '.join(['---'] * len(header)) + ' |'
        markdown_table.append(separator)
        
        # 数据行
        for row in rows[1:]:
            if len(row) == len(header):
                markdown_table.append('| ' + ' | '.join(row) + ' |')
            else:
                # 列数不匹配时补齐
                padded_row = row + [''] * (len(header) - len(row))
                markdown_table.append('| ' + ' | '.join(padded_row) + ' |')
        
        return '\n'.join(markdown_table)

def main():
    """
    主函数
    """
    if len(sys.argv) < 3:
        print("用法: python true_docx_converter.py <input.docx> <output.md>")
        return 1
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    converter = TrueDocxConverter()
    success = converter.convert(input_file, output_file)
    
    return 0 if success else 1

if __name__ == "__main__":
    sys.exit(main())
