#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
测试文档转换工具

此脚本用于测试doc_to_md_converter.py工具，
特别是转换docs目录中的"程序员称谓衍进数智工程师.docx"文件。
"""

import os
import sys
from src.utils.doc_to_md_converter import convert_to_md


def test_conversion():
    """
    测试文档转换功能
    """
    # 定义测试文件路径
    docx_path = os.path.join('docs', '程序员称谓衍进数智工程师.docx')
    md_path = os.path.join('docs', '程序员称谓衍进数智工程师.md')
    
    # 检查测试文件是否存在
    if not os.path.exists(docx_path):
        print(f"错误: 测试文件不存在 - {docx_path}")
        return False
    
    # 转换文件
    print(f"开始转换: {docx_path}")
    success = convert_to_md(docx_path, md_path)
    
    if success:
        # 检查输出文件是否存在
        if os.path.exists(md_path):
            print(f"成功: 转换完成，输出文件 - {md_path}")
            # 显示转换后的文件大小
            md_size = os.path.getsize(md_path)
            print(f"输出文件大小: {md_size} 字节")
            return True
        else:
            print(f"错误: 转换后文件不存在 - {md_path}")
            return False
    else:
        print("错误: 转换失败")
        return False


def main():
    """
    主函数
    """
    print("测试文档转换工具")
    print("=" * 50)
    
    # 执行测试
    success = test_conversion()
    
    print("=" * 50)
    if success:
        print("测试通过: 文档转换工具正常工作")
        return 0
    else:
        print("测试失败: 文档转换工具存在问题")
        return 1


if __name__ == "__main__":
    sys.exit(main())
