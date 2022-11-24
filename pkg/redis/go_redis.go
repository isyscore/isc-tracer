package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/isyscore/isc-gobase/isc"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/internal/trace"
)

var redisContextKey = "gobase-redis-trace-key"

type TracerRedisHook struct {
}

func (*TracerRedisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if !trace.RedisTraceSwitch {
		return ctx, nil
	}

	tracer := trace.ClientStartTrace(_const.REDIS, "【redis】: "+cmd.Name())
	ctx = context.WithValue(ctx, redisContextKey, tracer)
	return ctx, nil
}

func (*TracerRedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if !trace.RedisTraceSwitch {
		return nil
	}

	tracer, ok := ctx.Value(redisContextKey).(*trace.Tracer)
	if !ok || tracer == nil {
		return nil
	}

	resultMap := map[string]any{}
	result := _const.OK
	// 记录error
	err := cmd.Err()
	if err != nil {
		resultMap["err"] = err.Error()
		result = _const.ERROR
	}

	args, err := json.Marshal(cmd.Args())
	if err != nil {
		resultMap["err"] = err.Error()
		result = _const.ERROR
	}

	resultMap["cmd"] = cmd.Name()
	resultMap["fullName"] = cmd.FullName()
	resultMap["parameters"] = string(args)

	trace.EndTrace(tracer, 0, result, isc.ToJsonString(resultMap))
	return nil
}

func (*TracerRedisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (*TracerRedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
