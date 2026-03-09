import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useKeyboardShortcuts, type Shortcut } from '../useKeyboardShortcuts'
import { createPinia, setActivePinia } from 'pinia'

describe('useKeyboardShortcuts', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should register and trigger keyboard shortcuts', async () => {
    const mockHandler = vi.fn()
    
    const shortcuts: Shortcut[] = [
      {
        key: 's',
        ctrl: true,
        handler: mockHandler,
        description: 'Save file'
      }
    ]

    useKeyboardShortcuts(shortcuts)

    // Simulate Ctrl+S
    const event = new KeyboardEvent('keydown', {
      key: 's',
      ctrlKey: true,
      bubbles: true,
      cancelable: true
    })
    
    window.dispatchEvent(event)

    expect(mockHandler).toHaveBeenCalled()
  })

  it('should handle shift modifier', async () => {
    const mockHandler = vi.fn()
    
    const shortcuts: Shortcut[] = [
      {
        key: 's',
        ctrl: true,
        shift: true,
        handler: mockHandler,
        description: 'Save all files'
      }
    ]

    useKeyboardShortcuts(shortcuts)

    // Simulate Ctrl+Shift+S
    const event = new KeyboardEvent('keydown', {
      key: 's',
      ctrlKey: true,
      shiftKey: true,
      bubbles: true,
      cancelable: true
    })
    
    window.dispatchEvent(event)

    expect(mockHandler).toHaveBeenCalled()
  })

  it('should not trigger when modifiers do not match', async () => {
    const mockHandler = vi.fn()
    
    const shortcuts: Shortcut[] = [
      {
        key: 's',
        ctrl: true,
        handler: mockHandler,
        description: 'Save file'
      }
    ]

    useKeyboardShortcuts(shortcuts)

    // Simulate just 'S' without Ctrl
    const event = new KeyboardEvent('keydown', {
      key: 's',
      bubbles: true,
      cancelable: true
    })
    
    window.dispatchEvent(event)

    expect(mockHandler).not.toHaveBeenCalled()
  })

  it('should handle multiple shortcuts', async () => {
    const saveHandler = vi.fn()
    const saveAllHandler = vi.fn()
    
    const shortcuts: Shortcut[] = [
      {
        key: 's',
        ctrl: true,
        handler: saveHandler,
        description: 'Save file'
      },
      {
        key: 's',
        ctrl: true,
        shift: true,
        handler: saveAllHandler,
        description: 'Save all files'
      }
    ]

    useKeyboardShortcuts(shortcuts)

    // Simulate Ctrl+S
    const saveEvent = new KeyboardEvent('keydown', {
      key: 's',
      ctrlKey: true,
      bubbles: true,
      cancelable: true
    })
    window.dispatchEvent(saveEvent)

    expect(saveHandler).toHaveBeenCalled()
    expect(saveAllHandler).not.toHaveBeenCalled()

    saveHandler.mockClear()

    // Simulate Ctrl+Shift+S
    const saveAllEvent = new KeyboardEvent('keydown', {
      key: 's',
      ctrlKey: true,
      shiftKey: true,
      bubbles: true,
      cancelable: true
    })
    window.dispatchEvent(saveAllEvent)

    expect(saveAllHandler).toHaveBeenCalled()
    expect(saveHandler).not.toHaveBeenCalled()
  })

  it('should return shortcuts array', () => {
    const shortcuts: Shortcut[] = [
      {
        key: 's',
        ctrl: true,
        handler: vi.fn(),
        description: 'Save file'
      },
      {
        key: 's',
        ctrl: true,
        shift: true,
        handler: vi.fn(),
        description: 'Save all files'
      }
    ]

    const result = useKeyboardShortcuts(shortcuts)

    expect(result.shortcuts).toHaveLength(2)
    expect(result.shortcuts[0].description).toBe('Save file')
    expect(result.shortcuts[1].description).toBe('Save all files')
  })
})
