# Zero Trust Vulnerability Scanner (ZTVS)

[![Go Report Card](https://goreportcard.com/badge/github.com/mosesgameli/ztvs)](https://goreportcard.com/report/github.com/mosesgameli/ztvs)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**ZTVS** is a cross-platform, plugin-based host security scanner designed for high performance and strict isolation. It implements a Zero Trust execution model, ensuring security checks run in isolated environments with minimal privileges.

---

## 🚀 Key Features

*   **⚡ High-Concurrency Scanning**: Parallel check execution via a Go-powered engine.
*   **🛡️ Zero Trust Architecture**: Out-of-process plugin isolation with JSON-RPC 2.0 communication.
*   **🔌 Easy Extensibility**: Simple-to-use Go SDK for building new custom security checks.
*   **📊 Multiple Output Formats**: Supports Terminal (Pretty-print), JSON, and **SARIF** for CI/CD integration.
*   **🔍 Cross-Path Discovery**: Automatically discovers plugins in local, user, and system directories.

---

## 🛠️ Installation

### Prerequisites
*   [Go 1.24+](https://golang.org/dl/)
*   `make` (optional, for builds)

### Build from Source
```bash
git clone git@github.com:mosesgameli/ztvs.git
cd ztvs
make build
```

---

## 📖 Usage

Run a standard scan on the local host:
```bash
./zt scan
```

### Advanced Options
Export findings in JSON format for automation:
```bash
./zt --format json scan
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
| `plugins/` | First-party security plugins (e.g., `plugin-os`) |

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
The project is currently in **Phase 2 (Reporting & Concurrency)**.
View the full [Detailed Roadamp](_/work/roadmap.md) and [Phased Delivery Plan](_/work/phases.md) for more details.

---

## 📄 License
Distributed under the MIT License. See `LICENSE` for more information.
