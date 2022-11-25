package trace

import (
	"github.com/isyscore/isc-gobase/config"
	baseFile "github.com/isyscore/isc-gobase/file"
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

var traceChannel = make(chan *Tracer, 2048)

func SendTraceLog(tracer *Tracer) {
	traceChannel <- tracer
}

func init() {
	//path := "logs/middleware/trace/trace.log"
	path := "logs" + string(os.PathSeparator) + "middleware" + string(os.PathSeparator) + "trace" + string(os.PathSeparator) + "trace.log"

	if !baseFile.FileExists(path) {
		baseFile.CreateFile(path)
	}
	logFile := getTraceLogFile(path)

	go func() {
		for range time.NewTicker(time.Hour * 24).C {
			if logFile != nil {
				//_ = logFile.Truncate(0)
				logger.Warn("定时删除文件")
				logFile.Close()
				baseFile.DeleteFile(path)

				logFile, _ = os.OpenFile(path, os.O_RDWR, 0666)
			}
		}
	}()
	go func() {
		for tracer := range traceChannel {
			if logFile != nil {
				l := newTraceLog(tracer)
				if !baseFile.FileExists(path) {
					logFile.Close()
					baseFile.CreateFile(path)
					logFile = getTraceLogFile(path)
				}

				if _, err := logFile.WriteString(l); err != nil {
					logger.Error("%v", err.Error())
				}
			}
		}
	}()
}

func getTraceLogFile(path string) *os.File {
	logFile, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		logger.Warn("OpenFile err:", err)
	} else {
		_, err = logFile.Seek(0, 2)
		if err != nil {
			logger.Warn("Seek err:", err)
		}
	}
	return logFile
}

func newTraceLog(tracer *Tracer) string {
	s := ""
	s += logVersion + SPLIT
	s += _const.DEFAULT_PROFILES_ACTIVE + SPLIT
	s += strconv.FormatInt(tracer.StartTime, 10) + SPLIT
	s += tracer.TraceId + SPLIT
	s += tracer.RpcId + SPLIT
	s += strconv.FormatInt(int64(tracer.Endpoint), 10) + SPLIT
	s += strconv.FormatInt(int64(tracer.TraceType), 10) + SPLIT
	s += replaceSplit(trimNull(tracer.TraceName)) + SPLIT
	s += replaceSplit(config.GetValueStringDefault("base.application.name", _const.DEFAULT_APP_NAME)) + SPLIT
	s += replaceSplit(util.GetLocalIp()) + SPLIT
	s += replaceSplit(trimNull(tracer.RemoteIp)) + SPLIT
	s += strconv.FormatInt(int64(tracer.status), 10) + SPLIT
	s += strconv.FormatInt(int64(tracer.Size), 10) + SPLIT
	s += strconv.FormatInt(tracer.endTime-tracer.StartTime, 10) + SPLIT
	s += replaceSplit(trimNull(tracer.message)) + SPLIT
	//用户id
	userId := tracer.AttrMap[_const.TRACE_HEAD_USER_ID]
	if userId == "" {
		userId = tracer.AttrMap[_const.A_USER_ID]
	}
	s += trimNull(userId) + SPLIT + "\r\n"
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
