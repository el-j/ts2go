package module

import (
	"strings"
	"unicode"
)

// SymbolVisibility handles Go symbol visibility based on TypeScript exports
type SymbolVisibility struct {
	registry *ExportRegistry
	module   *Module
}

// NewSymbolVisibility creates a new symbol visibility handler
func NewSymbolVisibility(registry *ExportRegistry, filePath string) *SymbolVisibility {
	return &SymbolVisibility{
		registry: registry,
		module:   registry.GetModule(filePath),
	}
}

// IsExported checks if a symbol is exported from the module
func (s *SymbolVisibility) IsExported(symbolName string) bool {
	if s.module == nil {
		return false
	}

	for _, exp := range s.module.Exports {
		if exp.Name == symbolName || exp.LocalName == symbolName {
			return true
		}
		// Check for default export
		if exp.Type == ExportDefault && (symbolName == "default" || exp.LocalName == symbolName) {
			return true
		}
	}

	return false
}

// GetGoSymbolName converts a TypeScript symbol name to Go symbol name
// Exported symbols become PascalCase, unexported remain camelCase
func (s *SymbolVisibility) GetGoSymbolName(tsSymbolName string) string {
	if s.IsExported(tsSymbolName) {
		return ToPascalCase(tsSymbolName)
	}
	return toCamelCase(tsSymbolName)
}

// ToPascalCase converts a name to PascalCase (public in Go)
func ToPascalCase(name string) string {
	if name == "" {
		return name
	}

	// Handle snake_case and kebab-case
	if strings.Contains(name, "_") || strings.Contains(name, "-") {
		return convertDelimitedToPascal(name)
	}

	// Handle already PascalCase
	if unicode.IsUpper(rune(name[0])) {
		return name
	}

	// Convert first letter to uppercase
	runes := []rune(name)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// toCamelCase converts a name to camelCase (private in Go)
func toCamelCase(name string) string {
	if name == "" {
		return name
	}

	// Handle snake_case and kebab-case
	if strings.Contains(name, "_") || strings.Contains(name, "-") {
		return convertDelimitedToCamel(name)
	}

	// Handle already camelCase
	if unicode.IsLower(rune(name[0])) {
		return name
	}

	// Convert first letter to lowercase
	runes := []rune(name)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// convertDelimitedToPascal converts snake_case or kebab-case to PascalCase
func convertDelimitedToPascal(name string) string {
	// Split by underscore or hyphen
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == '-'
	})

	result := ""
	for _, part := range parts {
		if len(part) > 0 {
			result += strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
		}
	}

	return result
}

// convertDelimitedToCamel converts snake_case or kebab-case to camelCase
func convertDelimitedToCamel(name string) string {
	pascal := convertDelimitedToPascal(name)
	if len(pascal) == 0 {
		return pascal
	}

	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// GetExportedSymbols returns all exported symbols from the module
func (s *SymbolVisibility) GetExportedSymbols() []Export {
	if s.module == nil {
		return nil
	}
	return s.module.Exports
}

// GetDefaultExport returns the default export if it exists
func (s *SymbolVisibility) GetDefaultExport() *Export {
	if s.module == nil {
		return nil
	}

	for i := range s.module.Exports {
		if s.module.Exports[i].Type == ExportDefault {
			return &s.module.Exports[i]
		}
	}

	return nil
}

// ShouldExportSymbol checks if a symbol should be exported in Go
func (s *SymbolVisibility) ShouldExportSymbol(tsSymbolName string) bool {
	return s.IsExported(tsSymbolName)
}

// GetExportName returns the Go export name for a TypeScript symbol
// For default exports, it returns the type name or a sensible default
func (s *SymbolVisibility) GetExportName(tsSymbolName string, symbolType string) string {
	// Check if it's a default export
	if s.module != nil {
		for _, exp := range s.module.Exports {
			if exp.Type == ExportDefault && exp.LocalName == tsSymbolName {
				// For default exports, use the type name if available
				if symbolType != "" && symbolType != "default" {
					return ToPascalCase(symbolType)
				}
				// Otherwise use the local name
				return ToPascalCase(exp.LocalName)
			}
		}
	}

	return s.GetGoSymbolName(tsSymbolName)
}

// GetImportAlias generates an import alias if needed to avoid conflicts
func GetImportAlias(packagePath string, existingAliases map[string]bool) string {
	// Extract package name from path
	parts := strings.Split(packagePath, "/")
	baseName := parts[len(parts)-1]

	// If no conflict, no alias needed
	if !existingAliases[baseName] {
		return ""
	}

	// Generate alias with number suffix
	for i := 2; i < 100; i++ {
		alias := baseName + string(rune('0'+i))
		if !existingAliases[alias] {
			return alias
		}
	}

	// Fallback
	return baseName + "_imported"
}

// QualifySymbol adds package prefix to a symbol if needed
func QualifySymbol(symbol, packageAlias string) string {
	if packageAlias == "" {
		return symbol
	}
	return packageAlias + "." + symbol
}
