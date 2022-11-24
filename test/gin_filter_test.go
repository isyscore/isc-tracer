package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	baseHttp "github.com/isyscore/isc-gobase/http"
	"github.com/isyscore/isc-gobase/server"
	"github.com/isyscore/isc-gobase/server/rsp"
	"github.com/isyscore/isc-tracer/internal/trace"
	_ "github.com/isyscore/isc-tracer/pkg"
	"testing"
)

func TestTraceFilter(t *testing.T) {
	trace.OsTraceSwitch = true
	trace.HttpTraceSwitch = true

	server.Get("/test", test)

	server.Run()
}

func test(c *gin.Context) {
	rsp.Success(c, "成功")
}

func TestGetSimple(t *testing.T) {
	_, _, data, _ := baseHttp.GetSimple("http://localhost:8081/api/test")

	if data == nil {
		fmt.Println("返回值：nil")
	}
	fmt.Println("返回值：" + string(data.([]byte)))
}
