package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/server"
	"github.com/isyscore/isc-gobase/server/rsp"
	"github.com/isyscore/isc-tracer/config"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/pkg"
	"runtime/debug"
)

const (
	API_PREFIX = "/api"
)

var (
	excludes    = []string{"/system/status"}
	copyAttrMap = map[string]string{
		_const.TRACE_HEAD_REMOTE_APPNAME: _const.TRACE_HEAD_REMOTE_APPNAME,
		_const.TRACE_HEAD_REMOTE_IP:      _const.TRACE_HEAD_REMOTE_IP,
		_const.TRACE_HEAD_USER_ID:        _const.A_USER_ID,
		_const.TRACE_HEAD_USER_NAME:      _const.A_USER_NAME,
	}
)

func init() {
	server.AddGinHandlers(TraceFilter())
}

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
		traceConfig := config.GetConfig()
		if !traceConfig.Enable || isExclude(c) {
			c.Next()
			return
		}
		// 开始追踪
		tracer := pkg.ServerStartTrace(_const.HTTP, c.Request.RequestURI)
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
			pkg.ServerEndTrace(tracer, blw.body.Len(), code, msg)
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

//func record(c *gin.Context, msg string, start time.Time) {
//	uri := c.Request.RequestURI
//	if strings.HasPrefix(uri, API_PREFIX) {
//		traceName := fmt.Sprintf("<%s>%s", c.Request.Method, uri)
//		//异常记录
//		if msg != "" {
//			context := ""
//			rt := time.Since(start).Seconds()
//			// 根据code判断except还是warn
//			logExcept(traceName, context, msg, rt)
//		} else {
//			// log metric
//		}
//	}
//}
//
//func logExcept(name string, context string, msg string, rt float64) {
//	// 通过grpc发送给pivot服务端
//
//}
