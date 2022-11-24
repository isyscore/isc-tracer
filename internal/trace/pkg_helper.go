package trace

import (
	"github.com/isyscore/isc-gobase/store"
	_const "github.com/isyscore/isc-tracer/internal/const"
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

var OsTraceSwitch bool
var HttpTraceSwitch bool
var DatabaseTraceSwitch bool
var RedisTraceSwitch bool
var EtcdTraceSwitch bool

func init() {
	OsTraceSwitch = false
	HttpTraceSwitch = false
	DatabaseTraceSwitch = false
	RedisTraceSwitch = false
	EtcdTraceSwitch = false
}

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

func ServerStartTrace(traceType _const.TraceTypeEnum, traceName string) *Tracer {
	header := store.GetHeader()
	remoteAddr := store.GetRemoteAddr()

	tracerId := header.Get(_const.TRACE_HEAD_ID)
	frontIP := ""
	if tracerId == "" {
		tracerId = util.GenerateTraceId()
		frontIP = GetFrontIP(header, remoteAddr)
	}
	rpcId := header.Get(_const.TRACE_HEAD_RPC_ID)
	tracer := StartTrace(tracerId, rpcId, traceType, traceName, _const.SERVER)
	if tracer == nil {
		return nil
	}
	if frontIP != "" {
		tracer.RemoteIp = frontIP
	}
	// 往当前上下文添加远程端属性
	putAttr(tracer, header)
	return tracer
}

func ServerEndTrace(tracer *Tracer, responseSize int, status _const.TraceStatusEnum, message string) {
	endTrace(tracer, responseSize, status, message)
}

func endTrace(tracer *Tracer, responseSize int, status _const.TraceStatusEnum, message string) {
	tracer.Size = responseSize
	tracer.EndTrace(status, message)
}

func putAttr(tracer *Tracer, head http.Header) {
	if tracer.AttrMap == nil {
		tracer.AttrMap = make(map[string]string)
	}
	for key, copyKey := range copyAttrMap {
		if v := head.Get(key); v != "" {
			tracer.AttrMap[copyKey] = v
		}
	}
}
