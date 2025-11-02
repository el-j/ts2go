module github.com/yourusername/ts2go/pkg/cli

go 1.22.5

require (
	github.com/yourusername/ts2go/internal/analyzer v0.0.0
	github.com/yourusername/ts2go/internal/mapper v0.0.0
	github.com/yourusername/ts2go/internal/project v0.0.0
	github.com/yourusername/ts2go/internal/transpiler v0.0.0
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace (
	github.com/yourusername/ts2go/internal/analyzer => ../../internal/analyzer
	github.com/yourusername/ts2go/internal/mapper => ../../internal/mapper
	github.com/yourusername/ts2go/internal/project => ../../internal/project
	github.com/yourusername/ts2go/internal/transpiler => ../../internal/transpiler
)
