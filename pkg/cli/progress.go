package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// ProgressLevel defines the verbosity level for progress reporting
type ProgressLevel int

const (
	// ProgressQuiet shows no progress output
	ProgressQuiet ProgressLevel = iota
	// ProgressNormal shows standard progress output
	ProgressNormal
	// ProgressVerbose shows detailed progress output
	ProgressVerbose
)

// ProgressReporter handles progress reporting for CLI operations
type ProgressReporter struct {
	level       ProgressLevel
	writer      io.Writer
	startTime   time.Time
	totalSteps  int
	currentStep int
}

// NewProgressReporter creates a new progress reporter
func NewProgressReporter(level ProgressLevel) *ProgressReporter {
	return &ProgressReporter{
		level:     level,
		writer:    os.Stdout,
		startTime: time.Now(),
	}
}

// SetLevel changes the progress level
func (p *ProgressReporter) SetLevel(level ProgressLevel) {
	p.level = level
}

// Start begins a new operation with the given name
func (p *ProgressReporter) Start(operation string, totalSteps int) {
	p.startTime = time.Now()
	p.totalSteps = totalSteps
	p.currentStep = 0

	if p.level >= ProgressNormal {
		fmt.Fprintf(p.writer, "\n%s\n", operation)
		fmt.Fprintf(p.writer, "%s\n", strings.Repeat("=", len(operation)))
	}
}

// Step reports progress for a single step
func (p *ProgressReporter) Step(description string) {
	p.currentStep++

	if p.level >= ProgressNormal {
		percentage := float64(p.currentStep) / float64(p.totalSteps) * 100
		fmt.Fprintf(p.writer, "[%d/%d] (%.0f%%) %s\n", p.currentStep, p.totalSteps, percentage, description)
	}
}

// StepVerbose reports progress with verbose output
func (p *ProgressReporter) StepVerbose(description string, details string) {
	p.currentStep++

	if p.level >= ProgressVerbose {
		percentage := float64(p.currentStep) / float64(p.totalSteps) * 100
		fmt.Fprintf(p.writer, "[%d/%d] (%.0f%%) %s\n", p.currentStep, p.totalSteps, percentage, description)
		if details != "" {
			fmt.Fprintf(p.writer, "    %s\n", details)
		}
	} else if p.level >= ProgressNormal {
		// Show progress every 10 steps in normal mode
		if p.currentStep%10 == 0 || p.currentStep == p.totalSteps {
			percentage := float64(p.currentStep) / float64(p.totalSteps) * 100
			fmt.Fprintf(p.writer, "Progress: %d/%d (%.0f%%)\n", p.currentStep, p.totalSteps, percentage)
		}
	}
}

// Success reports a successful operation
func (p *ProgressReporter) Success(description string) {
	if p.level >= ProgressNormal {
		fmt.Fprintf(p.writer, "  ✓ %s\n", description)
	}
}

// Warning reports a warning
func (p *ProgressReporter) Warning(description string) {
	if p.level >= ProgressNormal {
		fmt.Fprintf(p.writer, "  ⚠️  %s\n", description)
	}
}

// Error reports an error
func (p *ProgressReporter) Error(description string) {
	if p.level >= ProgressQuiet {
		fmt.Fprintf(p.writer, "  ✗ %s\n", description)
	}
}

// Info reports informational output
func (p *ProgressReporter) Info(format string, args ...interface{}) {
	if p.level >= ProgressNormal {
		fmt.Fprintf(p.writer, "  ")
		fmt.Fprintf(p.writer, format, args...)
		fmt.Fprintf(p.writer, "\n")
	}
}

// Verbose reports verbose output
func (p *ProgressReporter) Verbose(format string, args ...interface{}) {
	if p.level >= ProgressVerbose {
		fmt.Fprintf(p.writer, "    ")
		fmt.Fprintf(p.writer, format, args...)
		fmt.Fprintf(p.writer, "\n")
	}
}

// Complete marks the operation as complete and shows timing
func (p *ProgressReporter) Complete(message string) {
	if p.level >= ProgressNormal {
		elapsed := time.Since(p.startTime)
		fmt.Fprintf(p.writer, "\n%s\n", strings.Repeat("=", 50))
		fmt.Fprintf(p.writer, "%s\n", message)
		fmt.Fprintf(p.writer, "Completed in %v\n", elapsed.Round(time.Millisecond))
		fmt.Fprintf(p.writer, "%s\n\n", strings.Repeat("=", 50))
	}
}

// ProgressBar displays a visual progress bar
type ProgressBar struct {
	total   int
	current int
	width   int
	writer  io.Writer
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int, width int) *ProgressBar {
	return &ProgressBar{
		total:  total,
		width:  width,
		writer: os.Stdout,
	}
}

// Update updates the progress bar
func (pb *ProgressBar) Update(current int) {
	pb.current = current
	pb.Render()
}

// Increment increments the progress bar
func (pb *ProgressBar) Increment() {
	pb.current++
	pb.Render()
}

// Render renders the progress bar
func (pb *ProgressBar) Render() {
	percentage := float64(pb.current) / float64(pb.total)
	filled := int(percentage * float64(pb.width))

	bar := strings.Repeat("█", filled) + strings.Repeat("░", pb.width-filled)

	fmt.Fprintf(pb.writer, "\r[%s] %d/%d (%.0f%%)", bar, pb.current, pb.total, percentage*100)

	if pb.current >= pb.total {
		fmt.Fprintf(pb.writer, "\n")
	}
}

// Clear clears the progress bar from the terminal
func (pb *ProgressBar) Clear() {
	fmt.Fprintf(pb.writer, "\r%s\r", strings.Repeat(" ", pb.width+20))
}

// Spinner displays a simple spinner animation
type Spinner struct {
	frames []string
	index  int
	writer io.Writer
	active bool
}

// NewSpinner creates a new spinner
func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		writer: os.Stdout,
	}
}

// Start starts the spinner with the given message
func (s *Spinner) Start(message string) {
	s.active = true
	s.index = 0
	fmt.Fprintf(s.writer, "%s %s", s.frames[s.index], message)
}

// Tick advances the spinner animation
func (s *Spinner) Tick() {
	if !s.active {
		return
	}
	s.index = (s.index + 1) % len(s.frames)
	fmt.Fprintf(s.writer, "\r%s", s.frames[s.index])
}

// Stop stops the spinner and clears it
func (s *Spinner) Stop() {
	if !s.active {
		return
	}
	s.active = false
	fmt.Fprintf(s.writer, "\r%s\r", strings.Repeat(" ", 50))
}
