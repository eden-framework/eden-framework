package validatetpl

import (
	"regexp"
)

const (
	InvalidPasswordType  = "密码类型错误"
	InvalidPasswordValue = "无效的密码"
)

var (
	passwordValidRegx = regexp.MustCompile(`^[0-9A-Za-z]{8,16}$`)
)

func ValidatePassword(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidPasswordType
	}

	if len(s) <= 0 {
		return false, InvalidPasswordValue
	}

	if !passwordValidRegx.MatchString(s) {
		return false, InvalidPasswordValue
	}
	return true, ""
}

func ValidatePasswordOrEmpty(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidPasswordType
	}

	if s != "" && !passwordValidRegx.MatchString(s) {
		return false, InvalidPasswordValue
	}
	return true, ""
}
