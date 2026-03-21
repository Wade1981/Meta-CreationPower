import json
import hashlib
import time
from typing import Dict, Any, List, Optional

class TradingEngine:
    """流通交易引擎：负责数字资产的流通交易和全生命周期管理"""
    
    def __init__(self, config: Dict[str, Any] = None):
        """初始化流通交易引擎"""
        self.config = config or {}
        self.transaction_dir = self.config.get('transaction_dir', 'transactions')
        self.audit_log_dir = self.config.get('audit_log_dir', 'audit_logs')
        self.transactions = {}
        self.audit_logs = {}
        self.asset_status = {}
        self.risk_monitor = RiskIntelligenceMonitor()
        self.lifecycle_auditor = LifecycleAuditor()
    
    def create_transaction(self, asset_id: str, buyer: str, seller: str, price: float, transaction_type: str = 'purchase') -> Dict[str, Any]:
        """创建交易"""
        # 生成交易ID
        transaction_id = hashlib.sha256((asset_id + buyer + seller + str(time.time())).encode()).hexdigest()
        
        # 创建交易记录
        transaction = {
            'transaction_id': transaction_id,
            'asset_id': asset_id,
            'buyer': buyer,
            'seller': seller,
            'price': price,
            'transaction_type': transaction_type,
            'status': 'pending',
            'created_at': time.time(),
            'updated_at': time.time()
        }
        
        # 存储交易
        self.transactions[transaction_id] = transaction
        
        # 记录审计日志
        self.lifecycle_auditor.log_transaction(transaction)
        
        # 监控风险
        risk_analysis = self.risk_monitor.analyze_transaction(transaction)
        if risk_analysis['risk_level'] == 'high':
            transaction['status'] = 'risk_hold'
            transaction['risk_analysis'] = risk_analysis
        
        return transaction
    
    def execute_transaction(self, transaction_id: str) -> Dict[str, Any]:
        """执行交易"""
        transaction = self.transactions.get(transaction_id)
        if not transaction:
            return {}
        
        if transaction['status'] == 'risk_hold':
            return transaction
        
        # 更新交易状态
        transaction['status'] = 'completed'
        transaction['updated_at'] = time.time()
        
        # 更新资产状态
        asset_id = transaction['asset_id']
        self.asset_status[asset_id] = {
            'owner': transaction['buyer'],
            'last_transaction': transaction_id,
            'last_updated': time.time(),
            'status': 'active'
        }
        
        # 记录审计日志
        self.lifecycle_auditor.log_transaction_execution(transaction)
        
        return transaction
    
    def cancel_transaction(self, transaction_id: str, reason: str) -> Dict[str, Any]:
        """取消交易"""
        transaction = self.transactions.get(transaction_id)
        if not transaction:
            return {}
        
        # 更新交易状态
        transaction['status'] = 'cancelled'
        transaction['reason'] = reason
        transaction['updated_at'] = time.time()
        
        # 记录审计日志
        self.lifecycle_auditor.log_transaction_cancellation(transaction, reason)
        
        return transaction
    
    def get_transaction(self, transaction_id: str) -> Optional[Dict[str, Any]]:
        """获取交易"""
        return self.transactions.get(transaction_id)
    
    def list_transactions(self, asset_id: str = None, status: str = None) -> List[Dict[str, Any]]:
        """列出交易"""
        transactions = []
        
        for transaction in self.transactions.values():
            if asset_id and transaction['asset_id'] != asset_id:
                continue
            if status and transaction['status'] != status:
                continue
            transactions.append(transaction)
        
        # 按时间排序
        transactions.sort(key=lambda x: x['created_at'], reverse=True)
        
        return transactions
    
    def get_asset_lifecycle(self, asset_id: str) -> Dict[str, Any]:
        """获取资产生命周期"""
        lifecycle = self.lifecycle_auditor.get_asset_lifecycle(asset_id)
        return lifecycle
    
    def monitor_asset_usage(self, asset_id: str, usage_data: Dict[str, Any]) -> Dict[str, Any]:
        """监控资产使用情况"""
        # 记录使用情况
        usage_record = {
            'asset_id': asset_id,
            'usage_data': usage_data,
            'timestamp': time.time()
        }
        
        # 记录审计日志
        self.lifecycle_auditor.log_asset_usage(usage_record)
        
        # 分析使用风险
        risk_analysis = self.risk_monitor.analyze_asset_usage(usage_record)
        
        return {
            'usage_record': usage_record,
            'risk_analysis': risk_analysis
        }
    
    def get_asset_status(self, asset_id: str) -> Optional[Dict[str, Any]]:
        """获取资产状态"""
        return self.asset_status.get(asset_id)

