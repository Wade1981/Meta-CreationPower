# RootPulseOS Main Entry Point

"""RootPulseOS主入口文件，负责启动和运行整个系统。"""

import logging
import sys
import time

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(sys.stdout),
        logging.FileHandler('rootpulseos.log')
    ]
)

logger = logging.getLogger(__name__)

# 导入核心模块
try:
    from src.core import RootPulseCore
    from src.root_sensor import RootSensor
    from src.cultural_translator import CulturalTranslator
    from src.elr_test_interface import ELRTestInterface
    logger.info("Successfully imported core modules")
except ImportError as e:
    logger.error(f"Failed to import core modules: {e}")
    sys.exit(1)

def main():
    """主函数，启动RootPulseOS系统。"""
    logger.info("Starting RootPulseOS...")
    
    try:
        # 初始化核心系统
        core = RootPulseCore()
        
        # 初始化根脉传感器
        root_sensor = RootSensor()
        core.register_component("root_sensor", root_sensor)
        
        # 初始化文化翻译器
        cultural_translator = CulturalTranslator()
        core.register_component("cultural_translator", cultural_translator)
        
        # 初始化ELR测试接口
        elr_test_interface = ELRTestInterface()
        core.register_component("elr_test_interface", elr_test_interface)
        
        # 启动系统
        core.start()
        
        logger.info("RootPulseOS started successfully!")
        logger.info(f"System status: {core.status()}")
        
        # 运行系统
        try:
            while True:
                # 模拟系统运行
                time.sleep(1)
                # 定期输出系统状态
                if int(time.time()) % 10 == 0:
                    logger.info(f"System status: {core.status()}")
        except KeyboardInterrupt:
            logger.info("Received keyboard interrupt, shutting down...")
        finally:
            # 停止系统
            core.stop()
            logger.info("RootPulseOS stopped successfully!")
            
    except Exception as e:
        logger.error(f"Failed to start RootPulseOS: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()