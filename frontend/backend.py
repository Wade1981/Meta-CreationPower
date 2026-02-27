from flask import Flask, jsonify, request
import json
import hashlib
import time
import requests

app = Flask(__name__)

# 模拟数据
assets = {
    'asset-001': {
        'asset_id': 'asset-001',
        'asset_type': 'general',
        'status': 'active',
        'created_at': '2026-02-20T10:00:00Z',
        'last_modified': '2026-02-20T10:00:00Z'
    },
    'asset-002': {
        'asset_id': 'asset-002',
        'asset_type': 'algorithm',
        'status': 'active',
        'created_at': '2026-02-21T11:30:00Z',
        'last_modified': '2026-02-21T11:30:00Z'
    },
    'asset-003': {
        'asset_id': 'asset-003',
        'asset_type': 'feature',
        'status': 'inactive',
        'created_at': '2026-02-22T14:20:00Z',
        'last_modified': '2026-02-22T14:20:00Z'
    }
}

workflows = [
    {
        'execution_id': 'exec-001',
        'workflow_id': 'workflow-001',
        'status': 'completed',
        'started_at': '2026-02-23T09:00:00Z',
        'completed_at': '2026-02-23T09:05:00Z'
    },
    {
        'execution_id': 'exec-002',
        'workflow_id': 'workflow-002',
        'status': 'running',
        'started_at': '2026-02-24T10:30:00Z',
        'completed_at': None
    }
]

# 仪表盘数据
dashboard_data = {
    'total_assets': len(assets),
    'active_processes': 3,
    'total_transactions': 45,
    'system_health': 'healthy'
}

# 系统状态数据
system_status = {
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
    'containers': [
        {
            'container_id': 'container-001',
            'status': 'running',
            'asset_id': 'asset-001'
        },
        {
            'container_id': 'container-002',
            'status': 'running',
            'asset_id': 'asset-002'
        }
    ]
}

# 允许跨域访问
@app.after_request
def after_request(response):
    response.headers.add('Access-Control-Allow-Origin', '*')
    response.headers.add('Access-Control-Allow-Headers', 'Content-Type,Authorization')
    response.headers.add('Access-Control-Allow-Methods', 'GET,PUT,POST,DELETE,OPTIONS')
    return response

# 获取仪表盘数据
@app.route('/api/dashboard', methods=['GET'])
def get_dashboard():
    return jsonify({'success': True, 'data': dashboard_data})

# 获取资产列表
@app.route('/api/assets', methods=['GET'])
def get_assets():
    asset_type = request.args.get('type')
    status = request.args.get('status')
    
    filtered_assets = list(assets.values())
    
    if asset_type:
        filtered_assets = [asset for asset in filtered_assets if asset['asset_type'] == asset_type]
    
    if status:
        filtered_assets = [asset for asset in filtered_assets if asset['status'] == status]
    
    return jsonify({'success': True, 'data': filtered_assets, 'total': len(filtered_assets)})

# 获取单个资产
@app.route('/api/assets/<asset_id>', methods=['GET'])
def get_asset(asset_id):
    asset = assets.get(asset_id)
    if not asset:
        return jsonify({'success': False, 'message': 'Asset not found'})
    return jsonify({'success': True, 'data': asset})

# 创建资产
@app.route('/api/assets', methods=['POST'])
def create_asset():
    data = request.json
    asset_type = data.get('type', 'general')
    asset_content = data.get('content', '')
    
    # 生成资产ID
    asset_id = 'asset-' + hashlib.sha256((str(data) + str(time.time())).encode()).hexdigest()[:8]
    
    new_asset = {
        'asset_id': asset_id,
        'asset_type': asset_type,
        'status': 'active',
        'created_at': time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime()),
        'last_modified': time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())
    }
    
    assets[asset_id] = new_asset
    dashboard_data['total_assets'] = len(assets)
    
    return jsonify({'success': True, 'data': new_asset, 'asset_id': asset_id})

# 更新资产
@app.route('/api/assets/<asset_id>', methods=['PUT'])
def update_asset(asset_id):
    asset = assets.get(asset_id)
    if not asset:
        return jsonify({'success': False, 'message': 'Asset not found'})
    
    data = request.json
    asset.update(data)
    asset['last_modified'] = time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())
    
    return jsonify({'success': True, 'message': 'Asset updated successfully', 'asset_id': asset_id})

# 删除资产
@app.route('/api/assets/<asset_id>', methods=['DELETE'])
def delete_asset(asset_id):
    asset = assets.get(asset_id)
    if not asset:
        return jsonify({'success': False, 'message': 'Asset not found'})
    
    del assets[asset_id]
    dashboard_data['total_assets'] = len(assets)
    
    return jsonify({'success': True, 'message': 'Asset deleted successfully', 'asset_id': asset_id})

