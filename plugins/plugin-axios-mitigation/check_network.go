package main

import (
	"context"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mosesgameli/ztvs-sdk-go/sdk"
)

type NetworkCheck struct{}

func (c *NetworkCheck) ID() string {
	return "axios_c2_network"
}

func (c *NetworkCheck) Name() string {
	return "Axios C2 Network Connection Probe"
}

func (c *NetworkCheck) Run(ctx context.Context) (*sdk.Finding, error) {
	finding := &sdk.Finding{
		ID:          "F-AXIOS-003",
		Severity:    "critical",
		Title:       "Active connection to Axios C2 server detected",
		Evidence:    make(map[string]interface{}),
		Remediation: "Terminate the process and report the incident. Rotate secrets and check for lateral movement.",
	}

	found := false
	c2_ip := "142.11.206.73"
	c2_host := "sfrclak.com"

	// Network connection check
	switch runtime.GOOS {
	case "darwin", "linux":
		// Use lsof or netstat
		cmd := exec.Command("lsof", "-i", "-n", "-P")
		out, err := cmd.Output()
		if err == nil {
			if strings.Contains(string(out), c2_ip) || strings.Contains(string(out), c2_host) {
				finding.Evidence["connection"] = "Detected active connection via lsof"
				found = true
			}
		}
	case "windows":
		// Use Get-NetTCPConnection
		cmd := exec.Command("powershell", "-Command", "Get-NetTCPConnection | Where-Object { $_.RemoteAddress -eq '"+c2_ip+"' }")
		out, err := cmd.Output()
		if err == nil && len(out) > 0 {
			finding.Evidence["connection"] = "Detected active connection via Get-NetTCPConnection"
			found = true
		}
	}

	if found {
		return finding, nil
	}

	return nil, nil
}
