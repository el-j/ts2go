import { onMounted, onUnmounted } from 'vue'

export interface Shortcut {
  key: string
  ctrl?: boolean
  shift?: boolean
  alt?: boolean
  meta?: boolean
  handler: () => void
  description: string
}

export function useKeyboardShortcuts(shortcuts: Shortcut[]) {
  const handleKeyDown = (event: KeyboardEvent) => {
    for (const shortcut of shortcuts) {
      const keyMatches = event.key.toLowerCase() === shortcut.key.toLowerCase()
      const ctrlMatches = !!shortcut.ctrl === (event.ctrlKey || event.metaKey)
      const shiftMatches = !!shortcut.shift === event.shiftKey
      const altMatches = !!shortcut.alt === event.altKey

      if (keyMatches && ctrlMatches && shiftMatches && altMatches) {
        event.preventDefault()
        shortcut.handler()
        break
      }
    }
  }

  onMounted(() => {
    window.addEventListener('keydown', handleKeyDown)
  })

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown)
  })

  return {
    shortcuts
  }
}

export const DEFAULT_SHORTCUTS: Shortcut[] = [
  {
    key: 's',
    ctrl: true,
    description: 'Save or Transpile',
    handler: () => {} // Will be overridden
  },
  {
    key: 'o',
    ctrl: true,
    description: 'Open File',
    handler: () => {} // Will be overridden
  },
  {
    key: 'w',
    ctrl: true,
    description: 'Close Tab',
    handler: () => {} // Will be overridden
  },
  {
    key: 'f',
    ctrl: true,
    description: 'Find',
    handler: () => {} // Will be overridden
  },
  {
    key: 'h',
    ctrl: true,
    description: 'Replace',
    handler: () => {} // Will be overridden
  },
  {
    key: '/',
    ctrl: true,
    description: 'Show Shortcuts',
    handler: () => {} // Will be overridden
  }
]
