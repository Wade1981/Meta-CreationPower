import http.server
import socketserver

PORT = 9004

class MyHTTPRequestHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/health':
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(b'{"status":"ok","service":"python-server"}')
        else:
            super().do_GET()

with socketserver.TCPServer(("127.0.0.1", PORT), MyHTTPRequestHandler) as httpd:
    print(f"Python server starting on 127.0.0.1:{PORT}")
    httpd.serve_forever()
