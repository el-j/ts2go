package main

import (
	"fmt"
	"os"

	"github.com/el-j/ts2go/internal/transpiler"
	"github.com/el-j/ts2go/pkg/cli"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	var err error
	switch command {
	case "convert", "c":
		err = convertCommand(args)
	case "transpile", "t":
		err = cli.TranspileCommand(args)
	case "analyze", "a":
		err = cli.AnalyzeCommand(args)
	case "ui":
		err = cli.UICommand(args)
	case "help", "h", "-h", "--help":
		printUsage()
		return
	case "version", "v", "-v", "--version":
		printVersion()
		return
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// convertCommand handles simple single-file transpilation (legacy --in/--out flags)
func convertCommand(args []string) error {
	var inputFile, outputFile string

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--in", "-i":
			if i+1 < len(args) {
				inputFile = args[i+1]
				i++
			}
		case "--out", "-o":
			if i+1 < len(args) {
				outputFile = args[i+1]
				i++
			}
		}
	}

	if inputFile == "" || outputFile == "" {
		return fmt.Errorf("usage: ts2go convert --in <input.ts> --out <output.go>")
	}

	fmt.Printf("Transpiling: %s -> %s\n", inputFile, outputFile)
	if err := transpiler.Transpile(inputFile, outputFile); err != nil {
		return fmt.Errorf("transpilation failed: %w", err)
	}

	fmt.Println("✓ Transpilation successful")
	return nil
}

func printUsage() {
	fmt.Println("TS2Go - TypeScript to Go Transpiler")
	fmt.Printf("Version: %s\n\n", version)
	fmt.Println("Usage:")
	fmt.Println("  ts2go <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  convert, c     Convert a single TypeScript file to Go")
	fmt.Println("                 ts2go convert --in input.ts --out output.go")
	fmt.Println()
	fmt.Println("  transpile, t   Transpile an entire TypeScript project")
	fmt.Println("                 ts2go transpile <project-dir> [--out <output-dir>] [--verbose] [--watch]")
	fmt.Println()
	fmt.Println("  analyze, a     Analyze a TypeScript project's dependencies")
	fmt.Println("                 ts2go analyze <project-dir>")
	fmt.Println()
	fmt.Println("  ui             Start web-based UI")
	fmt.Println("    Options:")
	fmt.Println("      --port, -p <port>    Port number (default: 8080)")
	fmt.Println("      --open, -o           Open browser automatically")
	fmt.Println()
	fmt.Println("  help, h        Show this help message")
	fmt.Println("  version, v     Show version information")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Convert single file")
	fmt.Println("  ts2go convert --in app.ts --out app.go")
	fmt.Println()
	fmt.Println("  # Transpile entire project")
	fmt.Println("  ts2go transpile ./my-project --out ./output")
	fmt.Println()
	fmt.Println("  # Analyze project dependencies")
	fmt.Println("  ts2go analyze ./my-project")
	fmt.Println()
	fmt.Println("  # Start web UI")
	fmt.Println("  ts2go ui --port 8080 --open")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/el-j/ts2go")
}

func printVersion() {
	fmt.Printf("ts2go version %s\n", version)
	fmt.Println("TypeScript to Go Transpiler")
	fmt.Println("https://github.com/el-j/ts2go")
}
