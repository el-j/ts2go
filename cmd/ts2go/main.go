package main

import (
"fmt"
"os"

"github.com/el-j/ts2go/internal/transpiler"
"github.com/el-j/ts2go/pkg/cli"
)

func main() {
if len(os.Args) < 2 {
printUsage()
os.Exit(1)
}

command := os.Args[1]
args := os.Args[2:]

switch command {
case "convert", "c":
if err := convertCommand(args); err != nil {
fmt.Fprintf(os.Stderr, "Error: %v\n", err)
os.Exit(1)
}
case "transpile", "t":
if err := cli.TranspileCommand(args); err != nil {
fmt.Fprintf(os.Stderr, "Error: %v\n", err)
os.Exit(1)
}
case "analyze", "a":
if err := cli.AnalyzeCommand(args); err != nil {
fmt.Fprintf(os.Stderr, "Error: %v\n", err)
os.Exit(1)
}
case "help", "h", "-h", "--help":
printUsage()
case "version", "v", "-v", "--version":
printVersion()
default:
fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
printUsage()
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
fmt.Println("ts2go - TypeScript to Go Transpiler")
fmt.Println()
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
fmt.Println("For more information, visit: https://github.com/el-j/ts2go")
}

func printVersion() {
fmt.Println("ts2go version 0.1.0")
fmt.Println("TypeScript to Go Transpiler")
fmt.Println("https://github.com/el-j/ts2go")
}
