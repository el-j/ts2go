# File Handling & Auto-Save Features

This guide covers the new file handling features added to the TS2Go Desktop UI, including auto-save, keyboard shortcuts, and backup management.

## Table of Contents

1. [Auto-Save](#auto-save)
2. [Keyboard Shortcuts](#keyboard-shortcuts)
3. [Unsaved Changes Protection](#unsaved-changes-protection)
4. [File Backups](#file-backups)
5. [Settings Configuration](#settings-configuration)

---

## Auto-Save

The desktop UI now includes an automatic save feature that saves your files after a period of inactivity.

### How It Works

- When you edit a file, a timer starts
- After the configured delay (default: 3 seconds), the file is automatically saved
- If you continue editing, the timer resets
- Auto-save respects the settings you configure

### Enabling/Disabling Auto-Save

1. Navigate to **Settings** (gear icon in the top-right)
2. Go to the **Application** tab
3. Toggle the **Auto Save** switch
4. Adjust the **Auto Save Delay** slider (1-10 seconds)

### Visual Indicators

- **Unsaved indicator**: A dot (•) appears next to the file name when there are unsaved changes
- **Auto-save notification**: Files are saved silently in the background (no popup)

---

## Keyboard Shortcuts

The UI supports comprehensive keyboard shortcuts to speed up your workflow.

### File Operations

| Shortcut | Action | Description |
|----------|--------|-------------|
| `Ctrl+S` / `Cmd+S` | Save File | Save the currently open file |
| `Ctrl+Shift+S` / `Cmd+Shift+S` | Save All | Save all open files with unsaved changes |

### Editor Operations

| Shortcut | Action | Description |
|----------|--------|-------------|
| `Ctrl+F` / `Cmd+F` | Find | Open find dialog in editor |
| `Ctrl+H` / `Cmd+H` | Replace | Open find and replace dialog |
| `Ctrl+/` / `Cmd+/` | Toggle Comment | Comment/uncomment selected lines |
| `Ctrl+D` / `Cmd+D` | Duplicate Line | Duplicate current line or selection |
| `Alt+Up/Down` | Move Line | Move current line up or down |
| `Ctrl+Z` / `Cmd+Z` | Undo | Undo last change |
| `Ctrl+Shift+Z` / `Cmd+Shift+Z` | Redo | Redo last undone change |

### Transpilation

| Shortcut | Action | Description |
|----------|--------|-------------|
| `Ctrl+T` / `Cmd+T` | Transpile Current | Transpile the currently open file |
| `Ctrl+Shift+T` / `Cmd+Shift+T` | Transpile All | Transpile all files in the project |

### Navigation

| Shortcut | Action | Description |
|----------|--------|-------------|
| `Ctrl+P` / `Cmd+P` | Quick Open | Open file by name |
| `Ctrl+B` / `Cmd+B` | Toggle Sidebar | Show/hide the file tree sidebar |
| `Ctrl+Tab` | Next Tab | Switch to next open file |
| `Ctrl+Shift+Tab` | Previous Tab | Switch to previous open file |

### Help

| Shortcut | Action | Description |
|----------|--------|-------------|
| `?` | Keyboard Shortcuts | Show all keyboard shortcuts |
| `Ctrl+Shift+P` / `Cmd+Shift+P` | Command Palette | Open command palette |

### Viewing All Shortcuts

Click the **?** (help) button in the top-right corner to view a complete list of keyboard shortcuts organized by category.

---

## Unsaved Changes Protection

The UI protects you from accidentally losing unsaved work.

### Features

1. **Browser Tab Close Warning**
   - When you try to close the browser tab with unsaved changes, you'll see a confirmation dialog
   - This works even if you accidentally hit `Cmd+W` or close the window

2. **Unsaved Changes Dialog**
   - Shows a list of all files with unsaved changes
   - Provides three options:
     - **Save All**: Save all files and proceed
     - **Don't Save**: Discard all changes and proceed
     - **Cancel**: Return to editing

### When It Appears

- Closing the browser tab/window
- Navigating away from the project view
- Refreshing the page

---

## File Backups

The UI can automatically create backup copies of your files before saving.

### How It Works

1. Before saving a file, a backup is created in the backup directory
2. Backups are timestamped to track versions: `filename.ts.2024-11-16T14-30-00.backup`
3. Old backups are kept (automatic cleanup not yet implemented)

### Enabling Backups

1. Navigate to **Settings** → **Application**
2. Toggle **Enable Backups**
3. Configure the **Backup Location** (default: `./.backups`)
4. Click the folder icon to browse and select a directory

### Backup File Naming

```
<original-filename>.<timestamp>.backup
```

Example: `main.ts.2024-11-16T14-30-00-123Z.backup`

### Restoring from Backup

To restore a file from backup:

1. Navigate to your backup directory
2. Find the backup file you want to restore
3. Copy the content
4. Paste it into your current file in the editor

---

## Settings Configuration

### Application Settings

Navigate to **Settings** → **Application** to configure:

#### Auto Save Settings

- **Auto Save**: Enable/disable automatic file saving
- **Auto Save Delay**: Time to wait after last edit before saving (1000-10000 ms)
  - Recommended: 3000 ms (3 seconds)
  - Shorter delays = more frequent saves, less chance of data loss
  - Longer delays = fewer disk writes, better for slow storage

#### Backup Settings

- **Enable Backups**: Create backup copies when saving files
- **Backup Location**: Directory where backups are stored
  - Default: `./.backups` (relative to project)
  - Can be absolute path: `/Users/username/ts2go-backups`

#### Other Settings

- **Theme**: Light, Dark, or System (follows OS preference)
- **Font Size**: Editor font size (12-18px)
- **Default Output Directory**: Where transpiled Go files are saved

---

## Tips & Best Practices

### Auto-Save

1. **Enable for peace of mind**: Auto-save protects against crashes and accidental closures
2. **Adjust delay**: If you find auto-save triggers too often, increase the delay
3. **Pair with backups**: Use both auto-save and backups for maximum protection

### Keyboard Shortcuts

1. **Learn incrementally**: Start with `Ctrl+S` and `Ctrl+Shift+S`, add more as needed
2. **Reference panel**: Press `?` anytime to see all shortcuts
3. **Muscle memory**: Practice shortcuts to speed up your workflow

### Backups

1. **Enable during critical work**: Turn on backups when working on important changes
2. **Cleanup regularly**: Backup files accumulate; delete old ones periodically
3. **External backups**: Backups are local only; use git or external backup for safety

### Unsaved Changes

1. **Trust the warning**: If you see the unsaved changes dialog, review it carefully
2. **Save All is safe**: When in doubt, choose "Save All" to keep all changes
3. **Don't Save is destructive**: Only use "Don't Save" if you're certain you want to discard changes

---

## Troubleshooting

### Auto-Save Not Working

**Symptom**: Files don't auto-save after the delay

**Solutions**:
1. Check if **Auto Save** is enabled in Settings → Application
2. Verify the file shows unsaved indicator (•) next to filename
3. Check browser console for errors (F12 → Console tab)
4. Try manually saving with `Ctrl+S` first

### Keyboard Shortcuts Not Working

**Symptom**: Pressing shortcuts doesn't trigger actions

**Solutions**:
1. Check if focus is in the editor (click in the editor area)
2. Verify you're using the correct modifiers (`Ctrl` on Windows/Linux, `Cmd` on macOS)
3. Check for conflicts with browser shortcuts (some shortcuts may be captured by browser)
4. Try in a different browser to rule out browser-specific issues

### Backup Files Not Created

**Symptom**: No backup files in backup directory

**Solutions**:
1. Verify **Enable Backups** is turned on in Settings
2. Check if backup directory exists and is writable
3. Look for backup files in the configured location
4. Check browser console for permission errors

### Unsaved Changes Dialog Not Appearing

**Symptom**: Can close tab without warning despite unsaved changes

**Solutions**:
1. This feature requires browser support for `beforeunload` event
2. Some browsers (Safari, older Firefox) may not show the dialog
3. Try in Chrome or Edge for best compatibility
4. Manually save before closing as a precaution

---

## Advanced Usage

### Custom Backup Strategy

Create a custom backup workflow:

1. Set **Backup Location** to a cloud-synced folder (Dropbox, Google Drive)
2. Enable **Auto Save** with a longer delay (5-10 seconds)
3. Backups are automatically synced to the cloud
4. Access backup history from any device

### Keyboard Shortcut Cheat Sheet

Print or keep this reference handy:

```
Essential Shortcuts:
- Ctrl+S: Save
- Ctrl+Shift+S: Save All
- Ctrl+T: Transpile
- ?: Help

Editor Power User:
- Ctrl+F: Find
- Ctrl+H: Replace
- Ctrl+/: Comment
- Alt+Up/Down: Move Lines
```

### Integration with Git

Best practice workflow:

1. Enable **Auto Save** for continuous local saves
2. Enable **Backups** for rollback capability
3. Commit to git regularly for version history
4. Use `.gitignore` to exclude backup directory:
   ```
   .backups/
   ```

---

## Feedback & Issues

If you encounter any issues or have suggestions for improvements:

1. Check the troubleshooting section above
2. Review the [GitHub Issues](https://github.com/el-j/ts2go/issues)
3. Open a new issue with:
   - Steps to reproduce
   - Expected vs actual behavior
   - Browser and OS information
   - Screenshots if applicable

---

## Version History

### v2.1.0 (Current)
- ✨ Added auto-save with configurable delay
- ✨ Added comprehensive keyboard shortcuts
- ✨ Added unsaved changes protection dialog
- ✨ Added backup system with timestamped files
- ✨ Added settings UI for all new features
- ✨ Added keyboard shortcuts reference panel
- 🐛 Fixed file tree state persistence

### Previous Versions
See [CHANGELOG.md](../../CHANGELOG.md) for complete history.
