package validatetpl

import (
	"regexp"
)

const (
	InvalidOrgInsitituteCodeType  = "组织机构代码类型错误"
	InvalidOrgInsitituteCodeValue = "无效的组织机构代码"
)

var (
	org_insititute_code_regexp = regexp.MustCompile(`^\w{8}-\w$`)
)

func ValidateOrgInsitituteCode(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidOrgInsitituteCodeType
	}
	if !org_insititute_code_regexp.MatchString(s) {
		return false, InvalidOrgInsitituteCodeValue
	}
	return true, ""
}

func ValidateOrgInsitituteCodeOrEmpty(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidOrgInsitituteCodeType
	}
	if s != "" && !org_insititute_code_regexp.MatchString(s) {
		return false, InvalidOrgInsitituteCodeValue
	}
	return true, ""
}
