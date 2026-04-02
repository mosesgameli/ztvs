# Zero Trust Vulnerability Scanner (ZTVS): Product Overview

> **For Technology Leadership** | Organizational Security Posture & Adoption Brief

## Executive Summary

The Zero Trust Vulnerability Scanner (ZTVS) is an open-source, enterprise-grade host security scanning platform built on a Zero Trust execution model. It continuously audits every host in your organization — developer workstations, CI/CD runners, and production servers alike — detecting configuration weaknesses, supply chain compromises, and active indicators of compromise (IOCs) before they escalate into incidents.

Unlike legacy monolithic security agents, ZTVS treats every security check as an untrusted, isolated process. No single vulnerability in a scanning rule can compromise the scanner itself or pivot laterally across your infrastructure. This design is not a feature — it is a foundational architectural constraint.

ZTVS is production-ready today, with native support for **Windows, macOS, and Linux**, one-command deployment, and enterprise output formats for immediate SIEM and CI/CD integration.

## The Threat Landscape That Motivated ZTVS

Software supply chain attacks have become the defining threat of the mid-2020s. The 2026 Axios Supply Chain RAT attack — where a widely trusted JavaScript HTTP library was hijacked to deliver a Remote Access Trojan to millions of developer machines and production servers — demonstrated a new class of risk that traditional endpoint agents are structurally unable to detect:

- Standard antivirus cannot distinguish between a legitimate npm package and a maliciously modified one
- SAST/DAST tools scan code at build time but miss runtime host-level IOCs planted after deployment
- No existing enterprise agent had a detection rule for this specific attack within 72 hours of disclosure

ZTVS ships a first-party `plugin-axios-mitigation` plugin that actively detects all three attack vectors of this class: compromised dependency manifests (`axios_dependency_audit`), host-level IOC artifacts (`axios_host_ioc`), and C2 callback network activity (`axios_c2_network`). It ships the response, not just the detection.

## Core Value Proposition for Technology Leadership

| Concern | ZTVS Response |
|---|---|
| **Supply chain integrity** | Cryptographic verification of every installed dependency and plugin binary (SHA-256) |
| **Developer machine posture** | Continuous host auditing with on-install automation hooks across all major package managers |
| **Regulatory & compliance alignment** | SARIF 2.1.0 output maps directly to CISA Zero Trust Maturity Model pillars |
| **Vendor lock-in risk** | MIT/Apache 2.0 licensed; fully open source; no proprietary agents or cloud dependencies |
| **Custom security rules** | First-party SDK enables internal teams to write organization-specific checks in Go or Python |
| **SIEM integration** | Structured JSON and SARIF output integrates with Splunk, Elastic, Microsoft Sentinel, and GitHub Advanced Security |
| **Operational overhead** | One-command installation; automatic plugin updates via a managed remote registry |

## Architecture: Why Zero Trust Matters Here

Traditional security agents run as monolithic processes. A single crashing rule, a vulnerable parser, or a compromised detection module can take down the entire agent — or worse, be exploited as a foothold. ZTVS enforces a fundamentally different model:

**Every security check runs as a separate, isolated child process.** The core engine communicates with each check over a restricted JSON-RPC 2.0 channel on standard I/O. Each check must declare, in a cryptographically verified manifest (`plugin.yaml`), the exact capabilities it requires:

```
capabilities:
  - read_files
  - execute_commands
  - network_access
```

The host engine enforces these declarations against a central policy file (`~/.ztvs/config.yaml`). A check that attempts to access the network without declaring `network_access` is denied at the OS call level. A check that crashes is isolated to its own process and does not interrupt the scan or affect other checks.

This architecture satisfies the **Least Privilege** and **Assume Breach** pillars of the NIST SP 800-207 Zero Trust standard — applied to the security tooling layer itself.

## Plugin Ecosystem: Security Coverage Out of the Box

ZTVS ships with three production-quality first-party plugins covering the most critical org-wide risk surfaces:

