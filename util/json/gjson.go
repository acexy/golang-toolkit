package json

import (
	"github.com/acexy/golang-toolkit/util/coll"
	"github.com/tidwall/gjson"
)

// GJsonValue 结构体
type GJsonValue struct {
	gr gjson.Result
}

// GJsonObject 结构体
type GJsonObject struct {
	m map[string]gjson.Result
}

// GetRawJson 获取指定路径的原始json字符串
func GetRawJson(json, path string) string {
	return gjson.Get(json, path).Raw
}

// GetStringValue 获取指定json结构中的字符串值
func GetStringValue(json, path string) (string, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return "", false
	}
	return result.String(), true
}

// GetIntValue 获取指定json结构中的int值
func GetIntValue(json, path string) (int64, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return 0, false
	}
	return result.Int(), true
}

// GetUintValue 获取指定json结构中的uint值
func GetUintValue(json, path string) (uint64, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return 0, false
	}
	return result.Uint(), true
}

// GetFloatValue 获取指定json结构中的float值
func GetFloatValue(json, path string) (float64, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return 0, false
	}
	return result.Float(), true
}

// GetBoolValue 获取指定json结构中的bool值
func GetBoolValue(json, path string) (bool, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return false, false
	}
	return result.Bool(), true
}

// GetArrayValue 获取指定json结构中的数组
func GetArrayValue(json, path string) ([]*GJsonValue, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return nil, false
	}
	array := result.Array()
	return coll.SliceCollect(array, func(item gjson.Result) *GJsonValue {
		return &GJsonValue{item}
	}), true
}

// GetObject 获取指定json结构中的对象
func GetObject(json, path string) (*GJsonObject, bool) {
	result := gjson.Get(json, path)
	if !result.Exists() {
		return nil, false
	}
	return &GJsonObject{result.Map()}, true
}

// Get 获取指定json结构中的值
func (g *GJsonObject) Get(jsonKey string) (*GJsonValue, bool) {
	v, ok := g.m[jsonKey]
	if !ok {
		return nil, false
	}
	return &GJsonValue{v}, true
}

// GetRawJson 获取指定json结构中的原始json字符串
func (g *GJsonObject) GetRawJson(jsonKey string) (string, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return "", false
	}
	return value.RawJsonString(), true
}

// GetStringValue 获取指定json结构中的字符串值
func (g *GJsonObject) GetStringValue(jsonKey string) (string, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return "", false
	}
	return value.StringValue()
}

// GetStringResult 获取指定json结构中的字符串值
func (g *GJsonObject) GetStringResult(jsonKey string) string {
	r, _ := g.GetStringValue(jsonKey)
	return r
}

// GetIntValue 获取指定json结构中的int值
func (g *GJsonObject) GetIntValue(jsonKey string) (int64, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return 0, false
	}
	return value.IntValue()
}

// GetIntResult 获取指定json结构中的int值
func (g *GJsonObject) GetIntResult(jsonKey string) int64 {
	r, _ := g.GetIntValue(jsonKey)
	return r
}

// GetUintValue 获取指定json结构中的uint值
func (g *GJsonObject) GetUintValue(jsonKey string) (uint64, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return 0, false
	}
	return value.UintValue()
}

// GetUintResult 获取指定json结构中的uint值
func (g *GJsonObject) GetUintResult(jsonKey string) uint64 {
	r, _ := g.GetUintValue(jsonKey)
	return r
}

// GetFloatValue 获取指定json结构中的float值
func (g *GJsonObject) GetFloatValue(jsonKey string) (float64, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return 0, false
	}
	return value.FloatValue()
}

// GetFloatResult 获取指定json结构中的float值
func (g *GJsonObject) GetFloatResult(jsonKey string) float64 {
	r, _ := g.GetFloatValue(jsonKey)
	return r
}

// GetBoolValue 获取指定json结构中的bool值
func (g *GJsonObject) GetBoolValue(jsonKey string) (bool, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return false, false
	}
	return value.BoolValue()
}

