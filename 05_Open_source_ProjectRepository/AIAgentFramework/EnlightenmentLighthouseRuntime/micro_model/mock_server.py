#!/usr/bin/env python3
"""
Mock ELR API server for testing
"""

import http.server
import json
import socketserver

PORT = 8082

class MockELRHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        # API endpoints
        if self.path == '/api':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "status": "ok",
                "message": "ELR API is running",
                "version": "1.0.0",
                "endpoints": [
                    "/api/models",
                    "/api/containers",
                    "/api/sandbox"
                ]
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
        elif self.path == '/api/models':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "status": "ok",
                "models": [
                    {
                        "name": "fish-speech",
                        "version": "1.0.0",
                        "status": "loaded",
                        "path": "model/models/fish-speech"
                    }
                ]
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
        elif self.path == '/api/containers':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "status": "ok",
                "containers": [
                    {
                        "id": "elr-1234567890",
                        "name": "test-container",
                        "status": "running"
                    }
                ]
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
        elif self.path == '/api/sandbox':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "status": "ok",
                "sandbox": {
                    "status": "running",
                    "models": ["fish-speech"]
                }
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
        else:
            self.send_response(404)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "status": "error",
                "message": "Endpoint not found"
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))

    def do_POST(self):
        # Handle POST requests
        if self.path == '/api/models/load':
            content_length = int(self.headers['Content-Length'])
            post_data = self.rfile.read(content_length)
            data = json.loads(post_data)
            
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "status": "ok",
                "message": f"Model {data.get('model')} loaded successfully"
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))
        else:
            self.send_response(404)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            response = {
                "status": "error",
                "message": "Endpoint not found"
            }
            self.wfile.write(json.dumps(response).encode('utf-8'))

def start_server():
    """Start the mock ELR API server"""
    print(f"Starting mock ELR API server on port {PORT}...")
    print(f"API endpoints available at: http://localhost:{PORT}/api")
    print("Press Ctrl+C to stop the server")
    
    with socketserver.TCPServer(("", PORT), MockELRHandler) as httpd:
        try:
            httpd.serve_forever()
        except KeyboardInterrupt:
            print("\nServer stopped.")

if __name__ == "__main__":
    start_server()
