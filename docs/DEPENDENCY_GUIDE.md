# Dependency Resolution Implementation Guide

## Overview

This guide provides concrete implementation details for the dependency resolution system described in the roadmap.

## Architecture

```
┌─────────────────┐
│  TypeScript     │
│  Project        │
└────────┬────────┘
         │
         ├─ package.json
         ├─ tsconfig.json
         └─ src/**/*.ts
                │
                ▼
┌────────────────────────────────────┐
│  Dependency Analyzer               │
│  - Parse package.json              │
│  - Build dependency graph          │
│  - Classify each dependency        │
└────────────┬───────────────────────┘
             │
     ┌───────┴───────┐
     ▼               ▼
┌─────────┐   ┌──────────────┐
│ Mapping │   │  Transpiler  │
│ Database│◄──┤  Engine      │
└─────────┘   └──────┬───────┘
                     │
        ┌────────────┼────────────┐
        ▼            ▼            ▼
   ┌────────┐  ┌─────────┐  ┌────────┐
   │ Runtime│  │ Go Equiv│  │ Custom │
   │ Wrapper│  │ Mapping │  │ Transp.│
   └────────┘  └─────────┘  └────────┘
                     │
                     ▼
             ┌───────────────┐
             │  Go Project   │
             │  with go.mod  │
             └───────────────┘
```

## Step-by-Step Implementation

### Step 1: Create Dependency Mapping Database

**File: `internal/deps/mappings.go`**

