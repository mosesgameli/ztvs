# ZTVS Strategic Roadmap

This document provides a high-level overview of the Zero Trust Vulnerability Scanner (ZTVS) development milestones and their current status.

## 📊 Phase Summary

| Phase | Title | Focus | Status |
| :--- | :--- | :--- | :--- |
| **0** | **Foundation** | Core architecture & directory layout | ✅ Completed |
| **1** | **Reliability** | Handshake, Discovery, & Timeout support | ✅ Completed |
| **2** | **Reporting** | Concurrency & Structured (JSON/SARIF) output | ✅ Completed |
| **3** | **Security** | Capability enforcement & Plugin sandboxing | ✅ Completed |
| **4** | **Platform** | Remote registry & Agent mode | ✅ Completed |

---

## ✅ Development Checklist

### Phase 0: Foundation
- [x] Repository initialization (`go mod init`)
- [x] Standardized Go directory layout (`internal/`, `pkg/`, `cmd/`)
- [x] Core CLI entrypoint (`zt scan`)
- [x] Initial JSON-RPC over stdio out-of-process execution engine
- [x] Basic Plugin SDK for Go developers
- [x] Reference `plugin-os` implementation

### Phase 1: MVP Reliability & Protocol
- [x] Mandatory Handshake negotiation
- [x] API Version verification (v1 enforcement)
- [x] Multi-path plugin discovery (`./plugins`, `~/.zt/plugins`, `/usr/local/lib/`)
- [x] Per-check context timeouts (default 30s)
- [x] Standardized RPC error propagation

### Phase 2: Enhanced Reporting & Concurrency
- [x] Parallel execution of checks via worker pool
- [x] Structured Terminal summary (pretty-printing findings)
- [x] JSON output reporter
- [x] SARIF export for CI/CD integration
- [ ] Global configuration management (`~/.ztvs/config.yaml`)

### Phase 3: Zero Trust & Security
- [x] Capability-based permission declaration in plugin manifests
- [x] Mandatory capability enforcement by the Host
- [x] Plugin binary integrity (SHA-256) validation
- [x] Basic execution sandboxing (process isolation)

### Phase 4: Scaling & Distribution
- [x] Centralized Plugin Registry (`zt plugin install`)
- [x] Periodic background audit mode (`zt agent`)
- [x] Host configuration management (`~/.ztvs/config.yaml`)
- [x] Automated plugin manifest checksumming
- [x] Policy-as-Code engine (severity thresholds and filtering)
- [x] Cloud-native metadata discovery

---

## 📈 Next Priority: Future Roadmap
The immediate focus is expanding the plugin ecosystem and integrating with cloud-native security orchestration platforms.
