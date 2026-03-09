// Mock Monaco Editor for testing
export const editor = {
	create: () => ({
		getValue: () => "mock value",
		setValue: () => {},
		updateOptions: () => {},
		getModel: () => null,
		onDidChangeModelContent: () => ({ dispose: () => {} }),
		dispose: () => {},
	}),
	setModelLanguage: () => {},
	setTheme: () => {},
};
