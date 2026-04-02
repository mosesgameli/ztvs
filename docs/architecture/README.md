🔙 [Back to Home](../README.md)

# ZTVS Architecture Overview

Zero Trust Vulnerability Scanner (ZTVS) is built with a host-and-plugin architecture where the core engine manages the lifecycle, execution, and reporting of isolated security checks.

## Core Principles

1.  **Process Isolation**: All plugins execute as separate child processes. A crash or compromise in one plugin does not affect the host or other plugins.
2.  **Least Privilege**: Every plugin must declare its required capabilities (e.g., `read_files`, `execute_commands`) in a manifest. The host enforces these through a central policy engine.
3.  **JSON-RPC Protocol**: Communication between the host and plugins uses a standard JSON-RPC 2.0 protocol over `stdin`/`stdout`.
4.  **Language Agnostic**: Because the protocol is based on standard stdio and JSON, plugins can be written in any language (Go, Python, Rust, etc.).

## Component Diagram

### ZTVS CLI Core (`internal/cli`)
-   **Config Loader**: Reads global settings and policies from `~/.ztvs/config.yaml`.
-   **Engine**: Orchestrates the scanning process, manages the worker pool, and aggregates findings.
-   **Plugin Host**: Handles process spawning, handshake negotiation, and capability validation.
-   **Reporting Engine**: Normalizes findings and emits reports in Terminal, JSON, or SARIF formats.

### Plugin Infrastructure
-   **Manifest (`plugin.yaml`)**: Declares the plugin's metadata, API version, and required system capabilities (e.g., `network_access`).
-   **Runner Subsystem**: ZTVS uses a polyglot runner to execute plugins natively. For `Go` and `Rust`, binary execution is used. For `Python` and `Node.js`, the runner dynamically binds to the host runtime interpreter (`python3` / `node`).
-   **SDKs**: First-party Go and Python libraries that handle JSON-RPC boilerplate. These are maintained in separate repositories:
    -   [ztvs-sdk-go](https://github.com/mosesgameli/ztvs-sdk-go)
    -   [ztvs-sdk-python](https://github.com/mosesgameli/ztvs-sdk-python)

## Plugin Registry and Lifecycle
ZTVS relies on `plugins.ztvs.dev` for managed plugin distribution:
1.  **Verification**: Plugins are downloaded, securely built, and matched against SHA-256 binary checksums.
2.  **Atomic Swaps**: When executing `zt plugin update`, the host downloads and compiles the new payload to a temporary cache, performing a fast atomic rename to prevent execution failures during live scans via locking.

## Execution Flow

1.  **Discovery**: Host searches for plugin binaries in standard locations (e.g., `./plugins`, `~/.ztvs/plugins`).
2.  **Integrity Check**: Host verifies the SHA-256 checksum of each plugin binary.
3.  **Handshake**: Host spawns the plugin and sends a `handshake` request to verify the API version and supported checks.
4.  **Enforcement**: Host checks the plugin's requested capabilities against the local security policy.
5.  **Scan**: Engine parallelizes the execution of `run_check` calls across all plugins.
6.  **Report**: Findings are aggregated and emitted in the requested format.
