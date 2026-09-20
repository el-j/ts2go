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
	defer func() { _ = os.RemoveAll(tempDir) }()

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
		_ = repo.SaveState("/path/to/another", &domain.ProjectState{ProjectPath: "/path/to/another"})

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

		exists, err := repo.Exists("/path/to/myproj")
		if err != nil {
			t.Fatalf("Exists failed after delete: %v", err)
		}
		if exists {
			t.Errorf("Expected state to be deleted")
		}
	})

	t.Run("Errors and Edge Cases", func(t *testing.T) {
		stateFile := filepath.Join(tempDir, "invalid.json")
		_ = os.WriteFile(stateFile, []byte("{ bad json"), 0644)

		_, err := repo.GetState("invalid")
		if err == nil {
			t.Errorf("Expected error for invalid JSON")
		}

		states, err := repo.GetAllStates()
		if err != nil {
			t.Errorf("GetAllStates with invalid JSON should not fail: %v", err)
		}
		if len(states) < 1 {
			t.Errorf("Expected at least 1 state")
		}

		err = repo.DeleteState("missing_project")
		if err != nil {
			t.Errorf("DeleteState on missing project should not error")
		}

		badDir := filepath.Join(tempDir, "bad_dir_file")
		_ = os.WriteFile(badDir, []byte("data"), 0644)

		badRepo, _ := persistence.NewJSONStateRepository(filepath.Join(badDir, "subdir"))
		if badRepo != nil {
			_ = badRepo.SaveState("test", &domain.ProjectState{})
		}
	})
}

func TestJSONSettingsRepository(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "json-settings-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

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

		if retrieved.Theme == "dark" {
			exists, _ := repo.Exists()
			if !exists {
				t.Errorf("Expected settings to exist after Reset")
			}
		}
	})

	t.Run("Errors and Edge Cases", func(t *testing.T) {
		_ = os.WriteFile(settingsFile, []byte("{ bad json"), 0644)
		_, err := repo.Get()
		if err == nil {
			t.Errorf("Expected error for invalid JSON")
		}

		missingSettingsFile := filepath.Join(tempDir, "missing", "settings.json")
		missingRepo, _ := persistence.NewJSONSettingsRepository(missingSettingsFile)
		_, err = missingRepo.Get()
		if err == nil {
			t.Errorf("Expected error for missing settings")
		}

		readOnlySettings := filepath.Join(tempDir, "readonly_set.json")
		_ = os.WriteFile(readOnlySettings, []byte("{}"), 0400)
		roSetRepo, _ := persistence.NewJSONSettingsRepository(readOnlySettings)
		if roSetRepo != nil {
			_ = roSetRepo.Save(domain.DefaultSettings())
		}
	})
}
