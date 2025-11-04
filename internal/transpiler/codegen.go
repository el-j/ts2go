package transpiler

import (
	"fmt"
	"strings"
	"unicode"
)

// CodeGenerator generates Go code from AST nodes
type CodeGenerator struct {
	output                    strings.Builder
	indent                    int
	currentFunctionReturnType string          // Track the current function's return type for type assertions
	currentReceiverVar        string          // Track the current method's receiver variable for "this" replacement
	currentClassMembers       map[string]bool // Track private members of current class (name -> isPrivate)

	// Module system support
	module     interface{} // *module.Module - using interface{} to avoid circular dependency
	resolver   interface{} // *module.ImportResolver
	visibility interface{} // *module.SymbolVisibility
	isEntry    bool        // Is this the entry point (main package)?
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		indent:                    0,
		currentFunctionReturnType: "",
		currentReceiverVar:        "",
	}
}

// Generate produces Go code from an AST
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
func (g *CodeGenerator) generateStatement(node *ASTNode) error {
	switch node.Kind {
	case InterfaceDeclaration:
		return g.generateInterface(node)
	case TypeAliasDeclaration:
		return g.generateTypeAlias(node)
	case EnumDeclaration:
		return g.generateEnum(node)
	case ClassDeclaration:
		return g.generateClass(node)
	case FunctionDeclaration:
		return g.generateFunction(node)
	case VariableStatement, "FirstStatement":
		return g.generateVariableStatement(node)
	case ExpressionStatement:
		return g.generateExpressionStatement(node)
	case ReturnStatement:
		return g.generateReturnStatement(node)
	case IfStatement:
		return g.generateIfStatement(node)
	case ForStatement:
		return g.generateForStatement(node)
	case ForOfStatement:
		return g.generateForOfStatement(node)
	case ForInStatement:
		return g.generateForInStatement(node)
	case WhileStatement:
		return g.generateWhileStatement(node)
	case SwitchStatement:
		return g.generateSwitchStatement(node)
	case BreakStatement:
		return g.generateBreakStatement(node)
	case ContinueStatement:
		return g.generateContinueStatement(node)
	default:
		// Skip unknown statements for now
		return nil
	}
}

