package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-gobase/server"
	"github.com/isyscore/isc-tracer/config"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/internal/trace"
	"github.com/isyscore/isc-tracer/pkg"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
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

func TraceFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		//todo 接口记录
		traceConfig := config.GetConfig()
		if !traceConfig.Enable || isExclude(c) {
			c.Next()
			return
		}

		// 开始追踪
		tracer := pkg.ServerStartTrace(_const.HTTP, c.Request.RequestURI)
		c.Writer.Header().Set(_const.TRACE_HEAD_ID, c.GetHeader(_const.TRACE_HEAD_ID))

		defer func() {
			msg := ""
			if err := recover(); err != nil {
				msg = string(debug.Stack())
			}
			pkg.ServerEndTrace(tracer, int(c.Request.ContentLength)+c.Writer.Size(), _const.ParseHttpStatus(c.Writer.Status()), msg)
		}()
		c.Next()
	}
}

func record(c *gin.Context, msg string, start time.Time) {
	uri := c.Request.RequestURI
	if strings.HasPrefix(uri, API_PREFIX) {
		traceName := fmt.Sprintf("<%s>%s", c.Request.Method, uri)
		//异常记录
		if msg != "" {
			context := ""
			rt := time.Since(start).Seconds()
			//todo 根据code判断except还是warn
			logExcept(traceName, context, msg, rt)
		} else {
			//todo log metric
		}
	}
}

func logExcept(name string, context string, msg string, rt float64) {
	//todo 通过grpc发送给pivot服务端

}

func putAttr(tracer *trace.Tracer, head http.Header) {
	if tracer.AttrMap == nil {
		tracer.AttrMap = make(map[string]string)
	}
	for key, copyKey := range copyAttrMap {
		if v := head.Get(key); v != "" {
			tracer.AttrMap[copyKey] = v
		}
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

//func getFrontIP(req *http.Request) string {
//	ip := req.Header.Get("X-Forwarded-For")
//	if ip != "" && strings.EqualFold(ip, "unKnown") {
//		//多次反向代理后会有多个ip值，第一个ip才是真实ip
//		if i := strings.Index(ip, ","); i != -1 {
//			return ip[:i]
//		}
//		return ip
//	}
//	ip = req.Header.Get("X-Real-IP")
//	if ip != "" && strings.EqualFold(ip, "unKnown") {
//		return ip
//	}
//	return req.RemoteAddr
//}
