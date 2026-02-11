// AIAgent Network Server - 网络服务器

class NetworkServer {
  constructor(config) {
    this.config = config || {};
    this.servers = new Map();
    this.routes = new Map();
    this.middlewares = [];
    this.isRunning = false;
  }

  async initialize() {
    console.log('Initializing AIAgent Network Server...');
    // 初始化服务器组件
    this.httpServer = new HTTPServer(this.config.http || {});
    this.httpsServer = new HTTPSServer(this.config.https || {});
    this.tcpServer = new TCPServer(this.config.tcp || {});
    this.udpServer = new UDPServer(this.config.udp || {});
    this.websocketServer = new WebSocketServer(this.config.websocket || {});

    // 初始化各个服务器
    await this.initializeServers();

    // 注册默认路由
    await this.registerDefaultRoutes();

    this.isRunning = true;
    console.log('AIAgent Network Server initialized successfully');
    this.logServerStatus();
    return true;
  }

  async initializeServers() {
    // 初始化HTTP服务器
    if (this.config.http?.enabled !== false) {
      await this.httpServer.initialize();
      this.servers.set('http', this.httpServer);
    }

    // 初始化HTTPS服务器
    if (this.config.https?.enabled) {
      await this.httpsServer.initialize();
      this.servers.set('https', this.httpsServer);
    }

    // 初始化TCP服务器
    if (this.config.tcp?.enabled !== false) {
      await this.tcpServer.initialize();
      this.servers.set('tcp', this.tcpServer);
    }

    // 初始化UDP服务器
    if (this.config.udp?.enabled) {
      await this.udpServer.initialize();
      this.servers.set('udp', this.udpServer);
    }

    // 初始化WebSocket服务器
    if (this.config.websocket?.enabled !== false) {
      await this.websocketServer.initialize();
      this.servers.set('websocket', this.websocketServer);
    }
  }

  async registerDefaultRoutes() {
    // 注册默认路由
    await this.registerRoute('GET', '/health', async (req, res) => {
      return {
        status: 'ok',
        timestamp: Date.now(),
        servers: Array.from(this.servers.keys()),
        routes: Object.keys(this.routes)
      };
    });

    await this.registerRoute('GET', '/info', async (req, res) => {
      return {
        name: 'AIAgent Network Server',
        version: '1.0.0',
        servers: this.getServerStatuses()
      };
    });

    await this.registerRoute('POST', '/api/echo', async (req, res) => {
      return {
        status: 'success',
        message: 'Echo service',
        data: req.body
      };
    });
  }

  async registerRoute(method, path, handler) {
    console.log(`Registering route: ${method} ${path}`);
    const routeKey = `${method.toUpperCase()}:${path}`;
    this.routes[routeKey] = handler;

    // 注册到各个服务器
    for (const [name, server] of this.servers) {
      if (server.registerRoute) {
        await server.registerRoute(method, path, handler);
      }
    }

    return { success: true, route: routeKey };
  }

  async use(middleware) {
    console.log('Registering middleware...');
    this.middlewares.push(middleware);

    // 应用到各个服务器
    for (const [name, server] of this.servers) {
      if (server.use) {
        await server.use(middleware);
      }
    }

    return { success: true };
  }

  async start() {
    console.log('Starting AIAgent Network Server...');
    
    // 启动各个服务器
    for (const [name, server] of this.servers) {
      if (server.start) {
        await server.start();
      }
    }

    this.isRunning = true;
    console.log('AIAgent Network Server started successfully');
    this.logServerStatus();
    return true;
  }

  async stop() {
    console.log('Stopping AIAgent Network Server...');
    
    // 停止各个服务器
    for (const [name, server] of this.servers) {
      if (server.stop) {
        await server.stop();
      }
    }

    this.isRunning = false;
    console.log('AIAgent Network Server stopped successfully');
    return true;
  }

  getServerStatuses() {
    // 获取服务器状态
    const statuses = {};
    for (const [name, server] of this.servers) {
      statuses[name] = server.getStatus();
    }
    return statuses;
  }

  logServerStatus() {
    // 记录服务器状态
    console.log('Server status:');
    for (const [name, server] of this.servers) {
      const status = server.getStatus();
      console.log(`  ${name}: ${status.status} on ${status.address}:${status.port}`);
    }
  }

  async shutdown() {
    console.log('Shutting down AIAgent Network Server...');
    await this.stop();
    console.log('AIAgent Network Server shutdown successfully');
    return true;
  }
}

class HTTPServer {
  constructor(config) {
    this.config = config || {};
    this.port = config.port || 8080;
    this.host = config.host || '0.0.0.0';
    this.routes = new Map();
    this.middlewares = [];
    this.isRunning = false;
  }

