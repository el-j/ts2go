package transpiler

import (
	"fmt"
	"strings"
	"unicode"
)

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
	case "NullKeyword":
		return "nil", nil
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
	case TypeOfExpression:
		return g.generateTypeOfExpression(node)
	case DeleteExpression:
		return g.generateDeleteExpression(node)
	case AwaitExpression:
		return g.generateAwaitExpression(node)
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
			// Save current output and create temporary buffer
			oldOutput := g.output.String()
			g.output.Reset()

			// Generate block statements
			if node.Body.Statements != nil {
				for _, stmt := range node.Body.Statements {
					if err := g.generateStatement(&stmt); err != nil {
						g.output.Reset()
						g.output.WriteString(oldOutput)
						return "", fmt.Errorf("generating arrow function statement: %w", err)
					}
				}
			}

			// Get generated code and restore output
			bodyCode = strings.TrimSpace(g.output.String())
			g.output.Reset()
			g.output.WriteString(oldOutput)
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
					g.trackImport("fmt")
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
		return "map[string]interface{}{}", nil
	}

	// Check if any properties are spread elements
	hasSpread := false
	for _, prop := range node.Properties {
		if prop.Kind == SpreadElement || prop.Kind == SpreadAssignment {
			hasSpread = true
			break
		}
	}

	// If there are spread elements, we need to merge maps
	if hasSpread {
		// Generate code to create and merge maps
		parts := []string{}
		currentProps := []string{}

		for _, prop := range node.Properties {
			if prop.Kind == SpreadElement || prop.Kind == SpreadAssignment {
				// If we have accumulated properties, add them as a map literal
				if len(currentProps) > 0 {
					parts = append(parts, fmt.Sprintf("map[string]interface{}{%s}", strings.Join(currentProps, ", ")))
					currentProps = []string{}
				}

				// Add the spread expression
				if prop.Expression != nil {
					expr, err := g.generateExpression(prop.Expression)
					if err != nil {
						return "", err
					}
					parts = append(parts, expr)
				}
			} else if prop.Kind == "PropertyAssignment" {
				// Regular property
				propName := prop.Name
				var value string
				var err error
				if prop.Initializer != nil {
					value, err = g.generateExpression(prop.Initializer)
					if err != nil {
						return "", err
					}
				}
				currentProps = append(currentProps, fmt.Sprintf(`"%s": %s`, propName, value))
			}
		}

		// Add any remaining properties
		if len(currentProps) > 0 {
			parts = append(parts, fmt.Sprintf("map[string]interface{}{%s}", strings.Join(currentProps, ", ")))
		}

		// If only one part, return it directly
		if len(parts) == 1 {
			return parts[0], nil
		}

		// Use a helper function to merge maps
		// For now, generate inline merge code
		return g.generateObjectMerge(parts), nil
	}

	// No spread elements - generate simple object literal
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

			fields = append(fields, fmt.Sprintf(`"%s": %s`, propName, value))
		}
	}

	// Return as a map literal
	return fmt.Sprintf("map[string]interface{}{%s}", strings.Join(fields, ", ")), nil
}

// generateObjectMerge generates code to merge multiple objects/maps
func (g *CodeGenerator) generateObjectMerge(parts []string) string {
	// Create a temporary variable for the merged result
	tempVar := fmt.Sprintf("_merge_%d", g.tempVarCounter)
	g.tempVarCounter++

	// Generate inline function that creates and merges maps
	mergeCode := fmt.Sprintf("func() map[string]interface{} { %s := make(map[string]interface{}); ", tempVar)

	for _, part := range parts {
		mergeCode += fmt.Sprintf("for k, v := range %s { %s[k] = v }; ", part, tempVar)
	}

	mergeCode += fmt.Sprintf("return %s }()", tempVar)

	return mergeCode
}

