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

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Agent.Interval != "1h" {
		t.Errorf("expected interval 1h, got %s", cfg.Agent.Interval)
	}
	if len(cfg.Policy.AllowedCapabilities) == 0 {
		t.Error("expected allowed capabilities, got empty list")
	}
}

func TestConfigDir(t *testing.T) {
	tempHome := t.TempDir()
	origHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tempHome)
	defer func() { _ = os.Setenv("HOME", origHome) }()

	dir := ConfigDir()
	expected := filepath.Join(tempHome, ".ztvs")
	if dir != expected {
		t.Errorf("expected %s, got %s", expected, dir)
	}
}

func TestLoadSave(t *testing.T) {
	tempHome := t.TempDir()
	origHome := os.Getenv("HOME")
	_ = os.Setenv("HOME", tempHome)
	defer func() { _ = os.Setenv("HOME", origHome) }()

	// 1. Test Load non-existent
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.Agent.Interval != "1h" {
		t.Errorf("expected default interval 1h, got %s", cfg.Agent.Interval)
	}

	// 2. Test Save
	cfg.Agent.Interval = "30m"
	err = cfg.Save()
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// 3. Test Load existing
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() existing error: %v", err)
	}
	if loaded.Agent.Interval != "30m" {
		t.Errorf("expected 30m, got %s", loaded.Agent.Interval)
	}

	// 4. Test Load corrupt (manual file overwrite)
	configPath := filepath.Join(tempHome, ".ztvs", "config.yaml")
	err = os.WriteFile(configPath, []byte("invalid: yaml: :"), 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	_, err = Load()
	if err == nil {
		t.Error("expected Load() to fail with invalid yaml")
	}

	// 5. Test Load Read error (permission denied)
	err = os.Chmod(configPath, 0000)
	if err != nil {
		t.Fatalf("Chmod failed: %v", err)
	}
	defer func() { _ = os.Chmod(configPath, 0644) }()
	_, err = Load()
	if err == nil {
		t.Error("expected Load() to fail with permission denied")
	}

	// 6. Test Save MkdirAll error
	// Create a file where the .ztvs directory should be
	os.RemoveAll(filepath.Join(tempHome, ".ztvs"))
	err = os.WriteFile(filepath.Join(tempHome, ".ztvs"), []byte("not a dir"), 0644)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	err = cfg.Save()
	if err == nil {
		t.Error("expected Save() to fail when directory path is a file")
	}
}
