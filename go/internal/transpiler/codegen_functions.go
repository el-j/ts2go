package transpiler

import (
	"fmt"
	"strings"
)

func (g *CodeGenerator) generateFunction(node *ASTNode) error {
	funcName := toPascalCase(node.Name)

	// Check if function is async
	isAsync := false
	if node.Modifiers != nil {
		for _, mod := range node.Modifiers {
			if mod.Kind == AsyncKeyword {
				isAsync = true
				break
			}
		}
	}

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

			// For rest parameters, convert to Go variadic parameters
			if isRestParam {
				// If type is array type like []T, extract T for variadic ...T
				if strings.HasPrefix(paramType, "[]") {
					elementType := strings.TrimPrefix(paramType, "[]")
					params = append(params, fmt.Sprintf("%s ...%s", paramName, elementType))
				} else {
					// Default to ...interface{}
					params = append(params, fmt.Sprintf("%s ...interface{}", paramName))
				}
			} else {
				params = append(params, fmt.Sprintf("%s %s", paramName, paramType))
			}
		}
	}

	// Handle default parameters by creating optional variadic wrapper
	// We'll use a pattern where optional parameters become variadic at the end
	finalParams := []string{}
	defaultInits := []struct {
		paramName    string
		defaultValue string
		paramType    string
	}{}

	if len(defaultParams) > 0 && node.Parameters != nil {
		// Find which parameters have defaults
		for i, param := range node.Parameters {
			if i >= len(params) {
				break
			}

			// Check if this parameter has a default
			hasDefault := false
			var defaultValue string
			for _, dp := range defaultParams {
				if dp.name == param.Name {
					hasDefault = true
					defaultValue = dp.defaultValue
					break
				}
			}

			// Extract type from "name type" format
			parts := strings.Split(params[i], " ")
			if len(parts) < 2 {
				finalParams = append(finalParams, params[i])
				continue
			}

			paramName := parts[0]
			paramType := strings.Join(parts[1:], " ")

			if hasDefault {
				// Convert to optional variadic parameter
				finalParams = append(finalParams, fmt.Sprintf("_%s_opt ...%s", paramName, paramType))
				defaultInits = append(defaultInits, struct {
					paramName    string
					defaultValue string
					paramType    string
				}{paramName, defaultValue, paramType})
			} else {
				finalParams = append(finalParams, params[i])
			}
		}
	} else {
		finalParams = params
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

	// For async functions, wrap return type in channel
	if isAsync {
		if returnType == "" {
			returnType = "chan interface{}"
		} else {
			returnType = fmt.Sprintf("chan %s", returnType)
		}
	}

	// Write function signature
	signature := fmt.Sprintf("func %s(%s)", funcName, strings.Join(finalParams, ", "))
	if returnType != "" {
		signature += " " + returnType
	}
	g.writeLine(signature + " {")
	g.indent++

	// Add initializations for optional parameters with defaults
	for _, init := range defaultInits {
		g.writeLine(fmt.Sprintf("%s := %s", init.paramName, init.defaultValue))
		g.writeLine(fmt.Sprintf("if len(_%s_opt) > 0 {", init.paramName))
		g.indent++
		g.writeLine(fmt.Sprintf("%s = _%s_opt[0]", init.paramName, init.paramName))
		g.indent--
		g.writeLine("}")
	}

	// For async functions, create result channel
	if isAsync {
		g.writeLine("resultCh := make(chan interface{}, 1)")
		g.writeLine("go func() {")
		g.indent++
	}

	// Set the current function return type for context-aware generation
	oldReturnType := g.currentFunctionReturnType
	g.currentFunctionReturnType = returnType
	defer func() { g.currentFunctionReturnType = oldReturnType }()

	// Generate function body
	if node.Body != nil {
		if err := g.generateBlock(node.Body); err != nil {
			return err
		}
	}

	// For async functions, close the goroutine and return channel
	if isAsync {
		g.indent--
		g.writeLine("}()")
		g.writeLine("return resultCh")
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
