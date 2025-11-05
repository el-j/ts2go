# Release Verification Checklist

Use this checklist when creating or verifying a release.

## Pre-Release Checks

- [ ] All tests passing on target branch
- [ ] CHANGELOG.md updated (stable releases only)
- [ ] VERSION file updated (if using auto-version)
- [ ] Dependencies up to date and secure
- [ ] Documentation reflects new features

## Release Artifacts Verification

### CLI Binaries (5 for stable/beta, 3 for alpha)
- [ ] `ts2go-<version>-linux-amd64.tar.gz`
- [ ] `ts2go-<version>-linux-arm64.tar.gz` (stable/beta only)
- [ ] `ts2go-<version>-darwin-amd64.tar.gz` (stable/beta only)
- [ ] `ts2go-<version>-darwin-arm64.tar.gz`
- [ ] `ts2go-<version>-windows-amd64.zip`

### Desktop Applications (all release types)
- [ ] macOS: `TS2Go-Desktop_<version>_universal.dmg`
- [ ] Linux: `TS2Go-Desktop_<version>_amd64.AppImage`
- [ ] Linux: `TS2Go-Desktop_<version>_amd64.deb`
- [ ] Linux: `TS2Go-Desktop_<version>_amd64.rpm`
- [ ] Windows: `TS2Go-Desktop_<version>_x64_en-US.msi`
- [ ] Windows: `TS2Go-Desktop_<version>_x64-setup.exe`

### Other Files
- [ ] `checksums.txt` with SHA256 hashes for all CLI binaries

### Docker (stable only)
- [ ] Image published to ghcr.io/el-j/ts2go
- [ ] Tags: version, major.minor, major, latest

## Release Metadata

- [ ] Release notes generated correctly
- [ ] Tag name follows format (v1.0.0 / v0.1.0-beta.N / v0.0.0-alpha.branch.sha)
- [ ] Release marked as prerelease (beta/alpha) or stable correctly
- [ ] Release published (not draft)

## Workflow Jobs Status

- [ ] create-release: ✓ Success
- [ ] build-cli (all matrix jobs): ✓ Success
- [ ] build-desktop (all matrix jobs): ✓ Success
- [ ] generate-checksums: ✓ Success
- [ ] publish-docker (stable only): ✓ Success

## Post-Release Verification

- [ ] Download and test CLI binary on at least one platform
- [ ] Download and test desktop app on at least one platform
- [ ] Verify checksums match downloaded files
- [ ] Check release announcement (if applicable)
- [ ] Monitor for issues reported by early adopters

## Common Issues

| Issue | Solution |
|-------|----------|
| Missing desktop artifacts | Check tauri-action used `releaseId` parameter |
| Missing CLI artifacts | Check each platform's matrix job completed |
| Incomplete checksums | Ensure generate-checksums ran after all builds |
| Wrong release type | Verify prerelease flag set correctly |
| Missing LICENSE in archive | Ensure LICENSE file exists in repo root |

## Emergency Rollback

If critical issues found after release:

1. Mark release as draft (hide from users)
2. Create hotfix branch
3. Fix issue and test thoroughly
4. Create new release with patch version
5. Delete problematic release if necessary

## Notes

- All builds run in parallel within their job
- Desktop builds wait for CLI builds to complete
- Checksums generated after all other artifacts uploaded
- Each release type has platform-specific optimizations (alpha builds fewer platforms)
