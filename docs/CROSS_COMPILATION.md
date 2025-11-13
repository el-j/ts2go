# Cross-Compilation Guide for TS2Go Desktop

This guide explains how to build TS2Go Desktop for Windows, Linux, and macOS from a single development machine.

## Quick Start

### Build Commands

```bash
# Build for current platform only
make build-desktop

# Build for specific platforms
make build-desktop-windows  # Windows x86_64
make build-desktop-linux    # Linux x86_64

# Build for all platforms at once
make build-desktop-all
```

## Platform Setup

### Building on macOS (Current Platform)

#### 1. Install Rust Targets

```bash
# Windows target (x86_64)
rustup target add x86_64-pc-windows-msvc

# Linux target (x86_64)
rustup target add x86_64-unknown-linux-gnu

# Verify installation
rustup target list --installed
```

#### 2. Install Cross-Compilation Tools

**For Windows builds:**
```bash
# Install LLVM (required for Windows cross-compilation)
brew install llvm

# Add to your shell profile (~/.zshrc or ~/.bash_profile):
export PATH="/opt/homebrew/opt/llvm/bin:$PATH"
```

**For Linux builds:**
```bash
# Install cross-compilation toolchain
brew tap messense/macos-cross-toolchains
brew install x86_64-unknown-linux-gnu

# Or use cross (easier, Docker-based)
cargo install cross
```

#### 3. Platform-Specific Node.js Binaries

The Makefile bundles your system's Node.js by default. For true cross-platform builds, you'll need platform-specific Node.js binaries:

**Option A: Download manually**
- Windows: https://nodejs.org/dist/v20.11.1/node-v20.11.1-win-x64.zip
- Linux: https://nodejs.org/dist/v20.11.1/node-v20.11.1-linux-x64.tar.xz
- macOS: Already bundled from your system

Extract and place in `desktop-ui/src-tauri/bin/`:
- `node.exe` for Windows
- `node` (Linux binary) for Linux

**Option B: Use pkg to create standalone binaries**
```bash
npm install -g pkg
# Creates self-contained Node.js binaries for all platforms
```

### Building on Linux

#### 1. Install Rust Targets

```bash
# Windows target
rustup target add x86_64-pc-windows-msvc

# macOS target (Intel)
rustup target add x86_64-apple-darwin

# macOS target (Apple Silicon)
rustup target add aarch64-apple-darwin
```

#### 2. Install Cross-Compilation Tools

**For Windows:**
```bash
# Install MinGW toolchain
sudo apt-get install mingw-w64
```

**For macOS:**
```bash
# Use osxcross (requires macOS SDK)
# See: https://github.com/tpoechtrager/osxcross
```

### Building on Windows

#### 1. Install Rust Targets

```powershell
# Linux target
rustup target add x86_64-unknown-linux-gnu

# macOS targets
rustup target add x86_64-apple-darwin
rustup target add aarch64-apple-darwin
```

#### 2. Install WSL2 for Linux Builds

```powershell
wsl --install
# Use WSL2 Ubuntu for building Linux targets
```

## Tauri Configuration

The `tauri.conf.json` is already configured to bundle resources correctly:

```json
{
  "bundle": {
    "resources": [
      "bin/ts2go-cli*",
      "bin/parser",
      "bin/node*",
      "bin/mappings"
    ]
  }
}
```

This ensures all dependencies are bundled regardless of target platform.

## Common Issues & Solutions

### Issue: "linker not found" on Windows target

**Solution:**
```bash
# Install LLVM with Windows support
brew install llvm
export PATH="/opt/homebrew/opt/llvm/bin:$PATH"

# Or use Wine for Windows executables
brew install wine-stable
```

### Issue: "cannot find -lc" on Linux target

**Solution:**
```bash
# Install cross-compilation toolchain
brew tap messense/macos-cross-toolchains
brew install x86_64-unknown-linux-gnu

# Set linker in ~/.cargo/config.toml:
[target.x86_64-unknown-linux-gnu]
linker = "x86_64-unknown-linux-gnu-gcc"
```

### Issue: Node.js binary incompatible with target platform

