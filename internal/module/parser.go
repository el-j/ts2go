package module

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Parser parses TypeScript imports and exports
type Parser struct {
	registry *ExportRegistry
}

// NewParser creates a new module parser
func NewParser() *Parser {
	return &Parser{
		registry: NewExportRegistry(),
	}
}

// GetRegistry returns the export registry
func (p *Parser) GetRegistry() *ExportRegistry {
	return p.registry
}

// ParseFile parses a TypeScript file and extracts imports/exports
func (p *Parser) ParseFile(filePath string) (*Module, error) {
	// Find parser directory - look for it relative to the workspace
	// This needs to be configurable or detected from the environment
	parserDir := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filePath))), "internal", "transpiler", "parser")

	// Check if parser directory exists, if not try alternate location
	if _, err := os.Stat(parserDir); os.IsNotExist(err) {
		// Try from current working directory
		cwd, _ := os.Getwd()
		parserDir = filepath.Join(cwd, "..", "transpiler", "parser")

		// If still not found, use absolute path (for tests)
		if _, err := os.Stat(parserDir); os.IsNotExist(err) {
			parserDir = "/Users/rex-fab-alt/Documents/code/playground/ts2go/internal/transpiler/parser"
		}
	}

	cmd := exec.Command("node", "-e", fmt.Sprintf(`
		const ts = require('typescript');
		const fs = require('fs');
		
		const filePath = %s;
		const source = fs.readFileSync(filePath, 'utf8');
		const sourceFile = ts.createSourceFile(filePath, source, ts.ScriptTarget.Latest, true);
		
		const imports = [];
		const exports = [];
		
		function visit(node) {
			// Import declarations
			if (ts.isImportDeclaration(node)) {
				const moduleSpecifier = node.moduleSpecifier.text;
				
				if (node.importClause) {
					// Default import: import foo from './module'
					if (node.importClause.name) {
						imports.push({
							name: node.importClause.name.text,
							importedName: 'default',
							type: 'default',
							source: moduleSpecifier,
							isType: false,
							position: sourceFile.getLineAndCharacterOfPosition(node.pos).line + 1
						});
					}
					
					// Named bindings
					if (node.importClause.namedBindings) {
						// Namespace import: import * as foo from './module'
						if (ts.isNamespaceImport(node.importClause.namedBindings)) {
							imports.push({
								name: node.importClause.namedBindings.name.text,
								importedName: '*',
								type: 'namespace',
								source: moduleSpecifier,
								isType: false,
								position: sourceFile.getLineAndCharacterOfPosition(node.pos).line + 1
							});
						}
						// Named imports: import { foo, bar as baz } from './module'
						else if (ts.isNamedImports(node.importClause.namedBindings)) {
							for (const element of node.importClause.namedBindings.elements) {
								imports.push({
									name: element.name.text,
									importedName: element.propertyName ? element.propertyName.text : element.name.text,
									type: 'named',
									source: moduleSpecifier,
									isType: false,
									position: sourceFile.getLineAndCharacterOfPosition(node.pos).line + 1
								});
							}
						}
					}
				} else {
					// Side-effect import: import './module'
					imports.push({
						name: '',
						importedName: '',
						type: 'side-effect',
						source: moduleSpecifier,
						isType: false,
						position: sourceFile.getLineAndCharacterOfPosition(node.pos).line + 1
					});
				}
			}
			
			// Export declarations
			if (ts.isExportDeclaration(node)) {
				if (node.exportClause) {
					if (ts.isNamedExports(node.exportClause)) {
						// Named exports: export { foo, bar as baz }
						for (const element of node.exportClause.elements) {
							exports.push({
								name: element.name.text,
								localName: element.propertyName ? element.propertyName.text : element.name.text,
								type: 'named',
								source: node.moduleSpecifier ? node.moduleSpecifier.text : '',
								isType: false,
								position: sourceFile.getLineAndCharacterOfPosition(node.pos).line + 1
							});
						}
					} else if (ts.isNamespaceExport(node.exportClause)) {
						// Re-export as namespace: export * as foo from './module'
						exports.push({
							name: node.exportClause.name.text,
							localName: '*',
							type: 'all-as',
							source: node.moduleSpecifier ? node.moduleSpecifier.text : '',
							isType: false,
							position: sourceFile.getLineAndCharacterOfPosition(node.pos).line + 1
						});
					}
				} else if (node.moduleSpecifier) {
					// Re-export all: export * from './module'
					exports.push({
						name: '*',
						localName: '*',
						type: 'all',
						source: node.moduleSpecifier.text,
						isType: false,
						position: sourceFile.getLineAndCharacterOfPosition(node.pos).line + 1
					});
				}
			}
			
			// Export assignments: export default foo
			if (ts.isExportAssignment(node)) {
				exports.push({
					name: 'default',
					localName: node.expression.text || 'default',
					type: 'default',
					source: '',
					isType: false,
					position: sourceFile.getLineAndCharacterOfPosition(node.pos).line + 1
				});
			}
			
			// Exported declarations: export interface Foo, export function bar()
			if (node.modifiers && node.modifiers.some(m => m.kind === ts.SyntaxKind.ExportKeyword)) {
				let name = '';
				let isType = false;
				let exportType = 'named';
				
				if (ts.isFunctionDeclaration(node) || ts.isClassDeclaration(node)) {
					name = node.name ? node.name.text : '';
				} else if (ts.isInterfaceDeclaration(node) || ts.isTypeAliasDeclaration(node)) {
					name = node.name.text;
					isType = true;
				} else if (ts.isVariableStatement(node)) {
					// Handle: export const foo = ...
					for (const decl of node.declarationList.declarations) {
						if (ts.isIdentifier(decl.name)) {
							exports.push({
								name: decl.name.text,
								localName: decl.name.text,
								type: 'named',
								source: '',
								isType: false,
								position: sourceFile.getLineAndCharacterOfPosition(node.pos).line + 1
							});
						}
					}
					return; // Already handled
				} else if (ts.isEnumDeclaration(node)) {
					name = node.name.text;
				}
				
				// Check for default export
				if (node.modifiers.some(m => m.kind === ts.SyntaxKind.DefaultKeyword)) {
					exportType = 'default';
					name = name || 'default';
				}
				
				if (name) {
					exports.push({
						name: name,
						localName: name,
						type: exportType,
						source: '',
						isType: isType,
						position: sourceFile.getLineAndCharacterOfPosition(node.pos).line + 1
					});
				}
			}
			
			ts.forEachChild(node, visit);
		}
		
		visit(sourceFile);
		
		console.log(JSON.stringify({ imports, exports }, null, 2));
	`, fmt.Sprintf(`'%s'`, filePath)))

	cmd.Dir = parserDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to parse TypeScript file: %w\nOutput: %s", err, string(output))
	}

	// Parse the JSON output
	var result struct {
		Imports []Import `json:"imports"`
		Exports []Export `json:"exports"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse parser output: %w\nOutput: %s", err, string(output))
	}

	// Extract dependencies (unique import sources)
	depMap := make(map[string]bool)
	for _, imp := range result.Imports {
		if imp.Source != "" && isRelativeImport(imp.Source) {
			depMap[imp.Source] = true
		}
	}

	dependencies := make([]string, 0, len(depMap))
	for dep := range depMap {
		dependencies = append(dependencies, dep)
	}

	module := &Module{
		Path:         filePath,
		Imports:      result.Imports,
		Exports:      result.Exports,
		Dependencies: dependencies,
	}

	p.registry.AddModule(module)
	return module, nil
}

// isRelativeImport checks if an import is a relative path
func isRelativeImport(source string) bool {
	return strings.HasPrefix(source, "./") || strings.HasPrefix(source, "../")
}

// ParseProject parses all TypeScript files in a project
func (p *Parser) ParseProject(files []string) error {
	for _, file := range files {
		if _, err := p.ParseFile(file); err != nil {
			return fmt.Errorf("failed to parse %s: %w", file, err)
		}
	}
	return nil
}
