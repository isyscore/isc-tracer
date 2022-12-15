package trace

import (
	"github.com/isyscore/isc-gobase/goid"
	"github.com/isyscore/isc-gobase/isc"
)

// 协程上存储tracer
var tracerStorage = goid.NewLocalStorage()

func createCurrentTracerIfAbsent() *Tracer {
	l := tracerStorage.Get()
	if l == nil {
		return &Tracer{}
	}
	tracerMap := l.(isc.OrderMap[string, *Tracer])
	if tracerMap.Size() != 0 {
		return tracerMap.GetValue(tracerMap.Size() - 1)
	}
	return &Tracer{}
}

func GetCurrentTracer() *Tracer {
	l := tracerStorage.Get()
	if l == nil {
		return nil
	}
	tracerMap := l.(isc.OrderMap[string, *Tracer])
	if tracerMap.Size() != 0 {
		return tracerMap.GetValue(tracerMap.Size() - 1)
	}
	return nil
}

func setTrace(tracer *Tracer) {
	if tracer == nil {
		return
	}
	l := tracerStorage.Get()
	if l == nil {
		tracerStorage.Set(isc.NewOrderMap[string, *Tracer]())
		l = tracerStorage.Get()
	}
	dict := l.(isc.OrderMap[string, *Tracer])
	dict.Put(tracer.RpcId, tracer)
}

func deleteTrace(rpcId string) {
	l := tracerStorage.Get()
	if l == nil {
		return
	}
	dict := l.(isc.OrderMap[string, *Tracer])
	dict.Delete(rpcId)
	if dict.Size() == 0 {
		tracerStorage.Del()
	} else {
		tracerStorage.Set(dict)
	}
}
