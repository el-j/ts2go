# Go Installation Guide

**For TS2Go Users** - Quick reference to install Go compiler on your system

---

## Why Do I Need Go?

TS2Go transpiles your TypeScript code to Go, but you need the Go compiler installed to:
- ✅ **Build** executable binaries from transpiled code
- ✅ **Test** your Go code with `go test`
- ✅ **Run** your projects with `go run`
- ✅ **Format** code with `go fmt` (automatic in TS2Go)

Without Go installed, you can still transpile TypeScript → Go, but you can't compile or run it.

---

## Quick Install

### macOS

**Option 1: Homebrew (Recommended)**
```bash
brew install go
```

**Option 2: Official Installer**
1. Download: https://go.dev/dl/
2. Choose `go1.21.x.darwin-arm64.pkg` (Apple Silicon) or `go1.21.x.darwin-amd64.pkg` (Intel)
3. Run the installer
4. Verify: `go version`

### Windows

**Option 1: Official Installer**
1. Download: https://go.dev/dl/
2. Choose `go1.21.x.windows-amd64.msi`
3. Run the installer (adds to PATH automatically)
4. Open new terminal and verify: `go version`

**Option 2: Chocolatey**
```powershell
choco install golang
```

### Linux

**Option 1: Package Manager**

Ubuntu/Debian:
```bash
sudo apt update
sudo apt install golang-go
```

Fedora:
```bash
sudo dnf install golang
```

Arch:
```bash
sudo pacman -S go
```

**Option 2: Official Tarball**
```bash
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

---

## Verify Installation

After installing, open a **new terminal** and run:

```bash
go version
```

**Expected output:**
```
go version go1.21.5 darwin/arm64
```

If you see this, Go is installed correctly! 🎉

---

## Configure TS2Go to Use Go

### Option 1: System Go (Default)

If Go is in your system PATH, TS2Go will automatically detect it:

1. Open TS2Go
2. Go to **Settings** (⚙️ icon)
3. Click **Project** tab
4. Look for **Go Compiler Configuration**
5. Select **"System Go (from PATH)"**
6. Click **"Detect Go Installation"**
7. You should see: ✅ **Go Compiler Detected** with version and path

### Option 2: Custom Go Path

If you installed Go to a custom location:

1. Open TS2Go
2. Go to **Settings** → **Project** tab
3. Select **"Custom Path"** in dropdown
4. Click the 📁 folder icon
5. Browse to your Go binary:
   - **macOS/Linux:** `/usr/local/go/bin/go` or `/opt/go/bin/go`
   - **Windows:** `C:\Go\bin\go.exe`
6. Click **"Detect Go Installation"**
7. Verify detection succeeded

---

## Troubleshooting

### "Go not found" Error

**Problem:** TS2Go shows "Go Compiler Not Found"

**Solutions:**
1. **Verify Go is installed:**
   ```bash
   go version
   ```
   If this fails, Go is not installed or not in PATH

2. **Add Go to PATH:**
   
   **macOS/Linux:**
   ```bash
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
   source ~/.zshrc
   ```
   
   **Windows:**
   - Search → "Environment Variables"
   - Edit "Path" under User variables
   - Add: `C:\Go\bin`
   - Restart terminal

3. **Use Custom Path in TS2Go:**
   - Settings → Project → Custom Path
   - Browse to Go binary directly

### "Permission Denied" Error

**Problem:** Go binary is not executable

**Solution (macOS/Linux):**
```bash
chmod +x /usr/local/go/bin/go
```

### Go Version Too Old

**Problem:** TS2Go requires Go 1.18+ (for generics)

**Solution:** Upgrade Go:
```bash
# macOS
brew upgrade go

# Linux
sudo apt update && sudo apt upgrade golang-go

# Or download latest from https://go.dev/dl/
```

### Multiple Go Versions

**Problem:** Have multiple Go installations

**Solution:**
1. Check which Go is active:
   ```bash
   which go
   go version
   ```
2. In TS2Go, use Custom Path to select specific version
3. Or update PATH to prioritize desired version

---

## Recommended Go Version

**Minimum:** Go 1.18 (required for generics)  
**Recommended:** Go 1.21+ (latest stable)  
**Download:** https://go.dev/dl/

---

## Testing Your Setup

After installing and configuring Go in TS2Go:

1. **Load a TypeScript project** in TS2Go
2. **Transpile** the project (🚀 Transpile All)
3. **Build** the project (🔨 Build Project)
4. **Run** the project (▶️ Run Project)

If all steps succeed, your Go setup is perfect! ✅

---

## Common Installation Paths

**Where Go is typically installed:**

| OS | Package Manager | Path |
|----|----------------|------|
| macOS | Homebrew | `/opt/homebrew/bin/go` (Apple Silicon)<br>`/usr/local/bin/go` (Intel) |
| macOS | Official | `/usr/local/go/bin/go` |
| Linux | apt/dnf | `/usr/bin/go` |
| Linux | Official | `/usr/local/go/bin/go` |
| Windows | Official | `C:\Go\bin\go.exe` |
| Windows | Chocolatey | `C:\ProgramData\chocolatey\bin\go.exe` |

---

## Need More Help?

### Official Go Documentation
- **Installation Guide:** https://go.dev/doc/install
- **Getting Started:** https://go.dev/doc/tutorial/getting-started
- **Download Page:** https://go.dev/dl/

### TS2Go Settings
1. Open TS2Go
2. Click ⚙️ **Settings** icon
3. Go to **Project** tab
4. Scroll to **Go Compiler Configuration**
5. Click **"Detect Go Installation"** to test

### Check TS2Go Logs
If detection fails, check the Output panel in ProjectView for detailed error messages.

---

**Last Updated:** November 15, 2025  
**TS2Go Version:** 0.1.1+  
**Minimum Go Version:** 1.18  
**Recommended Go Version:** 1.21+
