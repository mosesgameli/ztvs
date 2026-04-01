package pluginhost

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/pkg/registry"
	"github.com/mosesgameli/ztvs/pkg/rpc"
	"github.com/mosesgameli/ztvs/pkg/sdk"
	"gopkg.in/yaml.v3"
)

type PluginInfo struct {
	BinaryPath string
	Manifest   *sdk.Manifest
	Enabled    bool
}

type Host struct {
	paths    []string
	plugins  map[string]*PluginInfo
	lockfile *registry.Lockfile
}

func New() *Host {
	configDir := config.ConfigDir()
	lockPath := filepath.Join(configDir, "plugins.lock")
	lf, _ := registry.LoadLockfile(lockPath)

	return &Host{
		paths: []string{
			"./plugins",
			filepath.Join(configDir, "plugins"),
			"/usr/local/lib/zt/plugins",
		},
		plugins:  make(map[string]*PluginInfo),
		lockfile: lf,
	}
}

func (h *Host) Discover(ctx context.Context) ([]string, error) {
	var discovered []string

	for _, path := range h.paths {
		entries, err := os.ReadDir(path)
		if err != nil {
			continue // Skip missing directories
		}

		for _, entry := range entries {
			if entry.IsDir() {
				// 1. Located the plugin directory
				pluginDir := filepath.Join(path, entry.Name())

				// 2. Load manifest (Mandatory in Phase 1)
				manifestPath := filepath.Join(pluginDir, "plugin.yaml")
				manifest, err := h.loadManifest(manifestPath)
				if err != nil {
					// STRICT VALIDATION: Ignore plugins without valid manifest
					continue
				}

				// 3. Look for binary
				pluginBin := filepath.Join(pluginDir, entry.Name())
				if info, err := os.Stat(pluginBin); err == nil && !info.IsDir() {
					// 4. Check lockfile (Enabled status)
					enabled := true
					if lock, ok := h.lockfile.Get(manifest.Name); ok {
						enabled = lock.Enabled
					}

					h.plugins[pluginBin] = &PluginInfo{
						BinaryPath: pluginBin,
						Manifest:   manifest,
						Enabled:    enabled,
					}
					discovered = append(discovered, pluginBin)
				}
			}
		}
	}

	return discovered, nil
}

func (h *Host) loadManifest(path string) (*sdk.Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m sdk.Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse manifest %s: %v", path, err)
	}

	return &m, nil
}

func (h *Host) GetPluginInfo(pluginPath string) (*PluginInfo, bool) {
	info, ok := h.plugins[pluginPath]
	return info, ok
}

func (h *Host) GetManifest(pluginPath string) (*sdk.Manifest, bool) {
	info, ok := h.plugins[pluginPath]
	if !ok {
		return nil, false
	}
	return info.Manifest, true
}

func (h *Host) RunCheck(
	ctx context.Context,
	pluginPath string,
	checkID string,
) (*rpc.RunCheckResponse, error) {
	return h.runCheckProcess(ctx, pluginPath, checkID)
}

func (h *Host) Lockfile() *registry.Lockfile {
	return h.lockfile
}
