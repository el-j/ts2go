package domain_test

import (
	"testing"
	"time"

	"github.com/el-j/ts2go/core/domain"
)

func TestNewProject(t *testing.T) {
	path := "/path/to/project"
	name := "myproject"
	moduleName := "github.com/user/myproject"

	project := domain.NewProject(path, name, moduleName)

	if project.Path != path {
		t.Errorf("Expected path %s, got %s", path, project.Path)
	}
	if project.Name != name {
		t.Errorf("Expected name %s, got %s", name, project.Name)
	}
	if project.GoModuleName != moduleName {
		t.Errorf("Expected module name %s, got %s", moduleName, project.GoModuleName)
	}
	if project.Files == nil {
		t.Error("Expected Files to be initialized, got nil")
	}
	if len(project.Files) != 0 {
		t.Errorf("Expected 0 files, got %d", len(project.Files))
	}
	if project.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if project.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestProjectAddFile(t *testing.T) {
	project := domain.NewProject("/path", "test", "test/module")
	originalUpdatedAt := project.UpdatedAt

	// Sleep briefly to ensure timestamp changes
	time.Sleep(10 * time.Millisecond)

	file := domain.NewFile("/path/test.ts", "const x = 1;")
	project.AddFile(*file)

	if len(project.Files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(project.Files))
	}
	if project.Files[0].Path != file.Path {
		t.Errorf("Expected file path %s, got %s", file.Path, project.Files[0].Path)
	}
	if !project.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated after adding file")
	}
}

func TestProjectGetFileByPath(t *testing.T) {
	project := domain.NewProject("/path", "test", "test/module")

	file1 := domain.NewFile("/path/test1.ts", "const x = 1;")
	file2 := domain.NewFile("/path/test2.ts", "const y = 2;")

	project.AddFile(*file1)
	project.AddFile(*file2)

	// Test finding existing file
	found := project.GetFileByPath("/path/test1.ts")
	if found == nil {
		t.Error("Expected to find file, got nil")
	} else if found.Path != file1.Path {
		t.Errorf("Expected path %s, got %s", file1.Path, found.Path)
	}

	// Test non-existent file
	notFound := project.GetFileByPath("/path/nonexistent.ts")
	if notFound != nil {
		t.Error("Expected nil for non-existent file, got a file")
	}
}

func TestProjectTypeScriptFileCount(t *testing.T) {
	project := domain.NewProject("/path", "test", "test/module")

	tsFile1 := domain.NewFile("/path/test1.ts", "const x = 1;")
	tsFile2 := domain.NewFile("/path/test2.tsx", "const y = 2;")
	goFile := domain.NewFile("/path/main.go", "package main")
	jsFile := domain.NewFile("/path/script.js", "var z = 3;")

	project.AddFile(*tsFile1)
	project.AddFile(*tsFile2)
	project.AddFile(*goFile)
	project.AddFile(*jsFile)

	count := project.TypeScriptFileCount()
	if count != 2 {
		t.Errorf("Expected 2 TypeScript files, got %d", count)
	}
}

func TestProjectValidate(t *testing.T) {
	// Test basic project creation - models are simple POCOs without validation
	project := domain.NewProject("/valid/path", "validname", "github.com/user/repo")

	if project.Path != "/valid/path" {
		t.Error("Project path not set correctly")
	}
	if project.Name != "validname" {
		t.Error("Project name not set correctly")
	}
	if project.GoModuleName != "github.com/user/repo" {
		t.Error("Project module name not set correctly")
	}
}
