package validatetpl

import (
	"regexp"
)

const (
	InvalidPlateNoType  = "车牌号类型错误"
	InvalidPlateNoValue = "无效的车牌号"
)

var (
	plateNoRegexp = regexp.MustCompile(`^(京|津|沪|渝|蒙|新|藏|宁|桂|港|澳|黑|吉|辽|晋|冀|青|鲁|豫|苏|皖|浙|闽|赣|湘|鄂|粤|琼|甘|陕|贵|云|川)[A-Z][0-9A-Z]{4,5}[0-9A-Z挂港澳学警领试]$`)
)

func ValidatePlateNo(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidPlateNoType
	}

	if s != "" && !plateNoRegexp.MatchString(s) {
		return false, InvalidPlateNoValue
	}

	return true, ""
}
