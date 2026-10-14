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

// RestyClient wraps Resty client settings and proxy configuration.
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

// multiProxyTransport selects a dedicated proxy connection pool for each request.
type multiProxyTransport struct {
	template   *http.Transport
	proxies    []string
	transports map[string]*http.Transport
	choose     ChooseProxy
}

// RoundTrip selects a proxy and uses its dedicated connection pool.
func (m *multiProxyTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	proxy := m.choose.Choose(request, m.proxies)
	transport, ok := m.transports[proxy]
	if !ok {
		return nil, fmt.Errorf("choose proxy: proxy %q is not in the valid proxy pool", proxy)
	}
	// Record the selected proxy on this request rather than shared client state.
	if trace, ok := request.Context().Value(proxyContextKey{}).(*proxyTrace); ok {
		trace.setProxy(proxy)
	}
	return transport.RoundTrip(request)
}

// CloseIdleConnections closes idle connections in every proxy pool.
func (m *multiProxyTransport) CloseIdleConnections() {
	for _, transport := range m.transports {
		transport.CloseIdleConnections()
	}
}

// RawRestyClient returns the underlying Resty client.
func (r *RestyClient) RawRestyClient() *resty.Client {
	return r.r
}

// RestyRequest wraps a request and its request-scoped proxy trace.
type RestyRequest struct {
	request *resty.Request
	proxy   string
	trace   *proxyTrace
}

// Response wraps a response and the proxy used for that request.
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

// RawRestyResponse returns the underlying Resty response.
func (r *Response) RawRestyResponse() *resty.Response {
	return r.response
}

// RawResponse returns the underlying HTTP response.
func (r *Response) RawResponse() *http.Response {
	return r.response.RawResponse
}

// Proxy returns the selected proxy, or an empty string for a direct request.
func (r *Response) Proxy() string {
	if r == nil {
		return ""
	}
	return r.proxy
}

// Body returns the buffered response body.
func (r *Response) Body() []byte {
	return r.response.Body()
}

// RawBody returns the underlying body stream when raw response handling is enabled.
func (r *Response) RawBody() io.ReadCloser {
	return r.response.RawBody()
}

// String returns the buffered response body as a string.
func (r *Response) String() string {
	return r.response.String()
}

// Status returns the HTTP status text.
func (r *Response) Status() string {
	return r.response.Status()
}

// StatusCode returns the HTTP status code.
func (r *Response) StatusCode() int {
	return r.response.StatusCode()
}

// Proto returns the HTTP protocol version.
func (r *Response) Proto() string {
	return r.response.Proto()
}

// Header returns response headers.
func (r *Response) Header() http.Header {
	return r.response.Header()
}

// Cookies returns response cookies.
func (r *Response) Cookies() []*http.Cookie {
	return r.response.Cookies()
}

// Result returns the bound successful-response result.
func (r *Response) Result() any {
	return r.response.Result()
}

// Error returns the bound error-response result.
func (r *Response) Error() any {
	return r.response.Error()
}

// IsSuccess reports whether the status is between 200 and 299.
func (r *Response) IsSuccess() bool {
	return r.response.IsSuccess()
}

// IsError reports whether the status is at least 400.
func (r *Response) IsError() bool {
	return r.response.IsError()
}

// Time returns the request duration.
func (r *Response) Time() time.Duration {
	return r.response.Time()
}

// ReceivedAt returns the response receipt time.
func (r *Response) ReceivedAt() time.Time {
	return r.response.ReceivedAt()
}

// Size returns the response body size in bytes.
func (r *Response) Size() int64 {
	return r.response.Size()
}

// RestyMethod holds a request method, URL, and body configuration.
type RestyMethod struct {
	request *RestyRequest
	method  string
	url     string
}

// NewRestyClient creates a client with an optional proxy URL.
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

// NewRestyClientWithMultiProxy creates a proxy-pool client with random selection by default; configuration errors are logged.
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

// ChooseProxy selects a configured proxy for each transport request.
type ChooseProxy interface {
	// Choose returns a proxy from the supplied pool.
	Choose(request *http.Request, all []string) string
}

type randomChoose struct {
}

// Choose randomly selects a proxy, returning empty for an empty pool.
func (r *randomChoose) Choose(_ *http.Request, all []string) string {
	if len(all) == 0 {
		return ""
	}
	return all[random.RandInt(len(all)-1)]
}

// Client configuration.

// SetProxies configures a proxy pool and logs configuration errors.
func (r *RestyClient) SetProxies(proxyUrls []string, choose ...ChooseProxy) {
	if err := r.ConfigureProxies(proxyUrls, choose...); err != nil {
		logger.Logrus().Warningln(err)
	}
}

// ConfigureProxies configures dedicated proxy transports and returns configuration errors.
func (r *RestyClient) ConfigureProxies(proxyUrls []string, choose ...ChooseProxy) error {
	if len(proxyUrls) == 0 {
		return nil
	}

	// Copy and validate the pool so caller slice mutations cannot affect requests.
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

// SetProxy switches to a single proxy and closes old proxy-pool idle connections.
func (r *RestyClient) SetProxy(proxy string) *RestyClient {
	if multiTransport, ok := r.r.GetClient().Transport.(*multiProxyTransport); ok {
		r.r.SetTransport(multiTransport.template.Clone())
		multiTransport.CloseIdleConnections()
	}
	r.r.SetProxy(proxy)
	r.proxy = proxy
	return r
}

// DisableTLSVerify disables server certificate verification.
func (r *RestyClient) DisableTLSVerify() *RestyClient {
	if multiTransport, ok := r.r.GetClient().Transport.(*multiProxyTransport); ok {
		multiTransport.template.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true, // Disable server certificate verification.
		}
		for _, transport := range multiTransport.transports {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true, // Disable server certificate verification.
			}
		}
		return r
	}
	r.r.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true, // Disable server certificate verification.
	})
	return r
}

