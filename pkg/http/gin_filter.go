package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/server/rsp"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/internal/trace"
	"runtime/debug"
)

const (
	API_PREFIX = "/api"
)

var (
	excludes = []string{"/system/status"}
)

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
		if isExclude(c) {
			c.Next()
			return
		}

		if !trace.HttpTraceSwitch {
			c.Next()
			return
		}
		// 开始追踪
		tracer := trace.ServerStartTrace(_const.HTTP, fmt.Sprintf("<%s>%s", c.Request.Method, c.Request.RequestURI))
		c.Writer.Header().Set(_const.TRACE_HEAD_ID, c.GetHeader(_const.TRACE_HEAD_ID))
		// 重写writer,用于获取response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		defer func() {
			code := _const.OK
			var msg string

			var response rsp.DataResponse[any]

			if err := recover(); err != nil {
				code = _const.ERROR
				msg = string(debug.Stack())
			} else if httpStatus := c.Writer.Status(); httpStatus >= 300 {
				code = _const.ERROR
				msg = fmt.Sprintf("httpStatus: %d", httpStatus)
			} else if err := json.Unmarshal([]byte(blw.body.String()), &response); err != nil {
				code = _const.WARNING
				msg = err.Error()
			} else {
				if response.Code != 0 {
					code = _const.ERROR
				}
				msg = response.Message
			}
			// 结束追踪
			trace.ServerEndTrace(tracer, blw.body.Len(), code, msg)
		}()
		c.Next()
	}
}

func isExclude(context *gin.Context) bool {
	uri := context.Request.RequestURI

	for _, exclude := range excludes {
		if exclude == uri {
			return true
		}
	}
	return false
}
