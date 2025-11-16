# TS2Go Desktop - User Guide

**Version:** 0.1.1  
**Last Updated:** November 14, 2025

## Table of Contents

1. [Introduction](#introduction)
2. [Getting Started](#getting-started)
3. [Features](#features)
4. [Settings](#settings)
5. [Keyboard Shortcuts](#keyboard-shortcuts)
6. [Example Gallery](#example-gallery)
7. [Troubleshooting](#troubleshooting)

---

## Introduction

TS2Go Desktop is a graphical user interface for the TS2Go TypeScript-to-Go transpiler. It provides an intuitive way to transpile TypeScript code to Go without using the command line.

### Key Features

- 🎨 **Modern UI** - Clean, professional interface built with Vue 3 and PrimeVue
- ✨ **Real-time Transpilation** - See your Go code as you write TypeScript
- 📊 **Visual Feedback** - Progress tracking, logs, and error messages
- ⚙️ **Customizable** - Comprehensive settings for personalization
- 📖 **Example Gallery** - Learn from built-in examples
- ⌨️ **Keyboard Shortcuts** - Work faster with keyboard navigation
- 🌓 **Dark/Light Theme** - Choose your preferred color scheme

---

## Getting Started

### Installation

#### From Source
```bash
# Clone the repository
git clone https://github.com/el-j/ts2go.git
cd ts2go/desktop-ui

# Install dependencies
npm install

# Run development server
npm run tauri:dev

# Build production version
npm run tauri:build
```

#### Pre-built Binaries
Download the latest release for your platform:
- **Windows**: `TS2Go-Setup.exe`
- **macOS**: `TS2Go.dmg`
- **Linux**: `ts2go.AppImage` or `ts2go.deb`

### First Launch

1. **Launch the application**
2. **Select a theme** (Settings → Application → Theme)
3. **Explore the Example Gallery** to see what's possible
4. **Start transpiling!**

---

## Features

### 1. Code Editor

The split-pane editor provides:
- **Left Pane**: TypeScript input with syntax highlighting
- **Right Pane**: Generated Go code output
- **Resizable Panes**: Drag the divider to resize
- **Monaco Editor**: Full VS Code editor capabilities

#### Using the Editor

1. Type or paste TypeScript code in the left pane
2. Click **Transpile** Button or press `Ctrl+S`
3. View generated Go code in the right pane
4. Check logs for any warnings or errors

### 2. Project Management

Transpile entire TypeScript projects with full Go module support:

**Features:**
- Open TypeScript projects from your file system
- Browse and select individual files
- Multi-file transpilation with dependency tracking
- **Automatic Go code formatting** with `go fmt`
- Project-wide build, test, and run operations

**Workflow:**
1. Navigate to **Projects** tab
2. Click **Open Project** and select your TypeScript project folder
3. Review the file tree and select files to transpile
4. Click **Transpile All** to convert the entire project
5. Generated Go code is automatically formatted for readability
6. Use **Build**, **Test**, or **Run** buttons to work with the output

**State Persistence** ⭐ NEW in v0.1.1

TS2Go now remembers your transpilation state! When you return to a project, you'll see:

- 📦 **"Previous Transpilation Restored"** banner
- ⏱️ Timestamp of last transpilation
- 📊 Number of files transpiled
- ✅ Build/Test/Run buttons remain enabled

**How It Works:**
- State is automatically saved after successful transpilation
- Persists across app restarts
- Verified before restoration (checks if output directory still exists)
- Click **"Clear State"** to manually reset

**Benefits:**
- No need to re-transpile just to use Build/Test/Run
- Resume work immediately after reopening the app
- State is project-specific (each project tracked separately)

**Go Code Formatting:**

All transpiled Go code is automatically formatted using `go fmt` to ensure:
- ✅ Consistent style following Go conventions
- ✅ Professional, readable code output
- ✅ Standard indentation and spacing
- ✅ Clean diffs for version control

**Phase 0a (Current):** Requires Go installed on your system  
**Phase 0b (Future):** Will use bundled Go compiler (no installation needed)

If formatting fails (e.g., Go not installed), you'll see a warning but transpilation will still complete successfully. The code will be functional but unformatted.

### 3. Example Gallery

Access pre-built examples to learn TypeScript-to-Go transpilation:

**Basic Examples:**
- Interface definitions
- Function declarations
- Enum definitions
- Control flow structures

**Intermediate Examples:**
- Classes with inheritance
- Advanced type usage

**Advanced Examples:**
- Async/await with promises
- Complex control flows

#### Using Examples

1. Navigate to **Examples** in the sidebar
2. Filter by category (Basic, Intermediate, Advanced)
3. Click on an example card
4. View the TypeScript and Go code side-by-side
5. Click **Open in Editor** to edit and experiment

### 4. Log Viewer

Monitor transpilation progress and errors:
- **Color-coded levels**: Info (blue), Success (green), Warning (yellow), Error (red)
- **Filtering**: Show/hide specific log levels
- **Search**: Find specific messages
- **Export**: Save logs to file
- **Auto-scroll**: Automatically scroll to latest entry

#### Using Logs

1. Click the list icon (📋) in the toolbar to show logs
2. Use the filter dropdown to show specific levels
3. Search using the search box
4. Click **Export** to save logs
5. Click **Clear** to remove all logs

---

## Settings

Access settings by clicking **Settings** in the sidebar or pressing `Ctrl+,`.

### Application Settings

**Theme**
- **Light**: Bright, high-contrast theme
- **Dark**: Easy on the eyes for night coding
- **System**: Match your OS theme

**Font Size**
- Small (12px)
- Medium (14px) - Default
- Large (16px)
- Extra Large (18px)

**Auto Save**
- Enable/disable automatic saving of changes

**Default Output Directory**
- Where transpiled Go files are saved
- Click folder icon to browse

### Project Settings

**Go Module Name**
- Override the default Go module name
- Example: `github.com/username/project`

**Exclude Patterns**
- Comma-separated glob patterns for files to exclude
- Default: `node_modules, **/*.test.ts, dist`

**Include Patterns**
- Comma-separated glob patterns for files to include
- Default: `**/*.ts, **/*.tsx`

**Go Compiler Configuration** ⭐ NEW in v0.1.1

Configure which Go compiler TS2Go should use for Build, Test, Run, and Format operations:

**Go Binary Source:**
- **System Go (from PATH)** - Default. Uses Go installed on your system
- **Custom Path** - Specify a custom Go binary location

**Features:**
- 🔍 **Auto-detection** - Instantly checks if Go is available
- ✅ **Version display** - Shows detected Go version and path
- 📁 **File picker** - Browse to select custom Go binary
- ⚠️ **Clear errors** - Helpful messages if Go not found

**How to Configure:**

1. Open **Settings** → **Project** tab
2. Locate **Go Compiler Configuration** section
3. Select your preferred source:
   - Choose **"System Go"** if you have Go in your PATH
   - Choose **"Custom Path"** to specify a different Go installation
4. Click **"Detect Go Installation"** to verify
5. You should see: ✅ **Go Compiler Detected** with version info

**If Go Is Not Detected:**

See the [Go Installation Guide](../../docs/GO_INSTALLATION_GUIDE.md) for step-by-step instructions to install Go on your system.

**Quick Install:**
- **macOS:** `brew install go`
- **Windows:** Download from https://go.dev/dl/
- **Linux:** `sudo apt install golang-go`

After installing, return to Settings and click "Detect Go Installation" again.

### Editor Settings

**Tab Size**
- 2, 4, or 8 spaces per tab
- Default: 4 spaces

**Word Wrap**
- Wrap long lines in the editor
- Default: Enabled

**Line Numbers**
- Show line numbers in the editor
- Default: Enabled

**Minimap**
- Show code overview minimap
- Default: Enabled

**Auto Format on Save**
- Automatically format code when saving
- Default: Enabled

### Resetting Settings

Click **Reset to Defaults** Button to restore all settings to their default values.

---

## Keyboard Shortcuts

Speed up your workflow with keyboard shortcuts.

### Global Shortcuts

| Keys | Action |
|------|--------|
| `Ctrl+S` | Transpile current file |
| `Ctrl+O` | Open file dialog |
| `Ctrl+L` | Toggle log viewer |
| `Ctrl+/` | Show keyboard shortcuts |
| `Ctrl+,` | Open settings |

### Editor Shortcuts

| Keys | Action |
|------|--------|
| `Ctrl+F` | Find in file |
| `Ctrl+H` | Find and replace |
| `Ctrl+W` | Close current tab |
| `Ctrl+Shift+W` | Close all tabs |

### Navigation Shortcuts

| Keys | Action |
|------|--------|
| `Ctrl+1` | Go to Home |
| `Ctrl+2` | Go to Projects |
| `Ctrl+3` | Go to Editor |
| `Ctrl+4` | Go to Examples |
| `Ctrl+5` | Go to Settings |

**Note:** On macOS, use `Cmd` instead of `Ctrl`.

---

## Example Gallery

### Available Examples

#### 1. Interface Definition
Learn how interfaces are transpiled to Go structs.
- **Category**: Basic
- **Concepts**: Interfaces, types, JSON tags

#### 2. Class with Inheritance
See how TypeScript classes map to Go structs and methods.
- **Category**: Intermediate
- **Concepts**: Classes, constructors, methods, inheritance

#### 3. Function Declarations
Understand function transpilation with typed parameters.
- **Category**: Basic
- **Concepts**: Functions, parameters, return types

#### 4. Enum Definitions
Learn enum transpilation for both numeric and string enums.
- **Category**: Basic
- **Concepts**: Enums, constants, switch statements

#### 5. Async/Await with Promises
See how async functions map to Go goroutines and channels.
- **Category**: Advanced
- **Concepts**: Async/await, promises, goroutines, channels

#### 6. Control Flow Structures
Learn transpilation of if/else, loops, and switch statements.
- **Category**: Basic
- **Concepts**: If/else, for loops, switch statements

### Using the Gallery

1. **Browse Examples**: Scroll through the gallery
2. **Filter by Category**: Click Basic, Intermediate, or Advanced
3. **View Details**: Click an example card to see full code
4. **Compare**: See TypeScript and Go side-by-side
5. **Experiment**: Click "Open in Editor" to modify and test

---

## Troubleshooting

### Common Issues

#### 1. Transpilation Fails

**Symptoms:**
- Error messages in logs
- No Go code generated

**Solutions:**
- Check TypeScript syntax errors
- Ensure all types are properly defined
- Review error messages in the log viewer
- Try a simpler example first

#### 2. Settings Not Saving

**Symptoms:**
- Settings reset after restart

**Solutions:**
- Check browser localStorage is enabled (if web version)
- Verify file system permissions
- Try resetting to defaults and re-applying

#### 3. Dark Theme Not Working

**Symptoms:**
- Theme doesn't change
- Colors look wrong

**Solutions:**
- Refresh the application
- Check Settings → Application → Theme
- Try switching between themes
- Verify system theme settings (if using System mode)

#### 4. Keyboard Shortcuts Not Working

**Symptoms:**
- Shortcuts don't trigger actions

**Solutions:**
- Check if focus is in the correct element
- Try clicking in the editor first
- Verify shortcuts in the help dialog (Ctrl+/)
- Check for conflicting system shortcuts

#### 5. Go Compiler Not Found

**Symptoms:**
- Error: "Go Compiler Not Found"
- Build, Test, or Run buttons fail
- Message: "Configure Go in Settings to use Build, Test, and Run features"

**Solutions:**

1. **Check if Go is installed:**
   ```bash
   go version
   ```
   If this command fails, Go is not installed.

2. **Install Go:**
   - See the detailed [Go Installation Guide](../../docs/GO_INSTALLATION_GUIDE.md)
   - **macOS:** `brew install go`
   - **Windows:** Download from https://go.dev/dl/
   - **Linux:** `sudo apt install golang-go`

3. **Configure Go in TS2Go:**
   - Open **Settings** → **Project** tab
   - Locate **Go Compiler Configuration**
   - Select **"System Go (from PATH)"**
   - Click **"Detect Go Installation"**
   - Verify detection succeeds with green ✅

4. **Use Custom Go Path** (if installed in non-standard location):
   - Settings → Project → Go Compiler Configuration
   - Select **"Custom Path"**
   - Click folder icon 📁
   - Browse to your Go binary (e.g., `/usr/local/go/bin/go`)
   - Click **"Detect Go Installation"**

5. **Verify Go is in PATH** (for System Go):
   - **macOS/Linux**: Run `which go` in terminal
   - **Windows**: Run `where go` in Command Prompt
   - If not found, add Go to your system PATH

**Still Having Issues?**
- Check the Output panel logs for detailed error messages
- Verify Go version is 1.18+ (required for generics): `go version`
- Ensure Go binary has execute permissions (macOS/Linux): `chmod +x /path/to/go`

#### 6. Go Code Formatting Warnings

**Symptoms:**
- Warning message: "Failed to format Go code"
- Transpilation succeeds but code looks unformatted

**Solutions:**
- This is the same issue as #5 - Go not found
- Follow steps above to install and configure Go
- Formatting uses the same Go detection system
- After configuring Go, transpile again to see formatted output

3. **Restart the application** after installing Go

4. **Check Go installation**:
   - Navigate to Settings (once Go config UI is available in Phase 2)
   - Click "Detect" to verify Go installation
   - Should show Go version and path

**Note:** Even if formatting fails, your transpiled Go code will still work correctly. Formatting is for readability and style consistency, not functionality.

**Future Enhancement:** Phase 3 will bundle Go compiler with the app, eliminating the need for separate Go installation.

### Getting Help

If you encounter issues not covered here:

1. **Check the logs**: Enable verbose logging in settings
2. **Review documentation**: Visit the GitHub repository
3. **Report issues**: Open an issue on GitHub with:
   - Description of the problem
   - Steps to reproduce
   - Error messages from logs
   - Your operating system and version

### Performance Tips

**For Large Projects:**
- Use exclude patterns to skip unnecessary files
- Enable verbose logging only when debugging
- Close unused tabs in the editor
- Clear logs periodically

**For Slow Transpilation:**
- Check your TypeScript code complexity
- Verify system resources (CPU, memory)
- Close other resource-intensive applications

---

## Additional Resources

### Documentation
- [TS2Go Main Documentation](../README.md)
- [API Reference](../docs/API_REFERENCE.md)
- [Architecture Overview](../docs/ARCHITECTURE.md)
- [Go Installation Guide](../docs/GO_INSTALLATION_GUIDE.md) ⭐ NEW
- [Testing Guide](../TESTING_GUIDE_PHASE1_AND_2.md) ⭐ NEW

### Examples
- [Example Projects](../examples/)
- [Test Fixtures](../tests/fixtures/)

### Community
- [GitHub Repository](https://github.com/el-j/ts2go)
- [Issue Tracker](https://github.com/el-j/ts2go/issues)
- [Discussions](https://github.com/el-j/ts2go/discussions)

---

## Version History

### v0.1.1 (Current - November 15, 2025)

**Phase 0: Go Formatting**
- ✨ Automatic Go code formatting with `go fmt`
- ✅ All transpiled Go code follows standard Go formatting conventions
- ⚠️ Graceful error handling for formatting failures
- 📊 UI feedback showing formatted file counts

**Phase 1: State Persistence**
- ✨ **NEW:** Transpilation state persistence with localStorage
- 📦 "Previous Transpilation Restored" banner in ProjectView
- ⏱️ Shows timestamp and file count of last transpilation
- ✅ Build/Test/Run buttons remain enabled after app restart
- 🔄 Per-project state tracking
- 🧹 Manual "Clear State" button

**Phase 2: Go Configuration**
- ✨ **NEW:** Go Compiler Configuration UI in Settings
- 🔍 Smart Go detection with 3-tier priority (custom → bundled → system)
- ⚙️ Go Binary Source dropdown (System/Custom)
- 📁 File picker for custom Go path selection
- ✅ Instant Go detection with version/path display
- 🎯 Auto-detection on Settings page load
- ⚠️ Clear error messages for Go not found

**Infrastructure Improvements**
- 🛠️ All 6 go commands (Build, Test, Run) use smart Go detection
- 🚫 No more silent failures when Go not in PATH
- 📝 Improved error handling with actionable messages
- 🔧 2 new Tauri commands: detect_go_installation, check_directory_exists
- 📚 Comprehensive [Go Installation Guide](../docs/GO_INSTALLATION_GUIDE.md)

### v0.1.0 (November 2025)
- Initial release of Desktop UI
- Basic transpilation support
- Settings panel with theme customization
- Example gallery with 6 examples
- Keyboard shortcuts support
- Log viewer with filtering
- Monaco editor integration

---

**Need more help?** Open an issue on [GitHub](https://github.com/el-j/ts2go/issues) or check the [documentation](../docs/).
