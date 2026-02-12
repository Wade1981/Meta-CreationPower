"""
声部识别层 (Collaborative Sonic Map)
功能：识别与注册参与协同的"声部"（碳基用户与硅基智能体），建立动态的协同声部图谱
"""

from dataclasses import dataclass
from typing import Dict, List, Optional, Tuple
import json
import uuid


@dataclass
class Voice:
    """声部基类"""
    voice_id: str
    name: str
    voice_type: str  # 'carbon' 或 'silicon'
    capability_vector: Dict[str, float]  # 能力向量
    intention_vector: Dict[str, float]  # 意图向量
    description: str = ""
    created_at: str = ""
    last_active: str = ""


class CollaborativeSonicMap:
    """协同声部图谱"""
    
    def __init__(self):
        self.voices: Dict[str, Voice] = {}
        self.voice_map: Dict[str, List[str]] = {}  # 声部关系映射
    
    def register_voice(self, name: str, voice_type: str, 
                      capability_vector: Dict[str, float], 
                      intention_vector: Dict[str, float], 
                      description: str = "") -> Voice:
        """
        注册新声部
        
        Args:
            name: 声部名称
            voice_type: 声部类型 ('carbon' 或 'silicon')
            capability_vector: 能力向量
            intention_vector: 意图向量
            description: 声部描述
        
        Returns:
            注册的声部对象
        """
        voice_id = str(uuid.uuid4())
        voice = Voice(
            voice_id=voice_id,
            name=name,
            voice_type=voice_type,
            capability_vector=capability_vector,
            intention_vector=intention_vector,
            description=description,
            created_at="",  # 实际实现中应使用时间戳
            last_active=""
        )
        
        self.voices[voice_id] = voice
        if voice_type not in self.voice_map:
            self.voice_map[voice_type] = []
        self.voice_map[voice_type].append(voice_id)
        
        return voice
    
    def get_voice(self, voice_id: str) -> Optional[Voice]:
        """
        获取声部信息
        
        Args:
            voice_id: 声部ID
        
        Returns:
            声部对象，不存在则返回None
        """
        return self.voices.get(voice_id)
    
    def get_voices_by_type(self, voice_type: str) -> List[Voice]:
        """
        按类型获取声部列表
        
        Args:
            voice_type: 声部类型
        
        Returns:
            声部对象列表
        """
        voice_ids = self.voice_map.get(voice_type, [])
        return [self.voices[vid] for vid in voice_ids if vid in self.voices]
    
    def update_voice_activity(self, voice_id: str):
        """
        更新声部活动状态
        
        Args:
            voice_id: 声部ID
        """
        if voice_id in self.voices:
            self.voices[voice_id].last_active = ""  # 实际实现中应使用时间戳
    
    def remove_voice(self, voice_id: str):
        """
        移除声部
        
        Args:
            voice_id: 声部ID
        """
        if voice_id in self.voices:
            voice_type = self.voices[voice_id].voice_type
            self.voices.pop(voice_id)
            if voice_type in self.voice_map:
                self.voice_map[voice_type].remove(voice_id)
    
    def get_voice_map(self) -> Dict:
        """
        获取声部图谱
        
        Returns:
            声部图谱字典
        """
        return {
            "voices": {vid: {
                "name": v.name,
                "type": v.voice_type,
                "capabilities": v.capability_vector,
                "intentions": v.intention_vector,
                "description": v.description
            } for vid, v in self.voices.items()},
            "voice_map": self.voice_map
        }
    
    def save_to_file(self, file_path: str):
        """
        保存声部图谱到文件
        
        Args:
            file_path: 文件路径
        """
        # 实际实现中应序列化保存
        pass
    
    def load_from_file(self, file_path: str):
        """
        从文件加载声部图谱
        
        Args:
            file_path: 文件路径
        """
        # 实际实现中应反序列化加载
        pass
