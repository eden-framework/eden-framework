package transform

import (
	"sort"

	"github.com/eden-framework/eden-framework/pkg/courier/status_error"
)

type ErrMsgMap map[string]string

func (errMsgMap ErrMsgMap) Set(keyPath string, errMsg string) {
	errMsgMap[keyPath] = errMsg
}

func (errMsgMap ErrMsgMap) Merge(otherErrMsgMap ErrMsgMap) ErrMsgMap {
	if otherErrMsgMap == nil {
		return errMsgMap
	}
	for field, msg := range otherErrMsgMap {
		errMsgMap[field] = msg
	}
	return errMsgMap
}

func (errMsgMap ErrMsgMap) ErrorFieldsIn(in string, parentField string) (errFields status_error.ErrorFields) {
	if len(errMsgMap) == 0 {
		return
	}

	for field := range errMsgMap {
		if parentField != "" {
			field = parentField + "." + field
		}
		errFields = append(errFields, status_error.NewErrorField(in, field, errMsgMap[field]))
	}

	sort.Sort(errFields)
	return
}
