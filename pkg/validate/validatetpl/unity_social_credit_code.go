/*
规则参考：https://zh.wikisource.org/zh-hans/GB_32100-2015_%E6%B3%95%E4%BA%BA%E5%92%8C%E5%85%B6%E4%BB%96%E7%BB%84%E7%BB%87%E7%BB%9F%E4%B8%80%E7%A4%BE%E4%BC%9A%E4%BF%A1%E7%94%A8%E4%BB%A3%E7%A0%81%E7%BC%96%E7%A0%81%E8%A7%84%E5%88%99
*/

package validatetpl

import (
	"fmt"
	"regexp"
)

const (
	InvalidUnitySocialCreditCodeType             = "社会信用代码类型错误"
	InvalidUnitySocialCreditCodeValue            = "无效的社会信用代码"
	UnitySocialCreditCodeValueVerifyCodeNotMatch = "社会信用代码校验码不匹配"
	UnitySocialCreditCodeValueWeightNotFound     = "社会信用代码权重值不存在"
)

var unity_social_credit_code_regexp *regexp.Regexp
var code2value map[string]int
var positionWeights []int

func init() {
	first := `(1|5|9|Y)`
	second := `(1|2|3|9)`
	thirdToEight := `(\d{6})`
	nineToEighteen := `([0-9A-Z^IOZSV]{10})`
	str := fmt.Sprintf(`^%s%s%s%s$`, first, second, thirdToEight, nineToEighteen)
	unity_social_credit_code_regexp = regexp.MustCompile(str)
	code2value = map[string]int{
		`0`: 0,
		`1`: 1,
		`2`: 2,
		`3`: 3,
		`4`: 4,
		`5`: 5,
		`6`: 6,
		`7`: 7,
		`8`: 8,
		`9`: 9,
		`A`: 10,
		`B`: 11,
		`C`: 12,
		`D`: 13,
		`E`: 14,
		`F`: 15,
		`G`: 16,
		`H`: 17,
		`J`: 18,
		`K`: 19,
		`L`: 20,
		`M`: 21,
		`N`: 22,
		`P`: 23,
		`Q`: 24,
		`R`: 25,
		`T`: 26,
		`U`: 27,
		`W`: 28,
		`X`: 29,
		`Y`: 30,
	}
	positionWeights = []int{1, 3, 9, 27, 19, 26, 16, 17, 20, 29, 25, 13, 8, 24, 10, 30, 28}
}

func CheckoutUnitySocialCreditCode(array string) (bool, string) {
	var r int
	for index, i := range array[:17] {
		c, ok := code2value[string(i)]
		if !ok {
			return false, UnitySocialCreditCodeValueWeightNotFound
		}
		r += c * positionWeights[index]
	}
	r = 31 - r%31
	if r == 31 {
		r = 0
	}
	if r == 30 {
		r = code2value[`Y`]
	}

	checkCode, ok := code2value[string(array[17])]
	if !ok {
		return false, UnitySocialCreditCodeValueWeightNotFound
	}

	if r != checkCode {
		return false, UnitySocialCreditCodeValueVerifyCodeNotMatch
	}

	return true, ""
}

func ValidateUnitySocialCreditCode(v interface{}) (bool, string) {
	array, ok := v.(string)
	if !ok {
		return false, InvalidUnitySocialCreditCodeType
	}
	if !unity_social_credit_code_regexp.MatchString(array) {
		return false, InvalidUnitySocialCreditCodeValue
	}

	// Note: 存在不能通过校验的统一社会信用码，
	// 目前只做正则验证即可
	// if ok, str := CheckoutUnitySocialCreditCode(array); !ok {
	// 	return ok, str
	// }

	return true, ""
}

func ValidateUnitySocialCreditCodeOrEmpty(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidUnitySocialCreditCodeType
	}
	if s != "" && !unity_social_credit_code_regexp.MatchString(s) {
		return false, InvalidUnitySocialCreditCodeValue
	}
	return true, ""
}
