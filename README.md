# Zero Trust Vulnerability Scanner (ZTVS)

[![Go Report Card](https://goreportcard.com/badge/github.com/mosesgameli/ztvs)](https://goreportcard.com/report/github.com/mosesgameli/ztvs)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**ZTVS** is a next-generation, cross-platform host security scanner built from the ground up on a **Zero Trust execution model**. Traditional security agents run monolithic, highly privileged processes where a single compromised or crashing rule can expose the entire host. ZTVS isolates every security check (plugin) enforcing strict least-privilege capability models at execution time. 

Built for modern infrastructure, ZTVS natively supports polyglot plugins, allowing specialized security logic to be written in Go, Rust, Python, Node.js, and Java, all operating deterministically over our JSON-RPC 2.0 Subprocess Protocol.

---

## 🚀 Key Features

*   **🛡️ Zero Trust Architecture**: Out-of-process isolation. Plugins execute as untrusted child processes; crashes, hangs, or compromises cannot affect the Host engine or other plugins.
*   **⚡ High-Concurrency Core Engine**: Go-powered, heavily parallelized worker pool optimizing execution time across hundreds of granular checks.
*   **🛠️ Polyglot Runner Support**: Write checks in your language of choice. ZTVS dynamically binds to runtimes (`python3`, `node`) or executes natively compiled binaries (Go, Rust).
*   **🔒 Strict Capability Controls**: Plugins must define exact capabilities in `plugin.yaml` (e.g., `network_access`, `read_fs`). Granular enforcement via `~/.ztvs/config.yaml` policies.
*   **🚨 Sophisticated Supply Chain Auditing**: Native support for deep host dependency validation, mitigating targeted compromises like the 2026 Axios Supply Chain RAT attack out of the box (`plugin-axios-mitigation`).
*   **🔌 Comprehensive Plugin SDKs**: Fast, integrated developer experience with our first-party Go SDK (`pkg/sdk`).
*   **📡 Agent & One-Shot Modes**: Run dynamically as an isolated scheduled background daemon (`zt agent`) or on-demand (`zt scan`).
*   **📦 Manifest-Driven Remote Registry**: Discover, cryptographically verify (via SHA-256 integrity), and install plugins safely using `zt plugin install`. Live atomic updates prevent failures during execution.
*   **📊 Enterprise Reporting**: Fully structured output integrations for SIEM and CI/CD, including Terminal pretty-print, JSON, and full **SARIF 2.1.0**.

---

## 🧱 Architecture Details

ZTVS follows a deliberate Host + Plugin model:

1. **Host Engine (`internal/`):** Manages discovery, policy arbitration, scheduling, protocol negotiation, and reporting.
2. **Plugins (`plugins/`):** Ephemeral, untrusted executables.
3. **Communication (`pkg/rpc`):** Fast JSON-RPC 2.0 messaging over standard `stdin`/`stdout`. Custom log streams are routed over `stderr`.
4. **Sandboxing Constraints:** We prevent rogue plugins from lateral movement by tracking explicit permissions.

For rigorous details on the internal layout, request flows, and state machines, see [ZTVS Architecture](docs/architecture.md) and [Protocol Specification](docs/protocol.md).

---

## 🛠️ Installation

### Linux & macOS

```sh
curl -fsSL https://raw.githubusercontent.com/mosesgameli/ztvs/main/install.sh | sh
```
*Installs `zt` binary to `/usr/local/bin` and provisions the system plugin dir at `~/.ztvs/plugins`.*

> **Pin a specific version:** 
> `ZTVS_VERSION=v1.0.0 curl -fsSL ... | sh`

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/mosesgameli/ztvs/main/install.ps1 | iex
```
*Installs `zt.exe` to `%LOCALAPPDATA%\Programs\ztvs` (no admin required) and updates user `PATH`.*

### Compilation from Source

Requires [Go 1.24+](https://golang.org/dl/) and standard build tools:

```bash
git clone git@github.com:mosesgameli/ztvs.git
cd ztvs
make build     # Builds host and all first-party plugins
make install   # Places binaries in standard paths
```

---

## 📖 Usage Guide

Evaluate an entire host targeting local configurations, deployed artifacts, and active processes:

```bash
# Initialize environments and fetch base plugin sets
make init

# Standard verbose manual scan
zt scan

# Output results in SARIF for GitHub Advanced Security ingestion
zt --format sarif scan > report.sarif

# Continuous Background Auditing
zt agent
```

### Plugin Lifecycle Management

Leverage the remote registry to safely augment your ZTVS policies:

```bash
# List local plugins and active capabilities
zt plugin list

# Download, verify, and activate a remote plugin
zt plugin install plugin-axios-mitigation

# Run an atomic live update of all installed plugins
zt plugin update
```

---

## 🔌 Developing Plugins

The architectural goal of ZTVS is to democratize security rules. A basic Go plugin requires fewer than 25 lines of code:

### 1. Implement `sdk.Check`

```go
package main

import (
    "context"
    "github.com/mosesgameli/ztvs/pkg/sdk"
)

type MyCheck struct{}

func (c *MyCheck) ID() string   { return "custom_sec_check" }
func (c *MyCheck) Name() string { return "Custom Host Analysis" }

func (c *MyCheck) Run(ctx context.Context) (*sdk.Finding, error) {
    // Zero Trust Logic Here...
    return &sdk.Finding{
        ID: "F-CUSTOM-001",
        Severity: "high",
        Title: "Excessive file permissions detected",
        Description: "Critical key files are globally writable.",
    }, nil
}
```

### 2. Boot the Payload

```go
func main() {
    sdk.Run(sdk.Metadata{
        Name: "plugin-custom-checks",
        Version: "1.0.0",
        APIVersion: 1, // Must match JSON-RPC spec version
    }, []sdk.Check{&MyCheck{}})
}
```

View the detailed guide at [Plugin Developer Docs](.agents/skills/plugin-dev/SKILL.md) and [RFC Specifications](docs/).

---

## 🗺️ Roadmap & Operational Vision

While the core parallel JSON-RPC engine is stable, advanced functionality is under development:
- [Detailed Roadmap](_/work/roadmap.md)
- **Deep OS Sandboxing**: Migrating to `seccomp` (Linux), `Job Objects/AppContainers` (Windows), and `Sandbox.kext` (macOS).
- **Cryptographic Plugin Trust**: Enforcing full x509/PKI signing of remote plugins vs standard SHA-256.

## 📄 License
Source code distributed under the **MIT License**. See [`LICENSE`](LICENSE) for details.
