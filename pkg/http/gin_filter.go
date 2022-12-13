package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/server/rsp"
	_const2 "github.com/isyscore/isc-tracer/const"
	trace2 "github.com/isyscore/isc-tracer/trace"
	"strings"

	"runtime/debug"
)

var (
	excludes = []string{"/system/status"}
)

var httpContextKey = "gobase-http-context-key"

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func TraceFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if IsExclude(c.Request.RequestURI) {
			c.Next()
			return
		}

		if !trace2.TracerIsEnable() {
			c.Next()
			return
		}

		// 开始追踪
		tracer := trace2.ServerStartTrace(_const2.HTTP, fmt.Sprintf("<%s>%s", c.Request.Method, c.Request.RequestURI))
		c.Writer.Header().Set(_const2.TRACE_HEAD_ID, tracer.TraceId)
		// 重写writer,用于获取response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		defer func() {
			code := _const2.OK
			var msg string

			var response rsp.ResponseBase

			if err := recover(); err != nil {
				code = _const2.ERROR
				msg = string(debug.Stack())
			} else if httpStatus := c.Writer.Status(); httpStatus >= 300 {
				code = _const2.ERROR
				msg = fmt.Sprintf("httpStatus: %d", httpStatus)
			} else if err := json.Unmarshal([]byte(blw.body.String()), &response); err != nil {
				code = _const2.WARNING
				msg = err.Error()
			} else {
				// 取code
				if response.Code != 0 && response.Code != 200 {
					code = _const2.ERROR
				}
				msg = isc.ToJsonString(response)
			}
			// 结束追踪
			trace2.EndTrace(tracer, code, msg, blw.body.Len())
		}()
		c.Next()
	}
}

func IsExclude(uri string) bool {
	for _, exclude := range excludes {
		if strings.Contains(uri, exclude) {
			return true
		}
	}

	excludes := config.GetValueArray("tracer.http.excludes-url")
	for _, excludeUri := range excludes {
		if strings.Contains(uri, excludeUri.(string)) {
			return true
		}
	}
	return false
}
