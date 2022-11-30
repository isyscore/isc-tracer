package test

import (
	"github.com/isyscore/isc-tracer/internal/trace"
	"testing"
)

//func TestClientStartTraceWithHeader(t *testing.T) {
//	header := &http.Header{}
//	tr := pkg.ClientStartTraceWithHeader(header, "")
//	time.Sleep(time.Second)
//	trace.EndTrace(tr, _const.OK, "", 0)
//
//	t.Log(header)
//
//	time.Sleep(time.Second * 2)
//}

func TestTracer(t *testing.T) {
	tr := trace.Tracer{}
	t.Log(tr.ChildRpcSeq.Inc())
}
