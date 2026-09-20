package transpiler

import (
	"fmt"
)

// generateIfStatement generates code for if/else statements
func (g *CodeGenerator) generateIfStatement(node *ASTNode) error {
	var conditionNode *ASTNode

	if node.Expression != nil {
		conditionNode = node.Expression
	} else if len(node.Children) > 0 {
		conditionNode = &node.Children[0]
	} else {
		return fmt.Errorf("if statement missing condition")
	}

	condition, err := g.generateExpression(conditionNode)
	if err != nil {
		return fmt.Errorf("generating if condition: %w", err)
	}

	g.writeLine(fmt.Sprintf("if %s {", condition))
	g.indent++

	if len(node.Children) > 1 {
		thenBlock := &node.Children[1]
		if err := g.generateStatementBlock(thenBlock); err != nil {
			return fmt.Errorf("generating if then block: %w", err)
		}
	}

	g.indent--

	if len(node.Children) > 2 {
		elseBlock := &node.Children[2]

		if elseBlock.Kind == IfStatement {
			for i := 0; i < g.indent; i++ {
				g.output.WriteString("\t")
			}
			g.output.WriteString("} else ")

			if len(elseBlock.Children) > 0 {
				elseCond, err := g.generateExpression(&elseBlock.Children[0])
				if err != nil {
					return fmt.Errorf("generating else-if condition: %w", err)
				}
				fmt.Fprintf(&g.output, "if %s {\n", elseCond)

				g.indent++
				if len(elseBlock.Children) > 1 {
					if err := g.generateStatementBlock(&elseBlock.Children[1]); err != nil {
						return fmt.Errorf("generating else-if block: %w", err)
					}
				}
				g.indent--

				if len(elseBlock.Children) > 2 {
					furtherElse := &elseBlock.Children[2]
					if furtherElse.Kind == IfStatement {
						for i := 0; i < g.indent; i++ {
							g.output.WriteString("\t")
						}
						g.output.WriteString("} else ")
						g.output.WriteString("{\n")
						g.indent++
						if err := g.generateStatement(furtherElse); err != nil {
							return err
						}
						g.indent--
						g.writeLine("}")
					} else {
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
		}

		g.writeLine("} else {")
		g.indent++

		if err := g.generateStatementBlock(elseBlock); err != nil {
			return fmt.Errorf("generating else block: %w", err)
		}

		g.indent--
		g.writeLine("}")
	} else {
		g.writeLine("}")
	}

	return nil
}

// generateSwitchStatement generates code for switch statements
func (g *CodeGenerator) generateSwitchStatement(node *ASTNode) error {
	if len(node.Children) < 2 {
		return fmt.Errorf("switch statement requires 2 children (expression, case block)")
	}

	switchExpr, err := g.generateExpression(&node.Children[0])
	if err != nil {
		return fmt.Errorf("generating switch expression: %w", err)
	}

	g.writeLine(fmt.Sprintf("switch %s {", switchExpr))

	caseBlock := &node.Children[1]
	if caseBlock.Kind != CaseBlock {
		return fmt.Errorf("expected CaseBlock, got %s", caseBlock.Kind)
	}

	for _, clause := range caseBlock.Children {
		switch clause.Kind {
		case CaseClause:
			if err := g.generateCaseClause(&clause); err != nil {
				return fmt.Errorf("generating case clause: %w", err)
			}
		case DefaultClause:
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
	if len(node.Children) < 1 {
		return fmt.Errorf("case clause requires at least 1 child (case expression)")
	}

	caseValue, err := g.generateExpression(&node.Children[0])
	if err != nil {
		return fmt.Errorf("generating case value: %w", err)
	}

	g.writeLine(fmt.Sprintf("case %s:", caseValue))
	g.indent++

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
	g.writeLine("default:")
	g.indent++

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

// generateTryStatement generates code for try/catch/finally statements using defer/recover
func (g *CodeGenerator) generateTryStatement(node *ASTNode) error {
	if len(node.Children) == 0 {
		return fmt.Errorf("try statement has no children")
	}

	var tryBlock *ASTNode
	var catchClause *ASTNode
	var finallyBlock *ASTNode

	for i := range node.Children {
		child := &node.Children[i]
		switch child.Kind {
		case Block:
			if tryBlock == nil {
				tryBlock = child
			} else {
				finallyBlock = child
			}
		case CatchClause:
			catchClause = child
		}
	}

	if tryBlock == nil {
		return fmt.Errorf("try statement missing try block")
	}

	g.writeLine("func() {")
	g.indent++

	if catchClause != nil {
		errorVarName := "err"
		if catchClause.Children != nil {
			for _, child := range catchClause.Children {
				if child.Kind == VariableDeclaration && child.Name != "" {
					errorVarName = child.Name
					break
				}
			}
		}

		g.writeLine("defer func() {")
		g.indent++
		g.writeLine(fmt.Sprintf("if %s := recover(); %s != nil {", errorVarName, errorVarName))
		g.indent++

		var catchBlock *ASTNode
		if catchClause.Children != nil {
			for i := range catchClause.Children {
				if catchClause.Children[i].Kind == Block {
					catchBlock = &catchClause.Children[i]
					break
				}
			}
		}

		if catchBlock != nil {
			if err := g.generateBlock(catchBlock); err != nil {
				return fmt.Errorf("generating catch block: %w", err)
			}
		}

		g.indent--
		g.writeLine("}")
		g.indent--
		g.writeLine("}()")
	}

	if finallyBlock != nil {
		g.writeLine("defer func() {")
		g.indent++
		if err := g.generateBlock(finallyBlock); err != nil {
			return fmt.Errorf("generating finally block: %w", err)
		}
		g.indent--
		g.writeLine("}()")
	}

	if err := g.generateBlock(tryBlock); err != nil {
		return fmt.Errorf("generating try block: %w", err)
	}

	g.indent--
	g.writeLine("}()")

	return nil
}
