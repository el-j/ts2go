import { ref } from "vue";
import { invoke } from "@tauri-apps/api/core";
import { useWorkspaceStore } from "@/stores/workspace";
import { useSettingsStore } from "@/stores/settings";
import { useToast } from "primevue/usetoast";
import type { OpenFile } from "@/stores/workspace";

export interface SaveOptions {
	showDialog?: boolean;
	createBackup?: boolean;
	validateContent?: boolean;
}

export function useSaveFile() {
	const workspace = useWorkspaceStore();
	const settings = useSettingsStore();
	const toast = useToast();
	const isSaving = ref(false);

	/**
	 * Save a single file
	 */
	async function saveFile(
		file: OpenFile,
		options: SaveOptions = {},
	): Promise<boolean> {
		if (isSaving.value) return false;

		try {
			isSaving.value = true;

			// Create backup if requested or enabled in settings
			if (options.createBackup || settings.settings.enableBackups) {
				await createBackup(file);
			}

			// Validate content if requested
			if (options.validateContent) {
				const isValid = await validateSyntax(file);
				if (!isValid && options.showDialog) {
					const proceed = confirm("File has syntax errors. Save anyway?");
					if (!proceed) return false;
				}
			}

			// Save file to disk
			await invoke("write_file", {
				path: file.path,
				content: file.content,
			});

			// Mark file as saved in workspace
			workspace.markFileSaved(file.path);

			if (options.showDialog) {
				toast.add({
					severity: "success",
					summary: "File Saved",
					detail: `${file.name} saved successfully`,
					life: 3000,
				});
			}

			return true;
		} catch (error) {
			console.error("Failed to save file:", error);

			toast.add({
				severity: "error",
				summary: "Save Failed",
				detail: `Failed to save ${file.name}: ${error}`,
				life: 5000,
			});

			return false;
		} finally {
			isSaving.value = false;
		}
	}

	/**
	 * Save all dirty files
	 */
	async function saveAll(): Promise<{ saved: number; failed: number }> {
		const dirtyFiles = workspace.openFiles.filter((f) => f.isDirty);

		if (dirtyFiles.length === 0) {
			toast.add({
				severity: "info",
				summary: "Nothing to Save",
				detail: "No unsaved changes",
				life: 3000,
			});
			return { saved: 0, failed: 0 };
		}

		let saved = 0;
		let failed = 0;

		for (const file of dirtyFiles) {
			const success = await saveFile(file, { showDialog: false });
			if (success) {
				saved++;
			} else {
				failed++;
			}
		}

		toast.add({
			severity: saved > 0 && failed === 0 ? "success" : "warn",
			summary: "Save All Complete",
			detail: `Saved ${saved} file(s)${failed > 0 ? `, ${failed} failed` : ""}`,
			life: 4000,
		});

		return { saved, failed };
	}

	/**
	 * Create a backup of the file before saving
	 */
	async function createBackup(file: OpenFile): Promise<void> {
		// Use backup location from settings
		const backupDir = settings.settings.backupLocation || "./.backups";
		const fileName = file.name;
		const timestamp = new Date().toISOString().replace(/[:.]/g, "-");
		const backupPath = `${backupDir}/${fileName}.${timestamp}.backup`;

		try {
			await invoke("write_file", {
				path: backupPath,
				content: file.content,
			});
		} catch (error) {
			console.error("Failed to create backup:", error);
			// Don't throw - backup failure shouldn't prevent save
		}
	}

	/**
	 * Validate file syntax (basic check)
	 */
	async function validateSyntax(file: OpenFile): Promise<boolean> {
		const ext = file.name.split(".").pop()?.toLowerCase();

		// Only validate JSON files for now
		if (ext === "json") {
			try {
				JSON.parse(file.content);
				return true;
			} catch {
				return false;
			}
		}

		// Validate TypeScript syntax by seeing if it parses correctly
		if (ext === "ts" || ext === "tsx") {
			try {
				await invoke("transpile_code", {
					code: file.content,
					filename: file.name,
				});
				return true;
			} catch {
				return false;
			}
		}

		// Validate Go syntax by doing a test build
		if (ext === "go") {
			try {
				const res = await invoke<{ success: boolean }>("build_go_file", {
					code: file.content,
					outputPath: `${file.path}.bin`, // Output path is ignored by Go for non-main packages, but required by command
				});
				return res.success;
			} catch {
				return false;
			}
		}

		// Default to true for unsupported or generic files
		return true;
	}

	/**
	 * Check if there are unsaved changes
	 */
	function hasUnsavedChanges(): boolean {
		return workspace.hasUnsavedChanges;
	}

	/**
	 * Get list of dirty files
	 */
	function getDirtyFiles(): OpenFile[] {
		return workspace.openFiles.filter((f) => f.isDirty);
	}

	return {
		isSaving,
		saveFile,
		saveAll,
		hasUnsavedChanges,
		getDirtyFiles,
	};
}
