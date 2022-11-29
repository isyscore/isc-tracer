package trace

import (
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/goid"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/isyscore/isc-gobase/store"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/util"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const ROOT_RPC_ID = "0"

// 携程上存储tracer
var localStore = goid.NewLocalStorage()

var (
	copyAttrMap = map[string]string{
		_const.TRACE_HEAD_REMOTE_APPNAME: _const.TRACE_HEAD_REMOTE_APPNAME,
		_const.TRACE_HEAD_REMOTE_IP:      _const.TRACE_HEAD_REMOTE_IP,
		_const.TRACE_HEAD_USER_ID:        _const.A_USER_ID,
		_const.TRACE_HEAD_USER_NAME:      _const.A_USER_NAME,
	}
)

var (
	SwitchTrace         = false
	SwitchTraceDatabase = false
	SwitchTraceRedis    = false
	SwitchTraceEtcd     = false
)

type Tracer struct {
	// TraceId 调用链ID,一旦初始化,不能修改
	TraceId string
	// RpcId 调用顺序，依次为0 → 0.1 → 0.1.1,1 -> 1.1 -> 1.1.1 ...
	RpcId string
	// TraceType 链路跟踪类型
	TraceType _const.TraceTypeEnum
	/**
	 * 名称
	 * 可以是一个 http url
	 * 可以是一个rpc的 service.name
	 * 可以是一个MQ的 send.{topic}.{partition}
	 * 可以是访问redis的 get.{namespace}.{key}
	 */
	TraceName string
	// Endpoint 跟踪类型
	Endpoint _const.EndpointEnum
	// status 跟踪结果
	status _const.TraceStatusEnum

	// RemoteStatus 远程调用结果
	RemoteStatus _const.TraceStatusEnum
	// RemoteIp 远程调用IP,即下游(Client)或上游(Server)ip
	RemoteIp string
	// message 调用返回或异常信息
	message string
	// Size 响应体大小
	Size int

	// StartTime 当前span开始时间
	StartTime int64
	// endTime 当前span结束时间
	endTime int64
	// 是否采样
	Sampled bool
	// bizData 响应数据
	bizData map[string]any
	// 是否已经结束
	Ended bool
	// AttrMap 请求参数
	AttrMap map[string]string
}

func doStartTrace(traceId string, rpcId string, traceType _const.TraceTypeEnum, traceName string, endpoint _const.EndpointEnum) *Tracer {
	if !TracerIsEnable() {
		return nil
	}
	tracer := createCurrentTracerIfAbsent()
	if tracer.Ended {
		if tracer.TraceId == traceId {
			return tracer
		}
	} else if tracer.TraceId != "" {
		return tracer
	}

	tracer = &Tracer{
		TraceId:   traceId,
		RpcId:     rpcId,
		TraceType: traceType,
		TraceName: traceName,
		Endpoint:  endpoint,
		Sampled:   true,
		StartTime: time.Now().UnixMilli(),
	}
	tracer.startTrace()

	localStore.Set(tracer)
	return tracer
}

func (tracer *Tracer) startTrace() {
	tracer.Sampled = true
	tracer.StartTime = time.Now().UnixMilli()
	tracer.AttrMap = make(map[string]string)
}

func (tracer *Tracer) EndTrace(status _const.TraceStatusEnum, message string, responseSize int) {
	if !TracerIsEnable() || tracer.Ended {
		return
	}
	if tracer.TraceId == "" || tracer.RpcId == "" || tracer.StartTime == 0 {
		//log.Println("tracer's traceId is nil,will be not append tracer info")
		return
	}
	tracer.Ended = true
	if tracer.getStatus() == _const.OK && !tracer.Sampled {
		return
	}

	tracer.endTime = time.Now().UnixMilli()
	tracer.status = status
	tracer.Size = responseSize
	if message != "" {
		tracer.message = message
	}
	SendTraceLog(tracer)

	localStore.Del()
}

func createCurrentTracerIfAbsent() *Tracer {
	l := localStore.Get()
	if l == nil {
		tracer := &Tracer{}
		localStore.Set(tracer)
		return tracer
	}
	return l.(*Tracer)
}

func (tracer *Tracer) getStatus() _const.TraceStatusEnum {
	if tracer.status != _const.OK {
		return tracer.status
	}
	if tracer.RemoteStatus != _const.OK {
		return tracer.RemoteStatus
	}
	return _const.OK
}

func TracerIsEnable() bool {
	return config.GetValueBoolDefault("tracer.enable", true) && SwitchTrace
}

func StartTrace(traceType _const.TraceTypeEnum, endPoint _const.EndpointEnum, traceName string, header *http.Header) *Tracer {
	if !TracerIsEnable() {
		return nil
	}
	remoteAddr := store.GetRemoteAddr()
	if header == nil {
		h := store.GetHeader()
		header = &h
	}
	tracerId := header.Get(_const.TRACE_HEAD_ID)
	frontIP := ""
	if tracerId == "" {
		tracerId = util.GenerateTraceId()
		frontIP = GetFrontIP(header, remoteAddr)
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

	if *header != nil {
		header.Set(_const.TRACE_HEAD_ID, tracerId)
		header.Set(_const.TRACE_HEAD_RPC_ID, rpcId)
	}

	store.RequestHeadAdd(_const.TRACE_HEAD_ID, tracerId)
	store.RequestHeadAdd(_const.TRACE_HEAD_RPC_ID, rpcId)

	logger.PutMdc(_const.TRACE_HEAD_ID, tracerId)
	logger.PutMdc(_const.TRACE_HEAD_RPC_ID, rpcId)

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

func EndTrace(tracer *Tracer, status _const.TraceStatusEnum, message string, responseSize int) {
	tracer.EndTrace(status, message, responseSize)
}

func putAttr(tracer *Tracer, head *http.Header) {
	if tracer.AttrMap == nil {
		tracer.AttrMap = make(map[string]string)
	}
	for key, copyKey := range copyAttrMap {
		if v := head.Get(key); v != "" {
			tracer.AttrMap[copyKey] = v
		}
	}
}

func GetFrontIP(head *http.Header, remoteAddr string) string {
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
