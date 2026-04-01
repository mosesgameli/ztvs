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
)

func AgentCommand(cfg *config.Config) {
	fmt.Println("🚀 ZTVS Audit Agent starting...")

	r := report.NewTerminal()
	eng := engine.New(cfg, r)
	eng.Interactive = false

	interval, err := time.ParseDuration(cfg.Agent.Interval)
	if err != nil {
		fmt.Printf("invalid agent interval %s: %v. falling back to 1h\n", cfg.Agent.Interval, err)
		interval = time.Hour
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\n🛑 Shutdown signal received. Exiting...")
		cancel()
	}()

	if err := eng.RunLoop(ctx, interval); err != nil && err != context.Canceled {
		fmt.Printf("agent error: %v\n", err)
		os.Exit(1)
	}
}
