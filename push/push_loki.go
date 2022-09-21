package push

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/isyscore/isc-tracer/conf"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"sync"
	"time"
)

type Endpoints struct {
	push  string
	query string
	ready string
}

type LokiClient struct {
	url            string
	endpoints      Endpoints
	currentMessage jsonMessage
	streams        chan *jsonStream
	quit           chan struct{}
	batchCounter   int
	maxBatch       int
	maxWaitTime    time.Duration
	wait           sync.WaitGroup
}

type returnedJSON struct {
	Status interface{}
	Data   struct {
		ResultType string
		Result     []struct {
			Stream interface{}
			Values [][]string
		}
		Stats interface{}
	}
}

func (client *LokiClient) AddStream(messages []Message) {
	labels := make(map[string]string)
	labels["job"] = "tracelogs"
	client.AddStreamWithLabels(labels, messages)
}

func (client *LokiClient) AddStreamWithLabels(labels map[string]string, messages []Message) {
	var vals []jsonValue
	for i := range messages {
		var val jsonValue
		val[0] = messages[i].Time
		val[1] = messages[i].Message
		vals = append(vals, val)
	}
	stream := &jsonStream{
		Stream: labels,
		Values: vals,
	}
	client.streams <- stream
	// println("add message to stream channel success")
}

func (client *LokiClient) Query(queryString string) ([]Message, error) {
	response, err := http.Get(client.url + client.endpoints.query + "?query=" + queryString)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []Message{}, err
	}
	var answer returnedJSON
	_ = json.Unmarshal(body, &answer)
	var values []Message
	for i := range answer.Data.Result {
		for j := range answer.Data.Result[i].Values {
			msg := Message{
				Time:    answer.Data.Result[i].Values[j][0],
				Message: answer.Data.Result[i].Values[j][1],
			}
			values = append(values, msg)
		}
	}
	return values, nil
}

var httpClient http.Client

var lokiPushStrategy *LokiClient

func InitLockPushStrategy() *LokiClient {
	if lokiPushStrategy == nil {
		lokiPushStrategy = NewLokiClient()
	}
	return lokiPushStrategy
}
func NewLokiClient() *LokiClient {
	client := &LokiClient{
		url:         conf.Conf.Loki.Host,
		maxBatch:    conf.Conf.Loki.MaxBatch,
		maxWaitTime: time.Duration(conf.Conf.Loki.MaxWaitTime) * time.Second,
		quit:        make(chan struct{}),
		streams:     make(chan *jsonStream),
	}
	client.initEndpoints()
	client.wait.Add(1)
	hc := http.Client{
		Timeout: 2 * time.Second,
	}
	httpClient = hc
	go client.run()
	return client
}

func (client *LokiClient) initEndpoints() {
	client.endpoints.push = "/loki/api/v1/push"
	client.endpoints.query = "/loki/api/v1/query"
	client.endpoints.ready = "/ready"
}

func (client *LokiClient) Shutdown() {
	close(client.quit)
	client.wait.Wait()
}

func (client *LokiClient) run() {
	batchCounter := 0
	maxWait := time.NewTimer(client.maxWaitTime)

	defer func() {
		if batchCounter > 0 {
			_ = client.send()
		}
		client.wait.Done()
	}()

	for {
		select {
		case <-client.quit:
			return
		case stream := <-client.streams:
			client.currentMessage.Streams = append(client.currentMessage.Streams, *stream)
			batchCounter++
			if batchCounter == client.maxBatch {
				_ = client.send()
				batchCounter = 0
				client.currentMessage.Streams = []jsonStream{}
				maxWait.Reset(client.maxWaitTime)
			}
		case <-maxWait.C:
			if batchCounter > 0 {
				_ = client.send()
				client.currentMessage.Streams = []jsonStream{}
				batchCounter = 0
			}
			maxWait.Reset(client.maxWaitTime)
		}
	}
}

var traceCtx = httptrace.WithClientTrace(context.Background(), &httptrace.ClientTrace{})

// send Encodes the messages and sends them to loki
func (client *LokiClient) send() error {
	//log.Info().Msgf("执行发送任务")
	length := len(client.currentMessage.Streams)
	if length == 0 {
		return nil
	}
	defer func() {
		client.currentMessage.Streams = []jsonStream{}
	}()
	str, err := json.Marshal(client.currentMessage)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(traceCtx, "POST", client.url+client.endpoints.push, bytes.NewBuffer(str))
	if err != nil {
		return err
	}
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(str)))
	response, err := httpClient.Do(req)
	if err != nil {
		return err
	} else if response != nil && response.StatusCode != 204 {
		// println("上传条数", length)
		return err
	}
	// println("成功上传条数", length)
	return nil
}