```go
package deps

type MappingType string

const (
    MappingTypeRuntime      MappingType = "runtime"
    MappingTypeStdlib       MappingType = "stdlib"
    MappingTypeEquivalent   MappingType = "equivalent"
    MappingTypeTranspile    MappingType = "transpile"
    MappingTypeUnsupported  MappingType = "unsupported"
)

type PackageMapping struct {
    NPMPackage   string
    GoModule     string
    Type         MappingType
    APIRewrites  map[string]string
    Notes        string
    Suggestion   string
}

var PackageMappings = map[string]PackageMapping{
    // Node.js Built-ins
    "fs": {
        NPMPackage: "fs",
        GoModule:   "github.com/ts2go/runtime/fs",
        Type:       MappingTypeRuntime,
        APIRewrites: map[string]string{
            "readFileSync":  "fs.ReadFile",
            "writeFileSync": "fs.WriteFile",
            "existsSync":    "fs.Exists",
        },
    },
    
    "path": {
        NPMPackage: "path",
        GoModule:   "path/filepath",
        Type:       MappingTypeStdlib,
        APIRewrites: map[string]string{
            "path.join":     "filepath.Join",
            "path.dirname":  "filepath.Dir",
            "path.basename": "filepath.Base",
            "path.extname":  "filepath.Ext",
        },
    },
    
    "http": {
        NPMPackage: "http",
        GoModule:   "net/http",
        Type:       MappingTypeStdlib,
        APIRewrites: map[string]string{
            "http.createServer": "http.Server",
            "http.get":          "http.Get",
            "http.request":      "http.NewRequest",
        },
    },
    
    // Popular Libraries
    "axios": {
        NPMPackage: "axios",
        GoModule:   "github.com/go-resty/resty/v2",
        Type:       MappingTypeEquivalent,
        APIRewrites: map[string]string{
            "axios.get":    "client.R().Get",
            "axios.post":   "client.R().Post",
            "axios.put":    "client.R().Put",
            "axios.delete": "client.R().Delete",
        },
        Notes: "Requires creating a client instance first",
    },
    
    "express": {
        NPMPackage: "express",
        GoModule:   "github.com/gin-gonic/gin",
        Type:       MappingTypeEquivalent,
        APIRewrites: map[string]string{
            "express()":      "gin.Default()",
            "app.get":        "router.GET",
            "app.post":       "router.POST",
            "app.put":        "router.PUT",
            "app.delete":     "router.DELETE",
            "app.use":        "router.Use",
            "app.listen":     "router.Run",
            "req.params":     "c.Param",
            "req.query":      "c.Query",
            "req.body":       "c.BindJSON",
            "res.json":       "c.JSON",
            "res.send":       "c.String",
            "res.status":     "c.Status",
        },
        Notes: "Context-based API requires restructuring",
    },
    
    "lodash": {
        NPMPackage: "lodash",
        GoModule:   "github.com/samber/lo",
        Type:       MappingTypeEquivalent,
        APIRewrites: map[string]string{
            "_.map":       "lo.Map",
            "_.filter":    "lo.Filter",
            "_.find":      "lo.Find",
            "_.forEach":   "lo.ForEach",
            "_.reduce":    "lo.Reduce",
            "_.uniq":      "lo.Uniq",
            "_.flatten":   "lo.Flatten",
            "_.groupBy":   "lo.GroupBy",
            "_.sortBy":    "lo.Sort",
        },
    },
    
    "moment": {
        NPMPackage: "moment",
        GoModule:   "time",
        Type:       MappingTypeStdlib,
        APIRewrites: map[string]string{
            "moment()":           "time.Now()",
            "moment.format":      "time.Format",
            "moment.add":         "time.Add",
            "moment.subtract":    "time.Sub",
            "moment.diff":        "time.Sub",
            "moment.isBefore":    "time.Before",
            "moment.isAfter":     "time.After",
        },
        Notes: "Go time format uses reference time: Mon Jan 2 15:04:05 MST 2006",
    },
    
    "uuid": {
        NPMPackage: "uuid",
        GoModule:   "github.com/google/uuid",
        Type:       MappingTypeEquivalent,
        APIRewrites: map[string]string{
            "uuid.v4":      "uuid.New().String()",
            "uuid.v1":      "uuid.NewUUID().String()",
            "uuid.parse":   "uuid.Parse",
        },
    },
    
    "bcrypt": {
        NPMPackage: "bcrypt",
        GoModule:   "golang.org/x/crypto/bcrypt",
        Type:       MappingTypeStdlib,
        APIRewrites: map[string]string{
            "bcrypt.hash":    "bcrypt.GenerateFromPassword",
            "bcrypt.compare": "bcrypt.CompareHashAndPassword",
        },
    },
    
    "dotenv": {
        NPMPackage: "dotenv",
        GoModule:   "github.com/joho/godotenv",
        Type:       MappingTypeEquivalent,
        APIRewrites: map[string]string{
            "dotenv.config": "godotenv.Load",
        },
    },
    
    "joi": {
        NPMPackage: "joi",
        GoModule:   "github.com/go-playground/validator/v10",
        Type:       MappingTypeEquivalent,
        Notes: "Validation logic requires restructuring to struct tags",
    },
    
    // Transpilable Libraries
    "validator": {
        NPMPackage: "validator",
        Type:       MappingTypeTranspile,
        Notes:      "Pure TypeScript, can be transpiled directly",
    },
    
    // Unsupported
    "react": {
        NPMPackage: "react",
        Type:       MappingTypeUnsupported,
        Notes:      "Frontend framework, not applicable to Go",
        Suggestion: "Use templ, Go templates, or keep frontend separate",
    },
    
    "vue": {
        NPMPackage: "vue",
        Type:       MappingTypeUnsupported,
        Notes:      "Frontend framework, not applicable to Go",
        Suggestion: "Keep frontend separate or use HTMX with Go templates",
    },
}

// GetMapping retrieves the mapping for an npm package
func GetMapping(npmPackage string) (PackageMapping, bool) {
    mapping, exists := PackageMappings[npmPackage]
    return mapping, exists
}

// SuggestAlternative suggests a Go alternative for unsupported packages
func SuggestAlternative(npmPackage string) string {
    if mapping, exists := GetMapping(npmPackage); exists {
        if mapping.Type == MappingTypeUnsupported && mapping.Suggestion != "" {
            return mapping.Suggestion
        }
        if mapping.GoModule != "" {
            return mapping.GoModule
        }
    }
    return "No direct equivalent found. Consider manual implementation."
}
```

