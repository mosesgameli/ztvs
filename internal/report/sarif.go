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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/mosesgameli/ztvs-sdk-go/rpc"
)

// Simplified SARIF structures
type SARIFReport struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []SARIFRun `json:"runs"`
}

type SARIFRun struct {
	Tool    SARIFTool     `json:"tool"`
	Results []SARIFResult `json:"results"`
}

type SARIFTool struct {
	Driver SARIFDriver `json:"driver"`
}

type SARIFDriver struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type SARIFResult struct {
	RuleID  string          `json:"ruleId"`
	Message SARIFMessage    `json:"message"`
	Level   string          `json:"level,omitempty"`
	Places  []SARIFLocation `json:"locations,omitempty"`
}

type SARIFMessage struct {
	Text string `json:"text"`
}

type SARIFLocation struct {
	PhysicalLocation SARIFPhysicalLocation `json:"physicalLocation"`
}

type SARIFPhysicalLocation struct {
	ArtifactLocation SARIFArtifactLocation `json:"artifactLocation"`
}

type SARIFArtifactLocation struct {
	URI string `json:"uri"`
}

type SARIFReporter struct {
	report SARIFReport
	output io.Writer
	mu     sync.RWMutex
}

func NewSARIF() *SARIFReporter {
	return &SARIFReporter{
		report: SARIFReport{
			Version: "2.1.0",
			Schema:  "https://schemastore.azurewebsites.net/schemas/json/sarif-2.1.0-rtm.5.json",
			Runs: []SARIFRun{
				{
					Tool: SARIFTool{
						Driver: SARIFDriver{
							Name:    "ZTVS",
							Version: "1.0.0",
						},
					},
					Results: []SARIFResult{},
				},
			},
		},
		output: os.Stdout,
	}
}

func (r *SARIFReporter) SetOutput(w io.Writer) {
	r.output = w
}

func (r *SARIFReporter) AddFinding(pluginName string, finding *rpc.Finding) {
	if finding == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	level := "warning"
	if finding.Severity == "critical" || finding.Severity == "high" {
		level = "error"
	} else if finding.Severity == "info" {
		level = "note"
	}

	result := SARIFResult{
		RuleID: finding.CheckID,
		Message: SARIFMessage{
			Text: fmt.Sprintf("%s: %s", finding.Title, finding.Description),
		},
		Level: level,
	}

	// Add evidence location if available
	if file, ok := finding.Evidence["file"].(string); ok {
		result.Places = []SARIFLocation{
			{
				PhysicalLocation: SARIFPhysicalLocation{
					ArtifactLocation: SARIFArtifactLocation{
						URI: file,
					},
				},
			},
		}
	}

	r.report.Runs[0].Results = append(r.report.Runs[0].Results, result)
}

func (r *SARIFReporter) Flush() error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	encoder := json.NewEncoder(r.output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(r.report)
}
