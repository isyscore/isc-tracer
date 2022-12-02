package trace

import "github.com/isyscore/isc-gobase/goid"

// 协程上存储tracer
var localStore = goid.NewLocalStorage()

func createCurrentTracerIfAbsent() *Tracer {
	l := localStore.Get()
	if l == nil {
		return &Tracer{}
	}
	return l.(*Tracer)
}

func GetCurrentTracer() *Tracer {
	l := localStore.Get()
	if l == nil {
		return nil
	}
	return l.(*Tracer)
}

func setTrace(tracer *Tracer) {
	localStore.Set(tracer)
}
func deleteTrace() {
	localStore.Del()
}
