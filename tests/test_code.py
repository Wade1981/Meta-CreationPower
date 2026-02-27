#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
测试 Meta-CreationPower 项目代码是否可以正常运行
"""

import sys
import os

# 添加项目根目录到 Python 路径
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))


def test_voice_recognition():
    """测试声部识别层"""
    print("\n=== 测试声部识别层 ===")
    try:
        from src.layers.voice_recognition.voice_recognition import CollaborativeSonicMap
        print("✅ 成功导入 CollaborativeSonicMap")
        
        # 测试创建对象
        sonic_map = CollaborativeSonicMap()
        print("✅ 成功创建 CollaborativeSonicMap 对象")
        
        # 测试注册碳基声部
        carbon_voice = sonic_map.register_voice(
            name="X54先生",
            voice_type="carbon",
            capability_vector={"创意生成": 0.9, "逻辑分析": 0.7, "情感共鸣": 0.9},
            intention_vector={"探索性": 0.8, "完美性": 0.7, "效率": 0.6}
        )
        print(f"✅ 成功注册碳基声部，ID: {carbon_voice.voice_id}")
        
        # 测试注册硅基声部
        silicon_voice = sonic_map.register_voice(
            name="豆包",
            voice_type="silicon",
            capability_vector={"创意生成": 0.8, "逻辑分析": 0.9, "情感共鸣": 0.5},
            intention_vector={"探索性": 0.7, "完美性": 0.8, "效率": 0.9}
        )
        print(f"✅ 成功注册硅基声部，ID: {silicon_voice.voice_id}")
        
        # 测试获取声部
        retrieved_voice = sonic_map.get_voice(carbon_voice.voice_id)
        print(f"✅ 成功获取声部: {retrieved_voice.name}")
        
        # 测试按类型获取声部
        carbon_voices = sonic_map.get_voices_by_type("carbon")
        silicon_voices = sonic_map.get_voices_by_type("silicon")
        print(f"✅ 成功按类型获取声部: 碳基 {len(carbon_voices)} 个，硅基 {len(silicon_voices)} 个")
        
        return True
    except Exception as e:
        print(f"❌ 测试失败: {e}")
        import traceback
        traceback.print_exc()
        return False


def test_counterpoint_design():
    """测试协奏设计层"""
    print("\n=== 测试协奏设计层 ===")
    try:
        from src.layers.counterpoint_design.counterpoint_design import CounterpointDesigner
        print("✅ 成功导入 CounterpointDesigner")
        
        # 测试创建对象
        designer = CounterpointDesigner()
        print("✅ 成功创建 CounterpointDesigner 对象")
        
        # 测试创建协同路径
        path = designer.create_counterpoint_path(
            name="创意写作协同",
            pattern_type="staggered_complement",
            participating_voices=["carbon_1", "silicon_1"],
            creation_theme="探索人工智能与人类创造力的边界"
        )
        print(f"✅ 成功创建协同路径，ID: {path.path_id}")
        print(f"  模式类型: {path.pattern_type}")
        print(f"  参与声部: {path.participating_voices}")
        print(f"  步骤数量: {len(path.steps)}")
        
        # 测试执行路径步骤
        result = designer.execute_path_step(
            path_id=path.path_id,
            step_index=0,
            voice_id="carbon_1",
            inputs={"idea": "人工智能与人类协作的未来"}
        )
        print(f"✅ 成功执行路径步骤: {result.get('outputs', {}).get('message', '')}")
        
        # 测试获取适合的模式
        suitable_patterns = designer.get_suitable_patterns("概念设计")
        print(f"✅ 成功获取适合的模式: {[p['name'] for p in suitable_patterns]}")
        
        return True
    except Exception as e:
        print(f"❌ 测试失败: {e}")
        import traceback
        traceback.print_exc()
        return False


def test_steady_execution():
    """测试静定执行层"""
    print("\n=== 测试静定执行层 ===")
    try:
        from src.layers.steady_execution.steady_execution import SteadyExecutor
        print("✅ 成功导入 SteadyExecutor")
        
        # 测试创建对象
        executor = SteadyExecutor()
        print("✅ 成功创建 SteadyExecutor 对象")
        
        # 测试提交任务
        task_id = executor.submit_task(
            name="测试任务",
            task_type="test",
            payload={"test": "data"}
        )
        print(f"✅ 成功提交任务，ID: {task_id}")
        
        # 测试获取任务状态
        import time
        time.sleep(0.1)  # 等待任务执行
        status = executor.get_task_status(task_id)
        print(f"✅ 成功获取任务状态: {status.get('status')}")
        
        # 测试执行协同路径
        steps = [
            {"step": 1, "role": "carbon", "action": "提出模糊概念"},
            {"step": 2, "role": "silicon", "action": "生成百种变体"}
        ]
        voice_map = {
            "carbon": "carbon_1",
            "silicon": "silicon_1"
        }
        execution_result = executor.execute_counterpoint_path(
            path_id="test_path",
            steps=steps,
            voice_map=voice_map
        )
        print(f"✅ 成功执行协同路径，执行ID: {execution_result.get('execution_id')}")
        print(f"  成功: {execution_result.get('success')}")
        
        return True
    except Exception as e:
        print(f"❌ 测试失败: {e}")
        import traceback
        traceback.print_exc()
        return False


def test_main():
    """测试主程序"""
    print("\n=== 测试主程序 ===")
    try:
        from src.main import main
        print("✅ 成功导入 main 函数")
        print("执行主程序...")
        main()
        print("✅ 主程序执行成功")
        return True
    except Exception as e:
        print(f"❌ 测试失败: {e}")
        import traceback
        traceback.print_exc()
        return False


def run_all_tests():
    """运行所有测试"""
    print("开始测试 Meta-CreationPower 项目代码...")
    print("=" * 60)
    
    tests = [
        test_voice_recognition,
        test_counterpoint_design,
        test_steady_execution,
        test_main
    ]
    
    passed_tests = 0
    total_tests = len(tests)
    
    for test in tests:
        if test():
            passed_tests += 1
        print()
    
    # 打印测试结果
    print("=" * 60)
    print("测试结果总结")
    print("=" * 60)
    print(f"总测试数: {total_tests}")
    print(f"通过测试数: {passed_tests}")
    print(f"失败测试数: {total_tests - passed_tests}")
    print(f"测试通过率: {(passed_tests / total_tests) * 100:.1f}%")
    
    if passed_tests == total_tests:
        print("\n✅ 所有测试通过，项目代码可以正常运行！")
        return True
    else:
        print("\n❌ 部分测试失败，请检查错误信息并修复问题！")
        return False


if __name__ == "__main__":
    success = run_all_tests()
    sys.exit(0 if success else 1)
