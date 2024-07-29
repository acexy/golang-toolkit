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

// NonZeroFieldValue 返回结构体的非零字段的值
func NonZeroFieldValue(value interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, errors.New("inputStruct must be a struct or a pointer to a struct")
	}
	nonZeroValue := make(map[string]interface{})
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !isZeroValue(field) {
			fieldName := typ.Field(i).Name
			nonZeroValue[fieldName] = field.Interface()
		}
	}

	return nonZeroValue, nil
}

// DeepCopy 深拷贝 源和目标需要是同类型
func DeepCopy(src interface{}) interface{} {
	if src == nil {
		return nil
	}
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.New(srcVal.Type()).Elem()
	deepCopyRecursive(srcVal, dstVal)
	return dstVal.Interface()
}

func deepCopyRecursive(src, dst reflect.Value) {
	switch src.Kind() {
	case reflect.Ptr:
		if !src.IsNil() {
			dst.Set(reflect.New(src.Elem().Type()))
			deepCopyRecursive(src.Elem(), dst.Elem())
		}
	case reflect.Interface:
		if !src.IsNil() {
			newValue := reflect.New(src.Elem().Type()).Elem()
			deepCopyRecursive(src.Elem(), newValue)
			dst.Set(newValue)
		}
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			deepCopyRecursive(src.Field(i), dst.Field(i))
		}
	case reflect.Slice:
		if !src.IsNil() {
			dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
			for i := 0; i < src.Len(); i++ {
				deepCopyRecursive(src.Index(i), dst.Index(i))
			}
		}
	case reflect.Array:
		for i := 0; i < src.Len(); i++ {
			deepCopyRecursive(src.Index(i), dst.Index(i))
		}
	case reflect.Map:
		if !src.IsNil() {
			dst.Set(reflect.MakeMapWithSize(src.Type(), src.Len()))
			for _, key := range src.MapKeys() {
				newKey := reflect.New(key.Type()).Elem()
				newValue := reflect.New(src.MapIndex(key).Type()).Elem()
				deepCopyRecursive(key, newKey)
				deepCopyRecursive(src.MapIndex(key), newValue)
				dst.SetMapIndex(newKey, newValue)
			}
		}
	case reflect.Chan:
	case reflect.Func:
	default:
		dst.Set(src)
	}
}

func isZeroValue(v reflect.Value) bool {
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}
