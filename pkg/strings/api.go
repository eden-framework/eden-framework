package str

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
)

func ConvertFromStr(s string, v reflect.Value) error {
	if v.CanAddr() {
		if textUnmarshal, ok := v.Addr().Interface().(encoding.TextUnmarshaler); ok {
			return textUnmarshal.UnmarshalText([]byte(s))
		}
	}
	return StdStringifier.Unmarshal(s, v)
}

func ConvertToStr(v interface{}) (string, error) {
	if stringer, ok := v.(fmt.Stringer); ok {
		return stringer.String(), nil
	}
	rt := reflect.TypeOf(v)
	switch rt.Kind() {
	case reflect.String:
		return v.(string), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return fmt.Sprintf("%d", v), nil
	case reflect.Bool:
		return strconv.FormatBool(v.(bool)), nil
	case reflect.Float32:
		return strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64), nil
	default:
		return "", UnSupportTypeError
	}
}
