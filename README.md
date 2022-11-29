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
  # 采集总开关；默认关闭
  enable: false
  # 数据库相关配置
  database:
    # 是否启用（只有tracer.enable开启情况下才生效）；默认关闭
    enable: false
  # redis相关配置    
  redis:
    # 是否启用（只有tracer.enable开启情况下才生效）；默认关闭
    enable: false
  # etcd相关配置  
  etcd:
    # 是否启用（只有tracer.enable开启情况下才生效）；默认关闭
    enable: false
  # 动态调试功能
  debug:
    # 是否启用，默认关闭
    enable: false
    # 账号获取地址，默认：http://isc-core-back-service:31300
    account: http://http://isc-core-back-service:31300
  # 不采集的http的url配置
  http:
    url:
      excludes:
        - /api/xxx/xxx/xxxx
        - /api/xxx/xxx/xxxx
```

## 代码使用
请在main.go方法这里引入如下的包
```go
import _ "github.com/isyscore/isc-tracer"
```

## 自定义埋点
包trace提供如下的方法
```go
// 发起方埋点：start
func ClientStartTraceWithRequest(req *http.Request) *Tracer {}

// 发起方埋点：start
func ClientStartTrace(traceType _const.TraceTypeEnum, traceName string) *Tracer {}

// 接收方埋点：start
func ServerStartTrace(traceType _const.TraceTypeEnum, traceName string) *Tracer {}


// 发起方/接收方埋点：start
func StartTrace(traceType _const.TraceTypeEnum, traceName string, endPoint _const.EndpointEnum) *Tracer {}

// 发起方/接收方埋点：end
func EndTrace(tracer *Tracer, responseSize int, status _const.TraceStatusEnum, message string) {}
```
### 发起方（服务端）埋点示例：
```go
// 执行前
tracerClient := trace.ClientStartTrace(xxxx, "xxx-name")

// 业务执行
// ......

// 执行后
tracerClient.EndTrace(tracerClient, 0, xxxx, "xxxx")
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
