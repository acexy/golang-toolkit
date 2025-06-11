package reflect

import (
	"errors"
	"reflect"
)

// processStructFields 处理结构体字段的通用方法
func processStructFields(value interface{}, filter func(field reflect.Value) bool, process func(fieldName string, field reflect.Value)) error {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return errors.New("inputStruct must be a struct or a pointer to a struct")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		// 跳过未导出的字段
		if !field.CanInterface() {
			continue
		}
		if filter == nil || filter(field) {
			fieldName := typ.Field(i).Name
			process(fieldName, field)
		}
	}
	return nil
}

// AllFieldName 返回结构体的所有字段
func AllFieldName(value interface{}) ([]string, error) {
	var allFields []string
	err := processStructFields(value,
		nil, // 不需要过滤
		func(fieldName string, field reflect.Value) {
			allFields = append(allFields, fieldName)
		})
	return allFields, err
}

// AllFieldValue 返回结构体的所有字段的值
func AllFieldValue(value interface{}) (map[string]interface{}, error) {
	allValue := make(map[string]interface{})
	err := processStructFields(value,
		nil, // 不需要过滤
		func(fieldName string, field reflect.Value) {
			allValue[fieldName] = field.Interface()
		})
	return allValue, err
}

// NonZeroFieldName 返回结构体的非零字段
func NonZeroFieldName(value interface{}) ([]string, error) {
	var nonZeroFields []string
	err := processStructFields(value,
		func(field reflect.Value) bool {
			return !isZeroValue(field)
		},
		func(fieldName string, field reflect.Value) {
			nonZeroFields = append(nonZeroFields, fieldName)
		})
	return nonZeroFields, err
}

// NonZeroFieldValue 返回结构体的非零字段的值
func NonZeroFieldValue(value interface{}) (map[string]interface{}, error) {
	nonZeroValue := make(map[string]interface{})
	err := processStructFields(value,
		func(field reflect.Value) bool {
			return !isZeroValue(field)
		},
		func(fieldName string, field reflect.Value) {
			nonZeroValue[fieldName] = field.Interface()
		})
	return nonZeroValue, err
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
	if !v.IsValid() || !v.CanInterface() {
		return true // 无效或不可导出的字段视为零值
	}
	// 根据类型检查零值
	zero := reflect.Zero(v.Type())
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		return v.Len() == 0 // 切片、数组、映射或通道的零值为长度为0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil() // 指针和接口零值为nil
	default:
		// 对于其他类型，使用通用比较
		return reflect.DeepEqual(v.Interface(), zero.Interface())
	}
}
