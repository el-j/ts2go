package ports

import "github.com/el-j/ts2go/pkg/core/domain"

// StateRepository is what the core uses to persist project state.
// This is a DRIVEN PORT (outbound) - a dependency the core needs.
// Implementations might use JSON files, SQLite, or any other storage.
type StateRepository interface {
	// GetState retrieves state for a specific project
	GetState(projectPath string) (*domain.ProjectState, error)

	// SaveState saves state for a specific project
	SaveState(projectPath string, state *domain.ProjectState) error

	// GetAllStates retrieves all project states
	GetAllStates() (map[string]*domain.ProjectState, error)

	// DeleteState deletes state for a specific project
	DeleteState(projectPath string) error

	// Exists checks if state exists for a project
	Exists(projectPath string) (bool, error)
}

// SettingsRepository is what the core uses to persist application settings.
// This is a DRIVEN PORT (outbound) - a dependency the core needs.
type SettingsRepository interface {
	// Get retrieves application settings
	Get() (*domain.Settings, error)

	// Save saves application settings
	Save(settings *domain.Settings) error

	// Exists checks if settings exist
	Exists() (bool, error)

	// Reset resets settings to defaults
	Reset() error
}
