package trace

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/goid"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/logger"
	baseNet "github.com/isyscore/isc-gobase/system/net"
	_const "github.com/isyscore/isc-tracer/const"
	"github.com/isyscore/isc-tracer/pivot"
	"github.com/isyscore/isc-tracer/util"
	"github.com/robfig/cron"
	"google.golang.org/grpc"
)

var reqQueue *isc.Queue
var logCron *cron.Cron
var serverIsHealth = false
var serverAdminIsHealth = false
var serverService pivot.PivotServiceClient

func init() {
	if !TracerIsEnable() {
		return
	}

	reqQueue = isc.NewQueue()
	if logCron != nil {
		logCron.Stop()
	} else {
		logCron = cron.New()
	}

	serverUrl := config.GetValueStringDefault("tracer.server.url", "isc-pivot-platform:31108")
	conn, err := grpc.Dial(serverUrl, grpc.WithInsecure())
	if err != nil {
		logger.Error("连接服务端失败: server：%s；错误信息：%v", serverUrl, err.Error())
		return
	}
	serverService = pivot.NewPivotServiceClient(conn)

	// 检查服务端健康状况
	CheckServerHealth()

	// 每5秒检查pivot健康情况
	err = logCron.AddFunc("0/5 * * * * ?", CheckServerHealth)
	err = logCron.AddFunc("0/5 * * * * ?", CheckAdminServerHealth)
	// 每3秒上报tracer信息
	err = logCron.AddFunc("0/3 * * * * ?", UploadTracer)
	if err != nil {
		logger.Error("定时任务添加失败 error : %v", err)
		return
	}
	logCron.Start()
}

func IsHealth() bool {
	return serverIsHealth
}

func IsHealthOfAdmin() bool {
	return serverAdminIsHealth
}

func UploadTracer() {
	goid.Go(func() {
		numLeft := reqQueue.Num()
		for numLeft != 0 {
			logger.Group("tracer").Debug("数据准备发送，index=%v", numLeft)
			dataReq, numLeftTem := reqQueue.Poll()
			numLeft = numLeftTem
			if dataReq == nil {
				continue
			}
			// 发送到远端
			tracer := dataReq.(*Tracer)
			_, err := serverService.CollectTracer(context.Background(), changeToGrpcTracerRequest(tracer))
			if nil != err {
				logger.Error("链路上报服务端失败, traceId:%s, rpcId:%s, %v", tracer.TraceId, tracer.RpcId, err.Error())
			}
		}
	})
}

func CheckServerHealth() {
	serverUrl := config.GetValueStringDefault("tracer.server.url", "isc-pivot-platform:31108")
	if baseNet.IpPortAvailable(serverUrl) {
		serverIsHealth = true
	} else {
		serverIsHealth = false
	}
}

func CheckAdminServerHealth() {
	serverUrl := config.GetValueStringDefault("tracer.server.admin-url", "http://isc-pivot-platform:31107")
	if baseNet.IpPortAvailable(serverUrl) {
		serverAdminIsHealth = true
	} else {
		serverAdminIsHealth = false
	}
}

func SendTracerToServer(tracer *Tracer) {
	// 直接加入到队列
	reqQueue.Offer(tracer)

	num := reqQueue.Num()
	thresholdValue := config.GetValueInt32Default("pivot.trace.queue.load-size", 1000)
	if num < thresholdValue {
		return
	}

	logger.Info("超过队列阈值【%v】直接发送", thresholdValue)
	// 处理请求队列
	goid.Go(func() {
		UploadTracer()
	})

	// 判断队列个数是否大于等于触发的阈值，是则触发处理
	maxValue := config.GetValueInt32Default("pivot.trace.queue.max-size", 4096)
	if num < maxValue {
		return
	}

	// 大于最大值，则阻塞等待
	logger.Warn("【！！！！】超过队列最大值【%v】，阻塞等待", maxValue)
	dataFinish := make(chan int)
	goid.Go(func() {
		UploadTracer()
		dataFinish <- 1
	})
	select {
	case _ = <-dataFinish:
		{
			return
		}
	}
}

func changeToGrpcTracerRequest(pTracer *Tracer) *pivot.TraceLogRequest {
	userId := pTracer.AttrMap[_const.TRACE_HEAD_USER_ID]
	if userId == "" {
		userId = pTracer.AttrMap[_const.A_USER_ID]
	}
	return &pivot.TraceLogRequest{
		TraceId: pTracer.TraceId,
		RpcId: pTracer.RpcId,
		TraceType: isc.ToInt32(pTracer.TraceType),
		TraceName: pTracer.TraceName,
		Endpoint: isc.ToInt32(pTracer.Endpoint),
		Status: isc.ToInt32(pTracer.Status),
		RemoteStatus: isc.ToInt32(pTracer.RemoteStatus),
		RemoteIp:     pTracer.RemoteIp,
		Message:      pTracer.Message,
		Size:         pTracer.Size,
		StartTime:    pTracer.StartTime,
		EndTime:      pTracer.EndTime,
		Sampled:      pTracer.Sampled,
		BizData:      generateBytesMap(pTracer.bizData),
		Ended:        pTracer.Ended,
		AttrMap:      pTracer.AttrMap,

		ProfilesActive: _const.DEFAULT_PROFILES_ACTIVE,
		AppName:        getAppName(),
		Ip:             util.GetLocalIp(),
		Rt:             int32(pTracer.EndTime - pTracer.StartTime),
		UserId:         userId,
		Sql:            pTracer.AttrMap[_const.A_CMD],
	}
}

func generateBytesMap(anyMap map[string]any) map[string][]byte {
	resultMap := map[string][]byte{}
	for k, v := range anyMap {
		resultV, err := getBytes(v)
		if nil != err {
			continue
		}
		resultMap[k] = resultV
	}
	return resultMap
}

func getBytes(key any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
