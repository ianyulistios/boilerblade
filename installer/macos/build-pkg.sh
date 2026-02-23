#!/usr/bin/env bash
# Build Boilerblade .pkg installer (macOS)
# Run from repo root: ./installer/macos/build-pkg.sh
# Requires: pkgbuild, productbuild (built-in on macOS)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
BINARY_NAME="boilerblade"
VERSION="1.0.0"
ARCH=$(uname -m)
SOURCE_BIN="$PROJECT_ROOT/bin/$BINARY_NAME"
SOURCE_BIN_DARWIN="$PROJECT_ROOT/bin/$BINARY_NAME-darwin-$ARCH"
PKG_NAME="boilerblade"
OUT_PKG="$PROJECT_ROOT/bin/boilerblade-$VERSION.pkg"
STAGING="$PROJECT_ROOT/bin/pkg-staging"

echo "Building Boilerblade .pkg installer"

# Use darwin binary if present, else current binary
if [ -f "$SOURCE_BIN_DARWIN" ]; then
  SRC="$SOURCE_BIN_DARWIN"
elif [ -f "$SOURCE_BIN" ]; then
  SRC="$SOURCE_BIN"
else
  echo "Binary not found. Building for darwin/$ARCH..."
  (cd "$PROJECT_ROOT" && GOOS=darwin GOARCH=$ARCH go build -o "bin/$BINARY_NAME-darwin-$ARCH" ./cmd/cli)
  SRC="$PROJECT_ROOT/bin/$BINARY_NAME-darwin-$ARCH"
fi

rm -rf "$STAGING"
mkdir -p "$STAGING/usr/local/bin"
cp "$SRC" "$STAGING/usr/local/bin/$BINARY_NAME"
chmod 755 "$STAGING/usr/local/bin/$BINARY_NAME"

# Build .pkg (installs to /usr/local/bin)
pkgbuild --identifier com.boilerblade.cli \
  --version "$VERSION" \
  --root "$STAGING" \
  --install-location "/" \
  "$OUT_PKG"

rm -rf "$STAGING"

echo "Created: $OUT_PKG"
echo "Install by double-clicking the .pkg or: sudo installer -pkg $OUT_PKG -target /"
echo "Then run: boilerblade help"
