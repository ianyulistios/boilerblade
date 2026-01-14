#!/bin/bash

# Build script for Linux/Mac
BINARY_NAME="boilerblade"
CMD_DIR="cmd/cli"
BUILD_DIR="bin"

echo "Building $BINARY_NAME..."
mkdir -p $BUILD_DIR
go build -o $BUILD_DIR/$BINARY_NAME ./$CMD_DIR

if [ $? -eq 0 ]; then
    echo "✓ Binary built successfully: $BUILD_DIR/$BINARY_NAME"
    echo ""
    echo "To install globally, run:"
    echo "  sudo cp $BUILD_DIR/$BINARY_NAME /usr/local/bin/"
    echo ""
    echo "Or add to PATH:"
    echo "  export PATH=\$PATH:\$(pwd)/$BUILD_DIR"
else
    echo "✗ Build failed"
    exit 1
fi
