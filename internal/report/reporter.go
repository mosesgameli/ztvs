package report

import (
	"github.com/mosesgameli/ztvs-sdk-go/rpc"
)

// Reporter defines the interface for different output formats
type Reporter interface {
	// AddFinding adds a single finding to the report
	AddFinding(pluginName string, finding *rpc.Finding)
	// Flush finalizes the report and writes it to the output
	Flush() error
}
