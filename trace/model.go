package trace

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/isyscore/isc-tracer/conf"
	"github.com/isyscore/isc-tracer/push"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TraceTypeEnum 标明链路跟踪类型
type TraceTypeEnum int

// EndpointEnum 表明当前节点类型
type EndpointEnum int

// TraceStatusEnum 标明当前trace的结果
type TraceStatusEnum int

const (
	ROOT TraceTypeEnum = iota
	HTTP
	DUBBO
	MYSQL
	ROCKETMQ
	REDIS
	KAFKA
	IDS
	MQTT
	ORACLE
	ELASTIC
	ZOOKEEPER
	HBASE
	HADOOP
	FLINK
	SPARK
	KUDU
	HIVE
	STORM
	CONFIG
)
const (
	CLIENT EndpointEnum = iota
	SERVER
)
const (
	OK TraceStatusEnum = iota
	ERROR
	WARNING
	TIMEOUT
)

const (
	T_HEADER_TRACEID = "T-Head-TraceId"
	T_HEADER_RPCID   = "T-Head-Rpcid"
)

// Parameter 请求入参
type Parameter struct {
	// Name 参数名称
	Name string `json:"name"`
	// In 参数位置
	In string `json:"in"`
	// Value 参数值，如果是文件，则只显示文件名
	Value []string `json:"Value"`
}

type ServerTracer struct {
	*Tracer
	ClientRpcId string `json:"client_rpc_id"`
}

type ClientTracer struct {
	*Tracer
}

type Tracer struct {
	// TraceId 调用链ID,一旦初始化,不能修改
	TraceId string
	// RpcId 调用顺序，依次为0 → 0.1 → 0.1.1,1 -> 1.1 -> 1.1.1 ...
	RpcId string
	// TraceType 链路跟踪类型
	TraceType TraceTypeEnum
	// TraceName 链路跟踪名称
	TraceName string
	// Endpoint 跟踪类型
	Endpoint EndpointEnum
	// status 跟踪结果
	status TraceStatusEnum
	// RemoteStatus 远程调用结果
	RemoteStatus TraceStatusEnum
	// RemoteIp 远程调用IP
	RemoteIp string
	// message 调用返回或异常信息
	message string
	// Size 响应体大小
	Size int
	// startTime 当前span开始时间
	startTime int64
	// endTime 当前span结束时间
	endTime int64
	sampled bool
	// bizData 响应数据
	bizData map[string]interface{}
	Ended   bool
	// AttrMap 请求参数
	AttrMap []Parameter
	// ServiceName 当前服务名称
	ServiceName string
}

// NewServerTracer 开启服务端跟踪
func NewServerTracer(req *http.Request) *ServerTracer {
	tracer := New(req)
	serverTracer := &ServerTracer{tracer, ""}
	serverTracer.Endpoint = SERVER
	return serverTracer
}

// NewServerTracerWithoutReq 开启服务端跟踪,此用于服务端定时任务类请求
func NewServerTracerWithoutReq() *ServerTracer {
	tracer := &Tracer{
		TraceId:     LocalIdCreate.GenerateTraceId(),
		sampled:     true,
		ServiceName: conf.Conf.ServiceName,
		startTime:   time.Now().UnixMilli(),
		RpcId:       "0",
		TraceType:   HTTP,
		RemoteIp:    GetLocalIp(),
		TraceName:   "<default>_server",
	}
	return &ServerTracer{tracer, ""}
}

var clientTracerLock sync.Mutex

