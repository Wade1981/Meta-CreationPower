import hashlib
import base64
import json
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import padding
from typing import Dict, Any, Optional

class EncryptionUtils:
    """数据加密工具类"""
    
    @staticmethod
    def generate_key(password: str, salt: str = 'default_salt', key_length: int = 32) -> bytes:
        """生成加密密钥
        
        Args:
            password: 密码
            salt: 盐值
            key_length: 密钥长度，默认为32字节（256位）
            
        Returns:
            加密密钥
        """
        # 使用PBKDF2生成密钥（简化实现）
        key = hashlib.pbkdf2_hmac(
            'sha256',
            password.encode('utf-8'),
            salt.encode('utf-8'),
            100000,  # 迭代次数
            dklen=key_length
        )
        return key
    
    @staticmethod
    def encrypt_aes(data: str, key: bytes) -> str:
        """使用AES加密数据
        
        Args:
            data: 要加密的数据
            key: 加密密钥
            
        Returns:
            加密后的数据（base64编码）
        """
        if not data:
            return ""
        
        try:
            # 生成随机IV
            iv = hashlib.sha256(key).digest()[:16]
            
            # 填充数据
            padder = padding.PKCS7(128).padder()
            padded_data = padder.update(data.encode('utf-8')) + padder.finalize()
            
            # 创建加密器
            cipher = Cipher(
                algorithms.AES(key),
                modes.CBC(iv),
                backend=default_backend()
            )
            encryptor = cipher.encryptor()
            
            # 加密数据
            encrypted_data = encryptor.update(padded_data) + encryptor.finalize()
            
            # 组合IV和加密数据
            combined = iv + encrypted_data
            
            # 转换为base64编码
            encoded_data = base64.b64encode(combined).decode('utf-8')
            
            return encoded_data
        except Exception as e:
            print(f"Encryption error: {e}")
            return ""
    
    @staticmethod
    def decrypt_aes(encrypted_data: str, key: bytes) -> str:
        """使用AES解密数据
        
        Args:
            encrypted_data: 加密后的数据（base64编码）
            key: 加密密钥
            
        Returns:
            解密后的数据
        """
        if not encrypted_data:
            return ""
        
        try:
            # 解码base64
            decoded_data = base64.b64decode(encrypted_data.encode('utf-8'))
            
            # 提取IV
            iv = decoded_data[:16]
            ciphertext = decoded_data[16:]
            
            # 创建解密器
            cipher = Cipher(
                algorithms.AES(key),
                modes.CBC(iv),
                backend=default_backend()
            )
            decryptor = cipher.decryptor()
            
            # 解密数据
            padded_data = decryptor.update(ciphertext) + decryptor.finalize()
            
            # 移除填充
            unpadder = padding.PKCS7(128).unpadder()
            data = unpadder.update(padded_data) + unpadder.finalize()
            
            return data.decode('utf-8')
        except Exception as e:
            print(f"Decryption error: {e}")
            return ""
    
    @staticmethod
    def encrypt_json(data: Dict[str, Any], password: str) -> str:
        """加密JSON数据
        
        Args:
            data: 要加密的JSON数据
            password: 加密密码
            
        Returns:
            加密后的数据（base64编码）
        """
        try:
            # 转换为JSON字符串
            json_str = json.dumps(data, ensure_ascii=False)
            
            # 生成密钥
            key = EncryptionUtils.generate_key(password)
            
            # 加密数据
            encrypted_data = EncryptionUtils.encrypt_aes(json_str, key)
            
            return encrypted_data
        except Exception as e:
            print(f"JSON encryption error: {e}")
            return ""
    
    @staticmethod
    def decrypt_json(encrypted_data: str, password: str) -> Optional[Dict[str, Any]]:
        """解密JSON数据
        
        Args:
            encrypted_data: 加密后的数据（base64编码）
            password: 加密密码
            
        Returns:
            解密后的JSON数据
        """
        try:
            # 生成密钥
            key = EncryptionUtils.generate_key(password)
            
            # 解密数据
            decrypted_str = EncryptionUtils.decrypt_aes(encrypted_data, key)
            
            # 转换为JSON对象
            if decrypted_str:
                return json.loads(decrypted_str)
            return None
        except Exception as e:
            print(f"JSON decryption error: {e}")
            return None
    
    @staticmethod
    def hash_data(data: str, algorithm: str = 'sha256') -> str:
        """哈希数据
        
        Args:
            data: 要哈希的数据
            algorithm: 哈希算法，默认为sha256
            
        Returns:
            哈希值
        """
        if algorithm == 'sha256':
            return hashlib.sha256(data.encode('utf-8')).hexdigest()
        elif algorithm == 'sha512':
            return hashlib.sha512(data.encode('utf-8')).hexdigest()
        elif algorithm == 'md5':
            return hashlib.md5(data.encode('utf-8')).hexdigest()
        else:
            return hashlib.sha256(data.encode('utf-8')).hexdigest()
