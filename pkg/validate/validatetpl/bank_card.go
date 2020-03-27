package validatetpl

import (
	"regexp"
)

const (
	InvalidBankCardType  = "银行卡类型错误"
	InvalidBankCardValue = "无效的银行卡"
)

var (
	bankCardRegexp = regexp.MustCompile(`^\d{12,19}$`)
)

func ValidateBankCard(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidBankCardType
	}

	if !bankCardRegexp.MatchString(s) {
		return false, InvalidBankCardValue
	}

	return true, ""
}
