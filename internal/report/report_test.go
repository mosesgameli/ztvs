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
	"bytes"
	"encoding/json"
	"testing"

	"github.com/mosesgameli/ztvs-sdk-go/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONReporter(t *testing.T) {
	r := NewJSON()
	finding := &rpc.Finding{
		ID:       "TEST-01",
		Title:    "Test Vulnerability",
		Severity: "HIGH",
	}
	r.AddFinding("test-plugin", finding)

	var buf bytes.Buffer
	r.SetOutput(&buf)
	err := r.Flush()
	require.NoError(t, err)

	var output map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &output)
	require.NoError(t, err)

	findings := output["findings"].(map[string]interface{})
	assert.Contains(t, findings, "test-plugin")
}

func TestTerminalReporter_VariousSeverities(t *testing.T) {
	r := NewTerminal()
	r.AddFinding("p1", &rpc.Finding{ID: "F1", Title: "T1", Severity: "CRITICAL", Description: "D1"})
	r.AddFinding("p2", &rpc.Finding{ID: "F2", Title: "T2", Severity: "HIGH", Description: "D2"})
	r.AddFinding("p3", &rpc.Finding{ID: "F3", Title: "T3", Severity: "MEDIUM", Description: "D3"})
	r.AddFinding("p4", &rpc.Finding{ID: "F4", Title: "T4", Severity: "LOW", Description: "D4"})

	var buf bytes.Buffer
	r.SetOutput(&buf)
	err := r.Flush()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "CRITICAL")
	assert.Contains(t, out, "HIGH")
	assert.Contains(t, out, "MED")
	assert.Contains(t, out, "LOW")
	assert.Contains(t, out, "COMPROMISED")
}

func TestTerminalReporter_Empty(t *testing.T) {
	r := NewTerminal()
	var buf bytes.Buffer
	r.SetOutput(&buf)
	err := r.Flush()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "No vulnerabilities found")
}

func TestTerminalReporter_Vulnerable(t *testing.T) {
	r := NewTerminal()
	r.AddFinding("p1", &rpc.Finding{ID: "F1", Title: "T1", Severity: "HIGH", Description: "D1"})

	var buf bytes.Buffer
	r.SetOutput(&buf)
	err := r.Flush()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "VULNERABLE")
}

func TestSARIFReporter(t *testing.T) {
	r := NewSARIF()
	finding := &rpc.Finding{
		ID:          "TEST-01",
		CheckID:     "TEST-01",
		Title:       "Test Vulnerability",
		Description: "This is a test description",
		Severity:    "MEDIUM",
	}
	r.AddFinding("test-plugin", finding)

	var buf bytes.Buffer
	r.SetOutput(&buf)
	err := r.Flush()
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "$schema")
	assert.Contains(t, out, "TEST-01")

	// Test with path for higher coverage
	r.AddFinding("p2", &rpc.Finding{
		ID:      "F2",
		CheckID: "C2",
		Evidence: map[string]interface{}{
			"file": "main.go",
		},
	})
	buf.Reset()
	r.Flush()
	assert.Contains(t, buf.String(), "main.go")

	// Test: Severities for full coverage
	r.AddFinding("p3", &rpc.Finding{Severity: "critical"})
	r.AddFinding("p4", &rpc.Finding{Severity: "info"})
	r.AddFinding("p5", nil) // Covered by defensive check
	buf.Reset()
	r.Flush()
	assert.Contains(t, buf.String(), "error") // critical maps to error
	assert.Contains(t, buf.String(), "note")  // info maps to note
}

func TestReport_EvidenceComplex(t *testing.T) {
	// 1. JSON with complex evidence
	rj := NewJSON()
	rj.AddFinding("p1", &rpc.Finding{
		ID: "F1",
		Evidence: map[string]interface{}{
			"nested": map[string]string{"key": "val"},
			"list":   []int{1, 2, 3},
		},
	})
	var buf bytes.Buffer
	rj.SetOutput(&buf)
	rj.Flush()
	assert.Contains(t, buf.String(), "nested")
	assert.Contains(t, buf.String(), "list")

	// 2. Terminal with nil finding (defensive check)
	rt := NewTerminal()
	rt.AddFinding("p1", nil)
	buf.Reset()
	rt.SetOutput(&buf)
	rt.Flush()
	assert.Contains(t, buf.String(), "No vulnerabilities found")
	
	// 3. SARIF with complex evidence
	rs := NewSARIF()
	rs.AddFinding("p1", &rpc.Finding{
		ID:          "F1",
		CheckID:     "C1",
		Title:       "Test Title",
		Description: "Test Desc",
		Evidence: map[string]interface{}{
			"snippet": "code",
			"nested":  map[string]string{"foo": "bar"},
		},
	})
	buf.Reset()
	rs.SetOutput(&buf)
	rs.Flush()
	assert.Contains(t, buf.String(), "C1")
	assert.Contains(t, buf.String(), "Test Title")
}
