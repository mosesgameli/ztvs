---
description: scan
---

description: Build and run a vulnerability scan along with security quality gates.

1. Build the engine and all plugins:
// turbo
`make build`

2. Run security scanning (govulncheck):
// turbo
`go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...`

3. Run static analysis (gosec):
// turbo
`go install github.com/securego/gosec/v2/cmd/gosec@latest && gosec ./...`

4. Execute a standard terminal-formatted scan:
// turbo
`./zt scan`

5. Generate a SARIF report for security dashboards:
// turbo
`./zt --format sarif scan`
