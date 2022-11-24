package trace

import (
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/goid"
	_const "github.com/isyscore/isc-tracer/internal/const"
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
