package main

import (
	"context"
	"github.com/mosesgameli/ztvs-sdk-go/sdk"
)

type GoExampleCheck struct{}

func (c *GoExampleCheck) ID() string   { return "go_check" }
func (c *GoExampleCheck) Name() string { return "Go Example Check" }

func (c *GoExampleCheck) Run(ctx context.Context) (*sdk.Finding, error) {
	return &sdk.Finding{
		ID:          "F-GO-001",
		Severity:    "info",
		Title:       "Go Plugin Running",
		Description: "The first-party Go SDK example plugin is executing correctly.",
		Evidence:    map[string]any{"sdk": "ztvs/pkg/sdk"},
		Remediation: "No action required.",
	}, nil
}

func main() {
	meta := sdk.Metadata{
		Name:       "example-go",
		Version:    "1.0.0",
		APIVersion: 1,
	}

	sdk.Run(meta, []sdk.Check{&GoExampleCheck{}})
}
