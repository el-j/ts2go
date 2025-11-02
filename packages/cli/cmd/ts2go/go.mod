module github.com/el-j/ts2go/packages/cli/cmd/ts2go

go 1.22.5

require (
	github.com/el-j/ts2go/packages/core/internal/transpiler v0.0.0
	github.com/el-j/ts2go/packages/cli/pkg/cli v0.0.0
)

replace github.com/el-j/ts2go/packages/core/internal/transpiler => ../../../core/internal/transpiler

replace github.com/el-j/ts2go/packages/cli/pkg/cli => ../../pkg/cli

replace github.com/el-j/ts2go/packages/core/internal/analyzer => ../../../core/internal/analyzer

replace github.com/el-j/ts2go/packages/core/internal/mapper => ../../../core/internal/mapper

replace github.com/el-j/ts2go/packages/core/internal/project => ../../../core/internal/project

replace github.com/el-j/ts2go/packages/core/internal/optimizer => ../../../core/internal/optimizer

replace github.com/el-j/ts2go/packages/cli/pkg/optimizer => ../../pkg/optimizer
