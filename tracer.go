package isc_tracer

import (
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/debug"
	"github.com/isyscore/isc-gobase/extend/etcd"
	"github.com/isyscore/isc-gobase/extend/orm"
	baseRedis "github.com/isyscore/isc-gobase/extend/redis"
	"github.com/isyscore/isc-gobase/goid"
	baseHttp "github.com/isyscore/isc-gobase/http"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/listener"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/isyscore/isc-gobase/server"
	_const2 "github.com/isyscore/isc-tracer/const"
	pkgEtcd "github.com/isyscore/isc-tracer/pkg/etcd"
	pkgHttp "github.com/isyscore/isc-tracer/pkg/http"
	pkgOrm "github.com/isyscore/isc-tracer/pkg/orm"
	pkgRedis "github.com/isyscore/isc-tracer/pkg/redis"
	trace2 "github.com/isyscore/isc-tracer/trace"
	"net/http"
	"time"
)

const (
	SWITCH_OS_TRACE          = "debug.os.trace"
	SWITCH_OS_TRACE_DATABASE = "debug.os.trace.database"
	SWITCH_OS_TRACE_REDIS    = "debug.os.trace.redis"
	SWITCH_OS_TRACE_ETCD     = "debug.os.trace.etcd"
	SWITCH_OS_TRACE_KAFKA    = "debug.os.trace.kafka"
	SWITCH_OS_TRACE_EMQX     = "debug.os.trace.emqx"
)

func init() {
	server.AddGinHandlers(pkgHttp.TraceFilter())
	orm.AddGormHook(&pkgOrm.TracerGormHook{})
	orm.AddXormHook(&pkgOrm.TracerXormHook{})
	baseRedis.AddRedisHook(&pkgRedis.TracerRedisHook{})
	etcd.AddEtcdHook(&pkgEtcd.TracerEtcdHook{})
	baseHttp.AddHook(&pkgHttp.TracerHttpHook{})

	// 应用启动完成
	listener.AddListener(listener.EventOfServerRunFinish, func(event listener.BaseEvent) {
		goid.Go(func() {
			registerAppName()
		})

		goid.Go(func() {
			registerWatch()
		})
	})

	// 应用退出
	listener.AddListener(listener.EventOfServerStop, func(event listener.BaseEvent) {
		trace2.Close()
	})
}

func registerAppName() {
	pivotUrl := config.GetValueStringDefault("tracer.server.admin-url", "http://isc-pivot-platform:31107")
	header := http.Header{}
	parameter := map[string]string{}
	body := map[string]any{
		"appCode": config.GetValueStringDefault("base.application.name", _const2.DEFAULT_APP_NAME),
	}

	for i := 0; i < 100; i++ {
		if !trace2.IsHealthOfAdmin() {
			time.Sleep(5 * time.Second)
			continue
		}
		_, _, _, err := baseHttp.Put(pivotUrl+"/api/app/operation-center/tracer/register", header, parameter, body)
		if err != nil {
			logger.Warn("注册pivot异常，重试，%v", err.Error())
			time.Sleep(5 * time.Second)
		} else {
			logger.Info("注册pivot成功")
			break
		}
	}
}

func registerWatch() {
	if !config.GetValueBoolDefault("tracer.debug.enable", true) {
		return
	}

	// 获取etcd账号和密码
	corebackAddr := config.GetValueStringDefault("tracer.debug.account", _const2.CORE_BACK_ADDRESS)
	header := http.Header{}
	parameterMap := map[string]string{}
	_, _, data, err := baseHttp.GetOfStandard(corebackAddr+"/api/core/back/account/etcd", header, parameterMap)
	if err != nil {
		logger.Error("获取core-back服务异常，调试模式不支持：%v", err.Error())
		return
	}

	etcdAccountMap := map[string]any{}
	err = isc.DataToObject(data, &etcdAccountMap)
	if err != nil {
		logger.Error("转换失败：%v", err.Error())
		return
	}

	etcdEndpoints, exist := etcdAccountMap["endpoints"]
	etcdUser, _ := etcdAccountMap["username"]
	etcdPassword, _ := etcdAccountMap["password"]

	if !exist {
		logger.Error("etcd账号获取为空")
		return
	}

	var endpoints []string
	for _, endpoint := range etcdEndpoints.([]interface{}) {
		endpoints = append(endpoints, endpoint.(string))
	}

	debug.InitWithParameter(endpoints, etcdUser.(string), etcdPassword.(string))
	debug.AddWatcher(SWITCH_OS_TRACE, func(key string, value string) {
		logger.Info("配置最新值：key:【%v】, value：【%v】", key, value)
		trace2.SwitchTrace = isc.ToBool(value)
	})
	debug.AddWatcher(SWITCH_OS_TRACE_DATABASE, func(key string, value string) {
		logger.Info("配置最新值：key:【%v】, value：【%v】", key, value)
		trace2.SwitchTraceDatabase = isc.ToBool(value)
	})
	debug.AddWatcher(SWITCH_OS_TRACE_REDIS, func(key string, value string) {
		logger.Info("配置最新值：key:【%v】, value：【%v】", key, value)
		trace2.SwitchTraceRedis = isc.ToBool(value)
	})
	debug.AddWatcher(SWITCH_OS_TRACE_ETCD, func(key string, value string) {
		logger.Info("配置最新值：key:【%v】, value：【%v】", key, value)
		trace2.SwitchTraceEtcd = isc.ToBool(value)
	})
	debug.StartWatch()
}
