package transpiler

import (
	"fmt"
	"strings"
)

func (g *CodeGenerator) generateFunction(node *ASTNode) error {
	funcName := toPascalCase(node.Name)

	// Build parameter list
	params := []string{}
	if node.Parameters != nil {
		for _, param := range node.Parameters {
			paramName := param.Name
			paramType := "interface{}"
			if param.Type != nil {
				var err error
				paramType, err = g.generateType(param.Type)
				if err != nil {
					return err
				}
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
