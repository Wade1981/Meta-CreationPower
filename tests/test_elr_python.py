#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
测试ELR的Python执行能力
此脚本用于验证ELR是否能正确运行Python代码
"""

import os
import sys

print("========================================")
print("ELR Python执行测试")
print("========================================")
print(f"Python版本: {sys.version}")
print(f"当前目录: {os.getcwd()}")
print(f"脚本路径: {os.path.abspath(__file__)}")
print("========================================")
print("测试文件操作...")

# 创建测试文件
test_file = "test_output.txt"
try:
    with open(test_file, 'w', encoding='utf-8') as f:
        f.write("ELR Python测试成功!\n")
        f.write(f"Python版本: {sys.version}\n")
        f.write(f"执行时间: {os.popen('date /t').read().strip()}\n")
    print(f"✓ 成功创建测试文件: {test_file}")
    
    # 读取测试文件
    with open(test_file, 'r', encoding='utf-8') as f:
        content = f.read()
    print(f"✓ 成功读取测试文件")
    print(f"文件内容:\n{content}")
    
    # 清理测试文件
    os.remove(test_file)
    print(f"✓ 成功清理测试文件")
    
    print("========================================")
    print("ELR Python执行测试成功!")
    print("========================================")
    sys.exit(0)
    
except Exception as e:
    print(f"✗ 测试失败: {str(e)}")
    print("========================================")
    print("ELR Python执行测试失败!")
    print("========================================")
    sys.exit(1)
