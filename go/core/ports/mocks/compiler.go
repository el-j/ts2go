package mocks

import (
	"github.com/el-j/ts2go/core/domain"
	"github.com/el-j/ts2go/core/ports"
)

// MockGoCompiler is a mock implementation of ports.GoCompiler for testing
type MockGoCompiler struct {
	// DetectResult is the GoVersion to return from Detect
	DetectResult *domain.GoVersion

	// DetectError is the error to return from Detect
	DetectError error

	// BuildResult is the output to return from Build
	BuildResult string

	// BuildError is the error to return from Build
	BuildError error

	// BuildCalled tracks if Build was called
	BuildCalled bool

	// BuildArgs tracks the arguments passed to Build
	BuildArgs struct {
		SourcePath string
		OutputPath string
		Options    ports.BuildOptions
	}

	// RunResult tracks the result to return from Run
	RunResult struct {
		Stdout   string
		Stderr   string
		ExitCode int
	}

	// RunError is the error to return from Run
	RunError error

	// TestResult is the output to return from Test
	TestResult string

	// TestError is the error to return from Test
	TestError error

	// FormatResult is the formatted code to return
	FormatResult string

	// FormatError is the error to return from Format
	FormatError error

	// ModInitError is the error to return from ModInit
	ModInitError error

	// ModTidyError is the error to return from ModTidy
	ModTidyError error

	// GetError is the error to return from Get
	GetError error
}

// NewMockGoCompiler creates a new mock Go compiler
func NewMockGoCompiler() *MockGoCompiler {
	return &MockGoCompiler{
		DetectResult: &domain.GoVersion{
			Version: "go1.21.5",
			Path:    "/usr/local/go/bin/go",
			GOROOT:  "/usr/local/go",
			GOPATH:  "/Users/test/go",
			OS:      "darwin",
			Arch:    "arm64",
			Found:   true,
		},
	}
}

// Detect returns the configured Go version
func (m *MockGoCompiler) Detect() (*domain.GoVersion, error) {
	if m.DetectError != nil {
		return nil, m.DetectError
	}
	return m.DetectResult, nil
}

// Build simulates building a Go package
func (m *MockGoCompiler) Build(sourcePath, outputPath string, options ports.BuildOptions) (string, error) {
	m.BuildCalled = true
	m.BuildArgs.SourcePath = sourcePath
	m.BuildArgs.OutputPath = outputPath
	m.BuildArgs.Options = options

	if m.BuildError != nil {
		return "", m.BuildError
	}
	return m.BuildResult, nil
}

// Run simulates running a Go program
func (m *MockGoCompiler) Run(sourcePath string, args []string) (stdout, stderr string, exitCode int, err error) {
	if m.RunError != nil {
		return "", "", 1, m.RunError
	}
	return m.RunResult.Stdout, m.RunResult.Stderr, m.RunResult.ExitCode, nil
}

// Test simulates running Go tests
func (m *MockGoCompiler) Test(packagePath string, options ports.TestOptions) (string, error) {
	if m.TestError != nil {
		return "", m.TestError
	}
	return m.TestResult, nil
}

// Format simulates formatting Go code
func (m *MockGoCompiler) Format(path string) (string, error) {
	if m.FormatError != nil {
		return "", m.FormatError
	}
	return m.FormatResult, nil
}

// ModInit simulates initializing a Go module
func (m *MockGoCompiler) ModInit(path, moduleName string) error {
	return m.ModInitError
}

// ModTidy simulates tidying a Go module
func (m *MockGoCompiler) ModTidy(path string) error {
	return m.ModTidyError
}

// Get simulates downloading a Go package
func (m *MockGoCompiler) Get(packagePath string) error {
	return m.GetError
}

// SetBuildSuccess configures the mock to return a successful build
func (m *MockGoCompiler) SetBuildSuccess(output string) {
	m.BuildResult = output
	m.BuildError = nil
}

// SetBuildFailure configures the mock to return a build failure
func (m *MockGoCompiler) SetBuildFailure(err error) {
	m.BuildError = err
	m.BuildResult = ""
}

// SetTestSuccess configures the mock to return successful tests
func (m *MockGoCompiler) SetTestSuccess(output string) {
	m.TestResult = output
	m.TestError = nil
}

// SetTestFailure configures the mock to return test failures
func (m *MockGoCompiler) SetTestFailure(err error) {
	m.TestError = err
	m.TestResult = ""
}
