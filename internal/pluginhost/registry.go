package pluginhost

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/pkg/registry"
)

type Registry struct {
	BaseURL string
	git     *Git
}

func NewRegistry() *Registry {
	return &Registry{
		BaseURL: "https://plugins.ztvs.dev",
		git:     NewGit(),
	}
}

func (r *Registry) FetchIndex(ctx context.Context) (*registry.Index, error) {
	configDir := config.ConfigDir()
	cacheDir := filepath.Join(configDir, "cache")
	indexPath := filepath.Join(cacheDir, "index.json")

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}

	// 1. Fetch index.json
	indexURL := fmt.Sprintf("%s/index.json", r.BaseURL)
	if err := DownloadFile(indexPath, indexURL); err != nil {
		// If offline, try to load from cache
		return r.loadLocalIndex(indexPath)
	}

	// 2. Fetch signature
	sigPath := indexPath + ".sig"
	sigURL := indexURL + ".sig"
	if err := DownloadFile(sigPath, sigURL); err != nil {
		return nil, fmt.Errorf("failed to fetch index signature: %v", err)
	}

	// 3. Verify signature
	data, _ := os.ReadFile(indexPath)
	sig, _ := os.ReadFile(sigPath)
	if err := VerifySignature(data, sig); err != nil {
		return nil, fmt.Errorf("registry index signature verification failed: %v", err)
	}

	return r.loadLocalIndex(indexPath)
}

func (r *Registry) loadLocalIndex(path string) (*registry.Index, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var idx registry.Index
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	return &idx, nil
}

func (r *Registry) Search(ctx context.Context, query string) ([]registry.PluginMetadata, error) {
	idx, err := r.FetchIndex(ctx)
	if err != nil {
		return nil, err
	}

	var results []registry.PluginMetadata
	query = strings.ToLower(query)
	for _, m := range idx.Plugins {
		if strings.Contains(strings.ToLower(m.Name), query) {
			results = append(results, m)
		}
	}
	return results, nil
}

func (r *Registry) GetInfo(ctx context.Context, name string) (*registry.PluginMetadata, error) {
	idx, err := r.FetchIndex(ctx)
	if err != nil {
		return nil, err
	}

	for _, m := range idx.Plugins {
		if m.Name == name {
			return &m, nil
		}
	}
	return nil, fmt.Errorf("plugin %s not found in registry", name)
}

// Install downloads, builds, and registers a plugin.
func (r *Registry) Install(ctx context.Context, pluginName string, host *Host) error {
	// 1. Get Metadata
	meta, err := r.GetInfo(ctx, pluginName)
	if err != nil {
		return err
	}

	installDir := filepath.Join(config.ConfigDir(), "plugins", pluginName)
	if _, err := os.Stat(installDir); err == nil {
		fmt.Printf("Plugin %s is already installed. Use 'zt plugin update' to update.\n", pluginName)
		return nil
	}

	fmt.Printf("Found plugin %s at %s\n", pluginName, meta.Repo)

	// 2. Git Clone
	fmt.Printf("Cloning repository...\n")
	if err := r.git.Clone(meta.Repo, installDir); err != nil {
		return err
	}

	// 3. Build & Verify
	fmt.Printf("Building plugin (Go)...\n")
	binPath := filepath.Join(installDir, pluginName)
	if err := r.buildPlugin(installDir, pluginName); err != nil {
		fmt.Printf("Build failed: %v. Cleaning up...\n", err)
		r.git.Remove(installDir) // User-requested: Remove on failure
		return err
	}

	// 4. Verify Integrity
	fmt.Printf("Verifying integrity...\n")
	if err := VerifyIntegrity(binPath, meta.Checksum); err != nil {
		fmt.Printf("Integrity check failed: %v. Cleaning up...\n", err)
		r.git.Remove(installDir)
		return err
	}

	// 5. Update Lockfile
	fmt.Printf("Updating lockfile...\n")
	lf := host.Lockfile()
	lf.Set(pluginName, registry.PluginLock{
		Version:  meta.LatestVersion,
		Enabled:  true,
		Checksum: meta.Checksum,
	})
	if err := lf.Save(); err != nil {
		return fmt.Errorf("failed to save lockfile: %v", err)
	}

	fmt.Printf("✓ Plugin %s installed successfully!\n", pluginName)
	return nil
}

func (r *Registry) buildPlugin(dir, name string) error {
	cmd := exec.Command("go", "build", "-o", name, ".")
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("build error: %v\nOutput: %s", err, string(output))
	}
	return nil
}

func (r *Registry) CheckAndUpdateAll(ctx context.Context, host *Host, mode string) error {
	if mode == "locked" {
		return nil
	}

	fmt.Printf("Checking for plugin updates (mode: %s)...\n", mode)
	
	// Set 15-second timeout as requested
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	idx, err := r.FetchIndex(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch registry index: %v", err)
	}

	lf := host.Lockfile()
	plugins := lf.All()

	for name, lock := range plugins {
		// Find in index
		var meta *registry.PluginMetadata
		for _, m := range idx.Plugins {
			if m.Name == name {
				meta = &m
				break
			}
		}

		if meta == nil {
			continue
		}

		needsUpdate := false
		if mode == "always" {
			needsUpdate = true
		} else if meta.LatestVersion != lock.Version {
			needsUpdate = true
		}

		if needsUpdate {
			fmt.Printf("Updating plugin %s: %s -> %s\n", name, lock.Version, meta.LatestVersion)
			if err := r.PerformAtomicUpdate(ctx, name, host, meta); err != nil {
				fmt.Printf("Failed to update %s: %v\n", name, err)
			}
		}
	}

	return nil
}

func (r *Registry) PerformAtomicUpdate(ctx context.Context, name string, host *Host, meta *registry.PluginMetadata) error {
	configDir := config.ConfigDir()
	tmpDir := filepath.Join(configDir, "cache", "tmp", name)
	finalDir := filepath.Join(configDir, "plugins", name)

	_ = os.RemoveAll(tmpDir)
	if err := os.MkdirAll(filepath.Dir(tmpDir), 0755); err != nil {
		return err
	}

	// 1. Clone to tmp
	if err := r.git.Clone(meta.Repo, tmpDir); err != nil {
		return err
	}

	// 2. Build
	if err := r.buildPlugin(tmpDir, name); err != nil {
		_ = os.RemoveAll(tmpDir)
		return err
	}

	// 3. Verify
	binPath := filepath.Join(tmpDir, name)
	if err := VerifyIntegrity(binPath, meta.Checksum); err != nil {
		_ = os.RemoveAll(tmpDir)
		return err
	}

	// 4. Atomic Swap
	// On Unix, os.Rename is atomic. 
	// We need to move the old one aside or just replace it if we don't care about rollback for now.
	// RFC says: "rename swap"
	oldDir := finalDir + ".old"
	_ = os.RemoveAll(oldDir)
	
	if _, err := os.Stat(finalDir); err == nil {
		if err := os.Rename(finalDir, oldDir); err != nil {
			return err
		}
	}

	if err := os.Rename(tmpDir, finalDir); err != nil {
		// Try to restore old one
		_ = os.Rename(oldDir, finalDir)
		return err
	}
	_ = os.RemoveAll(oldDir)

	// 5. Update Lockfile
	lf := host.Lockfile()
	lf.Set(name, registry.PluginLock{
		Version:  meta.LatestVersion,
		Enabled:  true,
		Checksum: meta.Checksum,
	})
	return lf.Save()
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
