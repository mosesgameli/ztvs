 
trigger: always_on
 

# ZTVS Project Rules & Standards

These rules are NORMATIVE for all development on the Zero Trust Vulnerability Scanner (ZTVS).

## Core Principles
1.  **Go First**: The core engine and SDK MUST be written in Go (target version: 1.24+).
2.  **Zero Trust Execution**: Plugins MUST be executed as separate OS processes and communicate via JSON-RPC 2.0 over stdio.
3.  **No In-Process Plugins**: Shared libraries (`.so`, `.dll`) are prohibited for plugin execution.
4.  **Capability-based Security**: Plugins MUST declare required capabilities (e.g., `read_files`) in their manifests.

## Coding Standards
1.  **Licensing**:
    - Every new `.go`, `.sh`, and `.ps1` file MUST start with the official Apache License 2.0 header.
2.  **Internal vs Pkg**: 
    - Use `internal/` for logic private to the ZTVS host (e.g., `pluginhost`, `engine`).
    - Use `pkg/` for code intended for external consumption (e.g., `sdk`, `rpc` types).
3.  **Error Handling**:
    - Propagate errors with context using `fmt.Errorf("context: %w", err)`.
4.  **Documentation**:
    - All exported symbols in `pkg/` and core `internal/` packages MUST have Go docstrings.
5.  **Concurrency**:
    - Use `sync.WaitGroup` and worker pools for parallel scanning.

## Testing
1.  **Coverage**: Unit test coverage SHOULD be ≥95%. PRs MUST include tests for new logic.
2.  **Binary Testing**: Always verify that `make build` succeeds before committing.
3.  **Protocol Testing**: Use simulated plugins to verify handshake and timeout logic.
