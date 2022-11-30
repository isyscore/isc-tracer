package pkg

import (
	"fmt"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/internal/trace"
	"net/http"
)

// ClientStartTrace
// 开启客户端跟踪(如前端访问某个后端接口a, 接口a内访问其他接口b, 此时a访问b称为客户端, b接口内为服务端)
func ClientStartTrace(traceType _const.TraceTypeEnum, traceName string) *trace.Tracer {
	if !trace.TracerIsEnable() {
		return nil
	}
	return trace.StartTrace(traceType, _const.CLIENT, traceName, nil)
}

func ClientStartTraceWithRequest(req *http.Request) *trace.Tracer {
	if !trace.TracerIsEnable() {
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
	return trace.StartTrace(_const.HTTP, _const.CLIENT, fmt.Sprintf("<%s>%s", method, uri), req)
}

func ServerStartTrace(traceType _const.TraceTypeEnum, traceName string) *trace.Tracer {
	if !trace.TracerIsEnable() {
		return nil
	}
	return trace.StartTrace(traceType, _const.SERVER, traceName, nil)
}

func ServerStartTraceWithRequest(traceType _const.TraceTypeEnum, traceName string, request *http.Request) *trace.Tracer {
	if !trace.TracerIsEnable() {
		return nil
	}
	return trace.StartTrace(traceType, _const.SERVER, traceName, request)
}

func EndTraceOk(tracer *trace.Tracer, message string, responseSize int) {
	tracer.EndTrace(_const.OK, message, responseSize)
}

func EndTraceTimeout(tracer *trace.Tracer, message string, responseSize int) {
	tracer.EndTrace(_const.TIMEOUT, message, responseSize)
}

func EndTraceWarn(tracer *trace.Tracer, message string, responseSize int) {
	tracer.EndTrace(_const.WARNING, message, responseSize)
}

func EndTraceError(tracer *trace.Tracer, message string, responseSize int) {
	tracer.EndTrace(_const.ERROR, message, responseSize)
}
