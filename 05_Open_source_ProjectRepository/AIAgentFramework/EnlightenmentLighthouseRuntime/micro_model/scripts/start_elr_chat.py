#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
启动ELR容器协同对话微模型
功能：在ELR容器中运行elr_chat_model对话模型
"""

import sys
import os

# 切换到micro_model目录
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from examples.elr_chat_model import ELRChatModel

def main():
    """主函数"""
    print("=====================================")
    print("ELR容器协同对话微模型启动脚本")
    print("=====================================")
    print("模型：elr_chat_model")
    print("目标：local")
    print("")
    
    # 初始化模型
    model = ELRChatModel()
    
    print("=== ELR Interactive Chat (Python Mode) ===")
    print("欢迎使用ELR容器协同对话微模型！")
    print("您可以用英文或中文与模型对话。")
    print("输入 ',exit' 或 ',quit' 结束对话。")
    print("输入 ',help' 查看可用命令。")
    print("")
    
    # 主对话循环
    while True:
        try:
            user_input = input("你: ")
            if user_input.lower() in [',exit', ',quit']:
                print("模型: 再见！期待与您再次对话。")
                break
            response = model.predict(user_input)
            print(f"模型: {response}")
            print("")
        except KeyboardInterrupt:
            print("\n模型: 对话已中断，再见！")
            break
        except Exception as e:
            print(f"模型: 发生错误: {e}")
            print("")
    
    print("=====================================")
    print("对话结束")
    print("=====================================")

if __name__ == "__main__":
    main()
