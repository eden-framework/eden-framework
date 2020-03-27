package validatetpl

import (
	"fmt"
	"reflect"
)

func NewRangeValidateUint8(min uint8, max uint8, exc_min, exc_max bool) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Uint8 {
			value := uint8(reflect.ValueOf(v).Uint())
			if value < min || value > max || (value == min && exc_min) || (value == max && exc_max) {
				lb, rb := getBrackets(exc_min, exc_max)
				return false, fmt.Sprintf(INT_VALUE_NOT_IN_RANGE, lb, min, max, rb, value)
			}
			return true, ""
		}
		return false, TYPE_NOT_UINT8
	}
}
func NewEnumValidateUint8(enum_values ...uint8) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Uint8 {
			value := uint8(reflect.ValueOf(v).Uint())
			for _, enum_value := range enum_values {
				if value == uint8(enum_value) {
					return true, ""
				}
			}
			return false, fmt.Sprintf(INT_VALUE_NOT_IN_ENUM, enum_values, value)

		}
		return false, TYPE_NOT_UINT8
	}
}
func NewRangeValidateInt8(min int8, max int8, exc_min, exc_max bool) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Int8 {
			value := int8(reflect.ValueOf(v).Int())
			if value < min || value > max || (value == min && exc_min) || (value == max && exc_max) {
				lb, rb := getBrackets(exc_min, exc_max)
				return false, fmt.Sprintf(INT_VALUE_NOT_IN_RANGE, lb, min, max, rb, value)
			}
			return true, ""
		}
		return false, TYPE_NOT_INT8
	}
}
func NewEnumValidateInt8(enum_values ...int8) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Int8 {
			value := int8(reflect.ValueOf(v).Int())
			for _, enum_value := range enum_values {
				if value == int8(enum_value) {
					return true, ""
				}
			}
			return false, fmt.Sprintf(INT_VALUE_NOT_IN_ENUM, enum_values, value)
		}
		return false, TYPE_NOT_INT8
	}
}
