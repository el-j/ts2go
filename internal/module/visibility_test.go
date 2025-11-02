package module

import (
	"testing"
)

func TestSymbolVisibility_IsExported(t *testing.T) {
	registry := NewExportRegistry()

	module := &Module{
		Path: "/project/src/models/user.ts",
		Exports: []Export{
			{Name: "User", LocalName: "User", Type: ExportNamed},
			{Name: "createUser", LocalName: "createUser", Type: ExportNamed},
			{Name: "default", LocalName: "UserService", Type: ExportDefault},
		},
	}

	registry.AddModule(module)
	visibility := NewSymbolVisibility(registry, "/project/src/models/user.ts")

	tests := []struct {
		name     string
		symbol   string
		expected bool
	}{
		{"exported interface", "User", true},
		{"exported function", "createUser", true},
		{"default export", "UserService", true},
		{"unexported symbol", "privateHelper", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := visibility.IsExported(tt.symbol)
			if result != tt.expected {
				t.Errorf("IsExported(%s) = %v, want %v", tt.symbol, result, tt.expected)
			}
		})
	}
}

func TestSymbolVisibility_GetGoSymbolName(t *testing.T) {
	registry := NewExportRegistry()

	module := &Module{
		Path: "/project/src/models/user.ts",
		Exports: []Export{
			{Name: "User", LocalName: "User", Type: ExportNamed},
			{Name: "createUser", LocalName: "createUser", Type: ExportNamed},
		},
	}

	registry.AddModule(module)
	visibility := NewSymbolVisibility(registry, "/project/src/models/user.ts")

	tests := []struct {
		name     string
		symbol   string
		expected string
	}{
		{"exported interface - already PascalCase", "User", "User"},
		{"exported function - needs PascalCase", "createUser", "CreateUser"},
		{"unexported function - remains camelCase", "privateHelper", "privateHelper"},
		{"unexported variable", "count", "count"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := visibility.GetGoSymbolName(tt.symbol)
			if result != tt.expected {
				t.Errorf("GetGoSymbolName(%s) = %s, want %s", tt.symbol, result, tt.expected)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"already PascalCase", "User", "User"},
		{"camelCase", "createUser", "CreateUser"},
		{"lowercase", "user", "User"},
		{"snake_case", "user_profile", "UserProfile"},
		{"kebab-case", "user-profile", "UserProfile"},
		{"multiple words", "get_user_by_id", "GetUserById"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToPascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToPascalCase(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"already camelCase", "createUser", "createUser"},
		{"PascalCase", "User", "user"},
		{"lowercase", "user", "user"},
		{"snake_case", "user_profile", "userProfile"},
		{"kebab-case", "user-profile", "userProfile"},
		{"multiple words", "get_user_by_id", "getUserById"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("toCamelCase(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSymbolVisibility_GetExportedSymbols(t *testing.T) {
	registry := NewExportRegistry()

	module := &Module{
		Path: "/project/src/models/user.ts",
		Exports: []Export{
			{Name: "User", LocalName: "User", Type: ExportNamed},
			{Name: "createUser", LocalName: "createUser", Type: ExportNamed},
			{Name: "default", LocalName: "UserService", Type: ExportDefault},
		},
	}

	registry.AddModule(module)
	visibility := NewSymbolVisibility(registry, "/project/src/models/user.ts")

	symbols := visibility.GetExportedSymbols()

	if len(symbols) != 3 {
		t.Errorf("Expected 3 exported symbols, got %d", len(symbols))
	}
}

func TestSymbolVisibility_GetDefaultExport(t *testing.T) {
	registry := NewExportRegistry()

	module := &Module{
		Path: "/project/src/models/user.ts",
		Exports: []Export{
			{Name: "User", LocalName: "User", Type: ExportNamed},
			{Name: "default", LocalName: "UserService", Type: ExportDefault},
		},
	}

	registry.AddModule(module)
	visibility := NewSymbolVisibility(registry, "/project/src/models/user.ts")

	defaultExport := visibility.GetDefaultExport()

	if defaultExport == nil {
		t.Fatal("Expected default export, got nil")
	}

	if defaultExport.LocalName != "UserService" {
		t.Errorf("Expected default export LocalName UserService, got %s", defaultExport.LocalName)
	}
}

func TestSymbolVisibility_GetDefaultExport_NotFound(t *testing.T) {
	registry := NewExportRegistry()

	module := &Module{
		Path: "/project/src/models/user.ts",
		Exports: []Export{
			{Name: "User", LocalName: "User", Type: ExportNamed},
		},
	}

	registry.AddModule(module)
	visibility := NewSymbolVisibility(registry, "/project/src/models/user.ts")

	defaultExport := visibility.GetDefaultExport()

	if defaultExport != nil {
		t.Error("Expected no default export, got one")
	}
}

func TestSymbolVisibility_ShouldExportSymbol(t *testing.T) {
	registry := NewExportRegistry()

	module := &Module{
		Path: "/project/src/models/user.ts",
		Exports: []Export{
			{Name: "User", LocalName: "User", Type: ExportNamed},
		},
	}

	registry.AddModule(module)
	visibility := NewSymbolVisibility(registry, "/project/src/models/user.ts")

	if !visibility.ShouldExportSymbol("User") {
		t.Error("Should export User")
	}

	if visibility.ShouldExportSymbol("privateHelper") {
		t.Error("Should not export privateHelper")
	}
}

func TestSymbolVisibility_GetExportName(t *testing.T) {
	registry := NewExportRegistry()

	module := &Module{
		Path: "/project/src/models/user.ts",
		Exports: []Export{
			{Name: "User", LocalName: "User", Type: ExportNamed},
			{Name: "default", LocalName: "config", Type: ExportDefault},
		},
	}

	registry.AddModule(module)
	visibility := NewSymbolVisibility(registry, "/project/src/models/user.ts")

	tests := []struct {
		name       string
		symbol     string
		symbolType string
		expected   string
	}{
		{"named export", "User", "", "User"},
		{"default export with type", "config", "Config", "Config"},
		{"default export without type", "config", "", "Config"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := visibility.GetExportName(tt.symbol, tt.symbolType)
			if result != tt.expected {
				t.Errorf("GetExportName(%s, %s) = %s, want %s", tt.symbol, tt.symbolType, result, tt.expected)
			}
		})
	}
}

func TestGetImportAlias(t *testing.T) {
	tests := []struct {
		name            string
		packagePath     string
		existingAliases map[string]bool
		expectedEmpty   bool
	}{
		{
			name:            "no conflict",
			packagePath:     "github.com/user/project/models",
			existingAliases: map[string]bool{},
			expectedEmpty:   true,
		},
		{
			name:            "conflict exists",
			packagePath:     "github.com/user/project/models",
			existingAliases: map[string]bool{"models": true},
			expectedEmpty:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetImportAlias(tt.packagePath, tt.existingAliases)
			if tt.expectedEmpty && result != "" {
				t.Errorf("Expected empty alias, got %s", result)
			}
			if !tt.expectedEmpty && result == "" {
				t.Error("Expected non-empty alias")
			}
		})
	}
}

func TestQualifySymbol(t *testing.T) {
	tests := []struct {
		name         string
		symbol       string
		packageAlias string
		expected     string
	}{
		{"no alias", "User", "", "User"},
		{"with alias", "User", "models", "models.User"},
		{"with alias2", "CreateUser", "usermodels", "usermodels.CreateUser"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := QualifySymbol(tt.symbol, tt.packageAlias)
			if result != tt.expected {
				t.Errorf("QualifySymbol(%s, %s) = %s, want %s", tt.symbol, tt.packageAlias, result, tt.expected)
			}
		})
	}
}
