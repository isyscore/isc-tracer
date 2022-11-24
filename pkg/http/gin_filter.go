package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/server/rsp"
	"github.com/isyscore/isc-gobase/store"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/internal/trace"
	"net/http"
	"strings"
	"unsafe"

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

type TracerHttpHook struct {
}

func (*TracerHttpHook) Before(ctx context.Context, req *http.Request) context.Context {
	if !trace.HttpTraceSwitch {
		return ctx
	}

	srcHead := store.GetHeader()
	for headKey, srcHs := range srcHead {
		for _, srcH := range srcHs {
			req.Header.Add(headKey, srcH)
		}
	}

	tracer := trace.ServerStartTrace(_const.HTTP, "【http】: <"+req.Method+">"+req.URL.Path)
	ctx = context.WithValue(ctx, httpContextKey, tracer)
	return ctx
}

func (*TracerHttpHook) After(ctx context.Context, rsp *http.Response, rspCode int, rspData any, err error) {
	if !trace.HttpTraceSwitch {
		return
	}

	tracer, ok := ctx.Value(httpContextKey).(*trace.Tracer)
	if !ok || tracer == nil {
		return
	}

	resultMap := map[string]any{}
	result := _const.OK

	if rspCode >= 300 {
		result = _const.ERROR
		if err != nil {
			resultMap["err"] = err.Error()
		}
	} else {
		bodyStr := string(rspData.([]byte))

		if strings.HasPrefix(bodyStr, "{") && strings.HasSuffix(bodyStr, "}") {
			bodys := map[string]any{}
			_ = isc.StrToObject(bodyStr, &bodys)
			code, existCode := bodys["code"]
			msg, _ := bodys["message"]
			if existCode && code != 0 && code != 200 {
				resultMap["errCode"] = code
				resultMap["errMsg"] = msg

				trace.EndTrace(tracer, isc.ToInt(unsafe.Sizeof(rspData)), _const.ERROR, isc.ToJsonString(resultMap))
				return
			}
		}
	}
	trace.EndTrace(tracer, 0, result, isc.ToJsonString(resultMap))
	return
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
		tracer := trace.ServerStartTrace(_const.HTTP, fmt.Sprintf("【http】: <%s>%s", c.Request.Method, c.Request.RequestURI))
		c.Writer.Header().Set(_const.TRACE_HEAD_ID, c.GetHeader(_const.TRACE_HEAD_ID))
		// 重写writer,用于获取response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		defer func() {
			code := _const.OK
			var msg string

			var response rsp.ResponseBase

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
			trace.EndTrace(tracer, blw.body.Len(), code, msg)
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
