package transpiler

import (
	"fmt"
)

// generateStatementBlock generates code for a block of statements
func (g *CodeGenerator) generateStatementBlock(node *ASTNode) error {
	if node.Kind == Block {
		if node.Statements != nil {
			for _, stmt := range node.Statements {
				if err := g.generateStatement(&stmt); err != nil {
					return err
				}
			}
		}
	} else {
		if err := g.generateStatement(node); err != nil {
			return err
		}
	}
	return nil
}

// generateForStatement generates code for traditional for loops
func (g *CodeGenerator) generateForStatement(node *ASTNode) error {
	if len(node.Children) < 3 {
		return fmt.Errorf("for statement requires 3 children (condition, increment, body)")
	}

	var initPart string
	if node.Initializer != nil {
		if len(node.Initializer.Children) > 0 {
			decl := &node.Initializer.Children[0]
			if decl.Name != "" && decl.Initializer != nil {
				init, err := g.generateExpression(decl.Initializer)
				if err != nil {
					return fmt.Errorf("generating for init: %w", err)
				}
				initPart = fmt.Sprintf("%s := %s", decl.Name, init)
			}
		}
	}

	condPart, err := g.generateExpression(&node.Children[0])
	if err != nil {
		return fmt.Errorf("generating for condition: %w", err)
	}

	incPart, err := g.generateExpression(&node.Children[1])
	if err != nil {
		return fmt.Errorf("generating for increment: %w", err)
	}

	g.writeLine(fmt.Sprintf("for %s; %s; %s {", initPart, condPart, incPart))
	g.indent++

	bodyNode := &node.Children[2]
	if err := g.generateStatementBlock(bodyNode); err != nil {
		return fmt.Errorf("generating for body: %w", err)
	}

	g.indent--
	g.writeLine("}")

	return nil
}

// generateForOfStatement generates code for for...of loops
func (g *CodeGenerator) generateForOfStatement(node *ASTNode) error {
	varName := ""
	if node.Initializer != nil && len(node.Initializer.Children) > 0 {
		varDecl := &node.Initializer.Children[0]
		varName = varDecl.Name
	}

	if varName == "" {
		return fmt.Errorf("for-of: could not determine variable name")
	}

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

	g.writeLine(fmt.Sprintf("for _, %s := range %s {", varName, iterable))
	g.indent++

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
func (g *CodeGenerator) generateForInStatement(node *ASTNode) error {
	varName := ""
	if node.Initializer != nil && len(node.Initializer.Children) > 0 {
		varDecl := &node.Initializer.Children[0]
		varName = varDecl.Name
	}

	if varName == "" {
		return fmt.Errorf("for-in: could not determine variable name")
	}

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

	g.writeLine(fmt.Sprintf("for %s := range %s {", varName, object))
	g.indent++

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
func (g *CodeGenerator) generateWhileStatement(node *ASTNode) error {
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

	g.writeLine(fmt.Sprintf("for %s {", condition))
	g.indent++

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
