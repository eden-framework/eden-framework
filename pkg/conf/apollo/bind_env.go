package apollo

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var g_separator string = "_" //环境变量使用的分隔符

func SetSeparator(sep string) {
	g_separator = sep
}
func GetSeparator() string {
	return g_separator
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

func BindEnv(config interface{}, prefix string) {
	rv := reflect.ValueOf(config)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("can not bind env to invalid pointer")
	}

	rv = indirect(rv)
	switch rv.Kind() {
	case reflect.Struct:
		bindEnvStruct(rv, prefix)
	case reflect.Slice:
		bindEnvSlice(rv, prefix)
	default:
		panic("Config is not a struct or slice")
	}
}

func bindEnvSlice(rv reflect.Value, prefix string) {
	for i := 0; i < rv.Len(); i++ {
		v_item := rv.Index(i)
		env_key := getEnvKey(prefix, strconv.FormatInt(int64(i), 10))
		env_value, exist := os.LookupEnv(env_key)
		if !exist && v_item.Kind() != reflect.Struct && v_item.Kind() != reflect.Slice {
			continue
		}
		switch v_item.Kind() {
		case reflect.Bool:
			setBoolValue(v_item, env_value)
		case reflect.Int:
			setIntValue(v_item, env_value)
		case reflect.Uint:
			setUintValue(v_item, env_value)
		case reflect.Int8:
			setInt8Value(v_item, env_value)
		case reflect.Uint8:
			setUint8Value(v_item, env_value)
		case reflect.Int16:
			setInt16Value(v_item, env_value)
		case reflect.Uint16:
			setUint16Value(v_item, env_value)
		case reflect.Int32:
			setInt32Value(v_item, env_value)
		case reflect.Uint32:
			setUint32Value(v_item, env_value)
		case reflect.Int64:
			setInt64Value(v_item, env_value)
		case reflect.Uint64:
			setUint64Value(v_item, env_value)
		case reflect.Float32:
			setFloat32Value(v_item, env_value)
		case reflect.Float64:
			setFloat64Value(v_item, env_value)
		case reflect.String:
			setStringValue(v_item, env_value)
		case reflect.Struct:
			bindEnvStruct(v_item, env_key)
		case reflect.Slice:
			bindEnvSlice(v_item, env_key)
		default:
			continue
		}
	}
}

func bindEnvStruct(rv reflect.Value, prefix string) {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		t_field := rt.Field(i)
		v_field := rv.Field(i)
		field_name := t_field.Name
		env_key := getEnvKey(prefix, field_name)
		env_value, exist := os.LookupEnv(env_key)
		if !exist && v_field.Kind() != reflect.Struct && v_field.Kind() != reflect.Slice {
			continue
		}
		switch v_field.Kind() {
		case reflect.Struct:
			bindEnvStruct(v_field, env_key)
		case reflect.Bool:
			setBoolValue(v_field, env_value)
		case reflect.Int:
			setIntValue(v_field, env_value)
		case reflect.Uint:
			setUintValue(v_field, env_value)
		case reflect.Int8:
			setInt8Value(v_field, env_value)
		case reflect.Uint8:
			setUint8Value(v_field, env_value)
		case reflect.Int16:
			setInt16Value(v_field, env_value)
		case reflect.Uint16:
			setUint16Value(v_field, env_value)
		case reflect.Int32:
			setInt32Value(v_field, env_value)
		case reflect.Uint32:
			setUint32Value(v_field, env_value)
		case reflect.Int64:
			setInt64Value(v_field, env_value)
		case reflect.Uint64:
			setUint64Value(v_field, env_value)
		case reflect.Float32:
			setFloat32Value(v_field, env_value)
		case reflect.Float64:
			setFloat64Value(v_field, env_value)
		case reflect.String:
			setStringValue(v_field, env_value)

		case reflect.Slice:
			bindEnvSlice(v_field, env_key)
		default:
			continue
		}
	}
}

func getEnvKey(prefix, field_name string) string {
	if prefix == "" {
		return strings.ToUpper(field_name)
	}
	return strings.ToUpper(fmt.Sprintf("%s%s%s", prefix, g_separator, field_name))
}

func setBoolValue(v reflect.Value, str_value string) {
	str_value = strings.ToLower(str_value)
	switch str_value {
	case "t":
		fallthrough
	case "true":
		v.SetBool(true)
	case "f":
		fallthrough
	case "false":
		v.SetBool(false)
	default:
	}
}

func setIntValue(v reflect.Value, str_value string) {
	if num, err := strconv.ParseInt(str_value, 10, 0); err == nil {
		v.Set(reflect.ValueOf(int(num)))
	}
}

func setUintValue(v reflect.Value, str_value string) {
	if num, err := strconv.ParseUint(str_value, 10, 0); err == nil {
		v.Set(reflect.ValueOf(uint(num)))
	}
}

func setInt8Value(v reflect.Value, str_value string) {
	if num, err := strconv.ParseInt(str_value, 10, 8); err == nil {
		v.Set(reflect.ValueOf(int8(num)))
	}
}

func setUint8Value(v reflect.Value, str_value string) {
	if num, err := strconv.ParseUint(str_value, 10, 8); err == nil {
		v.Set(reflect.ValueOf(uint8(num)))
	}
}

func setInt16Value(v reflect.Value, str_value string) {
	if num, err := strconv.ParseInt(str_value, 10, 16); err == nil {
		v.Set(reflect.ValueOf(int16(num)))
	}
}

func setUint16Value(v reflect.Value, str_value string) {
	if num, err := strconv.ParseUint(str_value, 10, 16); err == nil {
		v.Set(reflect.ValueOf(uint16(num)))
	}
}
func setInt32Value(v reflect.Value, str_value string) {
	if num, err := strconv.ParseInt(str_value, 10, 32); err == nil {
		v.Set(reflect.ValueOf(int32(num)))
	}
}

func setUint32Value(v reflect.Value, str_value string) {
	if num, err := strconv.ParseUint(str_value, 10, 32); err == nil {
		v.Set(reflect.ValueOf(uint32(num)))
	}
}
func setInt64Value(v reflect.Value, str_value string) {
	if num, err := strconv.ParseInt(str_value, 10, 64); err == nil {
		v.Set(reflect.ValueOf(int64(num)))
	}
}

func setUint64Value(v reflect.Value, str_value string) {
	if num, err := strconv.ParseUint(str_value, 10, 64); err == nil {
		v.Set(reflect.ValueOf(uint64(num)))
	}
}

func setFloat32Value(v reflect.Value, str_value string) {
	if num, err := strconv.ParseFloat(str_value, 32); err == nil {
		v.Set(reflect.ValueOf(float32(num)))
	}
}

func setFloat64Value(v reflect.Value, str_value string) {
	if num, err := strconv.ParseFloat(str_value, 64); err == nil {
		v.Set(reflect.ValueOf(float64(num)))
	}
}

func setStringValue(v reflect.Value, str_value string) {
	v.SetString(str_value)
}
