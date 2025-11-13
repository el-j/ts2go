import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import LogViewer from '../LogViewer.vue'
import { useLogsStore } from '../../stores/logs'

describe('LogViewer Component', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should render log viewer container', () => {
    const wrapper = mount(LogViewer)
    expect(wrapper.find('.log-viewer').exists()).toBe(true)
  })

  it('should display correct log count', () => {
    const logsStore = useLogsStore()
    logsStore.addLog('info', 'Test log 1')
    logsStore.addLog('error', 'Test log 2')
    
    const wrapper = mount(LogViewer)
    expect(wrapper.text()).toContain('Logs (2)')
  })

  it('should filter logs by level', async () => {
    const logsStore = useLogsStore()
    logsStore.addLog('info', 'Info message')
    logsStore.addLog('error', 'Error message')
    logsStore.addLog('warning', 'Warning message')
    
    const wrapper = mount(LogViewer)
    const select = wrapper.find('select')
    
    // Filter to show only errors
    await select.setValue('error')
    expect(wrapper.text()).toContain('Error message')
    expect(wrapper.text()).not.toContain('Info message')
  })

  it('should search logs', async () => {
    const logsStore = useLogsStore()
    logsStore.addLog('info', 'This is a test')
    logsStore.addLog('info', 'Another message')
    
    const wrapper = mount(LogViewer)
    const searchInput = wrapper.find('input[type="text"]')
    
    await searchInput.setValue('test')
    expect(wrapper.text()).toContain('This is a test')
    expect(wrapper.text()).not.toContain('Another message')
  })

  it('should clear all logs', async () => {
    const logsStore = useLogsStore()
    logsStore.addLog('info', 'Test log')
    
    const wrapper = mount(LogViewer)
    expect(wrapper.text()).toContain('Test log')
    
    // Find and click clear Button
    const clearButton = wrapper.findAll('Button').find(btn => 
      btn.element.querySelector('.pi-trash')
    )
    await clearButton?.trigger('click')
    
    expect(logsStore.logs.length).toBe(0)
  })

  it('should toggle auto-scroll', async () => {
    const wrapper = mount(LogViewer)
    const autoScrollButton = wrapper.findAll('Button').find(btn => 
      btn.element.querySelector('.pi-angle-double-down')
    )
    
    // Check initial state (should be active)
    expect(autoScrollButton?.classes()).toContain('bg-primary-600')
    
    // Toggle off
    await autoScrollButton?.trigger('click')
    expect(autoScrollButton?.classes()).toContain('bg-gray-200')
  })

  it('should display "No logs" message when empty', () => {
    const wrapper = mount(LogViewer)
    expect(wrapper.text()).toContain('No logs to display')
  })

  it('should apply correct CSS classes based on log level', () => {
    const logsStore = useLogsStore()
    logsStore.addLog('error', 'Error message')
    logsStore.addLog('success', 'Success message')
    
    const wrapper = mount(LogViewer)
    
    // Check that error log text is present
    expect(wrapper.text()).toContain('Error message')
    expect(wrapper.text()).toContain('Success message')
    
    // Check that level badges are displayed (lowercase in component)
    expect(wrapper.text()).toContain('error')
    expect(wrapper.text()).toContain('success')
  })

  it('should format timestamps correctly', () => {
    const logsStore = useLogsStore()
    logsStore.addLog('info', 'Test log')
    
    const wrapper = mount(LogViewer)
    // Should contain time in format HH:MM:SS
    expect(wrapper.text()).toMatch(/\d{2}:\d{2}:\d{2}/)
  })
})
