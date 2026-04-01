package cli

import (
	"context"
	"os"
	"path/filepath"
	"fmt"

	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/internal/pluginhost"
	"github.com/mosesgameli/ztvs/pkg/registry"
	"github.com/mosesgameli/ztvs/pkg/sdk"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var host = pluginhost.New()

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage scanner plugins",
}

var pluginInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize ZTVS configuration (~/.ztvs/config.yaml)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.DefaultConfig()
		if err := cfg.Save(); err != nil {
			pterm.Error.Printf("initialization error: %v\n", err)
			os.Exit(1)
		}
		pterm.Success.Printf("ZTVS initialized at %s\n", config.ConfigDir())
	},
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed plugins",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		plugins, err := host.Discover(ctx)
		if err != nil {
			pterm.Error.Printf("discovery error: %v\n", err)
			os.Exit(1)
		}

		pterm.DefaultHeader.WithFullWidth().Printf("Installed Plugins (%d)", len(plugins))

		var tableData [][]string
		tableData = append(tableData, []string{"Name", "Version", "Path", "Status"})

		for _, p := range plugins {
			if manifest, ok := host.GetManifest(p); ok {
				status := "Enabled"
				if lock, ok := host.Lockfile().Get(manifest.Name); ok {
					if !lock.Enabled {
						status = pterm.FgYellow.Sprint("Disabled")
					} else {
						status = pterm.FgGreen.Sprint("Enabled")
					}
				}
				tableData = append(tableData, []string{
					pterm.FgCyan.Sprint(manifest.Name),
					manifest.Version,
					p,
					status,
				})
			} else {
				tableData = append(tableData, []string{pterm.FgRed.Sprint("unknown"), "", p, ""})
			}
		}

		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	},
}

var pluginSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for remote plugins",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := ""
		if len(args) > 0 {
			query = args[0]
		}
		registry := pluginhost.NewRegistry()
		ctx := context.Background()
		results, err := registry.Search(ctx, query)
		if err != nil {
			pterm.Error.Printf("search error: %v\n", err)
			os.Exit(1)
		}

		pterm.DefaultHeader.WithFullWidth().Printf("Found Plugins (%d)", len(results))

		var tableData [][]string
		tableData = append(tableData, []string{"Name", "Latest Version", "Status"})
		for _, r := range results {
			tableData = append(tableData, []string{pterm.FgCyan.Sprint(r.Name), r.LatestVersion, r.AuditStatus})
		}
		pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	},
}

var pluginInfoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show remote plugin details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		registry := pluginhost.NewRegistry()
		ctx := context.Background()
		meta, err := registry.GetInfo(ctx, name)
		if err != nil {
			pterm.Error.Printf("info error: %v\n", err)
			os.Exit(1)
		}

		pterm.DefaultHeader.Printf("Plugin: %s", meta.Name)
		pterm.Println(pterm.LightBlue("Version: ") + meta.LatestVersion)
		pterm.Println(pterm.LightBlue("Repository: ") + meta.Repo)
		pterm.Println(pterm.LightBlue("Audit Status: ") + meta.AuditStatus)
		pterm.Println(pterm.LightBlue("Checksum: ") + meta.Checksum)
	},
}

var pluginEnableCmd = &cobra.Command{
	Use:   "enable <name>",
	Short: "Enable a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		togglePlugin(args[0], true)
	},
}

var pluginDisableCmd = &cobra.Command{
	Use:   "disable <name>",
	Short: "Disable a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		togglePlugin(args[0], false)
	},
}

func togglePlugin(name string, enabled bool) {
	lf := host.Lockfile()
	lock, ok := lf.Get(name)
	if !ok {
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
			pterm.Error.Printf("plugin %s not found\n", name)
			os.Exit(1)
		}
	} else {
		lock.Enabled = enabled
	}

	lf.Set(name, lock)
	if err := lf.Save(); err != nil {
		pterm.Error.Printf("Error saving lockfile: %v\n", err)
		os.Exit(1)
	}
	
	action := "disabled"
	if enabled {
		action = "enabled"
	}
	pterm.Success.Printf("Plugin %s %s\n", name, action)
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install <name>",
	Short: "Install a remote plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		registryClient := pluginhost.NewRegistry()
		ctx := context.Background()

		spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing plugin '%s'...", name))
		if err := registryClient.Install(ctx, name, host); err != nil {
			spinner.Fail(fmt.Sprintf("installation error: %v", err))
			os.Exit(1)
		}
		spinner.Success(fmt.Sprintf("Successfully installed plugin '%s'", name))

		configDir := config.ConfigDir()
		manifestPath := filepath.Join(configDir, "plugins", name, "plugin.yaml")
		data, err := os.ReadFile(manifestPath)
		if err == nil {
			var m sdk.Manifest
			if err := yaml.Unmarshal(data, &m); err == nil {
				pterm.Info.Printf("Plugin '%s' requires the following capabilities:\n", m.Name)
				var items []pterm.BulletListItem
				for _, cap := range m.Capabilities {
					items = append(items, pterm.BulletListItem{Level: 0, Text: cap})
				}
				pterm.DefaultBulletList.WithItems(items).Render()
				pterm.Println()
				pterm.Warning.Println("These capabilities are recorded and will be enforced by the ZTVS policy engine during scans.")
			}
		}
	},
}

func init() {
	pluginCmd.AddCommand(pluginInitCmd, pluginListCmd, pluginSearchCmd, pluginInfoCmd, pluginEnableCmd, pluginDisableCmd, pluginInstallCmd)
	rootCmd.AddCommand(pluginCmd)
}
