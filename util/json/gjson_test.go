package json

import (
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)

var jsonString = `
{
    "name": "John Doe",
    "age": 30,
    "address": {
        "street": "123 Main St",
        "city": "Anytown",
        "state": "CA",
        "zip": "12345"
    },
    "phone": [
        {
            "type": "home",
            "number": "555-555-1234"
        },
        {
            "type": "work",
            "number": "555-555-5678",
			"isIphone": true
        }
	]
}
`

func TestGetValue(t *testing.T) {
	fmt.Println(GetStringValue(jsonString, "name"))
	fmt.Println(GetIntValue(jsonString, "address.zip"))
	fmt.Println(GetStringValue(jsonString, "address.zip1"))
	fmt.Println(GetStringValue(jsonString, "phone.1.type"))
	fmt.Println(GetBoolValue(jsonString, "phone.0.isIphone"))
	fmt.Println(GetArrayValue(jsonString, "phone"))
	fmt.Println(GetRawJson(jsonString, "address"))
	object, _ := GetObject(jsonString, "address")
	fmt.Println(object.GetStringValue("street"))
}

func TestNewGJson(t *testing.T) {
	g := NewGJson(jsonString)
	fmt.Println(g.Get("phone.0.number").StringValue())
	value, _ := g.Get("phone").ArrayValue()
	for _, v := range value {
		fmt.Println(v.RawJsonString())
	}
}

func TestForeach(t *testing.T) {
	str := `[
      [
        1499040000000,
        "0.01634790",
        "0.80000000",
        "0.01575800",
        "0.01577100",
        "148976.11427815",
        1499644799999,
        "2434.19055334",
        308,
        "1756.87402397",
        "28.46694368",
        "17928899.62484339"
      ],
[
        1499040000000,
        "0.01634790",
        "0.80000000",
        "0.01575800",
        "0.01577100",
        "148976.11427815",
        1499644799999,
        "2434.19055334",
        308,
        "1756.87402397",
        "28.46694368",
        "17928899.62484339"
      ]
    ]`
	gjson.Parse(str).ForEach(func(key, value gjson.Result) bool {
		fmt.Println("key:", key.String())
		//openTime := value.Get("0").Int()
		openPrice := value.Get("1").String()
		//high := value.Get("2").String()
		//low := value.Get("3").String()
		closePrice := value.Get("4").String()
		volume := value.Get("5").String()

		fmt.Println("Open:", openPrice, "Close:", closePrice, "Volume:", volume)
		return true // continue
	})

	NewGJson(str).Foreach(func(key, value gjson.Result) bool {
		fmt.Println("key:", key)
		return true
	})

}
