# RFC-002: Plugin SDK & JSON-RPC Protocol Specification

**Project:** Zero Trust Vulnerability Scanner (ZTVS)
**Status:** Draft v1
**Depends on:** RFC-001 Core Architecture
**Language:** Go
**Date:** 2026-04-01

---

# 1. Abstract

This RFC defines the **plugin developer SDK**, **wire protocol**, **process lifecycle**, and **JSON-RPC schemas** for third-party and first-party plugins in ZTVS.

This document standardizes:

* plugin executable contract
* manifest format
* host ↔ plugin communication
* handshake negotiation
* check execution protocol
* capability declaration
* error model
* version compatibility
* Go SDK interfaces
* testing contract

This RFC is normative for all plugins.

---

# 2. Goals

The plugin system MUST provide:

* stable API compatibility
* language-independent plugin support
* process isolation
* version negotiation
* deterministic results
* strict timeout control
* future marketplace compatibility

---

# 3. Non-Goals

This RFC does NOT define:

* plugin signing marketplace workflow
* remote plugin fetching
* cloud-hosted plugins
* policy DSL
* report output formats

Those are covered in later RFCs.

---

# 4. Design Principles

---

## 4.1 Language Agnostic

Plugins SHOULD be writable in any language.

Examples:

* Go
* Rust
* Python
* Node.js

The only hard requirement is protocol compliance.

---

## 4.2 Process Isolation

Plugins MUST run as separate OS processes.

In-process shared libraries are explicitly prohibited.

---

## 4.3 Backward Compatibility

The wire protocol MUST use semantic versioning.

Major version mismatches MUST fail handshake.

---

# 5. Transport Layer

Transport SHALL be:

**JSON-RPC 2.0 over stdin/stdout**

This is the normative transport.

Reference:

* host writes JSON requests to plugin stdin
* plugin writes JSON responses to stdout
* stderr reserved for diagnostics/logging

---

# 6. Process Lifecycle

```text id="rfc002-lifecycle"
spawn
→ handshake
→ manifest validation
→ capability negotiation
→ enumerate checks
→ execute checks
→ collect findings
→ terminate
```

---

# 7. Plugin Binary Contract

A plugin MUST be an executable binary.

Example:

```bash id="rfc002-bin"
plugin-os
plugin-network
plugin-package
```

The host launches:

```bash id="rfc002-launch"
plugin-os --rpc
```

This flag is REQUIRED.

---

# 8. Manifest Specification

Each plugin MUST include a manifest file.

Default name:

```text id="rfc002-manifest"
plugin.yaml
```

---

## 8.1 Schema

```yaml id="rfc002-manifest-example"
name: plugin-os
version: 1.2.0
sdk_version: 1.0.0
api_version: 1
author: zt-security
platforms:
  - linux
  - darwin
  - windows
capabilities:
  - read_files
  - execute_commands
checks:
  - ssh_config
  - firewall_status
entrypoint: ./plugin-os
```

---

## 8.2 Required Fields

Required:

* name
* version
* api_version
* entrypoint
* checks

---

# 9. JSON-RPC Envelope

All messages MUST follow JSON-RPC 2.0.

---

## 9.1 Request

```json id="rfc002-jsonrpc-request"
{
  "jsonrpc": "2.0",
  "id": "req-001",
  "method": "handshake",
  "params": {}
}
```

---

## 9.2 Response

```json id="rfc002-jsonrpc-response"
{
  "jsonrpc": "2.0",
  "id": "req-001",
  "result": {}
}
```

---

## 9.3 Error

```json id="rfc002-jsonrpc-error"
{
  "jsonrpc": "2.0",
  "id": "req-001",
  "error": {
    "code": 4001,
    "message": "unsupported api version"
  }
}
```

---

# 10. Protocol Methods

Defined methods:

* handshake
* health
* enumerate
* run_check
* self_test
* shutdown

---

# 11. Handshake Method

Handshake MUST be the first call.

---

## Request

```json id="rfc002-handshake-request"
{
  "method": "handshake",
  "params": {
    "host_version": "1.0.0",
    "api_version": 1
  }
}
```

---

## Response

```json id="rfc002-handshake-response"
{
  "name": "plugin-os",
  "version": "1.0.0",
  "api_version": 1,
  "sdk_version": "1.0.0",
  "checks_supported": [
    "ssh_config",
    "firewall_status"
  ]
}
```

---

## Failure Rules

Handshake MUST fail when:

* api version mismatch
* invalid manifest
* checksum mismatch
* missing capabilities declaration

---

# 12. Health Check Method

Used for host liveness validation.

---

## Request

```json id="rfc002-health-request"
{
  "method": "health"
}
```

---

## Response

```json id="rfc002-health-response"
{
  "status": "healthy"
}
```

---

# 13. Enumerate Method

Returns supported checks.

---

## Request

```json id="rfc002-enumerate-request"
{
  "method": "enumerate"
}
```

---

## Response

