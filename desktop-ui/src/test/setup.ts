import { vi } from 'vitest'

// Mock Tauri API for testing
vi.mock('@tauri-apps/api', () => ({
  invoke: vi.fn()
}))

vi.mock('@tauri-apps/plugin-shell', () => ({
  open: vi.fn()
}))
