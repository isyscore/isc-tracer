package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/server/rsp"
	c "github.com/isyscore/isc-tracer/config"
	"os"
	"os/signal"
	"testing"
)

func TestJson(t *testing.T) {
	v := map[string]any{
		"code":    0,
		"message": "成功",
		"data":    123,
	}
	b, _ := json.Marshal(v)
	var response rsp.ResponseBase
	json.Unmarshal(b, &response)
	t.Log(response)
}
func TestTraceFilter(t *testing.T) {
	serviceConfig := &c.Config{
		ServiceName: "test",
		Enable:      true,
	}
	c.ServerConfig = serviceConfig

	engine := gin.Default()
	engine.Use(TraceFilter())
	// http://localhost:8080/api/test
	group := engine.Group("/api")
	group.GET("/test", test)

	err := engine.Run(":8080")
	if err != nil {
		return
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func test(c *gin.Context) {
	dict := make(map[string]any)
	dict["code"] = 0
	dict["message"] = "成功"
	c.JSON(200, dict)
}