// generateInterface converts a TS interface to a Go struct
func (g *CodeGenerator) generateInterface(node *ASTNode) error {
	structName := toPascalCase(node.Name)
	g.writeLine(fmt.Sprintf("type %s struct {", structName))
	g.indent++

	if node.Members != nil {
		for _, member := range node.Members {
			if member.Kind == PropertySignature {
				fieldName := toPascalCase(member.Name)
				fieldType, err := g.generateType(member.Type)
				if err != nil {
					return err
				}
				jsonTag := fmt.Sprintf("`json:\"%s\"`", member.Name)
				g.writeLine(fmt.Sprintf("%s %s %s", fieldName, fieldType, jsonTag))
			}
		}
	}

	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// generateTypeAlias converts a TS type alias to a Go type
func (g *CodeGenerator) generateTypeAlias(node *ASTNode) error {
	typeName := toPascalCase(node.Name)

	// Check if this is a union type
	if node.Type != nil && node.Type.Kind == UnionType {
		return g.generateUnionType(typeName, node.Type)
	}

	// Regular type alias
	typeValue, err := g.generateType(node.Type)
	if err != nil {
		return err
	}
	g.writeLine(fmt.Sprintf("type %s = %s", typeName, typeValue))
	g.writeLine("")
	return nil
}

// generateEnum converts a TS enum to Go const declarations
func (g *CodeGenerator) generateEnum(node *ASTNode) error {
	enumName := toPascalCase(node.Name)

	// Check if this is a string enum (has string values) or numeric enum
	isStringEnum := g.isStringEnum(node)

	if isStringEnum {
		return g.generateStringEnum(enumName, node)
	} else {
		return g.generateNumericEnum(enumName, node)
	}
}

// isStringEnum determines if an enum has string values
func (g *CodeGenerator) isStringEnum(node *ASTNode) bool {
	if node.Members != nil {
		for _, member := range node.Members {
			if member.Kind == EnumMember && member.Initializer != nil {
				// Check if initializer is a string literal
				if member.Initializer.Kind == StringLiteral {
					return true
				}
			}
		}
	}
	return false
}

// generateNumericEnum generates a numeric enum using iota
func (g *CodeGenerator) generateNumericEnum(enumName string, node *ASTNode) error {
	g.writeLine(fmt.Sprintf("type %s int", enumName))
	g.writeLine("")
	g.writeLine("const (")

	if node.Members != nil {
		for i, member := range node.Members {
			if member.Kind == EnumMember {
				memberName := fmt.Sprintf("%s%s", enumName, toPascalCase(member.Name))
				if i == 0 {
					g.writeLine(fmt.Sprintf("%s %s = iota", memberName, enumName))
				} else if member.Initializer != nil {
					// Handle explicit values
					if member.Initializer.Kind == NumericLiteral {
						g.writeLine(fmt.Sprintf("%s = %s", memberName, member.Initializer.Text))
					}
				} else {
					g.writeLine(fmt.Sprintf("%s", memberName))
				}
			}
		}
	}

	g.writeLine(")")
	g.writeLine("")

	// Add String() method
	g.generateEnumStringMethod(enumName, node)
	return nil
}

// generateStringEnum generates a string enum
func (g *CodeGenerator) generateStringEnum(enumName string, node *ASTNode) error {
	g.writeLine(fmt.Sprintf("type %s string", enumName))
	g.writeLine("")
	g.writeLine("const (")

	if node.Members != nil {
		for _, member := range node.Members {
			if member.Kind == EnumMember {
				memberName := fmt.Sprintf("%s%s", enumName, toPascalCase(member.Name))
				if member.Initializer != nil && member.Initializer.Kind == StringLiteral {
					g.writeLine(fmt.Sprintf(`%s %s = "%s"`, memberName, enumName, member.Initializer.Text))
				}
			}
		}
	}

	g.writeLine(")")
	g.writeLine("")

	// Add String() method
	g.generateEnumStringMethod(enumName, node)
	return nil
}

// generateEnumStringMethod adds a String() method for enum debugging
func (g *CodeGenerator) generateEnumStringMethod(enumName string, node *ASTNode) {
	g.writeLine(fmt.Sprintf("func (e %s) String() string {", enumName))
	g.indent++
	g.writeLine("switch e {")

	if node.Members != nil {
		for _, member := range node.Members {
			if member.Kind == EnumMember {
				memberName := fmt.Sprintf("%s%s", enumName, toPascalCase(member.Name))
				if member.Initializer != nil && member.Initializer.Kind == StringLiteral {
					g.writeLine(fmt.Sprintf(`case %s:`, memberName))
					g.indent++
					g.writeLine(fmt.Sprintf(`return "%s"`, member.Initializer.Text))
					g.indent--
				} else {
					// For numeric enums, return the name
					g.writeLine(fmt.Sprintf("case %s:", memberName))
					g.indent++
					g.writeLine(fmt.Sprintf(`return "%s"`, member.Name))
					g.indent--
				}
			}
		}
	}

	g.writeLine("default:")
	g.indent++
	g.writeLine(`return "Unknown"`)
	g.indent--
	g.writeLine("}")
	g.indent--
	g.writeLine("}")
	g.writeLine("")
}

// generateClass converts a TS class to a Go struct with methods
func (g *CodeGenerator) generateClass(node *ASTNode) error {
	className := toPascalCase(node.Name)

	// Initialize member tracking for this class
	g.currentClassMembers = make(map[string]bool)
	defer func() { g.currentClassMembers = nil }()

	// First pass: collect member visibility information
	if node.Members != nil {
		for i := range node.Members {
			member := &node.Members[i]
			if member.Kind == PropertyDeclaration || member.Kind == MethodDeclaration {
				g.currentClassMembers[member.Name] = isPrivate(member)
			}
		}
	}

	// Check for base class (inheritance)
	var baseClassName string
	if node.HeritageClauses != nil && len(node.HeritageClauses) > 0 {
		for _, clause := range node.HeritageClauses {
			// Find the extends clause
			if clause.Types != nil && len(clause.Types) > 0 {
				// Get the base class name from ExpressionWithTypeArguments
				typeArg := &clause.Types[0]
				if len(typeArg.Children) > 0 && typeArg.Children[0].Kind == Identifier {
					baseClassName = toPascalCase(typeArg.Children[0].Text)
					break
				}
			}
		}
	}

	// Step 1: Generate the struct definition
	g.writeLine(fmt.Sprintf("type %s struct {", className))
	g.indent++

	// If there's a base class, embed it first
	if baseClassName != "" {
		g.writeLine(fmt.Sprintf("*%s // Embedded base class", baseClassName))
	}

	// Generate fields from PropertyDeclarations
	var constructorNode *ASTNode
	var methods []ASTNode
	var staticFields []ASTNode
	var staticMethods []ASTNode

	if node.Members != nil {
		for i := range node.Members {
			member := &node.Members[i] // Take address of slice element, not loop variable
			if member.Kind == PropertyDeclaration {
				// Skip static fields in struct - they'll be package-level vars
				if isStatic(member) {
					staticFields = append(staticFields, *member)
					continue
				}

				// Use camelCase for private fields, PascalCase for public
				fieldName := toPascalCase(member.Name)
				if isPrivate(member) {
					fieldName = toCamelCase(member.Name)
				}

				fieldType := "interface{}"
				if member.Type != nil {
					var err error
					fieldType, err = g.generateType(member.Type)
					if err != nil {
						return err
					}
				}
				g.writeLine(fmt.Sprintf("%s %s", fieldName, fieldType))
			} else if member.Kind == Constructor {
				constructorNode = member
			} else if member.Kind == MethodDeclaration {
				// Separate static and instance methods
				if isStatic(member) {
					staticMethods = append(staticMethods, *member)
				} else {
					methods = append(methods, *member)
				}
			} else if member.Kind == GetAccessor || member.Kind == SetAccessor {
				// Getters and setters are treated as methods
				methods = append(methods, *member)
			}
		}
	}

	g.indent--
	g.writeLine("}")
	g.writeLine("")

	// Generate static fields as package-level variables
	for _, field := range staticFields {
		fieldName := className + toPascalCase(field.Name)
		fieldType := "interface{}"
		if field.Type != nil {
			var err error
			fieldType, err = g.generateType(field.Type)
			if err != nil {
				return err
			}
		}

		// Generate with initializer if present
		if field.Initializer != nil {
			initValue, err := g.generateExpression(field.Initializer)
			if err != nil {
				return err
			}
			g.writeLine(fmt.Sprintf("var %s %s = %s", fieldName, fieldType, initValue))
		} else {
			g.writeLine(fmt.Sprintf("var %s %s", fieldName, fieldType))
		}
	}
	if len(staticFields) > 0 {
		g.writeLine("")
	}

	// Step 2: Generate constructor function
	if constructorNode != nil {
		if err := g.generateConstructorWithBase(className, baseClassName, constructorNode); err != nil {
			return err
		}
	}

	// Step 3: Generate instance methods
	for _, method := range methods {
		if err := g.generateMethod(className, &method); err != nil {
			return err
		}
	}

	// Step 4: Generate static methods as package-level functions
	for _, method := range staticMethods {
		if err := g.generateStaticMethod(className, &method); err != nil {
			return err
		}
	}

	return nil
}

// generateConstructorWithBase generates a New* constructor function with optional base class
func (g *CodeGenerator) generateConstructorWithBase(className string, baseClassName string, node *ASTNode) error {
	// Build parameter list
	params := []string{}
	if node.Parameters != nil {
		for _, param := range node.Parameters {
			paramName := param.Name
			paramType := "interface{}"
			if param.Type != nil {
				var err error
				paramType, err = g.generateType(param.Type)
				if err != nil {
					return err
				}
			}
			params = append(params, fmt.Sprintf("%s %s", paramName, paramType))
		}
	}

	// Write constructor signature
	g.writeLine(fmt.Sprintf("func New%s(%s) *%s {", className, strings.Join(params, ", "), className))
	g.indent++

	// Generate constructor body
	// Track super() call parameters for base class initialization
	var superArgs []string
	hasSuperCall := false

	// Check for super() call in constructor body
	if node.Body != nil && node.Body.Statements != nil {
		for _, stmt := range node.Body.Statements {
			if g.isSuperCall(&stmt) {
				hasSuperCall = true
				// Extract super() arguments
				superArgs = g.extractSuperArgs(&stmt)
				break
			}
		}
	}

	// Create instance - if there's a base class, initialize it
	if baseClassName != "" && hasSuperCall {
		g.writeLine(fmt.Sprintf("instance := &%s{", className))
		g.indent++
		g.writeLine(fmt.Sprintf("%s: New%s(%s),", baseClassName, baseClassName, strings.Join(superArgs, ", ")))
		g.indent--
		g.writeLine("}")
	} else {
		g.writeLine(fmt.Sprintf("instance := &%s{}", className))
	}

	// Generate body statements (assignments to this.field), skipping super() call
	if node.Body != nil && node.Body.Statements != nil {
		for _, stmt := range node.Body.Statements {
			// Skip super() call as we've already handled it
			if g.isSuperCall(&stmt) {
				continue
			}
			if err := g.generateConstructorStatement(&stmt, "instance"); err != nil {
				return err
			}
		}
	}

	g.writeLine("return instance")
	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// isSuperCall checks if a statement is a super() call
func (g *CodeGenerator) isSuperCall(stmt *ASTNode) bool {
	if stmt.Kind != ExpressionStatement {
		return false
	}
	if len(stmt.Children) == 0 || stmt.Children[0].Kind != CallExpression {
		return false
	}
	callExpr := &stmt.Children[0]
	if len(callExpr.Children) == 0 {
		return false
	}
	return callExpr.Children[0].Kind == SuperKeyword
}

// extractSuperArgs extracts argument expressions from super() call
func (g *CodeGenerator) extractSuperArgs(stmt *ASTNode) []string {
	args := []string{}
	if stmt.Kind != ExpressionStatement || len(stmt.Children) == 0 {
		return args
	}
	callExpr := &stmt.Children[0]
	if callExpr.Kind != CallExpression || len(callExpr.Children) < 2 {
		return args
	}

	// Arguments start from index 1 (index 0 is super keyword)
	for i := 1; i < len(callExpr.Children); i++ {
		arg := &callExpr.Children[i]
		argStr, err := g.generateExpression(arg)
		if err == nil {
			args = append(args, argStr)
		}
	}
	return args
}

// generateConstructorStatement handles statements in constructor (mainly this.x = y assignments)
func (g *CodeGenerator) generateConstructorStatement(node *ASTNode, instanceVar string) error {
	// Set receiver var context for "this" replacement
	oldReceiverVar := g.currentReceiverVar
	g.currentReceiverVar = instanceVar
	defer func() { g.currentReceiverVar = oldReceiverVar }()

	if node.Kind == ExpressionStatement {
		// Check if this is a "this.field = value" assignment
		// ExpressionStatement has the expression in the Expression property
		var binary *ASTNode
		if node.Expression != nil && node.Expression.Kind == BinaryExpression {
			binary = node.Expression
		} else if len(node.Children) > 0 && node.Children[0].Kind == BinaryExpression {
			binary = &node.Children[0]
		}

		if binary != nil && len(binary.Children) >= 2 {
			left := &binary.Children[0]
			right := &binary.Children[len(binary.Children)-1]

			// Check if left side is "this.something"
			// PropertyAccessExpression has ThisKeyword in expression property
			var isThisAccess bool
			if left.Kind == PropertyAccessExpression {
				if left.Expression != nil && left.Expression.Kind == ThisKeyword {
					isThisAccess = true
				} else if len(left.Children) > 0 && left.Children[0].Kind == ThisKeyword {
					isThisAccess = true
				}
			}

			if isThisAccess {
				// Use the member visibility map to determine casing
				fieldName := left.Name
				if g.currentClassMembers != nil && g.currentClassMembers[fieldName] {
					// Private member - use camelCase
					fieldName = toCamelCase(fieldName)
				} else {
					// Public member - use PascalCase
					fieldName = toPascalCase(fieldName)
				}
				rightExpr, err := g.generateExpression(right)
				if err != nil {
					return err
				}
				g.writeLine(fmt.Sprintf("%s.%s = %s", instanceVar, fieldName, rightExpr))
				return nil
			}
		}
	}
	// For other statements, use default generation
	return g.generateStatement(node)
}

// generateMethod generates a method with receiver
func (g *CodeGenerator) generateMethod(className string, node *ASTNode) error {
	// Handle getters and setters specially
	if node.Kind == GetAccessor {
		return g.generateGetter(className, node)
	} else if node.Kind == SetAccessor {
		return g.generateSetter(className, node)
	}

	// Use camelCase for private methods, PascalCase for public
	methodName := toPascalCase(node.Name)
	if isPrivate(node) {
		methodName = toCamelCase(node.Name)
	}

	// Build parameter list
	params := []string{}
	if node.Parameters != nil {
		for _, param := range node.Parameters {
			paramName := param.Name
			paramType := "interface{}"
			if param.Type != nil {
				var err error
				paramType, err = g.generateType(param.Type)
				if err != nil {
					return err
				}
			}
			params = append(params, fmt.Sprintf("%s %s", paramName, paramType))
		}
	}

	// Determine return type
	returnType := ""
	if node.Type != nil {
		var err error
		returnType, err = g.generateType(node.Type)
		if err != nil {
			return err
		}
	}

	// Write method signature with receiver
	receiverVar := strings.ToLower(string(className[0])) // Use first letter of class name
	signature := fmt.Sprintf("func (%s *%s) %s(%s)", receiverVar, className, methodName, strings.Join(params, ", "))
	if returnType != "" {
		signature += " " + returnType
	}
	g.writeLine(signature + " {")
	g.indent++

	// Set context for return type and receiver var
	oldReturnType := g.currentFunctionReturnType
	oldReceiverVar := g.currentReceiverVar
	g.currentFunctionReturnType = returnType
	g.currentReceiverVar = receiverVar
	defer func() {
		g.currentFunctionReturnType = oldReturnType
		g.currentReceiverVar = oldReceiverVar
	}()

	// Generate method body
	if node.Body != nil {
		if err := g.generateMethodBody(node.Body, receiverVar); err != nil {
			return err
		}
	}

	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// generateGetter generates a getter method (Get prefix)
func (g *CodeGenerator) generateGetter(className string, node *ASTNode) error {
	// Create method name with "Get" prefix
	methodName := "Get" + toPascalCase(node.Name)

	// Determine return type
	returnType := "interface{}"
	if node.Type != nil {
		var err error
		returnType, err = g.generateType(node.Type)
		if err != nil {
			return err
		}
	}

	// Write getter method signature with receiver
	receiverVar := strings.ToLower(string(className[0]))
	g.writeLine(fmt.Sprintf("func (%s *%s) %s() %s {", receiverVar, className, methodName, returnType))
	g.indent++

	// Set context for return type and receiver var
	oldReturnType := g.currentFunctionReturnType
	oldReceiverVar := g.currentReceiverVar
	g.currentFunctionReturnType = returnType
	g.currentReceiverVar = receiverVar
	defer func() {
		g.currentFunctionReturnType = oldReturnType
		g.currentReceiverVar = oldReceiverVar
	}()

	// Generate getter body
	if node.Body != nil {
		if err := g.generateMethodBody(node.Body, receiverVar); err != nil {
			return err
		}
	}

	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// generateSetter generates a setter method (Set prefix)
func (g *CodeGenerator) generateSetter(className string, node *ASTNode) error {
	// Create method name with "Set" prefix
	methodName := "Set" + toPascalCase(node.Name)

	// Build parameter (setter always has one parameter)
	paramName := "value"
	paramType := "interface{}"
	if node.Parameters != nil && len(node.Parameters) > 0 {
		paramName = node.Parameters[0].Name
		if node.Parameters[0].Type != nil {
			var err error
			paramType, err = g.generateType(node.Parameters[0].Type)
			if err != nil {
				return err
			}
		}
	}

	// Write setter method signature with receiver
	receiverVar := strings.ToLower(string(className[0]))
	g.writeLine(fmt.Sprintf("func (%s *%s) %s(%s %s) {", receiverVar, className, methodName, paramName, paramType))
	g.indent++

	// Set context for receiver var
	oldReceiverVar := g.currentReceiverVar
	g.currentReceiverVar = receiverVar
	defer func() {
		g.currentReceiverVar = oldReceiverVar
	}()

	// Generate setter body
	if node.Body != nil {
		if err := g.generateMethodBody(node.Body, receiverVar); err != nil {
			return err
		}
	}

	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// generateStaticMethod generates a static method as a package-level function
func (g *CodeGenerator) generateStaticMethod(className string, node *ASTNode) error {
	// Use camelCase for private static methods, PascalCase for public
	methodName := className + toPascalCase(node.Name)
	if isPrivate(node) {
		methodName = toCamelCase(className + toPascalCase(node.Name))
	}

	// Build parameter list
	params := []string{}
	if node.Parameters != nil {
		for _, param := range node.Parameters {
			paramName := param.Name
			paramType := "interface{}"
			if param.Type != nil {
				var err error
				paramType, err = g.generateType(param.Type)
				if err != nil {
					return err
				}
			}
			params = append(params, fmt.Sprintf("%s %s", paramName, paramType))
		}
	}

	// Determine return type
	returnType := ""
	if node.Type != nil {
		var err error
		returnType, err = g.generateType(node.Type)
		if err != nil {
			return err
		}
	}

	// Write function signature (no receiver for static methods)
	signature := fmt.Sprintf("func %s(%s)", methodName, strings.Join(params, ", "))
	if returnType != "" {
		signature += " " + returnType
	}
	g.writeLine(signature + " {")
	g.indent++

	// Set context for return type
	oldReturnType := g.currentFunctionReturnType
	g.currentFunctionReturnType = returnType
	defer func() {
		g.currentFunctionReturnType = oldReturnType
	}()

	// Generate method body (no receiver variable for static methods)
	if node.Body != nil && node.Body.Statements != nil {
		for _, stmt := range node.Body.Statements {
			if err := g.generateStatement(&stmt); err != nil {
				return err
			}
		}
	}

	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// generateMethodBody generates the body of a method, replacing "this" with receiver variable
func (g *CodeGenerator) generateMethodBody(bodyNode *ASTNode, receiverVar string) error {
	if bodyNode.Statements != nil {
		for _, stmt := range bodyNode.Statements {
			if err := g.generateMethodStatement(&stmt, receiverVar); err != nil {
				return err
			}
		}
	}
	return nil
}

// generateMethodStatement handles statements in methods, replacing "this" references
func (g *CodeGenerator) generateMethodStatement(node *ASTNode, receiverVar string) error {
	// For now, use the standard statement generation
	// We'll need to handle "this" in expressions
	return g.generateStatement(node)
}

// needsOptionalAccess checks if the AST uses optional chaining
func (g *CodeGenerator) needsOptionalAccess(node *ASTNode) bool {
	if node == nil {
		return false
	}

	if node.Kind == PropertyAccessExpression && node.QuestionDot {
		return true
	}

	// Check children recursively
	for _, child := range node.Children {
		if g.needsOptionalAccess(&child) {
			return true
		}
	}
	if node.Body != nil && g.needsOptionalAccess(node.Body) {
		return true
	}
	for _, stmt := range node.Statements {
		if g.needsOptionalAccess(&stmt) {
			return true
		}
	}
	for _, decl := range node.Declarations {
		if g.needsOptionalAccess(&decl) {
			return true
		}
	}
	if node.Initializer != nil && g.needsOptionalAccess(node.Initializer) {
		return true
	}

	return false
}

// needsNullishCoalesce checks if the AST uses nullish coalescing
func (g *CodeGenerator) needsNullishCoalesce(node *ASTNode) bool {
	if node == nil {
		return false
	}

	if node.Kind == BinaryExpression && node.Operator == QuestionQuestionToken {
		return true
	}

	// Check children recursively
	for _, child := range node.Children {
		if g.needsNullishCoalesce(&child) {
			return true
		}
	}
	if node.Body != nil && g.needsNullishCoalesce(node.Body) {
		return true
	}
	for _, stmt := range node.Statements {
		if g.needsNullishCoalesce(&stmt) {
			return true
		}
	}
	for _, decl := range node.Declarations {
		if g.needsNullishCoalesce(&decl) {
			return true
		}
	}
	if node.Initializer != nil && g.needsNullishCoalesce(node.Initializer) {
		return true
	}

	return false
}

// generateTupleType converts a TS tuple type to a Go struct
func (g *CodeGenerator) generateTupleType(node *ASTNode) (string, error) {
	// For now, generate an inline struct with numbered fields
	// In the future, this could generate a named struct for reusability
	fields := []string{}

	if node.Elements != nil {
		for i, element := range node.Elements {
			fieldType, err := g.generateType(&element)
			if err != nil {
				return "", err
			}
			fields = append(fields, fmt.Sprintf("Field%d %s", i, fieldType))
		}
	}

	if len(fields) == 0 {
		return "struct{}", nil
	}

	return fmt.Sprintf("struct{ %s }", strings.Join(fields, "; ")), nil
}

// generateUnionType creates a discriminated union struct for TypeScript union types
func (g *CodeGenerator) generateUnionType(typeName string, unionNode *ASTNode) error {
	// Generate the discriminated union struct
	g.writeLine(fmt.Sprintf("type %s struct {", typeName))
	g.indent++
	g.writeLine("Type  string")
	g.writeLine("Value interface{}")
	g.indent--
	g.writeLine("}")
	g.writeLine("")

	// Generate constructor functions for each union variant
	if unionNode.Types != nil {
		for _, unionType := range unionNode.Types {
			if err := g.generateUnionConstructor(typeName, &unionType); err != nil {
				return err
			}
		}
	}

	// Generate type guard methods
	if unionNode.Types != nil {
		for _, unionType := range unionNode.Types {
			if err := g.generateTypeGuard(typeName, &unionType); err != nil {
				return err
			}
		}
	}

	return nil
}

// generateUnionConstructor creates a constructor function for a union variant
func (g *CodeGenerator) generateUnionConstructor(unionTypeName string, variantType *ASTNode) error {
	variantName := g.getTypeName(variantType)
	constructorName := fmt.Sprintf("New%s%s", unionTypeName, variantName)

	g.writeLine(fmt.Sprintf("func %s(value %s) %s {", constructorName, g.getGoType(variantType), unionTypeName))
	g.indent++
	g.writeLine(fmt.Sprintf("return %s{Type: \"%s\", Value: value}", unionTypeName, variantName))
	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// generateTypeGuard creates a type guard method for a union variant
func (g *CodeGenerator) generateTypeGuard(unionTypeName string, variantType *ASTNode) error {
	variantName := g.getTypeName(variantType)
	methodName := fmt.Sprintf("Is%s", variantName)

	g.writeLine(fmt.Sprintf("func (u %s) %s() bool {", unionTypeName, methodName))
	g.indent++
	g.writeLine(fmt.Sprintf("return u.Type == \"%s\"", variantName))
	g.indent--
	g.writeLine("}")
	g.writeLine("")

	// Also generate a getter method
	getterName := fmt.Sprintf("As%s", variantName)
	g.writeLine(fmt.Sprintf("func (u %s) %s() %s {", unionTypeName, getterName, g.getGoType(variantType)))
	g.indent++
	g.writeLine(fmt.Sprintf("return u.Value.(%s)", g.getGoType(variantType)))
	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// getTypeName extracts a readable name from a type node
func (g *CodeGenerator) getTypeName(node *ASTNode) string {
	switch node.Kind {
	case TypeReference:
		if len(node.Children) > 0 && node.Children[0].Kind == Identifier {
			return toPascalCase(node.Children[0].Text)
		}
		if node.Text != "" {
			return toPascalCase(node.Text)
		}
	case LiteralType:
		// For literal types, use the string value
		if len(node.Children) > 0 && node.Children[0].Kind == StringLiteral {
			return toPascalCase(node.Children[0].Text)
		}
		return "Literal"
	case StringKeyword:
		return "String"
	case NumberKeyword:
		return "Number"
	case BooleanKeyword:
		return "Bool"
	default:
		return "Unknown"
	}
	return "Unknown"
}

// getGoType converts a type node to its Go equivalent
func (g *CodeGenerator) getGoType(node *ASTNode) string {
	switch node.Kind {
	case TypeReference:
		if len(node.Children) > 0 && node.Children[0].Kind == Identifier {
			return toPascalCase(node.Children[0].Text)
		}
		if node.Text != "" {
			return toPascalCase(node.Text)
		}
	case LiteralType:
		// String literals are still strings in Go
		return "string"
	case StringKeyword:
		return "string"
	case NumberKeyword:
		return "float64"
	case BooleanKeyword:
		return "bool"
	default:
		return "interface{}"
	}
	return "interface{}"
}

// generateFunction converts a TS function to a Go function
func (g *CodeGenerator) generateFunction(node *ASTNode) error {
	funcName := toPascalCase(node.Name)

	// Build parameter list
	params := []string{}
	if node.Parameters != nil {
		for _, param := range node.Parameters {
			paramName := param.Name
			paramType := "interface{}"
			if param.Type != nil {
				var err error
				paramType, err = g.generateType(param.Type)
				if err != nil {
					return err
				}
			}
			params = append(params, fmt.Sprintf("%s %s", paramName, paramType))
		}
	}

	// Determine return type
	returnType := ""
	if node.Type != nil {
		var err error
		returnType, err = g.generateType(node.Type)
		if err != nil {
			return err
		}
	}

	// Write function signature
	signature := fmt.Sprintf("func %s(%s)", funcName, strings.Join(params, ", "))
	if returnType != "" {
		signature += " " + returnType
	}
	g.writeLine(signature + " {")
	g.indent++

	// Set the current function return type for context-aware generation
	oldReturnType := g.currentFunctionReturnType
	g.currentFunctionReturnType = returnType
	defer func() { g.currentFunctionReturnType = oldReturnType }()

	// Generate function body
	if node.Body != nil {
		if err := g.generateBlock(node.Body); err != nil {
			return err
		}
	}

	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// generateBlock generates code for a block statement
func (g *CodeGenerator) generateBlock(node *ASTNode) error {
	if node.Statements != nil {
		for _, stmt := range node.Statements {
			if err := g.generateStatement(&stmt); err != nil {
				return err
			}
		}
	}
	return nil
}

// generateVariableStatement generates code for variable declarations
func (g *CodeGenerator) generateVariableStatement(node *ASTNode) error {
	if node.Declarations != nil {
		for _, decl := range node.Declarations {
			varName := decl.Name
			if decl.Initializer != nil {
				init, err := g.generateExpression(decl.Initializer)
				if err != nil {
					return err
				}
				g.writeLine(fmt.Sprintf("%s := %s", varName, init))
			} else if decl.Type != nil {
				varType, err := g.generateType(decl.Type)
				if err != nil {
					return err
				}
				g.writeLine(fmt.Sprintf("var %s %s", varName, varType))
			}
		}
	}
	return nil
}

// generateExpressionStatement generates code for expression statements
func (g *CodeGenerator) generateExpressionStatement(node *ASTNode) error {
	// ExpressionStatement has the expression in the Expression property
	var exprNode *ASTNode
	if node.Expression != nil {
		exprNode = node.Expression
	} else if len(node.Children) > 0 {
		exprNode = &node.Children[0]
	}

	if exprNode != nil {
		expr, err := g.generateExpression(exprNode)
		if err != nil {
			return err
		}
		g.writeLine(expr)
	}
	return nil
}

// generateReturnStatement generates code for return statements
func (g *CodeGenerator) generateReturnStatement(node *ASTNode) error {
	// Check for return expression in Expression property first, then Children
	var returnExpr *ASTNode
	if node.Expression != nil {
		returnExpr = node.Expression
	} else if len(node.Children) > 0 {
		returnExpr = &node.Children[0]
	}

	if returnExpr == nil {
		g.writeLine("return")
		return nil
	}

	// Generate the expression to return
	expr, err := g.generateExpression(returnExpr)
	if err != nil {
		return err
	}

	// If expression uses nullishCoalesce and we know the return type, add type assertion
	if strings.Contains(expr, "nullishCoalesce(") && g.currentFunctionReturnType != "" && g.currentFunctionReturnType != "interface{}" {
		// Add type assertion to match the function's return type
		expr = fmt.Sprintf("%s.(%s)", expr, g.currentFunctionReturnType)
	}

	g.writeLine(fmt.Sprintf("return %s", expr))
	return nil
}

// generateIfStatement generates code for if/else statements
func (g *CodeGenerator) generateIfStatement(node *ASTNode) error {
	// Get the condition expression (first child)
	if len(node.Children) == 0 {
		return fmt.Errorf("if statement missing condition")
	}

	condition, err := g.generateExpression(&node.Children[0])
	if err != nil {
		return fmt.Errorf("generating if condition: %w", err)
	}

	// Generate: if condition {
	g.writeLine(fmt.Sprintf("if %s {", condition))
	g.indent++

	// Generate the then block (second child)
	if len(node.Children) > 1 {
		thenBlock := &node.Children[1]
		if err := g.generateStatementBlock(thenBlock); err != nil {
			return fmt.Errorf("generating if then block: %w", err)
		}
	}

	g.indent--

	// Check for else clause (third child)
	if len(node.Children) > 2 {
		elseBlock := &node.Children[2]

		// Check if it's an else-if (IfStatement) or else block
		if elseBlock.Kind == IfStatement {
			// Generate: } else if ... {
			// We need to handle else-if inline, so write the closing brace and else on same line
			for i := 0; i < g.indent; i++ {
				g.output.WriteString("\t")
			}
			g.output.WriteString("} else ")

			// For else-if, we need to recursively generate the if condition inline
			// Get the condition from the else-if
			if len(elseBlock.Children) > 0 {
				elseCond, err := g.generateExpression(&elseBlock.Children[0])
				if err != nil {
					return fmt.Errorf("generating else-if condition: %w", err)
				}
				g.output.WriteString(fmt.Sprintf("if %s {\n", elseCond))

				// Generate the else-if then block
				g.indent++
				if len(elseBlock.Children) > 1 {
					if err := g.generateStatementBlock(&elseBlock.Children[1]); err != nil {
						return fmt.Errorf("generating else-if block: %w", err)
					}
				}
				g.indent--

				// Check for further else/else-if
				if len(elseBlock.Children) > 2 {
					// Recursively handle more else-if or final else
					furtherElse := &elseBlock.Children[2]
					if furtherElse.Kind == IfStatement {
						// Continue the chain
						for i := 0; i < g.indent; i++ {
							g.output.WriteString("\t")
						}
						g.output.WriteString("} else ")
						// This gets complex, let's simplify for now
						g.output.WriteString("{\n")
						g.indent++
						if err := g.generateStatement(furtherElse); err != nil {
							return err
						}
						g.indent--
						g.writeLine("}")
					} else {
						// Final else
						g.writeLine("} else {")
						g.indent++
						if err := g.generateStatementBlock(furtherElse); err != nil {
							return fmt.Errorf("generating final else: %w", err)
						}
						g.indent--
						g.writeLine("}")
					}
				} else {
					g.writeLine("}")
				}
			}
			return nil
		} else {
			// Regular else block
			g.writeLine("} else {")
			g.indent++

			if err := g.generateStatementBlock(elseBlock); err != nil {
				return fmt.Errorf("generating else block: %w", err)
			}

			g.indent--
			g.writeLine("}")
		}
	} else {
		// No else clause, just close the if
		g.writeLine("}")
	}

	return nil
}

// generateStatementBlock generates code for a block of statements
func (g *CodeGenerator) generateStatementBlock(node *ASTNode) error {
	// If it's a Block node, process its statements
	if node.Kind == Block {
		if node.Statements != nil {
			for _, stmt := range node.Statements {
				if err := g.generateStatement(&stmt); err != nil {
					return err
				}
			}
		}
	} else {
		// Single statement (no braces in TypeScript)
		if err := g.generateStatement(node); err != nil {
			return err
		}
	}
	return nil
}

// generateForStatement generates code for traditional for loops
// TypeScript: for (let i = 0; i < 10; i++) { }
// Go: for i := 0; i < 10; i++ { }
func (g *CodeGenerator) generateForStatement(node *ASTNode) error {
	// Traditional for loop has:
	// Initializer property (VariableDeclarationList)
	// Children[0] = condition (BinaryExpression)
	// Children[1] = incrementor (PostfixUnaryExpression)
	// Children[2] = body (Block)

	if len(node.Children) < 3 {
		return fmt.Errorf("for statement requires 3 children (condition, increment, body)")
	}

	// Generate initializer
	var initPart string
	if node.Initializer != nil {
		// Initializer is a VariableDeclarationList with Children
		if len(node.Initializer.Children) > 0 {
			decl := &node.Initializer.Children[0] // VariableDeclaration
			if decl.Name != "" && decl.Initializer != nil {
				init, err := g.generateExpression(decl.Initializer)
				if err != nil {
					return fmt.Errorf("generating for init: %w", err)
				}
				initPart = fmt.Sprintf("%s := %s", decl.Name, init)
			}
		}
	}

	// Generate condition
	condPart, err := g.generateExpression(&node.Children[0])
	if err != nil {
		return fmt.Errorf("generating for condition: %w", err)
	}

	// Generate incrementor
	incPart, err := g.generateExpression(&node.Children[1])
	if err != nil {
		return fmt.Errorf("generating for increment: %w", err)
	}

	// Generate: for init; condition; increment {
	g.writeLine(fmt.Sprintf("for %s; %s; %s {", initPart, condPart, incPart))
	g.indent++

	// Generate body
	bodyNode := &node.Children[2]
	if err := g.generateStatementBlock(bodyNode); err != nil {
		return fmt.Errorf("generating for body: %w", err)
	}

	g.indent--
	g.writeLine("}")

	return nil
}

// generateForOfStatement generates code for for...of loops
// TypeScript: for (const item of array) { }
// Go: for _, item := range array { }
func (g *CodeGenerator) generateForOfStatement(node *ASTNode) error {
	// for...of structure:
	// Initializer = VariableDeclarationList with variable name
	// Expression = iterable expression (array/string/etc)
	// Body or Children[last] = body Block

	// Get variable name from initializer
	varName := ""
	if node.Initializer != nil && len(node.Initializer.Children) > 0 {
		varDecl := &node.Initializer.Children[0] // VariableDeclaration
		varName = varDecl.Name
	}

	if varName == "" {
		return fmt.Errorf("for-of: could not determine variable name")
	}

	// Get iterable - check both Expression property and Children[0]
	var iterable string
	var err error
	if node.Expression != nil {
		iterable, err = g.generateExpression(node.Expression)
	} else if len(node.Children) > 0 {
		iterable, err = g.generateExpression(&node.Children[0])
	} else {
		return fmt.Errorf("for-of: no iterable expression found")
	}
	if err != nil {
		return fmt.Errorf("generating for-of iterable: %w", err)
	}

	// Generate: for _, item := range array {
	g.writeLine(fmt.Sprintf("for _, %s := range %s {", varName, iterable))
	g.indent++

	// Generate body - check both Body property and Children (last child is usually the body)
	var bodyNode *ASTNode
	if node.Body != nil {
		bodyNode = node.Body
	} else if len(node.Children) > 0 {
		bodyNode = &node.Children[len(node.Children)-1]
	} else {
		return fmt.Errorf("for-of: no body found")
	}

	if err := g.generateStatementBlock(bodyNode); err != nil {
		return fmt.Errorf("generating for-of body: %w", err)
	}

	g.indent--
	g.writeLine("}")

	return nil
}

// generateForInStatement generates code for for...in loops
// TypeScript: for (const key in object) { }
// Go: for key := range object { }
func (g *CodeGenerator) generateForInStatement(node *ASTNode) error {
	// for...in structure same as for...of
	// Initializer = VariableDeclarationList with variable name
	// Expression = object expression
	// Body or Children[last] = body Block

	// Get variable name from initializer
	varName := ""
	if node.Initializer != nil && len(node.Initializer.Children) > 0 {
		varDecl := &node.Initializer.Children[0]
		varName = varDecl.Name
	}

	if varName == "" {
		return fmt.Errorf("for-in: could not determine variable name")
	}

	// Get object - check both Expression property and Children
	var object string
	var err error
	if node.Expression != nil {
		object, err = g.generateExpression(node.Expression)
	} else if len(node.Children) > 0 {
		object, err = g.generateExpression(&node.Children[0])
	} else {
		return fmt.Errorf("for-in: no object expression found")
	}
	if err != nil {
		return fmt.Errorf("generating for-in object: %w", err)
	}

	// Generate: for key := range object {
	g.writeLine(fmt.Sprintf("for %s := range %s {", varName, object))
	g.indent++

	// Generate body - check both Body property and Children
	var bodyNode *ASTNode
	if node.Body != nil {
		bodyNode = node.Body
	} else if len(node.Children) > 0 {
		bodyNode = &node.Children[len(node.Children)-1]
	} else {
		return fmt.Errorf("for-in: no body found")
	}

	if err := g.generateStatementBlock(bodyNode); err != nil {
		return fmt.Errorf("generating for-in body: %w", err)
	}

	g.indent--
	g.writeLine("}")

	return nil
}

// generateWhileStatement generates code for while loops
// TypeScript: while (condition) { }
// Go: for condition { }
func (g *CodeGenerator) generateWhileStatement(node *ASTNode) error {
	// while loop has:
	// Expression = condition
	// Body or Children[last] = body

	// Generate condition - check Expression property or Children[0]
	var condition string
	var err error
	if node.Expression != nil {
		condition, err = g.generateExpression(node.Expression)
	} else if len(node.Children) > 0 {
		condition, err = g.generateExpression(&node.Children[0])
	} else {
		return fmt.Errorf("while statement: no condition found")
	}
	if err != nil {
		return fmt.Errorf("generating while condition: %w", err)
	}

	// Go uses 'for' for while loops
	g.writeLine(fmt.Sprintf("for %s {", condition))
	g.indent++

	// Generate body - check Body property or Children[0] (only 1 child - the body block)
	var bodyNode *ASTNode
	if node.Body != nil {
		bodyNode = node.Body
	} else if len(node.Children) > 0 {
		bodyNode = &node.Children[0]
	} else {
		return fmt.Errorf("while statement: no body found")
	}

	if err := g.generateStatementBlock(bodyNode); err != nil {
		return fmt.Errorf("generating while body: %w", err)
	}

	g.indent--
	g.writeLine("}")

	return nil
}

// generateBreakStatement generates code for break statements
func (g *CodeGenerator) generateBreakStatement(node *ASTNode) error {
	// Check if there's a label (Children[0] would be the label identifier)
	if len(node.Children) > 0 {
		label := node.Children[0].Text
		g.writeLine(fmt.Sprintf("break %s", label))
	} else {
		g.writeLine("break")
	}
	return nil
}

// generateContinueStatement generates code for continue statements
func (g *CodeGenerator) generateContinueStatement(node *ASTNode) error {
	// Check if there's a label
	if len(node.Children) > 0 {
		label := node.Children[0].Text
		g.writeLine(fmt.Sprintf("continue %s", label))
	} else {
		g.writeLine("continue")
	}
	return nil
}

// generateSwitchStatement generates code for switch statements
// TypeScript: switch (expr) { case value: statements; break; default: statements; }
// Go: switch expr { case value: statements; default: statements; }
func (g *CodeGenerator) generateSwitchStatement(node *ASTNode) error {
	// Switch statement structure:
	// Children[0] = expression being switched on
	// Children[1] = CaseBlock containing CaseClause and DefaultClause nodes

	if len(node.Children) < 2 {
		return fmt.Errorf("switch statement requires 2 children (expression, case block)")
	}

	// Generate the switch expression
	switchExpr, err := g.generateExpression(&node.Children[0])
	if err != nil {
		return fmt.Errorf("generating switch expression: %w", err)
	}

	g.writeLine(fmt.Sprintf("switch %s {", switchExpr))

	// Process the case block
	caseBlock := &node.Children[1]
	if caseBlock.Kind != CaseBlock {
		return fmt.Errorf("expected CaseBlock, got %s", caseBlock.Kind)
	}

	// Generate each case/default clause
	for _, clause := range caseBlock.Children {
		if clause.Kind == CaseClause {
			if err := g.generateCaseClause(&clause); err != nil {
				return fmt.Errorf("generating case clause: %w", err)
			}
		} else if clause.Kind == DefaultClause {
			if err := g.generateDefaultClause(&clause); err != nil {
				return fmt.Errorf("generating default clause: %w", err)
			}
		}
	}

	g.writeLine("}")
	return nil
}

// generateCaseClause generates code for a case clause in a switch statement
func (g *CodeGenerator) generateCaseClause(node *ASTNode) error {
	// CaseClause structure:
	// Children[0] = case expression (the value to match)
	// Statements property contains the statements to execute

	if len(node.Children) < 1 {
		return fmt.Errorf("case clause requires at least 1 child (case expression)")
	}

	// Generate the case value
	caseValue, err := g.generateExpression(&node.Children[0])
	if err != nil {
		return fmt.Errorf("generating case value: %w", err)
	}

	g.writeLine(fmt.Sprintf("case %s:", caseValue))
	g.indent++

	// Generate the statements in this case
	if node.Statements != nil {
		for _, stmt := range node.Statements {
			if err := g.generateStatement(&stmt); err != nil {
				return fmt.Errorf("generating case statement: %w", err)
			}
		}
	}

	g.indent--
	return nil
}

// generateDefaultClause generates code for the default clause in a switch statement
func (g *CodeGenerator) generateDefaultClause(node *ASTNode) error {
	// DefaultClause has only a Statements property

	g.writeLine("default:")
	g.indent++

	// Generate the statements in the default case
	if node.Statements != nil {
		for _, stmt := range node.Statements {
			if err := g.generateStatement(&stmt); err != nil {
				return fmt.Errorf("generating default statement: %w", err)
			}
		}
	}

	g.indent--
	return nil
}

// generateExpression generates code for an expression
func (g *CodeGenerator) generateExpression(node *ASTNode) (string, error) {
	switch node.Kind {
	case CallExpression:
		return g.generateCallExpression(node)
	case NewExpression:
		return g.generateNewExpression(node)
	case ThisKeyword:
		// Replace "this" with the current receiver variable
		if g.currentReceiverVar != "" {
			return g.currentReceiverVar, nil
		}
		return "this", nil
	case Identifier:
		return node.Text, nil
	case StringLiteral:
		return fmt.Sprintf(`"%s"`, node.Text), nil
	case NumericLiteral:
		return node.Text, nil
	case "FirstLiteralToken":
		// Numeric literal
		return node.Text, nil
	case "TrueKeyword":
		return "true", nil
	case "FalseKeyword":
		return "false", nil
	case BinaryExpression:
		return g.generateBinaryExpression(node)
	case PropertyAccessExpression:
		return g.generatePropertyAccess(node)
	case ConditionalExpression:
		return g.generateConditionalExpression(node)
	case ArrowFunction:
		return g.generateArrowFunction(node)
	case TemplateExpression:
		return g.generateTemplateExpression(node)
	case "PostfixUnaryExpression":
		return g.generatePostfixUnaryExpression(node)
	case "PrefixUnaryExpression":
		return g.generatePrefixUnaryExpression(node)
	case "ObjectLiteralExpression":
		return g.generateObjectLiteral(node)
	case "ArrayLiteralExpression":
		return g.generateArrayLiteral(node)
	default:
		return "/* unsupported expression */", nil
	}
}

// generateConditionalExpression generates code for ternary operator (condition ? true : false)
// TypeScript: condition ? trueExpr : falseExpr
// Go: uses if-else since Go doesn't have ternary
func (g *CodeGenerator) generateConditionalExpression(node *ASTNode) (string, error) {
	// TypeScript ConditionalExpression structure may vary
	// Try to get condition, whenTrue, whenFalse from Children array
	if len(node.Children) < 3 {
		// Fallback: maybe it's structured differently
		return "", fmt.Errorf("conditional expression has insufficient children: %d", len(node.Children))
	}

	// Parse children - TypeScript typically stores as:
	// Children[0] = condition
	// Children[1] = whenTrue (after '?')
	// Children[2] = whenFalse (after ':')
	condition, err := g.generateExpression(&node.Children[0])
	if err != nil {
		return "", fmt.Errorf("generating ternary condition: %w", err)
	}

	whenTrue, err := g.generateExpression(&node.Children[1])
	if err != nil {
		return "", fmt.Errorf("generating ternary whenTrue: %w", err)
	}

	whenFalse, err := g.generateExpression(&node.Children[2])
	if err != nil {
		return "", fmt.Errorf("generating ternary whenFalse: %w", err)
	}

	// Generate Go code using an immediately-invoked function
	// This allows ternary to be used as an expression
	result := fmt.Sprintf("func() interface{} { if %s { return %s } else { return %s } }()",
		condition, whenTrue, whenFalse)

	return result, nil
}

// generateArrowFunction generates code for arrow functions
// TypeScript: (a, b) => a + b or (x) => { return x * 2; }
// Go: func(a, b) return_type { return a + b }
func (g *CodeGenerator) generateArrowFunction(node *ASTNode) (string, error) {
	// Build parameter list
	params := []string{}
	if node.Parameters != nil {
		for _, param := range node.Parameters {
			paramName := param.Name
			paramType := "interface{}"
			if param.Type != nil {
				var err error
				paramType, err = g.generateType(param.Type)
				if err != nil {
					return "", fmt.Errorf("generating arrow function param type: %w", err)
				}
			}
			params = append(params, fmt.Sprintf("%s %s", paramName, paramType))
		}
	}

	// Determine return type
	returnType := "interface{}"
	if node.Type != nil {
		var err error
		returnType, err = g.generateType(node.Type)
		if err != nil {
			return "", fmt.Errorf("generating arrow function return type: %w", err)
		}
	}
	// For void return type, use empty string
	if returnType == "" {
		returnType = "interface{}"
	}

	// Generate function body
	var bodyCode string
	if node.Body != nil {
		if node.Body.Kind == Block {
			// Block body with statements
			// For arrow functions with blocks, we need to capture the generated code
			// Save current output and create temporary buffer
			oldOutput := g.output
			g.output = strings.Builder{}
			oldIndent := g.indent
			g.indent = 0

			// Generate the block statements
			if node.Body.Statements != nil {
				for _, stmt := range node.Body.Statements {
					if err := g.generateStatement(&stmt); err != nil {
						g.output = oldOutput
						g.indent = oldIndent
						return "", fmt.Errorf("generating arrow function block statement: %w", err)
					}
				}
			}

			// Capture generated code and restore output
			bodyCode = strings.TrimSpace(g.output.String())
			g.output = oldOutput
			g.indent = oldIndent
		} else {
			// Expression body: expr (implicit return)
			expr, err := g.generateExpression(node.Body)
			if err != nil {
				return "", fmt.Errorf("generating arrow function expression: %w", err)
			}
			bodyCode = "return " + expr
		}
	}

	// Generate: func(params) returnType { body }
	if strings.Contains(bodyCode, "\n") {
		// Multi-line body
		return fmt.Sprintf("func(%s) %s {\n\t%s\n}", strings.Join(params, ", "), returnType, bodyCode), nil
	}
	return fmt.Sprintf("func(%s) %s { %s }", strings.Join(params, ", "), returnType, bodyCode), nil
}

// generateTemplateExpression generates code for template literals
// TypeScript: `Hello ${name}!`
// Go: "Hello " + name + "!"
func (g *CodeGenerator) generateTemplateExpression(node *ASTNode) (string, error) {
	// Template expressions have:
	// - Head (TemplateHead) - the initial string part
	// - TemplateSpans in Children - alternating expressions and string parts

	parts := []string{}

	// Add the head if it exists
	if node.Head != nil && node.Head.Text != "" {
		// Escape special characters in template strings
		text := strings.ReplaceAll(node.Head.Text, "\n", "\\n")
		text = strings.ReplaceAll(text, "\t", "\\t")
		text = strings.ReplaceAll(text, "\"", "\\\"")
		parts = append(parts, fmt.Sprintf(`"%s"`, text))
	}

	// Process template spans (expression + text pairs)
	if node.Children != nil {
		for _, child := range node.Children {
			if child.Kind == "TemplateSpan" {
				// TemplateSpan has:
				// - Expression (the ${...} part)
				// - Literal (the text after the expression)

				// Generate the expression
				if child.Expression != nil {
					expr, err := g.generateExpression(child.Expression)
					if err != nil {
						return "", fmt.Errorf("generating template expression: %w", err)
					}
					// Convert expression to string using fmt.Sprint
					parts = append(parts, fmt.Sprintf("fmt.Sprint(%s)", expr))
				}

				// Add the literal text part if it exists
				if child.Literal != nil && child.Literal.Text != "" {
					// Escape special characters
					text := strings.ReplaceAll(child.Literal.Text, "\n", "\\n")
					text = strings.ReplaceAll(text, "\t", "\\t")
					text = strings.ReplaceAll(text, "\"", "\\\"")
					parts = append(parts, fmt.Sprintf(`"%s"`, text))
				}
			}
		}
	}

	// If only one part, return it directly
	if len(parts) == 1 {
		return parts[0], nil
	}

	// Join parts with + operator
	if len(parts) == 0 {
		return `""`, nil
	}

	return strings.Join(parts, " + "), nil
}

// generateObjectLiteral generates code for object literals
func (g *CodeGenerator) generateObjectLiteral(node *ASTNode) (string, error) {
	if node.Properties == nil || len(node.Properties) == 0 {
		return "{}", nil
	}

	// Generate struct field initializations
	fields := []string{}
	for _, prop := range node.Properties {
		if prop.Kind == "PropertyAssignment" {
			// Property name
			propName := prop.Name

			// Property value
			var value string
			var err error
			if prop.Initializer != nil {
				value, err = g.generateExpression(prop.Initializer)
				if err != nil {
					return "", err
				}
			}

			fields = append(fields, fmt.Sprintf("%s: %s", toPascalCase(propName), value))
		}
	}

	// Return as a struct initialization
	// If we have a known return type, prepend it to the literal
	if g.currentFunctionReturnType != "" && g.currentFunctionReturnType != "interface{}" {
		return fmt.Sprintf("%s{%s}", g.currentFunctionReturnType, strings.Join(fields, ", ")), nil
	}
	return fmt.Sprintf("{%s}", strings.Join(fields, ", ")), nil
}

// generateArrayLiteral generates code for array literals
func (g *CodeGenerator) generateArrayLiteral(node *ASTNode) (string, error) {
	if node.Elements == nil || len(node.Elements) == 0 {
		return "[]interface{}{}", nil
	}

	// Generate array elements
	elements := []string{}
	for _, elem := range node.Elements {
		value, err := g.generateExpression(&elem)
		if err != nil {
			return "", err
		}
		elements = append(elements, value)
	}

	// For tuple types used as inline values, generate as struct
	// This is a heuristic - if used in a context expecting a tuple type,
	// we return it as an inline struct initialization
	// For now, return as array slice
	return fmt.Sprintf("[]interface{}{%s}", strings.Join(elements, ", ")), nil
}

// generatePropertyAccess generates code for property access (e.g., person.name or Color.Red)
func (g *CodeGenerator) generatePropertyAccess(node *ASTNode) (string, error) {
	// PropertyAccessExpression has the object in the Expression property
	var objectNode *ASTNode
	if node.Expression != nil {
		objectNode = node.Expression
	} else if len(node.Children) > 0 {
		objectNode = &node.Children[0]
	} else {
		return "", fmt.Errorf("invalid property access: no object found")
	}

	object, err := g.generateExpression(objectNode)
	if err != nil {
		return "", err
	}

	property := node.Name
	if property == "" {
		return "", fmt.Errorf("property access missing name")
	}

	// For this.field, check if field is private using currentClassMembers
	if objectNode.Kind == ThisKeyword && g.currentClassMembers != nil {
		if g.currentClassMembers[property] {
			// Private member - use camelCase
			property = toCamelCase(property)
		} else {
			// Public member - use PascalCase
			property = toPascalCase(property)
		}
	} else {
		// Convert to PascalCase for struct field access
		property = toPascalCase(property)
	}

	// Handle optional chaining (?.) - generate nil-safe access
	if node.QuestionDot {
		return fmt.Sprintf("optionalAccess(%s, \"%s\")", object, property), nil
	}

	// Check if this is a static member access (e.g., Person.species -> PersonSpecies)
	// or enum member access (e.g., Color.Red -> ColorRed)
	// This is a heuristic: if the object is a simple identifier starting with uppercase,
	// it might be a class/enum type with static members
	if objectNode.Kind == Identifier {
		objectName := objectNode.Text
		if len(objectName) > 0 && unicode.IsUpper(rune(objectName[0])) {
			// Static member access - convert ClassName.member to ClassNameMember
			return fmt.Sprintf("%s%s", toPascalCase(objectName), property), nil
		}
	}

	return fmt.Sprintf("%s.%s", object, property), nil
}

// generateCallExpression generates code for function calls
func (g *CodeGenerator) generateCallExpression(node *ASTNode) (string, error) {
	// CallExpression has the function in Expression property, arguments in Children
	var funcNode *ASTNode
	if node.Expression != nil {
		funcNode = node.Expression
	} else if len(node.Children) > 0 {
		funcNode = &node.Children[0]
	} else {
		return "()", nil
	}

	// Check if function is PropertyAccessExpression (e.g., console.log)
	if funcNode.Kind == "PropertyAccessExpression" {
		// PropertyAccessExpression has object in Expression property
		var objNode *ASTNode
		if funcNode.Expression != nil {
			objNode = funcNode.Expression
		} else if len(funcNode.Children) > 0 {
			objNode = &funcNode.Children[0]
		}

		if objNode != nil {
			obj := objNode.Text
			prop := funcNode.Name

			// Special case for console.log
			if obj == "console" && prop == "log" {
				args := []string{}
				// Arguments are in node.Children
				for _, arg := range node.Children {
					argStr, err := g.generateExpression(&arg)
					if err != nil {
						return "", err
					}
					args = append(args, argStr)
				}
				return fmt.Sprintf("fmt.Println(%s)", strings.Join(args, ", ")), nil
			}
		}
	}

	// Generic function call
	funcExpr, err := g.generateExpression(funcNode)
	if err != nil {
		return "", err
	}

	// Convert function names to PascalCase only if they start with uppercase
	// (i.e., they're exported functions, not local variables)
	if funcNode.Kind == Identifier && len(funcExpr) > 0 {
		// Only PascalCase if the original starts with uppercase (exported function)
		firstChar := funcExpr[0]
		if firstChar >= 'A' && firstChar <= 'Z' {
			funcExpr = toPascalCase(funcExpr)
		}
		// Otherwise keep it as-is (local variable or parameter)
	}

	args := []string{}
	// Arguments are in node.Children
	for _, arg := range node.Children {
		argStr, err := g.generateExpression(&arg)
		if err != nil {
			return "", err
		}
		args = append(args, argStr)
	}

	return fmt.Sprintf("%s(%s)", funcExpr, strings.Join(args, ", ")), nil
}

// generateNewExpression generates code for "new ClassName()" expressions
func (g *CodeGenerator) generateNewExpression(node *ASTNode) (string, error) {
	// Get the class name from Expression property
	className := ""
	if node.Expression != nil {
		if node.Expression.Kind == Identifier {
			className = toPascalCase(node.Expression.Text)
		} else {
			classExpr, err := g.generateExpression(node.Expression)
			if err != nil {
				return "", err
			}
			className = classExpr
		}
	} else if len(node.Children) > 0 {
		// Fallback to children if Expression not available
		if node.Children[0].Kind == Identifier {
			className = toPascalCase(node.Children[0].Text)
		} else {
			classExpr, err := g.generateExpression(&node.Children[0])
			if err != nil {
				return "", err
			}
			className = classExpr
		}
	} else {
		return "", fmt.Errorf("invalid new expression: no class name found")
	}

	// Get arguments from children (all children are arguments)
	args := []string{}
	for _, arg := range node.Children {
		argStr, err := g.generateExpression(&arg)
		if err != nil {
			return "", err
		}
		args = append(args, argStr)
	}

	// Generate New* constructor call
	return fmt.Sprintf("New%s(%s)", className, strings.Join(args, ", ")), nil
}

// generateBinaryExpression generates code for binary operations
func (g *CodeGenerator) generateBinaryExpression(node *ASTNode) (string, error) {
	if len(node.Children) < 2 {
		return "", fmt.Errorf("invalid binary expression")
	}

	left, err := g.generateExpression(&node.Children[0])
	if err != nil {
		return "", err
	}

	right, err := g.generateExpression(&node.Children[len(node.Children)-1])
	if err != nil {
		return "", err
	}

	// Check for nullish coalescing (??)
	if node.Operator == QuestionQuestionToken {
		// Use helper function for nullish coalescing
		// Note: Returns interface{} - type assertions may be needed based on usage context
		return fmt.Sprintf("nullishCoalesce(%s, %s)", left, right), nil
	}

	// Map operator from AST constant to Go operator
	operator := g.mapOperator(node.Operator)

	// If operator is generic like "FirstBinaryOperator", try to determine from middle child token
	if (operator == "" || node.Operator == "FirstBinaryOperator") && len(node.Children) >= 3 {
		middleChild := &node.Children[1]
		if middleChild.Text != "" {
			operator = middleChild.Text
		}
		// Also check for common operator kinds
		if middleChild.Kind == "LessThanToken" || middleChild.Kind == "FirstBinaryOperator" {
			operator = "<"
		} else if middleChild.Kind == "GreaterThanToken" {
			operator = ">"
		} else if middleChild.Kind == "LessThanEqualsToken" {
			operator = "<="
		} else if middleChild.Kind == "GreaterThanEqualsToken" {
			operator = ">="
		}
	}

	// Default to plus if we couldn't determine the operator
	if operator == "" {
		operator = "+"
	}

	return fmt.Sprintf("%s %s %s", left, operator, right), nil
}

// generatePostfixUnaryExpression generates code for postfix unary expressions (i++, i--)
func (g *CodeGenerator) generatePostfixUnaryExpression(node *ASTNode) (string, error) {
	if len(node.Children) < 1 {
		return "", fmt.Errorf("invalid postfix unary expression")
	}

	operand, err := g.generateExpression(&node.Children[0])
	if err != nil {
		return "", err
	}

	// TypeScript has i++ and i--
	// The operator is usually in the node kind or we can detect from the syntax
	// Most common: i++ (increment) and i-- (decrement)
	// Assume ++ for now (we can enhance this later if needed)
	return operand + "++", nil
}

// generatePrefixUnaryExpression generates code for prefix unary expressions (++i, --i, !x, -x)
func (g *CodeGenerator) generatePrefixUnaryExpression(node *ASTNode) (string, error) {
	if len(node.Children) < 1 {
		return "", fmt.Errorf("invalid prefix unary expression")
	}

	operand, err := g.generateExpression(&node.Children[0])
	if err != nil {
		return "", err
	}

	// Determine the operator based on context
	// If operand is a number literal or the operator spans 1 char before operand, it's likely - or !
	// Otherwise it's likely ++ or --
	operatorLen := node.Children[0].Pos - node.Pos

	switch operatorLen {
	case 1:
		// Single character operator: -, +, !, ~
		// Check operand type to guess which one
		if node.Children[0].Kind == "TrueKeyword" || node.Children[0].Kind == "FalseKeyword" {
			return "!" + operand, nil
		}
		// Assume negation for numbers
		return "-" + operand, nil
	case 2:
		// Two character operator: ++, --
		return "++" + operand, nil
	default:
		// Default to negation
		return "-" + operand, nil
	}
}

// mapOperator maps AST operator constants to Go operators
func (g *CodeGenerator) mapOperator(op string) string {
	switch op {
	case "PlusToken":
		return "+"
	case "MinusToken":
		return "-"
	case "AsteriskToken":
		return "*"
	case "SlashToken":
		return "/"
	case "PercentToken":
		return "%"
	case "EqualsEqualsToken", "EqualsEqualsEqualsToken":
		return "=="
	case "ExclamationEqualsToken", "ExclamationEqualsEqualsToken":
		return "!="
	case "LessThanToken":
		return "<"
	case "LessThanEqualsToken":
		return "<="
	case "GreaterThanToken":
		return ">"
	case "GreaterThanEqualsToken":
		return ">="
	case "AmpersandAmpersandToken":
		return "&&"
	case "BarBarToken":
		return "||"
	case "AmpersandToken":
		return "&"
	case "BarToken":
		return "|"
	case "CaretToken":
		return "^"
	case "FirstAssignment":
		return "="
	default:
		return ""
	}
}

// generateType converts a TS type to a Go type
func (g *CodeGenerator) generateType(node *ASTNode) (string, error) {
	if node == nil {
		return "interface{}", nil
	}

	switch node.Kind {
	case StringKeyword:
		return "string", nil
	case NumberKeyword:
		return "float64", nil
	case BooleanKeyword:
		return "bool", nil
	case VoidKeyword:
		return "", nil // void functions have no return type in Go
	case ArrayType:
		if len(node.Children) > 0 {
			elemType, err := g.generateType(&node.Children[0])
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("[]%s", elemType), nil
		}
		return "[]interface{}", nil
	case TupleType:
		return g.generateTupleType(node)
	case TypeLiteral:
		// Inline object type -> struct
		return g.generateInlineStruct(node)
	case TypeReference:
		// Reference to a named type (interface, type alias, etc.)
		if len(node.Children) > 0 && node.Children[0].Kind == Identifier {
			return toPascalCase(node.Children[0].Text), nil
		}
		return "interface{}", nil
	case UnionType:
		// Union types are handled at the type alias level
		return "interface{}", nil
	default:
		return "interface{}", nil
	}
}

// generateInlineStruct converts a TypeLiteral to an inline struct definition
func (g *CodeGenerator) generateInlineStruct(node *ASTNode) (string, error) {
	if node.Members == nil || len(node.Members) == 0 {
		return "struct{}", nil
	}

	fields := []string{}
	for _, member := range node.Members {
		if member.Kind == PropertySignature {
			fieldName := toPascalCase(member.Name)
			fieldType, err := g.generateType(member.Type)
			if err != nil {
				return "", err
			}
			jsonTag := fmt.Sprintf("`json:\"%s\"`", member.Name)
			fields = append(fields, fmt.Sprintf("%s %s %s", fieldName, fieldType, jsonTag))
		}
	}

	return fmt.Sprintf("struct{ %s }", strings.Join(fields, "; ")), nil
}

// writeLine writes a line with proper indentation
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

// toCamelCase converts PascalCase to camelCase (for private members)
func toCamelCase(s string) string {
	if s == "" {
		return s
	}

	// If already starts with lowercase or underscore, keep as is (it's already private-style)
	runes := []rune(s)
	if unicode.IsLower(runes[0]) || runes[0] == '_' {
		return s
	}

	// Convert first character to lowercase
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// hasModifier checks if a node has a specific modifier keyword
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
	return hasModifier(node, PrivateKeyword)
}

// isStatic checks if a node has a static modifier
func isStatic(node *ASTNode) bool {
	return hasModifier(node, StaticKeyword)
}

// getPackageDeclaration returns the package declaration using module context
func (g *CodeGenerator) getPackageDeclaration() string {
	if g.isEntry {
		return "package main"
	}

	// Use reflection to avoid import cycle
	// This is safe because we set these fields from multipackage.go
	if g.module != nil {
		// Access PackageName field via type assertion
		type moduleInterface interface {
			GetPackageName() string
		}
		if m, ok := g.module.(moduleInterface); ok {
			return fmt.Sprintf("package %s", m.GetPackageName())
		}
		// Fallback: use reflection-style access
		return "package main"
	}

	return "package main"
}

// getModuleImports generates import block using module resolver
func (g *CodeGenerator) getModuleImports() string {
	if g.resolver == nil {
		return ""
	}

	// Use reflection to avoid import cycle
	type resolverInterface interface {
		GenerateImportBlock() (string, error)
	}

	if r, ok := g.resolver.(resolverInterface); ok {
		if block, err := r.GenerateImportBlock(); err == nil {
			return block
		}
	}

	return ""
}

// getGoSymbolName returns the Go symbol name using visibility rules
func (g *CodeGenerator) getGoSymbolName(tsName string) string {
	if g.visibility == nil {
		return tsName
	}

	// Use reflection to avoid import cycle
	type visibilityInterface interface {
		GetGoSymbolName(string) string
	}

	if v, ok := g.visibility.(visibilityInterface); ok {
		return v.GetGoSymbolName(tsName)
	}

	return tsName
}

// NewCodeGeneratorWithModule creates a code generator with module context
// This allows the code generator to use module information for imports, exports, and visibility
func NewCodeGeneratorWithModule(mod, resolver, visibility interface{}, isEntry bool) *CodeGenerator {
	gen := NewCodeGenerator()
	gen.module = mod
	gen.resolver = resolver
	gen.visibility = visibility
	gen.isEntry = isEntry
	return gen
}
