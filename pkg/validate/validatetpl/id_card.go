package validatetpl

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	InvalidIDCardNoType  = "身份证号类型错误"
	InvalidIDCardNoValue = "无效的身份证号"
)

var weights []int64 = []int64{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

var expectedLastChars = "10X98765432"

var (
	id_card_regexp = regexp.MustCompile(`^\d{17}(\d|x|X)$`)
)

func ValidateIDCardNo(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidIDCardNoType
	}
	if !id_card_regexp.MatchString(s) {
		return false, InvalidIDCardNoValue
	}
	if !calculateAndCompareLastChar(s) {
		return false, InvalidIDCardNoValue
	}
	return true, ""
}

func ValidateIDCardNoOrEmpty(v interface{}) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, InvalidIDCardNoType
	}
	if s != "" && !id_card_regexp.MatchString(s) {
		return false, InvalidIDCardNoValue
	}
	if s != "" && !calculateAndCompareLastChar(s) {
		return false, InvalidIDCardNoValue
	}
	return true, ""
}

func calculateAndCompareLastChar(idCard string) bool {
	var sum int64
	var idCardLength = len(idCard)
	for i := 0; i < idCardLength-1; i++ {
		number, err := strconv.ParseInt(string(idCard[i]), 10, 64)
		if err != nil {
			return false
		}
		sum += number * weights[i]
	}
	n := sum % 11
	expectedLastChar := expectedLastChars[n]
	if string(expectedLastChar) != strings.ToUpper(string(idCard[idCardLength-1])) {
		return false
	}
	return true
}
