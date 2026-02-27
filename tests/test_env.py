#!/usr/bin/env python3
"""
测试Python环境和文件编码
"""

import sys
import os

print(f"Python版本: {sys.version}")
print(f"当前目录: {os.getcwd()}")
print(f"Python路径: {sys.path}")

# 测试文件编码
try:
    with open('src/main.py', 'r', encoding='utf-8') as f:
        content = f.read()
    print("✓ 成功读取main.py文件")
    print(f"文件长度: {len(content)} 字符")
except Exception as e:
    print(f"✗ 读取文件失败: {e}")

# 测试简单导入
try:
    import src
    print("✓ 成功导入src模块")
except Exception as e:
    print(f"✗ 导入src模块失败: {e}")

# 测试main模块导入
try:
    from src import main
    print("✓ 成功导入src.main模块")
except Exception as e:
    print(f"✗ 导入src.main模块失败: {e}")
