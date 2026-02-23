#!/usr/bin/env bash
# Build Boilerblade .deb package (Debian/Ubuntu)
# Run from repo root: ./installer/linux/build-deb.sh
# Requires: dpkg-deb

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
BINARY_NAME="boilerblade"
VERSION="1.0.0"
ARCH="${GOARCH:-amd64}"
SOURCE_BIN="$PROJECT_ROOT/bin/$BINARY_NAME"
SOURCE_BIN_LINUX="$PROJECT_ROOT/bin/$BINARY_NAME-linux-$ARCH"
PKG_NAME="boilerblade"
PKG_DIR="$PROJECT_ROOT/bin/deb-build"
OUT_DEB="$PROJECT_ROOT/bin/${PKG_NAME}_${VERSION}_${ARCH}.deb"

echo "Building Boilerblade .deb package"

# Use linux binary if present (from make build-all), else current binary
if [ -f "$SOURCE_BIN_LINUX" ]; then
  SRC="$SOURCE_BIN_LINUX"
elif [ -f "$SOURCE_BIN" ]; then
  SRC="$SOURCE_BIN"
else
  echo "Binary not found. Building for linux/$ARCH..."
  (cd "$PROJECT_ROOT" && GOOS=linux GOARCH=$ARCH go build -o "bin/$BINARY_NAME-linux-$ARCH" ./cmd/cli)
  SRC="$PROJECT_ROOT/bin/$BINARY_NAME-linux-$ARCH"
fi

rm -rf "$PKG_DIR"
mkdir -p "$PKG_DIR/DEBIAN"
mkdir -p "$PKG_DIR/usr/local/bin"

cp "$SRC" "$PKG_DIR/usr/local/bin/$BINARY_NAME"
chmod 755 "$PKG_DIR/usr/local/bin/$BINARY_NAME"

DEB_ARCH="amd64"
[ "$ARCH" = "arm64" ] && DEB_ARCH="arm64"

cat > "$PKG_DIR/DEBIAN/control" << EOF
Package: $PKG_NAME
Version: $VERSION
Section: devel
Priority: optional
Architecture: $DEB_ARCH
Maintainer: Boilerblade
Description: Boilerblade CLI - Go boilerplate generator
 Create projects and generate code (model, repository, usecase, handler, dto, consumer, migration).
 Use: boilerblade new my-api, boilerblade make all -name=Product
EOF

dpkg-deb -b "$PKG_DIR" "$OUT_DEB"
rm -rf "$PKG_DIR"

echo "Created: $OUT_DEB"
echo "Install with: sudo dpkg -i $OUT_DEB"
echo "Then run: boilerblade help"
