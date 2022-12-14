package etcd

import (
	"context"
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/isc"
	_const2 "github.com/isyscore/isc-tracer/const"
	trace2 "github.com/isyscore/isc-tracer/trace"
	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
	etcdClientV3 "go.etcd.io/etcd/client/v3"
	"reflect"
)

var etcdContextKey = "gobase-etcd-context-key"

type TracerEtcdHook struct {
}

func (pHook *TracerEtcdHook) Before(ctx context.Context, op etcdClientV3.Op) context.Context {
	if !TracerEtcdIsEnable() {
		return ctx
	}

	tracer := trace2.ClientStartTrace(_const2.ETCD, "【etcd.io】: "+getCmd(op))
	ctx = context.WithValue(ctx, etcdContextKey, tracer)
	return ctx
}

func (pHook *TracerEtcdHook) After(ctx context.Context, op etcdClientV3.Op, pRsp any, err error) {
	if !TracerEtcdIsEnable() {
		return
	}

	tracer, ok := ctx.Value(etcdContextKey).(*trace2.Tracer)
	if !ok || tracer == nil {
		return
	}

	resultMap := map[string]any{}
	result := _const2.OK
	// 记录error
	if err != nil {
		resultMap["err"] = err.Error()
		result = _const2.ERROR
	}

	resultMap["req"] = isc.ToJsonString(toRequestOp(op))

	trace2.EndTrace(tracer, result, isc.ToJsonString(resultMap), 0)
	return
}

func toRequestOp(op etcdClientV3.Op) *pb.RequestOp {
	if op.IsGet() {
		return &pb.RequestOp{Request: &pb.RequestOp_RequestRange{RequestRange: toRangeRequest(op)}}
	} else if op.IsPut() {
		r := &pb.PutRequest{
			Key:    op.KeyBytes(),
			Value:  op.ValueBytes(),
			Lease:  int64(isc.GetPrivateFieldValue(reflect.ValueOf(&op), "leaseID").(etcdClientV3.LeaseID)),
			PrevKv: isc.GetPrivateFieldValue(reflect.ValueOf(&op), "prevKV").(bool),
		}
		return &pb.RequestOp{Request: &pb.RequestOp_RequestPut{RequestPut: r}}
	} else if op.IsDelete() {
		r := &pb.DeleteRangeRequest{
			Key:      op.KeyBytes(),
			RangeEnd: op.RangeBytes(),
			PrevKv:   isc.GetPrivateFieldValue(reflect.ValueOf(&op), "prevKV").(bool),
		}
		return &pb.RequestOp{Request: &pb.RequestOp_RequestDeleteRange{RequestDeleteRange: r}}
	}
	return nil
}

func toRangeRequest(op etcdClientV3.Op) *pb.RangeRequest {
	if !op.IsGet() {
		return nil
	}
	r := &pb.RangeRequest{
		Key:               op.KeyBytes(),
		RangeEnd:          op.RangeBytes(),
		Revision:          op.Rev(),
		Serializable:      op.IsSerializable(),
		KeysOnly:          op.IsKeysOnly(),
		CountOnly:         op.IsCountOnly(),
		MinModRevision:    op.MinModRev(),
		MaxModRevision:    op.MaxModRev(),
		MinCreateRevision: op.MinCreateRev(),
		MaxCreateRevision: op.MaxCreateRev(),
	}
	return r
}

func getCmd(op etcdClientV3.Op) string {
	if op.IsGet() {
		return "get"
	} else if op.IsPut() {
		return "put"
	} else if op.IsDelete() {
		return "delete"
	}
	return ""
}

func TracerEtcdIsEnable() bool {
	return config.GetValueBoolDefault("tracer.etcd.enable", true) && trace2.SwitchTraceEtcd
}
