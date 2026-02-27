#!/usr/bin/env python3
"""
测试 Meta-CreationPower 项目的各个模块是否可以正常运行
"""

import sys
import os

# 添加项目根目录到 Python 路径
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

def test_imports():
    """测试所有模块是否可以正常导入"""
    print("\n=== 测试模块导入 ===")
    
    modules_to_test = [
        "src",
        "src.main",
        "src.layers",
        "src.layers.voice_recognition",
        "src.layers.voice_recognition.voice_recognition",
        "src.layers.counterpoint_design",
        "src.layers.counterpoint_design.counterpoint_design",
        "src.layers.steady_execution",
        "src.layers.steady_execution.steady_execution",
        "src.layers.meta_protocol",
        "src.layers.meta_protocol.meta_protocol",
        "src.layers.consensus_crystal",
        "src.layers.consensus_crystal.consensus_crystal",
        "src.mechanisms",
        "src.utils"
    ]
    
    for module_name in modules_to_test:
        try:
            __import__(module_name)
            print(f"✓ 成功导入: {module_name}")
        except Exception as e:
            print(f"✗ 导入失败: {module_name} - {e}")
            return False
    
    return True

def test_voice_recognition():
    """测试声部识别模块"""
    print("\n=== 测试声部识别模块 ===")
    
    try:
        from src.layers.voice_recognition.voice_recognition import CollaborativeSonicMap
        
        # 创建协同声部图谱
        csm = CollaborativeSonicMap()
        
        # 注册碳基声部
        carbon_voice = csm.register_voice(
            name="测试用户",
            voice_type="carbon",
            capability_vector={"创意": 0.9, "审美": 0.8, "逻辑": 0.7},
            intention_vector={"探索": 0.8, "表达": 0.9},
            description="测试用碳基声部"
        )
        print(f"✓ 成功注册碳基声部: {carbon_voice.name}")
        
        # 注册硅基声部
        silicon_voice = csm.register_voice(
            name="测试AI",
            voice_type="silicon",
            capability_vector={"计算": 0.99, "记忆": 0.99, "逻辑": 0.98},
            intention_vector={"执行": 0.95, "优化": 0.9},
            description="测试用硅基声部"
        )
        print(f"✓ 成功注册硅基声部: {silicon_voice.name}")
        
        # 获取声部信息
        retrieved_voice = csm.get_voice(carbon_voice.voice_id)
        if retrieved_voice:
            print(f"✓ 成功获取声部信息: {retrieved_voice.name}")
        else:
            print("✗ 获取声部信息失败")
            return False
        
        # 按类型获取声部
        carbon_voices = csm.get_voices_by_type("carbon")
        silicon_voices = csm.get_voices_by_type("silicon")
        print(f"✓ 碳基声部数量: {len(carbon_voices)}")
        print(f"✓ 硅基声部数量: {len(silicon_voices)}")
        
        # 获取声部图谱
        voice_map = csm.get_voice_map()
        if voice_map:
            print("✓ 成功获取声部图谱")
        else:
            print("✗ 获取声部图谱失败")
            return False
        
        return True
        
    except Exception as e:
        print(f"✗ 声部识别模块测试失败: {e}")
        return False

def test_counterpoint_design():
    """测试对位设计模块"""
    print("\n=== 测试对位设计模块 ===")
    
    try:
        from src.layers.counterpoint_design.counterpoint_design import CounterpointDesigner
        
        # 创建协奏设计师
        designer = CounterpointDesigner()
        
        # 获取适合的模式
        suitable_patterns = designer.get_suitable_patterns("概念设计")
        print(f"✓ 获取到 {len(suitable_patterns)} 个适合概念设计的模式")
        
        # 创建协同路径
        path = designer.create_counterpoint_path(
            name="测试协同路径",
            pattern_type="staggered_complement",
            participating_voices=["voice1", "voice2"],
            creation_theme="测试创作主题"
        )
        print(f"✓ 成功创建协同路径: {path.name}")
        
        # 执行路径步骤
        result = designer.execute_path_step(
            path_id=path.path_id,
            step_index=0,
            voice_id="voice1",
            inputs={"concept": "测试概念"}
        )
        if result.get("success"):
            print("✓ 成功执行路径步骤")
        else:
            print(f"✗ 执行路径步骤失败: {result.get('error')}")
            return False
        
        # 验证协同路径
        valid, message = designer.validate_counterpoint_path(path)
        if valid:
            print(f"✓ 协同路径验证通过: {message}")
        else:
            print(f"✗ 协同路径验证失败: {message}")
            return False
        
        # 模拟执行协同路径
        simulation_results = designer.simulate_counterpoint_execution(path.path_id)
        print(f"✓ 成功模拟执行协同路径，共执行 {len(simulation_results)} 个步骤")
        
        return True
        
    except Exception as e:
        print(f"✗ 对位设计模块测试失败: {e}")
        return False

def test_steady_execution():
    """测试静定执行模块"""
    print("\n=== 测试静定执行模块 ===")
    
    try:
        from src.layers.steady_execution.steady_execution import SteadyExecutor
        
        # 创建静定执行器
        executor = SteadyExecutor()
        
        # 提交测试任务
        task_id = executor.submit_task(
            name="测试任务",
            task_type="test",
            payload={"test_data": "test_value"}
        )
        print(f"✓ 成功提交测试任务，任务ID: {task_id}")
        
        # 获取任务状态
        import time
        time.sleep(0.1)  # 等待任务执行完成
        
        task_status = executor.get_task_status(task_id)
        print(f"✓ 任务状态: {task_status.get('status')}")
        
        # 测试执行协同路径
        execution_result = executor.execute_counterpoint_path(
            path_id="test_path",
            steps=[
                {"action": "测试步骤1", "role": "carbon"},
                {"action": "测试步骤2", "role": "silicon"}
            ],
            voice_map={"carbon": "voice1", "silicon": "voice2"}
        )
        print(f"✓ 成功执行协同路径，执行ID: {execution_result.get('execution_id')}")
        
        # 获取执行统计信息
        stats = executor.get_execution_stats()
        print(f"✓ 执行统计信息: 队列大小={stats.get('queue_size')}, 活跃任务={stats.get('active_tasks')}")
        
        # 获取系统健康状态
        health = executor.get_system_health()
        print(f"✓ 系统健康状态: {health.get('status')}")
        
        # 关闭执行器
        executor.shutdown()
        print("✓ 成功关闭执行器")
        
        return True
        
    except Exception as e:
        print(f"✗ 静定执行模块测试失败: {e}")
        return False

def test_main():
    """测试主程序"""
    print("\n=== 测试主程序 ===")
    
    try:
        from src.main import main
        
        # 执行主函数
        main()
        print("✓ 主程序执行成功")
        
        return True
        
    except Exception as e:
        print(f"✗ 主程序测试失败: {e}")
        return False

def main():
    """主测试函数"""
    print("开始测试 Meta-CreationPower 项目...")
    
    # 运行所有测试
    tests = [
        test_imports,
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
        else:
            print(f"测试 {test.__name__} 失败")
    
    # 打印测试结果
    print(f"\n=== 测试结果 ===")
    print(f"总测试数: {total_tests}")
    print(f"通过测试数: {passed_tests}")
    print(f"失败测试数: {total_tests - passed_tests}")
    
    if passed_tests == total_tests:
        print("\n🎉 所有测试通过！Meta-CreationPower 项目可以正常运行。")
        return 0
    else:
        print("\n❌ 部分测试失败，需要修复项目中的问题。")
        return 1

if __name__ == "__main__":
    sys.exit(main())
