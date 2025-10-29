module github.com/ts2go/tests

go 1.21

require (
	github.com/ts2go/transpiler v0.0.0
	github.com/ts2go/runtime v0.0.0
)

replace (
	github.com/ts2go/transpiler => ../internal/transpiler
	github.com/ts2go/runtime => ../runtime
)