package transform

import (
	"encoding"
	"encoding/json"
	"fmt"
	"go/ast"
	"reflect"

	"github.com/eden-framework/eden-framework/pkg/validate"
	"github.com/eden-framework/reflectx"
	"github.com/eden-framework/strings"
)

type Validator func() (isValid bool, msg string)

func NewScanner() *Scanner {
	return &Scanner{}
}

type Scanner struct {
	walker    PathWalker
	errMsgMap ErrMsgMap
}

func (vs *Scanner) Validate(rv reflect.Value, tpe reflect.Type) (bool, ErrMsgMap) {
	vs.marshalAndValidate(rv, tpe, "", false, "", "")
	return len(vs.errMsgMap) == 0, vs.errMsgMap
}

func (vs *Scanner) setErrMsg(path string, msg string) {
	if vs.errMsgMap == nil {
		vs.errMsgMap = ErrMsgMap{}
	}
	vs.errMsgMap[path] = msg
}

func (vs *Scanner) marshalAndValidate(rv reflect.Value, tpe reflect.Type, defaultValue string, required bool, tagValidate string, tagErrMsg string) {
	v := rv.Interface()
	if rv.Kind() != reflect.Ptr {
		v = rv.Addr().Interface()
	}

	if _, ok := v.(encoding.TextUnmarshaler); ok {
		errMsg := MarshalAndValidate(rv, tpe, defaultValue, required, tagValidate, tagErrMsg)
		if errMsg != "" {
			vs.setErrMsg(vs.walker.String(), errMsg)
		}
		return
	}

	if _, ok := v.(json.Unmarshaler); ok {
		errMsg := MarshalAndValidate(rv, tpe, defaultValue, required, tagValidate, tagErrMsg)
		if errMsg != "" {
			vs.setErrMsg(vs.walker.String(), errMsg)
		}
		return
	}

	tpe = reflectx.IndirectType(tpe)

	switch tpe.Kind() {
	case reflect.Struct:
		if rv.Kind() == reflect.Ptr {
			errMsg := MarshalAndValidate(rv, tpe, defaultValue, required, tagValidate, tagErrMsg)
			if errMsg != "" {
				vs.setErrMsg(vs.walker.String(), errMsg)
				return
			}

			if rv.IsNil() {
				return
			}
		}

		rv = reflectx.Indirect(rv)

		for i := 0; i < tpe.NumField(); i++ {
			field := tpe.Field(i)
			if !ast.IsExported(field.Name) {
				continue
			}

			jsonTag, exists, jsonFlags := GetTagJSON(&field)
			if (exists && jsonTag != "-") || field.Anonymous {
				if !field.Anonymous {
					vs.walker.Enter(GetStructFieldDisplayName(&field))
				}

				tagValidate, _ := GetTagValidate(&field)
				defaultValue, hasDefaultValue := GetTagDefault(&field)
				tagErrMsg, _ := GetTagErrMsg(&field)

				required := true
				if hasOmitempty, ok := jsonFlags["omitempty"]; ok {
					required = !hasOmitempty
				} else {
					// todo don't use non-default as required
					required = !hasDefaultValue
				}

				vs.marshalAndValidate(rv.Field(i), field.Type, defaultValue, required, tagValidate, tagErrMsg)

				if rv.NumMethod() > 0 {
					validateHook := rv.MethodByName(fmt.Sprintf("Validate%s", field.Name))
					if validateHook.IsValid() {
						if validateFn, ok := validateHook.Interface().(func() string); ok {
							msg := validateFn()
							if msg != "" {
								vs.setErrMsg(vs.walker.String(), msg)
							}
						}
					}
				}

				if !field.Anonymous {
					vs.walker.Exit()
				}
			}
		}
	case reflect.Slice, reflect.Array:
		if tagValidate != "" {
			isValid, msg := validate.ValidateItem(tagValidate, rv.Interface(), tagErrMsg)
			if !isValid {
				vs.setErrMsg(vs.walker.String(), msg)
			}
		}
		for i := 0; i < rv.Len(); i++ {
			vs.walker.Enter(i)
			vs.marshalAndValidate(rv.Index(i), tpe.Elem(), "", false, "", tagErrMsg)
			vs.walker.Exit()
		}
	default:
		errMsg := MarshalAndValidate(rv, tpe, defaultValue, required, tagValidate, tagErrMsg)
		if errMsg != "" {
			vs.setErrMsg(vs.walker.String(), errMsg)
		}
	}
}

var ErrMsgForRequired = "缺失必填字段"

func MarshalAndValidate(
	rv reflect.Value, tpe reflect.Type,
	defaultValue string, isRequired bool,
	tagValidate string, tagErrMsg string,
) string {
	isPtr := rv.Kind() == reflect.Ptr

	if isPtr {
		if rv.IsNil() {
			if isRequired {
				return ErrMsgForRequired
			}

			// when not required，should set value
			if tpe.Kind() != reflect.Struct && defaultValue != "" && rv.CanSet() {
				rv.Set(reflect.New(reflectx.IndirectType(tpe)))
				rv = reflect.Indirect(rv)

				err := str.ConvertFromStr(defaultValue, rv)
				if err != nil {
					return fmt.Sprintf("%s can't set wrong default value %s", rv.Type().Name(), defaultValue)
				}
			}
		}

		if tagValidate != "" {
			rv = reflect.Indirect(rv)
			isValid, msg := validate.ValidateItem(tagValidate, rv.Interface(), tagErrMsg)
			if !isValid {
				return msg
			}
		}

		return ""
	}

	rv = reflect.Indirect(rv)
	isEmptyValue := reflectx.IsEmptyValue(rv)

	if isEmptyValue && isRequired {
		return ErrMsgForRequired
	}

	// only empty value can set default
	if isEmptyValue && defaultValue != "" && rv.CanSet() {
		err := str.ConvertFromStr(defaultValue, rv)
		if err != nil {
			return fmt.Sprintf("%s can't set wrong default value %s", rv.Type().Name(), defaultValue)
		}
	}

	if tagValidate != "" {
		isValid, msg := validate.ValidateItem(tagValidate, rv.Interface(), tagErrMsg)
		if !isValid {
			return msg
		}
	}

	return ""
}

type PathWalker struct {
	path []interface{}
}

func (pw *PathWalker) Enter(i interface{}) {
	pw.path = append(pw.path, i)
}

func (pw *PathWalker) Exit() {
	pw.path = pw.path[:len(pw.path)-1]
}

func (pw *PathWalker) Paths() []interface{} {
	return pw.path
}

func (pw *PathWalker) String() string {
	pathString := ""
	for i := 0; i < len(pw.path); i++ {
		switch pw.path[i].(type) {
		case string:
			if pathString != "" {
				pathString += "."
			}
			pathString += pw.path[i].(string)
		case int:
			pathString += fmt.Sprintf("[%d]", pw.path[i].(int))
		}
	}
	return pathString
}
