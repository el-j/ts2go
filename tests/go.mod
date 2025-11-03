module github.com/ts2go/tests

go 1.22.5

toolchain go1.24.9

replace (
	github.com/ts2go/runtime => ../runtime
	github.com/ts2go/transpiler => ../internal/transpiler
)
