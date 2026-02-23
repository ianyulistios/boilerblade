# Boilerblade Global Installer for Windows
# Installs to C:\boilerblade\bin and adds it to your user PATH so you can run "boilerblade" from CMD, Git Bash, and PowerShell.
# Run with: PowerShell -ExecutionPolicy Bypass -File install.ps1

$InstallDir = "C:\boilerblade\bin"
$BinaryName = "boilerblade.exe"

# Find project root (directory containing go.mod and cmd/cli)
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = $ScriptDir
$BinPath = Join-Path $ProjectRoot "bin"
$SourceExe = Join-Path $BinPath $BinaryName

Write-Host "Boilerblade Installer (Windows)" -ForegroundColor Cyan
Write-Host ""

# Build if binary doesn't exist
if (-not (Test-Path $SourceExe)) {
    Write-Host "Binary not found. Building..." -ForegroundColor Yellow
    Push-Location $ProjectRoot
    try {
        if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
            Write-Host "Error: Go is not installed or not in PATH. Install Go from https://go.dev/dl/" -ForegroundColor Red
            exit 1
        }
        New-Item -ItemType Directory -Force -Path $BinPath | Out-Null
        go build -o $SourceExe ./cmd/cli
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Error: Build failed." -ForegroundColor Red
            exit 1
        }
        Write-Host "Build successful." -ForegroundColor Green
    } finally {
        Pop-Location
    }
} else {
    Write-Host "Using existing binary: $SourceExe" -ForegroundColor Green
}

# Create install directory
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

# Copy binary
Copy-Item -Path $SourceExe -Destination (Join-Path $InstallDir $BinaryName) -Force
Write-Host "Installed to: $InstallDir" -ForegroundColor Green

# Add to user PATH if not already present
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    Write-Host "Added to user PATH: $InstallDir" -ForegroundColor Green
    Write-Host ""
    Write-Host "IMPORTANT: Close and reopen CMD, Git Bash, and PowerShell for 'boilerblade' to be recognized." -ForegroundColor Yellow
} else {
    Write-Host "Already in user PATH." -ForegroundColor Green
}

# Refresh current session so user can try in same window
$env:Path = [Environment]::GetEnvironmentVariable("Path", "User") + ";" + [Environment]::GetEnvironmentVariable("Path", "Machine")

Write-Host ""
Write-Host "Done. You can now run (in a new terminal):" -ForegroundColor Cyan
Write-Host "  boilerblade new my-api"
Write-Host "  boilerblade make all -name=Product"
Write-Host "  boilerblade make migration -name=add_orders_table"
Write-Host ""
