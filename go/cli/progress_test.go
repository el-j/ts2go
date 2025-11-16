package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestProgressReporter_Normal(t *testing.T) {
	buf := &bytes.Buffer{}
	progress := NewProgressReporter(ProgressNormal)
	progress.writer = buf

	progress.Start("Test Operation", 3)
	progress.Step("Step 1")
	progress.Step("Step 2")
	progress.Step("Step 3")
	progress.Complete("Done")

	output := buf.String()

	// Check that output contains expected elements
	checks := []string{
		"Test Operation",
		"[1/3]",
		"[2/3]",
		"[3/3]",
		"Done",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Output should contain '%s'\nGot: %s", check, output)
		}
	}
}

func TestProgressReporter_Verbose(t *testing.T) {
	buf := &bytes.Buffer{}
	progress := NewProgressReporter(ProgressVerbose)
	progress.writer = buf

	progress.StepVerbose("Processing file", "details about file")

	output := buf.String()

	if !strings.Contains(output, "Processing file") {
		t.Errorf("Verbose output should contain step description")
	}

	if !strings.Contains(output, "details about file") {
		t.Errorf("Verbose output should contain details")
	}
}

func TestProgressReporter_Quiet(t *testing.T) {
	buf := &bytes.Buffer{}
	progress := NewProgressReporter(ProgressQuiet)
	progress.writer = buf

	progress.Start("Test", 1)
	progress.Step("Step")
	progress.Info("Info message")
	progress.Success("Success message")

	output := buf.String()

	// In quiet mode, these should not appear
	if strings.Contains(output, "Step") || strings.Contains(output, "Info") || strings.Contains(output, "Success") {
		t.Errorf("Quiet mode should suppress normal output, got: %s", output)
	}
}

func TestProgressReporter_ErrorsAlwaysShow(t *testing.T) {
	buf := &bytes.Buffer{}
	progress := NewProgressReporter(ProgressQuiet)
	progress.writer = buf

	progress.Error("Critical error")

	output := buf.String()

	// Errors should show even in quiet mode
	if !strings.Contains(output, "Critical error") {
		t.Errorf("Errors should show in quiet mode")
	}
}

func TestProgressReporter_SuccessWarningError(t *testing.T) {
	buf := &bytes.Buffer{}
	progress := NewProgressReporter(ProgressNormal)
	progress.writer = buf

	progress.Success("All good")
	progress.Warning("Be careful")
	progress.Error("Something failed")

	output := buf.String()

	checks := map[string]string{
		"success": "✓",
		"warning": "⚠️",
		"error":   "✗",
	}

	for name, symbol := range checks {
		if !strings.Contains(output, symbol) {
			t.Errorf("%s should show symbol %s", name, symbol)
		}
	}
}

func TestProgressBar_Basic(t *testing.T) {
	buf := &bytes.Buffer{}
	bar := NewProgressBar(10, 20)
	bar.writer = buf

	bar.Update(5)

	output := buf.String()

	// Should show progress
	if !strings.Contains(output, "5/10") {
		t.Errorf("Progress bar should show current/total")
	}

	if !strings.Contains(output, "50%") {
		t.Errorf("Progress bar should show percentage")
	}
}

func TestProgressBar_Increment(t *testing.T) {
	buf := &bytes.Buffer{}
	bar := NewProgressBar(3, 10)
	bar.writer = buf

	bar.Increment()
	bar.Increment()
	bar.Increment()

	output := buf.String()

	// Should reach 100%
	if !strings.Contains(output, "100%") {
		t.Errorf("Progress bar should reach 100%%")
	}
}

func TestProgressBar_FullWidth(t *testing.T) {
	buf := &bytes.Buffer{}
	bar := NewProgressBar(100, 10)
	bar.writer = buf

	bar.Update(100)

	output := buf.String()

	// Should show full bar (10 filled characters)
	if !strings.Contains(output, strings.Repeat("█", 10)) {
		t.Errorf("Progress bar should be fully filled")
	}
}

func TestProgressBar_EmptyBar(t *testing.T) {
	buf := &bytes.Buffer{}
	bar := NewProgressBar(100, 10)
	bar.writer = buf

	bar.Update(0)

	output := buf.String()

	// Should show empty bar (10 empty characters)
	if !strings.Contains(output, strings.Repeat("░", 10)) {
		t.Errorf("Progress bar should be empty")
	}
}

func TestSpinner_StartStop(t *testing.T) {
	buf := &bytes.Buffer{}
	spinner := NewSpinner()
	spinner.writer = buf

	spinner.Start("Loading")

	output := buf.String()

	// Should show initial frame and message
	if !strings.Contains(output, "Loading") {
		t.Errorf("Spinner should show message")
	}

	buf.Reset()
	spinner.Stop()

	// Stop should clear the spinner
	// (we just verify it doesn't panic)
}

func TestSpinner_Tick(t *testing.T) {
	buf := &bytes.Buffer{}
	spinner := NewSpinner()
	spinner.writer = buf

	spinner.Start("Loading")
	initialFrame := spinner.frames[0]

	buf.Reset()
	spinner.Tick()

	// After tick, should show next frame
	// (we just verify it produces output)
	if buf.Len() == 0 {
		t.Errorf("Tick should produce output")
	}

	// Verify frame changed
	if spinner.index == 0 {
		t.Errorf("Tick should advance frame index")
	}

	_ = initialFrame // Use the variable
}

func TestProgressReporter_SetLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	progress := NewProgressReporter(ProgressNormal)
	progress.writer = buf

	progress.Info("Normal message")
	normalOutput := buf.String()

	buf.Reset()
	progress.SetLevel(ProgressQuiet)
	progress.Info("Quiet message")
	quietOutput := buf.String()

	if normalOutput == "" {
		t.Errorf("Normal mode should show info")
	}

	if quietOutput != "" {
		t.Errorf("Quiet mode should suppress info")
	}
}
