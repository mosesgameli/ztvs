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
)

var execCommand = exec.Command

// Git defines the interface for git operations.
type Git interface {
	Clone(url, dest string) error
	Pull(dest string) error
	Remove(dest string) error
}

type gitImpl struct{}

func NewGit() Git {
	return &gitImpl{}
}

func (g *gitImpl) Clone(url, dest string) error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git not found in PATH: %v", err)
	}

	cmd := execCommand("git", "clone", url, dest)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

func (g *gitImpl) Pull(dest string) error {
	cmd := execCommand("git", "-C", dest, "pull")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git pull failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

func (g *gitImpl) Remove(dest string) error {
	return os.RemoveAll(dest)
}
