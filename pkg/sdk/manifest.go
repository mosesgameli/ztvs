package sdk

// Manifest represents the static metadata and security requirements for a plugin.
// This is usually stored in a plugin.yaml file alongside the plugin binary.
type Manifest struct {
	Name            string   `yaml:"name"`
	Version         string   `yaml:"version"`
	APIVersion      int      `yaml:"api_version"`
	Capabilities    []string `yaml:"capabilities"`
	Checksum        string   `yaml:"checksum,omitempty"` // sha256 hex string
	ChecksSupported []string `yaml:"checks_supported,omitempty"`
}