# 获取系统状态
@app.route('/api/system/status', methods=['GET'])
def get_system_status():
    return jsonify({'success': True, 'data': {
        'dashboard': dashboard_data,
        'modules': system_status['modules'],
        'containers': system_status['containers']
    }})

# 执行工作流
@app.route('/api/workflows', methods=['POST'])
def execute_workflow():
    data = request.json
    workflow_id = data.get('workflow_id')
    
    if not workflow_id:
        return jsonify({'success': False, 'message': 'Workflow ID is required'})
    
    # 生成执行ID
    execution_id = 'exec-' + hashlib.sha256((workflow_id + str(time.time())).encode()).hexdigest()[:8]
    
    new_workflow = {
        'execution_id': execution_id,
        'workflow_id': workflow_id,
        'status': 'completed',
        'started_at': time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime()),
        'completed_at': time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())
    }
    
    workflows.append(new_workflow)
    
    return jsonify({'success': True, 'data': new_workflow})

# 获取工作流执行历史
@app.route('/api/workflows', methods=['GET'])
def get_workflows():
    return jsonify({'success': True, 'data': workflows})

# ELR API 基础URL
ELR_API_BASE_URL = 'http://localhost:8080'  # 默认ELR API端口

# ELR API 客户端
def elr_api_request(endpoint, method='GET', data=None):
    url = f"{ELR_API_BASE_URL}{endpoint}"
    try:
        if method == 'GET':
            response = requests.get(url)
        elif method == 'POST':
            response = requests.post(url, json=data)
        elif method == 'PUT':
            response = requests.put(url, json=data)
        elif method == 'DELETE':
            response = requests.delete(url)
        else:
            return {'success': False, 'message': 'Invalid method'}
        
        response.raise_for_status()
        return {'success': True, 'data': response.json()}
    except requests.exceptions.RequestException as e:
        return {'success': False, 'message': str(e)}

# ELR 模型相关端点
@app.route('/api/elr/models', methods=['GET'])
def get_elr_models():
    return jsonify(elr_api_request('/api/models'))

@app.route('/api/elr/models/<model_id>', methods=['GET'])
def get_elr_model(model_id):
    return jsonify(elr_api_request(f'/api/models/{model_id}'))

@app.route('/api/elr/models', methods=['POST'])
def download_elr_model():
    data = request.json
    return jsonify(elr_api_request('/api/models/download', 'POST', data))

@app.route('/api/elr/models/<model_id>', methods=['DELETE'])
def delete_elr_model(model_id):
    return jsonify(elr_api_request(f'/api/models/{model_id}', 'DELETE'))

@app.route('/api/elr/models/<model_id>', methods=['PUT'])
def update_elr_model(model_id):
    data = request.json
    return jsonify(elr_api_request(f'/api/models/{model_id}', 'PUT', data))

@app.route('/api/elr/models/run', methods=['POST'])
def run_elr_model():
    data = request.json
    return jsonify(elr_api_request('/api/models/run', 'POST', data))

# ELR 容器相关端点
@app.route('/api/elr/containers', methods=['GET'])
def get_elr_containers():
    return jsonify(elr_api_request('/api/containers'))

@app.route('/api/elr/containers/<container_name>', methods=['GET'])
def get_elr_container(container_name):
    return jsonify(elr_api_request(f'/api/containers/{container_name}'))

@app.route('/api/elr/containers', methods=['POST'])
def create_elr_container():
    data = request.json
    return jsonify(elr_api_request('/api/containers/create', 'POST', data))

@app.route('/api/elr/containers/start', methods=['POST'])
def start_elr_container():
    data = request.json
    return jsonify(elr_api_request('/api/containers/start', 'POST', data))

@app.route('/api/elr/containers/stop', methods=['POST'])
def stop_elr_container():
    data = request.json
    return jsonify(elr_api_request('/api/containers/stop', 'POST', data))

@app.route('/api/elr/containers/<container_name>', methods=['DELETE'])
def delete_elr_container(container_name):
    return jsonify(elr_api_request(f'/api/containers/{container_name}', 'DELETE'))

# ELR 沙箱相关端点
@app.route('/api/elr/sandbox/status/<container>', methods=['GET'])
def get_sandbox_status(container):
    return jsonify(elr_api_request(f'/api/sandbox/status/{container}'))

@app.route('/api/elr/sandbox/execute', methods=['POST'])
def execute_sandbox_command():
    data = request.json
    return jsonify(elr_api_request('/api/sandbox/execute', 'POST', data))

# ELR 监控相关端点
@app.route('/api/elr/monitor/metrics', methods=['GET'])
def get_elr_metrics():
    return jsonify(elr_api_request('/api/monitor/metrics'))

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)