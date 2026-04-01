package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: manifest-sync <plugin-dir> [bin-ext]")
		os.Exit(1)
	}

	pluginDir := os.Args[1]
	binExt := ""
	if len(os.Args) > 2 {
		binExt = os.Args[2]
	}
	pluginName := filepath.Base(pluginDir)
	binPath := filepath.Join(pluginDir, pluginName+binExt)
	manifestPath := filepath.Join(pluginDir, "plugin.yaml")

	// 1. Check if binary exists
	f, err := os.Open(binPath)
	if err != nil {
		fmt.Printf("Error: binary not found at %s. Build it first.\n", binPath)
		os.Exit(1)
	}
	defer f.Close()

	// 2. Compute SHA-256
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		fmt.Printf("Error computing hash: %v\n", err)
		os.Exit(1)
	}
	checksum := hex.EncodeToString(h.Sum(nil))

	// 3. Update manifest
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		fmt.Printf("Error reading manifest: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(string(data), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, "checksum:") {
			lines[i] = fmt.Sprintf("checksum: %s", checksum)
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, fmt.Sprintf("checksum: %s", checksum))
	}

	if err := os.WriteFile(manifestPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		fmt.Printf("Error writing manifest: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated %s with checksum %s\n", manifestPath, checksum)
}
