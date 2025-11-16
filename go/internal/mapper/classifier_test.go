package mapper

import (
	"testing"
)

func TestClassifyPackage(t *testing.T) {
	db := &MappingDatabase{
		Mappings: []Mapping{
			{Npm: "fs", Go: "runtime/fs", Type: MappingTypeRuntime, Status: StatusSupported, Complexity: ComplexitySimple},
			{Npm: "axios", Go: "github.com/go-resty/resty/v2", Type: MappingTypeEquivalent, Status: StatusSupported, Complexity: ComplexityMedium},
			{Npm: "express", Go: "github.com/gin-gonic/gin", Type: MappingTypeFramework, Status: StatusPartial, Complexity: ComplexityHigh},
			{Npm: "react", Type: MappingTypeUnsupported, Status: StatusUnsupported},
		},
	}
	db.index = make(map[string]*Mapping)
	for i := range db.Mappings {
		db.index[db.Mappings[i].Npm] = &db.Mappings[i]
	}

	classifier := NewClassifier(db)

	tests := []struct {
		name              string
		npmPackage        string
		expectedClass     PackageClass
		expectedSupported bool
		expectedConfMin   float64 // Minimum expected confidence
	}{
		{
			name:              "runtime package",
			npmPackage:        "fs",
			expectedClass:     PackageClassRuntime,
			expectedSupported: true,
			expectedConfMin:   0.9,
		},
		{
			name:              "equivalent package",
			npmPackage:        "axios",
			expectedClass:     PackageClassEquivalent,
			expectedSupported: true,
			expectedConfMin:   0.8,
		},
		{
			name:              "framework package",
			npmPackage:        "express",
			expectedClass:     PackageClassFramework,
			expectedSupported: true,
			expectedConfMin:   0.5,
		},
		{
			name:              "unsupported package",
			npmPackage:        "react",
			expectedClass:     PackageClassUnsupported,
			expectedSupported: false,
			expectedConfMin:   0.0,
		},
		{
			name:              "local file",
			npmPackage:        "./models/user",
			expectedClass:     PackageClassTranspilable,
			expectedSupported: true,
			expectedConfMin:   1.0,
		},
		{
			name:              "parent directory import",
			npmPackage:        "../utils/helper",
			expectedClass:     PackageClassTranspilable,
			expectedSupported: true,
			expectedConfMin:   1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			classification, err := classifier.ClassifyPackage(tt.npmPackage)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if classification.Class != tt.expectedClass {
				t.Errorf("Expected class %s, got %s", tt.expectedClass, classification.Class)
			}

			if classification.IsSupported != tt.expectedSupported {
				t.Errorf("Expected IsSupported=%v, got %v", tt.expectedSupported, classification.IsSupported)
			}

			if classification.Confidence < tt.expectedConfMin {
				t.Errorf("Expected confidence >= %.2f, got %.2f", tt.expectedConfMin, classification.Confidence)
			}

			t.Logf("Package: %s, Class: %s, Confidence: %.2f",
				classification.Package, classification.Class, classification.Confidence)
		})
	}
}

func TestClassifyUnknownPackage(t *testing.T) {
	db := &MappingDatabase{
		Mappings: []Mapping{},
	}
	db.index = make(map[string]*Mapping)

	classifier := NewClassifier(db)

	tests := []struct {
		name              string
		npmPackage        string
		expectedClass     PackageClass
		expectedSupported bool
	}{
		{
			name:              "type definitions",
			npmPackage:        "@types/node",
			expectedClass:     PackageClassTranspilable,
			expectedSupported: true,
		},
		{
			name:              "react plugin",
			npmPackage:        "react-router-dom",
			expectedClass:     PackageClassUnsupported,
			expectedSupported: false,
		},
		{
			name:              "webpack loader",
			npmPackage:        "webpack-dev-server",
			expectedClass:     PackageClassUnsupported,
			expectedSupported: false,
		},
		{
			name:              "babel plugin",
			npmPackage:        "@babel/preset-env",
			expectedClass:     PackageClassUnsupported,
			expectedSupported: false,
		},
		{
			name:              "unknown utility",
			npmPackage:        "some-unknown-package",
			expectedClass:     PackageClassTranspilable,
			expectedSupported: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			classification, err := classifier.ClassifyPackage(tt.npmPackage)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if classification.Class != tt.expectedClass {
				t.Errorf("Expected class %s, got %s", tt.expectedClass, classification.Class)
			}

			if classification.IsSupported != tt.expectedSupported {
				t.Errorf("Expected IsSupported=%v, got %v", tt.expectedSupported, classification.IsSupported)
			}

			t.Logf("Unknown package: %s, Guessed class: %s, Confidence: %.2f",
				classification.Package, classification.Class, classification.Confidence)
		})
	}
}

