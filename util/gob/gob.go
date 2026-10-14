package gob

import (
	"bytes"
	"encoding/gob"
)

// Encode serializes a value to GOB bytes.
func Encode(data any) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode decodes GOB bytes into a writable destination.
func Decode(bs []byte, result any) error {
	dec := gob.NewDecoder(bytes.NewBuffer(bs))
	if err := dec.Decode(result); err != nil {
		return err
	}
	return nil
}
