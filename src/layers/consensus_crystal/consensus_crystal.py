"""
凝华沉淀层 (Consensus Crystal)
功能：将成功的协同经验，沉淀为可复用的共识晶体 (Consensus Crystal)
"""

from dataclasses import dataclass
from typing import Dict, List, Optional, Tuple, Any
import json
import uuid
import time
import os


@dataclass
class ConsensusCrystal:
    """
    共识晶体基类
    结构化的协同模板
    """
    crystal_id: str
    name: str
    description: str
    participating_voices: List[Dict[str, Any]]  # 参与声部及其能力配置
    counterpoint_pattern: str  # 使用的对位模式
    steps: List[Dict[str, Any]]  # 具体步骤
    decision_points: List[Dict[str, Any]]  # 关键决策点
    satisfaction_score: float  # 碳基满意度
    flow_duration: float  # 心流体验时长（分钟）
    micro_rules: List[str]  # 发现或验证的微小新规则
    creation_theme: str  # 创作主题
    created_at: float
    updated_at: float
    tags: List[str]  # 标签


class CrystalRepository:
    """
    共识晶体仓库
    管理和存储共识晶体
    """
    
    def __init__(self, storage_path: str = "./crystals"):
        """
        初始化共识晶体仓库
        
        Args:
            storage_path: 存储路径
        """
        self.storage_path = storage_path
        self.crystals: Dict[str, ConsensusCrystal] = {}
        
        # 创建存储目录
        os.makedirs(self.storage_path, exist_ok=True)
        
        # 加载已有的共识晶体
        self._load_crystals()
    
    def _load_crystals(self):
        """
        加载已有的共识晶体
        """
        try:
            for filename in os.listdir(self.storage_path):
                if filename.endswith(".json"):
                    file_path = os.path.join(self.storage_path, filename)
                    try:
                        with open(file_path, "r", encoding="utf-8") as f:
                            data = json.load(f)
                            crystal = ConsensusCrystal(**data)
                            self.crystals[crystal.crystal_id] = crystal
                    except Exception as e:
                        print(f"加载晶体文件失败: {file_path} - {str(e)}")
        except Exception as e:
            print(f"加载晶体失败: {str(e)}")
    
    def create_crystal(self, 
                     name: str, 
                     description: str, 
                     participating_voices: List[Dict[str, Any]],
                     counterpoint_pattern: str,
                     steps: List[Dict[str, Any]],
                     decision_points: List[Dict[str, Any]],
                     satisfaction_score: float,
                     flow_duration: float,
                     micro_rules: List[str],
                     creation_theme: str,
                     tags: List[str] = None) -> ConsensusCrystal:
        """
        创建共识晶体
        
        Args:
            name: 晶体名称
            description: 晶体描述
            participating_voices: 参与声部及其能力配置
            counterpoint_pattern: 使用的对位模式
            steps: 具体步骤
            decision_points: 关键决策点
            satisfaction_score: 碳基满意度
            flow_duration: 心流体验时长（分钟）
            micro_rules: 发现或验证的微小新规则
            creation_theme: 创作主题
            tags: 标签（可选）
        
        Returns:
            共识晶体对象
        """
        crystal_id = str(uuid.uuid4())
        timestamp = time.time()
        
        crystal = ConsensusCrystal(
            crystal_id=crystal_id,
            name=name,
            description=description,
            participating_voices=participating_voices,
            counterpoint_pattern=counterpoint_pattern,
            steps=steps,
            decision_points=decision_points,
            satisfaction_score=satisfaction_score,
            flow_duration=flow_duration,
            micro_rules=micro_rules,
            creation_theme=creation_theme,
            created_at=timestamp,
            updated_at=timestamp,
            tags=tags or []
        )
        
        self.crystals[crystal_id] = crystal
        self._save_crystal(crystal)
        
        return crystal
    
    def _save_crystal(self, crystal: ConsensusCrystal):
        """
        保存共识晶体到文件
        
        Args:
            crystal: 共识晶体对象
        """
        try:
            file_path = os.path.join(self.storage_path, f"{crystal.crystal_id}.json")
            with open(file_path, "w", encoding="utf-8") as f:
                # 将dataclass转换为字典
                crystal_dict = {
                    "crystal_id": crystal.crystal_id,
                    "name": crystal.name,
                    "description": crystal.description,
                    "participating_voices": crystal.participating_voices,
                    "counterpoint_pattern": crystal.counterpoint_pattern,
                    "steps": crystal.steps,
                    "decision_points": crystal.decision_points,
                    "satisfaction_score": crystal.satisfaction_score,
                    "flow_duration": crystal.flow_duration,
                    "micro_rules": crystal.micro_rules,
                    "creation_theme": crystal.creation_theme,
                    "created_at": crystal.created_at,
                    "updated_at": crystal.updated_at,
                    "tags": crystal.tags
                }
                json.dump(crystal_dict, f, ensure_ascii=False, indent=2)
        except Exception as e:
            print(f"保存晶体失败: {str(e)}")
    
    def get_crystal(self, crystal_id: str) -> Optional[ConsensusCrystal]:
        """
        获取共识晶体
        
        Args:
            crystal_id: 晶体ID
        
        Returns:
            共识晶体对象，不存在则返回None
        """
        return self.crystals.get(crystal_id)
    
    def get_all_crystals(self) -> List[ConsensusCrystal]:
        """
        获取所有共识晶体
        
        Returns:
            共识晶体列表
        """
        return list(self.crystals.values())
    
    def search_crystals(self, 
                       query: str = "", 
                       tags: List[str] = None, 
                       min_satisfaction: float = 0) -> List[ConsensusCrystal]:
        """
        搜索共识晶体
        
        Args:
            query: 搜索关键词
            tags: 标签过滤
            min_satisfaction: 最低满意度
        
        Returns:
            匹配的共识晶体列表
        """
        results = []
        
        for crystal in self.crystals.values():
            # 满意度过滤
            if crystal.satisfaction_score < min_satisfaction:
                continue
            
            # 标签过滤
            if tags:
                if not any(tag in crystal.tags for tag in tags):
                    continue
            
            # 关键词搜索
            if query:
                search_text = f"{crystal.name} {crystal.description} {crystal.creation_theme}"
                if query.lower() not in search_text.lower():
                    continue
            
            results.append(crystal)
        
        # 按满意度排序
        results.sort(key=lambda c: c.satisfaction_score, reverse=True)
        
        return results
    
    def update_crystal(self, crystal_id: str, 
                      updates: Dict[str, Any]) -> Optional[ConsensusCrystal]:
        """
        更新共识晶体
        
        Args:
            crystal_id: 晶体ID
            updates: 更新内容
        
        Returns:
            更新后的共识晶体对象，不存在则返回None
        """
        if crystal_id not in self.crystals:
            return None
        
        crystal = self.crystals[crystal_id]
        
        # 更新字段
        for key, value in updates.items():
            if hasattr(crystal, key):
                setattr(crystal, key, value)
        
        # 更新时间戳
        crystal.updated_at = time.time()
        
        # 保存更新
        self._save_crystal(crystal)
        
        return crystal
    
    def delete_crystal(self, crystal_id: str) -> bool:
        """
        删除共识晶体
        
        Args:
            crystal_id: 晶体ID
        
        Returns:
            是否删除成功
        """
        if crystal_id not in self.crystals:
            return False
        
        # 删除文件
        file_path = os.path.join(self.storage_path, f"{crystal_id}.json")
        if os.path.exists(file_path):
            try:
                os.remove(file_path)
            except Exception as e:
                print(f"删除晶体文件失败: {str(e)}")
        
        # 从内存中删除
        self.crystals.pop(crystal_id)
        
        return True
    
    def get_crystal_stats(self) -> Dict[str, Any]:
        """
        获取晶体统计信息
        
        Returns:
            统计信息字典
        """
        if not self.crystals:
            return {
                "total_crystals": 0,
                "average_satisfaction": 0,
                "average_flow_duration": 0,
                "most_common_tags": [],
                "most_used_patterns": []
            }
        
        # 计算平均满意度
        total_satisfaction = sum(c.satisfaction_score for c in self.crystals.values())
        average_satisfaction = total_satisfaction / len(self.crystals)
        
        # 计算平均心流时长
        total_flow = sum(c.flow_duration for c in self.crystals.values())
        average_flow = total_flow / len(self.crystals)
        
        # 统计最常用的标签
        tag_counts = {}
        for crystal in self.crystals.values():
            for tag in crystal.tags:
                tag_counts[tag] = tag_counts.get(tag, 0) + 1
        most_common_tags = sorted(tag_counts.items(), key=lambda x: x[1], reverse=True)[:5]
        
        # 统计最常用的对位模式
        pattern_counts = {}
        for crystal in self.crystals.values():
            pattern_counts[crystal.counterpoint_pattern] = pattern_counts.get(crystal.counterpoint_pattern, 0) + 1
        most_used_patterns = sorted(pattern_counts.items(), key=lambda x: x[1], reverse=True)[:5]
        
        return {
            "total_crystals": len(self.crystals),
            "average_satisfaction": round(average_satisfaction, 2),
            "average_flow_duration": round(average_flow, 2),
            "most_common_tags": [tag for tag, _ in most_common_tags],
            "most_used_patterns": [pattern for pattern, _ in most_used_patterns]
        }
    
    def export_crystal(self, crystal_id: str, export_path: str) -> bool:
        """
        导出共识晶体
        
        Args:
            crystal_id: 晶体ID
            export_path: 导出路径
        
        Returns:
            是否导出成功
        """
        crystal = self.get_crystal(crystal_id)
        if not crystal:
            return False
        
        try:
            with open(export_path, "w", encoding="utf-8") as f:
                crystal_dict = {
                    "crystal_id": crystal.crystal_id,
                    "name": crystal.name,
                    "description": crystal.description,
                    "participating_voices": crystal.participating_voices,
                    "counterpoint_pattern": crystal.counterpoint_pattern,
                    "steps": crystal.steps,
                    "decision_points": crystal.decision_points,
                    "satisfaction_score": crystal.satisfaction_score,
                    "flow_duration": crystal.flow_duration,
                    "micro_rules": crystal.micro_rules,
                    "creation_theme": crystal.creation_theme,
                    "created_at": crystal.created_at,
                    "updated_at": crystal.updated_at,
                    "tags": crystal.tags
                }
                json.dump(crystal_dict, f, ensure_ascii=False, indent=2)
            return True
        except Exception as e:
            print(f"导出晶体失败: {str(e)}")
            return False
    
    def import_crystal(self, import_path: str) -> Optional[ConsensusCrystal]:
        """
        导入共识晶体
        
        Args:
            import_path: 导入路径
        
        Returns:
            导入的共识晶体对象，失败则返回None
        """
        try:
            with open(import_path, "r", encoding="utf-8") as f:
                data = json.load(f)
            
            # 创建新的晶体ID，避免冲突
            data["crystal_id"] = str(uuid.uuid4())
            data["created_at"] = time.time()
            data["updated_at"] = time.time()
            
            crystal = ConsensusCrystal(**data)
            self.crystals[crystal.crystal_id] = crystal
            self._save_crystal(crystal)
            
            return crystal
        except Exception as e:
            print(f"导入晶体失败: {str(e)}")
            return None
    
    def generate_crystal_from_execution(self, 
                                       execution_results: Dict[str, Any],
                                       satisfaction_score: float,
                                       flow_duration: float,
                                       micro_rules: List[str]) -> ConsensusCrystal:
        """
        从执行结果生成共识晶体
        
        Args:
            execution_results: 执行结果
            satisfaction_score: 满意度
            flow_duration: 心流时长
            micro_rules: 微小新规则
        
        Returns:
            生成的共识晶体
        """
        # 实际实现中应从执行结果中提取信息
        # 这里仅创建一个示例晶体
        return self.create_crystal(
            name=f"从执行生成的晶体 {time.strftime('%Y%m%d_%H%M%S')}",
            description="从协同执行结果自动生成的共识晶体",
            participating_voices=[],
            counterpoint_pattern="",
            steps=[],
            decision_points=[],
            satisfaction_score=satisfaction_score,
            flow_duration=flow_duration,
            micro_rules=micro_rules,
            creation_theme="",
            tags=["autogenerated"]
        )
