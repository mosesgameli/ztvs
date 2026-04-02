#!/usr/bin/env sh
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# ZTVS uninstaller
# Removes the zt binary and ~/.ztvs configuration directory.
#
# Usage (one-liner):
#   curl -fsSL https://raw.githubusercontent.com/mosesgameli/ztvs/main/uninstall.sh | sh

set -e

INSTALL_BIN="/usr/local/bin"
CONFIG_DIR="$HOME/.ztvs"

# ── helpers ────────────────────────────────────────────────────────────────
red()   { printf '\033[0;31m%s\033[0m\n' "$*"; }
green() { printf '\033[0;32m%s\033[0m\n' "$*"; }
bold()  { printf '\033[1m%s\033[0m\n'   "$*"; }
info()  { printf '  \033[0;34m→\033[0m %s\n' "$*"; }

main() {
  bold ""
  bold "  Uninstalling Zero Trust Vulnerability Scanner..."
  printf '\n'

  # ── remove host binary ──────────────────────────────────────────────────
  if [ -f "${INSTALL_BIN}/zt" ]; then
    info "Removing ${INSTALL_BIN}/zt ..."
    if [ -w "$INSTALL_BIN" ]; then
      rm -f "${INSTALL_BIN}/zt"
    else
      sudo rm -f "${INSTALL_BIN}/zt"
    fi
  else
    info "Binary ${INSTALL_BIN}/zt not found. Skipping."
  fi

  # ── remove config & plugins ──────────────────────────────────────────────
  if [ -d "$CONFIG_DIR" ]; then
    info "Removing directory ${CONFIG_DIR} ..."
    rm -rf "$CONFIG_DIR"
  else
    info "Config directory ${CONFIG_DIR} not found. Skipping."
  fi

  printf '\n'
  green "✓ ZTVS has been completely uninstalled from your system."
  printf '\n'
}

main "$@"
