#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Lumina Runtime Container - Python示例
"""

import os
import sys
import time

# 打印欢迎信息
print("====================================")
print("Hello from Lumina Runtime Container!")
print("====================================")
print(f"Language: Python {sys.version}")
print(f"Container Version: 1.0.0")
print(f"Current directory: {os.getcwd()}")
print("====================================")

# 演示基本功能
print("\nBasic Python Features:")

# 变量和输出
number = 42
message = "The answer to life, the universe, and everything"
print(f"Integer variable: {number}")
print(f"String variable: {message}")

# 循环
print("\nLoop demonstration:")
for i in range(1, 6):
    print(f"Iteration {i}")

# 列表
print("\nList demonstration:")
languages = ["C++", "Python", "Java", "JavaScript", "Go"]
for lang in languages:
    print(f"Supported language: {lang}")

# 函数
print("\nFunction demonstration:")

def square(x):
    """计算平方"""
    return x * x

for i in range(1, 4):
    print(f"Square of {i} is {square(i)}")

# 尝试导入AI框架
print("\nAI Framework Demonstration:")

try:
    import tensorflow as tf
    print(f"TensorFlow version: {tf.__version__}")
    # 简单的TensorFlow示例
    a = tf.constant([[1.0, 2.0], [3.0, 4.0]])
    b = tf.constant([[5.0, 6.0], [7.0, 8.0]])
    c = tf.matmul(a, b)
    print(f"TensorFlow matrix multiplication result:\n{c.numpy()}")
except ImportError:
    print("TensorFlow not installed")

try:
    import torch
    print(f"PyTorch version: {torch.__version__}")
    # 简单的PyTorch示例
    x = torch.tensor([[1.0, 2.0], [3.0, 4.0]])
    y = torch.tensor([[5.0, 6.0], [7.0, 8.0]])
    z = torch.matmul(x, y)
    print(f"PyTorch matrix multiplication result:\n{z.numpy()}")
except ImportError:
    print("PyTorch not installed")

try:
    import numpy as np
    print(f"NumPy version: {np.__version__}")
    # 简单的NumPy示例
    arr = np.array([1, 2, 3, 4, 5])
    print(f"NumPy array: {arr}")
    print(f"NumPy array mean: {np.mean(arr)}")
except ImportError:
    print("NumPy not installed")

# 结束信息
print("\n====================================")
print("Python example completed successfully!")
print("====================================")
