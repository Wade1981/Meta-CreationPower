# 对称数智财务系统配置文件

class Config:
    """系统配置类"""
    # 系统基本配置
    SYSTEM_NAME = "对称数智财务系统"
    SYSTEM_VERSION = "1.0.0"
    DEBUG = True
    
    # 数据配置
    DATA_DIR = "data"
    LOG_DIR = "logs"
    
    # 熵值计算配置
    ENTROPY_CALCULATION_INTERVAL = 60  # 熵值计算间隔（秒）
    HEALTH_ENTROPY_THRESHOLD_LOW = 30  # 低熵区阈值
    HEALTH_ENTROPY_THRESHOLD_HIGH = 70  # 高熵区阈值
    
    # 五维时间线配置
    FIVE_DIMENSIONS = [
        "宏观调控",
        "市场动态",
        "决策模型",
        "风控模型",
        "执行管理"
    ]
    DIMENSION_ENTROPY_DEVIATION_THRESHOLD = 0.05  # 维度熵变速率偏差阈值
    SINGLE_DIMENSION_ENTROPY_LIMIT = 0.3  # 单维度熵值占比上限
    
    # 强化学习配置
    RL_LEARNING_RATE = 0.001
    RL_DISCOUNT_FACTOR = 0.99
    RL_EXPLORATION_RATE = 0.1
    RL_BATCH_SIZE = 64
    RL_MEMORY_SIZE = 10000
    
    # 熵变权重配置
    ENTROPY_WEIGHTS = {
        "total_entropy": 0.5,
        "health_entropy": 0.3,
        "net_income": 0.2
    }
    
    # ELR容器配置
    ELR_CONTAINER_NAME = "symmetric-finance-system"
    ELR_CONTAINER_VERSION = "1.0.0"
    ELR_REQUIREMENTS = [
        "numpy",
        "pandas",
        "scikit-learn",
        "tensorflow",
        "matplotlib"
    ]

# 实例化配置对象
config = Config()
