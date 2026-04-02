// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package report

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mosesgameli/ztvs-sdk-go/rpc"
	"github.com/pterm/pterm"
)

type TerminalReporter struct {
	findings map[string][]*rpc.Finding
	output   io.Writer
}

func NewTerminal() *TerminalReporter {
	return &TerminalReporter{
		findings: make(map[string][]*rpc.Finding),
		output:   os.Stdout,
	}
}

func (r *TerminalReporter) SetOutput(w io.Writer) {
	r.output = w
}

func (r *TerminalReporter) AddFinding(pluginName string, finding *rpc.Finding) {
	r.findings[pluginName] = append(r.findings[pluginName], finding)
}

func (r *TerminalReporter) Flush() error {
	pterm.SetDefaultOutput(r.output)
	if len(r.findings) == 0 {
		pterm.DefaultSection.Println("Scan Results")
		pterm.Success.Println("No vulnerabilities found! System is clean.")
		return nil
	}

	totalFindings := 0
	criticals := 0
	highs := 0

	for plugin, findings := range r.findings {
		pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).WithMargin(2).Printf("Source: %s", plugin)
		fmt.Println()

		for _, f := range findings {
			totalFindings++
			sev := strings.ToUpper(f.Severity)

			var sevBadge string
			var accentStyle *pterm.Style

			switch sev {
			case "CRITICAL":
				criticals++
				sevBadge = pterm.BgRed.Sprint(pterm.FgBlack.Sprint(" ! CRITICAL "))
				accentStyle = pterm.NewStyle(pterm.FgRed, pterm.Bold)
			case "HIGH":
				highs++
				sevBadge = pterm.BgLightRed.Sprint(pterm.FgBlack.Sprint(" HIGH "))
				accentStyle = pterm.NewStyle(pterm.FgLightRed)
			case "MEDIUM":
				sevBadge = pterm.BgYellow.Sprint(pterm.FgBlack.Sprint(" MED "))
				accentStyle = pterm.NewStyle(pterm.FgYellow)
			default:
				sevBadge = pterm.BgCyan.Sprint(pterm.FgBlack.Sprint(" LOW "))
				accentStyle = pterm.NewStyle(pterm.FgCyan)
			}

			cardContent := fmt.Sprintf("%s [%s] %s\n\n", sevBadge, f.ID, accentStyle.Sprint(f.Title))
			cardContent += pterm.LightWhite(f.Description) + "\n"
			if f.Remediation != "" {
				cardContent += "\n" + pterm.BgGreen.Sprint(pterm.FgBlack.Sprint(" FIX ")) + " " + pterm.FgGreen.Sprint(f.Remediation)
			}

			_ = pterm.DefaultPanel.WithPadding(1).WithPanels(pterm.Panels{
				{{Data: cardContent}},
			}).Render()
			fmt.Println()
		}
	}

	// Final Summary Dashboard
	summaryTitle := pterm.Bold.Sprint("AUDIT SUMMARY")
	stats, _ := pterm.DefaultTable.WithData([][]string{
		{"Metric", "Value"},
		{"Total Plugins Reporting", fmt.Sprintf("%d", len(r.findings))},
		{"Total Findings Identified", fmt.Sprintf("%d", totalFindings)},
		{"Critical Alerts", pterm.FgRed.Sprint(criticals)},
		{"High Probability Risks", pterm.FgLightRed.Sprint(highs)},
	}).Srender()

	status := pterm.Success.Sprint("SECURE")
	if criticals > 0 {
		status = pterm.BgRed.Sprint(pterm.FgBlack.Sprint(" COMPROMISED "))
	} else if highs > 0 {
		status = pterm.FgRed.Sprint("VULNERABLE")
	}

	_ = pterm.DefaultPanel.WithPadding(2).WithPanels(pterm.Panels{
		{{Data: fmt.Sprintf("%s\n\n%s\nSystem Status: %s", summaryTitle, stats, status)}},
	}).Render()

	return nil
}
