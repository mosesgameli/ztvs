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
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	tmpDir := t.TempDir()
	
	t.Run("Usage_Error", func(t *testing.T) {
		err := run([]string{"manifest-sync"})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "usage:")
	})

	t.Run("Binary_Not_Found", func(t *testing.T) {
		err := run([]string{"manifest-sync", tmpDir})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "binary not found")
	})

	t.Run("Success", func(t *testing.T) {
		pluginDir := filepath.Join(tmpDir, "my-plugin")
		os.MkdirAll(pluginDir, 0755)
		
		binPath := filepath.Join(pluginDir, "my-plugin")
		os.WriteFile(binPath, []byte("binary-content"), 0755)
		
		manifestPath := filepath.Join(pluginDir, "plugin.yaml")
		os.WriteFile(manifestPath, []byte("name: my-plugin\nversion: 1.0.0"), 0644)
		
		err := run([]string{"manifest-sync", pluginDir})
		assert.NoError(t, err)
		
		// Verify checksum was added
		content, err := os.ReadFile(manifestPath)
		require.NoError(t, err)
		assert.Contains(t, string(content), "checksum:")
	})

	t.Run("Update_Existing_Checksum", func(t *testing.T) {
		pluginDir := filepath.Join(tmpDir, "existing-plugin")
		os.MkdirAll(pluginDir, 0755)
		
		binPath := filepath.Join(pluginDir, "existing-plugin")
		os.WriteFile(binPath, []byte("new-content"), 0755)
		
		manifestPath := filepath.Join(pluginDir, "plugin.yaml")
		os.WriteFile(manifestPath, []byte("name: existing\nchecksum: old-hash"), 0644)
		
		err := run([]string{"manifest-sync", pluginDir})
		assert.NoError(t, err)
		
		content, err := os.ReadFile(manifestPath)
		require.NoError(t, err)
		assert.NotContains(t, string(content), "old-hash")
		assert.Contains(t, string(content), "checksum:")
	})

	t.Run("Manifest_Read_Error", func(t *testing.T) {
		pluginDir := filepath.Join(tmpDir, "no-manifest")
		os.MkdirAll(pluginDir, 0755)
		os.WriteFile(filepath.Join(pluginDir, "no-manifest"), []byte("bin"), 0755)
		
		err := run([]string{"manifest-sync", pluginDir})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reading manifest")
	})
}

func TestMainCall(t *testing.T) {
	// Setup a valid environment
	tmpDir := t.TempDir()
	pluginDir := filepath.Join(tmpDir, "p1")
	os.MkdirAll(pluginDir, 0755)
	os.WriteFile(filepath.Join(pluginDir, "p1"), []byte("bin"), 0755)
	os.WriteFile(filepath.Join(pluginDir, "plugin.yaml"), []byte("name: p1"), 0644)

	// Override os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"manifest-sync", pluginDir}

	// We can't easily catch os.Exit(1) without a helper process,
	// but we can call main and if it doesn't exit, it succeeded.
	// In a real CI, we might skip this or use a wrapper.
	main()
}

func TestMainCall_Error(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		cmd := exec.Command(os.Args[0], "-test.run=TestMainCall_Error")
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
		err := cmd.Run()
		assert.Error(t, err) // Should exit 1
		return
	}
	
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"manifest-sync"} // Missing args
	main()
}
