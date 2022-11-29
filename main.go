package isc_tracer

import (
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/debug"
	"github.com/isyscore/isc-gobase/extend/etcd"
	"github.com/isyscore/isc-gobase/extend/orm"
	baseRedis "github.com/isyscore/isc-gobase/extend/redis"
	baseHttp "github.com/isyscore/isc-gobase/http"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/listener"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/isyscore/isc-gobase/server"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/internal/trace"
	pkgEtcd "github.com/isyscore/isc-tracer/pkg/etcd"
	pkgHttp "github.com/isyscore/isc-tracer/pkg/http"
	pkgOrm "github.com/isyscore/isc-tracer/pkg/orm"
	pkgRedis "github.com/isyscore/isc-tracer/pkg/redis"
	"net/http"
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
		register()
	})

	// 应用退出
	listener.AddListener(listener.EventOfServerStop, func(event listener.BaseEvent) {
		trace.Close()
	})
}

func register() {
	if !config.GetValueBoolDefault("tracer.debug.enable", true) {
		return
	}

	// 获取etcd账号和密码
	corebackAddr := config.GetValueStringDefault("tracer.debug.account", _const.CORE_BACK_ADDRESS)
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
		trace.SwitchTrace = isc.ToBool(value)
	})
	debug.AddWatcher(SWITCH_OS_TRACE_DATABASE, func(key string, value string) {
		logger.Info("配置最新值：key:【%v】, value：【%v】", key, value)
		trace.SwitchTraceDatabase = isc.ToBool(value)
	})
	debug.AddWatcher(SWITCH_OS_TRACE_REDIS, func(key string, value string) {
		logger.Info("配置最新值：key:【%v】, value：【%v】", key, value)
		trace.SwitchTraceRedis = isc.ToBool(value)
	})
	debug.AddWatcher(SWITCH_OS_TRACE_ETCD, func(key string, value string) {
		logger.Info("配置最新值：key:【%v】, value：【%v】", key, value)
		trace.SwitchTraceEtcd = isc.ToBool(value)
	})
	debug.StartWatch()
}
