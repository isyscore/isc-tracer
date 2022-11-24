package trace

import (
	"fmt"
	"github.com/isyscore/isc-gobase/store"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/util"
	"net/http"
	"strconv"
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
	return ClientStartTrace(_const.HTTP, fmt.Sprintf("【http】: <%s>%s", method, uri))
}

// ClientStartTrace
// 开启客户端跟踪(如前端访问某个后端接口a, 接口a内访问其他接口b, 此时a访问b称为客户端, b接口内为服务端)
func ClientStartTrace(traceType _const.TraceTypeEnum, traceName string) *Tracer {
	return StartTrace(traceType, traceName, _const.CLIENT)
}

// ServerStartTrace
/**
traceName 名称
 可以是一个 http url
 可以是一个rpc的 service.name
 可以是一个MQ的 send.{topic}.{partition}
 可以是访问redis的 get.{namespace}.{key}
*/
func ServerStartTrace(traceType _const.TraceTypeEnum, traceName string) *Tracer {
	return StartTrace(traceType, traceName, _const.SERVER)
}
func StartTrace(traceType _const.TraceTypeEnum, traceName string, endPoint _const.EndpointEnum) *Tracer {
	header := store.GetHeader()
	remoteAddr := store.GetRemoteAddr()

	tracerId := header.Get(_const.TRACE_HEAD_ID)
	frontIP := ""
	if tracerId == "" {
		tracerId = util.GenerateTraceId()
		frontIP = GetFrontIP(header, remoteAddr)
		if header != nil {
			header.Set(_const.TRACE_HEAD_ID, tracerId)
		}
	}

	rpcId := header.Get(_const.TRACE_HEAD_RPC_ID)
	if rpcId == "" {
		rpcId = ROOT_RPC_ID
	} else {
		// 获取最后一位 +1
		splits := strings.Split(rpcId, ".")
		lastOne, _ := strconv.Atoi(splits[len(splits)-1])
		lastOne += 1
		splits[len(splits)-1] = strconv.Itoa(lastOne)
		rpcId = strings.Join(splits, ".")
	}
	if header != nil {
		header.Set(_const.TRACE_HEAD_RPC_ID, rpcId)
	}

	tracer := doStartTrace(tracerId, rpcId, traceType, traceName, endPoint)
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

func EndTrace(tracer *Tracer, responseSize int, status _const.TraceStatusEnum, message string) {
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
