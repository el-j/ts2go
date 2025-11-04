<template>
  <div ref="editorContainer" class="code-editor"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import * as monaco from 'monaco-editor'
import { useSettingsStore } from '@/stores/settings'

interface Props {
  modelValue: string
  language?: string
  readonly?: boolean
  theme?: 'vs' | 'vs-dark' | 'hc-black'
}

const props = withDefaults(defineProps<Props>(), {
  language: 'typescript',
  readonly: false,
  theme: 'vs-dark'
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const editorContainer = ref<HTMLElement>()
let editor: monaco.editor.IStandaloneCodeEditor | null = null
const settingsStore = useSettingsStore()

// Expose editor actions for external use
defineExpose({
  showFind: () => {
    editor?.getAction('actions.find')?.run()
  },
  showReplace: () => {
    editor?.getAction('editor.action.startFindReplaceAction')?.run()
  },
  getEditor: () => editor
})

onMounted(() => {
  if (!editorContainer.value) return

  editor = monaco.editor.create(editorContainer.value, {
    value: props.modelValue,
    language: props.language,
    theme: props.theme,
    readOnly: props.readonly,
    automaticLayout: true,
    minimap: { enabled: settingsStore.settings.minimap },
    fontSize: settingsStore.settings.fontSize,
    lineNumbers: settingsStore.settings.lineNumbers ? 'on' : 'off',
    folding: true,
    scrollBeyondLastLine: false,
    wordWrap: settingsStore.settings.wordWrap ? 'on' : 'off',
    tabSize: settingsStore.settings.tabSize,
    find: {
      addExtraSpaceOnTop: false,
      autoFindInSelection: 'never',
      seedSearchStringFromSelection: 'always'
    }
  })

  editor.onDidChangeModelContent(() => {
    if (editor) {
      emit('update:modelValue', editor.getValue())
    }
  })

  // Watch settings changes and update editor
  watch(() => settingsStore.settings, (newSettings) => {
    if (editor) {
      editor.updateOptions({
        minimap: { enabled: newSettings.minimap },
        fontSize: newSettings.fontSize,
        lineNumbers: newSettings.lineNumbers ? 'on' : 'off',
        wordWrap: newSettings.wordWrap ? 'on' : 'off',
        tabSize: newSettings.tabSize
      })
    }
  }, { deep: true })
})

watch(() => props.modelValue, (newValue) => {
  if (editor && newValue !== editor.getValue()) {
    editor.setValue(newValue)
  }
})

watch(() => props.language, (newLanguage) => {
  if (editor) {
    const model = editor.getModel()
    if (model) {
      monaco.editor.setModelLanguage(model, newLanguage)
    }
  }
})

watch(() => props.theme, (newTheme) => {
  if (editor) {
    monaco.editor.setTheme(newTheme)
  }
})

onBeforeUnmount(() => {
  editor?.dispose()
})
</script>

<style scoped>
.code-editor {
  width: 100%;
  height: 100%;
  min-height: 400px;
}
</style>