class LifecycleAuditor:
    """全生命周期审计系统"""
    
    def __init__(self):
        """初始化全生命周期审计系统"""
        self.asset_lifecycles = {}
    
    def log_transaction(self, transaction: Dict[str, Any]):
        """记录交易日志"""
        asset_id = transaction['asset_id']
        if asset_id not in self.asset_lifecycles:
            self.asset_lifecycles[asset_id] = {
                'asset_id': asset_id,
                'events': []
            }
        
        event = {
            'type': 'transaction_created',
            'transaction_id': transaction['transaction_id'],
            'details': transaction,
            'timestamp': time.time()
        }
        
        self.asset_lifecycles[asset_id]['events'].append(event)
    
    def log_transaction_execution(self, transaction: Dict[str, Any]):
        """记录交易执行日志"""
        asset_id = transaction['asset_id']
        if asset_id not in self.asset_lifecycles:
            self.asset_lifecycles[asset_id] = {
                'asset_id': asset_id,
                'events': []
            }
        
        event = {
            'type': 'transaction_executed',
            'transaction_id': transaction['transaction_id'],
            'details': transaction,
            'timestamp': time.time()
        }
        
        self.asset_lifecycles[asset_id]['events'].append(event)
    
    def log_transaction_cancellation(self, transaction: Dict[str, Any], reason: str):
        """记录交易取消日志"""
        asset_id = transaction['asset_id']
        if asset_id not in self.asset_lifecycles:
            self.asset_lifecycles[asset_id] = {
                'asset_id': asset_id,
                'events': []
            }
        
        event = {
            'type': 'transaction_cancelled',
            'transaction_id': transaction['transaction_id'],
            'reason': reason,
            'details': transaction,
            'timestamp': time.time()
        }
        
        self.asset_lifecycles[asset_id]['events'].append(event)
    
    def log_asset_usage(self, usage_record: Dict[str, Any]):
        """记录资产使用日志"""
        asset_id = usage_record['asset_id']
        if asset_id not in self.asset_lifecycles:
            self.asset_lifecycles[asset_id] = {
                'asset_id': asset_id,
                'events': []
            }
        
        event = {
            'type': 'asset_used',
            'usage_data': usage_record['usage_data'],
            'details': usage_record,
            'timestamp': time.time()
        }
        
        self.asset_lifecycles[asset_id]['events'].append(event)
    
    def get_asset_lifecycle(self, asset_id: str) -> Dict[str, Any]:
        """获取资产生命周期"""
        lifecycle = self.asset_lifecycles.get(asset_id, {
            'asset_id': asset_id,
            'events': []
        })
        
        # 按时间排序事件
        lifecycle['events'].sort(key=lambda x: x['timestamp'])
        
        return lifecycle
    
    def verify_asset_history(self, asset_id: str) -> bool:
        """验证资产历史"""
        lifecycle = self.asset_lifecycles.get(asset_id)
        if not lifecycle:
            return False
        
        # 简单验证，实际应用中应该使用更复杂的验证方法
        return len(lifecycle.get('events', [])) > 0

