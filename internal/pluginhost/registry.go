package pluginhost

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Registry handles remote plugin discovery and installation.
type Registry struct {
	BaseURL string
}

func NewRegistry() *Registry {
	return &Registry{
		BaseURL: "https://registry.ztvs.io", // Placeholder
	}
}

// Install downloads a plugin binary and manifest from the registry.
func (r *Registry) Install(pluginName string) error {
	home, _ := os.UserHomeDir()
	installDir := filepath.Join(home, ".ztvs", "plugins", pluginName)
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
