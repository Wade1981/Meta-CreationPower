// AIAgentFramework - Main entry point

// Import core modules
const { RunningService } = require('./worker/running-services');
const { AIBotDesigner } = require('./worker/bot-designer');
const { AgentConsole } = require('./worker/agent-console');
const { AgentApp } = require('./worker/agent-app');

const { ModelIDE } = require('./training/model-ide');
const { TrainingDataManager } = require('./training/training-data');
const { TrainingService } = require('./training/training-services');

const { NetworkProtocol } = require('./network/network-protocol');
const { NetworkNode } = require('./network/network-node');
const { NetworkServer } = require('./network/network-server');
const { NetworkClient } = require('./network/network-client');
const { MailService } = require('./network/network-mail');

const { BlockStorage } = require('./storage/block-storage');

class AIAgent {
  constructor(config) {
    this.config = config || {};
    this.name = this.config.name || 'AIAgent';
    this.modules = new Map();
  }

  async initialize() {
    console.log(`Initializing ${this.name}...`);
    // Initialize core modules
    if (this.config.worker) {
      const runningService = new RunningService(this.config.worker);
      await runningService.initialize();
      this.modules.set('runningService', runningService);
    }

    if (this.config.network) {
      const networkClient = new NetworkClient(this.config.network);
      await networkClient.initialize();
      this.modules.set('networkClient', networkClient);
    }

    if (this.config.storage) {
      const blockStorage = new BlockStorage(this.config.storage);
      await blockStorage.initialize();
      this.modules.set('blockStorage', blockStorage);
    }

    console.log(`${this.name} initialized successfully`);
    return true;
  }

  async execute(task, context) {
    console.log(`Executing task: ${task} with context:`, context);
    // Task execution logic goes here
    return { success: true, result: 'Task executed' };
  }

  async shutdown() {
    console.log(`Shutting down ${this.name}...`);
    // Shutdown all modules
    for (const [name, module] of this.modules) {
      if (module.shutdown) {
        await module.shutdown();
      }
    }
    console.log(`${this.name} shutdown successfully`);
    return true;
  }

  getModule(moduleName) {
    return this.modules.get(moduleName);
  }
}

class AgentFramework {
  constructor(config = {}) {
    this.config = config;
    this.agents = new Map();
    this.services = new Map();
    this.modules = new Map();
  }

  async initialize() {
    console.log('Initializing AIAgentFramework...');
    
    // Initialize core services
    if (this.config.worker) {
      const runningService = new RunningService(this.config.worker);
      await runningService.initialize();
      this.services.set('runningService', runningService);
    }

    if (this.config.training) {
      const trainingService = new TrainingService(this.config.training);
      await trainingService.initialize();
      this.services.set('trainingService', trainingService);
    }

    if (this.config.network) {
      const networkProtocol = new NetworkProtocol(this.config.network);
      await networkProtocol.initialize();
      this.services.set('networkProtocol', networkProtocol);

      const networkServer = new NetworkServer(this.config.network);
      await networkServer.initialize();
      this.services.set('networkServer', networkServer);
    }

    if (this.config.storage) {
      const blockStorage = new BlockStorage(this.config.storage);
      await blockStorage.initialize();
      this.services.set('blockStorage', blockStorage);
    }

    console.log('AIAgentFramework initialized successfully');
    return true;
  }

  async registerAgent(name, agentConfig) {
    console.log(`Registering agent: ${name}`);
    const agent = new AIAgent(agentConfig);
    await agent.initialize();
    this.agents.set(name, agent);
    console.log(`Agent ${name} registered successfully`);
    return { success: true, agentId: name };
  }

  getAgent(name) {
    return this.agents.get(name);
  }

  getService(serviceName) {
    return this.services.get(serviceName);
  }

  async start() {
    console.log('Starting AIAgentFramework...');
    
    // Start network server
    const networkServer = this.services.get('networkServer');
    if (networkServer) {
      await networkServer.start();
    }

    console.log('AIAgentFramework started successfully');
    return true;
  }

  async stop() {
    console.log('Stopping AIAgentFramework...');
    
    // Stop network server
    const networkServer = this.services.get('networkServer');
    if (networkServer) {
      await networkServer.stop();
    }

    console.log('AIAgentFramework stopped successfully');
    return true;
  }

  async shutdown() {
    console.log('Shutting down AIAgentFramework...');
    
    // Shutdown all agents
    for (const [name, agent] of this.agents) {
      await agent.shutdown();
    }

    // Shutdown all services
    for (const [name, service] of this.services) {
      if (service.shutdown) {
        await service.shutdown();
      }
    }

    console.log('AIAgentFramework shutdown successfully');
    return true;
  }
}

// Export all modules
module.exports = {
  // Core classes
  AIAgent,
  AgentFramework,
  
  // Worker modules
  RunningService,
  AIBotDesigner,
  AgentConsole,
  AgentApp,
  
  // Training modules
  ModelIDE,
  TrainingDataManager,
  TrainingService,
  
  // Network modules
  NetworkProtocol,
  NetworkNode,
  NetworkServer,
  NetworkClient,
  MailService,
  
  // Storage modules
  BlockStorage
};