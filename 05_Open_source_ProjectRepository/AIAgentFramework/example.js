// Example usage of AIAgentFramework

const { AIAgent, AgentFramework } = require('./src/index');

async function runExample() {
  // Create framework instance
  const framework = new AgentFramework();

  // Create and register agents
  const assistantAgent = new AIAgent({ name: 'AssistantAgent' });
  const toolAgent = new AIAgent({ name: 'ToolAgent' });

  framework.registerAgent('assistant', assistantAgent);
  framework.registerAgent('tool', toolAgent);

  // Initialize all agents
  await framework.initializeAll();

  // Execute tasks with agents
  const assistant = framework.getAgent('assistant');
  const tool = framework.getAgent('tool');

  const assistantResult = await assistant.execute('greet', { user: 'John' });
  console.log('Assistant result:', assistantResult);

  const toolResult = await tool.execute('calculate', { operation: 'add', numbers: [5, 3] });
  console.log('Tool result:', toolResult);

  // Shutdown all agents
  await framework.shutdownAll();

  console.log('Example completed successfully!');
}

// Run the example
runExample().catch(console.error);