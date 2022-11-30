package trace

import "github.com/isyscore/isc-gobase/goid"

// 携程上存储tracer
var localStore = goid.NewLocalStorage()

func createCurrentTracerIfAbsent() *Tracer {
	l := localStore.Get()
	if l == nil {
		return &Tracer{}
	}
	for _, tracer := range l.(map[string]*Tracer) {
		return tracer
	}
	return &Tracer{}
}

func setTrace(rpcId string, tracer *Tracer) {
	l := localStore.Get()
	if l == nil {
		n := make(map[string]*Tracer)
		localStore.Set(n)
		l = localStore.Get()
	}
	dict := l.(map[string]*Tracer)
	dict[rpcId] = tracer
}
func deleteTrace(rpcId string) {
	l := localStore.Get()
	if l == nil {
		return
	}
	dict := l.(map[string]*Tracer)
	delete(dict, rpcId)
	if len(dict) == 0 {
		localStore.Del()
	}
}
