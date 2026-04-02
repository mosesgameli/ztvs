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
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscover(t *testing.T) {
	ctx := context.Background()

	t.Run("discover valid plugin", func(t *testing.T) {
		h := New()
		tmpDir := t.TempDir()
		h.paths = []string{tmpDir}

		pluginDir := filepath.Join(tmpDir, "test-plugin")
		err := os.MkdirAll(pluginDir, 0755)
		require.NoError(t, err)

		manifestPath := filepath.Join(pluginDir, "plugin.yaml")
		content := `
name: test-plugin
version: 1.0.0
runtime:
  type: binary
  entrypoint: main-bin
`
		err = os.WriteFile(manifestPath, []byte(content), 0644)
		require.NoError(t, err)

		// Create mock entrypoint binary (for BinaryRunner validation)
		binPath := filepath.Join(pluginDir, "main-bin")
		err = os.WriteFile(binPath, []byte("#!/bin/sh\necho ok"), 0755)
		require.NoError(t, err)

		discovered, err := h.Discover(ctx)
		assert.NoError(t, err)
		assert.Len(t, discovered, 1)
		assert.Contains(t, discovered[0], "main-bin")

		info, ok := h.GetPluginInfo(discovered[0])
		assert.True(t, ok)
		assert.Equal(t, "test-plugin", info.Manifest.Name)
	})

	t.Run("ignore missing manifest", func(t *testing.T) {
		h := New()
		tmpDir := t.TempDir()
		h.paths = []string{tmpDir}

		pluginDir := filepath.Join(tmpDir, "invalid-plugin")
		err := os.MkdirAll(pluginDir, 0755)
		require.NoError(t, err)

		discovered, err := h.Discover(ctx)
		assert.NoError(t, err)
		assert.Empty(t, discovered)
	})

	t.Run("ignore unsupported runtime", func(t *testing.T) {
		h := New()
		tmpDir := t.TempDir()
		h.paths = []string{tmpDir}

		pluginDir := filepath.Join(tmpDir, "unsupported-plugin")
		err := os.MkdirAll(pluginDir, 0755)
		require.NoError(t, err)

		manifestPath := filepath.Join(pluginDir, "plugin.yaml")
		content := `
name: unsupported-plugin
version: 1.0.0
runtime:
  type: cobol
  entrypoint: main
`
		err = os.WriteFile(manifestPath, []byte(content), 0644)
		require.NoError(t, err)

		discovered, err := h.Discover(ctx)
		assert.NoError(t, err)
		assert.Empty(t, discovered)
	})
}
