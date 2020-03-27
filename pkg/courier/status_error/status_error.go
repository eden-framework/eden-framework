package status_error

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

// code, key, msg, desc, canBeTalkError
var RegexpStatusError = regexp.MustCompile(`@httpError\(([0-9]+),(.+),"(.+)?","(.+)?",(false|true)\);`)

func ParseString(s string) *StatusError {
	matched := RegexpStatusError.FindStringSubmatch(s)
	if len(matched) != 6 {
		return nil
	}

	code, _ := strconv.ParseInt(matched[1], 10, 64)
	canBeTalkErr, _ := strconv.ParseBool(matched[5])

	return &StatusError{
		Code:           code,
		Key:            matched[2],
		Msg:            matched[3],
		Desc:           matched[4],
		CanBeErrorTalk: canBeTalkErr,
	}
}

func FromError(err error) *StatusError {
	if err == nil {
		return nil
	}

	if statusErrCode, ok := err.(StatusErrorCode); ok {
		return statusErrCode.StatusError()
	}

	if statusError, ok := err.(*StatusError); ok {
		return statusError
	}

	return UnknownError.StatusError().WithDesc(err.Error())
}

func (statusErr *StatusError) String() string {
	return fmt.Sprintf(`@httpError(%d,%s,"%s","%s",%v);`, statusErr.Code, statusErr.Key, statusErr.Msg, statusErr.Desc, statusErr.CanBeErrorTalk)
}

type StatusError struct {
	// 错误 Key
	Key string `json:"key" xml:"key"`
	// 错误代码
	Code int64 `json:"code" xml:"code"`
	// 错误信息
	Msg string `json:"msg" xml:"msg"`
	// 详细描述
	Desc string `json:"desc" xml:"desc"`
	// 是否能作为错误话术
	CanBeErrorTalk bool `json:"canBeTalkError" xml:"canBeTalkError"`
	// 错误溯源
	Source []string `json:"source" xml:"source"`
	// 请求 ID
	ID string `json:"id" xml:"id"`
	// 出错字段
	ErrorFields ErrorFields `json:"errorFields" xml:"errorFields"`
}

func (statusErr *StatusError) Is(err error) bool {
	return FromError(statusErr).Code == statusErr.Code
}

func (statusErr *StatusError) Error() string {
	return fmt.Sprintf("%v[%s][%d][%s%s] %s", statusErr.Source, statusErr.Key, statusErr.Code, statusErr.Msg, statusErr.ErrorFields, statusErr.Desc)
}

func (statusErr *StatusError) Status() int {
	strCode := fmt.Sprintf("%d", statusErr.Code)
	if len(strCode) < 3 {
		return 0
	}
	status, _ := strconv.Atoi(strCode[:3])
	return status
}

// deprecated
func (statusErr *StatusError) HttpCode() int {
	return statusErr.Status()
}

func (statusErr StatusError) WithErrTalk() *StatusError {
	statusErr.CanBeErrorTalk = true
	return &statusErr
}

func (statusErr StatusError) WithMsg(msg string) *StatusError {
	statusErr.Msg = msg
	return &statusErr
}

func (statusErr StatusError) WithoutErrTalk() *StatusError {
	statusErr.CanBeErrorTalk = false
	return &statusErr
}

func (statusErr StatusError) WithDesc(desc string) *StatusError {
	statusErr.Desc = desc
	return &statusErr
}

func (statusErr StatusError) WithID(id string) *StatusError {
	statusErr.ID = id
	return &statusErr
}

func (statusErr StatusError) WithSource(sourceName string) *StatusError {
	length := len(statusErr.Source)
	if length == 0 || statusErr.Source[length-1] != sourceName {
		statusErr.Source = append(statusErr.Source, sourceName)
	}
	return &statusErr
}

func (statusErr StatusError) WithErrorField(in string, field string, msg string) *StatusError {
	statusErr.ErrorFields = append(statusErr.ErrorFields, NewErrorField(in, field, msg))
	return &statusErr
}

func (statusErr StatusError) WithErrorFields(errorFields ...*ErrorField) *StatusError {
	statusErr.ErrorFields = append(statusErr.ErrorFields, errorFields...)
	return &statusErr
}

func NewErrorField(in string, field string, msg string) *ErrorField {
	return &ErrorField{
		In:    in,
		Field: field,
		Msg:   msg,
	}
}

type ErrorField struct {
	// 出错字段路径
	// 这个信息为一个 json 字符串,方便客户端进行定位错误原因
	// 例如输入中 {"name":{ "alias" : "test"}} 中的alias出错,则返回 "name.alias"
	// 如果alias是数组, 且第2个元素的a字段错误,则返回"name.alias[2].a"
	Field string `json:"field" xml:"field"`
	// 错误信息
	Msg string `json:"msg" xml:"msg"`
	// 错误字段位置
	// body, query, header, path, formData
	In string `json:"in" xml:"in"`
}

func (s ErrorField) String() string {
	return "[" + s.In + "]" + s.Field + " " + s.Msg
}

type ErrorFields []*ErrorField

func (fields ErrorFields) String() string {
	if len(fields) == 0 {
		return ""
	}
	buf := &bytes.Buffer{}
	buf.WriteString("<")
	for _, f := range fields {
		buf.WriteString(f.String())
	}
	buf.WriteString(">")
	return buf.String()
}

func (fields ErrorFields) Sort() ErrorFields {
	sort.Sort(fields)
	return fields
}

func (fields ErrorFields) Len() int {
	return len(fields)
}
func (fields ErrorFields) Swap(i, j int) {
	fields[i], fields[j] = fields[j], fields[i]
}
func (fields ErrorFields) Less(i, j int) bool {
	return fields[i].Field < fields[j].Field
}
