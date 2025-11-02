package transpiler

import (
	"fmt"
	"strings"
)

// generateType converts a TypeScript type node to its Go equivalent
func (g *CodeGenerator) generateType(node *ASTNode) (string, error) {
	if node == nil {
		return "interface{}", nil
	}

	switch node.Kind {
	case StringKeyword:
		return "string", nil
	case NumberKeyword:
		return "float64", nil
	case BooleanKeyword:
		return "bool", nil
	case VoidKeyword:
		return "", nil // void functions have no return type in Go
	case ArrayType:
		if len(node.Children) > 0 {
			elemType, err := g.generateType(&node.Children[0])
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("[]%s", elemType), nil
		}
		return "[]interface{}", nil
	case TupleType:
		return g.generateTupleType(node)
	case TypeLiteral:
		// Inline object type -> struct
		return g.generateInlineStruct(node)
	case TypeReference:
		// Reference to a named type (interface, type alias, etc.)
		if len(node.Children) > 0 && node.Children[0].Kind == Identifier {
			return toPascalCase(node.Children[0].Text), nil
		}
		return "interface{}", nil
	case UnionType:
		// Union types are handled at the type alias level
		return "interface{}", nil
	default:
		return "interface{}", nil
	}
}

// generateTupleType generates an inline struct for a tuple type
func (g *CodeGenerator) generateTupleType(node *ASTNode) (string, error) {
	// For now, generate an inline struct with numbered fields
	// In the future, this could generate a named struct for reusability
	fields := []string{}

	if node.Elements != nil {
		for i, element := range node.Elements {
			fieldType, err := g.generateType(&element)
			if err != nil {
				return "", err
			}
			fields = append(fields, fmt.Sprintf("Field%d %s", i, fieldType))
		}
	}

	if len(fields) == 0 {
		return "struct{}", nil
	}

	return fmt.Sprintf("struct{ %s }", strings.Join(fields, "; ")), nil
}

// generateInlineStruct converts a TypeLiteral to an inline struct definition
func (g *CodeGenerator) generateInlineStruct(node *ASTNode) (string, error) {
	if node.Members == nil || len(node.Members) == 0 {
		return "struct{}", nil
	}

	fields := []string{}
	for _, member := range node.Members {
		if member.Kind == PropertySignature {
			fieldName := toPascalCase(member.Name)
			fieldType, err := g.generateType(member.Type)
			if err != nil {
				return "", err
			}
			jsonTag := fmt.Sprintf("`json:\"%s\"`", member.Name)
			fields = append(fields, fmt.Sprintf("%s %s %s", fieldName, fieldType, jsonTag))
		}
	}

	return fmt.Sprintf("struct{ %s }", strings.Join(fields, "; ")), nil
}

// generateUnionType creates a discriminated union struct for TypeScript union types
func (g *CodeGenerator) generateUnionType(typeName string, unionNode *ASTNode) error {
	// Generate the discriminated union struct
	g.writeLine(fmt.Sprintf("type %s struct {", typeName))
	g.indent++
	g.writeLine("Type  string")
	g.writeLine("Value interface{}")
	g.indent--
	g.writeLine("}")
	g.writeLine("")

	// Generate constructor functions for each union variant
	if unionNode.Types != nil {
		for _, unionType := range unionNode.Types {
			if err := g.generateUnionConstructor(typeName, &unionType); err != nil {
				return err
			}
		}
	}

	// Generate type guard methods
	if unionNode.Types != nil {
		for _, unionType := range unionNode.Types {
			if err := g.generateTypeGuard(typeName, &unionType); err != nil {
				return err
			}
		}
	}

	return nil
}

// generateUnionConstructor creates a constructor function for a union variant
func (g *CodeGenerator) generateUnionConstructor(unionTypeName string, variantType *ASTNode) error {
	variantName := g.getTypeName(variantType)
	constructorName := fmt.Sprintf("New%s%s", unionTypeName, variantName)

	g.writeLine(fmt.Sprintf("func %s(value %s) %s {", constructorName, g.getGoType(variantType), unionTypeName))
	g.indent++
	g.writeLine(fmt.Sprintf("return %s{Type: \"%s\", Value: value}", unionTypeName, variantName))
	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// generateTypeGuard creates a type guard method for a union variant
func (g *CodeGenerator) generateTypeGuard(unionTypeName string, variantType *ASTNode) error {
	variantName := g.getTypeName(variantType)
	methodName := fmt.Sprintf("Is%s", variantName)

	g.writeLine(fmt.Sprintf("func (u %s) %s() bool {", unionTypeName, methodName))
	g.indent++
	g.writeLine(fmt.Sprintf("return u.Type == \"%s\"", variantName))
	g.indent--
	g.writeLine("}")
	g.writeLine("")

	// Also generate a getter method
	getterName := fmt.Sprintf("As%s", variantName)
	g.writeLine(fmt.Sprintf("func (u %s) %s() %s {", unionTypeName, getterName, g.getGoType(variantType)))
	g.indent++
	g.writeLine(fmt.Sprintf("return u.Value.(%s)", g.getGoType(variantType)))
	g.indent--
	g.writeLine("}")
	g.writeLine("")
	return nil
}

// getTypeName extracts a readable name from a type node
func (g *CodeGenerator) getTypeName(node *ASTNode) string {
	switch node.Kind {
	case TypeReference:
		if len(node.Children) > 0 && node.Children[0].Kind == Identifier {
			return toPascalCase(node.Children[0].Text)
		}
		if node.Text != "" {
			return toPascalCase(node.Text)
		}
	case LiteralType:
		// For literal types, use the string value
		if len(node.Children) > 0 && node.Children[0].Kind == StringLiteral {
			return toPascalCase(node.Children[0].Text)
		}
		return "Literal"
	case StringKeyword:
		return "String"
	case NumberKeyword:
		return "Number"
	case BooleanKeyword:
		return "Bool"
	default:
		return "Unknown"
	}
	return "Unknown"
}

// getGoType converts a type node to its Go equivalent
func (g *CodeGenerator) getGoType(node *ASTNode) string {
	switch node.Kind {
	case TypeReference:
		if len(node.Children) > 0 && node.Children[0].Kind == Identifier {
			return toPascalCase(node.Children[0].Text)
		}
		if node.Text != "" {
			return toPascalCase(node.Text)
		}
	case LiteralType:
		// String literals are still strings in Go
		return "string"
	case StringKeyword:
		return "string"
	case NumberKeyword:
		return "float64"
	case BooleanKeyword:
		return "bool"
	default:
		return "interface{}"
	}
	return "interface{}"
}
