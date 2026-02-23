# Boilerblade Global Installer

Install Boilerblade so you can run `boilerblade` from **CMD**, **Git Bash**, and **Terminal** from any directory (like Composer for Laravel).

## Native installers (.msi / .deb / .pkg)

If you have a pre-built installer package:

| Platform | File | Install |
|----------|------|--------|
| **Windows** | `boilerblade-1.0.0-amd64.msi` | Double-click or `msiexec /i boilerblade-1.0.0-amd64.msi`. Installs to C:\Program Files\boilerblade and adds to system PATH. |
| **Linux** | `boilerblade_1.0.0_amd64.deb` | `sudo dpkg -i boilerblade_1.0.0_amd64.deb`. Installs to `/usr/local/bin`. |
| **macOS** | `boilerblade-1.0.0.pkg` | Double-click or `sudo installer -pkg boilerblade-1.0.0.pkg -target /`. Installs to `/usr/local/bin`. |

To **build** these installers from source, see [installer/README.md](installer/README.md) (requires WiX on Windows for .msi).

---

## Script-based install (no installer package)

Alternatively, run the install script from the project root (builds the binary if needed and adds it to PATH):

## Prerequisites

- **Go 1.24+** – [Download Go](https://go.dev/dl/)
- (Optional) Build the binary first: from project root run `go build -o bin/boilerblade ./cmd/cli` (Windows: `bin/boilerblade.exe`). The installer will build for you if the binary is missing.

---

## Windows (CMD, Git Bash, PowerShell)

1. Open **PowerShell** (Run as current user; no admin required).
2. Go to the Boilerblade project directory:
   ```powershell
   cd D:\RND\boilerblade
   ```
3. Run the installer (allow script execution if prompted):
   ```powershell
   PowerShell -ExecutionPolicy Bypass -File install.ps1
   ```
4. **Close and reopen** CMD, Git Bash, and PowerShell so the updated PATH is loaded.
5. From any directory you can run:
   ```bash
   boilerblade new my-api
   boilerblade make all -name=Product
   boilerblade make migration -name=add_orders_table
   ```

**Install location:** `C:\boilerblade\bin`  
The installer adds this folder to your **user** PATH so it works in CMD, Git Bash, and PowerShell without admin rights.

---

## macOS and Linux (Terminal)

1. Open a terminal and go to the Boilerblade project directory:
   ```bash
   cd /path/to/boilerblade
   ```
2. Run the installer (user install, no sudo):
   ```bash
   chmod +x install.sh
   ./install.sh
   ```
3. Load the updated PATH (or open a new terminal):
   ```bash
   export PATH="$HOME/.local/bin:$PATH"
   ```
4. From any directory:
   ```bash
   boilerblade new my-api
   boilerblade make all -name=Product
   boilerblade make migration -name=add_orders_table
   ```

**Install location (default):** `~/.local/bin`  
The script adds this to your shell config (`.zshrc`, `.bashrc`, or `.profile`) if needed.

### Global install (optional, requires sudo)

To install for all users under `/usr/local/bin`:

```bash
./install.sh --global
```

Then run `boilerblade` from any terminal; no PATH change needed.

---

## Verify installation

From a **new** terminal (any directory):

```bash
boilerblade help
boilerblade version
```

If you see the help message and version, the global install is working.

---

## Uninstall

- **Windows:** Remove `C:\boilerblade\bin` from your user PATH (Settings → Environment variables), then delete the folder `C:\boilerblade`.
- **macOS / Linux:** Delete `~/.local/bin/boilerblade` (or `/usr/local/bin/boilerblade` if you used `--global`). Remove the “Boilerblade / local bin” lines from `~/.zshrc`, `~/.bashrc`, or `~/.profile` if you no longer need `~/.local/bin` on PATH.
