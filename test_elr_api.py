#!/usr/bin/env python3
"""
测试ELR API服务
"""

import http.server
import socketserver
import json
import threading
import time

class ELRTestHandler(http.server.BaseHTTPRequestHandler):
    """测试ELR API处理器"""
    
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
        elif self.path == '/health':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "status": "ok",
                "timestamp": int(time.time()),
                "service": "micro-model-server"
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

def start_test_server():
    """启动测试服务器"""
    PORT = 8080
    Handler = ELRTestHandler
    
    with socketserver.TCPServer(("", PORT), Handler) as httpd:
        print(f"测试ELR API服务启动在 http://localhost:{PORT}")
        httpd.serve_forever()

if __name__ == "__main__":
    # 启动测试服务器
    server_thread = threading.Thread(target=start_test_server)
    server_thread.daemon = True
    server_thread.start()
    
    print("测试ELR API服务已启动")
    print("按Ctrl+C退出")
    
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("退出测试服务")
