package cli

import (
	"flag"
	"fmt"
	"os"

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
		fmt.Println("\nflags:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	command := flag.Arg(0)

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

		eng := engine.New(r)
		if err := eng.Scan(); err != nil {
			fmt.Printf("scan failed: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("unknown command: %s\n", command)
		os.Exit(1)
	}
}
