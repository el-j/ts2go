package ports

import "github.com/el-j/ts2go/core/domain"

// GoCompiler is what the core uses to interact with the Go compiler.
// This is a DRIVEN PORT (outbound) - a dependency the core needs.
// Implementations will use exec.Command, but the core doesn't know that.
type GoCompiler interface {
	// Detect detects the Go installation and returns version information
	Detect() (*domain.GoVersion, error)

	// Build builds a Go package or binary
	Build(sourcePath, outputPath string, options BuildOptions) (string, error)

	// Run executes a Go program
	Run(sourcePath string, args []string) (stdout, stderr string, exitCode int, err error)

	// Test runs Go tests
	Test(packagePath string, options TestOptions) (string, error)

	// Format formats Go source code using gofmt
	Format(path string) (string, error)

	// ModInit initializes a new Go module
	ModInit(path, moduleName string) error

	// ModTidy cleans up go.mod and go.sum
	ModTidy(path string) error

	// Get downloads and installs packages
	Get(packagePath string) error
}

// BuildOptions contains options for building Go code
type BuildOptions struct {
	Tags     []string // Build tags
	GOOS     string   // Target OS
	GOARCH   string   // Target architecture
	Optimize bool     // Enable optimizations
	Race     bool     // Enable race detector
	Verbose  bool     // Verbose output
}

// TestOptions contains options for running Go tests
type TestOptions struct {
	Coverage bool   // Enable coverage
	Verbose  bool   // Verbose output
	Short    bool   // Run short tests only
	Pattern  string // Test name pattern
	Timeout  string // Test timeout (e.g., "10m")
	Race     bool   // Enable race detector
}
