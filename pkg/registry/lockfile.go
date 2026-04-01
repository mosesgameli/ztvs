package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type PluginLock struct {
	Version  string `yaml:"version"`
	Enabled  bool   `yaml:"enabled"`
	Checksum string `yaml:"checksum,omitempty"`
}

type Lockfile struct {
	Version string                `yaml:"version"`
	Plugins map[string]PluginLock `yaml:"plugins"`
	path    string
	mu      sync.RWMutex
}

func NewLockfile(path string) *Lockfile {
	return &Lockfile{
		Version: "1.0",
		Plugins: make(map[string]PluginLock),
		path:    path,
	}
}

func LoadLockfile(path string) (*Lockfile, error) {
	lf := NewLockfile(path)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return lf, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, lf); err != nil {
		return nil, fmt.Errorf("parse lockfile %s: %v", path, err)
	}

	if lf.Plugins == nil {
		lf.Plugins = make(map[string]PluginLock)
	}

	return lf, nil
}

func (lf *Lockfile) Save() error {
	lf.mu.RLock()
	defer lf.mu.RUnlock()

	data, err := yaml.Marshal(lf)
	if err != nil {
		return err
	}

	dir := filepath.Dir(lf.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(lf.path, data, 0644)
}

func (lf *Lockfile) Get(name string) (PluginLock, bool) {
	lf.mu.RLock()
	defer lf.mu.RUnlock()
	lock, ok := lf.Plugins[name]
	return lock, ok
}

func (l *Lockfile) Set(name string, lock PluginLock) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Plugins[name] = lock
}

func (l *Lockfile) All() map[string]PluginLock {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Return a copy to prevent external modification
	tmp := make(map[string]PluginLock)
	for k, v := range l.Plugins {
		tmp[k] = v
	}
	return tmp
}

func (lf *Lockfile) Remove(name string) {
	lf.mu.Lock()
	defer lf.mu.Unlock()
	delete(lf.Plugins, name)
}
