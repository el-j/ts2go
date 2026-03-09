import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useLogsStore } from '../logs'

describe('Logs Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.useFakeTimers()
  })

  it('should initialize with empty logs', () => {
    const store = useLogsStore()
    
    expect(store.logs).toEqual([])
  })

  it('should add info log', () => {
    const store = useLogsStore()
    const timestamp = new Date('2025-01-01T00:00:00.000Z')
    vi.setSystemTime(timestamp)
    
    store.addLog('info', 'Test info message')
    
    expect(store.logs).toHaveLength(1)
    expect(store.logs[0].level).toBe('info')
    expect(store.logs[0].message).toBe('Test info message')
    expect(store.logs[0].timestamp).toEqual(timestamp)
  })

  it('should add multiple logs with different levels', () => {
    const store = useLogsStore()
    
    store.addLog('info', 'Info message')
    store.addLog('warning', 'Warning message')
    store.addLog('error', 'Error message')
    store.addLog('success', 'Success message')
    
    expect(store.logs).toHaveLength(4)
    expect(store.logs[0].level).toBe('info')
    expect(store.logs[1].level).toBe('warning')
    expect(store.logs[2].level).toBe('error')
    expect(store.logs[3].level).toBe('success')
  })

  it('should generate unique IDs for each log', () => {
    const store = useLogsStore()
    
    store.addLog('info', 'Message 1')
    vi.advanceTimersByTime(1)
    store.addLog('info', 'Message 2')
    
    expect(store.logs[0].id).not.toBe(store.logs[1].id)
  })

  it('should clear all logs', () => {
    const store = useLogsStore()
    
    store.addLog('info', 'Message 1')
    store.addLog('error', 'Message 2')
    store.addLog('warning', 'Message 3')
    
    expect(store.logs).toHaveLength(3)
    
    store.clearLogs()
    
    expect(store.logs).toEqual([])
  })

  it('should maintain log order', () => {
    const store = useLogsStore()
    
    store.addLog('info', 'First')
    vi.advanceTimersByTime(100)
    store.addLog('warning', 'Second')
    vi.advanceTimersByTime(100)
    store.addLog('error', 'Third')
    
    expect(store.logs[0].message).toBe('First')
    expect(store.logs[1].message).toBe('Second')
    expect(store.logs[2].message).toBe('Third')
  })
})
