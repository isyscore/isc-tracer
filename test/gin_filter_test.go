package test

import (
	"github.com/gin-gonic/gin"
	baseHttp "github.com/isyscore/isc-gobase/http"
	"github.com/isyscore/isc-gobase/listener"
	"github.com/isyscore/isc-gobase/server"
	"github.com/isyscore/isc-gobase/server/rsp"
	_ "github.com/isyscore/isc-tracer"
	"testing"
)

// 使用环境变量：base.profiles.active=http
func TestTraceFilter(t *testing.T) {
	server.Get("/test", test)
	server.Get("/test/err", testErr)
	server.Get("/test/send", send)
	server.Get("/test/receive", receive)

	listener.AddListener(listener.EventOfServerRunFinish, func(event listener.BaseEvent) {
		baseHttp.GetSimple("http://localhost:8082/api/test/send")
	})
	server.Run()

	//访问 http://localhost:8082/api/test/send
}

func test(c *gin.Context) {
	t := &testing.T{}
	TestEtcd(t)

	baseHttp.GetSimple("http://localhost:8082/api/test/err")
	rsp.SuccessOfStandard(c, nil)
}

func testErr(c *gin.Context) {
	t := &testing.T{}
	TestXorm(t)
	rsp.FailedOfStandard(c, 103222, "xxx业务的配置异常")
}

func send(c *gin.Context) {
	t := &testing.T{}
	TestRedis(t)

	baseHttp.GetSimple("http://localhost:8082/api/test/receive")
	rsp.SuccessOfStandard(c, nil)
}

func receive(c *gin.Context) {
	t := &testing.T{}
	TestGorm(t)

	baseHttp.GetSimple("http://localhost:8082/api/test")
	rsp.SuccessOfStandard(c, "ok")
}
