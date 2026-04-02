const readline = require('readline');

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
  terminal: false
});

rl.on('line', (line) => {
  try {
    const request = JSON.parse(line);
    const { id, method } = request;

    if (method === 'handshake') {
      const response = {
        jsonrpc: '2.0',
        id: id,
        result: {
          name: 'node-fixture',
          version: '1.0.0',
          api_version: 1,
          checks_supported: ['node_test_check']
        }
      };
      process.stdout.write(JSON.stringify(response) + '\n');
    } else if (method === 'run_check') {
      const { params } = request;
      if (params.check_id === 'node_test_check') {
        const response = {
          jsonrpc: '2.0',
          id: id,
          result: {
            status: 'pass',
            finding: {
              id: 'F-JS-001',
              severity: 'info',
              title: 'Node.js Fixture Executed',
              description: 'This is a successful result from the raw Node.js polyglot fixture.'
            }
          }
        };
        process.stdout.write(JSON.stringify(response) + '\n');
      }
    }
  } catch (err) {
    process.stderr.write(`Error: ${err.message}\n`);
  }
});
