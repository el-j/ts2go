package transpiler

import (
	"fmt"
)

// generateTypeOfExpression generates code for typeof operator
// TypeScript: typeof x
// Go: reflect.TypeOf(x).String()
func (g *CodeGenerator) generateTypeOfExpression(node *ASTNode) (string, error) {
	if node.Expression == nil {
		return "", fmt.Errorf("typeof expression missing operand")
	}

	expr, err := g.generateExpression(node.Expression)
	if err != nil {
		return "", fmt.Errorf("generating typeof operand: %w", err)
	}

	return fmt.Sprintf("reflect.TypeOf(%s).String()", expr), nil
}

// generateDeleteExpression generates code for delete operator
// TypeScript: delete obj.prop or delete obj["key"]
// Go: Generates delete() for maps, comment for others
func (g *CodeGenerator) generateDeleteExpression(node *ASTNode) (string, error) {
	if node.Expression == nil {
		return "", fmt.Errorf("delete expression missing operand")
	}

	expr := node.Expression
	switch expr.Kind {
	case PropertyAccessExpression:
		obj, err := g.generateExpression(expr.Expression)
		if err != nil {
			return "", err
		}
		propName := expr.Name
		return fmt.Sprintf("delete(%s, \"%s\")", obj, propName), nil
	case "ElementAccessExpression":
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

	operand, err := g.generateExpression(expr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/* delete %s - not supported in Go */", operand), nil
}
