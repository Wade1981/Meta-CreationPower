# 模型集成模块

import sys
import os

# 添加项目根目录到Python路径
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..')))

from src.models.entropy_balance_model import EntropyBalanceModel
from src.models.five_dimension_entropy_model import FiveDimensionEntropyModel
from src.models.health_entropy_model import HealthEntropyModel
from src.models.rl_entropy_optimization_model import RLEntropyOptimizationModel
from src.models.four_process_entropy_verification_model import FourProcessEntropyVerificationModel

class ModelIntegration:
    """模型集成类"""
    
    def __init__(self, config):
        """初始化模型集成
        
        Args:
            config: 配置对象
        """
        self.config = config
        self.entropy_balance_model = EntropyBalanceModel(config)
        self.five_dimension_model = FiveDimensionEntropyModel(config)
        self.health_entropy_model = HealthEntropyModel(config)
        self.rl_optimization_model = RLEntropyOptimizationModel(config)
        self.four_process_model = FourProcessEntropyVerificationModel(config)
    
    def calculate_entropy_balance(self, internal_entropy, external_entropy):
        """计算熵变平衡
        
        Args:
            internal_entropy: 内部自发熵增
            external_entropy: 外部输入的有效负熵流
            
        Returns:
            熵变平衡结果
        """
        data = {
            'internal_entropy': internal_entropy,
            'external_entropy': external_entropy
        }
        return self.entropy_balance_model.predict(data)
    
    def check_five_dimension_balance(self, dimension_data):
        """检查五维时间线平衡
        
        Args:
            dimension_data: 各维度的数据
            
        Returns:
            五维时间线平衡结果
        """
        data = {
            'dimension_data': dimension_data
        }
        return self.five_dimension_model.predict(data)
    
    def evaluate_health_entropy(self, financial_indicators):
        """评估财务健康熵值
        
        Args:
            financial_indicators: 核心财务指标
            
        Returns:
            财务健康熵值评估结果
        """
        data = {
            'financial_indicators': financial_indicators
        }
        return self.health_entropy_model.predict(data)
    
    def generate_entropy_reduction_strategy(self, state):
        """生成熵减策略
        
        Args:
            state: 当前状态
            
        Returns:
            熵减策略
        """
        data = {
            'state': state
        }
        return self.rl_optimization_model.predict(data)
    
    def verify_four_process(self, process_data):
        """验证四流程
        
        Args:
            process_data: 四流程数据
            
        Returns:
            四流程验证结果
        """
        data = {
            'process_data': process_data
        }
        return self.four_process_model.predict(data)
    
    def run_full_analysis(self, data):
        """运行完整分析
        
        Args:
            data: 包含所有必要数据的字典
            
        Returns:
            完整分析结果
        """
        results = {}
        
        # 1. 计算熵变平衡
        internal_entropy = data.get('internal_entropy', 0)
        external_entropy = data.get('external_entropy', 0)
        results['entropy_balance'] = self.calculate_entropy_balance(internal_entropy, external_entropy)
        
        # 2. 检查五维时间线平衡
        dimension_data = data.get('dimension_data', {})
        results['five_dimension_balance'] = self.check_five_dimension_balance(dimension_data)
        
        # 3. 评估财务健康熵值
        financial_indicators = data.get('financial_indicators', {})
        results['health_entropy'] = self.evaluate_health_entropy(financial_indicators)
        
        # 4. 生成熵减策略
        state = data.get('state', None)
        if state:
            results['entropy_reduction_strategy'] = self.generate_entropy_reduction_strategy(state)
        
        # 5. 验证四流程
        process_data = data.get('process_data', {})
        if process_data:
            results['four_process_verification'] = self.verify_four_process(process_data)
        
        return results