func TestClassifyPackages(t *testing.T) {
	db := &MappingDatabase{
		Mappings: []Mapping{
			{Npm: "fs", Go: "runtime/fs", Type: MappingTypeRuntime, Status: StatusSupported},
			{Npm: "axios", Go: "resty", Type: MappingTypeEquivalent, Status: StatusSupported},
			{Npm: "react", Type: MappingTypeUnsupported, Status: StatusUnsupported},
		},
	}
	db.index = make(map[string]*Mapping)
	for i := range db.Mappings {
		db.index[db.Mappings[i].Npm] = &db.Mappings[i]
	}

	classifier := NewClassifier(db)

	packages := []string{"fs", "axios", "react", "./local"}
	classifications, err := classifier.ClassifyPackages(packages)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(classifications) != 4 {
		t.Errorf("Expected 4 classifications, got %d", len(classifications))
	}

	// Verify each classification
	expectedClasses := []PackageClass{
		PackageClassRuntime,
		PackageClassEquivalent,
		PackageClassUnsupported,
		PackageClassTranspilable,
	}

	for i, classification := range classifications {
		if classification.Class != expectedClasses[i] {
			t.Errorf("Package %d: expected class %s, got %s",
				i, expectedClasses[i], classification.Class)
		}
	}
}

func TestGetSupportedPackages(t *testing.T) {
	db := &MappingDatabase{
		Mappings: []Mapping{
			{Npm: "fs", Go: "runtime/fs", Type: MappingTypeRuntime, Status: StatusSupported},
			{Npm: "axios", Go: "resty", Type: MappingTypeEquivalent, Status: StatusSupported},
			{Npm: "express", Go: "gin", Type: MappingTypeFramework, Status: StatusPartial},
			{Npm: "react", Type: MappingTypeUnsupported, Status: StatusUnsupported},
		},
	}
	db.index = make(map[string]*Mapping)
	for i := range db.Mappings {
		db.index[db.Mappings[i].Npm] = &db.Mappings[i]
	}

	classifier := NewClassifier(db)

	packages := []string{"fs", "axios", "express", "react"}
	supported, err := classifier.GetSupportedPackages(packages)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should get fs, axios, express (not react)
	if len(supported) != 3 {
		t.Errorf("Expected 3 supported packages, got %d", len(supported))
	}

	for _, classification := range supported {
		if !classification.IsSupported {
			t.Errorf("Package %s should be supported", classification.Package)
		}
	}
}

func TestGetUnsupportedPackages(t *testing.T) {
	db := &MappingDatabase{
		Mappings: []Mapping{
			{Npm: "fs", Go: "runtime/fs", Type: MappingTypeRuntime, Status: StatusSupported},
			{Npm: "react", Type: MappingTypeUnsupported, Status: StatusUnsupported},
			{Npm: "vue", Type: MappingTypeUnsupported, Status: StatusUnsupported},
		},
	}
	db.index = make(map[string]*Mapping)
	for i := range db.Mappings {
		db.index[db.Mappings[i].Npm] = &db.Mappings[i]
	}

	classifier := NewClassifier(db)

	packages := []string{"fs", "react", "vue"}
	unsupported, err := classifier.GetUnsupportedPackages(packages)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should get react, vue (not fs)
	if len(unsupported) != 2 {
		t.Errorf("Expected 2 unsupported packages, got %d", len(unsupported))
	}

	for _, classification := range unsupported {
		if classification.IsSupported {
			t.Errorf("Package %s should be unsupported", classification.Package)
		}
	}
}

