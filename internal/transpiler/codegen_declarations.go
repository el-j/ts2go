package transpiler

import "fmt"

// generateInterface converts a TypeScript interface to a Go struct
func (g *CodeGenerator) generateInterface(node *ASTNode) error {
	structName := toPascalCase(node.Name)
	g.writeLine(fmt.Sprintf("type %s struct {", structName))
	g.indent++

	if node.Members != nil {
		for _, member := range node.Members {
			if member.Kind == PropertySignature {
				fieldName := toPascalCase(member.Name)
				fieldType, err := g.generateType(member.Type)
				if err != nil {
					return err
				}
				jsonTag := fmt.Sprintf("`json:\"%s\"`", member.Name)
				g.writeLine(fmt.Sprintf("%s %s %s", fieldName, fieldType, jsonTag))
			}
		}
	}

	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// generateTypeAlias converts a TypeScript type alias to a Go type
func (g *CodeGenerator) generateTypeAlias(node *ASTNode) error {
	typeName := toPascalCase(node.Name)

	// Check if this is a union type
	if node.Type != nil && node.Type.Kind == UnionType {
		return g.generateUnionType(typeName, node.Type)
	}

	// Regular type alias
	typeValue, err := g.generateType(node.Type)
	if err != nil {
		return err
	}
	g.writeLine(fmt.Sprintf("type %s = %s", typeName, typeValue))
	g.writeLine("")
	return nil
}

// generateEnum converts a TypeScript enum to Go const declarations
func (g *CodeGenerator) generateEnum(node *ASTNode) error {
	enumName := toPascalCase(node.Name)

	// Check if this is a string enum (has string values) or numeric enum
	isStringEnum := g.isStringEnum(node)

	if isStringEnum {
		return g.generateStringEnum(enumName, node)
	} else {
		return g.generateNumericEnum(enumName, node)
	}
}

// isStringEnum determines if an enum has string values
func (g *CodeGenerator) isStringEnum(node *ASTNode) bool {
	if node.Members != nil {
		for _, member := range node.Members {
			if member.Kind == EnumMember && member.Initializer != nil {
				// Check if initializer is a string literal
				if member.Initializer.Kind == StringLiteral {
					return true
				}
			}
		}
	}
	return false
}

// generateNumericEnum generates a numeric enum using iota
func (g *CodeGenerator) generateNumericEnum(enumName string, node *ASTNode) error {
	g.writeLine(fmt.Sprintf("type %s int", enumName))
	g.writeLine("")
	g.writeLine("const (")

	if node.Members != nil {
		for i, member := range node.Members {
			if member.Kind == EnumMember {
				memberName := fmt.Sprintf("%s%s", enumName, toPascalCase(member.Name))
				if i == 0 {
					g.writeLine(fmt.Sprintf("%s %s = iota", memberName, enumName))
				} else if member.Initializer != nil {
					// Handle explicit values
					if member.Initializer.Kind == NumericLiteral {
						g.writeLine(fmt.Sprintf("%s = %s", memberName, member.Initializer.Text))
					}
				} else {
					g.writeLine(fmt.Sprintf("%s", memberName))
				}
			}
		}
	}

	g.writeLine(")")
	g.writeLine("")

	// Add String() method
	g.generateEnumStringMethod(enumName, node)
	return nil
}

// generateStringEnum generates a string enum
func (g *CodeGenerator) generateStringEnum(enumName string, node *ASTNode) error {
	g.writeLine(fmt.Sprintf("type %s string", enumName))
	g.writeLine("")
	g.writeLine("const (")

	if node.Members != nil {
		for _, member := range node.Members {
			if member.Kind == EnumMember {
				memberName := fmt.Sprintf("%s%s", enumName, toPascalCase(member.Name))
				if member.Initializer != nil && member.Initializer.Kind == StringLiteral {
					g.writeLine(fmt.Sprintf(`%s %s = "%s"`, memberName, enumName, member.Initializer.Text))
				}
			}
		}
	}

	g.writeLine(")")
	g.writeLine("")

	// Add String() method
	g.generateEnumStringMethod(enumName, node)
	return nil
}

// generateEnumStringMethod adds a String() method for enum debugging
func (g *CodeGenerator) generateEnumStringMethod(enumName string, node *ASTNode) {
	g.writeLine(fmt.Sprintf("func (e %s) String() string {", enumName))
	g.indent++
	g.writeLine("switch e {")

	if node.Members != nil {
		for _, member := range node.Members {
			if member.Kind == EnumMember {
				memberName := fmt.Sprintf("%s%s", enumName, toPascalCase(member.Name))
				if member.Initializer != nil && member.Initializer.Kind == StringLiteral {
					g.writeLine(fmt.Sprintf(`case %s:`, memberName))
					g.indent++
					g.writeLine(fmt.Sprintf(`return "%s"`, member.Initializer.Text))
					g.indent--
				} else {
					// For numeric enums, return the name
					g.writeLine(fmt.Sprintf("case %s:", memberName))
					g.indent++
					g.writeLine(fmt.Sprintf(`return "%s"`, member.Name))
					g.indent--
				}
			}
		}
	}

	g.writeLine("default:")
	g.indent++
	g.writeLine(`return "Unknown"`)
	g.indent--
	g.writeLine("}")
	g.indent--
	g.writeLine("}")
	g.writeLine("")
}
