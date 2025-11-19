# TS2Go SaaS Web Frontend

This is the web frontend for the TS2Go SaaS platform, using the shared UI package.

## Structure

- UI: Imports from `@ts2go/ui-shared` package
- Backend: Connects to SaaS backend API via HTTP

## Development

```bash
# Start web development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Configuration

The web app uses the shared UI package built with `--mode web` flag.
It connects to the backend API at `http://localhost:8080` (configurable via `VITE_API_URL`).

## Authentication

Unlike the desktop app, the web version requires authentication:
- Users must register/login with email and password
- JWT tokens are stored in localStorage
- Protected routes require valid authentication
