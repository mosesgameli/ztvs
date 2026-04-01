# Skill: Developing ZTVS Plugins

This skill enables the development and integration of new vulnerability scanning plugins for the Zero Trust Vulnerability Scanner (ZTVS).

## Overview
ZTVS plugins are separate Go executables that communicate with the Host via JSON-RPC 2.0 over `stdin` and `stdout`.

## 🛠️ Step-by-Step Implementation

### 1. Project Setup
Plugins MUST reside in `plugins/<plugin-name>/`.
Create a `main.go` and a separate file for each security check (e.g., `ssh_check.go`).

### 2. Implementation using SDK
Use `github.com/mosesgameli/ztvs/pkg/sdk` to simplify the build process.

```go
package main

import (
    "github.com/mosesgameli/ztvs/pkg/sdk"
    "context"
)

func main() {
    sdk.Run(sdk.Metadata{
        Name:       "plugin-example",
        Version:    "1.0.0",
        APIVersion: 1,
    }, []sdk.Check{
        &MyCheck{},
    })
}

type MyCheck struct{}

func (c *MyCheck) ID() string   { return "example_check" }
func (c *MyCheck) Name() string { return "Example Security Check" }

func (c *MyCheck) Run(ctx context.Context) (*sdk.Finding, error) {
    // Implement security logic here
    return &sdk.Finding{
        ID:          "F-EX-001",
        Severity:    "medium",
        Title:       "Example finding",
        Description: "Something was found",
    }, nil
}
```

### 3. Protocol Requirements
- **Handshake**: The SDK handles the handshake automatically. Ensure `APIVersion: 1` is used.
- **RPC Support**: Plugins MUST support the `--rpc` flag (handled by `sdk.Run`).
- **Isolation**: Never use global state that could interfere with parallel check execution.

### 4. Build & Distribution
1. Build the plugin: `go build -o plugin-example ./plugins/plugin-example`
2. Test discovery: Run `./zt scan` and ensure the new plugin is listed.

## 🛡️ Security Best Practices
- **Read-Only**: Plugins SHOULD avoid writing to the disk.
- **Evidence**: Always include `Evidence` (e.g., file paths, command output) in the finding.
- **Timeouts**: The Host enforces a 30s timeout per check; keep logic efficient.
