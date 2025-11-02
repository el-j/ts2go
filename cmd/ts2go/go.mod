module github.com/yourusername/ts2go/cmd/ts2go

go 1.22.5

require (
	github.com/yourusername/ts2go/internal/transpiler v0.0.0
	github.com/yourusername/ts2go/pkg/cli v0.0.0
)

replace github.com/yourusername/ts2go/internal/transpiler => ../../internal/transpiler

replace github.com/yourusername/ts2go/pkg/cli => ../../pkg/cli

replace github.com/yourusername/ts2go/internal/analyzer => ../../internal/analyzer

replace github.com/yourusername/ts2go/internal/mapper => ../../internal/mapper

replace github.com/yourusername/ts2go/internal/project => ../../internal/project

replace github.com/yourusername/ts2go/internal/optimizer => ../../internal/optimizer

replace github.com/yourusername/ts2go/pkg/optimizer => ../../pkg/optimizer
