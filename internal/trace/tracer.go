package trace

import (
	"fmt"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/goid"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/util"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const ROOT_RPC_ID = "0"

// 携程上存储tracer
var localStore = goid.NewLocalStorage()

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

func StartTrace(traceId string, rpcId string, traceType _const.TraceTypeEnum, traceName string, endpoint _const.EndpointEnum) *Tracer {
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

	if rpcId == "" {
		rpcId = ROOT_RPC_ID
	}
	tracer = &Tracer{
		TraceId:   traceId,
		RpcId:     rpcId,
		TraceType: traceType,
		Sampled:   true,
	}
	tracer.startTrace(traceName, endpoint)
	return tracer
}

func (tracer *Tracer) startTrace(traceName string, endpoint _const.EndpointEnum) {
	tracer.TraceName = traceName
	tracer.Endpoint = endpoint
	tracer.StartTime = time.Now().UnixMilli()
}

func (tracer *Tracer) EndTrace(status _const.TraceStatusEnum, message string) {
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

// NewServerTracer 开启服务端跟踪
func NewServerTracer(req *http.Request) *Tracer {
	if !TracerIsEnable() {
		return nil
	}
	tracer := New(req)
	tracer.Endpoint = _const.SERVER
	return tracer
}

// NewServerTracerWithoutReq 开启服务端跟踪,此用于服务端定时任务类请求
func NewServerTracerWithoutReq() *Tracer {
	if !TracerIsEnable() {
		return nil
	}
	tracer := &Tracer{
		TraceId:   util.GenerateTraceId(),
		Sampled:   true,
		TraceName: config.GetValueString("base.application.name"),
		StartTime: time.Now().UnixMilli(),
		RpcId:     "0",
		TraceType: _const.HTTP,
		RemoteIp:  util.GetLocalIp(),
	}
	return tracer
}

var clientTracerLock sync.Mutex

func (tracer *Tracer) NewClientWithHeader(header *http.Header) *Tracer {
	if !TracerIsEnable() {
		return nil
	}

	clientTracerLock.Lock()
	defer clientTracerLock.Unlock()
	rpcId := tracer.RpcId
	if rpcId == "" {
		rpcId = tracer.RpcId
		rpcId += ".1"
	} else {
		// 获取最后一位 +1
		splits := strings.Split(rpcId, ".")
		lastOne, _ := strconv.Atoi(splits[len(splits)-1])
		lastOne += 1
		splits[len(splits)-1] = strconv.Itoa(lastOne)
		rpcId = strings.Join(splits, ".")
	}
	tracer.RpcId = rpcId
	// fixme TraceName和Size 需要手动写入
	clientTracer := &Tracer{
		TraceId:   tracer.TraceId,
		Sampled:   true,
		TraceName: config.GetValueString("base.application.name"),
		StartTime: time.Now().UnixMilli(),
		RpcId:     rpcId,
		TraceType: _const.HTTP,
		RemoteIp:  util.GetLocalIp(),
	}
	header.Set(_const.TRACE_HEAD_ID, tracer.TraceId)
	header.Set(_const.TRACE_HEAD_RPC_ID, rpcId)
	return clientTracer
}

// NewClientTracer 开启客户端跟踪
func (tracer *Tracer) NewClientTracer(req *http.Request) *Tracer {
	if !TracerIsEnable() {
		return nil
	}

	clientTracerLock.Lock()
	defer clientTracerLock.Unlock()
	rpcId := tracer.RpcId
	if rpcId == "" {
		rpcId = tracer.RpcId
		rpcId += ".1"
	} else {
		// 获取最后一位 +1
		splits := strings.Split(rpcId, ".")
		lastOne, _ := strconv.Atoi(splits[len(splits)-1])
		lastOne += 1
		splits[len(splits)-1] = strconv.Itoa(lastOne)
		rpcId = strings.Join(splits, ".")
	}

	clientTracer := NewWithRpcId(req, rpcId)
	clientTracer.TraceId = tracer.TraceId
	clientTracer.Endpoint = _const.CLIENT
	tracer.RpcId = rpcId
	return clientTracer
}

// NewWithRpcId 自定义rpcId
func NewWithRpcId(req *http.Request, rpcId string) *Tracer {
	if !TracerIsEnable() {
		return nil
	}

	tracer := New(req)
	req.Header.Set(_const.TRACE_HEAD_RPC_ID, rpcId)
	tracer.RpcId = rpcId
	return tracer
}

func New(req *http.Request) *Tracer {
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
	strLength := req.Header.Get("Content-Length")
	if strLength == "" {
		strLength = "0"
	}
	length, _ := strconv.Atoi(strLength)
	return &Tracer{
		TraceId:   getOrCreateTraceId(req),
		Sampled:   true,
		StartTime: time.Now().UnixMilli(),
		RpcId:     getAndIncreaseRpcId(req),
		TraceType: _const.HTTP,
		RemoteIp:  req.RemoteAddr,
		TraceName: fmt.Sprintf("<%s>%s", method, uri),
		AttrMap:   make(map[string]string),
		Size:      length,
	}
}

func getOrCreateTraceId(req *http.Request) string {
	traceId := req.Header.Get(_const.TRACE_HEAD_ID)
	if traceId == "" {
		traceId = util.GenerateTraceId()
		if req.Header != nil {
			req.Header.Set(_const.TRACE_HEAD_ID, traceId)
		}
	}
	return traceId
}

func getAndIncreaseRpcId(req *http.Request) string {
	rpcId := req.Header.Get(_const.TRACE_HEAD_RPC_ID)
	if rpcId == "" {
		rpcId = "0"
	}
	if req.Header != nil {
		req.Header.Set(_const.TRACE_HEAD_RPC_ID, rpcId)
	}
	return rpcId
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
	return config.GetValueBoolDefault("tracer.enable", true) && OsTraceSwitch
}
