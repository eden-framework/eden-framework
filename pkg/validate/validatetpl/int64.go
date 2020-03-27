package validatetpl

import (
	"fmt"
	"reflect"
)

func NewRangeValidateUint64(min uint64, max uint64, exc_min, exc_max bool) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Uint64 {
			value := uint64(reflect.ValueOf(v).Uint())
			if value < min || value > max || (value == min && exc_min) || (value == max && exc_max) {
				lb, rb := getBrackets(exc_min, exc_max)
				return false, fmt.Sprintf(INT_VALUE_NOT_IN_RANGE, lb, min, max, rb, value)
			}
			return true, ""
		}
		return false, TYPE_NOT_UINT64
	}
}
func NewRangeValidateInt64(min int64, max int64, exc_min, exc_max bool) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Int64 {
			value := int64(reflect.ValueOf(v).Int())
			if value < min || value > max || (value == min && exc_min) || (value == max && exc_max) {
				lb, rb := getBrackets(exc_min, exc_max)
				return false, fmt.Sprintf(INT_VALUE_NOT_IN_RANGE, lb, min, max, rb, value)
			}
			return true, ""
		}
		return false, TYPE_NOT_INT64
	}
}
func NewEnumValidateUint64(enum_values ...uint64) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Uint64 {
			value := uint64(reflect.ValueOf(v).Uint())
			for _, enum_value := range enum_values {
				if value == uint64(enum_value) {
					return true, ""
				}
			}
			return false, fmt.Sprintf(INT_VALUE_NOT_IN_ENUM, enum_values, value)
		}
		return false, TYPE_NOT_UINT64
	}
}
func NewEnumValidateInt64(enum_values ...int64) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Int64 {
			value := int64(reflect.ValueOf(v).Int())
			for _, enum_value := range enum_values {
				if value == int64(enum_value) {
					return true, ""
				}
			}
			return false, fmt.Sprintf(INT_VALUE_NOT_IN_ENUM, enum_values, value)
		}
		return false, TYPE_NOT_INT64
	}
}
