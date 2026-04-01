package main

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mosesgameli/ztvs/pkg/sdk"
)

type DependencyCheck struct{}

func (c *DependencyCheck) ID() string {
	return "axios_dependency_audit"
}

func (c *DependencyCheck) Name() string {
	return "Axios Dependency Audit"
}

func (c *DependencyCheck) Run(ctx context.Context) (*sdk.Finding, error) {
	finding := &sdk.Finding{
		ID:          "F-AXIOS-001",
		Severity:    "critical",
		Title:       "Compromised Axios version detected",
		Evidence:    make(map[string]interface{}),
		Remediation: "Remove axios versions 1.14.1 or 0.30.4. Audit plain-crypto-js in your dependency tree. Rotate all secrets.",
	}

	foundInLockfile := false

	// 1. Scan current and parent directories for lockfiles (Immediate Workspace)
	lockfiles := []string{"package-lock.json", "yarn.lock", "bun.lockb", "bun.lock"}
	cwd, _ := os.Getwd()
	dir := cwd
	for {
		for _, lf := range lockfiles {
			path := filepath.Join(dir, lf)
			if _, err := os.Stat(path); err == nil {
				if c.scanLockfile(path, finding) {
					foundInLockfile = true
				}
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// 2. System-wide scan for plain-crypto-js using native search tools
	nodeModulesPaths := c.findNodeModules(ctx)
	for _, path := range nodeModulesPaths {
		if ctx.Err() != nil {
			break
		}
		// Check for malicious dependency within found node_modules
		target := filepath.Join(path, "plain-crypto-js")
		if _, err := os.Stat(target); err == nil {
			finding.Title = "Malicious dependency 'plain-crypto-js' found on system"
			finding.Evidence["malicious_package_path"] = target
			foundInLockfile = true
		}
	}

	if foundInLockfile {
		return finding, nil
	}

	return nil, nil
}

func (c *DependencyCheck) findNodeModules(ctx context.Context) []string {
	var paths []string

	// Tier 1: Spotlight (macOS)
	if runtime.GOOS == "darwin" {
		cmd := exec.CommandContext(ctx, "mdfind", "kMDItemFSName == 'node_modules' && kMDItemContentType == 'public.folder'")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if trimmed := strings.TrimSpace(line); trimmed != "" {
					paths = append(paths, trimmed)
				}
			}
			return paths
		}
	}

	// Tier 2: Locate (Linux/Unix)
	if runtime.GOOS != "windows" {
		// Try locate but filter specifically for directory endings
		cmd := exec.CommandContext(ctx, "locate", "-r", "/node_modules$")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if trimmed := strings.TrimSpace(line); trimmed != "" {
					paths = append(paths, trimmed)
				}
			}
			if len(paths) > 0 {
				return paths
			}
		}
	}

	// Tier 3: Find (Comprehensive Fallback)
	// We run this as a slower background walk if native indexed search failed
	// but we prioritize home directory and /usr
	searchRoots := []string{"/usr/local/lib", "/usr/lib"}
	home, _ := os.UserHomeDir()
	if home != "" {
		searchRoots = append(searchRoots, home)
	}

	for _, root := range searchRoots {
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if d.IsDir() && d.Name() == "node_modules" {
				paths = append(paths, path)
				return filepath.SkipDir
			}
			return nil
		})
	}

	return paths
}

func (c *DependencyCheck) scanLockfile(path string, finding *sdk.Finding) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	maliciousVersions := []string{"1.14.1", "0.30.4"}
	maliciousDep := "plain-crypto-js"

	detected := false
	for scanner.Scan() {
		line := scanner.Text()
		for _, v := range maliciousVersions {
			if strings.Contains(line, "axios") && strings.Contains(line, v) {
				finding.Evidence["lockfile"] = path
				finding.Evidence["detected_version"] = v
				detected = true
			}
		}
		if strings.Contains(line, maliciousDep) {
			finding.Evidence["lockfile"] = path
			finding.Evidence["malicious_dependency"] = maliciousDep
			detected = true
		}
	}
	return detected
}
