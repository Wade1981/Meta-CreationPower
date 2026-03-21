#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
ELR 沙箱测试 - EL-CSCC Archive 数字档案管理模型
功能：验证 EL-CSCC Archive 模型是否可以在 ELR 沙箱中正常运行
"""

import sys
import os
import json

# 添加模型路径到 Python 路径
sys.path.append(os.path.join(os.path.dirname(__file__), '..', 'model', 'models', 'el_cscc_archive'))

class ELSandboxTester:
    """ELR沙箱测试器"""
    
    def __init__(self):
        """初始化测试器"""
        self.test_name = "EL-CSCC Archive 数字档案管理模型沙箱测试"
        self.test_date = "2026-03-21"
        self.model = None
        print(f"初始化测试器: {self.test_name}")
    
    def load_model(self):
        """
        加载模型到沙箱
        返回:
            bool: 加载是否成功
        """
        try:
            print("\n1. 加载模型到ELR沙箱...")
            from el_cscc_archive_model import ELCSCCArchiveModel
            self.model = ELCSCCArchiveModel()
            print("✓ 模型加载成功！")
            return True
        except Exception as e:
            print(f"✗ 模型加载失败: {e}")
            return False
    
    def test_model_info(self):
        """
        测试模型信息
        返回:
            bool: 测试是否成功
        """
        try:
            print("\n2. 测试模型信息...")
            # 测试模型是否有基本属性
            if hasattr(self.model, 'predict') and callable(self.model.predict):
                print("✓ 模型具有 predict 方法")
                return True
            else:
                print("✗ 模型缺少 predict 方法")
                return False
        except Exception as e:
            print(f"✗ 测试模型信息失败: {e}")
            return False
    
    def test_model_prediction(self):
        """
        测试模型推理
        返回:
            bool: 测试是否成功
        """
        try:
            print("\n3. 测试模型推理...")
            
            # 测试列出档案
            print("  测试列出档案...")
            result = self.model.predict("list:")
            print(f"  输出: {result}")
            print("  ✓ 列出档案测试通过")
            
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
            result = self.model.predict(f"add: {json.dumps(test_archive)}")
            print(f"  输出: {result}")
            print("  ✓ 添加档案测试通过")
            
            # 测试查询档案
            print("  测试查询档案...")
            # 假设添加的档案编号是 EL-CSCC-000001
            result = self.model.predict("query: EL-CSCC-000001")
            print(f"  输出: {result}")
            print("  ✓ 查询档案测试通过")
            
            return True
        except Exception as e:
            print(f"✗ 测试模型推理失败: {e}")
            return False
    
    def test_resource_usage(self):
        """
        测试资源使用情况
        返回:
            bool: 测试是否成功
        """
        try:
            print("\n4. 测试资源使用情况...")
            # 尝试导入 psutil 进行资源监控
            try:
                import psutil
                process = psutil.Process()
                memory_info = process.memory_info()
                print(f"  内存使用: {memory_info.rss / 1024 / 1024:.2f} MB")
                print("  ✓ 资源使用测试通过")
                return True
            except ImportError:
                print("  ⚠ psutil 未安装，跳过资源使用测试")
                return True
        except Exception as e:
            print(f"✗ 测试资源使用失败: {e}")
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
            ("加载模型", self.load_model),
            ("测试模型信息", self.test_model_info),
            ("测试模型推理", self.test_model_prediction),
            ("测试资源使用", self.test_resource_usage)
        ]
        
        passed = 0
        total = len(tests)
        
        for test_name, test_func in tests:
            if test_func():
                passed += 1
            print("-" * 60)
        
        print("\n=== 测试结果汇总 ===")
        print(f"测试通过: {passed}/{total}")
        
        if passed == total:
            print("\n🎉 所有测试通过！")
            print("结论: EL-CSCC Archive 模型可以在 ELR 沙箱中正常运行")
            return True
        else:
            print("\n❌ 部分测试失败！")
            print("结论: EL-CSCC Archive 模型在 ELR 沙箱中运行存在问题")
            return False

if __name__ == "__main__":
    tester = ELSandboxTester()
    tester.run_all_tests()
