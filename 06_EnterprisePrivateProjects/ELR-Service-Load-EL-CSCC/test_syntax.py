#!/usr/bin/env python3
"""
Test script to check syntax of all Python files in the project.
"""

import os
import sys
import ast


def check_syntax(file_path):
    """
    Check syntax of a Python file.
    
    Args:
        file_path (str): Path to Python file
        
    Returns:
        bool: True if syntax is correct, False otherwise
    """
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        ast.parse(content)
        print(f"✓ {file_path}: Syntax is correct")
        return True
    except SyntaxError as e:
        print(f"✗ {file_path}: Syntax error - {e}")
        return False
    except Exception as e:
        print(f"✗ {file_path}: Error - {e}")
        return False


def main():
    """Check syntax of all Python files in the project."""
    project_root = os.path.dirname(os.path.abspath(__file__))
    python_files = []
    
    # Find all Python files
    for root, dirs, files in os.walk(project_root):
        for file in files:
            if file.endswith('.py'):
                python_files.append(os.path.join(root, file))
    
    # Check syntax of each file
    print(f"Checking syntax of {len(python_files)} Python files...")
    print("=" * 80)
    
    success_count = 0
    for file_path in python_files:
        if check_syntax(file_path):
            success_count += 1
    
    print("=" * 80)
    print(f"Summary: {success_count}/{len(python_files)} files have correct syntax")
    
    if success_count == len(python_files):
        print("All files have correct syntax!")
        return 0
    else:
        print("Some files have syntax errors!")
        return 1


if __name__ == "__main__":
    sys.exit(main())
