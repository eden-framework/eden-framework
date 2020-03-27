package validatetpl

import (
	"fmt"
	"reflect"
)

func NewRangeValidateUint16(min uint16, max uint16, exc_min, exc_max bool) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Uint16 {
			value := uint16(reflect.ValueOf(v).Uint())
			if value < min || value > max || (value == min && exc_min) || (value == max && exc_max) {
				lb, rb := getBrackets(exc_min, exc_max)
				return false, fmt.Sprintf(INT_VALUE_NOT_IN_RANGE, lb, min, max, rb, value)
			}
			return true, ""
		}
		return false, TYPE_NOT_UINT16
	}
}
func NewRangeValidateInt16(min int16, max int16, exc_min, exc_max bool) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Int16 {
			value := int16(reflect.ValueOf(v).Int())
			if value < min || value > max || (value == min && exc_min) || (value == max && exc_max) {
				lb, rb := getBrackets(exc_min, exc_max)
				return false, fmt.Sprintf(INT_VALUE_NOT_IN_RANGE, lb, min, max, rb, value)
			}
			return true, ""
		}
		return false, TYPE_NOT_INT16
	}
}
func NewEnumValidateUint16(enum_values ...uint16) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Uint16 {
			value := uint16(reflect.ValueOf(v).Uint())
			for _, enum_value := range enum_values {
				if value == uint16(enum_value) {
					return true, ""
				}
			}
			return false, fmt.Sprintf(INT_VALUE_NOT_IN_ENUM, enum_values, value)
		}
		return false, TYPE_NOT_UINT16
	}
}
func NewEnumValidateInt16(enum_values ...int16) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Int16 {
			value := int16(reflect.ValueOf(v).Int())
			for _, enum_value := range enum_values {
				if value == int16(enum_value) {
					return true, ""
				}
			}
			return false, fmt.Sprintf(INT_VALUE_NOT_IN_ENUM, enum_values, value)

		}
		return false, TYPE_NOT_INT16
	}
}
