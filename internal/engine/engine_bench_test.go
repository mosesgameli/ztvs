package engine

import (
	"io"
	"testing"

	"github.com/mosesgameli/ztvs-sdk-go/rpc"
	"github.com/mosesgameli/ztvs/internal/config"
)

// MockReporter satisfies report.Reporter without doing anything
type MockReporter struct{}

func (m *MockReporter) AddFinding(p string, f *rpc.Finding) {}
func (m *MockReporter) Flush() error                         { return nil }
func (m *MockReporter) Write(w io.Writer) error              { return nil }

func BenchmarkScan(b *testing.B) {
	cfg := &config.Config{}
	reporter := &MockReporter{}
	e := New(cfg, reporter)
	e.Interactive = false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Note: This benchmarks the Scan orchestration. 
		// Real performance depends on plugin execution overhead.
		_ = e.Scan()
	}
}
