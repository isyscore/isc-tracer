package http

import (
	"context"
	"github.com/isyscore/isc-gobase/isc"
	"github.com/isyscore/isc-gobase/store"
	_const2 "github.com/isyscore/isc-tracer/const"
	trace2 "github.com/isyscore/isc-tracer/trace"
	"net/http"
	"strings"
	"unsafe"
)

type TracerHttpHook struct {
}

func (*TracerHttpHook) Before(ctx context.Context, req *http.Request) (context.Context, http.Header) {
	if !trace2.TracerIsEnable() {
		return ctx, req.Header
	}

	newHeader := req.Header.Clone()
	if req != nil {
		for headKey, srcHs := range req.Header {
			for _, srcH := range srcHs {
				newHeader.Set(headKey, srcH)
			}
		}

		if url := req.URL; url != nil && IsExclude(url.Path) {
			return ctx, newHeader
		}
	}

	tracer := trace2.ClientStartTraceWithRequest(req)
	newHeader.Set(_const2.TRACE_HEAD_ID, isc.ToString(store.Get(_const2.TRACE_HEAD_ID)))
	newHeader.Set(_const2.TRACE_HEAD_RPC_ID, isc.ToString(store.Get(_const2.TRACE_HEAD_RPC_ID)))

	ctx = context.WithValue(ctx, httpContextKey, tracer)
	return ctx, newHeader
}

func (*TracerHttpHook) After(ctx context.Context, rsp *http.Response, rspCode int, rspData any, err error) {
	if !trace2.TracerIsEnable() {
		return
	}

	tracer, ok := ctx.Value(httpContextKey).(*trace2.Tracer)
	if !ok || tracer == nil {
		return
	}

	resultMap := map[string]any{}
	result := _const2.OK

	if rspCode >= 300 {
		result = _const2.ERROR
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
					resultMap["code"] = code
					resultMap["message"] = msg

					trace2.EndTrace(tracer, _const2.ERROR, isc.ToJsonString(resultMap), isc.ToInt(unsafe.Sizeof(rspData)))
					return
				}
			}
		}
	}
	trace2.EndTrace(tracer, result, isc.ToJsonString(resultMap), 0)
	return
}
