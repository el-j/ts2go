package mocks

import (
	"fmt"

	"github.com/el-j/ts2go/core/domain"
)

// MockStateRepository is a mock implementation of ports.StateRepository for testing
type MockStateRepository struct {
	// States maps project paths to their state
	States map[string]*domain.ProjectState

	// GetError is the error to return from GetState
	GetError error

	// SaveError is the error to return from SaveState
	SaveError error

	// DeleteError is the error to return from DeleteState
	DeleteError error

	// GetCalled tracks if GetState was called
	GetCalled bool

	// SaveCalled tracks if SaveState was called
	SaveCalled bool

	// LastSavedState tracks the most recent state saved
	LastSavedState *domain.ProjectState
}

// NewMockStateRepository creates a new mock state repository
func NewMockStateRepository() *MockStateRepository {
	return &MockStateRepository{
		States: make(map[string]*domain.ProjectState),
	}
}

// GetState retrieves state for a project
func (m *MockStateRepository) GetState(projectPath string) (*domain.ProjectState, error) {
	m.GetCalled = true

	if m.GetError != nil {
		return nil, m.GetError
	}

	state, exists := m.States[projectPath]
	if !exists {
		return nil, fmt.Errorf("state not found for project: %s", projectPath)
	}

	return state, nil
}

// SaveState saves state for a project
func (m *MockStateRepository) SaveState(projectPath string, state *domain.ProjectState) error {
	m.SaveCalled = true
	m.LastSavedState = state

	if m.SaveError != nil {
		return m.SaveError
	}

	m.States[projectPath] = state
	return nil
}

// GetAllStates retrieves all project states
func (m *MockStateRepository) GetAllStates() (map[string]*domain.ProjectState, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	return m.States, nil
}

// DeleteState deletes state for a project
func (m *MockStateRepository) DeleteState(projectPath string) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}

	delete(m.States, projectPath)
	return nil
}

// Exists checks if state exists for a project
func (m *MockStateRepository) Exists(projectPath string) (bool, error) {
	_, exists := m.States[projectPath]
	return exists, nil
}

// AddState adds a state to the mock repository (helper for tests)
func (m *MockStateRepository) AddState(projectPath string, state *domain.ProjectState) {
	m.States[projectPath] = state
}

// MockSettingsRepository is a mock implementation of ports.SettingsRepository for testing
type MockSettingsRepository struct {
	// Settings is the settings to return
	Settings *domain.Settings

	// GetError is the error to return from Get
	GetError error

	// SaveError is the error to return from Save
	SaveError error

	// ResetError is the error to return from Reset
	ResetError error

	// GetCalled tracks if Get was called
	GetCalled bool

	// SaveCalled tracks if Save was called
	SaveCalled bool

	// LastSavedSettings tracks the most recent settings saved
	LastSavedSettings *domain.Settings
}

// NewMockSettingsRepository creates a new mock settings repository
func NewMockSettingsRepository() *MockSettingsRepository {
	return &MockSettingsRepository{
		Settings: &domain.Settings{},
	}
}

// Get retrieves application settings
func (m *MockSettingsRepository) Get() (*domain.Settings, error) {
	m.GetCalled = true

	if m.GetError != nil {
		return nil, m.GetError
	}

	return m.Settings, nil
}

// Save saves application settings
func (m *MockSettingsRepository) Save(settings *domain.Settings) error {
	m.SaveCalled = true
	m.LastSavedSettings = settings

	if m.SaveError != nil {
		return m.SaveError
	}

	m.Settings = settings
	return nil
}

// Exists checks if settings exist
func (m *MockSettingsRepository) Exists() (bool, error) {
	return m.Settings != nil, nil
}

// Reset resets settings to defaults
func (m *MockSettingsRepository) Reset() error {
	if m.ResetError != nil {
		return m.ResetError
	}

	m.Settings = &domain.Settings{}
	return nil
}

// SetSettings sets the settings to return (helper for tests)
func (m *MockSettingsRepository) SetSettings(settings *domain.Settings) {
	m.Settings = settings
}
