# TS2Go Next Steps

**Last Updated:** November 13, 2025  
**Current Version:** 0.7.0-beta  
**Focus:** Immediate priorities for the next 2-4 weeks

---

## 🎯 Overview

This document outlines the immediate next steps for TS2Go development. For the complete roadmap to v1.0.0, see [ROADMAP_TO_1.0.0.md](ROADMAP_TO_1.0.0.md).

**Current State:**
- ✅ Core transpilation working (75-80%)
- ✅ Modern JS features implemented (Phase 2 complete)
- ✅ Desktop UI navigation and settings complete
- 🔄 Desktop UI at 60% completion (Phase 3 Week 6 complete)
- 📋 Advanced features pending (Phase 3 Week 7-8)

---

## 📋 Immediate Priorities (Next 2 Weeks)

### Week 1: Desktop UI - Advanced Editor Features

**Goal:** Complete syntax error highlighting and improve editor experience

#### 1. Syntax Error Highlighting
**Status:** Not Started  
**Priority:** HIGH  
**Estimated Time:** 3-4 days

**Tasks:**
- [ ] Integrate TypeScript language server into CodeEditor.vue
- [ ] Display inline error markers in Monaco Editor
- [ ] Show error messages on hover
- [ ] Add error count badge to status bar
- [ ] Highlight problematic code sections

**Files:**
- `desktop-ui/src/components/CodeEditor.vue`
- `desktop-ui/src/utils/diagnostics.ts` (new)
- `desktop-ui/src/stores/editor.ts`

**Success Criteria:**
- TypeScript syntax errors show inline in editor
- Hover shows detailed error messages
- Status bar displays error count
- Errors update in real-time as user types

---

#### 2. Editor Performance Improvements
**Status:** Not Started  
**Priority:** MEDIUM  
**Estimated Time:** 2 days

**Tasks:**
- [ ] Debounce transpilation on input (300ms delay)
- [ ] Add loading indicator during transpilation
- [ ] Optimize Monaco Editor configuration
- [ ] Add keyboard shortcuts for common actions
- [ ] Improve editor responsiveness for large files

**Success Criteria:**
- No lag when typing in large files (>1000 lines)
- Smooth scrolling and editing experience
- Clear feedback during transpilation

---

### Week 2: Desktop UI - Multi-File Project Support

**Goal:** Enable working with multi-file TypeScript projects

#### 3. File Tree View
**Status:** Not Started  
**Priority:** HIGH  
**Estimated Time:** 3 days

**Tasks:**
- [ ] Create FileTree.vue component with collapsible folders
- [ ] Integrate with project store
- [ ] Add file icons based on file type
- [ ] Support folder expansion/collapse
- [ ] Add file selection and navigation
- [ ] Show file count in status bar

**Files:**
- `desktop-ui/src/components/FileTree.vue` (new)
- `desktop-ui/src/stores/workspace.ts` (new)
- `desktop-ui/src/views/ProjectView.vue`

**Success Criteria:**
- File tree displays project structure
- Clicking file opens it in editor
- Visual feedback for selected file
- Smooth navigation between files

---

#### 4. Multi-File Tabs
**Status:** Not Started  
**Priority:** HIGH  
**Estimated Time:** 2 days

**Tasks:**
- [ ] Create FileTabs.vue component
- [ ] Support multiple open files
- [ ] Add close buttons on tabs
- [ ] Support tab switching with keyboard (Ctrl+Tab)
- [ ] Show unsaved indicator on modified files
- [ ] Add "close all" and "close others" actions

**Files:**
- `desktop-ui/src/components/FileTabs.vue` (new)
- `desktop-ui/src/stores/workspace.ts`

**Success Criteria:**
- Multiple files can be open simultaneously
- Easy switching between open files
- Clear indication of active file
- Unsaved changes are tracked

---

## 🚀 Next Phase (Weeks 3-4)

### Phase 3 Week 8: Build & Run Features

**Goal:** Enable building and running transpiled Go code from the UI

#### 5. Build Integration
**Status:** Planned  
**Priority:** HIGH  
**Estimated Time:** 3 days

