package pluginhost

import (
	"context"
	"os"
	"path/filepath"

	"github.com/mosesgameli/ztvs/pkg/rpc"
)

type Host struct {
	paths []string
}

func New() *Host {
	home, _ := os.UserHomeDir()
	return &Host{
		paths: []string{
			"./plugins",
			filepath.Join(home, ".zt", "plugins"),
			"/usr/local/lib/zt/plugins",
		},
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
				// Each plugin should be in its own directory
				// Look for a binary with the same name as the directory
				pluginBin := filepath.Join(path, entry.Name(), entry.Name())
				if info, err := os.Stat(pluginBin); err == nil && !info.IsDir() {
					discovered = append(discovered, pluginBin)
				}
			}
		}
	}

	return discovered, nil
}

func (h *Host) RunCheck(
	ctx context.Context,
	pluginPath string,
	checkID string,
) (*rpc.RunCheckResponse, error) {
	return h.runCheckProcess(ctx, pluginPath, checkID)
}
