import json
import hashlib
import time
from typing import Dict, Any, List, Optional

class DigitalAssetControlCenter:
    """数字资产管控中心：整合所有模块，提供统一的管理接口"""
    
    def __init__(self, config: Dict[str, Any] = None):
        """初始化数字资产管控中心"""
        self.config = config or {}
        self.center_id = hashlib.sha256(str(time.time()).encode()).hexdigest()
        self.assets = {}
        self.processes = {}
        self.workflows = {}
        self.dashboard = {
            'total_assets': 0,
            'active_processes': 0,
            'total_transactions': 0,
            'system_health': 'healthy'
        }
        
        # 导入各个模块
        from ..core.raw_material_engine.raw_material_engine import RawMaterialEngine
        from ..core.feature_asset_engine.feature_asset_engine import FeatureAssetEngine
        from ..core.algorithm_asset_engine.algorithm_asset_engine import AlgorithmAssetEngine
        from ..core.asset_packaging_engine.asset_packaging_engine import AssetPackagingEngine
        from ..core.trading_engine.trading_engine import TradingEngine
        from ..elr_integration.elr_integration import ELRIntegration
        from ..utils.compression.compression import CompressionUtils
        from ..utils.encryption.encryption import EncryptionUtils
        from ..utils.file_transfer.file_transfer import FileTransferUtils, CollaborativeNetwork
        from ..ide_integration.ide_integration import IDEIntegration
        
        # 初始化各个模块
        self.raw_material_engine = RawMaterialEngine()
        self.feature_asset_engine = FeatureAssetEngine()
        self.algorithm_asset_engine = AlgorithmAssetEngine()
        self.asset_packaging_engine = AssetPackagingEngine()
        self.trading_engine = TradingEngine()
        self.elr_integration = ELRIntegration()
        self.compression_utils = CompressionUtils()
        self.encryption_utils = EncryptionUtils()
        self.file_transfer_utils = FileTransferUtils()
        self.collaborative_network = CollaborativeNetwork()
        self.ide_integration = IDEIntegration()
    
    def create_asset(self, asset_data: Dict[str, Any], asset_type: str = 'general') -> Dict[str, Any]:
        """创建数字资产
        
        Args:
            asset_data: 资产数据
            asset_type: 资产类型
            
        Returns:
            资产信息
        """
        # 生成资产ID
        asset_id = hashlib.sha256((str(asset_data) + asset_type + str(time.time())).encode()).hexdigest()
        
        # 处理原始数据
        processed_data = self.raw_material_engine.process_data(
            asset_data.get('content', ''),
            asset_data.get('data_type', 'text'),
            asset_data.get('structure_type', 'unstructured')
        )
        
        # 提取特征资产
        feature_asset = self.feature_asset_engine.create_feature_asset(processed_data)
        
        # 创建算法资产
        algorithm_asset = self.algorithm_asset_engine.create_algorithm_asset(
            [feature_asset],
            asset_type
        )
        
        # 封装资产
        market_data = asset_data.get('market_data', {})
        usage_rules = asset_data.get('usage_rules', {})
        packaged_asset = self.asset_packaging_engine.package_asset(
            [feature_asset],
            [algorithm_asset],
            market_data,
            usage_rules
        )
        
        # 部署到ELR容器
        container = self.elr_integration.create_asset_container(asset_id)
        deployment = self.elr_integration.deploy_asset_to_container(
            asset_id,
            container['container_id'],
            packaged_asset
        )
        
        # 完整资产信息
        asset = {
            'asset_id': asset_id,
            'asset_type': asset_type,
            'processed_data': processed_data,
            'feature_asset': feature_asset,
            'algorithm_asset': algorithm_asset,
            'packaged_asset': packaged_asset,
            'container': container,
            'deployment': deployment,
            'created_at': time.time(),
            'last_modified': time.time(),
            'status': 'active'
        }
        
        # 存储资产
        self.assets[asset_id] = asset
        
        # 更新仪表盘
        self._update_dashboard()
        
        return asset
    
    def get_asset(self, asset_id: str) -> Dict[str, Any]:
        """获取资产
        
        Args:
            asset_id: 资产ID
            
        Returns:
            资产信息
        """
        asset = self.assets.get(asset_id)
        if not asset:
            return {
                'success': False,
                'message': 'Asset not found'
            }
        
        return asset
    
    def update_asset(self, asset_id: str, update_data: Dict[str, Any]) -> Dict[str, Any]:
        """更新资产
        
        Args:
            asset_id: 资产ID
            update_data: 更新数据
            
        Returns:
            更新结果
        """
        asset = self.assets.get(asset_id)
        if not asset:
            return {
                'success': False,
                'message': 'Asset not found'
            }
        
        # 更新资产数据
        asset.update(update_data)
        asset['last_modified'] = time.time()
        
        return {
            'success': True,
            'message': 'Asset updated successfully',
            'asset_id': asset_id
        }
    
    def delete_asset(self, asset_id: str) -> Dict[str, Any]:
        """删除资产
        
        Args:
            asset_id: 资产ID
            
        Returns:
            删除结果
        """
        asset = self.assets.get(asset_id)
        if not asset:
            return {
                'success': False,
                'message': 'Asset not found'
            }
        
        # 停止容器
        container_id = asset.get('container', {}).get('container_id')
        if container_id:
            self.elr_integration.stop_container(container_id)
        
        # 删除资产
        del self.assets[asset_id]
        
        # 更新仪表盘
        self._update_dashboard()
        
        return {
            'success': True,
            'message': 'Asset deleted successfully',
            'asset_id': asset_id
        }
    
    def list_assets(self, asset_type: str = None, status: str = None) -> List[Dict[str, Any]]:
        """列出资产
        
        Args:
            asset_type: 资产类型（可选）
            status: 资产状态（可选）
            
        Returns:
            资产列表
        """
        assets = []
        
        for asset in self.assets.values():
            if asset_type and asset['asset_type'] != asset_type:
                continue
            if status and asset['status'] != status:
                continue
            assets.append(asset)
        
        # 按创建时间排序
        assets.sort(key=lambda x: x['created_at'], reverse=True)
        
        return assets
    
    def execute_workflow(self, workflow_id: str, workflow_data: Dict[str, Any]) -> Dict[str, Any]:
        """执行工作流
        
        Args:
            workflow_id: 工作流ID
            workflow_data: 工作流数据
            
        Returns:
            执行结果
        """
        # 生成工作流执行ID
        execution_id = hashlib.sha256((workflow_id + str(time.time())).encode()).hexdigest()
        
        # 执行工作流（示例实现）
        execution = {
            'execution_id': execution_id,
            'workflow_id': workflow_id,
            'status': 'running',
            'started_at': time.time(),
            'workflow_data': workflow_data
        }
        
        # 存储执行记录
        self.processes[execution_id] = execution
        
        # 模拟工作流执行
        # 实际应用中应该根据工作流定义执行具体步骤
        execution['status'] = 'completed'
        execution['completed_at'] = time.time()
        execution['result'] = {'success': True, 'message': 'Workflow executed successfully'}
        
        return execution
    
    def get_dashboard(self) -> Dict[str, Any]:
        """获取仪表盘
        
        Returns:
            仪表盘数据
        """
        return self.dashboard
    
    def _update_dashboard(self):
        """更新仪表盘数据"""
        # 更新资产数量
        self.dashboard['total_assets'] = len(self.assets)
        
        # 更新活跃进程数
        active_processes = [p for p in self.processes.values() if p['status'] == 'running']
        self.dashboard['active_processes'] = len(active_processes)
        
        # 更新交易数量
        transactions = self.trading_engine.list_transactions()
        self.dashboard['total_transactions'] = len(transactions)
        
        # 更新系统健康状态
        self.dashboard['system_health'] = 'healthy'
        self.dashboard['last_updated'] = time.time()
    
    def get_system_status(self) -> Dict[str, Any]:
        """获取系统状态
        
        Returns:
            系统状态
        """
        status = {
            'center_id': self.center_id,
            'timestamp': time.time(),
            'dashboard': self.dashboard,
            'modules': {
                'raw_material_engine': 'active',
                'feature_asset_engine': 'active',
                'algorithm_asset_engine': 'active',
                'asset_packaging_engine': 'active',
                'trading_engine': 'active',
                'elr_integration': 'active',
                'compression_utils': 'active',
                'encryption_utils': 'active',
                'file_transfer_utils': 'active',
                'collaborative_network': 'active',
                'ide_integration': 'active'
            },
            'containers': self.elr_integration.list_containers(),
            'ide_status': self.ide_integration.detect_ide()
        }
        
        return status

