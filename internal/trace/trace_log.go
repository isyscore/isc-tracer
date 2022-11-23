package trace

import (
	"github.com/isyscore/isc-tracer/config"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/util"
	"strconv"
	"strings"
)

const (
	logVersion = "0"
	SPLIT      = "|"
	NULL_TAG   = "-"
)

var logFile =
var traceChannel = make(chan *Tracer, 2048)

func SendTraceLog(tracer *Tracer) {
	traceChannel <- tracer
}

func init() {

	go func() {
		for tracer := range traceChannel {
			traceLog := newTraceLog(tracer)

		}
	}()

}

func newTraceLog(tracer *Tracer) string {
	s := ""
	s += logVersion + SPLIT
	s += _const.DEFAULT_PROFILES_ACTIVE + SPLIT
	s += strconv.FormatInt(tracer.StartTime, 10) + SPLIT
	s += tracer.TraceId + SPLIT
	s += tracer.RpcId + SPLIT
	s += string(rune(tracer.Endpoint)) + SPLIT
	s += string(rune(tracer.TraceType)) + SPLIT
	s += replaceSplit(trimNull(tracer.TraceName)) + SPLIT
	s += replaceSplit(config.GetConfig().ServiceName) + SPLIT
	s += replaceSplit(util.GetLocalIp()) + SPLIT
	s += replaceSplit(trimNull(tracer.RemoteIp)) + SPLIT
	s += string(rune(tracer.status)) + SPLIT
	s += string(rune(tracer.Size)) + SPLIT
	s += strconv.FormatInt(tracer.endTime-tracer.StartTime, 10) + SPLIT
	s += replaceSplit(trimNull(tracer.message)) + SPLIT
	return s
}

func replaceSplit(str string) string {
	//todo (char) 28; //不可见字符
	s := strings.ReplaceAll(str, "\n", "")
	return strings.ReplaceAll(s, "|", "")
}
func trimNull(str string) string {
	if str == "" {
		return NULL_TAG
	}
	return str
}
