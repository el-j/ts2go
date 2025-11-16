package domain

import "time"

// ProjectState represents the persistent state of a transpilation project
type ProjectState struct {
	ProjectPath          string                      `json:"projectPath"`
	LastTranspilation    *TranspilationResult        `json:"lastTranspilation,omitempty"`
	TranspilationHistory []TranspilationHistoryEntry `json:"transpilationHistory"`
	FileStates           map[string]FileState        `json:"fileStates"`
	CreatedAt            time.Time                   `json:"createdAt"`
	UpdatedAt            time.Time                   `json:"updatedAt"`
}

// TranspilationHistoryEntry represents a historical transpilation
type TranspilationHistoryEntry struct {
	Timestamp      time.Time     `json:"timestamp"`
	Success        bool          `json:"success"`
	FilesProcessed int           `json:"filesProcessed"`
	Duration       time.Duration `json:"duration"`
}

// FileState represents the state of a single file
type FileState struct {
	SourcePath     string    `json:"sourcePath"`
	TargetPath     string    `json:"targetPath"`
	LastTranspiled time.Time `json:"lastTranspiled"`
	Hash           string    `json:"hash"` // For detecting changes
	Success        bool      `json:"success"`
}

// Settings represents application-wide settings
type Settings struct {
	// Application settings
	Theme            string `json:"theme"` // "light", "dark", "system"
	FontSize         int    `json:"fontSize"`
	AutoSave         bool   `json:"autoSave"`
	AutoSaveDelay    int    `json:"autoSaveDelay"` // milliseconds
	EnableBackups    bool   `json:"enableBackups"`
	BackupLocation   string `json:"backupLocation"`
	DefaultOutputDir string `json:"defaultOutputDir"`

	// Project settings
	GoModuleName    string `json:"goModuleName"`
	ExcludePatterns string `json:"excludePatterns"`
	IncludePatterns string `json:"includePatterns"`

	// Editor settings
	TabSize          int  `json:"tabSize"`
	WordWrap         bool `json:"wordWrap"`
	LineNumbers      bool `json:"lineNumbers"`
	Minimap          bool `json:"minimap"`
	AutoFormatOnSave bool `json:"autoFormatOnSave"`

	// Go configuration
	GoBinarySource     string `json:"goBinarySource"` // "bundled", "system", "custom"
	CustomGoBinaryPath string `json:"customGoBinaryPath"`
}

// NewProjectState creates a new ProjectState
func NewProjectState(projectPath string) *ProjectState {
	now := time.Now()
	return &ProjectState{
		ProjectPath:          projectPath,
		TranspilationHistory: make([]TranspilationHistoryEntry, 0),
		FileStates:           make(map[string]FileState),
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

// UpdateFromResult updates the state from a transpilation result
func (s *ProjectState) UpdateFromResult(result *TranspilationResult) {
	s.LastTranspilation = result
	s.UpdatedAt = time.Now()

	// Add to history
	entry := TranspilationHistoryEntry{
		Timestamp:      result.CompletedAt,
		Success:        result.Success,
		FilesProcessed: result.FilesProcessed,
		Duration:       result.Duration,
	}
	s.TranspilationHistory = append(s.TranspilationHistory, entry)

	// Keep only last 50 entries
	if len(s.TranspilationHistory) > 50 {
		s.TranspilationHistory = s.TranspilationHistory[len(s.TranspilationHistory)-50:]
	}

	// Update file states
	for _, fileResult := range result.FileResults {
		s.FileStates[fileResult.SourcePath] = FileState{
			SourcePath:     fileResult.SourcePath,
			TargetPath:     fileResult.TargetPath,
			LastTranspiled: time.Now(),
			Success:        fileResult.Success,
		}
	}
}

// DefaultSettings returns default application settings
func DefaultSettings() *Settings {
	return &Settings{
		Theme:              "system",
		FontSize:           14,
		AutoSave:           true,
		AutoSaveDelay:      3000,
		EnableBackups:      true,
		BackupLocation:     "./.backups",
		DefaultOutputDir:   "./output",
		GoModuleName:       "",
		ExcludePatterns:    "node_modules, **/*.test.ts, dist",
		IncludePatterns:    "**/*.ts, **/*.tsx",
		TabSize:            4,
		WordWrap:           true,
		LineNumbers:        true,
		Minimap:            true,
		AutoFormatOnSave:   true,
		GoBinarySource:     "system",
		CustomGoBinaryPath: "",
	}
}

// NewDefaultSettings is an alias for DefaultSettings
func NewDefaultSettings() *Settings {
	return DefaultSettings()
}
