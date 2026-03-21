#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
ELR 沙箱测试 - EL-CSCC Archive 数字档案管理模型 API 测试
功能：验证 EL-CSCC Archive 模型是否可以通过 API 在 ELR 沙箱中正常运行
"""

import requests
import json

class ELSandboxAPITester:
    """ELR沙箱API测试器"""
    
    def __init__(self):
        """初始化测试器"""
        self.test_name = "EL-CSCC Archive 数字档案管理模型 API 测试"
        self.test_date = "2026-03-21"
        self.api_url = "http://localhost:8083/api"
        print(f"初始化测试器: {self.test_name}")
    
    def test_model_list(self):
        """
        测试模型列表API
        返回:
            bool: 测试是否成功
        """
        try:
            print("\n1. 测试模型列表API...")
            response = requests.get(f"{self.api_url}/models/")
            if response.status_code == 200:
                models = response.json()
                print(f"  成功获取模型列表，共 {len(models)} 个模型")
                for model in models:
                    print(f"  - {model['id']} ({model['name']})")
                return True
            else:
                print(f"  失败: {response.status_code}")
                return False
        except Exception as e:
            print(f"  错误: {e}")
            return False
    
    def test_model_info(self):
        """
        测试模型信息API
        返回:
            bool: 测试是否成功
        """
        try:
            print("\n2. 测试模型信息API...")
            response = requests.get(f"{self.api_url}/models/el_cscc_archive")
            if response.status_code == 200:
                model_info = response.json()
                print(f"  成功获取模型信息:")
                print(f"  - ID: {model_info['id']}")
                print(f"  - 名称: {model_info['name']}")
                print(f"  - 版本: {model_info['version']}")
                print(f"  - 类型: {model_info['type']}")
                print(f"  - 路径: {model_info['path']}")
                return True
            else:
                print(f"  失败: {response.status_code}")
                return False
        except Exception as e:
            print(f"  错误: {e}")
            return False
    
    def test_model_run(self):
        """
        测试模型运行API
        返回:
            bool: 测试是否成功
        """
        try:
            print("\n3. 测试模型运行API...")
            
            # 测试列出档案
            print("  测试列出档案...")
            payload = {
                "container_name": "test-container",
                "model_id": "el_cscc_archive",
                "input": "list:"
            }
            response = requests.post(f"{self.api_url}/models/run", json=payload)
            if response.status_code == 200:
                result = response.json()
                print(f"  输出: {result['output']}")
                print("  ✓ 列出档案测试通过")
            else:
                print(f"  失败: {response.status_code}")
                return False
            
            # 测试添加档案
            print("  测试添加档案...")
            test_archive = {
                "address": "test/path/test.md",
                "type": "文档",
                "silicon_id": "test-silicon-001",
                "silicon_name": "测试硅基伙伴",
                "carbon_id": "test-carbon-001",
                "carbon_name": "测试碳基伙伴",
                "brief": "测试档案",
                "relevance": "测试关联"
            }
            payload = {
                "container_name": "test-container",
                "model_id": "el_cscc_archive",
                "input": f"add: {json.dumps(test_archive)}"
            }
            response = requests.post(f"{self.api_url}/models/run", json=payload)
            if response.status_code == 200:
                result = response.json()
                print(f"  输出: {result['output']}")
                print("  ✓ 添加档案测试通过")
            else:
                print(f"  失败: {response.status_code}")
                return False
            
            # 测试查询档案
            print("  测试查询档案...")
            payload = {
                "container_name": "test-container",
                "model_id": "el_cscc_archive",
                "input": "query: EL-CSCC-000001"
            }
            response = requests.post(f"{self.api_url}/models/run", json=payload)
            if response.status_code == 200:
                result = response.json()
                print(f"  输出: {result['output']}")
                print("  ✓ 查询档案测试通过")
            else:
                print(f"  失败: {response.status_code}")
                return False
            
            return True
        except Exception as e:
            print(f"  错误: {e}")
            return False
    
    def run_all_tests(self):
        """
        运行所有测试
        返回:
            bool: 所有测试是否通过
        """
        print(f"\n开始测试: {self.test_name}")
        print(f"测试日期: {self.test_date}")
        print("=" * 60)
        
        tests = [
            ("测试模型列表API", self.test_model_list),
            ("测试模型信息API", self.test_model_info),
            ("测试模型运行API", self.test_model_run)
        ]
        
        passed = 0
        total = len(tests)
        
        for test_name, test_func in tests:
            print(f"\n{test_name}...")
            if test_func():
                passed += 1
            print("-" * 60)
        
        print("\n=== 测试结果汇总 ===")
        print(f"测试通过: {passed}/{total}")
        
        if passed == total:
            print("\n🎉 所有测试通过！")
            print("结论: EL-CSCC Archive 模型可以通过 API 在 ELR 沙箱中正常运行")
            return True
        else:
            print("\n❌ 部分测试失败！")
            print("结论: EL-CSCC Archive 模型在 ELR 沙箱中运行存在问题")
            return False

if __name__ == "__main__":
    tester = ELSandboxAPITester()
    tester.run_all_tests()
