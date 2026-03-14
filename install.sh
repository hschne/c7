#!/bin/sh
# Install c7 from GitHub Releases.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/hschne/c7/main/install.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/hschne/c7/main/install.sh | sh -s -- --local
#
# Options:
#   --local       Install to ~/.local/bin instead of /usr/local/bin
#
# Environment variables:
#   C7_VERSION    Version to install (default: latest)
#   C7_INSTALL    Install directory (default: /usr/local/bin, overridden by --local)
#
set -e

REPO="hschne/c7"
VERSION="${C7_VERSION:-latest}"
INSTALL_DIR="${C7_INSTALL:-/usr/local/bin}"

main() {
  parse_args "$@"

  os="$(detect_os)"
  arch="$(detect_arch)"

  if [ "$VERSION" = "latest" ]; then
    VERSION="$(latest_version)"
  fi

  version_bare="${VERSION#v}"
  asset="c7_${version_bare}_${os}_${arch}.tar.gz"
  url="https://github.com/${REPO}/releases/download/${VERSION}/${asset}"

  tmpdir="$(mktemp -d)"
  trap 'rm -rf "$tmpdir"' EXIT

  echo "Downloading c7 ${VERSION} for ${os}/${arch}..."
  download "$url" "$tmpdir/$asset"

  echo "Extracting..."
  tar -xzf "$tmpdir/$asset" -C "$tmpdir"

  echo "Installing to ${INSTALL_DIR}/c7..."
  install_binary "$tmpdir/c7" "$INSTALL_DIR/c7"

  echo "Done. Run 'c7 --version' to verify."
}

parse_args() {
  while [ "$#" -gt 0 ]; do
    case "$1" in
      --local)
        INSTALL_DIR="${HOME}/.local/bin"
        ;;
      -h|--help)
        usage
        exit 0
        ;;
      *)
        echo "Unknown option: $1" >&2
        usage >&2
        exit 1
        ;;
    esac
    shift
  done
}

usage() {
  cat <<EOF
Usage: install.sh [--local]

Options:
  --local   Install to ~/.local/bin instead of /usr/local/bin

Environment variables:
  C7_VERSION  Version to install (default: latest)
  C7_INSTALL  Install directory (default: /usr/local/bin)
EOF
}

detect_os() {
  case "$(uname -s)" in
    Linux*)  echo "linux" ;;
    Darwin*) echo "darwin" ;;
    *)       echo "Unsupported OS: $(uname -s)" >&2; exit 1 ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    aarch64|arm64) echo "arm64" ;;
    *) echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
  esac
}

latest_version() {
  url="https://api.github.com/repos/${REPO}/releases/latest"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/'
  elif command -v wget >/dev/null 2>&1; then
    wget -qO- "$url" | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/'
  else
    echo "Error: curl or wget is required" >&2
    exit 1
  fi
}

download() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL -o "$2" "$1"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$2" "$1"
  else
    echo "Error: curl or wget is required" >&2
    exit 1
  fi
}

install_binary() {
  src="$1"
  dest="$2"

  mkdir -p "$(dirname "$dest")"

  if [ -w "$(dirname "$dest")" ]; then
    install -m 755 "$src" "$dest"
  else
    echo "Need elevated permissions to install to $(dirname "$dest")."
    sudo install -d -m 755 "$(dirname "$dest")"
    sudo install -m 755 "$src" "$dest"
  fi
}

main "$@"
