package transform

import (
	"go/ast"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/profzone/eden-framework/pkg/courier/status_error"
	"github.com/profzone/eden-framework/pkg/reflectx"
)

func ParameterGroupFromValue(v interface{}) (m *ParameterGroup) {
	return ParameterGroupFromReflectValue(reflect.ValueOf(v))
}

func ParameterGroupFromReflectValue(rv reflect.Value) (group *ParameterGroup) {
	rv = reflect.Indirect(rv)
	MustStructValue("ParameterGroupFromReflectValue", rv)
	group = &ParameterGroup{}
	group.Value = rv
	group.Scan()
	return group
}

type ParameterGroup struct {
	Parameters ParameterMap
	Context    ParameterMap
	FormData   ParameterMap
	Body       *ParameterMeta
	Value      reflect.Value
}

func (m *ParameterGroup) ValidateNoBodyByHook() (bool, status_error.ErrorFields) {
	rv := m.Value
	tpe := rv.Type()

	errFields := status_error.ErrorFields{}

	for i := 0; i < rv.NumMethod(); i++ {
		method := rv.Method(i)
		validateMethod := tpe.Method(i)
		validatePrefix := "Validate"
		fieldName := strings.TrimPrefix(validateMethod.Name, validatePrefix)
		if validatePrefix+fieldName == validateMethod.Name {
			if parameterMeta, ok := m.Parameters.Get(fieldName); ok {
				if validateFn, ok := method.Interface().(func() string); ok {
					msg := validateFn()
					if msg != "" {
						errFields = append(errFields, status_error.NewErrorField(parameterMeta.In, parameterMeta.Name, msg))
					}

				}
			}
		}
	}

	return len(errFields) == 0, errFields
}

func (m *ParameterGroup) Scan() {
	m.Parameters = ParameterMap{}
	m.FormData = ParameterMap{}
	m.Context = ParameterMap{}
	m.scan(m.Value)
}

func (m *ParameterGroup) scanFormData(field *reflect.StructField, rv reflect.Value, tagInFlags map[string]bool) {
	indirectRv := reflect.Indirect(rv)

	tpe := reflectx.IndirectType(field.Type)

	if indirectRv.Kind() != reflect.Struct {
		if _, hasTagIn, tagInFlags := GetTagIn(field); hasTagIn {
			m.FormData.Add(NewParameterMeta(field, rv, "formData", tagInFlags))
		}
		return
	}

	for i := 0; i < tpe.NumField(); i++ {
		f := tpe.Field(i)

		if !ast.IsExported(f.Name) {
			continue
		}

		fieldValue := indirectRv.Field(i)

		if f.Anonymous {
			m.scanFormData(&f, fieldValue, tagInFlags)
			continue
		}

		m.FormData.Add(NewParameterMeta(&f, fieldValue, "formData", tagInFlags))
	}
}

func (m *ParameterGroup) scan(rv reflect.Value) {
	rv = reflect.Indirect(rv)

	if rv.Kind() != reflect.Struct {
		return
	}

	tpe := rv.Type()
	for i := 0; i < tpe.NumField(); i++ {
		field := tpe.Field(i)
		if !ast.IsExported(field.Name) {
			continue
		}

		fieldValue := rv.Field(i)
		tagIn, hasTagIn, tagInFlags := GetTagIn(&field)

		if !hasTagIn {
			if field.Anonymous {
				m.scan(fieldValue)
				continue
			}
			logrus.Panicf("%s.%s has no \"in\" tag", tpe.Name(), field.Name)
		}

		parameterMeta := NewParameterMeta(&field, fieldValue, tagIn, tagInFlags)

		switch parameterMeta.In {
		case "formData":
			m.scanFormData(&field, fieldValue, tagInFlags)
		case "body":
			m.Body = parameterMeta
		case "context":
			// todo removed context
			m.Context.Add(parameterMeta)
		default:
			m.Parameters.Add(parameterMeta)
		}
	}
}
