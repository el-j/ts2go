// Mock Monaco Editor for testing
export const editor = {
  create: () => ({
    getValue: () => 'mock value',
    setValue: () => {},
    getModel: () => null,
    onDidChangeModelContent: () => ({ dispose: () => {} }),
    dispose: () => {}
  }),
  setModelLanguage: () => {},
  setTheme: () => {}
}
