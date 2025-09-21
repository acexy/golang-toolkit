package reflect

import (
	"errors"
	"fmt"
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

// SetFieldValue 将 map[string]any 中的值赋给 structPtr 指向的结构体中对应字段
func SetFieldValue(structPtr any, values map[string]any, checkMissingField ...bool) error {
	v := reflect.ValueOf(structPtr)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("input must be a non-nil pointer to a struct")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("input must be a pointer to a struct")
	}

	for fieldName, val := range values {
		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			if len(checkMissingField) > 0 && checkMissingField[0] {
				return fmt.Errorf("field %s does not exist in struct", fieldName)
			}
			continue // 忽略不存在的字段
		}
		if !field.CanSet() {
			return fmt.Errorf("field %s cannot be set (may be unexported)", fieldName)
		}
		value := reflect.ValueOf(val)
		fieldType := field.Type()

		// 自动兼容基本类型转换
		if value.Type().ConvertibleTo(fieldType) {
			field.Set(value.Convert(fieldType))
			continue
		}

		// 尝试处理不同数值类型间的兼容
		if converted, ok := tryConvertNumber(value, fieldType); ok {
			field.Set(converted)
			continue
		}

		return fmt.Errorf("cannot assign value of type %s to field %s (type %s)", value.Type(), fieldName, fieldType)
	}

	return nil
}

// tryConvertNumber 尝试在不同的数字类型之间转换（int/uint/float）
func tryConvertNumber(val reflect.Value, targetType reflect.Type) (reflect.Value, bool) {
	if !val.IsValid() {
		return reflect.Value{}, false
	}
	kind := val.Kind()
	tKind := targetType.Kind()
	if isNumericKind(kind) && isNumericKind(tKind) {
		// 转换为 float64 作为中介
		var floatVal float64
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			floatVal = float64(val.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			floatVal = float64(val.Uint())
		case reflect.Float32, reflect.Float64:
			floatVal = val.Float()
		default:
			return reflect.Value{}, false
		}
		// 转换为目标类型
		switch tKind {
		case reflect.Int:
			return reflect.ValueOf(int(floatVal)).Convert(targetType), true
		case reflect.Int8:
			return reflect.ValueOf(int8(floatVal)).Convert(targetType), true
		case reflect.Int16:
			return reflect.ValueOf(int16(floatVal)).Convert(targetType), true
		case reflect.Int32:
			return reflect.ValueOf(int32(floatVal)).Convert(targetType), true
		case reflect.Int64:
			return reflect.ValueOf(int64(floatVal)).Convert(targetType), true
		case reflect.Uint:
			return reflect.ValueOf(uint(floatVal)).Convert(targetType), true
		case reflect.Uint8:
			return reflect.ValueOf(uint8(floatVal)).Convert(targetType), true
		case reflect.Uint16:
			return reflect.ValueOf(uint16(floatVal)).Convert(targetType), true
		case reflect.Uint32:
			return reflect.ValueOf(uint32(floatVal)).Convert(targetType), true
		case reflect.Uint64:
			return reflect.ValueOf(uint64(floatVal)).Convert(targetType), true
		case reflect.Float32:
			return reflect.ValueOf(float32(floatVal)).Convert(targetType), true
		case reflect.Float64:
			return reflect.ValueOf(floatVal).Convert(targetType), true
		default:
			return reflect.Value{}, false
		}
	}
	return reflect.Value{}, false
}

func isNumericKind(kind reflect.Kind) bool {
	return kind >= reflect.Int && kind <= reflect.Float64
}
