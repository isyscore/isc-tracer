package test

import (
	"github.com/isyscore/isc-tracer/conf"
	"github.com/isyscore/isc-tracer/trace"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type fields struct {
	Tracer      *trace.Tracer
	clientRpcId string
}

type args struct {
	url          string
	header       http.Header
	parameterMap map[string]string
	body         any
	httpRequest  *http.Request
}

type testCase struct {
	name    string
	fields  fields
	args    args
	want    []byte
	wantErr bool
}

var caseDelete = []testCase{
	// TODO: Add test cases.
	{
		name: "删除接口测试",
		fields: fields{
			&trace.Tracer{
				TraceId:     trace.LocalIdCreate.GenerateTraceId(),
				sampled:     true,
				ServiceName: conf.Conf.ServiceName,
				startTime:   time.Now().UnixMilli(),
				RpcId:       "0",
				TraceType:   trace.HTTP,
				RemoteIp:    trace.GetLocalIp(),
				TraceName:   "<default>_server",
			},
			"",
		},
		args: args{
			url:          "http://10.30.30.78:38080/api/apix/execute",
			header:       map[string][]string{"id": {"kucs"}},
			parameterMap: map[string]string{"name": "库陈胜"},
		},
		wantErr: true,
	},
}

var caseGet = []testCase{}
var caseHead = []testCase{}
var casePatch = []testCase{}
var casePost = []testCase{}
var casePut = []testCase{}
var caseCall = []testCase{}

func TestServerTracer_Delete(t *testing.T) {
	for _, tt := range caseDelete {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.Delete(tt.args.url, tt.args.header, tt.args.parameterMap)
			t.Logf("clientId=%s", server.ClientRpcId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_DeleteOfStandard(t *testing.T) {
	for _, tt := range caseDelete {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.DeleteOfStandard(tt.args.url, tt.args.header, tt.args.parameterMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_DeleteSimple(t *testing.T) {
	for _, tt := range caseDelete {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.DeleteSimple(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteSimple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteSimple() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_DeleteSimpleOfStandard(t *testing.T) {
	for _, tt := range caseDelete {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.DeleteSimpleOfStandard(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteSimpleOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteSimpleOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_Get(t *testing.T) {
	for _, tt := range caseGet {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.Get(tt.args.url, tt.args.header, tt.args.parameterMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_GetOfStandard(t *testing.T) {
	for _, tt := range caseGet {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.GetOfStandard(tt.args.url, tt.args.header, tt.args.parameterMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_GetSimple(t *testing.T) {
	for _, tt := range caseGet {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.GetSimple(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSimple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSimple() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_GetSimpleOfStandard(t *testing.T) {
	for _, tt := range caseGet {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.GetSimpleOfStandard(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSimpleOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSimpleOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_Head(t *testing.T) {
	for _, tt := range caseHead {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			if err := server.Head(tt.args.url, tt.args.header, tt.args.parameterMap); (err != nil) != tt.wantErr {
				t.Errorf("Head() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerTracer_HeadSimple(t *testing.T) {
	for _, tt := range caseHead {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			if err := server.HeadSimple(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("HeadSimple() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerTracer_Patch(t *testing.T) {
	for _, tt := range casePatch {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.Patch(tt.args.url, tt.args.header, tt.args.parameterMap, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Patch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Patch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PatchOfStandard(t *testing.T) {
	for _, tt := range casePatch {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.PatchOfStandard(tt.args.url, tt.args.header, tt.args.parameterMap, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PatchSimple(t *testing.T) {
	for _, tt := range casePatch {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.PatchSimple(tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchSimple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchSimple() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PatchSimpleOfStandard(t *testing.T) {
	for _, tt := range casePatch {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.PatchSimpleOfStandard(tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PatchSimpleOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchSimpleOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_Post(t *testing.T) {
	for _, tt := range casePost {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.Post(tt.args.url, tt.args.header, tt.args.parameterMap, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Post() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PostOfStandard(t *testing.T) {
	for _, tt := range casePost {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.PostOfStandard(tt.args.url, tt.args.header, tt.args.parameterMap, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PostSimple(t *testing.T) {
	for _, tt := range casePost {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.PostSimple(tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostSimple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostSimple() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PostSimpleOfStandard(t *testing.T) {
	for _, tt := range casePost {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.PostSimpleOfStandard(tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostSimpleOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostSimpleOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_Put(t *testing.T) {
	for _, tt := range casePut {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.Put(tt.args.url, tt.args.header, tt.args.parameterMap, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Put() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PutOfStandard(t *testing.T) {
	for _, tt := range casePut {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.PutOfStandard(tt.args.url, tt.args.header, tt.args.parameterMap, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PutSimple(t *testing.T) {
	for _, tt := range casePut {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.PutSimple(tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutSimple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutSimple() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_PutSimpleOfStandard(t *testing.T) {
	for _, tt := range casePut {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.PutSimpleOfStandard(tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutSimpleOfStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutSimpleOfStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_call(t *testing.T) {
	for _, tt := range caseCall {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.call(tt.args.httpRequest, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("call() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("call() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerTracer_callIgnoreReturn(t *testing.T) {
	for _, tt := range caseCall {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			if err := server.callIgnoreReturn(tt.args.httpRequest, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("callIgnoreReturn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerTracer_callToStandard(t *testing.T) {
	for _, tt := range caseCall {
		t.Run(tt.name, func(t *testing.T) {
			server := &trace.ServerTracer{
				Tracer:      tt.fields.Tracer,
				ClientRpcId: tt.fields.clientRpcId,
			}
			_, _, got, err := server.callToStandard(tt.args.httpRequest, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("callToStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("callToStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}
