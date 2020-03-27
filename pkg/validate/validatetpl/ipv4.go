package validatetpl

import (
	"regexp"
)

const (
	InvalidIPv4Type  = "IPv4类型错误"
	InvalidIPv4Value = "无效的IPv4"
)

var (
	reIpAddress = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
)

func ValidateIPv4(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidIPv4Type
	}

	if !reIpAddress.MatchString(s) {
		return false, InvalidIPv4Value
	}
	return true, ""
}

func ValidateIPv4OrEmpty(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidIPv4Type
	}

	if s != "" && !reIpAddress.MatchString(s) {
		return false, InvalidIPv4Value
	}
	return true, ""
}
