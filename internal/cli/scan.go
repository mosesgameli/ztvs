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
	"os"

	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/internal/engine"
	"github.com/mosesgameli/ztvs/internal/report"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var formatFlag string

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run vulnerability scans across all enabled plugins",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			pterm.Error.Printf("config error: %v\n", err)
			os.Exit(1)
		}

		var r report.Reporter
		switch formatFlag {
		case "json":
			r = report.NewJSON()
		case "sarif":
			r = report.NewSARIF()
		default:
			r = report.NewTerminal()
		}

		eng := engine.New(cfg, r)
		eng.Interactive = true

		pterm.Info.Println("Initializing global security scan...")

		if err := eng.Scan(); err != nil {
			pterm.Error.Printf("Scan failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	scanCmd.Flags().StringVar(&formatFlag, "format", "terminal", "Output format (terminal, json, sarif)")
	rootCmd.AddCommand(scanCmd)
}
