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

type GoBaseRedisHook struct {
}

func (*GoBaseRedisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if !trace.RedisTraceSwitch {
		return ctx, nil
	}

	tracer := trace.ServerStartTrace(_const.REDIS, cmd.Name())
	ctx = context.WithValue(ctx, redisContextKey, tracer)
	return ctx, nil
}

func (*GoBaseRedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
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

	// todo 返回值暂时未知，先不写
	trace.EndTrace(tracer, 0, result, isc.ToJsonString(resultMap))
	return nil
}

func (*GoBaseRedisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (*GoBaseRedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
