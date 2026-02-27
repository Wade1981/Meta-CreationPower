"""
 Word转Markdown转换器

 版本: 1.0

 作者: 启蒙灯塔起源团 - 代码织梦者参考

 描述: 纯Python实现，不依赖任何外部包，支持.docx/.doc格式转换

 """

import zipfile
import xml.etree.ElementTree as ET
import re
import os
import struct
from io import BytesIO
from typing import Dict, List, Tuple, Optional, Union


class WordToMarkdownConverter:
    """
    Word转Markdown转换器

    支持功能：
    1. 标题识别 (H1-H6)
    2. 文本格式化 (粗体、斜体、下划线、删除线)
    3. 列表处理 (有序和无序)
    4. 链接和图片提取
    5. 表格转换
    6. 代码块处理

    """

    # Word XML命名空间
    NAMESPACES = {
        'w': 'http://schemas.openxmlformats.org/wordprocessingml/2006/main',
        'wp': 'http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing',
        'a': 'http://schemas.openxmlformats.org/drawingml/2006/main',
        'pic': 'http://schemas.openxmlformats.org/drawingml/2006/picture',
        'r': 'http://schemas.openxmlformats.org/officeDocument/2006/relationships',
        'v': 'urn:schemas-microsoft-com:vml'
    }

    def __init__(self, word_path: str):
        """
        初始化转换器

        Args:
            word_path: Word文件路径(.docx或.doc)

        """
        self.word_path = word_path
        self.file_ext = os.path.splitext(word_path)[1].lower()
        self.markdown_lines = []
        self.relationships = {}  # 存储文档中的关系映射
        self.image_counter = 0   # 图片计数器
        self.list_styles = {}    # 列表样式
        self.current_list_level = 0
        self.in_table = False
        self.table_rows = []
        self.current_row = []

    def convert(self) -> str:
        """
        执行转换

        Returns:
            转换后的Markdown字符串

        """
        if self.file_ext == '.docx':
            return self._convert_docx()
        elif self.file_ext == '.doc':
            return self._convert_doc()
        else:
            raise ValueError(f"不支持的文件格式: {self.file_ext}")

    def _convert_docx(self) -> str:
        """转换.docx文件"""
        try:
            with zipfile.ZipFile(self.word_path, 'r') as docx_zip:
                # 1. 读取主文档
                document_content = docx_zip.read('word/document.xml').decode('utf-8')

                # 2. 读取样式（用于标题识别）
                styles_content = None
                if 'word/styles.xml' in docx_zip.namelist():
                    styles_content = docx_zip.read('word/styles.xml').decode('utf-8')

                # 3. 读取关系（用于图片和链接）
                relationships_content = None
                if 'word/_rels/document.xml.rels' in docx_zip.namelist():
                    rels_content = docx_zip.read('word/_rels/document.xml.rels').decode('utf-8')
                    self.relationships = self._parse_relationships(rels_content)

                # 4. 解析并转换
                self._parse_docx_xml(document_content, styles_content)

        except Exception as e:
            raise Exception(f"读取.docx文件失败: {str(e)}")

        return '\n'.join(self.markdown_lines)

    def _convert_doc(self) -> str:
        """转换.doc文件（二进制格式，简化处理）"""
        try:
            with open(self.word_path, 'rb') as f:
                content = f.read()

            # .doc格式复杂，这里提供简化版的文本提取
            # 注意：这是基础版本，复杂格式可能无法完美转换
            text_content = self._extract_text_from_doc(content)

            # 基础格式转换
            lines = text_content.split('\n')
            md_lines = []
            in_list = False
            list_indent = 0

            for line in lines:
                line = line.strip()
                if not line:
                    if in_list:
                        md_lines.append('')
                    continue

                # 简单标题检测（基于行长度和结尾标点）
                if len(line) < 100 and line[-1] not in '.。!！?？':
                    if line.startswith('第') and ('章' in line or '节' in line):
                        md_lines.append(f'# {line}')
                        continue

                # 列表检测
                if re.match(r'^[\d•·\-*]', line.lstrip()):
                    if not in_list:
                        md_lines.append('')
                    in_list = True
                    # 简单缩进处理
                    indent = len(line) - len(line.lstrip())
                    prefix = '  ' * (indent // 4)

                    if re.match(r'^\d+[\.、]', line.lstrip()):
                        md_lines.append(f'{prefix}1. {line.lstrip()[2:]}')
                    else:
                        char = line.lstrip()[0]
                        md_lines.append(f'{prefix}* {line.lstrip()[1:]}')
                else:
                    if in_list:
                        md_lines.append('')
                        in_list = False
                    md_lines.append(line)

            return '\n'.join(md_lines)

        except Exception as e:
            raise Exception(f"读取.doc文件失败: {str(e)}")

    def _parse_relationships(self, rels_xml: str) -> Dict[str, str]:
        """解析关系XML"""
        relationships = {}
        try:
            root = ET.fromstring(rels_xml)
            for rel in root.findall('.//{http://schemas.openxmlformats.org/package/2006/relationships}Relationship'):
                rel_id = rel.get('Id')
                rel_type = rel.get('Type')
                target = rel.get('Target')
                if rel_id and target:
                    relationships[rel_id] = target
        except:
            pass
        return relationships

    def _parse_docx_xml(self, document_xml: str, styles_xml: Optional[str] = None):
        """解析Word XML文档"""
        # 解析样式
        style_mapping = self._parse_styles(styles_xml) if styles_xml else {}

        # 解析主文档
        root = ET.fromstring(document_xml)

        # 注册命名空间
        for prefix, uri in self.NAMESPACES.items():
            ET.register_namespace(prefix, uri)

        # 提取所有段落
        body = root.find('w:body', self.NAMESPACES)
        if body is None:
            return

        for element in body:
            # 处理段落
            if element.tag.endswith('p'):
                self._process_paragraph(element, style_mapping)
            # 处理表格
            elif element.tag.endswith('tbl'):
                self._process_table(element)

    def _parse_styles(self, styles_xml: str) -> Dict[str, str]:
        """解析样式XML，建立样式ID到名称的映射"""
        style_mapping = {}
        try:
            root = ET.fromstring(styles_xml)

            for style in root.findall('.//w:style', self.NAMESPACES):
                style_id = style.get('{http://schemas.openxmlformats.org/wordprocessingml/2006/main}styleId')
                style_name_elem = style.find('w:name', self.NAMESPACES)
                if style_name_elem is not None:
                    style_name = style_name_elem.get('{http://schemas.openxmlformats.org/wordprocessingml/2006/main}val')
                    if style_id and style_name:
                        style_mapping[style_id] = style_name

        except:
            pass

        return style_mapping

    def _process_paragraph(self, p_elem, style_mapping: Dict[str, str]):
        """处理段落元素"""
        # 检查是否在列表中
        num_pr = p_elem.find('.//w:numPr', self.NAMESPACES)

        if num_pr is not None:
            # 处理列表项
            self._process_list_item(p_elem, num_pr)
            return

        # 检查是否为标题
        p_style = p_elem.find('w:pPr/w:pStyle', self.NAMESPACES)
        if p_style is not None:
            style_id = p_style.get('{http://schemas.openxmlformats.org/wordprocessingml/2006/main}val')
            if style_id and style_id in style_mapping:
                style_name = style_mapping[style_id]
                # 判断是否为标题样式
                if style_name.startswith('Heading'):
                    try:
                        level = int(style_name.replace('Heading', '').strip())
                        if 1 <= level <= 6:
                            text = self._extract_text_from_paragraph(p_elem)
                            if text:
                                self.markdown_lines.append(f"{'#' * level} {text}")
                                return
                    except:
                        pass

        # 普通段落
        text = self._extract_text_from_paragraph(p_elem)
        if text:
            # 检查是否为代码块（简单检测）
            if self._looks_like_code(text):
                self.markdown_lines.append(f'```\n{text}\n```')
            else:
                self.markdown_lines.append(text)

    def _extract_text_from_paragraph(self, p_elem) -> str:
        """从段落元素中提取文本（包含格式）"""
        text_parts = []

        # 遍历所有文本运行
        for r_elem in p_elem.findall('.//w:r', self.NAMESPACES):
            r_text = self._extract_text_from_run(r_elem)
            if r_text:
                # 检查运行格式
                formatted_text = self._apply_format_to_run(r_elem, r_text)
                text_parts.append(formatted_text)

        # 合并文本
        result = ''.join(text_parts)

        # 处理超链接
        result = self._process_hyperlinks(p_elem, result)

        return result.strip()

    def _extract_text_from_run(self, r_elem) -> str:
        """从运行元素中提取文本"""
        text = ''
        t_elems = r_elem.findall('.//w:t', self.NAMESPACES)

        for t_elem in t_elems:
            if t_elem.text:
                text += t_elem.text

        # 处理制表符
        tab_elems = r_elem.findall('.//w:tab', self.NAMESPACES)
        if tab_elems:
            text += '\t' * len(tab_elems)

        # 处理换行
        br_elems = r_elem.findall('.//w:br', self.NAMESPACES)
        if br_elems:
            text += '\n' * len(br_elems)

        return text

    def _apply_format_to_run(self, r_elem, text: str) -> str:
        """应用格式到文本"""
        if not text:
            return text

        r_pr = r_elem.find('w:rPr', self.NAMESPACES)
        if r_pr is None:
            return text

        # 粗体
        b_elem = r_pr.find('w:b', self.NAMESPACES)
        b_cs_elem = r_pr.find('w:bCs', self.NAMESPACES)
        if b_elem is not None or b_cs_elem is not None:
            text = f'**{text}**'

        # 斜体
        i_elem = r_pr.find('w:i', self.NAMESPACES)
        i_cs_elem = r_pr.find('w:iCs', self.NAMESPACES)
        if i_elem is not None or i_cs_elem is not None:
            text = f'*{text}*'

        # 下划线
        u_elem = r_pr.find('w:u', self.NAMESPACES)
        if u_elem is not None:
            text = f'<u>{text}</u>'  # Markdown无标准下划线语法

        # 删除线
        strike_elem = r_pr.find('w:strike', self.NAMESPACES)
        del_elem = r_pr.find('w:del', self.NAMESPACES)
        if strike_elem is not None or del_elem is not None:
            text = f'~~{text}~~'

        return text

    def _process_hyperlinks(self, p_elem, text: str) -> str:
        """处理超链接"""
        # 查找段落中的所有超链接
        links = p_elem.findall('.//w:hyperlink', self.NAMESPACES)

        for link in links:
            link_id = link.get('{http://schemas.openxmlformats.org/officeDocument/2006/relationships}id')
            if link_id in self.relationships:
                url = self.relationships[link_id]
                # 提取链接文本
                link_text_elem = link.find('.//w:r/w:t', self.NAMESPACES)
                if link_text_elem is not None and link_text_elem.text:
                    link_text = link_text_elem.text
                    # 替换文本中的链接
                    if link_text in text:
                        text = text.replace(link_text, f'[{link_text}]({url})')

        return text

    def _process_list_item(self, p_elem, num_pr):
        """处理列表项"""
        # 提取列表文本
        text = self._extract_text_from_paragraph(p_elem)
        if not text:
            return

        # 确定列表级别和类型
        ilvl_elem = num_pr.find('w:ilvl', self.NAMESPACES)
        num_id_elem = num_pr.find('w:numId', self.NAMESPACES)

        level = 0
        if ilvl_elem is not None:
            level_attr = ilvl_elem.get('{http://schemas.openxmlformats.org/wordprocessingml/2006/main}val')
            if level_attr:
                try:
                    level = int(level_attr)
                except:
                    level = 0

        # 简单判断有序/无序列表
        is_ordered = False
        if num_id_elem is not None:
            # 这里可以更复杂地判断，简化处理
            num_id = num_id_elem.get('{http://schemas.openxmlformats.org/wordprocessingml/2006/main}val')
            if num_id and num_id.isdigit():
                is_ordered = int(num_id) % 2 == 1  # 简单假设

        # 生成Markdown列表项
        indent = '  ' * level
        prefix = f'1.' if is_ordered else '*'

        self.markdown_lines.append(f'{indent}{prefix} {text}')

    def _process_table(self, tbl_elem):
        """处理表格"""
        rows = []

        # 提取所有行
        for tr_elem in tbl_elem.findall('.//w:tr', self.NAMESPACES):
            row_cells = []

            # 提取所有单元格
            for tc_elem in tr_elem.findall('.//w:tc', self.NAMESPACES):
                cell_text = []

                # 提取单元格内的所有段落文本
                for p_elem in tc_elem.findall('.//w:p', self.NAMESPACES):
                    p_text = self._extract_text_from_paragraph(p_elem)
                    if p_text:
                        cell_text.append(p_text)

                row_cells.append(' '.join(cell_text))

            if row_cells:
                rows.append(row_cells)

        if not rows:
            return

        # 生成Markdown表格
        self.markdown_lines.append('')  # 空行

        # 表头
        header = rows[0]
        self.markdown_lines.append('| ' + ' | '.join(header) + ' |')

        # 分隔行
        separator = '| ' + ' | '.join(['---'] * len(header)) + ' |'
        self.markdown_lines.append(separator)

        # 数据行
        for row in rows[1:]:
            if len(row) == len(header):
                self.markdown_lines.append('| ' + ' | '.join(row) + ' |')
            else:
                # 如果列数不匹配，补齐
                padded_row = row + [''] * (len(header) - len(row))
                self.markdown_lines.append('| ' + ' | '.join(padded_row) + ' |')

        self.markdown_lines.append('')  # 空行

    def _extract_text_from_doc(self, content: bytes) -> str:
        """从.doc文件中提取文本（简化版）"""
        text = ''

        # .doc文件是OLE复合文档，这里使用简单的文本提取
        # 注意：这是基础版本，复杂格式会丢失

        # 尝试查找文本片段
        try:
            # 将二进制转换为字符串（忽略非文本部分）
            content_str = content.decode('utf-8', errors='ignore')

            # 提取看起来像文本的部分
            lines = content_str.split('\x00')
            for line in lines:
                # 过滤控制字符和不可打印字符
                clean_line = ''.join(c for c in line if 32 <= ord(c) < 127 or c in '\n\r\t')
                if clean_line.strip():
                    text += clean_line + '\n'
        except:
            # 如果UTF-8解码失败，尝试其他编码
            try:
                content_str = content.decode('gbk', errors='ignore')
                lines = content_str.split('\x00')
                for line in lines:
                    clean_line = ''.join(c for c in line if 32 <= ord(c) < 127 or c in '\n\r\t')
                    if clean_line.strip():
                        text += clean_line + '\n'
            except:
                text = "无法解析.doc文件内容，请考虑转换为.docx格式"

        return text

    def _looks_like_code(self, text: str) -> bool:
        """简单检测文本是否像代码"""
        # 检测代码的常见特征
        code_indicators = [
            ('{', 0.1),  # 花括号
            ('}', 0.1),
            (';', 0.1),  # 分号
            ('=', 0.1),  # 等号
            ('(', 0.1),  # 括号
            (')', 0.1),
            ('def ', 0.3),  # 函数定义
            ('function', 0.3),
            ('class ', 0.3),  # 类定义
            ('import ', 0.3),  # 导入
            ('var ', 0.2),  # 变量声明
            ('let ', 0.2),
            ('const ', 0.2),
        ]

        score = 0
        text_lower = text.lower()

        for indicator, weight in code_indicators:
            if indicator in text or indicator in text_lower:
                score += weight

        # 检查缩进（行以空格开头）
        lines = text.split('\n')
        indented_lines = sum(1 for line in lines if line.startswith('    ') or line.startswith('\t'))
        if len(lines) > 0:
            indent_ratio = indented_lines / len(lines)
            score += indent_ratio * 0.5

        return score > 0.5


# 使用示例和工具函数
def convert_word_to_markdown(word_path: str, output_path: Optional[str] = None) -> str:
    """
    转换Word文件为Markdown

    Args:
        word_path: Word文件路径
        output_path: 可选的输出文件路径，如果为None则只返回字符串

    Returns:
        Markdown字符串

    """
    try:
        converter = WordToMarkdownConverter(word_path)
        markdown_content = converter.convert()

        if output_path:
            with open(output_path, 'w', encoding='utf-8') as f:
                f.write(markdown_content)
            print(f"转换完成！已保存到: {output_path}")

        return markdown_content

    except Exception as e:
        error_msg = f"转换失败: {str(e)}"
        print(error_msg)
        return error_msg


def batch_convert(word_dir: str, output_dir: str, pattern: str = "*.docx"):
    """
    批量转换Word文件为Markdown

    Args:
        word_dir: Word文件所在目录
        output_dir: 输出目录
        pattern: 文件匹配模式

    """
    import glob

    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    word_files = glob.glob(os.path.join(word_dir, pattern))

    for word_file in word_files:
        try:
            filename = os.path.basename(word_file)
            output_filename = os.path.splitext(filename)[0] + '.md'
            output_path = os.path.join(output_dir, output_filename)

            print(f"正在转换: {filename}")
            convert_word_to_markdown(word_file, output_path)

        except Exception as e:
            print(f"转换 {word_file} 失败: {str(e)}")


# 命令行接口
if __name__ == "__main__":
    import sys
    import argparse

    parser = argparse.ArgumentParser(description='Word转Markdown转换器')
    parser.add_argument('input', help='输入Word文件路径或目录')
    parser.add_argument('-o', '--output', help='输出Markdown文件路径或目录')
    parser.add_argument('-b', '--batch', action='store_true', help='批量转换模式')
    parser.add_argument('-p', '--pattern', default="*.docx", help='批量转换时的文件匹配模式')

    args = parser.parse_args()

    if args.batch:
        if not args.output:
            print("错误：批量转换需要指定输出目录")
            sys.exit(1)

        batch_convert(args.input, args.output, args.pattern)
    else:
        output_path = args.output
        if not output_path:
            # 默认输出到同名.md文件
            base_name = os.path.splitext(args.input)[0]
            output_path = base_name + '.md'

        convert_word_to_markdown(args.input, output_path)


# 提供给代码织梦者的参考实现说明
"""
===============================================================
Word转Markdown转换器 - 代码织梦者参考实现
===============================================================

设计理念：
1. 零依赖：仅使用Python标准库，无需安装任何第三方包
2. 兼容性：支持.docx和基础.doc格式
3. 可扩展：模块化设计，易于添加新功能
4. 容错性：良好的错误处理和异常管理

核心组件：
1. WordToMarkdownConverter: 主转换类
2. XML解析：处理.docx的XML结构
3. 样式映射：识别标题、列表等格式
4. 格式转换：将Word格式转换为Markdown语法
5. 表格处理：转换表格为Markdown表格

使用示例：
1. 单个文件转换：
   python word_to_md.py input.docx -o output.md
    
2. 批量转换：
   python word_to_md.py ./docs -b -o ./markdown

扩展建议：
1. 图片提取：从.docx的media文件夹提取图片
2. 样式映射：更精确的样式识别
3. 复杂表格：支持合并单元格
4. 批注处理：提取Word批注
5. 目录生成：自动生成文档目录

注意事项：
1. .doc格式支持有限，建议先转为.docx
2. 复杂格式可能无法完美转换
3. 某些特殊字符可能需要额外处理
4. 性能考虑：大文件可能较慢

代码结构说明：
- _convert_docx: 处理.docx文件
- _convert_doc: 处理.doc文件（基础版）
- _parse_docx_xml: 解析Word XML
- _process_paragraph: 处理段落
- _process_table: 处理表格
- _apply_format_to_run: 应用文本格式

此代码为"启蒙灯塔起源团"知识库贡献，
遵循"和清寂静"原则设计，力求清晰、简洁、实用。
"""
