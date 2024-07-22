package reflect

import (
	"errors"
	"reflect"
)

// NonZeroField 返回结构体的非零字段
func NonZeroField(value interface{}) ([]string, error) {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, errors.New("inputStruct must be a struct or a pointer to a struct")
	}
	var nonZeroFields []string
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !isZeroValue(field) {
			fieldName := typ.Field(i).Name
			nonZeroFields = append(nonZeroFields, fieldName)
		}
	}
	return nonZeroFields, nil
}

func isZeroValue(v reflect.Value) bool {
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}
