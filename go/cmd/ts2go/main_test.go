package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestPrintVersion(t *testing.T) {
	out := captureStdout(printVersion)
	if !strings.Contains(out, "ts2go version") {
		t.Errorf("expected version output, got: %s", out)
	}
	if !strings.Contains(out, "TypeScript to Go Transpiler") {
		t.Errorf("expected header in version output, got: %s", out)
	}
}

func TestPrintUsage(t *testing.T) {
	out := captureStdout(printUsage)
	requiredKeywords := []string{
		"convert",
		"transpile",
		"analyze",
		"build",
		"test",
		"state",
		"settings",
		"ui",
		"help",
		"version",
	}
	for _, kw := range requiredKeywords {
		if !strings.Contains(out, kw) {
			t.Errorf("expected usage output to contain command %q", kw)
		}
	}
}

func TestConvertCommandErrors(t *testing.T) {
	// Missing arguments
	if err := convertCommand([]string{}); err == nil {
		t.Errorf("expected error for empty args")
	}

	// Missing out flag
	if err := convertCommand([]string{"--in", "foo.ts"}); err == nil {
		t.Errorf("expected error when --out is missing")
	}

	// Missing in flag
	if err := convertCommand([]string{"--out", "foo.go"}); err == nil {
		t.Errorf("expected error when --in is missing")
	}

	// Nonexistent input file
	if err := convertCommand([]string{"-i", "/nonexistent/input.ts", "-o", "/tmp/out.go"}); err == nil {
		t.Errorf("expected error for nonexistent input file")
	}
}