```json id="rfc002-enumerate-response"
{
  "checks": [
    {
      "id": "ssh_config",
      "name": "SSH Configuration Check",
      "severity_default": "high",
      "platforms": ["linux", "darwin"]
    }
  ]
}
```

---

# 14. Run Check Method

This is the primary execution RPC.

---

## Request

```json id="rfc002-run-request"
{
  "method": "run_check",
  "params": {
    "check_id": "ssh_config",
    "target": {
      "type": "local_host",
      "os": "linux"
    },
    "context": {
      "timeout_seconds": 30
    }
  }
}
```

---

## Response

```json id="rfc002-run-response"
{
  "status": "fail",
  "finding": {
    "id": "F-001",
    "severity": "high",
    "title": "Root login enabled",
    "description": "PermitRootLogin is enabled",
    "evidence": {
      "file": "/etc/ssh/sshd_config",
      "value": "PermitRootLogin yes"
    },
    "remediation": "Set PermitRootLogin no"
  }
}
```

---

# 15. Finding Schema

This MUST conform to RFC-001.

---

## JSON Schema

```json id="rfc002-finding"
{
  "id": "string",
  "check_id": "string",
  "severity": "critical|high|medium|low|info",
  "title": "string",
  "description": "string",
  "evidence": {},
  "remediation": "string"
}
```

---

# 16. Error Codes

Reserved codes:

| Code | Meaning                 |
| ---- | ----------------------- |
| 4001 | unsupported api version |
| 4002 | invalid check id        |
| 4003 | capability denied       |
| 4004 | timeout                 |
| 5000 | internal plugin error   |

---

# 17. Capability Contract

Capabilities MUST be declared both in manifest and runtime metadata.

---

## Request Example

```json id="rfc002-capability"
{
  "required_capabilities": [
    "read_files",
    "execute_commands"
  ]
}
```

---

## Enforcement Rule

Host SHALL reject execution if requested capability is denied.

---

# 18. Timeout Semantics

Each run_check MUST receive a timeout.

Plugins MUST respect cancellation.

Go SDK SHALL expose context propagation.

Example:

```go id="rfc002-timeout-go"
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()
```

---

# 19. Go SDK Specification

Package:

```text id="rfc002-sdk-path"
/pkg/sdk
```

---

## 19.1 Plugin Interface

```go id="rfc002-go-plugin"
type Plugin interface {
    Metadata() Metadata
    Checks() []Check
}
```

---

## 19.2 Check Interface

```go id="rfc002-go-check"
type Check interface {
    ID() string
    Name() string
    Platforms() []string
    RequiredCapabilities() []Capability
    Run(ctx context.Context, host HostContext) (*Finding, error)
}
```

---

## 19.3 Metadata

```go id="rfc002-go-metadata"
type Metadata struct {
    Name       string
    Version    string
    APIVersion int
}
```

---

# 20. SDK Bootstrap Example

```go id="rfc002-bootstrap"
func main() {
    sdk.Run(&OSPlugin{})
}
```

---

# 21. Example Plugin

```go id="rfc002-example-plugin"
type SSHCheck struct{}

func (c *SSHCheck) ID() string {
    return "ssh_config"
}

func (c *SSHCheck) Name() string {
    return "SSH Configuration"
}

func (c *SSHCheck) Platforms() []string {
    return []string{"linux", "darwin"}
}

func (c *SSHCheck) RequiredCapabilities() []sdk.Capability {
    return []sdk.Capability{
        sdk.CapReadFiles,
    }
}

func (c *SSHCheck) Run(
    ctx context.Context,
    host sdk.HostContext,
) (*sdk.Finding, error) {
    return &sdk.Finding{
        ID: "F-001",
        Severity: sdk.High,
        Title: "Root login enabled",
    }, nil
}
```

---

# 22. Plugin Test Contract

Each plugin MUST implement:

```text id="rfc002-test-contract"
self_test
```

This SHALL verify:

* dependencies available
* OS compatibility
* required binaries present
* permission access possible

---

# 23. Shutdown Contract

Plugins MUST gracefully terminate.

---

## Request

```json id="rfc002-shutdown"
{
  "method": "shutdown"
}
```

---

# 24. Logging

stderr is reserved for logs.

Allowed levels:

* debug
* info
* warn
* error

stdout MUST remain protocol-only.

---

# 25. Compatibility Matrix

| Host | Plugin | Result     |
| ---- | ------ | ---------- |
| v1.x | v1.x   | compatible |
| v1.x | v2.x   | reject     |
| v2.x | v1.x   | reject     |

---

# 26. Security Requirements

Plugins MUST NOT:

* open network sockets unless capability granted
* write outside temp directory
* spawn child processes without declaration
* persist credentials

---

# 27. Future RFCs

Next recommended:

* RFC-003 Policy-as-Code
* RFC-004 Signed Plugin Registry
* RFC-005 Remote Agent Protocol

---

# Final Recommendation

This RFC gives you a production-grade plugin ecosystem contract.

It is stable enough to begin implementation immediately for:

* SDK package
* host process manager
* first-party plugins
* third-party developer docs
