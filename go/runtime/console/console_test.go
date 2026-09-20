package console

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

func TestLog(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	Log("hello", "world", 42)

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, "hello world 42") {
		t.Errorf("expected 'hello world 42', got %q", output)
	}
}

func TestError(t *testing.T) {
	var buf bytes.Buffer
	oldOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(oldOutput)

	Error("failed operation", 500)

	output := buf.String()
	if !strings.Contains(output, "failed operation 500") {
		t.Errorf("expected 'failed operation 500', got %q", output)
	}
}

func TestWarn(t *testing.T) {
	var buf bytes.Buffer
	oldOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(oldOutput)

	Warn("deprecated feature")

	output := buf.String()
	if !strings.Contains(output, "deprecated feature") {
		t.Errorf("expected 'deprecated feature', got %q", output)
	}
}
