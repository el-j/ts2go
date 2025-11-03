package module

// ExportType represents the type of export
type ExportType string

const (
	// ExportNamed is a named export: export { foo, bar }
	ExportNamed ExportType = "named"
	// ExportDefault is a default export: export default foo
	ExportDefault ExportType = "default"
	// ExportAll is a re-export all: export * from './foo'
	ExportAll ExportType = "all"
	// ExportAllAs is a re-export all as namespace: export * as foo from './bar'
	ExportAllAs ExportType = "all-as"
)

// ImportType represents the type of import
type ImportType string

const (
	// ImportNamed is a named import: import { foo, bar } from './module'
	ImportNamed ImportType = "named"
	// ImportDefault is a default import: import foo from './module'
	ImportDefault ImportType = "default"
	// ImportNamespace is a namespace import: import * as foo from './module'
	ImportNamespace ImportType = "namespace"
	// ImportSideEffect is a side-effect import: import './module'
	ImportSideEffect ImportType = "side-effect"
)

// Export represents a symbol exported from a module
type Export struct {
	// Name is the exported symbol name
	Name string
	// Type is the export type (named, default, all, all-as)
	Type ExportType
	// LocalName is the local name in the module (for renamed exports)
	LocalName string
	// Source is the source module for re-exports (e.g., './other')
	Source string
	// IsType indicates if this is a type-only export
	IsType bool
	// Position is the line number where the export occurs
	Position int
}

// Import represents a symbol imported into a module
type Import struct {
	// Name is the imported symbol name (local name)
	Name string
	// ImportedName is the name in the source module (for renamed imports)
	ImportedName string
	// Type is the import type (named, default, namespace, side-effect)
	Type ImportType
	// Source is the module being imported (e.g., './user', 'axios')
	Source string
	// IsType indicates if this is a type-only import
	IsType bool
	// Position is the line number where the import occurs
	Position int
}

// Module represents a TypeScript module with its imports and exports
type Module struct {
	// Path is the file path of the module
	Path string
	// PackageName is the Go package name for this module
	PackageName string
	// Imports is the list of imports in this module
	Imports []Import
	// Exports is the list of exports from this module
	Exports []Export
	// Dependencies is a list of module paths this module depends on
	Dependencies []string
}

// GetPackageName returns the package name (for interface compatibility)
func (m *Module) GetPackageName() string {
	return m.PackageName
}

// ExportRegistry tracks all exports across the project
type ExportRegistry struct {
	// modules maps file path to Module
	modules map[string]*Module
}

// NewExportRegistry creates a new export registry
func NewExportRegistry() *ExportRegistry {
	return &ExportRegistry{
		modules: make(map[string]*Module),
	}
}

// AddModule adds a module to the registry
func (r *ExportRegistry) AddModule(module *Module) {
	r.modules[module.Path] = module
}

// GetModule retrieves a module by path
func (r *ExportRegistry) GetModule(path string) *Module {
	return r.modules[path]
}

// GetExports returns all exports from a module
func (r *ExportRegistry) GetExports(path string) []Export {
	if module := r.modules[path]; module != nil {
		return module.Exports
	}
	return nil
}

// GetImports returns all imports from a module
func (r *ExportRegistry) GetImports(path string) []Import {
	if module := r.modules[path]; module != nil {
		return module.Imports
	}
	return nil
}

// AllModules returns all registered modules
func (r *ExportRegistry) AllModules() []*Module {
	modules := make([]*Module, 0, len(r.modules))
	for _, module := range r.modules {
		modules = append(modules, module)
	}
	return modules
}
