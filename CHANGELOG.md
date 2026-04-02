# Changelog

All notable changes to the Zero Trust Vulnerability Scanner (ZTVS) will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Standardized project governance docs (`CONTRIBUTING.md`, `SECURITY.md`, etc.).
- Restructured `docs/` for better navigation and onboarding.
- Initialized `test/` directory skeleton for unified testing.
- Relocated installer scripts to a dedicated `scripts/` directory.

### Changed
- Moved `install.sh`, `install.ps1`, `uninstall.sh`, and `uninstall.ps1` to `/scripts/*`.
- Updated `README.md` installation URLs.
- Merged `sdk-go.md` with `plugin-guide.md` into `docs/contributor/plugins.md`.

## [0.1.0] - 2026-04-01

### Added
- Initial core engine with parallel JSON-RPC worker pool.
- Support for out-of-process Zero Trust execution.
- Capability-based security policy enforcement.
- Remote plugin registry and atomic update mechanism.
- SARIF 2.1.0 reporting support.
- First-party Go and Python SDKs.
