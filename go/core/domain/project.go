package domain

import "time"

// Project represents a TypeScript to Go transpilation project
type Project struct {
	Path         string    `json:"path"`
	Name         string    `json:"name"`
	GoModuleName string    `json:"goModuleName"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Files        []File    `json:"files"`
}

// NewProject creates a new Project instance
func NewProject(path, name, goModuleName string) *Project {
	now := time.Now()
	return &Project{
		Path:         path,
		Name:         name,
		GoModuleName: goModuleName,
		CreatedAt:    now,
		UpdatedAt:    now,
		Files:        make([]File, 0),
	}
}

// AddFile adds a file to the project
func (p *Project) AddFile(file File) {
	p.Files = append(p.Files, file)
	p.UpdatedAt = time.Now()
}

// GetFileByPath returns a file by its path
func (p *Project) GetFileByPath(path string) *File {
	for i := range p.Files {
		if p.Files[i].Path == path {
			return &p.Files[i]
		}
	}
	return nil
}

// TypeScriptFileCount returns the number of TypeScript files in the project
func (p *Project) TypeScriptFileCount() int {
	count := 0
	for _, file := range p.Files {
		if file.IsTypeScript() {
			count++
		}
	}
	return count
}
