package reflectx

import (
	"fmt"
	"reflect"
)

func SetValue(rv reflect.Value, tpe reflect.Type, rightValue interface{}) error {
	if rightValue == nil {
		return nil
	}

	tpe = IndirectType(tpe)
	rightValueType := IndirectType(reflect.TypeOf(rightValue))
	if tpe != rightValueType {
		panic(fmt.Errorf("%s cannot set %s", tpe.String(), rightValueType.String()))
	}

	if rv.IsValid() && rv.Kind() == reflect.Ptr && rv.IsNil() && rv.CanSet() {
		rv.Set(reflect.New(tpe))
	}

	indirectRv := reflect.Indirect(rv)
	if indirectRv.CanSet() {
		indirectRv.Set(reflect.Indirect(reflect.ValueOf(rightValue)))
	}

	return nil
}
