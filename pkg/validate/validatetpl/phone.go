package validatetpl

import (
	"regexp"
)

const (
	InvalidPhoneNoType  = "手机号类型错误"
	InvalidPhoneNoValue = "无效的手机号"
)

var (
	phoneValidRegx = regexp.MustCompile("^1[0-9]{10}$")
)

func ValidatePhone(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidPhoneNoType
	}

	if len(s) <= 0 {
		return false, InvalidPhoneNoValue
	}

	if !phoneValidRegx.MatchString(s) {
		return false, InvalidPhoneNoValue
	}
	return true, ""
}

func ValidatePhoneOrEmpty(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidPhoneNoType
	}

	if s != "" && !phoneValidRegx.MatchString(s) {
		return false, InvalidPhoneNoValue
	}
	return true, ""
}
