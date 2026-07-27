package httpclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	toolkitError "github.com/acexy/golang-toolkit/error"
	"github.com/acexy/golang-toolkit/logger"
	"github.com/acexy/golang-toolkit/math/random"
	"github.com/acexy/golang-toolkit/util/coll"
	"github.com/go-resty/resty/v2"
)

// RestyClient resty客户端
type RestyClient struct {
	r     *resty.Client
	proxy string
}

type proxyContextKey struct{}

type proxyTrace struct {
	mu    sync.RWMutex
	proxy string
}

func (p *proxyTrace) setProxy(proxy string) {
	p.mu.Lock()
	p.proxy = proxy
	p.mu.Unlock()
}

func (p *proxyTrace) proxyURL() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.proxy
}

// multiProxyTransport 为每个代理维护独立的连接池，并在每次请求前执行代理选择策略。
type multiProxyTransport struct {
	template   *http.Transport
	proxies    []string
	transports map[string]*http.Transport
	choose     ChooseProxy
}

func (m *multiProxyTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	proxy := m.choose.Choose(request, m.proxies)
	transport, ok := m.transports[proxy]
	if !ok {
		return nil, fmt.Errorf("choose proxy: proxy %q is not in the valid proxy pool", proxy)
	}
	// 将实际选中的代理绑定到本次请求，确保并发请求之间互不影响。
	if trace, ok := request.Context().Value(proxyContextKey{}).(*proxyTrace); ok {
		trace.setProxy(proxy)
	}
	return transport.RoundTrip(request)
}

// CloseIdleConnections 关闭所有代理 Transport 的空闲连接。
func (m *multiProxyTransport) CloseIdleConnections() {
	for _, transport := range m.transports {
		transport.CloseIdleConnections()
	}
}

// RawRestyClient 获取原始restyClient实例
func (r *RestyClient) RawRestyClient() *resty.Client {
	return r.r
}

// RestyRequest resty请求对象
type RestyRequest struct {
	request *resty.Request
	proxy   string
	trace   *proxyTrace
}

// Response HTTP响应对象
type Response struct {
	response *resty.Response
	proxy    string
}

func wrapResponse(response *resty.Response, configuredProxy string, trace *proxyTrace) *Response {
	if response == nil {
		return nil
	}
	actualProxy := configuredProxy
	if trace != nil {
		if proxy := trace.proxyURL(); proxy != "" {
			actualProxy = proxy
		}
	}
	return &Response{
		response: response,
		proxy:    actualProxy,
	}
}

// RawRestyResponse 获取原始 Resty 响应对象。
func (r *Response) RawRestyResponse() *resty.Response {
	return r.response
}

// RawResponse 获取原始 HTTP 响应对象。
func (r *Response) RawResponse() *http.Response {
	return r.response.RawResponse
}

// Proxy 获取本次请求实际使用的代理，未使用代理时返回空字符串。
func (r *Response) Proxy() string {
	if r == nil {
		return ""
	}
	return r.proxy
}

// Body 获取响应体。
func (r *Response) Body() []byte {
	return r.response.Body()
}

// RawBody 获取未读取的原始响应体。
func (r *Response) RawBody() io.ReadCloser {
	return r.response.RawBody()
}

// String 以字符串形式获取响应体。
func (r *Response) String() string {
	return r.response.String()
}

// Status 获取响应状态。
func (r *Response) Status() string {
	return r.response.Status()
}

// StatusCode 获取响应状态码。
func (r *Response) StatusCode() int {
	return r.response.StatusCode()
}

// Proto 获取响应协议版本。
func (r *Response) Proto() string {
	return r.response.Proto()
}

// Header 获取响应头。
func (r *Response) Header() http.Header {
	return r.response.Header()
}

// Cookies 获取响应 Cookie。
func (r *Response) Cookies() []*http.Cookie {
	return r.response.Cookies()
}

// Result 获取成功响应绑定的结果。
func (r *Response) Result() interface{} {
	return r.response.Result()
}

// Error 获取错误响应绑定的结果。
func (r *Response) Error() interface{} {
	return r.response.Error()
}

// IsSuccess 判断响应状态码是否为 200 至 299。
func (r *Response) IsSuccess() bool {
	return r.response.IsSuccess()
}

// IsError 判断响应状态码是否大于等于 400。
func (r *Response) IsError() bool {
	return r.response.IsError()
}

