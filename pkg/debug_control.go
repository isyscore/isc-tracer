package pkg

import (
	"github.com/isyscore/isc-gobase/debug"
	"github.com/isyscore/isc-gobase/extend/etcd"
	"github.com/isyscore/isc-gobase/extend/orm"
	baseRedis "github.com/isyscore/isc-gobase/extend/redis"
	"github.com/isyscore/isc-gobase/http"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/listener"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/isyscore/isc-gobase/server"
	"github.com/isyscore/isc-tracer/internal/trace"
	pkgEtcd "github.com/isyscore/isc-tracer/pkg/etcd"
	pkgHttp "github.com/isyscore/isc-tracer/pkg/http"
	pkgOrm "github.com/isyscore/isc-tracer/pkg/orm"
	pkgRedis "github.com/isyscore/isc-tracer/pkg/redis"
)

const (
	SWITCH_OS_TRACE          = "debug.os.trace"
	SWITCH_OS_TRACE_HTTP     = "debug.os.trace.http"
	SWITCH_OS_TRACE_DATABASE = "debug.os.trace.database"
	SWITCH_OS_TRACE_REDIS    = "debug.os.trace.redis"
	SWITCH_OS_TRACE_ETCD     = "debug.os.trace.etcd"
	SWITCH_OS_TRACE_KAFKA    = "debug.os.trace.kafka"
	SWITCH_OS_TRACE_EMQX     = "debug.os.trace.emqx"
)

func init() {
	server.AddGinHandlers(pkgHttp.TraceFilter())
	orm.AddGormHook(&pkgOrm.GobaseGormHook{})
	orm.AddXormHook(&pkgOrm.GobaseXormHook{})
	baseRedis.AddRedisHook(&pkgRedis.GoBaseRedisHook{})
	etcd.AddEtcdHook(&pkgEtcd.TracerEtcdHook{})
	http.AddHook(&pkgHttp.TracerHttpHook{})

	// 应用启动完成
	listener.AddListener(listener.EventOfServerRunFinish, func(event listener.BaseEvent) {
		register()
	})
}

func register() {
	debug.Init()
	debug.AddWatcher(SWITCH_OS_TRACE, func(key string, value string) {
		logger.Info("配置最新值：key:【%v】, value：【%v】", key, value)
		trace.OsTraceSwitch = isc.ToBool(value)
	})
	debug.AddWatcher(SWITCH_OS_TRACE_HTTP, func(key string, value string) {
		logger.Info("配置最新值：key:【%v】, value：【%v】", key, value)
		trace.HttpTraceSwitch = isc.ToBool(value)
	})
	debug.AddWatcher(SWITCH_OS_TRACE_DATABASE, func(key string, value string) {
		logger.Info("配置最新值：key:【%v】, value：【%v】", key, value)
		trace.DatabaseTraceSwitch = isc.ToBool(value)
	})
	debug.AddWatcher(SWITCH_OS_TRACE_REDIS, func(key string, value string) {
		logger.Info("配置最新值：key:【%v】, value：【%v】", key, value)
		trace.RedisTraceSwitch = isc.ToBool(value)
	})
	debug.AddWatcher(SWITCH_OS_TRACE_ETCD, func(key string, value string) {
		logger.Info("配置最新值：key:【%v】, value：【%v】", key, value)
		trace.EtcdTraceSwitch = isc.ToBool(value)
	})
	debug.StartWatch()
}
