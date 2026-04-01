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
