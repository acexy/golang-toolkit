package json

import (
	stdjson "encoding/json"
)

// ToString 转json字符串 忽略任何错误
func ToString(any any) (json string) {
	json, _ = ToStringError(any)
	return
}

// ToStringError 转json字符串 返回任何错误
func ToStringError(any any) (jsonString string, err error) {
	jsonBytes, err := ToBytesError(any)
	if err != nil {
		return
	}
	jsonString = string(jsonBytes)
	return
}

// ToStringPanic 转json字符串 任何错误将触发panic
func ToStringPanic(any any) (jsonString string) {
	jsonString, err := ToStringError(any)
	if err != nil {
		panic(err)
	}
	return
}

// ToStringFormat 转json字符串并格式化输出 忽略任何错误
func ToStringFormat(any any) (jsonFormat string) {
	jsonFormat, _ = ToStringFormatError(any)
	return
}

// ToStringFormatError 转json字符串并格式化输出 返回任何错误
func ToStringFormatError(any any) (string, error) {
	jsonBytes, err := stdjson.MarshalIndent(any, "", "  ")
	return string(jsonBytes), err
}

// ToBytes 转json字节 忽略任何错误
func ToBytes(any any) (bytes []byte) {
	bytes, _ = ToBytesError(any)
	return
}

// ToBytesError 转json字节 返回任何错误
func ToBytesError(any any) ([]byte, error) {
	return stdjson.Marshal(any)
}

// ToBytesPanic 转json字节 任何错误将触发panic
func ToBytesPanic(any any) (jsonBytes []byte) {
	jsonBytes, err := ToBytesError(any)
	if err != nil {
		panic(err)
	}
	return
}

// ParseBytes 将byte数据转化成对象 忽略任何错误
func ParseBytes(bytes []byte, any any) {
	_ = ParseBytesError(bytes, any)
}

// ParseBytesError 将byte数据转化成对象 返回任何错误
func ParseBytesError(bytes []byte, any any) error {
	return stdjson.Unmarshal(bytes, any)
}

// ParseBytesPanic 将byte数据转化成对象 任何错误将触发panic
func ParseBytesPanic(bytes []byte, any any) {
	err := ParseBytesError(bytes, any)
	if err != nil {
		panic(err)
	}
}

// ParseString 转对象 忽略任何错误
func ParseString(jsonString string, any any) {
	_ = ParseStringError(jsonString, any)
}

// ParseStringError 转对象 返回任何错误
func ParseStringError(jsonString string, any any) error {
	return ParseBytesError([]byte(jsonString), any)
}

// ParseStringPanic 转对象 任何错误将触发panic
func ParseStringPanic(jsonString string, any any) {
	err := ParseStringError(jsonString, any)
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
	jsonBytes, err := ToBytesError(originData)
	if err != nil {
		return err
	}
	err = ParseBytesError(jsonBytes, targetStruct)
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