  async initialize() {
    console.log(`Initializing HTTP Server on ${this.host}:${this.port}`);
    // 模拟HTTP服务器初始化
    return true;
  }

  async registerRoute(method, path, handler) {
    const routeKey = `${method.toUpperCase()}:${path}`;
    this.routes.set(routeKey, handler);
    return true;
  }

  async use(middleware) {
    this.middlewares.push(middleware);
    return true;
  }

  async start() {
    console.log(`Starting HTTP Server on ${this.host}:${this.port}`);
    this.isRunning = true;
    return true;
  }

  async stop() {
    console.log('Stopping HTTP Server...');
    this.isRunning = false;
    return true;
  }

  getStatus() {
    return {
      status: this.isRunning ? 'running' : 'stopped',
      address: this.host,
      port: this.port,
      protocol: 'http'
    };
  }
}

class HTTPSServer {
  constructor(config) {
    this.config = config || {};
    this.port = config.port || 8443;
    this.host = config.host || '0.0.0.0';
    this.cert = config.cert;
    this.key = config.key;
    this.routes = new Map();
    this.isRunning = false;
  }

  async initialize() {
    console.log(`Initializing HTTPS Server on ${this.host}:${this.port}`);
    // 模拟HTTPS服务器初始化
    return true;
  }

  async registerRoute(method, path, handler) {
    const routeKey = `${method.toUpperCase()}:${path}`;
    this.routes.set(routeKey, handler);
    return true;
  }

  async start() {
    console.log(`Starting HTTPS Server on ${this.host}:${this.port}`);
    this.isRunning = true;
    return true;
  }

  async stop() {
    console.log('Stopping HTTPS Server...');
    this.isRunning = false;
    return true;
  }

  getStatus() {
    return {
      status: this.isRunning ? 'running' : 'stopped',
      address: this.host,
      port: this.port,
      protocol: 'https'
    };
  }
}

class TCPServer {
  constructor(config) {
    this.config = config || {};
    this.port = config.port || 8081;
    this.host = config.host || '0.0.0.0';
    this.maxConnections = config.maxConnections || 1000;
    this.connections = new Set();
    this.isRunning = false;
  }

  async initialize() {
    console.log(`Initializing TCP Server on ${this.host}:${this.port}`);
    // 模拟TCP服务器初始化
    return true;
  }

  async start() {
    console.log(`Starting TCP Server on ${this.host}:${this.port}`);
    this.isRunning = true;
    return true;
  }

  async stop() {
    console.log('Stopping TCP Server...');
    this.isRunning = false;
    return true;
  }

  getStatus() {
    return {
      status: this.isRunning ? 'running' : 'stopped',
      address: this.host,
      port: this.port,
      protocol: 'tcp',
      connections: this.connections.size,
      maxConnections: this.maxConnections
    };
  }
}

class UDPServer {
  constructor(config) {
    this.config = config || {};
    this.port = config.port || 8082;
    this.host = config.host || '0.0.0.0';
    this.isRunning = false;
  }

  async initialize() {
    console.log(`Initializing UDP Server on ${this.host}:${this.port}`);
    // 模拟UDP服务器初始化
    return true;
  }

  async start() {
    console.log(`Starting UDP Server on ${this.host}:${this.port}`);
    this.isRunning = true;
    return true;
  }

  async stop() {
    console.log('Stopping UDP Server...');
    this.isRunning = false;
    return true;
  }

  getStatus() {
    return {
      status: this.isRunning ? 'running' : 'stopped',
      address: this.host,
      port: this.port,
      protocol: 'udp'
    };
  }
}

class WebSocketServer {
  constructor(config) {
    this.config = config || {};
    this.port = config.port || 8083;
    this.host = config.host || '0.0.0.0';
    this.clients = new Set();
    this.isRunning = false;
  }

  async initialize() {
    console.log(`Initializing WebSocket Server on ${this.host}:${this.port}`);
    // 模拟WebSocket服务器初始化
    return true;
  }

  async start() {
    console.log(`Starting WebSocket Server on ${this.host}:${this.port}`);
    this.isRunning = true;
    return true;
  }

  async stop() {
    console.log('Stopping WebSocket Server...');
    this.isRunning = false;
    return true;
  }

  getStatus() {
    return {
      status: this.isRunning ? 'running' : 'stopped',
      address: this.host,
      port: this.port,
      protocol: 'websocket',
      clients: this.clients.size
    };
  }
}

module.exports = {
  NetworkServer,
  HTTPServer,
  HTTPSServer,
  TCPServer,
  UDPServer,
  WebSocketServer
};