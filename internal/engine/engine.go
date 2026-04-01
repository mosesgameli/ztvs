package engine

import (
	"context"
	"fmt"
	"log"
	"time"

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
		return fmt.Errorf("discovery: %v", err)
	}

	for _, p := range plugins {
		// 1. Handshake
		meta, err := e.host.Handshake(ctx, p)
		if err != nil {
			log.Printf("Plugin %s failed handshake: %v", p, err)
			continue
		}

		if meta.APIVersion != 1 {
			log.Printf("Plugin %s has unsupported API version: %d", meta.Name, meta.APIVersion)
			continue
		}

		log.Printf("Running checks for plugin: %s (%s)", meta.Name, meta.Version)

		// 2. Run Checks
		for _, checkID := range meta.ChecksSupported {
			// Per-check timeout
			checkCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			res, err := e.host.RunCheck(checkCtx, p, checkID)
			if err != nil {
				log.Printf("  Check %s failed: %v", checkID, err)
				continue
			}

			if res.Finding != nil {
				fmt.Printf("[%s] %s: %s\n", res.Finding.Severity, res.Finding.Title, res.Finding.Description)
			}
		}
	}

	return nil
}
