🔙 [Back to Home](../README.md)

# ZTVS JSON-RPC Protocol Specification

Zero Trust Vulnerability Scanner (ZTVS) uses JSON-RPC 2.0 over standard streams (`stdin`/`stdout`) for communication between the host and plugins.

-   **Host**: Writes JSON requests to the plugin's `stdin`.
-   **Plugin**: Writes JSON responses to its `stdout`.
-   **Diagnostics**: All non-protocol logging and errors should be written to `stderr`.

## JSON-RPC Envelope

All messages must follow the JSON-RPC 2.0 specification.

### Request

```json
{
  "jsonrpc": "2.0",
  "id": "req-001",
  "method": "handshake",
  "params": {}
}
```

### Response

```json
{
  "jsonrpc": "2.0",
  "id": "req-001",
  "result": {}
}
```

### Error

```json
{
  "jsonrpc": "2.0",
  "id": "req-001",
  "error": {
    "code": 4001,
    "message": "unsupported api version"
  }
}
```

## Protocol Methods

### 1. `handshake`

Handshake **must** be the first call between the host and the plugin.

#### Request Params
- `host_version` (string): The version of the ZTVS host.
- `api_version` (integer): The supported protocol version (currently `1`).

#### Result
- `name` (string): Plugin identifier.
- `version` (string): Plugin version.
- `api_version` (integer): Plugin's supported protocol version.
- `checks_supported` (array of strings): List of check IDs the plugin can run.

---

### 2. `run_check`

Execution of a specific security check.

#### Request Params
- `check_id` (string): The ID of the check to run (e.g., `ssh_config`).

#### Result
- `status` (string): Status of the check (`pass`, `fail`, `error`, `skipped`).
- `finding` (object, nullable): The security finding details. If the check completes successfully without discovering any vulnerabilities, the plugin **must** return `null` for the finding instead of crashing or returning an empty object.

#### Finding Schema
```json
{
  "id": "F-001",
  "check_id": "ssh_config",
  "severity": "high",
  "title": "Root login enabled",
  "description": "PermitRootLogin is enabled in sshd_config",
  "evidence": {
    "file": "/etc/ssh/sshd_config",
    "value": "PermitRootLogin yes"
  },
  "remediation": "Set PermitRootLogin no in sshd_config"
}
```

---

### 3. `shutdown` (Future)

Plugin should gracefully terminate.

---

## Error Codes

| Code | Meaning                 | Description                                      |
| :--- | :---------------------- | :----------------------------------------------- |
| 4001 | unsupported api version | Host and plugin protocol versions are incompatible. |
| 4002 | invalid check id        | The requested check ID does not exist in the plugin. |
| 4003 | capability denied       | Plugin attempted a check without required permissions. |
| 5000 | internal plugin error   | Unexpected unhandled error inside the plugin.    |
