# RFC: Zero Trust Vulnerability Scanner (ZTVS)

**Status:** Draft v1
**Language:** Go
**Authors:** ChatGPT / Product Architecture Draft
**Date:** 2026-04-01

---

# 1. Abstract

This RFC defines the architecture, interfaces, extensibility model, security posture, packaging, and operational semantics for **ZTVS (Zero Trust Vulnerability Scanner)** — a cross-platform, plugin-based system vulnerability scanning tool.

The system is designed with the following first principles:

* **single-binary installable core**
* **cross-platform support**
* **plugin-first extensibility**
* **zero trust execution assumptions**
* **evidence-backed findings**
* **least-privilege runtime**
* **deterministic machine-readable output**

The implementation language for the core engine SHALL be **Go**.

---

# 2. Goals

The system MUST support:

* Windows
* Linux
* macOS

The system MUST provide:

* easy installation
* plugin isolation
* version-safe extension APIs
* fast parallel execution
* structured output formats
* future SaaS / agent deployment compatibility

---

# 3. Non-Goals (v1)

The following are explicitly excluded from v1:

* remote active exploitation
* penetration testing automation
* EDR/AV replacement
* live memory forensics
* kernel drivers
* cloud workload runtime enforcement
* web application vulnerability scanning

This RFC focuses on **host posture validation and system-level security checks**.

---

# 4. Design Principles

---

## 4.1 Zero Trust Principles

The system SHALL assume that:

* host state may be compromised
* plugins may fail or behave maliciously
* results require evidence
* privileges must be explicitly declared
* network access must be denied by default

---

## 4.2 Least Privilege

Every plugin MUST declare required capabilities.

Example:

```yaml
capabilities:
  - read_files
  - execute_commands
  - inspect_processes
```

The host SHALL decide whether these capabilities are granted.

---

## 4.3 Crash Isolation

Plugins SHALL execute **out-of-process**.

A plugin crash MUST NOT crash the scanner host.

---

## 4.4 Deterministic Output

For identical host state and identical plugin versions, results SHOULD be deterministic.

---

# 5. High-Level Architecture

```text
+------------------------------------------------------+
|                    ZTVS CLI CORE                     |
|------------------------------------------------------|
| CLI | Config | Policy | Plugin Host | Report Engine |
+-----------------------+------------------------------+
                        |
                        | JSON-RPC over stdio
                        |
        +---------------+----------------+
        |               |                |
        v               v                v
+---------------+ +---------------+ +---------------+
| plugin-os     | | plugin-net    | | plugin-pkg    |
+---------------+ +---------------+ +---------------+
```

---

# 6. Component Architecture

---

## 6.1 Core Engine

Package:

```text
/internal/engine
```

Responsibilities:

* load config
* discover plugins
* schedule checks
* enforce policies
* aggregate findings
* generate reports

Core interfaces:

```go
type Engine interface {
    Scan(ctx context.Context, target Target) (*ScanResult, error)
}
```

---

## 6.2 CLI Layer

Package:

```text
/cmd/zt
```

Primary commands:

```bash
zt scan
zt plugin list
zt plugin install
zt report
zt policy validate
```

Example:

```bash
zt scan --profile cis-lite --format sarif
```

---

## 6.3 Plugin Host

Package:

```text
/internal/pluginhost
```

Responsibilities:

* plugin discovery
* version negotiation
* capability validation
* process lifecycle
* timeout enforcement
* memory limits
* handshake verification

---

## 6.4 Reporting Engine

Package:

```text
/internal/report
```

Supported outputs:

* terminal text
* JSON
* SARIF
* YAML

Future:

* HTML
* PDF
* SIEM connectors

---

# 7. Plugin Architecture

---

## 7.1 Model

Plugins SHALL be **separate executables**.

This is a normative requirement.

Dynamic shared libraries SHALL NOT be used in v1.

Reasoning:

* cross-platform stability
* version-safe
* crash isolation
* security boundary

---

## 7.2 Plugin Lifecycle

Lifecycle:

```text
discover → handshake → validate → execute → terminate
```

---

## 7.3 Discovery

Default locations:

Linux/macOS:

```text
/usr/local/lib/zt/plugins
~/.zt/plugins
./plugins
```

Windows:

```text
C:\Program Files\zt\plugins
%USERPROFILE%\.zt\plugins
```

---

## 7.4 Manifest

Each plugin MUST provide a manifest.

Example:

```yaml
name: plugin-os
version: 1.0.0
api_version: 1
platforms:
  - linux
  - windows
checks:
  - ssh_config
  - firewall_status
capabilities:
  - read_files
  - execute_commands
signature: sha256:xxxx
```

---

## 7.5 Handshake Protocol

Transport: **JSON-RPC over stdio**

Initial handshake request:

```json
{
  "method": "handshake",
  "params": {
    "host_version": "1.0.0",
    "api_version": 1
  }
}
```

Response:

```json
{
  "name": "plugin-os",
  "version": "1.0.0",
  "api_version": 1,
  "supported_checks": [
    "ssh_config",
    "firewall_status"
  ]
}
```

---

# 8. Plugin SDK

Package:

```text
/pkg/sdk
```

---

## 8.1 Core Interface

```go
type Check interface {
    ID() string
    Name() string
    Description() string
    Platforms() []string
    Run(ctx context.Context, host HostContext) (*Finding, error)
}
```

---

## 8.2 Plugin Interface

```go
type Plugin interface {
    Name() string
    Version() string
    Checks() []Check
}
```

