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
	"bytes"
	"testing"
	"github.com/pterm/pterm"
	"github.com/stretchr/testify/assert"
)

func TestRoot_Execute(t *testing.T) {
	pterm.DisableStyling()
	defer pterm.EnableStyling()

	// Test --help
	cmd := NewRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--help"})
	
	err := cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, out.String(), "isolated nodes")

	// Test --version (Built-in Cobra flag when Version is set)
	cmd = NewRootCmd()
	out.Reset()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--version"})
	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, out.String(), "zt dev")

	// Test invalid command returns error (using unknown flag to trigger cobra error)
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.SetArgs([]string{"--invalid-flag"})
	err = rootCmd.Execute()
	assert.Error(t, err)
}

func TestExecute_Error(t *testing.T) {
	oldExit := osExit
	var exitCode int
	osExit = func(code int) {
		exitCode = code
	}
	defer func() { osExit = oldExit }()

	// Trigger error by setting an invalid flag on the global rootCmd
	rootCmd.SetArgs([]string{"--invalid-flag-for-execute"})
	Execute()
	assert.Equal(t, 1, exitCode)
}

func TestRoot_Subcommands(t *testing.T) {
	// We use the global rootCmd to test if subcommands are registered via init()
	assert.NotNil(t, rootCmd)
	
	// Check for 'scan' subcommand
	scan, _, _ := rootCmd.Find([]string{"scan"})
	assert.NotNil(t, scan)
	assert.Equal(t, "scan", scan.Name())

	// Check for 'plugin' subcommand
	plugin, _, _ := rootCmd.Find([]string{"plugin"})
	assert.NotNil(t, plugin)
	assert.Equal(t, "plugin", plugin.Name())
}
