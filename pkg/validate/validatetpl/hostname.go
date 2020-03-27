package validatetpl

import (
	"regexp"
)

const (
	InvalidHostnameType  = "Hostname 类型错误"
	InvalidHostnameValue = "Hostname 类型值错误"
)

func init() {
	AddValidateFunc("@hostname", ValidateHostname)
}

var (
	reHostname      = regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	reIpAddressLike = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)\.(\d+)$`)
)

func ValidateHostname(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidHostnameType
	}
	if !reHostname.MatchString(s) {
		return false, InvalidHostnameValue
	}
	if reIpAddressLike.MatchString(s) {
		if !reIpAddress.MatchString(s) {
			return false, InvalidHostnameValue
		}
	}
	return true, ""
}
