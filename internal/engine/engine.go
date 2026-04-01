package engine

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/internal/pluginhost"
	"github.com/mosesgameli/ztvs/internal/policy"
	"github.com/mosesgameli/ztvs/internal/report"
)

type Engine struct {
	host     *pluginhost.Host
	reporter report.Reporter
	policy   *policy.Policy
	mutex    sync.Mutex
}

func New(cfg *config.Config, r report.Reporter) *Engine {
	return &Engine{
		host:     pluginhost.New(),
		reporter: r,
		policy:   policy.New(cfg),
	}
}

func (e *Engine) RunLoop(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Starting background audit agent (interval: %v)", interval)

	// First run
	if err := e.Scan(); err != nil {
		log.Printf("Initial scan error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			log.Printf("Periodic audit started...")
			if err := e.Scan(); err != nil {
				log.Printf("Audit scan error: %v", err)
			}
		}
	}
}

func (e *Engine) Scan() error {
	ctx := context.Background()

	plugins, err := e.host.Discover(ctx)
	if err != nil {
		return fmt.Errorf("discovery: %v", err)
	}

	var wg sync.WaitGroup

	for _, p := range plugins {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			e.scanPlugin(ctx, path)
		}(p)
	}

	wg.Wait()
	return e.reporter.Flush()
}

func (e *Engine) scanPlugin(ctx context.Context, p string) {
	// 1. Get Manifest & Enforce Policy (Phase 3)
	manifest, ok := e.host.GetManifest(p)
	if !ok {
		log.Printf("Security alert: Plugin at %s has no manifest. Skipping.", p)
		return
	}

	if err := e.policy.IsAllowed(manifest.Name, manifest.Capabilities); err != nil {
		log.Printf("Policy rejection: %v. Skipping plugin %s.", err, manifest.Name)
		return
	}

	// 2. Handshake (Active verification)
	meta, err := e.host.Handshake(ctx, p)
	if err != nil {
		log.Printf("Plugin %s failed handshake: %v", p, err)
		return
	}

	if meta.APIVersion != 1 {
		log.Printf("Plugin %s has unsupported API version: %d", meta.Name, meta.APIVersion)
		return
	}

	// 3. Run Checks
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
			e.mutex.Lock()
			e.reporter.AddFinding(meta.Name, res.Finding)
			e.mutex.Unlock()
		}
	}
}
