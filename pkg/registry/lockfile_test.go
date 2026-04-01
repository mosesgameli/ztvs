package registry

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLockfile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "lockfile-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

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
}
