package copy

import (
	"github.com/acexy/golang-toolkit/util/json"
	"github.com/jinzhu/copier"
	"reflect"
	"time"
)

// Copy 拓展copier功能
func Copy(toValue interface{}, fromValue interface{}) error {
	err := copier.Copy(toValue, fromValue)
	if err != nil {
		return err
	}
	return convertExtendsFields(toValue, fromValue)
}

func convertExtendsFields(dst interface{}, src interface{}) error {
	dstValue := reflect.ValueOf(dst).Elem()
	srcValue := reflect.ValueOf(src).Elem()

	for i := 0; i < dstValue.NumField(); i++ {
		dstField := dstValue.Field(i)
		srcField := srcValue.Field(i)

		if dstField.CanSet() && srcField.IsValid() && !srcField.IsZero() {
			if !srcField.IsZero() {
				if srcField.Type() == reflect.TypeOf(time.Time{}) && dstField.Type() == reflect.TypeOf(json.Timestamp{}) {
					dstField.Set(reflect.ValueOf(json.Timestamp{Time: srcField.Interface().(time.Time)}))
				} else if srcField.Type() == reflect.PtrTo(reflect.TypeOf(time.Time{})) && dstField.Type() == reflect.PtrTo(reflect.TypeOf(json.Timestamp{})) {
					if !srcField.IsNil() {
						timeValue := *srcField.Interface().(*time.Time)
						dstField.Set(reflect.ValueOf(&json.Timestamp{Time: timeValue}))
					}
				} else if srcField.Type() == reflect.PtrTo(reflect.TypeOf(time.Time{})) && dstField.Type() == reflect.TypeOf(json.Timestamp{}) {
					if !srcField.IsNil() {
						timeValue := *srcField.Interface().(*time.Time)
						dstField.Set(reflect.ValueOf(json.Timestamp{Time: timeValue}))
					}
				} else if srcField.Type() == reflect.TypeOf(time.Time{}) && dstField.Type() == reflect.PtrTo(reflect.TypeOf(json.Timestamp{})) {
					timeValue := srcField.Interface().(time.Time)
					dstField.Set(reflect.ValueOf(&json.Timestamp{Time: timeValue}))
				} else if srcField.Type() == reflect.TypeOf(json.Timestamp{}) && dstField.Type() == reflect.TypeOf(time.Time{}) {
					dstField.Set(reflect.ValueOf(srcField.Interface().(json.Timestamp).Time))
				} else if srcField.Type() == reflect.PtrTo(reflect.TypeOf(json.Timestamp{})) && dstField.Type() == reflect.PtrTo(reflect.TypeOf(time.Time{})) {
					if !srcField.IsNil() {
						timestampValue := srcField.Interface().(*json.Timestamp)
						dstField.Set(reflect.ValueOf(&timestampValue.Time))
					}
				} else if srcField.Type() == reflect.TypeOf(json.Timestamp{}) && dstField.Type() == reflect.PtrTo(reflect.TypeOf(time.Time{})) {
					timeValue := srcField.Interface().(json.Timestamp).Time
					dstField.Set(reflect.ValueOf(&timeValue))
				} else if srcField.Type() == reflect.PtrTo(reflect.TypeOf(json.Timestamp{})) && dstField.Type() == reflect.TypeOf(time.Time{}) {
					if !srcField.IsNil() {
						timeValue := srcField.Interface().(*json.Timestamp).Time
						dstField.Set(reflect.ValueOf(timeValue))
					}
				}
			}
		}
	}
	return nil
}
