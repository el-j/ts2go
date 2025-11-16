package services

import (
	"fmt"
	"time"

	"github.com/el-j/ts2go/pkg/core/domain"
	"github.com/el-j/ts2go/pkg/core/ports"
)

// StateServiceImpl implements the StateService port
type StateServiceImpl struct {
	stateRepo    ports.StateRepository
	settingsRepo ports.SettingsRepository
}

// NewStateService creates a new state service
func NewStateService(stateRepo ports.StateRepository, settingsRepo ports.SettingsRepository) ports.StateService {
	return &StateServiceImpl{
		stateRepo:    stateRepo,
		settingsRepo: settingsRepo,
	}
}

// GetProjectState retrieves the state for a specific project
func (s *StateServiceImpl) GetProjectState(projectPath string) (*domain.ProjectState, error) {
	if projectPath == "" {
		return nil, domain.NewValidationError("project path cannot be empty")
	}

	state, err := s.stateRepo.GetState(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load project state: %w", err)
	}

	return state, nil
}

// SaveProjectState saves the current state of a project
func (s *StateServiceImpl) SaveProjectState(state *domain.ProjectState) error {
	if state == nil {
		return domain.NewValidationError("state cannot be nil")
	}
	if state.ProjectPath == "" {
		return domain.NewValidationError("project path cannot be empty")
	}

	// Update timestamp
	state.UpdatedAt = time.Now()

	// Save to repository
	if err := s.stateRepo.SaveState(state.ProjectPath, state); err != nil {
		return fmt.Errorf("failed to save project state: %w", err)
	}

	return nil
}

// GetAllProjectStates retrieves all project states
func (s *StateServiceImpl) GetAllProjectStates() (map[string]*domain.ProjectState, error) {
	states, err := s.stateRepo.GetAllStates()
	if err != nil {
		return nil, fmt.Errorf("failed to load all project states: %w", err)
	}

	return states, nil
}

// DeleteProjectState deletes the state for a specific project
func (s *StateServiceImpl) DeleteProjectState(projectPath string) error {
	if projectPath == "" {
		return domain.NewValidationError("project path cannot be empty")
	}

	if err := s.stateRepo.DeleteState(projectPath); err != nil {
		return fmt.Errorf("failed to delete project state: %w", err)
	}

	return nil
}

// GetSettings retrieves the global settings
func (s *StateServiceImpl) GetSettings() (*domain.Settings, error) {
	settings, err := s.settingsRepo.Get()
	if err != nil {
		// If settings don't exist, return defaults
		if err == domain.ErrNotFound {
			return domain.NewDefaultSettings(), nil
		}
		return nil, fmt.Errorf("failed to load settings: %w", err)
	}

	return settings, nil
}

// SaveSettings saves application settings
func (s *StateServiceImpl) SaveSettings(settings *domain.Settings) error {
	if settings == nil {
		return domain.NewValidationError("settings cannot be nil")
	}

	// Validate settings
	if err := s.validateSettings(settings); err != nil {
		return err
	}

	if err := s.settingsRepo.Save(settings); err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}

	return nil
}

// validateSettings validates settings values
func (s *StateServiceImpl) validateSettings(settings *domain.Settings) error {
	if settings.AutoSave && settings.AutoSaveDelay < 100 {
		return domain.NewValidationError("auto-save delay must be at least 100ms")
	}

	// Add more validation as needed

	return nil
}
