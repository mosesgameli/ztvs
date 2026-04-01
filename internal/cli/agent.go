package cli

import (
	"context"
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
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			pterm.Error.Printf("config error: %v\n", err)
			os.Exit(1)
		}

		pterm.Success.Println("🚀 ZTVS Audit Agent starting...")

		r := report.NewTerminal()
		eng := engine.New(cfg, r)
		eng.Interactive = false

		interval, err := time.ParseDuration(cfg.Agent.Interval)
		if err != nil {
			pterm.Warning.Printf("invalid agent interval %s: %v. falling back to 1h\n", cfg.Agent.Interval, err)
			interval = time.Hour
		}

		ctx, cancel := context.WithCancel(context.Background())
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
			pterm.Error.Printf("agent error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
}
