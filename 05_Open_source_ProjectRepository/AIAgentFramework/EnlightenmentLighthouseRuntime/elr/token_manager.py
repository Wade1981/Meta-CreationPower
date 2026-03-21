#!/usr/bin/env python3
"""
ELR容器令牌管理器
负责生成、验证和管理ELR容器的访问令牌
"""

import json
import time
import os
import hashlib
import uuid

class TokenManager:
    """令牌管理器"""
    
    def __init__(self, token_file="elr_token.json"):
        """初始化令牌管理器"""
        self.token_file = token_file
        self.token_data = self.load_tokens()
    
    def load_tokens(self):
        """从文件加载令牌数据"""
        if os.path.exists(self.token_file):
            try:
                with open(self.token_file, 'r', encoding='utf-8') as f:
                    return json.load(f)
            except Exception as e:
                print(f"Error loading token file: {e}")
        return {
            "tokens": [],
            "last_updated": time.time()
        }
    
    def save_tokens(self):
        """保存令牌数据到文件"""
        try:
            self.token_data["last_updated"] = time.time()
            with open(self.token_file, 'w', encoding='utf-8') as f:
                json.dump(self.token_data, f, indent=2, ensure_ascii=False)
            return True
        except Exception as e:
            print(f"Error saving token file: {e}")
            return False
    
    def generate_token(self, description="ELR Container Token"):
        """生成新令牌"""
        token_id = str(uuid.uuid4())
        token_secret = hashlib.sha256(str(uuid.uuid4()).encode()).hexdigest()
        token = {
            "id": token_id,
            "secret": token_secret,
            "description": description,
            "created_at": time.time(),
            "expires_at": time.time() + (7 * 24 * 3600),  # 7天过期
            "status": "active"
        }
        
        self.token_data["tokens"].append(token)
        self.save_tokens()
        
        # 返回完整令牌（id + secret）
        return f"{token_id}.{token_secret}"
    
    def validate_token(self, token):
        """验证令牌"""
        if not token:
            return False, "Token is required"
        
        try:
            token_parts = token.split('.')
            if len(token_parts) != 2:
                return False, "Invalid token format"
            
            token_id, token_secret = token_parts
            
            for t in self.token_data["tokens"]:
                if t["id"] == token_id and t["secret"] == token_secret:
                    # 检查令牌是否过期
                    if time.time() > t["expires_at"]:
                        return False, "Token has expired"
                    # 检查令牌状态
                    if t["status"] != "active":
                        return False, "Token is not active"
                    return True, "Token is valid"
            
            return False, "Token not found"
        except Exception as e:
            return False, f"Error validating token: {e}"
    
    def refresh_token(self, old_token, description="Refreshed ELR Container Token"):
        """刷新令牌"""
        # 验证旧令牌
        valid, message = self.validate_token(old_token)
        if not valid:
            return None, message
        
        # 禁用旧令牌
        token_parts = old_token.split('.')
        token_id = token_parts[0]
        
        for i, t in enumerate(self.token_data["tokens"]):
            if t["id"] == token_id:
                self.token_data["tokens"][i]["status"] = "revoked"
                break
        
        # 生成新令牌
        new_token = self.generate_token(description)
        return new_token, "Token refreshed successfully"
    
    def list_tokens(self):
        """列出所有令牌"""
        tokens = []
        for t in self.token_data["tokens"]:
            token_info = {
                "id": t["id"],
                "description": t["description"],
                "created_at": t["created_at"],
                "expires_at": t["expires_at"],
                "status": t["status"],
                "expired": time.time() > t["expires_at"]
            }
            tokens.append(token_info)
        return tokens
    
    def revoke_token(self, token_id):
        """撤销令牌"""
        for i, t in enumerate(self.token_data["tokens"]):
            if t["id"] == token_id:
                self.token_data["tokens"][i]["status"] = "revoked"
                self.save_tokens()
                return True, "Token revoked successfully"
        return False, "Token not found"

if __name__ == "__main__":
    # 测试令牌管理器
    tm = TokenManager()
    
    # 生成新令牌
    token = tm.generate_token("Test Token")
    print(f"Generated token: {token}")
    
    # 验证令牌
    valid, message = tm.validate_token(token)
    print(f"Token validation: {valid}, {message}")
    
    # 列出令牌
    tokens = tm.list_tokens()
    print(f"Tokens: {json.dumps(tokens, indent=2)}")
    
    # 刷新令牌
    new_token, message = tm.refresh_token(token, "Refreshed Test Token")
    print(f"Refreshed token: {new_token}, {message}")
    
    # 验证新令牌
    valid, message = tm.validate_token(new_token)
    print(f"New token validation: {valid}, {message}")
    
    # 验证旧令牌
    valid, message = tm.validate_token(token)
    print(f"Old token validation: {valid}, {message}")
