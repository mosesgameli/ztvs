# Developer Setup & Testing Guide

This guide provides detailed instructions for setting up your local development environment and running the ZTVS test suite.

## Environment Setup

### macOS
1.  **Homebrew**: Ensure [Homebrew](https://brew.sh/) is installed.
2.  **Go**: `brew install go` (recommend 1.26.1+).
3.  **GCC**: `xcode-select --install` (required for some Go builds).

### Linux (Ubuntu/Debian)
1.  **Go**: Follow the [official install guide](https://go.dev/doc/install).
2.  **Build Tools**: `sudo apt update && sudo apt install build-essential`.

### Windows
1.  **Go**: Download the MSI from [golang.org](https://go.dev/dl/).
2.  **Make**: Install `make` via [Chocolatey](https://chocolatey.org/) (`choco install make`) or use WSL2.

 

## Workspace Initialization

After cloning the repository, initialize your environment:

```bash
cd ztvs/zt
make build    # Pre-builds the host engine
make init     # Provisions ~/.ztvs/ and sets up the plugin directory
```

 

## The Testing Pyramid

ZTVS follows a strict testing hierarchy to ensure reliability across the host engine and multi-language plugins.

### 1. Unit Tests
Location: `internal/**/*_test.go`, `pkg/**/*_test.go`
Focus on individual functions, registry parsing, and CLI command logic.
```bash
go test ./... -v
```

### 2. Contract Tests
Location: `test/contracts/`
Focus on ensuring the JSON-RPC interface remains stable. These tests use JSON fixtures to validate that the host and plugins follow the same protocol specification.

### 3. Integration Tests
Location: `test/integration/`
Validate the interaction between the host and external plugin processes. These tests often use a "mock" plugin to verify handshake and timeout logic.

### 4. End-to-End (E2E) Tests
Location: `test/e2e/`
Full CLI walk-throughs:
- `zt scan` (with various output formats)
- `zt plugin install`
- `zt agent`

 

## Makefile Reference

| Target | Description |
| :--- | :--- |
| `make build` | Compiles the `zt` host engine to the root. |
| `make init` | Seeds the default configuration and plugin paths. |
| `make sync_manifests` | Synchronizes local plugin manifests with the remote registry. |
| `make clean` | Removes local binaries and temporary artifacts. |
