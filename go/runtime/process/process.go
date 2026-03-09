package process

import (
	"os"
	"runtime"
	"strings"
)

// Process represents the Node.js process object
type Process struct {
	// Argv holds command-line arguments
	Argv []string
	// Env holds environment variables
	Env map[string]string
	// Platform is the operating system platform
	Platform string
	// Arch is the CPU architecture
	Arch string
	// Version is the Go version (equivalent to Node version)
	Version string
	// Pid is the process ID
	Pid int
	// ExecPath is the path to the executable
	ExecPath string
}

// Default is the global process instance
var Default *Process

func init() {
	Default = NewProcess()
}

// NewProcess creates a new Process instance
func NewProcess() *Process {
	execPath, _ := os.Executable()

	// Get environment variables as map
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}

	return &Process{
		Argv:     os.Args,
		Env:      envMap,
		Platform: runtime.GOOS,
		Arch:     runtime.GOARCH,
		Version:  runtime.Version(),
		Pid:      os.Getpid(),
		ExecPath: execPath,
	}
}

// Cwd returns the current working directory
func (p *Process) Cwd() (string, error) {
	return os.Getwd()
}

// Chdir changes the current working directory
func (p *Process) Chdir(dir string) error {
	return os.Chdir(dir)
}

// Exit terminates the process with the given exit code
func (p *Process) Exit(code int) {
	os.Exit(code)
}

// Getenv returns the value of an environment variable
func (p *Process) Getenv(key string) string {
	return os.Getenv(key)
}

// Setenv sets the value of an environment variable
func (p *Process) Setenv(key, value string) error {
	if err := os.Setenv(key, value); err != nil {
		return err
	}
	p.Env[key] = value
	return nil
}

// Unsetenv removes an environment variable
func (p *Process) Unsetenv(key string) error {
	if err := os.Unsetenv(key); err != nil {
		return err
	}
	delete(p.Env, key)
	return nil
}

// Umask sets the process umask (Unix only)
func (p *Process) Umask(mask int) int {
	// Note: Go doesn't provide a cross-platform umask function
	// This is a placeholder that would need platform-specific implementation
	return 0
}

// GetUID returns the user ID (Unix only)
func (p *Process) GetUID() int {
	return os.Getuid()
}

// GetGID returns the group ID (Unix only)
func (p *Process) GetGID() int {
	return os.Getgid()
}

// Getgroups returns supplementary group IDs (Unix only)
func (p *Process) Getgroups() ([]int, error) {
	return os.Getgroups()
}

// GetEUID returns the effective user ID (Unix only)
func (p *Process) GetEUID() int {
	return os.Geteuid()
}

// GetEGID returns the effective group ID (Unix only)
func (p *Process) GetEGID() int {
	return os.Getegid()
}

// Uptime returns the number of seconds the process has been running
// Note: This is a simplified version - real Node.js tracks actual uptime
func (p *Process) Uptime() float64 {
	// Would need to track start time to implement accurately
	return 0.0
}

// Hrtime returns high-resolution real time in [seconds, nanoseconds]
func (p *Process) Hrtime() [2]int64 {
	// This is a simplified version
	// Real implementation would use time.Now().UnixNano()
	return [2]int64{0, 0}
}

// MemoryUsage returns memory usage information
func (p *Process) MemoryUsage() MemoryUsage {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemoryUsage{
		RSS:          int64(m.Sys),
		HeapTotal:    int64(m.HeapSys),
		HeapUsed:     int64(m.HeapAlloc),
		External:     0, // Go doesn't track this separately
		ArrayBuffers: 0, // Not applicable in Go
	}
}

// MemoryUsage represents memory usage information
type MemoryUsage struct {
	RSS          int64 // Resident Set Size
	HeapTotal    int64 // Total heap size
	HeapUsed     int64 // Used heap size
	External     int64 // C++ objects bound to JS objects
	ArrayBuffers int64 // ArrayBuffers and SharedArrayBuffers
}

// CPUUsage returns CPU usage information
func (p *Process) CPUUsage() CPUUsage {
	// This is a placeholder - accurate CPU usage tracking would require
	// platform-specific code and tracking deltas
	return CPUUsage{
		User:   0,
		System: 0,
	}
}

// CPUUsage represents CPU usage information in microseconds
type CPUUsage struct {
	User   int64 // User CPU time
	System int64 // System CPU time
}

// --- Package-level convenience functions ---

// Cwd returns the current working directory
func Cwd() (string, error) {
	return Default.Cwd()
}

// Chdir changes the current working directory
func Chdir(dir string) error {
	return Default.Chdir(dir)
}

// Exit terminates the process with the given exit code
func Exit(code int) {
	Default.Exit(code)
}

// Env returns environment variables map
func Env() map[string]string {
	return Default.Env
}

// Argv returns command-line arguments
func Argv() []string {
	return Default.Argv
}

// Platform returns the operating system platform
func Platform() string {
	return Default.Platform
}

// Arch returns the CPU architecture
func Arch() string {
	return Default.Arch
}

// Version returns the runtime version
func Version() string {
	return Default.Version
}

// Pid returns the process ID
func Pid() int {
	return Default.Pid
}

// ExecPath returns the path to the executable
func ExecPath() string {
	return Default.ExecPath
}

// Getenv returns the value of an environment variable
func Getenv(key string) string {
	return Default.Getenv(key)
}

// Setenv sets the value of an environment variable
func Setenv(key, value string) error {
	return Default.Setenv(key, value)
}

// GetMemoryUsage returns memory usage information
func GetMemoryUsage() MemoryUsage {
	return Default.MemoryUsage()
}

// GetCPUUsage returns CPU usage information
func GetCPUUsage() CPUUsage {
	return Default.CPUUsage()
}
