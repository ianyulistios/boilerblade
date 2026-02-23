#!/usr/bin/env bash
# Boilerblade Global Installer for macOS and Linux
# Installs to ~/.local/bin (default, no sudo) or /usr/local/bin (with --global).
# Adds the install dir to PATH in your shell profile so "boilerblade" works in Terminal.
#
# Usage:
#   ./install.sh           # install to ~/.local/bin (recommended)
#   ./install.sh --global  # install to /usr/local/bin (requires sudo)

set -e

INSTALL_GLOBAL=false
if [ "${1:-}" = "--global" ]; then
  INSTALL_GLOBAL=true
fi

# Project root (directory containing install.sh and go.mod)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"
BINARY_NAME="boilerblade"
SOURCE_BIN="$PROJECT_ROOT/bin/$BINARY_NAME"

echo "Boilerblade Installer (macOS / Linux)"
echo ""

# Build if binary doesn't exist
if [ ! -f "$SOURCE_BIN" ]; then
  echo "Binary not found. Building..."
  if ! command -v go >/dev/null 2>&1; then
    echo "Error: Go is not installed or not in PATH. Install Go from https://go.dev/dl/"
    exit 1
  fi
  mkdir -p "$PROJECT_ROOT/bin"
  (cd "$PROJECT_ROOT" && go build -o "bin/$BINARY_NAME" ./cmd/cli)
  echo "Build successful."
else
  echo "Using existing binary: $SOURCE_BIN"
fi

if [ "$INSTALL_GLOBAL" = true ]; then
  INSTALL_DIR="/usr/local/bin"
  echo "Installing globally to $INSTALL_DIR (requires sudo)"
  sudo mkdir -p "$INSTALL_DIR"
  sudo cp "$SOURCE_BIN" "$INSTALL_DIR/$BINARY_NAME"
  sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
  INSTALL_DIR="$HOME/.local/bin"
  echo "Installing to $INSTALL_DIR (user)"
  mkdir -p "$INSTALL_DIR"
  cp "$SOURCE_BIN" "$INSTALL_DIR/$BINARY_NAME"
  chmod +x "$INSTALL_DIR/$BINARY_NAME"

  # Add ~/.local/bin to PATH in one shell profile (so boilerblade works in new terminals)
  PATH_LINE='export PATH="$HOME/.local/bin:$PATH"'
  for profile in "$HOME/.zshrc" "$HOME/.bashrc" "$HOME/.profile"; do
    if [ -f "$profile" ]; then
      if grep -qF '.local/bin' "$profile" 2>/dev/null; then
        break
      fi
      echo "" >> "$profile"
      echo "# Boilerblade / local bin" >> "$profile"
      echo "$PATH_LINE" >> "$profile"
      echo "Added PATH to $profile"
      break
    fi
  done
fi

echo ""
echo "Installed to: $INSTALL_DIR"
echo ""
if [ "$INSTALL_GLOBAL" = false ]; then
  echo "Run this in your terminal (or open a new one) so 'boilerblade' is found:"
  echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
  echo ""
fi
echo "Then you can run:"
echo "  boilerblade new my-api"
echo "  boilerblade make all -name=Product"
echo "  boilerblade make migration -name=add_orders_table"
echo ""
