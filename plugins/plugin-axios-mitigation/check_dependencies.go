package main

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
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

	// 1. Scan current and parent directories for lockfiles
	lockfiles := []string{"package-lock.json", "yarn.lock", "bun.lockb", "bun.lock"}
	cwd, _ := os.Getwd()
	foundInLockfile := false

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

	// 2. System-wide scan for plain-crypto-js in node_modules
	// We scan common global paths and the user's home directory
	commonPaths := []string{
		"/usr/local/lib/node_modules",
		"/usr/lib/node_modules",
	}

	home, _ := os.UserHomeDir()
	if home != "" {
		commonPaths = append(commonPaths, home)
	}

	for _, p := range commonPaths {
		_ = filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // skip errors
			}
			if info.IsDir() && info.Name() == "node_modules" {
				// Check for plain-crypto-js
				target := filepath.Join(path, "plain-crypto-js")
				if _, err := os.Stat(target); err == nil {
					finding.Title = "Malicious dependency 'plain-crypto-js' found on system"
					finding.Evidence["malicious_package_path"] = target
					foundInLockfile = true
				}
				return filepath.SkipDir // Don't recurse into node_modules
			}
			// Limit depth for home directory scan to avoid extreme slowdowns
			if p == home {
				rel, _ := filepath.Rel(home, path)
				if strings.Count(rel, string(os.PathSeparator)) > 3 {
					return filepath.SkipDir
				}
			}
			return nil
		})
	}

	if foundInLockfile {
		return finding, nil
	}

	return nil, nil
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
