package test

import (
	"context"
	"github.com/isyscore/isc-gobase/extend/redis"
	tracer "github.com/isyscore/isc-tracer"
	"github.com/isyscore/isc-tracer/internal/trace"
	"testing"
)

func TestRedis(t *testing.T) {
	trace.OsTraceSwitch = true
	trace.RedisTraceSwitch = true

	tracer.Init()

	redisCli, err := redis.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	key := "test"
	cmd := redisCli.Set(ctx, key, "233", 0)
	if cmd.Err() != nil {
		t.Fatal(cmd.Err())
	}
	getCmd := redisCli.Get(ctx, key)
	if getCmd.Err() != nil {
		t.Fatal(getCmd.Err())
	}
	t.Log(getCmd.Val())
}