**Tasks:**
- [ ] Integrate Go compiler in Rust backend
- [ ] Create BuildOutput.vue component
- [ ] Display build progress and results
- [ ] Show compilation errors with file/line links
- [ ] Add "Build Project" button to UI
- [ ] Generate go.mod for projects

**Success Criteria:**
- Can build Go code from UI
- Build output is clearly displayed
- Build errors link back to source files
- Build artifacts are tracked

---

#### 6. Run & Debug
**Status:** Planned  
**Priority:** MEDIUM  
**Estimated Time:** 2 days

**Tasks:**
- [ ] Create RuntimeConsole.vue component
- [ ] Execute compiled Go binaries
- [ ] Capture and display stdout/stderr
- [ ] Add stop/restart controls
- [ ] Support command-line arguments
- [ ] Show execution time and exit code

**Success Criteria:**
- Can run transpiled code from UI
- Output is displayed in real-time
- Clear controls for managing execution
- Graceful handling of crashes

---

## 📊 Progress Tracking

### Completion Status

| Feature | Status | Completion | Target Date |
|---------|--------|------------|-------------|
| Syntax Error Highlighting | 📋 Planned | 0% | Nov 18, 2025 |
| Editor Performance | 📋 Planned | 0% | Nov 20, 2025 |
| File Tree View | 📋 Planned | 0% | Nov 23, 2025 |
| Multi-File Tabs | 📋 Planned | 0% | Nov 25, 2025 |
| Build Integration | 📋 Planned | 0% | Nov 28, 2025 |
| Run & Debug | 📋 Planned | 0% | Nov 30, 2025 |

### Version Targets

- **v0.8.0-beta:** Phase 3 Week 7-8 complete (late November)
- **v0.9.0-rc:** Phase 4 complete (mid-December)
- **v1.0.0:** Phase 5 complete (late December)

---

## 🔧 Technical Considerations

### Dependencies
- **Monaco Editor:** Already integrated, enhance with language server
- **TypeScript Language Server:** Need to integrate for error checking
- **Go Compiler:** Available via system, need Rust bindings
- **Tauri IPC:** Use for build/run communication

### Performance Goals
- Editor response time: <100ms
- Transpilation time: <1s for 10K lines
- File tree rendering: <500ms for 1000 files
- Build time: System-dependent (Go compiler speed)

### User Experience
- Progressive disclosure: Show advanced features when needed
- Clear feedback: Loading states, progress indicators
- Error recovery: Graceful handling of build/run failures
- Keyboard shortcuts: Power user efficiency

---

## 📚 Related Documentation

- [ROADMAP_TO_1.0.0.md](ROADMAP_TO_1.0.0.md) - Complete roadmap to stable release
- [PROJECT_STATUS.md](PROJECT_STATUS.md) - Current project status and metrics
- [KNOWN_ISSUES.md](KNOWN_ISSUES.md) - Known limitations and workarounds
- [Desktop UI Developer Guide](desktop-ui/DEVELOPER_GUIDE.md) - Desktop UI development

---

## 🤝 Contributing

Interested in helping with these next steps? See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

**High-Priority Help Needed:**
1. TypeScript language server integration
2. Monaco Editor performance optimization
3. File tree component implementation
4. Go compiler integration with Rust/Tauri

---

## ✅ Definition of Done

Each feature is considered complete when:
- [ ] Implementation code is written and tested
- [ ] Unit tests pass (80%+ coverage)
- [ ] Integration tests pass
- [ ] UI is responsive and performant
- [ ] Documentation is updated
- [ ] Code review is complete
- [ ] No regressions in existing features

---

## 🎯 Success Metrics

**By End of Phase 3 (Week 8):**
- Desktop UI completion: 90%+
- User can work with multi-file projects
- Full edit-transpile-build-run cycle works
- Performance meets targets
- Zero critical bugs

**User Satisfaction:**
- Smooth editing experience
- Clear error messages
- Fast feedback loops
- Professional UI/UX

---

*This document is updated weekly. Last review: November 13, 2025*
