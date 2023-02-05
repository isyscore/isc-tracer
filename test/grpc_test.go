package test

import (
	"context"
	"fmt"
	"github.com/isyscore/isc-tracer/pivot"
	"google.golang.org/grpc"
	"testing"
)

func TestGrpcCollectTracer(t *testing.T) {
	//建立链接
	// 连接服务器
	url := "localhost:9091"
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("连接服务端失败: %s", err)
		return
	}
	defer conn.Close()

	pivotService := pivot.NewPivotServiceClient(conn)

	ctx := context.Background()

	tracer := &pivot.TracerRequest{
		TraceId: "tracer_id",
		RpcId: "rpc_id",
		TraceType: 1,
		TraceName: "tracer_name",
		Endpoint: 1,
		Status: 1,
		RemoteStatus: 1,
		RemoteIp: "remote_ip",
		Message: "message",
		Size: 12,
		StartTime: 38172391872,
		EndTime: 38172391872,
		Sampled: true,
		BizData: map[string][]byte{
			"k1": []byte{32, 54 , 32,1, 32},
		},
		Ended: true,
		AttrMap: map[string]string{
			"k1": "v1",
		},
		ContextMap: map[string][]byte{
			"k1": []byte{32, 54 , 32,1, 32},
		},
		ThreadMode: 1,
	}
	pivotService.CollectTracer(ctx, tracer)
}
