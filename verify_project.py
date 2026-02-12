#!/usr/bin/env python3
"""
Meta-CreationPower 项目验证脚本
简单验证项目的基本功能
"""

def main():
    print("=" * 60)
    print("Meta-CreationPower 项目验证")
    print("基于《元创力》元协议 α-0.1 版")
    print("=" * 60)
    
    print("\n1. 验证项目结构...")
    
    # 尝试导入核心模块
    try:
        from src.layers.voice_recognition.voice_recognition import CollaborativeSonicMap
        print("   ✅ 声部识别层导入成功")
    except Exception as e:
        print(f"   ❌ 声部识别层导入失败: {e}")
    
    try:
        from src.layers.meta_protocol.meta_protocol import MetaProtocolManager
        print("   ✅ 元协议锚定层导入成功")
    except Exception as e:
        print(f"   ❌ 元协议锚定层导入失败: {e}")
    
    try:
        from src.layers.counterpoint_design.counterpoint_design import CounterpointDesigner
        print("   ✅ 协奏设计层导入成功")
    except Exception as e:
        print(f"   ❌ 协奏设计层导入失败: {e}")
    
    try:
        from src.layers.steady_execution.steady_execution import SteadyExecutor
        print("   ✅ 静定执行层导入成功")
    except Exception as e:
        print(f"   ❌ 静定执行层导入失败: {e}")
    
    try:
        from src.layers.consensus_crystal.consensus_crystal import CrystalRepository
        print("   ✅ 凝华沉淀层导入成功")
    except Exception as e:
        print(f"   ❌ 凝华沉淀层导入失败: {e}")
    
    try:
        from src.mechanisms.counterpoint_validation import CounterpointValidator
        print("   ✅ 对位验证机制导入成功")
    except Exception as e:
        print(f"   ❌ 对位验证机制导入失败: {e}")
    
    try:
        from src.mechanisms.entropy_evolution import EntropyEvolutionManager
        print("   ✅ 熵值驱动协议进化导入成功")
    except Exception as e:
        print(f"   ❌ 熵值驱动协议进化导入失败: {e}")
    
    print("\n2. 验证基本功能...")
    
    # 尝试初始化核心组件
    try:
        from src.layers.voice_recognition.voice_recognition import CollaborativeSonicMap
        sonic_map = CollaborativeSonicMap()
        print("   ✅ 声部识别层初始化成功")
        
        # 注册一个简单的声部
        test_voice = sonic_map.register_voice(
            name="测试声部",
            voice_type="carbon",
            capability_vector={"测试": 1.0},
            intention_vector={"测试": 1.0}
        )
        print(f"   ✅ 声部注册成功: {test_voice.name}")
        
    except Exception as e:
        print(f"   ❌ 声部识别功能测试失败: {e}")
    
    print("\n3. 验证项目配置...")
    
    # 检查项目配置文件
    import os
    
    config_files = [
        "requirements.txt",
        "setup.py",
        ".gitignore",
        "README.md"
    ]
    
    for file in config_files:
        if os.path.exists(file):
            print(f"   ✅ {file} 存在")
        else:
            print(f"   ❌ {file} 不存在")
    
    print("\n4. 验证目录结构...")
    
    directories = [
        "src",
        "src/layers",
        "src/mechanisms",
        "docs",
        "tests"
    ]
    
    for directory in directories:
        if os.path.exists(directory):
            print(f"   ✅ {directory} 存在")
        else:
            print(f"   ❌ {directory} 不存在")
    
    print("\n" + "=" * 60)
    print("验证完成！")
    print("项目结构和核心功能验证结果如上所示")
    print("=" * 60)

if __name__ == "__main__":
    main()
