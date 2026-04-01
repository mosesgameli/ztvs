package cli

import (
	"fmt"
	"os"

	"github.com/mosesgameli/ztvs/internal/engine"
)

func Execute() {
	if len(os.Args) < 2 {
		fmt.Println("usage: zt <command>")
		fmt.Println("\ncommands:")
		fmt.Println("  scan    Run vulnerability scans")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "scan":
		eng := engine.New()
		if err := eng.Scan(); err != nil {
			fmt.Printf("scan failed: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
