module github.com/el-j/ts2go/pkg/cli

go 1.22.5

require (
	github.com/el-j/ts2go/internal/analyzer v0.0.0
	github.com/el-j/ts2go/internal/mapper v0.0.0
	github.com/el-j/ts2go/internal/project v0.0.0
	github.com/el-j/ts2go/internal/transpiler v0.0.0
	github.com/fsnotify/fsnotify v1.9.0
)

require (
	golang.org/x/sys v0.13.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/el-j/ts2go/internal/analyzer => ../../internal/analyzer
	github.com/el-j/ts2go/internal/mapper => ../../internal/mapper
	github.com/el-j/ts2go/internal/project => ../../internal/project
	github.com/el-j/ts2go/internal/transpiler => ../../internal/transpiler
)
