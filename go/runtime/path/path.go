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

// IsAbsolute returns whether path is absolute
func IsAbsolute(path string) bool {
	return filepath.IsAbs(path)
}

// Normalize cleans and normalizes a path
func Normalize(path string) string {
	return filepath.Clean(path)
}

// Resolve resolves a sequence of paths or path segments into an absolute path
func Resolve(parts ...string) string {
	joined := filepath.Join(parts...)
	if filepath.IsAbs(joined) {
		return filepath.Clean(joined)
	}
	abs, err := filepath.Abs(joined)
	if err != nil {
		return joined
	}
	return abs
}
