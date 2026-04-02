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

package pluginhost

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGit_Clone(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()

	execCommand = func(command string, args ...string) *exec.Cmd {
		return exec.Command(os.Args[0], append([]string{"-test.run=TestGitHelperProcess", "--", command}, args...)...)
	}

	g := NewGit()
	err := g.Clone("http://example.com", t.TempDir())
	assert.NoError(t, err)
}

func TestGit_Pull(t *testing.T) {
	oldExec := execCommand
	defer func() { execCommand = oldExec }()

	execCommand = func(command string, args ...string) *exec.Cmd {
		return exec.Command(os.Args[0], append([]string{"-test.run=TestGitHelperProcess", "--", command}, args...)...)
	}

	g := NewGit()
	err := g.Pull(t.TempDir())
	assert.NoError(t, err)
}

func TestGit_Remove(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGit()
	err := g.Remove(tmpDir)
	assert.NoError(t, err)
	_, err = os.Stat(tmpDir)
	assert.True(t, os.IsNotExist(err))
}

func TestGitHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintf(os.Stdout, "git success")
	os.Exit(0)
}
