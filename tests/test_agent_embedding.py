#!/usr/bin/env python3
"""
智能体嵌入包测试脚本

测试Meta-CreationPower智能体嵌入包的核心功能
"""

import sys
import os

# 添加项目根目录到Python路径
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from agent_embedding import EmbeddingAPI
from agent_embedding.utils import setup_logger

def test_basic_functionality():
    """
    测试基本功能
    """
    print("=" * 80)
    print("测试1: 基本功能测试")
    print("=" * 80)
    
    try:
        # 初始化嵌入API
        api = EmbeddingAPI(
            agent_name="测试智能体",
            agent_description="用于测试的智能体"
        )
        print("✓ 初始化嵌入API成功")
        
        # 注册碳基伙伴
        partner_id = api.register_carbon_partner(
            partner_name="测试伙伴",
            partner_description="用于测试的碳基伙伴"
        )
        print(f"✓ 注册碳基伙伴成功，ID: {partner_id}")
        
        # 快速启动协同
        result = api.quick_start_collaboration(
            partner_name="测试伙伴",
            theme="测试主题",
            collaboration_type="staggered_complement"
        )
        print(f"✓ 快速启动协同成功，结果: {'成功' if result['success'] else '失败'}")
        
        # 计算系统健康状态
        health_status = api.calculate_system_health()
        print(f"✓ 计算系统健康状态成功，熵值: {health_status['entropy_score']:.2f}")
        
        # 获取智能体信息
        agent_info = api.get_agent_info()
        print(f"✓ 获取智能体信息成功，名称: {agent_info['name']}")
        
        # 获取碳基伙伴
        partners = api.get_carbon_partners()
        print(f"✓ 获取碳基伙伴成功，数量: {len(partners)}")
        
        print("\n✓ 基本功能测试通过！")
        return True
        
    except Exception as e:
        print(f"✗ 基本功能测试失败: {str(e)}")
        return False

def test_collaboration_types():
    """
    测试不同类型的协同
    """
    print("\n" + "=" * 80)
    print("测试2: 协同类型测试")
    print("=" * 80)
    
    try:
        # 初始化嵌入API
        api = EmbeddingAPI(
            agent_name="类型测试智能体",
            agent_description="测试不同协同类型的智能体"
        )
        
        # 注册碳基伙伴
        api.register_carbon_partner(
            partner_name="类型测试伙伴",
            partner_description="用于测试不同协同类型的碳基伙伴"
        )
        
        # 测试错位互补模式
        result1 = api.create_staggered_complement_collaboration(
            partner_name="类型测试伙伴",
            collaboration_name="错位互补测试",
            creation_theme="测试错位互补模式"
        )
        print(f"✓ 错位互补模式测试成功，结果: {'成功' if result1['success'] else '失败'}")
        
        # 测试卡农式推进模式
        result2 = api.create_canon_progression_collaboration(
            partner_name="类型测试伙伴",
            collaboration_name="卡农式推进测试",
            creation_theme="测试卡农式推进模式"
        )
        print(f"✓ 卡农式推进模式测试成功，结果: {'成功' if result2['success'] else '失败'}")
        
        # 测试赋格式交织模式
        result3 = api.create_fugue_interweaving_collaboration(
            partner_name="类型测试伙伴",
            collaboration_name="赋格式交织测试",
            creation_theme="测试赋格式交织模式"
        )
        print(f"✓ 赋格式交织模式测试成功，结果: {'成功' if result3['success'] else '失败'}")
        
        print("\n✓ 协同类型测试通过！")
        return True
        
    except Exception as e:
        print(f"✗ 协同类型测试失败: {str(e)}")
        return False