// Time 获取请求耗时。
func (r *Response) Time() time.Duration {
	return r.response.Time()
}

// ReceivedAt 获取收到响应的时间。
func (r *Response) ReceivedAt() time.Time {
	return r.response.ReceivedAt()
}

// Size 获取响应体大小。
func (r *Response) Size() int64 {
	return r.response.Size()
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
// 每个代理使用独立的 Transport，所有请求都会触发 chooseProxy
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

	var chooseFn ChooseProxy
	if len(choose) > 0 {
		chooseFn = choose[0]
	} else {
		chooseFn = &randomChoose{}
	}

	var template *http.Transport
	currentTransport := r.r.GetClient().Transport
	if multiTransport, ok := currentTransport.(*multiProxyTransport); ok {
		template = multiTransport.template.Clone()
	} else {
		transport, ok := currentTransport.(*http.Transport)
		if !ok {
			return fmt.Errorf("configure proxies: unsupported transport type %T", currentTransport)
		}
		template = transport.Clone()
	}
	template.Proxy = nil

	transports := make(map[string]*http.Transport, len(proxyCache))
	for proxy, proxyURL := range proxyCache {
		transport := template.Clone()
		transport.Proxy = http.ProxyURL(proxyURL)
		transports[proxy] = transport
	}

	multiTransport := &multiProxyTransport{
		template:   template,
		proxies:    validProxies,
		transports: transports,
		choose:     chooseFn,
	}
	r.r.SetTransport(multiTransport)
	r.proxy = ""
	if previous, ok := currentTransport.(*multiProxyTransport); ok {
		previous.CloseIdleConnections()
	}
	return nil
}

// SetProxy 设置代理
func (r *RestyClient) SetProxy(proxy string) *RestyClient {
	if multiTransport, ok := r.r.GetClient().Transport.(*multiProxyTransport); ok {
		r.r.SetTransport(multiTransport.template.Clone())
		multiTransport.CloseIdleConnections()
	}
	r.r.SetProxy(proxy)
	r.proxy = proxy
	return r
}

// DisableTLSVerify 禁用TLS验证
func (r *RestyClient) DisableTLSVerify() *RestyClient {
	if multiTransport, ok := r.r.GetClient().Transport.(*multiProxyTransport); ok {
		multiTransport.template.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true, // 不验证证书签名
		}
		for _, transport := range multiTransport.transports {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true, // 不验证证书签名
			}
		}
		return r
	}
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
	trace := &proxyTrace{}
	request := r.r.R()
	request.SetContext(context.WithValue(context.Background(), proxyContextKey{}, trace))
	return &RestyRequest{
		request: request,
		proxy:   r.proxy,
		trace:   trace,
	}
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
	r.request.SetContext(context.WithValue(ctx, proxyContextKey{}, r.trace))
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
func (m *RestyMethod) Execute() (*Response, error) {
	switch m.method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete,
		http.MethodHead, http.MethodOptions, http.MethodPatch:
		response, err := m.request.request.Execute(m.method, m.url)
		return wrapResponse(response, m.request.proxy, m.request.trace), err
	}
	return nil, toolkitError.ErrUnsupportedHTTPMethod
}

// 常用的快捷请求方法，默认使用resty.R()

// Get 发起 GET 请求。
func (r *RestyRequest) Get(url string) (*Response, error) {
	response, err := r.request.Get(url)
	return wrapResponse(response, r.proxy, r.trace), err
}

// Post 发起 POST 请求。
func (r *RestyRequest) Post(url string) (*Response, error) {
	response, err := r.request.Post(url)
	return wrapResponse(response, r.proxy, r.trace), err
}

// PostForm 发起表单 POST 请求。
func (r *RestyRequest) PostForm(url string, formEncode map[string]string) (*Response, error) {
	var response *resty.Response
	var err error
	if len(formEncode) == 0 {
		response, err = r.request.Post(url)
	} else {
		response, err = r.request.SetFormData(formEncode).Post(url)
	}
	return wrapResponse(response, r.proxy, r.trace), err
}

// PostJSON 发起 JSON POST 请求。
func (r *RestyRequest) PostJSON(url string, jsonString string, charset ...string) (*Response, error) {
	r.request.SetBody(jsonString)
	r.request.SetHeader(HeaderContentType, getContentType(ContentTypeJSON, charset...))
	response, err := r.request.Post(url)
	return wrapResponse(response, r.proxy, r.trace), err
}
