#!/usr/bin/env python3
"""
ELR API服务器
提供ELR容器的HTTP API接口，支持文件上传功能
"""

import http.server
import socketserver
import json
import threading
import time
import os
import shutil
import urllib.parse
import subprocess
import platform

class ELRAPIHandler(http.server.BaseHTTPRequestHandler):
    """ELR API处理器"""
    
    def do_GET(self):
        """处理GET请求"""
        if self.path == '/api/status':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "status": "running",
                "message": "ELR API服务运行正常"
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
        elif self.path == '/api/containers':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = [
                {
                    "name": "test-container",
                    "status": "created"
                },
                {
                    "name": "python-app",
                    "status": "running"
                }
            ]
            self.wfile.write(json.dumps(response).encode('utf-8'))
        elif self.path == '/api/files':
            self.list_files()
        elif self.path == '/api/resources':
            self.get_resources()
        elif self.path == '/health':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "status": "ok",
                "timestamp": int(time.time()),
                "service": "elr-api-server"
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
        else:
            self.send_response(404)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "error": "Not found"
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
    
    def do_POST(self):
        """处理POST请求"""
        if self.path == '/api/upload':
            self.handle_file_upload()
        else:
            self.send_response(404)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "error": "Not found"
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
    
    def do_DELETE(self):
        """处理DELETE请求"""
        if self.path.startswith('/api/files/'):
            filename = self.path.split('/api/files/')[1]
            self.delete_file(filename)
        else:
            self.send_response(404)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "error": "Not found"
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
    
    def handle_file_upload(self):
        """处理文件上传"""
        try:
            # 解析Content-Type头，获取boundary
            content_type = self.headers.get('Content-Type', '')
            if not content_type.startswith('multipart/form-data'):
                self.send_response(400)
                self.send_header('Content-type', 'application/json')
                self.end_headers()
                response = {
                    "success": False,
                    "error": "Invalid Content-Type"
                }
                self.wfile.write(json.dumps(response).encode('utf-8'))
                return
            
            # 提取boundary
            boundary = content_type.split('boundary=')[1].strip()
            
            # 读取请求体
            content_length = int(self.headers.get('Content-Length', 0))
            body = self.rfile.read(content_length)
            
            # 解析multipart/form-data
            parts = body.split(b'--' + boundary.encode())
            
            # 寻找文件部分
            file_data = None
            filename = None
            
            for part in parts:
                if b'Content-Disposition: form-data' in part:
                    # 提取文件名
                    lines = part.split(b'\r\n')
                    for line in lines:
                        if b'filename="' in line:
                            filename = line.split(b'filename="')[1].split(b'"')[0].decode('utf-8')
                            break
                    
                    # 提取文件数据
                    if filename:
                        # 找到文件数据的开始位置
                        data_start = part.find(b'\r\n\r\n') + 4
                        if data_start > 3:
                            file_data = part[data_start:].rstrip(b'\r\n')
                            break
            
            if not filename or not file_data:
                self.send_response(400)
                self.send_header('Content-type', 'application/json')
                self.end_headers()
                response = {
                    "success": False,
                    "error": "No file uploaded"
                }
                self.wfile.write(json.dumps(response).encode('utf-8'))
                return
            
            # 确保上传目录存在
            upload_dir = 'uploads'
            if not os.path.exists(upload_dir):
                os.makedirs(upload_dir)
            
            # 保存文件
            filepath = os.path.join(upload_dir, filename)
            
            with open(filepath, 'wb') as f:
                f.write(file_data)
            
            # 处理文件（根据文件类型进行不同处理）
            file_type = self.detect_file_type(filename)
            process_result = self.process_file(filepath, file_type)
            
            # 返回成功响应
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "success": True,
                "message": f"File uploaded successfully: {filename}",
                "file_type": file_type,
                "filepath": filepath,
                "process_result": process_result
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
            
        except Exception as e:
            self.send_response(500)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "success": False,
                "error": str(e)
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
    
    def detect_file_type(self, filename):
        """检测文件类型"""
        ext = os.path.splitext(filename)[1].lower()
        
        # 模型文件
        if ext in ['.pt', '.pth', '.onnx', '.model']:
            return 'model'
        # Python文件
        elif ext in ['.py', '.pyc']:
            return 'python'
        # 配置文件
        elif ext in ['.json', '.yaml', '.yml', '.config']:
            return 'config'
        # 图像文件
        elif ext in ['.png', '.jpg', '.jpeg', '.gif', '.bmp']:
            return 'image'
        # 音频文件
        elif ext in ['.wav', '.mp3', '.flac', '.ogg']:
            return 'audio'
        # 视频文件
        elif ext in ['.mp4', '.avi', '.mov', '.mkv']:
            return 'video'
        # 文档文件
        elif ext in ['.txt', '.doc', '.docx', '.pdf']:
            return 'document'
        # 其他文件
        else:
            return 'other'
    
    def list_files(self):
        """列出上传的文件"""
        try:
            upload_dir = 'uploads'
            if not os.path.exists(upload_dir):
                files = []
            else:
                files = []
                for filename in os.listdir(upload_dir):
                    filepath = os.path.join(upload_dir, filename)
                    if os.path.isfile(filepath):
                        file_info = {
                            "name": filename,
                            "type": self.detect_file_type(filename),
                            "size": os.path.getsize(filepath),
                            "path": filepath
                        }
                        files.append(file_info)
            
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "success": True,
                "files": files
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
        except Exception as e:
            self.send_response(500)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "success": False,
                "error": str(e)
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
    
    def delete_file(self, filename):
        """删除文件"""
        try:
            upload_dir = 'uploads'
            filepath = os.path.join(upload_dir, filename)
            
            if not os.path.exists(filepath):
                self.send_response(404)
                self.send_header('Content-type', 'application/json')
                self.end_headers()
                response = {
                    "success": False,
                    "error": "File not found"
                }
                self.wfile.write(json.dumps(response).encode('utf-8'))
                return
            
            os.remove(filepath)
            
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "success": True,
                "message": f"File deleted successfully: {filename}"
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
        except Exception as e:
            self.send_response(500)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "success": False,
                "error": str(e)
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
    
    def process_file(self, filepath, file_type):
        """根据文件类型处理文件"""
        try:
            if file_type == 'model':
                # 处理模型文件
                return self.process_model_file(filepath)
            elif file_type == 'python':
                # 处理Python文件
                return self.process_python_file(filepath)
            elif file_type == 'config':
                # 处理配置文件
                return self.process_config_file(filepath)
            elif file_type == 'image':
                # 处理图像文件
                return self.process_image_file(filepath)
            elif file_type == 'audio':
                # 处理音频文件
                return self.process_audio_file(filepath)
            elif file_type == 'video':
                # 处理视频文件
                return self.process_video_file(filepath)
            else:
                # 其他文件类型
                return f"File processed as {file_type} type"
        except Exception as e:
            return f"Error processing file: {str(e)}"
    
    def process_model_file(self, filepath):
        """处理模型文件"""
        # 确保模型目录存在
        model_dir = 'models'
        if not os.path.exists(model_dir):
            os.makedirs(model_dir)
        
        # 复制模型文件到模型目录
        model_name = os.path.basename(filepath)
        model_dest = os.path.join(model_dir, model_name)
        shutil.copy2(filepath, model_dest)
        
        # 自动装载模型
        load_result = self.load_model(model_dest, model_name)
        
        return f"Model file loaded to: {model_dest}. {load_result}"
    
    def load_model(self, model_path, model_name):
        """自动装载模型"""
        try:
            # 这里可以添加实际的模型加载逻辑
            # 例如使用模型框架加载模型到内存
            # 记录模型状态
            if not hasattr(self, 'loaded_models'):
                self.loaded_models = {}
            
            # 模拟模型加载
            self.loaded_models[model_name] = {
                'path': model_path,
                'loaded_at': time.time(),
                'status': 'loaded'
            }
            
            return f"Model {model_name} automatically loaded"
        except Exception as e:
            return f"Error loading model: {str(e)}"
    
    def process_python_file(self, filepath):
        """处理Python文件"""
        # 确保组件目录存在
        component_dir = 'components'
        if not os.path.exists(component_dir):
            os.makedirs(component_dir)
        
        # 复制Python文件到组件目录
        component_name = os.path.basename(filepath)
        component_dest = os.path.join(component_dir, component_name)
        shutil.copy2(filepath, component_dest)
        
        # 自动装载组件
        load_result = self.load_component(component_dest, component_name)
        
        return f"Python file processed as component: {component_dest}. {load_result}"
    
    def load_component(self, component_path, component_name):
        """自动装载组件"""
        try:
            # 记录组件状态
            if not hasattr(self, 'loaded_components'):
                self.loaded_components = {}
            
            # 模拟组件加载
            self.loaded_components[component_name] = {
                'path': component_path,
                'loaded_at': time.time(),
                'status': 'loaded'
            }
            
            return f"Component {component_name} automatically loaded"
        except Exception as e:
            return f"Error loading component: {str(e)}"
    
    def process_config_file(self, filepath):
        """处理配置文件"""
        # 这里可以添加配置文件加载逻辑
        # 例如读取配置，应用到系统中
        return f"Config file processed: {filepath}"
    
    def process_image_file(self, filepath):
        """处理图像文件"""
        # 确保数字资产目录存在
        asset_dir = 'assets/images'
        if not os.path.exists(asset_dir):
            os.makedirs(asset_dir)
        
        # 复制图像文件到资产目录
        asset_name = os.path.basename(filepath)
        asset_dest = os.path.join(asset_dir, asset_name)
        shutil.copy2(filepath, asset_dest)
        
        # 自动装载数字资产
        load_result = self.load_asset(asset_dest, asset_name, 'image')
        
        return f"Image file processed as digital asset: {asset_dest}. {load_result}"
    
    def load_asset(self, asset_path, asset_name, asset_type):
        """自动装载数字资产"""
        try:
            # 记录资产状态
            if not hasattr(self, 'loaded_assets'):
                self.loaded_assets = {}
            
            # 模拟资产加载
            self.loaded_assets[asset_name] = {
                'path': asset_path,
                'type': asset_type,
                'loaded_at': time.time(),
                'status': 'loaded'
            }
            
            return f"Digital asset {asset_name} automatically loaded"
        except Exception as e:
            return f"Error loading asset: {str(e)}"
    
    def process_audio_file(self, filepath):
        """处理音频文件"""
        # 确保数字资产目录存在
        asset_dir = 'assets/audio'
        if not os.path.exists(asset_dir):
            os.makedirs(asset_dir)
        
        # 复制音频文件到资产目录
        asset_name = os.path.basename(filepath)
        asset_dest = os.path.join(asset_dir, asset_name)
        shutil.copy2(filepath, asset_dest)
        
        # 自动装载数字资产
        load_result = self.load_asset(asset_dest, asset_name, 'audio')
        
        return f"Audio file processed as digital asset: {asset_dest}. {load_result}"
    
    def process_video_file(self, filepath):
        """处理视频文件"""
        # 确保数字资产目录存在
        asset_dir = 'assets/video'
        if not os.path.exists(asset_dir):
            os.makedirs(asset_dir)
        
        # 复制视频文件到资产目录
        asset_name = os.path.basename(filepath)
        asset_dest = os.path.join(asset_dir, asset_name)
        shutil.copy2(filepath, asset_dest)
        
        # 自动装载数字资产
        load_result = self.load_asset(asset_dest, asset_name, 'video')
        
        return f"Video file processed as digital asset: {asset_dest}. {load_result}"
    
    def get_resources(self):
        """获取系统资源使用情况"""
        try:
            resources = {
                "memory": self.get_memory_usage(),
                "disk": self.get_disk_usage(),
                "cpu": self.get_cpu_usage(),
                "gpu": self.get_gpu_usage(),
                "system": {
                    "platform": platform.platform(),
                    "python_version": platform.python_version(),
                    "timestamp": int(time.time())
                }
            }
            
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "success": True,
                "resources": resources
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
        except Exception as e:
            self.send_response(500)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "success": False,
                "error": str(e)
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
    
    def get_memory_usage(self):
        """获取内存使用情况"""
        try:
            if platform.system() == 'Windows':
                # Windows系统
                output = subprocess.check_output(['wmic', 'OS', 'get', 'FreePhysicalMemory,TotalVisibleMemorySize', '/value']).decode('utf-8')
                lines = output.strip().split('\n')
                memory_info = {}
                for line in lines:
                    if '=' in line:
                        key, value = line.split('=', 1)
                        memory_info[key.strip()] = int(value.strip())
                
                total = memory_info.get('TotalVisibleMemorySize', 0) * 1024
                free = memory_info.get('FreePhysicalMemory', 0) * 1024
                used = total - free
                usage_percent = (used / total * 100) if total > 0 else 0
                
                return {
                    "total": total,
                    "used": used,
                    "free": free,
                    "usage_percent": round(usage_percent, 2)
                }
            else:
                # Linux/macOS系统
                import psutil
                memory = psutil.virtual_memory()
                return {
                    "total": memory.total,
                    "used": memory.used,
                    "free": memory.available,
                    "usage_percent": round(memory.percent, 2)
                }
        except Exception as e:
            return {
                "error": str(e),
                "total": 0,
                "used": 0,
                "free": 0,
                "usage_percent": 0
            }
    
    def get_disk_usage(self):
        """获取磁盘使用情况"""
        try:
            if platform.system() == 'Windows':
                # Windows系统
                output = subprocess.check_output(['wmic', 'logicaldisk', 'get', 'DeviceID,Size,FreeSpace', '/value']).decode('utf-8')
                lines = output.strip().split('\n')
                disks = []
                disk_info = {}
                for line in lines:
                    if '=' in line:
                        key, value = line.split('=', 1)
                        key = key.strip()
                        value = value.strip()
                        if key == 'DeviceID':
                            if disk_info:
                                disks.append(disk_info)
                            disk_info = {'device': value}
                        elif key == 'Size' and value:
                            disk_info['total'] = int(value)
                        elif key == 'FreeSpace' and value:
                            disk_info['free'] = int(value)
                if disk_info:
                    disks.append(disk_info)
                
                # 计算每个磁盘的使用情况
                for disk in disks:
                    if 'total' in disk and 'free' in disk:
                        disk['used'] = disk['total'] - disk['free']
                        disk['usage_percent'] = round((disk['used'] / disk['total'] * 100), 2) if disk['total'] > 0 else 0
                
                return disks
            else:
                # Linux/macOS系统
                import psutil
                disks = []
                for partition in psutil.disk_partitions():
                    try:
                        usage = psutil.disk_usage(partition.mountpoint)
                        disks.append({
                            "device": partition.device,
                            "mountpoint": partition.mountpoint,
                            "total": usage.total,
                            "used": usage.used,
                            "free": usage.free,
                            "usage_percent": round(usage.percent, 2)
                        })
                    except:
                        pass
                return disks
        except Exception as e:
            return [{
                "error": str(e),
                "total": 0,
                "used": 0,
                "free": 0,
                "usage_percent": 0
            }]
    
    def get_cpu_usage(self):
        """获取CPU使用情况"""
        try:
            if platform.system() == 'Windows':
                # Windows系统
                output = subprocess.check_output(['wmic', 'cpu', 'get', 'LoadPercentage', '/value']).decode('utf-8')
                lines = output.strip().split('\n')
                for line in lines:
                    if 'LoadPercentage=' in line:
                        usage = int(line.split('=', 1)[1].strip())
                        return {
                            "usage_percent": usage
                        }
                return {
                    "usage_percent": 0
                }
            else:
                # Linux/macOS系统
                import psutil
                return {
                    "usage_percent": round(psutil.cpu_percent(interval=1), 2)
                }
        except Exception as e:
            return {
                "error": str(e),
                "usage_percent": 0
            }
    
    def get_gpu_usage(self):
        """获取GPU使用情况"""
        try:
            if platform.system() == 'Windows':
                # Windows系统 - 尝试使用nvidia-smi
                try:
                    output = subprocess.check_output(['nvidia-smi', '--query-gpu=utilization.gpu,memory.used,memory.total', '--format=csv,noheader,nounits']).decode('utf-8')
                    gpus = []
                    for line in output.strip().split('\n'):
                        if line:
                            gpu_util, mem_used, mem_total = line.split(',')
                            gpus.append({
                                "gpu_utilization": float(gpu_util.strip()),
                                "memory_used": int(mem_used.strip()),
                                "memory_total": int(mem_total.strip()),
                                "memory_usage_percent": round((int(mem_used.strip()) / int(mem_total.strip()) * 100), 2) if int(mem_total.strip()) > 0 else 0
                            })
                    return gpus
                except:
                    return [{
                        "error": "GPU not found or nvidia-smi not available",
                        "gpu_utilization": 0,
                        "memory_used": 0,
                        "memory_total": 0,
                        "memory_usage_percent": 0
                    }]
            else:
                # Linux系统 - 尝试使用nvidia-smi
                try:
                    output = subprocess.check_output(['nvidia-smi', '--query-gpu=utilization.gpu,memory.used,memory.total', '--format=csv,noheader,nounits']).decode('utf-8')
                    gpus = []
                    for line in output.strip().split('\n'):
                        if line:
                            gpu_util, mem_used, mem_total = line.split(',')
                            gpus.append({
                                "gpu_utilization": float(gpu_util.strip()),
                                "memory_used": int(mem_used.strip()),
                                "memory_total": int(mem_total.strip()),
                                "memory_usage_percent": round((int(mem_used.strip()) / int(mem_total.strip()) * 100), 2) if int(mem_total.strip()) > 0 else 0
                            })
                    return gpus
                except:
                    return [{
                        "error": "GPU not found or nvidia-smi not available",
                        "gpu_utilization": 0,
                        "memory_used": 0,
                        "memory_total": 0,
                        "memory_usage_percent": 0
                    }]
        except Exception as e:
            return [{
                "error": str(e),
                "gpu_utilization": 0,
                "memory_used": 0,
                "memory_total": 0,
                "memory_usage_percent": 0
            }]

def start_server(port=8080):
    """启动API服务器"""
    handler = ELRAPIHandler
    with socketserver.TCPServer(("", port), handler) as httpd:
        print(f"ELR API服务器启动在 http://localhost:{port}")
        print("可用端点:")
        print("  GET  /api/status      - 获取ELR状态")
        print("  GET  /api/containers  - 获取容器列表")
        print("  POST /api/upload      - 上传文件")
        print("  GET  /api/files       - 列出上传的文件")
        print("  DELETE /api/files/{name} - 删除文件")
        print("  GET  /health          - 健康检查")
        httpd.serve_forever()

if __name__ == "__main__":
    start_server()
