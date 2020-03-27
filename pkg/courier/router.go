package courier

import (
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

type IServe interface {
	Serve(router *Router) error
}

func Run(router *Router, serves ...IServe) {
	errs := make(chan error)

	for i := range serves {
		s := serves[i]
		go func() {
			if err := s.Serve(router); err != nil {
				errs <- err
			}
		}()
	}

	select {
	case err := <-errs:
		logrus.Errorf("%s", err.Error())
	}
}

type Router struct {
	operators []IOperator
	parent    *Router
	children  map[*Router]bool
}

func NewRouter(operators ...IOperator) *Router {
	return &Router{
		operators: operators,
	}
}

func (router *Router) Register(r *Router) {
	if router.children == nil {
		router.children = map[*Router]bool{}
	}
	r.parent = router
	router.children[r] = true
}

func (router *Router) Route() *Route {
	parent := router.parent
	operators := router.operators

	for parent != nil {
		operators = append(parent.operators, operators...)
		parent = parent.parent
	}

	return &Route{
		last:      router.children == nil,
		Operators: operators,
	}
}

func (router *Router) Routes() (routes []*Route) {
	for child := range router.children {
		route := child.Route()
		if route.last && len(route.Operators) > 0 {
			routes = append(routes, route)
		}
		if child.children != nil {
			routes = append(routes, child.Routes()...)
		}
	}
	return
}

type Route struct {
	Operators []IOperator
	last      bool
}

func (route *Route) EffectiveOperators() (operators []IOperator, operatorTypeNames []string) {
	for _, operator := range route.Operators {
		if _, isEmptyOperator := operator.(IEmptyOperator); !isEmptyOperator {
			operatorType := reflect.TypeOf(operator)
			modify := ""
			if stringer, ok := operator.(fmt.Stringer); ok {
				modify = stringer.String()
			}
			if operatorType.Kind() == reflect.Ptr {
				operatorTypeNames = append(operatorTypeNames, "*"+operatorType.Elem().Name()+modify)
			} else {
				operatorTypeNames = append(operatorTypeNames, operatorType.Name()+modify)
			}
			operators = append(operators, operator)
		}
	}
	return
}
