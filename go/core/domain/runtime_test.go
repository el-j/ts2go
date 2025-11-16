package domain_test

import (
	"testing"

	"github.com/el-j/ts2go/core/domain"
)

func TestNewBuildResult(t *testing.T) {
	outputPath := "/output/binary"

	result := domain.NewBuildResult(outputPath)

	if result.OutputPath != outputPath {
		t.Errorf("Expected output path %s, got %s", outputPath, result.OutputPath)
	}
	if result.Output != "" {
		t.Error("Expected empty output initially")
	}
}
