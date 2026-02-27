from src.layers.voice_recognition.voice_recognition import CollaborativeSonicMap

def test_voice_registration():
    """测试声部注册功能"""
    print("测试声部注册功能...")
    sonic_map = CollaborativeSonicMap()
    
    # 注册碳基声部
    carbon_voice = sonic_map.register_voice(
        name="X54先生",
        voice_type="carbon",
        capability_vector={"创意生成": 0.9, "逻辑分析": 0.7, "情感共鸣": 0.9},
        intention_vector={"探索性": 0.8, "完美性": 0.7, "效率": 0.6}
    )
    
    # 注册硅基声部
    silicon_voice = sonic_map.register_voice(
        name="豆包",
        voice_type="silicon",
        capability_vector={"创意生成": 0.8, "逻辑分析": 0.9, "情感共鸣": 0.5},
        intention_vector={"探索性": 0.7, "完美性": 0.8, "效率": 0.9}
    )
    
    print(f"碳基声部 ID: {carbon_voice.voice_id}")
    print(f"硅基声部 ID: {silicon_voice.voice_id}")
    print("✅ 声部注册测试通过！")
    return True

if __name__ == "__main__":
    test_voice_registration()
