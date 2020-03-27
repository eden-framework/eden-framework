package validatetpl

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func GenerateRangeValidateFunc(str_type, str_param string, exc_min, exc_max bool) func(v interface{}) (bool, string) {
	switch str_type {
	case "uint8":
		params := strings.Split(str_param, ",")
		if len(params) == 2 {
			min, err := strconv.ParseUint(params[0], 10, 8)
			if err != nil {
				min = 0
			}
			max, err := strconv.ParseUint(params[1], 10, 8)
			if err != nil {
				max = MAX_UINT8
			}
			return NewRangeValidateUint8(uint8(min), uint8(max), exc_min, exc_max)
		}
	case "int8":
		params := strings.Split(str_param, ",")
		if len(params) == 2 {
			min, err := strconv.ParseInt(params[0], 10, 8)
			if err != nil {
				min = 0
			}
			max, err := strconv.ParseInt(params[1], 10, 8)
			if err != nil {
				max = MAX_INT8
			}
			return NewRangeValidateInt8(int8(min), int8(max), exc_min, exc_max)
		}
	case "uint16":
		params := strings.Split(str_param, ",")
		if len(params) == 2 {
			min, err := strconv.ParseUint(params[0], 10, 16)
			if err != nil {
				min = 0
			}
			max, err := strconv.ParseUint(params[1], 10, 16)
			if err != nil {
				max = MAX_UINT16
			}
			return NewRangeValidateUint16(uint16(min), uint16(max), exc_min, exc_max)
		}
	case "int16":
		params := strings.Split(str_param, ",")
		if len(params) == 2 {
			min, err := strconv.ParseInt(params[0], 10, 16)
			if err != nil {
				min = 0
			}
			max, err := strconv.ParseInt(params[1], 10, 16)
			if err != nil {
				max = MAX_INT16
			}
			return NewRangeValidateInt16(int16(min), int16(max), exc_min, exc_max)
		}
	case "uint32":
		params := strings.Split(str_param, ",")
		if len(params) == 2 {
			min, err := strconv.ParseUint(params[0], 10, 32)
			if err != nil {
				min = 0
			}
			max, err := strconv.ParseUint(params[1], 10, 32)
			if err != nil {
				max = MAX_UINT32
			}
			return NewRangeValidateUint32(uint32(min), uint32(max), exc_min, exc_max)
		}
	case "int32":
		params := strings.Split(str_param, ",")
		if len(params) == 2 {
			min, err := strconv.ParseInt(params[0], 10, 32)
			if err != nil {
				min = 0
			}
			max, err := strconv.ParseInt(params[1], 10, 32)
			if err != nil {
				max = MAX_INT32
			}
			return NewRangeValidateInt32(int32(min), int32(max), exc_min, exc_max)
		}
	case "uint64":
		params := strings.Split(str_param, ",")
		if len(params) == 2 {
			min, err := strconv.ParseUint(params[0], 10, 64)
			if err != nil {
				min = 0
			}
			max, err := strconv.ParseUint(params[1], 10, 64)
			if err != nil {
				max = MAX_UINT64
			}
			return NewRangeValidateUint64(uint64(min), uint64(max), exc_min, exc_max)
		}
	case "int64":
		params := strings.Split(str_param, ",")
		if len(params) == 2 {
			min, err := strconv.ParseInt(params[0], 10, 64)
			if err != nil {
				min = 0
			}
			max, err := strconv.ParseInt(params[1], 10, 64)
			if err != nil {
				max = MAX_INT64
			}
			return NewRangeValidateInt64(int64(min), int64(max), exc_min, exc_max)
		}
	case "float32":
		params := strings.Split(str_param, ",")
		if len(params) == 2 {
			min, err := strconv.ParseFloat(params[0], 32)
			if err != nil {
				min = 0
			}
			max, err := strconv.ParseFloat(params[1], 32)
			if err != nil {
				panic(fmt.Sprintf("general range validate func failed[err:%s, tag:%s]", Float32NoDefaultMaxError, str_type))
			}
			return NewRangeValidateFloat32(float32(min), float32(max), exc_min, exc_max)
		}
	case "float64":
		params := strings.Split(str_param, ",")
		if len(params) == 2 {
			min, err := strconv.ParseFloat(params[0], 64)
			if err != nil {
				min = 0
			}
			max, err := strconv.ParseFloat(params[1], 64)
			if err != nil {
				panic(fmt.Sprintf("general range validate func failed[err:%s, tag:%s]", Float64NoDefaultMaxError, str_type))
			}
			return NewRangeValidateFloat64(float64(min), float64(max), exc_min, exc_max)
		}
	case "string":
		params := strings.Split(str_param, ",")
		if len(params) == 2 {
			min, err := strconv.ParseInt(params[0], 10, 0)
			if err != nil {
				min = 0
			}

			var max int64 = 0
			if params[1] == UNLIMIT {
				max = STRING_UNLIMIT_VALUE
			} else {
				max, err = strconv.ParseInt(params[1], 10, 0)
				if err != nil {
					max = DEFAULT_MAX_STRING_LENGTH
				}
			}
			return NewValidateStringLength(int(min), int(max))
		}
	case "char":
		params := strings.Split(str_param, ",")
		if len(params) == 2 {
			min, err := strconv.ParseInt(params[0], 10, 0)
			if err != nil {
				min = 0
			}
			max, err := strconv.ParseInt(params[1], 10, 0)
			return NewValidateChar(int(min), int(max))
		}
	case "regexp":
		reg, err := regexp.Compile(str_param)
		if err == nil {
			return NewValidateStringRegExp(reg)
		}
	default:
		panic(fmt.Sprintf("general range validate func failed[err:%s, type:%s]", InvalidTypeStringError, str_type))
	}
	panic(fmt.Sprintf("general range validate func failed[err:%s, type:%s]", InvalidParamStringError, str_param))
}

