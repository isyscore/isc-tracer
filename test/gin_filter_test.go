package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	baseHttp "github.com/isyscore/isc-gobase/http"
	"github.com/isyscore/isc-gobase/server"
	"github.com/isyscore/isc-gobase/server/rsp"
	_ "github.com/isyscore/isc-tracer"
	"testing"
	"time"
)

func TestTraceFilter(t *testing.T) {
	server.Get("/test", test)
	server.Get("/test/err", testErr)

	go server.Run()

	_, _, data, _ := baseHttp.GetSimple("http://localhost:8082/api/test/err")
	if data == nil {
		fmt.Println("返回值：nil")
		return
	}
	fmt.Println("返回值：" + string(data.([]byte)))

	time.Sleep(10000000)
}

func test(c *gin.Context) {
	dict := make(map[string]any)
	dict["code"] = 0
	dict["message"] = "成功"
	rsp.Success(c, dict)
}

func testErr(c *gin.Context) {
	rsp.FailedOfStandard(c, 103222, "xxx业务的配置异常")
}
