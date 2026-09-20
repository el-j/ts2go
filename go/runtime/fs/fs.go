package fs

import (
	"os"
)

// ReadFile reads a file and returns its contents as string
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes data to a file with default permissions
func WriteFile(path, data string) error {
	return os.WriteFile(path, []byte(data), 0644)
}

// Exists checks if a file or directory exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Unlink deletes a file (equivalent to fs.unlink)
func Unlink(path string) error {
	return os.Remove(path)
}

// Mkdir creates a directory including any necessary parents
func Mkdir(path string) error {
	return os.MkdirAll(path, 0755)
}

// Stat returns file information
func Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
