package transpiler

import (
	"strings"
	"unicode"
)

// writeLine writes a line of code with proper indentation
func (g *CodeGenerator) writeLine(line string) {
	for i := 0; i < g.indent; i++ {
		g.output.WriteString("\t")
	}
	g.output.WriteString(line)
	g.output.WriteString("\n")
}

// toPascalCase converts camelCase or snake_case to PascalCase
func toPascalCase(s string) string {
	if s == "" {
		return s
	}

	// Convert first character to uppercase
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// toCamelCase converts PascalCase or snake_case to camelCase
func toCamelCase(s string) string {
	if s == "" {
		return s
	}

	// Split by underscore for snake_case
	parts := strings.Split(s, "_")
	if len(parts) > 1 {
		result := strings.ToLower(parts[0])
		for i := 1; i < len(parts); i++ {
			result += toPascalCase(parts[i])
		}
		return result
	}

	// Just lowercase first character
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// hasModifier checks if a node has a specific modifier
func hasModifier(node *ASTNode, modifierKind string) bool {
	if node.Children == nil {
		return false
	}
	for _, child := range node.Children {
		if child.Kind == modifierKind {
			return true
		}
	}
	return false
}

// isPrivate checks if a node has a private modifier
func isPrivate(node *ASTNode) bool {
	return hasModifier(node, "PrivateKeyword")
}

// isStatic checks if a node has a static modifier
func isStatic(node *ASTNode) bool {
	return hasModifier(node, "StaticKeyword")
}

// getPackageDeclaration returns the package declaration for module mode
func (g *CodeGenerator) getPackageDeclaration() string {
	if g.module == nil {
		return "package main"
	}

	// Use reflection to access the Module's PackageName
	// We have to use interface{} to avoid circular dependency
	type moduleWithPackageName interface {
		GetPackageName() string
	}

	if m, ok := g.module.(moduleWithPackageName); ok {
		pkgName := m.GetPackageName()
		if g.isEntry {
			return "package main"
		}
		return "package " + pkgName
	}

	return "package main"
}

// getModuleImports returns the import block for module mode
func (g *CodeGenerator) getModuleImports() string {
	if g.resolver == nil {
		return ""
	}

	// Use reflection to get imports
	type resolverWithImports interface {
		GetImports() []string
	}

	if r, ok := g.resolver.(resolverWithImports); ok {
		imports := r.GetImports()
		if len(imports) == 0 {
			return ""
		}

		var result strings.Builder
		result.WriteString("import (\n")
		for _, imp := range imports {
			result.WriteString("\t\"" + imp + "\"\n")
		}
		result.WriteString(")\n")
		return result.String()
	}

	return ""
}

// getGoSymbolName converts TypeScript symbol name to Go symbol name based on visibility
func (g *CodeGenerator) getGoSymbolName(tsName string) string {
	if g.visibility == nil {
		// Default behavior - public by default
		return toPascalCase(tsName)
	}

	// Use reflection to check visibility
	type visibilityChecker interface {
		IsExported(string) bool
	}

	if v, ok := g.visibility.(visibilityChecker); ok {
		if v.IsExported(tsName) {
			return toPascalCase(tsName)
		}
		return toCamelCase(tsName)
	}

	return toPascalCase(tsName)
}

// needsOptionalAccess checks if the AST tree contains any optional chaining
func (g *CodeGenerator) needsOptionalAccess(node *ASTNode) bool {
	if node == nil {
		return false
	}

	// Check current node
	if node.Kind == "PropertyAccessExpression" && node.QuestionDot {
		return true
	}

	// Check children recursively
	for _, child := range node.Children {
		if g.needsOptionalAccess(&child) {
			return true
		}
	}

	// Check statements
	if node.Statements != nil {
		for _, stmt := range node.Statements {
			if g.needsOptionalAccess(&stmt) {
				return true
			}
		}
	}

	return false
}

// needsNullishCoalesce checks if the AST tree contains any nullish coalescing
func (g *CodeGenerator) needsNullishCoalesce(node *ASTNode) bool {
	if node == nil {
		return false
	}

	// Check if this is a nullish coalescing operator
	if node.Kind == "BinaryExpression" && node.Operator == "??" {
		return true
	}

	// Check children recursively
	for _, child := range node.Children {
		if g.needsNullishCoalesce(&child) {
			return true
		}
	}

	// Check statements
	if node.Statements != nil {
		for _, stmt := range node.Statements {
			if g.needsNullishCoalesce(&stmt) {
				return true
			}
		}
	}

	return false
}

// trackImport adds an import to the tracking map
func (g *CodeGenerator) trackImport(pkg string) {
	if g.imports == nil {
		g.imports = make(map[string]bool)
	}
	g.imports[pkg] = true
}

// getTrackedImports returns a slice of all tracked imports
func (g *CodeGenerator) getTrackedImports() []string {
	result := []string{}
	for imp := range g.imports {
		result = append(result, imp)
	}
	return result
}
