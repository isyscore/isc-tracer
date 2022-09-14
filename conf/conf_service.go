package conf

var Conf = &ServiceConf{
	ServiceName: "default",
	Using:       "loki",
	Loki: lokiConf{
		Host:        "http://loki-service:3100",
		MaxBatch:    64,
		MaxWaitTime: 1,
	},
}

type lokiConf struct {
	Host        string `json:"host" yaml:"host"`
	MaxBatch    int    `json:"max_batch" yaml:"maxBatch"`
	MaxWaitTime int64  `json:"max_wait_time" yaml:"maxWaitTime"`
}

// ServiceConf 当前服务配置信息
type ServiceConf struct {
	ServiceName string `json:"service_name" yaml:"serviceName"`
	// push 策略，枚举值:loki、local，默认loki
	Using string   `json:"using" yaml:"using"`
	Loki  lokiConf `json:"loki" yaml:"loki"`
}
