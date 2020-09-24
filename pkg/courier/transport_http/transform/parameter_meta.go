package transform

import (
	"encoding"
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/eden-framework/eden-framework/pkg/courier/status_error"
	"github.com/eden-framework/eden-framework/pkg/reflectx"
	"github.com/eden-framework/eden-framework/pkg/strings"
)

func NewParameterMeta(field *reflect.StructField, rv reflect.Value, tagIn string, tagInFlags TagFlags) *ParameterMeta {
	name, nameFlags := GetParameterDisplayName(field)

	p := &ParameterMeta{
		Name:    name,
		In:      tagIn,
		InFlags: tagInFlags,
		Field:   field,
		Value:   rv,
	}

	style, _, styleFlags := GetTagStyle(field)
	p.Style = style
	p.StyleFlags = styleFlags

	if p.Style == "" {
		switch p.In {
		case "query", "cookie":
			p.Style = "form"
			// todo set explode as default
		case "path", "header":
			p.Style = "simple"
		}
	}

	p.Format, _, _ = GetTagFmt(field)
	if p.Format == "" {
		p.Format = "json"
	}

	defaultValue, hasDefaultTag := GetTagDefault(field)
	p.DefaultValue = defaultValue

	p.Required = true
	if hasOmitempty, ok := nameFlags["omitempty"]; ok {
		p.Required = !hasOmitempty
	} else {
		// todo don't use non-default as required
		p.Required = !hasDefaultTag
	}

	tagValidate, _ := GetTagValidate(field)
	p.TagValidate = tagValidate

	tagErrMsg, _ := GetTagErrMsg(field)
	p.TagErrMsg = tagErrMsg

	return p
}

type ParameterMeta struct {
	Name   string
	In     string
	Format string

	InFlags TagFlags
	Field   *reflect.StructField
	Value   reflect.Value

	Style      string
	StyleFlags TagFlags

	DefaultValue string
	Required     bool

	TagValidate string
	TagErrMsg   string
}

func (p *ParameterMeta) IsExplode() bool {
	return p.StyleFlags != nil && p.StyleFlags["explode"]
}

func (p *ParameterMeta) IsMultipart() bool {
	return p.InFlags != nil && p.InFlags["multipart"]
}

// for parameter
func (p *ParameterMeta) UnmarshalStringAndValidate(ss ...string) error {
	return p.UnmarshalAndValidate(BytesList(ResolveCommaSplitValues(p.Field.Type, ss...)...)...)
}

func (p *ParameterMeta) UnmarshalAndValidate(dataList ...[]byte) error {
	if err := p.Unmarshal(dataList...); err != nil {
		if statusErr, ok := err.(*status_error.StatusError); ok {
			return statusErr
		}
		return status_error.InvalidField.StatusError().WithErrorField(p.In, p.Name, err.Error())
	}
	if valid, errFields := p.Validate(); !valid {
		return status_error.InvalidField.StatusError().WithErrorFields(errFields...)
	}
	return nil
}

// for body
func (p *ParameterMeta) UnmarshalFromReader(reader io.Reader) error {
	if reader == nil {
		return status_error.ReadFailed.StatusError()
	}
	data, readErr := ioutil.ReadAll(reader)
	if readErr != nil {
		return status_error.ReadFailed.StatusError()
	}
	if len(data) == 0 {
		return status_error.InvalidBodyStruct.StatusError().WithDesc("empty body")
	}
	return p.UnmarshalAndValidate(data)
}

func (p *ParameterMeta) Unmarshal(datalist ...[]byte) error {
	return ContentUnmarshal(p.Value, p.Field.Type, p.In, p.Name, p.Format, datalist...)
}

func (p *ParameterMeta) Marshal() (dataList [][]byte, err error) {
	return ContentMarshal(p.Value, p.In, p.Format)
}

var rtypeEncodingTextMarshaler = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()

func (p *ParameterMeta) Validate() (bool, status_error.ErrorFields) {
	errMsgMap := ErrMsgMap{}
	parentField := ""

	typ := reflectx.IndirectType(p.Field.Type)

	kind := typ.Kind()
	if (kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Array) && !typ.Implements(rtypeEncodingTextMarshaler) {
		validateScan := NewScanner()
		isValid, errMsgs := validateScan.Validate(p.Value, p.Field.Type)
		if !isValid {
			errMsgMap = errMsgMap.Merge(errMsgs)
		}
		if p.In == "formData" {
			parentField = p.Name
		}
	} else {
		errMsg := MarshalAndValidate(
			p.Value, p.Field.Type,
			p.DefaultValue, p.Required,
			p.TagValidate, p.TagErrMsg,
		)

		if errMsg != "" {
			errMsgMap[p.Name] = errMsg
		}
	}

	return len(errMsgMap) == 0, errMsgMap.ErrorFieldsIn(p.In, parentField)
}

