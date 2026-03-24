#!/usr/bin/env python3
"""
使用site-packages路径构建ELR-Desktop-Assistant
"""

import subprocess
import os
import sys

# 添加site-packages到Python路径
site_packages = "C:\\Users\\Administrator\\AppData\\Local\\Packages\\PythonSoftwareFoundation.Python.3.13_qbz5n2kfra8p0\\LocalCache\\local-packages\\Python313\\site-packages"
sys.path.insert(0, site_packages)

# 构建命令
cmd = [
    'python',
    '-c',
    'import pyinstaller; pyinstaller.__main__.run(["--name", "ELRDesktopAssistant", "--windowed", "--icon", "icons/elr_icon.png", "--onefile", "elr_desktop_assistant.py"])'
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
