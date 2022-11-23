package test

import (
	"github.com/isyscore/isc-gobase/goid"
	"github.com/isyscore/isc-tracer/internal/trace"
	"testing"
)

func TestName(t *testing.T) {
	//tracer := Tracer{}
	//tracer.AttrMap["k"] = "123"

	localStore := goid.NewLocalStorage()
	t2 := localStore.Get().(*trace.Tracer)
	t.Log(t2)
}
