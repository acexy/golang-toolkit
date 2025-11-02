package httpclient

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/acexy/golang-toolkit/logger"
	"github.com/acexy/golang-toolkit/math/random"
	"github.com/acexy/golang-toolkit/util/coll"
	"github.com/go-resty/resty/v2"
)

// RestyClient resty客户端
type RestyClient struct {
	r *resty.Client
}

// RawRestyClient 获取原始restyClient实例
func (r *RestyClient) RawRestyClient() *resty.Client {
	return r.r
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
// proxyHttpHost 可以指定代理 如 http://localhost:7890
func NewRestyClient(proxyHttpHost ...string) *RestyClient {
	var client = &RestyClient{
		r: resty.New(),
	}
	client.r.SetLogger(logger.Logrus())
	if len(proxyHttpHost) > 0 {
		client.SetProxy(proxyHttpHost[0])
	}
	return client
}

// NewRestyClientWithMultiProxy 创建一个多代理实例，该实例下的请求将通过策略通过代理
// chooseProxy 可以指定选择代理的策略 默认为随机
// 受连接复用的影响，并不是每个请求都会触发chooseProxy
func NewRestyClientWithMultiProxy(multiProxy []string, choose ...ChooseProxy) *RestyClient {
	if len(multiProxy) < 2 {
		logger.Logrus().Warningln("multiProxies must be greater than 2")
		return nil
	}
	client := &RestyClient{
		r: resty.New(),
	}
	client.SetProxies(multiProxy, choose...)
	client.r.SetLogger(logger.Logrus())
	return client
}

// ChooseProxy 多代理模式下的选择代理策略
type ChooseProxy interface {
	// Choose 选择代理
	Choose(request *http.Request, all []string) string
}

type randomChoose struct {
}

func (r *randomChoose) Choose(_ *http.Request, all []string) string {
	return all[random.RandInt(len(all)-1)]
}

// r 公共属性设置

// SetProxies 设置代理池
func (r *RestyClient) SetProxies(proxyUrls []string, choose ...ChooseProxy) {
	// 为原始client设置已设置代理的标识
	transport, err := r.r.SetProxy(proxyUrls[0]).Transport()
	if err != nil {
		logger.Logrus().Warningln(err)
		return
	}
	var proxyCache map[string]*url.URL // 转换代理地址
	proxyCache = coll.SliceFilterToMap(proxyUrls, func(proxy string) (string, *url.URL, bool) {
		pURL, err := url.Parse(proxy)
		if err != nil {
			logger.Logrus().Errorln("parse proxy url error", proxy, err)
			return "", nil, false
		}
		return proxy, pURL, true
	})
	var chooseFn ChooseProxy
	if len(choose) > 0 {
		chooseFn = choose[0]
	} else {
		chooseFn = &randomChoose{}
	}
	transport.Proxy = func(request *http.Request) (*url.URL, error) {
		proxyUrl := chooseFn.Choose(request, proxyUrls)
		return proxyCache[proxyUrl], nil
	}
}

// SetProxy 设置代理
func (r *RestyClient) SetProxy(proxy string) *RestyClient {
	r.r.SetProxy(proxy)
	return r
}

// DisableTLSVerify 禁用TLS验证
func (r *RestyClient) DisableTLSVerify() *RestyClient {
	r.r.SetTransport(&http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 不验证证书签名
		},
	})
	return r
}

// DisableAllAutoRedirect 禁用所有自动重定向
func (r *RestyClient) DisableAllAutoRedirect() *RestyClient {
	r.r.SetRedirectPolicy(resty.NoRedirectPolicy())
	return r
}

// SetBaseUrl 设置BaseUrl
func (r *RestyClient) SetBaseUrl(baseUrl string) *RestyClient {
	r.r.SetBaseURL(baseUrl)
	return r
}

// SetTimeout 设置超时时间
func (r *RestyClient) SetTimeout(timeout time.Duration) *RestyClient {
	r.r.SetTimeout(timeout)
	return r
}

// SetHeader 设置请求头
func (r *RestyClient) SetHeader(key, value string) *RestyClient {
	r.r.SetHeader(key, value)
	return r
}

// SetHeaders 设置请求头
func (r *RestyClient) SetHeaders(headers map[string]string) *RestyClient {
	r.r.SetHeaders(headers)
	return r
}

// R 获取Request实例
func (r *RestyClient) R() *RestyRequest {
	return &RestyRequest{request: r.r.R(), client: r}
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
	defer func() {
		if outputFile == nil {
			return
		}
		_ = outputFile.Close()
	}()
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

func (m *RestyMethod) SetQueryValues(query url.Values) *RestyMethod {
	m.request.request.QueryParam = query
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

// 常用的快捷请求方法，默认使用resty.R()

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

func (r *RestyRequest) PostJson(url string, jsonString string, charset ...string) (*resty.Response, error) {
	r.request.SetBody(jsonString)
	r.request.SetHeader(HeadContentType, getContentType(ContentTypeJson, charset...))
	return r.request.Post(url)
}
