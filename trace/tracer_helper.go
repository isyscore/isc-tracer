package trace

import (
	"fmt"
	_const "github.com/isyscore/isc-tracer/const"
	"net/http"
)

// ClientStartTrace
// 开启客户端跟踪(如前端访问某个后端接口a, 接口a内访问其他接口b, 此时a访问b称为客户端, b接口内为服务端)
func ClientStartTrace(traceType _const.TraceTypeEnum, traceName string) *Tracer {
	if !TracerIsEnable() {
		return nil
	}
	return StartTrace(traceType, _const.CLIENT, traceName, nil)
}

func ClientStartTraceWithHeader(header *http.Header, traceName string) *Tracer {
	if !TracerIsEnable() {
		return nil
	}
	if traceName == "" {
		traceName = "<default>_server"
	}
	return StartTraceWithHeader(_const.HTTP, _const.CLIENT, traceName, header)
}

func ClientStartTraceWithRequest(req *http.Request) *Tracer {
	if !TracerIsEnable() {
		return nil
	}

	method := req.Method
	if method == "" {
		method = "nil"
	}
	uri := "nil"
	if url := req.URL; url != nil {
		if uri = url.Path; len(uri) == 0 {
			uri = url.String()
		}
	}
	return StartTrace(_const.HTTP, _const.CLIENT, fmt.Sprintf("<%s>%s", method, uri), req)
}

func ServerStartTrace(traceType _const.TraceTypeEnum, traceName string) *Tracer {
	if !TracerIsEnable() {
		return nil
	}
	return StartTrace(traceType, _const.SERVER, traceName, nil)
}

func ServerStartTraceWithRequest(traceType _const.TraceTypeEnum, traceName string, request *http.Request) *Tracer {
	if !TracerIsEnable() {
		return nil
	}
	return StartTrace(traceType, _const.SERVER, traceName, request)
}

func EndTraceOk(tracer *Tracer, message string, responseSize int) {
	if tracer == nil {
		return
	}
	tracer.EndTrace(_const.OK, message, responseSize)
}

func EndTraceTimeout(tracer *Tracer, message string, responseSize int) {
	if tracer == nil {
		return
	}
	tracer.EndTrace(_const.TIMEOUT, message, responseSize)
}

func EndTraceWarn(tracer *Tracer, message string, responseSize int) {
	if tracer == nil {
		return
	}
	tracer.EndTrace(_const.WARNING, message, responseSize)
}

func EndTraceError(tracer *Tracer, message string, responseSize int) {
	if tracer == nil {
		return
	}
	tracer.EndTrace(_const.ERROR, message, responseSize)
}

func DiscardTrace(tracer *Tracer) {
	if tracer == nil {
		return
	}
	deleteTrace(tracer.RpcId)
}
