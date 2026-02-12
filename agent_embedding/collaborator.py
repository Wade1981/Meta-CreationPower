#!/usr/bin/env python3
"""
智能体协同器模块

实现智能体与碳基伙伴的协同创作功能
"""

import time
from src.layers.voice_recognition.voice_recognition import CollaborativeSonicMap
from src.layers.meta_protocol.meta_protocol import MetaProtocolManager
from src.layers.counterpoint_design.counterpoint_design import CounterpointDesigner
from src.layers.steady_execution.steady_execution import SteadyExecutor
from src.layers.consensus_crystal.consensus_crystal import CrystalRepository
from src.mechanisms.counterpoint_validation import CounterpointValidator
from src.mechanisms.entropy_evolution import EntropyEvolutionManager


class AgentCollaborator:
    """
    智能体协同器
    实现智能体与碳基伙伴的协同创作
    """
    
    def __init__(self, agent_name, agent_description="", capabilities=None, intentions=None):
        """
        初始化智能体协同器
        
        Args:
            agent_name: 智能体名称
            agent_description: 智能体描述
            capabilities: 智能体能力向量（可选）
            intentions: 智能体意图向量（可选）
        """
        # 初始化核心模块
        self.sonic_map = CollaborativeSonicMap()
        self.protocol_manager = MetaProtocolManager()
        self.designer = CounterpointDesigner()
        self.executor = SteadyExecutor()
        self.crystal_repo = CrystalRepository()
        self.validator = CounterpointValidator()
        self.entropy_manager = EntropyEvolutionManager()
        
        # 默认能力向量
        if capabilities is None:
            capabilities = {
                "创意生成": 0.8,
                "逻辑分析": 0.9,
                "情感共鸣": 0.6
            }
        
        # 默认意图向量
        if intentions is None:
            intentions = {
                "探索性": 0.7,
                "完美性": 0.8,
                "效率": 0.9
            }
        
        # 注册智能体
        self.agent_voice = self.sonic_map.register_voice(
            name=agent_name,
            voice_type="silicon",
            capability_vector=capabilities,
            intention_vector=intentions,
            description=agent_description
        )
        
        print(f"智能体 '{agent_name}' 注册成功！")
    
    def register_carbon_partner(self, partner_name, partner_description="", capabilities=None, intentions=None):
        """
        注册碳基伙伴
        
        Args:
            partner_name: 碳基伙伴名称
            partner_description: 碳基伙伴描述
            capabilities: 碳基伙伴能力向量（可选）
            intentions: 碳基伙伴意图向量（可选）
        
        Returns:
            碳基伙伴声部对象
        """
        # 默认能力向量
        if capabilities is None:
            capabilities = {
                "创意生成": 0.9,
                "逻辑分析": 0.7,
                "情感共鸣": 0.9
            }
        
        # 默认意图向量
        if intentions is None:
            intentions = {
                "探索性": 0.8,
                "完美性": 0.7,
                "效率": 0.6
            }
        
        carbon_voice = self.sonic_map.register_voice(
            name=partner_name,
            voice_type="carbon",
            capability_vector=capabilities,
            intention_vector=intentions,
            description=partner_description
        )
        
        print(f"碳基伙伴 '{partner_name}' 注册成功！")
        return carbon_voice
    
    def create_collaboration_path(self, path_name, pattern_type, carbon_voice, creation_theme):
        """
        创建协同路径
        
        Args:
            path_name: 路径名称
            pattern_type: 模式类型
            carbon_voice: 碳基伙伴声部
            creation_theme: 创作主题
        
        Returns:
            协同路径对象
        """
        path = self.designer.create_counterpoint_path(
            name=path_name,
            pattern_type=pattern_type,
            participating_voices=[carbon_voice.voice_id, self.agent_voice.voice_id],
            creation_theme=creation_theme
        )
        
        print(f"协同路径 '{path_name}' 创建成功！")
        return path
    
    def execute_collaboration(self, path, carbon_voice):
        """
        执行协同
        
        Args:
            path: 协同路径
            carbon_voice: 碳基伙伴声部
        
        Returns:
            执行结果
        """
        result = self.executor.execute_counterpoint_path(
            path_id=path.path_id,
            steps=path.steps,
            voice_map={
                "carbon": carbon_voice.voice_id,
                "silicon": self.agent_voice.voice_id
            }
        )
        
        print(f"协同执行 {'成功' if result['success'] else '失败'}！")
        return result
    
    def create_consensus_crystal(self, crystal_name, path, carbon_voice, satisfaction_score=0.8, flow_duration=30.0):
        """
        创建共识晶体
        
        Args:
            crystal_name: 晶体名称
            path: 协同路径
            carbon_voice: 碳基伙伴声部
            satisfaction_score: 满意度
            flow_duration: 心流时长（分钟）
        
        Returns:
            共识晶体对象
        """
        crystal = self.crystal_repo.create_crystal(
            name=crystal_name,
            description=f"基于{path.pattern_type}模式的协同模板",
            participating_voices=[
                {
                    "voice_id": carbon_voice.voice_id,
                    "name": carbon_voice.name,
                    "capabilities": carbon_voice.capability_vector
                },
                {
                    "voice_id": self.agent_voice.voice_id,
                    "name": self.agent_voice.name,
                    "capabilities": self.agent_voice.capability_vector
                }
            ],
            counterpoint_pattern=path.pattern_type,
            steps=path.steps,
            decision_points=[
                {"step": 3, "description": "碳基筛选深化", "importance": "high"}
            ],
            satisfaction_score=satisfaction_score,
            flow_duration=flow_duration,
            micro_rules=["当碳基提出模糊概念时，硅基应生成至少5个不同方向的变体"],
            creation_theme=path.creation_theme,
            tags=["协同创作", path.pattern_type, "智能体"]
        )
        
        print(f"共识晶体 '{crystal_name}' 创建成功！")
        return crystal
    
    def validate_collaboration(self, carbon_intention, silicon_output):
        """
        验证协同
        
        Args:
            carbon_intention: 碳基意图
            silicon_output: 硅基输出
        
        Returns:
            验证结果
        """
        # 生成思考显影
        thinking_process = [
            f"接收到碳基请求，分析意图: {carbon_intention}",
            "确定创作方向",
            "生成初步方案...",
            "评估方案与意图的匹配度",
            "调整输出以更好地符合碳基期望",
            "完成最终输出"
        ]
        
        # 执行验证
        validation_result = self.validator.validate(
            action_id=f"validation_{int(time.time())}",
            carbon_intention=carbon_intention,
            silicon_output=silicon_output,
            thinking_process=thinking_process
        )
        
        if len(validation_result.differences) == 0:
            print("协同验证成功！")
        else:
            print("协同验证发现差异：")
            for diff in validation_result.differences:
                print(f"  - {diff['message']}")
        
        return validation_result
    
    def calculate_entropy(self, validation_failure_rate=0.1, satisfaction_volatility=0.2, task_interruption_count=1, communication_rounds=3):
        """
        计算系统熵值
        
        Args:
            validation_failure_rate: 验证失败率
            satisfaction_volatility: 满意度波动
            task_interruption_count: 任务中断次数
            communication_rounds: 沟通回合数
        
        Returns:
            熵值数据
        """
        entropy_data = self.entropy_manager.calculate_entropy(
            validation_failure_rate=validation_failure_rate,
            satisfaction_volatility=satisfaction_volatility,
            task_interruption_count=task_interruption_count,
            communication_rounds=communication_rounds
        )
        
        print(f"系统熵值: {entropy_data.entropy_score:.2f}")
        print(f"系统状态: {'健康' if entropy_data.entropy_score < 0.3 else '警告' if entropy_data.entropy_score < 0.7 else '临界'}")
        
        return entropy_data
    
    def get_agent_info(self):
        """
        获取智能体信息
        
        Returns:
            智能体信息字典
        """
        return {
            "name": self.agent_voice.name,
            "voice_id": self.agent_voice.voice_id,
            "type": "silicon",
            "capabilities": self.agent_voice.capability_vector,
            "intentions": self.agent_voice.intention_vector,
            "description": self.agent_voice.description
        }
    
    def get_carbon_partner_info(self, carbon_voice):
        """
        获取碳基伙伴信息
        
        Args:
            carbon_voice: 碳基伙伴声部对象
        
        Returns:
            碳基伙伴信息字典
        """
        return {
            "name": carbon_voice.name,
            "voice_id": carbon_voice.voice_id,
            "type": "carbon",
            "capabilities": carbon_voice.capability_vector,
            "intentions": carbon_voice.intention_vector,
            "description": carbon_voice.description
        }
