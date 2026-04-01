# Zero Trust Vulnerability Scanner (ZTVS)

[![Go Report Card](https://goreportcard.com/badge/github.com/mosesgameli/ztvs)](https://goreportcard.com/report/github.com/mosesgameli/ztvs)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**ZTVS** is a cross-platform, plugin-based host security scanner designed for high performance and strict isolation. It implements a Zero Trust execution model, ensuring security checks run in isolated environments with minimal privileges.

---

## 🚀 Key Features

*   **⚡ High-Concurrency Scanning**: Parallel check execution via a Go-powered engine.
*   **🛡️ Zero Trust Architecture**: Out-of-process plugin isolation with JSON-RPC 2.0 communication.
*   **🛠️ Polyglot Support**: Native execution for plugins written in Go, Python, Node.js, Rust, and Java.
*   **🔒 Strict Capabilities**: Plugin permissions (e.g., `network_access`, `execute_commands`) are strictly enforced by global policy configurations (`~/.ztvs/config.yaml`).
*   **🚨 Supply Chain & IOC Auditing**: Deep dependency validation (e.g., `plugin-axios-mitigation` for active supply chain compromises and RAT detection).
*   **🔌 Easy Extensibility**: Simple-to-use Go SDK for building new custom security checks.
*   **📡 Agent Mode**: Periodic background auditing for continuous visibility.
*   **📦 Plugin Registry**: Manifest-driven discovery with secure remote installation and atomic updates.
*   **📊 Multiple Output Formats**: Supports Terminal (Pretty-print), JSON, and **SARIF** for CI/CD integration.
*   **🔍 Cross-Path Discovery**: Automatically discovers plugins in local, user, and system directories (`~/.ztvs/plugins`).

---

## 🛠️ Installation

### Linux & macOS

```sh
curl -fsSL https://raw.githubusercontent.com/mosesgameli/ztvs/main/install.sh | sh
```

Installs `zt` to `/usr/local/bin` and first-party plugins to `~/.ztvs/plugins`.

> **Pin a version**: `ZTVS_VERSION=v1.0.0 curl -fsSL https://raw.githubusercontent.com/mosesgameli/ztvs/main/install.sh | sh`

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/mosesgameli/ztvs/main/install.ps1 | iex
```

Installs `zt.exe` to `%LOCALAPPDATA%\Programs\ztvs` (no admin required), seeds first-party plugins to `%USERPROFILE%\.ztvs\plugins`, and adds the install directory to your user `PATH`.

> **Pin a version**: `$env:ZTVS_VERSION="v1.0.0"; irm https://raw.githubusercontent.com/mosesgameli/ztvs/main/install.ps1 | iex`

Both scripts:
1. Auto-detect your architecture (amd64 / arm64).
2. Download the correct pre-built archive from [GitHub Releases](https://github.com/mosesgameli/ztvs/releases).
3. Bootstrap `~/.ztvs/config.yaml` (or `%USERPROFILE%\.ztvs\config.yaml`) on first run.

### Prerequisites (build from source only)
*   [Go 1.24+](https://golang.org/dl/)
*   `make` (Linux/macOS, optional)

### Build from Source
```bash
git clone git@github.com:mosesgameli/ztvs.git
cd ztvs
make build
```

### Uninstallation

To completely remove ZTVS from your system alongside its configuration and downloaded plugins:

**Linux & macOS:**
```sh
curl -fsSL https://raw.githubusercontent.com/mosesgameli/ztvs/main/uninstall.sh | sh
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/mosesgameli/ztvs/main/uninstall.ps1 | iex
```

---


## 📖 Usage

Run a standard scan on the local host:
```bash
# Initialize ZTVS environment
make init

# Build host and all first-party plugins
make build

# Run a manual scan
./zt scan

# Run a scan with JSON output
./zt --format json scan

# Manage plugins locally
./zt plugin list

# Install a plugin from the ZTVS remote registry
./zt plugin install plugin-axios-mitigation

# Automatically update all installed plugins
./zt plugin update

# Start the continuous audit agent
./zt agent
```

Generate a SARIF report for GitHub Advanced Security:
```bash
./zt --format sarif scan
```

---

## 🧱 Project Structure

| Directory | Description |
| :--- | :--- |
| `cmd/zt` | CLI Entrypoint |
| `internal/` | Core Engine, Plugin Host, and Reporting logic |
| `pkg/sdk` | SDK for plugin developers |
| `pkg/rpc` | JSON-RPC 2.0 message definitions |
| `plugins/` | First-party security plugins (e.g., `plugin-os`, `plugin-axios-mitigation`) and polyglot tests (`dummy-python`) |

---

## 🔌 Developing Plugins

Creating a new ZTVS plugin is easy with the included SDK. 

### 1. Define your check
```go
package main

import (
    "context"
    "github.com/mosesgameli/ztvs/pkg/sdk"
)

type MyCheck struct{}

func (c *MyCheck) ID() string   { return "custom_check" }
func (c *MyCheck) Name() string { return "My Security Check" }

func (c *MyCheck) Run(ctx context.Context) (*sdk.Finding, error) {
    // Implement your detection logic here
    return &sdk.Finding{
        ID: "F-CUSTOM-001",
        Severity: "high",
        Title: "Securtity vulnerability found!",
    }, nil
}
```

### 2. Register and Run
```go
func main() {
    sdk.Run(sdk.Metadata{
        Name: "plugin-custom",
        Version: "1.0.0",
        APIVersion: 1,
    }, []sdk.Check{&MyCheck{}})
}
```

Refer to the [.agents/skills/plugin-dev/SKILL.md](.agents/skills/plugin-dev/SKILL.md) for a full guide.

---

## 🗺️ Roadmap
The core development phases are now completed. ZTVS is a production-ready Zero Trust vulnerability platform.
View the full [Detailed Roadamp](_/work/roadmap.md) and [Phased Delivery Plan](_/work/phases.md) for future vision.

---

## 📄 License
Distributed under the MIT License. See `LICENSE` for more information.
