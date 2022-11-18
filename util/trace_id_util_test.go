package util

import "testing"

func TestGenerateTraceId(t *testing.T) {
	t.Log(GenerateTraceId())
}
