module github.com/el-j/ts2go/internal/orchestrator

go 1.22.5

toolchain go1.24.9

require (
	github.com/el-j/ts2go/internal/module v0.0.0
	github.com/el-j/ts2go/internal/project v0.0.0
	github.com/el-j/ts2go/internal/transpiler v0.0.0
)

replace (
	github.com/el-j/ts2go/internal/module => ../module
	github.com/el-j/ts2go/internal/project => ../project
	github.com/el-j/ts2go/internal/transpiler => ../transpiler
)
