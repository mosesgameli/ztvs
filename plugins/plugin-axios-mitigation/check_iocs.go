package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mosesgameli/ztvs/pkg/sdk"
)

type HostIOCCheck struct{}

func (c *HostIOCCheck) ID() string {
	return "axios_host_ioc"
}

func (c *HostIOCCheck) Name() string {
	return "Axios Host IOC Signature Scan"
}

func (c *HostIOCCheck) Run(ctx context.Context) (*sdk.Finding, error) {
	finding := &sdk.Finding{
		ID:          "F-AXIOS-002",
		Severity:    "critical",
		Title:       "Axios RAT Indicators of Compromise detected",
		Evidence:    make(map[string]interface{}),
		Remediation: "Identify the malicious process and terminate it. Rebuild the system from a clean image.",
	}

	found := false

	// File IOCs
	var targetFile string
	switch runtime.GOOS {
	case "darwin":
		targetFile = "/Library/Caches/com.apple.act.mond"
	case "windows":
		targetFile = os.Getenv("PROGRAMDATA") + "\\wt.exe"
	case "linux":
		targetFile = "/tmp/ld.py"
	}

	if targetFile != "" {
		if _, err := os.Stat(targetFile); err == nil {
			absPath, err := filepath.Abs(targetFile)
			if err != nil {
				absPath = targetFile
			}
			finding.Title = fmt.Sprintf("Axios RAT file detected: %s", absPath)
			finding.Description = fmt.Sprintf("Malicious file found at full path: %s", absPath)
			finding.Evidence["malicious_file"] = absPath
			found = true
		}
	}

	// Windows process check: PowerShell masquerading as wt.exe
	if runtime.GOOS == "windows" {
		if c.checkWindowsProcess(finding) {
			found = true
		}
	}

	if found {
		return finding, nil
	}

	return nil, nil
}

func (c *HostIOCCheck) checkWindowsProcess(finding *sdk.Finding) bool {
	// Look for wt.exe which is actually PowerShell
	cmd := exec.Command("powershell", "-Command", "Get-Process wt -ErrorAction SilentlyContinue | Select-Object Path")
	out, err := cmd.Output()
	if err == nil && len(out) > 0 {
		path := strings.TrimSpace(string(out))
		if strings.Contains(path, "powershell") {
			finding.Title = fmt.Sprintf("Axios RAT process detected: wt.exe (PowerShell masquerade) at %s", path)
			finding.Description = fmt.Sprintf("Malicious process found at full path: %s", path)
			finding.Evidence["malicious_process"] = "wt.exe (PowerShell masquerade)"
			finding.Evidence["process_path"] = path
			return true
		}
	}
	return false
}
