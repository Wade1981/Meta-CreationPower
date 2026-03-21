import zlib
import base64
import json
from typing import Dict, Any, Optional

class CompressionUtils:
    """数据压缩工具类"""
    
    @staticmethod
    def compress_data(data: str, compression_level: int = 6) -> str:
        """压缩数据
        
        Args:
            data: 要压缩的数据
            compression_level: 压缩级别，1-9，默认为6
            
        Returns:
            压缩后的数据（base64编码）
        """
        if not data:
            return ""
        
        # 压缩数据
        compressed_data = zlib.compress(data.encode('utf-8'), compression_level)
        
        # 转换为base64编码
        encoded_data = base64.b64encode(compressed_data).decode('utf-8')
        
        return encoded_data
    
    @staticmethod
    def decompress_data(compressed_data: str) -> str:
        """解压缩数据
        
        Args:
            compressed_data: 压缩后的数据（base64编码）
            
        Returns:
            解压缩后的数据
        """
        if not compressed_data:
            return ""
        
        try:
            # 解码base64
            decoded_data = base64.b64decode(compressed_data.encode('utf-8'))
            
            # 解压缩
            decompressed_data = zlib.decompress(decoded_data).decode('utf-8')
            
            return decompressed_data
        except Exception as e:
            print(f"Decompression error: {e}")
            return ""
    
    @staticmethod
    def compress_json(data: Dict[str, Any], compression_level: int = 6) -> str:
        """压缩JSON数据
        
        Args:
            data: 要压缩的JSON数据
            compression_level: 压缩级别，1-9，默认为6
            
        Returns:
            压缩后的数据（base64编码）
        """
        json_str = json.dumps(data, ensure_ascii=False)
        return CompressionUtils.compress_data(json_str, compression_level)
    
    @staticmethod
    def decompress_json(compressed_data: str) -> Optional[Dict[str, Any]]:
        """解压缩JSON数据
        
        Args:
            compressed_data: 压缩后的数据（base64编码）
            
        Returns:
            解压缩后的JSON数据
        """
        try:
            decompressed_str = CompressionUtils.decompress_data(compressed_data)
            if decompressed_str:
                return json.loads(decompressed_str)
            return None
        except Exception as e:
            print(f"JSON decompression error: {e}")
            return None
    
    @staticmethod
    def calculate_compression_ratio(original_data: str, compressed_data: str) -> float:
        """计算压缩率
        
        Args:
            original_data: 原始数据
            compressed_data: 压缩后的数据
            
        Returns:
            压缩率（压缩后大小/原始大小）
        """
        if not original_data:
            return 0.0
        
        original_size = len(original_data.encode('utf-8'))
        compressed_size = len(compressed_data.encode('utf-8'))
        
        return compressed_size / original_size if original_size > 0 else 0.0
