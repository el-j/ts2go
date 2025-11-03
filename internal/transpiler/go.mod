module github.com/el-j/ts2go/internal/transpiler

go 1.22.5

require github.com/el-j/ts2go/pkg/optimizer v0.0.0

replace (
	github.com/el-j/ts2go/internal/module => ../module
	github.com/el-j/ts2go/internal/project => ../project
	github.com/el-j/ts2go/pkg/optimizer => ../../pkg/optimizer
)
