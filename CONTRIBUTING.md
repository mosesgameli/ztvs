# Contributing to ZTVS

Thank you for your interest in contributing to the Zero Trust Vulnerability Scanner (ZTVS)! We welcome contributions from the community to help make host security more robust and isolated.

## Core Principles

1.  **Go First**: The core engine and SDK are written in Go (target version: 1.24+).
2.  **Zero Trust Execution**: Plugins must execute as separate OS processes and communicate via JSON-RPC 2.0 over stdio.
3.  **No In-Process Plugins**: Shared libraries (`.so`, `.dll`) are prohibited for plugin execution.
4.  **Capability-based Security**: Plugins must declare required capabilities (e.g., `read_files`) in their manifests.

 

## Detailed Contributor Guides

For a deep dive into specific topics, see our technical documentation:

- [**Developer Setup & Testing**](docs/contributor/development.md): Detailed environment setup, architecture-specific instructions, and test pyramid overview.
- [**Coding Style Guide**](docs/contributor/style-guide.md): In-depth Go standards, concurrency patterns, and internal directory topology.
- [**Contributor Workflows**](docs/contributor/workflows.md): Detailed PR processes, branching strategy, and conventional commit standards.
- [**Plugin Development Guide**](docs/contributor/plugins.md): Protocol requirements and Go SDK reference for building new security checks.

 

## Getting Started

### Prerequisites

- [Go 1.26.1+](https://golang.org/dl/)
- `make`
- `git`

### Local Setup

1.  Clone the repository:
    ```bash
    git clone https://github.com/mosesgameli/ztvs.git
    cd ztvs/zt
    ```

2.  Build the host engine:
    ```bash
    make build
    ```

3.  Initialize the local environment:
    ```bash
    make init
    ```

### Running Tests

Execute the unit test suite:
```bash
go test ./...
```

For integration and end-to-end tests, see the `test/` directory.

## Development Workflow

### Branching Strategy

We use a simple branching model. Create a branch for your changes:
- `feat/summary` for new features.
- `fix/summary` for bug fixes.
- `docs/summary` for documentation.

### Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/).

Format: `<type>(<scope>): <description>`

Common types:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `chore`: Updating build tasks, package manager configs, etc.

### Pull Requests

1.  All code MUST be merged via Pull Requests. Direct pushes to `main` are discouraged.
2.  Ensure `make build` and all tests pass before submitting.
3.  Describe *what* changed and *why* in the PR body.

## Coding Standards

- **Internal vs Pkg**:
    - Use `internal/` for logic private to the ZTVS host (e.g., `pluginhost`, `engine`).
    - Use `pkg/` for code intended for external consumption (e.g., `sdk`, `rpc` types).
- **Error Handling**:
    - Propagate errors with context using `fmt.Errorf("context: %w", err)`.
    - Use JSON-RPC error codes (4001: Version mismatch, 4002: Check not found) for plugin communication.
- **Concurrency**:
    - Use `sync.WaitGroup` and worker pools for parallel scanning.
    - All shared resources (e.g., Reporters) MUST be accessed through mutexes.

## Licensing

By contributing to ZTVS, you agree that your contributions will be licensed under the Apache License 2.0.
