import json
import hashlib
import time
from typing import Dict, Any, List, Optional

class FeatureAssetEngine:
    """特征资产引擎：负责将原始数据转化为特征资产"""
    
    def __init__(self, config: Dict[str, Any] = None):
        """初始化特征资产引擎"""
        self.config = config or {}
        self.federated_learning_enabled = self.config.get('federated_learning_enabled', True)
        self.blockchain_integration = self.config.get('blockchain_integration', True)
        self.feature_dimension = self.config.get('feature_dimension', 256)
        self.nodes = self.config.get('nodes', ['local_node'])
        self.feature_store = {}
    
    def generate_feature_vector(self, data: Dict[str, Any], data_type: str) -> List[float]:
        """生成特征向量"""
        # 基于数据内容生成特征向量
        # 实际应用中应该使用更复杂的特征提取算法
        content_str = str(data.get('content', ''))
        hash_value = hashlib.sha256(content_str.encode()).hexdigest()
        
        # 生成固定维度的特征向量
        feature_vector = []
        for i in range(self.feature_dimension):
            # 基于哈希值生成特征
            feature_value = (int(hash_value[i % len(hash_value)], 16) / 15.0 - 0.5) * 2
            feature_vector.append(feature_value)
        
        return feature_vector
    
    def extract_local_features(self, data: Dict[str, Any], node_id: str) -> Dict[str, Any]:
        """在本地节点提取特征"""
        data_type = data.get('type', 'unknown')
        feature_vector = self.generate_feature_vector(data, data_type)
        
        local_features = {
            'node_id': node_id,
            'data_type': data_type,
            'feature_vector': feature_vector,
            'timestamp': time.time(),
            'metadata': data.get('metadata', {})
        }
        
        return local_features
    
    def aggregate_features(self, local_features_list: List[Dict[str, Any]]) -> Dict[str, Any]:
        """聚合跨节点特征"""
        if not local_features_list:
            return {}
        
        # 简单的特征聚合（实际应用中应该使用更复杂的联邦学习算法）
        aggregated_vector = [0.0] * self.feature_dimension
        for local_features in local_features_list:
            vector = local_features.get('feature_vector', [])
            if len(vector) == self.feature_dimension:
                for i in range(self.feature_dimension):
                    aggregated_vector[i] += vector[i]
        
        # 归一化
        norm = sum(x**2 for x in aggregated_vector) ** 0.5
        if norm > 0:
            aggregated_vector = [x / norm for x in aggregated_vector]
        
        aggregated_features = {
            'aggregated_vector': aggregated_vector,
            'node_count': len(local_features_list),
            'timestamp': time.time(),
            'source_nodes': [f['node_id'] for f in local_features_list]
        }
        
        return aggregated_features
    
    def generate_blockchain_fingerprint(self, features: Dict[str, Any]) -> str:
        """生成区块链指纹"""
        # 生成包含特征信息的字符串
        fingerprint_data = {
            'feature_vector': features.get('aggregated_vector', []),
            'timestamp': features.get('timestamp', time.time()),
            'node_count': features.get('node_count', 1),
            'source_nodes': features.get('source_nodes', [])
        }
        
        # 生成哈希作为指纹
        fingerprint_str = json.dumps(fingerprint_data, sort_keys=True)
        fingerprint = hashlib.sha256(fingerprint_str.encode()).hexdigest()
        
        return fingerprint
    
    def create_feature_asset(self, data: Dict[str, Any]) -> Dict[str, Any]:
        """创建特征资产"""
        # 1. 在各节点提取本地特征
        local_features_list = []
        for node in self.nodes:
            local_features = self.extract_local_features(data, node)
            local_features_list.append(local_features)
        
        # 2. 聚合特征
        aggregated_features = self.aggregate_features(local_features_list)
        
        # 3. 生成区块链指纹
        fingerprint = self.generate_blockchain_fingerprint(aggregated_features)
        
        # 4. 创建特征资产
        feature_asset = {
            'asset_id': fingerprint,
            'type': 'feature_asset',
            'original_data_type': data.get('type', 'unknown'),
            'features': aggregated_features,
            'blockchain_fingerprint': fingerprint,
            'metadata': {
                'creation_time': time.time(),
                'data_size': len(str(data)),
                'source_nodes': self.nodes
            }
        }
        
        # 存储特征资产
        self.feature_store[fingerprint] = feature_asset
        
        return feature_asset
    
    def get_feature_asset(self, asset_id: str) -> Optional[Dict[str, Any]]:
        """获取特征资产"""
        return self.feature_store.get(asset_id)
    
    def validate_feature_asset(self, asset: Dict[str, Any]) -> bool:
        """验证特征资产"""
        required_fields = ['asset_id', 'type', 'features', 'blockchain_fingerprint', 'metadata']
        for field in required_fields:
            if field not in asset:
                return False
        
        # 验证区块链指纹
        computed_fingerprint = self.generate_blockchain_fingerprint(asset['features'])
        if computed_fingerprint != asset['blockchain_fingerprint']:
            return False
        
        return True
    
    def list_feature_assets(self) -> List[str]:
        """列出所有特征资产ID"""
        return list(self.feature_store.keys())
    
    def get_asset_count(self) -> int:
        """获取资产数量"""
        return len(self.feature_store)

