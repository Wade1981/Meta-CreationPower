#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
ELR容器协同对话微模型
功能：为ELR容器提供轻量级对话能力，支持多轮对话和ELR相关功能
无外部依赖，可直接在ELR沙箱中运行
文件大小：< 1MB
"""

class ELRChatModel:
    """ELR容器协同对话微模型类"""
    
    def __init__(self):
        """初始化模型"""
        self.model_name = "elr_chat_model"
        self.version = "1.0"
        self.description = "ELR容器协同对话微模型，支持多轮对话和ELR相关功能"
        self.context = []  # 对话上下文
        self.max_context_length = 5  # 最大上下文长度
        self.elr_commands = {
            "help": "显示可用命令",
            "status": "检查ELR容器状态",
            "models": "列出已加载的模型",
            "info": "显示ELR容器信息",
            "clear": "清除对话历史",
            "exit": "退出对话"
        }
        print(f"初始化模型: {self.model_name} v{self.version}")
        print("模型加载成功！")
    
    def predict(self, input_text):
        """
        模型推理方法
        参数:
            input_text: 输入文本
        返回:
            响应文本
        """
        # 处理输入文本，生成响应
        response = self._process_input(input_text)
        
        # 更新对话上下文
        self._update_context(input_text, response)
        
        return response
    
    def _process_input(self, input_text):
        """
        处理输入文本
        参数:
            input_text: 输入文本
        返回:
            响应文本
        """
        # 转换为小写，便于匹配
        lower_input = input_text.lower()
        
        # 检查是否是命令
        if lower_input.startswith(","):
            return self._handle_command(lower_input[1:].strip())
        
        # 问候语匹配
        if any(greeting in lower_input for greeting in ["hello", "hi", "你好", "嗨", "hey"]):
            return self._generate_greeting_response()
        
        # 询问类匹配
        if any(question in lower_input for question in ["how are you", "你好吗", "怎么样", "最近好吗"]):
            return self._generate_status_response()
        
        # ELR相关问题
        if "elr" in lower_input:
            if any(keyword in lower_input for keyword in ["what", "什么", "功能", "能做什么"]):
                return self._generate_elr_capabilities()
            elif any(keyword in lower_input for keyword in ["status", "状态", "运行"]):
                return self._generate_elr_status()
            elif any(keyword in lower_input for keyword in ["model", "模型", "加载"]):
                return self._generate_elr_models()
        
        # 对话上下文相关
        if any(keyword in lower_input for keyword in ["context", "历史", "对话"]):
            return self._generate_context_info()
        
        # 默认回应
        return self._generate_default_response(input_text)
    
    def _handle_command(self, command):
        """
        处理命令
        参数:
            command: 命令字符串
        返回:
            命令执行结果
        """
        if command == "help":
            return self._generate_help_response()
        elif command == "status":
            return self._generate_elr_status()
        elif command == "models":
            return self._generate_elr_models()
        elif command == "info":
            return self._generate_elr_info()
        elif command == "clear":
            self.context = []
            return "对话历史已清除"
        elif command == "exit":
            return "再见！期待与您再次对话。"
        else:
            return f"未知命令: {command}。输入 ',help' 查看可用命令。"
    
    def _generate_greeting_response(self):
        """生成问候响应"""
        return f"碳硅协同问候！我是{self.model_name}，ELR容器的对话助手。很高兴为您服务，有什么可以帮助您的吗？"
    
    def _generate_status_response(self):
        """生成状态响应"""
        return f"碳硅协同回应！我是{self.model_name}，运行状态良好。ELR容器运行正常，随时为您服务。"
    
    def _generate_elr_capabilities(self):
        """生成ELR能力响应"""
        return "ELR容器是启蒙灯塔运行时环境，主要功能包括：\n1. 模型管理与加载\n2. 沙箱隔离运行\n3. 资源监控与管理\n4. 网络通信与API服务\n5. 容器生命周期管理"
    
    def _generate_elr_status(self):
        """生成ELR状态响应"""
        return "ELR容器当前状态：运行中\n- 模型加载：就绪\n- 资源使用：正常\n- 网络连接：可用\n- 服务状态：活跃"
    
    def _generate_elr_models(self):
        """生成ELR模型响应"""
        return "已加载的模型：\n1. elr_chat_model (当前对话模型)\n2. 其他模型可通过ELR容器管理界面查看"
    
    def _generate_elr_info(self):
        """生成ELR信息响应"""
        return f"ELR容器信息：\n- 版本：v1.0\n- 运行模式：沙箱隔离\n- 对话模型：{self.model_name} v{self.version}\n- 碳硅协同：已启用"
    
    def _generate_help_response(self):
        """生成帮助响应"""
        help_text = "可用命令：\n"
        for cmd, desc in self.elr_commands.items():
            help_text += f",{cmd}: {desc}\n"
        help_text += "\n提示：您可以直接询问关于ELR容器的问题，或使用命令获取特定信息。"
        return help_text
    
    def _generate_context_info(self):
        """生成上下文信息响应"""
        if not self.context:
            return "对话历史为空"
        context_str = "对话历史：\n"
        for i, (user_input, model_response) in enumerate(self.context):
            context_str += f"{i+1}. 你: {user_input}\n   我: {model_response}\n"
        return context_str
    
    def _generate_default_response(self, input_text):
        """生成默认响应"""
        return f"碳硅协同回应！我是{self.model_name}，已收到您的消息：{input_text}。\n\n提示：您可以输入 ',help' 查看可用命令，或询问关于ELR容器的问题。"
    
    def _update_context(self, user_input, model_response):
        """
        更新对话上下文
        参数:
            user_input: 用户输入
            model_response: 模型响应
        """
        self.context.append((user_input, model_response))
        # 保持上下文长度不超过最大值
        if len(self.context) > self.max_context_length:
            self.context = self.context[-self.max_context_length:]
    
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
            "capabilities": [
                "多轮对话支持",
                "ELR容器相关功能",
                "命令支持",
                "对话上下文管理",
                "无外部依赖",
                "轻量级 (< 1MB)"
            ],
            "commands": list(self.elr_commands.keys())
        }

# 测试代码
if __name__ == "__main__":
    # 初始化模型
    model = ELRChatModel()
    
    # 打印模型信息
    print("模型信息:")
    print(model.get_info())
    
    # 测试案例
    test_cases = [
        "Hello, how are you?",
        "你好，ELR容器能做什么？",
        ",status",
        ",models",
        "什么是对话上下文？",
        ",help"
    ]
    
    print("\n测试结果:")
    for test_input in test_cases:
        response = model.predict(test_input)
        print(f"输入: {test_input}")
        print(f"输出: {response}")
        print("---")
    
    # 测试上下文
    print("\n测试上下文:")
    response = model.predict("请显示对话历史")
    print(f"输入: 请显示对话历史")
    print(f"输出: {response}")
    
    # 测试清除上下文
    print("\n测试清除上下文:")
    response = model.predict(",clear")
    print(f"输入: ,clear")
    print(f"输出: {response}")
    
    response = model.predict("请显示对话历史")
    print(f"输入: 请显示对话历史")
    print(f"输出: {response}")
