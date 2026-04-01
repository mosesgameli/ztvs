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
