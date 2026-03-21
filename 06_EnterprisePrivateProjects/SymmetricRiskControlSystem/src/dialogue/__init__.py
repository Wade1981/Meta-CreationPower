# 对称风险管控系统对话式交互模块
# 支持查询系统状态和运行机制

import os
import sys
import json
import time
import argparse
from datetime import datetime

class DialogueInterface:
    """对话式交互界面类"""
    
    def __init__(self):
        """初始化对话式交互界面"""
        self.config_file = os.path.join(os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__)))), "config", "system_config.json")
        self.log_file = os.path.join(os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__)))), "logs", "dialogue.log")
        self.config = self._load_config()
        self._ensure_log_directory()
    
    def _ensure_log_directory(self):
        """确保日志目录存在"""
        log_dir = os.path.dirname(self.log_file)
        if not os.path.exists(log_dir):
            os.makedirs(log_dir)
    
    def _load_config(self):
        """加载系统配置"""
        try:
            if os.path.exists(self.config_file):
                with open(self.config_file, 'r', encoding='utf-8') as f:
                    return json.load(f)
            else:
                return self._get_default_config()
        except Exception as e:
            self._log_error(f"加载配置文件失败: {str(e)}")
            return self._get_default_config()
    
    def _get_default_config(self):
        """获取默认配置"""
        return {
            "version": "1.0.0",
            "container_name": "symmetric-risk-control",
            "log_level": "info",
            "risk_assessment": {
                "enabled": True,
                "threshold": 0.7,
                "update_interval": 60
            },
            "risk_control": {
                "enabled": True,
                "max_exposure": 0.5,
                "hedging_ratio": 0.8
            },
            "dialogue": {
                "enabled": True,
                "language": "zh-CN"
            }
        }
    
    def _log(self, message, level="info"):
        """记录日志"""
        timestamp = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        log_entry = f"[{timestamp}] [{level.upper()}] {message}"
        print(log_entry)
        try:
            with open(self.log_file, 'a', encoding='utf-8') as f:
                f.write(log_entry + '\n')
        except Exception:
            pass
    
    def _log_error(self, message):
        """记录错误日志"""
        self._log(message, "error")
    
    def _log_info(self, message):
        """记录信息日志"""
        self._log(message, "info")
    
    def get_status(self):
        """获取系统状态"""
        try:
            status = {
                "timestamp": datetime.now().isoformat(),
                "system": {
                    "version": self.config.get("version", "1.0.0"),
                    "status": "运行中",
                    "uptime": self._get_uptime(),
                    "container_name": self.config.get("container_name", "symmetric-risk-control")
                },
                "risk_assessment": {
                    "enabled": self.config.get("risk_assessment", {}).get("enabled", True),
                    "threshold": self.config.get("risk_assessment", {}).get("threshold", 0.7),
                    "update_interval": self.config.get("risk_assessment", {}).get("update_interval", 60),
                    "risk_level": self._calculate_risk_level(),
                    "last_assessment": self._get_last_assessment_time()
                },
                "risk_control": {
                    "enabled": self.config.get("risk_control", {}).get("enabled", True),
                    "max_exposure": self.config.get("risk_control", {}).get("max_exposure", 0.5),
                    "hedging_ratio": self.config.get("risk_control", {}).get("hedging_ratio", 0.8),
                    "active_strategies": self._get_active_strategies()
                },
                "resources": {
                    "cpu_usage": self._get_cpu_usage(),
                    "memory_usage": self._get_memory_usage(),
                    "disk_usage": self._get_disk_usage()
                }
            }
            
            return self._format_status(status)
        except Exception as e:
            self._log_error(f"获取系统状态失败: {str(e)}")
            return "系统状态获取失败，请检查日志获取详细信息。"
    
    def get_mechanism(self):
        """获取运行机制"""
        try:
            mechanism = {
                "overview": "对称风险管控系统基于ELR容器运行，采用多层风险管控架构，实现全方位的风险识别、评估和管控。",
                "architecture": "系统采用碳硅协同设计理念，由风险识别、风险评估、风险管控和策略生成四个核心模块组成。",
                "process": [
                    "1. 数据采集: 从企业各系统实时采集运营数据",
                    "2. 风险识别: 基于多维度指标分析识别潜在风险",
                    "3. 风险评估: 使用机器学习模型评估风险等级",
                    "4. 策略生成: 根据风险评估结果生成对冲策略",
                    "5. 执行监控: 实时监控策略执行效果",
                    "6. 反馈优化: 基于执行结果持续优化风险模型"
                ],
                "features": [
                    "- 实时风险监控: 24/7不间断监控企业运营风险",
                    "- 智能风险评估: 基于机器学习的风险等级评估",
                    "- 自动策略生成: 根据风险等级自动生成对冲策略",
                    "- 容器隔离保护: 基于ELR容器的安全隔离",
                    "- 对话式交互: 支持查询系统状态和运行机制"
                ],
                "elr_integration": "系统通过ELR容器的exec命令执行，实现无源代码部署，确保核心算法安全。"
            }
            
            return self._format_mechanism(mechanism)
        except Exception as e:
            self._log_error(f"获取运行机制失败: {str(e)}")
            return "运行机制获取失败，请检查日志获取详细信息。"
    
    def get_report(self):
        """获取风险分析报告"""
        try:
            report = {
                "title": "对称风险管控系统 - 风险分析报告",
                "generated_at": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
                "summary": "本报告提供对称风险管控系统的风险分析结果和建议。",
                "risk_analysis": {
                    "current_risk_level": self._calculate_risk_level(),
                    "risk_trend": "稳定",
                    "key_risk_factors": [
                        "市场波动风险",
                        "信用风险",
                        "流动性风险",
                        "操作风险"
                    ],
                    "risk_distribution": {
                        "low": 40,
                        "medium": 35,
                        "high": 20,
                        "critical": 5
                    }
                },
                "recommendations": [
                    "1. 优化风险对冲策略，提高对冲效率",
                    "2. 调整风险限额，适应市场变化",
                    "3. 加强流动性风险管理，提高资金使用效率",
                    "4. 定期更新风险模型，适应新的市场环境"
                ],
                "conclusion": "系统运行正常，风险水平可控。建议继续保持当前的风险管控策略，并根据市场变化及时调整。"
            }
            
            return self._format_report(report)
        except Exception as e:
            self._log_error(f"获取风险分析报告失败: {str(e)}")
            return "风险分析报告获取失败，请检查日志获取详细信息。"
    
    def get_config(self, param=None):
        """获取系统配置"""
        try:
            if param:
                # 获取特定参数
                value = self._get_nested_value(self.config, param.split('.'))
                if value is not None:
                    return f"配置参数 '{param}': {value}"
                else:
                    return f"配置参数 '{param}' 不存在"
            else:
                # 获取所有配置
                return self._format_config(self.config)
        except Exception as e:
            self._log_error(f"获取系统配置失败: {str(e)}")
            return "系统配置获取失败，请检查日志获取详细信息。"
    
    def _get_nested_value(self, data, keys):
        """获取嵌套字典的值"""
        for key in keys:
            if isinstance(data, dict) and key in data:
                data = data[key]
            else:
                return None
        return data
    
    def _format_status(self, status):
        """格式化系统状态输出"""
        output = []
        output.append("========================================")
        output.append("对称风险管控系统 - 状态信息")
        output.append("========================================")
        output.append(f"生成时间: {status['timestamp']}")
        output.append("")
        
        # 系统信息
        output.append("[系统信息]")
        output.append(f"版本: {status['system']['version']}")
        output.append(f"状态: {status['system']['status']}")
        output.append(f"运行时间: {status['system']['uptime']}")
        output.append(f"容器名称: {status['system']['container_name']}")
        output.append("")
        
        # 风险评估
        output.append("[风险评估]")
        output.append(f"启用状态: {'已启用' if status['risk_assessment']['enabled'] else '已禁用'}")
        output.append(f"风险阈值: {status['risk_assessment']['threshold']}")
        output.append(f"更新间隔: {status['risk_assessment']['update_interval']}秒")
        output.append(f"当前风险等级: {status['risk_assessment']['risk_level']}")
        output.append(f"上次评估时间: {status['risk_assessment']['last_assessment']}")
        output.append("")
        
        # 风险管控
        output.append("[风险管控]")
        output.append(f"启用状态: {'已启用' if status['risk_control']['enabled'] else '已禁用'}")
        output.append(f"最大风险暴露: {status['risk_control']['max_exposure']}")
        output.append(f"对冲比例: {status['risk_control']['hedging_ratio']}")
        output.append(f"活跃策略数: {status['risk_control']['active_strategies']}")
        output.append("")
        
        # 资源使用
        output.append("[资源使用]")
        output.append(f"CPU使用率: {status['resources']['cpu_usage']}")
        output.append(f"内存使用率: {status['resources']['memory_usage']}")
        output.append(f"磁盘使用率: {status['resources']['disk_usage']}")
        output.append("")
        output.append("========================================")
        
        return '\n'.join(output)
    
    def _format_mechanism(self, mechanism):
        """格式化运行机制输出"""
        output = []
        output.append("========================================")
        output.append("对称风险管控系统 - 运行机制")
        output.append("========================================")
        output.append("")
        
        output.append("[系统概述]")
        output.append(mechanism['overview'])
        output.append("")
        
        output.append("[系统架构]")
        output.append(mechanism['architecture'])
        output.append("")
        
        output.append("[运行流程]")
        for step in mechanism['process']:
            output.append(step)
        output.append("")
        
        output.append("[核心功能]")
        for feature in mechanism['features']:
            output.append(feature)
        output.append("")
        
        output.append("[ELR集成]")
        output.append(mechanism['elr_integration'])
        output.append("")
        output.append("========================================")
        
        return '\n'.join(output)
    
    def _format_report(self, report):
        """格式化风险分析报告输出"""
        output = []
        output.append("========================================")
        output.append(report['title'])
        output.append("========================================")
        output.append(f"生成时间: {report['generated_at']}")
        output.append("")
        
        output.append("[报告摘要]")
        output.append(report['summary'])
        output.append("")
        
        output.append("[风险分析]")
        output.append(f"当前风险等级: {report['risk_analysis']['current_risk_level']}")
        output.append(f"风险趋势: {report['risk_analysis']['risk_trend']}")
        output.append("")
        output.append("关键风险因素:")
        for factor in report['risk_analysis']['key_risk_factors']:
            output.append(f"- {factor}")
        output.append("")
        output.append("风险分布:")
        output.append(f"低风险: {report['risk_analysis']['risk_distribution']['low']}%")
        output.append(f"中风险: {report['risk_analysis']['risk_distribution']['medium']}%")
        output.append(f"高风险: {report['risk_analysis']['risk_distribution']['high']}%")
        output.append(f"临界风险: {report['risk_analysis']['risk_distribution']['critical']}%")
        output.append("")
        
        output.append("[建议措施]")
        for recommendation in report['recommendations']:
            output.append(recommendation)
        output.append("")
        
        output.append("[结论]")
        output.append(report['conclusion'])
        output.append("")
        output.append("========================================")
        
        return '\n'.join(output)
    
    def _format_config(self, config):
        """格式化系统配置输出"""
        output = []
        output.append("========================================")
        output.append("对称风险管控系统 - 配置信息")
        output.append("========================================")
        output.append("")
        
        output.append("[基本配置]")
        output.append(f"版本: {config.get('version', '1.0.0')}")
        output.append(f"容器名称: {config.get('container_name', 'symmetric-risk-control')}")
        output.append(f"日志级别: {config.get('log_level', 'info')}")
        output.append("")
        
        output.append("[风险评估配置]")
        risk_assessment = config.get('risk_assessment', {})
        output.append(f"启用状态: {'已启用' if risk_assessment.get('enabled', True) else '已禁用'}")
        output.append(f"风险阈值: {risk_assessment.get('threshold', 0.7)}")
        output.append(f"更新间隔: {risk_assessment.get('update_interval', 60)}秒")
        output.append("")
        
        output.append("[风险管控配置]")
        risk_control = config.get('risk_control', {})
        output.append(f"启用状态: {'已启用' if risk_control.get('enabled', True) else '已禁用'}")
        output.append(f"最大风险暴露: {risk_control.get('max_exposure', 0.5)}")
        output.append(f"对冲比例: {risk_control.get('hedging_ratio', 0.8)}")
        output.append("")
        
        output.append("[对话式交互配置]")
        dialogue = config.get('dialogue', {})
        output.append(f"启用状态: {'已启用' if dialogue.get('enabled', True) else '已禁用'}")
        output.append(f"语言: {dialogue.get('language', 'zh-CN')}")
        output.append("")
        output.append("========================================")
        
        return '\n'.join(output)
    
    def _get_uptime(self):
        """获取系统运行时间"""
        # 模拟系统运行时间
        return "2天 14小时 30分钟"
    
    def _calculate_risk_level(self):
        """计算当前风险等级"""
        # 模拟风险等级计算
        import random
        risk_levels = ["低", "中", "高", "临界"]
        return random.choice(risk_levels)
    
    def _get_last_assessment_time(self):
        """获取上次风险评估时间"""
        # 模拟上次评估时间
        return datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    
    def _get_active_strategies(self):
        """获取活跃策略数量"""
        # 模拟活跃策略数
        import random
        return random.randint(5, 15)
    
    def _get_cpu_usage(self):
        """获取CPU使用率"""
        # 模拟CPU使用率
        import random
        return f"{random.randint(10, 60)}%"
    
    def _get_memory_usage(self):
        """获取内存使用率"""
        # 模拟内存使用率
        import random
        return f"{random.randint(30, 70)}%"
    
    def _get_disk_usage(self):
        """获取磁盘使用率"""
        # 模拟磁盘使用率
        import random
        return f"{random.randint(20, 50)}%"

def main():
    """主函数"""
    parser = argparse.ArgumentParser(description='对称风险管控系统对话式交互工具')
    parser.add_argument('command', choices=['status', 'mechanism', 'report', 'config'], help='要执行的命令')
    parser.add_argument('--param', help='配置参数路径 (仅用于 config 命令)')
    
    args = parser.parse_args()
    
    dialogue = DialogueInterface()
    
    try:
        if args.command == 'status':
            print(dialogue.get_status())
        elif args.command == 'mechanism':
            print(dialogue.get_mechanism())
        elif args.command == 'report':
            print(dialogue.get_report())
        elif args.command == 'config':
            print(dialogue.get_config(args.param))
    except Exception as e:
        print(f"执行命令失败: {str(e)}")
        sys.exit(1)

if __name__ == '__main__':
    main()
