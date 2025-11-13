# GitHub Actions Workflows

## build-desktop.yml

Automated multi-platform desktop app builds.

### Triggers
- **Push** to `main` or `develop` branches
- **Pull requests** to `main`
- **Tags** starting with `v*` (e.g., `v0.5.1`)
- **Manual** workflow dispatch

### Platforms Built
- macOS (Apple Silicon) - `aarch64-apple-darwin`
- Windows (x64) - `x86_64-pc-windows-msvc`
- Linux (x64) - `x86_64-unknown-linux-gnu`

### Artifacts
Each platform build produces:
- **macOS**: `.app` bundle (zipped for release)
- **Windows**: `.msi` installer
- **Linux**: `.AppImage` and `.deb` packages

### Viewing Build Results

1. Go to **Actions** tab in GitHub
2. Click on latest workflow run
3. Scroll to **Artifacts** section
4. Download `ts2go-desktop-{platform}` artifact

### Creating a Release

```bash
# Tag the release
git tag v0.5.2
git push origin v0.5.2

# Workflow automatically:
# 1. Builds all platforms
# 2. Creates GitHub release
# 3. Uploads platform binaries
```

### Local Testing

To test the workflow locally before pushing:

```bash
# Install act (GitHub Actions local runner)
brew install act

# Run workflow locally
act push
```

### Modifying the Workflow

Key sections:
- **matrix.include**: Add/modify build platforms
- **Install dependencies**: Platform-specific system packages
- **Upload artifacts**: Retention and artifact names
- **Release assets**: Files to attach to GitHub releases

### Troubleshooting

**Build fails on Linux:**
- Check `libwebkit2gtk-4.1-dev` is available
- May need to update package names for newer Ubuntu

**Build fails on Windows:**
- Ensure MSVC toolchain is correctly set up
- Check Node.js binary architecture matches target

**Artifacts missing:**
- Verify `release/stable/{platform}/` path is correct
- Check build logs for bundle location

## Adding More Workflows

Create additional workflows for:
- **CI tests**: `test.yml` - Run tests on all platforms
- **Release notes**: `release-notes.yml` - Auto-generate changelog
- **Documentation**: `docs.yml` - Build and deploy docs
