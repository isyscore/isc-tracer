package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	baseRedis "github.com/isyscore/isc-gobase/extend/redis"
	"github.com/isyscore/isc-gobase/isc"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/internal/trace"
	"github.com/isyscore/isc-tracer/pkg"
)

var redisContextKey = "gobase-redis-trace-key"

type GoBaseRedisHook struct {
}

func init() {
	baseRedis.AddRedisHook(&GoBaseRedisHook{})
}

func (*GoBaseRedisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	// 这里是关键，通过 envoy 传过来的 header 解析出父 span，如果没有，则会创建新的根 span
	//zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	//spanCtx, err := zipkinPropagator.Extract(opentracing.HTTPHeadersCarrier(server.GetHeader()))
	//if err != nil {
	//	logger.Warn("span 解析失败, 错误原因: %v", err)
	//	return ctx, err
	//}

	tracer := pkg.ServerStartTrace(_const.REDIS, cmd.Name())
	ctx = context.WithValue(ctx, redisContextKey, tracer)
	return ctx, nil
}

func (*GoBaseRedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
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
	pkg.ServerEndTrace(tracer, 0, result, isc.ToJsonString(resultMap))
	return nil
}

func (*GoBaseRedisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (*GoBaseRedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
