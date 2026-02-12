#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Hello World example for Enlightenment Lighthouse Runtime (ELR)
This example demonstrates how to run a Python application in ELR
"""

import os
import sys
import time

# Print welcome message
print("====================================")
print("Hello from Enlightenment Lighthouse Runtime!")
print("====================================")
print(f"Language: Python {sys.version}")
print(f"Runtime: ELR")
print("====================================")

# Demonstrate basic Python features
print("\nBasic Python Features:")

# Variables and output
number = 42
message = "The answer to life, the universe, and everything"
print(f"Integer variable: {number}")
print(f"String variable: {message}")

# Loop
print("\nLoop demonstration:")
for i in range(1, 6):
    print(f"Iteration {i}")

# List
print("\nList demonstration:")
languages = ["C++", "Python", "Java", "JavaScript", "Go"]
for lang in languages:
    print(f"Supported language: {lang}")

# Function
print("\nFunction demonstration:")

def square(x):
    """Calculate square of a number"""
    return x * x

for i in range(1, 4):
    print(f"Square of {i} is {square(i)}")

# Environment variables
print("\nEnvironment variables:")
elr_container_id = os.environ.get("ELR_CONTAINER_ID", "Not set (running outside ELR)")
print(f"ELR_CONTAINER_ID: {elr_container_id}")

# Try to import common libraries
print("\nLibrary demonstration:")

try:
    import numpy as np
    print("NumPy: Available")
    # Simple NumPy example
    arr = np.array([1, 2, 3, 4, 5])
    print(f"NumPy array: {arr}")
    print(f"Array mean: {np.mean(arr)}")
except ImportError:
    print("NumPy: Not available")

try:
    import requests
    print("Requests: Available")
except ImportError:
    print("Requests: Not available")

# End message
print("\n====================================")
print("Python example completed successfully!")
print("====================================")
