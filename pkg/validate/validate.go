package validate

import (
	"github.com/eden-framework/eden-framework/pkg/validate/validatetpl"
)

type ValidateFn func(v interface{}) (bool, string)

func AddValidateFunc(name string, validateFunc ValidateFn) {
	validatetpl.AddValidateFunc(name, validateFunc)
}

func GetValidateFunc(name string) (validateFunc ValidateFn, ok bool) {
	return validatetpl.GetValidateFunc(name)
}

func ValidateItem(tagValidate string, v interface{}, tagErrMsg ...string) (valid bool, message string) {
	var validateFn ValidateFn
	var ok bool
	if validateFn, ok = GetValidateFunc(tagValidate); !ok {
		validateFn = validatetpl.GenerateValidateFuncByTag(tagValidate)
		if validateFn != nil {
			AddValidateFunc(tagValidate, validateFn)
		}
	}
	if validateFn != nil {
		valid, message = validateFn(v)
		if !valid {
			if len(tagErrMsg) > 0 && tagErrMsg[0] != "" {
				message = tagErrMsg[0]
			}
		}
		return
	}
	return true, ""
}
