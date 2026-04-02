# Contributor Workflows Guide

This guide provides deep technical context on the project's branching models, commit standards, and Pull Request (PR) lifecycle.

## 1. Branching Model

We follow a structured branching model. All official development occurs in the `main` branch. Contributors must work on feature-specific branches:

-   **`feat/`**: Dedicated to implementing new functionality.
-   **`fix/`**: Dedicated to bug fixes.
-   **`docs/`**: Dedicated to documentation improvements.
-   **`chore/`**: Dedicated to maintenance tasks (e.g., dependencies, build scripts).

### Cleanup
Avoid working out of `main` or long-lived feature branches. Always sync your local branch with `me/main` (or `origin/main`) before submitting a PR.

 

## 2. Commit Standards

ZTVS follows the [Conventional Commits](https://www.conventionalcommits.org/) specification.

### Structure
```text
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Common Scopes
-   `engine`: For core host scanner changes.
-   `pluginhost`: For registry and plugin execution changes.
-   `sdk`: For changes in `pkg/sdk`.
-   `registry`: For plugin manifest and metadata updates.
-   `repo`: For project infrastructure (e.g., scripts, README).

 

## 3. Pull Request (PR) Lifecycle

### 1. Self-Check
Before opening a PR, ensure:
-   `make build` succeeds.
-   `go test ./...` passes.
-   All new code includes doc comments for exported types/functions.

### 2. Opening the PR
-   Use a descriptive title that matches your commit message.
-   Fill out the PR template (if available) or provide a clear "Changes Made" summary and "How to Test" instructions.
-   Include relevant issue numbers (e.g., `Closes #123`).

### 3. Review Process
-   Maintainers will review your code for logic, style, and security.
-   Respect the review feedback and respond to comments.
-   Push follow-up commits to the same branch; GitHub will automatically update the PR.

### 4. Merge
-   Maintainers will merge your PR once all quality gates pass and at least one approval is received.
-   Squash and merge is the preferred merge method to keep the history clean.
