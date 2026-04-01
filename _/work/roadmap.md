# ZTVS Strategic Roadmap

This document provides a high-level overview of the Zero Trust Vulnerability Scanner (ZTVS) development milestones and their current status.

## 📊 Phase Summary

| Phase | Title | Focus | Status |
| :--- | :--- | :--- | :--- |
| **0** | **Foundation** | Core architecture & directory layout | ✅ Completed |
| **1** | **Reliability** | Handshake, Discovery, & Timeout support | ✅ Completed |
| **2** | **Reporting** | Concurrency & Structured (JSON/SARIF) output | 📅 Planned |
| **3** | **Security** | Capability enforcement & Plugin sandboxing | 📅 Planned |
| **4** | **Platform** | Remote registry & Agent mode | 📅 Planned |

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
- [ ] Parallel execution of checks via worker pool
- [ ] Structured Terminal summary (pretty-printing findings)
- [ ] JSON output reporter
- [ ] SARIF export for CI/CD integration
- [ ] Global configuration management (`~/.zt/config.yaml`)

### Phase 3: Zero Trust & Security
- [ ] Capability-based permission declaration in plugin manifests
- [ ] Mandatory capability enforcement by the Host
- [ ] Plugin binary signing and SHA-256 verification
- [ ] Basic OS-level sandboxing for plugin processes

### Phase 4: Scaling & Distribution
- [ ] Centralized Plugin Registry (`zt plugin install`)
- [ ] Agent mode for periodic background scanning
- [ ] Policy-as-Code engine (severity thresholds and filtering)
- [ ] Cloud-native metadata discovery

---

## 📈 Next Priority: Concurrency & Reports (Phase 2)
The immediate focus is improving scanning performance by running checks in parallel and providing machine-readable outputs for automation.
