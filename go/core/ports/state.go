package ports

import "github.com/el-j/ts2go/core/domain"

// StateService manages persistent state for projects and application settings.
// This is a DRIVING PORT (inbound) - it's the API of our core logic.
type StateService interface {
	// GetProjectState retrieves the state for a specific project
	GetProjectState(projectPath string) (*domain.ProjectState, error)

	// SaveProjectState saves the state for a specific project
	SaveProjectState(state *domain.ProjectState) error

	// GetAllProjectStates retrieves all project states
	GetAllProjectStates() (map[string]*domain.ProjectState, error)

	// DeleteProjectState deletes the state for a specific project
	DeleteProjectState(projectPath string) error

	// GetSettings retrieves application settings
	GetSettings() (*domain.Settings, error)

	// SaveSettings saves application settings
	SaveSettings(settings *domain.Settings) error
}
