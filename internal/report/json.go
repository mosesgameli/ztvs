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
	"io"
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
	output io.Writer
}

func NewJSON() *JSONReporter {
	return &JSONReporter{
		report: JSONReport{
			Timestamp: time.Now(),
			Summary:   make(map[string]int),
			Findings:  make(map[string][]*rpc.Finding),
		},
		output: os.Stdout,
	}
}

func (r *JSONReporter) AddFinding(pluginName string, finding *rpc.Finding) {
	r.report.Findings[pluginName] = append(r.report.Findings[pluginName], finding)
	r.report.Summary[finding.Severity]++
}

func (r *JSONReporter) Flush() error {
	encoder := json.NewEncoder(r.output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(r.report)
}

func (r *JSONReporter) SetOutput(w io.Writer) {
	r.output = w
}
