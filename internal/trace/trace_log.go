package trace

import (
	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/logger"
	_const "github.com/isyscore/isc-tracer/internal/const"
	"github.com/isyscore/isc-tracer/util"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	logVersion = "0"
	SPLIT      = "|"
	NULL_TAG   = "-"
)

var logFile = ""
var traceChannel = make(chan *Tracer, 2048)

func SendTraceLog(tracer *Tracer) {
	traceChannel <- tracer
}

func init() {
	path := "./logs/middleware/trace/trace.log"
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	go func() {
		for range time.NewTicker(time.Hour * 24).C {
			_ = file.Truncate(0)
		}
	}()
	go func() {
		for tracer := range traceChannel {
			l := newTraceLog(tracer)
			_, err := file.Write([]byte(l))
			if err != nil {
				logger.Error("%v", err.Error())
			}
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
	s += replaceSplit(config.GetValueString("base.application.name")) + SPLIT
	s += replaceSplit(util.GetLocalIp()) + SPLIT
	s += replaceSplit(trimNull(tracer.RemoteIp)) + SPLIT
	s += string(rune(tracer.status)) + SPLIT
	s += string(rune(tracer.Size)) + SPLIT
	s += strconv.FormatInt(tracer.endTime-tracer.StartTime, 10) + SPLIT
	s += replaceSplit(trimNull(tracer.message)) + SPLIT
	//用户id
	userId := tracer.AttrMap[_const.TRACE_HEAD_USER_ID]
	if userId == "" {
		userId = tracer.AttrMap[_const.A_USER_ID]
	}
	s += userId + SPLIT
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
