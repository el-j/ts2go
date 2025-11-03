module github.com/el-j/ts2go/cmd/ts2go

go 1.22.5

require github.com/el-j/ts2go/pkg/cli v0.0.0

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)

replace github.com/el-j/ts2go/pkg/cli => ../../pkg/cli
