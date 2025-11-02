module github.com/el-j/ts2go/packages/cli/pkg/cli

go 1.22.5

require (
	github.com/el-j/ts2go/packages/core/internal/analyzer v0.0.0
	github.com/el-j/ts2go/packages/core/internal/mapper v0.0.0
	github.com/el-j/ts2go/packages/core/internal/project v0.0.0
	github.com/el-j/ts2go/packages/core/internal/transpiler v0.0.0
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace (
	github.com/el-j/ts2go/packages/core/internal/analyzer => ../../../../core/internal/analyzer
	github.com/el-j/ts2go/packages/core/internal/mapper => ../../../../core/internal/mapper
	github.com/el-j/ts2go/packages/core/internal/project => ../../../../core/internal/project
	github.com/el-j/ts2go/packages/core/internal/transpiler => ../../../../core/internal/transpiler
)
