# CI/CD Workflow Optimization Plan

## Current Issues

### Workflow Duplication Analysis
1. **ci-main.yml** - Runs on `main` push
2. **ci-develop.yml** - Runs on `develop` push + PRs to main
3. **ci-feature.yml** - Runs on `feature/**` and `copilot/**` branches
4. **build-desktop.yml** - Runs on main, develop, tags, PRs
5. **release.yml** - Runs on main push + tags
6. **release-beta.yml** - Runs on develop push
7. **release-alpha.yml** - Runs on feature branches
8. **desktop-ui.yml** - Runs on PRs (desktop path changes)
9. **auto-version.yml** - Runs on main push

**Problems:**
- Main branch push triggers: ci-main, build-desktop, release, auto-version (4 workflows!)
- Develop push triggers: ci-develop, build-desktop, release-beta (3 workflows!)
- Feature push triggers: ci-feature, build-desktop, release-alpha (3 workflows!)
- PR triggers: ci-develop/ci-main, build-desktop, desktop-ui (2-3 workflows!)

### Cost Impact
- **Redundant builds**: Desktop built in ci-* workflows AND build-desktop
- **Duplicate releases**: release.yml + build-desktop both create releases on tags
- **Wasted cache**: Each workflow rebuilds dependencies separately
- **Concurrent waste**: No concurrency groups = parallel duplicate runs

---

## Optimized Architecture

### Workflow Structure (5 workflows total)

```
┌─────────────────────────────────────────────────────┐
│  1. ci.yml - Smart CI (PRs + Branch Pushes)        │
│     - Run tests, linting, type checks               │
│     - Cache: Go, Node, Rust                          │
│     - Concurrency: per PR/branch                     │
│     - Triggers: PRs, push to any branch              │
└─────────────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────┐
│  2. auto-version.yml - Version Bump (Main Only)    │
│     - Bump VERSION file                              │
│     - Create git tag (v1.2.3)                        │
│     - Triggers: push to main (except [skip-release]) │
└─────────────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────┐
│  3. release.yml - Release Creation (Tags Only)     │
│     - Create GitHub release                          │
│     - Generate changelog                             │
│     - Trigger desktop + web builds                   │
│     - Triggers: on tag creation (v*)                 │
└─────────────────────────────────────────────────────┘
              │
        ┌─────┴─────┐
        │           │
        ▼           ▼
┌──────────────┐  ┌──────────────────────────────────┐
│ 4. Desktop   │  │ 5. Web Deploy (Main Only)        │
│    Build     │  │    - Build ui-shared (web mode)  │
│    (Tag)     │  │    - Deploy to GitHub Pages       │
│              │  │    - Triggers: push to main       │
│    - Mac     │  └──────────────────────────────────┘
│    - Windows │
│    - Linux   │
│    - Assets  │
└──────────────┘
```

### Workflow Details

#### 1. **ci.yml** - Unified CI Testing
```yaml
Triggers:
  - pull_request (all targets)
  - push (all branches except main)
  
Jobs:
  - lint-and-test:
      Go: fmt, vet, staticcheck, tests
      TypeScript: build ui-shared (both modes)
      Cache: Go modules, Node modules, Rust
  
Concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: true
  
Cost savings: Replaces ci-main, ci-develop, ci-feature (3→1)
```

#### 2. **auto-version.yml** - Keep Current (Already Optimized)
```yaml
Triggers:
  - push to main (only when [skip-release] NOT in commit)
  
Jobs:
  - Bump VERSION file
  - Create tag (triggers release.yml)
  
Concurrency:
  group: version-${{ github.sha }}
  
No changes needed
```

#### 3. **release.yml** - Streamlined Release
```yaml
Triggers:
  - tags: v*
  
Jobs:
  - create-release: Generate changelog, create GitHub release
  - call-desktop-build: Reusable workflow call
  - call-web-deploy: Reusable workflow call (only if tag from main)
  
Concurrency:
  group: release-${{ github.ref }}
  
Cost savings: Removes release-alpha, release-beta, build-desktop duplication
```

#### 4. **build-desktop.yml** - Tag-Only Desktop Builds
```yaml
Triggers:
  - workflow_call (from release.yml)
  - workflow_dispatch (manual)
  
Jobs:
  - build-macos:
      Target: universal (arm64 + x86_64)
      Cache: Rust, Go, Node
      Output: DMG, .app.zip
      
  - build-windows:
      Target: x86_64-msvc
      Cache: Rust, Go, Node
      Output: MSI, .exe
      
  - build-linux:
      Target: x86_64-gnu
      Cache: Rust, Go, Node
      Output: AppImage, .deb
      
  - upload-release-assets:
      Needs: [build-macos, build-windows, build-linux]
      
Concurrency:
  group: desktop-build-${{ github.ref }}
  
Cost savings: No more builds on every PR/push
```

