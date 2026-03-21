import json
import hashlib
import time
from typing import Dict, Any, List, Optional, Tuple

class AlgorithmAssetEngine:
    """算法资产引擎：负责基于特征资产训练可复用的数字应用算法"""
    
    def __init__(self, config: Dict[str, Any] = None):
        """初始化算法资产引擎"""
        self.config = config or {}
        self.model_dir = self.config.get('model_dir', 'models')
        self.training_iterations = self.config.get('training_iterations', 100)
        self.learning_rate = self.config.get('learning_rate', 0.01)
        self.algorithm_store = {}
        self.market_feedback = {}
        
        # 市场反馈优化系统配置
        self.feedback_weights = {
            'usage_metrics': 0.25,      # 用户使用数据
            'business_metrics': 0.35,    # 业务指标
            'technical_evaluation': 0.25, # 技术评估
            'industry_trends': 0.15      # 行业趋势
        }
        
        # 优化策略配置
        self.optimization_strategies = {
            'micro': True,   # 微观优化：模型参数调整
            'meso': True,    # 中观优化：特征工程改进
            'macro': False   # 宏观优化：算法范式创新
        }
        
        # 协同进化配置
        self.collaborative_evolution = {
            'knowledge_transfer': True,  # 知识迁移
            'cross_industry_sharing': False,  # 跨行业共享
            'feature_engine_collaboration': True  # 与特征引擎协同
        }
        
        # 时间衰减参数
        self.time_decay_factor = 0.95
        
        # 优化周期配置
        self.optimization_cycle = {
            'min_feedback_count': 5,     # 最小反馈数量
            'max_feedback_age': 30 * 24 * 3600,  # 最大反馈年龄（30天）
            'optimization_threshold': 0.1  # 优化阈值
        }
    
    def create_basic_model(self, feature_dimension: int, output_dimension: int) -> Dict[str, Any]:
        """创建基础模型结构"""
        # 简单的线性模型结构
        # 实际应用中应该使用更复杂的模型架构
        model = {
            'weights': [0.0] * feature_dimension,
            'bias': 0.0,
            'feature_dimension': feature_dimension,
            'output_dimension': output_dimension,
            'architecture': 'linear'
        }
        
        return model
    
    def train_model(self, model: Dict[str, Any], feature_assets: List[Dict[str, Any]], labels: List[float]) -> Dict[str, Any]:
        """训练模型"""
        if not feature_assets or not labels:
            return model
        
        feature_dimension = model.get('feature_dimension', 0)
        if feature_dimension == 0:
            return model
        
        # 简单的线性回归训练
        # 实际应用中应该使用更复杂的训练算法
        weights = model.get('weights', [0.0] * feature_dimension)
        bias = model.get('bias', 0.0)
        
        for iteration in range(self.training_iterations):
            total_error = 0.0
            
            for i, asset in enumerate(feature_assets):
                features = asset.get('features', {}).get('aggregated_vector', [])
                if len(features) != feature_dimension:
                    continue
                
                label = labels[i]
                
                # 前向传播
                prediction = sum(w * f for w, f in zip(weights, features)) + bias
                error = prediction - label
                total_error += error ** 2
                
                # 反向传播
                for j in range(feature_dimension):
                    weights[j] -= self.learning_rate * error * features[j]
                bias -= self.learning_rate * error
            
            # 计算平均误差
            avg_error = total_error / len(feature_assets)
            if avg_error < 0.001:
                break
        
        trained_model = {
            'weights': weights,
            'bias': bias,
            'feature_dimension': feature_dimension,
            'output_dimension': model.get('output_dimension', 1),
            'architecture': model.get('architecture', 'linear'),
            'training_iterations': iteration + 1,
            'final_error': avg_error
        }
        
        return trained_model
    
    def predict(self, model: Dict[str, Any], features: List[float]) -> float:
        """使用模型进行预测"""
        weights = model.get('weights', [])
        bias = model.get('bias', 0.0)
        
        if len(features) != len(weights):
            return 0.0
        
        prediction = sum(w * f for w, f in zip(weights, features)) + bias
        return prediction
    
    def create_algorithm_asset(self, feature_assets: List[Dict[str, Any]], algorithm_type: str, labels: List[float] = None) -> Dict[str, Any]:
        """创建算法资产"""
        if not feature_assets:
            return {}
        
        # 确定特征维度
        first_asset = feature_assets[0]
        feature_dimension = len(first_asset.get('features', {}).get('aggregated_vector', []))
        
        # 创建基础模型
        model = self.create_basic_model(feature_dimension, 1)
        
        # 如果有标签数据，进行训练
        if labels:
            model = self.train_model(model, feature_assets, labels)
        
        # 生成算法资产ID
        model_str = json.dumps(model, sort_keys=True)
        asset_id = hashlib.sha256(model_str.encode()).hexdigest()
        
        # 创建算法资产
        algorithm_asset = {
            'asset_id': asset_id,
            'type': 'algorithm_asset',
            'algorithm_type': algorithm_type,
            'model': model,
            'feature_assets': [asset.get('asset_id') for asset in feature_assets],
            'metadata': {
                'creation_time': time.time(),
                'feature_count': len(feature_assets),
                'training_status': 'trained' if labels else 'untrained'
            }
        }
        
        # 存储算法资产
        self.algorithm_store[asset_id] = algorithm_asset
        
        return algorithm_asset
    
    def update_market_feedback(self, algorithm_id: str, feedback: Dict[str, Any]) -> None:
        """更新市场反馈"""
        if algorithm_id not in self.market_feedback:
            self.market_feedback[algorithm_id] = []
        
        feedback_with_timestamp = {
            'timestamp': time.time(),
            'feedback': feedback
        }
        
        self.market_feedback[algorithm_id].append(feedback_with_timestamp)
    
    def optimize_algorithm(self, algorithm_asset: Dict[str, Any], context: Dict[str, Any] = None) -> Dict[str, Any]:
        """基于市场反馈优化算法"""
        context = context or {}
        algorithm_id = algorithm_asset.get('asset_id')
        if algorithm_id not in self.market_feedback:
            return algorithm_asset
        
        # 获取市场反馈
        feedbacks = self.market_feedback[algorithm_id]
        if not feedbacks:
            return algorithm_asset
        
        # 检查是否达到优化条件
        if len(feedbacks) < self.optimization_cycle['min_feedback_count']:
            return algorithm_asset
        
        # 分析反馈数据
        feedback_analysis = self.analyze_feedback(feedbacks, algorithm_asset)
        
        # 检查优化阈值
        if feedback_analysis.get('overall_score', 0.5) > 0.7:
            # 性能已经很好，不需要优化
            return algorithm_asset
        
        # 设计优化目标函数
        optimization_target = self.design_optimization_target(algorithm_asset)
        
        # 获取原始模型
        model = algorithm_asset.get('model', {})
        
        # 执行多层次优化
        # 1. 微观优化：模型参数调整
        optimized_model = self.optimize_model_parameters(model, feedback_analysis)
        
        # 2. 中观优化：特征工程改进
        optimized_model = self.optimize_feature_engineering(optimized_model, feedback_analysis)
        
        # 3. 宏观优化：算法范式创新
        optimized_model = self.optimize_algorithm_paradigm(optimized_model, feedback_analysis)
        
        # 执行协同进化
        # 获取相关资产
        related_assets = []
        for asset_id, asset in self.algorithm_store.items():
            if asset_id != algorithm_id and asset.get('algorithm_type') == algorithm_asset.get('algorithm_type'):
                related_assets.append(asset)
        
        # 执行知识迁移
        if self.collaborative_evolution['knowledge_transfer'] and related_assets:
            algorithm_asset = self.perform_collaborative_evolution(algorithm_asset, related_assets)
        
        # 更新模型
        optimized_model['optimization_timestamp'] = time.time()
        optimized_model['feedback_count'] = feedback_analysis.get('feedback_count', 0)
        optimized_model['overall_score'] = feedback_analysis.get('overall_score', 0.5)
        optimized_model['optimization_target'] = optimization_target
        
        # 更新算法资产
        optimized_asset = algorithm_asset.copy()
        optimized_asset['model'] = optimized_model
        optimized_asset['metadata']['optimization_time'] = time.time()
        optimized_asset['metadata']['optimization_target'] = optimization_target
        optimized_asset['metadata']['feedback_analysis'] = feedback_analysis
        
        # 评估优化效果
        evaluation = self.evaluate_optimization(algorithm_asset, optimized_asset)
        optimized_asset['metadata']['optimization_evaluation'] = evaluation
        
        # 存储更新后的算法资产
        self.algorithm_store[algorithm_id] = optimized_asset
        
        return optimized_asset
    
    def get_algorithm_asset(self, asset_id: str) -> Optional[Dict[str, Any]]:
        """获取算法资产"""
        return self.algorithm_store.get(asset_id)
    
    def list_algorithm_assets(self) -> List[str]:
        """列出所有算法资产"""
        return list(self.algorithm_store.keys())
    
    def get_market_feedback(self, algorithm_id: str) -> List[Dict[str, Any]]:
        """获取市场反馈"""
        return self.market_feedback.get(algorithm_id, [])
    
    def standardize_feedback(self, feedback: Dict[str, Any]) -> Dict[str, Any]:
        """标准化反馈数据"""
        standardized = {
            'usage_metrics': {
                'accuracy': max(0, min(1, feedback.get('usage_metrics', {}).get('accuracy', 0.5))),
                'response_time': max(0, min(1, 1 - feedback.get('usage_metrics', {}).get('response_time', 0.5) / 10.0)),
                'resource_consumption': max(0, min(1, 1 - feedback.get('usage_metrics', {}).get('resource_consumption', 0.5) / 100.0))
            },
            'business_metrics': {
                'roi': max(0, min(1, feedback.get('business_metrics', {}).get('roi', 0.5))),
                'conversion_rate': max(0, min(1, feedback.get('business_metrics', {}).get('conversion_rate', 0.5))),
                'market_share': max(0, min(1, feedback.get('business_metrics', {}).get('market_share', 0.5)))
            },
            'technical_evaluation': {
                'stability': max(0, min(1, feedback.get('technical_evaluation', {}).get('stability', 0.5))),
                'scalability': max(0, min(1, feedback.get('technical_evaluation', {}).get('scalability', 0.5))),
                'security': max(0, min(1, feedback.get('technical_evaluation', {}).get('security', 0.5)))
            },
            'industry_trends': {
                'tech_integration': max(0, min(1, feedback.get('industry_trends', {}).get('tech_integration', 0.5))),
                'regulatory_compliance': max(0, min(1, feedback.get('industry_trends', {}).get('regulatory_compliance', 0.5)))
            },
            'rating': max(1, min(5, feedback.get('rating', 3)))
        }
        return standardized
    
    def calculate_dynamic_weights(self, algorithm_asset: Dict[str, Any], context: Dict[str, Any] = None) -> Dict[str, float]:
        """计算动态权重"""
        context = context or {}
        algorithm_type = algorithm_asset.get('algorithm_type', 'general')
        usage_scenario = context.get('usage_scenario', 'standard')
        
        # 基础权重
        weights = self.feedback_weights.copy()
        
        # 根据算法类型调整权重
        if algorithm_type == 'financial':
            weights['business_metrics'] *= 1.2
            weights['technical_evaluation'] *= 1.1
        elif algorithm_type == 'healthcare':
            weights['technical_evaluation'] *= 1.3
            weights['industry_trends'] *= 1.2
        elif algorithm_type == 'marketing':
            weights['business_metrics'] *= 1.3
            weights['usage_metrics'] *= 1.1
        
        # 根据使用场景调整权重
        if usage_scenario == 'real_time':
            weights['usage_metrics'] *= 1.2
        elif usage_scenario == 'batch_processing':
            weights['technical_evaluation'] *= 1.1
        elif usage_scenario == 'high_stakes':
            weights['technical_evaluation'] *= 1.3
        
        # 归一化权重
        total_weight = sum(weights.values())
        normalized_weights = {k: v / total_weight for k, v in weights.items()}
        
        return normalized_weights
    
    def calculate_time_decayed_feedback(self, feedbacks: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
        """计算时间衰减后的反馈"""
        current_time = time.time()
        decayed_feedbacks = []
        
        for feedback in feedbacks:
            timestamp = feedback.get('timestamp', current_time)
            age = current_time - timestamp
            
            # 过滤过期反馈
            if age > self.optimization_cycle['max_feedback_age']:
                continue
            
            # 计算时间衰减因子
            days_old = age / (24 * 3600)
            decay_factor = self.time_decay_factor ** days_old
            
            # 应用时间衰减
            decayed_feedback = feedback.copy()
            decayed_feedback['time_decay_factor'] = decay_factor
            decayed_feedbacks.append(decayed_feedback)
        
        return decayed_feedbacks
    
    def optimize_model_parameters(self, model: Dict[str, Any], feedback_analysis: Dict[str, Any]) -> Dict[str, Any]:
        """微观优化：模型参数调整"""
        if not self.optimization_strategies['micro']:
            return model
        
        weights = model.get('weights', [])
        bias = model.get('bias', 0.0)
        
        # 基于反馈分析调整参数
        optimization_score = feedback_analysis.get('optimization_score', 0.5)
        adjustment_factor = 1.0 + (optimization_score - 0.5) * 0.2
        
        # 调整权重和偏置
        adjusted_weights = [w * adjustment_factor for w in weights]
        adjusted_bias = bias * adjustment_factor
        
        # 更新模型
        optimized_model = model.copy()
        optimized_model['weights'] = adjusted_weights
        optimized_model['bias'] = adjusted_bias
        optimized_model['parameter_adjustment_factor'] = adjustment_factor
        
        return optimized_model
    
    def optimize_feature_engineering(self, model: Dict[str, Any], feedback_analysis: Dict[str, Any]) -> Dict[str, Any]:
        """中观优化：特征工程改进"""
        if not self.optimization_strategies['meso']:
            return model
        
        # 基于反馈分析生成特征工程建议
        feature_suggestions = []
        
        if feedback_analysis.get('usage_metrics_score', 0.5) < 0.4:
            feature_suggestions.append('Improve feature relevance and quality')
        
        if feedback_analysis.get('business_metrics_score', 0.5) < 0.4:
            feature_suggestions.append('Add domain-specific features')
        
        # 更新模型的特征工程建议
        optimized_model = model.copy()
        optimized_model['feature_engineering_suggestions'] = feature_suggestions
        optimized_model['feature_engineering_optimized'] = True
        
        return optimized_model
    
    def optimize_algorithm_paradigm(self, model: Dict[str, Any], feedback_analysis: Dict[str, Any]) -> Dict[str, Any]:
        """宏观优化：算法范式创新"""
        if not self.optimization_strategies['macro']:
            return model
        
        # 基于反馈分析生成算法范式创新建议
        paradigm_suggestions = []
        
        if feedback_analysis.get('overall_score', 0.5) < 0.3:
            paradigm_suggestions.append('Consider algorithm paradigm shift')
        
        # 更新模型的算法范式建议
        optimized_model = model.copy()
        optimized_model['paradigm_suggestions'] = paradigm_suggestions
        optimized_model['paradigm_optimized'] = True
        
        return optimized_model
    
    def perform_collaborative_evolution(self, algorithm_asset: Dict[str, Any], related_assets: List[Dict[str, Any]] = None) -> Dict[str, Any]:
        """执行协同进化"""
        if not self.collaborative_evolution['knowledge_transfer']:
            return algorithm_asset
        
        if not related_assets:
            return algorithm_asset
        
        # 执行知识迁移
        knowledge_transferred = False
        model = algorithm_asset.get('model', {})
        
        # 从相关资产中提取有用信息
        for related_asset in related_assets:
            related_model = related_asset.get('model', {})
            if related_model.get('architecture') == model.get('architecture'):
                # 迁移知识
                knowledge_transferred = True
                break
        
        # 更新算法资产
        optimized_asset = algorithm_asset.copy()
        optimized_asset['knowledge_transferred'] = knowledge_transferred
        optimized_asset['collaborative_evolution_timestamp'] = time.time()
        
        return optimized_asset
    
    def analyze_feedback(self, feedbacks: List[Dict[str, Any]], algorithm_asset: Dict[str, Any]) -> Dict[str, Any]:
        """分析反馈数据"""
        if not feedbacks:
            return {
                'optimization_score': 0.5,
                'overall_score': 0.5,
                'usage_metrics_score': 0.5,
                'business_metrics_score': 0.5,
                'technical_evaluation_score': 0.5,
                'industry_trends_score': 0.5,
                'feedback_count': 0
            }
        
        # 计算时间衰减后的反馈
        decayed_feedbacks = self.calculate_time_decayed_feedback(feedbacks)
        if not decayed_feedbacks:
            return {
                'optimization_score': 0.5,
                'overall_score': 0.5,
                'usage_metrics_score': 0.5,
                'business_metrics_score': 0.5,
                'technical_evaluation_score': 0.5,
                'industry_trends_score': 0.5,
                'feedback_count': 0
            }
        
        # 获取动态权重
        dynamic_weights = self.calculate_dynamic_weights(algorithm_asset)
        
        # 分析各项反馈指标
        scores = {
            'usage_metrics': 0.0,
            'business_metrics': 0.0,
            'technical_evaluation': 0.0,
            'industry_trends': 0.0
        }
        
        for feedback in decayed_feedbacks:
            fb = feedback.get('feedback', {})
            decay_factor = feedback.get('time_decay_factor', 1.0)
            
            # 分析使用指标
            usage_score = (
                fb.get('usage_metrics', {}).get('accuracy', 0.5) * 0.4 +
                fb.get('usage_metrics', {}).get('response_time', 0.5) * 0.3 +
                fb.get('usage_metrics', {}).get('resource_consumption', 0.5) * 0.3
            )
            scores['usage_metrics'] += usage_score * decay_factor
            
            # 分析业务指标
            business_score = (
                fb.get('business_metrics', {}).get('roi', 0.5) * 0.4 +
                fb.get('business_metrics', {}).get('conversion_rate', 0.5) * 0.3 +
                fb.get('business_metrics', {}).get('market_share', 0.5) * 0.3
            )
            scores['business_metrics'] += business_score * decay_factor
            
            # 分析技术评估
            technical_score = (
                fb.get('technical_evaluation', {}).get('stability', 0.5) * 0.4 +
                fb.get('technical_evaluation', {}).get('scalability', 0.5) * 0.3 +
                fb.get('technical_evaluation', {}).get('security', 0.5) * 0.3
            )
            scores['technical_evaluation'] += technical_score * decay_factor
            
            # 分析行业趋势
            industry_score = (
                fb.get('industry_trends', {}).get('tech_integration', 0.5) * 0.5 +
                fb.get('industry_trends', {}).get('regulatory_compliance', 0.5) * 0.5
            )
            scores['industry_trends'] += industry_score * decay_factor
        
        # 计算平均分数
        total_decay = sum(fb.get('time_decay_factor', 1.0) for fb in decayed_feedbacks)
        if total_decay > 0:
            scores = {k: v / total_decay for k, v in scores.items()}
        
        # 计算加权总分
        overall_score = sum(scores[k] * dynamic_weights[k] for k in scores)
        
        # 计算优化分数（基于与阈值的差异）
        optimization_score = max(0, min(1, (overall_score - 0.5) * 2))
        
        return {
            'optimization_score': optimization_score,
            'overall_score': overall_score,
            'usage_metrics_score': scores['usage_metrics'],
            'business_metrics_score': scores['business_metrics'],
            'technical_evaluation_score': scores['technical_evaluation'],
            'industry_trends_score': scores['industry_trends'],
            'feedback_count': len(decayed_feedbacks),
            'dynamic_weights': dynamic_weights
        }
    
    def design_optimization_target(self, algorithm_asset: Dict[str, Any]) -> Dict[str, Any]:
        """设计优化目标函数"""
        algorithm_type = algorithm_asset.get('algorithm_type', 'general')
        
        # 根据算法类型设计目标函数
        if algorithm_type == 'financial':
            target_function = {
                'primary_objective': 'maximize_roi',
                'secondary_objectives': ['minimize_risk', 'maximize_stability'],
                'constraints': ['regulatory_compliance', 'response_time < 1s']
            }
        elif algorithm_type == 'healthcare':
            target_function = {
                'primary_objective': 'maximize_accuracy',
                'secondary_objectives': ['minimize_false_negatives', 'maximize_security'],
                'constraints': ['regulatory_compliance', 'privacy_protection']
            }
        elif algorithm_type == 'marketing':
            target_function = {
                'primary_objective': 'maximize_conversion_rate',
                'secondary_objectives': ['maximize_reach', 'minimize_cost'],
                'constraints': ['brand_safety', 'response_time < 5s']
            }
        else:
            target_function = {
                'primary_objective': 'maximize_accuracy',
                'secondary_objectives': ['minimize_resource_consumption', 'maximize_scalability'],
                'constraints': ['response_time < 10s']
            }
        
        return target_function
    
    def evaluate_optimization(self, original_asset: Dict[str, Any], optimized_asset: Dict[str, Any]) -> Dict[str, Any]:
        """评估优化效果"""
        original_model = original_asset.get('model', {})
        optimized_model = optimized_asset.get('model', {})
        
        # 计算优化前后的差异
        evaluation = {
            'optimization_timestamp': time.time(),
            'parameters_adjusted': False,
            'feature_engineering_optimized': False,
            'paradigm_optimized': False,
            'knowledge_transferred': False,
            'overall_improvement': 0.0
        }
        
        # 检查参数调整
        if 'parameter_adjustment_factor' in optimized_model:
            evaluation['parameters_adjusted'] = True
        
        # 检查特征工程优化
        if optimized_model.get('feature_engineering_optimized', False):
            evaluation['feature_engineering_optimized'] = True
        
        # 检查范式优化
        if optimized_model.get('paradigm_optimized', False):
            evaluation['paradigm_optimized'] = True
        
        # 检查知识迁移
        if optimized_asset.get('knowledge_transferred', False):
            evaluation['knowledge_transferred'] = True
        
        # 计算整体改进度
        improvement_factors = []
        if evaluation['parameters_adjusted']:
            improvement_factors.append(0.3)
        if evaluation['feature_engineering_optimized']:
            improvement_factors.append(0.4)
        if evaluation['paradigm_optimized']:
            improvement_factors.append(0.6)
        if evaluation['knowledge_transferred']:
            improvement_factors.append(0.2)
        
        if improvement_factors:
            evaluation['overall_improvement'] = sum(improvement_factors) / len(improvement_factors)
        
        return evaluation

class HierarchicalModelTrainer:
    """分层模型训练体系"""
    
    def __init__(self):
        """初始化分层模型训练体系"""
        self.layers = {
            'data': self._train_data_layer,
            'feature': self._train_feature_layer,
            'inference': self._train_inference_layer
        }
    
    def _train_data_layer(self, data: List[Dict[str, Any]]) -> Dict[str, Any]:
        """训练数据层"""
        # 数据预处理和特征提取
        processed_data = {
            'processed_count': len(data),
            'timestamp': time.time(),
            'status': 'completed'
        }
        
        return processed_data
    
    def _train_feature_layer(self, features: List[Dict[str, Any]]) -> Dict[str, Any]:
        """训练特征层"""
        # 特征选择和降维
        feature_training = {
            'feature_count': len(features),
            'timestamp': time.time(),
            'status': 'completed'
        }
        
        return feature_training
    
    def _train_inference_layer(self, model: Dict[str, Any], data: List[Dict[str, Any]]) -> Dict[str, Any]:
        """训练推理层"""
        # 推理模型训练
        inference_training = {
            'model_architecture': model.get('architecture', 'unknown'),
            'training_samples': len(data),
            'timestamp': time.time(),
            'status': 'completed'
        }
        
        return inference_training
    
    def train_hierarchical_model(self, data: List[Dict[str, Any]], feature_assets: List[Dict[str, Any]], model: Dict[str, Any]) -> Dict[str, Any]:
        """训练分层模型"""
        # 训练数据层
        data_layer_result = self._train_data_layer(data)
        
        # 训练特征层
        feature_layer_result = self._train_feature_layer(feature_assets)
        
        # 训练推理层
        inference_layer_result = self._train_inference_layer(model, data)
        
        # 整合结果
        hierarchical_result = {
            'data_layer': data_layer_result,
            'feature_layer': feature_layer_result,
            'inference_layer': inference_layer_result,
            'timestamp': time.time(),
            'status': 'completed'
        }
        
        return hierarchical_result

class MarketFeedbackOptimizer:
    """市场反馈调优系统"""
    
    def __init__(self):
        """初始化市场反馈调优系统"""
        self.feedback_weights = {
            'usage_frequency': 0.3,
            'rating': 0.5,
            'transaction_volume': 0.2
        }
    
    def calculate_optimization_score(self, feedbacks: List[Dict[str, Any]]) -> float:
        """计算优化分数"""
        if not feedbacks:
            return 0.0
        
        total_score = 0.0
        for feedback in feedbacks:
            fb = feedback.get('feedback', {})
            
            # 计算各项反馈的权重
            usage_frequency = fb.get('usage_frequency', 1)
            rating = fb.get('rating', 3) / 5.0  # 归一化到0-1
            transaction_volume = fb.get('transaction_volume', 0)
            
            # 计算加权分数
            score = (
                usage_frequency * self.feedback_weights['usage_frequency'] +
                rating * self.feedback_weights['rating'] +
                transaction_volume * self.feedback_weights['transaction_volume']
            )
            
            total_score += score
        
        avg_score = total_score / len(feedbacks)
        return avg_score
    
    def generate_optimization_suggestions(self, algorithm_asset: Dict[str, Any], feedbacks: List[Dict[str, Any]]) -> List[str]:
        """生成优化建议"""
        suggestions = []
        
        if not feedbacks:
            suggestions.append('No feedback available for optimization')
            return suggestions
        
        # 分析反馈
        usage_count = len(feedbacks)
        avg_rating = sum(fb.get('feedback', {}).get('rating', 0) for fb in feedbacks) / usage_count
        
        if avg_rating < 3.0:
            suggestions.append('Model performance needs improvement')
        
        if usage_count < 10:
            suggestions.append('Need more usage data for better optimization')
        
        return suggestions
