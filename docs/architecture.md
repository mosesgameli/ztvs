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

### Plugin Architecture
-   **Manifest (`plugin.yaml`)**: Declares the plugin's metadata, supported checks, and required capabilities.
-   **SDK (`pkg/sdk`)**: A first-party Go library that handles the JSON-RPC boilerplate for plugin developers.

## Execution Flow

1.  **Discovery**: Host searches for plugin binaries in standard locations (e.g., `./plugins`, `~/.ztvs/plugins`).
2.  **Integrity Check**: Host verifies the SHA-256 checksum of each plugin binary.
3.  **Handshake**: Host spawns the plugin and sends a `handshake` request to verify the API version and supported checks.
4.  **Enforcement**: Host checks the plugin's requested capabilities against the local security policy.
5.  **Scan**: Engine parallelizes the execution of `run_check` calls across all plugins.
6.  **Report**: Findings are aggregated and emitted in the requested format.
