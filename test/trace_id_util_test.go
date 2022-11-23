package test

import (
	"github.com/isyscore/isc-tracer/util"
	"testing"
)

func TestGenerateTraceId(t *testing.T) {
	t.Log(util.GenerateTraceId())
}
