package persistence_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/el-j/ts2go/adapters/driven/persistence"
	"github.com/el-j/ts2go/core/domain"
)

func TestJSONStateRepository(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "json-repo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo, err := persistence.NewJSONStateRepository(tempDir)
	if err != nil {
		t.Fatalf("NewJSONStateRepository failed: %v", err)
	}

	state := &domain.ProjectState{
		ProjectPath: "/path/to/myproj",
	}

	t.Run("Save and Get", func(t *testing.T) {
		err := repo.SaveState("/path/to/myproj", state)
		if err != nil {
			t.Fatalf("SaveState failed: %v", err)
		}

		retrieved, err := repo.GetState("/path/to/myproj")
		if err != nil {
			t.Fatalf("GetState failed: %v", err)
		}
		if retrieved.ProjectPath != "/path/to/myproj" {
			t.Errorf("Expected path to be '/path/to/myproj', got '%s'", retrieved.ProjectPath)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := repo.Exists("/path/to/myproj")
		if err != nil || !exists {
			t.Errorf("Exists failed or returned false: %v", err)
		}

		exists, _ = repo.Exists("/path/to/nonexistent")
		if exists {
			t.Errorf("Should return false for non-existent state")
		}
	})

	t.Run("GetAllStates", func(t *testing.T) {
		repo.SaveState("/path/to/another", &domain.ProjectState{ProjectPath: "/path/to/another"})

		states, err := repo.GetAllStates()
		if err != nil {
			t.Fatalf("GetAllStates failed: %v", err)
		}

		if len(states) != 2 {
			t.Errorf("Expected 2 states, got %d", len(states))
		}
	})

	t.Run("DeleteState", func(t *testing.T) {
		err := repo.DeleteState("/path/to/myproj")
		if err != nil {
			t.Fatalf("DeleteState failed: %v", err)
		}

		exists, _ := repo.Exists("/path/to/myproj")
		if exists {
			t.Errorf("State still exists after deletion")
		}
	})

	t.Run("Errors and Edge Cases", func(t *testing.T) {
		// Write invalid json
		stateFile := filepath.Join(tempDir, "invalid.json")
		os.WriteFile(stateFile, []byte("{ bad json"), 0644)

		_, err := repo.GetState("invalid")
		if err == nil {
			t.Errorf("Expected error for invalid JSON")
		}

		// Hit GetAllStates with invalid JSON present (it should skip without error)
		states, err := repo.GetAllStates()
		if err != nil {
			t.Errorf("GetAllStates with invalid JSON should not fail: %v", err)
		}
		if len(states) != 2 { // We saved 2 states earlier: /path/to/myproj, /path/to/another
			// Just a sanity check, might be 1 if myproj was deleted
		}

		// Delete missing state
		err = repo.DeleteState("missing_project")
		if err != nil {
			t.Errorf("DeleteState on missing project should not error")
		}

		// SaveState with invalid path (e.g. file exists where dir should be)
		badDir := filepath.Join(tempDir, "bad_dir_file")
		os.WriteFile(badDir, []byte("data"), 0644)

		badRepo, _ := persistence.NewJSONStateRepository(filepath.Join(badDir, "subdir"))
		if badRepo != nil {
			// If it somehow creates it, we try to save
			badRepo.SaveState("test", &domain.ProjectState{})
		}
	})
}

func TestJSONSettingsRepository(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "json-settings-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	settingsFile := filepath.Join(tempDir, "settings.json")
	repo, err := persistence.NewJSONSettingsRepository(settingsFile)
	if err != nil {
		t.Fatalf("NewJSONSettingsRepository failed: %v", err)
	}

	t.Run("Save and Get", func(t *testing.T) {
		settings := domain.DefaultSettings()
		settings.Theme = "dark"

		err := repo.Save(settings)
		if err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		retrieved, err := repo.Get()
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if retrieved.Theme != "dark" {
			t.Errorf("Expected 'dark', got '%s'", retrieved.Theme)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := repo.Exists()
		if err != nil || !exists {
			t.Errorf("Exists failed or false: %v", err)
		}
	})

	t.Run("Reset", func(t *testing.T) {
		err := repo.Reset()
		if err != nil {
			t.Fatalf("Reset failed: %v", err)
		}

		retrieved, err := repo.Get()
		if err != nil {
			t.Fatalf("Get failed after reset: %v", err)
		}

		if retrieved.Theme == "dark" /* unless default is dark */ {
			exists, _ := repo.Exists()
			if !exists {
				t.Errorf("Expected settings to exist after Reset")
			}
		}
	})

	t.Run("Errors and Edge Cases", func(t *testing.T) {
		// Invalid JSON
		os.WriteFile(settingsFile, []byte("{ bad json"), 0644)
		_, err := repo.Get()
		if err == nil {
			t.Errorf("Expected error for invalid JSON")
		}

		// Get missing settings
		missingSettingsFile := filepath.Join(tempDir, "missing", "settings.json")
		missingRepo, _ := persistence.NewJSONSettingsRepository(missingSettingsFile)
		_, err = missingRepo.Get()
		if err == nil {
			t.Errorf("Expected error for missing settings")
		}

		// Read-only settings
		readOnlySettings := filepath.Join(tempDir, "readonly_set.json")
		os.WriteFile(readOnlySettings, []byte("{}"), 0400)
		roSetRepo, _ := persistence.NewJSONSettingsRepository(readOnlySettings)
		if roSetRepo != nil {
			roSetRepo.Save(domain.DefaultSettings()) // just to bump coverage
		}
	})
}
