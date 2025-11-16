package mocks

import (
	"fmt"

	"github.com/el-j/ts2go/pkg/core/domain"
	"github.com/el-j/ts2go/pkg/core/ports"
)

// MockFileSystem is a mock implementation of ports.FileSystem for testing
type MockFileSystem struct {
	// Files maps file paths to their content
	Files map[string][]byte

	// Directories tracks which paths are directories
	Directories map[string]bool

	// FileInfo tracks metadata for files
	FileInfo map[string]ports.FileInfo

	// ShouldFailRead causes ReadFile to fail for specific paths
	ShouldFailRead map[string]error

	// ShouldFailWrite causes WriteFile to fail for specific paths
	ShouldFailWrite map[string]error

	// ScanResult is the result to return from ScanDirectory
	ScanResult []domain.File

	// ScanError is the error to return from ScanDirectory
	ScanError error
}

// NewMockFileSystem creates a new mock filesystem
func NewMockFileSystem() *MockFileSystem {
	return &MockFileSystem{
		Files:           make(map[string][]byte),
		Directories:     make(map[string]bool),
		FileInfo:        make(map[string]ports.FileInfo),
		ShouldFailRead:  make(map[string]error),
		ShouldFailWrite: make(map[string]error),
		ScanResult:      make([]domain.File, 0),
	}
}

// ReadFile reads a file from the mock filesystem
func (m *MockFileSystem) ReadFile(path string) ([]byte, error) {
	if err, exists := m.ShouldFailRead[path]; exists {
		return nil, err
	}

	content, exists := m.Files[path]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	return content, nil
}

// WriteFile writes a file to the mock filesystem
func (m *MockFileSystem) WriteFile(path string, data []byte) error {
	if err, exists := m.ShouldFailWrite[path]; exists {
		return err
	}

	m.Files[path] = data
	return nil
}

// ScanDirectory returns the configured scan result
func (m *MockFileSystem) ScanDirectory(path string, includePatterns, excludePatterns []string) ([]domain.File, error) {
	if m.ScanError != nil {
		return nil, m.ScanError
	}
	return m.ScanResult, nil
}

// DirectoryExists checks if a directory exists in the mock
func (m *MockFileSystem) DirectoryExists(path string) (bool, error) {
	return m.Directories[path], nil
}

// FileExists checks if a file exists in the mock
func (m *MockFileSystem) FileExists(path string) (bool, error) {
	_, exists := m.Files[path]
	return exists, nil
}

// PathExists checks if a path exists (file or directory)
func (m *MockFileSystem) PathExists(path string) (bool, error) {
	if m.Directories[path] {
		return true, nil
	}
	_, exists := m.Files[path]
	return exists, nil
}

// CreateDirectory creates a directory in the mock
func (m *MockFileSystem) CreateDirectory(path string) error {
	m.Directories[path] = true
	return nil
}

// DeleteFile deletes a file from the mock
func (m *MockFileSystem) DeleteFile(path string) error {
	delete(m.Files, path)
	return nil
}

// CopyFile copies a file in the mock
func (m *MockFileSystem) CopyFile(src, dst string) error {
	content, exists := m.Files[src]
	if !exists {
		return fmt.Errorf("source file not found: %s", src)
	}
	m.Files[dst] = content
	return nil
}

// GetFileInfo returns file info from the mock
func (m *MockFileSystem) GetFileInfo(path string) (ports.FileInfo, error) {
	info, exists := m.FileInfo[path]
	if !exists {
		return ports.FileInfo{}, fmt.Errorf("file not found: %s", path)
	}
	return info, nil
}

// AddFile adds a file to the mock filesystem (helper for tests)
func (m *MockFileSystem) AddFile(path string, content string) {
	m.Files[path] = []byte(content)
}

// AddDirectory adds a directory to the mock filesystem (helper for tests)
func (m *MockFileSystem) AddDirectory(path string) {
	m.Directories[path] = true
}

// SetScanResult sets the result for ScanDirectory (helper for tests)
func (m *MockFileSystem) SetScanResult(files []domain.File) {
	m.ScanResult = files
}
