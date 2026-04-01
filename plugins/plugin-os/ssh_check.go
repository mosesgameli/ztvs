package main

import (
	"context"

	"github.com/mosesgameli/ztvs/pkg/sdk"
)

type SSHCheck struct{}

func (c *SSHCheck) ID() string {
	return "ssh_config"
}

func (c *SSHCheck) Name() string {
	return "SSH Configuration Check"
}

func (c *SSHCheck) Run(
	ctx context.Context,
) (*sdk.Finding, error) {

	// For MVP, we simulate a finding.
	// In a real implementation, we'd read /etc/ssh/sshd_config
	return &sdk.Finding{
		ID:          "F-SSH-001",
		Severity:    "high",
		Title:       "Root login enabled",
		Description: "PermitRootLogin yes found",
		Evidence: map[string]interface{}{
			"file":  "/etc/ssh/sshd_config",
			"value": "PermitRootLogin yes",
		},
		Remediation: "Set PermitRootLogin no",
	}, nil
}
