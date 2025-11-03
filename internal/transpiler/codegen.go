package transpiler

import (
	"strings"
)

// CodeGenerator generates Go code from AST nodes
type CodeGenerator struct {
	output                    strings.Builder
	indent                    int
	currentFunctionReturnType string          // Track the current function's return type for type assertions
	currentReceiverVar        string          // Track the current method's receiver variable for "this" replacement
	currentClassMembers       map[string]bool // Track private members of current class (name -> isPrivate)
	tempVarCounter            int             // Counter for generating unique temporary variables

	// Module system support
	module     interface{} // *module.Module - using interface{} to avoid circular dependency
	resolver   interface{} // *module.ImportResolver
	visibility interface{} // *module.SymbolVisibility
	isEntry    bool        // Is this the entry point (main package)?
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{
		indent:                    0,
		currentFunctionReturnType: "",
		currentReceiverVar:        "",
	}
}

// NewCodeGeneratorWithModule creates a code generator with module system support
func NewCodeGeneratorWithModule(mod, resolver, visibility interface{}, isEntry bool) *CodeGenerator {
	return &CodeGenerator{
		indent:                    0,
		currentFunctionReturnType: "",
		currentReceiverVar:        "",
		module:                    mod,
		resolver:                  resolver,
		visibility:                visibility,
		isEntry:                   isEntry,
	}
}
