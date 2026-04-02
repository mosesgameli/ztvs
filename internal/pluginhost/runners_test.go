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
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprint(os.Stdout, "{\"jsonrpc\":\"2.0\",\"id\":\"1\",\"result\":{\"status\":\"ok\"}}")
	os.Exit(0)
}

func TestBinaryRunner(t *testing.T) {
	r := &BinaryRunner{}
	assert.Equal(t, "Binary", r.Name())
	assert.True(t, r.Supports("binary"))
	assert.False(t, r.Supports("python"))

	t.Run("Validate", func(t *testing.T) {
		tmpDir := t.TempDir()
		binPath := filepath.Join(tmpDir, "test-bin")
		
		// 1. Missing file
		err := r.Validate(binPath)
		assert.Error(t, err)

		// 2. Not executable
		err = os.WriteFile(binPath, []byte("test"), 0644)
		require.NoError(t, err)
		err = r.Validate(binPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not executable")

		// 3. Valid executable
		err = os.Chmod(binPath, 0755)
		require.NoError(t, err)
		err = r.Validate(binPath)
		assert.NoError(t, err)
	})

	t.Run("Execute", func(t *testing.T) {
		tmpDir := t.TempDir()
		mockPath := filepath.Join(tmpDir, "mock-plugin.sh")
		content := "#!/bin/sh\necho '{\"jsonrpc\":\"2.0\",\"id\":\"1\",\"result\":{\"status\":\"ok\"}}'\n"
		err := os.WriteFile(mockPath, []byte(content), 0755)
		require.NoError(t, err)

		ctx := context.Background()
		out, err := r.Execute(ctx, mockPath, []byte("{}"))
		assert.NoError(t, err)
		assert.Contains(t, string(out), "jsonrpc")
	})
}

func TestPythonRunner(t *testing.T) {
	r := &PythonRunner{}
	assert.Contains(t, r.Name(), "Python")
	assert.True(t, r.Supports("python"))

	t.Run("Validate - Missing file", func(t *testing.T) {
		err := r.Validate("non-existent.py")
		assert.Error(t, err)
	})
}

func TestNodeRunner(t *testing.T) {
	r := &NodeRunner{}
	assert.Equal(t, "Node.js", r.Name())
	assert.True(t, r.Supports("node"))
}

func TestJavaRunner(t *testing.T) {
	r := &JavaRunner{}
	assert.Equal(t, "Java", r.Name())
	assert.True(t, r.Supports("java"))
	
	err := r.Validate("any.jar")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not yet supported")

	_, err = r.Execute(context.Background(), "any.jar", nil)
	assert.Error(t, err)
}
