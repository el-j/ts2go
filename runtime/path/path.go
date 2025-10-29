package path

import (
	"path/filepath"
)

// Join joins path segments
func Join(parts ...string) string {
	return filepath.Join(parts...)
}

// Dirname returns the directory name of a path
func Dirname(path string) string {
	return filepath.Dir(path)
}

// Basename returns the base name of a path
func Basename(path string) string {
	return filepath.Base(path)
}

// Extname returns the extension of a path
func Extname(path string) string {
	return filepath.Ext(path)
}
