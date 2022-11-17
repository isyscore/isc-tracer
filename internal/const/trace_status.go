package _const

// TraceStatusEnum 标明当前trace的结果
type TraceStatusEnum int

const (
	OK TraceStatusEnum = iota
	ERROR
	WARNING
	TIMEOUT
)
