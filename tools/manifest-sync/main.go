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

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: manifest-sync <plugin-dir> [bin-ext]")
	}

	pluginDir := args[1]
	binExt := ""
	if len(args) > 2 {
		binExt = args[2]
	}
	pluginName := filepath.Base(pluginDir)
	binPath := filepath.Join(pluginDir, pluginName+binExt)
	manifestPath := filepath.Join(pluginDir, "plugin.yaml")

	// 1. Check if binary exists
	f, err := os.Open(binPath)
	if err != nil {
		return fmt.Errorf("binary not found at %s. Build it first", binPath)
	}
	defer f.Close()

	// 2. Compute SHA-256
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("computing hash: %v", err)
	}
	checksum := hex.EncodeToString(h.Sum(nil))

	// 3. Update manifest
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("reading manifest: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, "checksum:") {
			lines[i] = fmt.Sprintf("checksum: %s", checksum)
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, fmt.Sprintf("checksum: %s", checksum))
	}

	if err := os.WriteFile(manifestPath, []byte(strings.Join(lines, "\n")), 0600); err != nil {
		return fmt.Errorf("writing manifest: %v", err)
	}

	fmt.Printf("Updated %s with checksum %s\n", manifestPath, checksum)
	return nil
}
