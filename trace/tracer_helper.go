package trace

import (
	"fmt"
	_const2 "github.com/isyscore/isc-tracer/const"
	"net/http"
)

// ClientStartTrace
// 开启客户端跟踪(如前端访问某个后端接口a, 接口a内访问其他接口b, 此时a访问b称为客户端, b接口内为服务端)
func ClientStartTrace(traceType _const2.TraceTypeEnum, traceName string) *Tracer {
	if !TracerIsEnable() {
		return nil
	}
	return StartTrace(traceType, _const2.CLIENT, traceName, nil)
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
	return StartTrace(_const2.HTTP, _const2.CLIENT, fmt.Sprintf("<%s>%s", method, uri), req)
}

func ServerStartTrace(traceType _const2.TraceTypeEnum, traceName string) *Tracer {
	if !TracerIsEnable() {
		return nil
	}
	return StartTrace(traceType, _const2.SERVER, traceName, nil)
}

func ServerStartTraceWithRequest(traceType _const2.TraceTypeEnum, traceName string, request *http.Request) *Tracer {
	if !TracerIsEnable() {
		return nil
	}
	return StartTrace(traceType, _const2.SERVER, traceName, request)
}

func EndTraceOk(tracer *Tracer, message string, responseSize int) {
	tracer.EndTrace(_const2.OK, message, responseSize)
}

func EndTraceTimeout(tracer *Tracer, message string, responseSize int) {
	tracer.EndTrace(_const2.TIMEOUT, message, responseSize)
}

func EndTraceWarn(tracer *Tracer, message string, responseSize int) {
	tracer.EndTrace(_const2.WARNING, message, responseSize)
}

func EndTraceError(tracer *Tracer, message string, responseSize int) {
	tracer.EndTrace(_const2.ERROR, message, responseSize)
}
