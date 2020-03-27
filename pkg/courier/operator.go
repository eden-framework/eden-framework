package courier

import (
	"context"
	"fmt"
	"reflect"
)

type IOperator interface {
	Output(ctx context.Context) (interface{}, error)
}

type IContextProvider interface {
	IOperator
	ContextKey() string
}

type IEmptyOperator interface {
	IOperator
	NoOutput() bool
}

type EmptyOperator struct {
}

func (EmptyOperator) NoOutput() bool {
	return true
}

func (EmptyOperator) Output(ctx context.Context) (interface{}, error) {
	return nil, nil
}

func Group(path string) *GroupOperator {
	return &GroupOperator{
		path: path,
	}
}

type GroupOperator struct {
	IEmptyOperator
	path string
}

func (g *GroupOperator) Path() string {
	return g.path
}

func GetOperatorMeta(op IOperator, last bool) OperatorMeta {
	opMeta := OperatorMeta{}
	opMeta.IsLast = last
	if !opMeta.IsLast {
		ctxKey, ok := op.(IContextProvider)
		if !ok {
			panic(fmt.Sprintf("Operator %#v as middleware should has method `ContextKey() string`", op))
		}
		opMeta.ContextKey = ctxKey.ContextKey()
	}
	opMeta.Operator = op
	opMeta.Type = typeOfOperator(op)
	return opMeta
}

type OperatorMeta struct {
	IsLast     bool
	ContextKey string
	Operator   IOperator
	Type       reflect.Type
}

func ToOperatorMetaList(ops ...IOperator) (opMetas []OperatorMeta) {
	length := len(ops)
	for i, op := range ops {
		opMetas = append(opMetas, GetOperatorMeta(op, i == length-1))
	}
	return opMetas
}

func typeOfOperator(op IOperator) reflect.Type {
	tpe := reflect.TypeOf(op)
	for tpe.Kind() == reflect.Ptr {
		tpe = tpe.Elem()
	}
	return tpe
}
