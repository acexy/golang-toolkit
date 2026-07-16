package httpclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	toolkitError "github.com/acexy/golang-toolkit/error"
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
		logger.Logrus().Warningln("multiProxies must contain at least 2 proxies")
		return NewRestyClient(multiProxy...)
	}
	client := &RestyClient{
		r: resty.New(),
	}
	if err := client.ConfigureProxies(multiProxy, choose...); err != nil {
		logger.Logrus().Warningln(err)
	}
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
	if len(all) == 0 {
		return ""
	}
	return all[random.RandInt(len(all)-1)]
}

// r 公共属性设置

// SetProxies 设置代理池
func (r *RestyClient) SetProxies(proxyUrls []string, choose ...ChooseProxy) {
	if err := r.ConfigureProxies(proxyUrls, choose...); err != nil {
		logger.Logrus().Warningln(err)
	}
}

// ConfigureProxies 设置代理池，并向调用方返回配置错误。
func (r *RestyClient) ConfigureProxies(proxyUrls []string, choose ...ChooseProxy) error {
	if len(proxyUrls) == 0 {
		return nil
	}

	// 复制并过滤代理列表，避免调用方后续修改切片影响并发请求。
	proxies := append([]string(nil), proxyUrls...)
	proxyCache := coll.SliceFilterToMap(proxies, func(proxy string) (string, *url.URL, bool) {
		pURL, err := url.Parse(proxy)
		if err != nil {
			logger.Logrus().Errorln("parse proxy url error", proxy, err)
			return "", nil, false
		}
		if pURL.Scheme == "" || pURL.Host == "" {
			logger.Logrus().Errorln("parse proxy url error", proxy, "scheme or host is empty")
			return "", nil, false
		}
		return proxy, pURL, true
	})
	validProxies := make([]string, 0, len(proxyCache))
	for _, proxy := range proxies {
		if _, ok := proxyCache[proxy]; ok {
			validProxies = append(validProxies, proxy)
		}
	}
	if len(validProxies) == 0 {
		return fmt.Errorf("configure proxies: no valid proxy URL")
	}

	// SetProxy 用于初始化并标记 Resty 的代理 Transport。
	transport, err := r.r.SetProxy(validProxies[0]).Transport()
	if err != nil {
		return fmt.Errorf("configure proxies: %w", err)
	}
	var chooseFn ChooseProxy
	if len(choose) > 0 {
		chooseFn = choose[0]
	} else {
		chooseFn = &randomChoose{}
	}
	transport.Proxy = func(request *http.Request) (*url.URL, error) {
		proxyURL := chooseFn.Choose(request, validProxies)
		parsedProxy, ok := proxyCache[proxyURL]
		if !ok {
			return nil, fmt.Errorf("choose proxy: proxy %q is not in the valid proxy pool", proxyURL)
		}
		return parsedProxy, nil
	}
	return nil
}

// SetProxy 设置代理
func (r *RestyClient) SetProxy(proxy string) *RestyClient {
	r.r.SetProxy(proxy)
	return r
}

// DisableTLSVerify 禁用TLS验证
func (r *RestyClient) DisableTLSVerify() *RestyClient {
	r.r.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true, // 不验证证书签名
	})
	return r
}

// DisableAllAutoRedirect 禁用所有自动重定向
func (r *RestyClient) DisableAllAutoRedirect() *RestyClient {
	r.r.SetRedirectPolicy(resty.NoRedirectPolicy())
	return r
}

// SetBaseURL 设置BaseURL
func (r *RestyClient) SetBaseURL(baseURL string) *RestyClient {
	r.r.SetBaseURL(baseURL)
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
	return &RestyRequest{request: r.r.R()}
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
	r.request.SetOutput(filepath)
	return r
}

// WithContext 设置请求上下文。
func (r *RestyRequest) WithContext(ctx context.Context) *RestyRequest {
	r.request.SetContext(ctx)
	return r
}

// SetHeaders 批量设置请求头。
func (r *RestyRequest) SetHeaders(headers map[string]string) *RestyRequest {
	r.request.SetHeaders(headers)
	return r
}

// SetHeader 设置请求头。
func (r *RestyRequest) SetHeader(key, value string) *RestyRequest {
	r.request.SetHeader(key, value)
	return r
}

// Method 设置请求方法
func (r *RestyRequest) Method(httpMethod string, url string) *RestyMethod {
	return &RestyMethod{
		request: r,
		method:  httpMethod,
		url:     url,
	}
}

// ConfigureRequest 使用回调配置原始 Resty 请求。
func (m *RestyMethod) ConfigureRequest(configure func(request *resty.Request)) *RestyMethod {
	if configure != nil {
		configure(m.request.request)
	}
	return m
}

// SetRequestBody 设置请求体及其内容类型。
func (m *RestyMethod) SetRequestBody(body interface{}, contentType string) *RestyMethod {
	m.request.request.SetBody(body)
	m.request.SetHeader(HeaderContentType, contentType)
	return m
}

// SetQueryValues 设置 URL 查询参数。
func (m *RestyMethod) SetQueryValues(query url.Values) *RestyMethod {
	m.request.request.SetQueryParamsFromValues(query)
	return m
}

// SetPathValues 设置 URL 路径参数。
func (m *RestyMethod) SetPathValues(pathParams map[string]string) *RestyMethod {
	m.request.request.SetPathParams(pathParams)
	return m
}

// SetBodyJSON 设置 JSON 字符串请求体。
func (m *RestyMethod) SetBodyJSON(bodyJSON string, charset ...string) *RestyMethod {
	m.SetRequestBody(bodyJSON, getContentType(ContentTypeJSON, charset...))
	return m
}

// SetBodyForm 设置表单请求体。
func (m *RestyMethod) SetBodyForm(formEncode map[string]string) *RestyMethod {
	m.request.request.SetFormData(formEncode)
	return m
}

// Execute 执行请求
func (m *RestyMethod) Execute() (*resty.Response, error) {
	switch m.method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete,
		http.MethodHead, http.MethodOptions, http.MethodPatch:
		return m.request.request.Execute(m.method, m.url)
	}
	return nil, toolkitError.ErrUnsupportedHTTPMethod
}

// 常用的快捷请求方法，默认使用resty.R()

// Get 发起 GET 请求。
func (r *RestyRequest) Get(url string) (*resty.Response, error) {
	return r.request.Get(url)
}

// Post 发起 POST 请求。
func (r *RestyRequest) Post(url string) (*resty.Response, error) {
	return r.request.Post(url)
}

// PostForm 发起表单 POST 请求。
func (r *RestyRequest) PostForm(url string, formEncode map[string]string) (*resty.Response, error) {
	if len(formEncode) == 0 {
		return r.request.Post(url)
	}
	return r.request.SetFormData(formEncode).Post(url)
}

// PostJSON 发起 JSON POST 请求。
func (r *RestyRequest) PostJSON(url string, jsonString string, charset ...string) (*resty.Response, error) {
	r.request.SetBody(jsonString)
	r.request.SetHeader(HeaderContentType, getContentType(ContentTypeJSON, charset...))
	return r.request.Post(url)
}
