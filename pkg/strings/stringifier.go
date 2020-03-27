package str

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	UnSupportTypeError = errors.New("un support type")
)

type StringUnmarshal func(s string, v reflect.Value) (bool, error)

type Stringifier struct {
	stringUnmarshalList []StringUnmarshal
}

func (stringifier *Stringifier) Register(stringUnmarshalList ...StringUnmarshal) {
	stringifier.stringUnmarshalList = append(stringifier.stringUnmarshalList, stringUnmarshalList...)
}

func (stringifier Stringifier) Unmarshal(str string, rv reflect.Value) error {
	rv = reflect.Indirect(rv)
	for _, stringUnmarshal := range stringifier.stringUnmarshalList {
		matched, err := stringUnmarshal(str, rv)
		if matched {
			if err != nil {
				return fmt.Errorf("unmarshal failed for %v: %s", stringUnmarshal, err.Error())
			}
			return nil
		}
	}
	return stringifier.UnmarshalBuiltIn(str, rv)
}

func (stringifier Stringifier) UnmarshalBuiltIn(str string, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.String:
		rv.SetString(str)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intV, err := strconv.ParseInt(str, 10, getBitSize(rv))
		if err != nil {
			return err
		}
		rv.Set(reflect.ValueOf(intV).Convert(rv.Type()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintV, err := strconv.ParseUint(str, 10, getBitSize(rv))
		if err != nil {
			return err
		}
		rv.Set(reflect.ValueOf(uintV).Convert(rv.Type()))
	case reflect.Float32, reflect.Float64:
		bitSize := getBitSize(rv)
		floatV, err := strconv.ParseFloat(str, bitSize)
		if err != nil {
			return err
		}
		rv.Set(reflect.ValueOf(floatV).Convert(rv.Type()))
	case reflect.Bool:
		boolV, err := strconv.ParseBool(str)
		if err != nil {
			return err
		}
		rv.SetBool(boolV)
	case reflect.Slice, reflect.Array:
		if str != "" {
			subStrValues := strings.Split(str, ",")
			elemType := rv.Type().Elem()
			tmpSlice := reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0)
			for _, subStrValue := range subStrValues {
				elemValue := reflect.New(elemType)
				err := stringifier.Unmarshal(subStrValue, elemValue)
				if err != nil {
					return err
				}
				tmpSlice = reflect.Append(tmpSlice, elemValue.Elem())
			}
			rv.Set(tmpSlice)
		}
	default:
		return UnSupportTypeError
	}
	return nil
}

func getBitSize(v reflect.Value) int {
	switch v.Kind() {
	case reflect.Int, reflect.Uint:
		return 32
	case reflect.Int8, reflect.Uint8:
		return 8
	case reflect.Int16, reflect.Uint16:
		return 16
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return 32
	case reflect.Int64, reflect.Uint64, reflect.Float64:
		return 64
	default:
		panic("only int, uint and float can support getBitSize")
	}
}
