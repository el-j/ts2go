package gocompiler

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/el-j/ts2go/core/domain"
	"github.com/el-j/ts2go/core/ports"
)

// SystemGoCompiler implements the GoCompiler port using system Go installation
type SystemGoCompiler struct {
	goBinary string
}

// NewSystemGoCompiler creates a new system Go compiler adapter
func NewSystemGoCompiler() (ports.GoCompiler, error) {
	// Find Go binary
	goBinary, err := exec.LookPath("go")
	if err != nil {
		return nil, domain.ErrGoNotFound()
	}

	return &SystemGoCompiler{
		goBinary: goBinary,
	}, nil
}

// Detect detects the Go installation and returns version information
func (c *SystemGoCompiler) Detect() (*domain.GoVersion, error) {
	// Get version
	versionCmd := exec.Command(c.goBinary, "version")
	versionOutput, err := versionCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get Go version: %w", err)
	}

	version := strings.TrimSpace(string(versionOutput))

	// Get environment
	envCmd := exec.Command(c.goBinary, "env", "GOROOT", "GOPATH")
	envOutput, err := envCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get Go environment: %w", err)
	}

	envLines := strings.Split(strings.TrimSpace(string(envOutput)), "\n")
	goroot := ""
	gopath := ""
	if len(envLines) >= 1 {
		goroot = envLines[0]
	}
	if len(envLines) >= 2 {
		gopath = envLines[1]
	}

	return &domain.GoVersion{
		Version: version,
		Path:    c.goBinary,
		GOROOT:  goroot,
		GOPATH:  gopath,
		OS:      runtime.GOOS,
		Arch:    runtime.GOARCH,
		Found:   true,
	}, nil
}

// Build builds a Go package or binary
func (c *SystemGoCompiler) Build(sourcePath, outputPath string, options ports.BuildOptions) (string, error) {
	args := []string{"build"}

	// Add output path
	if outputPath != "" {
		args = append(args, "-o", outputPath)
	}

	// Add build tags
	if len(options.Tags) > 0 {
		args = append(args, "-tags", strings.Join(options.Tags, ","))
	}

	// Add optimization flags
	if !options.Optimize {
		args = append(args, "-gcflags", "all=-N -l")
	}

	// Add race detector
	if options.Race {
		args = append(args, "-race")
	}

	// Add verbose output
	if options.Verbose {
		args = append(args, "-v")
	}

	// Add source path (can be package or file)
	args = append(args, sourcePath)

	// Execute build
	cmd := exec.Command(c.goBinary, args...)
	cmd.Dir = filepath.Dir(sourcePath)

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		return outputStr, fmt.Errorf("build failed: %w", err)
	}

	return outputStr, nil
}

// Run executes a Go program
func (c *SystemGoCompiler) Run(sourcePath string, args []string) (stdout, stderr string, exitCode int, err error) {
	cmdArgs := []string{"run", sourcePath}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command(c.goBinary, cmdArgs...)
	cmd.Dir = filepath.Dir(sourcePath)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()
	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
		return stdout, stderr, exitCode, err
	}

	return stdout, stderr, 0, nil
}

// Test runs Go tests
func (c *SystemGoCompiler) Test(packagePath string, options ports.TestOptions) (string, error) {
	args := []string{"test"}

	// Add coverage
	if options.Coverage {
		args = append(args, "-cover")
	}

	// Add verbose output
	if options.Verbose {
		args = append(args, "-v")
	}

	// Add short flag
	if options.Short {
		args = append(args, "-short")
	}

	// Add race detector
	if options.Race {
		args = append(args, "-race")
	}

	// Add timeout
	if options.Timeout != "" {
		args = append(args, "-timeout", options.Timeout)
	}

	// Add test pattern
	if options.Pattern != "" {
		args = append(args, "-run", options.Pattern)
	}

	// Add package path
	args = append(args, packagePath)

	// Execute tests
	cmd := exec.Command(c.goBinary, args...)
	if packagePath != "." && packagePath != "" {
		cmd.Dir = packagePath
	}

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		return outputStr, fmt.Errorf("tests failed: %w", err)
	}

	return outputStr, nil
}

// Format formats Go source code using gofmt
func (c *SystemGoCompiler) Format(path string) (string, error) {
	// Use gofmt command
	cmd := exec.Command("gofmt", "-w", path)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("format failed: %w", err)
	}

	return string(output), nil
}

// ModInit initializes a new Go module
func (c *SystemGoCompiler) ModInit(path, moduleName string) error {
	cmd := exec.Command(c.goBinary, "mod", "init", moduleName)
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mod init failed: %s: %w", string(output), err)
	}

	return nil
}

// ModTidy cleans up go.mod and go.sum
func (c *SystemGoCompiler) ModTidy(path string) error {
	cmd := exec.Command(c.goBinary, "mod", "tidy")
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mod tidy failed: %s: %w", string(output), err)
	}

	return nil
}

// Get downloads and installs packages
func (c *SystemGoCompiler) Get(packagePath string) error {
	cmd := exec.Command(c.goBinary, "get", packagePath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go get failed: %s: %w", string(output), err)
	}

	return nil
}
