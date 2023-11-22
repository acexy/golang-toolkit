package http

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
	"net/http"
	"net/url"
	"time"
)

type RestyClient struct {
	client *resty.Client
}

type restyRequest struct {
	request *resty.Request
}

type restyMethod struct {
	request *resty.Request
	method  string
	url     string
}

// NewRestyClient 创建一个httpClient对象
// proxyHttpHost 可以指定代理 localhost:7890
func NewRestyClient(proxyHttpHost ...string) *RestyClient {
	var client = &RestyClient{}
	if len(proxyHttpHost) > 0 {
		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: func(*http.Request) (*url.URL, error) {
					return &url.URL{Scheme: "http", Host: proxyHttpHost[0]}, nil
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
func (c *RestyClient) R() *restyRequest {
	return &restyRequest{request: c.client.R()}
}

// restyRequest 设置

func (r *restyRequest) SetReturnStruct(any interface{}) *restyRequest {
	r.request.SetResult(any)
	return r
}

func (r *restyRequest) SetQueryValues(any interface{}) *restyRequest {
	queryParam, _ := query.Values(any)
	r.request.QueryParam = queryParam
	return r
}

func (r *restyRequest) SetPathValues(pathParams map[string]string) *restyRequest {
	r.request.PathParams = pathParams
	return r
}

func (r *restyRequest) WithContext(ctx context.Context) *restyRequest {
	r.request.SetContext(ctx)
	return r
}

func (r *restyRequest) SetHeaders(headers map[string]string) *restyRequest {
	r.request.SetHeaders(headers)
	return r
}

func (r *restyRequest) SetHeader(key, value string) *restyRequest {
	r.request.SetHeader(key, value)
	return r
}

// M set Method
func (r *restyRequest) M(httpMethod string, url string) *restyMethod {
	return &restyMethod{
		request: r.request,
		method:  httpMethod,
		url:     url,
	}
}

func (m *restyMethod) SetRequestBody(bodyString *string, contentType ContentType) *restyMethod {
	m.request.SetBody(&bodyString)
	m.request.SetHeader(HeadContentType, string(contentType))
	return m
}

func (m *restyMethod) SetBodyJson(bodyJson *string) *restyMethod {
	m.SetRequestBody(bodyJson, ContentTypeJson)
	return m
}

func (m *restyMethod) SetBodyForm(formEncode map[string]string) *restyMethod {
	m.request.SetFormData(formEncode)
	return m
}

// E Execution
func (m *restyMethod) E() (*resty.Response, error) {
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
	return nil, nil
}

func (r *restyRequest) Get(url string) (*resty.Response, error) {
	return r.request.Get(url)
}

func (r *restyRequest) PostForm(url string, formEncode map[string]string) (*resty.Response, error) {
	if len(formEncode) == 0 {
		return r.request.Post(url)
	}
	return r.request.SetFormData(formEncode).Post(url)
}

func (r *restyRequest) PostJson(url string, jsonString string) (*resty.Response, error) {
	r.request.SetBody(jsonString)
	r.request.SetHeader(HeadContentType, string(ContentTypeJson))
	return r.request.Post(url)
}
