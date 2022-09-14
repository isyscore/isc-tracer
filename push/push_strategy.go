package push

import "github.com/isyscore/isc-tracer/conf"

type jsonValue [2]string

type jsonStream struct {
	Stream map[string]string `json:"stream"`
	Values []jsonValue       `json:"values"`
}

type Message struct {
	Message string
	Time    string
}

type jsonMessage struct {
	Streams []jsonStream `json:"streams"`
}

type SendStrategy interface {
	AddStream(messages []Message)
	AddStreamWithLabels(labels map[string]string, messages []Message)
	Query(queryString string) ([]Message, error)
}

func GetStrategy() SendStrategy {
	if conf.Conf.Using == "loki" {
		return InitLockPushStrategy()
	}
	return nil
}
