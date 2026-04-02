# Go Coding Style Guide for ZTVS

This document outlines the coding standards and project topology used within the Zero Trust Vulnerability Scanner (ZTVS) repository.

## 1. Project Topology

We follow standard Go practices for repository organization:

-   **/cmd/zt**: The main entry point for the ZTVS CLI. Keep this directory minimal; it should only handle command-line flag parsing and invoke the core engine.
-   **/internal/**: Contains logic private to the ZTVS host engine. This includes the plugin runner, policy engine, and scheduler. Code in `internal/` cannot be imported by other projects.
-   **/pkg/**: Contains code intended for public use. This includes JSON-RPC 2.0 type definitions, SDK interfaces, and registry schemas.
-   **/registry/**: Local storage for plugin manifests and trust verification logic.

## 2. Naming Conventions

-   **Interfaces**: Keep interfaces small (e.g., `Check`, `Reporter`). Use the `-er` suffix where appropriate.
-   **Variables**: Use short, descriptive names. For local variables, `ctx`, `err`, `req`, `cfg` are preferred.
-   **Exported Types**: Always provide a doc comment for exported functions, types, and constants.

## 3. Concurrency & Context

ZTVS is a high-concurrency engine. Follow these rules for safe execution:

### Worker Pools
Use the worker pool pattern for parallel scanning. Individual security checks should be independent and shouldn't share mutable state.

### Context Propagation
Always propagate `context.Context` through your call stack. This ensures that:
-   Timeouts from the CLI (e.g., `--timeout 30s`) are respected.
-   Plugin processes are correctly terminated if the host is interrupted (`SIGINT`).
-   Goroutines are cleaned up properly.

### Mutexes
All shared resources (e.g., the final `Finding` list, `Registry` cache) must be protected by a `sync.Mutex`. Hold locks for the shortest time possible.

## 4. Error Handling

### Contextual Errors
Avoid returning "naked" errors. Provide context using `fmt.Errorf`:
```go
if err := r.loadIndex(); err != nil {
    return fmt.Errorf("failed to load registry index: %w", err)
}
```

### JSON-RPC Errors
When communicating with plugins, use the established JSON-RPC error codes:

| Code | Meaning |
| :--- | :--- |
| `-32700` | Parse error |
| `-32601` | Method not found |
| `4001` | Version mismatch |
| `4002` | Check not found |

## 5. Security & Isolation

-   **Zero Trust**: Never trust output from a plugin. Always validate JSON schemas and sanitize evidence before reporting.
-   **Path Sanitization**: Ensure that plugins cannot bypass capability checks by using relative paths (`../`).
-   **Memory Management**: Be mindful of large plugin outputs (base64 evidence). Use buffered readers or limits where appropriate.