// GetBoolResult 获取指定json结构中的bool值
func (g *GJsonObject) GetBoolResult(jsonKey string) bool {
	r, _ := g.GetBoolValue(jsonKey)
	return r
}

// GetArrayValue 获取指定json结构中的数组
func (g *GJsonObject) GetArrayValue(jsonKey string) ([]*GJsonValue, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return nil, false
	}
	return value.ArrayValue()
}

// GetObject 获取指定json结构中的对象
func (g *GJsonObject) GetObject(jsonKey string) (*GJsonObject, bool) {
	value, ok := g.Get(jsonKey)
	if !ok {
		return nil, false
	}
	return value.GJsonObjectValue()
}

// NewGJson 创建一个GJsonValue
func NewGJson(json string) *GJsonValue {
	return &GJsonValue{
		gjson.Parse(json),
	}
}

// NewGJsonBytes 创建一个GJsonValue
func NewGJsonBytes(json []byte) *GJsonValue {
	return &GJsonValue{
		gjson.ParseBytes(json),
	}
}

// Get 获取指定json结构中的值
func (g *GJsonValue) Get(path string) *GJsonValue {
	return &GJsonValue{
		g.gr.Get(path),
	}
}

func (g *GJsonValue) Foreach(fn func(key, value gjson.Result) bool) {
	g.gr.ForEach(func(key, value gjson.Result) bool {
		return fn(key, value)
	})
}

// StringValue 获取指定json结构中的字符串值
func (g *GJsonValue) StringValue() (string, bool) {
	if !g.gr.Exists() {
		return "", false
	}
	return g.gr.String(), true
}

// StringResult 获取指定json结构中的字符串值
func (g *GJsonValue) StringResult() string {
	r, _ := g.StringValue()
	return r
}

// IntValue 获取指定json结构中的int值
func (g *GJsonValue) IntValue() (int64, bool) {
	if !g.gr.Exists() {
		return 0, false
	}
	return g.gr.Int(), true
}

// IntResult 获取指定json结构中的int值
func (g *GJsonValue) IntResult() int64 {
	r, _ := g.IntValue()
	return r
}

// UintValue 获取指定json结构中的uint值
func (g *GJsonValue) UintValue() (uint64, bool) {
	if !g.gr.Exists() {
		return 0, false
	}
	return g.gr.Uint(), true
}

// UintResult 获取指定json结构中的uint值
func (g *GJsonValue) UintResult() uint64 {
	r, _ := g.UintValue()
	return r
}

// FloatValue 获取指定json结构中的float值
func (g *GJsonValue) FloatValue() (float64, bool) {
	if !g.gr.Exists() {
		return 0, false
	}
	return g.gr.Float(), true
}

// FloatResult 获取指定json结构中的float值
func (g *GJsonValue) FloatResult() float64 {
	r, _ := g.FloatValue()
	return r
}

// BoolValue 获取指定json结构中的bool值
func (g *GJsonValue) BoolValue() (bool, bool) {
	if !g.gr.Exists() {
		return false, false
	}
	return g.gr.Bool(), true
}

// BoolResult 获取指定json结构中的bool值
func (g *GJsonValue) BoolResult() bool {
	r, _ := g.BoolValue()
	return r
}

// ArrayValue 获取指定json结构中的数组
func (g *GJsonValue) ArrayValue() ([]*GJsonValue, bool) {
	if !g.gr.Exists() {
		return nil, false
	}
	array := g.gr.Array()
	return coll.SliceCollect(array, func(item gjson.Result) *GJsonValue {
		return &GJsonValue{item}
	}), true
}

// RawJsonString 获取指定json结构中的原始json字符串
func (g *GJsonValue) RawJsonString() string {
	return g.gr.Raw
}

// GJsonObjectValue 获取指定json结构中的对象
func (g *GJsonValue) GJsonObjectValue() (*GJsonObject, bool) {
	if !g.gr.Exists() {
		return nil, false
	}
	return &GJsonObject{g.gr.Map()}, true
}
