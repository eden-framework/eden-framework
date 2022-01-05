package application

import (
	"gitee.com/eden-framework/apollo"
	"gitee.com/eden-framework/reflectx"
	"github.com/sirupsen/logrus"
	"reflect"
)

type Option func(app *Application)

func WithInitializer(strict bool, init ...func() error) Option {
	return func(app *Application) {
		app.onInit = init
		app.onInitStrict = strict
	}
}

func WithApollo(conf *apollo.ApolloBaseConfig) Option {
	return func(app *Application) {
		app.apolloConfig = conf
	}
}

func WithConfig(conf ...interface{}) Option {
	return func(app *Application) {
		configFieldNames := make(map[string]struct{})
		for _, c := range app.envConfig {
			typ := reflectx.IndirectType(reflect.TypeOf(c))
			for j := 0; j < typ.NumField(); j++ {
				field := typ.Field(j)
				configFieldNames[field.Name] = struct{}{}
			}
		}

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

		app.envConfig = append(app.envConfig, conf...)
	}
}
