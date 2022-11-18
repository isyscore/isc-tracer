package _const

// TraceStatusEnum 标明当前trace的结果
type TraceStatusEnum int

const (
	OK TraceStatusEnum = iota
	ERROR
	WARNING
	TIMEOUT
)

func ParseHttpStatus(status int) TraceStatusEnum {
	if status >= 200 && status < 300 {
		return OK
	}
	return ERROR
}
