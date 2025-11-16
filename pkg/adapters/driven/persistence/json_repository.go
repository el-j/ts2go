package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/el-j/ts2go/pkg/core/domain"
	"github.com/el-j/ts2go/pkg/core/ports"
)

// JSONStateRepository implements StateRepository using JSON files
type JSONStateRepository struct {
	baseDir string
	mu      sync.RWMutex
}

// NewJSONStateRepository creates a new JSON-based state repository
func NewJSONStateRepository(baseDir string) (ports.StateRepository, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %w", err)
	}

	return &JSONStateRepository{
		baseDir: baseDir,
	}, nil
}

// GetState retrieves state for a specific project
func (r *JSONStateRepository) GetState(projectPath string) (*domain.ProjectState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stateFile := r.getStateFilePath(projectPath)

	data, err := os.ReadFile(stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var state domain.ProjectState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %w", err)
	}

	return &state, nil
}

// SaveState saves state for a specific project
func (r *JSONStateRepository) SaveState(projectPath string, state *domain.ProjectState) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stateFile := r.getStateFilePath(projectPath)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(stateFile), 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// GetAllStates retrieves all project states
func (r *JSONStateRepository) GetAllStates() (map[string]*domain.ProjectState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	states := make(map[string]*domain.ProjectState)

	// Read all JSON files in the base directory
	entries, err := os.ReadDir(r.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return states, nil
		}
		return nil, fmt.Errorf("failed to read state directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(r.baseDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue // Skip files that can't be read
		}

		var state domain.ProjectState
		if err := json.Unmarshal(data, &state); err != nil {
			continue // Skip files that can't be parsed
		}

		states[state.ProjectPath] = &state
	}

	return states, nil
}

// DeleteState deletes state for a specific project
func (r *JSONStateRepository) DeleteState(projectPath string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stateFile := r.getStateFilePath(projectPath)

	if err := os.Remove(stateFile); err != nil {
		if os.IsNotExist(err) {
			return nil // Already deleted
		}
		return fmt.Errorf("failed to delete state file: %w", err)
	}

	return nil
}

// Exists checks if state exists for a project
func (r *JSONStateRepository) Exists(projectPath string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stateFile := r.getStateFilePath(projectPath)

	_, err := os.Stat(stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// getStateFilePath generates the file path for a project's state
func (r *JSONStateRepository) getStateFilePath(projectPath string) string {
	// Create a safe filename from the project path
	filename := filepath.Base(projectPath) + ".json"
	return filepath.Join(r.baseDir, filename)
}

// JSONSettingsRepository implements SettingsRepository using JSON file
type JSONSettingsRepository struct {
	filePath string
	mu       sync.RWMutex
}

// NewJSONSettingsRepository creates a new JSON-based settings repository
func NewJSONSettingsRepository(filePath string) (ports.SettingsRepository, error) {
	// Create parent directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create settings directory: %w", err)
	}

	return &JSONSettingsRepository{
		filePath: filePath,
	}, nil
}

// Get retrieves application settings
func (r *JSONSettingsRepository) Get() (*domain.Settings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to read settings file: %w", err)
	}

	var settings domain.Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse settings file: %w", err)
	}

	return &settings, nil
}

// Save saves application settings
func (r *JSONSettingsRepository) Save(settings *domain.Settings) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	return nil
}

// Exists checks if settings exist
func (r *JSONSettingsRepository) Exists() (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, err := os.Stat(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Reset resets settings to defaults
func (r *JSONSettingsRepository) Reset() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	defaults := domain.DefaultSettings()
	return r.Save(defaults)
}
