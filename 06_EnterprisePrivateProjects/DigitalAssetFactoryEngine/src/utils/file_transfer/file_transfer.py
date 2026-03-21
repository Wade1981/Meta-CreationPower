import json
import hashlib
import time
import socket
import threading
from typing import Dict, Any, List, Optional

class FileTransferUtils:
    """文件传输工具类"""
    
    @staticmethod
    def calculate_file_hash(file_path: str, algorithm: str = 'sha256') -> str:
        """计算文件哈希值
        
        Args:
            file_path: 文件路径
            algorithm: 哈希算法，默认为sha256
            
        Returns:
            文件哈希值
        """
        try:
            if algorithm == 'sha256':
                hash_obj = hashlib.sha256()
            elif algorithm == 'sha512':
                hash_obj = hashlib.sha512()
            else:
                hash_obj = hashlib.sha256()
            
            with open(file_path, 'rb') as f:
                while chunk := f.read(8192):
                    hash_obj.update(chunk)
            
            return hash_obj.hexdigest()
        except Exception as e:
            print(f"File hash calculation error: {e}")
            return ""
    
    @staticmethod
    def create_file_metadata(file_path: str, asset_id: str = None) -> Dict[str, Any]:
        """创建文件元数据
        
        Args:
            file_path: 文件路径
            asset_id: 资产ID
            
        Returns:
            文件元数据
        """
        import os
        
        metadata = {
            'file_path': file_path,
            'file_name': os.path.basename(file_path),
            'file_size': os.path.getsize(file_path) if os.path.exists(file_path) else 0,
            'file_hash': FileTransferUtils.calculate_file_hash(file_path),
            'asset_id': asset_id,
            'created_at': time.time()
        }
        
        return metadata
    
    @staticmethod
    def send_file(file_path: str, host: str, port: int) -> Dict[str, Any]:
        """发送文件
        
        Args:
            file_path: 文件路径
            host: 目标主机
            port: 目标端口
            
        Returns:
            发送结果
        """
        try:
            # 创建socket连接
            with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
                s.connect((host, port))
                
                # 发送文件元数据
                metadata = FileTransferUtils.create_file_metadata(file_path)
                metadata_json = json.dumps(metadata).encode('utf-8')
                s.sendall(len(metadata_json).to_bytes(4, byteorder='big'))
                s.sendall(metadata_json)
                
                # 发送文件数据
                with open(file_path, 'rb') as f:
                    while chunk := f.read(8192):
                        s.sendall(chunk)
                
                # 接收确认
                response = s.recv(1024)
                
                return {
                    'success': True,
                    'message': 'File sent successfully',
                    'metadata': metadata
                }
        except Exception as e:
            print(f"File send error: {e}")
            return {
                'success': False,
                'message': str(e)
            }
    
    @staticmethod
    def receive_file(save_dir: str, host: str, port: int) -> Dict[str, Any]:
        """接收文件
        
        Args:
            save_dir: 保存目录
            host: 监听主机
            port: 监听端口
            
        Returns:
            接收结果
        """
        try:
            import os
            
            # 确保保存目录存在
            os.makedirs(save_dir, exist_ok=True)
            
            # 创建socket服务器
            with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
                s.bind((host, port))
                s.listen(1)
                
                print(f"Waiting for file on {host}:{port}...")
                conn, addr = s.accept()
                print(f"Connected by {addr}")
                
                with conn:
                    # 接收元数据
                    metadata_len = int.from_bytes(conn.recv(4), byteorder='big')
                    metadata_json = conn.recv(metadata_len)
                    metadata = json.loads(metadata_json.decode('utf-8'))
                    
                    # 保存文件
                    file_path = os.path.join(save_dir, metadata['file_name'])
                    with open(file_path, 'wb') as f:
                        while True:
                            data = conn.recv(8192)
                            if not data:
                                break
                            f.write(data)
                    
                    # 发送确认
                    conn.sendall(b'File received successfully')
                    
                    return {
                        'success': True,
                        'message': 'File received successfully',
                        'metadata': metadata,
                        'saved_path': file_path
                    }
        except Exception as e:
            print(f"File receive error: {e}")
            return {
                'success': False,
                'message': str(e)
            }

class CollaborativeNetwork:
    """协同网络模块"""
    
    def __init__(self, config: Dict[str, Any] = None):
        """初始化协同网络模块"""
        self.config = config or {}
        self.nodes = self.config.get('nodes', [])
        self.node_id = self.config.get('node_id', hashlib.sha256(str(time.time()).encode()).hexdigest()[:8])
        self.shared_assets = {}
        self.lock = threading.Lock()
    
    def add_node(self, node_info: Dict[str, Any]) -> Dict[str, Any]:
        """添加节点
        
        Args:
            node_info: 节点信息
            
        Returns:
            添加结果
        """
        with self.lock:
            # 检查节点是否已存在
            for node in self.nodes:
                if node.get('node_id') == node_info.get('node_id'):
                    return {
                        'success': False,
                        'message': 'Node already exists'
                    }
            
            # 添加节点
            node_info['joined_at'] = time.time()
            self.nodes.append(node_info)
            
            return {
                'success': True,
                'message': 'Node added successfully',
                'node_info': node_info
            }
    
    def remove_node(self, node_id: str) -> Dict[str, Any]:
        """移除节点
        
        Args:
            node_id: 节点ID
            
        Returns:
            移除结果
        """
        with self.lock:
            # 查找并移除节点
            for i, node in enumerate(self.nodes):
                if node.get('node_id') == node_id:
                    removed_node = self.nodes.pop(i)
                    return {
                        'success': True,
                        'message': 'Node removed successfully',
                        'node_info': removed_node
                    }
            
            return {
                'success': False,
                'message': 'Node not found'
            }
    
    def share_asset(self, asset_id: str, asset_data: Dict[str, Any]) -> Dict[str, Any]:
        """共享资产
        
        Args:
            asset_id: 资产ID
            asset_data: 资产数据
            
        Returns:
            共享结果
        """
        with self.lock:
            # 存储共享资产
            self.shared_assets[asset_id] = {
                'asset_id': asset_id,
                'asset_data': asset_data,
                'shared_by': self.node_id,
                'shared_at': time.time(),
                'nodes': []
            }
            
            return {
                'success': True,
                'message': 'Asset shared successfully',
                'asset_id': asset_id
            }
    
    def get_shared_asset(self, asset_id: str) -> Dict[str, Any]:
        """获取共享资产
        
        Args:
            asset_id: 资产ID
            
        Returns:
            资产数据
        """
        with self.lock:
            asset = self.shared_assets.get(asset_id)
            if not asset:
                return {
                    'success': False,
                    'message': 'Asset not found'
                }
            
            return {
                'success': True,
                'asset_data': asset['asset_data'],
                'shared_by': asset['shared_by'],
                'shared_at': asset['shared_at']
            }
    
    def list_shared_assets(self) -> List[Dict[str, Any]]:
        """列出共享资产
        
        Returns:
            共享资产列表
        """
        with self.lock:
            assets = []
            for asset_id, asset in self.shared_assets.items():
                assets.append({
                    'asset_id': asset_id,
                    'shared_by': asset['shared_by'],
                    'shared_at': asset['shared_at']
                })
            return assets
    
    def list_nodes(self) -> List[Dict[str, Any]]:
        """列出节点
        
        Returns:
            节点列表
        """
        return self.nodes
