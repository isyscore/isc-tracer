# isc-tracer

该框架是基于isc-gobase之上开发的用于链路搜集的sdk

该框架目前支持以下的相关客户端埋点。以下客户端全部都要是基于isc-gobase提供的客户端才行，否则请用户自行埋点

- http
- orm
  - gorm
  - xorm  
- redis
  - go-redis
- etcd
  - go-etcd
  
相关isc-gobase的客户端接入请见这里 [isc-gobase/extend](https://github.com/isyscore/isc-gobase/tree/feature/trace/extend)


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
```

代码引入
请在main.go方法这里引入如下的包
```go
import _ "github.com/isyscore/isc-tracer"
```

## 自定义埋点
```go

```
