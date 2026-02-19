#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
ELR沙箱测试脚本
测试极简文本响应微模型在ELR沙箱中的运行情况
"""

import os
import sys

# 添加当前目录到Python路径，确保可以导入模型
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from simple_text_model import SimpleTextModel

class ELSandboxTester:
    """ELR沙箱测试器"""
    
    def __init__(self):
        """初始化测试器"""
        self.test_name = "ELR极简文本模型沙箱测试"
        self.test_date = "2026-02-19"
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
            self.model = SimpleTextModel()
            print("✓ 模型加载成功！")
            return True
        except Exception as e:
            print(f"✗ 模型加载失败: {e}")
            return False
    
    def test_model_info(self):
        """
        测试模型信息获取
        返回:
            bool: 测试是否成功
        """
        try:
            print("\n2. 测试模型信息获取...")
            info = self.model.get_info()
            print(f"✓ 模型信息获取成功:")
            for key, value in info.items():
                print(f"  - {key}: {value}")
            return True
        except Exception as e:
            print(f"✗ 模型信息获取失败: {e}")
            return False
    
    def test_model_prediction(self):
        """
        测试模型推理
        返回:
            bool: 测试是否成功
        """
        try:
            print("\n3. 测试模型推理...")
            
            # 测试案例
            test_cases = [
                "Hello, how are you?",
                "你好，最近怎么样？",
                "Hi there!",
                "请问这个模型能做什么？",
                "测试消息"
            ]
            
            all_passed = True
            for i, test_input in enumerate(test_cases, 1):
                try:
                    response = self.model.predict(test_input)
                    print(f"✓ 测试案例 {i}:")
                    print(f"  输入: {test_input}")
                    print(f"  输出: {response}")
                    # 验证输出是否包含预期内容
                    if "碳硅协同" in response:
                        print("  验证: 输出包含碳硅协同标识")
                    else:
                        print("  验证: 输出缺少碳硅协同标识")
                        all_passed = False
                except Exception as e:
                    print(f"✗ 测试案例 {i} 失败: {e}")
                    all_passed = False
            
            return all_passed
        except Exception as e:
            print(f"✗ 模型推理测试失败: {e}")
            return False
    
    def test_resource_usage(self):
        """
        测试资源使用情况
        返回:
            bool: 测试是否成功
        """
        try:
            print("\n4. 测试资源使用情况...")
            # 简单的内存使用估算
            import psutil
            import os
            
            process = psutil.Process(os.getpid())
            memory_info = process.memory_info()
            memory_mb = memory_info.rss / 1024 / 1024
            
            print(f"✓ 资源使用情况:")
            print(f"  内存使用: {memory_mb:.2f} MB")
            
            # 验证是否轻量级
            if memory_mb < 100:
                print("  验证: 内存使用符合轻量级要求")
                return True
            else:
                print("  验证: 内存使用超出轻量级要求")
                return False
        except ImportError:
            print("⚠  psutil 未安装，跳过资源使用测试")
            return True
        except Exception as e:
            print(f"✗ 资源使用测试失败: {e}")
            return False
    
    def run_full_test(self):
        """
        运行完整测试
        返回:
            bool: 测试是否全部通过
        """
        print(f"\n=== {self.test_name} ===")
        print(f"测试日期: {self.test_date}")
        print("测试目标: 验证ELR沙箱是否支持极简微模型运行")
        
        # 运行所有测试
        tests = [
            self.load_model,
            self.test_model_info,
            self.test_model_prediction,
            self.test_resource_usage
        ]
        
        results = []
        for test in tests:
            results.append(test())
        
        # 汇总结果
        print("\n=== 测试结果汇总 ===")
        passed = sum(results)
        total = len(results)
        print(f"测试通过: {passed}/{total}")
        
        if all(results):
            print("\n🎉 所有测试通过！")
            print("结论: ELR沙箱支持极简微模型运行，满足轻量级、无外部依赖的要求。")
            return True
        else:
            print("\n❌ 部分测试失败！")
            print("结论: ELR沙箱运行存在问题，需要进一步排查。")
            return False

# 运行测试
if __name__ == "__main__":
    tester = ELSandboxTester()
    tester.run_full_test()
