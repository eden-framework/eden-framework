package application

import (
	"github.com/profzone/eden-framework/pkg/conf/apollo"
	"github.com/sirupsen/logrus"
	"reflect"
)

type ApplicationOption func(app *Application)

func WithApollo(conf *apollo.ApolloBaseConfig) ApplicationOption {
	return func(app *Application) {
		app.apolloConfig = conf
		app.envConfig = append(app.envConfig, conf)
	}
}

func WithConfig(conf ...interface{}) ApplicationOption {
	for i, c := range conf {
		typ := reflect.TypeOf(c)
		if typ.Kind() != reflect.Ptr {
			logrus.Panicf("the [%d] config must be a ptr value", i)
		}
	}
	return func(app *Application) {
		app.envConfig = append(app.envConfig, conf...)
	}
}