func GenerateEnumValidateFunc(str_type, str_param string) func(v interface{}) (bool, string) {
	switch str_type {
	case "uint8":
		params := strings.Split(str_param, ",")
		enum_values := []uint8{}
		for _, param := range params {
			num, err := strconv.ParseUint(param, 10, 8)
			if err != nil {
				panic(fmt.Sprintf("general uint8 enum validate func failed[err:%s, param:%s]", InvalidParamStringError, str_param))
			}
			enum_values = append(enum_values, uint8(num))
		}
		return NewEnumValidateUint8(enum_values...)
	case "int8":
		params := strings.Split(str_param, ",")
		enum_values := []int8{}
		for _, param := range params {
			num, err := strconv.ParseInt(param, 10, 8)
			if err != nil {
				panic(fmt.Sprintf("general int8 enum validate func failed[err:%s, param:%s]", InvalidParamStringError, str_param))
			}
			enum_values = append(enum_values, int8(num))
		}
		return NewEnumValidateInt8(enum_values...)
	case "uint16":
		params := strings.Split(str_param, ",")
		enum_values := []uint16{}
		for _, param := range params {
			num, err := strconv.ParseUint(param, 10, 16)
			if err != nil {
				panic(fmt.Sprintf("general uint16 enum validate func failed[err:%s, param:%s]", InvalidParamStringError, str_param))
			}
			enum_values = append(enum_values, uint16(num))
		}
		return NewEnumValidateUint16(enum_values...)
	case "int16":
		params := strings.Split(str_param, ",")
		enum_values := []int16{}
		for _, param := range params {
			num, err := strconv.ParseInt(param, 10, 16)
			if err != nil {
				panic(fmt.Sprintf("general int16 enum validate func failed[err:%s, param:%s]", InvalidParamStringError, str_param))
			}
			enum_values = append(enum_values, int16(num))
		}
		return NewEnumValidateInt16(enum_values...)
	case "uint32":
		params := strings.Split(str_param, ",")
		enum_values := []uint32{}
		for _, param := range params {
			num, err := strconv.ParseUint(param, 10, 32)
			if err != nil {
				panic(fmt.Sprintf("general uint32 enum validate func failed[err:%s, param:%s]", InvalidParamStringError, str_param))
			}
			enum_values = append(enum_values, uint32(num))
		}
		return NewEnumValidateUint32(enum_values...)
	case "int32":
		params := strings.Split(str_param, ",")
		enum_values := []int32{}
		for _, param := range params {
			num, err := strconv.ParseInt(param, 10, 32)
			if err != nil {
				panic(fmt.Sprintf("general int32 enum validate func failed[err:%s, param:%s]", InvalidParamStringError, str_param))
			}
			enum_values = append(enum_values, int32(num))
		}
		return NewEnumValidateInt32(enum_values...)
	case "uint64":
		params := strings.Split(str_param, ",")
		enum_values := []uint64{}
		for _, param := range params {
			num, err := strconv.ParseUint(param, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("general uint64 enum validate func failed[err:%s, param:%s]", InvalidParamStringError, str_param))
			}
			enum_values = append(enum_values, uint64(num))
		}
		return NewEnumValidateUint64(enum_values...)
	case "int64":
		params := strings.Split(str_param, ",")
		enum_values := []int64{}
		for _, param := range params {
			num, err := strconv.ParseInt(param, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("general int64 enum validate func failed[err:%s, param:%s]", InvalidParamStringError, str_param))
			}
			enum_values = append(enum_values, int64(num))
		}
		return NewEnumValidateInt64(enum_values...)
	case "string":
		params := strings.Split(str_param, ",")
		return NewEnumValidateString(params...)
	default:
		panic(fmt.Sprintf("general enum validate func failed[err:%s, type:%s]", InvalidTypeStringError, str_type))
	}
}

