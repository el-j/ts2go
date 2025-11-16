package services_test

import (
	"fmt"
	"testing"

	"github.com/el-j/ts2go/core/ports"
	"github.com/el-j/ts2go/core/ports/mocks"
	"github.com/el-j/ts2go/core/services"
)

func TestNewGoRuntimeService(t *testing.T) {
	mockCompiler := mocks.NewMockGoCompiler()
	mockFS := mocks.NewMockFileSystem()

	service := services.NewGoRuntimeService(mockCompiler, mockFS)

	if service == nil {
		t.Error("Expected service to be created, got nil")
	}
}

func TestDetectGoInstallationSuccess(t *testing.T) {
	mockCompiler := mocks.NewMockGoCompiler()
	mockFS := mocks.NewMockFileSystem()

	service := services.NewGoRuntimeService(mockCompiler, mockFS)

	version, err := service.DetectGoInstallation()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if version == nil {
		t.Fatal("Expected version, got nil")
	}

	if version.Version != "go1.21.5" {
		t.Errorf("Expected version 'go1.21.5', got %s", version.Version)
	}

	if !version.Found {
		t.Error("Expected Found to be true")
	}
}

func TestDetectGoInstallationFailure(t *testing.T) {
	mockCompiler := mocks.NewMockGoCompiler()
	mockFS := mocks.NewMockFileSystem()

	// Configure mock to fail
	mockCompiler.DetectError = fmt.Errorf("go not found")

	service := services.NewGoRuntimeService(mockCompiler, mockFS)

	_, err := service.DetectGoInstallation()
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestBuildProjectEmptyPath(t *testing.T) {
	mockCompiler := mocks.NewMockGoCompiler()
	mockFS := mocks.NewMockFileSystem()

	service := services.NewGoRuntimeService(mockCompiler, mockFS)

	_, err := service.BuildProject("", "/output")
	if err == nil {
		t.Error("Expected error for empty project path, got nil")
	}
}

func TestBuildProjectNotFound(t *testing.T) {
	mockCompiler := mocks.NewMockGoCompiler()
	mockFS := mocks.NewMockFileSystem()

	service := services.NewGoRuntimeService(mockCompiler, mockFS)

	_, err := service.BuildProject("/nonexistent", "/output")
	if err == nil {
		t.Error("Expected error for non-existent project, got nil")
	}
}

func TestBuildProjectSuccess(t *testing.T) {
	mockCompiler := mocks.NewMockGoCompiler()
	mockFS := mocks.NewMockFileSystem()

	// Setup mock filesystem
	mockFS.AddDirectory("/project")

	// Configure successful build
	mockCompiler.SetBuildSuccess("Build successful")

	service := services.NewGoRuntimeService(mockCompiler, mockFS)

	result, err := service.BuildProject("/project", "/output/binary")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.Success {
		t.Error("Expected build success")
	}

	if result.OutputPath != "/output/binary" {
		t.Errorf("Expected output path '/output/binary', got %s", result.OutputPath)
	}

	if !mockCompiler.BuildCalled {
		t.Error("Expected Build to be called on compiler")
	}
}

func TestBuildProjectFailure(t *testing.T) {
	mockCompiler := mocks.NewMockGoCompiler()
	mockFS := mocks.NewMockFileSystem()

	// Setup mock filesystem
	mockFS.AddDirectory("/project")

	// Configure build failure
	mockCompiler.SetBuildFailure(fmt.Errorf("compilation error"))

	service := services.NewGoRuntimeService(mockCompiler, mockFS)

	result, err := service.BuildProject("/project", "/output/binary")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result == nil {
		t.Fatal("Expected result even with error, got nil")
	}

	if result.Success {
		t.Error("Expected build failure")
	}
}

func TestRunProjectEmptyPath(t *testing.T) {
	mockCompiler := mocks.NewMockGoCompiler()
	mockFS := mocks.NewMockFileSystem()

	service := services.NewGoRuntimeService(mockCompiler, mockFS)

	_, err := service.RunProject("", []string{})
	if err == nil {
		t.Error("Expected error for empty project path, got nil")
	}
}

func TestRunProjectSuccess(t *testing.T) {
	mockCompiler := mocks.NewMockGoCompiler()
	mockFS := mocks.NewMockFileSystem()

	// Setup mock filesystem
	mockFS.AddDirectory("/project")

	// Configure successful run
	mockCompiler.RunResult.Stdout = "Hello, World!"
	mockCompiler.RunResult.ExitCode = 0

	service := services.NewGoRuntimeService(mockCompiler, mockFS)

	result, err := service.RunProject("/project", []string{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.Success {
		t.Error("Expected run success")
	}

	if result.Stdout != "Hello, World!" {
		t.Errorf("Expected stdout 'Hello, World!', got %s", result.Stdout)
	}

	if result.ExitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", result.ExitCode)
	}
}

func TestRunProjectFailure(t *testing.T) {
	mockCompiler := mocks.NewMockGoCompiler()
	mockFS := mocks.NewMockFileSystem()

	// Setup mock filesystem
	mockFS.AddDirectory("/project")

	// Configure run failure
	mockCompiler.RunResult.Stderr = "runtime error"
	mockCompiler.RunResult.ExitCode = 1
	mockCompiler.RunError = fmt.Errorf("execution failed")

	service := services.NewGoRuntimeService(mockCompiler, mockFS)

	result, err := service.RunProject("/project", []string{})
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result == nil {
		t.Fatal("Expected result even with error, got nil")
	}

	if result.Success {
		t.Error("Expected run failure")
	}

	if result.ExitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", result.ExitCode)
	}
}

func TestTestProjectSuccess(t *testing.T) {
	mockCompiler := mocks.NewMockGoCompiler()
	mockFS := mocks.NewMockFileSystem()

	// Setup mock filesystem
	mockFS.AddDirectory("/project")

	// Configure successful test
	mockCompiler.SetTestSuccess("PASS\nok  \tpackage\t0.123s")

	service := services.NewGoRuntimeService(mockCompiler, mockFS)

	result, err := service.TestProject("/project", &ports.TestOptions{})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if !result.Success {
		t.Error("Expected test success")
	}
}

func TestTestProjectFailure(t *testing.T) {
	mockCompiler := mocks.NewMockGoCompiler()
	mockFS := mocks.NewMockFileSystem()

	// Setup mock filesystem
	mockFS.AddDirectory("/project")

	// Configure test failure
	mockCompiler.SetTestFailure(fmt.Errorf("tests failed"))

	service := services.NewGoRuntimeService(mockCompiler, mockFS)

	_, err := service.TestProject("/project", &ports.TestOptions{})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