### Step 2: Package.json Parser

**File: `internal/deps/parser.go`**

```go
package deps

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type PackageJSON struct {
    Name            string            `json:"name"`
    Version         string            `json:"version"`
    Main            string            `json:"main"`
    Dependencies    map[string]string `json:"dependencies"`
    DevDependencies map[string]string `json:"devDependencies"`
}

// ParsePackageJSON reads and parses package.json
func ParsePackageJSON(projectPath string) (*PackageJSON, error) {
    pkgPath := filepath.Join(projectPath, "package.json")
    
    data, err := os.ReadFile(pkgPath)
    if err != nil {
        return nil, err
    }
    
    var pkg PackageJSON
    if err := json.Unmarshal(data, &pkg); err != nil {
        return nil, err
    }
    
    return &pkg, nil
}

// GetAllDependencies returns all dependencies (including dev)
func (p *PackageJSON) GetAllDependencies() map[string]string {
    all := make(map[string]string)
    
    for name, version := range p.Dependencies {
        all[name] = version
    }
    
    for name, version := range p.DevDependencies {
        all[name] = version
    }
    
    return all
}
```

### Step 3: Dependency Analyzer

**File: `internal/deps/analyzer.go`**

```go
package deps

import (
    "fmt"
    "strings"
)

type DependencyInfo struct {
    Name       string
    Version    string
    Type       MappingType
    GoModule   string
    Strategy   string
    Notes      string
}

type AnalysisReport struct {
    Total        int
    Supported    int
    Transpilable int
    Unsupported  int
    Dependencies []DependencyInfo
    Warnings     []string
}

// AnalyzeDependencies analyzes all project dependencies
func AnalyzeDependencies(projectPath string) (*AnalysisReport, error) {
    pkg, err := ParsePackageJSON(projectPath)
    if err != nil {
        return nil, fmt.Errorf("failed to parse package.json: %w", err)
    }
    
    report := &AnalysisReport{
        Dependencies: []DependencyInfo{},
        Warnings:     []string{},
    }
    
    allDeps := pkg.GetAllDependencies()
    report.Total = len(allDeps)
    
    for name, version := range allDeps {
        info := analyzeSingleDependency(name, version)
        report.Dependencies = append(report.Dependencies, info)
        
        switch info.Type {
        case MappingTypeRuntime, MappingTypeStdlib, MappingTypeEquivalent:
            report.Supported++
        case MappingTypeTranspile:
            report.Transpilable++
        case MappingTypeUnsupported:
            report.Unsupported++
            report.Warnings = append(report.Warnings,
                fmt.Sprintf("Unsupported: %s - %s", name, info.Notes))
        }
    }
    
    return report, nil
}

func analyzeSingleDependency(name, version string) DependencyInfo {
    mapping, exists := GetMapping(name)
    
    if !exists {
        // Try to determine if it's transpilable
        if isLikelyTranspilable(name) {
            return DependencyInfo{
                Name:     name,
                Version:  version,
                Type:     MappingTypeTranspile,
                Strategy: "Automatic transpilation",
                Notes:    "Will attempt to transpile",
            }
        }
        
        return DependencyInfo{
            Name:     name,
            Version:  version,
            Type:     MappingTypeUnsupported,
            Strategy: "Manual intervention required",
            Notes:    "No mapping found",
        }
    }
    
    return DependencyInfo{
        Name:     name,
        Version:  version,
        Type:     mapping.Type,
        GoModule: mapping.GoModule,
        Strategy: string(mapping.Type),
        Notes:    mapping.Notes,
    }
}

func isLikelyTranspilable(packageName string) bool {
    // Heuristics for determining if a package can be transpiled
    transpilablePatterns := []string{
        "@types/",     // TypeScript type definitions
        "-types",      // Type packages
        "validator",   // Validation libraries
        "util",        // Utility libraries
    }
    
    for _, pattern := range transpilablePatterns {
        if strings.Contains(packageName, pattern) {
            return true
        }
    }
    
    return false
}

// GenerateReport creates a human-readable report
func (r *AnalysisReport) GenerateReport() string {
    var sb strings.Builder
    
    sb.WriteString("Dependency Analysis Report\n")
    sb.WriteString("==========================\n\n")
    
    sb.WriteString(fmt.Sprintf("Total Dependencies: %d\n", r.Total))
    sb.WriteString(fmt.Sprintf("Supported: %d (%.1f%%)\n", 
        r.Supported, float64(r.Supported)/float64(r.Total)*100))
    sb.WriteString(fmt.Sprintf("Transpilable: %d (%.1f%%)\n", 
        r.Transpilable, float64(r.Transpilable)/float64(r.Total)*100))
    sb.WriteString(fmt.Sprintf("Unsupported: %d (%.1f%%)\n\n", 
        r.Unsupported, float64(r.Unsupported)/float64(r.Total)*100))
    
    if len(r.Warnings) > 0 {
        sb.WriteString("Warnings:\n")
        for _, warning := range r.Warnings {
            sb.WriteString(fmt.Sprintf("  ⚠️  %s\n", warning))
        }
        sb.WriteString("\n")
    }
    
    sb.WriteString("Dependencies by Type:\n\n")
    
    for _, dep := range r.Dependencies {
        icon := getIconForType(dep.Type)
        sb.WriteString(fmt.Sprintf("%s %s@%s\n", icon, dep.Name, dep.Version))
        if dep.GoModule != "" {
            sb.WriteString(fmt.Sprintf("   → %s\n", dep.GoModule))
        }
        if dep.Notes != "" {
            sb.WriteString(fmt.Sprintf("   📝 %s\n", dep.Notes))
        }
        sb.WriteString("\n")
    }
    
    return sb.String()
}

func getIconForType(t MappingType) string {
    switch t {
    case MappingTypeRuntime:
        return "🔧"
    case MappingTypeStdlib:
        return "📦"
    case MappingTypeEquivalent:
        return "🔄"
    case MappingTypeTranspile:
        return "⚙️"
    case MappingTypeUnsupported:
        return "❌"
    default:
        return "❓"
    }
}
```

