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
	client := NewRestyClientWithMultiProxy([]string{"http://127.0.0.1:7890"})
	if client == nil {
		t.Fatal("expected fallback resty client")
	}
}
