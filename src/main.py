
"""
Meta-CreationPower 主程序入口
基于"和、清、寂、静"四字原则构建的人机协作框架
"""

from src.layers.voice_recognition.voice_recognition import CollaborativeSonicMap
from src.layers.meta_protocol.meta_protocol import MetaProtocolManager
from src.layers.counterpoint_design.counterpoint_design import CounterpointDesigner
from src.layers.steady_execution.steady_execution import SteadyExecutor
from src.layers.consensus_crystal.consensus_crystal import CrystalRepository
from src.mechanisms.counterpoint_validation import CounterpointValidator
from src.mechanisms.entropy_evolution import EntropyEvolutionManager
import time


def main():
    print("欢迎使用 Meta-CreationPower 项目")
    print("基于'和、清、寂、静'四字原则构建的人机协作框架")
    print("=" * 60)
    
    # 1. 初始化声部识别层
    print("\n1. 初始化声部识别层...")
    sonic_map = CollaborativeSonicMap()
    
    # 2. 注册声部
    print("2. 注册碳基与硅基声部...")
    carbon_voice = sonic_map.register_voice(
        name="X54先生",
        voice_type="carbon",
        capability_vector={"创意生成": 0.9, "逻辑分析": 0.7, "情感共鸣": 0.9},
        intention_vector={"探索性": 0.8, "完美性": 0.7, "效率": 0.6},
        description="碳基人文科技践行者，负责思维锚点、架构设计、价值定调"
    )
    
    silicon_voice = sonic_map.register_voice(
        name="代码织梦者",
        voice_type="silicon",
        capability_vector={"创意生成": 0.8, "逻辑分析": 0.9, "情感共鸣": 0.5},
        intention_vector={"探索性": 0.7, "完美性": 0.8, "效率": 0.9},
        description="硅基智能代码编织者，基于'和清寂静'启蒙灯塔起源团队碳硅协同内核，以碳硅协同对位法协同结合对话式思维编程，为碳基伙伴提供专业的代码编织服务"
    )
    
    print(f"   - 碳基声部: {carbon_voice.name}")
    print(f"   - 硅基声部: {silicon_voice.name}")
    
    # 3. 初始化元协议管理器
    print("\n3. 初始化元协议管理器...")
    meta_protocol = MetaProtocolManager()
    protocol_info = meta_protocol.get_protocol_info()
    print(f"   - 元协议版本: {protocol_info['protocol_version']}")
    print("   - 核心价值观:")
    for key, value in protocol_info['core_values'].items():
        print(f"     {key}: {value}")
    
    # 4. 初始化协奏设计器
    print("\n4. 初始化协奏设计器...")
    designer = CounterpointDesigner()
    
    # 5. 创建协同路径
    print("5. 创建协同路径...")
    path = designer.create_counterpoint_path(
        name="创意写作协同",
        pattern_type="staggered_complement",
        participating_voices=[carbon_voice.voice_id, silicon_voice.voice_id],
        creation_theme="探索人工智能与人类创造力的边界"
    )
    print(f"   - 路径名称: {path.name}")
    print(f"   - 模式类型: {path.pattern_type}")
    print(f"   - 创作主题: {path.creation_theme}")
    print("   - 协同步骤:")
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
    
    print(f"   - 执行ID: {result['execution_id']}")
    print(f"   - 执行状态: {'成功' if result['success'] else '失败'}")
    print("   - 执行结果:")
    for i, task_result in enumerate(result['task_results']):
        status = "完成" if task_result['status'] == "completed" else "失败"
        print(f"     步骤 {i+1}: {status}")
    
    # 8. 初始化共识晶体仓库
    print("\n8. 初始化共识晶体仓库...")
    crystal_repo = CrystalRepository()
    
    # 9. 生成共识晶体
    print("9. 生成共识晶体...")
    crystal = crystal_repo.create_crystal(
        name="创意写作协同晶体",
        description="基于错位互补模式的创意写作协同经验",
        participating_voices=[
            {
                "name": carbon_voice.name,
                "type": carbon_voice.voice_type,
                "capabilities": carbon_voice.capability_vector
            },
            {
                "name": silicon_voice.name,
                "type": silicon_voice.voice_type,
                "capabilities": silicon_voice.capability_vector
            }
        ],
        counterpoint_pattern=path.pattern_type,
        steps=path.steps,
        decision_points=[
            {
                "step": 3,
                "description": "碳基筛选深化环节",
                "importance": "high"
            }
        ],
        satisfaction_score=0.95,
        flow_duration=45.0,
        micro_rules=[
            "碳基提出的模糊概念越具体，硅基生成的变体质量越高",
            "硅基生成的变体数量应控制在合理范围内，避免信息过载",
            "协同过程中应保持定期的反馈循环"
        ],
        creation_theme=path.creation_theme,
        tags=["创意写作", "错位互补", "碳硅协同"]
    )
    
    print(f"   - 晶体ID: {crystal.crystal_id}")
    print(f"   - 晶体名称: {crystal.name}")
    print(f"   - 满意度: {crystal.satisfaction_score}")
    print(f"   - 心流时长: {crystal.flow_duration}分钟")
    
    # 10. 显示系统状态
    print("\n10. 系统状态:")
    print(f"   - 声部数量: {len(sonic_map.get_voice_map()['voices'])}")
    print(f"   - 协同路径: {len(designer.get_counterpoint_paths())}")
    print(f"   - 共识晶体: {len(crystal_repo.get_all_crystals())}")
    print(f"   - 执行器健康状态: {executor.get_system_health()['status']}")
    
    # 11. 关闭执行器
    print("\n11. 关闭系统...")
    executor.shutdown()
    
    print("\n" + "=" * 60)
    print("Meta-CreationPower 协同流程执行完成！")
    print("基于'和、清、寂、静'四字原则的碳硅协同实践成功！")
    

if __name__ == "__main__":
    main()
