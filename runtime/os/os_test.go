package os

import (
	"runtime"
	"testing"
)

func TestOS_Hostname(t *testing.T) {
	o := &OS{}
	hostname, err := o.Hostname()
	if err != nil {
		t.Fatalf("Hostname failed: %v", err)
	}

	if hostname == "" {
		t.Error("Hostname should not be empty")
	}

	t.Logf("Hostname: %s", hostname)
}

func TestOS_Tmpdir(t *testing.T) {
	o := &OS{}
	tmpdir := o.Tmpdir()

	if tmpdir == "" {
		t.Error("Tmpdir should not be empty")
	}

	t.Logf("Tmpdir: %s", tmpdir)
}

func TestOS_Homedir(t *testing.T) {
	o := &OS{}
	homedir, err := o.Homedir()
	if err != nil {
		t.Fatalf("Homedir failed: %v", err)
	}

	if homedir == "" {
		t.Error("Homedir should not be empty")
	}

	t.Logf("Homedir: %s", homedir)
}

func TestOS_Platform(t *testing.T) {
	o := &OS{}
	platform := o.Platform()

	if platform == "" {
		t.Error("Platform should not be empty")
	}

	// Should be one of: darwin, linux, windows, etc.
	validPlatforms := []string{"darwin", "linux", "windows", "freebsd", "openbsd", "netbsd", "dragonfly", "solaris", "aix", "android", "js"}
	found := false
	for _, p := range validPlatforms {
		if platform == p {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Platform %s is not recognized", platform)
	}

	t.Logf("Platform: %s", platform)
}

func TestOS_Arch(t *testing.T) {
	o := &OS{}
	arch := o.Arch()

	if arch == "" {
		t.Error("Arch should not be empty")
	}

	// Should be one of: amd64, arm64, 386, arm, etc.
	validArchs := []string{"amd64", "arm64", "386", "arm", "ppc64", "ppc64le", "mips", "mipsle", "mips64", "mips64le", "s390x", "wasm"}
	found := false
	for _, a := range validArchs {
		if arch == a {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Arch %s is not recognized", arch)
	}

	t.Logf("Arch: %s", arch)
}

func TestOS_Uptime(t *testing.T) {
	o := &OS{}
	uptime := o.Uptime()

	// For now, this is a placeholder that returns 0
	if uptime < 0 {
		t.Error("Uptime should not be negative")
	}

	t.Logf("Uptime: %.2f seconds", uptime)
}

func TestOS_Cpus(t *testing.T) {
	o := &OS{}
	cpus := o.Cpus()

	if len(cpus) == 0 {
		t.Error("Should have at least one CPU")
	}

	expectedCount := runtime.NumCPU()
	if len(cpus) != expectedCount {
		t.Errorf("Expected %d CPUs, got %d", expectedCount, len(cpus))
	}

	t.Logf("CPUs: %d cores", len(cpus))
}

func TestOS_Totalmem(t *testing.T) {
	o := &OS{}
	totalmem := o.Totalmem()

	if totalmem == 0 {
		t.Error("Totalmem should not be zero")
	}

	t.Logf("Total memory: %d bytes (%.2f GB)", totalmem, float64(totalmem)/(1024*1024*1024))
}

func TestOS_Freemem(t *testing.T) {
	o := &OS{}
	freemem := o.Freemem()

	// Free memory can be zero in some cases, so just check it's not negative
	t.Logf("Free memory: %d bytes (%.2f GB)", freemem, float64(freemem)/(1024*1024*1024))
}

func TestOS_Type(t *testing.T) {
	o := &OS{}
	osType := o.Type()

	if osType == "" {
		t.Error("Type should not be empty")
	}

	t.Logf("OS Type: %s", osType)
}

func TestOS_Endianness(t *testing.T) {
	o := &OS{}
	endianness := o.Endianness()

	if endianness != "LE" && endianness != "BE" {
		t.Errorf("Endianness should be 'LE' or 'BE', got %s", endianness)
	}

	t.Logf("Endianness: %s", endianness)
}

func TestEOL(t *testing.T) {
	eol := EOL()

	if runtime.GOOS == "windows" {
		if eol != "\r\n" {
			t.Errorf("Expected Windows EOL \\r\\n, got %q", eol)
		}
	} else {
		if eol != "\n" {
			t.Errorf("Expected Unix EOL \\n, got %q", eol)
		}
	}

	t.Logf("EOL: %q", eol)
}

// Test package-level functions
func TestPackageFunctions(t *testing.T) {
	// Test that package-level functions work
	hostname, err := Hostname()
	if err != nil {
		t.Errorf("Package Hostname() failed: %v", err)
	}
	if hostname == "" {
		t.Error("Package Hostname() returned empty string")
	}

	tmpdir := Tmpdir()
	if tmpdir == "" {
		t.Error("Package Tmpdir() returned empty string")
	}

	homedir, err := Homedir()
	if err != nil {
		t.Errorf("Package Homedir() failed: %v", err)
	}
	if homedir == "" {
		t.Error("Package Homedir() returned empty string")
	}

	platform := Platform()
	if platform == "" {
		t.Error("Package Platform() returned empty string")
	}

	arch := Arch()
	if arch == "" {
		t.Error("Package Arch() returned empty string")
	}

	uptime := Uptime()
	if uptime < 0 {
		t.Error("Package Uptime() returned negative value")
	}

	cpus := Cpus()
	if len(cpus) == 0 {
		t.Error("Package Cpus() returned empty array")
	}

	totalmem := Totalmem()
	if totalmem == 0 {
		t.Error("Package Totalmem() returned zero")
	}

	_ = Freemem() // Just check it doesn't panic

	osType := Type()
	if osType == "" {
		t.Error("Package Type() returned empty string")
	}

	endianness := Endianness()
	if endianness != "LE" && endianness != "BE" {
		t.Errorf("Package Endianness() returned invalid value: %s", endianness)
	}
}

func TestDefault(t *testing.T) {
	// Test that Default instance works
	if Default == nil {
		t.Fatal("Default instance should not be nil")
	}

	hostname, err := Default.Hostname()
	if err != nil {
		t.Errorf("Default.Hostname() failed: %v", err)
	}
	if hostname == "" {
		t.Error("Default.Hostname() returned empty string")
	}
}

// Test that the package doesn't conflict with standard 'os' package
func TestNoConflict(t *testing.T) {
	// This test ensures we can still use standard os functions
	// by importing as 'stdos' or using qualified names

	// Our OS type should work
	o := &OS{}
	platform := o.Platform()

	// Check it matches runtime.GOOS
	if platform != runtime.GOOS {
		t.Errorf("Platform() should return %s, got %s", runtime.GOOS, platform)
	}
}
