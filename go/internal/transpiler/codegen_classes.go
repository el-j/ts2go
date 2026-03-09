package transpiler

import (
	"fmt"
	"strings"
)

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
