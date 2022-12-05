package orm

import (
	"context"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/isc"
	_const2 "github.com/isyscore/isc-tracer/const"
	trace2 "github.com/isyscore/isc-tracer/trace"
	"strings"
)

const (
	traceContextGormKey = "gobase-gorm-trace-key"
)

type TracerGormHook struct {
}

func (*TracerGormHook) Before(ctx context.Context, driverName string, parameters map[string]any) (context.Context, error) {
	if !TracerDatabaseIsEnable() {
		return ctx, nil
	}

	query, exist := parameters["query"]
	if !exist {
		return ctx, nil
	}

	cmds := strings.SplitN(query.(string), " ", 2)
	tracer := trace2.ClientStartTrace(getSqlType(driverName), "【gorm】:"+cmds[0])
	return context.WithValue(ctx, traceContextGormKey, tracer), nil
}

func (*TracerGormHook) After(ctx context.Context, driverName string, parameters map[string]any) (context.Context, error) {
	if !TracerDatabaseIsEnable() {
		return ctx, nil
	}

	tracer, ok := ctx.Value(traceContextGormKey).(*trace2.Tracer)
	if !ok || tracer == nil {
		return ctx, nil
	}

	query, _ := parameters["query"]
	args, _ := parameters["args"]

	resultMap := map[string]any{}
	resultMap["database"] = driverName
	resultMap["sql"] = query
	//resultMap["parameters"] = args

	trace2.EndTrace(tracer, _const2.OK, isc.ToJsonString(resultMap), 0)
	return ctx, nil
}

func (*TracerGormHook) Err(ctx context.Context, driverName string, err error, parameters map[string]any) error {
	if !TracerDatabaseIsEnable() {
		return nil
	}

	tracer, ok := ctx.Value(traceContextGormKey).(*trace2.Tracer)
	if !ok || tracer == nil {
		return nil
	}

	query, _ := parameters["query"]
	//args, _ := parameters["args"]

	resultMap := map[string]any{}
	resultMap["database"] = driverName
	resultMap["sql"] = query
	//resultMap["parameters"] = args
	resultMap["err"] = err.Error()

	trace2.EndTrace(tracer, _const2.ERROR, isc.ToJsonString(resultMap), 0)
	return nil
}

func getSqlType(driverName string) _const2.TraceTypeEnum {
	driverName = strings.ToLower(driverName)
	switch driverName {
	case "mysql":
		return _const2.MYSQL
	case "postgresql":
		return _const2.POSTGRESQL
	case "sqlite":
		return _const2.SQLITE
	}
	return _const2.UNKNOWN
}

func TracerDatabaseIsEnable() bool {
	return config.GetValueBoolDefault("tracer.database.enable", true) && trace2.SwitchTraceDatabase
}
