---
name: open-source-compliance
description: Consolidates the standards for professional, gold-standard open-source development.
---

# Professional Open Source Compliance Skill

This skill provides the AGENT with a comprehensive model of high-quality open-source software implementation.

## Guiding Principles

- **Zero-Ambiguity**: The codebase must be self-documenting and legally clear (LICENSE headers everywhere).
- **Quality Gates**: Manual testing is insufficient; 95%+ automated coverage and security scanning are mandatory.
- **Enterprise-Ready**: The project should follow established patterns (Go internal/pkg, Conventional Commits).

## Key Standards (GOLD STATUS)

### 1. Licensing
- **Requirement**: ALL source files MUST include the Apache License 2.0 preamble.
- **Verification**: Use the [`/license`](.agents/workflows/license.md) workflow.

### 2. Testing & Quality
- **Requirement**: ≥95% unit test coverage for core logic (Internal/Pkg).
- **Requirement**: No vulnerabilities detected in dependencies or static analysis.
- **Verification**: Use the [`/test`](.agents/workflows/test.md) and [`/scan`](.agents/workflows/scan.md) workflows.

### 3. Documentation
- **Requirement**: Every exported symbol MUST be documented following Go conventions.
- **Requirement**: ADRs (RFCs) should be maintained for all major architecture decisions.

### 4. Commits & Workflows
- **Requirement**: Use Conventional Commits.
- **Requirement**: All PRs should be linked to an issue and verified via the [`/pr`](.agents/workflows/pr.md) workflow.

## Usage in Pair Programming

When solving coding tasks, the agent should proactively:
1.  Verify license headers on new files.
2.  Suggest unit test improvements to reach the 95% target.
3.  Ensure public APIs are properly documented.
