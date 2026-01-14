# Building and Installing Boilerblade CLI

Panduan untuk build dan install binary CLI generator.

## Quick Start

### Option 1: Install via Go Install (Recommended)

```bash
go install ./cmd/generate
```

Setelah install, pastikan `$GOPATH/bin` ada di PATH, lalu gunakan:

```bash
boilerblade -layer=all -name=Product
```

### Option 2: Build Binary

#### Using Makefile (Linux/Mac)

```bash
# Build binary
make build

# Install to GOPATH/bin
make install

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

#### Using Build Scripts

**Linux/Mac:**
```bash
chmod +x build.sh
./build.sh
```

**Windows:**
```cmd
build.bat
```

#### Manual Build

```bash
# Build for current platform
go build -o bin/boilerblade ./cmd/generate

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o bin/boilerblade-linux-amd64 ./cmd/generate
GOOS=windows GOARCH=amd64 go build -o bin/boilerblade-windows-amd64.exe ./cmd/generate
GOOS=darwin GOARCH=amd64 go build -o bin/boilerblade-darwin-amd64 ./cmd/generate
```

## Installation

### Linux/Mac

1. **Build binary:**
   ```bash
   make build
   # or
   ./build.sh
   ```

2. **Install globally:**
   ```bash
   sudo cp bin/boilerblade /usr/local/bin/
   ```

3. **Or add to PATH:**
   ```bash
   export PATH=$PATH:$(pwd)/bin
   # Add to ~/.bashrc or ~/.zshrc for permanent
   ```

4. **Verify installation:**
   ```bash
   boilerblade -version
   ```

### Windows

1. **Build binary:**
   ```cmd
   build.bat
   # or
   go build -o bin\boilerblade.exe .\cmd\generate
   ```

2. **Add to PATH:**
   - Copy `bin\boilerblade.exe` ke folder yang ada di PATH (e.g., `C:\Windows\System32\`)
   - Atau tambahkan `bin` folder ke PATH environment variable

3. **Verify installation:**
   ```cmd
   boilerblade -version
   ```

## Usage After Installation

Setelah binary terinstall, gunakan langsung:

```bash
# Generate all layers
boilerblade -layer=all -name=Product -fields="Name:string:required,Price:float64:required"

# Generate specific layer
boilerblade -layer=model -name=Order

# Show help
boilerblade -help

# Show version
boilerblade -version
```

## Troubleshooting

### Command not found

**Linux/Mac:**
- Pastikan `$GOPATH/bin` ada di PATH:
  ```bash
  echo $GOPATH
  export PATH=$PATH:$GOPATH/bin
  ```
- Atau gunakan full path: `$GOPATH/bin/boilerblade`

**Windows:**
- Pastikan folder binary ada di PATH
- Atau gunakan full path: `C:\path\to\bin\boilerblade.exe`

### Permission denied (Linux/Mac)

```bash
chmod +x bin/boilerblade
```

### Binary tidak ditemukan

1. Cek lokasi binary:
   ```bash
   # Linux/Mac
   which boilerblade
   
   # Windows
   where boilerblade
   ```

2. Pastikan binary sudah di-build:
   ```bash
   ls -la bin/boilerblade  # Linux/Mac
   dir bin\boilerblade.exe # Windows
   ```

## Cross-Platform Build

Untuk build binary untuk platform lain:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bin/boilerblade-linux-amd64 ./cmd/generate

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/boilerblade-windows-amd64.exe ./cmd/generate

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o bin/boilerblade-darwin-amd64 ./cmd/generate

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o bin/boilerblade-darwin-arm64 ./cmd/generate
```

## Development

Untuk development, gunakan `go run`:

```bash
go run cmd/generate/main.go -layer=all -name=Product
```

## File Structure

```
cmd/
└── generate/
    └── main.go          # CLI entry point

bin/                     # Build output (created after build)
└── boilerblade         # Binary executable

Makefile                # Build commands
build.sh                # Linux/Mac build script
build.bat               # Windows build script
```
