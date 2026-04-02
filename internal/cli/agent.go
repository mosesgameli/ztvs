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
	"os/signal"
	"syscall"
	"time"

	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/internal/engine"
	"github.com/mosesgameli/ztvs/internal/report"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Start the background audit agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgent(cmd.Context())
	},
}

func runAgent(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	pterm.Success.Println("🚀 ZTVS Audit Agent starting...")

	r := report.NewTerminal()
	eng := engine.New(cfg, host, r, registryClient)
	eng.Interactive = false

	interval, err := time.ParseDuration(cfg.Agent.Interval)
	if err != nil {
		pterm.Warning.Printf("invalid agent interval %s: %v. falling back to 1h\n", cfg.Agent.Interval, err)
		interval = time.Hour
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		pterm.Warning.Println("\n🛑 Shutdown signal received. Exiting...")
		cancel()
	}()

	if err := eng.RunLoop(ctx, interval); err != nil && err != context.Canceled {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(agentCmd)
}
