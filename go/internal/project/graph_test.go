package project

import (
	"testing"

	"github.com/el-j/ts2go/internal/analyzer"
	"github.com/el-j/ts2go/internal/mapper"
)

func TestBuildDependencyGraph(t *testing.T) {
	// Create a project with some files and imports
	project := &Project{
		RootDir: "/test/project",
		Files: []*SourceFile{
			{
				Path:         "/test/project/index.ts",
				RelativePath: "index.ts",
				Imports: []analyzer.Import{
					{Source: "./src/app", Type: analyzer.ImportTypeLocal},
				},
			},
			{
				Path:         "/test/project/src/app.ts",
				RelativePath: "src/app.ts",
				Imports: []analyzer.Import{
					{Source: "./utils/helper", Type: analyzer.ImportTypeLocal},
					{Source: "axios", Type: analyzer.ImportTypePackage},
				},
			},
			{
				Path:         "/test/project/src/utils/helper.ts",
				RelativePath: "src/utils/helper.ts",
				Imports:      []analyzer.Import{},
			},
		},
		Dependencies: make(map[string]*mapper.Classification),
	}

	// Build dependency graph
	graph, err := project.BuildDependencyGraph()
	if err != nil {
		t.Fatalf("BuildDependencyGraph failed: %v", err)
	}

	// Verify nodes
	if len(graph.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(graph.Nodes))
	}

	// Verify dependencies
	indexNode := graph.Nodes["index.ts"]
	if indexNode == nil {
		t.Fatal("index.ts node not found")
	}
	if len(indexNode.Dependencies) != 1 {
		t.Errorf("Expected index.ts to have 1 dependency, got %d", len(indexNode.Dependencies))
	}

	appNode := graph.Nodes["src/app.ts"]
	if appNode == nil {
		t.Fatal("src/app.ts node not found")
	}
	if len(appNode.Dependencies) != 1 { // Should only have helper, not axios
		t.Errorf("Expected src/app.ts to have 1 dependency, got %d", len(appNode.Dependencies))
	}
	if len(appNode.Dependents) != 1 {
		t.Errorf("Expected src/app.ts to have 1 dependent (index), got %d", len(appNode.Dependents))
	}

	helperNode := graph.Nodes["src/utils/helper.ts"]
	if helperNode == nil {
		t.Fatal("src/utils/helper.ts node not found")
	}
	if len(helperNode.Dependencies) != 0 {
		t.Errorf("Expected helper.ts to have 0 dependencies, got %d", len(helperNode.Dependencies))
	}
	if len(helperNode.Dependents) != 1 {
		t.Errorf("Expected helper.ts to have 1 dependent (app), got %d", len(helperNode.Dependents))
	}
}

func TestResolveBuildOrder(t *testing.T) {
	// Create a simple dependency graph
	project := &Project{
		RootDir: "/test/project",
		Files: []*SourceFile{
			{
				Path:         "/test/project/index.ts",
				RelativePath: "index.ts",
				Imports: []analyzer.Import{
					{Source: "./a", Type: analyzer.ImportTypeLocal},
					{Source: "./b", Type: analyzer.ImportTypeLocal},
				},
			},
			{
				Path:         "/test/project/a.ts",
				RelativePath: "a.ts",
				Imports: []analyzer.Import{
					{Source: "./c", Type: analyzer.ImportTypeLocal},
				},
			},
			{
				Path:         "/test/project/b.ts",
				RelativePath: "b.ts",
				Imports: []analyzer.Import{
					{Source: "./c", Type: analyzer.ImportTypeLocal},
				},
			},
			{
				Path:         "/test/project/c.ts",
				RelativePath: "c.ts",
				Imports:      []analyzer.Import{},
			},
		},
		Dependencies: make(map[string]*mapper.Classification),
	}

	graph, err := project.BuildDependencyGraph()
	if err != nil {
		t.Fatalf("BuildDependencyGraph failed: %v", err)
	}

	// Resolve build order
	order, err := graph.ResolveBuildOrder()
	if err != nil {
		t.Fatalf("ResolveBuildOrder failed: %v", err)
	}

	// Verify order (c must come before a and b, a and b before index)
	if len(order) != 4 {
		t.Errorf("Expected 4 files in order, got %d", len(order))
	}

	// Find positions in order
	positions := make(map[string]int)
	for i, file := range order {
		positions[file] = i
	}

	// Verify c comes before a and b
	if positions["c.ts"] >= positions["a.ts"] {
		t.Error("c.ts should come before a.ts")
	}
	if positions["c.ts"] >= positions["b.ts"] {
		t.Error("c.ts should come before b.ts")
	}

	// Verify a and b come before index
	if positions["a.ts"] >= positions["index.ts"] {
		t.Error("a.ts should come before index.ts")
	}
	if positions["b.ts"] >= positions["index.ts"] {
		t.Error("b.ts should come before index.ts")
	}
}

