package validatetpl

import (
	"fmt"
	"reflect"
	"strings"
)

func NewRangeValidateFloat32(min, max float32, exc_min, exc_max bool) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Float32 {
			value := float32(reflect.ValueOf(v).Float())
			if value < min || value > max || (value == min && exc_min) || (value == max && exc_max) {
				lb, rb := getBrackets(exc_min, exc_max)
				return false, fmt.Sprintf(FLOAT_VALUE_NOT_IN_RANGE, lb, min, max, rb, value)
			}
			return true, ""
		}
		return false, TYPE_NOT_FLOAT32
	}
}
func NewRangeValidateFloat64(min, max float64, exc_min, exc_max bool) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Float64 {
			value := float64(reflect.ValueOf(v).Float())
			if value < min || value > max || (value == min && exc_min) || (value == max && exc_max) {
				lb, rb := getBrackets(exc_min, exc_max)
				return false, fmt.Sprintf(FLOAT_VALUE_NOT_IN_RANGE, lb, min, max, rb, value)
			}
			return true, ""
		}
		return false, TYPE_NOT_FLOAT64
	}
}
func NewEnumValidateFloat32(enum_values ...float32) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Float32 {
			value := float32(reflect.ValueOf(v).Float())
			for _, enum_value := range enum_values {
				if value == float32(enum_value) {
					return true, ""
				}
			}
			return false, fmt.Sprintf(FLOAT_VALUE_NOT_IN_ENUM, enum_values, value)
		}
		return false, TYPE_NOT_FLOAT32
	}
}
func NewEnumValidateFloat64(enum_values ...float64) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Float64 {
			value := float64(reflect.ValueOf(v).Float())
			for _, enum_value := range enum_values {
				if value == float64(enum_value) {
					return true, ""
				}
			}
			return false, fmt.Sprintf(FLOAT_VALUE_NOT_IN_ENUM, enum_values, value)
		}
		return false, TYPE_NOT_FLOAT64
	}
}

func NewDecimalValidateFloat32(total_len, decimal_len int) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Float32 {
			value := float32(reflect.ValueOf(v).Float())
			digits := strings.Split(fmt.Sprintf("%v", value), ".")
			// 避免小数部分为0
			if len(digits) == 1 {
				digits = append(digits, "0")
			}
			// 正数部分最大长度
			integer_len := total_len - decimal_len
			// 整数部分，高位0不计算有效位
			integer_digits := strings.TrimLeft(digits[0], "0")
			// 小数部分，低位0不计算有效位
			decimal_digits := strings.TrimRight(digits[1], "0")
			if len(integer_digits) > integer_len || len(decimal_digits) > decimal_len {
				return false, fmt.Sprintf(FLOAT_VALUE_DIGIT_INVALID, decimal_len, total_len, value)
			}
			return true, ""
		}
		return false, TYPE_NOT_FLOAT32
	}
}

func NewDecimalValidateFloat64(total_len, decimal_len int) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if reflect.TypeOf(v).Kind() == reflect.Float64 {
			value := float64(reflect.ValueOf(v).Float())
			digits := strings.Split(fmt.Sprintf("%v", value), ".")
			// 避免小数部分为0
			if len(digits) == 1 {
				digits = append(digits, "0")
			}
			// 正数部分最大长度
			integer_len := total_len - decimal_len
			// 整数部分，高位0不计算有效位
			integer_digits := strings.TrimLeft(digits[0], "0")
			// 小数部分，低位0不计算有效位
			decimal_digits := strings.TrimRight(digits[1], "0")
			if len(integer_digits) > integer_len || len(decimal_digits) > decimal_len {
				return false, fmt.Sprintf(FLOAT_VALUE_DIGIT_INVALID, decimal_len, total_len, value)
			}
			return true, ""
		}
		return false, TYPE_NOT_FLOAT64
	}
}
