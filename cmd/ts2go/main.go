package main

import (
	"fmt"
	"os"

	"github.com/yourusername/ts2go/pkg/cli"
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
	case "transpile":
		err = cli.TranspileCommand(args)
	case "analyze":
		err = cli.AnalyzeCommand(args)
	case "ui":
		err = cli.UICommand(args)
	case "version", "--version", "-v":
		fmt.Printf("ts2go version %s\n", version)
		return
	case "help", "--help", "-h":
		printUsage()
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

func printUsage() {
	fmt.Println("TS2Go - TypeScript to Go Transpiler")
	fmt.Printf("Version: %s\n\n", version)
	fmt.Println("Usage:")
	fmt.Println("  ts2go <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  transpile <project-dir>  Transpile TypeScript project to Go")
	fmt.Println("    Options:")
	fmt.Println("      --out, -o <dir>      Output directory (default: output)")
	fmt.Println("      --verbose, -v        Verbose output")
	fmt.Println("      --quiet, -q          Quiet mode")
	fmt.Println("      --watch, -w          Watch mode for continuous transpilation")
	fmt.Println()
	fmt.Println("  analyze <project-dir>    Analyze TypeScript project")
	fmt.Println()
	fmt.Println("  ui                       Start web-based UI")
	fmt.Println("    Options:")
	fmt.Println("      --port, -p <port>    Port number (default: 8080)")
	fmt.Println("      --open, -o           Open browser automatically")
	fmt.Println()
	fmt.Println("  version                  Show version")
	fmt.Println("  help                     Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ts2go transpile ./my-project --out ./output")
	fmt.Println("  ts2go analyze ./my-project")
	fmt.Println("  ts2go ui --port 8080 --open")
}
