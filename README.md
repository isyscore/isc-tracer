# 简介
链路跟踪bintools/tracer为分布式应用提供了完整的调用链路还原、调用请求量统计、链路拓扑、应用依赖分析等工具，可以帮助开发者快速分析和诊断分布式应用架构下的性能瓶颈，提供微服务时代下的开发诊断效率
# 主要功能
+ 分布式调用链查询和诊断：追踪分布式架构中的所有微服务用户请求，并将它们汇总成分布式调用链
+ 分布式拓扑动态发现：用户的所有分布式微服务应用和相关产品可以通过链路追踪收集到分布式调用信息
+ 丰富的下游对接场景：收集的链路可直接用于日志分析，且可对接到系统管理-运维中心等下游分析平台。
# 快速开始
## Install
```bash
go get github.com/kuchensheng/bintools/tracer
```
## Example

//初始化配置信息
var Conf = &ServiceConf{
	//当前服务名
    ServiceName: "default",
	//链路跟踪保存策略，默认Loki保存
    Using:       "loki",
	//Loki保存策略的配置信息
    Loki: lokiConf{
		//Loki地址
        Host:        "http://loki-service:3100",
		//批量提交的最大值，默认64条
        MaxBatch:    512,
		//提交前最大的等待时间，单位秒，默认1秒
        MaxWaitTime: 1,
        },
}
//create a server tracer
"github.com/kuchensheng/bintools/tracer/conf"
"github.com/kuchensheng/bintools/tracer/push"
func testReq(req *http.Request)  {
	//开启服务端跟踪
    serverTracer := NewServerTracer(req)
	//如果是定时任务，或其他自发性的请求，用以下方式开启服务端跟踪
    //serverTracer1 := NewServerTracerWithoutReq()
    println("服务端其他业务请求")
    
    for i := 0; i < 3; i++ {
        println("作为客户端，向其他服务发起请求")
		req1 := &http.Request{}
		//clientTracer也支持仅有请求头的处理
        //clientTracer := serverTracer.NewClientWithHeader(header)
        //clientTracer.TraceName = "自定义traceName，默认:<Method>uri"
        //clientTracer.AttrMap = []Parameter{}
		//开启客户端跟踪
        clientTracer := serverTracer.NewClientTracer(req1)
		println("req1请求处理以及其他业务处理")
		//结束当前客户端请求跟踪
        clientTracer.EndTrace(OK, "i am danger")
    }
	//结束服务端跟踪
    serverTracer.EndTrace(OK, "i am not in danger")
}
tracer模块也提供了http请求封装,这些请求都被serverTracer所包裹。
分装包括了基本的GET|POST|PUT|DELETE请求
示例如下
```go
req := &http.Request{}
server := NewServerTracer(req)
//如果是自发请求，无需req也可创建serverTracer
//server := NewServerTracerWithoutReq()
url := "www.baidu.com"
header := map[string][]string{"id": {"kucs"}}
parameter := map[string]string{"name":"库陈胜"}
server.GetSimple("www.baidu.com")
server.Get(url,header,parameter)
//server所有业务完成后,结束处理
server.EndTraceOK()
```
## 上报的内容格式
```text
0|default|1662021286867|000100000182f8303bcd0a0070c54e68|1.1|1|1|<GET>http://localhost:8080?id=23|default|10.0.112.197|192.168.10.97|0|0|3|i am not in danger|[{"name":"isyscoreOS","in":"query"},{"name":"id","in":"query"},{"name":"isyscoreOS","in":"form"}]
0|default|1662021286867|000100000182f8303bcd0a0070c54e68|1.1.1|0|1|<GET>http://localhost:8080?id=23|default|10.0.112.197|192.168.10.97|0|0|0|i am danger|[{"name":"isyscoreOS","in":"query"},{"name":"id","in":"query"},{"name":"isyscoreOS","in":"form"}]
0|default|1662021286870|000100000182f8303bcd0a0070c54e68|1.1.2|0|1|<GET>http://localhost:8080?id=23|default|10.0.112.197|192.168.10.97|0|0|0|i am danger|[{"name":"isyscoreOS","in":"query"},{"name":"id","in":"query"},{"name":"isyscoreOS","in":"form"}]
0|default|1662021286870|000100000182f8303bcd0a0070c54e68|1.1.3|0|1|<GET>http://localhost:8080?id=23|default|10.0.112.197|192.168.10.97|0|0|0|i am danger|[{"name":"isyscoreOS","in":"query"},{"name":"id","in":"query"},{"name":"isyscoreOS","in":"form"}]
```
字段释义

 | 字段             | 描述                                          | 
|----------------|---------------------------------------------|
 | version        | 记录日志版本号，用于日志格式解析，默认：0                       |
 | profilesActive | 环境,固定值default                               |
| startTime      | 该日志的开始记录时间                                  |
| traceId        | 跟踪ID，在整条链路上传递                               |
| rpcId          | spanId，标识调用广度及深度                            | 
 | endpoint       | 表示日志打印端，1-服务端,0-客户端                         |
| version        | 日志格式版本，固定值”1“                               |
| traceName      | 跟踪名称，http类型的接口<Method>URI                   |
| traceType      | 跟踪类型，HTTP=1                                 |
| appName        | 当前endpoint的应用名称,如果没有，则默认default             |
| localIp        | 当前endpoint的IP地址                             |
| remoteIp       | 发起rpc的远端IP地址                                |
| status         | 本次跟踪的结果，0-成功，1-失败                           |
| size           | 本次请求跟踪请求的大小                                 |
| span           | 结束跟踪时计算从startTime到当前的时间差，即为本次跟踪的耗时          |
| message        | 结束跟踪时记录的简单信息，如OK，Exception等等，便于快速了解问题出在什么地方 |
| attrMap        | 本次请求的请求入参                                   |

# 数据如何上报？
见下回分解
