package _const

// TraceTypeEnum 标明链路跟踪类型
type TraceTypeEnum int

const (
	ROOT TraceTypeEnum = iota
	HTTP
	DUBBO
	MYSQL
	ROCKETMQ
	REDIS
	KAFKA
	IDS
	MQTT
	ORACLE
	ELASTIC
	ZOOKEEPER
	HBASE
	HADOOP
	FLINK
	SPARK
	KUDU
	HIVE
	STORM
	CONFIG
)
