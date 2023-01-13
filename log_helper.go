package api_client

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strings"
)

type ApiLog struct {
	Type            string
	URL             string
	Method          string
	RequestHeaders  string
	RequestBody     string
	ResponseHeaders string
	ResponseBody    string
	HTTPCode        string
	Time            string
	ElapsedTimeMs   string
}

func (l *ApiClient) ApiLog(method, client string, resp interface{}) *ApiLog {
	var (
		reqHeader, reqBody, resHeader, resBody, httpCode string
	)

	reqHeader = JMarshal(l.Headers)
	reqBody = JMarshal(l.Body)

	response, ok := resp.(*resty.Response)
	if ok {
		if resp != nil {
			if response.RawResponse != nil {
				resHeader = JMarshal(response.RawResponse.Header)
				httpCode = fmt.Sprint(response.StatusCode())
				resBody = string(response.Body())
			}
		}
	}
	out := &ApiLog{
		Type:            strings.ToUpper(client),
		URL:             l.Endpoint,
		Method:          method,
		RequestHeaders:  reqHeader,
		RequestBody:     reqBody,
		ResponseHeaders: resHeader,
		ResponseBody:    resBody,
		HTTPCode:        httpCode,
		Time:            l.Start,
		ElapsedTimeMs:   l.End,
	}

	return out
}

func JMarshal(data interface{}) string {
	dataByte, _ := json.Marshal(data)
	return string(dataByte)
}
