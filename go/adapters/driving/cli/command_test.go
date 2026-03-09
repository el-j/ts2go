package cli_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func skipIfGoNotInstalled(t *testing.T) {
	_, err := exec.LookPath("go")
	if err != nil {
		t.Skip("Go is not installed in PATH, skipping subprocess test")
	}
}

func TestCommand_Convert(t *testing.T) {
	skipIfGoNotInstalled(t)

	tempDir, err := os.MkdirTemp("", "cli-command-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a simple TS file
	tsFile := filepath.Join(tempDir, "simple.ts")
	err = os.WriteFile(tsFile, []byte("const hello = 'world';\nconsole.log(hello);"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test ts file: %v", err)
	}

	outGo := filepath.Join(tempDir, "output.go")

	// The command we want to spawn is 'go run github.com/el-j/ts2go/cmd/ts2go convert --in <tsFile> --out <outGo>'
	// We run this from the current module context, so go run should resolve the package correctly.
	cmd := exec.Command("go", "run", "github.com/el-j/ts2go/cmd/ts2go", "convert", "--in", tsFile, "--out", outGo)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("convert command failed: %v\nOutput:\n%s", err, string(output))
	}

	// Verify the output Go file exists
	if _, err := os.Stat(outGo); os.IsNotExist(err) {
		t.Errorf("output.go was not created by the convert command")
	}
}
