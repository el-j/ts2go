package fs

import (
	"path/filepath"
	"testing"
)

func TestFS(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	subDir := filepath.Join(tempDir, "nested", "dir")

	// Test Mkdir
	if err := Mkdir(subDir); err != nil {
		t.Fatalf("Mkdir failed: %v", err)
	}

	// Test WriteFile
	content := "Hello TS2Go File System"
	if err := WriteFile(testFile, content); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Test Exists
	if !Exists(testFile) {
		t.Errorf("expected %s to exist", testFile)
	}
	if Exists(filepath.Join(tempDir, "nonexistent.txt")) {
		t.Errorf("expected nonexistent.txt to not exist")
	}

	// Test ReadFile
	readBack, err := ReadFile(testFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if readBack != content {
		t.Errorf("expected %q, got %q", content, readBack)
	}

	// Test ReadFile error
	_, err = ReadFile(filepath.Join(tempDir, "nonexistent.txt"))
	if err == nil {
		t.Errorf("expected error reading nonexistent file")
	}

	// Test Stat
	info, err := Stat(testFile)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Size() != int64(len(content)) {
		t.Errorf("expected size %d, got %d", len(content), info.Size())
	}

	// Test Unlink
	if err := Unlink(testFile); err != nil {
		t.Fatalf("Unlink failed: %v", err)
	}
	if Exists(testFile) {
		t.Errorf("expected file to be unlinked")
	}
}
