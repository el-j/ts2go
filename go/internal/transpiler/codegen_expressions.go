package transpiler

import (
	"fmt"
	"strings"
)

// generateExpression generates code for expressions
func (g *CodeGenerator) generateExpression(node *ASTNode) (string, error) {
	if node == nil {
		return "", nil
	}

	switch node.Kind {
	case ThisKeyword:
		if g.currentReceiverVar != "" {
			return g.currentReceiverVar, nil
		}
		return "this", nil

	case Identifier:
		name := node.Text
		if name == "undefined" {
			return "nil", nil
		}
		if mappedName, ok := g.declaredFunctions[name]; ok {
			return mappedName, nil
		}
		return name, nil

	case StringLiteral, "NoSubstitutionTemplateLiteral":
		return fmt.Sprintf(`"%s"`, node.Text), nil

	case NumericLiteral, "FirstLiteralToken":
		return node.Text, nil

	case TrueKeyword:
		return "true", nil

	case FalseKeyword:
		return "false", nil

	case NullKeyword, "UndefinedKeyword", "undefined":
		return "nil", nil

	case TemplateExpression:
		return g.generateTemplateExpression(node)

	case ObjectLiteralExpression:
		return g.generateObjectLiteral(node)

	case ArrayLiteralExpression:
		return g.generateArrayLiteral(node)

	case PropertyAccessExpression:
		return g.generatePropertyAccess(node)

	case ElementAccessExpression:
		return g.generateElementAccess(node)

	case CallExpression:
		return g.generateCallExpression(node)

	case NewExpression:
		return g.generateNewExpression(node)

	case BinaryExpression:
		return g.generateBinaryExpression(node)

	case PostfixUnaryExpression:
		return g.generatePostfixUnaryExpression(node)

	case PrefixUnaryExpression:
		return g.generatePrefixUnaryExpression(node)

	case ConditionalExpression:
		return g.generateConditionalExpression(node)

	case ArrowFunction:
		return g.generateArrowFunction(node)

	case "ParenthesizedExpression":
		return g.generateParenthesizedExpression(node)

	case "NonNullExpression":
		return g.generateNonNullExpression(node)

	case "AsExpression":
		return g.generateAsExpression(node)

	case SpreadElement:
		return g.generateSpreadElement(node)

	case TypeOfExpression:
		return g.generateTypeOfExpression(node)

	case DeleteExpression:
		return g.generateDeleteExpression(node)

	case AwaitExpression:
		return g.generateAwaitExpression(node)

	default:
		return "", UnsupportedFeatureError("", 0, 0, fmt.Sprintf("unsupported expression: %s", node.Kind))
	}
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

// generateParenthesizedExpression generates code for parenthesized expressions
func (g *CodeGenerator) generateParenthesizedExpression(node *ASTNode) (string, error) {
	if node.Expression != nil {
		expr, err := g.generateExpression(node.Expression)
		if err != nil {
			return "", fmt.Errorf("generating parenthesized expression: %w", err)
		}
		return fmt.Sprintf("(%s)", expr), nil
	}
	if len(node.Children) > 0 {
		expr, err := g.generateExpression(&node.Children[0])
		if err != nil {
			return "", fmt.Errorf("generating parenthesized expression: %w", err)
		}
		return fmt.Sprintf("(%s)", expr), nil
	}
	return "()", nil
}

// generateNonNullExpression generates code for non-null assertion expressions (expr!)
func (g *CodeGenerator) generateNonNullExpression(node *ASTNode) (string, error) {
	if node.Expression != nil {
		return g.generateExpression(node.Expression)
	}
	if len(node.Children) > 0 {
		return g.generateExpression(&node.Children[0])
	}
	return "", fmt.Errorf("non-null expression missing inner expression")
}

// generateAsExpression generates code for type assertions (expr as Type)
func (g *CodeGenerator) generateAsExpression(node *ASTNode) (string, error) {
	var innerNode *ASTNode
	if node.Expression != nil {
		innerNode = node.Expression
	} else if len(node.Children) > 0 {
		innerNode = &node.Children[0]
	}

	if innerNode == nil {
		return "", fmt.Errorf("as expression missing inner expression")
	}

	expr, err := g.generateExpression(innerNode)
	if err != nil {
		return "", fmt.Errorf("generating as expression operand: %w", err)
	}

	if node.Type != nil {
		targetType, err := g.generateType(node.Type)
		if err == nil && targetType != "" && targetType != "interface{}" {
			if targetType == "string" || targetType == "int" || targetType == "float64" || targetType == "bool" {
				return fmt.Sprintf("%s(%s)", targetType, expr), nil
			}
			return fmt.Sprintf("%s.(%s)", expr, targetType), nil
		}
	}

	return expr, nil
}

// generateSpreadElement generates code for spread elements (...expr)
func (g *CodeGenerator) generateSpreadElement(node *ASTNode) (string, error) {
	var innerNode *ASTNode
	if node.Expression != nil {
		innerNode = node.Expression
	} else if len(node.Children) > 0 {
		innerNode = &node.Children[0]
	}

	if innerNode == nil {
		return "", fmt.Errorf("spread element missing expression")
	}

	expr, err := g.generateExpression(innerNode)
	if err != nil {
		return "", fmt.Errorf("generating spread expression: %w", err)
	}

	return fmt.Sprintf("%s...", expr), nil
}

// generateAwaitExpression generates code for await expression
// TypeScript: await promise
// Go: <-channel (receiving from a channel)
func (g *CodeGenerator) generateAwaitExpression(node *ASTNode) (string, error) {
	if node.Expression == nil {
		return "", fmt.Errorf("await expression missing operand")
	}

	expr, err := g.generateExpression(node.Expression)
	if err != nil {
		return "", fmt.Errorf("generating await operand: %w", err)
	}

	return fmt.Sprintf("(<-%s)", expr), nil
}

// isStructTypeName returns true if the type name represents a user-defined struct
func isStructTypeName(t string) bool {
	clean := strings.TrimPrefix(t, "*")
	if clean == "" {
		return false
	}
	switch clean {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "string", "bool", "byte", "rune",
		"error", "interface{}", "any":
		return false
	}
	if strings.HasPrefix(clean, "[]") || strings.HasPrefix(clean, "map[") {
		return false
	}
	first := clean[0]
	return first >= 'A' && first <= 'Z'
}