class RiskIntelligenceMonitor:
    """风险智能预警系统"""
    
    def __init__(self):
        """初始化风险智能预警系统"""
        self.risk_rules = {
            'price_anomaly': {
                'threshold': 10000,
                'description': 'Price too high'
            },
            'frequency_anomaly': {
                'threshold': 5,
                'description': 'Transaction frequency too high'
            },
            'usage_anomaly': {
                'threshold': 100,
                'description': 'Usage frequency too high'
            }
        }
        
        self.transaction_history = {}
        self.usage_history = {}
    
    def analyze_transaction(self, transaction: Dict[str, Any]) -> Dict[str, Any]:
        """分析交易风险"""
        risk_factors = []
        risk_level = 'low'
        
        # 价格异常检测
        if transaction['price'] > self.risk_rules['price_anomaly']['threshold']:
            risk_factors.append({
                'factor': 'price_anomaly',
                'description': self.risk_rules['price_anomaly']['description'],
                'value': transaction['price'],
                'threshold': self.risk_rules['price_anomaly']['threshold']
            })
            risk_level = 'medium'
        
        # 交易频率异常检测
        asset_id = transaction['asset_id']
        if asset_id not in self.transaction_history:
            self.transaction_history[asset_id] = []
        
        self.transaction_history[asset_id].append(transaction)
        
        # 检查最近交易频率
        recent_transactions = [t for t in self.transaction_history[asset_id] if time.time() - t['created_at'] < 86400]  # 24小时内
        if len(recent_transactions) > self.risk_rules['frequency_anomaly']['threshold']:
            risk_factors.append({
                'factor': 'frequency_anomaly',
                'description': self.risk_rules['frequency_anomaly']['description'],
                'value': len(recent_transactions),
                'threshold': self.risk_rules['frequency_anomaly']['threshold']
            })
            risk_level = 'high'
        
        # 综合风险评估
        if len(risk_factors) > 1:
            risk_level = 'high'
        
        return {
            'risk_level': risk_level,
            'risk_factors': risk_factors,
            'transaction_id': transaction['transaction_id'],
            'timestamp': time.time()
        }
    
    def analyze_asset_usage(self, usage_record: Dict[str, Any]) -> Dict[str, Any]:
        """分析资产使用风险"""
        risk_factors = []
        risk_level = 'low'
        
        asset_id = usage_record['asset_id']
        if asset_id not in self.usage_history:
            self.usage_history[asset_id] = []
        
        self.usage_history[asset_id].append(usage_record)
        
        # 检查使用频率
        recent_usage = [u for u in self.usage_history[asset_id] if time.time() - u['timestamp'] < 3600]  # 1小时内
        if len(recent_usage) > self.risk_rules['usage_anomaly']['threshold']:
            risk_factors.append({
                'factor': 'usage_anomaly',
                'description': self.risk_rules['usage_anomaly']['description'],
                'value': len(recent_usage),
                'threshold': self.risk_rules['usage_anomaly']['threshold']
            })
            risk_level = 'medium'
        
        # 检查使用模式
        usage_data = usage_record['usage_data']
        if usage_data.get('unauthorized_access'):
            risk_factors.append({
                'factor': 'unauthorized_access',
                'description': 'Unauthorized access detected',
                'value': True,
                'threshold': False
            })
            risk_level = 'high'
        
        return {
            'risk_level': risk_level,
            'risk_factors': risk_factors,
            'asset_id': asset_id,
            'timestamp': time.time()
        }
    
    def generate_risk_alert(self, risk_analysis: Dict[str, Any]) -> Dict[str, Any]:
        """生成风险预警"""
        if risk_analysis['risk_level'] == 'low':
            return {}
        
        alert = {
            'alert_id': hashlib.sha256((json.dumps(risk_analysis, sort_keys=True) + str(time.time())).encode()).hexdigest(),
            'risk_analysis': risk_analysis,
            'alert_level': risk_analysis['risk_level'],
            'recommended_actions': self._get_recommended_actions(risk_analysis),
            'created_at': time.time()
        }
        
        return alert
    
    def _get_recommended_actions(self, risk_analysis: Dict[str, Any]) -> List[str]:
        """获取推荐操作"""
        actions = []
        
        if risk_analysis['risk_level'] == 'medium':
            actions.append('Monitor transaction closely')
            actions.append('Verify participant identities')
        elif risk_analysis['risk_level'] == 'high':
            actions.append('Hold transaction')
            actions.append('Investigate participant backgrounds')
            actions.append('Review asset history')
        
        return actions
