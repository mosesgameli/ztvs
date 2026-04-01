🔙 [Back to Home](../README.md)

# ZTVS Go SDK Reference

The ZTVS Go SDK provides a robust, easy-to-use framework for building ZTVS-compatible security plugins in Go. It handles the JSON-RPC communication, input parsing, and output formatting.

## Installation

The SDK is available as a separate repository:

```bash
go get github.com/mosesgameli/ztvs-sdk-go/sdk
```

Add the following to your `main.go`:

```go
import "github.com/mosesgameli/ztvs-sdk-go/sdk"
```

## Basic Structure

A ZTVS plugin in Go is defined by its metadata and a collection of checks that implement the `Check` interface.

```go
type Check interface {
	ID() string
	Name() string
	Run(ctx context.Context) (*Finding, error)
}
```

## Example Plugin

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

## `Finding` Schema

| Field | Type | Description |
| :--- | :--- | :--- |
| `ID` | `string` | Unique identifier for findings (e.g., `F-001`). |
| `Severity` | `string` | Risk level (`critical`, `high`, `medium`, `low`, `info`). |
| `Title` | `string` | Short description of the issue. |
| `Description`| `string` | Full explanation of the vulnerability. |
| `Evidence` | `map[string]any`| Key-value data to support the finding. |
| `Remediation`| `string` | Instructions to fix the problem. |

## Running & Testing

The SDK is designed to be executed by the ZTVS host, but you can test it manually from the command line:

```bash
echo '{"jsonrpc":"2.0","id":"1","method":"handshake","params":{}}' | ./my-plugin
```
