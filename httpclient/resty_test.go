package httpclient

import (
	"fmt"
	"github.com/acexy/golang-toolkit/math/conversion"
	"github.com/acexy/golang-toolkit/util/json"
	"net/http"
	"testing"
	"time"
)

var client *RestyClient

func init() {
	client = NewRestyClient()
}

func TestPoxyClient(t *testing.T) {
	response, err := NewRestyClient("localhost:7890").R().Get("https://google.com")
	if err != nil {
		return
	}
	fmt.Println(response.RawResponse)
}

func TestGet(t *testing.T) {
	resp, err := client.SetTimeout(time.Second * 3).R().Get("https://github.com")
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	println(resp.String())

	resp, err = client.R().M(http.MethodGet, "https://github.com").E()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Println(resp.String())

	s := struct {
		Message  string `json:"message"`
		CityInfo struct {
			City string `json:"city"`
		} `json:"cityInfo"`
		Data struct {
			Shidu string `json:"shidu"`
		} `json:"data"`
	}{}

	req := client.R().
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36").
		SetReturnStruct(&s)

	resp, err = req.Get("http://t.weather.sojson.com/api/weather/city/101030100")
	if err != nil {
		return
	}
	fmt.Println(resp.String())
	fmt.Println(json.ToJson(s))

	resp, err = req.Get("http://t.weather.sojson.com/api/weather/city/101030100")
	if err != nil {
		return
	}
	fmt.Println(resp.String())
	fmt.Println(json.ToJson(s))
}

func TestProxy(t *testing.T) {
	restyClient := NewRestyClient()

	restyClient.SetTimeout(time.Second * 3)
	resp, err := restyClient.R().Get("https://www.google.com")
	if err != nil {
		t.Errorf("%v\n", err)
	} else {
		fmt.Println(conversion.FromBytes(resp.Body()))
	}

	restyClient.SetProxy("http://127.0.0.1:7890")
	resp, err = restyClient.R().M(http.MethodGet, "https://www.google.com").E()
	if err != nil {
		t.Errorf("%v\n", err)
	} else {
		fmt.Println(conversion.FromBytes(resp.Body()))
	}

	resp, err = restyClient.R().Get("https://www.google.com")
	if err != nil {
		t.Errorf("%v\n", err)
	} else {
		fmt.Println(conversion.FromBytes(resp.Body()))
	}

	restyClient.SetProxy("http://127.0.0.1:1234")
	resp, err = restyClient.R().M(http.MethodGet, "https://google.com").E()
	if err != nil {
		t.Errorf("%v\n", err)
	} else {
		fmt.Println(conversion.FromBytes(resp.Body()))
	}

}

func TestMultiProxy(t *testing.T) {
	restyClient := NewRestyClientWithMultiProxy([]string{
		"http://localhost:7890",
		"http://localhost:1234",
	})

	restyClient.SetTimeout(time.Second * 3)
	resp, err := restyClient.R().Get("https://www.google.com")
	if err != nil {
		t.Errorf("%v\n", err)
	} else {
		fmt.Println(conversion.FromBytes(resp.Body()))
	}

	resp, err = restyClient.R().M(http.MethodGet, "https://www.google.com").E()
	if err != nil {
		t.Errorf("%v\n", err)
	} else {
		fmt.Println(conversion.FromBytes(resp.Body()))
	}

	resp, err = restyClient.R().Get("https://www.google.com")
	if err != nil {
		t.Errorf("%v\n", err)
	} else {
		fmt.Println(conversion.FromBytes(resp.Body()))
	}

	resp, err = restyClient.R().M(http.MethodGet, "https://google.com").E()
	if err != nil {
		t.Errorf("%v\n", err)
	} else {
		fmt.Println(conversion.FromBytes(resp.Body()))
	}

}
