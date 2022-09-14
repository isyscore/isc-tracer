package trace

type TraceIdIntf interface {
	// GenerateTraceId 生成或获取到唯一traceId值
	GenerateTraceId() string
}
