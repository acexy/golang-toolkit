package http

import (
	"github.com/acexy/golang-toolkit/util"
	"testing"
	"time"
)

var client *RestyClient

func init() {
	client = NewRestyClient()
}

func TestGet(t *testing.T) {
	resp, err := client.SetTimeout(time.Second * 3).R().Get("https://github.com")
	if err != nil {
		return
	}

	println(resp.String())

	resp, err = client.R().M(HttpMethodGet, "https://github.com").E()
	if err != nil {
		return
	}

	println(resp.String())

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
	println(resp.String())
	println(util.ToJsonString(s))

	resp, err = req.Get("http://t.weather.sojson.com/api/weather/city/101030100")
	if err != nil {
		return
	}
	println(resp.String())
	println(util.ToJsonString(s))

}