### `plugin-os` — Host Operating System Hardening
A high-performance compiled Go plugin that audits core OS security configuration:

- **SSH hardening**: Detects permissive `sshd_config` settings (root login, password auth, empty passwords)
- **Password policy enforcement**: Validates system-level password strength and expiration controls
- **User account auditing**: Flags unauthorized accounts, privilege escalation vectors, and stale sessions

This plugin requires only `read_files` and `execute_commands` capabilities — it performs no network activity and writes nothing to disk.

### `plugin-axios-mitigation` — Supply Chain Attack Response
A three-vector Go plugin built as a direct, incident-response-grade answer to the 2026 Axios supply chain compromise:

- **`axios_dependency_audit`**: Scans all Node.js dependency manifests and lockfiles for compromised axios versions and the malicious `plain-crypto-js@4.2.1` payload
- **`axios_host_ioc`**: Scans the host filesystem for IOC artifacts (known malicious file paths, suspicious executables) left by the RAT dropper
- **`axios_c2_network`**: Probes for active outbound connections to known C2 infrastructure associated with the 2026 campaign

### `plugin-axios-github-scan` — GitHub Organization Exposure Audit
A Python-based plugin that extends scanning beyond the local host to your entire GitHub organization:

- **`axios-github-inventory-audit`**: Inventories every repository in your GitHub org, identifies those using vulnerable axios versions, correlates with recent CI activity to assess actual exposure, and detects the `plain-crypto-js` payload in lockfiles
- **`github-deployment-audit`**: Monitors deployment activity across production, staging, and other environments; flags failed deployment statuses and provides a deployment frequency summary for security review

This plugin requires a GitHub Personal Access Token with `repo`, `read:org`, and `actions:read` scopes.

## Extensibility: Build Your Own Security Rules

ZTVS is designed to democratize security rule authorship. Any team in your organization — security, platform, application — can write a custom check in **fewer than 25 lines of code**. First-party SDKs are available for the two most common internal language stacks:

**Go SDK** (`github.com/mosesgameli/ztvs-sdk-go`): Produces a natively compiled binary. Ideal for performance-sensitive checks, low-level system inspection, and checks that must run on endpoints without a language runtime.

**Python SDK** (`github.com/mosesgameli/ztvs-sdk-python`): Supports async execution and Pydantic-validated data models. Ideal for checks that call external APIs, parse structured data, or leverage the Python security ecosystem (e.g., `bandit`, `semgrep` integrations).

Both SDKs abstract the JSON-RPC protocol entirely. A developer implements one interface (`Run() → Finding`), and the SDK handles handshake, capability negotiation, error handling, and timeouts. Plugins that do not return within **30 seconds** are automatically killed by the host engine.

The plugin distribution model uses a cryptographically verified remote registry (`plugins.ztvs.dev`). Publishing a new internal check requires only a SHA-256 checksum and a `plugin.yaml` manifest — no proprietary signing infrastructure.

## Enterprise Integration and Reporting

ZTVS produces findings in three output formats, configurable at runtime:

| Format | Use Case |
|---|---|
| **Terminal (pretty-print)** | Developer-facing, interactive scanning and triage |
| **JSON** | Machine-readable pipeline integration, custom dashboards |
| **SARIF 2.1.0** | Native upload to GitHub Advanced Security, compatibility with most SAST/DAST platforms and enterprise SIEMs |

```bash
# Output SARIF for GitHub Advanced Security ingestion
zt --format sarif scan > report.sarif
```

SARIF is the industry-standard static analysis result format. ZTVS findings in SARIF can be ingested by:
- GitHub Advanced Security (code scanning alerts)
- Microsoft Defender for DevOps
- Splunk SOAR
- Elastic Security
- Azure DevOps

## Deployment Models for Org-Wide Adoption

ZTVS supports two distinct operational modes, enabling adoption across different org tiers:

