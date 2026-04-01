# ZTVS Phased Delivery Plan

This document outlines the detailed roadmap for the Zero Trust Vulnerability Scanner (ZTVS) project, from the current "Skeleton" state to a production-grade security tool.

## Phase 0: Foundation (COMPLETED)
- [x] Repository initialization (`go mod init`).
- [x] Directory structure layout (following `skeleton.md`).
- [x] CLI entrypoint & Engine orchestration skeleton.
- [x] Out-of-process JSON-RPC plugin execution model.
- [x] Plugin SDK for Go developers.
- [x] First-party reference plugin (`plugin-os`).

## Phase 1: MVP Reliability & Protocol (Next Priority)
Goal: Stabilize the communication protocol and ensure resource safety.

### P0 Tasks
- **Handshake Negotiation**: Enforce version checking (API v1) during plugin startup.
- **Plugin Discovery**: Implement multi-path discovery (`/usr/local/lib/zt/plugins`, `~/.zt/plugins`).
- **Timeout Support**: Connect `context.WithTimeout` from the Host to the plugin lifecycle.
- **Detailed Error Handling**: Properly propagate RPC error codes (e.g., 4001, 4002) to the CLI.

### P1 Tasks
- **Logger Integration**: Capture plugin `stderr` and route it to the host's log aggregator.
- **Result Normalization**: Ensure every finding follows the RFC-001 schema exactly.

## Phase 2: Enhanced Reporting & Concurrency
Goal: improve performance and make the output useful for integration.

### P0 Tasks
- **Concurrency Worker Pool**: Run multiple plugin checks in parallel (default: CPU cores * 2).
- **JSON Reporter**: Machine-readable scan results for automation.
- **Terminal Reporter**: Human-friendly, color-coded output for local usage.

### P1 Tasks
- **SARIF Integration**: Export findings in Static Analysis Results Interchange Format for GitHub/GitLab integration.
- **Basic Configuration**: Load `~/.zt/config.yaml` to manage enabled plugins and thresholds.

## Phase 3: Security & Trust
Goal: Implement the "Zero Trust" part of the project name.

### P0 Tasks
- **Capability Declaration**: Plugins must declare permissions in `plugin.yaml`.
- **Capability Enforcement**: Host rejects plugins requesting permissions that violate local policy.

### P1 Tasks
- **Signed Plugins**: Verify SHA-256 checksums and digital signatures before execution.
- **Process Sandboxing**: (Linux-first) Use `seccomp` or `namespaces` to restrict plugin OS access.

## Phase 4: Scaling & Platform
Goal: Support complex environments and large fleets.

### P0 Tasks
- **Policy-as-Code**: Filter findings based on severity or type using a policy engine (e.g., OPA/Rego).
- **Remote Registry**: Create `zt plugin install <name>` to fetch signed plugins from a central repo.

### P1 Tasks
- **Agent Mode**: Run ZTVS as a background service with periodic scan intervals.
- **Cloud Connectors**: Basic metadata discovery for AWS/GCP/Azure environments.

---

## Technical Debt & Maintenance
- **Unit Test Suite**: Achieve >80% coverage for the `internal/engine` and `pkg/sdk`.
- **Documentation**: auto-generate plugin documentation from manifests.
- **Cross-Compilation**: Setup CI/CD runners for Windows (amd64), Linux (amd64/arm64), and macOS (arm64).
