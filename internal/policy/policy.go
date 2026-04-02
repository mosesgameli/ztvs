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

package policy

import (
	"fmt"

	"github.com/mosesgameli/ztvs/internal/config"
)

// Policy defines the host's security rules for plugin execution.
type Policy struct {
	AllowedCapabilities []string
	BlockedCapabilities []string
}

func New(cfg *config.Config) *Policy {
	return &Policy{
		AllowedCapabilities: cfg.Policy.AllowedCapabilities,
		BlockedCapabilities: cfg.Policy.BlockedCapabilities,
	}
}

// IsAllowed checks if a plugin's requested capabilities are permitted by the host.
// It returns the specific capability that failed, and the error.
func (p *Policy) IsAllowed(pluginName string, requested []string) (string, error) {
	for _, cap := range requested {
		// 1. Check if explicitly blocked
		for _, blocked := range p.BlockedCapabilities {
			if cap == blocked {
				return cap, fmt.Errorf("plugin %s requested blocked capability: %s", pluginName, cap)
			}
		}

		// 2. Check if not in allowed list (if allowed list is defined)
		if len(p.AllowedCapabilities) > 0 {
			found := false
			for _, allowed := range p.AllowedCapabilities {
				if cap == allowed {
					found = true
					break
				}
			}
			if !found {
				return cap, fmt.Errorf("plugin %s requested unauthorized capability: %s", pluginName, cap)
			}
		}
	}
	return "", nil
}

// Reload updates the in-memory policy with fresh configuration
func (p *Policy) Reload(cfg *config.Config) {
	p.AllowedCapabilities = cfg.Policy.AllowedCapabilities
	p.BlockedCapabilities = cfg.Policy.BlockedCapabilities
}
