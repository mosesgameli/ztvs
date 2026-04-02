// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package engine

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/internal/pluginhost"
	"github.com/mosesgameli/ztvs/internal/policy"
	"github.com/mosesgameli/ztvs/internal/report"
	"github.com/pterm/pterm"
)

type Engine struct {
	host        pluginhost.PluginHost
	reporter    report.Reporter
	policy      *policy.Policy
	cfg         *config.Config
	mutex       sync.Mutex
	Interactive bool
	stdin       io.Reader
	registry    pluginhost.Registry
}

func New(cfg *config.Config, h pluginhost.PluginHost, r report.Reporter, reg pluginhost.Registry) *Engine {
	return &Engine{
		host:     h,
		reporter: r,
		policy:   policy.New(cfg),
		cfg:      cfg,
		stdin:    os.Stdin,
		registry: reg,
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

	var spinner *pterm.SpinnerPrinter
	if e.Interactive {
		spinner, _ = pterm.DefaultSpinner.Start("Checking for registry updates...")
	}

	// Phase 4: Atomic Auto-Updates
	reg := e.registry
	if err := reg.CheckAndUpdateAll(ctx, e.host, e.cfg.Update.Mode); err != nil {
		if spinner != nil {
			pterm.Warning.Printf("Update check bypassed: %v\n", err)
		}
	}

	// 1. Discover plugins
	if spinner != nil {
		spinner.UpdateText("Discovering installed nodes...")
	}
	plugins, err := e.host.Discover(ctx)
	if err != nil {
		if spinner != nil {
			spinner.Fail(fmt.Sprintf("Discovery error: %v", err))
		}
		return fmt.Errorf("discovery: %v", err)
	}

	// Run pre-flight checks sequentially to present prompts cleanly
	plugins = e.preflightCapabilityCheck(plugins)

	var wg sync.WaitGroup

	if spinner != nil {
		spinner.UpdateText(fmt.Sprintf("Auditing system with %d nodes...", len(plugins)))
	}

	for _, p := range plugins {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			if spinner != nil {
				if manifest, ok := e.host.GetManifest(path); ok {
					spinner.UpdateText(fmt.Sprintf("Node [%s]: Executing checks...", manifest.Name))
				}
			}
			e.scanPlugin(ctx, path)
		}(p)
	}

	wg.Wait()

	if spinner != nil {
		spinner.Success("System audit finalized")
	}

	return e.reporter.Flush()
}

func (e *Engine) preflightCapabilityCheck(plugins []string) []string {
	var allowed []string
	for _, p := range plugins {
		info, _ := e.host.GetPluginInfo(p)
		if info != nil && !info.Enabled {
			continue
		}

		manifest, ok := e.host.GetManifest(p)
		if !ok {
			log.Printf("Security alert: Plugin at %s has no manifest. Skipping.", p)
			continue
		}

		if cap, err := e.policy.IsAllowed(manifest.Name, manifest.Capabilities); err != nil {
			if e.Interactive {
				fmt.Printf("\n[SEC] Plugin '%s' requested blocked/unauthorized capability: %s\n", manifest.Name, cap)
				fmt.Printf("      Do you want to permanently grant this capability in your global config? [y/N]: ")

				reader := bufio.NewReader(e.stdin)
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))

				if response == "y" || response == "yes" {
					// 1. Remove from blocked
					var newBlocked []string
					for _, b := range e.cfg.Policy.BlockedCapabilities {
						if b != cap {
							newBlocked = append(newBlocked, b)
						}
					}
					e.cfg.Policy.BlockedCapabilities = newBlocked

					// 2. Add to allowed
					e.cfg.Policy.AllowedCapabilities = append(e.cfg.Policy.AllowedCapabilities, cap)

					// 3. Save to disk
					if saveErr := e.cfg.Save(); saveErr != nil {
						log.Printf("Failed to save config: %v", saveErr)
						continue
					}

					// 4. Reload in-memory
					e.policy.Reload(e.cfg)
					fmt.Printf("      Granted %s capability.\n\n", cap)

					allowed = append(allowed, p)
					continue
				}
			}
			log.Printf("Policy rejection: %v. Skipping plugin %s.", err, manifest.Name)
			continue
		}

		allowed = append(allowed, p)
	}
	return allowed
}

func (e *Engine) scanPlugin(ctx context.Context, p string) {
	manifest, _ := e.host.GetManifest(p)

	// Final verification in case of race constraints
	if _, err := e.policy.IsAllowed(manifest.Name, manifest.Capabilities); err != nil {
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
