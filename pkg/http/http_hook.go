package http

import (
	"context"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/store"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/internal/trace"
	"github.com/isyscore/isc-tracer/pkg"
	"net/http"
	"strings"
	"unsafe"
)

type TracerHttpHook struct {
}

func (*TracerHttpHook) Before(ctx context.Context, req *http.Request) context.Context {
	if !trace.TracerIsEnable() {
		return ctx
	}
	_srcHead := store.GetHeader()
	srcHead := *_srcHead
	for headKey, srcHs := range srcHead {
		for _, srcH := range srcHs {
			req.Header.Add(headKey, srcH)
		}
	}

	tracer := pkg.ClientStartTraceWithRequest(req)
	ctx = context.WithValue(ctx, httpContextKey, tracer)
	return ctx
}

func (*TracerHttpHook) After(ctx context.Context, rsp *http.Response, rspCode int, rspData any, err error) {
	if !trace.TracerIsEnable() {
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
		if rspData != nil {
			bodyStr := string(rspData.([]byte))

			if strings.HasPrefix(bodyStr, "{") && strings.HasSuffix(bodyStr, "}") {
				bodys := map[string]any{}
				_ = isc.StrToObject(bodyStr, &bodys)
				code, existCode := bodys["code"]
				msg, _ := bodys["message"]
				if existCode && isc.ToInt(code) != 0 && isc.ToInt(code) != 200 {
					resultMap["errCode"] = code
					resultMap["errMsg"] = msg

					trace.EndTrace(tracer, _const.ERROR, isc.ToJsonString(resultMap), isc.ToInt(unsafe.Sizeof(rspData)))
					return
				}
			}
		}
	}
	trace.EndTrace(tracer, result, isc.ToJsonString(resultMap), 0)
	return
}