func TestDetectCircularDependencies(t *testing.T) {
	// Create a project with circular dependencies
	project := &Project{
		RootDir: "/test/project",
		Files: []*SourceFile{
			{
				Path:         "/test/project/a.ts",
				RelativePath: "a.ts",
				Imports: []analyzer.Import{
					{Source: "./b", Type: analyzer.ImportTypeLocal},
				},
			},
			{
				Path:         "/test/project/b.ts",
				RelativePath: "b.ts",
				Imports: []analyzer.Import{
					{Source: "./c", Type: analyzer.ImportTypeLocal},
				},
			},
			{
				Path:         "/test/project/c.ts",
				RelativePath: "c.ts",
				Imports: []analyzer.Import{
					{Source: "./a", Type: analyzer.ImportTypeLocal},
				},
			},
		},
		Dependencies: make(map[string]*mapper.Classification),
	}

	graph, err := project.BuildDependencyGraph()
	if err != nil {
		t.Fatalf("BuildDependencyGraph failed: %v", err)
	}

	// Detect cycles
	cycles := graph.DetectCircularDependencies()
	if len(cycles) == 0 {
		t.Error("Expected to detect circular dependencies")
	}

	// Verify build order fails
	_, err = graph.ResolveBuildOrder()
	if err == nil {
		t.Error("Expected ResolveBuildOrder to fail with circular dependencies")
	}
}

func TestGetDependencyStats(t *testing.T) {
	project := &Project{
		RootDir: "/test/project",
		Files: []*SourceFile{
			{
				Path:         "/test/project/a.ts",
				RelativePath: "a.ts",
				Imports: []analyzer.Import{
					{Source: "./b", Type: analyzer.ImportTypeLocal},
					{Source: "./c", Type: analyzer.ImportTypeLocal},
				},
			},
			{
				Path:         "/test/project/b.ts",
				RelativePath: "b.ts",
				Imports: []analyzer.Import{
					{Source: "./c", Type: analyzer.ImportTypeLocal},
				},
			},
			{
				Path:         "/test/project/c.ts",
				RelativePath: "c.ts",
				Imports:      []analyzer.Import{},
			},
		},
		Dependencies: make(map[string]*mapper.Classification),
	}

	graph, err := project.BuildDependencyGraph()
	if err != nil {
		t.Fatalf("BuildDependencyGraph failed: %v", err)
	}

	stats := graph.GetDependencyStats()

	// Verify stats
	if stats["total_files"] != 3 {
		t.Errorf("Expected 3 total files, got %v", stats["total_files"])
	}

	if stats["total_deps"] != 3 {
		t.Errorf("Expected 3 total dependencies, got %v", stats["total_deps"])
	}

	if stats["max_deps"] != 2 {
		t.Errorf("Expected max 2 dependencies, got %v", stats["max_deps"])
	}

	if stats["max_deps_file"] != "a.ts" {
		t.Errorf("Expected a.ts to have max dependencies, got %v", stats["max_deps_file"])
	}
}

