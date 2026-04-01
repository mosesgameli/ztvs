package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/internal/engine"
	"github.com/mosesgameli/ztvs/internal/report"
)

func Execute() {
	format := flag.String("format", "terminal", "Output format (terminal, json, sarif)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("usage: zt [flags] <command>")
		fmt.Println("\ncommands:")
		fmt.Println("  scan    Run vulnerability scans")
		fmt.Println("  plugin  Manage scanner plugins")
		fmt.Println("  agent   Start the background audit agent")
		fmt.Println("\nflags:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	command := flag.Arg(0)

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("config error: %v\n", err)
		os.Exit(1)
	}

	switch command {
	case "scan":
		var r report.Reporter
		switch *format {
		case "json":
			r = report.NewJSON()
		case "sarif":
			r = report.NewSARIF()
		default:
			r = report.NewTerminal()
		}

		eng := engine.New(cfg, r)
		if err := eng.Scan(); err != nil {
			fmt.Printf("scan failed: %v\n", err)
			os.Exit(1)
		}
	case "plugin":
		PluginCommand()
	case "agent":
		AgentCommand(cfg)
	default:
		fmt.Printf("unknown command: %s\n", command)
		os.Exit(1)
	}
}
