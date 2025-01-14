package httpclient

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/acexy/golang-toolkit/logger"
	"github.com/acexy/golang-toolkit/math/random"
	"github.com/acexy/golang-toolkit/util/str"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
	"net/http"
	"os"
	"time"
)

// RestyClient resty客户端
type RestyClient struct {
	client *resty.Client

	// 是否设置了初始化代理，如果设置了不允许后续变更
	initProxy bool

	// 多代理模式
	multiProxy []string
	// 多代理选择器
	chooseProxy ChooseProxy
}

// RawRestyClient 获取原始restyClient实例
func (r *RestyClient) RawRestyClient() *resty.Client {
	return r.client
}

// RestyRequest resty请求对象
type RestyRequest struct {
	request *resty.Request
	client  *RestyClient
}

// RestyMethod resty请求方法对象
type RestyMethod struct {
	request *RestyRequest
	method  string
	url     string
}

// NewRestyClient 创建一个httpClient对象
// proxyHttpHost 可以指定代理 如 localhost:7890
func NewRestyClient(proxyHttpHost ...string) *RestyClient {
	var client = &RestyClient{}
	if len(proxyHttpHost) > 0 && str.HasText(proxyHttpHost[0]) {
		client.initProxy = true
		client.client = resty.New()
		client.client.SetProxy(proxyHttpHost[0])
	} else {
		client.client = resty.New()
	}
	client.client.SetLogger(logger.Logrus())
	return client
}

// NewRestyClientWithMultiProxy 创建一个多代理实例，该实例下的请求将通过策略通过代理
func NewRestyClientWithMultiProxy(multiProxy []string, choose ...ChooseProxy) *RestyClient {
	client := NewRestyClient()
	if len(multiProxy) > 0 {
		copied := make([]string, len(multiProxy))
		copy(copied, multiProxy)
		client.multiProxy = copied
		client.initProxy = true
	}
	if len(choose) > 0 {
		client.chooseProxy = choose[0]
	} else {
		client.chooseProxy = &randomChoose{}
	}
	client.client.SetLogger(logger.Logrus())
	return client
}

// ChooseProxy 多代理模式下的选择代理策略
type ChooseProxy interface {
	// Choose 选择代理
	Choose(all []string) string
}

type randomChoose struct {
}

func (r *randomChoose) Choose(all []string) string {
	return all[random.RandInt(len(all)-1)]
}

// client 公共属性设置

func (r *RestyClient) DisableTLSVerify() *RestyClient {
	r.client.SetTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 不验证证书签名
		},
	})
	return r
}

// SetBaseUrl 设置BaseUrl
func (r *RestyClient) SetBaseUrl(baseUrl string) *RestyClient {
	r.client.SetBaseURL(baseUrl)
	return r
}

// SetTimeout 设置超时时间
func (r *RestyClient) SetTimeout(timeout time.Duration) *RestyClient {
	r.client.SetTimeout(timeout)
	return r
}

// SetHeader 设置请求头
func (r *RestyClient) SetHeader(key, value string) *RestyClient {
	r.client.SetHeader(key, value)
	return r
}

// SetHeaders 设置请求头
func (r *RestyClient) SetHeaders(headers map[string]string) *RestyClient {
	r.client.SetHeaders(headers)
	return r
}

// SetProxy 设置代理
func (r *RestyClient) SetProxy(proxy string) *RestyClient {
	if r.initProxy {
		logger.Logrus().Warning("A global proxy is already set, operation ignored.")
		return r
	}
	r.client.SetProxy(proxy)
	return r
}

// R 获取Request实例
func (r *RestyClient) R() *RestyRequest {
	return &RestyRequest{request: r.client.R(), client: r}
}

// 对 restyRequest进行设置

// SetReturnStruct 使用默认响应Body内容与结构体绑定
// 仅支持响应码 200 - 299 内容类型为 JSON or XML时
func (r *RestyRequest) SetReturnStruct(any interface{}) *RestyRequest {
	r.request.SetResult(any)
	return r
}

// SetDownloadFile 将原始内容下载为文件
// filepath 文件完整路径(含文件名)
func (r *RestyRequest) SetDownloadFile(filepath string) *RestyRequest {
	outputFile, err := os.Create(filepath)
	if err != nil {
		logger.Logrus().WithError(err).Println("Failed to create output file", filepath)
	}
	defer outputFile.Close()
	r.request.SetOutput(filepath)
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
		request: r,
		method:  httpMethod,
		url:     url,
	}
}

func (m *RestyMethod) SetRawBody(raw func(raw *resty.Request)) *RestyMethod {
	raw(m.request.request)
	return m
}
func (m *RestyMethod) SetRequestBody(body interface{}, contentType string) *RestyMethod {
	m.request.request.SetBody(body)
	m.request.SetHeader(HeadContentType, contentType)
	return m
}

func (m *RestyMethod) SetQueryValues(any interface{}) *RestyMethod {
	queryParam, _ := query.Values(any)
	m.request.request.QueryParam = queryParam
	return m
}

func (m *RestyMethod) SetPathValues(pathParams map[string]string) *RestyMethod {
	m.request.request.PathParams = pathParams
	return m
}

func (m *RestyMethod) SetBodyJson(bodyJson string, charset ...string) *RestyMethod {
	m.SetRequestBody(bodyJson, getContentType(ContentTypeJson, charset...))
	return m
}

func (m *RestyMethod) SetBodyForm(formEncode map[string]string) *RestyMethod {
	m.request.request.SetFormData(formEncode)
	return m
}

// E Execution
func (m *RestyMethod) E() (*resty.Response, error) {
	setProxy(m.request)
	switch m.method {
	case http.MethodGet:
		return m.request.request.Get(m.url)
	case http.MethodPost:
		return m.request.request.Post(m.url)
	case http.MethodPut:
		return m.request.request.Put(m.url)
	case http.MethodDelete:
		return m.request.request.Delete(m.url)
	case http.MethodHead:
		return m.request.request.Head(m.url)
	case http.MethodOptions:
		return m.request.request.Options(m.url)
	case http.MethodPatch:
		return m.request.request.Patch(m.url)
	}
	return nil, errors.New("Unknown Method " + m.method)
}

func setProxy(r *RestyRequest) {
	if r.client.initProxy && len(r.client.multiProxy) > 0 {
		r.client.client.SetProxy(r.client.chooseProxy.Choose(r.client.multiProxy))
	}
}

// 常用的快捷请求方法，默认使用resty.R()

func (r *RestyRequest) Get(url string) (*resty.Response, error) {
	setProxy(r)
	return r.request.Get(url)
}

func (r *RestyRequest) Post(url string) (*resty.Response, error) {
	setProxy(r)
	return r.request.Post(url)
}

func (r *RestyRequest) PostForm(url string, formEncode map[string]string) (*resty.Response, error) {
	setProxy(r)
	if len(formEncode) == 0 {
		return r.request.Post(url)
	}
	return r.request.SetFormData(formEncode).Post(url)
}

func (r *RestyRequest) PostJson(url string, jsonString string, charset ...string) (*resty.Response, error) {
	setProxy(r)
	r.request.SetBody(jsonString)
	r.request.SetHeader(HeadContentType, getContentType(ContentTypeJson, charset...))
	return r.request.Post(url)
}