func TestAnalyzeDependencies(t *testing.T) {
	// Create a project with npm dependencies
	project := &Project{
		RootDir: "/test/project",
		Files: []*SourceFile{
			{
				Path:         "/test/project/index.ts",
				RelativePath: "index.ts",
				Imports: []analyzer.Import{
					{Source: "axios", Type: analyzer.ImportTypePackage},
					{Source: "lodash", Type: analyzer.ImportTypePackage},
					{Source: "@scope/package", Type: analyzer.ImportTypePackage},
					{Source: "fs", Type: analyzer.ImportTypeBuiltin},
				},
			},
		},
		Dependencies: make(map[string]*mapper.Classification),
	}

	// Create a mock classifier
	db := &mapper.MappingDatabase{
		Mappings: []mapper.Mapping{
			{
				Npm:    "axios",
				Go:     "github.com/go-resty/resty/v2",
				Type:   mapper.MappingTypeEquivalent,
				Status: mapper.StatusSupported,
			},
		},
	}
	classifier := mapper.NewClassifier(db)

	// Analyze dependencies
	err := project.AnalyzeDependencies(classifier)
	if err != nil {
		t.Fatalf("AnalyzeDependencies failed: %v", err)
	}

	// Verify dependencies were classified
	if len(project.Dependencies) != 3 { // axios, lodash, @scope/package (not fs)
		t.Errorf("Expected 3 dependencies, got %d", len(project.Dependencies))
	}

	// Verify axios was classified
	if _, ok := project.Dependencies["axios"]; !ok {
		t.Error("Expected axios to be in dependencies")
	}
}

func TestExtractPackageName(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{"regular package", "axios", "axios"},
		{"scoped package", "@scope/package", "@scope/package"},
		{"subpath import", "lodash/map", "lodash"},
		{"scoped with subpath", "@scope/package/subpath", "@scope/package"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPackageName(tt.source)
			if result != tt.expected {
				t.Errorf("extractPackageName(%q) = %q, want %q", tt.source, result, tt.expected)
			}
		})
	}
}

func TestResolveImportPath(t *testing.T) {
	tests := []struct {
		name       string
		rootDir    string
		sourceFile string
		importPath string
		expected   string
	}{
		{
			"relative same dir",
			"/project",
			"/project/src/app.ts",
			"./helper",
			"/project/src/helper.ts",
		},
		{
			"relative parent dir",
			"/project",
			"/project/src/app.ts",
			"../utils/helper",
			"/project/utils/helper.ts",
		},
		{
			"with extension",
			"/project",
			"/project/app.ts",
			"./helper.ts",
			"/project/helper.ts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveImportPath(tt.rootDir, tt.sourceFile, tt.importPath)
			if result != tt.expected {
				t.Errorf("resolveImportPath() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestGetSupportedPackages(t *testing.T) {
	project := &Project{
		Dependencies: map[string]*mapper.Classification{
			"axios": {
				Package: "axios",
				Class:   mapper.PackageClassEquivalent,
			},
			"lodash": {
				Package: "lodash",
				Class:   mapper.PackageClassUnsupported,
			},
			"fs": {
				Package: "fs",
				Class:   mapper.PackageClassBuiltin,
			},
		},
	}

	supported := project.GetSupportedPackages()

	// Should have axios and fs, not lodash
	if len(supported) != 2 {
		t.Errorf("Expected 2 supported packages, got %d", len(supported))
	}
}

func TestGetUnsupportedPackages(t *testing.T) {
	project := &Project{
		Dependencies: map[string]*mapper.Classification{
			"axios": {
				Package: "axios",
				Class:   mapper.PackageClassEquivalent,
			},
			"lodash": {
				Package: "lodash",
				Class:   mapper.PackageClassUnsupported,
			},
		},
	}

	unsupported := project.GetUnsupportedPackages()

	// Should have only lodash
	if len(unsupported) != 1 {
		t.Errorf("Expected 1 unsupported package, got %d", len(unsupported))
	}

	if unsupported[0].Package != "lodash" {
		t.Errorf("Expected lodash, got %s", unsupported[0].Package)
	}
}
