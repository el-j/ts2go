package transpiler

import (
	"fmt"
	"strings"
)

// generateTemplateExpression generates code for template literal expressions
// TypeScript: `Hello ${name}!`
// Go: "Hello " + name + "!"
func (g *CodeGenerator) generateTemplateExpression(node *ASTNode) (string, error) {
	parts := []string{}

	// Add the head if it exists
	if node.Head != nil && node.Head.Text != "" {
		text := strings.ReplaceAll(node.Head.Text, "\n", "\\n")
		text = strings.ReplaceAll(text, "\t", "\\t")
		text = strings.ReplaceAll(text, "\"", "\\\"")
		parts = append(parts, fmt.Sprintf(`"%s"`, text))
	}

	// Process template spans (expression + text pairs)
	if node.Children != nil {
		for _, child := range node.Children {
			if child.Kind == "TemplateSpan" {
				if child.Expression != nil {
					expr, err := g.generateExpression(child.Expression)
					if err != nil {
						return "", fmt.Errorf("generating template expression: %w", err)
					}
					g.trackImport("fmt")
					parts = append(parts, fmt.Sprintf("fmt.Sprint(%s)", expr))
				}

				if child.Literal != nil && child.Literal.Text != "" {
					text := strings.ReplaceAll(child.Literal.Text, "\n", "\\n")
					text = strings.ReplaceAll(text, "\t", "\\t")
					text = strings.ReplaceAll(text, "\"", "\\\"")
					parts = append(parts, fmt.Sprintf(`"%s"`, text))
				}
			}
		}
	}

	if len(parts) == 1 {
		return parts[0], nil
	}

	if len(parts) == 0 {
		return `""`, nil
	}

	return strings.Join(parts, " + "), nil
}

// generateObjectLiteral generates code for object literals
func (g *CodeGenerator) generateObjectLiteral(node *ASTNode) (string, error) {
	if len(node.Properties) == 0 {
		return "map[string]interface{}{}", nil
	}

	// Check if any properties are spread elements
	hasSpread := false
	for _, prop := range node.Properties {
		switch prop.Kind {
		case SpreadElement, SpreadAssignment:
			hasSpread = true
		}
		if hasSpread {
			break
		}
	}

	// If there are spread elements, we need to merge maps
	if hasSpread {
		parts := []string{}
		currentProps := []string{}

		for _, prop := range node.Properties {
			switch prop.Kind {
			case SpreadElement, SpreadAssignment:
				if len(currentProps) > 0 {
					parts = append(parts, fmt.Sprintf("map[string]interface{}{%s}", strings.Join(currentProps, ", ")))
					currentProps = []string{}
				}
				if prop.Expression != nil {
					expr, err := g.generateExpression(prop.Expression)
					if err != nil {
						return "", err
					}
					parts = append(parts, expr)
				}
			case "PropertyAssignment":
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

		if len(currentProps) > 0 {
			parts = append(parts, fmt.Sprintf("map[string]interface{}{%s}", strings.Join(currentProps, ", ")))
		}

		if len(parts) == 1 {
			return parts[0], nil
		}

		return g.generateObjectMerge(parts), nil
	}

	targetType := g.expectedType
	if targetType == "" && g.currentFunctionReturnType != "" {
		targetType = g.currentFunctionReturnType
	}

	if isStructTypeName(targetType) {
		isPtr := strings.HasPrefix(targetType, "*")
		structName := strings.TrimPrefix(targetType, "*")

		fields := []string{}
		for _, prop := range node.Properties {
			if prop.Kind == "PropertyAssignment" {
				propName := toPascalCase(prop.Name)
				var value string
				var err error
				if prop.Initializer != nil {
					value, err = g.generateExpression(prop.Initializer)
					if err != nil {
						return "", err
					}
				}
				fields = append(fields, fmt.Sprintf("%s: %s", propName, value))
			}
		}

		prefix := ""
		if isPtr {
			prefix = "&"
		}
		return fmt.Sprintf("%s%s{%s}", prefix, structName, strings.Join(fields, ", ")), nil
	}

	// No spread elements - generate simple object literal
	fields := []string{}
	for _, prop := range node.Properties {
		if prop.Kind == "PropertyAssignment" {
			propName := prop.Name
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

	return fmt.Sprintf("map[string]interface{}{%s}", strings.Join(fields, ", ")), nil
}

// generateObjectMerge generates code to merge multiple objects/maps
func (g *CodeGenerator) generateObjectMerge(parts []string) string {
	tempVar := fmt.Sprintf("_merge_%d", g.tempVarCounter)
	g.tempVarCounter++

	mergeCode := fmt.Sprintf("func() map[string]interface{} { %s := make(map[string]interface{}); ", tempVar)
	for _, part := range parts {
		mergeCode += fmt.Sprintf("for k, v := range %s { %s[k] = v }; ", part, tempVar)
	}
	mergeCode += fmt.Sprintf("return %s }()", tempVar)

	return mergeCode
}

// generateArrayLiteral generates code for array literals
func (g *CodeGenerator) generateArrayLiteral(node *ASTNode) (string, error) {
	if len(node.Elements) == 0 {
		return "[]interface{}{}", nil
	}

	hasSpread := false
	for _, elem := range node.Elements {
		if elem.Kind == SpreadElement {
			hasSpread = true
			break
		}
	}

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

	g.trackImport("github.com/el-j/ts2go/runtime/array")

	slices := []string{}
	currentLiteral := []string{}

	for _, elem := range node.Elements {
		if elem.Kind == SpreadElement {
			if len(currentLiteral) > 0 {
				slices = append(slices, fmt.Sprintf("[]interface{}{%s}", strings.Join(currentLiteral, ", ")))
				currentLiteral = []string{}
			}

			if elem.Expression != nil {
				expr, err := g.generateExpression(elem.Expression)
				if err != nil {
					return "", err
				}
				slices = append(slices, expr)
			}
		} else {
			value, err := g.generateExpression(&elem)
			if err != nil {
				return "", err
			}
			currentLiteral = append(currentLiteral, value)
		}
	}

	if len(currentLiteral) > 0 {
		slices = append(slices, fmt.Sprintf("[]interface{}{%s}", strings.Join(currentLiteral, ", ")))
	}

	if len(slices) == 1 {
		return slices[0], nil
	}

	return fmt.Sprintf("array.Concat(%s)", strings.Join(slices, ", ")), nil
}
