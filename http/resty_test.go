package http

import "testing"

func TestGet(t *testing.T) {
	client := NewRestyClient()
	resp, err := client.R().Get("https://github.com")
	if err != nil {
		return
	}

	print(resp.String())

	resp, err = client.R().M(HttpMethodGet, "https://github.com").E()
	if err != nil {
		return
	}

	print(resp.String())
}
