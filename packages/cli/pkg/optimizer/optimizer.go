package optimizer

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
)

// Optimizer performs various optimizations on generated Go code
type Optimizer struct {
	fset *token.FileSet
}

// NewOptimizer creates a new optimizer
func NewOptimizer() *Optimizer {
	return &Optimizer{
		fset: token.NewFileSet(),
	}
}

// Optimize applies all optimization passes to the generated Go code
func (o *Optimizer) Optimize(code string) (string, error) {
	// Parse the Go code
	file, err := parser.ParseFile(o.fset, "", code, parser.ParseComments)
	if err != nil {
		// If parsing fails, return original code (might be partial/invalid)
		return code, nil
	}

	// Apply optimization passes
	o.removeUnusedImports(file)
	o.removeUnusedVariables(file)
	o.removeUnusedFunctions(file)

	// Format the optimized code
	var buf strings.Builder
	if err := format.Node(&buf, o.fset, file); err != nil {
		return code, err
	}

	return buf.String(), nil
}

// removeUnusedImports removes imports that are not used in the code
func (o *Optimizer) removeUnusedImports(file *ast.File) {
	if file == nil {
		return
	}

	// Track which imports are actually used
	usedImports := make(map[string]bool)

	// Walk the AST to find all identifier references
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectorExpr:
			// Track package usage (e.g., fmt.Println)
			if ident, ok := x.X.(*ast.Ident); ok {
				usedImports[ident.Name] = true
			}
		case *ast.CallExpr:
			// Track direct function calls that might use imports
			if fun, ok := x.Fun.(*ast.Ident); ok {
				usedImports[fun.Name] = true
			}
		}
		return true
	})

	// Remove unused imports
	var newDecls []ast.Decl
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			newDecls = append(newDecls, decl)
			continue
		}

		// Filter import specs
		var newSpecs []ast.Spec
		for _, spec := range genDecl.Specs {
			importSpec := spec.(*ast.ImportSpec)
			
			// Get import name (either explicit or from path)
			var importName string
			if importSpec.Name != nil {
				importName = importSpec.Name.Name
			} else {
				// Extract package name from path
				path := strings.Trim(importSpec.Path.Value, `"`)
				parts := strings.Split(path, "/")
				importName = parts[len(parts)-1]
			}

			// Keep import if it's used or if it's a side-effect import (_)
			if usedImports[importName] || (importSpec.Name != nil && importSpec.Name.Name == "_") {
				newSpecs = append(newSpecs, spec)
			}
		}

		// Only keep the import declaration if there are remaining imports
		if len(newSpecs) > 0 {
			genDecl.Specs = newSpecs
			newDecls = append(newDecls, genDecl)
		}
	}

	file.Decls = newDecls
}

// removeUnusedVariables removes variable declarations that are never used
// Note: This is conservative and only removes obviously unused top-level variables
func (o *Optimizer) removeUnusedVariables(file *ast.File) {
	if file == nil {
		return
	}

	// Track variable usage across the file
	usedVars := make(map[string]bool)

	// Walk the AST to find variable usage (excluding declarations)
	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			// Mark identifier as used if it's not part of a declaration
			if x.Obj == nil {
				// This might be a reference to a variable
				usedVars[x.Name] = true
			}
		case *ast.CallExpr:
			// Check function arguments
			for _, arg := range x.Args {
				if ident, ok := arg.(*ast.Ident); ok {
					usedVars[ident.Name] = true
				}
			}
		case *ast.BinaryExpr:
			// Check binary expression operands
			if ident, ok := x.X.(*ast.Ident); ok {
				usedVars[ident.Name] = true
			}
			if ident, ok := x.Y.(*ast.Ident); ok {
				usedVars[ident.Name] = true
			}
		}
		return true
	})

	// Remove unused top-level variables
	// Be conservative: only remove if clearly unused
	var newDecls []ast.Decl
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			newDecls = append(newDecls, decl)
			continue
		}

		// Filter variable specs
		var newSpecs []ast.Spec
		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)
			keepSpec := false
			
			// Keep the spec if any name is used or exported
			for _, name := range valueSpec.Names {
				if name.Name == "_" {
					keepSpec = true
					break
				}
				// Keep exported variables
				if len(name.Name) > 0 && name.Name[0] >= 'A' && name.Name[0] <= 'Z' {
					keepSpec = true
					break
				}
				// Keep if used
				if usedVars[name.Name] {
					keepSpec = true
					break
				}
			}

			if keepSpec {
				newSpecs = append(newSpecs, spec)
			}
		}

		// Only keep the declaration if there are remaining specs
		if len(newSpecs) > 0 {
			genDecl.Specs = newSpecs
			newDecls = append(newDecls, genDecl)
		}
	}

	file.Decls = newDecls
}

// removeUnusedFunctions removes function declarations that are never called
func (o *Optimizer) removeUnusedFunctions(file *ast.File) {
	if file == nil {
		return
	}

	// Track all function declarations and their usage
	declaredFuncs := make(map[string]*ast.FuncDecl)
	usedFuncs := make(map[string]bool)

	// Always keep main, init, and exported functions
	alwaysKeep := map[string]bool{
		"main": true,
		"init": true,
	}

	// First pass: collect all function declarations
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			funcName := funcDecl.Name.Name
			declaredFuncs[funcName] = funcDecl
			
			// Keep exported functions (start with uppercase)
			if len(funcName) > 0 && funcName[0] >= 'A' && funcName[0] <= 'Z' {
				alwaysKeep[funcName] = true
			}
		}
	}

	// Second pass: find all function calls
	ast.Inspect(file, func(n ast.Node) bool {
		if callExpr, ok := n.(*ast.CallExpr); ok {
			switch fun := callExpr.Fun.(type) {
			case *ast.Ident:
				// Direct function call
				usedFuncs[fun.Name] = true
			case *ast.SelectorExpr:
				// Method call or package function
				if ident, ok := fun.X.(*ast.Ident); ok {
					usedFuncs[ident.Name] = true
				}
			}
		}
		return true
	})

	// Third pass: remove unused functions
	var newDecls []ast.Decl
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			newDecls = append(newDecls, decl)
			continue
		}

		funcName := funcDecl.Name.Name
		
		// Keep function if it's used or should always be kept
		if usedFuncs[funcName] || alwaysKeep[funcName] {
			newDecls = append(newDecls, decl)
		}
	}

	file.Decls = newDecls
}

// OptimizeImports is a convenience function to only optimize imports
func (o *Optimizer) OptimizeImports(code string) (string, error) {
	file, err := parser.ParseFile(o.fset, "", code, parser.ParseComments)
	if err != nil {
		return code, nil
	}

	o.removeUnusedImports(file)

	var buf strings.Builder
	if err := format.Node(&buf, o.fset, file); err != nil {
		return code, err
	}

	return buf.String(), nil
}