func GenerateDecimalValidateFunc(str_type, str_param string) func(v interface{}) (bool, string) {
	switch str_type {
	case "float32":
		params := strings.Split(str_param, ",")
		if len(params) == 2 {
			total_len, err := strconv.ParseUint(params[0], 10, 32)
			if err != nil {
				panic(fmt.Sprintf("general decimal validate func failed[err:%s, tag:%s]", InvalidTotalLenError, str_type))
			}
			decimal_len, err := strconv.ParseUint(params[1], 10, 32)
			if err != nil {
				panic(fmt.Sprintf("general decimal validate func failed[err:%s, tag:%s]", InvalidDecimalLenError, str_type))
			}
			if total_len > MAX_VALID_DIGIT_FLOAT32 {
				panic(fmt.Sprintf("general decimal validate func failed[float32 valid len should not greater than %d]", MAX_VALID_DIGIT_FLOAT32))
			} else if decimal_len > total_len {
				panic(fmt.Sprintf("general decimal validate func failed[decimal len should not greater than total len[total:%d, decimal:%d]]",
					total_len, decimal_len))
			}
			return NewDecimalValidateFloat32(int(total_len), int(decimal_len))
		}
	case "float64":
		params := strings.Split(str_param, ",")
		if len(params) == 2 {
			total_len, err := strconv.ParseUint(params[0], 10, 32)
			if err != nil {
				panic(fmt.Sprintf("general range validate func failed[err:%s, tag:%s]", InvalidTotalLenError, str_type))
			}
			decimal_len, err := strconv.ParseUint(params[1], 10, 32)
			if err != nil {
				panic(fmt.Sprintf("general range validate func failed[err:%s, tag:%s]", InvalidDecimalLenError, str_type))
			}
			if total_len > MAX_VALID_DIGIT_FLOAT64 {
				panic(fmt.Sprintf("general decimal validate func failed[float32 valid len should not greater than %d]", MAX_VALID_DIGIT_FLOAT32))
			} else if decimal_len > total_len {
				panic(fmt.Sprintf("general decimal validate func failed[decimal len should not greater than total len[total:%d, decimal:%d]]",
					total_len, decimal_len))
			}
			return NewDecimalValidateFloat64(int(total_len), int(decimal_len))
		}
	default:
		panic(fmt.Sprintf("general decimal validate func failed[err:%s, type:%s]", InvalidTypeStringError, str_type))
	}
	panic(fmt.Sprintf("general decimal validate func failed[err:%s, type:%s]", InvalidParamStringError, str_param))
}
func GenerateSliceValidateFunc(str_param, item_tag string) func(v interface{}) (bool, string) {
	params := strings.Split(str_param, ",")
	var min, max uint64
	var err error
	if len(params) == 2 {
		min, err = strconv.ParseUint(params[0], 10, 64)
		if err != nil {
			min = 0
		}
		max, err = strconv.ParseUint(params[1], 10, 64)
		if err != nil {
			max = 0
		}
	}
	return NewValidateSlice(min, max, GenerateValidateFuncByTag(item_tag))
}

