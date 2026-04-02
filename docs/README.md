# ZTVS Documentation

Welcome to the official documentation for **Zero Trust Vulnerability Scanner (ZTVS)**.

ZTVS is a cross-platform, plugin-based host security auditing tool designed with a Zero Trust philosophy. It executes plugins in isolated processes and communicates via a standard JSON-RPC protocol.

## For Technology Leadership

- **[Product Overview](./product-overview.md)**: Executive brief on ZTVS capabilities, architecture, and org-wide adoption.

## Getting Started

- **[Installation Guide](./installation.md)**: How to get ZTVS running on your system.
- **[Architecture Overview](./architecture.md)**: Understand the core engine and plugin model.
- **[CLI Reference](./commands.md)**: Detailed guide to `zt scan`, `zt plugin`, and `zt agent`.

## For Developers

- **[Plugin Development Guide](./plugin-guide.md)**: Learn how to build your own security checks in any of the five supported languages (Go, Python, Rust, JavaScript, Java) and install them independently.
- **[JSON-RPC Protocol](./protocol.md)**: Specification for host-plugin communication.
- **[Go SDK Reference](./sdk-go.md)**: Documentation for our first-party Go SDK.
- **[Python SDK (External)](https://github.com/mosesgameli/ztvs-sdk-python)**: Our official Python SDK for plugin development.

## Community & Contributing

- **[Contribution Guidelines](../CONTRIBUTING.md)**
- **[Code of Conduct](../CODE_OF_CONDUCT.md)**
- **[RFCs Index](../_/rfc/README.md)** (Internal access only)