### Step 4: Import Rewriter

**File: `internal/transpiler/imports.go`**

```go
package transpiler

import (
    "fmt"
    "strings"
    
    "github.com/ts2go/transpiler/internal/deps"
)

type ImportRewriter struct {
    mappings map[string]deps.PackageMapping
    imports  []string
}

func NewImportRewriter() *ImportRewriter {
    return &ImportRewriter{
        mappings: make(map[string]deps.PackageMapping),
        imports:  []string{},
    }
}

// RewriteImport converts a TypeScript import to Go import
func (r *ImportRewriter) RewriteImport(npmPackage string) (string, error) {
    mapping, exists := deps.GetMapping(npmPackage)
    
    if !exists {
        return "", fmt.Errorf("no mapping found for %s", npmPackage)
    }
    
    r.mappings[npmPackage] = mapping
    
    if mapping.GoModule != "" && !contains(r.imports, mapping.GoModule) {
        r.imports = append(r.imports, mapping.GoModule)
    }
    
    return mapping.GoModule, nil
}

// RewriteAPICall transforms an API call from npm package to Go equivalent
func (r *ImportRewriter) RewriteAPICall(npmPackage, call string) (string, error) {
    mapping, exists := r.mappings[npmPackage]
    
    if !exists {
        return call, fmt.Errorf("package %s not imported", npmPackage)
    }
    
    if rewrite, ok := mapping.APIRewrites[call]; ok {
        return rewrite, nil
    }
    
    // Default: convert to PascalCase
    return toPascalCase(call), nil
}

// GetAllImports returns all Go imports needed
func (r *ImportRewriter) GetAllImports() []string {
    return r.imports
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

func toPascalCase(s string) string {
    parts := strings.Split(s, ".")
    result := []string{}
    
    for _, part := range parts {
        if len(part) > 0 {
            result = append(result, strings.ToUpper(part[:1]) + part[1:])
        }
    }
    
    return strings.Join(result, ".")
}
```

