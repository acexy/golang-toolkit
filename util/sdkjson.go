package util

import "encoding/json"

// ToJson 转json字符串 忽略任何错误
func ToJson(any any) (json string) {
	json, _ = ToJsonError(any)
	return
}

// ToJsonBytes 转json字节 忽略任何错误
func ToJsonBytes(any any) (bytes []byte) {
	bytes, _ = ToJsonBytesError(any)
	return
}

// ParseJson 转对象 忽略任何错误
func ParseJson(jsonString string, any any) {
	_ = ParseJsonError(jsonString, any)
}

// ToJsonError 转json字符串 返回任何错误
func ToJsonError(any any) (jsonString string, err error) {
	bytes, err := ToJsonBytesError(any)
	if err != nil {
		return
	}
	jsonString = string(bytes)
	return
}

// ToJsonBytesError 转json字节 返回任何错误
func ToJsonBytesError(any any) ([]byte, error) {
	return json.Marshal(any)
}

// ParseJsonError 转对象 返回任何错误
func ParseJsonError(jsonString string, any any) error {
	return json.Unmarshal([]byte(jsonString), any)
}

// ToJsonPanic 转json字符串 任何错误将触发panic
func ToJsonPanic(any any) (jsonString string) {
	jsonString, err := ToJsonError(any)
	if err != nil {
		panic(err)
	}
	return
}

// ToJsonBytesPanic 转json字节 任何错误将触发panic
func ToJsonBytesPanic(any any) (bytes []byte) {
	bytes, err := ToJsonBytesError(any)
	if err != nil {
		panic(err)
	}
	return
}

// ParseJsonPanic 转对象 任何错误将触发panic
func ParseJsonPanic(jsonString string, any any) {
	err := ParseJsonError(jsonString, any)
	if err != nil {
		panic(err)
	}
}

// JsonCopyStruct 通过json序列化/反序列化将origin struct复制给target struct 忽略任何错误
func JsonCopyStruct(originData, targetStruct any) {
	_ = JsonCopyStructError(originData, targetStruct)
}

// JsonCopyStructError 通过json序列化/反序列化将origin struct复制给target struct 返回任何错误
func JsonCopyStructError(originData, targetStruct any) error {
	jsonString, err := ToJsonError(originData)
	if err != nil {
		return err
	}
	err = ParseJsonError(jsonString, targetStruct)
	if err != nil {
		return err
	}
	return nil
}

// JsonCopyStructPanic 通过json序列化/反序列化将origin struct复制给target struct 任何错误将触发panic
func JsonCopyStructPanic(originData, targetStruct any) {
	err := JsonCopyStructError(originData, targetStruct)
	if err != nil {
		panic(err)
	}
}