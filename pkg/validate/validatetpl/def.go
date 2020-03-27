package validatetpl

import (
	"errors"
	"reflect"
	"regexp"
)

const (
	MAX_UINT8  = 1<<8 - 1
	MAX_INT8   = 1<<7 - 1
	MAX_UINT16 = 1<<16 - 1
	MAX_INT16  = 1<<15 - 1
	MAX_UINT32 = 1<<32 - 1
	MAX_INT32  = 1<<31 - 1
	//MAX_UINT64 = 1<<64 - 1
	MAX_UINT64 = 1<<63 - 1 // 兼容mysql driver只支持MAX 1<<63-1
	//MAX_INT64 = 1<<63 - 1
	MAX_INT64 = 1<<53 - 1 // 兼容浏览器只支持MAX 2<<53-1
	// floa32最大有效位数
	MAX_VALID_DIGIT_FLOAT32 = 7
	// floa64最大有效位数
	MAX_VALID_DIGIT_FLOAT64 = 15
)
const (
	DEFAULT_MAX_STRING_LENGTH = 1024
)
const (
	STRING_CHARS_NOT_IN_RANGE   = "字符串字数不在[%d, %d]范围内,当前长度: %d"
	STRING_LENGHT_NOT_IN_RANGE  = "字符串长度不在[%d， %d]范围内，当前长度：%d"
	STRING_VALUE_NOT_IN_ENUM    = "字符串不在%v集合内，当前值：%s"
	STRING_NOT_MATCH_REGEXP     = "字符串不符合正则表达式[%s]"
	INT_VALUE_NOT_IN_RANGE      = "整形值不在%s%d, %d%s范围内，当前值：%d"
	INT_VALUE_NOT_IN_ENUM       = "整形值不在%v集合内，当前值：%d"
	FLOAT_VALUE_NOT_IN_RANGE    = "浮点值不在%s%v,%v%s范围内，当前值：%v"
	FLOAT_VALUE_NOT_IN_ENUM     = "浮点值不在%v集合内，当前值：%v"
	FLOAT_VALUE_DIGIT_INVALID   = "浮点值小数位必须为%d，总位数不能超过%d位，当前值：%v"
	SLICE_ELEM_NUM_NOT_IN_RANGE = "切片元素个数不在[%d， %d]范围内，当前个数：%d"
	SLICE_ELEM_INVALID          = "切片元素不满足校验[%s]"
)
const (
	TYPE_NOT_STRING  = "非string类型"
	TYPE_NOT_UINT8   = "非uint8类型"
	TYPE_NOT_INT8    = "非int8类型"
	TYPE_NOT_UINT16  = "非uint16类型"
	TYPE_NOT_INT16   = "非int16类型"
	TYPE_NOT_UINT32  = "非uint32类型"
	TYPE_NOT_INT32   = "非int32类型"
	TYPE_NOT_UINT64  = "非uint64类型"
	TYPE_NOT_INT64   = "非int64类型"
	TYPE_NOT_FLOAT32 = "非float32类型"
	TYPE_NOT_FLOAT64 = "非float64类型"
	TYPE_NOT_SLICE   = "非slice类型"
)

const (
	UNLIMIT              = "unlimit"
	STRING_UNLIMIT_VALUE = -1
)

var (
	InvalidTypeStringError   = errors.New("invalid type string error")
	InvalidParamStringError  = errors.New("invalid param string error")
	InvalidTagStringError    = errors.New("invalid tag string error")
	Float32NoDefaultMaxError = errors.New("float32 no default max error")
	Float64NoDefaultMaxError = errors.New("float64 no default max error")
	InvalidTotalLenError     = errors.New("invalid total length error")
	InvalidDecimalLenError   = errors.New("invalid decimal length error")
)

var regexp_range_except_float_tag = regexp.MustCompile(`^@(int8|uint8|int16|uint16|int32|uint32|int64|uint64|string|char)(\[|\()\-?\d*\,|unlimit\-?\d*(\]|\))$`)

var regexp_range_float_tag = regexp.MustCompile(`^@(float32|float64)(\<\d*\,\d*\>)?((\[|\()-?\d*(\.\d*)?\,-?\d*(\.\d*)?(\]|\)))?$`)

var regexp_reg_tag = regexp.MustCompile(`^@regexp\[.+\]`)

var regexp_enum_tag = regexp.MustCompile(`^@(uint8|int8|uint16|int16|uint32|int32|uint64|int64|float32|float64|string)\{\w*(\,\w+)*\}$`)

var regexp_array_tag = regexp.MustCompile(`^@array(\[|\()\-?\d*\,\-?\d*(\]|\))(:@.*)?$`)

func getBrackets(exc_min, exc_max bool) (left_brackets, right_brackets string) {
	left_brackets = "["
	right_brackets = "]"
	if exc_min {
		left_brackets = "("
	}
	if exc_max {
		right_brackets = ")"
	}
	return
}

func indirect(v reflect.Value) reflect.Value {
	for {
		if v.Kind() == reflect.Interface {
			e := v.Elem()
			if e.Kind() == reflect.Ptr {
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}
		v = v.Elem()
	}

	return v
}
