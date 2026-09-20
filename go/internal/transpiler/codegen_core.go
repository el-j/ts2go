package transpiler

import (
	"fmt"
	"strings"
)

func (g *CodeGenerator) Generate(node *ASTNode) (string, error) {
	g.output.Reset()
	g.indent = 0
	g.imports = make(map[string]bool) // Reset imports for this generation

	// Check if we need helper functions
	needsOptionalAccess := g.needsOptionalAccess(node)
	needsNullishCoalesce := g.needsNullishCoalesce(node)

	// Track imports needed for helper functions
	if needsOptionalAccess {
		g.trackImport("reflect")
	}

	// Generate body into the output builder (to track imports)
	// Add helper functions if needed
	if needsOptionalAccess {
		g.writeLine("// optionalAccess provides safe property access for optional chaining")
		g.writeLine("func optionalAccess(obj interface{}, field string) interface{} {")
		g.indent++
		g.writeLine("if obj == nil {")
		g.indent++
		g.writeLine("return nil")
		g.indent--
		g.writeLine("}")
		g.writeLine("// Use reflection to access the field safely")
		g.writeLine("v := reflect.ValueOf(obj)")
		g.writeLine("if v.Kind() == reflect.Ptr {")
		g.indent++
		g.writeLine("if v.IsNil() {")
		g.indent++
		g.writeLine("return nil")
		g.indent--
		g.writeLine("}")
		g.writeLine("v = v.Elem()")
		g.indent--
		g.writeLine("}")
		g.writeLine("if v.Kind() != reflect.Struct {")
		g.indent++
		g.writeLine("return nil")
		g.indent--
		g.writeLine("}")
		g.writeLine("fieldVal := v.FieldByName(field)")
		g.writeLine("if !fieldVal.IsValid() {")
		g.indent++
		g.writeLine("return nil")
		g.indent--
		g.writeLine("}")
		g.writeLine("return fieldVal.Interface()")
		g.indent--
		g.writeLine("}")
		g.writeLine("")
	}

	if needsNullishCoalesce {
		g.writeLine("// nullishCoalesce provides nullish coalescing (??) behavior")
		g.writeLine("func nullishCoalesce(left, right interface{}) interface{} {")
		g.indent++
		g.writeLine("if left != nil {")
		g.indent++
		g.writeLine("return left")
		g.indent--
		g.writeLine("}")
		g.writeLine("return right")
		g.indent--
		g.writeLine("}")
		g.writeLine("")
	}

	// Track declared functions and classes for casing consistency
	if node.Statements != nil {
		for _, stmt := range node.Statements {
			if stmt.Kind == FunctionDeclaration && stmt.Name != "" {
				g.declaredFunctions[stmt.Name] = toPascalCase(stmt.Name)
			} else if stmt.Kind == ClassDeclaration && stmt.Name != "" {
				g.declaredFunctions[stmt.Name] = toPascalCase(stmt.Name)
			}
		}
	}

	// Separate type declarations from executable statements
	typeDecls := []ASTNode{}
	execStmts := []ASTNode{}

	if node.Statements != nil {
		for _, stmt := range node.Statements {
			if stmt.Kind == InterfaceDeclaration || stmt.Kind == TypeAliasDeclaration || stmt.Kind == EnumDeclaration || stmt.Kind == ClassDeclaration || stmt.Kind == FunctionDeclaration {
				typeDecls = append(typeDecls, stmt)
			} else if (stmt.Kind == VariableStatement || stmt.Kind == FirstStatement) && isTopLevelConstant(&stmt) {
				// Pure constants at top-level can be emitted at package level
				typeDecls = append(typeDecls, stmt)
			} else {
				execStmts = append(execStmts, stmt)
			}
		}
	}

	// Generate type declarations first
	for _, stmt := range typeDecls {
		if err := g.generateStatement(&stmt); err != nil {
			return "", err
		}
	}

	// If there are executable statements, wrap them in main()
	if len(execStmts) > 0 {
		g.writeLine("func main() {")
		g.indent++
		for _, stmt := range execStmts {
			if err := g.generateStatement(&stmt); err != nil {
				return "", err
			}
		}
		g.indent--
		g.writeLine("}")
	}

	// Get the generated body
	body := g.output.String()

	// Now build the final output with package, imports, and body
	result := strings.Builder{}

	// Add package declaration
	if g.module != nil {
		result.WriteString(g.getPackageDeclaration())
	} else {
		result.WriteString("package main")
	}
	result.WriteString("\n\n")

	// Add imports (now we know what's needed)
	if g.resolver != nil {
		// Use module resolver to generate imports
		importBlock := g.getModuleImports()
		if importBlock != "" {
			result.WriteString(importBlock)
		}
	} else {
		// Use tracked imports
		trackedImports := g.getTrackedImports()
		if len(trackedImports) > 0 {
			if len(trackedImports) == 1 {
				result.WriteString(fmt.Sprintf("import \"%s\"\n\n", trackedImports[0]))
			} else {
				result.WriteString("import (\n")
				for _, imp := range trackedImports {
					result.WriteString(fmt.Sprintf("\t\"%s\"\n", imp))
				}
				result.WriteString(")\n\n")
			}
		}
	}

	// Add the body
	result.WriteString(body)

	return result.String(), nil
}

// collectImports analyzes the AST to determine needed imports
func (g *CodeGenerator) collectImports(node *ASTNode) []string {
	imports := make(map[string]bool)
	g.findImports(node, imports)

	result := []string{}
	for imp := range imports {
		result = append(result, imp)
	}
	return result
}

func (g *CodeGenerator) findImports(node *ASTNode, imports map[string]bool) {
	if node == nil {
		return
	}

	// Check for console.log -> needs fmt
	if node.Kind == CallExpression {
		if len(node.Children) > 0 && node.Children[0].Kind == "PropertyAccessExpression" {
			propAccess := &node.Children[0]
			if len(propAccess.Children) > 0 {
				obj := propAccess.Children[0].Text
				prop := propAccess.Name
				if obj == "console" && prop == "log" {
					imports["fmt"] = true
				}
			}
		}
	}

	// Recursively check children
	for _, child := range node.Children {
		g.findImports(&child, imports)
	}
	if node.Body != nil {
		g.findImports(node.Body, imports)
	}
	for _, stmt := range node.Statements {
		g.findImports(&stmt, imports)
	}
	for _, decl := range node.Declarations {
		g.findImports(&decl, imports)
	}
	if node.Initializer != nil {
		g.findImports(node.Initializer, imports)
	}
}

// isTopLevelConstant checks if a variable statement consists exclusively of constants with literal initializers
func isTopLevelConstant(stmt *ASTNode) bool {
	if !stmt.IsConst || len(stmt.Declarations) == 0 {
		return false
	}
	for _, decl := range stmt.Declarations {
		if decl.Initializer == nil || !isConstantLiteral(decl.Initializer) {
			return false
		}
	}
	return true
}

// isConstantLiteral checks if an AST node is a compile-time constant literal
func isConstantLiteral(node *ASTNode) bool {
	if node == nil {
		return false
	}
	switch node.Kind {
	case NumericLiteral, StringLiteral, "FirstLiteralToken", "TrueKeyword", "FalseKeyword", "NoSubstitutionTemplateLiteral":
		return true
	default:
		return false
	}
}

// generateStatement generates code for a statement
