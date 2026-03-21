"""Security utility module for ELR Service Load."""

import os
import hashlib
import hmac
from datetime import datetime, timedelta
import base64
import json


class SecurityManager:
    """Security manager for handling authentication and encryption."""

    def __init__(self, secret_key=None):
        """
        Initialize security manager.
        
        Args:
            secret_key (str): Secret key for signing tokens
        """
        self.secret_key = secret_key or os.environ.get('ELR_SECRET_KEY', 'default_secret_key')
        self.token_expiry = int(os.environ.get('ELR_TOKEN_EXPIRY', '3600'))  # Default 1 hour

    def generate_token(self, user_id):
        """
        Generate JWT-like token for authentication.
        
        Args:
            user_id (str): User ID
            
        Returns:
            str: Generated token
        """
        # Create payload
        payload = {
            'user_id': user_id,
            'exp': (datetime.utcnow() + timedelta(seconds=self.token_expiry)).timestamp(),
            'iat': datetime.utcnow().timestamp()
        }
        
        # Encode payload
        payload_json = json.dumps(payload)
        payload_encoded = base64.b64encode(payload_json.encode('utf-8')).decode('utf-8')
        
        # Generate signature
        signature = hmac.new(
            self.secret_key.encode('utf-8'),
            payload_encoded.encode('utf-8'),
            hashlib.sha256
        ).hexdigest()
        
        # Combine parts
        token = f"{payload_encoded}.{signature}"
        return token

    def validate_token(self, token):
        """
        Validate token for authentication.
        
        Args:
            token (str): Token to validate
            
        Returns:
            dict: Decoded payload if valid, None otherwise
        """
        try:
            # Split token parts
            payload_encoded, signature = token.split('.')
            
            # Verify signature
            expected_signature = hmac.new(
                self.secret_key.encode('utf-8'),
                payload_encoded.encode('utf-8'),
                hashlib.sha256
            ).hexdigest()
            
            if not hmac.compare_digest(signature, expected_signature):
                return None
            
            # Decode payload
            payload_json = base64.b64decode(payload_encoded).decode('utf-8')
            payload = json.loads(payload_json)
            
            # Check expiration
            if datetime.utcnow().timestamp() > payload.get('exp', 0):
                return None
            
            return payload
        except Exception as e:
            return None

    def hash_password(self, password):
        """
        Hash password for storage.
        
        Args:
            password (str): Plain text password
            
        Returns:
            str: Hashed password
        """
        salt = os.urandom(16).hex()
        hash_obj = hashlib.pbkdf2_hmac(
            'sha256',
            password.encode('utf-8'),
            salt.encode('utf-8'),
            100000
        )
        return f"{salt}${hash_obj.hex()}"

    def verify_password(self, password, hashed_password):
        """
        Verify password against stored hash.
        
        Args:
            password (str): Plain text password
            hashed_password (str): Stored hashed password
            
        Returns:
            bool: True if password is correct, False otherwise
        """
        try:
            salt, hash_part = hashed_password.split('$')
            hash_obj = hashlib.pbkdf2_hmac(
                'sha256',
                password.encode('utf-8'),
                salt.encode('utf-8'),
                100000
            )
            return hmac.compare_digest(hash_obj.hex(), hash_part)
        except Exception as e:
            return False

    def get_ssl_context(self, cert_file=None, key_file=None):
        """
        Get SSL context for HTTPS server.
        
        Args:
            cert_file (str): Path to certificate file
            key_file (str): Path to private key file
            
        Returns:
            ssl.SSLContext: SSL context if files provided, None otherwise
        """
        if cert_file and key_file and os.path.exists(cert_file) and os.path.exists(key_file):
            import ssl
            context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
            context.load_cert_chain(cert_file, key_file)
            return context
        return None
