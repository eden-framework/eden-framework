package files

import (
	"github.com/profzone/envconfig"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ConfigDefaultFile struct {
	EnvVars []envconfig.EnvVar
}

func NewConfigDefaultFile(envVars []envconfig.EnvVar) *ConfigDefaultFile {
	return &ConfigDefaultFile{EnvVars: envVars}
}

func (f *ConfigDefaultFile) String() string {
	e := make(map[string]string)

	e["GOENV"] = "DEV"

	for _, envVar := range f.EnvVars {
		e[envVar.Key] = envVar.GetValue(false)
	}

	bytes, err := yaml.Marshal(e)
	if err != nil {
		logrus.Panic(err)
	}
	return string(bytes)
}
