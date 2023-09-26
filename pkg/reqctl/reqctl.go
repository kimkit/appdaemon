package reqctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	urlib "net/url"
	"strings"
	"time"
)

func NewClient(timeout int) *http.Client {
	if timeout <= 0 {
		timeout = 1
	}
	return &http.Client{
		// Proxy: http.ProxyFromEnvironment,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(timeout) * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   time.Duration(timeout) * time.Second,
			ExpectContinueTimeout: time.Duration(timeout) * time.Second,
		},
		Timeout: time.Duration(timeout) * time.Second,
	}
}

type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Request    *http.Request
}

func buildResponse(resp *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       body,
		Request:    resp.Request,
	}, nil
}

func Get(client *http.Client, url string, data interface{}, header map[string]string) (*Response, error) {
	query := ""
	if data != nil {
		switch _data := data.(type) {
		case string:
			query = _data
		case []byte:
			query = string(_data)
		case map[string]string:
			values := urlib.Values{}
			for k, v := range _data {
				values.Set(k, v)
			}
			query = values.Encode()
		}
	}

	if len(query) > 0 {
		if strings.Contains(url, "?") {
			url = fmt.Sprintf("%s&%s", url, query)
		} else {
			url = fmt.Sprintf("%s?%s", url, query)
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return buildResponse(resp)
}

func Post(client *http.Client, url string, data interface{}, header map[string]string) (*Response, error) {
	var body []byte
	if data != nil {
		jsonReq := false
		for k, v := range header {
			if strings.ToLower(k) == "content-type" && v == "application/json" {
				jsonReq = true
				break
			}
		}
		if jsonReq {
			switch _data := data.(type) {
			case string:
				body = []byte(_data)
			case []byte:
				body = _data
			default:
				_body, err := json.Marshal(data)
				if err != nil {
					return nil, err
				}
				body = _body
			}
		} else {
			switch _data := data.(type) {
			case string:
				body = []byte(_data)
			case []byte:
				body = _data
			case map[string]string:
				values := urlib.Values{}
				for k, v := range _data {
					values.Set(k, v)
				}
				_body := values.Encode()
				body = []byte(_body)
			}
		}
	}

	// fmt.Println(string(body))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return buildResponse(resp)
}