func (server *ServerTracer) NewClientWithHeader(header *http.Header) *ClientTracer {
	clientTracerLock.Lock()
	defer clientTracerLock.Unlock()
	rpcId := server.ClientRpcId
	if rpcId == "" {
		rpcId = server.RpcId
		rpcId += ".1"
	} else {
		// 获取最后一位 +1
		splits := strings.Split(rpcId, ".")
		lastOne, _ := strconv.Atoi(splits[len(splits)-1])
		lastOne += 1
		splits[len(splits)-1] = strconv.Itoa(lastOne)
		rpcId = strings.Join(splits, ".")
	}
	server.ClientRpcId = rpcId
	// fixme TraceName和Size 需要手动写入
	clientTracer := &ClientTracer{&Tracer{
		TraceId:     server.TraceId,
		sampled:     true,
		ServiceName: conf.Conf.ServiceName,
		startTime:   time.Now().UnixMilli(),
		RpcId:       rpcId,
		TraceType:   HTTP,
		RemoteIp:    GetLocalIp(),
		TraceName:   "<default>_default",
	}}
	header.Set(T_HEADER_TRACEID, server.TraceId)
	header.Set(T_HEADER_RPCID, rpcId)
	return clientTracer
}

// NewClientTracer 开启客户端跟踪
func (server *ServerTracer) NewClientTracer(req *http.Request) *ClientTracer {
	clientTracerLock.Lock()
	defer clientTracerLock.Unlock()
	rpcId := server.ClientRpcId
	if rpcId == "" {
		rpcId = server.RpcId
		rpcId += ".1"
	} else {
		// 获取最后一位 +1
		splits := strings.Split(rpcId, ".")
		lastOne, _ := strconv.Atoi(splits[len(splits)-1])
		lastOne += 1
		splits[len(splits)-1] = strconv.Itoa(lastOne)
		rpcId = strings.Join(splits, ".")
	}

	clientTracer := &ClientTracer{NewWithRpcId(req, rpcId)}
	clientTracer.TraceId = server.TraceId
	clientTracer.Endpoint = CLIENT
	server.ClientRpcId = rpcId
	return clientTracer
}

// NewWithRpcId 自定义rpcId
func NewWithRpcId(req *http.Request, rpcId string) *Tracer {
	tracer := New(req)
	req.Header.Set(T_HEADER_RPCID, rpcId)
	tracer.RpcId = rpcId
	return tracer
}

func New(req *http.Request) *Tracer {
	method := req.Method
	if method == "" {
		method = "nil"
	}
	uri := "nil"
	if url := req.URL; url != nil {
		if uri = url.Path; len(uri) == 0 {
			uri = url.String()
		}
	}
	strLength := req.Header.Get("Content-Length")
	if strLength == "" {
		strLength = "0"
	}
	length, _ := strconv.Atoi(strLength)
	return &Tracer{
		TraceId:     getOrCreateTraceId(req),
		sampled:     true,
		ServiceName: conf.Conf.ServiceName,
		startTime:   time.Now().UnixMilli(),
		RpcId:       getAndIncreaseRpcId(req),
		TraceType:   HTTP,
		RemoteIp:    req.RemoteAddr,
		TraceName:   fmt.Sprintf("<%s>%s", method, uri),
		AttrMap:     parametersCollector(req),
		Size:        length,
	}
}

func (server *ServerTracer) EndServerTracer(status TraceStatusEnum, message string) {
	server.EndTrace(status, message)
}

func (tracer *ClientTracer) EndServerTracer(status TraceStatusEnum, message string) {
	tracer.EndTrace(status, message)
}

// EndTraceOk 快速记录成功请求的链路信息
func (tracer *Tracer) EndTraceOk() {
	tracer.EndTrace(OK, "")
}

// EndTraceError 快速记录失败请求的链路信息
func (tracer *Tracer) EndTraceError(err error) {
	tracer.EndTrace(ERROR, err.Error())
}

func (tracer *Tracer) EndTrace(status TraceStatusEnum, message string) {
	if tracer.Ended {
		log.Default().Println("tracer is ended,will be not append tracer info")
		return
	}
	if tracer.TraceId == "" {
		log.Println("tracer's traceId is nil,will be not append tracer info")
		return
	}
	if tracer.RpcId == "" {
		log.Println("tracer's rpcId is nil,will be not append tracer info")
		return
	}
	if !tracer.sampled {
		log.Println("tracer's sampled is false,will be not append tracer info")
		return
	}
	tracer.Ended = true
	tracer.endTime = time.Now().UnixMilli()
	tracer.status = status
	if message != "" {
		tracer.message = message
	}
	// 记录本地文件或丢到loki的发送队列中
	push.GetStrategy().AddStream([]push.Message{tracer.buildLog()})
}

