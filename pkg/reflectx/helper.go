package reflectx

import (
	"reflect"
)

func MustCanSetStruct(rv reflect.Value) {
	if !rv.CanSet() || rv.Type().Kind() != reflect.Struct {
		panic("need struct which can be set")
	}
}
