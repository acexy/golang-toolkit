package json

import (
	"github.com/acexy/golang-toolkit/util/coll"
	"github.com/tidwall/gjson"
)

// GJSONValue wraps a gjson query result.
type GJSONValue struct {
	gr gjson.Result
}

// GJSONObject wraps a JSON object for key-based access.
type GJSONObject struct {
	m map[string]gjson.Result
}

// GetRawJSON returns raw JSON at a gjson path.
func GetRawJSON(json, path string) string {
	return gjson.Get(json, path).Raw
}

// GetStringValue returns the converted string value and whether the path exists.
func GetStringValue(json, path string) (string, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return "", false
	}
	return result.String(), true
}

// GetIntValue returns the converted int value and whether the path exists.
func GetIntValue(json, path string) (int64, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return 0, false
	}
	return result.Int(), true
}

// GetUintValue returns the converted uint value and whether the path exists.
func GetUintValue(json, path string) (uint64, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return 0, false
	}
	return result.Uint(), true
}

// GetFloatValue returns the converted float value and whether the path exists.
func GetFloatValue(json, path string) (float64, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return 0, false
	}
	return result.Float(), true
}

// GetBoolValue returns the converted bool value and whether the path exists.
func GetBoolValue(json, path string) (bool, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return false, false
	}
	return result.Bool(), true
}

// GetArrayValue returns array items and whether the path contains an array.
func GetArrayValue(json, path string) ([]*GJSONValue, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() || !result.IsArray() {
		return nil, false
	}
	array := result.Array()
	return coll.SliceCollect(array, func(item gjson.Result) *GJSONValue {
		return &GJSONValue{item}
	}), true
}

// GetObject returns an object and whether the path contains an object.
func GetObject(json, path string) (*GJSONObject, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() || !result.IsObject() {
		return nil, false
	}
	return &GJSONObject{result.Map()}, true
}

// Get returns the named member and whether it exists.
func (g *GJSONObject) Get(jsonKey string) (*GJSONValue, bool) {
	if g == nil {
		return nil, false
	}
	v, ok := g.m[jsonKey]
	if !ok {
		return nil, false
	}
	return &GJSONValue{v}, true
}

// GetRawJSON returns the raw member JSON and whether it exists.
func (g *GJSONObject) GetRawJSON(jsonKey string) (string, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return "", false
	}
	return value.RawJSONString(), true
}

// GetStringValue returns the converted member value and whether it exists.
func (g *GJSONObject) GetStringValue(jsonKey string) (string, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return "", false
	}
	return value.StringValue()
}

// GetStringOrZero returns the converted member value, or its zero value when absent.
func (g *GJSONObject) GetStringOrZero(jsonKey string) string {
	r, _ := g.GetStringValue(jsonKey)
	return r
}

// GetIntValue returns the converted member value and whether it exists.
func (g *GJSONObject) GetIntValue(jsonKey string) (int64, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return 0, false
	}
	return value.IntValue()
}

// GetIntOrZero returns the converted member value, or its zero value when absent.
func (g *GJSONObject) GetIntOrZero(jsonKey string) int64 {
	r, _ := g.GetIntValue(jsonKey)
	return r
}

// GetUintValue returns the converted member value and whether it exists.
func (g *GJSONObject) GetUintValue(jsonKey string) (uint64, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return 0, false
	}
	return value.UintValue()
}

// GetUintOrZero returns the converted member value, or its zero value when absent.
func (g *GJSONObject) GetUintOrZero(jsonKey string) uint64 {
	r, _ := g.GetUintValue(jsonKey)
	return r
}

// GetFloatValue returns the converted member value and whether it exists.
func (g *GJSONObject) GetFloatValue(jsonKey string) (float64, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return 0, false
	}
	return value.FloatValue()
}

// GetFloatOrZero returns the converted member value, or its zero value when absent.
func (g *GJSONObject) GetFloatOrZero(jsonKey string) float64 {
	r, _ := g.GetFloatValue(jsonKey)
	return r
}

// GetBoolValue returns the converted member value and whether it exists.
func (g *GJSONObject) GetBoolValue(jsonKey string) (bool, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return false, false
	}
	return value.BoolValue()
}

// GetBoolOrZero returns the converted member value, or its zero value when absent.
func (g *GJSONObject) GetBoolOrZero(jsonKey string) bool {
	r, _ := g.GetBoolValue(jsonKey)
	return r
}

