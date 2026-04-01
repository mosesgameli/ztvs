package cli

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
)

// Version is the current version of the application, injected during build.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "zt",
	Version: Version,
	Short:   "Zero Trust Vulnerability Scanner",
	Long: func() string {
		logo, _ := pterm.DefaultBigText.WithLetters(
			putils.LettersFromStringWithStyle("ZT", pterm.NewStyle(pterm.FgCyan)),
			putils.LettersFromStringWithStyle("VS", pterm.NewStyle(pterm.FgMagenta))).
			Srender()
		return logo + "\n" + pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgMagenta)).WithMargin(10).Sprint("THE INDEPENDENT AUDIT PLATFORM") + "\n\n" +
			pterm.LightMagenta("ZTVS") + " is a high-fidelity, cross-platform security engine using isolated nodes for system auditing."
	}(),
}

// Execute triggers the overarching command line execution flow
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		pterm.Error.Prefix = pterm.Prefix{Text: "ALERT", Style: pterm.NewStyle(pterm.BgRed, pterm.FgBlack)}
		pterm.Error.Println(err)
		os.Exit(1)
	}
}
