package http

import (
	"github.com/gin-gonic/gin"
	c "github.com/isyscore/isc-tracer/config"
	"os"
	"os/signal"
	"testing"
)

func TestTraceFilter(t *testing.T) {
	serviceConfig := &c.Config{
		ServiceName: "test",
		Enable:      true,
	}
	c.ServerConfig = serviceConfig

	engine := gin.Default()
	engine.Use(TraceFilter)
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
