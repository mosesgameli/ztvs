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
func (p *Policy) IsAllowed(pluginName string, requested []string) error {
	for _, cap := range requested {
		// 1. Check if explicitly blocked
		for _, blocked := range p.BlockedCapabilities {
			if cap == blocked {
				return fmt.Errorf("plugin %s requested blocked capability: %s", pluginName, cap)
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
				return fmt.Errorf("plugin %s requested unauthorized capability: %s", pluginName, cap)
			}
		}
	}
	return nil
}
