package trace

import (
	"github.com/isyscore/isc-gobase/goid"
	"testing"
)

func TestName(t *testing.T) {
	//tracer := Tracer{}
	//tracer.AttrMap["k"] = "123"

	localStore := goid.NewLocalStorage()
	t2 := localStore.Get().(*Tracer)
	t.Log(t2)
}
