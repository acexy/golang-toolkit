package gob

import (
	"bytes"
	"encoding/gob"
)

// Encode 将制定数据进行gob编码
func Encode(data any) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode 将制定数据进行gob解码
func Decode(bs []byte, result any) error {
	dec := gob.NewDecoder(bytes.NewBuffer(bs))
	if err := dec.Decode(result); err != nil {
		return err
	}
	return nil
}
