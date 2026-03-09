package filesystem_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/el-j/ts2go/adapters/driven/filesystem"
)

func TestOSFileSystem(t *testing.T) {
	// Create sandbox temp dir
	tempDir, err := os.MkdirTemp("", "ts2go-fs-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // cleanup

	fs := filesystem.NewOSFileSystem()

	t.Run("WriteFile and ReadFile", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test-write.txt")
		data := []byte("hello world")

		err := fs.WriteFile(filePath, data)
		if err != nil {
			t.Fatalf("WriteFile failed: %v", err)
		}

		readData, err := fs.ReadFile(filePath)
		if err != nil {
			t.Fatalf("ReadFile failed: %v", err)
		}

		if string(readData) != string(data) {
			t.Errorf("Expected %s, got %s", string(data), string(readData))
		}

		// Error path
		_, err = fs.ReadFile(filepath.Join(tempDir, "missing.txt"))
		if err == nil {
			t.Errorf("ReadFile on missing file should error")
		}
	})

	t.Run("CreateDirectory and DirectoryExists", func(t *testing.T) {
		dirPath := filepath.Join(tempDir, "nested", "dir")

		err := fs.CreateDirectory(dirPath)
		if err != nil {
			t.Fatalf("CreateDirectory failed: %v", err)
		}

		exists, err := fs.DirectoryExists(dirPath)
		if err != nil || !exists {
			t.Errorf("DirectoryExists failed or returned false: %v", err)
		}
	})

	t.Run("FileExists", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "exist-test.txt")
		fs.WriteFile(filePath, []byte("data"))

		exists, err := fs.FileExists(filePath)
		if err != nil || !exists {
			t.Errorf("FileExists failed or returned false: %v", err)
		}

		exists, _ = fs.FileExists(filepath.Join(tempDir, "non-existent.txt"))
		if exists {
			t.Errorf("FileExists should return false for non-existent file")
		}
	})

	t.Run("PathExists", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "path-exist-test.txt")
		fs.WriteFile(filePath, []byte("data"))

		exists, err := fs.PathExists(filePath)
		if err != nil || !exists {
			t.Errorf("PathExists failed or returned false for file: %v", err)
		}

		exists, err = fs.PathExists(tempDir)
		if err != nil || !exists {
			t.Errorf("PathExists failed or returned false for dir: %v", err)
		}
	})

	t.Run("DeleteFile", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "delete-test.txt")
		fs.WriteFile(filePath, []byte("data"))

		err := fs.DeleteFile(filePath)
		if err != nil {
			t.Fatalf("DeleteFile failed: %v", err)
		}

		exists, _ := fs.FileExists(filePath)
		if exists {
			t.Errorf("File still exists after DeleteFile")
		}
	})

	t.Run("CopyFile", func(t *testing.T) {
		src := filepath.Join(tempDir, "src.txt")
		dst := filepath.Join(tempDir, "dst.txt")

		fs.WriteFile(src, []byte("copied content"))
		err := fs.CopyFile(src, dst)
		if err != nil {
			t.Fatalf("CopyFile failed: %v", err)
		}

		data, _ := fs.ReadFile(dst)
		if string(data) != "copied content" {
			t.Errorf("CopyFile content mismatch")
		}

		// Error path
		err = fs.CopyFile(filepath.Join(tempDir, "missing_src.txt"), dst)
		if err == nil {
			t.Errorf("CopyFile on missing source should error")
		}
	})

	t.Run("GetFileInfo", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "info.txt")
		fs.WriteFile(filePath, []byte("12345"))

		info, err := fs.GetFileInfo(filePath)
		if err != nil {
			t.Fatalf("GetFileInfo failed: %v", err)
		}

		if info.Name != "info.txt" {
			t.Errorf("Expected name %s, got %s", "info.txt", info.Name)
		}
		if info.Size != 5 {
			t.Errorf("Expected size 5, got %d", info.Size)
		}
		if info.IsDir {
			t.Errorf("Expected IsDir false")
		}
		if info.ModTime > time.Now().Unix()+5 {
			t.Errorf("ModTime seems incorrect")
		}
	})

	t.Run("ScanDirectory", func(t *testing.T) {
		scanDir := filepath.Join(tempDir, "scan")
		fs.CreateDirectory(scanDir)
		fs.WriteFile(filepath.Join(scanDir, "file1.ts"), []byte("data"))
		fs.WriteFile(filepath.Join(scanDir, "file2.go"), []byte("data"))
		fs.WriteFile(filepath.Join(scanDir, "ignore.txt"), []byte("data"))
		fs.CreateDirectory(filepath.Join(scanDir, "node_modules"))
		fs.WriteFile(filepath.Join(scanDir, "node_modules", "module.ts"), []byte("data"))

		// Test basic scan with include
		files, err := fs.ScanDirectory(scanDir, []string{"*.ts"}, nil)
		if err != nil {
			t.Fatalf("ScanDirectory failed: %v", err)
		}

		hasTSFile := false
		for _, f := range files {
			if f.Name == "file1.ts" || f.Name == "module.ts" {
				hasTSFile = true
			}
			if f.Name == "file2.go" || f.Name == "ignore.txt" {
				t.Errorf("ScanDirectory included unwanted file: %s", f.Name)
			}
		}
		if !hasTSFile {
			t.Errorf("ScanDirectory missed included file")
		}

		// Test with exclude
		files, err = fs.ScanDirectory(scanDir, []string{"*.ts"}, []string{"node_modules"})
		if err != nil {
			t.Fatalf("ScanDirectory with exclude failed: %v", err)
		}

		for _, f := range files {
			if f.Name == "module.ts" {
				t.Errorf("ScanDirectory did not exclude node_modules")
			}
		}

		// Additional edge cases
		fs.ScanDirectory(scanDir, []string{"**/*.go"}, []string{"**/*.txt"})
		fs.ScanDirectory(scanDir, []string{""}, []string{""})
	})
}
