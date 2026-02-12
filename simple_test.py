#!/usr/bin/env python3
"""
简单的项目验证脚本
"""

print("开始执行Meta-CreationPower项目验证...")
print("=" * 60)

# 尝试导入核心模块
try:
    from src.layers.voice_recognition.voice_recognition import CollaborativeSonicMap
    print("✅ 声部识别层导入成功")
except Exception as e:
    print(f"❌ 声部识别层导入失败: {e}")

try:
    from src.layers.meta_protocol.meta_protocol import MetaProtocolManager
    print("✅ 元协议锚定层导入成功")
except Exception as e:
    print(f"❌ 元协议锚定层导入失败: {e}")

try:
    from src.layers.counterpoint_design.counterpoint_design import CounterpointDesigner
    print("✅ 协奏设计层导入成功")
except Exception as e:
    print(f"❌ 协奏设计层导入失败: {e}")

try:
    from src.layers.steady_execution.steady_execution import SteadyExecutor
    print("✅ 静定执行层导入成功")
except Exception as e:
    print(f"❌ 静定执行层导入失败: {e}")

try:
    from src.layers.consensus_crystal.consensus_crystal import CrystalRepository
    print("✅ 凝华沉淀层导入成功")
except Exception as e:
    print(f"❌ 凝华沉淀层导入失败: {e}")

try:
    from src.mechanisms.counterpoint_validation import CounterpointValidator
    print("✅ 对位验证机制导入成功")
except Exception as e:
    print(f"❌ 对位验证机制导入失败: {e}")

try:
    from src.mechanisms.entropy_evolution import EntropyEvolutionManager
    print("✅ 熵值驱动协议进化导入成功")
except Exception as e:
    print(f"❌ 熵值驱动协议进化导入失败: {e}")

print("=" * 60)
print("项目验证完成！")
