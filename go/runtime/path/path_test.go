package path

import (
	"path/filepath"
	"testing"
)

func TestPath(t *testing.T) {
	// Join
	if got := Join("a", "b", "c"); got != filepath.Join("a", "b", "c") {
		t.Errorf("Join failed: got %q", got)
	}

	// Dirname
	if got := Dirname("/a/b/c.txt"); got != "/a/b" {
		t.Errorf("Dirname failed: got %q, expected /a/b", got)
	}

	// Basename
	if got := Basename("/a/b/c.txt"); got != "c.txt" {
		t.Errorf("Basename failed: got %q, expected c.txt", got)
	}

	// Extname
	if got := Extname("/a/b/c.txt"); got != ".txt" {
		t.Errorf("Extname failed: got %q, expected .txt", got)
	}

	// IsAbsolute
	if !IsAbsolute("/usr/local") {
		t.Errorf("expected /usr/local to be absolute")
	}
	if IsAbsolute("relative/path") {
		t.Errorf("expected relative/path to not be absolute")
	}

	// Normalize
	if got := Normalize("a/b/../c/./d"); got != filepath.Clean("a/b/../c/./d") {
		t.Errorf("Normalize failed: got %q", got)
	}

	// Resolve
	absExpected, _ := filepath.Abs("foo/bar")
	if got := Resolve("foo", "bar"); got != absExpected {
		t.Errorf("Resolve relative failed: got %q, expected %q", got, absExpected)
	}
	if got := Resolve("/absolute", "path"); got != filepath.Clean("/absolute/path") {
		t.Errorf("Resolve absolute failed: got %q", got)
	}
}
