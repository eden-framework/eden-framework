package courier

import (
	"reflect"
)

type IDefaultsMarshal interface {
	MarshalDefaults(v interface{})
}

type OperatorDecoder func(op IOperator, rv reflect.Value) error

func NewOperatorBy(opType reflect.Type, defaultOp IOperator, decodeOperator OperatorDecoder) (op IOperator, err error) {
	rv := reflect.New(opType)
	op = rv.Interface().(IOperator)

	if defaultsMarshal, ok := defaultOp.(IDefaultsMarshal); ok {
		defaultsMarshal.MarshalDefaults(op)
	}

	err = decodeOperator(op, rv)
	if err != nil {
		return
	}

	return
}
