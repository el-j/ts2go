package filesystem

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/el-j/ts2go/pkg/core/domain"
	"github.com/el-j/ts2go/pkg/core/ports"
)

// OSFileSystem implements the FileSystem port using the OS filesystem
type OSFileSystem struct{}

// NewOSFileSystem creates a new OS filesystem adapter
func NewOSFileSystem() ports.FileSystem {
	return &OSFileSystem{}
}

// ReadFile reads the entire content of a file
func (fs *OSFileSystem) ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WriteFile writes data to a file, creating it if it doesn't exist
func (fs *OSFileSystem) WriteFile(path string, data []byte) error {
	// Create parent directories if they don't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ScanDirectory recursively scans a directory and returns all files
func (fs *OSFileSystem) ScanDirectory(path string, includePatterns, excludePatterns []string) ([]domain.File, error) {
	var files []domain.File

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Check if directory should be excluded
			if shouldExclude(filePath, excludePatterns) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file matches patterns
		if shouldExclude(filePath, excludePatterns) {
			return nil
		}

		if len(includePatterns) > 0 && !shouldInclude(filePath, includePatterns) {
			return nil
		}

		// Create File domain entity
		file := domain.File{
			Path: filePath,
			Name: info.Name(),
		}

		files = append(files, file)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// DirectoryExists checks if a directory exists
func (fs *OSFileSystem) DirectoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

// FileExists checks if a file exists
func (fs *OSFileSystem) FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !info.IsDir(), nil
}

// PathExists checks if a path exists (file or directory)
func (fs *OSFileSystem) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateDirectory creates a directory and all parent directories
func (fs *OSFileSystem) CreateDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// DeleteFile deletes a file
func (fs *OSFileSystem) DeleteFile(path string) error {
	return os.Remove(path)
}

// CopyFile copies a file from src to dst
func (fs *OSFileSystem) CopyFile(src, dst string) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination directory if needed
	destDir := filepath.Dir(dst)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy content
	_, err = io.Copy(destFile, sourceFile)
	return err
}

// GetFileInfo returns file information
func (fs *OSFileSystem) GetFileInfo(path string) (ports.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return ports.FileInfo{}, err
	}

	return ports.FileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		ModTime: info.ModTime().Unix(),
		IsDir:   info.IsDir(),
	}, nil
}

// shouldExclude checks if a path should be excluded based on patterns
func shouldExclude(path string, patterns []string) bool {
	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		// Simple pattern matching
		if strings.Contains(path, pattern) {
			return true
		}

		// Glob pattern matching
		matched, _ := filepath.Match(pattern, filepath.Base(path))
		if matched {
			return true
		}
	}
	return false
}

// shouldInclude checks if a path should be included based on patterns
func shouldInclude(path string, patterns []string) bool {
	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		// Handle ** glob patterns
		if strings.Contains(pattern, "**") {
			ext := filepath.Ext(path)
			patternExt := filepath.Ext(pattern)
			if ext == patternExt {
				return true
			}
		}

		// Glob pattern matching
		matched, _ := filepath.Match(pattern, filepath.Base(path))
		if matched {
			return true
		}
	}
	return false
}
