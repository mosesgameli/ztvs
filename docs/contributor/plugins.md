# ZTVS Plugin Development & SDK Guide

This guide covers the core requirements for building ZTVS plugins and provides a reference for the first-party Go SDK.

 

## 1. Core Plugin Requirements

ZTVS is designed to be truly language agnostic. Any language that can execute on the target platform and communicate via the JSON-RPC standard is a first-class plugin developer.

### Protocol Standards
-   **Binary Executable**: The plugin must be an executable binary or script.
-   **JSON-RPC Stdio**: The plugin must read requests from `stdin` and write responses to `stdout`.
-   **No Stdout Logging**: `stdout` is reserved for protocol messages. All logs should be written to `stderr`.

### Implementation Checklist
-   [ ]   **Manifest (`plugin.yaml`)**: Every plugin must include a manifest declaring its name, version, and capabilities.
-   [ ]   **Handshake**: Implement the `handshake` method. This is the first call made by the host.
-   [ ]   **Check IDs**: Define unique check IDs that map to security validations.
-   [ ]   **Run Check**: Implement the `run_check` method to execute security logic and return status.
-   [ ]   **Capabilities**: Only perform actions for which you have declared capabilities. Users must explicitly whitelist required capabilities like `network_access` or `execute_commands` inside their global `~/.ztvs/config.yaml`.

 

## 2. Go SDK Reference

The ZTVS Go SDK provides a robust framework for building compatible plugins.

### Installation
The SDK is available as a separate repository:
```bash
go get github.com/mosesgameli/ztvs-sdk-go/sdk
```

### Basic Structure
A ZTVS plugin in Go implements the `Check` interface:
```go
type Check interface {
	ID() string
	Name() string
	Run(ctx context.Context) (*Finding, error)
}
```

### Example Plugin Implementation
```go
package main

import (
	"context"
	"github.com/mosesgameli/ztvs-sdk-go/sdk"
)

type MyCheck struct{}

func (c *MyCheck) ID() string   { return "my_check" }
func (c *MyCheck) Name() string { return "My Custom Check" }

func (c *MyCheck) Run(ctx context.Context) (*sdk.Finding, error) {
	return &sdk.Finding{
		ID: "F-001",
		Severity: "high",
		Title: "Vulnerability Found",
		Description: "Found a critical issue.",
		Evidence: map[string]any{"path": "/tmp/test"},
		Remediation: "Remove the file.",
	}, nil
}

func main() {
	meta := sdk.Metadata{
		Name: "my-plugin",
		Version: "1.0.0",
		APIVersion: 1,
	}

	sdk.Run(meta, []sdk.Check{&MyCheck{}})
}
```

 

## 3. `Finding` Schema

| Field | Type | Description |
| :--- | :--- | :--- |
| `ID` | `string` | Unique identifier for findings (e.g., `F-001`). |
| `Severity` | `string` | Risk level (`critical`, `high`, `medium`, `low`, `info`). |
| `Title` | `string` | Short description of the issue. |
| `Description`| `string` | Full explanation of the vulnerability. |
| `Evidence` | `map[string]any`| Key-value data to support the finding. |
| `Remediation`| `string` | Instructions to fix the problem. |

---

## 4. Testing & Validation

The SDK is designed to be executed by the ZTVS host, but you can test it manually:
```bash
echo '{"jsonrpc":"2.0","id":"1","method":"handshake","params":{}}' | ./my-plugin
```

For more details on the protocol, see the [Protocol Specification](../protocol/README.md).
