package test

import (
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/debug"
	"github.com/isyscore/isc-gobase/server"
	isc_tracer "github.com/isyscore/isc-tracer"
	"testing"
)

// 使用环境变量：base.profiles.active=switch
func TestSwitch(t *testing.T) {
	server.Get("/switch/trace/open", traceOpen)
	server.Get("/switch/trace/close", traceClose)
	server.Get("/switch/trace/database/open", traceDatabaseOpen)
	server.Get("/switch/trace/database/close", traceDatabaseClose)
	//server.Get("/switch/trace/redis/open", traceRedisOpen)
	//server.Get("/switch/trace/redis/close", traceRedisClose)
	server.Get("/switch/trace/etcd/open", traceEtcdOpen)
	server.Get("/switch/trace/etcd/close", traceEtcdClose)

	server.Run()
}

func traceOpen(c *gin.Context) {
	debug.Update(isc_tracer.SWITCH_OS_TRACE, "true")
}

func traceClose(c *gin.Context) {
	debug.Update(isc_tracer.SWITCH_OS_TRACE, "false")
}

func traceDatabaseOpen(c *gin.Context) {
	debug.Update(isc_tracer.SWITCH_OS_TRACE_DATABASE, "true")
}

func traceDatabaseClose(c *gin.Context) {
	debug.Update(isc_tracer.SWITCH_OS_TRACE_DATABASE, "false")
}

//func traceRedisOpen(c *gin.Context) {
//	debug.Update(isc_tracer.SWITCH_OS_TRACE_REDIS, "true")
//}

//func traceRedisClose(c *gin.Context) {
//	debug.Update(isc_tracer.SWITCH_OS_TRACE_REDIS, "false")
//}

func traceEtcdOpen(c *gin.Context) {
	debug.Update(isc_tracer.SWITCH_OS_TRACE_ETCD, "true")
}

func traceEtcdClose(c *gin.Context) {
	debug.Update(isc_tracer.SWITCH_OS_TRACE_ETCD, "false")
}
