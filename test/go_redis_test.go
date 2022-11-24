package test

import (
	"context"
	"github.com/isyscore/isc-gobase/extend/redis"
	_ "github.com/isyscore/isc-tracer/pkg"
	"testing"
)

func TestRedis(t *testing.T) {
	redis, err := redis.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	key := "test"
	redis.Set(ctx, key, 233, -1)
	v := redis.Get(ctx, key)
	t.Log(v)
}
