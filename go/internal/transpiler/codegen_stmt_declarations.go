package transpiler

import (
	"fmt"
	"strings"
)

// generateVariableStatement generates code for variable statements (let, const, var)
func (g *CodeGenerator) generateVariableStatement(node *ASTNode) error {
	if node.Declarations != nil {
		for _, decl := range node.Declarations {
			// Check if this is a destructuring declaration
			if decl.NameNode != nil {
				if err := g.generateDestructuring(decl.NameNode, decl.Initializer); err != nil {
					return err
				}
				continue
			}

			// Simple variable declaration
			varName := decl.Name
			if decl.Initializer != nil {
				oldExpected := g.expectedType
				if decl.Type != nil {
					varType, err := g.generateType(decl.Type)
					if err == nil && varType != "" {
						g.expectedType = varType
					}
				}
				init, err := g.generateExpression(decl.Initializer)
				g.expectedType = oldExpected
				if err != nil {
					return err
				}
				if node.IsConst && (decl.Initializer == nil || isConstantLiteral(decl.Initializer)) {
					g.writeLine(fmt.Sprintf("const %s = %s", varName, init))
				} else {
					g.writeLine(fmt.Sprintf("%s := %s", varName, init))
				}
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

// generateDestructuring generates code for destructuring assignments
// TypeScript: const {x, y} = obj  or  const [a, b] = arr
// Go: temp := obj; x := temp["x"]; y := temp["y"]  or  temp := arr; a := temp[0]; b := temp[1]
func (g *CodeGenerator) generateDestructuring(pattern *ASTNode, initializer *ASTNode) error {
	if initializer == nil {
		return fmt.Errorf("destructuring requires an initializer")
	}

	init, err := g.generateExpression(initializer)
	if err != nil {
		return fmt.Errorf("generating destructuring initializer: %w", err)
	}

	// Create a temporary variable to hold the initializer
	// This avoids duplicating complex expressions
	tempVar := fmt.Sprintf("_destructure_%d", g.tempVarCounter)
	g.tempVarCounter++
	g.writeLine(fmt.Sprintf("%s := %s", tempVar, init))

	switch pattern.Kind {
	case ObjectBindingPattern:
		if pattern.Elements != nil {
			for _, elem := range pattern.Elements {
				if elem.Name != "" {
					hasRest := false
					if elem.Children != nil {
						for _, child := range elem.Children {
							if child.Kind == "DotDotDotToken" {
								hasRest = true
								break
							}
						}
					}

					if hasRest {
						g.writeLine(fmt.Sprintf("%s := make(map[string]interface{})", elem.Name))
						g.writeLine(fmt.Sprintf("for k, v := range %s.(map[string]interface{}) {", tempVar))
						g.indent++
						g.writeLine(fmt.Sprintf("%s[k] = v", elem.Name))
						g.indent--
						g.writeLine("}")

						for _, otherElem := range pattern.Elements {
							if otherElem.Name != "" && otherElem.Name != elem.Name {
								g.writeLine(fmt.Sprintf("delete(%s, \"%s\")", elem.Name, otherElem.Name))
							}
						}
					} else {
						g.writeLine(fmt.Sprintf(`%s := %s.(map[string]interface{})["%s"]`, elem.Name, tempVar, elem.Name))
					}
				}
			}
		}

	case ArrayBindingPattern:
		if pattern.Elements != nil {
			for i, elem := range pattern.Elements {
				if elem.Name != "" {
					hasRest := false
					if elem.Children != nil {
						for _, child := range elem.Children {
							if child.Kind == "DotDotDotToken" {
								hasRest = true
								break
							}
						}
					}

					if hasRest {
						g.writeLine(fmt.Sprintf("%s := %s[%d:]", elem.Name, tempVar, i))
					} else {
						g.writeLine(fmt.Sprintf("%s := %s[%d]", elem.Name, tempVar, i))
					}
				}
			}
		}

	default:
		return fmt.Errorf("unsupported binding pattern: %s", pattern.Kind)
	}

	return nil
}

// generateReturnStatement generates code for return statements
func (g *CodeGenerator) generateReturnStatement(node *ASTNode) error {
	var expressionNode *ASTNode
	if node.Expression != nil {
		expressionNode = node.Expression
	} else if len(node.Children) > 0 {
		expressionNode = &node.Children[0]
	}

	if expressionNode == nil {
		g.writeLine("return")
		return nil
	}

	oldExpected := g.expectedType
	if g.currentFunctionReturnType != "" {
		g.expectedType = g.currentFunctionReturnType
	}
	expr, err := g.generateExpression(expressionNode)
	g.expectedType = oldExpected
	if err != nil {
		return err
	}

	if strings.Contains(expr, "nullishCoalesce(") && g.currentFunctionReturnType != "" && g.currentFunctionReturnType != "interface{}" {
		expr = fmt.Sprintf("%s.(%s)", expr, g.currentFunctionReturnType)
	} else if g.currentFunctionReturnType == "float64" && expressionNode.Kind == "Identifier" {
		expr = fmt.Sprintf("float64(%s)", expr)
	}

	g.writeLine(fmt.Sprintf("return %s", expr))
	return nil
}

// generateBreakStatement generates code for break statements
func (g *CodeGenerator) generateBreakStatement(node *ASTNode) error {
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
	if len(node.Children) > 0 {
		label := node.Children[0].Text
		g.writeLine(fmt.Sprintf("continue %s", label))
	} else {
		g.writeLine("continue")
	}
	return nil
}

// generateThrowStatement generates code for throw statements
func (g *CodeGenerator) generateThrowStatement(node *ASTNode) error {
	if node.Expression == nil {
		g.writeLine("panic(nil)")
		return nil
	}

	expr, err := g.generateExpression(node.Expression)
	if err != nil {
		return fmt.Errorf("generating throw expression: %w", err)
	}

	g.writeLine(fmt.Sprintf("panic(%s)", expr))
	return nil
}
