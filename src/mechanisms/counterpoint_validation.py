"""
对位验证机制 (Counterpoint Validation)
功能：确保碳硅协同符合"对位法"的实时检查与校准机制
"""

from dataclasses import dataclass
from typing import Dict, List, Optional, Tuple, Any
import json
import uuid
import time


@dataclass
class ValidationResult:
    """
    验证结果基类
    """
    validation_id: str
    action_id: str
    carbon_intention: Dict[str, Any]
    silicon_output: Dict[str, Any]
    thinking_process: List[str]  # 硅基思考链
    differences: List[Dict[str, Any]]  # 检测到的差异
    negotiation_status: str  # 协商状态：none, in_progress, completed
    negotiation_outcome: Optional[str]  # 协商结果
    protocol_evolution: bool  # 是否触发协议进化
    timestamp: float


class CounterpointValidator:
    """
    对位验证器
    实现碳硅协同的实时检查与校准
    """
    
    def __init__(self):
        """
        初始化对位验证器
        """
        self.validation_results: Dict[str, ValidationResult] = {}
        self.difference_threshold = 0.3  # 差异阈值
        self.evolution_trigger_count = 3  # 进化触发次数
        self.difference_counter: Dict[str, int] = {}  # 差异计数器
    
    def validate(self, action_id: str, 
                 carbon_intention: Dict[str, Any], 
                 silicon_output: Dict[str, Any], 
                 thinking_process: List[str]) -> ValidationResult:
        """
        执行对位验证
        
        Args:
            action_id: 动作ID
            carbon_intention: 碳基意图
            silicon_output: 硅基输出
            thinking_process: 硅基思考链
        
        Returns:
            验证结果
        """
        validation_id = str(uuid.uuid4())
        
        # 1. 思考显影：验证硅基是否清晰展示推理过程
        if not self._validate_thinking_visualization(thinking_process):
            return ValidationResult(
                validation_id=validation_id,
                action_id=action_id,
                carbon_intention=carbon_intention,
                silicon_output=silicon_output,
                thinking_process=thinking_process,
                differences=[{"type": "missing_thinking_process", "message": "硅基未清晰展示推理过程"}],
                negotiation_status="none",
                negotiation_outcome=None,
                protocol_evolution=False,
                timestamp=time.time()
            )
        
        # 2. 差异检测：对比碳基意图与硅基输出
        differences = self._detect_differences(carbon_intention, silicon_output)
        
        negotiation_status = "none"
        negotiation_outcome = None
        protocol_evolution = False
        
        # 3. 共识协商：当差异超过阈值，启动碳硅协商
        if len(differences) > 0:
            difference_score = self._calculate_difference_score(differences)
            if difference_score > self.difference_threshold:
                negotiation_status, negotiation_outcome = self._initiate_negotiation(
                    carbon_intention, silicon_output, differences
                )
                
                # 4. 协议进化：如果某种差异频繁出现，触发协议升级
                for diff in differences:
                    diff_type = diff.get("type", "unknown")
                    self.difference_counter[diff_type] = self.difference_counter.get(diff_type, 0) + 1
                    
                    if self.difference_counter[diff_type] >= self.evolution_trigger_count:
                        protocol_evolution = True
                        # 重置计数器
                        self.difference_counter[diff_type] = 0
        
        result = ValidationResult(
            validation_id=validation_id,
            action_id=action_id,
            carbon_intention=carbon_intention,
            silicon_output=silicon_output,
            thinking_process=thinking_process,
            differences=differences,
            negotiation_status=negotiation_status,
            negotiation_outcome=negotiation_outcome,
            protocol_evolution=protocol_evolution,
            timestamp=time.time()
        )
        
        self.validation_results[validation_id] = result
        
        return result
    
    def _validate_thinking_visualization(self, thinking_process: List[str]) -> bool:
        """
        验证思考显影
        
        Args:
            thinking_process: 硅基思考链
        
        Returns:
            是否符合要求
        """
        # 简单验证：思考链不能为空且长度至少为2
        return len(thinking_process) >= 2
    
    def _detect_differences(self, carbon_intention: Dict[str, Any], 
                           silicon_output: Dict[str, Any]) -> List[Dict[str, Any]]:
        """
        检测差异
        
        Args:
            carbon_intention: 碳基意图
            silicon_output: 硅基输出
        
        Returns:
            差异列表
        """
        differences = []
        
        # 检测意图与输出的差异
        # 这里仅实现基本的差异检测逻辑
        
        # 检测风格差异
        carbon_style = carbon_intention.get("style", "")
        silicon_style = silicon_output.get("style", "")
        if carbon_style and silicon_style and carbon_style != silicon_style:
            differences.append({
                "type": "style_mismatch",
                "message": f"风格不匹配: 期望 '{carbon_style}', 实际 '{silicon_style}'",
                "severity": "medium"
            })
        
        # 检测主题差异
        carbon_theme = carbon_intention.get("theme", "")
        silicon_theme = silicon_output.get("theme", "")
        if carbon_theme and silicon_theme:
            if carbon_theme.lower() not in silicon_theme.lower():
                differences.append({
                    "type": "theme_mismatch",
                    "message": f"主题不匹配: 期望包含 '{carbon_theme}', 实际 '{silicon_theme}'",
                    "severity": "high"
                })
        
        # 检测情感差异
        carbon_emotion = carbon_intention.get("emotion", "")
        silicon_emotion = silicon_output.get("emotion", "")
        if carbon_emotion and silicon_emotion and carbon_emotion != silicon_emotion:
            differences.append({
                "type": "emotion_mismatch",
                "message": f"情感不匹配: 期望 '{carbon_emotion}', 实际 '{silicon_emotion}'",
                "severity": "medium"
            })
        
        # 检测长度差异
        carbon_length = carbon_intention.get("length", 0)
        silicon_length = silicon_output.get("length", 0)
        if carbon_length > 0 and silicon_length > 0:
            length_ratio = abs(carbon_length - silicon_length) / max(carbon_length, silicon_length)
            if length_ratio > 0.5:
                differences.append({
                    "type": "length_mismatch",
                    "message": f"长度不匹配: 期望 {carbon_length}, 实际 {silicon_length}",
                    "severity": "low"
                })
        
        return differences
    
    def _calculate_difference_score(self, differences: List[Dict[str, Any]]) -> float:
        """
        计算差异得分
        
        Args:
            differences: 差异列表
        
        Returns:
            差异得分（0-1）
        """
        if not differences:
            return 0.0
        
        # 基于严重程度计算得分
        severity_scores = {
            "low": 0.1,
            "medium": 0.3,
            "high": 0.6
        }
        
        total_score = 0.0
        for diff in differences:
            severity = diff.get("severity", "medium")
            total_score += severity_scores.get(severity, 0.3)
        
        # 归一化得分
        max_possible_score = len(differences) * 0.6
        if max_possible_score == 0:
            return 0.0
        
        normalized_score = min(total_score / max_possible_score, 1.0)
        return normalized_score
    
    def _initiate_negotiation(self, carbon_intention: Dict[str, Any], 
                             silicon_output: Dict[str, Any], 
                             differences: List[Dict[str, Any]]) -> Tuple[str, str]:
        """
        启动共识协商
        
        Args:
            carbon_intention: 碳基意图
            silicon_output: 硅基输出
            differences: 检测到的差异
        
        Returns:
            (协商状态, 协商结果)
        """
        # 实际实现中应启动真实的协商过程
        # 这里仅模拟协商结果
        
        # 生成协商结果
        negotiation_outcome = ""
        for diff in differences:
            if diff["type"] == "style_mismatch":
                negotiation_outcome += f"调整风格以匹配碳基期望: {diff['message']}\n"
            elif diff["type"] == "theme_mismatch":
                negotiation_outcome += f"重新调整主题以符合碳基意图: {diff['message']}\n"
            elif diff["type"] == "emotion_mismatch":
                negotiation_outcome += f"调整情感表达: {diff['message']}\n"
            elif diff["type"] == "length_mismatch":
                negotiation_outcome += f"调整长度: {diff['message']}\n"
        
        return "completed", negotiation_outcome.strip()
    
    def get_validation_result(self, validation_id: str) -> Optional[ValidationResult]:
        """
        获取验证结果
        
        Args:
            validation_id: 验证ID
        
        Returns:
            验证结果对象，不存在则返回None
        """
        return self.validation_results.get(validation_id)
    
    def get_recent_validations(self, limit: int = 10) -> List[ValidationResult]:
        """
        获取最近的验证结果
        
        Args:
            limit: 限制数量
        
        Returns:
            验证结果列表
        """
        results = list(self.validation_results.values())
        results.sort(key=lambda r: r.timestamp, reverse=True)
        return results[:limit]
    
    def generate_thinking_visualization(self, silicon_output: Dict[str, Any]) -> List[str]:
        """
        生成硅基思考显影
        
        Args:
            silicon_output: 硅基输出
        
        Returns:
            思考链列表
        """
        # 实际实现中应从硅基系统获取真实的思考过程
        # 这里仅生成模拟的思考链
        
        thinking_process = [
            f"接收到碳基请求，分析意图: {json.dumps(silicon_output.get('intention', {}), ensure_ascii=False)}",
            f"确定创作风格: {silicon_output.get('style', '未知')}",
            f"生成初步方案...",
            f"评估方案与意图的匹配度",
            f"调整输出以更好地符合碳基期望",
            f"完成最终输出"
        ]
        
        return thinking_process
    
    def suggest_protocol_evolution(self, difference_type: str) -> Dict[str, Any]:
        """
        建议协议进化
        
        Args:
            difference_type: 差异类型
        
        Returns:
            进化建议
        """
        # 基于差异类型生成进化建议
        evolution_suggestions = {
            "style_mismatch": {
                "suggestion": "新增微规则: 当碳基指定风格时，硅基应提供三个不同风格强度的选项",
                "priority": "medium",
                "reason": "频繁出现风格不匹配问题，需要更灵活的风格调整机制"
            },
            "theme_mismatch": {
                "suggestion": "新增微规则: 当碳基标记为'抽象'时，硅基应提供三个不同抽象层次的选择",
                "priority": "high",
                "reason": "主题是创作的核心，需要更精确的主题匹配机制"
            },
            "emotion_mismatch": {
                "suggestion": "新增微规则: 硅基应提供情感强度滑块，允许碳基实时调整",
                "priority": "medium",
                "reason": "情感表达需要更精细的调整机制"
            },
            "length_mismatch": {
                "suggestion": "新增微规则: 硅基应根据碳基历史偏好自动调整输出长度",
                "priority": "low",
                "reason": "长度差异对创作质量影响较小，但可通过历史学习优化"
            }
        }
        
        return evolution_suggestions.get(difference_type, {
            "suggestion": f"针对 {difference_type} 类型的差异，建议新增相应的微规则",
            "priority": "medium",
            "reason": "该类型差异频繁出现，需要相应的解决方案"
        })
    
    def get_difference_stats(self) -> Dict[str, Any]:
        """
        获取差异统计信息
        
        Returns:
            差异统计信息字典
        """
        total_validations = len(self.validation_results)
        total_differences = sum(len(r.differences) for r in self.validation_results.values())
        
        difference_type_count = {}
        for result in self.validation_results.values():
            for diff in result.differences:
                diff_type = diff.get("type", "unknown")
                difference_type_count[diff_type] = difference_type_count.get(diff_type, 0) + 1
        
        return {
            "total_validations": total_validations,
            "total_differences": total_differences,
            "average_differences_per_validation": round(total_differences / max(total_validations, 1), 2),
            "difference_type_distribution": difference_type_count,
            "evolution_trigger_count": self.evolution_trigger_count
        }
    
    def reset(self):
        """
        重置验证器状态
        """
        self.validation_results.clear()
        self.difference_counter.clear()
    
    def export_validation_data(self, export_path: str) -> bool:
        """
        导出验证数据
        
        Args:
            export_path: 导出路径
        
        Returns:
            是否导出成功
        """
        try:
            export_data = []
            for result in self.validation_results.values():
                result_dict = {
                    "validation_id": result.validation_id,
                    "action_id": result.action_id,
                    "carbon_intention": result.carbon_intention,
                    "silicon_output": result.silicon_output,
                    "thinking_process": result.thinking_process,
                    "differences": result.differences,
                    "negotiation_status": result.negotiation_status,
                    "negotiation_outcome": result.negotiation_outcome,
                    "protocol_evolution": result.protocol_evolution,
                    "timestamp": result.timestamp
                }
                export_data.append(result_dict)
            
            with open(export_path, "w", encoding="utf-8") as f:
                json.dump(export_data, f, ensure_ascii=False, indent=2)
            
            return True
        except Exception as e:
            print(f"导出验证数据失败: {str(e)}")
            return False
