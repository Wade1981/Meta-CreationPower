# 模型基础类

import numpy as np
from abc import ABC, abstractmethod

class BaseModel(ABC):
    """模型基础类"""
    
    def __init__(self, config):
        """初始化模型
        
        Args:
            config: 配置对象
        """
        self.config = config
    
    @abstractmethod
    def predict(self, data):
        """预测方法
        
        Args:
            data: 输入数据
            
        Returns:
            预测结果
        """
        pass
    
    @abstractmethod
    def train(self, X, y):
        """训练方法
        
        Args:
            X: 特征数据
            y: 标签数据
        """
        pass
    
    def save(self, path):
        """保存模型
        
        Args:
            path: 保存路径
        """
        pass
    
    def load(self, path):
        """加载模型
        
        Args:
            path: 加载路径
        """
        pass
