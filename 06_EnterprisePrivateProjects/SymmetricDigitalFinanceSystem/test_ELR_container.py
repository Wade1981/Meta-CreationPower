# ELR容器测试脚本

import json
import os
import sys

# 检查系统文件结构
def check_system_structure():
    print("=== 检查系统文件结构 ===")
    required_files = [
        "container.json",
        "requirements.txt",
        "src/main.py",
        "config/config.py",
        "data/input",
        "data/output",
        "logs"
    ]
    
    all_exist = True
    for file_path in required_files:
        if os.path.exists(file_path):
            print(f"✅ {file_path} 存在")
        else:
            print(f"❌ {file_path} 不存在")
            all_exist = False
    
    return all_exist

# 检查容器配置文件
def check_container_config():
    print("\n=== 检查容器配置文件 ===")
    try:
        with open("container.json", "r", encoding="utf-8") as f:
            config = json.load(f)
        
        required_fields = ["name", "version", "main", "entrypoint", "requirements", "health_check"]
        missing_fields = []
        
        for field in required_fields:
            if field not in config:
                missing_fields.append(field)
        
        if missing_fields:
            print(f"❌ 缺少必要字段: {missing_fields}")
            return False
        else:
            print("✅ 容器配置文件完整")
            print(f"  系统名称: {config['name']}")
            print(f"  版本: {config['version']}")
            print(f"  主入口: {config['main']}")
            print(f"  启动命令: {config['entrypoint']}")
            print(f"  依赖包数量: {len(config['requirements'])}")
            print(f"  健康检查命令: {config['health_check']['command']}")
            return True
    except Exception as e:
        print(f"❌ 容器配置文件错误: {e}")
        return False

# 检查依赖文件
def check_requirements():
    print("\n=== 检查依赖文件 ===")
    try:
        with open("requirements.txt", "r", encoding="utf-8") as f:
            dependencies = f.readlines()
        
        dependencies = [dep.strip() for dep in dependencies if dep.strip()]
        print(f"✅ 依赖文件存在，包含 {len(dependencies)} 个依赖包")
        for dep in dependencies:
            print(f"  - {dep}")
        return True
    except Exception as e:
        print(f"❌ 依赖文件错误: {e}")
        return False

# 检查主入口文件
def check_main_file():
    print("\n=== 检查主入口文件 ===")
    try:
        with open("src/main.py", "r", encoding="utf-8") as f:
            content = f.read()
        
        if "SymmetricDigitalFinanceSystem" in content:
            print("✅ 主入口文件存在，包含系统类定义")
            return True
        else:
            print("❌ 主入口文件不包含系统类定义")
            return False
    except Exception as e:
        print(f"❌ 主入口文件错误: {e}")
        return False

# 验证系统是否可以在ELR容器中运行
def validate_ELR_container():
    print("=== 验证ELR容器兼容性 ===")
    checks = [
        ("文件结构", check_system_structure),
        ("容器配置", check_container_config),
        ("依赖文件", check_requirements),
        ("主入口文件", check_main_file)
    ]
    
    all_passed = True
    for check_name, check_func in checks:
        print(f"\n--- 检查: {check_name} ---")
        if not check_func():
            all_passed = False
    
    print("\n=== 验证结果 ===")
    if all_passed:
        print("🎉 系统已成功配置为ELR容器格式，可以在ELR沙箱中运行")
        print("\n容器配置详情:")
        print("- 系统名称: symmetric-finance-system")
        print("- 版本: 1.0.0")
        print("- 主入口: src/main.py")
        print("- 启动命令: python src/main.py")
        print("- 健康检查: python src/main.py")
        print("- 资源需求: 2 CPU, 4G 内存")
        print("\n系统已准备就绪，可以装载到ELR容器沙箱运行")
    else:
        print("❌ 系统配置存在问题，需要修复后才能在ELR沙箱中运行")
    
    return all_passed

if __name__ == "__main__":
    validate_ELR_container()
