package json

import (
	"bytes"
	"encoding/json"
)

// ToJson 转json字符串 忽略任何错误
func ToJson(any any) (json string) {
	json, _ = ToJsonError(any)
	return
}

// ToJsonError 转json字符串 返回任何错误
func ToJsonError(any any) (jsonString string, err error) {
	jsonBytes, err := ToJsonBytesError(any)
	if err != nil {
		return
	}
	jsonString = string(jsonBytes)
	return
}

// ToJsonFormat 转json字符串并格式化输出 忽略任何错误
func ToJsonFormat(any any) (jsonFormat string) {
	jsonFormat, _ = ToJsonFormatError(any)
	return
}

// ToJsonFormatError 转json字符串并格式化输出 返回任何错误
func ToJsonFormatError(any any) (string, error) {
	jsonString, err := ToJsonError(any)
	if err != nil {
		return "", err
	}
	var formattedJSON bytes.Buffer
	err = json.Indent(&formattedJSON, []byte(jsonString), "", "  ")
	if err == nil {
		return formattedJSON.String(), nil
	}
	return "", err
}

// ToJsonBytes 转json字节 忽略任何错误
func ToJsonBytes(any any) (bytes []byte) {
	bytes, _ = ToJsonBytesError(any)
	return
}

// ToJsonBytesError 转json字节 返回任何错误
func ToJsonBytesError(any any) ([]byte, error) {
	return json.Marshal(any)
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
func ToJsonBytesPanic(any any) (jsonBytes []byte) {
	jsonBytes, err := ToJsonBytesError(any)
	if err != nil {
		panic(err)
	}
	return
}

// ParseJson 转对象 忽略任何错误
func ParseJson(jsonString string, any any) {
	_ = ParseJsonError(jsonString, any)
}

// ParseJsonError 转对象 返回任何错误
func ParseJsonError(jsonString string, any any) error {
	return json.Unmarshal([]byte(jsonString), any)
}

// ParseJsonPanic 转对象 任何错误将触发panic
func ParseJsonPanic(jsonString string, any any) {
	err := ParseJsonError(jsonString, any)
	if err != nil {
		panic(err)
	}
}

// CopyStruct 通过json序列化/反序列化将origin struct复制给target struct 忽略任何错误
func CopyStruct(originData, targetStruct any) {
	_ = CopyStructError(originData, targetStruct)
}

// CopyStructError 通过json序列化/反序列化将origin struct复制给target struct 返回任何错误
func CopyStructError(originData, targetStruct any) error {
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

// CopyStructPanic 通过json序列化/反序列化将origin struct复制给target struct 任何错误将触发panic
func CopyStructPanic(originData, targetStruct any) {
	err := CopyStructError(originData, targetStruct)
	if err != nil {
		panic(err)
	}
}
