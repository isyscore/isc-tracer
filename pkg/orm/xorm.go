package orm

import (
	"context"
	"github.com/isyscore/isc-gobase/isc"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/internal/trace"
	"strings"
	"xorm.io/xorm/contexts"
)

const (
	traceContextXormKey = "tracer-xorm-trace-key"
)

type TracerXormHook struct {
}

func (*TracerXormHook) BeforeProcess(c *contexts.ContextHook, driverName string) (context.Context, error) {
	if !TracerDatabaseIsEnable() {
		return c.Ctx, nil
	}

	if c.SQL == "" {
		return c.Ctx, nil
	}

	ctx := c.Ctx
	sqlMetas := strings.SplitN(c.SQL, " ", 2)
	tracer := trace.ClientStartTrace(getSqlType(driverName), "【xorm】:"+sqlMetas[0])
	ctx = context.WithValue(ctx, traceContextXormKey, tracer)
	return ctx, nil
}

func (*TracerXormHook) AfterProcess(c *contexts.ContextHook, driverName string) error {
	if !TracerDatabaseIsEnable() {
		return nil
	}

	ctx := c.Ctx
	tracer, ok := ctx.Value(traceContextXormKey).(*trace.Tracer)
	if !ok || tracer == nil {
		return nil
	}

	resultMap := map[string]any{}
	result := _const.OK

	//b, _ := json.Marshal(c.Args)

	if c.Err != nil {
		resultMap["err"] = c.Err.Error()
		result = _const.ERROR
	}
	resultMap["database"] = driverName
	resultMap["sql"] = c.SQL
	//resultMap["parameters"] = string(b)

	trace.EndTrace(tracer, result, isc.ToJsonString(resultMap), 0)
	return nil
}
