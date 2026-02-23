.PHONY: build install clean help

# Binary name
BINARY_NAME=boilerblade
CMD_DIR=cmd/cli
BUILD_DIR=bin

# Build binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "✓ Binary built successfully: $(BUILD_DIR)/$(BINARY_NAME)"

# Install binary to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	@go install ./$(CMD_DIR)
	@echo "✓ $(BINARY_NAME) installed to $$GOPATH/bin"
	@echo "Make sure $$GOPATH/bin is in your PATH"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)
	@echo "✓ Binaries built for all platforms"

# Build native installers (.msi, .deb, .pkg) — run on the target OS (or WSL for .deb)
build-installer-windows:
	@powershell -ExecutionPolicy Bypass -File installer/windows/build-msi.ps1

build-installer-linux:
	@chmod +x installer/linux/build-deb.sh && ./installer/linux/build-deb.sh

build-installer-macos:
	@chmod +x installer/macos/build-pkg.sh && ./installer/macos/build-pkg.sh

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "✓ Cleaned"

# Show help
help:
	@echo "Boilerblade CLI Generator - Build Commands"
	@echo ""
	@echo "Usage:"
	@echo "  make build       - Build binary for current platform"
	@echo "  make install     - Install binary to GOPATH/bin"
	@echo "  make build-all   - Build binaries for all platforms"
	@echo "  make build-installer-windows - Build .msi (Windows, requires WiX)"
	@echo "  make build-installer-linux   - Build .deb (Linux)"
	@echo "  make build-installer-macos   - Build .pkg (macOS)"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make help        - Show this help message"
	@echo ""
	@echo "After installation, use:"
	@echo "  boilerblade new my-api"
	@echo "  boilerblade make all -name=Product"
