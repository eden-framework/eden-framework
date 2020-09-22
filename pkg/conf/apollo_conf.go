package conf

import (
	"github.com/profzone/eden-framework/pkg/conf/apollo"
	"os"
)

func FromApollo(apolloConfig *apollo.ApolloBaseConfig, conf []interface{}) {
	if apolloConfig == nil || os.Getenv("GOENV") == "LOCAL" {
		return
	}

	apollo.AssignConfWithDefault(*apolloConfig, conf...)
}
