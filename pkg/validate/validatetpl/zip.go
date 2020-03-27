package validatetpl

import (
	"regexp"
)

const (
	InvalidZipZhType  = "邮编类型错误"
	InvalidZipZhValue = "无效的邮编"
)

var (
	ZipZhValidRegx = regexp.MustCompile(`^[1-9]\d{5}$`)
)

func ValidateZipZh(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidZipZhType
	}

	if len(s) <= 0 {
		return false, InvalidZipZhValue
	}

	if !ZipZhValidRegx.MatchString(s) {
		return false, InvalidZipZhValue
	}
	return true, ""
}

func ValidateZipZhOrEmpty(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidZipZhType
	}

	if s != "" && !ZipZhValidRegx.MatchString(s) {
		return false, InvalidZipZhValue
	}
	return true, ""
}
