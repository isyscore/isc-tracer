package trace

import (
	"github.com/isyscore/isc-gobase/config"
	baseFile "github.com/isyscore/isc-gobase/file"
	"github.com/isyscore/isc-gobase/goid"
	"github.com/isyscore/isc-gobase/logger"
	_const2 "github.com/isyscore/isc-tracer/const"
	"github.com/isyscore/isc-tracer/util"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	logVersion = "0"
	SPLIT      = "|"
	NULL_TAG   = "-"
)

var traceChannel = make(chan *Tracer, 8182)

// var logFileWriter *bufio.Writer
var logFile *os.File
var lock sync.Mutex
var arr []string

func SendTraceLog(tracer *Tracer) {
	traceChannel <- tracer
}

func init() {
	arr = make([]string, 0)
	lock = sync.Mutex{}
	//path := "logs/middleware/trace/{ip}/trace.log"
	path := "logs" + string(os.PathSeparator) + "middleware" + string(os.PathSeparator) +
		"trace" + string(os.PathSeparator) + util.GetLocalIp() + string(os.PathSeparator) + "trace.log"

	if !baseFile.FileExists(path) {
		baseFile.CreateFile(path)
	}
	logFile = getTraceLogFile(path)
	//logFileWriter = bufio.NewWriter(logFile)

	goid.Go(func() {
		for range time.NewTicker(time.Hour * 24).C {
			if logFile != nil {
				//_ = logFile.Truncate(0)
				//logger.Warn("定时删除文件")

				//_ = logFileWriter.Flush()
				_ = logFile.Close()
				baseFile.DeleteFile(path)

				logFile, _ = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
				//logFileWriter = bufio.NewWriter(logFile)
			}
		}
	})
	goid.Go(func() {
		for range time.NewTicker(time.Second).C {
			if logFile != nil {
				lock.Lock()
				flush()
				lock.Unlock()
			}
		}
	})
	goid.Go(func() {
		for tracer := range traceChannel {
			if logFile != nil {
				l := newTraceLog(tracer)

				lock.Lock()
				arr = append(arr, l)
				if len(arr) >= 100 {
					flush()
				}
				lock.Unlock()
			}
		}
	})
}

func flush() {
	if len(arr) == 0 {
		return
	}
	writeContent := ""
	for _, traceLog := range arr {
		writeContent += traceLog
	}
	_, err := logFile.Write([]byte(writeContent))
	if err != nil {
		logger.Warn("文件写入异常, %v", err)
		return
	}
	arr = arr[:0]
	//if err := logFileWriter.Flush(); err != nil {
	//	logger.Warn("定时刷新文件异常, %v", err)
	//}
}
func Close() {
	logger.Warn("trace_log Close!")
	//if logFileWriter != nil {
	//	logFileWriter.Flush()
	//}
	flush()
	if logFile != nil {
		logFile.Close()
	}
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
	s += _const2.DEFAULT_PROFILES_ACTIVE + SPLIT
	s += strconv.FormatInt(tracer.StartTime, 10) + SPLIT
	s += tracer.TraceId + SPLIT
	s += tracer.RpcId + SPLIT
	s += strconv.FormatInt(int64(tracer.Endpoint), 10) + SPLIT
	s += strconv.FormatInt(int64(tracer.TraceType), 10) + SPLIT
	s += replaceSplit(trimNull(tracer.TraceName)) + SPLIT
	s += replaceSplit(getAppName()) + SPLIT
	s += replaceSplit(util.GetLocalIp()) + SPLIT
	s += replaceSplit(trimNull(tracer.RemoteIp)) + SPLIT
	s += strconv.FormatInt(int64(tracer.Status), 10) + SPLIT
	s += strconv.FormatInt(int64(tracer.Size), 10) + SPLIT
	s += strconv.FormatInt(tracer.EndTime-tracer.StartTime, 10) + SPLIT
	s += replaceSplit(trimNull(tracer.Message)) + SPLIT
	userId := tracer.AttrMap[_const2.TRACE_HEAD_USER_ID]
	if userId == "" {
		userId = tracer.AttrMap[_const2.A_USER_ID]
	}
	s += trimNull(userId) + SPLIT + "\r\n"
	return s
}

func getAppName() string {
	appName := config.GetValueString("base.application.name")
	if appName != "" {
		return appName
	}
	appName = os.Getenv("appName")
	if appName != "" {
		return appName
	}
	return _const2.DEFAULT_APP_NAME
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
