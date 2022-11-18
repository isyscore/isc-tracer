package tracing

import (
	"fmt"
	"sync"

	"github.com/isyscore/isc-gobase/config"
	"github.com/isyscore/isc-gobase/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

var GlobalTracing opentracing.Tracer
var traceConfig bool
var loadLock sync.Mutex

func InitTracing() error {
	loadLock.Lock()
	defer loadLock.Unlock()
	if traceConfig {
		return nil
	}

	serviceName := config.GetValueStringDefault("base.application.name", "gobase-default")
	collectorEndpoint := config.GetValueStringDefault("base.tracing.collector-endpoint", "http://isc-core-back-service:31300/api/core/back/v1/middle/spans")

	conf := jaegerConfig.Configuration{
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		ServiceName: serviceName,
		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: collectorEndpoint,
		},
	}

	tracer, _, err := conf.NewTracer(jaegerConfig.Logger(&baseJaegerLogger{}))
	if err != nil {
		logger.Warn("globalTracer 插件初始化失败, 错误原因: %v", err)
		return err
	}

	GlobalTracing = tracer
	opentracing.SetGlobalTracer(GlobalTracing)

	traceConfig = true
	return nil
}

type baseJaegerLogger struct{}

func (l *baseJaegerLogger) Error(msg string) {
	logger.Error("ERROR: %s", msg)
}

func (l *baseJaegerLogger) Infof(msg string, args ...interface{}) {
	printLog := config.GetValueBoolDefault("base.tracing.print-log", false)
	if printLog {
		logger.Info(msg, args...)
	}
}

func (l *baseJaegerLogger) Debugf(msg string, args ...interface{}) {
	logger.Debug(fmt.Sprintf("DEBUG: %s", msg), args...)
}

func GetHeaderWithKey(s string) string {
	return s
}
func GetHeader() opentracing.HTTPHeadersCarrier {
	return nil
}
