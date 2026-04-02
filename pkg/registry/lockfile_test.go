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

package registry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLockfile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lockfile-test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	lockPath := filepath.Join(tmpDir, "plugins.lock")
	lf := NewLockfile(lockPath)

	// Test Set/Get
	lf.Set("test-plugin", PluginLock{Version: "1.0.0", Enabled: true})
	lock, ok := lf.Get("test-plugin")
	if !ok || lock.Version != "1.0.0" || !lock.Enabled {
		t.Errorf("expected version 1.0.0 and enabled, got %v", lock)
	}

	// Test Save/Load
	if err := lf.Save(); err != nil {
		t.Fatal(err)
	}

	lf2, err := LoadLockfile(lockPath)
	if err != nil {
		t.Fatal(err)
	}

	lock2, ok := lf2.Get("test-plugin")
	if !ok || lock2.Version != "1.0.0" || !lock2.Enabled {
		t.Errorf("expected loaded version 1.0.0 and enabled, got %v", lock2)
	}

	// Test Disable
	lock2.Enabled = false
	lf2.Set("test-plugin", lock2)
	if err := lf2.Save(); err != nil {
		t.Fatal(err)
	}

	lf3, err := LoadLockfile(lockPath)
	if err != nil {
		t.Fatal(err)
	}
	lock3, _ := lf3.Get("test-plugin")
	if lock3.Enabled {
		t.Error("expected plugin to be disabled")
	}

	// Test All/Remove
	all := lf3.All()
	assert.Len(t, all, 1)

	lf3.Remove("test-plugin")
	all = lf3.All()
	assert.Empty(t, all)

	// Test non-existent file load
	lf4, err := LoadLockfile("non-existent.lock")
	require.NoError(t, err)
	assert.NotNil(t, lf4)
	assert.Empty(t, lf4.Plugins)

	// Test invalid YAML load
	invalidPath := filepath.Join(tmpDir, "invalid.lock")
	_ = os.WriteFile(invalidPath, []byte("invalid: yaml: :"), 0644)
	_, err = LoadLockfile(invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse lockfile")

	// Test Save error (Permission Denied)
	if os.Getuid() != 0 { // Skip if running as root
		readOnlyDir := filepath.Join(tmpDir, "readonly")
		_ = os.MkdirAll(readOnlyDir, 0555)
		badPath := filepath.Join(readOnlyDir, "plugins.lock")
		lfBad := NewLockfile(badPath)
		err = lfBad.Save()
		assert.Error(t, err)
	}

	// Test corrupt YAML for Load
	p := filepath.Join(tmpDir, "corrupt.lock")
	_ = os.WriteFile(p, []byte("!!invalid yaml"), 0644)
	_, err = LoadLockfile(p)
	assert.Error(t, err)

	// Test ReadFile error (Permission Denied)
	if os.Getuid() != 0 {
		readOnlyFile := filepath.Join(tmpDir, "noperm.lock")
		_ = os.WriteFile(readOnlyFile, []byte("data"), 0000)
		defer func() { _ = os.Chmod(readOnlyFile, 0644) }()
		_, err = LoadLockfile(readOnlyFile)
		assert.Error(t, err)
	}

	// Test Save error: MkdirAll failure
	conflictDir := filepath.Join(tmpDir, "conflict")
	_ = os.WriteFile(conflictDir, []byte("actually-a-file"), 0644)
	badPath2 := filepath.Join(conflictDir, "sub", "plugins.lock")
	lfBad2 := NewLockfile(badPath2)
	err = lfBad2.Save()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a directory")
}
