package transpiler

import (
	"fmt"
	"strings"
	"unicode"
)

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
	if objectNode.Kind == Identifier {
		objectName := objectNode.Text
		if len(objectName) > 0 && unicode.IsUpper(rune(objectName[0])) {
			return fmt.Sprintf("%s%s", toPascalCase(objectName), property), nil
		}
	}

	return fmt.Sprintf("%s.%s", object, property), nil
}

// generateCallExpression generates code for function calls
func (g *CodeGenerator) generateCallExpression(node *ASTNode) (string, error) {
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
				g.trackImport("github.com/el-j/ts2go/runtime/array")

				objExpr, err := g.generateExpression(objNode)
				if err != nil {
					return "", err
				}

				args := []string{}
				for _, arg := range node.Children {
					argStr, err := g.generateExpression(&arg)
					if err != nil {
						return "", err
					}
					args = append(args, argStr)
				}

				methodName := toPascalCase(prop)

				if prop == "push" || prop == "pop" || prop == "shift" || prop == "unshift" {
					return fmt.Sprintf("array.%s(&%s, %s)", methodName, objExpr, strings.Join(args, ", ")), nil
				}
				if len(args) > 0 {
					return fmt.Sprintf("array.%s(%s, %s)", methodName, objExpr, strings.Join(args, ", ")), nil
				}
				return fmt.Sprintf("array.%s(%s)", methodName, objExpr), nil
			}
		}
	}

	// Generic function call
	funcExpr, err := g.generateExpression(funcNode)
	if err != nil {
		return "", err
	}

	// Map declared function names to their PascalCase equivalents
	if mappedName, ok := g.declaredFunctions[funcExpr]; ok {
		funcExpr = mappedName
	} else if funcNode.Kind == Identifier && len(funcExpr) > 0 {
		firstChar := funcExpr[0]
		if firstChar >= 'A' && firstChar <= 'Z' {
			funcExpr = toPascalCase(funcExpr)
		}
	}

	args := []string{}
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

	args := []string{}
	for _, arg := range node.Children {
		argStr, err := g.generateExpression(&arg)
		if err != nil {
			return "", err
		}
		args = append(args, argStr)
	}

	return fmt.Sprintf("New%s(%s)", className, strings.Join(args, ", ")), nil
}

// generateElementAccess generates code for element access expressions
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

	var indexNode *ASTNode
	if node.ArgumentExpression != nil {
		indexNode = node.ArgumentExpression
	} else if len(node.Children) > 0 && node.Children[0].Kind != "undefined" {
		indexNode = &node.Children[0]
	}

	if indexNode != nil {
		index, err := g.generateExpression(indexNode)
		if err != nil {
			return "", fmt.Errorf("generating element access index: %w", err)
		}
		return fmt.Sprintf("%s[%s]", expr, index), nil
	}

	return "", fmt.Errorf("element access missing index")
}
