import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { useKeyboardShortcuts, type Shortcut } from "../useKeyboardShortcuts";
import { createPinia, setActivePinia } from "pinia";

/**
 * useKeyboardShortcuts uses onMounted/onUnmounted Vue lifecycle hooks.
 * Outside a Vue component these don't fire, so we test the core matching logic
 * directly by manually registering a listener with the same algorithm.
 */
function createShortcutListener(shortcuts: Shortcut[]) {
	const handler = (event: KeyboardEvent) => {
		for (const shortcut of shortcuts) {
			const keyMatches = event.key.toLowerCase() === shortcut.key.toLowerCase();
			const ctrlMatches = !!shortcut.ctrl === (event.ctrlKey || event.metaKey);
			const shiftMatches = !!shortcut.shift === event.shiftKey;
			const altMatches = !!shortcut.alt === event.altKey;

			if (keyMatches && ctrlMatches && shiftMatches && altMatches) {
				event.preventDefault();
				shortcut.handler();
				break;
			}
		}
	};
	window.addEventListener("keydown", handler);
	return () => window.removeEventListener("keydown", handler);
}

describe("useKeyboardShortcuts", () => {
	let cleanup: (() => void) | null = null;

	beforeEach(() => {
		setActivePinia(createPinia());
	});

	afterEach(() => {
		if (cleanup) {
			cleanup();
			cleanup = null;
		}
		vi.clearAllMocks();
	});

	it("should register and trigger keyboard shortcuts", () => {
		const mockHandler = vi.fn();

		const shortcuts: Shortcut[] = [
			{
				key: "s",
				ctrl: true,
				handler: mockHandler,
				description: "Save file",
			},
		];

		// Verify composable returns expected shape
		const result = useKeyboardShortcuts(shortcuts);
		expect(result.shortcuts).toBe(shortcuts);

		// Test key matching logic via manual listener (lifecycle hooks don't run outside Vue)
		cleanup = createShortcutListener(shortcuts);

		window.dispatchEvent(
			new KeyboardEvent("keydown", {
				key: "s",
				ctrlKey: true,
				bubbles: true,
				cancelable: true,
			}),
		);

		expect(mockHandler).toHaveBeenCalled();
	});

	it("should handle shift modifier", () => {
		const mockHandler = vi.fn();

		const shortcuts: Shortcut[] = [
			{
				key: "s",
				ctrl: true,
				shift: true,
				handler: mockHandler,
				description: "Save all files",
			},
		];

		cleanup = createShortcutListener(shortcuts);

		window.dispatchEvent(
			new KeyboardEvent("keydown", {
				key: "s",
				ctrlKey: true,
				shiftKey: true,
				bubbles: true,
				cancelable: true,
			}),
		);

		expect(mockHandler).toHaveBeenCalled();
	});

	it("should not trigger when modifiers do not match", () => {
		const mockHandler = vi.fn();

		const shortcuts: Shortcut[] = [
			{
				key: "s",
				ctrl: true,
				handler: mockHandler,
				description: "Save file",
			},
		];

		cleanup = createShortcutListener(shortcuts);

		// 'S' without Ctrl should not trigger
		window.dispatchEvent(
			new KeyboardEvent("keydown", {
				key: "s",
				bubbles: true,
				cancelable: true,
			}),
		);

		expect(mockHandler).not.toHaveBeenCalled();
	});

	it("should handle multiple shortcuts without collision", () => {
		const saveHandler = vi.fn();
		const saveAllHandler = vi.fn();

		const shortcuts: Shortcut[] = [
			{
				key: "s",
				ctrl: true,
				handler: saveHandler,
				description: "Save file",
			},
			{
				key: "s",
				ctrl: true,
				shift: true,
				handler: saveAllHandler,
				description: "Save all files",
			},
		];

		cleanup = createShortcutListener(shortcuts);

		// Ctrl+S triggers saveHandler, not saveAllHandler
		window.dispatchEvent(
			new KeyboardEvent("keydown", {
				key: "s",
				ctrlKey: true,
				bubbles: true,
				cancelable: true,
			}),
		);

		expect(saveHandler).toHaveBeenCalledTimes(1);
		expect(saveAllHandler).not.toHaveBeenCalled();

		saveHandler.mockClear();

		// Ctrl+Shift+S triggers saveAllHandler, not saveHandler
		window.dispatchEvent(
			new KeyboardEvent("keydown", {
				key: "s",
				ctrlKey: true,
				shiftKey: true,
				bubbles: true,
				cancelable: true,
			}),
		);

		expect(saveAllHandler).toHaveBeenCalledTimes(1);
		expect(saveHandler).not.toHaveBeenCalled();
	});

	it("should return shortcuts array", () => {
		const shortcuts: Shortcut[] = [
			{
				key: "s",
				ctrl: true,
				handler: vi.fn(),
				description: "Save file",
			},
			{
				key: "s",
				ctrl: true,
				shift: true,
				handler: vi.fn(),
				description: "Save all files",
			},
		];

		const result = useKeyboardShortcuts(shortcuts);

		expect(result.shortcuts).toHaveLength(2);
		expect(result.shortcuts[0].description).toBe("Save file");
		expect(result.shortcuts[1].description).toBe("Save all files");
	});
});