---

# 9. Capability Model

This is core to zero trust.

---

## 9.1 Capability Enum

```go
type Capability string

const (
    CapReadFiles        Capability = "read_files"
    CapExecCommands     Capability = "execute_commands"
    CapReadRegistry     Capability = "read_registry"
    CapInspectProcesses Capability = "inspect_processes"
    CapNetworkAccess    Capability = "network_access"
)
```

---

## 9.2 Enforcement

Host SHALL reject plugins requesting undeclared permissions.

Example:

```yaml
policy:
  deny:
    - network_access
```

---

# 10. Scan Execution Flow

```text
1. load config
2. discover plugins
3. validate signatures
4. handshake
5. detect host OS
6. resolve runnable checks
7. execute in parallel
8. normalize findings
9. apply policy
10. emit report
```

---

# 11. Parallel Execution

The engine MUST support bounded concurrency.

Example:

```go
workerPool := make(chan struct{}, maxParallel)
```

Recommended default:

```text
max_parallel = CPU cores * 2
```

Plugins SHALL run with:

* timeout
* panic recovery
* context cancellation

---

# 12. Data Model

---

## 12.1 Finding

```go
type Finding struct {
    ID          string            `json:"id"`
    CheckID     string            `json:"check_id"`
    Severity    Severity          `json:"severity"`
    Confidence  Confidence        `json:"confidence"`
    Title       string            `json:"title"`
    Description string            `json:"description"`
    Evidence    map[string]any    `json:"evidence"`
    Remediation string            `json:"remediation"`
    Tags        []string          `json:"tags"`
    Plugin      string            `json:"plugin"`
    Timestamp   time.Time         `json:"timestamp"`
}
```

---

## 12.2 Severity

```go
type Severity string

const (
    Critical Severity = "critical"
    High     Severity = "high"
    Medium   Severity = "medium"
    Low      Severity = "low"
    Info     Severity = "info"
)
```

---

## 12.3 Result Status

```go
type CheckStatus string

const (
    Pass          CheckStatus = "pass"
    Fail          CheckStatus = "fail"
    Warning       CheckStatus = "warning"
    Skipped       CheckStatus = "skipped"
    Error         CheckStatus = "error"
    NotApplicable CheckStatus = "not_applicable"
)
```

---

# 13. Reporting Formats

---

## 13.1 JSON

Primary machine format.

```json
{
  "summary": {
    "critical": 1,
    "high": 2
  },
  "findings": []
}
```

---

## 13.2 SARIF

Required for CI/CD integration.

This enables GitHub and DevSecOps pipelines.

---

## 13.3 Human Terminal Output

Example:

```text
CRITICAL  SSH root login enabled
HIGH      Firewall disabled
MEDIUM    Package openssl outdated
```

---

# 14. Configuration

Default file:

```text
~/.zt/config.yaml
```

---

## Example

```yaml
version: 1

plugins:
  enabled:
    - plugin-os
    - plugin-network

runtime:
  timeout_seconds: 30
  max_parallel_checks: 8

policy:
  fail_on:
    - critical
    - high
```

---

# 15. Security Requirements

---

## 15.1 Signed Plugins

All plugins SHOULD be signed.

Future MUST requirement.

---

## 15.2 Checksum Validation

Every plugin binary SHALL be verified before execution.

---

## 15.3 No Network by Default

Plugins MUST NOT access the network unless capability granted.

---

## 15.4 Execution Sandboxing

Recommended future:

Linux:

* seccomp
* namespaces

Windows:

* job objects

macOS:

* sandbox profiles

---

# 16. Installation & Distribution

---

## Linux/macOS

Primary distribution:

```bash
curl -fsSL https://install.zt.dev | sh
```

Secondary:

* Homebrew
* apt
* yum

---

## Windows

Primary:

* winget
* signed installer

---

## Binary Strategy

The core SHALL be a **single Go binary**.

Example builds:

```bash
GOOS=linux GOARCH=amd64
GOOS=windows GOARCH=amd64
GOOS=darwin GOARCH=arm64
```

---

# 17. Suggested Repository Layout

```text
zt/
├── cmd/zt
├── internal/engine
├── internal/pluginhost
├── internal/report
├── internal/policy
├── pkg/sdk
├── pkg/types
├── plugins
│   ├── plugin-os
│   ├── plugin-network
│   └── plugin-package
├── schemas
└── docs
```

---

# 18. First-Party Plugins (v1)

---

## plugin-os

Checks:

* SSH config
* firewall
* sudo policy
* insecure services
* file permissions

---

## plugin-network

Checks:

* open ports
* metadata exposure
* weak TLS settings
* localhost admin interfaces

---

## plugin-package

Checks:

* outdated packages
* unsupported OS versions
* CVE matches

---

# 19. Future RFCs

Planned follow-up RFCs:

* RFC-002 Plugin Marketplace
* RFC-003 Policy-as-Code
* RFC-004 Remote Agent Mode
* RFC-005 Cloud Asset Scanning
* RFC-006 SBOM + CVE Correlation

---

# 20. Recommended MVP Roadmap

Phase 1:

* core engine
* 3 plugins
* JSON + terminal output

Phase 2:

* SARIF
* CI integrations
* signed plugins

Phase 3:

* cloud mode
* dashboard backend
* agent fleet scanning

---

# Final Recommendation

This design is production-grade and highly extensible.

The **Go + out-of-process plugin model** is the correct long-term architecture for this product.
