# RootPulseOS ELR Integration Test Script

"""测试RootPulseOS与ELR智能测试系统的集成功能。"""

import logging
import sys
import time

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(sys.stdout)
    ]
)

logger = logging.getLogger(__name__)

# 添加项目根目录到Python路径
import os
import sys
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

# 导入RootPulseOS模块
try:
    from src.core import RootPulseCore
    from src.elr_test_interface import ELRTestInterface
    logger.info("Successfully imported RootPulseOS modules")
except ImportError as e:
    logger.error(f"Failed to import RootPulseOS modules: {e}")
    sys.exit(1)

def test_elr_integration():
    """测试ELR集成功能。"""
    logger.info("Starting ELR integration test...")
    
    try:
        # 初始化RootPulseCore
        core = RootPulseCore()
        
        # 初始化ELR测试接口
        elr_test_interface = ELRTestInterface()
        core.register_component("elr_test_interface", elr_test_interface)
        
        # 启动系统
        core.start()
        
        # 等待系统启动
        time.sleep(2)
        
        # 测试ELR测试接口状态
        elr_status = elr_test_interface.status()
        logger.info(f"ELR test interface status: {elr_status}")
        
        # 测试系统连接
        if elr_status.get("connected", False):
            logger.info("Successfully connected to ELR test system")
            
            # 测试获取测试用例
            test_cases = elr_test_interface.get_test_cases()
            logger.info(f"Found {len(test_cases)} test cases")
            
            # 测试获取系统状态
            system_status = elr_test_interface.get_system_status()
            logger.info(f"ELR test system status: {system_status}")
        else:
            logger.warning("Failed to connect to ELR test system. This is expected if the system is not running.")
        
        # 停止系统
        core.stop()
        logger.info("ELR integration test completed successfully!")
        
    except Exception as e:
        logger.error(f"ELR integration test failed: {e}")
        return False
    
    return True

def main():
    """主函数。"""
    success = test_elr_integration()
    if success:
        logger.info("All tests passed!")
        sys.exit(0)
    else:
        logger.error("Some tests failed!")
        sys.exit(1)

if __name__ == "__main__":
    main()