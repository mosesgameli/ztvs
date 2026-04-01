package cli

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zt",
	Short: "Zero Trust Vulnerability Scanner",
	Long:  pterm.LightMagenta("ZTVS") + " is a highly isolated, cross-platform vulnerability scanning platform based on independent plugins.",
}

// Execute triggers the overarching command line execution flow
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}
