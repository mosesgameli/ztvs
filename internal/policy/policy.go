package policy

import (
	"fmt"
)

// Policy defines the host's security rules for plugin execution.
type Policy struct {
	AllowedCapabilities []string
	BlockedCapabilities []string
}

func NewDefault() *Policy {
	return &Policy{
		AllowedCapabilities: []string{"read_files", "execute_commands", "system_info"},
		BlockedCapabilities: []string{"network_access", "write_files"},
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
