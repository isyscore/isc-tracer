package redis

import (
	"context"
	"encoding/json"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/isc"
	_const2 "github.com/isyscore/isc-tracer/const"
	trace2 "github.com/isyscore/isc-tracer/trace"
	"github.com/redis/go-redis/v9"
)

var redisContextKey = "gobase-redis-trace-key"

type TracerRedisHook struct {
}

func (*TracerRedisHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}
func (*TracerRedisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	if !TracerRedisIsEnable() {
		return next
	}
	return func(ctx context.Context, cmd redis.Cmder) error {
		tracer := trace2.ClientStartTrace(_const2.REDIS, "【go-redis】: "+cmd.Name())
		ctx = context.WithValue(ctx, redisContextKey, tracer)

		next(ctx, cmd)

		tracer, ok := ctx.Value(redisContextKey).(*trace2.Tracer)
		if !ok || tracer == nil {
			return nil
		}

		resultMap := map[string]any{}
		result := _const2.OK
		// 记录error
		err := cmd.Err()
		if err != nil {
			resultMap["err"] = err.Error()
			result = _const2.ERROR
		}

		args, err := json.Marshal(cmd.Args())
		if err != nil {
			resultMap["err"] = err.Error()
			result = _const2.ERROR
		}

		resultMap["cmd"] = cmd.Name()
		resultMap["fullName"] = cmd.FullName()
		resultMap["parameters"] = string(args)

		tracer.PutAttr("a-cmd", cmd.FullName())
		trace2.EndTrace(tracer, result, isc.ToJsonString(resultMap), 0)
		return nil
	}
}
func (*TracerRedisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}

func (*TracerRedisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if !TracerRedisIsEnable() {
		return ctx, nil
	}

	tracer := trace2.ClientStartTrace(_const2.REDIS, "【go-redis】: "+cmd.Name())
	ctx = context.WithValue(ctx, redisContextKey, tracer)
	return ctx, nil
}

func (*TracerRedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if !TracerRedisIsEnable() {
		return nil
	}

	tracer, ok := ctx.Value(redisContextKey).(*trace2.Tracer)
	if !ok || tracer == nil {
		return nil
	}

	resultMap := map[string]any{}
	result := _const2.OK
	// 记录error
	err := cmd.Err()
	if err != nil {
		resultMap["err"] = err.Error()
		result = _const2.ERROR
	}

	args, err := json.Marshal(cmd.Args())
	if err != nil {
		resultMap["err"] = err.Error()
		result = _const2.ERROR
	}

	resultMap["cmd"] = cmd.Name()
	resultMap["fullName"] = cmd.FullName()
	resultMap["parameters"] = string(args)

	tracer.PutAttr("a-cmd", cmd.FullName())
	trace2.EndTrace(tracer, result, isc.ToJsonString(resultMap), 0)
	return nil
}

func (*TracerRedisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (*TracerRedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

func TracerRedisIsEnable() bool {
	return config.GetValueBoolDefault("tracer.redis.enable", true) && trace2.SwitchTraceRedis
}
