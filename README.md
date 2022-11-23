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

## 配置
```yaml
base:
  application:
    # 服务名：用于埋点的服务名使用，如果不配置，则为默认：isc-tracer-default-name
    name: isc-biz-xxx
  
tracer:
  # 采集总开关；默认开启；
  enable: true
```

## 自定义埋点
```yaml
// todo
```
