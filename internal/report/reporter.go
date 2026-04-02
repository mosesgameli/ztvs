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
	"github.com/mosesgameli/ztvs-sdk-go/rpc"
	"io"
)

// Reporter defines the interface for different output formats
type Reporter interface {
	// AddFinding adds a single finding to the report
	AddFinding(pluginName string, finding *rpc.Finding)
	// Flush finalizes the report and writes it to the output
	Flush() error
	// SetOutput sets the output writer for Flush
	SetOutput(w io.Writer)
}
