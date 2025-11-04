module github.com/el-j/ts2go/internal/project

go 1.22.5

toolchain go1.24.9

require (
	github.com/el-j/ts2go/internal/analyzer v0.0.0
	github.com/el-j/ts2go/internal/mapper v0.0.0
)

require gopkg.in/yaml.v3 v3.0.1 // indirect

replace github.com/el-j/ts2go/internal/analyzer => ../analyzer

replace github.com/el-j/ts2go/internal/mapper => ../mapper
