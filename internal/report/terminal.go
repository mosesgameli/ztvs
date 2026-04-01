package report

import (
	"fmt"
	"strings"

	"github.com/mosesgameli/ztvs/pkg/rpc"
	"github.com/pterm/pterm"
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
		pterm.Success.Println("No vulnerabilities found! System is clean.")
		return nil
	}

	for plugin, findings := range r.findings {
		pterm.DefaultHeader.WithFullWidth().Printf("Plugin: %s (%d findings)", plugin, len(findings))

		var tableData [][]string
		tableData = append(tableData, []string{"Severity", "Title", "Description", "Remediation"})

		for _, f := range findings {
			sev := strings.ToUpper(f.Severity)
			var coloredSev string
			switch sev {
			case "CRITICAL":
				coloredSev = pterm.BgRed.Sprint(pterm.FgBlack.Sprint(" " + sev + " "))
			case "HIGH":
				coloredSev = pterm.FgRed.Sprint(sev)
			case "MEDIUM":
				coloredSev = pterm.FgYellow.Sprint(sev)
			case "LOW":
				coloredSev = pterm.FgCyan.Sprint(sev)
			default:
				coloredSev = pterm.FgLightBlue.Sprint(sev)
			}

			tableData = append(tableData, []string{
				coloredSev,
				f.Title,
				f.Description,
				f.Remediation,
			})
		}
		
		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
		fmt.Println()
	}
	return nil
}
