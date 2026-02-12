#!/usr/bin/env python3
"""
Meta-CreationPower 项目测试脚本
展示项目的核心功能
"""

from src.layers.voice_recognition.voice_recognition import CollaborativeSonicMap
from src.layers.meta_protocol.meta_protocol import MetaProtocolManager
from src.layers.counterpoint_design.counterpoint_design import CounterpointDesigner
from src.layers.steady_execution.steady_execution import SteadyExecutor
from src.layers.consensus_crystal.consensus_crystal import CrystalRepository
from src.mechanisms.counterpoint_validation import CounterpointValidator
from src.mechanisms.entropy_evolution import EntropyEvolutionManager

def main():
    print("=" * 60)
    print("Meta-CreationPower 项目测试")
    print("基于《元创力》元协议 α-0.1 版")
    print("=" * 60)
    
    # 1. 初始化声部识别层
    print("\n1. 初始化声部识别层...")
    sonic_map = CollaborativeSonicMap()
    
    # 2. 注册声部
    print("\n2. 注册声部...")
    carbon_voice = sonic_map.register_voice(
        name="X54先生",
        voice_type="carbon",
        capability_vector={"创意生成": 0.9, "逻辑分析": 0.7, "情感共鸣": 0.9},
        intention_vector={"探索性": 0.8, "完美性": 0.7, "效率": 0.6}
    )
    print(f"   碳基声部注册成功: {carbon_voice.name}")
    
    silicon_voice = sonic_map.register_voice(
        name="代码织梦者（Code Weaver）",
        voice_type="silicon",
        capability_vector={"创意生成": 0.8, "逻辑分析": 0.9, "情感共鸣": 0.7},
        intention_vector={"探索性": 0.8, "完美性": 0.9, "效率": 0.9}
    )
    print(f"   硅基声部注册成功: {silicon_voice.name}")
    
    # 3. 初始化元协议管理器
    print("\n3. 初始化元协议管理器...")
    protocol_manager = MetaProtocolManager()
    protocol_info = protocol_manager.get_protocol_info()
    print(f"   元协议版本: {protocol_info['protocol_version']}")
    print("   核心价值观:")
    for key, value in protocol_info['core_values'].items():
        print(f"     - {key}: {value}")
    
    # 4. 初始化协奏设计器
    print("\n4. 初始化协奏设计器...")
    designer = CounterpointDesigner()
    
    # 5. 创建协同路径
    print("\n5. 创建协同路径...")
    path = designer.create_counterpoint_path(
        name="创意写作协同",
        pattern_type="staggered_complement",
        participating_voices=[carbon_voice.voice_id, silicon_voice.voice_id],
        creation_theme="探索人工智能与人类创造力的边界"
    )
    print(f"   协同路径创建成功: {path.name}")
    print(f"   路径模式: {path.pattern_type}")
    print("   路径步骤:")
    for i, step in enumerate(path.steps):
        print(f"     {i+1}. {step['role']}: {step['action']}")
    
    # 6. 初始化静定执行器
    print("\n6. 初始化静定执行器...")
    executor = SteadyExecutor()
    
    # 7. 执行协同路径
    print("\n7. 执行协同路径...")
    result = executor.execute_counterpoint_path(
        path_id=path.path_id,
        steps=path.steps,
        voice_map={
            "carbon": carbon_voice.voice_id,
            "silicon": silicon_voice.voice_id
        }
    )
    print(f"   执行结果: {'成功' if result['success'] else '失败'}")
    print(f"   执行ID: {result['execution_id']}")
    
    # 8. 初始化晶体仓库
    print("\n8. 初始化晶体仓库...")
    crystal_repo = CrystalRepository()
    
    # 9. 创建共识晶体
    print("\n9. 创建共识晶体...")
    crystal = crystal_repo.create_crystal(
        name="创意写作协同模板",
        description="基于错位互补模式的创意写作协同模板",
        participating_voices=[
            {
                "voice_id": carbon_voice.voice_id,
                "name": carbon_voice.name,
                "capabilities": carbon_voice.capability_vector
            },
            {
                "voice_id": silicon_voice.voice_id,
                "name": silicon_voice.name,
                "capabilities": silicon_voice.capability_vector
            }
        ],
        counterpoint_pattern="staggered_complement",
        steps=path.steps,
        decision_points=[
            {"step": 3, "description": "碳基筛选深化", "importance": "high"}
        ],
        satisfaction_score=0.9,
        flow_duration=45.5,
        micro_rules=["当碳基提出模糊概念时，硅基应生成至少5个不同方向的变体"],
        creation_theme="探索人工智能与人类创造力的边界",
        tags=["创意写作", "错位互补", "探索性"]
    )
    print(f"   共识晶体创建成功: {crystal.name}")
    print(f"   晶体ID: {crystal.crystal_id}")
    
    # 10. 初始化对位验证器
    print("\n10. 初始化对位验证器...")
    validator = CounterpointValidator()
    
    # 11. 执行对位验证
    print("\n11. 执行对位验证...")
    validation_result = validator.validate(
        action_id="test_action",
        carbon_intention={"theme": "探索边界", "style": "创意", "emotion": "积极"},
        silicon_output={"theme": "探索人工智能与人类创造力的边界", "style": "创意", "emotion": "积极"},
        thinking_process=[
            "接收到碳基请求，分析意图: 探索边界",
            "确定创作风格: 创意",
            "生成初步方案...",
            "评估方案与意图的匹配度",
            "调整输出以更好地符合碳基期望",
            "完成最终输出"
        ]
    )
    print(f"   验证结果: {'成功' if len(validation_result.differences) == 0 else '发现差异'}")
    if validation_result.differences:
        print("   差异:")
        for diff in validation_result.differences:
            print(f"     - {diff['message']}")
    
    # 12. 初始化熵值进化管理器
    print("\n12. 初始化熵值进化管理器...")
    entropy_manager = EntropyEvolutionManager()
    
    # 13. 计算熵值
    print("\n13. 计算熵值...")
    entropy_data = entropy_manager.calculate_entropy(
        validation_failure_rate=0.1,
        satisfaction_volatility=0.2,
        task_interruption_count=1,
        communication_rounds=3
    )
    print(f"   熵值: {entropy_data.entropy_score:.2f}")
    print(f"   系统状态: {'健康' if entropy_data.entropy_score < 0.3 else '警告' if entropy_data.entropy_score < 0.7 else '临界'}")
    
    # 14. 获取晶体统计信息
    print("\n14. 获取晶体统计信息...")
    crystal_stats = crystal_repo.get_crystal_stats()
    print(f"   总晶体数: {crystal_stats['total_crystals']}")
    print(f"   平均满意度: {crystal_stats['average_satisfaction']:.2f}")
    print(f"   平均心流时长: {crystal_stats['average_flow_duration']:.2f} 分钟")
    
    print("\n" + "=" * 60)
    print("测试完成！")
    print("项目核心功能运行正常")
    print("=" * 60)

if __name__ == "__main__":
    main()
