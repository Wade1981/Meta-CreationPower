"""
熵值驱动协议进化 (Entropy Evolution)
功能：基于系统熵值自动触发协议升级的机制
"""

from dataclasses import dataclass
from typing import Dict, List, Optional, Tuple, Any
import json
import uuid
import time


@dataclass
class EntropyData:
    """
    熵值数据基类
    """
    entropy_id: str
    validation_failure_rate: float  # 对位验证的失败率
    satisfaction_volatility: float  # 人类满意度波动
    task_interruption_count: int  # 协同任务的中断次数
    communication_rounds: int  # 碳硅沟通的回合数
    entropy_score: float  # 计算的熵值（0-1）
    timestamp: float


@dataclass
class EvolutionProposal:
    """
    进化提案基类
    """
    proposal_id: str
    chaos_cause: str  # 混乱主因
    micro_rule: str  # 微小补丁规则
    priority: str  # 优先级：low, medium, high
    description: str  # 提案描述
    estimated_impact: float  # 预计影响（0-1）
    voting_status: str  # 投票状态：pending, passed, rejected
    carbon_vote: Optional[bool]  # 碳基投票
    silicon_vote: Optional[bool]  # 硅基投票
    timestamp: float


class EntropyEvolutionManager:
    """
    熵值进化管理器
    实现基于熵值的协议自动升级
    """
    
    def __init__(self):
        """
        初始化熵值进化管理器
        """
        self.entropy_threshold = 0.7  # 熵值阈值
        self.entropy_history: List[EntropyData] = []
        self.evolution_proposals: Dict[str, EvolutionProposal] = {}
        self.recent_entropy_window = 5  # 最近5次熵值计算作为窗口
        self.micro_rules: List[str] = []  # 已生效的微规则
    
    def calculate_entropy(self, 
                         validation_failure_rate: float, 
                         satisfaction_volatility: float, 
                         task_interruption_count: int, 
                         communication_rounds: int) -> EntropyData:
        """
        计算系统熵值
        
        Args:
            validation_failure_rate: 对位验证的失败率
            satisfaction_volatility: 人类满意度波动
            task_interruption_count: 协同任务的中断次数
            communication_rounds: 碳硅沟通的回合数
        
        Returns:
            熵值数据对象
        """
        entropy_id = str(uuid.uuid4())
        
        # 标准化各项指标
        # 验证失败率（0-1）已经标准化
        
        # 满意度波动（0-1）
        normalized_volatility = min(satisfaction_volatility, 1.0)
        
        # 任务中断次数（标准化为0-1）
        normalized_interruptions = min(task_interruption_count / 10, 1.0)  # 假设10次中断为最大值
        
        # 沟通回合数（标准化为0-1）
        normalized_rounds = min(communication_rounds / 20, 1.0)  # 假设20回合为最大值
        
        # 加权计算熵值
        weights = {
            "validation_failure": 0.4,
            "satisfaction_volatility": 0.3,
            "task_interruptions": 0.2,
            "communication_rounds": 0.1
        }
        
        entropy_score = (
            validation_failure_rate * weights["validation_failure"] +
            normalized_volatility * weights["satisfaction_volatility"] +
            normalized_interruptions * weights["task_interruptions"] +
            normalized_rounds * weights["communication_rounds"]
        )
        
        entropy_data = EntropyData(
            entropy_id=entropy_id,
            validation_failure_rate=validation_failure_rate,
            satisfaction_volatility=satisfaction_volatility,
            task_interruption_count=task_interruption_count,
            communication_rounds=communication_rounds,
            entropy_score=entropy_score,
            timestamp=time.time()
        )
        
        self.entropy_history.append(entropy_data)
        
        # 保持历史记录在合理范围内
        if len(self.entropy_history) > 100:
            self.entropy_history = self.entropy_history[-100:]
        
        return entropy_data
    
    def check_evolution_trigger(self) -> bool:
        """
        检查是否触发进化
        
        Returns:
            是否触发进化
        """
        if len(self.entropy_history) < self.recent_entropy_window:
            return False
        
        # 计算最近窗口内的平均熵值
        recent_entropies = self.entropy_history[-self.recent_entropy_window:]
        average_entropy = sum(data.entropy_score for data in recent_entropies) / len(recent_entropies)
        
        return average_entropy > self.entropy_threshold
    
    def analyze_chaos_cause(self) -> str:
        """
        分析混乱主因
        
        Returns:
            混乱主因
        """
        if not self.entropy_history:
            return "insufficient_data"
        
        # 分析最近的熵值数据，找出贡献最大的因素
        recent_data = self.entropy_history[-1]
        
        factors = [
            ("validation_failure", recent_data.validation_failure_rate * 0.4),
            ("satisfaction_volatility", recent_data.satisfaction_volatility * 0.3),
            ("task_interruptions", min(recent_data.task_interruption_count / 10, 1.0) * 0.2),
            ("communication_rounds", min(recent_data.communication_rounds / 20, 1.0) * 0.1)
        ]
        
        # 找出贡献最大的因素
        factors.sort(key=lambda x: x[1], reverse=True)
        main_cause = factors[0][0]
        
        return main_cause
    
    def generate_evolution_proposal(self) -> EvolutionProposal:
        """
        生成进化提案
        
        Returns:
            进化提案对象
        """
        proposal_id = str(uuid.uuid4())
        chaos_cause = self.analyze_chaos_cause()
        
        # 基于混乱主因生成微规则
        micro_rule_templates = {
            "validation_failure": {
                "rule": "当对位验证失败率超过30%时，系统应自动暂停并请求碳基澄清意图",
                "priority": "high",
                "description": "验证失败率过高表明碳硅之间存在严重的理解偏差"
            },
            "satisfaction_volatility": {
                "rule": "当人类满意度波动超过20%时，系统应提供更详细的选项和解释",
                "priority": "medium",
                "description": "满意度波动大表明系统响应不符合碳基期望"
            },
            "task_interruptions": {
                "rule": "系统应预测可能的任务中断点，并提前保存状态以减少中断影响",
                "priority": "medium",
                "description": "频繁的任务中断会破坏创作心流"
            },
            "communication_rounds": {
                "rule": "当碳硅沟通回合数超过15时，系统应自动总结当前状态并提出明确的下一步建议",
                "priority": "low",
                "description": "过多的沟通回合表明协作效率低下"
            },
            "insufficient_data": {
                "rule": "系统应建立更完善的数据收集机制，以准确评估协同状态",
                "priority": "low",
                "description": "数据不足无法准确分析系统状态"
            }
        }
        
        template = micro_rule_templates.get(chaos_cause, micro_rule_templates["insufficient_data"])
        
        proposal = EvolutionProposal(
            proposal_id=proposal_id,
            chaos_cause=chaos_cause,
            micro_rule=template["rule"],
            priority=template["priority"],
            description=template["description"],
            estimated_impact=0.7,  # 预计影响
            voting_status="pending",
            carbon_vote=None,
            silicon_vote=None,
            timestamp=time.time()
        )
        
        self.evolution_proposals[proposal_id] = proposal
        
        return proposal
    
    def vote_on_proposal(self, proposal_id: str, voter_type: str, vote: bool) -> bool:
        """
        对提案进行投票
        
        Args:
            proposal_id: 提案ID
            voter_type: 投票者类型：carbon, silicon
            vote: 投票结果：True（赞成）, False（反对）
        
        Returns:
            是否投票成功
        """
        proposal = self.evolution_proposals.get(proposal_id)
        if not proposal:
            return False
        
        if voter_type == "carbon":
            proposal.carbon_vote = vote
        elif voter_type == "silicon":
            proposal.silicon_vote = vote
        else:
            return False
        
        # 检查投票是否完成
        if proposal.carbon_vote is not None and proposal.silicon_vote is not None:
            # 确定投票结果
            if proposal.carbon_vote and proposal.silicon_vote:
                proposal.voting_status = "passed"
                # 将通过的微规则添加到规则库
                self.micro_rules.append(proposal.micro_rule)
            else:
                proposal.voting_status = "rejected"
        
        return True
    
    def get_proposal(self, proposal_id: str) -> Optional[EvolutionProposal]:
        """
        获取提案
        
        Args:
            proposal_id: 提案ID
        
        Returns:
            提案对象，不存在则返回None
        """
        return self.evolution_proposals.get(proposal_id)
    
    def get_recent_proposals(self, limit: int = 5) -> List[EvolutionProposal]:
        """
        获取最近的提案
        
        Args:
            limit: 限制数量
        
        Returns:
            提案列表
        """
        proposals = list(self.evolution_proposals.values())
        proposals.sort(key=lambda p: p.timestamp, reverse=True)
        return proposals[:limit]
    
    def get_entropy_history(self, limit: int = 10) -> List[EntropyData]:
        """
        获取熵值历史
        
        Args:
            limit: 限制数量
        
        Returns:
            熵值数据列表
        """
        recent_history = self.entropy_history[-limit:]
        recent_history.reverse()  # 最新的在前
        return recent_history
    
    def get_micro_rules(self) -> List[str]:
        """
        获取所有生效的微规则
        
        Returns:
            微规则列表
        """
        return self.micro_rules
    
    def reset_entropy(self):
        """
        重置熵值
        当协议进化成功后调用
        """
        # 清空最近的熵值历史，保留部分历史数据用于趋势分析
        if len(self.entropy_history) > 20:
            self.entropy_history = self.entropy_history[-20:]
        else:
            self.entropy_history.clear()
    
    def get_system_health(self) -> Dict[str, Any]:
        """
        获取系统健康状态
        
        Returns:
            系统健康状态字典
        """
        if not self.entropy_history:
            return {
                "status": "insufficient_data",
                "entropy_score": 0.0,
                "evolution_status": "stable"
            }
        
        latest_entropy = self.entropy_history[-1]
        entropy_score = latest_entropy.entropy_score
        
        if entropy_score < 0.3:
            status = "healthy"
        elif entropy_score < 0.7:
            status = "warning"
        else:
            status = "critical"
        
        evolution_status = "triggered" if self.check_evolution_trigger() else "stable"
        
        return {
            "status": status,
            "entropy_score": round(entropy_score, 2),
            "evolution_status": evolution_status,
            "latest_entropy_data": {
                "validation_failure_rate": latest_entropy.validation_failure_rate,
                "satisfaction_volatility": latest_entropy.satisfaction_volatility,
                "task_interruption_count": latest_entropy.task_interruption_count,
                "communication_rounds": latest_entropy.communication_rounds
            }
        }
    
    def export_evolution_data(self, export_path: str) -> bool:
        """
        导出进化数据
        
        Args:
            export_path: 导出路径
        
        Returns:
            是否导出成功
        """
        try:
            export_data = {
                "entropy_history": [
                    {
                        "entropy_id": data.entropy_id,
                        "validation_failure_rate": data.validation_failure_rate,
                        "satisfaction_volatility": data.satisfaction_volatility,
                        "task_interruption_count": data.task_interruption_count,
                        "communication_rounds": data.communication_rounds,
                        "entropy_score": data.entropy_score,
                        "timestamp": data.timestamp
                    }
                    for data in self.entropy_history
                ],
                "proposals": [
                    {
                        "proposal_id": p.proposal_id,
                        "chaos_cause": p.chaos_cause,
                        "micro_rule": p.micro_rule,
                        "priority": p.priority,
                        "description": p.description,
                        "estimated_impact": p.estimated_impact,
                        "voting_status": p.voting_status,
                        "carbon_vote": p.carbon_vote,
                        "silicon_vote": p.silicon_vote,
                        "timestamp": p.timestamp
                    }
                    for p in self.evolution_proposals.values()
                ],
                "micro_rules": self.micro_rules,
                "system_health": self.get_system_health()
            }
            
            with open(export_path, "w", encoding="utf-8") as f:
                json.dump(export_data, f, ensure_ascii=False, indent=2)
            
            return True
        except Exception as e:
            print(f"导出进化数据失败: {str(e)}")
            return False
