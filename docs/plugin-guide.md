🔙 [Back to Home](../README.md)

# ZTVS Plugin Development Guide

ZTVS is designed to be truly language agnostic. Any language that can execute on the target platform and communicate via the JSON-RPC status standard is a first-class plugin developer.

## Core Plugin Requirements

1.  **Binary Executable**: The plugin must be an executable binary or script.
2.  **JSON-RPC Stdio**: The plugin must read requests from `stdin` and write responses to `stdout`.
3.  **No Stdout Logging**: `stdout` is reserved for protocol messages. All logs should be written to `stderr`.

## Implementation Checklist

-   [ ]   **Manifest (`plugin.yaml`)**: Every plugin must include a manifest declaring its name, version, and capabilities.
-   [ ]   **Handshake**: Implement the `handshake` method. This is the first call made by the host.
-   [ ]   **Check IDs**: Define unique check IDs that mapped to security validations.
-   [ ]   **Run Check**: Implement the `run_check` method to execute security logic and return a passing or failing status.
-   [ ]   **Capabilities**: Only perform actions for which you have declared capabilities. Users must explicitly whitelist required capabilities like `network_access` or `execute_commands` inside their global `~/.ztvs/config.yaml` before your plugin can function.

## Example Handshake (Raw)

**Host Request:**
```json
{
  "jsonrpc": "2.0",
  "id": "handshake",
  "method": "handshake",
  "params": {
    "host_version": "1.0.0",
    "api_version": 1
  }
}
```

**Plugin Response:**
```json
{
  "jsonrpc": "2.0",
  "id": "handshake",
  "result": {
    "name": "my-plugin",
    "version": "1.0.0",
    "api_version": 1,
    "checks_supported": ["ssh_config", "port_scan"]
  }
}
```

## Tips for Success

-   **Environment Agnostic**: Avoid hardcoded paths; use environment variables or capabilities to detect system resources.
-   **Structured Evidence**: In the `evidence` field, return structured data (maps/objects) that the host can use for reporting.
-   **Graceful Errors**: If a check fails due to an unexpected error, return a standard JSON-RPC error rather than a panic.
-   **Polyglot Examples**: Check our [examples/](../examples/) directory for reference implementations in Python, Rust, Node.js, and Java.
