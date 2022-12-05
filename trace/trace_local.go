package trace

import "github.com/isyscore/isc-gobase/goid"

// 协程上存储tracer
var tracerStorage = goid.NewLocalStorage()

func createCurrentTracerIfAbsent() *Tracer {
	l := tracerStorage.Get()
	if l == nil {
		return &Tracer{}
	}
	for _, tracer := range l.(map[string]*Tracer) {
		return tracer
	}
	return &Tracer{}
}

func GetCurrentTracer() *Tracer {
	l := tracerStorage.Get()
	if l == nil {
		return nil
	}
	for _, tracer := range l.(map[string]*Tracer) {
		return tracer
	}
	return nil
}

func setTrace(tracer *Tracer) {
	if tracer == nil {
		return
	}
	l := tracerStorage.Get()
	if l == nil {
		tracerStorage.Set(make(map[string]*Tracer))
		l = tracerStorage.Get()
	}
	dict := l.(map[string]*Tracer)
	dict[tracer.RpcId] = tracer
}

func deleteTrace(rpcId string) {
	l := tracerStorage.Get()
	if l == nil {
		return
	}
	dict := l.(map[string]*Tracer)
	delete(dict, rpcId)
	if len(dict) == 0 {
		tracerStorage.Del()
	} else {
		tracerStorage.Set(dict)
	}
}
