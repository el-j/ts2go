module github.com/yourusername/ts2go/internal/project

go 1.21

require (
	github.com/yourusername/ts2go/internal/analyzer v0.0.0
	github.com/yourusername/ts2go/internal/mapper v0.0.0
)

replace github.com/yourusername/ts2go/internal/analyzer => ../analyzer

replace github.com/yourusername/ts2go/internal/mapper => ../mapper