// generateArrayLiteral generates code for array literals
func (g *CodeGenerator) generateArrayLiteral(node *ASTNode) (string, error) {
	if node.Elements == nil || len(node.Elements) == 0 {
		return "[]interface{}{}", nil
	}

	// Check if any elements are spread elements
	hasSpread := false
	for _, elem := range node.Elements {
		if elem.Kind == SpreadElement {
			hasSpread = true
			break
		}
	}

	// If no spread elements, use simple literal syntax
	if !hasSpread {
		elements := []string{}
		for _, elem := range node.Elements {
			value, err := g.generateExpression(&elem)
			if err != nil {
				return "", err
			}
			elements = append(elements, value)
		}
		return fmt.Sprintf("[]interface{}{%s}", strings.Join(elements, ", ")), nil
	}

	// If there are spread elements, we need to use array.Concat or append
	// Track import for array runtime
	g.trackImport("github.com/ts2go/runtime/array")

	// Build list of slices to concatenate
	slices := []string{}
	currentLiteral := []string{}

	for _, elem := range node.Elements {
		if elem.Kind == SpreadElement {
			// If we have accumulated literal elements, add them as a slice
			if len(currentLiteral) > 0 {
				slices = append(slices, fmt.Sprintf("[]interface{}{%s}", strings.Join(currentLiteral, ", ")))
				currentLiteral = []string{}
			}

			// Add the spread expression
			if elem.Expression != nil {
				expr, err := g.generateExpression(elem.Expression)
				if err != nil {
					return "", err
				}
				slices = append(slices, expr)
			}
		} else {
			// Regular element - accumulate for literal
			value, err := g.generateExpression(&elem)
			if err != nil {
				return "", err
			}
			currentLiteral = append(currentLiteral, value)
		}
	}

	// Add any remaining literal elements
	if len(currentLiteral) > 0 {
		slices = append(slices, fmt.Sprintf("[]interface{}{%s}", strings.Join(currentLiteral, ", ")))
	}

	// If only one slice, return it directly
	if len(slices) == 1 {
		return slices[0], nil
	}

	// Use array.Concat to merge all slices
	return fmt.Sprintf("array.Concat(%s)", strings.Join(slices, ", ")), nil
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
				g.trackImport("fmt")
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

			// Special case for array methods (filter, map, reduce, etc.)
			arrayMethods := map[string]bool{
				"filter": true, "map": true, "reduce": true,
				"find": true, "findIndex": true, "some": true,
				"every": true, "includes": true, "indexOf": true,
				"forEach": true, "push": true, "pop": true,
				"shift": true, "unshift": true, "reverse": true,
				"slice": true, "concat": true, "join": true,
			}

			if arrayMethods[prop] {
				// Track import for array runtime
				g.trackImport("github.com/ts2go/runtime/array")

				// Generate object expression
				objExpr, err := g.generateExpression(objNode)
				if err != nil {
					return "", err
				}

				// Generate arguments
				args := []string{}
				for _, arg := range node.Children {
					argStr, err := g.generateExpression(&arg)
					if err != nil {
						return "", err
					}
					args = append(args, argStr)
				}

				// Map JavaScript method names to Go function names
				methodName := toPascalCase(prop)

				// Handle special cases
				if prop == "push" || prop == "pop" || prop == "shift" || prop == "unshift" {
					// These methods modify the array in place, need pointer
					return fmt.Sprintf("array.%s(&%s, %s)", methodName, objExpr, strings.Join(args, ", ")), nil
				} else {
					// Other methods return new slices
					if len(args) > 0 {
						return fmt.Sprintf("array.%s(%s, %s)", methodName, objExpr, strings.Join(args, ", ")), nil
					} else {
						return fmt.Sprintf("array.%s(%s)", methodName, objExpr), nil
					}
				}
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

	// Check for instanceof operator
	if node.Operator == "InstanceOfKeyword" || node.Operator == "instanceof" {
		// TypeScript: obj instanceof ClassName
		// Go: Check type using reflect or type switch
		return fmt.Sprintf("/* instanceof check: %s is %s */ reflect.TypeOf(%s).String() == reflect.TypeOf((*%s)(nil)).Elem().String()", left, right, left, right), nil
	}

	// Check for in operator
	if node.Operator == "InKeyword" || node.Operator == "in" {
		// TypeScript: "key" in obj
		// Go: Check map key existence - _, exists := map[key]
		return fmt.Sprintf("func() bool { _, exists := %s[%s]; return exists }()", right, left), nil
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

// generateElementAccess generates code for array/object element access
// TypeScript: array[index] or object["key"]
// Go: array[index] or object["key"]
func (g *CodeGenerator) generateElementAccess(node *ASTNode) (string, error) {
	if node.Expression == nil {
		return "", fmt.Errorf("element access missing expression")
	}

	expr, err := g.generateExpression(node.Expression)
	if err != nil {
		return "", fmt.Errorf("generating element access expression: %w", err)
	}

	// The argument expression is the index/key
	if len(node.Children) > 0 && node.Children[0].Kind != "undefined" {
		index, err := g.generateExpression(&node.Children[0])
		if err != nil {
			return "", fmt.Errorf("generating element access index: %w", err)
		}
		return fmt.Sprintf("%s[%s]", expr, index), nil
	}

	return "", fmt.Errorf("element access missing index")
}

// generateSpreadElement generates code for spread operator
// TypeScript: ...args
// Go: args... (in function calls) or append/copy patterns
func (g *CodeGenerator) generateSpreadElement(node *ASTNode) (string, error) {
	if node.Expression == nil {
		return "", fmt.Errorf("spread element missing expression")
	}

	expr, err := g.generateExpression(node.Expression)
	if err != nil {
		return "", fmt.Errorf("generating spread expression: %w", err)
	}

	// In Go, spread is used differently depending on context
	// For now, append ... suffix for variadic calls
	return expr + "...", nil
}

// generateType converts a TS type to a Go type

// generateTypeOfExpression generates code for typeof operator
// TypeScript: typeof x
// Go: reflect.TypeOf(x).String() - requires import "reflect"
func (g *CodeGenerator) generateTypeOfExpression(node *ASTNode) (string, error) {
	if node.Expression == nil {
		return "", fmt.Errorf("typeof expression missing operand")
	}

	expr, err := g.generateExpression(node.Expression)
	if err != nil {
		return "", fmt.Errorf("generating typeof operand: %w", err)
	}

	// Generate Go code using reflect.TypeOf
	// Note: This returns the Go type, not JavaScript type strings
	return fmt.Sprintf("reflect.TypeOf(%s).String()", expr), nil
}

// generateDeleteExpression generates code for delete operator
// TypeScript: delete obj.prop or delete obj["key"]
// Go: No direct equivalent - generates delete() for maps, comment for others
func (g *CodeGenerator) generateDeleteExpression(node *ASTNode) (string, error) {
	if node.Expression == nil {
		return "", fmt.Errorf("delete expression missing operand")
	}

	// Check if it's a property access or element access
	expr := node.Expression
	if expr.Kind == PropertyAccessExpression {
		// delete obj.prop
		// Generate: delete(map, "key") - only works for maps
		obj, err := g.generateExpression(expr.Expression)
		if err != nil {
			return "", err
		}
		propName := expr.Name
		return fmt.Sprintf("delete(%s, \"%s\")", obj, propName), nil
	} else if expr.Kind == "ElementAccessExpression" {
		// delete obj["key"]
		obj, err := g.generateExpression(expr.Expression)
		if err != nil {
			return "", err
		}
		if len(expr.Children) > 0 {
			key, err := g.generateExpression(&expr.Children[0])
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("delete(%s, %s)", obj, key), nil
		}
	}

	// Fallback - delete on a simple expression
	operand, err := g.generateExpression(expr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/* delete %s - not supported in Go */", operand), nil
}

// generateAwaitExpression generates code for await expression
// TypeScript: await promise
// Go: <-channel (receiving from a channel)
func (g *CodeGenerator) generateAwaitExpression(node *ASTNode) (string, error) {
	if node.Expression == nil {
		return "", fmt.Errorf("await expression missing operand")
	}

	// Generate the expression being awaited (should be a Promise/channel)
	expr, err := g.generateExpression(node.Expression)
	if err != nil {
		return "", fmt.Errorf("generating await operand: %w", err)
	}

	// In Go, await translates to receiving from a channel
	// The expression should be a function call that returns a channel
	return fmt.Sprintf("(<-%s)", expr), nil
}
