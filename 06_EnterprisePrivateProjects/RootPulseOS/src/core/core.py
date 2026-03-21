# RootPulseOS Core Implementation

"""RootPulseOS核心实现，负责系统的初始化、协调和管理。"""

import logging
import time
from typing import Dict, List, Optional

class RootPulseCore:
    """RootPulseOS核心类，管理系统的各个组件和模块。"""
    
    def __init__(self):
        """初始化RootPulseCore实例。"""
        self.logger = logging.getLogger(__name__)
        self.logger.info("Initializing RootPulseCore...")
        
        self.components = {}
        self.running = False
        self.start_time = None
    
    def register_component(self, name: str, component):
        """注册系统组件。
        
        Args:
            name: 组件名称
            component: 组件实例
        """
        self.components[name] = component
        self.logger.info(f"Registered component: {name}")
    
    def start(self):
        """启动RootPulseOS系统。"""
        self.logger.info("Starting RootPulseOS...")
        self.running = True
        self.start_time = time.time()
        
        # 启动所有组件
        for name, component in self.components.items():
            if hasattr(component, "start"):
                try:
                    component.start()
                    self.logger.info(f"Started component: {name}")
                except Exception as e:
                    self.logger.error(f"Failed to start component {name}: {e}")
    
    def stop(self):
        """停止RootPulseOS系统。"""
        self.logger.info("Stopping RootPulseOS...")
        self.running = False
        
        # 停止所有组件
        for name, component in reversed(list(self.components.items())):
            if hasattr(component, "stop"):
                try:
                    component.stop()
                    self.logger.info(f"Stopped component: {name}")
                except Exception as e:
                    self.logger.error(f"Failed to stop component {name}: {e}")
    
    def is_running(self) -> bool:
        """检查系统是否正在运行。
        
        Returns:
            bool: 系统运行状态
        """
        return self.running
    
    def get_uptime(self) -> float:
        """获取系统运行时间。
        
        Returns:
            float: 运行时间（秒）
        """
        if self.start_time:
            return time.time() - self.start_time
        return 0
    
    def get_component(self, name: str):
        """获取指定组件。
        
        Args:
            name: 组件名称
            
        Returns:
            组件实例或None
        """
        return self.components.get(name)
    
    def get_all_components(self) -> Dict[str, object]:
        """获取所有组件。
        
        Returns:
            Dict[str, object]: 组件字典
        """
        return self.components
    
    def status(self) -> Dict[str, any]:
        """获取系统状态。
        
        Returns:
            Dict[str, any]: 系统状态信息
        """
        return {
            "running": self.running,
            "uptime": self.get_uptime(),
            "components": list(self.components.keys()),
            "version": "0.1.0"
        }