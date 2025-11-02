package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/ts2go/internal/mapper"
)

func main() {
	// Command flags
	lookupCmd := flag.NewFlagSet("lookup", flag.ExitOnError)
	lookupPackage := lookupCmd.String("package", "", "npm package name to lookup")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listType := listCmd.String("type", "all", "filter by type: runtime, stdlib, equivalent, framework, unsupported, all")

	summaryCmd := flag.NewFlagSet("summary", flag.ExitOnError)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Load mappings - try to find in workspace
	db, err := mapper.LoadDefaultMappings()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load mappings: %v\n", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "lookup":
		lookupCmd.Parse(os.Args[2:])
		if *lookupPackage == "" {
			fmt.Fprintf(os.Stderr, "Error: -package flag is required\n")
			lookupCmd.PrintDefaults()
			os.Exit(1)
		}
		lookupMapping(db, *lookupPackage)

	case "list":
		listCmd.Parse(os.Args[2:])
		listMappings(db, *listType)

	case "summary":
		summaryCmd.Parse(os.Args[2:])
		fmt.Println(db.Summary())

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Mapping Database CLI - Query npm to Go package mappings")
	fmt.Println("\nUsage:")
	fmt.Println("  mapping-cli lookup -package <npm-package>")
	fmt.Println("  mapping-cli list [-type runtime|stdlib|equivalent|framework|unsupported|all]")
	fmt.Println("  mapping-cli summary")
	fmt.Println("\nExamples:")
	fmt.Println("  mapping-cli lookup -package axios")
	fmt.Println("  mapping-cli list -type runtime")
	fmt.Println("  mapping-cli summary")
}

func lookupMapping(db *mapper.MappingDatabase, npmPackage string) {
	mapping, err := db.GetMapping(npmPackage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("NPM Package: %s\n", mapping.Npm)
	if mapping.Go != "" {
		fmt.Printf("Go Package: %s\n", mapping.Go)
	}
	fmt.Printf("Type: %s\n", mapping.Type)
	fmt.Printf("Status: %s\n", mapping.Status)
	if mapping.Complexity != "" {
		fmt.Printf("Complexity: %s\n", mapping.Complexity)
	}
	if mapping.Description != "" {
		fmt.Printf("Description: %s\n", mapping.Description)
	}

	if len(mapping.APIMappings) > 0 {
		fmt.Println("\nAPI Mappings:")
		for npmAPI, goAPI := range mapping.APIMappings {
			fmt.Printf("  %s → %s\n", npmAPI, goAPI)
		}
	}

	if mapping.Notes != "" {
		fmt.Printf("\nNotes: %s\n", mapping.Notes)
	}

	if mapping.Reason != "" {
		fmt.Printf("\nReason: %s\n", mapping.Reason)
	}

	if mapping.Suggestion != "" {
		fmt.Printf("Suggestion: %s\n", mapping.Suggestion)
	}

	if mapping.Example != "" {
		fmt.Printf("\nExample:\n%s\n", strings.TrimSpace(mapping.Example))
	}
}

func listMappings(db *mapper.MappingDatabase, filterType string) {
	var mappings []mapper.Mapping

	switch filterType {
	case "runtime":
		mappings = db.GetMappingsByType(mapper.MappingTypeRuntime)
	case "stdlib":
		mappings = db.GetMappingsByType(mapper.MappingTypeStdlib)
	case "equivalent":
		mappings = db.GetMappingsByType(mapper.MappingTypeEquivalent)
	case "framework":
		mappings = db.GetMappingsByType(mapper.MappingTypeFramework)
	case "unsupported":
		mappings = db.GetMappingsByType(mapper.MappingTypeUnsupported)
	case "all":
		mappings = db.Mappings
	default:
		fmt.Fprintf(os.Stderr, "Unknown type: %s\n", filterType)
		os.Exit(1)
	}

	if len(mappings) == 0 {
		fmt.Printf("No mappings found for type: %s\n", filterType)
		return
	}

	fmt.Printf("Mappings (%s): %d packages\n\n", filterType, len(mappings))

	// Group by status
	supported := []mapper.Mapping{}
	partial := []mapper.Mapping{}
	unsupported := []mapper.Mapping{}

	for _, m := range mappings {
		switch m.Status {
		case mapper.StatusSupported:
			supported = append(supported, m)
		case mapper.StatusPartial:
			partial = append(partial, m)
		case mapper.StatusUnsupported:
			unsupported = append(unsupported, m)
		}
	}

	if len(supported) > 0 {
		fmt.Printf("✅ Supported (%d):\n", len(supported))
		for _, m := range supported {
			if m.Go != "" {
				fmt.Printf("  %s → %s\n", m.Npm, m.Go)
			} else {
				fmt.Printf("  %s (unsupported)\n", m.Npm)
			}
		}
		fmt.Println()
	}

	if len(partial) > 0 {
		fmt.Printf("⚠️  Partial (%d):\n", len(partial))
		for _, m := range partial {
			fmt.Printf("  %s → %s\n", m.Npm, m.Go)
			if m.Notes != "" {
				fmt.Printf("    Note: %s\n", m.Notes)
			}
		}
		fmt.Println()
	}

	if len(unsupported) > 0 {
		fmt.Printf("❌ Unsupported (%d):\n", len(unsupported))
		for _, m := range unsupported {
			fmt.Printf("  %s - %s\n", m.Npm, m.Reason)
		}
	}
}
