package reflectx

import (
	"reflect"
)

func IndirectType(tpe reflect.Type) reflect.Type {
	for {
		if tpe.Kind() == reflect.Interface {
			e := tpe.Elem()
			if e.Kind() == reflect.Ptr {
				tpe = e
				continue
			}
		}
		if tpe.Kind() != reflect.Ptr {
			break
		}
		tpe = tpe.Elem()
	}
	return tpe
}

func Indirect(v reflect.Value) reflect.Value {
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

func New(tpe reflect.Type) reflect.Value {
	rv := reflect.New(tpe).Elem()
	if tpe.Kind() == reflect.Ptr {
		rv.Set(New(tpe.Elem()).Addr())
		return rv
	}
	return rv
}
