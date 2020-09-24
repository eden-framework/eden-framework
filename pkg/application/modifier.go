package application

import (
	"github.com/eden-framework/eden-framework/pkg/conf/apollo"
	"github.com/eden-framework/eden-framework/pkg/reflectx"
	"github.com/sirupsen/logrus"
	"reflect"
)

type Option func(app *Application)

func WithApollo(conf *apollo.ApolloBaseConfig) Option {
	return func(app *Application) {
		app.apolloConfig = conf
	}
}

func WithConfig(conf ...interface{}) Option {
	configFieldNames := make(map[string]struct{})
	for i, c := range conf {
		typ := reflect.TypeOf(c)
		if typ.Kind() != reflect.Ptr {
			logrus.Panicf("the [%d] config must be a ptr value", i)
		}

		typ = reflectx.IndirectType(typ)
		for j := 0; j < typ.NumField(); j++ {
			field := typ.Field(j)
			if _, ok := configFieldNames[field.Name]; ok {
				logrus.Panicf("the [%d] config field named [%s] is duplicated. can not define the same field name in the root of each config struct.", i, field.Name)
			} else {
				configFieldNames[field.Name] = struct{}{}
			}
		}
	}
	return func(app *Application) {
		app.envConfig = append(app.envConfig, conf...)
	}
}
