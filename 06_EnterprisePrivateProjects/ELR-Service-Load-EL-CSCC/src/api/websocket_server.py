"""WebSocket server module for ELR Service Load."""

import asyncio
import json
import logging
from websockets import serve


class WebSocketServer:
    """WebSocket server for real-time communication."""

    def __init__(self, host, port, archive_service, logger):
        """
        Initialize WebSocket server.
        
        Args:
            host (str): Server host
            port (int): Server port
            archive_service: Archive service instance
            logger: Logger instance
        """
        self.host = host
        self.port = port
        self.archive_service = archive_service
        self.logger = logger
        self.clients = set()

    async def handle_client(self, websocket):
        """
        Handle WebSocket client connection.
        
        Args:
            websocket: WebSocket connection
        """
        # Add client to set
        self.clients.add(websocket)
        self.logger.info(f"New WebSocket client connected. Total clients: {len(self.clients)}")

        try:
            async for message in websocket:
                # Process received message
                await self.process_message(websocket, message)
        except Exception as e:
            self.logger.error(f"WebSocket error: {e}")
        finally:
            # Remove client from set
            self.clients.remove(websocket)
            self.logger.info(f"WebSocket client disconnected. Total clients: {len(self.clients)}")

    async def process_message(self, websocket, message):
        """
        Process WebSocket message.
        
        Args:
            websocket: WebSocket connection
            message (str): Received message
        """
        try:
            # Parse JSON message
            data = json.loads(message)
            self.logger.debug(f"Received WebSocket message: {data}")

            # Handle different message types
            message_type = data.get('type')
            if message_type == 'ping':
                await websocket.send(json.dumps({'type': 'pong'}))
            elif message_type == 'get_archive':
                archive_id = data.get('archive_id')
                archive = self.archive_service.get_archive(archive_id)
                await websocket.send(json.dumps({'type': 'archive', 'data': archive}))
            elif message_type == 'list_archives':
                archives = self.archive_service.list_archives()
                await websocket.send(json.dumps({'type': 'archives', 'data': archives}))
            elif message_type == 'search':
                query = data.get('query')
                results = self.archive_service.search_archives(query)
                await websocket.send(json.dumps({'type': 'search_results', 'data': results}))
            else:
                await websocket.send(json.dumps({'type': 'error', 'message': 'Unknown message type'}))
        except json.JSONDecodeError:
            await websocket.send(json.dumps({'type': 'error', 'message': 'Invalid JSON format'}))
        except Exception as e:
            self.logger.error(f"Error processing WebSocket message: {e}")
            await websocket.send(json.dumps({'type': 'error', 'message': 'Internal server error'}))

    async def broadcast(self, message):
        """
        Broadcast message to all connected clients.
        
        Args:
            message (dict): Message to broadcast
        """
        if self.clients:
            try:
                message_json = json.dumps(message)
                await asyncio.gather(*[client.send(message_json) for client in self.clients])
            except Exception as e:
                self.logger.error(f"Error broadcasting message: {e}")

    async def start(self):
        """
        Start WebSocket server.
        """
        async with serve(self.handle_client, self.host, self.port):
            self.logger.info(f"WebSocket server started on {self.host}:{self.port}")
            await asyncio.Future()  # Run forever

    def stop(self):
        """
        Stop WebSocket server.
        """
        self.logger.info("WebSocket server stopping...")
        # Close all client connections
        for client in self.clients:
            try:
                asyncio.create_task(client.close())
            except Exception as e:
                self.logger.error(f"Error closing client connection: {e}")
