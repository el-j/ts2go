module github.com/el-j/ts2go/examples/multipackage-demo

go 1.21

require (
	github.com/el-j/ts2go/internal/orchestrator v0.0.0
	github.com/el-j/ts2go/internal/project v0.0.0
)

replace (
	github.com/el-j/ts2go/internal/orchestrator => ../../internal/orchestrator
	github.com/el-j/ts2go/internal/project => ../../internal/project
)
