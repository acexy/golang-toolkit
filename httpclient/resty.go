package httpclient

import (
	"context"
	"errors"
	"github.com/acexy/golang-toolkit/util/str"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
	"net/http"
	"net/url"
	"time"
)

type RestyClient struct {
	client *resty.Client
}

type RestyRequest struct {
	request *resty.Request
}

type RestyMethod struct {
	request *resty.Request
	method  string
	url     string
}

// NewRestyClient 创建一个httpClient对象
// proxyHttpHost 可以指定代理 如 localhost:7890
func NewRestyClient(proxyHttpHost ...string) *RestyClient {
	var client = &RestyClient{}
	if len(proxyHttpHost) > 0 && str.HasText(proxyHttpHost[0]) {
		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: func(*http.Request) (*url.URL, error) {
					return &url.URL{Scheme: "httpclient", Host: proxyHttpHost[0]}, nil
				},
			},
		}
		client.client = resty.NewWithClient(httpClient)
	} else {
		client.client = resty.New()
	}
	return client
}

// client 设置

func (c *RestyClient) SetBaseUrl(baseUrl string) *RestyClient {
	c.client.SetBaseURL(baseUrl)
	return c
}

func (c *RestyClient) SetTimeout(timeout time.Duration) *RestyClient {
	c.client.SetTimeout(timeout)
	return c
}

func (c *RestyClient) SetHeader(key, value string) *RestyClient {
	c.client.SetHeader(key, value)
	return c
}

func (c *RestyClient) SetHeaders(headers map[string]string) *RestyClient {
	c.client.SetHeaders(headers)
	return c
}

// R GetRequest
func (c *RestyClient) R() *RestyRequest {
	return &RestyRequest{request: c.client.R()}
}

// restyRequest 设置

func (r *RestyRequest) SetReturnStruct(any interface{}) *RestyRequest {
	r.request.SetResult(any)
	return r
}

func (r *RestyRequest) SetQueryValues(any interface{}) *RestyRequest {
	queryParam, _ := query.Values(any)
	r.request.QueryParam = queryParam
	return r
}

func (r *RestyRequest) SetPathValues(pathParams map[string]string) *RestyRequest {
	r.request.PathParams = pathParams
	return r
}

func (r *RestyRequest) WithContext(ctx context.Context) *RestyRequest {
	r.request.SetContext(ctx)
	return r
}

func (r *RestyRequest) SetHeaders(headers map[string]string) *RestyRequest {
	r.request.SetHeaders(headers)
	return r
}

func (r *RestyRequest) SetHeader(key, value string) *RestyRequest {
	r.request.SetHeader(key, value)
	return r
}

// M set Method
func (r *RestyRequest) M(httpMethod string, url string) *RestyMethod {
	return &RestyMethod{
		request: r.request,
		method:  httpMethod,
		url:     url,
	}
}

func (m *RestyMethod) SetRequestBody(bodyString *string, contentType ContentType) *RestyMethod {
	m.request.SetBody(&bodyString)
	m.request.SetHeader(HeadContentType, string(contentType))
	return m
}

func (m *RestyMethod) SetBodyJson(bodyJson *string) *RestyMethod {
	m.SetRequestBody(bodyJson, ContentTypeJson)
	return m
}

func (m *RestyMethod) SetBodyForm(formEncode map[string]string) *RestyMethod {
	m.request.SetFormData(formEncode)
	return m
}

// E Execution
func (m *RestyMethod) E() (*resty.Response, error) {
	switch m.method {
	case http.MethodGet:
		return m.request.Get(m.url)
	case http.MethodPost:
		return m.request.Post(m.url)
	case http.MethodPut:
		return m.request.Put(m.url)
	case http.MethodDelete:
		return m.request.Delete(m.url)
	case http.MethodHead:
		return m.request.Head(m.url)
	case http.MethodOptions:
		return m.request.Options(m.url)
	case http.MethodPatch:
		return m.request.Patch(m.url)
	}
	return nil, errors.New("Unknown Method " + m.method)
}

func (r *RestyRequest) Get(url string) (*resty.Response, error) {
	return r.request.Get(url)
}

func (r *RestyRequest) Post(url string) (*resty.Response, error) {
	return r.request.Post(url)
}

func (r *RestyRequest) PostForm(url string, formEncode map[string]string) (*resty.Response, error) {
	if len(formEncode) == 0 {
		return r.request.Post(url)
	}
	return r.request.SetFormData(formEncode).Post(url)
}

func (r *RestyRequest) PostJson(url string, jsonString string) (*resty.Response, error) {
	r.request.SetBody(jsonString)
	r.request.SetHeader(HeadContentType, string(ContentTypeJson))
	return r.request.Post(url)
}
