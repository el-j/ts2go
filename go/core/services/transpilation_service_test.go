package services_test

import (
	"testing"

	"github.com/el-j/ts2go/core/domain"
	"github.com/el-j/ts2go/core/ports"
	"github.com/el-j/ts2go/core/ports/mocks"
	"github.com/el-j/ts2go/core/services"
)

func TestNewTranspilationService(t *testing.T) {
	mockFS := mocks.NewMockFileSystem()
	mockCompiler := mocks.NewMockGoCompiler()

	service := services.NewTranspilationService(
		mockFS,
		mockCompiler,
		nil, // analyzer
		nil, // mapper
		nil, // codegen
	)

	if service == nil {
		t.Error("Expected service to be created, got nil")
	}
}

func TestTranspileProjectEmptyPath(t *testing.T) {
	mockFS := mocks.NewMockFileSystem()
	mockCompiler := mocks.NewMockGoCompiler()

	service := services.NewTranspilationService(
		mockFS,
		mockCompiler,
		nil,
		nil,
		nil,
	)

	options := ports.TranspileOptions{
		OutputPath: "/output",
	}

	_, err := service.TranspileProject("", options)
	if err == nil {
		t.Error("Expected error for empty project path, got nil")
	}
}

func TestTranspileProjectSuccess(t *testing.T) {
	mockFS := mocks.NewMockFileSystem()
	mockCompiler := mocks.NewMockGoCompiler()

	// Setup mock filesystem with TypeScript files
	tsFile1 := domain.NewFile("/project/src/main.ts", "const x = 1;")
	tsFile2 := domain.NewFile("/project/src/utils.ts", "export const y = 2;")
	mockFS.SetScanResult([]domain.File{*tsFile1, *tsFile2})

	service := services.NewTranspilationService(
		mockFS,
		mockCompiler,
		nil,
		nil,
		nil,
	)

	options := ports.TranspileOptions{
		OutputPath: "/output",
	}

	result, err := service.TranspileProject("/project", options)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ProjectPath != "/project" {
		t.Errorf("Expected project path '/project', got %s", result.ProjectPath)
	}

	if result.OutputPath != "/output" {
		t.Errorf("Expected output path '/output', got %s", result.OutputPath)
	}
}

func TestTranspileFileEmptyPath(t *testing.T) {
	mockFS := mocks.NewMockFileSystem()
	mockCompiler := mocks.NewMockGoCompiler()

	service := services.NewTranspilationService(
		mockFS,
		mockCompiler,
		nil,
		nil,
		nil,
	)

	options := ports.TranspileOptions{}

	_, err := service.TranspileFile("", options)
	if err == nil {
		t.Error("Expected error for empty file path, got nil")
	}
}

func TestTranspileFileNotFound(t *testing.T) {
	mockFS := mocks.NewMockFileSystem()
	mockCompiler := mocks.NewMockGoCompiler()

	service := services.NewTranspilationService(
		mockFS,
		mockCompiler,
		nil,
		nil,
		nil,
	)

	options := ports.TranspileOptions{}

	result, err := service.TranspileFile("/nonexistent.ts", options)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	if result == nil {
		t.Error("Expected result even with error, got nil")
	} else if result.Success {
		t.Error("Expected result.Success to be false")
	}
}

func TestAnalyzeProjectEmptyPath(t *testing.T) {
	mockFS := mocks.NewMockFileSystem()
	mockCompiler := mocks.NewMockGoCompiler()

	service := services.NewTranspilationService(
		mockFS,
		mockCompiler,
		nil,
		nil,
		nil,
	)

	_, err := service.AnalyzeProject("")
	if err == nil {
		t.Error("Expected error for empty project path, got nil")
	}
}

func TestAnalyzeProjectSuccess(t *testing.T) {
	mockFS := mocks.NewMockFileSystem()
	mockCompiler := mocks.NewMockGoCompiler()

	// Setup mock filesystem
	tsFile := domain.NewFile("/project/main.ts", "const x = 1;")
	mockFS.SetScanResult([]domain.File{*tsFile})

	service := services.NewTranspilationService(
		mockFS,
		mockCompiler,
		nil,
		nil,
		nil,
	)

	report, err := service.AnalyzeProject("/project")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if report == nil {
		t.Fatal("Expected report, got nil")
	}

	if report.ProjectPath != "/project" {
		t.Errorf("Expected project path '/project', got %s", report.ProjectPath)
	}
}
