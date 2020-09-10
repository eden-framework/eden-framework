package conf

import (
	"github.com/profzone/eden-framework/pkg/reflectx"
	"reflect"
)

type Initializer interface {
	Init()
}

func Initialize(c interface{}) {
	rv := reflectx.Indirect(reflect.ValueOf(c))
	for i := 0; i < rv.NumField(); i++ {
		value := rv.Field(i)
		if conf, ok := value.Interface().(Initializer); ok {
			conf.Init()
		}
	}
}
