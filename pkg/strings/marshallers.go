package str

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/eden-framework/eden-framework/pkg/reflectx"
)

var StdStringifier = Stringifier{}

func init() {
	StdStringifier.Register(UnmarshalJSONUnmarshaler, UnmarshalTimeDuration)
}

func UnmarshalTimeDuration(s string, v reflect.Value) (matched bool, err error) {
	tpe := reflectx.IndirectType(v.Type())
	if tpe.PkgPath() == "time" && tpe.Name() == "Duration" {
		matched = true
		d, errForParse := time.ParseDuration(s)
		if errForParse != nil {
			err = errForParse
			return
		}
		v.SetInt(int64(d))
	}
	return
}

func UnmarshalJSONUnmarshaler(s string, rv reflect.Value) (matched bool, err error) {
	if rv.CanAddr() {
		rv = rv.Addr()
	}
	if unmarshaler, ok := rv.Interface().(json.Unmarshaler); ok {
		matched = true
		tpe := rv.Type()
		errForUnmarshal := unmarshaler.UnmarshalJSON([]byte(s))
		if errForUnmarshal != nil || !reflect.TypeOf(unmarshaler).Elem().ConvertibleTo(tpe) {
			err = errForUnmarshal
			return
		}
		rv.Set(reflect.ValueOf(unmarshaler).Elem().Convert(tpe))
	}
	return
}
