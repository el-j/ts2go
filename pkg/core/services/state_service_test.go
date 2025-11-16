package services_test

import (
	"fmt"
	"testing"

	"github.com/el-j/ts2go/pkg/core/domain"
	"github.com/el-j/ts2go/pkg/core/ports/mocks"
	"github.com/el-j/ts2go/pkg/core/services"
)

func TestNewStateService(t *testing.T) {
	mockStateRepo := mocks.NewMockStateRepository()
	mockSettingsRepo := mocks.NewMockSettingsRepository()

	service := services.NewStateService(mockStateRepo, mockSettingsRepo)

	if service == nil {
		t.Error("Expected service to be created, got nil")
	}
}

func TestGetProjectStateEmptyPath(t *testing.T) {
	mockStateRepo := mocks.NewMockStateRepository()
	mockSettingsRepo := mocks.NewMockSettingsRepository()

	service := services.NewStateService(mockStateRepo, mockSettingsRepo)

	_, err := service.GetProjectState("")
	if err == nil {
		t.Error("Expected error for empty project path, got nil")
	}
}

func TestGetProjectStateSuccess(t *testing.T) {
	mockStateRepo := mocks.NewMockStateRepository()
	mockSettingsRepo := mocks.NewMockSettingsRepository()

	// Add state to mock repository
	expectedState := &domain.ProjectState{
		ProjectPath: "/project",
	}
	mockStateRepo.AddState("/project", expectedState)

	service := services.NewStateService(mockStateRepo, mockSettingsRepo)

	state, err := service.GetProjectState("/project")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if state == nil {
		t.Fatal("Expected state, got nil")
	}

	if state.ProjectPath != "/project" {
		t.Errorf("Expected project path '/project', got %s", state.ProjectPath)
	}

	if !mockStateRepo.GetCalled {
		t.Error("Expected GetState to be called on repository")
	}
}

func TestGetProjectStateNotFound(t *testing.T) {
	mockStateRepo := mocks.NewMockStateRepository()
	mockSettingsRepo := mocks.NewMockSettingsRepository()

	service := services.NewStateService(mockStateRepo, mockSettingsRepo)

	_, err := service.GetProjectState("/nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent project, got nil")
	}
}

func TestSaveProjectStateNilState(t *testing.T) {
	mockStateRepo := mocks.NewMockStateRepository()
	mockSettingsRepo := mocks.NewMockSettingsRepository()

	service := services.NewStateService(mockStateRepo, mockSettingsRepo)

	err := service.SaveProjectState(nil)
	if err == nil {
		t.Error("Expected error for nil state, got nil")
	}
}

func TestSaveProjectStateEmptyPath(t *testing.T) {
	mockStateRepo := mocks.NewMockStateRepository()
	mockSettingsRepo := mocks.NewMockSettingsRepository()

	service := services.NewStateService(mockStateRepo, mockSettingsRepo)

	state := &domain.ProjectState{
		ProjectPath: "",
	}

	err := service.SaveProjectState(state)
	if err == nil {
		t.Error("Expected error for empty project path, got nil")
	}
}

func TestSaveProjectStateSuccess(t *testing.T) {
	mockStateRepo := mocks.NewMockStateRepository()
	mockSettingsRepo := mocks.NewMockSettingsRepository()

	service := services.NewStateService(mockStateRepo, mockSettingsRepo)

	state := &domain.ProjectState{
		ProjectPath: "/project",
	}

	err := service.SaveProjectState(state)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !mockStateRepo.SaveCalled {
		t.Error("Expected SaveState to be called on repository")
	}

	if mockStateRepo.LastSavedState == nil {
		t.Error("Expected last saved state to be set")
	} else if mockStateRepo.LastSavedState.ProjectPath != "/project" {
		t.Errorf("Expected saved project path '/project', got %s", mockStateRepo.LastSavedState.ProjectPath)
	}
}

