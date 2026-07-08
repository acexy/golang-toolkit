package json

import (
	"github.com/acexy/golang-toolkit/util/coll"
	"github.com/tidwall/gjson"
)

// GJSONValue 表示一个 gjson 查询结果
type GJSONValue struct {
	gr gjson.Result
}

// GJSONObject 表示一个 JSON 对象
type GJSONObject struct {
	m map[string]gjson.Result
}

// GetRawJSON 获取指定路径的原始 JSON 字符串
func GetRawJSON(json, path string) string {
	return gjson.Get(json, path).Raw
}

// GetStringValue 获取指定 JSON 结构中的字符串值
func GetStringValue(json, path string) (string, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return "", false
	}
	return result.String(), true
}

// GetIntValue 获取指定 JSON 结构中的 int 值
func GetIntValue(json, path string) (int64, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return 0, false
	}
	return result.Int(), true
}

// GetUintValue 获取指定 JSON 结构中的 uint 值
func GetUintValue(json, path string) (uint64, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return 0, false
	}
	return result.Uint(), true
}

// GetFloatValue 获取指定 JSON 结构中的 float 值
func GetFloatValue(json, path string) (float64, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return 0, false
	}
	return result.Float(), true
}

// GetBoolValue 获取指定 JSON 结构中的 bool 值
func GetBoolValue(json, path string) (bool, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return false, false
	}
	return result.Bool(), true
}

// GetArrayValue 获取指定 JSON 结构中的数组
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

// GetObject 获取指定 JSON 结构中的对象
func GetObject(json, path string) (*GJSONObject, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() || !result.IsObject() {
		return nil, false
	}
	return &GJSONObject{result.Map()}, true
}

// Get 获取指定 JSON 结构中的值
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

// GetRawJSON 获取指定 JSON 结构中的原始 JSON 字符串
func (g *GJSONObject) GetRawJSON(jsonKey string) (string, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return "", false
	}
	return value.RawJSONString(), true
}

// GetStringValue 获取指定 JSON 结构中的字符串值
func (g *GJSONObject) GetStringValue(jsonKey string) (string, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return "", false
	}
	return value.StringValue()
}

// GetStringOrZero 获取指定 JSON 结构中的字符串值，不存在时返回零值
func (g *GJSONObject) GetStringOrZero(jsonKey string) string {
	r, _ := g.GetStringValue(jsonKey)
	return r
}

// GetIntValue 获取指定 JSON 结构中的 int 值
func (g *GJSONObject) GetIntValue(jsonKey string) (int64, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return 0, false
	}
	return value.IntValue()
}

// GetIntOrZero 获取指定 JSON 结构中的 int 值，不存在时返回零值
func (g *GJSONObject) GetIntOrZero(jsonKey string) int64 {
	r, _ := g.GetIntValue(jsonKey)
	return r
}

// GetUintValue 获取指定 JSON 结构中的 uint 值
func (g *GJSONObject) GetUintValue(jsonKey string) (uint64, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return 0, false
	}
	return value.UintValue()
}

// GetUintOrZero 获取指定 JSON 结构中的 uint 值，不存在时返回零值
func (g *GJSONObject) GetUintOrZero(jsonKey string) uint64 {
	r, _ := g.GetUintValue(jsonKey)
	return r
}

// GetFloatValue 获取指定 JSON 结构中的 float 值
func (g *GJSONObject) GetFloatValue(jsonKey string) (float64, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return 0, false
	}
	return value.FloatValue()
}

// GetFloatOrZero 获取指定 JSON 结构中的 float 值，不存在时返回零值
func (g *GJSONObject) GetFloatOrZero(jsonKey string) float64 {
	r, _ := g.GetFloatValue(jsonKey)
	return r
}

// GetBoolValue 获取指定 JSON 结构中的 bool 值
func (g *GJSONObject) GetBoolValue(jsonKey string) (bool, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return false, false
	}
	return value.BoolValue()
}

// GetBoolOrZero 获取指定 JSON 结构中的 bool 值，不存在时返回零值
func (g *GJSONObject) GetBoolOrZero(jsonKey string) bool {
	r, _ := g.GetBoolValue(jsonKey)
	return r
}

// GetArrayValue 获取指定 JSON 结构中的数组
func (g *GJSONObject) GetArrayValue(jsonKey string) ([]*GJSONValue, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return nil, false
	}
	return value.ArrayValue()
}

