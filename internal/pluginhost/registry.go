package pluginhost

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mosesgameli/ztvs/internal/config"
	"github.com/mosesgameli/ztvs/pkg/registry"
)

// Registry handles remote plugin discovery and installation.
type Registry struct {
	BaseURL string
}

func NewRegistry() *Registry {
	return &Registry{
		BaseURL: "https://plugins.ztvs.dev",
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

// Install downloads a plugin binary and manifest from the registry.
func (r *Registry) Install(pluginName string) error {
	installDir := filepath.Join(config.ConfigDir(), "plugins", pluginName)
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return err
	}

	fmt.Printf("Installing plugin %s from %s...\n", pluginName, r.BaseURL)

	// In a real implementation, we would:
	// 1. Download plugin.yaml
	// 2. Download binary
	// 3. Verify integrity (SHA-256) using the logic in integrity.go

	return fmt.Errorf("remote registry connectivity not implemented in this phase")
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
