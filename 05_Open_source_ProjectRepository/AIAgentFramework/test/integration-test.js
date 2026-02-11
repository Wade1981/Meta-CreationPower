// AIAgentFramework - Integration Test

const { AgentFramework } = require('../src/index');

async function runIntegrationTest() {
  console.log('=== AIAgentFramework Integration Test ===');

  try {
    // Create framework instance
    const framework = new AgentFramework({
      worker: {
        enabled: true
      },
      training: {
        enabled: true
      },
      network: {
        http: {
          port: 8080
        },
        websocket: {
          port: 8083
        }
      },
      storage: {
        localDirectory: './test-storage'
      }
    });

    console.log('1. Initializing framework...');
    await framework.initialize();
    console.log('✓ Framework initialized successfully');

    console.log('\n2. Registering agents...');
    await framework.registerAgent('assistant', {
      name: 'AssistantAgent',
      worker: {
        enabled: true
      },
      network: {
        enabled: true
      }
    });

    await framework.registerAgent('tool', {
      name: 'ToolAgent',
      worker: {
        enabled: true
      },
      storage: {
        enabled: true
      }
    });
    console.log('✓ Agents registered successfully');

    console.log('\n3. Starting framework...');
    await framework.start();
    console.log('✓ Framework started successfully');

    console.log('\n4. Testing services...');

    // Test running service
    const runningService = framework.getService('runningService');
    if (runningService) {
      console.log('  Testing RunningService...');
      const modelDeployment = await runningService.deployModel('test-model', {
        type: 'classification',
        parameters: {
          epochs: 10,
          batchSize: 32
        }
      });
      console.log('  ✓ Model deployed:', modelDeployment);

      const taskResult = await runningService.executeTask('test-task-1', {
        action: 'test',
        parameters: {
          message: 'Hello, world!'
        }
      });
      console.log('  ✓ Task executed:', taskResult);
    }

    // Test storage service
    const blockStorage = framework.getService('blockStorage');
    if (blockStorage) {
      console.log('  Testing BlockStorage...');
      const storeResult = await blockStorage.storeData('test-data-1', {
        message: 'Test data for storage'
      }, 'runtime');
      console.log('  ✓ Data stored:', storeResult);

      const retrieveResult = await blockStorage.retrieveData('test-data-1', 'runtime');
      console.log('  ✓ Data retrieved:', retrieveResult);

      const listResult = await blockStorage.listData('runtime');
      console.log('  ✓ Data listed:', listResult);
    }

    console.log('✓ Services tested successfully');

    console.log('\n5. Testing agents...');
    const assistantAgent = framework.getAgent('assistant');
    if (assistantAgent) {
      const agentResult = await assistantAgent.execute('greet', {
        user: 'Test User'
      });
      console.log('  ✓ Assistant agent executed:', agentResult);
    }

    const toolAgent = framework.getAgent('tool');
    if (toolAgent) {
      const agentResult = await toolAgent.execute('calculate', {
        operation: 'add',
        numbers: [5, 3]
      });
      console.log('  ✓ Tool agent executed:', agentResult);
    }
    console.log('✓ Agents tested successfully');

    console.log('\n6. Stopping framework...');
    await framework.stop();
    console.log('✓ Framework stopped successfully');

    console.log('\n7. Shutting down framework...');
    await framework.shutdown();
    console.log('✓ Framework shutdown successfully');

    console.log('\n=== Integration Test PASSED ===');
    console.log('All modules are working correctly!');

  } catch (error) {
    console.error('=== Integration Test FAILED ===');
    console.error('Error:', error);
    console.error('Stack:', error.stack);
    process.exit(1);
  }
}

// Run the test
runIntegrationTest();