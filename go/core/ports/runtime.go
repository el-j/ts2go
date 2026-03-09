package ports

import "github.com/el-j/ts2go/core/domain"

// GoRuntimeService defines the use cases for running Go code.
// This is a DRIVING PORT (inbound) - it's the API of our core logic.
type GoRuntimeService interface {
	// DetectGoInstallation detects if Go is installed and returns version info
	DetectGoInstallation() (*domain.GoVersion, error)

	// BuildProject builds a Go project
	BuildProject(projectPath string, outputPath string) (*domain.BuildResult, error)

	// RunProject runs a Go project
	RunProject(projectPath string, args []string) (*domain.RunResult, error)

	// TestProject runs tests for a Go project
	TestProject(projectPath string, options *TestOptions) (*domain.TestResult, error)

	// FormatProject formats all Go files in a project
	FormatProject(projectPath string) error
}
