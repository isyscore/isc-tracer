package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/server"
	"github.com/isyscore/isc-gobase/server/rsp"
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
	server.Get("/test", test)

	server.Run()
}

func test(c *gin.Context) {
	dict := make(map[string]any)
	dict["code"] = 0
	dict["message"] = "成功"
	c.JSON(200, dict)
}
