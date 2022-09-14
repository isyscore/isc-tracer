package trace

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var httpClient = createHTTPClient()

const (
	MaxIdleConns          int    = 100
	MaxIdleConnsPerHost   int    = 100
	IdleConnTimeout       int    = 90
	ContentTypeJson       string = "application/json; charset=utf-8"
	ContentTypeHtml       string = "text/html; charset=utf-8"
	ContentTypeText       string = "text/plain; charset=utf-8"
	ContentTypeCss        string = "text/css; charset=utf-8"
	ContentTypeJavaScript string = "application/x-javascript; charset=utf-8"
	ContentTypeJpeg       string = "image/jpeg"
	ContentTypePng        string = "image/png"
	ContentTypeGif        string = "image/gif"
	ContentTypeAll        string = "*/*"
)

type NetError struct {
	ErrMsg string
}

func (error *NetError) Error() string {
	return error.ErrMsg
}

type DataResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
		},

		Timeout: 20 * time.Second,
	}
	return client
}

func SetHttpClient(httpClientOuter *http.Client) {
	httpClient = httpClientOuter
}

func (server *ServerTracer) DoRequest(req *http.Request) (*http.Response, error) {
	clientTracer := server.NewClientTracer(req)
	resp, err := httpClient.Do(req)
	defer func(response *http.Response, errorInfo error, tracer *ClientTracer) {
		if err != nil {
			tracer.EndTraceError(err)
		} else if response.StatusCode > http.StatusIMUsed {
			tracer.EndTrace(WARNING, "访问失败，响应码="+strconv.Itoa(response.StatusCode))
		} else {
			tracer.EndTraceOk()
		}
	}(resp, err, clientTracer)
	return resp, err
}

// ------------------ get ------------------

func (server *ServerTracer) GetSimple(url string) (int, http.Header, []byte, error) {
	return server.Get(url, nil, nil)
}

func (server *ServerTracer) GetSimpleOfStandard(url string) (int, http.Header, []byte, error) {
	return server.GetOfStandard(url, nil, nil)
}

