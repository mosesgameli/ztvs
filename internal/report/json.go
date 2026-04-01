package report

import (
	"encoding/json"
	"os"
	"time"

	"github.com/mosesgameli/ztvs-sdk-go/rpc"
)

type JSONReport struct {
	Timestamp time.Time                 `json:"timestamp"`
	Summary   map[string]int            `json:"summary"`
	Findings  map[string][]*rpc.Finding `json:"findings"`
}

type JSONReporter struct {
	report JSONReport
}

func NewJSON() *JSONReporter {
	return &JSONReporter{
		report: JSONReport{
			Timestamp: time.Now(),
			Summary:   make(map[string]int),
			Findings:  make(map[string][]*rpc.Finding),
		},
	}
}

func (r *JSONReporter) AddFinding(pluginName string, finding *rpc.Finding) {
	r.report.Findings[pluginName] = append(r.report.Findings[pluginName], finding)
	r.report.Summary[finding.Severity]++
}

func (r *JSONReporter) Flush() error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(r.report)
}
