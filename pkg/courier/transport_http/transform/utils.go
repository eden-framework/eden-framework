package transform

import (
	"encoding"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/eden-framework/eden-framework/pkg/env"
	"github.com/eden-framework/reflectx"
	"github.com/eden-framework/strings"
)

type TagFlags map[string]bool

func GetTagValueAndFlags(tagString string) (v string, tagFlags TagFlags) {
	valueAndFlags := strings.Split(tagString, ",")
	v = valueAndFlags[0]
	tagFlags = TagFlags{}
	if len(valueAndFlags) > 1 {
		for _, flag := range valueAndFlags[1:] {
			tagFlags[flag] = true
		}
	}
	return
}

func GetStructFieldDisplayName(field *reflect.StructField) string {
	if jsonName, exists, _ := GetTagJSON(field); exists && jsonName != "" {
		return jsonName
	}
	return field.Name
}

func GetParameterDisplayName(field *reflect.StructField) (string, TagFlags) {
	if fieldName, exists, flags := GetTagName(field); exists && fieldName != "" {
		return fieldName, flags
	}
	if jsonName, exists, flags := GetTagJSON(field); exists && jsonName != "" {
		if !env.IsOnline() {
			//logrus.Warnf("%s `%s`, deprecated `json` tag for naming parameter, please use `name` tag instead", field.Name, field.Tag)
		}
		return jsonName, flags
	}
	return field.Name, TagFlags{}
}

func GetTagJSON(field *reflect.StructField) (name string, exists bool, tagFlags TagFlags) {
	name, exists = field.Tag.Lookup("json")
	if exists {
		name, tagFlags = GetTagValueAndFlags(name)
	}
	return
}

func GetTagName(field *reflect.StructField) (name string, exists bool, tagFlags TagFlags) {
	name, exists = field.Tag.Lookup("name")
	if exists {
		name, tagFlags = GetTagValueAndFlags(name)
	}
	return
}

func GetTagStyle(field *reflect.StructField) (name string, exists bool, tagFlags TagFlags) {
	name, exists = field.Tag.Lookup("style")
	if exists {
		name, tagFlags = GetTagValueAndFlags(name)
	}
	return
}

func GetTagFmt(field *reflect.StructField) (name string, exists bool, tagFlags TagFlags) {
	name, exists = field.Tag.Lookup("fmt")
	if exists {
		name, tagFlags = GetTagValueAndFlags(name)
	}
	return
}

func GetTagDefault(field *reflect.StructField) (string, bool) {
	return field.Tag.Lookup("default")
}

func GetTagValidate(field *reflect.StructField) (string, bool) {
	return field.Tag.Lookup("validate")
}

func GetTagErrMsg(field *reflect.StructField) (string, bool) {
	return field.Tag.Lookup("errMsg")
}

func GetTagIn(field *reflect.StructField) (tagIn string, hasIn bool, tagInFlags map[string]bool) {
	tagIn, hasIn = field.Tag.Lookup("in")
	if hasIn {
		tagIn, tagInFlags = GetTagValueAndFlags(tagIn)
	}
	return
}

func MustStructValue(functionName string, rv reflect.Value) {
	if rv.Kind() != reflect.Struct {
		logrus.Panicf("%s args must a struct value, but got %+v %+v", functionName, rv.Type(), rv.Interface())
	}
}

func Stringify(v interface{}) ([]byte, error) {
	if marshaler, ok := v.(json.Marshaler); ok {
		return Unquote(marshaler.MarshalJSON())
	}
	if textMarshaler, ok := v.(encoding.TextMarshaler); ok {
		return textMarshaler.MarshalText()
	}
	s, err := str.ConvertToStr(v)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

func Unquote(data []byte, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	if len(data) > 2 {
		if string(data[0]) == `"` && string(data[len(data)-1]) == `"` {
			return data[1 : len(data)-1], nil
		}
	}
	return data, err
}

func ResolveCommaSplitValues(tpe reflect.Type, ss ...string) (values []string) {
	tpe = reflectx.IndirectType(tpe)
	if tpe.Kind() == reflect.Slice || tpe.Kind() == reflect.Array {
		itemTpe := reflectx.IndirectType(tpe.Elem())
		if !(itemTpe.Kind() == reflect.Struct || itemTpe.Kind() == reflect.Map) {
			for _, s := range ss {
				values = append(values, strings.Split(s, ",")...)
			}
			return
		}
	}
	return ss
}

func LocateJSONPath(data []byte, offset int64) string {
	i := 0
	arrayPaths := map[string]bool{}
	arrayIdxSet := map[string]int{}
	pathWalker := &PathWalker{}

	markObjectKey := func() {
		jsonKey, l := nextString(data[i:])
		i += l

		if i < int(offset) && len(jsonKey) > 0 {
			key, _ := strconv.Unquote(string(jsonKey))
			pathWalker.Enter(key)
		}
	}

	markArrayIdx := func(path string) {
		if arrayPaths[path] {
			arrayIdxSet[path]++
		} else {
			arrayPaths[path] = true
		}
		pathWalker.Enter(arrayIdxSet[path])
	}

	for i < int(offset) {
		i += nextToken(data[i:])
		char := data[i]

		switch char {
		case '"':
			_, l := nextString(data[i:])
			i += l
		case '[', '{':
			i++

			if char == '[' {
				markArrayIdx(pathWalker.String())
			} else {
				markObjectKey()
			}
		case '}', ']', ',':
			i++
			pathWalker.Exit()

			if char == ',' {
				path := pathWalker.String()

				if _, ok := arrayPaths[path]; ok {
					markArrayIdx(path)
				} else {
					markObjectKey()
				}
			}
		default:
			i++
		}
	}

	return pathWalker.String()
}

func nextToken(data []byte) int {
	for i, c := range data {
		switch c {
		case ' ', '\n', '\r', '\t':
			continue
		default:
			return i
		}
	}
	return -1
}

func nextString(data []byte) (finalData []byte, l int) {
	quoteStartAt := -1
	for i, c := range data {
		switch c {
		case '"':
			if i > 0 && string(data[i-1]) == "\\" {
				continue
			}
			if quoteStartAt >= 0 {
				return data[quoteStartAt : i+1], i + 1
			} else {
				quoteStartAt = i
			}
		default:
			continue
		}
	}
	return nil, 0
}
