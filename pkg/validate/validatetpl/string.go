package validatetpl

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

func NewValidateChar(min int, max int) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if stringer, ok := v.(fmt.Stringer); ok {
			v = stringer.String()
		}
		if value, ok := v.(string); ok {
			word_len := utf8.RuneCount([]byte(value))
			if word_len < min || word_len > max {
				return false, fmt.Sprintf(STRING_CHARS_NOT_IN_RANGE, min, max, word_len)
			}
			return true, ""
		}
		return false, TYPE_NOT_STRING
	}
}

func NewValidateStringLength(min_len int, max_len int) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if stringer, ok := v.(fmt.Stringer); ok {
			v = stringer.String()
		}
		if value, ok := v.(string); ok {
			str_len := len(value)
			if str_len < min_len || (str_len > max_len && max_len != STRING_UNLIMIT_VALUE) {
				return false, fmt.Sprintf(STRING_LENGHT_NOT_IN_RANGE, min_len, max_len, str_len)
			}
			return true, ""
		}
		return false, TYPE_NOT_STRING
	}
}
func NewValidateStringRegExp(reg *regexp.Regexp) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if stringer, ok := v.(fmt.Stringer); ok {
			v = stringer.String()
		}
		if value, ok := v.(string); ok {
			if !reg.MatchString(value) {
				return false, fmt.Sprintf(STRING_NOT_MATCH_REGEXP, reg.String())
			}
			return true, ""
		}
		return false, TYPE_NOT_STRING
	}
}

func NewEnumValidateString(enum_values ...string) func(v interface{}) (bool, string) {
	return func(v interface{}) (bool, string) {
		if stringer, ok := v.(fmt.Stringer); ok {
			v = stringer.String()
		}
		value := fmt.Sprintf("%v", v)
		for _, enum_value := range enum_values {
			if value == enum_value {
				return true, ""
			}
		}
		return false, fmt.Sprintf(STRING_VALUE_NOT_IN_ENUM, enum_values, value)
	}
}
