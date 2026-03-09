import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import CodeEditor from '../CodeEditor.vue'

describe('CodeEditor Component', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should render editor container', () => {
    const wrapper = mount(CodeEditor, {
      props: {
        modelValue: 'test code'
      },
      global: {
        plugins: [createPinia()]
      }
    })

    expect(wrapper.find('.code-editor').exists()).toBe(true)
  })

  it('should accept modelValue prop', () => {
    const wrapper = mount(CodeEditor, {
      props: {
        modelValue: 'const x = 1;'
      },
      global: {
        plugins: [createPinia()]
      }
    })

    expect(wrapper.props('modelValue')).toBe('const x = 1;')
  })

  it('should accept language prop', () => {
    const wrapper = mount(CodeEditor, {
      props: {
        modelValue: '',
        language: 'go'
      },
      global: {
        plugins: [createPinia()]
      }
    })

    expect(wrapper.props('language')).toBe('go')
  })

  it('should default language to typescript', () => {
    const wrapper = mount(CodeEditor, {
      props: {
        modelValue: ''
      },
      global: {
        plugins: [createPinia()]
      }
    })

    expect(wrapper.props('language')).toBe('typescript')
  })

  it('should accept readonly prop', () => {
    const wrapper = mount(CodeEditor, {
      props: {
        modelValue: '',
        readonly: true
      },
      global: {
        plugins: [createPinia()]
      }
    })

    expect(wrapper.props('readonly')).toBe(true)
  })

  it('should accept theme prop', () => {
    const wrapper = mount(CodeEditor, {
      props: {
        modelValue: '',
        theme: 'vs'
      },
      global: {
        plugins: [createPinia()]
      }
    })

    expect(wrapper.props('theme')).toBe('vs')
  })

  it('should default theme to vs-dark', () => {
    const wrapper = mount(CodeEditor, {
      props: {
        modelValue: ''
      },
      global: {
        plugins: [createPinia()]
      }
    })

    expect(wrapper.props('theme')).toBe('vs-dark')
  })
})