func (server *ServerTracer) Get(url string, header http.Header, parameterMap map[string]string) (int, http.Header, []byte, error) {
	httpRequest, err := http.NewRequest("GET", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return server.call(httpRequest, url)
}

func (server *ServerTracer) GetOfStandard(url string, header http.Header, parameterMap map[string]string) (int, http.Header, []byte, error) {
	httpRequest, err := http.NewRequest("GET", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return server.callToStandard(httpRequest, url)
}

// ------------------ head ------------------

func (server *ServerTracer) HeadSimple(url string) error {
	return server.Head(url, nil, nil)
}

func (server *ServerTracer) Head(url string, header http.Header, parameterMap map[string]string) error {
	httpRequest, err := http.NewRequest("GET", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return server.callIgnoreReturn(httpRequest, url)
}

// ------------------ post ------------------

func (server *ServerTracer) PostSimple(url string, body any) (int, http.Header, []byte, error) {
	return server.Post(url, nil, nil, body)
}

func (server *ServerTracer) PostSimpleOfStandard(url string, body any) (int, http.Header, []byte, error) {
	return server.PostOfStandard(url, nil, nil, body)
}

func (server *ServerTracer) Post(url string, header http.Header, parameterMap map[string]string, body any) (int, http.Header, []byte, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("POST", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", ContentTypeJson)
	return server.call(httpRequest, url)
}

func (server *ServerTracer) PostOfStandard(url string, header http.Header, parameterMap map[string]string, body any) (int, http.Header, []byte, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("POST", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", ContentTypeJson)
	return server.callToStandard(httpRequest, url)
}

// ------------------ put ------------------

func (server *ServerTracer) PutSimple(url string, body any) (int, http.Header, []byte, error) {
	return server.Put(url, nil, nil, body)
}

func (server *ServerTracer) PutSimpleOfStandard(url string, body any) (int, http.Header, []byte, error) {
	return server.PutOfStandard(url, nil, nil, body)
}

func (server *ServerTracer) Put(url string, header http.Header, parameterMap map[string]string, body any) (int, http.Header, []byte, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("PUT", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", ContentTypeJson)
	return server.call(httpRequest, url)
}

func (server *ServerTracer) PutOfStandard(url string, header http.Header, parameterMap map[string]string, body any) (int, http.Header, []byte, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("PUT", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", ContentTypeJson)
	return server.callToStandard(httpRequest, url)
}

// ------------------ delete ------------------

func (server *ServerTracer) DeleteSimple(url string) (int, http.Header, []byte, error) {
	return server.Get(url, nil, nil)
}

func (server *ServerTracer) DeleteSimpleOfStandard(url string) (int, http.Header, []byte, error) {
	return server.GetOfStandard(url, nil, nil)
}

func (server *ServerTracer) Delete(url string, header http.Header, parameterMap map[string]string) (int, http.Header, []byte, error) {
	httpRequest, err := http.NewRequest("DELETE", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return server.call(httpRequest, url)
}

func (server *ServerTracer) DeleteOfStandard(url string, header http.Header, parameterMap map[string]string) (int, http.Header, []byte, error) {
	httpRequest, err := http.NewRequest("DELETE", urlWithParameter(url, parameterMap), nil)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}

	return server.callToStandard(httpRequest, url)
}

// ------------------ patch ------------------

func (server *ServerTracer) PatchSimple(url string, body any) (int, http.Header, []byte, error) {
	return server.Post(url, nil, nil, body)
}

func (server *ServerTracer) PatchSimpleOfStandard(url string, body any) (int, http.Header, []byte, error) {
	return server.PostOfStandard(url, nil, nil, body)
}

func (server *ServerTracer) Patch(url string, header http.Header, parameterMap map[string]string, body any) (int, http.Header, []byte, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("PATCH", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", ContentTypeJson)
	return server.call(httpRequest, url)
}

func (server *ServerTracer) PatchOfStandard(url string, header http.Header, parameterMap map[string]string, body any) (int, http.Header, []byte, error) {
	bytes, _ := json.Marshal(body)
	payload := strings.NewReader(string(bytes))
	httpRequest, err := http.NewRequest("PATCH", urlWithParameter(url, parameterMap), payload)
	if err != nil {
		log.Printf("NewRequest error(%v)\n", err)
		return -1, nil, nil, err
	}

	if header != nil {
		httpRequest.Header = header
	}
	httpRequest.Header.Add("Content-Type", ContentTypeJson)
	return server.callToStandard(httpRequest, url)
}

func (server *ServerTracer) call(httpRequest *http.Request, url string) (int, http.Header, []byte, error) {
	// 开始客户端跟踪
	clientTracer := server.NewClientTracer(httpRequest)
	if httpResponse, err := httpClient.Do(httpRequest); err != nil && httpResponse == nil {
		defer func(e error) {
			clientTracer.EndTrace(ERROR, e.Error())
		}(err)
		log.Printf("Error sending request to API endpoint. %+v", err)
		return -1, nil, nil, &NetError{ErrMsg: "Error sending request, url: " + url + ", err" + err.Error()}
	} else {
		if httpResponse == nil {
			clientTracer.EndTrace(OK, "")
			log.Printf("httpResponse is nil\n")
			return -1, nil, nil, nil
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Body close error(%v)", err)
			}
		}(httpResponse.Body)

		code := httpResponse.StatusCode
		headers := httpResponse.Header
		if code != http.StatusOK {
			body, _ := ioutil.ReadAll(httpResponse.Body)
			errMsg := &NetError{ErrMsg: "remote error, url: " + url + ", code " + strconv.Itoa(code) + ", message: " + string(body)}
			clientTracer.EndTrace(ERROR, errMsg.Error())
			return code, headers, nil, errMsg
		}

		// We have seen inconsistencies even when we get 200 OK response
		body, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			log.Printf("Couldn't parse response body(%v)", err)
			errMsg := &NetError{ErrMsg: "Couldn't parse response body, err: " + err.Error()}
			clientTracer.EndTrace(ERROR, errMsg.Error())
			return code, headers, nil, errMsg
		}
		respBodyLength := httpResponse.Header.Get("Content-Length")
		if length, _ := strconv.Atoi(respBodyLength); length < 64 {
			clientTracer.EndTrace(OK, string(body))
		} else {
			clientTracer.EndTrace(OK, "")
		}
		return code, headers, body, nil
	}
}

// ------------------ trace ------------------
// ------------------ options ------------------
// 暂时先不处理

func (server *ServerTracer) callIgnoreReturn(httpRequest *http.Request, url string) error {
	clientTracer := server.NewClientTracer(httpRequest)
	if httpResponse, err := httpClient.Do(httpRequest); err != nil && httpResponse == nil {
		log.Printf("Error sending request to API endpoint. %v", err)
		clientTracer.EndTrace(ERROR, err.Error())
		return &NetError{ErrMsg: "Error sending request, url: " + url + ", err" + err.Error()}
	} else {
		if httpResponse == nil {
			clientTracer.EndTrace(OK, "")
			log.Printf("httpResponse is nil\n")
			return nil
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Body close error(%v)", err)
			}
		}(httpResponse.Body)

		code := httpResponse.StatusCode
		if code != http.StatusOK {
			body, _ := ioutil.ReadAll(httpResponse.Body)
			clientTracer.EndTrace(ERROR, string(body))
			return &NetError{ErrMsg: "remote error, url: " + url + ", code " + strconv.Itoa(code) + ", message: " + string(body)}
		}
		clientTracer.EndTrace(OK, "")
		return nil
	}
}

func (server *ServerTracer) callToStandard(httpRequest *http.Request, url string) (int, http.Header, []byte, error) {
	return parseStandard(server.call(httpRequest, url))
}

func parseStandard(statusCode int, header http.Header, responseResult []byte, errs error) (int, http.Header, []byte, error) {
	if errs != nil {
		return statusCode, header, nil, errs
	}
	var standRsp DataResponse[any]
	err := json.Unmarshal(responseResult, &standRsp)
	if err != nil {
		return statusCode, header, nil, err
	}

	// 判断业务的失败信息
	if standRsp.Code != 0 && standRsp.Code != 200 {
		return statusCode, header, nil, &NetError{ErrMsg: fmt.Sprintf("remote err, bizCode=%d, message=%s", standRsp.Code, standRsp.Message)}
	}

	if data, err := json.Marshal(standRsp.Data); err != nil {
		return statusCode, header, nil, err
	} else {
		return statusCode, header, data, nil
	}
}

func urlWithParameter(url string, parameterMap map[string]string) string {
	if parameterMap == nil || len(parameterMap) == 0 {
		return url
	}

	url += "?"

	var parameters []string
	for key, value := range parameterMap {
		parameters = append(parameters, key+"="+value)
	}

	return url + strings.Join(parameters, "&")
}
