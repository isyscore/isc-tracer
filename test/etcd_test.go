package test

import (
	"context"
	"fmt"
	"github.com/isyscore/isc-gobase/extend/etcd"
	"github.com/isyscore/isc-gobase/time"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
)

// 使用环境变量：base.profiles.active=etcd
func TestEtcd(t *testing.T) {
	etcdClient, _ := etcd.NewEtcdClient()

	ctx := context.Background()
	etcdClient.Put(ctx, "test", time.TimeToStringYmdHms(time.Now()))

	rsp, _ := etcdClient.Get(ctx, "test")
	fmt.Println(rsp)

	etcdClient.Get(ctx, "test", func(pOp *clientv3.Op) {
		//fmt.Println("信息")
		//fmt.Println(isc.ToJsonString(&pOp))
	})
}
