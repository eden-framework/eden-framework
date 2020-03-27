package validatetpl

import (
	"regexp"
)

const (
	InvalidEmailType      = "邮箱类型错误"
	InvalidEmailValue     = "无效的邮箱"
	EmailMaxLength    int = 128
)

var (
	emailValidRegx = regexp.MustCompile(`^\w[-\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\.)+[A-Za-z]{2,14}$`)
)

func ValidateEmail(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidEmailType
	}

	length := len(s)
	if length <= 0 || length > EmailMaxLength {
		return false, InvalidEmailValue
	}

	if !emailValidRegx.MatchString(s) {
		return false, InvalidEmailValue
	}
	return true, ""
}

func ValidateEmailOrEmpty(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidEmailType
	}

	length := len(s)
	if length > EmailMaxLength {
		return false, InvalidEmailValue
	}

	if s != "" && !emailValidRegx.MatchString(s) {
		return false, InvalidEmailValue
	}
	return true, ""
}
