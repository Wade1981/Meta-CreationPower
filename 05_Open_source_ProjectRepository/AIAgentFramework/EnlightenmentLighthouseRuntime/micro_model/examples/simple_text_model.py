#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
极简文本响应微模型
功能：输入文本问候→输出碳硅协同标准友好回应
无外部依赖，可直接在ELR沙箱中运行
"""

class SimpleTextModel:
    """极简文本响应微模型类"""
    
    def __init__(self):
        """初始化模型"""
        self.model_name = "simple_text_model"
        self.version = "1.0"
        self.description = "极简文本响应微模型，无外部依赖"
        print(f"初始化模型: {self.model_name} v{self.version}")
    
    def predict(self, input_text):
        """
        模型推理方法
        参数:
            input_text: 输入文本
        返回:
            响应文本
        """
        # 处理输入文本，生成友好回应
        response = self._generate_response(input_text)
        return response
    
    def _generate_response(self, input_text):
        """
        生成响应文本
        参数:
            input_text: 输入文本
        返回:
            响应文本
        """
        # 转换为小写，便于匹配
        lower_input = input_text.lower()
        
        # 问候语匹配
        greetings = ["hello", "hi", "你好", "嗨", "hey"]
        for greeting in greetings:
            if greeting in lower_input:
                return f"碳硅协同问候！我是{self.model_name}，很高兴为您服务。您的问候已收到，祝您一天愉快！"
        
        # 询问类匹配
        questions = ["how are you", "你好吗", "怎么样", "最近好吗"]
        for question in questions:
            if question in lower_input:
                return f"碳硅协同回应！我是{self.model_name}，运行状态良好。感谢您的关心，祝您一切顺利！"
        
        # 默认回应
        return f"碳硅协同回应！我是{self.model_name}，已收到您的消息：{input_text}。随时为您服务！"
    
    def get_info(self):
        """
        获取模型信息
        返回:
            模型信息字典
        """
        return {
            "model_name": self.model_name,
            "version": self.version,
            "description": self.description,
            "capabilities": ["文本响应", "无外部依赖", "轻量级"]
        }

# 测试代码
if __name__ == "__main__":
    # 初始化模型
    model = SimpleTextModel()
    
    # 打印模型信息
    print("模型信息:")
    print(model.get_info())
    
    # 测试案例
    test_cases = [
        "Hello, how are you?",
        "你好，最近怎么样？",
        "Hi there!",
        "请问这个模型能做什么？",
        "测试消息"
    ]
    
    print("\n测试结果:")
    for test_input in test_cases:
        response = model.predict(test_input)
        print(f"输入: {test_input}")
        print(f"输出: {response}")
        print("---")
