package push

var localStore = &LocalStore{}

type LocalStore struct {
	FilePath string `json:"file_path"`
}

func (client *LocalStore) AddStream(messages []Message) {
	labels := make(map[string]string)
	labels["job"] = "tracelogs"
	client.AddStreamWithLabels(labels, messages)
}

func (client *LocalStore) AddStreamWithLabels(labels map[string]string, messages []Message) {
	var vals []jsonValue
	for i := range messages {
		var val jsonValue
		val[0] = messages[i].Time
		val[1] = messages[i].Message
		vals = append(vals, val)
	}
	// println("add message to stream channel success")
}

func (client *LocalStore) Query(queryString string) ([]Message, error) {
	return nil, nil
}
