package json

import (
	"testing"

	"github.com/tidwall/gjson"
)

var sampleJSONString = `
{
  "name": "John Doe",
  "age": 30,
  "score": 98.5,
  "enabled": true,
  "address": {
    "street": "123 Main St",
    "zip": "12345"
  },
  "phone": [
    {
      "type": "home",
      "number": "555-555-1234"
    },
    {
      "type": "work",
      "number": "555-555-5678"
    }
  ]
}
`

func TestTopLevelGetValue(t *testing.T) {
	if got, ok := GetStringValue(sampleJSONString, "name"); !ok || got != "John Doe" {
		t.Fatalf("GetStringValue() = %s, %v", got, ok)
	}
	if got, ok := GetIntValue(sampleJSONString, "age"); !ok || got != 30 {
		t.Fatalf("GetIntValue() = %d, %v", got, ok)
	}
	if got, ok := GetUintValue(sampleJSONString, "age"); !ok || got != 30 {
		t.Fatalf("GetUintValue() = %d, %v", got, ok)
	}
	if got, ok := GetFloatValue(sampleJSONString, "score"); !ok || got != 98.5 {
		t.Fatalf("GetFloatValue() = %f, %v", got, ok)
	}
	if got, ok := GetBoolValue(sampleJSONString, "enabled"); !ok || !got {
		t.Fatalf("GetBoolValue() = %v, %v", got, ok)
	}
	if _, ok := GetStringValue(sampleJSONString, "missing"); ok {
		t.Fatal("missing value should not exist")
	}
}

func TestTopLevelGetArrayAndObject(t *testing.T) {
	values, ok := GetArrayValue(sampleJSONString, "phone")
	if !ok || len(values) != 2 {
		t.Fatalf("GetArrayValue() len = %d, %v", len(values), ok)
	}
	if _, ok = GetArrayValue(sampleJSONString, "name"); ok {
		t.Fatal("non-array should not be returned as array")
	}

	object, ok := GetObject(sampleJSONString, "address")
	if !ok {
		t.Fatal("expected address object")
	}
	if got := object.GetStringOrZero("street"); got != "123 Main St" {
		t.Fatalf("unexpected street: %s", got)
	}
	if _, ok = GetObject(sampleJSONString, "name"); ok {
		t.Fatal("non-object should not be returned as object")
	}
}

func TestGJSONObjectGetValue(t *testing.T) {
	object, ok := GetObject(sampleJSONString, "address")
	if !ok {
		t.Fatal("expected object")
	}
	if got, ok := object.GetStringValue("zip"); !ok || got != "12345" {
		t.Fatalf("GetStringValue() = %s, %v", got, ok)
	}
	if got := object.GetStringOrZero("missing"); got != "" {
		t.Fatalf("missing value should be zero, got %s", got)
	}

	var nilObject *GJSONObject
	if _, ok := nilObject.Get("x"); ok {
		t.Fatal("nil object should not return value")
	}
}

func TestNewGJSON(t *testing.T) {
	g := NewGJSON(sampleJSONString)
	if got := g.Get("phone.0.number").StringOrZero(); got != "555-555-1234" {
		t.Fatalf("unexpected phone number: %s", got)
	}

	fromBytes := NewGJSONBytes([]byte(sampleJSONString))
	if got := fromBytes.Get("name").StringOrZero(); got != "John Doe" {
		t.Fatalf("unexpected name: %s", got)
	}
}

func TestGJSONValueTypes(t *testing.T) {
	g := NewGJSON(sampleJSONString)
	if got, ok := g.Get("age").IntValue(); !ok || got != 30 {
		t.Fatalf("IntValue() = %d, %v", got, ok)
	}
	if got, ok := g.Get("age").UintValue(); !ok || got != 30 {
		t.Fatalf("UintValue() = %d, %v", got, ok)
	}
	if got, ok := g.Get("score").FloatValue(); !ok || got != 98.5 {
		t.Fatalf("FloatValue() = %f, %v", got, ok)
	}
	if got, ok := g.Get("enabled").BoolValue(); !ok || !got {
		t.Fatalf("BoolValue() = %v, %v", got, ok)
	}
	if _, ok := g.Get("missing").StringValue(); ok {
		t.Fatal("missing value should not exist")
	}

	var nilValue *GJSONValue
	if got := nilValue.StringOrZero(); got != "" {
		t.Fatalf("nil StringOrZero should return empty string, got %s", got)
	}
}

func TestGJSONValueArrayAndObject(t *testing.T) {
	g := NewGJSON(sampleJSONString)
	array, ok := g.Get("phone").ArrayValue()
	if !ok || len(array) != 2 {
		t.Fatalf("ArrayValue() len = %d, %v", len(array), ok)
	}
	if _, ok = g.Get("name").ArrayValue(); ok {
		t.Fatal("non-array should not be returned as array")
	}

	object, ok := g.Get("address").ObjectValue()
	if !ok {
		t.Fatal("expected object")
	}
	raw, ok := object.GetRawJSON("zip")
	if !ok || raw != `"12345"` {
		t.Fatalf("unexpected raw json: %s, %v", raw, ok)
	}
	if _, ok = g.Get("name").ObjectValue(); ok {
		t.Fatal("non-object should not be returned as object")
	}
}

func TestGJSONForEach(t *testing.T) {
	g := NewGJSON(`[{"name":"a"},{"name":"b"}]`)
	count := 0
	g.ForEach(func(key, value gjson.Result) bool {
		count++
		return true
	})
	if count != 2 {
		t.Fatalf("unexpected foreach count: %d", count)
	}

	var nilValue *GJSONValue
	nilValue.ForEach(func(key, value gjson.Result) bool {
		t.Fatal("nil value should not iterate")
		return true
	})
}
