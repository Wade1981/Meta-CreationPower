#!/usr/bin/env python3
"""
智能体嵌入包使用示例

展示如何使用 Meta-CreationPower 智能体嵌入包实现与碳基伙伴的协同
"""

from agent_embedding import EmbeddingAPI
from agent_embedding.utils import setup_logger


def basic_example():
    """
    基本使用示例
    """
    print("=" * 80)
    print("智能体嵌入包 - 基本使用示例")
    print("=" * 80)
    
    # 初始化嵌入API
    print("\n1. 初始化智能体嵌入API...")
    api = EmbeddingAPI(
        agent_name="创意助手智能体",
        agent_description="专注于创意生成和内容创作的智能体"
    )
    
    # 注册碳基伙伴
    print("\n2. 注册碳基伙伴...")
    partner_id = api.register_carbon_partner(
        partner_name="X54先生",
        partner_description="碳基人文科技践行者"
    )
    print(f"碳基伙伴注册成功！ID: {partner_id}")
    
    # 创建错位互补模式的协同
    print("\n3. 创建错位互补模式的协同...")
    result = api.create_staggered_complement_collaboration(
        partner_name="X54先生",
        collaboration_name="创意写作协同",
        creation_theme="探索人工智能与人类创造力的边界"
    )
    print(f"协同执行结果: {'成功' if result['success'] else '失败'}")
    
    # 计算系统健康状态
    print("\n4. 计算系统健康状态...")
    health_status = api.calculate_system_health()
    print(f"系统熵值: {health_status['entropy_score']:.2f}")
    print(f"健康状态: {health_status['health_level']}")
    print("建议:")
    for recommendation in health_status['recommendations']:
        print(f"  - {recommendation}")
    
    # 获取智能体信息
    print("\n5. 获取智能体信息...")
    agent_info = api.get_agent_info()
    print(f"智能体名称: {agent_info['name']}")
    print(f"智能体类型: {agent_info['type']}")
    print("智能体能力:")
    for capability, value in agent_info['capabilities'].items():
        print(f"  - {capability}: {value}")
    
    print("\n" + "=" * 80)
    print("基本使用示例完成！")
    print("=" * 80)


def quick_start_example():
    """
    快速启动示例
    """
    print("\n" + "=" * 80)
    print("智能体嵌入包 - 快速启动示例")
    print("=" * 80)
    
    # 初始化嵌入API
    api = EmbeddingAPI(
        agent_name="快速助手智能体",
        agent_description="快速启动协同的智能体"
    )
    
    # 快速启动协同
    print("\n1. 快速启动协同...")
    result = api.quick_start_collaboration(
        partner_name="创意总监",
        theme="未来科技与人类生活",
        collaboration_type="staggered_complement"
    )
    print(f"快速协同执行结果: {'成功' if result['success'] else '失败'}")
    
    # 获取所有碳基伙伴
    print("\n2. 获取所有碳基伙伴...")
    partners = api.get_carbon_partners()
    print(f"已注册碳基伙伴数量: {len(partners)}")
    for partner in partners:
        print(f"  - {partner['name']} (ID: {partner['voice_id']})")
    
    print("\n" + "=" * 80)
    print("快速启动示例完成！")
    print("=" * 80)


def advanced_example():
    """
    高级使用示例
    """
    print("\n" + "=" * 80)
    print("智能体嵌入包 - 高级使用示例")
    print("=" * 80)
    
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
    print("\n1. 初始化高级智能体嵌入API...")
    api = EmbeddingAPI(
        agent_name="高级创意智能体",
        agent_description="具备高级创意能力的智能体",
        capabilities=custom_capabilities,
        intentions=custom_intentions
    )
    
    # 注册具有自定义能力的碳基伙伴
    print("\n2. 注册具有自定义能力的碳基伙伴...")
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
        partner_name="艺术总监",
        partner_description="具有丰富艺术经验的碳基伙伴",
        capabilities=partner_capabilities,
        intentions=partner_intentions
    )
    print(f"高级碳基伙伴注册成功！ID: {partner_id}")
    
    # 测试不同类型的协同
    print("\n3. 测试不同类型的协同...")
    
    # 测试错位互补模式
    print("\n3.1 测试错位互补模式...")
    result1 = api.create_staggered_complement_collaboration(
        partner_name="艺术总监",
        collaboration_name="艺术创意协同",
        creation_theme="数字艺术与传统美学的融合"
    )
    
    # 测试卡农式推进模式
    print("\n3.2 测试卡农式推进模式...")
    result2 = api.create_canon_progression_collaboration(
        partner_name="艺术总监",
        collaboration_name="流程化创作协同",
        creation_theme="多媒体内容的流程化创作"
    )
    
    # 测试赋格式交织模式
    print("\n3.3 测试赋格式交织模式...")
    result3 = api.create_fugue_interweaving_collaboration(
        partner_name="艺术总监",
        collaboration_name="复杂创意协同",
        creation_theme="多维度创意的交织融合"
    )
    
    print("\n协同测试结果:")
    print(f"  错位互补模式: {'成功' if result1['success'] else '失败'}")
    print(f"  卡农式推进模式: {'成功' if result2['success'] else '失败'}")
    print(f"  赋格式交织模式: {'成功' if result3['success'] else '失败'}")
    
    print("\n" + "=" * 80)
    print("高级使用示例完成！")
    print("=" * 80)


def validation_example():
    """
    协同验证示例
    """
    print("\n" + "=" * 80)
    print("智能体嵌入包 - 协同验证示例")
    print("=" * 80)
    
    # 初始化嵌入API
    api = EmbeddingAPI(
        agent_name="验证助手智能体",
        agent_description="专注于协同验证的智能体"
    )
    
    # 模拟碳基意图和硅基输出
    carbon_intention = {
        "theme": "环境保护与可持续发展",
        "style": "严肃专业",
        "emotion": "积极",
        "requirements": ["科学准确", "数据支持", "可行性建议"]
    }
    
    silicon_output = {
        "theme": "环境保护与可持续发展",
        "style": "严肃专业",
        "emotion": "积极",
        "content": "基于最新研究数据，我们提出以下可持续发展策略...",
        "recommendations": ["推广可再生能源", "减少碳排放", "加强环境教育"]
    }
    
    # 验证协同
    print("\n1. 验证协同...")
    validation_result = api.validate_collaboration(
        carbon_intention=carbon_intention,
        silicon_output=silicon_output
    )
    
    print("\n验证结果详情:")
    print(f"  验证ID: {validation_result.validation_id}")
    print(f"  差异数量: {len(validation_result.differences)}")
    if validation_result.differences:
        print("  差异详情:")
        for diff in validation_result.differences:
            print(f"    - {diff['message']}")
    else:
        print("  验证通过，无差异！")
    
    print("\n" + "=" * 80)
    print("协同验证示例完成！")
    print("=" * 80)


if __name__ == "__main__":
    # 运行所有示例
    basic_example()
    quick_start_example()
    advanced_example()
    validation_example()
    
    print("\n" + "=" * 80)
    print("所有示例运行完成！")
    print("=" * 80)
