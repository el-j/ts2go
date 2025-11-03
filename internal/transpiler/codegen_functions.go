package transpiler

import (
	"fmt"
	"strings"
)

func (g *CodeGenerator) generateFunction(node *ASTNode) error {
	funcName := toPascalCase(node.Name)

	// Build parameter list and track default parameters
	params := []string{}
	defaultParams := []struct {
		name         string
		defaultValue string
	}{}
	
	if node.Parameters != nil {
		for _, param := range node.Parameters {
			paramName := param.Name
			paramType := "interface{}"
			
			// Check for rest parameter (...args)
			isRestParam := false
			if param.Children != nil {
				for _, child := range param.Children {
					if child.Kind == "DotDotDotToken" {
						isRestParam = true
						break
					}
				}
			}
			
			if param.Type != nil {
				var err error
				paramType, err = g.generateType(param.Type)
				if err != nil {
					return err
				}
			}
			
			// Handle default parameters
			if param.Initializer != nil {
				defaultValue, err := g.generateExpression(param.Initializer)
				if err != nil {
					return err
				}
				defaultParams = append(defaultParams, struct {
					name         string
					defaultValue string
				}{paramName, defaultValue})
			}
			
			// For rest parameters, ensure the type is a slice
			if isRestParam && !strings.HasPrefix(paramType, "[]") {
				// Already a slice type, keep as is
			}
			
			params = append(params, fmt.Sprintf("%s %s", paramName, paramType))
		}
	}

	// Determine return type
	returnType := ""
	if node.Type != nil {
		var err error
		returnType, err = g.generateType(node.Type)
		if err != nil {
			return err
		}
	}

	// Write function signature
	signature := fmt.Sprintf("func %s(%s)", funcName, strings.Join(params, ", "))
	if returnType != "" {
		signature += " " + returnType
	}
	g.writeLine(signature + " {")
	g.indent++

	// Set the current function return type for context-aware generation
	oldReturnType := g.currentFunctionReturnType
	g.currentFunctionReturnType = returnType
	defer func() { g.currentFunctionReturnType = oldReturnType }()

	// Generate default parameter initialization at the start of the function
	// Note: Go doesn't support default parameters natively, so we'd need to use overloading
	// or optional parameters pattern. For now, we'll add a comment.
	if len(defaultParams) > 0 {
		g.writeLine("// Note: Default parameters not fully supported in Go")
		for _, dp := range defaultParams {
			g.writeLine(fmt.Sprintf("// Default: %s = %s", dp.name, dp.defaultValue))
		}
	}

	// Generate function body
	if node.Body != nil {
		if err := g.generateBlock(node.Body); err != nil {
			return err
		}
	}

	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// generateBlock generates code for a block statement
func (g *CodeGenerator) generateBlock(node *ASTNode) error {
	if node.Statements != nil {
		for _, stmt := range node.Statements {
			if err := g.generateStatement(&stmt); err != nil {
				return err
			}
		}
	}
	return nil
}

// generateVariableStatement generates code for variable declarations
