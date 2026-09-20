package transpiler

import (
	"fmt"
	"strings"
)

// generateMethod generates a method with receiver
func (g *CodeGenerator) generateMethod(className string, node *ASTNode) error {
	switch node.Kind {
	case GetAccessor:
		return g.generateGetter(className, node)
	case SetAccessor:
		return g.generateSetter(className, node)
	}

	methodName := toPascalCase(node.Name)
	if isPrivate(node) {
		methodName = toCamelCase(node.Name)
	}

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

	returnType := ""
	if node.Type != nil {
		var err error
		returnType, err = g.generateType(node.Type)
		if err != nil {
			return err
		}
	}

	receiverVar := strings.ToLower(string(className[0]))
	signature := fmt.Sprintf("func (%s *%s) %s(%s)", receiverVar, className, methodName, strings.Join(params, ", "))
	if returnType != "" {
		signature += " " + returnType
	}
	g.writeLine(signature + " {")
	g.indent++

	oldReturnType := g.currentFunctionReturnType
	oldReceiverVar := g.currentReceiverVar
	g.currentFunctionReturnType = returnType
	g.currentReceiverVar = receiverVar
	defer func() {
		g.currentFunctionReturnType = oldReturnType
		g.currentReceiverVar = oldReceiverVar
	}()

	if node.Body != nil {
		if err := g.generateMethodBody(node.Body, receiverVar); err != nil {
			return err
		}
	}

	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// generateGetter generates a getter method (Get prefix)
func (g *CodeGenerator) generateGetter(className string, node *ASTNode) error {
	methodName := "Get" + toPascalCase(node.Name)

	returnType := "interface{}"
	if node.Type != nil {
		var err error
		returnType, err = g.generateType(node.Type)
		if err != nil {
			return err
		}
	}

	receiverVar := strings.ToLower(string(className[0]))
	g.writeLine(fmt.Sprintf("func (%s *%s) %s() %s {", receiverVar, className, methodName, returnType))
	g.indent++

	oldReturnType := g.currentFunctionReturnType
	oldReceiverVar := g.currentReceiverVar
	g.currentFunctionReturnType = returnType
	g.currentReceiverVar = receiverVar
	defer func() {
		g.currentFunctionReturnType = oldReturnType
		g.currentReceiverVar = oldReceiverVar
	}()

	if node.Body != nil {
		if err := g.generateMethodBody(node.Body, receiverVar); err != nil {
			return err
		}
	}

	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// generateSetter generates a setter method (Set prefix)
func (g *CodeGenerator) generateSetter(className string, node *ASTNode) error {
	methodName := "Set" + toPascalCase(node.Name)

	paramName := "value"
	paramType := "interface{}"
	if len(node.Parameters) > 0 {
		paramName = node.Parameters[0].Name
		if node.Parameters[0].Type != nil {
			var err error
			paramType, err = g.generateType(node.Parameters[0].Type)
			if err != nil {
				return err
			}
		}
	}

	receiverVar := strings.ToLower(string(className[0]))
	g.writeLine(fmt.Sprintf("func (%s *%s) %s(%s %s) {", receiverVar, className, methodName, paramName, paramType))
	g.indent++

	oldReceiverVar := g.currentReceiverVar
	g.currentReceiverVar = receiverVar
	defer func() {
		g.currentReceiverVar = oldReceiverVar
	}()

	if node.Body != nil {
		if err := g.generateMethodBody(node.Body, receiverVar); err != nil {
			return err
		}
	}

	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// generateStaticMethod generates a static method as a package-level function
func (g *CodeGenerator) generateStaticMethod(className string, node *ASTNode) error {
	methodName := className + toPascalCase(node.Name)
	if isPrivate(node) {
		methodName = toCamelCase(className + toPascalCase(node.Name))
	}

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

	returnType := ""
	if node.Type != nil {
		var err error
		returnType, err = g.generateType(node.Type)
		if err != nil {
			return err
		}
	}

	signature := fmt.Sprintf("func %s(%s)", methodName, strings.Join(params, ", "))
	if returnType != "" {
		signature += " " + returnType
	}
	g.writeLine(signature + " {")
	g.indent++

	oldReturnType := g.currentFunctionReturnType
	g.currentFunctionReturnType = returnType
	defer func() {
		g.currentFunctionReturnType = oldReturnType
	}()

	if node.Body != nil && node.Body.Statements != nil {
		for _, stmt := range node.Body.Statements {
			if err := g.generateStatement(&stmt); err != nil {
				return err
			}
		}
	}

	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// generateMethodBody generates the body of a method, replacing "this" with receiver variable
func (g *CodeGenerator) generateMethodBody(bodyNode *ASTNode, receiverVar string) error {
	if bodyNode.Statements != nil {
		for _, stmt := range bodyNode.Statements {
			if err := g.generateMethodStatement(&stmt, receiverVar); err != nil {
				return err
			}
		}
	}
	return nil
}

// generateMethodStatement handles statements in methods, replacing "this" references
func (g *CodeGenerator) generateMethodStatement(node *ASTNode, receiverVar string) error {
	return g.generateStatement(node)
}
