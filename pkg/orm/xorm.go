package orm

import (
	"context"
	"encoding/json"
	"github.com/isyscore/isc-gobase/isc"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/internal/trace"
	"strings"
	"xorm.io/xorm/contexts"
)

const (
	traceContextXormKey = "tracer-xorm-trace-key"
)

type GobaseXormHook struct {
}

func (*GobaseXormHook) BeforeProcess(c *contexts.ContextHook) (context.Context, error) {
	if !trace.DatabaseTraceSwitch {
		return c.Ctx, nil
	}

	if c.SQL == "" {
		return c.Ctx, nil
	}

	ctx := c.Ctx
	sqlMetas := strings.SplitN(c.SQL, " ", 2)
	tracer := trace.ServerStartTrace(_const.MYSQL, "xorm: "+sqlMetas[0])
	ctx = context.WithValue(ctx, traceContextXormKey, tracer)
	return ctx, nil
}

func (*GobaseXormHook) AfterProcess(c *contexts.ContextHook) error {
	if !trace.DatabaseTraceSwitch {
		return nil
	}

	ctx := c.Ctx
	tracer, ok := ctx.Value(traceContextXormKey).(*trace.Tracer)
	if !ok || tracer == nil {
		return nil
	}

	resultMap := map[string]any{}
	result := _const.OK

	b, _ := json.Marshal(c.Args)

	if c.Err != nil {
		resultMap["err"] = c.Err.Error()
		result = _const.ERROR
	}
	resultMap["sql"] = c.SQL
	resultMap["parameters"] = string(b)

	trace.EndTrace(tracer, 0, result, isc.ToJsonString(resultMap))
	return nil
}
