package process

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNewProcess(t *testing.T) {
	p := NewProcess()

	if p == nil {
		t.Fatal("NewProcess returned nil")
	}

	if len(p.Argv) == 0 {
		t.Error("Expected Argv to have arguments")
	}

	if p.Platform == "" {
		t.Error("Expected Platform to be set")
	}

	if p.Arch == "" {
		t.Error("Expected Arch to be set")
	}

	if p.Pid == 0 {
		t.Error("Expected Pid to be set")
	}
}

func TestProcess_Cwd(t *testing.T) {
	p := NewProcess()

	cwd, err := p.Cwd()
	if err != nil {
		t.Fatalf("Cwd failed: %v", err)
	}

	if cwd == "" {
		t.Error("Expected non-empty current working directory")
	}

	t.Logf("Current working directory: %s", cwd)
}

func TestProcess_Chdir(t *testing.T) {
	p := NewProcess()

	// Get original directory
	originalDir, err := p.Cwd()
	if err != nil {
		t.Fatalf("Failed to get original directory: %v", err)
	}

	// Change to temp directory
	tmpDir := os.TempDir()
	if err := p.Chdir(tmpDir); err != nil {
		t.Fatalf("Chdir failed: %v", err)
	}

	// Verify we changed - need to resolve symlinks for comparison (macOS has /var -> /private/var)
	newDir, err := p.Cwd()
	if err != nil {
		t.Fatalf("Failed to get new directory: %v", err)
	}

	expectedDir, _ := filepath.EvalSymlinks(tmpDir)
	actualDir, _ := filepath.EvalSymlinks(newDir)

	if actualDir != expectedDir {
		t.Errorf("Expected directory %s, got %s", expectedDir, actualDir)
	}

	// Change back
	if err := p.Chdir(originalDir); err != nil {
		t.Fatalf("Failed to restore original directory: %v", err)
	}
}

func TestProcess_Env(t *testing.T) {
	p := NewProcess()

	if len(p.Env) == 0 {
		t.Error("Expected environment variables to be populated")
	}

	// Test getting a standard env var (PATH should exist on most systems)
	path := p.Getenv("PATH")
	if path == "" {
		// Try HOME as fallback
		home := p.Getenv("HOME")
		if home == "" && runtime.GOOS != "windows" {
			t.Error("Expected at least PATH or HOME to be set")
		}
	}
}

func TestProcess_SetenvAndGetenv(t *testing.T) {
	p := NewProcess()

	key := "TEST_VAR_12345"
	value := "test_value"

	// Set environment variable
	if err := p.Setenv(key, value); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}

	// Get it back
	result := p.Getenv(key)
	if result != value {
		t.Errorf("Expected %s, got %s", value, result)
	}

	// Check it's in the map
	if p.Env[key] != value {
		t.Errorf("Expected Env[%s] = %s, got %s", key, value, p.Env[key])
	}

	// Clean up
	p.Unsetenv(key)
}

func TestProcess_Unsetenv(t *testing.T) {
	p := NewProcess()

	key := "TEST_VAR_UNSET_12345"
	value := "test_value"

	// Set it first
	p.Setenv(key, value)

	// Unset it
	if err := p.Unsetenv(key); err != nil {
		t.Fatalf("Unsetenv failed: %v", err)
	}

	// Verify it's gone
	result := p.Getenv(key)
	if result != "" {
		t.Errorf("Expected empty string after unset, got %s", result)
	}

	// Check it's not in the map
	if _, exists := p.Env[key]; exists {
		t.Error("Expected key to be removed from Env map")
	}
}

func TestProcess_Platform(t *testing.T) {
	p := NewProcess()

	expectedPlatforms := []string{"linux", "darwin", "windows", "freebsd", "openbsd"}
	found := false
	for _, platform := range expectedPlatforms {
		if p.Platform == platform {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Unexpected platform: %s", p.Platform)
	}

	t.Logf("Platform: %s", p.Platform)
}

func TestProcess_Arch(t *testing.T) {
	p := NewProcess()

	expectedArchs := []string{"amd64", "386", "arm", "arm64"}
	found := false
	for _, arch := range expectedArchs {
		if p.Arch == arch {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Unexpected architecture: %s", p.Arch)
	}

	t.Logf("Architecture: %s", p.Arch)
}

func TestProcess_Version(t *testing.T) {
	p := NewProcess()

	if p.Version == "" {
		t.Error("Expected Version to be set")
	}

	if !strings.HasPrefix(p.Version, "go") {
		t.Errorf("Expected version to start with 'go', got %s", p.Version)
	}

	t.Logf("Version: %s", p.Version)
}

func TestProcess_Pid(t *testing.T) {
	p := NewProcess()

	if p.Pid <= 0 {
		t.Errorf("Expected positive PID, got %d", p.Pid)
	}

	// Verify it matches os.Getpid()
	if p.Pid != os.Getpid() {
		t.Errorf("Expected PID %d, got %d", os.Getpid(), p.Pid)
	}

	t.Logf("PID: %d", p.Pid)
}

func TestProcess_ExecPath(t *testing.T) {
	p := NewProcess()

	if p.ExecPath == "" {
		t.Error("Expected ExecPath to be set")
	}

	t.Logf("ExecPath: %s", p.ExecPath)
}

func TestProcess_MemoryUsage(t *testing.T) {
	p := NewProcess()

	mem := p.MemoryUsage()

	if mem.RSS == 0 {
		t.Error("Expected non-zero RSS")
	}

	if mem.HeapTotal == 0 {
		t.Error("Expected non-zero HeapTotal")
	}

	if mem.HeapUsed == 0 {
		t.Error("Expected non-zero HeapUsed")
	}

	t.Logf("Memory Usage: RSS=%d, HeapTotal=%d, HeapUsed=%d",
		mem.RSS, mem.HeapTotal, mem.HeapUsed)
}

func TestProcess_GetUID_Unix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping Unix-specific test on Windows")
	}

	p := NewProcess()
	uid := p.GetUID()

	// UID should be non-negative
	if uid < 0 {
		t.Errorf("Expected non-negative UID, got %d", uid)
	}

	t.Logf("UID: %d", uid)
}

func TestProcess_GetGID_Unix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping Unix-specific test on Windows")
	}

	p := NewProcess()
	gid := p.GetGID()

	// GID should be non-negative
	if gid < 0 {
		t.Errorf("Expected non-negative GID, got %d", gid)
	}

	t.Logf("GID: %d", gid)
}

func TestPackageFunctions(t *testing.T) {
	// Test package-level convenience functions

	cwd, err := Cwd()
	if err != nil {
		t.Fatalf("Cwd() failed: %v", err)
	}
	if cwd == "" {
		t.Error("Expected non-empty cwd")
	}

	env := Env()
	if len(env) == 0 {
		t.Error("Expected non-empty env")
	}

	argv := Argv()
	if len(argv) == 0 {
		t.Error("Expected non-empty argv")
	}

	platform := Platform()
	if platform == "" {
		t.Error("Expected non-empty platform")
	}

	arch := Arch()
	if arch == "" {
		t.Error("Expected non-empty arch")
	}

	version := Version()
	if version == "" {
		t.Error("Expected non-empty version")
	}

	pid := Pid()
	if pid <= 0 {
		t.Error("Expected positive pid")
	}
}

func TestDefault(t *testing.T) {
	if Default == nil {
		t.Fatal("Default process should be initialized")
	}

	if Default.Pid != os.Getpid() {
		t.Error("Default process should have current PID")
	}
}