def test_validation():
    """
    测试协同验证功能
    """
    print("\n" + "=" * 80)
    print("测试3: 协同验证测试")
    print("=" * 80)
    
    try:
        # 初始化嵌入API
        api = EmbeddingAPI(
            agent_name="验证测试智能体",
            agent_description="测试协同验证功能的智能体"
        )
        
        # 模拟碳基意图和硅基输出
        carbon_intention = {
            "theme": "测试主题",
            "style": "测试风格",
            "emotion": "积极",
            "requirements": ["测试要求1", "测试要求2"]
        }
        
        silicon_output = {
            "theme": "测试主题",
            "style": "测试风格",
            "emotion": "积极",
            "content": "测试内容",
            "recommendations": ["测试建议1", "测试建议2"]
        }
        
        # 验证协同
        validation_result = api.validate_collaboration(
            carbon_intention=carbon_intention,
            silicon_output=silicon_output
        )
        print(f"✓ 协同验证测试成功，验证ID: {validation_result.validation_id}")
        print(f"  差异数量: {len(validation_result.differences)}")
        
        print("\n✓ 协同验证测试通过！")
        return True
        
    except Exception as e:
        print(f"✗ 协同验证测试失败: {str(e)}")
        return False

def test_custom_capabilities():
    """
    测试自定义能力和意图向量
    """
    print("\n" + "=" * 80)
    print("测试4: 自定义能力测试")
    print("=" * 80)
    
    try:
        # 自定义能力和意图向量
        custom_capabilities = {
            "创意生成": 0.9,
            "逻辑分析": 0.8,
            "情感共鸣": 0.7,
            "技术实现": 0.95
        }
        
        custom_intentions = {
            "探索性": 0.8,
            "完美性": 0.9,
            "效率": 0.7,
            "创新性": 0.9
        }
        
        # 初始化嵌入API
        api = EmbeddingAPI(
            agent_name="自定义能力智能体",
            agent_description="具有自定义能力的智能体",
            capabilities=custom_capabilities,
            intentions=custom_intentions
        )
        print("✓ 初始化自定义能力智能体成功")
        
        # 注册具有自定义能力的碳基伙伴
        partner_capabilities = {
            "创意生成": 0.95,
            "逻辑分析": 0.7,
            "情感共鸣": 0.9,
            "艺术感知": 0.85
        }
        
        partner_intentions = {
            "探索性": 0.9,
            "完美性": 0.8,
            "效率": 0.6,
            "美学追求": 0.9
        }
        
        partner_id = api.register_carbon_partner(
            partner_name="自定义能力伙伴",
            partner_description="具有自定义能力的碳基伙伴",
            capabilities=partner_capabilities,
            intentions=partner_intentions
        )
        print(f"✓ 注册自定义能力碳基伙伴成功，ID: {partner_id}")
        
        # 获取智能体信息
        agent_info = api.get_agent_info()
        print(f"✓ 获取智能体信息成功，能力数量: {len(agent_info['capabilities'])}")
        
        print("\n✓ 自定义能力测试通过！")
        return True
        
    except Exception as e:
        print(f"✗ 自定义能力测试失败: {str(e)}")
        return False

def run_all_tests():
    """
    运行所有测试
    """
    print("\n" + "=" * 80)
    print("开始运行智能体嵌入包测试")
    print("=" * 80)
    
    # 运行所有测试
    tests = [
        test_basic_functionality,
        test_collaboration_types,
        test_validation,
        test_custom_capabilities
    ]
    
    passed_tests = 0
    total_tests = len(tests)
    
    for test in tests:
        if test():
            passed_tests += 1
        print()
    
    # 打印测试结果
    print("=" * 80)
    print("测试结果总结")
    print("=" * 80)
    print(f"总测试数: {total_tests}")
    print(f"通过测试数: {passed_tests}")
    print(f"失败测试数: {total_tests - passed_tests}")
    print(f"测试通过率: {(passed_tests / total_tests) * 100:.1f}%")
    
    if passed_tests == total_tests:
        print("\n🎉 所有测试通过！智能体嵌入包功能正常。")
        return True
    else:
        print("\n❌ 部分测试失败，请检查错误信息。")
        return False

if __name__ == "__main__":
    # 运行所有测试
    success = run_all_tests()
    
    # 根据测试结果设置退出码
    sys.exit(0 if success else 1)
