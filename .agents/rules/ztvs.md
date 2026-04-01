---
trigger: always_on
---

# ZTVS Project Rules & Standards

These rules are NORMATIVE for all development on the Zero Trust Vulnerability Scanner (ZTVS).

## Core Principles
1.  **Go First**: The core engine and SDK MUST be written in Go (target version: 1.24+).
2.  **Zero Trust Execution**: Plugins MUST be executed as separate OS processes and communicate via JSON-RPC 2.0 over stdio.
3.  **No In-Process Plugins**: Shared libraries (`.so`, `.dll`) are prohibited for plugin execution.
4.  **Capability-based Security**: Plugins MUST declare required capabilities (e.g., `read_files`) in their manifests.

## Coding Standards
1.  **Internal vs Pkg**: 
    - Use `internal/` for logic private to the ZTVS host (e.g., `pluginhost`, `engine`).
    - Use `pkg/` for code intended for external consumption (e.g., `sdk`, `rpc` types).
2.  **Error Handling**:
    - Propagate errors with context using `fmt.Errorf("context: %w", err)`.
    - Use JSON-RPC error codes (4001: Version mismatch, 4002: Check not found) for plugin communication.
3.  **Concurrency**:
    - Use `sync.WaitGroup` and worker pools for parallel scanning.
    - All shared resources (e.g., Reporters) MUST be accessed through mutexes.

## Testing
1.  **Binary Testing**: Always verify that `make build` succeeds before committing.
2.  **Protocol Testing**: Use simulated plugins to verify handshake and timeout logic.
