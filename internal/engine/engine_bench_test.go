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

package engine

import (
	"testing"

	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/internal/pluginhost"
)

func BenchmarkScan(b *testing.B) {
	cfg := &config.Config{}
	reporter := &MockReporter{}
	e := New(cfg, pluginhost.New(), reporter, pluginhost.NewRegistry())
	e.Interactive = false

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Note: This benchmarks the Scan orchestration.
		// Real performance depends on plugin execution overhead.
		_ = e.Scan()
	}
}
