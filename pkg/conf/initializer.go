package conf

import (
	"github.com/eden-framework/eden-framework/pkg/reflectx"
	"reflect"
)

type Initializer interface {
	Init()
}

func Initialize(config ...interface{}) {
	for _, c := range config {
		rv := reflectx.Indirect(reflect.ValueOf(c))
		for i := 0; i < rv.NumField(); i++ {
			value := rv.Field(i)
			if conf, ok := value.Interface().(Initializer); ok {
				conf.Init()
			}
		}
	}
}