func GenerateValidateFuncByTag(tag string) func(v interface{}) (bool, string) {
	if regexp_range_except_float_tag.MatchString(tag) || regexp_reg_tag.MatchString(tag) {
		param_start := strings.IndexAny(tag, "[(")
		param_end := strings.LastIndexAny(tag, "])")
		if param_start == -1 || param_end == -1 {
			panic(fmt.Sprintf("generate range validate func failed[err:%s, tag:%s]", InvalidTagStringError, tag))
		}

		lb := tag[param_start : param_start+1]
		rb := tag[param_end : param_end+1]
		str_type := tag[1:param_start]
		str_param := tag[param_start+1 : param_end]

		exc_min := false
		exc_max := false
		if lb == "(" {
			exc_min = true
		}
		if rb == ")" {
			exc_max = true
		}
		return GenerateRangeValidateFunc(str_type, str_param, exc_min, exc_max)
	} else if regexp_enum_tag.MatchString(tag) {
		param_start := strings.IndexAny(tag, "{")
		param_end := strings.LastIndexAny(tag, "}")
		if param_start == -1 || param_end == -1 {
			panic(fmt.Sprintf("generate enum validate func failed[err:%s, tag:%s]", InvalidTagStringError, tag))
		}
		str_type := tag[1:param_start]
		str_param := tag[param_start+1 : param_end]
		return GenerateEnumValidateFunc(str_type, str_param)
	} else if regexp_array_tag.MatchString(tag) {
		tags := strings.SplitN(tag, ":", 2)
		if len(tags) != 2 {
			//panic(fmt.Sprintf("generate slice validate func failed[err:%s, tag:%s]", InvalidTagStringError, tag))
			tags = append(tags, "@")
		}
		array_tag := tags[0]
		item_tag := tags[1]

		param_start := strings.IndexAny(array_tag, "[")
		param_end := strings.LastIndexAny(array_tag, "]")
		if param_start == -1 || param_end == -1 {
			panic(fmt.Sprintf("generate slice validate func failed[err:%s, tag:%s]", InvalidTagStringError, tag))
		}
		str_param := tag[param_start+1 : param_end]
		return GenerateSliceValidateFunc(str_param, item_tag)
	} else if regexp_range_float_tag.MatchString(tag) {
		valideFuncList := []func(v interface{}) (bool, string){}
		var param_start, param_end int
		var str_type string
		// float位数校验
		if param_start = strings.IndexAny(tag, "<"); param_start != -1 {
			param_end = strings.LastIndexAny(tag, ">")
			if param_end == -1 {
				panic(fmt.Sprintf("generate range validate func failed[err:%s, tag:%s]", InvalidTagStringError, tag))
			}

			str_type = tag[1:param_start]
			str_param := tag[param_start+1 : param_end]
			decimalValidFunc := GenerateDecimalValidateFunc(str_type, str_param)
			if decimalValidFunc != nil {
				valideFuncList = append(valideFuncList, decimalValidFunc)
			}
		}

		// float值范围校验
		if param_start = strings.IndexAny(tag, "[("); param_start != -1 {
			param_end = strings.LastIndexAny(tag, "])")
			if param_end == -1 {
				panic(fmt.Sprintf("generate range validate func failed[err:%s, tag:%s]", InvalidTagStringError, tag))
			}
			// type string
			if str_type == "" {
				str_type = tag[1:param_start]
			}
			// param string
			str_param := tag[param_start+1 : param_end]

			exc_min := false
			exc_max := false
			if lb := tag[param_start : param_start+1]; lb == "(" {
				exc_min = true
			}
			if rb := tag[param_end : param_end+1]; rb == ")" {
				exc_max = true
			}

			rangeValidFunc := GenerateRangeValidateFunc(str_type, str_param, exc_min, exc_max)
			if rangeValidFunc != nil {
				valideFuncList = append(valideFuncList, rangeValidFunc)
			}
		}
		return func(v interface{}) (bool, string) {
			for _, validateFunc := range valideFuncList {
				if ok, err_msg := validateFunc(v); !ok {
					return false, err_msg
				}
			}
			return true, ""
		}
	}

	// 非Auto-Generate类型，从presetValidateFuncMap查询
	if validateFunc, ok := GetValidateFunc(tag); ok {
		return validateFunc
	}
	return nil
}
