package httpclient

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

type roundRobinChooseProxy struct {
	count int
}

func (c *roundRobinChooseProxy) Choose(_ *http.Request, all []string) string {
	proxy := all[c.count%len(all)]
	c.count++
	return proxy
}

func TestMethodExecute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	response, err := NewRestyClient().R().Method(http.MethodGet, server.URL).Execute()
	if err != nil {
		t.Fatal(err)
	}
	if response.String() != "ok" {
		t.Fatalf("unexpected response: %s", response.String())
	}
	if response.Proxy() != "" {
		t.Fatalf("expected empty proxy, got %q", response.Proxy())
	}
	if response.RawRestyResponse() == nil || response.RawResponse() == nil {
		t.Fatal("expected raw response")
	}
}

func TestMethodExecuteUnsupportedMethod(t *testing.T) {
	_, err := NewRestyClient().R().Method("BAD", "http://example.com").Execute()
	if !errors.Is(err, toolkitError.ErrUnsupportedHTTPMethod) {
		t.Fatalf("expected ErrUnsupportedHTTPMethod, got %v", err)
	}
}

func TestSetQueryValuesAndPathValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/100" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("status") != "active" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	query := url.Values{}
	query.Set("status", "active")
	response, err := NewRestyClient().R().
		Method(http.MethodGet, server.URL+"/users/{id}").
		SetPathValues(map[string]string{"id": "100"}).
		SetQueryValues(query).
		Execute()
	if err != nil {
		t.Fatal(err)
	}
	if response.String() != "ok" {
		t.Fatalf("unexpected response: %s", response.String())
	}
}

func TestPostJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(HeaderContentType) != "application/json; charset=utf-8" {
			t.Fatalf("unexpected content type: %s", r.Header.Get(HeaderContentType))
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	response, err := NewRestyClient().R().PostJSON(server.URL, `{"name":"toolkit"}`, "utf-8")
	if err != nil {
		t.Fatal(err)
	}
	if response.String() != "ok" {
		t.Fatalf("unexpected response: %s", response.String())
	}
}

func TestSetDownloadFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("download content"))
	}))
	defer server.Close()

	filepath := filepath.Join(t.TempDir(), "download.txt")
	response, err := NewRestyClient().R().SetDownloadFile(filepath).Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if response.IsError() {
		t.Fatalf("unexpected response status: %s", response.Status())
	}
	content, err := os.ReadFile(filepath)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "download content" {
		t.Fatalf("unexpected file content: %s", string(content))
	}
}

func TestDisableTLSVerify(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	response, err := NewRestyClient().DisableTLSVerify().R().Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if response.String() != "ok" {
		t.Fatalf("unexpected response: %s", response.String())
	}
}

func TestNewRestyClientWithMultiProxyFallback(t *testing.T) {
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("fallback-proxy"))
	}))
	defer proxy.Close()

	response, err := NewRestyClientWithMultiProxy([]string{proxy.URL}).R().
		Get("http://fallback-proxy.test")
	if err != nil {
		t.Fatal(err)
	}
	if response.String() != "fallback-proxy" {
		t.Fatalf("unexpected response: %s", response.String())
	}
	if response.Proxy() != proxy.URL {
		t.Fatalf("expected proxy %q, got %q", proxy.URL, response.Proxy())
	}
}

func TestResponseProxy(t *testing.T) {
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer proxy.Close()

	response, err := NewRestyClient(proxy.URL).R().Get("http://single-proxy.test")
	if err != nil {
		t.Fatal(err)
	}
	if response.Proxy() != proxy.URL {
		t.Fatalf("expected proxy %q, got %q", proxy.URL, response.Proxy())
	}
}

func TestMultiProxyChoosesForEveryRoundTrip(t *testing.T) {
	newProxy := func(responseBody string) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(responseBody))
		}))
	}
	proxy1 := newProxy("proxy-1")
	defer proxy1.Close()
	proxy2 := newProxy("proxy-2")
	defer proxy2.Close()

	choose := &roundRobinChooseProxy{}
	client := NewRestyClientWithMultiProxy([]string{proxy1.URL, proxy2.URL}, choose)
	expectedResponses := []struct {
		body  string
		proxy string
	}{
		{body: "proxy-1", proxy: proxy1.URL},
		{body: "proxy-2", proxy: proxy2.URL},
		{body: "proxy-1", proxy: proxy1.URL},
	}
	for _, expected := range expectedResponses {
		response, err := client.R().Get("http://multi-proxy.test")
		if err != nil {
			t.Fatal(err)
		}
		if response.String() != expected.body {
			t.Fatalf("expected response %q, got %q", expected.body, response.String())
		}
		if response.Proxy() != expected.proxy {
			t.Fatalf("expected proxy %q, got %q", expected.proxy, response.Proxy())
		}
	}
	if choose.count != 3 {
		t.Fatalf("expected choose proxy to be called three times, got %d", choose.count)
	}
}
