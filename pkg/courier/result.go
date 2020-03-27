package courier

import (
	"fmt"
	"reflect"

	"github.com/profzone/eden-framework/pkg/courier/status_error"
)

type TUnmarshal func(data []byte, v interface{}) error

type Result struct {
	Err       error
	Meta      Metadata
	Data      []byte
	Unmarshal TUnmarshal
}

func (r Result) BindMeta(meta Metadata) *Result {
	for key, values := range r.Meta {
		meta.Set(key, values...)
	}
	return &r
}

func (r Result) Into(v interface{}) error {
	if r.Err != nil {
		return r.Err
	}
	if r.Unmarshal == nil {
		r.Unmarshal = UnmarshalBytes
	}
	if len(r.Data) > 0 {
		err := r.Unmarshal(r.Data, v)
		if err != nil {
			return status_error.InvalidStruct.StatusError().WithDesc(err.Error())
		}
	}
	return nil
}

func UnmarshalBytes(data []byte, v interface{}) error {
	if _, ok := v.(*[]byte); ok {
		reflect.Indirect(reflect.ValueOf(v)).SetBytes(data)
		return nil
	}
	return fmt.Errorf("target value %#v is not []byte, data: %v", v, data)
}
