# RootPulseOS Root Sensor Implementation

"""RootPulseOS根脉传感器实现，负责感知碳基世界的需求与变化。"""

import logging
import time
from typing import Dict, List, Optional, Any

class RootSensor:
    """根脉传感器类，负责感知碳基世界的需求与变化。"""
    
    def __init__(self):
        """初始化RootSensor实例。"""
        self.logger = logging.getLogger(__name__)
        self.logger.info("Initializing RootSensor...")
        
        self.sensors = {}
        self.running = False
    
    def register_sensor(self, name: str, sensor):
        """注册传感器。
        
        Args:
            name: 传感器名称
            sensor: 传感器实例
        """
        self.sensors[name] = sensor
        self.logger.info(f"Registered sensor: {name}")
    
    def start(self):
        """启动根脉传感器。"""
        self.logger.info("Starting RootSensor...")
        self.running = True
        
        # 启动所有传感器
        for name, sensor in self.sensors.items():
            if hasattr(sensor, "start"):
                try:
                    sensor.start()
                    self.logger.info(f"Started sensor: {name}")
                except Exception as e:
                    self.logger.error(f"Failed to start sensor {name}: {e}")
    
    def stop(self):
        """停止根脉传感器。"""
        self.logger.info("Stopping RootSensor...")
        self.running = False
        
        # 停止所有传感器
        for name, sensor in reversed(list(self.sensors.items())):
            if hasattr(sensor, "stop"):
                try:
                    sensor.stop()
                    self.logger.info(f"Stopped sensor: {name}")
                except Exception as e:
                    self.logger.error(f"Failed to stop sensor {name}: {e}")
    
    def sense(self) -> Dict[str, Any]:
        """感知碳基世界的需求与变化。
        
        Returns:
            Dict[str, Any]: 感知结果
        """
        if not self.running:
            self.logger.warning("RootSensor is not running")
            return {}
        
        results = {}
        for name, sensor in self.sensors.items():
            if hasattr(sensor, "sense"):
                try:
                    results[name] = sensor.sense()
                except Exception as e:
                    self.logger.error(f"Failed to sense with {name}: {e}")
                    results[name] = None
        
        return results
    
    def get_sensor(self, name: str):
        """获取指定传感器。
        
        Args:
            name: 传感器名称
            
        Returns:
            传感器实例或None
        """
        return self.sensors.get(name)
    
    def get_all_sensors(self) -> Dict[str, object]:
        """获取所有传感器。
        
        Returns:
            Dict[str, object]: 传感器字典
        """
        return self.sensors
    
    def status(self) -> Dict[str, any]:
        """获取根脉传感器状态。
        
        Returns:
            Dict[str, any]: 状态信息
        """
        return {
            "running": self.running,
            "sensors": list(self.sensors.keys())
        }