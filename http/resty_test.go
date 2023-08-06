package http

import (
	"fmt"
	"github.com/acexy/golang-toolkit/util"
	"net/http"
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
	fmt.Println(util.ToJson(s))

	resp, err = req.Get("http://t.weather.sojson.com/api/weather/city/101030100")
	if err != nil {
		return
	}
	fmt.Println(resp.String())
	fmt.Println(util.ToJson(s))

}
