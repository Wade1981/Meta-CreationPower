# RootPulseOS ELR Test Interface Implementation

"""RootPulseOS ELR测试接口实现，负责与ELR智能测试系统进行交互。"""

import logging
import requests
import json
from typing import Dict, List, Optional, Any

class ELRTestInterface:
    """ELR测试接口类，负责与ELR智能测试系统进行交互。"""
    
    def __init__(self, base_url: str = None):
        """初始化ELRTestInterface实例。
        
        Args:
            base_url: ELR测试系统API的基础URL，如果为None则从环境变量读取
        """
        import os
        
        self.logger = logging.getLogger(__name__)
        self.logger.info("Initializing ELRTestInterface...")
        
        # 从环境变量读取API URL
        self.base_url = base_url or os.environ.get("ELR_TEST_API_URL", "http://localhost:8000")
        self.logger.info(f"Using ELR test API URL: {self.base_url}")
        self.running = False
    
    def start(self):
        """启动ELR测试接口。"""
        self.logger.info("Starting ELRTestInterface...")
        self.running = True
        
        # 测试连接
        if self._test_connection():
            self.logger.info("Successfully connected to ELR test system")
        else:
            self.logger.warning("Failed to connect to ELR test system")
    
    def stop(self):
        """停止ELR测试接口。"""
        self.logger.info("Stopping ELRTestInterface...")
        self.running = False
    
    def _test_connection(self) -> bool:
        """测试与ELR测试系统的连接。
        
        Returns:
            bool: 连接是否成功
        """
        try:
            response = requests.get(f"{self.base_url}/health", timeout=5)
            return response.status_code == 200
        except Exception as e:
            self.logger.error(f"Connection test failed: {e}")
            return False
    
    def get_test_cases(self) -> List[Dict[str, Any]]:
        """获取所有测试用例。
        
        Returns:
            List[Dict[str, Any]]: 测试用例列表
        """
        if not self.running:
            self.logger.warning("ELRTestInterface is not running")
            return []
        
        try:
            response = requests.get(f"{self.base_url}/test_cases")
            if response.status_code == 200:
                return response.json()
            else:
                self.logger.error(f"Failed to get test cases: {response.status_code}")
                return []
        except Exception as e:
            self.logger.error(f"Error getting test cases: {e}")
            return []
    
    def run_test_case(self, test_case_id: str) -> Dict[str, Any]:
        """运行指定的测试用例。
        
        Args:
            test_case_id: 测试用例ID
            
        Returns:
            Dict[str, Any]: 测试结果
        """
        if not self.running:
            self.logger.warning("ELRTestInterface is not running")
            return {"status": "error", "message": "Interface not running"}
        
        try:
            response = requests.post(f"{self.base_url}/test_cases/{test_case_id}/run")
            if response.status_code == 200:
                return response.json()
            else:
                self.logger.error(f"Failed to run test case: {response.status_code}")
                return {"status": "error", "message": f"Failed with status code: {response.status_code}"}
        except Exception as e:
            self.logger.error(f"Error running test case: {e}")
            return {"status": "error", "message": str(e)}
    
    def get_test_results(self, test_case_id: str) -> Dict[str, Any]:
        """获取测试用例的结果。
        
        Args:
            test_case_id: 测试用例ID
            
        Returns:
            Dict[str, Any]: 测试结果
        """
        if not self.running:
            self.logger.warning("ELRTestInterface is not running")
            return {}
        
        try:
            response = requests.get(f"{self.base_url}/test_cases/{test_case_id}/results")
            if response.status_code == 200:
                return response.json()
            else:
                self.logger.error(f"Failed to get test results: {response.status_code}")
                return {}
        except Exception as e:
            self.logger.error(f"Error getting test results: {e}")
            return {}
    
    def create_test_plan(self, plan_name: str, test_cases: List[str]) -> Dict[str, Any]:
        """创建测试计划。
        
        Args:
            plan_name: 测试计划名称
            test_cases: 测试用例ID列表
            
        Returns:
            Dict[str, Any]: 创建结果
        """
        if not self.running:
            self.logger.warning("ELRTestInterface is not running")
            return {"status": "error", "message": "Interface not running"}
        
        try:
            payload = {
                "name": plan_name,
                "test_cases": test_cases
            }
            response = requests.post(f"{self.base_url}/test_plans", json=payload)
            if response.status_code == 200:
                return response.json()
            else:
                self.logger.error(f"Failed to create test plan: {response.status_code}")
                return {"status": "error", "message": f"Failed with status code: {response.status_code}"}
        except Exception as e:
            self.logger.error(f"Error creating test plan: {e}")
            return {"status": "error", "message": str(e)}
    
    def run_test_plan(self, plan_id: str) -> Dict[str, Any]:
        """运行测试计划。
        
        Args:
            plan_id: 测试计划ID
            
        Returns:
            Dict[str, Any]: 运行结果
        """
        if not self.running:
            self.logger.warning("ELRTestInterface is not running")
            return {"status": "error", "message": "Interface not running"}
        
        try:
            response = requests.post(f"{self.base_url}/test_plans/{plan_id}/run")
            if response.status_code == 200:
                return response.json()
            else:
                self.logger.error(f"Failed to run test plan: {response.status_code}")
                return {"status": "error", "message": f"Failed with status code: {response.status_code}"}
        except Exception as e:
            self.logger.error(f"Error running test plan: {e}")
            return {"status": "error", "message": str(e)}
    
    def get_system_status(self) -> Dict[str, Any]:
        """获取ELR测试系统的状态。
        
        Returns:
            Dict[str, Any]: 系统状态
        """
        if not self.running:
            self.logger.warning("ELRTestInterface is not running")
            return {"status": "error", "message": "Interface not running"}
        
        try:
            response = requests.get(f"{self.base_url}/status")
            if response.status_code == 200:
                return response.json()
            else:
                self.logger.error(f"Failed to get system status: {response.status_code}")
                return {"status": "error", "message": f"Failed with status code: {response.status_code}"}
        except Exception as e:
            self.logger.error(f"Error getting system status: {e}")
            return {"status": "error", "message": str(e)}
    
    def status(self) -> Dict[str, any]:
        """获取ELR测试接口状态。
        
        Returns:
            Dict[str, any]: 状态信息
        """
        return {
            "running": self.running,
            "base_url": self.base_url,
            "connected": self._test_connection()
        }