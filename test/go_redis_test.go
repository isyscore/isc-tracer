package test

import (
	"context"
	"github.com/isyscore/isc-gobase/extend/redis"
	"github.com/isyscore/isc-tracer/internal/trace"
	"testing"
	"time"
)

// 使用环境变量：base.profiles.active=redis
func TestRedis(t *testing.T) {
	trace.SwitchTrace = true
	trace.SwitchTraceRedis = true

	redisCli, err := redis.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 1; i++ {
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
	time.Sleep(time.Second * 2)

}
