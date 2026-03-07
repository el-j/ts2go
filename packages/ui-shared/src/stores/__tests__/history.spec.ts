import { describe, it, expect, beforeEach, vi } from "vitest";
import { setActivePinia, createPinia } from "pinia";
import { useHistoryStore } from "../history";

// Mock @tauri-apps/api/core so useBackendState falls through to localStorage path
vi.mock("@tauri-apps/api/core", () => ({ invoke: vi.fn() }));
// Ensure we appear as web (not Tauri) so localStorage branch is used
Object.defineProperty(window, "__TAURI__", {
	value: undefined,
	configurable: true,
});

// Mock useBackendState so save() is synchronous and not gated by isLoading
vi.mock("@/composables/useBackendState", () => ({
	useBackendState: <T>(key: string, defaultValue: T) => ({
		load: async () => {
			const raw = localStorage.getItem(`ts2go:${key}`);
			if (!raw || raw === "{}") return defaultValue;
			return JSON.parse(raw) as T;
		},
		save: async (data: T) => {
			localStorage.setItem(`ts2go:${key}`, JSON.stringify(data));
		},
	}),
	useBackendSettings: <T>(defaultValue: T) => ({
		load: async () => defaultValue,
		save: async () => {},
	}),
}));

describe("History Store", () => {
	beforeEach(() => {
		setActivePinia(createPinia());
		// Clear localStorage before each test
		localStorage.clear();
	});

	it("should initialize with empty builds", () => {
		const store = useHistoryStore();

		expect(store.builds).toEqual([]);
		expect(store.successfulBuilds).toEqual([]);
		expect(store.failedBuilds).toEqual([]);
		expect(store.successRate).toBe(0);
		expect(store.averageDuration).toBe(0);
	});

	it("should add a build", () => {
		const store = useHistoryStore();

		store.addBuild({
			projectPath: "/test/project",
			filesProcessed: 5,
			totalFiles: 5,
			duration: 1000,
			status: "success",
			errors: 0,
			warnings: 0,
		});

		expect(store.builds).toHaveLength(1);
		expect(store.builds[0].projectPath).toBe("/test/project");
		expect(store.builds[0].status).toBe("success");
		expect(store.builds[0]).toHaveProperty("id");
		expect(store.builds[0]).toHaveProperty("timestamp");
	});

	it("should calculate success rate correctly", () => {
		const store = useHistoryStore();

		store.addBuild({
			projectPath: "",
			filesProcessed: 1,
			totalFiles: 1,
			duration: 100,
			status: "success",
			errors: 0,
			warnings: 0,
		});

		store.addBuild({
			projectPath: "",
			filesProcessed: 0,
			totalFiles: 1,
			duration: 100,
			status: "failed",
			errors: 1,
			warnings: 0,
		});

		store.addBuild({
			projectPath: "",
			filesProcessed: 1,
			totalFiles: 1,
			duration: 100,
			status: "success",
			errors: 0,
			warnings: 0,
		});

		expect(store.builds).toHaveLength(3);
		expect(store.successfulBuilds).toHaveLength(2);
		expect(store.failedBuilds).toHaveLength(1);
		expect(store.successRate).toBeCloseTo(66.67, 1);
	});

	it("should calculate average duration correctly", () => {
		const store = useHistoryStore();

		store.addBuild({
			projectPath: "",
			filesProcessed: 1,
			totalFiles: 1,
			duration: 1000,
			status: "success",
			errors: 0,
			warnings: 0,
		});

		store.addBuild({
			projectPath: "",
			filesProcessed: 1,
			totalFiles: 1,
			duration: 2000,
			status: "success",
			errors: 0,
			warnings: 0,
		});

		store.addBuild({
			projectPath: "",
			filesProcessed: 1,
			totalFiles: 1,
			duration: 3000,
			status: "success",
			errors: 0,
			warnings: 0,
		});

		expect(store.averageDuration).toBe(2000);
	});

	it("should not include failed builds in average duration", () => {
		const store = useHistoryStore();

		store.addBuild({
			projectPath: "",
			filesProcessed: 1,
			totalFiles: 1,
			duration: 1000,
			status: "success",
			errors: 0,
			warnings: 0,
		});

		store.addBuild({
			projectPath: "",
			filesProcessed: 0,
			totalFiles: 1,
			duration: 5000,
			status: "failed",
			errors: 1,
			warnings: 0,
		});

		expect(store.averageDuration).toBe(1000);
	});

	it("should remove a build", () => {
		const store = useHistoryStore();

		store.addBuild({
			projectPath: "",
			filesProcessed: 1,
			totalFiles: 1,
			duration: 1000,
			status: "success",
			errors: 0,
			warnings: 0,
		});

		const buildId = store.builds[0].id;

		store.removeBuild(buildId);

		expect(store.builds).toHaveLength(0);
	});

	it("should clear all history", () => {
		const store = useHistoryStore();

		store.addBuild({
			projectPath: "",
			filesProcessed: 1,
			totalFiles: 1,
			duration: 1000,
			status: "success",
			errors: 0,
			warnings: 0,
		});

		store.addBuild({
			projectPath: "",
			filesProcessed: 1,
			totalFiles: 1,
			duration: 2000,
			status: "success",
			errors: 0,
			warnings: 0,
		});

		expect(store.builds).toHaveLength(2);

		store.clearHistory();

		expect(store.builds).toHaveLength(0);
	});

	it("should limit history to maxHistorySize", () => {
		const store = useHistoryStore();
		store.maxHistorySize = 3;

		for (let i = 0; i < 5; i++) {
			store.addBuild({
				projectPath: "",
				filesProcessed: 1,
				totalFiles: 1,
				duration: 1000,
				status: "success",
				errors: 0,
				warnings: 0,
			});
		}

		expect(store.builds).toHaveLength(3);
	});

	it("should persist to localStorage", async () => {
		const store = useHistoryStore();
		// Await the initial loadFromBackend so isLoading is cleared
		await Promise.resolve();
		await Promise.resolve();

		store.addBuild({
			projectPath: "/test",
			filesProcessed: 1,
			totalFiles: 1,
			duration: 1000,
			status: "success",
			errors: 0,
			warnings: 0,
		});
		// Await saveToBackend (async)
		await Promise.resolve();
		await Promise.resolve();

		// History store persists with prefix 'ts2go:' + key 'build-history'
		const stored = localStorage.getItem("ts2go:build-history");
		expect(stored).not.toBeNull();

		const parsed = JSON.parse(stored as string);
		expect(parsed).toHaveLength(1);
		expect(parsed[0].projectPath).toBe("/test");
	});

	it("should load from localStorage on init", async () => {
		// Set up data in localStorage before creating the store
		const mockBuild = {
			id: "test-id",
			timestamp: new Date().toISOString(),
			projectPath: "/test",
			filesProcessed: 1,
			totalFiles: 1,
			duration: 1000,
			status: "success",
			errors: 0,
			warnings: 0,
		};

		localStorage.setItem("ts2go:build-history", JSON.stringify([mockBuild]));

		// Create new store instance — loadFromBackend is called in setup
		setActivePinia(createPinia());
		const store = useHistoryStore();
		// Await async loadFromBackend
		await Promise.resolve();
		await Promise.resolve();

		expect(store.builds).toHaveLength(1);
		expect(store.builds[0].projectPath).toBe("/test");
		expect(store.builds[0].timestamp).toBeInstanceOf(Date);
	});
});