func BytesList(ss ...string) (dataList [][]byte) {
	for _, s := range ss {
		dataList = append(dataList, []byte(s))
	}
	return
}

func ContentMarshal(rv reflect.Value, in string, format string) (dataList [][]byte, err error) {
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return
	}

	rv = reflect.Indirect(rv)

	if marshal, ok := rv.Interface().(encoding.TextMarshaler); ok {
		data, marshalErr := marshal.MarshalText()
		if marshalErr != nil {
			err = marshalErr
			return
		}
		dataList = [][]byte{data}
		return
	}

	switch rv.Kind() {
	case reflect.Map, reflect.Struct:
		data, errForMarshal := GetContentTransformer(format).Marshal(rv.Interface())
		if errForMarshal != nil {
			err = errForMarshal
			return
		}
		dataList = [][]byte{data}
		return
	case reflect.Slice, reflect.Array:
		if data, ok := rv.Interface().([]byte); ok {
			dataList = [][]byte{data}
			return
		}
		if in == "body" {
			data, errForMarshal := GetContentTransformer(format).Marshal(rv.Interface())
			if errForMarshal != nil {
				err = errForMarshal
				return
			}
			dataList = [][]byte{data}
			return
		}
		for i := 0; i < rv.Len(); i++ {
			itemsList, errForStringify := ContentMarshal(rv.Index(i), in, format)
			if errForStringify != nil {
				err = errForStringify
				return
			}
			dataList = append(dataList, itemsList...)
		}
		return
	default:
		data, errForStringify := Stringify(rv.Interface())
		if errForStringify != nil {
			err = errForStringify
			return
		}
		dataList = [][]byte{data}
		return
	}
}

func ContentUnmarshal(rv reflect.Value, tpe reflect.Type, in string, name string, format string, dataList ...[]byte) error {
	if len(dataList) == 0 || len(dataList[0]) == 0 {
		return nil
	}

	if rv.Kind() == reflect.Ptr && rv.IsNil() && rv.CanSet() {
		rv.Set(reflect.New(reflectx.IndirectType(tpe)))
	}

	rv = reflect.Indirect(rv)

	if textUnmarshaler, ok := rv.Addr().Interface().(encoding.TextUnmarshaler); ok {
		if err := textUnmarshaler.UnmarshalText(dataList[0]); err != nil {
			return err
		}
		return nil
	}

	switch rv.Kind() {
	case reflect.Map, reflect.Struct:
		if rv.CanAddr() {
			rv = rv.Addr()
		}
		return structUnmarshal(in, name, format, dataList[0], rv.Interface())
	case reflect.Slice:
		if _, ok := rv.Interface().([]byte); ok {
			rv.SetBytes(dataList[0])
			return nil
		}

		if in == "body" {
			if rv.CanAddr() {
				rv = rv.Addr()
			}
			return structUnmarshal(in, name, format, dataList[0], rv.Interface())
		}

		sliceRv := reflect.MakeSlice(tpe, len(dataList), cap(dataList))

		itemType := rv.Type().Elem()
		for i, data := range dataList {
			err := ContentUnmarshal(sliceRv.Index(i), itemType, in, name, format, data)
			if err != nil {
				return err
			}
		}

		rv.Set(sliceRv)

		return nil
	case reflect.Array:
		if in == "body" {
			if rv.CanAddr() {
				rv = rv.Addr()
			}
			return structUnmarshal(in, name, format, dataList[0], rv.Interface())
		}
		itemType := rv.Type().Elem()
		for i, data := range dataList {
			err := ContentUnmarshal(rv.Index(i), itemType, in, name, format, data)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return str.ConvertFromStr(string(dataList[0]), rv)
	}
}

func structUnmarshal(in string, rootField string, format string, data []byte, v interface{}) error {
	err := GetContentTransformer(format).Unmarshal(data, v)
	if err != nil {
		statusError := status_error.InvalidBodyStruct.StatusError()
		if unmarshalTypeErr, ok := err.(*json.UnmarshalTypeError); ok {
			if unmarshalTypeErr.Type == nil {
				return statusError.WithErrorField(in, LocateJSONPath(data, unmarshalTypeErr.Offset), unmarshalTypeErr.Value)
			}
			return statusError.WithErrorField(in, LocateJSONPath(data, unmarshalTypeErr.Offset), unmarshalTypeErr.Error())
		}
		return statusError.WithErrorField(in, rootField, "参数格式错误").WithDesc(err.Error())
	}
	return nil
}
