package validatetpl

import (
	"fmt"
	"reflect"
)

func NewRangeValidateUint32(min uint32, max uint32, exc_min, exc_max bool) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Uint32 {
			value := uint32(reflect.ValueOf(v).Uint())
			if value < min || value > max || (value == min && exc_min) || (value == max && exc_max) {
				lb, rb := getBrackets(exc_min, exc_max)
				return false, fmt.Sprintf(INT_VALUE_NOT_IN_RANGE, lb, min, max, rb, value)
			}
			return true, ""
		}
		return false, TYPE_NOT_UINT32
	}
}
func NewRangeValidateInt32(min int32, max int32, exc_min, exc_max bool) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Int32 {
			value := int32(reflect.ValueOf(v).Int())
			if value < min || value > max || (value == min && exc_min) || (value == max && exc_max) {
				lb, rb := getBrackets(exc_min, exc_max)
				return false, fmt.Sprintf(INT_VALUE_NOT_IN_RANGE, lb, min, max, rb, value)
			}
			return true, ""
		}
		return false, TYPE_NOT_INT32
	}
}
func NewEnumValidateUint32(enum_values ...uint32) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Uint32 {
			value := uint32(reflect.ValueOf(v).Uint())
			for _, enum_value := range enum_values {
				if value == uint32(enum_value) {
					return true, ""
				}
			}
			return false, fmt.Sprintf(INT_VALUE_NOT_IN_ENUM, enum_values, value)
		}
		return false, TYPE_NOT_UINT32
	}
}
func NewEnumValidateInt32(enum_values ...int32) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Int32 {
			value := int32(reflect.ValueOf(v).Int())
			for _, enum_value := range enum_values {
				if value == int32(enum_value) {
					return true, ""
				}
			}
			return false, fmt.Sprintf(INT_VALUE_NOT_IN_ENUM, enum_values, value)
		}
		return false, TYPE_NOT_INT32
	}
}
