package engine

import (
	"context"

	"github.com/mosesgameli/ztvs/internal/pluginhost"
)

type Engine struct {
	host *pluginhost.Host
}

func New() *Engine {
	return &Engine{
		host: pluginhost.New(),
	}
}

func (e *Engine) Scan() error {
	ctx := context.Background()

	plugins, err := e.host.Discover(ctx)
	if err != nil {
		return err
	}

	for _, p := range plugins {
		_, err := e.host.RunCheck(ctx, p, "ssh_config")
		if err != nil {
			return err
		}
	}

	return nil
}
