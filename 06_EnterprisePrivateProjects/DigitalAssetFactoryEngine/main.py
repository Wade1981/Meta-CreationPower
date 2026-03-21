import json
import sys
import argparse
from src.control_center.control_center import DigitalAssetControlCenter

class DigitalAssetFactoryEngine:
    """数字资产工厂引擎主入口"""
    
    def __init__(self, config=None):
        """初始化数字资产工厂引擎"""
        self.config = config or {}
        self.control_center = DigitalAssetControlCenter(self.config)
        self.running = False
    
    def start(self):
        """启动数字资产工厂引擎"""
        print("Starting Digital Asset Factory Engine...")
        print(f"Control Center ID: {self.control_center.center_id}")
        
        # 启动各个模块
        print("Initializing modules...")
        
        # 检查系统状态
        status = self.control_center.get_system_status()
        print(f"System status: {status['dashboard']['system_health']}")
        print(f"Supported IDEs: {status['ide_status']}")
        
        self.running = True
        print("Digital Asset Factory Engine started successfully!")
        
        return status
    
    def stop(self):
        """停止数字资产工厂引擎"""
        print("Stopping Digital Asset Factory Engine...")
        self.running = False
        print("Digital Asset Factory Engine stopped successfully!")
        return {"status": "stopped"}
    
    def create_asset(self, asset_data):
        """创建数字资产"""
        if not self.running:
            return {"error": "Engine is not running"}
        
        asset = self.control_center.create_asset(asset_data)
        return asset
    
    def get_asset(self, asset_id):
        """获取数字资产"""
        if not self.running:
            return {"error": "Engine is not running"}
        
        asset = self.control_center.get_asset(asset_id)
        return asset
    
    def list_assets(self, asset_type=None, status=None):
        """列出数字资产"""
        if not self.running:
            return {"error": "Engine is not running"}
        
        assets = self.control_center.list_assets(asset_type, status)
        return assets
    
    def delete_asset(self, asset_id):
        """删除数字资产"""
        if not self.running:
            return {"error": "Engine is not running"}
        
        result = self.control_center.delete_asset(asset_id)
        return result
    
    def get_dashboard(self):
        """获取仪表盘"""
        if not self.running:
            return {"error": "Engine is not running"}
        
        dashboard = self.control_center.get_dashboard()
        return dashboard
    
    def get_system_status(self):
        """获取系统状态"""
        status = self.control_center.get_system_status()
        return status
    
    def execute_workflow(self, workflow_id, workflow_data):
        """执行工作流"""
        if not self.running:
            return {"error": "Engine is not running"}
        
        result = self.control_center.execute_workflow(workflow_id, workflow_data)
        return result

def main():
    """主函数"""
    parser = argparse.ArgumentParser(description="Digital Asset Factory Engine")
    parser.add_argument('--start', action='store_true', help='Start the engine')
    parser.add_argument('--stop', action='store_true', help='Stop the engine')
    parser.add_argument('--status', action='store_true', help='Get system status')
    parser.add_argument('--dashboard', action='store_true', help='Get dashboard')
    parser.add_argument('--create-asset', type=str, help='Create asset with JSON data')
    parser.add_argument('--get-asset', type=str, help='Get asset by ID')
    parser.add_argument('--list-assets', action='store_true', help='List assets')
    parser.add_argument('--delete-asset', type=str, help='Delete asset by ID')
    
    args = parser.parse_args()
    
    # 创建引擎实例
    engine = DigitalAssetFactoryEngine()
    
    # 处理命令行参数
    if args.start:
        status = engine.start()
        print(json.dumps(status, indent=2, ensure_ascii=False))
    elif args.stop:
        status = engine.stop()
        print(json.dumps(status, indent=2, ensure_ascii=False))
    elif args.status:
        status = engine.get_system_status()
        print(json.dumps(status, indent=2, ensure_ascii=False))
    elif args.dashboard:
        dashboard = engine.get_dashboard()
        print(json.dumps(dashboard, indent=2, ensure_ascii=False))
    elif args.create_asset:
        try:
            asset_data = json.loads(args.create_asset)
            engine.start()
            asset = engine.create_asset(asset_data)
            print(json.dumps(asset, indent=2, ensure_ascii=False))
        except json.JSONDecodeError as e:
            print(f"Invalid JSON: {e}")
    elif args.get_asset:
        engine.start()
        asset = engine.get_asset(args.get_asset)
        print(json.dumps(asset, indent=2, ensure_ascii=False))
    elif args.list_assets:
        engine.start()
        assets = engine.list_assets()
        print(json.dumps(assets, indent=2, ensure_ascii=False))
    elif args.delete_asset:
        engine.start()
        result = engine.delete_asset(args.delete_asset)
        print(json.dumps(result, indent=2, ensure_ascii=False))
    else:
        parser.print_help()

if __name__ == "__main__":
    main()
