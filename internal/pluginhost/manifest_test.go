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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadManifest(t *testing.T) {
	h := New()

	t.Run("valid manifest", func(t *testing.T) {
		tmpDir := t.TempDir()
		manifestPath := filepath.Join(tmpDir, "plugin.yaml")
		content := `
name: test-plugin
version: 1.0.0
runtime:
  type: go
  entrypoint: main
capabilities:
  - read_files
`
		err := os.WriteFile(manifestPath, []byte(content), 0644)
		require.NoError(t, err)

		m, err := h.loadManifest(manifestPath)
		assert.NoError(t, err)
		assert.Equal(t, "test-plugin", m.Name)
		assert.Equal(t, "1.0.0", m.Version)
		assert.Equal(t, "go", m.Runtime.Type)
		assert.Contains(t, m.Capabilities, "read_files")
	})

	t.Run("invalid yaml", func(t *testing.T) {
		tmpDir := t.TempDir()
		manifestPath := filepath.Join(tmpDir, "plugin.yaml")
		content := `
name: test-plugin: invalid
`
		err := os.WriteFile(manifestPath, []byte(content), 0644)
		require.NoError(t, err)

		_, err = h.loadManifest(manifestPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parse manifest")
	})

	t.Run("missing file", func(t *testing.T) {
		_, err := h.loadManifest("non-existent.yaml")
		assert.Error(t, err)
		assert.True(t, os.IsNotExist(err))
	})
}
