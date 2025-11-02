module github.com/yourusername/ts2go/internal/transpiler

go 1.22.5

require github.com/yourusername/ts2go/pkg/optimizer v0.0.0

replace (
	github.com/yourusername/ts2go/internal/module => ../module
	github.com/yourusername/ts2go/internal/project => ../project
	github.com/yourusername/ts2go/pkg/optimizer => ../../pkg/optimizer
)
