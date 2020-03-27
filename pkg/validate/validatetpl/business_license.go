package validatetpl

import (
	"regexp"
)

const (
	InvalidBusinessLicenseType  = "营业执照号类型错误"
	InvalidBusinessLicenseValue = "无效的营业执照号"
)

var (
	business_license_regexp = regexp.MustCompile(`^\d{15}$`)
)

func ValidateBusinessLicense(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidBusinessLicenseType
	}
	if !business_license_regexp.MatchString(s) {
		return false, InvalidBusinessLicenseValue
	}
	return true, ""
}

func ValidateBusinessLicenseOrEmpty(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidBusinessLicenseType
	}
	if s != "" && !business_license_regexp.MatchString(s) {
		return false, InvalidBusinessLicenseValue
	}
	return true, ""
}
