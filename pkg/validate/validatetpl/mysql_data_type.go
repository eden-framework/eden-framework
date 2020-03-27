package validatetpl

import (
	"regexp"
	"strings"
)

const (
	InvalidMySQLDataTypeType  = "MySQL数据类型错误"
	InvalidMySQLDataTypeValue = "无效的MySQL数据"
)

var (
	mysqlDataTypeRegexp = regexp.MustCompile(`^(tinyint\(8\)|smallint\(16\)|int\(32\)|bigint\(64\))( unsigned)?|varchar\(\d+\)$`)
)

func ValidateMySQLDataType(v interface{}) (bool, string) {
	s, ok := v.(string)
	s = strings.ToLower(s)
	if !ok {
		return false, InvalidMySQLDataTypeType
	}
	if !mysqlDataTypeRegexp.MatchString(s) {
		return false, InvalidMySQLDataTypeValue
	}
	return true, ""
}

func ValidateMySQLDataTypeOrEmpty(v interface{}) (bool, string) {
	s, ok := v.(string)
	s = strings.ToLower(s)
	if !ok {
		return false, InvalidMySQLDataTypeType
	}
	if s != "" && !mysqlDataTypeRegexp.MatchString(s) {
		return false, InvalidMySQLDataTypeValue
	}
	return true, ""
}