#### 5. **deploy-web.yml** - GitHub Pages + Backend
```yaml
Triggers:
  - push to main (paths: packages/ui-shared/**)
  - workflow_call (from release.yml)
  
Jobs:
  - build-web:
      - npm install (packages/ui-shared)
      - npm run build (web mode)
      - Output: dist-web/
      
  - deploy-pages:
      Needs: build-web
      - actions/upload-pages-artifact
      - actions/deploy-pages
      Environment: github-pages
      
  - deploy-backend (TODO):
      Strategy: Cloud Functions or Uberspace
      
Concurrency:
  group: web-deploy
  cancel-in-progress: true
  
New workflow for web deployment
```

---

## Backend Deployment Strategy

### Option A: Cloud Functions (Recommended for MVP)
**Platform**: Google Cloud Functions / AWS Lambda / Vercel Edge
- **Pros**: Auto-scaling, pay-per-use, zero maintenance
- **Cons**: Cold starts, 60s timeout limits
- **Cost**: ~$0-5/month for low traffic
- **Best for**: API endpoints, job queue triggers

### Option B: Uberspace (Recommended for Testing)
**Platform**: https://uberspace.de/
- **Pros**: €5/month, full control, persistent processes, websockets
- **Cons**: Manual deployment, single server
- **Cost**: €5/month flat
- **Best for**: Development, early testing, full SaaS backend

### Recommended Approach
1. **Phase 1**: Deploy static web to GitHub Pages (free)
2. **Phase 2**: Backend API to Uberspace (€5/month, easier deployment)
3. **Phase 3**: Migrate to Cloud Functions when traffic grows

**Backend Requirements:**
- PostgreSQL (Uberspace: built-in, Cloud: managed DB)
- Redis (Uberspace: built-in, Cloud: separate service)
- MinIO/S3 (Uberspace: local storage, Cloud: S3)
- WebSocket (Uberspace: direct support, Cloud: separate service)

**Deployment Workflow:**
```yaml
deploy-backend:
  runs-on: ubuntu-latest
  steps:
    - name: Deploy to Uberspace via SSH
      run: |
        ssh user@host 'cd ts2go-backend && git pull && supervisorctl restart ts2go'
```

---

## Migration Steps

### Phase 1: Consolidate CI Workflows
1. ✅ Create new `ci.yml` (replaces ci-main, ci-develop, ci-feature)
2. ✅ Add concurrency groups to all workflows
3. ✅ Delete redundant ci-* files

### Phase 2: Optimize Release Process
1. ✅ Remove release-alpha.yml, release-beta.yml
2. ✅ Update release.yml to call reusable workflows
3. ✅ Update build-desktop.yml to only run on workflow_call

### Phase 3: Add Web Deployment
1. ✅ Create deploy-web.yml for GitHub Pages
2. ✅ Add Uberspace backend deployment (manual SSH for now)
3. ✅ Update README with deployment instructions

### Phase 4: Verify & Test
1. ✅ Test PR workflow (should run ci.yml only)
2. ✅ Test main push (should run auto-version → release → desktop + web)
3. ✅ Verify no duplicate builds

---

## Cost Savings Estimate

### Before Optimization
- Main push: 4 workflows × 3 jobs × 10min = 120 minutes
- PR: 3 workflows × 2 jobs × 8min = 48 minutes
- Feature push: 3 workflows × 2 jobs × 10min = 60 minutes

**Monthly estimate** (50 PRs, 20 main pushes, 10 feature pushes):
- Total: (50×48) + (20×120) + (10×60) = 5,400 minutes
- Cost: ~$0.008/min = **$43.20/month**

### After Optimization
- Main push: 1 version + 1 release (desktop + web) = 35 minutes
- PR: 1 ci workflow × 1 job × 6min = 6 minutes
- Feature push: 1 ci workflow × 1 job × 6min = 6 minutes

**Monthly estimate** (same traffic):
- Total: (50×6) + (20×35) + (10×6) = 1,060 minutes
- Cost: ~$0.008/min = **$8.48/month**

**Savings: ~80% reduction ($34.72/month)**

---

## Implementation Checklist

- [ ] Create unified `ci.yml` workflow
- [ ] Add concurrency groups to all workflows
- [ ] Create reusable `build-desktop.yml` workflow
- [ ] Create `deploy-web.yml` for GitHub Pages
- [ ] Update `release.yml` to call reusable workflows
- [ ] Delete: ci-main.yml, ci-develop.yml, ci-feature.yml
- [ ] Delete: release-alpha.yml, release-beta.yml
- [ ] Update: desktop-ui.yml (remove Tauri build, keep tests)
- [ ] Add: Backend deployment docs (Uberspace setup)
- [ ] Test: Complete pipeline on test branch