func TestAnalyzeDependencies(t *testing.T) {
	db := &MappingDatabase{
		Mappings: []Mapping{
			{Npm: "fs", Go: "runtime/fs", Type: MappingTypeRuntime, Status: StatusSupported, Complexity: ComplexitySimple},
			{Npm: "axios", Go: "resty", Type: MappingTypeEquivalent, Status: StatusSupported, Complexity: ComplexityMedium},
			{Npm: "express", Go: "gin", Type: MappingTypeFramework, Status: StatusPartial, Complexity: ComplexityHigh},
			{Npm: "lodash", Go: "lo", Type: MappingTypeEquivalent, Status: StatusSupported, Complexity: ComplexitySimple},
			{Npm: "react", Type: MappingTypeUnsupported, Status: StatusUnsupported},
		},
	}
	db.index = make(map[string]*Mapping)
	for i := range db.Mappings {
		db.index[db.Mappings[i].Npm] = &db.Mappings[i]
	}

	classifier := NewClassifier(db)

	dependencies := map[string]string{
		"fs":      "builtin",
		"axios":   "^1.6.0",
		"express": "^4.18.0",
		"lodash":  "^4.17.21",
		"react":   "^18.0.0",
	}

	analysis, err := classifier.AnalyzeDependencies(dependencies)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if analysis.Total != 5 {
		t.Errorf("Expected 5 total dependencies, got %d", analysis.Total)
	}

	if analysis.Supported != 3 {
		t.Errorf("Expected 3 supported, got %d", analysis.Supported)
	}

	if analysis.Partial != 1 {
		t.Errorf("Expected 1 partial, got %d", analysis.Partial)
	}

	if analysis.Unsupported != 1 {
		t.Errorf("Expected 1 unsupported, got %d", analysis.Unsupported)
	}

	// Check by class
	if analysis.ByClass[PackageClassRuntime] != 1 {
		t.Errorf("Expected 1 runtime package, got %d", analysis.ByClass[PackageClassRuntime])
	}

	if analysis.ByClass[PackageClassEquivalent] != 2 {
		t.Errorf("Expected 2 equivalent packages, got %d", analysis.ByClass[PackageClassEquivalent])
	}

	if analysis.ByClass[PackageClassFramework] != 1 {
		t.Errorf("Expected 1 framework package, got %d", analysis.ByClass[PackageClassFramework])
	}

	if analysis.ByClass[PackageClassUnsupported] != 1 {
		t.Errorf("Expected 1 unsupported package, got %d", analysis.ByClass[PackageClassUnsupported])
	}

	t.Logf("\n%s", analysis.Summary())
}

func TestDependencyAnalysisSummary(t *testing.T) {
	analysis := &DependencyAnalysis{
		Total:       10,
		Supported:   6,
		Partial:     2,
		Unsupported: 2,
		ByClass: map[PackageClass]int{
			PackageClassRuntime:     2,
			PackageClassEquivalent:  4,
			PackageClassFramework:   2,
			PackageClassUnsupported: 2,
		},
	}

	summary := analysis.Summary()

	// Check that summary contains key information
	expectedStrings := []string{
		"Total Dependencies: 10",
		"Supported: 6",
		"Partial: 2",
		"Unsupported: 2",
		"Runtime: 2",
		"Equivalent: 4",
		"Framework: 2",
	}

	for _, expected := range expectedStrings {
		if !contains(summary, expected) {
			t.Errorf("Summary missing expected string: %s", expected)
		}
	}

	t.Logf("Summary:\n%s", summary)
}

func TestGetUnsupportedList(t *testing.T) {
	analysis := &DependencyAnalysis{
		Classifications: []*Classification{
			{Package: "fs", IsSupported: true},
			{Package: "axios", IsSupported: true},
			{Package: "react", IsSupported: false},
			{Package: "vue", IsSupported: false},
		},
	}

	unsupported := analysis.GetUnsupportedList()

	if len(unsupported) != 2 {
		t.Errorf("Expected 2 unsupported packages, got %d", len(unsupported))
	}

	expected := map[string]bool{"react": true, "vue": true}
	for _, pkg := range unsupported {
		if !expected[pkg] {
			t.Errorf("Unexpected package in unsupported list: %s", pkg)
		}
	}
}

func TestPackageClassString(t *testing.T) {
	tests := []struct {
		class    PackageClass
		expected string
	}{
		{PackageClassBuiltin, "builtin"},
		{PackageClassRuntime, "runtime"},
		{PackageClassEquivalent, "equivalent"},
		{PackageClassFramework, "framework"},
		{PackageClassTranspilable, "transpilable"},
		{PackageClassUnsupported, "unsupported"},
	}

	for _, tt := range tests {
		if tt.class.String() != tt.expected {
			t.Errorf("Expected %s, got %s", tt.expected, tt.class.String())
		}
	}
}
