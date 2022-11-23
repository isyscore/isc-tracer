package pkg

import (
	"github.com/isyscore/isc-gobase/server"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/internal/trace"
	"github.com/isyscore/isc-tracer/util"
	"net/http"
	"strings"
)

var (
	copyAttrMap = map[string]string{
		_const.TRACE_HEAD_REMOTE_APPNAME: _const.TRACE_HEAD_REMOTE_APPNAME,
		_const.TRACE_HEAD_REMOTE_IP:      _const.TRACE_HEAD_REMOTE_IP,
		_const.TRACE_HEAD_USER_ID:        _const.A_USER_ID,
		_const.TRACE_HEAD_USER_NAME:      _const.A_USER_NAME,
	}
)

func GetFrontIP(head http.Header, remoteAddr string) string {
	ip := head.Get("X-Forwarded-For")
	if ip != "" && strings.EqualFold(ip, "unKnown") {
		//多次反向代理后会有多个ip值，第一个ip才是真实ip
		if i := strings.Index(ip, ","); i != -1 {
			return ip[:i]
		}
		return ip
	}
	ip = head.Get("X-Real-IP")
	if ip != "" && strings.EqualFold(ip, "unKnown") {
		return ip
	}
	return remoteAddr
}

func ClientStartTrace(traceType _const.TraceTypeEnum, traceName string) *trace.Tracer {
	header := server.GetHeader()
	remoteAddr := server.GetRemoteAddr()

	tracerId := header.Get(_const.TRACE_HEAD_ID)
	frontIP := ""
	if tracerId == "" {
		tracerId = util.GenerateTraceId()
		frontIP = GetFrontIP(header, remoteAddr)
	}
	rpcId := header.Get(_const.TRACE_HEAD_RPC_ID)
	tracer := trace.StartTrace(tracerId, rpcId, traceType, traceName, _const.CLIENT)
	if frontIP != "" {
		tracer.RemoteIp = frontIP
	}
	putAttr(tracer, header)
	return tracer
}

func ClientEndTrace(tracer *trace.Tracer, responseSize int, status _const.TraceStatusEnum, message string) {
	endTrace(tracer, responseSize, status, message)
}

func ServerStartTrace(traceType _const.TraceTypeEnum, traceName string) *trace.Tracer {
	header := server.GetHeader()
	remoteAddr := server.GetRemoteAddr()

	tracerId := header.Get(_const.TRACE_HEAD_ID)
	frontIP := ""
	if tracerId == "" {
		tracerId = util.GenerateTraceId()
		frontIP = GetFrontIP(header, remoteAddr)
	}
	rpcId := header.Get(_const.TRACE_HEAD_RPC_ID)
	tracer := trace.StartTrace(tracerId, rpcId, traceType, traceName, _const.SERVER)
	if frontIP != "" {
		tracer.RemoteIp = frontIP
	}
	// 往当前上下文添加远程端属性
	putAttr(tracer, header)
	return tracer
}

func ServerEndTrace(tracer *trace.Tracer, responseSize int, status _const.TraceStatusEnum, message string) {
	endTrace(tracer, responseSize, status, message)
}

func endTrace(tracer *trace.Tracer, responseSize int, status _const.TraceStatusEnum, message string) {
	req := server.GetRequest()
	tracer.Size = int(req.ContentLength) + responseSize
	tracer.EndTrace(status, message)
}

func putAttr(tracer *trace.Tracer, head http.Header) {
	if tracer.AttrMap == nil {
		tracer.AttrMap = make(map[string]string)
	}
	for key, copyKey := range copyAttrMap {
		if v := head.Get(key); v != "" {
			tracer.AttrMap[copyKey] = v
		}
	}
}
