package _const

// TraceTypeEnum 标明链路跟踪类型
type TraceTypeEnum int

const (
	ROOT       TraceTypeEnum = iota
	HTTP                     // 1
	DUBBO                    // 2
	MYSQL                    // 3
	ROCKETMQ                 // 4
	REDIS                    // 5
	KAFKA                    // 6
	IDS                      // 7
	MQTT                     // 8
	ORACLE                   // 9
	ELASTIC                  // 10
	ZOOKEEPER                // 11
	HBASE                    // 12
	HADOOP                   // 13
	FLINK                    // 14
	SPARK                    // 15
	KUDU                     // 16
	HIVE                     // 17
	STORM                    // 18
	CONFIG                   // 19
	ETCD                     // 20
	POSTGRESQL               // 21
	SQLITE                   // 22
	UNKNOWN    = 100
)
