package services

import (
	"fmt"

	"github.com/el-j/ts2go/core/domain"
	"github.com/el-j/ts2go/core/ports"
)

// GoRuntimeServiceImpl implements the GoRuntimeService port
type GoRuntimeServiceImpl struct {
	compiler ports.GoCompiler
	fs       ports.FileSystem
}

// NewGoRuntimeService creates a new Go runtime service
func NewGoRuntimeService(compiler ports.GoCompiler, fs ports.FileSystem) ports.GoRuntimeService {
	return &GoRuntimeServiceImpl{
		compiler: compiler,
		fs:       fs,
	}
}

// DetectGoInstallation detects if Go is installed and returns version info
func (s *GoRuntimeServiceImpl) DetectGoInstallation() (*domain.GoVersion, error) {
	version, err := s.compiler.Detect()
	if err != nil {
		return nil, fmt.Errorf("failed to detect Go installation: %w", err)
	}
	return version, nil
}

// BuildProject builds a Go project
func (s *GoRuntimeServiceImpl) BuildProject(projectPath string, outputPath string) (*domain.BuildResult, error) {
	if projectPath == "" {
		return nil, domain.NewValidationError("project path cannot be empty")
	}

	// Check if project exists
	exists, err := s.fs.PathExists(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to check project path: %w", err)
	}
	if !exists {
		return nil, domain.NewValidationError("project path does not exist: %s", projectPath)
	}

	// Build using compiler
	opts := ports.BuildOptions{}

	output, err := s.compiler.Build(projectPath, outputPath, opts)
	if err != nil {
		return &domain.BuildResult{
			Success: false,
			Output:  output,
			Error:   err.Error(),
		}, err
	}

	return &domain.BuildResult{
		Success:    true,
		Output:     output,
		OutputPath: outputPath,
	}, nil
}

// RunProject runs a Go project
func (s *GoRuntimeServiceImpl) RunProject(projectPath string, args []string) (*domain.RunResult, error) {
	if projectPath == "" {
		return nil, domain.NewValidationError("project path cannot be empty")
	}

	// Check if project exists
	exists, err := s.fs.PathExists(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to check project path: %w", err)
	}
	if !exists {
		return nil, domain.NewValidationError("project path does not exist: %s", projectPath)
	}

	// Run using compiler
	stdout, stderr, exitCode, err := s.compiler.Run(projectPath, args)

	result := &domain.RunResult{
		Stdout:   stdout,
		Stderr:   stderr,
		ExitCode: exitCode,
		Success:  exitCode == 0,
	}

	if err != nil {
		result.Error = err.Error()
		return result, fmt.Errorf("run failed: %w", err)
	}

	return result, nil
}

// TestProject runs tests for a Go project
func (s *GoRuntimeServiceImpl) TestProject(projectPath string, options *ports.TestOptions) (*domain.TestResult, error) {
	if projectPath == "" {
		return nil, domain.NewValidationError("project path cannot be empty")
	}

	// Check if project exists
	exists, err := s.fs.PathExists(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to check project path: %w", err)
	}
	if !exists {
		return nil, domain.NewValidationError("project path does not exist: %s", projectPath)
	}

	// Use default options if not provided
	if options == nil {
		options = &ports.TestOptions{}
	}

	// Run tests using compiler
	output, err := s.compiler.Test(projectPath, *options)

	result := &domain.TestResult{
		Output:  output,
		Success: err == nil,
	}

	if err != nil {
		result.Error = err.Error()
		return result, fmt.Errorf("test failed: %w", err)
	}

	return result, nil
}

// FormatProject formats all Go files in a project
func (s *GoRuntimeServiceImpl) FormatProject(projectPath string) error {
	if projectPath == "" {
		return domain.NewValidationError("project path cannot be empty")
	}

	// Check if project exists
	exists, err := s.fs.PathExists(projectPath)
	if err != nil {
		return fmt.Errorf("failed to check project path: %w", err)
	}
	if !exists {
		return domain.NewValidationError("project path does not exist: %s", projectPath)
	}

	// Format using compiler
	_, err = s.compiler.Format(projectPath)
	if err != nil {
		return fmt.Errorf("format failed: %w", err)
	}

	return nil
}