// GetObject 获取指定 JSON 结构中的对象
func (g *GJSONObject) GetObject(jsonKey string) (*GJSONObject, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return nil, false
	}
	return value.ObjectValue()
}

// NewGJSON 创建一个 GJSONValue
func NewGJSON(json string) *GJSONValue {
	return &GJSONValue{
		gjson.Parse(json),
	}
}

// NewGJSONBytes 创建一个 GJSONValue
func NewGJSONBytes(json []byte) *GJSONValue {
	return &GJSONValue{
		gjson.ParseBytes(json),
	}
}

// Get 获取指定 JSON 结构中的值
func (g *GJSONValue) Get(path string) *GJSONValue {
	if g == nil {
		return &GJSONValue{}
	}
	return &GJSONValue{
		g.gr.Get(path),
	}
}

// ForEach 遍历 JSON 对象或数组
func (g *GJSONValue) ForEach(fn func(key, value gjson.Result) bool) {
	if g == nil || fn == nil {
		return
	}
	g.gr.ForEach(func(key, value gjson.Result) bool {
		return fn(key, value)
	})
}

// StringValue 获取指定 JSON 结构中的字符串值
func (g *GJSONValue) StringValue() (string, bool) {
	if g == nil || !g.gr.Exists() {
		return "", false
	}
	return g.gr.String(), true
}

// StringOrZero 获取指定 JSON 结构中的字符串值，不存在时返回零值
func (g *GJSONValue) StringOrZero() string {
	r, _ := g.StringValue()
	return r
}

// IntValue 获取指定 JSON 结构中的 int 值
func (g *GJSONValue) IntValue() (int64, bool) {
	if g == nil || !g.gr.Exists() {
		return 0, false
	}
	return g.gr.Int(), true
}

// IntOrZero 获取指定 JSON 结构中的 int 值，不存在时返回零值
func (g *GJSONValue) IntOrZero() int64 {
	r, _ := g.IntValue()
	return r
}

// UintValue 获取指定 JSON 结构中的 uint 值
func (g *GJSONValue) UintValue() (uint64, bool) {
	if g == nil || !g.gr.Exists() {
		return 0, false
	}
	return g.gr.Uint(), true
}

// UintOrZero 获取指定 JSON 结构中的 uint 值，不存在时返回零值
func (g *GJSONValue) UintOrZero() uint64 {
	r, _ := g.UintValue()
	return r
}

// FloatValue 获取指定 JSON 结构中的 float 值
func (g *GJSONValue) FloatValue() (float64, bool) {
	if g == nil || !g.gr.Exists() {
		return 0, false
	}
	return g.gr.Float(), true
}

// FloatOrZero 获取指定 JSON 结构中的 float 值，不存在时返回零值
func (g *GJSONValue) FloatOrZero() float64 {
	r, _ := g.FloatValue()
	return r
}

// BoolValue 获取指定 JSON 结构中的 bool 值
func (g *GJSONValue) BoolValue() (bool, bool) {
	if g == nil || !g.gr.Exists() {
		return false, false
	}
	return g.gr.Bool(), true
}

// BoolOrZero 获取指定 JSON 结构中的 bool 值，不存在时返回零值
func (g *GJSONValue) BoolOrZero() bool {
	r, _ := g.BoolValue()
	return r
}

// ArrayValue 获取指定 JSON 结构中的数组
func (g *GJSONValue) ArrayValue() ([]*GJSONValue, bool) {
	if g == nil || !g.gr.Exists() || !g.gr.IsArray() {
		return nil, false
	}
	array := g.gr.Array()
	return coll.SliceCollect(array, func(item gjson.Result) *GJSONValue {
		return &GJSONValue{item}
	}), true
}

// RawJSONString 获取指定 JSON 结构中的原始 JSON 字符串
func (g *GJSONValue) RawJSONString() string {
	if g == nil {
		return ""
	}
	return g.gr.Raw
}

// ObjectValue 获取指定 JSON 结构中的对象
func (g *GJSONValue) ObjectValue() (*GJSONObject, bool) {
	if g == nil || !g.gr.Exists() || !g.gr.IsObject() {
		return nil, false
	}
	return &GJSONObject{g.gr.Map()}, true
}
