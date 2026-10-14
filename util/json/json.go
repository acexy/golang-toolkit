package json

import (
	stdjson "encoding/json"
)

// ToString serializes a value as JSON, discarding errors.
func ToString(any any) (json string) {
	json, _ = ToStringError(any)
	return
}

// ToStringError serializes a value as JSON and returns errors.
func ToStringError(any any) (jsonString string, err error) {
	jsonBytes, err := ToBytesError(any)
	if err != nil {
		return
	}
	jsonString = string(jsonBytes)
	return
}

// ToStringPanic serializes a value as JSON and panics on failure.
func ToStringPanic(any any) (jsonString string) {
	jsonString, err := ToStringError(any)
	if err != nil {
		panic(err)
	}
	return
}

// ToStringFormat returns indented JSON, discarding errors.
func ToStringFormat(any any) (jsonFormat string) {
	jsonFormat, _ = ToStringFormatError(any)
	return
}

// ToStringFormatError returns indented JSON and serialization errors.
func ToStringFormatError(any any) (string, error) {
	jsonBytes, err := stdjson.MarshalIndent(any, "", "  ")
	return string(jsonBytes), err
}

// ToBytes serializes a value to JSON bytes, discarding errors.
func ToBytes(any any) (bytes []byte) {
	bytes, _ = ToBytesError(any)
	return
}

// ToBytesError serializes a value to JSON bytes and returns errors.
func ToBytesError(any any) ([]byte, error) {
	return stdjson.Marshal(any)
}

// ToBytesPanic serializes a value to JSON bytes and panics on failure.
func ToBytesPanic(any any) (jsonBytes []byte) {
	jsonBytes, err := ToBytesError(any)
	if err != nil {
		panic(err)
	}
	return
}

// ParseBytes decodes JSON bytes into the destination, discarding errors.
func ParseBytes(bytes []byte, any any) {
	_ = ParseBytesError(bytes, any)
}

// ParseBytesError decodes JSON bytes into the destination and returns errors.
func ParseBytesError(bytes []byte, any any) error {
	return stdjson.Unmarshal(bytes, any)
}

// ParseBytesPanic decodes JSON bytes into the destination and panics on failure.
func ParseBytesPanic(bytes []byte, any any) {
	err := ParseBytesError(bytes, any)
	if err != nil {
		panic(err)
	}
}

// ParseString decodes a JSON string into the destination, discarding errors.
func ParseString(jsonString string, any any) {
	_ = ParseStringError(jsonString, any)
}

// ParseStringError decodes a JSON string into the destination and returns errors.
func ParseStringError(jsonString string, any any) error {
	return ParseBytesError([]byte(jsonString), any)
}

// ParseStringPanic decodes a JSON string into the destination and panics on failure.
func ParseStringPanic(jsonString string, any any) {
	err := ParseStringError(jsonString, any)
	if err != nil {
		panic(err)
	}
}

// CopyStruct copies through JSON serialization, discarding errors.
func CopyStruct(originData, targetStruct any) {
	_ = CopyStructError(originData, targetStruct)
}

// CopyStructError copies through JSON serialization and returns errors.
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

// CopyStructPanic copies through JSON serialization and panics on failure.
func CopyStructPanic(originData, targetStruct any) {
	err := CopyStructError(originData, targetStruct)
	if err != nil {
		panic(err)
	}
}
