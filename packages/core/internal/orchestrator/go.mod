module github.com/yourusername/ts2go/internal/orchestrator

go 1.21

require (
	github.com/yourusername/ts2go/internal/module v0.0.0
	github.com/yourusername/ts2go/internal/project v0.0.0
	github.com/yourusername/ts2go/internal/transpiler v0.0.0
)

replace (
	github.com/yourusername/ts2go/internal/module => ../module
	github.com/yourusername/ts2go/internal/project => ../project
	github.com/yourusername/ts2go/internal/transpiler => ../transpiler
)