// DisableAllAutoRedirect disables automatic HTTP redirects.
func (r *RestyClient) DisableAllAutoRedirect() *RestyClient {
	r.r.SetRedirectPolicy(resty.NoRedirectPolicy())
	return r
}

// SetBaseURL sets the base URL for relative requests.
func (r *RestyClient) SetBaseURL(baseURL string) *RestyClient {
	r.r.SetBaseURL(baseURL)
	return r
}

// SetTimeout sets the client request timeout.
func (r *RestyClient) SetTimeout(timeout time.Duration) *RestyClient {
	r.r.SetTimeout(timeout)
	return r
}

// SetHeader sets a client-wide header.
func (r *RestyClient) SetHeader(key, value string) *RestyClient {
	r.r.SetHeader(key, value)
	return r
}

// SetHeaders sets client-wide headers.
func (r *RestyClient) SetHeaders(headers map[string]string) *RestyClient {
	r.r.SetHeaders(headers)
	return r
}

// R creates a request with its own proxy trace.
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

// Request configuration.

// SetReturnStruct binds successful JSON or XML response data to result.
func (r *RestyRequest) SetReturnStruct(result any) *RestyRequest {
	r.request.SetResult(result)
	return r
}

// SetDownloadFile writes response data to the supplied file path.
func (r *RestyRequest) SetDownloadFile(filepath string) *RestyRequest {
	r.request.SetOutput(filepath)
	return r
}

// WithContext sets the context while preserving the request proxy trace.
func (r *RestyRequest) WithContext(ctx context.Context) *RestyRequest {
	r.request.SetContext(context.WithValue(ctx, proxyContextKey{}, r.trace))
	return r
}

// SetHeaders sets request headers.
func (r *RestyRequest) SetHeaders(headers map[string]string) *RestyRequest {
	r.request.SetHeaders(headers)
	return r
}

// SetHeader sets a request header.
func (r *RestyRequest) SetHeader(key, value string) *RestyRequest {
	r.request.SetHeader(key, value)
	return r
}

// Method creates a builder for the method and URL.
func (r *RestyRequest) Method(httpMethod string, url string) *RestyMethod {
	return &RestyMethod{
		request: r,
		method:  httpMethod,
		url:     url,
	}
}

// ConfigureRequest applies a nonnil callback to the underlying Resty request.
func (m *RestyMethod) ConfigureRequest(configure func(request *resty.Request)) *RestyMethod {
	if configure != nil {
		configure(m.request.request)
	}
	return m
}

// SetRequestBody sets the body and Content-Type header.
func (m *RestyMethod) SetRequestBody(body any, contentType string) *RestyMethod {
	m.request.request.SetBody(body)
	m.request.SetHeader(HeaderContentType, contentType)
	return m
}

// SetQueryValues sets URL query parameters.
func (m *RestyMethod) SetQueryValues(query url.Values) *RestyMethod {
	m.request.request.SetQueryParamsFromValues(query)
	return m
}

// SetPathValues sets URL path parameters.
func (m *RestyMethod) SetPathValues(pathParams map[string]string) *RestyMethod {
	m.request.request.SetPathParams(pathParams)
	return m
}

// SetBodyJSON sets a JSON string body with an optional charset.
func (m *RestyMethod) SetBodyJSON(bodyJSON string, charset ...string) *RestyMethod {
	m.SetRequestBody(bodyJSON, getContentType(ContentTypeJSON, charset...))
	return m
}

// SetBodyForm sets a URL-encoded form body.
func (m *RestyMethod) SetBodyForm(formEncode map[string]string) *RestyMethod {
	m.request.request.SetFormData(formEncode)
	return m
}

// Execute executes a supported HTTP method or returns ErrUnsupportedHTTPMethod.
func (m *RestyMethod) Execute() (*Response, error) {
	switch m.method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete,
		http.MethodHead, http.MethodOptions, http.MethodPatch:
		response, err := m.request.request.Execute(m.method, m.url)
		return wrapResponse(response, m.request.proxy, m.request.trace), err
	}
	return nil, toolkitError.ErrUnsupportedHTTPMethod
}

// Convenience request methods.

// Get executes a GET request.
func (r *RestyRequest) Get(url string) (*Response, error) {
	response, err := r.request.Get(url)
	return wrapResponse(response, r.proxy, r.trace), err
}

// Post executes a POST request.
func (r *RestyRequest) Post(url string) (*Response, error) {
	response, err := r.request.Post(url)
	return wrapResponse(response, r.proxy, r.trace), err
}

// PostForm executes a POST with optional form fields.
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

// PostJSON executes a POST with a JSON string body and optional charset.
func (r *RestyRequest) PostJSON(url string, jsonString string, charset ...string) (*Response, error) {
	r.request.SetBody(jsonString)
	r.request.SetHeader(HeaderContentType, getContentType(ContentTypeJSON, charset...))
	response, err := r.request.Post(url)
	return wrapResponse(response, r.proxy, r.trace), err
}
