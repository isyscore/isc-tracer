package test

import (
	"fmt"
	tracer "github.com/isyscore/isc-tracer"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	baseHttp "github.com/isyscore/isc-gobase/http"
	"github.com/isyscore/isc-gobase/server"
	"github.com/isyscore/isc-gobase/server/rsp"
	"github.com/isyscore/isc-tracer/internal/trace"
)

func TestTraceFilter(t *testing.T) {
	trace.OsTraceSwitch = true
	trace.HttpTraceSwitch = true
	tracer.Init()

	server.Get("/test", test)
	server.Get("/test/err", testErr)

	server.Run()
}

func test(c *gin.Context) {
	rsp.Success(c, "成功")
}

func testErr(c *gin.Context) {
	rsp.FailedOfStandard(c, 103222, "xxx业务的配置异常")
}

func TestGetSimple(t *testing.T) {
	trace.OsTraceSwitch = true
	trace.HttpTraceSwitch = true
	tracer.Init()

	_, _, data, _ := baseHttp.GetSimple("http://localhost:8082/api/test")
	if data == nil {
		fmt.Println("返回值：nil")
		return
	}
	fmt.Println("返回值：" + string(data.([]byte)))

	_, _, data, _ = baseHttp.GetSimple("http://localhost:8082/api/test/err")
	if data == nil {
		fmt.Println("返回值：nil")
		return
	}
	fmt.Println("返回值：" + string(data.([]byte)))

	time.Sleep(30000)
}