func TestSaveProjectStateFailure(t *testing.T) {
	mockStateRepo := mocks.NewMockStateRepository()
	mockSettingsRepo := mocks.NewMockSettingsRepository()

	// Configure mock to fail
	mockStateRepo.SaveError = fmt.Errorf("save failed")

	service := services.NewStateService(mockStateRepo, mockSettingsRepo)

	state := &domain.ProjectState{
		ProjectPath: "/project",
	}

	err := service.SaveProjectState(state)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetAllProjectStatesSuccess(t *testing.T) {
	mockStateRepo := mocks.NewMockStateRepository()
	mockSettingsRepo := mocks.NewMockSettingsRepository()

	// Add states to mock repository
	mockStateRepo.AddState("/project1", &domain.ProjectState{ProjectPath: "/project1"})
	mockStateRepo.AddState("/project2", &domain.ProjectState{ProjectPath: "/project2"})

	service := services.NewStateService(mockStateRepo, mockSettingsRepo)

	states, err := service.GetAllProjectStates()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(states) != 2 {
		t.Errorf("Expected 2 states, got %d", len(states))
	}

	if states["/project1"] == nil {
		t.Error("Expected state for /project1")
	}
	if states["/project2"] == nil {
		t.Error("Expected state for /project2")
	}
}

func TestDeleteProjectStateEmptyPath(t *testing.T) {
	mockStateRepo := mocks.NewMockStateRepository()
	mockSettingsRepo := mocks.NewMockSettingsRepository()

	service := services.NewStateService(mockStateRepo, mockSettingsRepo)

	err := service.DeleteProjectState("")
	if err == nil {
		t.Error("Expected error for empty project path, got nil")
	}
}

func TestDeleteProjectStateSuccess(t *testing.T) {
	mockStateRepo := mocks.NewMockStateRepository()
	mockSettingsRepo := mocks.NewMockSettingsRepository()

	// Add state to mock repository
	mockStateRepo.AddState("/project", &domain.ProjectState{ProjectPath: "/project"})

	service := services.NewStateService(mockStateRepo, mockSettingsRepo)

	err := service.DeleteProjectState("/project")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify state was deleted
	exists, _ := mockStateRepo.Exists("/project")
	if exists {
		t.Error("Expected state to be deleted")
	}
}

func TestGetSettingsSuccess(t *testing.T) {
	mockStateRepo := mocks.NewMockStateRepository()
	mockSettingsRepo := mocks.NewMockSettingsRepository()

	// Setup mock settings
	expectedSettings := &domain.Settings{}
	mockSettingsRepo.SetSettings(expectedSettings)

	service := services.NewStateService(mockStateRepo, mockSettingsRepo)

	settings, err := service.GetSettings()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if settings == nil {
		t.Error("Expected settings, got nil")
	}

	if !mockSettingsRepo.GetCalled {
		t.Error("Expected Get to be called on settings repository")
	}
}

func TestSaveSettingsNilSettings(t *testing.T) {
	mockStateRepo := mocks.NewMockStateRepository()
	mockSettingsRepo := mocks.NewMockSettingsRepository()

	service := services.NewStateService(mockStateRepo, mockSettingsRepo)

	err := service.SaveSettings(nil)
	if err == nil {
		t.Error("Expected error for nil settings, got nil")
	}
}

func TestSaveSettingsSuccess(t *testing.T) {
	mockStateRepo := mocks.NewMockStateRepository()
	mockSettingsRepo := mocks.NewMockSettingsRepository()

	service := services.NewStateService(mockStateRepo, mockSettingsRepo)

	settings := &domain.Settings{}

	err := service.SaveSettings(settings)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !mockSettingsRepo.SaveCalled {
		t.Error("Expected Save to be called on settings repository")
	}

	if mockSettingsRepo.LastSavedSettings == nil {
		t.Error("Expected last saved settings to be set")
	}
}
