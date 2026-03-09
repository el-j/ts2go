package domain

import "time"

// GoVersion represents information about an installed Go compiler
type GoVersion struct {
	Version string `json:"version"` // e.g., "go1.21.5"
	Path    string `json:"path"`    // Path to go binary
	GOROOT  string `json:"goroot"`  // GOROOT environment variable
	GOPATH  string `json:"gopath"`  // GOPATH environment variable
	OS      string `json:"os"`      // Target OS
	Arch    string `json:"arch"`    // Target architecture
	Found   bool   `json:"found"`   // Whether Go was found
}

// BuildResult represents the result of building a Go project
type BuildResult struct {
	Success    bool          `json:"success"`
	OutputPath string        `json:"outputPath"`
	Duration   time.Duration `json:"duration"`
	Output     string        `json:"output"`
	Error      string        `json:"error,omitempty"`
	Warnings   []string      `json:"warnings,omitempty"`
	BuildTime  time.Time     `json:"buildTime"`
}

// RunResult represents the result of running a Go program
type RunResult struct {
	Success     bool          `json:"success"`
	ExitCode    int           `json:"exitCode"`
	Stdout      string        `json:"stdout"`
	Stderr      string        `json:"stderr"`
	Error       string        `json:"error,omitempty"`
	Duration    time.Duration `json:"duration"`
	StartedAt   time.Time     `json:"startedAt"`
	CompletedAt time.Time     `json:"completedAt"`
}

// TestResult represents the result of running Go tests
type TestResult struct {
	Success        bool          `json:"success"`
	TotalTests     int           `json:"totalTests"`
	PassedTests    int           `json:"passedTests"`
	FailedTests    int           `json:"failedTests"`
	SkippedTests   int           `json:"skippedTests"`
	Coverage       float64       `json:"coverage"` // 0-100
	Duration       time.Duration `json:"duration"`
	Output         string        `json:"output"`
	Error          string        `json:"error,omitempty"`
	FailureDetails []TestFailure `json:"failureDetails,omitempty"`
}

// TestFailure represents a single test failure
type TestFailure struct {
	TestName string `json:"testName"`
	Package  string `json:"package"`
	Message  string `json:"message"`
	Location string `json:"location"`
}

// NewBuildResult creates a new BuildResult
func NewBuildResult(outputPath string) *BuildResult {
	return &BuildResult{
		OutputPath: outputPath,
		Warnings:   make([]string, 0),
		BuildTime:  time.Now(),
	}
}

// NewRunResult creates a new RunResult
func NewRunResult() *RunResult {
	return &RunResult{
		StartedAt: time.Now(),
	}
}

// Complete marks the run as complete
func (r *RunResult) Complete() {
	r.CompletedAt = time.Now()
	r.Duration = r.CompletedAt.Sub(r.StartedAt)
}

// NewTestResult creates a new TestResult
func NewTestResult() *TestResult {
	return &TestResult{
		FailureDetails: make([]TestFailure, 0),
	}
}
