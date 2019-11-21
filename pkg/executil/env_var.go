package executil

import (
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

var reEnvVar = regexp.MustCompile("(\\$?\\$)\\{?([A-Z0-9_]+)\\}?")

type EnvVars map[string]string

func (envVars EnvVars) AddEnvVar(key, value string) {
	envVars[key] = value
}

func (envVars EnvVars) LoadFromEnviron() {
	for _, keyPair := range os.Environ() {
		keyValue := strings.Split(keyPair, "=")
		envVars.AddEnvVar(keyValue[0], keyValue[1])
	}
}

func (envVars EnvVars) Parse(s string) string {
	result := reEnvVar.ReplaceAllStringFunc(s, func(str string) string {
		matched := reEnvVar.FindAllStringSubmatch(str, -1)[0]

		// skip $${ }
		if matched[1] == "$$" {
			return "${" + matched[2] + "}"
		}

		if value, ok := envVars[matched[2]]; ok {
			return value
		}

		logrus.Errorf("Missing environment variable ${%s}", matched[2])
		return "${" + matched[2] + "}"
	})

	return result
}
