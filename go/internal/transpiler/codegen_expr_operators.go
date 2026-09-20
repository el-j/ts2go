package transpiler

import (
	"fmt"
)

// generateConditionalExpression generates code for ternary operator (condition ? whenTrue : whenFalse)
func (g *CodeGenerator) generateConditionalExpression(node *ASTNode) (string, error) {
	var exprChildren []ASTNode
	for _, ch := range node.Children {
		if ch.Kind != "QuestionToken" && ch.Kind != "ColonToken" {
			exprChildren = append(exprChildren, ch)
		}
	}

	if len(exprChildren) < 3 {
		return "", fmt.Errorf("conditional expression has insufficient children: %d", len(exprChildren))
	}

	condition, err := g.generateExpression(&exprChildren[0])
	if err != nil {
		return "", fmt.Errorf("generating ternary condition: %w", err)
	}

	whenTrue, err := g.generateExpression(&exprChildren[1])
	if err != nil {
		return "", fmt.Errorf("generating ternary whenTrue: %w", err)
	}

	whenFalse, err := g.generateExpression(&exprChildren[2])
	if err != nil {
		return "", fmt.Errorf("generating ternary whenFalse: %w", err)
	}

	retType := "interface{}"
	if g.expectedType != "" && g.expectedType != "interface{}" {
		retType = g.expectedType
	}

	result := fmt.Sprintf("func() %s { if %s { return %s } else { return %s } }()",
		retType, condition, whenTrue, whenFalse)

	return result, nil
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

	if node.Operator == QuestionQuestionToken {
		expr := fmt.Sprintf("nullishCoalesce(%s, %s)", left, right)
		if g.expectedType != "" && g.expectedType != "interface{}" {
			expr = fmt.Sprintf("%s.(%s)", expr, g.expectedType)
		}
		return expr, nil
	}

	if node.Operator == "InstanceOfKeyword" || node.Operator == "instanceof" {
		return fmt.Sprintf("/* instanceof check: %s is %s */ reflect.TypeOf(%s).String() == reflect.TypeOf((*%s)(nil)).Elem().String()", left, right, left, right), nil
	}

	if node.Operator == "InKeyword" || node.Operator == "in" {
		return fmt.Sprintf("func() bool { _, exists := %s[%s]; return exists }()", right, left), nil
	}

	operator := g.mapOperator(node.Operator)

	if (operator == "" || node.Operator == "FirstBinaryOperator") && len(node.Children) >= 3 {
		middleChild := &node.Children[1]
		if middleChild.Text != "" {
			operator = middleChild.Text
		}
		switch middleChild.Kind {
		case "LessThanToken", "FirstBinaryOperator":
			operator = "<"
		case "GreaterThanToken":
			operator = ">"
		case "LessThanEqualsToken":
			operator = "<="
		case "GreaterThanEqualsToken":
			operator = ">="
		}
	}

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

	return operand + "++", nil
}

// generatePrefixUnaryExpression generates code for prefix unary expressions (++i, --i, !x, -x)
func (g *CodeGenerator) generatePrefixUnaryExpression(node *ASTNode) (string, error) {
	if len(node.Children) < 1 {
		return "", fmt.Errorf("invalid prefix unary expression")
	}

	var opNode *ASTNode
	var operandNode *ASTNode

	if len(node.Children) >= 2 && isUnaryOperatorToken(node.Children[0].Kind) {
		opNode = &node.Children[0]
		operandNode = &node.Children[1]
	} else {
		operandNode = &node.Children[0]
	}

	operand, err := g.generateExpression(operandNode)
	if err != nil {
		return "", err
	}

	op := ""
	if opNode != nil {
		switch opNode.Kind {
		case "ExclamationToken":
			op = "!"
		case "PlusPlusToken":
			op = "++"
		case "MinusMinusToken":
			op = "--"
		case "MinusToken":
			op = "-"
		case "PlusToken":
			op = "+"
		case "TildeToken":
			op = "^"
		}
	}
	if op == "" && node.Operator != "" {
		switch node.Operator {
		case "ExclamationToken", "!":
			op = "!"
		case "PlusPlusToken", "++":
			op = "++"
		case "MinusMinusToken", "--":
			op = "--"
		case "MinusToken", "-":
			op = "-"
		case "PlusToken", "+":
			op = "+"
		}
	}

	if op != "" {
		return op + operand, nil
	}

	operatorLen := operandNode.Pos - node.Pos
	switch operatorLen {
	case 1:
		if operandNode.Kind == "TrueKeyword" || operandNode.Kind == "FalseKeyword" {
			return "!" + operand, nil
		}
		return "-" + operand, nil
	case 2:
		return "++" + operand, nil
	default:
		return "-" + operand, nil
	}
}

func isUnaryOperatorToken(kind string) bool {
	switch kind {
	case "ExclamationToken", "PlusPlusToken", "MinusMinusToken", "MinusToken", "PlusToken", "TildeToken":
		return true
	default:
		return false
	}
}

// mapOperator maps AST operator constants to Go operators
func (g *CodeGenerator) mapOperator(op string) string {
	switch op {
	case "EqualsToken", "=":
		return "="
	case "PlusEqualsToken", "FirstCompoundAssignment", "+=":
		return "+="
	case "MinusEqualsToken", "-=":
		return "-="
	case "AsteriskEqualsToken", "*=":
		return "*="
	case "SlashEqualsToken", "/=":
		return "/="
	case "PercentEqualsToken", "%=":
		return "%="
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
	case "BarToken":
		return "|"
	case "AmpersandToken":
		return "&"
	case "CaretToken":
		return "^"
	case "LessThanLessThanToken":
		return "<<"
	case "GreaterThanGreaterThanToken":
		return ">>"
	default:
		return ""
	}
}
