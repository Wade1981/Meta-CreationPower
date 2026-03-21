# 财务健康熵值评估模型

import numpy as np
from .base_model import BaseModel

class HealthEntropyModel(BaseModel):
    """财务健康熵值评估模型"""
    
    def __init__(self, config):
        """初始化模型
        
        Args:
            config: 配置对象
        """
        super().__init__(config)
        self.threshold_low = config.HEALTH_ENTROPY_THRESHOLD_LOW
        self.threshold_high = config.HEALTH_ENTROPY_THRESHOLD_HIGH
        self.adjustment_factor = 1.0  # 财务场景适配调整系数
    
    def calculate_health_entropy(self, financial_indicators):
        """计算财务健康熵值
        
        Args:
            financial_indicators: 核心财务指标
            
        Returns:
            health_entropy: 财务健康熵值
        """
        # 提取核心财务指标
        indicators = []
        for key, value in financial_indicators.items():
            if isinstance(value, (int, float)):
                indicators.append(value)
        
        if not indicators:
            return 0
        
        # 计算概率分布（简化处理，实际应基于历史数据和行业阈值）
        values = np.array(indicators)
        # 归一化处理
        min_val = np.min(values)
        max_val = np.max(values)
        if max_val > min_val:
            normalized_values = (values - min_val) / (max_val - min_val)
        else:
            normalized_values = np.zeros_like(values)
        
        # 计算概率
        probabilities = normalized_values / np.sum(normalized_values) if np.sum(normalized_values) > 0 else np.ones_like(normalized_values) / len(normalized_values)
        
        # 计算熵值
        entropy = -np.sum(probabilities * np.log(probabilities + 1e-10))
        health_entropy = self.adjustment_factor * entropy
        
        return health_entropy
    
    def get_health_level(self, health_entropy):
        """获取健康等级
        
        Args:
            health_entropy: 财务健康熵值
            
        Returns:
            level: 健康等级
            description: 等级描述
        """
        if health_entropy <= self.threshold_low:
            return "低熵区", "财务健康度优秀，系统高度有序，风险完全可控"
        elif health_entropy <= self.threshold_high:
            return "中熵区", "财务健康度良好，系统基本有序，局部风险需关注"
        else:
            return "高熵区", "财务健康度红色预警，系统无序度高，风险快速累积，需立即干预"
    
    def predict(self, data):
        """预测财务健康状态
        
        Args:
            data: 包含财务指标的数据
            
        Returns:
            预测结果
        """
        financial_indicators = data.get('financial_indicators', {})
        health_entropy = self.calculate_health_entropy(financial_indicators)
        level, description = self.get_health_level(health_entropy)
        
        return {
            'health_entropy': health_entropy,
            'health_level': level,
            'description': description,
            'financial_indicators': financial_indicators
        }
    
    def train(self, X, y):
        """训练模型（此模型为规则模型，无需训练）
        
        Args:
            X: 特征数据
            y: 标签数据
        """
        pass
