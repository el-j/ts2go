# Quick Setup: Cross-Platform Builds

## TL;DR - Current Platform Only

```bash
# Build for your current platform (fastest)
make build-desktop

# Output: release/stable/macos/ (or windows/linux)
```

## Multi-Platform Builds (Recommended: GitHub Actions)

Due to platform-specific toolchain requirements, we **recommend using GitHub Actions** for multi-platform builds:

1. **Push to GitHub**
2. **Workflow runs automatically** (see `.github/workflows/build-desktop.yml`)
3. **Download artifacts** for all platforms

### Why GitHub Actions?

- ✅ No local setup required
- ✅ Builds all platforms in parallel
- ✅ Automated on every push/PR
- ✅ Release artifacts automatically uploaded

## Local Multi-Platform Builds (Advanced)

### Option 1: Native Builds (Recommended)

Build on each platform natively:

**On macOS:**
```bash
make build-desktop  # → release/stable/macos/
```

**On Windows:**
```bash
make build-desktop  # → release/stable/windows/
```

**On Linux:**
```bash
make build-desktop  # → release/stable/linux/
```

### Option 2: Docker-based Linux Builds (from macOS/Windows)

```bash
# One-time setup
cargo install cross
# Ensure Docker Desktop is running

# Build for Linux
make build-desktop-linux
```

### Option 3: Experimental Windows Cross-Compilation (from macOS)

```bash
# One-time setup (not fully tested)
cargo install cargo-xwin
brew install --cask wine-stable

# Attempt Windows build
make build-desktop-windows
```

## Build All (with graceful fallbacks)

```bash
make build-desktop-all
# ✅ Builds current platform
# ⚠️  Skips others with helpful messages
```

## Quick Commands

```bash
# See all build options
make help

# Development mode (fast, hot-reload)
make dev-desktop

# Debug build (with symbols, current platform)
make build-desktop-debug

# Production build (current platform)
make build-desktop

# Try building all platforms (graceful fallback)
make build-desktop-all
```

## CI/CD - Automated Multi-Platform Builds

The included GitHub Actions workflow (`.github/workflows/build-desktop.yml`) automatically builds for all platforms:

- **Triggers:** Push to main/develop, tags, or manual dispatch
- **Platforms:** macOS (Apple Silicon), Windows (x64), Linux (x64)
- **Outputs:** Artifacts uploaded for each platform
- **On release tags:** Automatically creates GitHub release with all binaries

**To use:**
1. Push code to GitHub
2. Check Actions tab for build progress
3. Download artifacts from completed workflow

## Expected Build Times

- **macOS**: ~2-3 minutes
- **Windows**: ~3-5 minutes (cross-compile)
- **Linux**: ~3-5 minutes (cross-compile)
- **All platforms**: ~8-12 minutes

## Common Issues

### "linker not found" error
```bash
brew install llvm
export PATH="/opt/homebrew/opt/llvm/bin:$PATH"
```

### Wrong Node.js architecture
Download platform-specific binaries from [nodejs.org](https://nodejs.org/dist/v20.11.1/)

### Go CLI is wrong architecture
The Makefile auto-detects your platform. For explicit cross-compilation:
```bash
GOOS=windows GOARCH=amd64 go build -o ts2go-cli.exe ./cmd/ts2go
GOOS=linux GOARCH=amd64 go build -o ts2go-cli ./cmd/ts2go
```

## Full Documentation

See [CROSS_COMPILATION.md](./CROSS_COMPILATION.md) for:
- Detailed platform setup instructions
- CI/CD configuration examples
- Troubleshooting guide
- Platform-specific bundling notes

## Output Artifacts

### macOS
- `TS2Go-Desktop-v0.5.1.app` (139 MB)
- Universal Binary (Intel + Apple Silicon)

### Windows
- `TS2Go-Desktop-v0.5.1.msi` (~150 MB)
- Includes bundled Node.js runtime

### Linux
- `TS2Go-Desktop-v0.5.1.AppImage` (~140 MB)
- `TS2Go-Desktop-v0.5.1.deb` (~140 MB)
- Portable and installable formats

## Development Workflow

```bash
# 1. Make code changes
# 2. Test locally
make dev-desktop

# 3. Build for testing
make build-desktop-debug

# 4. Build release for all platforms
make build-desktop-all

# 5. Distribute
# Files in: release/stable/{macos,windows,linux}/
```

## CI/CD Integration

GitHub Actions workflow example:

```yaml
name: Build All Platforms
on: [push, release]

jobs:
  build:
    strategy:
      matrix:
        os: [macos-latest, ubuntu-latest, windows-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v3
      - uses: dtolnay/rust-toolchain@stable
      - run: make build-desktop
      - uses: actions/upload-artifact@v3
        with:
          name: ts2go-${{ matrix.os }}
          path: release/stable/
```

## Need Help?

- 📖 Full docs: [CROSS_COMPILATION.md](./CROSS_COMPILATION.md)
- 🐛 Issues: [GitHub Issues](https://github.com/el-j/ts2go/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/el-j/ts2go/discussions)
