package reflect

import (
	"errors"
	"fmt"
	"reflect"
)

// processStructFields visits accessible struct fields with optional filtering.
func processStructFields(value any, filter func(field reflect.Value) bool, process func(fieldName string, field reflect.Value)) error {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return errors.New("inputStruct must be a struct or a pointer to a struct")
	}

	for structField, field := range val.Fields() {
		// Skip inaccessible fields.
		if !field.CanInterface() {
			continue
		}
		if filter == nil || filter(field) {
			process(structField.Name, field)
		}
	}
	return nil
}

// AllFieldName returns accessible struct field names.
func AllFieldName(value any) ([]string, error) {
	var allFields []string
	err := processStructFields(value,
		nil, // No field filter.
		func(fieldName string, field reflect.Value) {
			allFields = append(allFields, fieldName)
		})
	return allFields, err
}

// AllFieldValue returns accessible struct fields as a map.
func AllFieldValue(value any) (map[string]any, error) {
	allValue := make(map[string]any)
	err := processStructFields(value,
		nil, // No field filter.
		func(fieldName string, field reflect.Value) {
			allValue[fieldName] = field.Interface()
		})
	return allValue, err
}

// NonZeroFieldName returns accessible nonzero field names; empty collections count as zero.
func NonZeroFieldName(value any) ([]string, error) {
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

// NonZeroFieldValue returns accessible nonzero fields; empty collections count as zero.
func NonZeroFieldValue(value any) (map[string]any, error) {
	nonZeroValue := make(map[string]any)
	err := processStructFields(value,
		func(field reflect.Value) bool {
			return !isZeroValue(field)
		},
		func(fieldName string, field reflect.Value) {
			nonZeroValue[fieldName] = field.Interface()
		})
	return nonZeroValue, err
}

// DeepCopy recursively copies supported values without cycle detection.
// Unexported struct fields can panic; channels and functions are not copied.
func DeepCopy(src any) any {
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
		for structField, field := range src.Fields() {
			deepCopyRecursive(field, dst.FieldByIndex(structField.Index))
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
		return true // Treat invalid or inaccessible fields as zero.
	}
	// Apply the existing kind-specific zero rules.
	zero := reflect.Zero(v.Type())
	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		return v.Len() == 0 // Empty collections count as zero, including arrays and channels.
	case reflect.Ptr, reflect.Interface:
		return v.IsNil() // Nil pointers and interfaces count as zero.
	default:
		// Compare other kinds with their typed zero value.
		return reflect.DeepEqual(v.Interface(), zero.Interface())
	}
}

// SetFieldValue assigns convertible values to a struct pointer, optionally rejecting missing fields.
// Nil map values are unsupported and can panic.
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
			continue // Ignore missing fields unless strict checking was requested.
		}
		if !field.CanSet() {
			return fmt.Errorf("field %s cannot be set (may be unexported)", fieldName)
		}
		value := reflect.ValueOf(val)
		fieldType := field.Type()

		// Apply conversions supported by reflection.
		if value.Type().ConvertibleTo(fieldType) {
			field.Set(value.Convert(fieldType))
			continue
		}

		// Try numeric-kind conversion as a fallback.
		if converted, ok := tryConvertNumber(value, fieldType); ok {
			field.Set(converted)
			continue
		}

		return fmt.Errorf("cannot assign value of type %s to field %s (type %s)", value.Type(), fieldName, fieldType)
	}

	return nil
}

// tryConvertNumber converts between numeric kinds through float64.
func tryConvertNumber(val reflect.Value, targetType reflect.Type) (reflect.Value, bool) {
	if !val.IsValid() {
		return reflect.Value{}, false
	}
	kind := val.Kind()
	tKind := targetType.Kind()
	if isNumericKind(kind) && isNumericKind(tKind) {
		// Convert through float64.
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
		// Convert to the destination type.
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

// GetFieldValue returns an accessible named struct field or a validation error.
func GetFieldValue(value any, fieldName string) (any, error) {
	if value == nil {
		return nil, fmt.Errorf("input must be a struct or pointer to struct")
	}
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input must be a struct or pointer to struct")
	}
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("field %s does not exist in struct", fieldName)
	}
	if !field.CanInterface() {
		return nil, fmt.Errorf("field %s is unexported", fieldName)
	}
	return field.Interface(), nil
}
