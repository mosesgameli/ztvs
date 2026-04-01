#!/usr/bin/env node

const readline = require('readline');

const rl = readline.createInterface({
  input: process.stdin,
  terminal: false
});

rl.on('line', (line) => {
  try {
    const req = JSON.parse(line);
    const { method, id } = req;

    if (method === 'handshake') {
      const resp = {
        jsonrpc: "2.0",
        id: id,
        result: {
          name: "example-nodejs",
          version: "1.0.0",
          api_version: 1,
          checks_supported: ["nodejs_check"]
        }
      };
      console.log(JSON.stringify(resp));
    } else if (method === 'run_check') {
      const resp = {
        jsonrpc: "2.0",
        id: id,
        result: {
          status: "pass",
          finding: {
            id: "F-JS-001",
            check_id: "nodejs_check",
            severity: "info",
            title: "Node.js Plugin Running",
            description: "Node.js-based plugin is communicating correctly.",
            evidence: { "runtime": "node.js", "version": process.version },
            remediation: "None"
          }
        }
      };
      console.log(JSON.stringify(resp));
    }
  } catch (err) {
    process.stderr.write(`Error: ${err.message}\n`);
  }
});
