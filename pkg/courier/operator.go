package courier

import (
	"context"
	"fmt"
	"net/url"
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
	IEmptyOperator
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

type OperatorWithParams interface {
	OperatorParams() map[string][]string
}

type DefaultsSetter interface {
	SetDefaults()
}

type ContextProvider interface {
	IOperator
	ContextKey() string
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
	opMeta.Type = typeOfOperator(reflect.TypeOf(op))
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

func NewOperatorFactory(op IOperator, last bool) *OperatorFactory {
	opType := typeOfOperator(reflect.TypeOf(op))
	if opType.Kind() != reflect.Struct {
		panic(fmt.Errorf("operator must be a struct type, got %#v", op))
	}

	meta := &OperatorFactory{}
	meta.IsLast = last

	meta.Operator = op

	if _, isOperatorWithoutOutput := op.(IEmptyOperator); isOperatorWithoutOutput {
		meta.NoOutput = true
	}

	meta.Type = typeOfOperator(reflect.TypeOf(op))

	if operatorWithParams, ok := op.(OperatorWithParams); ok {
		meta.Params = operatorWithParams.OperatorParams()
	}

	if !meta.IsLast {
		if ctxKey, ok := op.(ContextProvider); ok {
			meta.ContextKey = ctxKey.ContextKey()
		} else {
			meta.ContextKey = meta.Type.String()
		}
	}

	return meta
}

func typeOfOperator(tpe reflect.Type) reflect.Type {
	for tpe.Kind() == reflect.Ptr {
		return typeOfOperator(tpe.Elem())
	}
	return tpe
}

type OperatorFactory struct {
	Type       reflect.Type
	ContextKey string
	NoOutput   bool
	Params     url.Values
	IsLast     bool
	Operator   IOperator
}

func (o *OperatorFactory) String() string {
	if o.Params != nil {
		return o.Type.String() + "?" + o.Params.Encode()
	}
	return o.Type.String()
}

func (o *OperatorFactory) New() IOperator {
	rv := reflect.New(o.Type)
	op := rv.Interface().(IOperator)

	if defaultsSetter, ok := op.(DefaultsSetter); ok {
		defaultsSetter.SetDefaults()
	}

	return op
}
