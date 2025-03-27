package json

import (
	"fmt"
	"testing"
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
