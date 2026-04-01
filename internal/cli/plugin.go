package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/internal/pluginhost"
	"github.com/mosesgameli/ztvs/pkg/registry"
	"github.com/mosesgameli/ztvs/pkg/sdk"
	"gopkg.in/yaml.v3"
)

func PluginCommand() {
	if len(os.Args) < 3 {
		fmt.Println("usage: zt plugin <subcommand>")
		fmt.Println("\nsubcommands:")
		fmt.Println("  init      Initialize ZTVS configuration (~/.ztvs/config.yaml)")
		fmt.Println("  list      List installed plugins")
		fmt.Println("  search    Search for remote plugins")
		fmt.Println("  info      Show remote plugin details")
		fmt.Println("  enable    Enable a plugin")
		fmt.Println("  disable   Disable a plugin")
		fmt.Println("  install   Install a remote plugin (simulation)")
		os.Exit(1)
	}

	subcommand := os.Args[2]
	host := pluginhost.New()

	switch subcommand {
	case "init":
		cfg := config.DefaultConfig()
		if err := cfg.Save(); err != nil {
			fmt.Printf("initialization error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ ZTVS initialized at %s\n", config.ConfigDir())
	case "list":
		ctx := context.Background()
		plugins, err := host.Discover(ctx)
		if err != nil {
			fmt.Printf("discovery error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("--- Installed Plugins (%d) ---\n", len(plugins))
		for _, p := range plugins {
			if manifest, ok := host.GetManifest(p); ok {
				status := "Enabled"
				if lock, ok := host.Lockfile().Get(manifest.Name); ok {
					if !lock.Enabled {
						status = "Disabled"
					}
				}
				fmt.Printf("[%s] v%s (%s) - %s\n", manifest.Name, manifest.Version, p, status)
			} else {
				fmt.Printf("[unknown] %s\n", p)
			}
		}
	case "search":
		query := ""
		if len(os.Args) >= 4 {
			query = os.Args[3]
		}
		registry := pluginhost.NewRegistry()
		ctx := context.Background()
		results, err := registry.Search(ctx, query)
		if err != nil {
			fmt.Printf("search error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("--- Found Plugins (%d) ---\n", len(results))
		for _, r := range results {
			fmt.Printf("[%s] v%s - %s\n", r.Name, r.LatestVersion, r.AuditStatus)
		}
	case "info":
		if len(os.Args) < 4 {
			fmt.Println("usage: zt plugin info <name>")
			os.Exit(1)
		}
		name := os.Args[3]
		registry := pluginhost.NewRegistry()
		ctx := context.Background()
		meta, err := registry.GetInfo(ctx, name)
		if err != nil {
			fmt.Printf("info error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Plugin: %s\n", meta.Name)
		fmt.Printf("Version: %s\n", meta.LatestVersion)
		fmt.Printf("Repository: %s\n", meta.Repo)
		fmt.Printf("Audit Status: %s\n", meta.AuditStatus)
		fmt.Printf("Checksum: %s\n", meta.Checksum)
	case "enable", "disable":
		if len(os.Args) < 4 {
			fmt.Printf("usage: zt plugin %s <name>\n", subcommand)
			os.Exit(1)
		}
		name := os.Args[3]
		enabled := subcommand == "enable"

		lf := host.Lockfile()
		lock, ok := lf.Get(name)
		if !ok {
			// Try to find it first to get the version
			ctx := context.Background()
			plugins, _ := host.Discover(ctx)
			found := false
			for _, p := range plugins {
				if m, ok := host.GetManifest(p); ok && m.Name == name {
					lock = registry.PluginLock{
						Version: m.Version,
						Enabled: enabled,
					}
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("Error: plugin %s not found\n", name)
				os.Exit(1)
			}
		} else {
			lock.Enabled = enabled
		}

		lf.Set(name, lock)
		if err := lf.Save(); err != nil {
			fmt.Printf("Error saving lockfile: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✓ Plugin %s %sd\n", name, subcommand)
	case "install":
		if len(os.Args) < 4 {
			fmt.Println("usage: zt plugin install <name>")
			os.Exit(1)
		}
		name := os.Args[3]
		registryClient := pluginhost.NewRegistry()
		ctx := context.Background()

		// 1. Get info to show capabilities before install (if possible)
		// For now, we'll clone first then show capabilities from the manifest on disk.
		// In a real registry, the index.json might include capabilities.

		fmt.Printf("Attempting to install plugin: %s\n", name)

		// 2. Install (which clones and builds)
		if err := registryClient.Install(ctx, name, host); err != nil {
			fmt.Printf("installation error: %v\n", err)
			os.Exit(1)
		}

		// 3. Post-install capability check
		// (Normally we'd do this BEFORE build/install, but for Phase 3 we'll do it post-clone)
		configDir := config.ConfigDir()
		manifestPath := filepath.Join(configDir, "plugins", name, "plugin.yaml")
		data, err := os.ReadFile(manifestPath)
		if err == nil {
			var m sdk.Manifest
			if err := yaml.Unmarshal(data, &m); err == nil {
				fmt.Printf("\nPlugin '%s' requires the following capabilities:\n", m.Name)
				for _, cap := range m.Capabilities {
					fmt.Printf("  - %s\n", cap)
				}
				fmt.Printf("\n[NOTE] These capabilities are recorded and will be enforced by the ZTVS policy engine during scans.\n")
			}
		}
	default:
		fmt.Printf("unknown plugin subcommand: %s\n", subcommand)
		os.Exit(1)
	}
}
