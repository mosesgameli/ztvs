package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/internal/pluginhost"
)

func PluginCommand() {
	if len(os.Args) < 3 {
		fmt.Println("usage: zt plugin <subcommand>")
		fmt.Println("\nsubcommands:")
		fmt.Println("  init      Initialize ZTVS configuration (~/.ztvs/config.yaml)")
		fmt.Println("  list      List installed plugins")
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
		fmt.Println("✓ ZTVS initialized at ~/.ztvs/config.yaml")
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
				fmt.Printf("[%s] v%s (%s)\n", manifest.Name, manifest.Version, p)
			} else {
				fmt.Printf("[unknown] %s\n", p)
			}
		}
	case "install":
		if len(os.Args) < 4 {
			fmt.Println("usage: zt plugin install <name>")
			os.Exit(1)
		}
		registry := pluginhost.NewRegistry()
		if err := registry.Install(os.Args[3]); err != nil {
			fmt.Printf("installation error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("unknown plugin subcommand: %s\n", subcommand)
		os.Exit(1)
	}
}
