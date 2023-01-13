package api_client

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

type ApiClient struct {
	Endpoint string
	Headers  map[string]string
	Body     interface{}
	Resty    *resty.Client
	Start    string
	End      string
}

func New() *ApiClient {
	client := resty.New()
	return &ApiClient{
		Headers: map[string]string{},
		Resty:   client,
	}
}

func (a *ApiClient) SetEndPoint(endpoint string) {
	a.Endpoint = endpoint
}

func (a *ApiClient) SetHeader(header map[string]string) {
	a.Headers = header
}
func (a *ApiClient) SetBody(body interface{}) {
	a.Body = body
}

func (a *ApiClient) SetRequest(endpoint string, header map[string]string, body interface{}) {
	a.Endpoint = endpoint
	a.Headers = header
	a.Body = body
}

func (a *ApiClient) Post(response interface{}) (interface{}, error) {
	a.Start = time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("Endpoint : ", a.Endpoint)
	fmt.Println("Header : ", a.Headers)
	fmt.Println("Request : ", Stringify(a.Body))

	data, err := a.Resty.SetPreRequestHook(a.BeforeRequest).R().SetBody(a.Body).Post(a.Endpoint)

	if err != nil {
		logs := a.ApiLog("Post", a.Endpoint, data)
		return logs, status.Error(codes.Internal, err.Error())
	}

	fmt.Println("Res Header : ", data.RawResponse.Header)
	fmt.Println("Res Body : ", string(data.Body()))

	if data.RawResponse.StatusCode != 200 {
		logs := a.ApiLog("Post", a.Endpoint, data)
		return logs, status.Error(codes.Internal, data.Status())
	}

	var body = data.Body()
	if err = json.Unmarshal(body, response); err != nil {
		logs := a.ApiLog("Post", a.Endpoint, data)
		return logs, status.Error(codes.Internal, data.Status())
	}
	logs := a.ApiLog("Post", a.Endpoint, data)
	return logs, nil
}

func (a *ApiClient) Get(response interface{}) (interface{}, error) {
	a.Start = time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("Endpoint : ", a.Endpoint)
	fmt.Println("Header : ", a.Headers)

	data, err := a.Resty.SetPreRequestHook(a.BeforeRequest).R().Get(a.Endpoint)

	if err != nil {
		logs := a.ApiLog("Get", a.Endpoint, data)
		return logs, status.Error(codes.Internal, err.Error())
	}

	fmt.Println("Res Header : ", data.RawResponse.Header)
	fmt.Println("Res Body : ", string(data.Body()))

	if data.RawResponse.StatusCode != 200 {
		logs := a.ApiLog("Get", a.Endpoint, data)
		return logs, status.Error(codes.Internal, data.Status())
	}

	var body = data.Body()
	if err = json.Unmarshal(body, response); err != nil {
		logs := a.ApiLog("Get", a.Endpoint, data)
		return logs, status.Error(codes.Internal, data.Status())
	}
	logs := a.ApiLog("Get", a.Endpoint, data)
	return logs, nil
}

func Stringify(data interface{}) string {
	dataByte, _ := json.Marshal(data)
	return string(dataByte)
}

func (a *ApiClient) BeforeRequest(r *resty.Client, h *http.Request) error {
	for k, v := range a.Headers {
		h.Header[k] = []string{v}
	}
	return nil
}
