package test

import (
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/internal/trace"
	"github.com/isyscore/isc-tracer/pkg"
	"net/http"
	"testing"
	"time"
)

func TestClientStartTraceWithHeader(t *testing.T) {
	header := &http.Header{}
	tr := pkg.ClientStartTraceWithHeader(header, "")
	time.Sleep(time.Second)
	trace.EndTrace(tr, _const.OK, "", 0)

	t.Log(header)

	time.Sleep(time.Second * 2)
}

func TestTracer(t *testing.T) {
	tr := trace.Tracer{}
	t.Log(tr.ChildRpcSeq.Inc())
}
