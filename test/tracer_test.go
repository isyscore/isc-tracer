package test

import (
	trace2 "github.com/isyscore/isc-tracer/trace"
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
	tr := trace2.Tracer{}
	t.Log(tr.ChildRpcSeq.Inc())
}
