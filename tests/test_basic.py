print("开始测试...")

# 测试基本Python功能
print("Python版本:", __import__('sys').version)

# 测试模块导入
try:
    from src.layers.voice_recognition.voice_recognition import CollaborativeSonicMap
    print("✅ 成功导入 CollaborativeSonicMap")
    
    # 测试创建对象
    sonic_map = CollaborativeSonicMap()
    print("✅ 成功创建 CollaborativeSonicMap 对象")
    
except ImportError as e:
    print("❌ 导入失败:", e)
except Exception as e:
    print("❌ 执行失败:", e)

print("测试完成!")
