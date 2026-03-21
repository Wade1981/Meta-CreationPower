"""HTTP server module for ELR Service Load."""

import http.server
import json
import urllib.parse
from utils.security import SecurityManager


class RequestHandler(http.server.BaseHTTPRequestHandler):
    """HTTP request handler for archive service API."""

    def __init__(self, *args, archive_service=None, logger=None, security_manager=None, **kwargs):
        """
        Initialize request handler with archive service and security manager.
        
        Args:
            archive_service: Archive service instance
            logger: Logger instance
            security_manager: Security manager instance
        """
        self.archive_service = archive_service
        self.logger = logger
        self.security_manager = security_manager or SecurityManager()
        super().__init__(*args, **kwargs)

    def _send_response(self, status_code, data):
        """
        Send JSON response.
        
        Args:
            status_code (int): HTTP status code
            data (dict): Response data
        """
        self.send_response(status_code)
        self.send_header('Content-type', 'application/json')
        self.send_header('Access-Control-Allow-Origin', '*')
        self.end_headers()
        self.wfile.write(json.dumps(data, ensure_ascii=False).encode('utf-8'))

    def _parse_body(self):
        """
        Parse request body as JSON.
        
        Returns:
            dict: Parsed JSON data or empty dict
        """
        content_length = int(self.headers.get('Content-Length', 0))
        if content_length > 0:
            body = self.rfile.read(content_length)
            try:
                return json.loads(body.decode('utf-8'))
            except json.JSONDecodeError:
                return {}
        return {}

    def _authenticate(self):
        """
        Authenticate request using Authorization header.
        
        Returns:
            dict: User payload if authenticated, None otherwise
        """
        auth_header = self.headers.get('Authorization')
        if not auth_header:
            return None
        
        try:
            token_type, token = auth_header.split(' ', 1)
            if token_type.lower() != 'bearer':
                return None
            
            payload = self.security_manager.validate_token(token)
            return payload
        except Exception as e:
            return None

    def do_GET(self):
        """Handle GET requests."""
        parsed_path = urllib.parse.urlparse(self.path)
        path = parsed_path.path
        query = urllib.parse.parse_qs(parsed_path.query)

        # Health check endpoint
        if path == '/health':
            self._send_response(200, {'status': 'ok'})
            return

        # Get archive by ID
        if path.startswith('/archives/'):
            archive_id = path.split('/')[-1]
            archive = self.archive_service.get_archive(archive_id)
            if archive:
                self._send_response(200, archive)
            else:
                self._send_response(404, {'error': 'Archive not found'})
            return

        # List all archives
        if path == '/archives':
            archives = self.archive_service.list_archives()
            self._send_response(200, archives)
            return

        # Search archives
        if path == '/search' and 'q' in query:
            query_term = query['q'][0]
            results = self.archive_service.search_archives(query_term)
            self._send_response(200, results)
            return

        # Default 404
        self._send_response(404, {'error': 'Not found'})

    def do_POST(self):
        """Handle POST requests."""
        parsed_path = urllib.parse.urlparse(self.path)
        path = parsed_path.path

        # Login endpoint
        if path == '/login':
            login_data = self._parse_body()
            username = login_data.get('username')
            password = login_data.get('password')
            
            # Simple authentication for demo purposes
            # In production, this should check against a user database
            if username == 'admin' and password == 'password':
                token = self.security_manager.generate_token(username)
                self._send_response(200, {'token': token, 'user_id': username})
            else:
                self._send_response(401, {'error': 'Invalid credentials'})
            return

        # Add new archive (requires authentication)
        if path == '/archives':
            # Check authentication
            user_payload = self._authenticate()
            if not user_payload:
                self._send_response(401, {'error': 'Unauthorized'})
                return
            
            archive_data = self._parse_body()
            if archive_data:
                added_archive = self.archive_service.add_archive(archive_data)
                self._send_response(201, added_archive)
            else:
                self._send_response(400, {'error': 'Invalid archive data'})
            return

        # Default 404
        self._send_response(404, {'error': 'Not found'})

    def do_PUT(self):
        """Handle PUT requests."""
        parsed_path = urllib.parse.urlparse(self.path)
        path = parsed_path.path

        # Update archive (requires authentication)
        if path.startswith('/archives/'):
            # Check authentication
            user_payload = self._authenticate()
            if not user_payload:
                self._send_response(401, {'error': 'Unauthorized'})
                return
            
            archive_id = path.split('/')[-1]
            archive_data = self._parse_body()
            if archive_data:
                updated_archive = self.archive_service.update_archive(archive_id, archive_data)
                if updated_archive:
                    self._send_response(200, updated_archive)
                else:
                    self._send_response(404, {'error': 'Archive not found'})
            else:
                self._send_response(400, {'error': 'Invalid archive data'})
            return

        # Default 404
        self._send_response(404, {'error': 'Not found'})

    def do_DELETE(self):
        """Handle DELETE requests."""
        parsed_path = urllib.parse.urlparse(self.path)
        path = parsed_path.path

        # Delete archive (requires authentication)
        if path.startswith('/archives/'):
            # Check authentication
            user_payload = self._authenticate()
            if not user_payload:
                self._send_response(401, {'error': 'Unauthorized'})
                return
            
            archive_id = path.split('/')[-1]
            deleted = self.archive_service.delete_archive(archive_id)
            if deleted:
                self._send_response(200, {'status': 'deleted'})
            else:
                self._send_response(404, {'error': 'Archive not found'})
            return

        # Default 404
        self._send_response(404, {'error': 'Not found'})


class HTTPServer:
    """HTTP server for archive service API."""

    def __init__(self, host, port, archive_service, logger, security_manager=None, cert_file=None, key_file=None):
        """
        Initialize HTTP server.
        
        Args:
            host (str): Server host
            port (int): Server port
            archive_service: Archive service instance
            logger: Logger instance
            security_manager: Security manager instance
            cert_file (str): Path to certificate file for HTTPS
            key_file (str): Path to private key file for HTTPS
        """
        self.host = host
        self.port = port
        self.archive_service = archive_service
        self.logger = logger
        self.security_manager = security_manager or SecurityManager()
        self.cert_file = cert_file
        self.key_file = key_file

    def start(self):
        """
        Start HTTP server.
        """
        # Create server with custom request handler
        def handler(*args, **kwargs):
            RequestHandler(*args, archive_service=self.archive_service, logger=self.logger, 
                          security_manager=self.security_manager, **kwargs)

        server = http.server.HTTPServer((self.host, self.port), handler)
        
        # Add SSL context if certificate files provided
        if self.cert_file and self.key_file:
            ssl_context = self.security_manager.get_ssl_context(self.cert_file, self.key_file)
            if ssl_context:
                server.socket = ssl_context.wrap_socket(server.socket, server_side=True)
                self.logger.info(f"HTTPS server started on {self.host}:{self.port}")
            else:
                self.logger.warning("Invalid certificate files, starting HTTP server instead")
                self.logger.info(f"HTTP server started on {self.host}:{self.port}")
        else:
            self.logger.info(f"HTTP server started on {self.host}:{self.port}")
        
        try:
            server.serve_forever()
        except KeyboardInterrupt:
            self.logger.info("HTTP server stopped")
            server.shutdown()
