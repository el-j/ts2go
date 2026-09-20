package transpiler

import (
	"fmt"
)

func (g *CodeGenerator) generateClass(node *ASTNode) error {
	className := toPascalCase(node.Name)

	// Initialize member tracking for this class
	g.currentClassMembers = make(map[string]bool)
	defer func() { g.currentClassMembers = nil }()

	// First pass: collect member visibility information
	if node.Members != nil {
		for i := range node.Members {
			member := &node.Members[i]
			if member.Kind == PropertyDeclaration || member.Kind == MethodDeclaration {
				g.currentClassMembers[member.Name] = isPrivate(member)
			}
		}
	}

	// Check for base class (inheritance)
	var baseClassName string
	if len(node.HeritageClauses) > 0 {
		for _, clause := range node.HeritageClauses {
			if len(clause.Types) > 0 {
				typeArg := &clause.Types[0]
				if len(typeArg.Children) > 0 && typeArg.Children[0].Kind == Identifier {
					baseClassName = toPascalCase(typeArg.Children[0].Text)
					break
				}
			}
		}
	}

	// Step 1: Generate the struct definition
	g.writeLine(fmt.Sprintf("type %s struct {", className))
	g.indent++

	// If there's a base class, embed it first
	if baseClassName != "" {
		g.writeLine(fmt.Sprintf("*%s // Embedded base class", baseClassName))
	}

	// Generate fields from PropertyDeclarations
	var constructorNode *ASTNode
	var methods []ASTNode
	var staticFields []ASTNode
	var staticMethods []ASTNode

	if node.Members != nil {
		for i := range node.Members {
			member := &node.Members[i]
			if member.Kind == PropertyDeclaration {
				if isStatic(member) {
					staticFields = append(staticFields, *member)
					continue
				}

				fieldName := toPascalCase(member.Name)
				if isPrivate(member) {
					fieldName = toCamelCase(member.Name)
				}

				fieldType := "interface{}"
				if member.Type != nil {
					var err error
					fieldType, err = g.generateType(member.Type)
					if err != nil {
						return err
					}
				}
				g.writeLine(fmt.Sprintf("%s %s", fieldName, fieldType))
			} else if member.Kind == Constructor {
				constructorNode = member
			} else if member.Kind == MethodDeclaration {
				if isStatic(member) {
					staticMethods = append(staticMethods, *member)
				} else {
					methods = append(methods, *member)
				}
			} else if member.Kind == GetAccessor || member.Kind == SetAccessor {
				methods = append(methods, *member)
			}
		}
	}

	g.indent--
	g.writeLine("}")
	g.writeLine("")

	// Generate static fields as package-level variables
	for _, field := range staticFields {
		fieldName := className + toPascalCase(field.Name)
		fieldType := "interface{}"
		if field.Type != nil {
			var err error
			fieldType, err = g.generateType(field.Type)
			if err != nil {
				return err
			}
		}

		if field.Initializer != nil {
			initValue, err := g.generateExpression(field.Initializer)
			if err != nil {
				return err
			}
			g.writeLine(fmt.Sprintf("var %s %s = %s", fieldName, fieldType, initValue))
		} else {
			g.writeLine(fmt.Sprintf("var %s %s", fieldName, fieldType))
		}
	}
	if len(staticFields) > 0 {
		g.writeLine("")
	}

	// Step 2: Generate constructor function
	if constructorNode != nil {
		if err := g.generateConstructorWithBase(className, baseClassName, constructorNode); err != nil {
			return err
		}
	}

	// Step 3: Generate instance methods
	for _, method := range methods {
		if err := g.generateMethod(className, &method); err != nil {
			return err
		}
	}

	// Step 4: Generate static methods as package-level functions
	for _, method := range staticMethods {
		if err := g.generateStaticMethod(className, &method); err != nil {
			return err
		}
	}

	return nil
}
