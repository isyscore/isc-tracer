# isc-tracer

该框架是基于isc-gobase之上开发的用于链路搜集的sdk。目前支持以下协议的埋点，客户端全部都要是基于isc-gobase提供的客户端才行，其他一些协议请自行埋点

- http
- orm
  - gorm
  - xorm
- redis
  - go-redis
- etcd
  - go-etcd
  
相关isc-gobase的客户端接入请见这里 [isc-gobase/extend](https://github.com/isyscore/isc-gobase/tree/feature/trace/extend)


## 下载
```go
go get github.com/isyscore/isc-tracer
```

## 配置
```yaml
base:
  application:
    # 服务名：用于埋点的服务名使用，如果不配置，则为默认：isc-tracer-default-name
    name: isc-biz-xxx
  
tracer:
  # 采集总开关；默认开启；注意：关闭后可能会造成链路不连续
  enable: true
  # 数据库相关配置
  database:
    # 是否启用（只有tracer.enable开启情况下才有效）；默认开启
    enable: true
  # redis相关配置    
  redis:
    # 是否启用（只有tracer.enable开启情况下才有效）；默认开启
    enable: true
  # etcd相关配置  
  etcd:
    # 是否启用（只有tracer.enable开启情况下才有效）；默认开启
    
    enable: true
  # 动态调试功能
  debug:
    # 是否启用，默认开启
    enable: true
    # 账号获取地址，默认：http://isc-core-back-service:31300
    account: http://isc-core-back-service:31300
  # 不采集的http的url配置
  http:
    excludes-url:
      - /api/xxx/xxx/xxxx
      - /api/xxx/xxx/xxxx
  # 用户名注册位置
  url:
    pivot: http://isc-pivot-platform:31107
  server:
    # 链路的服务端接收地址，这里采用grpc，这里默认如下
    url: isc-pivot-platform:31107
```

## 代码使用
请在main.go方法这里引入如下的包
```go
import _ "github.com/isyscore/isc-tracer"
```

## 自定义埋点
包trace提供如下的方法
```go
import "github.com/isyscore/isc-tracer/trace"

// 发起方埋点：start
func ClientStartTrace(traceType _const2.TraceTypeEnum, traceName string) *Tracer {}

// 发起方埋点：start
func ClientStartTraceWithRequest(req *http.Request) *Tracer {}

// 接收方埋点：start
func ServerStartTrace(traceType _const2.TraceTypeEnum, traceName string) *Tracer {}

// 接收方埋点：start
func ServerStartTraceWithRequest(traceType _const2.TraceTypeEnum, traceName string, request *http.Request) *Tracer {}


// 发起方/接收方埋点：start
func StartTrace(traceType _const2.TraceTypeEnum, endPoint _const2.EndpointEnum, traceName string, request *http.Request) *Tracer {}

// 发起方/接收方埋点：end
func EndTrace(tracer *Tracer, status _const2.TraceStatusEnum, message string, responseSize int) {}
func EndTraceOk(tracer *Tracer, message string, responseSize int) {}
func EndTraceTimeout(tracer *Tracer, message string, responseSize int) {}
func EndTraceWarn(tracer *Tracer, message string, responseSize int) {}
func EndTraceError(tracer *Tracer, message string, responseSize int) {}
```
### 发起方（服务端）埋点示例：
```go
// 执行前
tracerClient := trace.ClientStartTrace(xxxx, "xxx-name")

// 业务执行
// ......

// 执行后
tracer.EndTrace(tracerClient, 0, xxxx, "xxxx")
```
### 接收方（服务端）埋点示例：
```go
// 执行前埋点
tracerServer := trace.ServerStartTrace(_const.HTTP, "xxx")

// 业务执行
// ......

// 业务执行完后埋点
trace.EndTrace(tracerServer, size, status, "xxx")
```

### 项目开发者阅读
### 安装
记住以下两个必须都安装
```properties
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

### grpc代码生成
项目主目录
```properties
protoc --go_out=. --go-grpc_out=. ./protobuf/pivot.proto
```
