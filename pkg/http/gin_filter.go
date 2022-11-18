package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/isyscore/isc-tracer/config"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/internal/trace"
	"github.com/isyscore/isc-tracer/util"
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

func filter(c *gin.Context) {
	//todo 接口记录
	traceConfig := config.GetConfig()
	if !traceConfig.Enable || isExclude(c) {
		c.Next()
		return
	}
	//start := time.Now()

	tracerId := c.GetHeader(_const.TRACE_HEAD_ID)
	frontIP := ""
	if tracerId == "" {
		tracerId = util.GenerateTraceId()
		frontIP = getFrontIP(c.Request)
	}
	rpcId := c.GetHeader(_const.TRACE_HEAD_RPC_ID)
	// 开始追踪
	tracer := trace.StartTrace(tracerId, rpcId, _const.HTTP, c.Request.RequestURI)
	if frontIP != "" {
		tracer.RemoteIp = frontIP
	}
	// 往当前上下文添加远程端属性
	putAttr(tracer, c.Request)

	c.Writer.Header().Set(_const.TRACE_HEAD_ID, tracerId)

	defer func() {
		msg := ""
		if err := recover(); err != nil {
			msg = string(debug.Stack())
		}
		//todo 解析返回值,拿到code和msg
		code := _const.ParseHttpStatus(c.Writer.Status())
		//record(c, msg, start)
		tracer.Size = int(c.Request.ContentLength) + c.Writer.Size()
		tracer.EndTrace(code, msg)
	}()
	c.Next()
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

func putAttr(tracer *trace.Tracer, request *http.Request) {
	if tracer.AttrMap == nil {
		tracer.AttrMap = make(map[string]string)
	}
	for key, copyKey := range copyAttrMap {
		if v := request.Header.Get(key); v != "" {
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

func getFrontIP(req *http.Request) string {
	ip := req.Header.Get("X-Forwarded-For")
	if ip != "" && strings.EqualFold(ip, "unKnown") {
		//多次反向代理后会有多个ip值，第一个ip才是真实ip
		if i := strings.Index(ip, ","); i != -1 {
			return ip[:i]
		}
		return ip
	}
	ip = req.Header.Get("X-Real-IP")
	if ip != "" && strings.EqualFold(ip, "unKnown") {
		return ip
	}
	return req.RemoteAddr
}
