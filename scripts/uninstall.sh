#!/usr/bin/env sh
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
