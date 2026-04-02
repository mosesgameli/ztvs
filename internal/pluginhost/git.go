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

// Git handles git operations for plugin distribution.
type Git struct{}

func NewGit() *Git {
	return &Git{}
}

// Clone clones a git repository to the specified destination.
func (g *Git) Clone(url, dest string) error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git not found in PATH: %v", err)
	}

	cmd := exec.Command("git", "clone", url, dest)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git clone failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

// Pull updates an existing git repository in the specified destination.
func (g *Git) Pull(dest string) error {
	cmd := exec.Command("git", "-C", dest, "pull")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git pull failed: %v\nOutput: %s", err, string(output))
	}
	return nil
}

// Remove cleans up a plugin directory (e.g., on build failure).
func (g *Git) Remove(dest string) error {
	return os.RemoveAll(dest)
}
