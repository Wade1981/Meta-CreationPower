print("开始测试 voice_recognition 模块...")

# 测试模块导入
try:
    from src.layers.voice_recognition.voice_recognition import CollaborativeSonicMap
    print("✅ 成功导入 CollaborativeSonicMap")
except Exception as e:
    print("❌ 导入失败:", e)
    exit(1)

# 测试创建对象
try:
    sonic_map = CollaborativeSonicMap()
    print("✅ 成功创建 CollaborativeSonicMap 对象")
except Exception as e:
    print("❌ 创建对象失败:", e)
    exit(1)

# 测试注册碳基声部
try:
    carbon_voice = sonic_map.register_voice(
        name="X54先生",
        voice_type="carbon",
        capability_vector={"创意生成": 0.9, "逻辑分析": 0.7, "情感共鸣": 0.9},
        intention_vector={"探索性": 0.8, "完美性": 0.7, "效率": 0.6}
    )
    print("✅ 成功注册碳基声部，ID:", carbon_voice.voice_id)
except Exception as e:
    print("❌ 注册碳基声部失败:", e)

# 测试注册硅基声部
try:
    silicon_voice = sonic_map.register_voice(
        name="豆包",
        voice_type="silicon",
        capability_vector={"创意生成": 0.8, "逻辑分析": 0.9, "情感共鸣": 0.5},
        intention_vector={"探索性": 0.7, "完美性": 0.8, "效率": 0.9}
    )
    print("✅ 成功注册硅基声部，ID:", silicon_voice.voice_id)
except Exception as e:
    print("❌ 注册硅基声部失败:", e)

print("测试完成!")