### Step 5: CLI Integration

**File: `cmd/ts2go/analyze.go`**

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
    "github.com/ts2go/transpiler/internal/deps"
)

var analyzeCmd = &cobra.Command{
    Use:   "analyze [project-path]",
    Short: "Analyze project dependencies",
    Long: `Analyze a TypeScript project's dependencies and determine
which can be automatically transpiled or mapped to Go equivalents.`,
    Args: cobra.MaximumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        projectPath := "."
        if len(args) > 0 {
            projectPath = args[0]
        }
        
        verbose, _ := cmd.Flags().GetBool("verbose")
        jsonOutput, _ := cmd.Flags().GetBool("json")
        
        report, err := deps.AnalyzeDependencies(projectPath)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }
        
        if jsonOutput {
            // Output as JSON
            data, _ := json.MarshalIndent(report, "", "  ")
            fmt.Println(string(data))
        } else {
            // Human-readable output
            fmt.Println(report.GenerateReport())
            
            if verbose {
                fmt.Println("\nDetailed Dependency Information:")
                for _, dep := range report.Dependencies {
                    fmt.Printf("\n%s@%s:\n", dep.Name, dep.Version)
                    fmt.Printf("  Type: %s\n", dep.Type)
                    fmt.Printf("  Strategy: %s\n", dep.Strategy)
                    if dep.GoModule != "" {
                        fmt.Printf("  Go Module: %s\n", dep.GoModule)
                    }
                    if dep.Notes != "" {
                        fmt.Printf("  Notes: %s\n", dep.Notes)
                    }
                }
            }
        }
        
        // Exit with error code if there are unsupported dependencies
        if report.Unsupported > 0 {
            fmt.Fprintf(os.Stderr, "\n⚠️  Warning: %d unsupported dependencies found\n", 
                report.Unsupported)
            os.Exit(1)
        }
    },
}

func init() {
    analyzeCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
    analyzeCmd.Flags().BoolP("json", "j", false, "JSON output")
    rootCmd.AddCommand(analyzeCmd)
}
```

### Step 6: Usage Example

```bash
# Analyze a project
ts2go analyze ./my-project

# Output:
# Dependency Analysis Report
# ==========================
# 
# Total Dependencies: 15
# Supported: 10 (66.7%)
# Transpilable: 3 (20.0%)
# Unsupported: 2 (13.3%)
# 
# Warnings:
#   ⚠️  Unsupported: react - Frontend framework, not applicable to Go
# 
# Dependencies by Type:
# 
# 📦 express@4.18.0
#    → github.com/gin-gonic/gin
#    📝 Context-based API requires restructuring
# 
# 🔄 axios@1.4.0
#    → github.com/go-resty/resty/v2
# 
# 🔄 lodash@4.17.21
#    → github.com/samber/lo
# 
# ❌ react@18.0.0
#    📝 Frontend framework, not applicable to Go

# Transpile with dependency resolution
ts2go transpile ./my-project --with-deps --out ./go-project

# This will:
# 1. Analyze all dependencies
# 2. Map to Go equivalents
# 3. Rewrite import statements
# 4. Transform API calls
# 5. Generate go.mod with all dependencies
```

## Next Steps

1. Implement the mapping database (start with top 20 packages)
2. Create the analyzer CLI command
3. Integrate with main transpiler
4. Add API call rewriting
5. Test with real projects
6. Expand mapping database based on usage

## Contributing

To add a new package mapping:

1. Add entry to `PackageMappings` in `mappings.go`
2. Document API rewrites
3. Add test case
4. Update documentation

See `CONTRIBUTING.md` for details.
