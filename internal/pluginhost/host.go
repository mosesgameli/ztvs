package pluginhost

import (
	"context"

	"github.com/mosesgameli/ztvs/pkg/rpc"
)

type Host struct{}

func New() *Host {
	return &Host{}
}

func (h *Host) Discover(ctx context.Context) ([]string, error) {
	// For MVP, look in relative plugins directory or simulate discovery
	return []string{
		"./plugins/plugin-os/plugin-os",
	}, nil
}

func (h *Host) RunCheck(
	ctx context.Context,
	pluginPath string,
	checkID string,
) (*rpc.RunCheckResponse, error) {
	return h.runCheckProcess(ctx, pluginPath, checkID)
}
