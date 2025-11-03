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