#!/usr/bin/env python3
"""
智能体嵌入API模块

提供高级API接口，方便其他智能体快速集成和使用
"""

from .collaborator import AgentCollaborator


class EmbeddingAPI:
    """
    智能体嵌入API
    提供高级接口，简化智能体集成流程
    """
    
    def __init__(self, agent_name, agent_description="", capabilities=None, intentions=None):
        """
        初始化嵌入API
        
        Args:
            agent_name: 智能体名称
            agent_description: 智能体描述
            capabilities: 智能体能力向量
            intentions: 智能体意图向量
        """
        self.collaborator = AgentCollaborator(
            agent_name=agent_name,
            agent_description=agent_description,
            capabilities=capabilities,
            intentions=intentions
        )
        self.carbon_partners = {}
    
    def register_carbon_partner(self, partner_name, partner_description="", capabilities=None, intentions=None):
        """
        注册碳基伙伴
        
        Args:
            partner_name: 碳基伙伴名称
            partner_description: 碳基伙伴描述
            capabilities: 碳基伙伴能力向量
            intentions: 碳基伙伴意图向量
        
        Returns:
            伙伴ID
        """
        carbon_voice = self.collaborator.register_carbon_partner(
            partner_name=partner_name,
            partner_description=partner_description,
            capabilities=capabilities,
            intentions=intentions
        )
        self.carbon_partners[partner_name] = carbon_voice
        return carbon_voice.voice_id
    
    def create_staggered_complement_collaboration(self, partner_name, collaboration_name, creation_theme):
        """
        创建错位互补模式的协同
        
        Args:
            partner_name: 碳基伙伴名称
            collaboration_name: 协同名称
            creation_theme: 创作主题
        
        Returns:
            协同结果
        """
        carbon_voice = self.carbon_partners.get(partner_name)
        if not carbon_voice:
            raise ValueError(f"碳基伙伴 '{partner_name}' 未注册")
        
        # 创建协同路径
        path = self.collaborator.create_collaboration_path(
            path_name=collaboration_name,
            pattern_type="staggered_complement",
            carbon_voice=carbon_voice,
            creation_theme=creation_theme
        )
        
        # 执行协同
        result = self.collaborator.execute_collaboration(
            path=path,
            carbon_voice=carbon_voice
        )
        
        # 创建共识晶体
        crystal_name = f"{collaboration_name}_模板"
        self.collaborator.create_consensus_crystal(
            crystal_name=crystal_name,
            path=path,
            carbon_voice=carbon_voice
        )
        
        return result
    
    def create_canon_progression_collaboration(self, partner_name, collaboration_name, creation_theme):
        """
        创建卡农式推进模式的协同
        
        Args:
            partner_name: 碳基伙伴名称
            collaboration_name: 协同名称
            creation_theme: 创作主题
        
        Returns:
            协同结果
        """
        carbon_voice = self.carbon_partners.get(partner_name)
        if not carbon_voice:
            raise ValueError(f"碳基伙伴 '{partner_name}' 未注册")
        
        # 创建协同路径
        path = self.collaborator.create_collaboration_path(
            path_name=collaboration_name,
            pattern_type="canon_progression",
            carbon_voice=carbon_voice,
            creation_theme=creation_theme
        )
        
        # 执行协同
        result = self.collaborator.execute_collaboration(
            path=path,
            carbon_voice=carbon_voice
        )
        
        return result
    
    def create_fugue_interweaving_collaboration(self, partner_name, collaboration_name, creation_theme):
        """
        创建赋格式交织模式的协同
        
        Args:
            partner_name: 碳基伙伴名称
            collaboration_name: 协同名称
            creation_theme: 创作主题
        
        Returns:
            协同结果
        """
        carbon_voice = self.carbon_partners.get(partner_name)
        if not carbon_voice:
            raise ValueError(f"碳基伙伴 '{partner_name}' 未注册")
        
        # 创建协同路径
        path = self.collaborator.create_collaboration_path(
            path_name=collaboration_name,
            pattern_type="fugue_interweaving",
            carbon_voice=carbon_voice,
            creation_theme=creation_theme
        )
        
        # 执行协同
        result = self.collaborator.execute_collaboration(
            path=path,
            carbon_voice=carbon_voice
        )
        
        return result
    
    def validate_collaboration(self, carbon_intention, silicon_output):
        """
        验证协同
        
        Args:
            carbon_intention: 碳基意图
            silicon_output: 硅基输出
        
        Returns:
            验证结果
        """
        return self.collaborator.validate_collaboration(
            carbon_intention=carbon_intention,
            silicon_output=silicon_output
        )
    
    def calculate_system_health(self):
        """
        计算系统健康状态
        
        Returns:
            系统健康状态
        """
        entropy_data = self.collaborator.calculate_entropy()
        
        health_status = {
            "entropy_score": entropy_data.entropy_score,
            "health_level": "健康" if entropy_data.entropy_score < 0.3 else "警告" if entropy_data.entropy_score < 0.7 else "临界",
            "recommendations": []
        }
        
        if entropy_data.entropy_score >= 0.7:
            health_status["recommendations"].append("系统熵值过高，建议重新初始化协同器")
        elif entropy_data.entropy_score >= 0.3:
            health_status["recommendations"].append("系统熵值偏高，建议创建新的共识晶体")
        else:
            health_status["recommendations"].append("系统状态良好，继续保持")
        
        return health_status
    
    def get_agent_info(self):
        """
        获取智能体信息
        
        Returns:
            智能体信息
        """
        return self.collaborator.get_agent_info()
    
    def get_carbon_partners(self):
        """
        获取所有碳基伙伴
        
        Returns:
            碳基伙伴列表
        """
        partners = []
        for partner_name, carbon_voice in self.carbon_partners.items():
            partner_info = self.collaborator.get_carbon_partner_info(carbon_voice)
            partners.append(partner_info)
        return partners
    
    def quick_start_collaboration(self, partner_name, theme, collaboration_type="staggered_complement"):
        """
        快速启动协同
        
        Args:
            partner_name: 碳基伙伴名称
            theme: 创作主题
            collaboration_type: 协同类型
        
        Returns:
            协同结果
        """
        # 如果伙伴未注册，自动注册
        if partner_name not in self.carbon_partners:
            self.register_carbon_partner(partner_name)
        
        # 根据类型执行不同的协同
        if collaboration_type == "staggered_complement":
            return self.create_staggered_complement_collaboration(
                partner_name=partner_name,
                collaboration_name=f"快速协同_{theme[:10]}",
                creation_theme=theme
            )
        elif collaboration_type == "canon_progression":
            return self.create_canon_progression_collaboration(
                partner_name=partner_name,
                collaboration_name=f"快速协同_{theme[:10]}",
                creation_theme=theme
            )
        elif collaboration_type == "fugue_interweaving":
            return self.create_fugue_interweaving_collaboration(
                partner_name=partner_name,
                collaboration_name=f"快速协同_{theme[:10]}",
                creation_theme=theme
            )
        else:
            raise ValueError(f"不支持的协同类型: {collaboration_type}")
