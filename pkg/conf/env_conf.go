package conf

import (
	"github.com/profzone/envconfig"
	"github.com/sirupsen/logrus"
)

func FromEnv(confPrefix string, conf []interface{}) []envconfig.EnvVar {
	var envVars = make([]envconfig.EnvVar, 0)
	for _, c := range conf {
		err := envconfig.Process(confPrefix, c)
		if err != nil {
			logrus.Panic(err)
		}
		envconfig.Usage(confPrefix, c)

		envs, err := envconfig.GatherInfo(confPrefix, c)
		if err != nil {
			logrus.Panic(err)
		}

		envVars = append(envVars, envs...)
	}

	return envVars
}