// GetArrayValue returns array items and whether the member contains an array.
func (g *GJSONObject) GetArrayValue(jsonKey string) ([]*GJSONValue, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return nil, false
	}
	return value.ArrayValue()
}

// GetObject returns a member object and whether it has object shape.
func (g *GJSONObject) GetObject(jsonKey string) (*GJSONObject, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return nil, false
	}
	return value.ObjectValue()
}

// NewGJSON parses a JSON string for gjson queries.
func NewGJSON(json string) *GJSONValue {
	return &GJSONValue{
		gjson.Parse(json),
	}
}

// NewGJSONBytes parses JSON bytes for gjson queries.
func NewGJSONBytes(json []byte) *GJSONValue {
	return &GJSONValue{
		gjson.ParseBytes(json),
	}
}

// Get queries a nested path, returning an empty result for a nil receiver.
func (g *GJSONValue) Get(path string) *GJSONValue {
	if g == nil {
		return &GJSONValue{}
	}
	return &GJSONValue{
		g.gr.Get(path),
	}
}

// ForEach visits object or array entries until the callback returns false.
func (g *GJSONValue) ForEach(fn func(key, value gjson.Result) bool) {
	if g == nil || fn == nil {
		return
	}
	g.gr.ForEach(func(key, value gjson.Result) bool {
		return fn(key, value)
	})
}

// StringValue returns the converted value and whether the result exists.
func (g *GJSONValue) StringValue() (string, bool) {
	if g == nil || !g.gr.Exists() {
		return "", false
	}
	return g.gr.String(), true
}

// StringOrZero returns the converted value, or its zero value when absent.
func (g *GJSONValue) StringOrZero() string {
	r, _ := g.StringValue()
	return r
}

// IntValue returns the converted value and whether the result exists.
func (g *GJSONValue) IntValue() (int64, bool) {
	if g == nil || !g.gr.Exists() {
		return 0, false
	}
	return g.gr.Int(), true
}

// IntOrZero returns the converted value, or its zero value when absent.
func (g *GJSONValue) IntOrZero() int64 {
	r, _ := g.IntValue()
	return r
}

// UintValue returns the converted value and whether the result exists.
func (g *GJSONValue) UintValue() (uint64, bool) {
	if g == nil || !g.gr.Exists() {
		return 0, false
	}
	return g.gr.Uint(), true
}

// UintOrZero returns the converted value, or its zero value when absent.
func (g *GJSONValue) UintOrZero() uint64 {
	r, _ := g.UintValue()
	return r
}

// FloatValue returns the converted value and whether the result exists.
func (g *GJSONValue) FloatValue() (float64, bool) {
	if g == nil || !g.gr.Exists() {
		return 0, false
	}
	return g.gr.Float(), true
}

// FloatOrZero returns the converted value, or its zero value when absent.
func (g *GJSONValue) FloatOrZero() float64 {
	r, _ := g.FloatValue()
	return r
}

// BoolValue returns the converted value and whether the result exists.
func (g *GJSONValue) BoolValue() (bool, bool) {
	if g == nil || !g.gr.Exists() {
		return false, false
	}
	return g.gr.Bool(), true
}

// BoolOrZero returns the converted value, or its zero value when absent.
func (g *GJSONValue) BoolOrZero() bool {
	r, _ := g.BoolValue()
	return r
}

// ArrayValue returns array items and whether the result contains an array.
func (g *GJSONValue) ArrayValue() ([]*GJSONValue, bool) {
	if g == nil || !g.gr.Exists() || !g.gr.IsArray() {
		return nil, false
	}
	array := g.gr.Array()
	return coll.SliceCollect(array, func(item gjson.Result) *GJSONValue {
		return &GJSONValue{item}
	}), true
}

// RawJSONString returns the raw result JSON, or empty for a nil receiver.
func (g *GJSONValue) RawJSONString() string {
	if g == nil {
		return ""
	}
	return g.gr.Raw
}

// ObjectValue returns the object and whether the result has object shape.
func (g *GJSONValue) ObjectValue() (*GJSONObject, bool) {
	if g == nil || !g.gr.Exists() || !g.gr.IsObject() {
		return nil, false
	}
	return &GJSONObject{g.gr.Map()}, true
}
