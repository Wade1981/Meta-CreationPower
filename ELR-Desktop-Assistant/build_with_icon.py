#!/usr/bin/env python3
"""
使用PyInstaller构建ELR-Desktop-Assistant并设置图标
"""

import subprocess
import os

# 构建命令
cmd = [
    'python', '-m', 'pyinstaller',
    '--name', 'ELRDesktopAssistant',
    '--windowed',
    '--icon', 'icons/elr_icon.png',
    '--onefile',
    'elr_desktop_assistant.py'
]

print(f"构建命令: {' '.join(cmd)}")

# 运行构建
try:
    result = subprocess.run(cmd, cwd=os.path.dirname(os.path.abspath(__file__)), capture_output=True, text=True)
    print(f"构建输出: {result.stdout}")
    if result.stderr:
        print(f"构建错误: {result.stderr}")
    print(f"构建返回码: {result.returncode}")
    if result.returncode == 0:
        print("构建成功！")
    else:
        print("构建失败！")
except Exception as e:
    print(f"构建过程中出错: {e}")
