# RootPulseOS Cultural Translator Implementation

"""RootPulseOS文化翻译器实现，实现碳硅之间的文化与语言桥梁。"""

import logging
from typing import Dict, List, Optional, Any

class CulturalTranslator:
    """文化翻译器类，实现碳硅之间的文化与语言桥梁。"""
    
    def __init__(self):
        """初始化CulturalTranslator实例。"""
        self.logger = logging.getLogger(__name__)
        self.logger.info("Initializing CulturalTranslator...")
        
        self.translators = {}
        self.running = False
    
    def register_translator(self, name: str, translator):
        """注册翻译器。
        
        Args:
            name: 翻译器名称
            translator: 翻译器实例
        """
        self.translators[name] = translator
        self.logger.info(f"Registered translator: {name}")
    
    def start(self):
        """启动文化翻译器。"""
        self.logger.info("Starting CulturalTranslator...")
        self.running = True
        
        # 启动所有翻译器
        for name, translator in self.translators.items():
            if hasattr(translator, "start"):
                try:
                    translator.start()
                    self.logger.info(f"Started translator: {name}")
                except Exception as e:
                    self.logger.error(f"Failed to start translator {name}: {e}")
    
    def stop(self):
        """停止文化翻译器。"""
        self.logger.info("Stopping CulturalTranslator...")
        self.running = False
        
        # 停止所有翻译器
        for name, translator in reversed(list(self.translators.items())):
            if hasattr(translator, "stop"):
                try:
                    translator.stop()
                    self.logger.info(f"Stopped translator: {name}")
                except Exception as e:
                    self.logger.error(f"Failed to stop translator {name}: {e}")
    
    def translate(self, content: Any, source: str, target: str) -> Any:
        """翻译内容。
        
        Args:
            content: 要翻译的内容
            source: 源语言/文化
            target: 目标语言/文化
            
        Returns:
            Any: 翻译后的内容
        """
        if not self.running:
            self.logger.warning("CulturalTranslator is not running")
            return content
        
        # 找到合适的翻译器
        for name, translator in self.translators.items():
            if hasattr(translator, "can_translate") and translator.can_translate(source, target):
                try:
                    return translator.translate(content, source, target)
                except Exception as e:
                    self.logger.error(f"Failed to translate with {name}: {e}")
        
        self.logger.warning(f"No translator found for {source} to {target}")
        return content
    
    def get_translator(self, name: str):
        """获取指定翻译器。
        
        Args:
            name: 翻译器名称
            
        Returns:
            翻译器实例或None
        """
        return self.translators.get(name)
    
    def get_all_translators(self) -> Dict[str, object]:
        """获取所有翻译器。
        
        Returns:
            Dict[str, object]: 翻译器字典
        """
        return self.translators
    
    def status(self) -> Dict[str, any]:
        """获取文化翻译器状态。
        
        Returns:
            Dict[str, any]: 状态信息
        """
        return {
            "running": self.running,
            "translators": list(self.translators.keys())
        }