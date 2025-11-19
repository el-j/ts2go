# TS2Go Desktop App

This is the desktop application entry point that uses the shared UI package.

## Structure

- UI: Imports from `@ts2go/ui-shared` package
- Backend: Uses Tauri backend from `desktop/tauri-backend/`

## Development

```bash
# Start Tauri development mode
npm run tauri:dev

# Build desktop application
npm run tauri:build
```

## Configuration

The Tauri backend configuration is located in `desktop/tauri-backend/tauri.conf.json`.
The UI is built from the shared package with `--mode tauri` flag.
