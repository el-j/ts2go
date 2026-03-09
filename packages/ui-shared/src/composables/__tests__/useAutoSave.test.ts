import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { useAutoSave } from "../useAutoSave";
import { useSettingsStore } from "@/stores/settings";
import { createPinia, setActivePinia } from "pinia";
import type { OpenFile } from "@/stores/workspace";

// Stable spy shared across mock calls — avoids stale reference issues when
// the test calls useSaveFile() separately to retrieve the function.
const saveFileSpy = vi.fn().mockResolvedValue(true);

vi.mock("../useSaveFile", () => ({
	useSaveFile: () => ({
		saveFile: saveFileSpy,
	}),
}));

describe("useAutoSave", () => {
	beforeEach(() => {
		setActivePinia(createPinia());
		vi.useFakeTimers();
		saveFileSpy.mockClear();
	});

	afterEach(() => {
		vi.useRealTimers();
		vi.clearAllMocks();
	});

	it("should schedule auto-save for a dirty file", async () => {
		const { scheduleSave } = useAutoSave();
		const settings = useSettingsStore();

		settings.settings.autoSave = true;
		settings.settings.autoSaveDelay = 3000;

		const mockFile: OpenFile = {
			path: "/test/file.ts",
			name: "file.ts",
			content: "const x = 1;",
			isDirty: true,
		};

		scheduleSave(mockFile);

		vi.advanceTimersByTime(3000);
		await vi.runAllTimersAsync();

		expect(saveFileSpy).toHaveBeenCalledWith(mockFile, { showDialog: false });
	});

	it("should respect auto-save delay setting", async () => {
		const { scheduleSave } = useAutoSave();
		const settings = useSettingsStore();

		settings.settings.autoSave = true;
		settings.settings.autoSaveDelay = 5000;

		const mockFile: OpenFile = {
			path: "/test/file.ts",
			name: "file.ts",
			content: "const x = 1;",
			isDirty: true,
		};

		scheduleSave(mockFile);

		// Advance only 4999ms — timer should NOT have fired yet
		vi.advanceTimersByTime(4999);
		expect(saveFileSpy).not.toHaveBeenCalled();

		// Advance final 1ms — timer fires now
		vi.advanceTimersByTime(1);
		await Promise.resolve(); // let the async setTimeout callback resolve
		expect(saveFileSpy).toHaveBeenCalled();
	});

	it("should not schedule save when auto-save is disabled", () => {
		const { scheduleSave } = useAutoSave();
		const settings = useSettingsStore();

		settings.settings.autoSave = false;

		const mockFile: OpenFile = {
			path: "/test/file.ts",
			name: "file.ts",
			content: "const x = 1;",
			isDirty: true,
		};

		scheduleSave(mockFile);
		vi.advanceTimersByTime(3000);

		expect(vi.getTimerCount()).toBe(0);
	});

	it("should cancel existing timer when scheduling new save", async () => {
		const { scheduleSave } = useAutoSave();
		const settings = useSettingsStore();

		settings.settings.autoSave = true;
		settings.settings.autoSaveDelay = 3000;

		const mockFile: OpenFile = {
			path: "/test/file.ts",
			name: "file.ts",
			content: "const x = 1;",
			isDirty: true,
		};

		// Schedule first save
		scheduleSave(mockFile);
		vi.advanceTimersByTime(1000);

		// Schedule again (should cancel first)
		scheduleSave(mockFile);
		vi.advanceTimersByTime(3000);
		await vi.runAllTimersAsync();

		// Should only be called once (second timer)
		expect(saveFileSpy).toHaveBeenCalledTimes(1);
	});

	it("should cancel all pending saves", () => {
		const { scheduleSave, cancelAllSaves } = useAutoSave();
		const settings = useSettingsStore();

		settings.settings.autoSave = true;
		settings.settings.autoSaveDelay = 3000;

		const file1: OpenFile = {
			path: "/test/file1.ts",
			name: "file1.ts",
			content: "const x = 1;",
			isDirty: true,
		};

		const file2: OpenFile = {
			path: "/test/file2.ts",
			name: "file2.ts",
			content: "const y = 2;",
			isDirty: true,
		};

		scheduleSave(file1);
		scheduleSave(file2);

		expect(vi.getTimerCount()).toBe(2);

		cancelAllSaves();

		expect(vi.getTimerCount()).toBe(0);
	});

	it("should cancel save for specific file", () => {
		const { scheduleSave, cancelSave } = useAutoSave();
		const settings = useSettingsStore();

		settings.settings.autoSave = true;
		settings.settings.autoSaveDelay = 3000;

		const file1: OpenFile = {
			path: "/test/file1.ts",
			name: "file1.ts",
			content: "const x = 1;",
			isDirty: true,
		};

		const file2: OpenFile = {
			path: "/test/file2.ts",
			name: "file2.ts",
			content: "const y = 2;",
			isDirty: true,
		};

		scheduleSave(file1);
		scheduleSave(file2);

		cancelSave("/test/file1.ts");

		expect(vi.getTimerCount()).toBe(1);
	});

	it("should not save file that is no longer dirty", async () => {
		const { scheduleSave } = useAutoSave();
		const settings = useSettingsStore();

		settings.settings.autoSave = true;
		settings.settings.autoSaveDelay = 3000;

		const mockFile: OpenFile = {
			path: "/test/file.ts",
			name: "file.ts",
			content: "const x = 1;",
			isDirty: true,
		};

		scheduleSave(mockFile);

		// Mark file as clean before timer fires
		mockFile.isDirty = false;

		vi.advanceTimersByTime(3000);
		await vi.runAllTimersAsync();

		expect(saveFileSpy).not.toHaveBeenCalled();
	});
});
