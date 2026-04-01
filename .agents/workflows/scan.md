description: Build and run a vulnerability scan across all discovered plugins.

1. Build the engine and all plugins:
// turbo
`make build`

2. Execute a standard terminal-formatted scan:
// turbo
`./zt scan`

3. Generate a JSON report:
// turbo
`./zt --format json scan`

4. Generate a SARIF report for security dashboards:
// turbo
`./zt --format sarif scan`
