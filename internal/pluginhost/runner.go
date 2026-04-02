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

package pluginhost

import (
	"context"
	"fmt"
)

// Runner defines the interface for different plugin execution environments.
type Runner interface {
	Name() string
	Supports(runtimeType string) bool
	Validate(entrypoint string) error
	Execute(ctx context.Context, entrypoint string, stdin []byte) ([]byte, error)
}

// RunnerRegistry manages the available plugin runners.
type RunnerRegistry struct {
	runners []Runner
}

func NewRunnerRegistry() *RunnerRegistry {
	return &RunnerRegistry{
		runners: []Runner{},
	}
}

func (r *RunnerRegistry) Register(runner Runner) {
	r.runners = append(r.runners, runner)
}

func (r *RunnerRegistry) GetRunner(runtimeType string) (Runner, error) {
	for _, runner := range r.runners {
		if runner.Supports(runtimeType) {
			return runner, nil
		}
	}
	return nil, fmt.Errorf("no runner found for runtime type: %s", runtimeType)
}
