#!/usr/bin/env python3
"""
工具函数模块

提供智能体嵌入包的工具函数
"""

import logging
import os

def setup_logger(agent_name, log_level=logging.INFO):
    """
    设置日志
    
    Args:
        agent_name: 智能体名称
        log_level: 日志级别
    
    Returns:
        日志记录器
    """
    # 创建日志目录
    log_dir = os.path.join(os.path.dirname(__file__), "logs")
    os.makedirs(log_dir, exist_ok=True)
    
    # 创建日志记录器
    logger = logging.getLogger(agent_name)
    logger.setLevel(log_level)
    
    # 避免重复添加处理器
    if not logger.handlers:
        # 创建文件处理器
        log_file = os.path.join(log_dir, f"{agent_name}.log")
        file_handler = logging.FileHandler(log_file, encoding="utf-8")
        file_handler.setLevel(log_level)
        
        # 创建控制台处理器
        console_handler = logging.StreamHandler()
        console_handler.setLevel(log_level)
        
        # 设置日志格式
        formatter = logging.Formatter(
            '%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )
        file_handler.setFormatter(formatter)
        console_handler.setFormatter(formatter)
        
        # 添加处理器
        logger.addHandler(file_handler)
        logger.addHandler(console_handler)
    
    return logger

def validate_config(config):
    """
    验证配置
    
    Args:
        config: 配置字典
    
    Returns:
        验证结果
    """
    required_keys = ["agent_name"]
    
    for key in required_keys:
        if key not in config:
            return False, f"缺少必要配置项: {key}"
    
    return True, "配置验证通过"

def ensure_directory(path):
    """
    确保目录存在
    
    Args:
        path: 目录路径
    """
    os.makedirs(path, exist_ok=True)

def get_project_root():
    """
    获取项目根目录
    
    Returns:
        项目根目录路径
    """
    current_dir = os.path.dirname(__file__)
    return os.path.abspath(os.path.join(current_dir, ".."))

def format_timestamp():
    """
    格式化时间戳
    
    Returns:
        格式化的时间戳字符串
    """
    import datetime
    return datetime.datetime.now().strftime("%Y%m%d_%H%M%S")
