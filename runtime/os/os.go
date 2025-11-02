package os

import (
	"os"
	"runtime"
	"unsafe"
)

// OS provides Node.js os module functionality
type OS struct{}

// Hostname returns the hostname of the operating system
func (o *OS) Hostname() (string, error) {
	return os.Hostname()
}

// Tmpdir returns the operating system's default directory for temporary files
func (o *OS) Tmpdir() string {
	return os.TempDir()
}

// Homedir returns the home directory of the current user
func (o *OS) Homedir() (string, error) {
	return os.UserHomeDir()
}

// Platform returns the operating system platform
func (o *OS) Platform() string {
	return runtime.GOOS
}

// Arch returns the CPU architecture
func (o *OS) Arch() string {
	return runtime.GOARCH
}

// Uptime returns the system uptime in seconds
// Note: This is a simplified version - accurate uptime requires platform-specific code
func (o *OS) Uptime() float64 {
	// This is a placeholder that returns 0
	// A real implementation would use platform-specific syscalls
	return 0
}

// Cpus returns information about each CPU/core
func (o *OS) Cpus() []CPUInfo {
	numCPU := runtime.NumCPU()
	cpus := make([]CPUInfo, numCPU)

	for i := 0; i < numCPU; i++ {
		cpus[i] = CPUInfo{
			Model: runtime.GOARCH, // Simplified - real impl would read /proc/cpuinfo or equivalent
			Speed: 0,              // Not available in standard Go
			Times: CPUTimes{
				User: 0,
				Nice: 0,
				Sys:  0,
				Idle: 0,
				Irq:  0,
			},
		}
	}

	return cpus
}

// CPUInfo represents information about a CPU/core
type CPUInfo struct {
	Model string
	Speed int64
	Times CPUTimes
}

// CPUTimes represents CPU time information
type CPUTimes struct {
	User int64
	Nice int64
	Sys  int64
	Idle int64
	Irq  int64
}

// Totalmem returns total system memory in bytes
func (o *OS) Totalmem() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// This is not accurate - Go doesn't expose total system memory directly
	// A real implementation would use platform-specific syscalls
	return m.Sys
}

// Freemem returns free system memory in bytes
func (o *OS) Freemem() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// This is not accurate - Go doesn't expose free system memory directly
	return m.Sys - m.HeapAlloc
}

// Type returns the operating system name
func (o *OS) Type() string {
	return runtime.GOOS
}

// Release returns the operating system release
func (o *OS) Release() string {
	// Not available in standard Go without platform-specific code
	return ""
}

// Endianness returns the CPU endianness
func (o *OS) Endianness() string {
	// Detect endianness
	var i int32 = 0x01020304
	u := (*[4]byte)(unsafe.Pointer(&i))
	if u[0] == 0x04 {
		return "LE" // Little Endian
	}
	return "BE" // Big Endian
}

// Global OS instance
var Default = &OS{}

// --- Package-level convenience functions ---

// Hostname returns the hostname of the operating system
func Hostname() (string, error) {
	return Default.Hostname()
}

// Tmpdir returns the operating system's default directory for temporary files
func Tmpdir() string {
	return Default.Tmpdir()
}

// Homedir returns the home directory of the current user
func Homedir() (string, error) {
	return Default.Homedir()
}

// Platform returns the operating system platform
func Platform() string {
	return Default.Platform()
}

// Arch returns the CPU architecture
func Arch() string {
	return Default.Arch()
}

// Uptime returns the system uptime in seconds
func Uptime() float64 {
	return Default.Uptime()
}

// Cpus returns information about each CPU/core
func Cpus() []CPUInfo {
	return Default.Cpus()
}

// Totalmem returns total system memory in bytes
func Totalmem() uint64 {
	return Default.Totalmem()
}

// Freemem returns free system memory in bytes
func Freemem() uint64 {
	return Default.Freemem()
}

// Type returns the operating system name
func Type() string {
	return Default.Type()
}

// Release returns the operating system release
func Release() string {
	return Default.Release()
}

// Endianness returns the CPU endianness
func Endianness() string {
	return Default.Endianness()
}

// EOL returns the end-of-line marker for the current OS
func EOL() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
