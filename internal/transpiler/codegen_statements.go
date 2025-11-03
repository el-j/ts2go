package transpiler

import (
	"fmt"
	"strings"
)

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
	case TryStatement:
		return g.generateTryStatement(node)
	case ThrowStatement:
		return g.generateThrowStatement(node)
	default:
		// Skip unknown statements for now
		return nil
	}
}

// generateInterface converts a TS interface to a Go struct
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

// generateDestructuring generates code for destructuring assignments
// TypeScript: const {x, y} = obj  or  const [a, b] = arr
// Go: x := obj["x"].(type); y := obj["y"].(type)  or  a := arr[0]; b := arr[1]
func (g *CodeGenerator) generateDestructuring(pattern *ASTNode, initializer *ASTNode) error {
	if initializer == nil {
		return fmt.Errorf("destructuring requires an initializer")
	}
	
	init, err := g.generateExpression(initializer)
	if err != nil {
		return fmt.Errorf("generating destructuring initializer: %w", err)
	}
	
	switch pattern.Kind {
	case ObjectBindingPattern:
		// Object destructuring: const {x, y, z} = obj
		// Generate: x := obj.(map[string]interface{})["x"]; y := obj.(map[string]interface{})["y"]
		if pattern.Elements != nil {
			for _, elem := range pattern.Elements {
				if elem.Name != "" {
					// Simple property: {x} means extract property "x" from object
					g.writeLine(fmt.Sprintf(`%s := %s.(map[string]interface{})["%s"]`, elem.Name, init, elem.Name))
				}
			}
		}
		
	case ArrayBindingPattern:
		// Array destructuring: const [a, b, c] = arr
		// Generate: a := arr.([]interface{})[0]; b := arr.([]interface{})[1]
		if pattern.Elements != nil {
			for i, elem := range pattern.Elements {
				if elem.Name != "" {
					g.writeLine(fmt.Sprintf("%s := %s.([]interface{})[%d]", elem.Name, init, i))
				}
			}
		}
		
	default:
		return fmt.Errorf("unsupported binding pattern: %s", pattern.Kind)
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
	if len(node.Children) == 0 {
		g.writeLine("return")
		return nil
	}

	// Generate the expression to return
	expr, err := g.generateExpression(&node.Children[0])
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
	// Get the condition expression
	var conditionNode *ASTNode
	
	// Check for expression field first (newer parser format)
	if node.Expression != nil {
		conditionNode = node.Expression
	} else if len(node.Children) > 0 {
		// Fall back to first child
		conditionNode = &node.Children[0]
	} else {
		return fmt.Errorf("if statement missing condition")
	}

	condition, err := g.generateExpression(conditionNode)
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

// generateTryStatement generates code for try/catch/finally statements
// TypeScript: try { } catch (e) { } finally { }
// Go: Uses defer and recover pattern
func (g *CodeGenerator) generateTryStatement(node *ASTNode) error {
if node.Children == nil || len(node.Children) == 0 {
return fmt.Errorf("try statement has no children")
}

// Find the try block, catch clause, and finally block
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
// This is the finally block
finallyBlock = child
}
case CatchClause:
catchClause = child
}
}

if tryBlock == nil {
return fmt.Errorf("try statement missing try block")
}

// Generate the try block with defer/recover pattern
g.writeLine("func() {")
g.indent++

// If there's a catch clause, add defer with recover
if catchClause != nil {
// Extract error variable name from catch clause
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

// Find the catch block
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

// If there's a finally block, add another defer
if finallyBlock != nil {
g.writeLine("defer func() {")
g.indent++
if err := g.generateBlock(finallyBlock); err != nil {
return fmt.Errorf("generating finally block: %w", err)
}
g.indent--
g.writeLine("}()")
}

// Generate the try block content
if err := g.generateBlock(tryBlock); err != nil {
return fmt.Errorf("generating try block: %w", err)
}

g.indent--
g.writeLine("}()")

return nil
}

// generateThrowStatement generates code for throw statements
// TypeScript: throw new Error("message")
// Go: panic(errors.New("message")) or panic(err)
func (g *CodeGenerator) generateThrowStatement(node *ASTNode) error {
if node.Expression == nil {
g.writeLine("panic(nil)")
return nil
}

// Generate the expression to throw
expr, err := g.generateExpression(node.Expression)
if err != nil {
return fmt.Errorf("generating throw expression: %w", err)
}

// If it's a new Error(...), convert to errors.New(...)
// For simplicity, we'll just panic with the expression
g.writeLine(fmt.Sprintf("panic(%s)", expr))

return nil
}
