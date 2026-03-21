# 对称数智财务系统主入口

import json
import sys
import os

# 添加项目根目录到Python路径
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from config.config import config
from src.modules.model_integration import ModelIntegration

class SymmetricDigitalFinanceSystem:
    """对称数智财务系统"""
    
    def __init__(self):
        """初始化系统"""
        self.config = config
        self.model_integration = ModelIntegration(config)
        print(f"初始化{config.SYSTEM_NAME} v{config.SYSTEM_VERSION}")
    
    def run(self, input_data):
        """运行系统
        
        Args:
            input_data: 输入数据
            
        Returns:
            系统运行结果
        """
        try:
            # 运行完整分析
            results = self.model_integration.run_full_analysis(input_data)
            return results
        except Exception as e:
            return {
                'error': str(e),
                'status': '系统运行失败'
            }
    
    def health_check(self):
        """健康检查
        
        Returns:
            健康检查结果
        """
        try:
            # 简单的健康检查
            test_data = {
                'internal_entropy': 0.1,
                'external_entropy': 0.1,
                'financial_indicators': {
                    'cash_flow': 100000,
                    'revenue': 500000,
                    'expenses': 400000,
                    'assets': 1000000,
                    'liabilities': 500000
                },
                'dimension_data': {
                    '宏观调控': [0.1, 0.2, 0.3],
                    '市场动态': [0.2, 0.3, 0.4],
                    '决策模型': [0.3, 0.4, 0.5],
                    '风控模型': [0.4, 0.5, 0.6],
                    '执行管理': [0.5, 0.6, 0.7]
                }
            }
            results = self.run(test_data)
            return {
                'status': '系统健康',
                'test_results': results
            }
        except Exception as e:
            return {
                'status': '系统异常',
                'error': str(e)
            }

if __name__ == "__main__":
    # 初始化系统
    system = SymmetricDigitalFinanceSystem()
    
    # 健康检查
    health_result = system.health_check()
    print("健康检查结果:")
    print(json.dumps(health_result, ensure_ascii=False, indent=2))
    
    # 示例运行
    if len(sys.argv) > 1:
        input_file = sys.argv[1]
        if os.path.exists(input_file):
            with open(input_file, 'r', encoding='utf-8') as f:
                input_data = json.load(f)
            results = system.run(input_data)
            print("系统运行结果:")
            print(json.dumps(results, ensure_ascii=False, indent=2))
        else:
            print(f"输入文件不存在: {input_file}")
    else:
        print("请提供输入数据文件路径")