class APIServer:
    """API服务器：提供RESTful API接口"""
    
    def __init__(self, control_center: DigitalAssetControlCenter):
        """初始化API服务器
        
        Args:
            control_center: 数字资产管控中心
        """
        self.control_center = control_center
        self.routes = {
            'GET /assets': self.get_assets,
            'GET /assets/{asset_id}': self.get_asset,
            'POST /assets': self.create_asset,
            'PUT /assets/{asset_id}': self.update_asset,
            'DELETE /assets/{asset_id}': self.delete_asset,
            'GET /dashboard': self.get_dashboard,
            'GET /system/status': self.get_system_status,
            'POST /workflows': self.execute_workflow
        }
    
    def get_assets(self, request: Dict[str, Any]) -> Dict[str, Any]:
        """获取资产列表
        
        Args:
            request: 请求数据
            
        Returns:
            响应数据
        """
        asset_type = request.get('query', {}).get('type')
        status = request.get('query', {}).get('status')
        
        assets = self.control_center.list_assets(asset_type, status)
        
        return {
            'success': True,
            'data': assets,
            'total': len(assets)
        }
    
    def get_asset(self, request: Dict[str, Any]) -> Dict[str, Any]:
        """获取资产
        
        Args:
            request: 请求数据
            
        Returns:
            响应数据
        """
        asset_id = request.get('params', {}).get('asset_id')
        asset = self.control_center.get_asset(asset_id)
        
        if 'success' in asset and not asset['success']:
            return asset
        
        return {
            'success': True,
            'data': asset
        }
    
    def create_asset(self, request: Dict[str, Any]) -> Dict[str, Any]:
        """创建资产
        
        Args:
            request: 请求数据
            
        Returns:
            响应数据
        """
        asset_data = request.get('body', {})
        asset_type = request.get('body', {}).get('type', 'general')
        
        asset = self.control_center.create_asset(asset_data, asset_type)
        
        return {
            'success': True,
            'data': asset,
            'asset_id': asset['asset_id']
        }
    
    def update_asset(self, request: Dict[str, Any]) -> Dict[str, Any]:
        """更新资产
        
        Args:
            request: 请求数据
            
        Returns:
            响应数据
        """
        asset_id = request.get('params', {}).get('asset_id')
        update_data = request.get('body', {})
        
        result = self.control_center.update_asset(asset_id, update_data)
        
        return result
    
    def delete_asset(self, request: Dict[str, Any]) -> Dict[str, Any]:
        """删除资产
        
        Args:
            request: 请求数据
            
        Returns:
            响应数据
        """
        asset_id = request.get('params', {}).get('asset_id')
        result = self.control_center.delete_asset(asset_id)
        
        return result
    
    def get_dashboard(self, request: Dict[str, Any]) -> Dict[str, Any]:
        """获取仪表盘
        
        Args:
            request: 请求数据
            
        Returns:
            响应数据
        """
        dashboard = self.control_center.get_dashboard()
        
        return {
            'success': True,
            'data': dashboard
        }
    
    def get_system_status(self, request: Dict[str, Any]) -> Dict[str, Any]:
        """获取系统状态
        
        Args:
            request: 请求数据
            
        Returns:
            响应数据
        """
        status = self.control_center.get_system_status()
        
        return {
            'success': True,
            'data': status
        }
    
    def execute_workflow(self, request: Dict[str, Any]) -> Dict[str, Any]:
        """执行工作流
        
        Args:
            request: 请求数据
            
        Returns:
            响应数据
        """
        workflow_id = request.get('body', {}).get('workflow_id')
        workflow_data = request.get('body', {})
        
        result = self.control_center.execute_workflow(workflow_id, workflow_data)
        
        return {
            'success': True,
            'data': result
        }
