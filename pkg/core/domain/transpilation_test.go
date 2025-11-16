package domain_test

import (
	"testing"

	"github.com/el-j/ts2go/pkg/core/domain"
)

func TestNewTranspilationResult(t *testing.T) {
	projectPath := "/path/to/project"
	outputPath := "/output/path"

	result := domain.NewTranspilationResult(projectPath, outputPath)

	if result.ProjectPath != projectPath {
		t.Errorf("Expected project path %s, got %s", projectPath, result.ProjectPath)
	}
	if result.OutputPath != outputPath {
		t.Errorf("Expected output path %s, got %s", outputPath, result.OutputPath)
	}
	if result.Success != false {
		t.Error("Expected initial success to be false")
	}
	if result.FilesProcessed != 0 {
		t.Errorf("Expected 0 files processed, got %d", result.FilesProcessed)
	}
	if result.Errors == nil {
		t.Error("Expected errors slice to be initialized")
	}
	if result.Warnings == nil {
		t.Error("Expected warnings slice to be initialized")
	}
}

func TestTranspilationResultAddFileResult(t *testing.T) {
	result := domain.NewTranspilationResult("/path", "/output")

	fileResult := domain.FileTranspilationResult{
		SourcePath: "/path/test.ts",
		TargetPath: "/output/test.go",
		Success:    true,
		GoCode:     "package main",
	}

	result.AddFileResult(fileResult)

	if result.FilesProcessed != 1 {
		t.Errorf("Expected 1 file processed, got %d", result.FilesProcessed)
	}
	if result.FilesSuccess != 1 {
		t.Errorf("Expected 1 file success, got %d", result.FilesSuccess)
	}
	if result.FilesFailed != 0 {
		t.Errorf("Expected 0 files failed, got %d", result.FilesFailed)
	}
}

func TestTranspilationResultAddFailedFileResult(t *testing.T) {
	result := domain.NewTranspilationResult("/path", "/output")

	fileResult := domain.FileTranspilationResult{
		SourcePath: "/path/test.ts",
		Success:    false,
		Error:      "transpilation error",
	}

	result.AddFileResult(fileResult)

	if result.FilesProcessed != 1 {
		t.Errorf("Expected 1 file processed, got %d", result.FilesProcessed)
	}
	if result.FilesFailed != 1 {
		t.Errorf("Expected 1 file failed, got %d", result.FilesFailed)
	}
	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}
}

func TestTranspilationResultComplete(t *testing.T) {
	result := domain.NewTranspilationResult("/path", "/output")

	// Add successful file
	result.AddFileResult(domain.FileTranspilationResult{
		SourcePath: "/path/test.ts",
		Success:    true,
	})

	result.Complete()

	if !result.Success {
		t.Error("Expected success to be true when no failures")
	}
	if result.CompletedAt.IsZero() {
		t.Error("Expected CompletedAt to be set")
	}
	if result.Duration == 0 {
		t.Error("Expected Duration to be set")
	}
}