class FederatedFeatureDistillation:
    """联邦特征蒸馏网络"""
    
    def __init__(self, nodes: List[str]):
        """初始化联邦特征蒸馏网络"""
        self.nodes = nodes
        self.distillation_rounds = 10
    
    def distill_features(self, local_features_list: List[Dict[str, Any]]) -> Dict[str, Any]:
        """执行特征蒸馏"""
        # 简单的特征蒸馏实现
        # 实际应用中应该使用更复杂的知识蒸馏算法
        
        if not local_features_list:
            return {}
        
        # 初始化蒸馏后的特征
        feature_dimension = len(local_features_list[0].get('feature_vector', []))
        distilled_vector = [0.0] * feature_dimension
        
        # 多轮蒸馏
        for round_idx in range(self.distillation_rounds):
            # 每轮更新权重
            weights = [1.0 / len(local_features_list)] * len(local_features_list)
            
            # 加权聚合
            for i, local_features in enumerate(local_features_list):
                vector = local_features.get('feature_vector', [])
                if len(vector) == feature_dimension:
                    for j in range(feature_dimension):
                        distilled_vector[j] += vector[j] * weights[i]
            
            # 归一化
            norm = sum(x**2 for x in distilled_vector) ** 0.5
            if norm > 0:
                distilled_vector = [x / norm for x in distilled_vector]
        
        distilled_features = {
            'distilled_vector': distilled_vector,
            'rounds': self.distillation_rounds,
            'node_count': len(local_features_list),
            'timestamp': time.time()
        }
        
        return distilled_features

class BlockchainFingerprintSystem:
    """区块链指纹系统"""
    
    def __init__(self, blockchain_config: Dict[str, Any] = None):
        """初始化区块链指纹系统"""
        self.config = blockchain_config or {}
        self.chain_id = self.config.get('chain_id', 'local_chain')
        self.network = self.config.get('network', 'local')
        self.fingerprints = {}
    
    def create_fingerprint(self, data: Dict[str, Any]) -> str:
        """创建指纹"""
        # 生成包含链信息的指纹数据
        fingerprint_data = {
            'chain_id': self.chain_id,
            'network': self.network,
            'data': data,
            'timestamp': time.time()
        }
        
        # 生成哈希作为指纹
        fingerprint_str = json.dumps(fingerprint_data, sort_keys=True)
        fingerprint = hashlib.sha256(fingerprint_str.encode()).hexdigest()
        
        # 存储指纹
        self.fingerprints[fingerprint] = fingerprint_data
        
        return fingerprint
    
    def verify_fingerprint(self, fingerprint: str, data: Dict[str, Any]) -> bool:
        """验证指纹"""
        # 重新计算指纹并比较
        test_fingerprint = self.create_fingerprint(data)
        return test_fingerprint == fingerprint
    
    def get_fingerprint_data(self, fingerprint: str) -> Optional[Dict[str, Any]]:
        """获取指纹数据"""
        return self.fingerprints.get(fingerprint)
    
    def list_fingerprints(self) -> List[str]:
        """列出所有指纹"""
        return list(self.fingerprints.keys())
