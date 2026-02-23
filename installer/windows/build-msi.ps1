# Build Boilerblade Windows .msi installer (requires WiX Toolset 3.x)
# WiX: https://wixtoolset.org/docs/wix3/
# Run from repo root: .\installer\windows\build-msi.ps1

$ErrorActionPreference = "Stop"
$ProjectRoot = Split-Path -Parent (Split-Path -Parent (Split-Path -Parent $PSScriptRoot))
$BinDir = Join-Path $ProjectRoot "bin"
$ExePath = Join-Path $BinDir "boilerblade.exe"
$InstallerDir = Join-Path $ProjectRoot "installer\windows"
$WxsPath = Join-Path $InstallerDir "boilerblade.wxs"
$Version = "1.0.0"
$OutMsi = Join-Path $ProjectRoot "bin\boilerblade-$Version-amd64.msi"

Write-Host "Building Boilerblade .msi installer" -ForegroundColor Cyan

if (-not (Test-Path $ExePath)) {
    Write-Host "Binary not found. Building..." -ForegroundColor Yellow
    Push-Location $ProjectRoot
    try {
        New-Item -ItemType Directory -Force -Path $BinDir | Out-Null
        go build -o $ExePath ./cmd/cli
        if ($LASTEXITCODE -ne 0) { exit 1 }
    } finally { Pop-Location }
}

$wixPath = $env:WIX
if (-not $wixPath) {
    $wixPath = "${env:ProgramFiles(x86)}\WiX Toolset v3.11\bin"
    if (-not (Test-Path $wixPath)) { $wixPath = "$env:ProgramFiles\WiX Toolset v3.11\bin" }
}
if (-not (Test-Path $wixPath)) {
    Write-Host "Error: WiX Toolset not found. Install from https://wixtoolset.org/docs/wix3/" -ForegroundColor Red
    Write-Host "  Or set WIX to the WiX bin directory (e.g. C:\Program Files (x86)\WiX Toolset v3.11\bin)" -ForegroundColor Red
    exit 1
}

$candle = Join-Path $wixPath "candle.exe"
$light = Join-Path $wixPath "light.exe"
$outDir = Join-Path $InstallerDir "obj"
New-Item -ItemType Directory -Force -Path $outDir | Out-Null

Write-Host "Running candle..."
& $candle -out (Join-Path $outDir "boilerblade.wixobj") -dSourceDir=$BinDir $WxsPath
if ($LASTEXITCODE -ne 0) { exit 1 }

Write-Host "Running light..."
& $light -out $OutMsi (Join-Path $outDir "boilerblade.wixobj")
if ($LASTEXITCODE -ne 0) { exit 1 }

Write-Host "Created: $OutMsi" -ForegroundColor Green
Write-Host "Run the .msi to install (adds C:\Program Files\boilerblade to system PATH)." -ForegroundColor Cyan
