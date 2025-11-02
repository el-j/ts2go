package transpiler

import (
	"fmt"
	"strings"
)

func (g *CodeGenerator) Generate(node *ASTNode) (string, error) {
	g.output.Reset()
	g.indent = 0

	// Check if we need helper functions
	needsOptionalAccess := g.needsOptionalAccess(node)
	needsNullishCoalesce := g.needsNullishCoalesce(node)

	// Add package declaration
	if g.module != nil {
		// Use module package name
		g.writeLine(g.getPackageDeclaration())
	} else {
		g.writeLine("package main")
	}
	g.writeLine("")

	// Add imports
	var importBlock string
	if g.resolver != nil {
		// Use module resolver to generate imports
		importBlock = g.getModuleImports()
	} else {
		// Fallback to simple import collection
		imports := g.collectImports(node)
		if len(imports) > 0 {
			var builder strings.Builder
			builder.WriteString("import (\n")
			for _, imp := range imports {
				builder.WriteString(fmt.Sprintf("\t\"%s\"\n", imp))
			}
			builder.WriteString(")")
			importBlock = builder.String()
		}
	}

	if importBlock != "" {
		g.writeLine(importBlock)
		g.writeLine("")
	}

	// Add helper functions if needed
	if needsOptionalAccess {
		g.writeLine("import \"reflect\"")
		g.writeLine("")
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

	// Separate type declarations from executable statements
	typeDecls := []ASTNode{}
	execStmts := []ASTNode{}

	if node.Statements != nil {
		for _, stmt := range node.Statements {
			// DEBUG
			fmt.Printf("DEBUG: Processing statement kind=%s\n", stmt.Kind)
			if stmt.Kind == InterfaceDeclaration || stmt.Kind == TypeAliasDeclaration || stmt.Kind == EnumDeclaration || stmt.Kind == ClassDeclaration || stmt.Kind == FunctionDeclaration {
				typeDecls = append(typeDecls, stmt)
			} else {
				execStmts = append(execStmts, stmt)
			}
		}
	}
	fmt.Printf("DEBUG: typeDecls=%d, execStmts=%d\n", len(typeDecls), len(execStmts))

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

	return g.output.String(), nil
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

// generateStatement generates code for a statement
