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

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
)

// Version is the current version of the application, injected during build.
var Version = "dev"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "zt",
		Version: Version,
		Short:   "Zero Trust Vulnerability Scanner",
		Args:    cobra.NoArgs,
		Long: func() string {
			logo, _ := pterm.DefaultBigText.WithLetters(
				putils.LettersFromStringWithStyle("ZT", pterm.NewStyle(pterm.FgCyan)),
				putils.LettersFromStringWithStyle("VS", pterm.NewStyle(pterm.FgMagenta))).
				Srender()
			return logo + "\n" + pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgMagenta)).WithMargin(10).Sprint("THE INDEPENDENT AUDIT PLATFORM") + "\n\n" +
				pterm.LightMagenta("ZTVS") + " is a high-fidelity, cross-platform security engine using isolated nodes for system auditing."
		}(),
	}
	cmd.SetVersionTemplate("{{.Name}} {{.Version}}\n")
	return cmd
}

var rootCmd = NewRootCmd()
var osExit = os.Exit

// Execute triggers the overarching command line execution flow
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		pterm.Error.Prefix = pterm.Prefix{Text: "ALERT", Style: pterm.NewStyle(pterm.BgRed, pterm.FgBlack)}
		pterm.Error.Println(err)
		osExit(1)
	}
}