**Solution:**
Download platform-specific Node.js binaries:
```bash
# Windows (from macOS/Linux)
curl -O https://nodejs.org/dist/v20.11.1/node-v20.11.1-win-x64.zip
unzip node-v20.11.1-win-x64.zip
cp node-v20.11.1-win-x64/node.exe desktop-ui/src-tauri/bin/

# Linux (from macOS/Windows)
curl -O https://nodejs.org/dist/v20.11.1/node-v20.11.1-linux-x64.tar.xz
tar xf node-v20.11.1-linux-x64.tar.xz
cp node-v20.11.1-linux-x64/bin/node desktop-ui/src-tauri/bin/
```

### Issue: Go CLI binary is macOS-only

The Makefile builds the Go CLI for your current platform. For cross-platform builds:

```bash
# Build CLI for Windows
GOOS=windows GOARCH=amd64 go build -o ts2go-cli.exe ./cmd/ts2go

# Build CLI for Linux
GOOS=linux GOARCH=amd64 go build -o ts2go-cli ./cmd/ts2go

# Build CLI for macOS
GOOS=darwin GOARCH=amd64 go build -o ts2go-cli ./cmd/ts2go
GOOS=darwin GOARCH=arm64 go build -o ts2go-cli ./cmd/ts2go
```

## Using Docker for Isolated Builds

For maximum reliability, use Docker containers:

```bash
# Use Tauri's official Docker images
docker run --rm -v $(pwd):/app -w /app taurijs/tauri-cli:latest \
  tauri build --target x86_64-unknown-linux-gnu
```

## CI/CD Setup

For automated cross-platform builds, see `.github/workflows/` examples:

### GitHub Actions Matrix Build

```yaml
strategy:
  matrix:
    platform: [macos-latest, ubuntu-latest, windows-latest]
    
runs-on: ${{ matrix.platform }}
steps:
  - uses: actions/checkout@v3
  - uses: actions/setup-node@v3
  - uses: dtolnay/rust-toolchain@stable
  - run: make build-desktop
```

## Output Locations

After successful builds:

```
release/stable/
├── macos/
│   └── TS2Go-Desktop-v0.5.1.app
├── windows/
│   └── TS2Go-Desktop-v0.5.1.msi
└── linux/
    ├── TS2Go-Desktop-v0.5.1.AppImage
    └── TS2Go-Desktop-v0.5.1.deb
```

## Recommended Workflow

### For Development (Current Platform Only)
```bash
make dev-desktop  # Fast, hot-reload enabled
```

### For Testing (Debug Build)
```bash
make build-desktop-debug  # Unoptimized, with symbols
```

### For Release (All Platforms)
```bash
# One-time setup
rustup target add x86_64-pc-windows-msvc
rustup target add x86_64-unknown-linux-gnu

# Build all platforms
make build-desktop-all
```

## Platform-Specific Notes

### macOS
- **Universal Binary**: To support both Intel and Apple Silicon, build separately:
  ```bash
  rustup target add x86_64-apple-darwin aarch64-apple-darwin
  ```
- **Code Signing**: Required for distribution outside App Store
- **Notarization**: Required for Gatekeeper approval

### Windows
- **MSI Installer**: Default bundle format
- **NSIS**: Alternative installer (configure in `tauri.conf.json`)
- **Code Signing**: Optional but recommended for SmartScreen bypass

### Linux
- **AppImage**: Universal format, no installation required
- **DEB**: Debian/Ubuntu package format
- **RPM**: Red Hat/Fedora package format (configure in `tauri.conf.json`)

## Troubleshooting

### Enable Verbose Output
```bash
cd desktop-ui && npm run tauri:build -- --verbose
```

### Check Tauri CLI Version
```bash
npm list @tauri-apps/cli
# Should be: @tauri-apps/cli@2.9.4
```

### Verify Rust Toolchain
```bash
rustc --version
cargo --version
rustup show
```

## Additional Resources

- [Tauri Cross-Compilation Guide](https://tauri.app/v1/guides/building/cross-platform)
- [Rust Cross-Compilation Guide](https://rust-lang.github.io/rustup/cross-compilation.html)
- [Go Cross-Compilation](https://golang.org/doc/install/source#environment)
- [TS2Go Build Issues](https://github.com/el-j/ts2go/issues)
