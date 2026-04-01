package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Agent  AgentConfig  `yaml:"agent"`
	Policy PolicyConfig `yaml:"policy"`
	Update UpdateConfig `yaml:"update"`
}

type UpdateConfig struct {
	Mode string `yaml:"mode"` // safe, always, locked
}

type AgentConfig struct {
	Interval string `yaml:"interval"` // e.g., 1h
}

type PolicyConfig struct {
	AllowedCapabilities []string `yaml:"allowed_capabilities"`
	BlockedCapabilities []string `yaml:"blocked_capabilities"`
}

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ztvs")
}

func Load() (*Config, error) {
	configPath := filepath.Join(ConfigDir(), "config.yaml")

	// 1. Create default config if missing
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// 2. Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %v", configPath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %v", configPath, err)
	}

	return &cfg, nil
}

func DefaultConfig() *Config {
	return &Config{
		Agent: AgentConfig{
			Interval: "1h",
		},
		Policy: PolicyConfig{
			AllowedCapabilities: []string{"read_files", "execute_commands", "system_info"},
			BlockedCapabilities: []string{"network_access", "write_files"},
		},
		Update: UpdateConfig{
			Mode: "safe",
		},
	}
}

func (c *Config) Save() error {
	configDir := ConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.yaml")
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
