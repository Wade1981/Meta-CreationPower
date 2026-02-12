#!/usr/bin/env python3
"""
Meta-CreationPower 项目安装脚本
基于《元创力》元协议 α-0.1 版
"""

from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as f:
    long_description = f.read()

setup(
    name="meta-creationpower",
    version="0.1.0",
    description="基于《元创力》元协议的碳硅协同实践平台",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/X54/Meta-CreationPower",
    author="X54先生 & 启蒙灯塔团队",
    author_email="",
    license="MIT",
    packages=find_packages(),
    python_requires=">=3.8",
    install_requires=[],
    extras_require={
        "dev": [
            "pytest>=7.0.0",
            "black>=23.0.0",
            "flake8>=6.0.0",
            "mypy>=1.0.0",
        ],
    },
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Topic :: Software Development :: Libraries :: Application Frameworks",
        "Topic :: Software Development :: Human Machine Interfaces",
    ],
)
