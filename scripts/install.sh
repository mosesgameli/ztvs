#!/usr/bin/env sh
# ZTVS installer — downloads the latest release for your platform,
# installs the zt binary to /usr/local/bin, and seeds first-party
# plugins into ~/.ztvs/plugins.
#
# Usage (one-liner):
#   curl -fsSL https://raw.githubusercontent.com/mosesgameli/ztvs/main/install.sh | sh

set -e

REPO="mosesgameli/ztvs"
INSTALL_BIN="${ZTVS_INSTALL_BIN:-/usr/local/bin}"
ZTVS_HOME="${ZTVS_HOME:-$HOME/.ztvs}"
PLUGIN_DIR="$ZTVS_HOME/plugins"

# ── helpers ────────────────────────────────────────────────────────────────
red()   { printf '\033[0;31m%s\033[0m\n' "$*"; }
green() { printf '\033[0;32m%s\033[0m\n' "$*"; }
bold()  { printf '\033[1m%s\033[0m\n'   "$*"; }
info()  { printf '  \033[0;34m→\033[0m %s\n' "$*"; }

die() { red "error: $*" >&2; exit 1; }

need_cmd() { command -v "$1" >/dev/null 2>&1 || die "required command not found: $1"; }

# ── detect platform ────────────────────────────────────────────────────────
detect_os() {
  case "$(uname -s)" in
    Linux*)   echo linux ;;
    Darwin*)  echo darwin ;;
    MINGW*|MSYS*|CYGWIN*) echo windows ;;
    *) die "unsupported OS: $(uname -s)" ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo amd64 ;;
    arm64|aarch64) echo arm64 ;;
    *) die "unsupported architecture: $(uname -m)" ;;
  esac
}

# ── resolve latest version ─────────────────────────────────────────────────
latest_version() {
  need_cmd curl
  curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' \
    | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/'
}

# ── main ───────────────────────────────────────────────────────────────────
main() {
  bold ""
  bold "  ███████╗████████╗██╗   ██╗███████╗"
  bold "     ███╔╝╚══██╔══╝██║   ██║██╔════╝"
  bold "    ███╔╝    ██║   ██║   ██║███████╗"
  bold "   ███╔╝     ██║   ╚██╗ ██╔╝╚════██║"
  bold "  ███████╗   ██║    ╚████╔╝ ███████║"
  bold "  ╚══════╝   ╚═╝     ╚═══╝  ╚══════╝"
  bold "  Zero Trust Vulnerability Scanner"
  printf '\n'

  need_cmd curl
  need_cmd tar

  OS=$(detect_os)
  ARCH=$(detect_arch)

  VERSION="${ZTVS_VERSION:-$(latest_version)}"
  [ -n "$VERSION" ] || die "could not determine latest version"

  info "Detected platform : ${OS}/${ARCH}"
  info "Installing version: ${VERSION}"

  if [ "$OS" = "windows" ]; then
    ARCHIVE="ztvs-${VERSION}-${OS}-${ARCH}.zip"
    need_cmd unzip
  else
    ARCHIVE="ztvs-${VERSION}-${OS}-${ARCH}.tar.gz"
  fi

  DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"
  TMP_DIR="$(mktemp -d)"
  trap 'rm -rf "$TMP_DIR"' EXIT

  info "Downloading ${ARCHIVE} ..."
  curl -fsSL --progress-bar "$DOWNLOAD_URL" -o "${TMP_DIR}/${ARCHIVE}" \
    || die "download failed — check that release ${VERSION} exists for ${OS}/${ARCH}"

  info "Extracting ..."
  if [ "$OS" = "windows" ]; then
    unzip -q "${TMP_DIR}/${ARCHIVE}" -d "$TMP_DIR"
  else
    tar -xzf "${TMP_DIR}/${ARCHIVE}" -C "$TMP_DIR"
  fi

  # ── install host binary ──────────────────────────────────────────────────
  BIN_EXT=""
  [ "$OS" = "windows" ] && BIN_EXT=".exe"

  info "Installing zt${BIN_EXT} → ${INSTALL_BIN}/"
  if [ -w "$INSTALL_BIN" ]; then
    cp "${TMP_DIR}/zt${BIN_EXT}" "${INSTALL_BIN}/zt${BIN_EXT}"
  else
    sudo cp "${TMP_DIR}/zt${BIN_EXT}" "${INSTALL_BIN}/zt${BIN_EXT}"
  fi
  chmod +x "${INSTALL_BIN}/zt${BIN_EXT}"

  # ── install plugins ──────────────────────────────────────────────────────
  mkdir -p "$PLUGIN_DIR"

  for plugin in plugin-os plugin-axios-mitigation; do
    BIN="${TMP_DIR}/${plugin}${BIN_EXT}"
    YAML="${TMP_DIR}/${plugin}.yaml"

    if [ -f "$BIN" ]; then
      info "Installing plugin: ${plugin}"
      cp "$BIN"  "${PLUGIN_DIR}/${plugin}${BIN_EXT}"
      chmod +x   "${PLUGIN_DIR}/${plugin}${BIN_EXT}"
      [ -f "$YAML" ] && cp "$YAML" "${PLUGIN_DIR}/${plugin}.yaml"
    fi
  done

  # ── bootstrap config if first install ────────────────────────────────────
  if [ ! -f "${ZTVS_HOME}/config.yaml" ]; then
    info "Bootstrapping ${ZTVS_HOME}/config.yaml ..."
    "${INSTALL_BIN}/zt${BIN_EXT}" plugin init 2>/dev/null || true
  fi

  printf '\n'
  green "✓ ZTVS ${VERSION} installed successfully!"
  info  "Run 'zt scan' to start your first audit."
  printf '\n'
}

main "$@"
