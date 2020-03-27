package validatetpl

import (
	"sync"
)

type ValidateFnMap struct {
	sync.Map
}

var validateFns = ValidateFnMap{}

func (m *ValidateFnMap) Store(name string, validateFn func(v interface{}) (bool, string)) *ValidateFnMap {
	m.Map.Store(name, validateFn)
	return m
}

func (m *ValidateFnMap) Load(name string) (validateFunc func(v interface{}) (bool, string), ok bool) {
	v, ok := m.Map.Load(name)
	if ok {
		validateFunc = v.(func(v interface{}) (bool, string))
	}
	return
}

func AddValidateFunc(name string, validateFunc func(v interface{}) (bool, string)) {
	validateFns.Store(name, validateFunc)
}

func GetValidateFunc(name string) (validateFunc func(v interface{}) (bool, string), ok bool) {
	return validateFns.Load(name)
}

func init() {
	AddValidateFunc("@phoneNo", ValidatePhone)
	AddValidateFunc("@phoneNoOrEmpty", ValidatePhoneOrEmpty)
	AddValidateFunc("@idCardNo", ValidateIDCardNo)
	AddValidateFunc("@idCardNoOrEmpty", ValidateIDCardNoOrEmpty)
	AddValidateFunc("@ipv4", ValidateIPv4)
	AddValidateFunc("@ipv4OrEmpty", ValidateIPv4OrEmpty)
	AddValidateFunc("@httpUrl", ValidateHttpUrl)
	AddValidateFunc("@httpUrlOrEmpty", ValidateHttpUrlOrEmpty)
	AddValidateFunc("@email", ValidateEmail)
	AddValidateFunc("@emailOrEmpty", ValidateEmailOrEmpty)
	AddValidateFunc("@password", ValidatePassword)
	AddValidateFunc("@passwordOrEmpty", ValidatePasswordOrEmpty)
	AddValidateFunc("@zipZh", ValidateZipZh)
	AddValidateFunc("@zipZhOrEmpty", ValidateZipZhOrEmpty)
	AddValidateFunc("@unitySocialCreditCode", ValidateUnitySocialCreditCode)
	AddValidateFunc("@unitySocialCreditCodeOrEmpty", ValidateUnitySocialCreditCodeOrEmpty)
	AddValidateFunc("@businessLicense", ValidateBusinessLicense)
	AddValidateFunc("@businessLicenseOrEmpty", ValidateBusinessLicenseOrEmpty)
	AddValidateFunc("@orgInsitituteCode", ValidateOrgInsitituteCode)
	AddValidateFunc("@orgInsitituteCodeOrEmpty", ValidateOrgInsitituteCodeOrEmpty)
	AddValidateFunc("@mysqlDataType", ValidateMySQLDataType)
	AddValidateFunc("@mysqlDataTypeOrEmpty", ValidateMySQLDataTypeOrEmpty)
	AddValidateFunc("@crontab", ValidateCrontab)
	AddValidateFunc("@crontabOrEmpty", ValidateCrontabOrEmpty)
	AddValidateFunc("@bankCard", ValidateBankCard)
	AddValidateFunc("@plateNo", ValidatePlateNo)
}
