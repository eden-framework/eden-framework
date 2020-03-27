package status_error

import (
	"bytes"
	"fmt"
)

var StatusErrorCodes = StatusErrorCodeMap{}

type StatusErrorCodeMap map[int64]StatusError

func (statusErrorCodeMap StatusErrorCodeMap) Register(key string, code int64, msg string, desc string, canBeTalkError bool) {
	statusErrorCodeMap[code] = StatusError{
		Key:            key,
		Code:           code,
		Msg:            msg,
		Desc:           desc,
		CanBeErrorTalk: canBeTalkError,
	}
}

func (statusErrorCodeMap StatusErrorCodeMap) Merge(targetStatusErrorCodeMap StatusErrorCodeMap) StatusErrorCodeMap {
	for code, statusError := range targetStatusErrorCodeMap {
		statusErrorCodeMap[code] = statusError
	}
	return statusErrorCodeMap
}

func (statusErrorCodeMap StatusErrorCodeMap) String() string {
	buffer := bytes.Buffer{}
	for _, statusError := range statusErrorCodeMap {
		buffer.WriteString(statusError.String() + "\n")
	}
	return buffer.String()
}

type StatusErrorCode int64

func (code StatusErrorCode) Is(err error) bool {
	return FromError(err).Code == int64(code)
}

func (code StatusErrorCode) StatusError() *StatusError {
	statusError, ok := StatusErrorCodes[int64(code)]
	if !ok {
		panic(fmt.Errorf("%d is not registered to statusErrorCodes", int64(code)))
	}
	return &statusError
}

// deprecated
func (code StatusErrorCode) ToError() *StatusError {
	return code.StatusError()
}

func (code StatusErrorCode) Error() string {
	statusError := code.StatusError()
	return statusError.Error()
}

func DemoErr() error {
	return InvalidSecret
}