func (tracer *Tracer) buildLog() push.Message {
	var strItem []string
	result := &push.Message{
		Time: strconv.FormatInt(tracer.endTime, 10) + "000000",
	}
	strItem = append(strItem, "0", "default", strconv.FormatInt(tracer.startTime, 10), tracer.TraceId,
		tracer.RpcId, strconv.Itoa(int(tracer.Endpoint)), strconv.Itoa(int(tracer.TraceType)), tracer.TraceName,
		tracer.ServiceName, GetLocalIp(), tracer.RemoteIp, strconv.Itoa(int(tracer.status)), strconv.Itoa(tracer.Size),
		strconv.FormatInt(tracer.endTime-tracer.startTime, 10), tracer.message)
	if tracer.AttrMap != nil {
		if data, err := json.Marshal(tracer.AttrMap); err != nil {
			// do nothing
		} else {
			strItem = append(strItem, string(data))
		}
	}
	result.Message = strings.Join(strItem, "|")
	return *result
}

type localIp struct {
	LocalIp string
}

var li *localIp

func GetLocalIp() string {
	if li != nil {
		return li.LocalIp
	}
	li = &localIp{
		LocalIp: "127.0.0.1",
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Println("获取本地地址异常", err)
		return li.LocalIp
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := fmt.Sprintf(ipnet.IP.String())
				li = &localIp{ip}
				return li.LocalIp
			}

		}
	}
	return li.LocalIp
}

func parametersCollector(req *http.Request) []Parameter {
	cloneRequest := req.Clone(context.TODO())
	// 读取请求参数
	_ = cloneRequest.ParseForm()
	var parameters []Parameter
	for s, ss := range cloneRequest.Form {
		parameters = append(parameters, Parameter{
			Name:  s,
			In:    "query",
			Value: ss,
		})
	}
	for s, ss := range cloneRequest.PostForm {
		parameters = append(parameters, Parameter{
			Name:  s,
			In:    "form",
			Value: ss,
		})
	}
	if multipartForms := cloneRequest.MultipartForm; multipartForms != nil {
		for s, ss := range multipartForms.Value {
			parameters = append(parameters, Parameter{
				Name:  s,
				In:    "multiform",
				Value: ss,
			})
		}
		for s, headers := range multipartForms.File {
			parameters = append(parameters, Parameter{
				Name: s,
				In:   "multiform",
				Value: func(hs []*multipart.FileHeader) []string {
					var fileNames []string
					for _, h := range hs {
						fileNames = append(fileNames, h.Filename)
					}
					return fileNames
				}(headers),
			})
		}
	}

	for _, cookie := range cloneRequest.Cookies() {
		parameters = append(parameters, Parameter{
			Name:  cookie.Name,
			In:    "cookie",
			Value: []string{cookie.Value},
		})
	}
	if req.Body != nil {
		if data, err := ioutil.ReadAll(req.Body); err != nil {
			//do nothing
		} else {
			parameters = append(parameters, Parameter{
				Name:  "请求体",
				In:    "body",
				Value: []string{string(data)},
			})
		}
	}
	return parameters
}

func getOrCreateTraceId(req *http.Request) string {
	traceId := req.Header.Get(T_HEADER_TRACEID)
	if traceId == "" {
		traceId = LocalIdCreate.GenerateTraceId()
		if req.Header != nil {
			req.Header.Set(T_HEADER_TRACEID, traceId)
		}
	}
	return traceId
}

func getAndIncreaseRpcId(req *http.Request) string {
	rpcId := req.Header.Get(T_HEADER_RPCID)
	if rpcId == "" {
		rpcId = "0"
	}
	if req.Header != nil {
		req.Header.Set(T_HEADER_RPCID, rpcId)
	}
	return rpcId
}
