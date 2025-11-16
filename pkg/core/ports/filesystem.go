package ports

import "github.com/el-j/ts2go/pkg/core/domain"

// FileSystem is what the core uses to access files.
// This is a DRIVEN PORT (outbound) - a dependency the core needs.
// Implementations will use os.* functions, but the core doesn't know that.
type FileSystem interface {
	// ReadFile reads the entire content of a file
	ReadFile(path string) ([]byte, error)

	// WriteFile writes data to a file, creating it if it doesn't exist
	WriteFile(path string, data []byte) error

	// ScanDirectory recursively scans a directory and returns all files
	// It respects include/exclude patterns if provided
	ScanDirectory(path string, includePatterns, excludePatterns []string) ([]domain.File, error)

	// DirectoryExists checks if a directory exists
	DirectoryExists(path string) (bool, error)

	// FileExists checks if a file exists
	FileExists(path string) (bool, error)

	// PathExists checks if a path exists (file or directory)
	PathExists(path string) (bool, error)

	// CreateDirectory creates a directory and all parent directories
	CreateDirectory(path string) error

	// DeleteFile deletes a file
	DeleteFile(path string) error

	// CopyFile copies a file from src to dst
	CopyFile(src, dst string) error

	// GetFileInfo returns file information (size, mod time, etc.)
	GetFileInfo(path string) (FileInfo, error)
}

// FileInfo represents file metadata
type FileInfo struct {
	Name    string
	Size    int64
	ModTime int64 // Unix timestamp
	IsDir   bool
}
