#!/usr/bin/env python3
"""
ELR Service Load for EL-CSCC Archive

This service load provides a RESTful API for the EL-CSCC Archive digital archive management system.
It runs as a long-running service that can be loaded into ELR containers.
"""

import sys
import os

# Add src directory to Python path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), 'src'))

from config.config import load_config
from service.archive_service import ArchiveService
from api.http_server import HTTPServer
from utils.logger import setup_logger
from utils.security import SecurityManager

# Try to import WebSocket server, but make it optional
try:
    import asyncio
    from api.websocket_server import WebSocketServer
    WEBSOCKET_AVAILABLE = True
except ImportError:
    WEBSOCKET_AVAILABLE = False


def load_full_config():
    """
    Load full configuration including WebSocket and security settings.
    
    Returns:
        dict: Full configuration dictionary
    """
    config = load_config()
    
    # Add WebSocket settings
    config['websocket_host'] = os.environ.get('ELR_WEBSOCKET_HOST', config['host'])
    config['websocket_port'] = int(os.environ.get('ELR_WEBSOCKET_PORT', '8001'))
    
    # Add security settings
    config['secret_key'] = os.environ.get('ELR_SECRET_KEY', 'default_secret_key')
    config['cert_file'] = os.environ.get('ELR_CERT_FILE')
    config['key_file'] = os.environ.get('ELR_KEY_FILE')
    
    return config


# Define start_websocket_server only if WebSocket is available
if WEBSOCKET_AVAILABLE:
    async def start_websocket_server(websocket_server):
        """
        Start WebSocket server in asyncio event loop.
        
        Args:
            websocket_server: WebSocket server instance
        """
        await websocket_server.start()


def main():
    """Main function to start the service."""
    try:
        # Setup logger
        logger = setup_logger()
        logger.info("Starting ELR Service Load for EL-CSCC Archive")

        # Load configuration
        config = load_full_config()
        logger.info(f"Configuration loaded: {config}")

        # Initialize security manager
        security_manager = SecurityManager(config['secret_key'])
        logger.info("Security manager initialized")

        # Initialize archive service
        archive_service = ArchiveService(config['archive_file'])
        logger.info("Archive service initialized")

        # Start WebSocket server if available
        if WEBSOCKET_AVAILABLE:
            try:
                import asyncio
                # Create WebSocket server
                websocket_server = WebSocketServer(
                    config['websocket_host'], 
                    config['websocket_port'], 
                    archive_service, 
                    logger
                )

                # Start WebSocket server in background
                loop = asyncio.get_event_loop()
                loop.create_task(start_websocket_server(websocket_server))
                logger.info(f"WebSocket server starting on {config['websocket_host']}:{config['websocket_port']}")
            except Exception as e:
                logger.warning(f"Failed to start WebSocket server: {e}")
        else:
            logger.warning("WebSocket server not available (websockets library not installed)")

        # Create and start HTTP server
        server = HTTPServer(
            config['host'], 
            config['port'], 
            archive_service, 
            logger, 
            security_manager, 
            config['cert_file'], 
            config['key_file']
        )
        logger.info(f"Starting HTTP server on {config['host']}:{config['port']}")
        server.start()

    except Exception as e:
        print(f"Error starting service: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