### On-Demand Mode (`zt scan`)
A single, blocking scan of the local host. Suitable for:
- Developer workstations (triggered manually or via CI/CD job step)
- Pre-commit or pre-push hooks
- One-time incident response sweeps

### Continuous Agent Mode (`zt agent`)
A lightweight background daemon that runs as an isolated scheduled process. Suitable for:
- Production servers and long-lived infrastructure
- CI/CD build runners
- Shared developer environments

### Automatic Scanning on Package Install
ZTVS can be configured to trigger `zt scan` automatically whenever a developer installs packages, intercepting supply chain threats the moment new dependencies land on a machine. This hooks integrate with every major package manager:

| Package Manager | Trigger Events |
|---|---|
| `npm` / `pnpm` / `yarn` | `install`, `add`, `ci` |
| `pip` / `pip3` | `install` |
| `brew` | `install`, `upgrade` |
| `cargo` | `install`, `add` |
| `go` | `get`, `install` |

Shell configuration snippets for Bash, Zsh, and PowerShell are provided in the main repository README.

## Installation and Rollout

### Single-Host Installation (Developer Onboarding)

**Linux & macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/mosesgameli/ztvs/main/scripts/install.sh | sh
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/mosesgameli/ztvs/main/scripts/install.ps1 | iex
```

Both scripts install the `zt` binary and seed the default first-party plugins in a single step. A pinnable version flag (`ZTVS_VERSION=v1.0.0`) supports controlled rollouts through MDM or infrastructure-as-code tooling.

### Org-Wide Rollout Considerations

| Deployment Method | Recommended For |
|---|---|
| MDM (Jamf, Intune, Workspace ONE) | Developer workstations, macOS/Windows fleets |
| Ansible / Chef / Puppet | Linux server fleets, CI/CD runner configuration |
| Dockerfile / base image layer | Containerized CI/CD environments |
| GitHub Actions workflow step | Repository-scoped scanning in CI |

Plugin policy can be centrally managed by distributing a standard `~/.ztvs/config.yaml` through your configuration management tooling, enforcing which plugins are permitted to use which capabilities across the entire fleet.

## Plugin Lifecycle and Governance

All official plugins are distributed through the ZTVS remote registry with the following governance controls:

1. **Cryptographic Verification**: Every plugin binary is matched against a published SHA-256 checksum before execution
2. **Capability Declaration**: Plugins must explicitly declare all required system permissions in `plugin.yaml`
3. **Policy Enforcement**: Admins define capability allow-lists in `config.yaml`; the host engine rejects any plugin attempting undeclared access
4. **Atomic Updates**: `zt plugin update` downloads and compiles to a temporary cache, performing an atomic rename to prevent scan failures during live updates
5. **Audit Trail**: All plugin capability grants and denials are loggable for compliance review

## Licensing and Open Source Posture

| Component | License |
|---|---|
| ZTVS Host Engine | Apache License 2.0 |
| ZTVS Go SDK | MIT |
| ZTVS Python SDK | MIT |
| First-Party Plugins | MIT |

The entire ZTVS ecosystem is open source. There are no proprietary binaries, no required cloud connectivity, and no telemetry. Organizations with air-gapped environments can build from source and self-host the plugin registry.

## Recommended Next Steps

1. **Proof of Concept**: Run a single-host pilot scan on a representative developer workstation using `zt scan`
2. **Plugin Inventory**: Review the three first-party plugins against your current security coverage gaps
3. **Custom Rule Assessment**: Identify two or three organization-specific security checks that could be built using the Go or Python SDK
4. **SIEM Integration**: Evaluate SARIF output ingestion into your existing security tooling stack
5. **Fleet Rollout Plan**: Define the MDM or configuration management strategy for `zt agent` deployment to server and developer fleets
6. **Policy Configuration**: Draft the `config.yaml` capability policy appropriate for your security posture

For technical deep-dives, see the [Architecture Overview](./architecture/README.md) and the [Protocol Specification](./protocol/README.md).
