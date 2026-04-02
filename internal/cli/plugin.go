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

package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mosesgameli/ztvs-sdk-go/sdk"
	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/internal/pluginhost"
	"github.com/mosesgameli/ztvs/pkg/registry"
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
		if err := runPluginInit(); err != nil {
			pterm.Error.Printf("initialization error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runPluginInit() error {
	cfg := config.DefaultConfig()
	if err := cfg.Save(); err != nil {
		return err
	}
	pterm.Success.Printf("ZTVS initialized at %s\n", config.ConfigDir())
	return nil
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed plugins",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPluginList(); err != nil {
			pterm.Error.Printf("list error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runPluginList() error {
	ctx := context.Background()
	plugins, err := host.Discover(ctx)
	if err != nil {
		return err
	}

	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgMagenta)).WithMargin(2).Printf("LOCAL NODE REGISTRY (%d)", len(plugins))
	fmt.Println()

	var tableData [][]string
	tableData = append(tableData, []string{"Node Name", "Version", "Status", "Filesystem Path"})

	for _, p := range plugins {
		if manifest, ok := host.GetManifest(p); ok {
			status := pterm.FgGreen.Sprint("● ACTIVE")
			if lock, ok := host.Lockfile().Get(manifest.Name); ok {
				if !lock.Enabled {
					status = pterm.FgYellow.Sprint("○ DISABLED")
				}
			}
			tableData = append(tableData, []string{
				pterm.FgCyan.Sprint(manifest.Name),
				pterm.LightWhite(manifest.Version),
				status,
				pterm.FgGray.Sprint(p),
			})
		} else {
			tableData = append(tableData, []string{pterm.FgRed.Sprint("unknown"), "", "", p})
		}
	}

	table, _ := pterm.DefaultTable.WithHasHeader().WithData(tableData).Srender()
	_ = pterm.DefaultPanel.WithPadding(1).WithPanels(pterm.Panels{
		{{Data: table}},
	}).Render()
	return nil
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
		if err := runPluginSearch(pluginhost.NewRegistry(), query); err != nil {
			pterm.Error.Printf("search error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runPluginSearch(reg pluginhost.Registry, query string) error {
	ctx := context.Background()
	results, err := reg.Search(ctx, query)
	if err != nil {
		return err
	}

	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).WithMargin(2).Printf("REMOTE CATALOG SEARCH (%d)", len(results))
	fmt.Println()

	var tableData [][]string
	tableData = append(tableData, []string{"Plugin", "Version", "Audit Level"})
	for _, r := range results {
		auditLevel := pterm.BgBlue.Sprint(pterm.FgBlack.Sprint(" " + r.AuditStatus + " "))
		tableData = append(tableData, []string{pterm.FgCyan.Sprint(r.Name), r.LatestVersion, auditLevel})
	}
	table, _ := pterm.DefaultTable.WithHasHeader().WithData(tableData).Srender()
	_ = pterm.DefaultPanel.WithPadding(1).WithPanels(pterm.Panels{
		{{Data: table}},
	}).Render()
	return nil
}

var pluginInfoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Show remote plugin details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPluginInfo(pluginhost.NewRegistry(), args[0]); err != nil {
			pterm.Error.Printf("info error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runPluginInfo(reg pluginhost.Registry, name string) error {
	ctx := context.Background()
	meta, err := reg.GetInfo(ctx, name)
	if err != nil {
		return err
	}

	pterm.DefaultSection.Printf("Node Manifest: %s", meta.Name)

	infoData := [][]string{
		{"Field", "Value"},
		{"Version", pterm.FgMagenta.Sprint(meta.LatestVersion)},
		{"Remote Repository", pterm.FgCyan.Sprint(meta.Repo)},
		{"Security Audit", pterm.BgGreen.Sprint(pterm.FgBlack.Sprint(" " + meta.AuditStatus + " "))},
		{"Integrity Hash", pterm.FgGray.Sprint(meta.Checksum)},
	}
	_ = pterm.DefaultTable.WithHasHeader().WithData(infoData).Render()
	return nil
}

var pluginEnableCmd = &cobra.Command{
	Use:   "enable <name>",
	Short: "Enable a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPluginToggle(args[0], true); err != nil {
			pterm.Error.Printf("enable error: %v\n", err)
			os.Exit(1)
		}
	},
}

var pluginDisableCmd = &cobra.Command{
	Use:   "disable <name>",
	Short: "Disable a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPluginToggle(args[0], false); err != nil {
			pterm.Error.Printf("disable error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runPluginToggle(name string, enabled bool) error {
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
			return fmt.Errorf("plugin %s not found", name)
		}
	} else {
		lock.Enabled = enabled
	}

	lf.Set(name, lock)
	if err := lf.Save(); err != nil {
		return fmt.Errorf("error saving lockfile: %w", err)
	}

	action := "disabled"
	if enabled {
		action = "enabled"
	}
	pterm.Success.Printf("Plugin %s %s\n", name, action)
	return nil
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install <name>",
	Short: "Install a remote plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPluginInstall(pluginhost.NewRegistry(), args[0]); err != nil {
			pterm.Error.Printf("install error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runPluginInstall(reg pluginhost.Registry, name string) error {
	ctx := context.Background()

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Installing plugin '%s'...", name))
	if err := reg.Install(ctx, name, host); err != nil {
		spinner.Fail(fmt.Sprintf("installation error: %v", err))
		return err
	}
	spinner.Success(fmt.Sprintf("Successfully installed plugin '%s'", name))

	configDir := config.ConfigDir()
	// #nosec G304 - name is a plugin name, and configDir is a trusted root.
	manifestPath := filepath.Join(configDir, "plugins", filepath.Clean(name), "plugin.yaml")
	data, err := os.ReadFile(manifestPath)
	if err == nil {
		var m sdk.Manifest
		if err := yaml.Unmarshal(data, &m); err == nil {
			pterm.Info.Printf("Plugin '%s' requires the following capabilities:\n", m.Name)
			var items []pterm.BulletListItem
			for _, cap := range m.Capabilities {
				items = append(items, pterm.BulletListItem{Level: 0, Text: cap})
			}
			_ = pterm.DefaultBulletList.WithItems(items).Render()
			pterm.Println()
			pterm.Warning.Println("These capabilities are recorded and will be enforced by the ZTVS policy engine during scans.")
		}
	}
	return nil
}

var pluginUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update installed plugins to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runPluginUpdate(pluginhost.NewRegistry()); err != nil {
			pterm.Error.Printf("update error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runPluginUpdate(reg pluginhost.Registry) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	ctx := context.Background()

	if err := reg.CheckAndUpdateAll(ctx, host, cfg.Update.Mode); err != nil {
		return err
	}
	pterm.Success.Println("All plugins are up to date.")
	return nil
}

func init() {
	pluginCmd.AddCommand(pluginInitCmd, pluginListCmd, pluginSearchCmd, pluginInfoCmd, pluginEnableCmd, pluginDisableCmd, pluginInstallCmd, pluginUpdateCmd)
	rootCmd.AddCommand(pluginCmd)
}
