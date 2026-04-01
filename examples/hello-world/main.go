package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mosesgameli/ztvs/pkg/sdk"
)

// HelloWorldCheck implements a simple check that greets the host.
type HelloWorldCheck struct{}

func (c *HelloWorldCheck) ID() string   { return "hello_world" }
func (c *HelloWorldCheck) Name() string { return "Hello World Greet" }

func (c *HelloWorldCheck) Run(ctx context.Context) (*sdk.Finding, error) {
	hostname, _ := os.Hostname()

	return &sdk.Finding{
		ID:          "F-HELLO-001",
		Severity:    "info",
		Title:       "Greetings from Go Plugin",
		Description: fmt.Sprintf("The Hello World plugin is running on host: %s", hostname),
		Evidence: map[string]any{
			"runtime": "go",
			"message": "Hello from the polyglot-ready plugin system!",
		},
		Remediation: "No action required. This finding confirms the plugin is operational.",
	}, nil
}

func main() {
	// Plugin Metadata matches the plugin.yaml manifest
	meta := sdk.Metadata{
		Name:       "hello-world",
		Version:    "1.0.0",
		APIVersion: 1,
	}

	// sdk.Run handles the RPC handshake and check execution loop.
	// It expects --rpc flag to be passed by the ZTVS host.
	sdk.Run(meta, []sdk.Check{&HelloWorldCheck{}})
}
