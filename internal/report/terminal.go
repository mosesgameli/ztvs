package report

import (
	"fmt"
	"strings"

	"github.com/mosesgameli/ztvs/pkg/rpc"
)

type TerminalReporter struct {
	findings map[string][]*rpc.Finding
}

func NewTerminal() *TerminalReporter {
	return &TerminalReporter{
		findings: make(map[string][]*rpc.Finding),
	}
}

func (r *TerminalReporter) AddFinding(pluginName string, finding *rpc.Finding) {
	r.findings[pluginName] = append(r.findings[pluginName], finding)
}

func (r *TerminalReporter) Flush() error {
	if len(r.findings) == 0 {
		fmt.Println("No findings found.")
		return nil
	}

	for plugin, findings := range r.findings {
		fmt.Printf("\n--- Plugin: %s (%d findings) ---\n", plugin, len(findings))
		for _, f := range findings {
			severity := strings.ToUpper(f.Severity)
			fmt.Printf("[%s] %s: %s\n", severity, f.Title, f.Description)
			if f.Remediation != "" {
				fmt.Printf("      Remediation: %s\n", f.Remediation)
			}
		}
	}
	fmt.Println()
	return nil
}
