package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/yourusername/ts2go/internal/transpiler"
)

// Watcher monitors TypeScript files and automatically transpiles them
type Watcher struct {
	watcher      *fsnotify.Watcher
	projectDir   string
	outputDir    string
	progress     *ProgressReporter
	fileMap      map[string]string // maps input file to output file
	debounceTime time.Duration
	mutex        sync.Mutex
	pending      map[string]time.Time
	stopChan     chan bool
}

// NewWatcher creates a new file watcher
func NewWatcher(projectDir, outputDir string, progress *ProgressReporter) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &Watcher{
		watcher:      watcher,
		projectDir:   projectDir,
		outputDir:    outputDir,
		progress:     progress,
		fileMap:      make(map[string]string),
		debounceTime: 500 * time.Millisecond,
		pending:      make(map[string]time.Time),
		stopChan:     make(chan bool),
	}, nil
}

// AddFile adds a TypeScript file to watch
func (w *Watcher) AddFile(inputPath, outputPath string) error {
	w.fileMap[inputPath] = outputPath

	// Watch the file's directory
	dir := filepath.Dir(inputPath)
	err := w.watcher.Add(dir)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to watch directory %s: %w", dir, err)
	}

	return nil
}

// Start begins watching for file changes
func (w *Watcher) Start() error {
	w.progress.Info("👀 Watching for changes... (Press Ctrl+C to stop)")

	// Start debounce ticker
	ticker := time.NewTicker(w.debounceTime)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return nil
			}

			// Only handle write and create events for .ts files
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				if strings.HasSuffix(event.Name, ".ts") {
					w.scheduleTranspile(event.Name)
				}
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return nil
			}
			w.progress.Warning(fmt.Sprintf("Watch error: %v", err))

		case <-ticker.C:
			// Process pending transpilations
			w.processPending()

		case <-w.stopChan:
			return nil
		}
	}
}

// Stop stops the watcher
func (w *Watcher) Stop() error {
	close(w.stopChan)
	return w.watcher.Close()
}

// scheduleTranspile schedules a file for transpilation with debouncing
func (w *Watcher) scheduleTranspile(filePath string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Check if this file is in our watch list
	if _, ok := w.fileMap[filePath]; !ok {
		return
	}

	// Schedule for processing
	w.pending[filePath] = time.Now()
}

// processPending processes all pending transpilations
func (w *Watcher) processPending() {
	w.mutex.Lock()
	toProcess := make([]string, 0, len(w.pending))
	now := time.Now()

	for filePath, scheduledTime := range w.pending {
		if now.Sub(scheduledTime) >= w.debounceTime {
			toProcess = append(toProcess, filePath)
			delete(w.pending, filePath)
		}
	}
	w.mutex.Unlock()

	// Process files outside the lock
	for _, filePath := range toProcess {
		w.transpileFile(filePath)
	}
}

// transpileFile transpiles a single file
func (w *Watcher) transpileFile(inputPath string) {
	outputPath, ok := w.fileMap[inputPath]
	if !ok {
		return
	}

	relPath, err := filepath.Rel(w.projectDir, inputPath)
	if err != nil {
		relPath = inputPath
	}

	w.progress.Verbose("🔄 Transpiling %s", relPath)

	err = transpiler.Transpile(inputPath, outputPath)
	if err != nil {
		w.progress.Error(fmt.Sprintf("Failed to transpile %s: %v", relPath, err))
	} else {
		w.progress.Success(fmt.Sprintf("Transpiled %s", relPath))
	}
}

// WatchProject watches an entire project directory
func WatchProject(projectDir, outputDir string, progressLevel ProgressLevel) error {
	progress := NewProgressReporter(progressLevel)

	// Create absolute paths
	projectDir, err := filepath.Abs(projectDir)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}

	outputDir, err = filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	progress.Start("Setting up watch mode", 2)

	// Create watcher
	watcher, err := NewWatcher(projectDir, outputDir, progress)
	if err != nil {
		return err
	}
	defer watcher.Stop()

	progress.Step("Scanning for TypeScript files...")

	// Find all .ts files
	tsFiles := make([]string, 0)
	err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip node_modules and hidden directories
		if info.IsDir() {
			name := info.Name()
			if name == "node_modules" || name == ".git" || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Only watch .ts files (not .d.ts)
		if strings.HasSuffix(path, ".ts") && !strings.HasSuffix(path, ".d.ts") {
			tsFiles = append(tsFiles, path)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scan project: %w", err)
	}

	progress.Success(fmt.Sprintf("Found %d TypeScript files", len(tsFiles)))
	progress.Step("Setting up file watchers...")

	// Add all files to watcher
	for _, inputPath := range tsFiles {
		relPath, err := filepath.Rel(projectDir, inputPath)
		if err != nil {
			continue
		}

		outputPath := filepath.Join(outputDir, relPath)
		outputPath = outputPath[:len(outputPath)-3] + ".go"

		// Create output directory
		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			progress.Warning(fmt.Sprintf("Failed to create directory for %s", relPath))
			continue
		}

		if err := watcher.AddFile(inputPath, outputPath); err != nil {
			progress.Warning(fmt.Sprintf("Failed to watch %s: %v", relPath, err))
		}
	}

	progress.Success(fmt.Sprintf("Watching %d files", len(tsFiles)))

	// Initial transpilation
	if progressLevel >= ProgressNormal {
		fmt.Println("\nPerforming initial transpilation...")
	}

	for _, inputPath := range tsFiles {
		outputPath := watcher.fileMap[inputPath]
		if err := transpiler.Transpile(inputPath, outputPath); err != nil {
			relPath, _ := filepath.Rel(projectDir, inputPath)
			progress.Verbose("⚠️  %s: %v", relPath, err)
		}
	}

	if progressLevel >= ProgressNormal {
		fmt.Println("✓ Initial transpilation complete")
	}

	// Start watching
	return watcher.Start()
}
