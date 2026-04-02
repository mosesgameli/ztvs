#!/usr/bin/env bash

set -euo pipefail

# Installer Smoke Test for ZTVS
# Validates the scripts/install.sh script in a mocked environment.

# 1. Setup Mock Environment
export HOME_MOCK=$(mktemp -d)
export BIN_DIR_MOCK="$HOME_MOCK/bin"
mkdir -p "$BIN_DIR_MOCK"

# Mock local profile path
export PATH="$BIN_DIR_MOCK:$PATH"

echo "Running installer smoke test..."
echo "Mock Home: $HOME_MOCK"

# 2. Execute Installer
# We need to mock the 'curl' download by pointing to the local file
# or just running the local file directly with mocked environment variables.
export ZTVS_INSTALL_BIN="$BIN_DIR_MOCK"
export ZTVS_HOME="$HOME_MOCK/.ztvs"

# Run install.sh
bash ./scripts/install.sh

# 3. Verify Local Installation
if [[ ! -d "$ZTVS_HOME" ]]; then
    echo "Error: Installation directory $ZTVS_HOME not created."
    exit 1
fi

if [[ ! -d "$ZTVS_HOME/plugins" ]]; then
    echo "Error: Plugins directory not created."
    exit 1
fi

# 4. Success
echo "✓ Installer smoke test passed!"
# Clean up
rm -rf "$HOME_MOCK"
