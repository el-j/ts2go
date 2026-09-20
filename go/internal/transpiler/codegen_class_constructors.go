package transpiler

import (
	"fmt"
	"strings"
)

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

	// Track super() call parameters for base class initialization
	var superArgs []string
	hasSuperCall := false

	// Check for super() call in constructor body
	if node.Body != nil && node.Body.Statements != nil {
		for _, stmt := range node.Body.Statements {
			if g.isSuperCall(&stmt) {
				hasSuperCall = true
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
	oldReceiverVar := g.currentReceiverVar
	g.currentReceiverVar = instanceVar
	defer func() { g.currentReceiverVar = oldReceiverVar }()

	if node.Kind == ExpressionStatement {
		var binary *ASTNode
		if node.Expression != nil && node.Expression.Kind == BinaryExpression {
			binary = node.Expression
		} else if len(node.Children) > 0 && node.Children[0].Kind == BinaryExpression {
			binary = &node.Children[0]
		}

		if binary != nil && len(binary.Children) >= 2 {
			left := &binary.Children[0]
			right := &binary.Children[len(binary.Children)-1]

			var isThisAccess bool
			if left.Kind == PropertyAccessExpression {
				if left.Expression != nil && left.Expression.Kind == ThisKeyword {
					isThisAccess = true
				} else if len(left.Children) > 0 && left.Children[0].Kind == ThisKeyword {
					isThisAccess = true
				}
			}

			if isThisAccess {
				fieldName := left.Name
				if g.currentClassMembers != nil && g.currentClassMembers[fieldName] {
					fieldName = toCamelCase(fieldName)
				} else {
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
	return g.generateStatement(node)
}
