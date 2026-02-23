# Boilerblade native installers

Build platform-specific installers so users can install Boilerblade like Composer (globally in CMD, Git Bash, Terminal).

| Platform | Format | Install location | Build script |
|----------|--------|------------------|--------------|
| Windows  | `.msi` | C:\Program Files\boilerblade, added to system PATH | `installer\windows\build-msi.ps1` |
| Linux    | `.deb` | /usr/local/bin | `installer/linux/build-deb.sh` |
| macOS    | `.pkg` | /usr/local/bin | `installer/macos/build-pkg.sh` |

## Prerequisites

- **All:** Go 1.24+ and this repo built (or the script will build the binary).
- **Windows .msi:** [WiX Toolset 3.x](https://wixtoolset.org/docs/wix3/) installed (e.g. v3.11). The script looks for `WiX Toolset v3.11\bin` or use env `WIX` to point to the bin folder.
- **Linux .deb:** `dpkg-deb` (standard on Debian/Ubuntu).
- **macOS .pkg:** `pkgbuild` (built-in on macOS).

## Build installers

From the **repository root**:

### Windows (.msi)

```powershell
.\installer\windows\build-msi.ps1
```

Output: `bin\boilerblade-1.0.0-amd64.msi`  
Install: double-click the .msi or `msiexec /i bin\boilerblade-1.0.0-amd64.msi` (admin). Then open a new CMD/Git Bash/PowerShell and run `boilerblade`.

### Linux (.deb)

```bash
chmod +x installer/linux/build-deb.sh
./installer/linux/build-deb.sh
```

Output: `bin/boilerblade_1.0.0_amd64.deb`  
Install: `sudo dpkg -i bin/boilerblade_1.0.0_amd64.deb`, then `boilerblade help`.

For arm64: `GOARCH=arm64 ./installer/linux/build-deb.sh` → `boilerblade_1.0.0_arm64.deb`.

### macOS (.pkg)

```bash
chmod +x installer/macos/build-pkg.sh
./installer/macos/build-pkg.sh
```

Output: `bin/boilerblade-1.0.0.pkg`  
Install: double-click the .pkg or `sudo installer -pkg bin/boilerblade-1.0.0.pkg -target /`, then `boilerblade help`.

## Makefile targets

From repo root:

```bash
make build-installer-windows   # .msi (requires WiX on Windows)
make build-installer-linux     # .deb (on Linux or WSL)
make build-installer-macos     # .pkg (on macOS)
```

These run the scripts above; the binary is built automatically if missing.
